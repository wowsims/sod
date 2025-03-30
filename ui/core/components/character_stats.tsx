import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import * as Mechanics from '../constants/mechanics.js';
import { Player } from '../player.js';
import { InputType, ItemSlot, PseudoStat, Spec, Stat, WeaponType } from '../proto/common.js';
import { Stats, UnitStat } from '../proto_utils/stats.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { Component } from './component.js';
import { NumberPicker } from './number_picker';

export type StatMods = { talents?: Stats; buffs?: Stats };
export type DisplayStat = {
	stat: UnitStat;
	notEditable?: boolean;
};

const statGroups = new Map<string, Array<DisplayStat>>([
	['Primary', [{ stat: UnitStat.fromStat(Stat.StatHealth) }, { stat: UnitStat.fromStat(Stat.StatMana) }]],
	[
		'Attributes',
		[
			{ stat: UnitStat.fromStat(Stat.StatStrength) },
			{ stat: UnitStat.fromStat(Stat.StatAgility) },
			{ stat: UnitStat.fromStat(Stat.StatStamina) },
			{ stat: UnitStat.fromStat(Stat.StatIntellect) },
			{ stat: UnitStat.fromStat(Stat.StatSpirit) },
		],
	],
	[
		'Physical',
		[
			{ stat: UnitStat.fromStat(Stat.StatAttackPower) },
			{ stat: UnitStat.fromStat(Stat.StatFeralAttackPower) },
			{ stat: UnitStat.fromStat(Stat.StatRangedAttackPower) },
			{ stat: UnitStat.fromStat(Stat.StatMeleeHit) },
			{ stat: UnitStat.fromStat(Stat.StatExpertise) },
			{ stat: UnitStat.fromStat(Stat.StatMeleeCrit) },
			{ stat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatMeleeSpeedMultiplier) },
			{ stat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatRangedSpeedMultiplier) },
			{ stat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatBonusPhysicalDamage) },
		],
	],
	[
		'Spell',
		[
			{ stat: UnitStat.fromStat(Stat.StatSpellPower) },
			{ stat: UnitStat.fromStat(Stat.StatSpellDamage) },
			{ stat: UnitStat.fromStat(Stat.StatArcanePower) },
			{ stat: UnitStat.fromStat(Stat.StatFirePower) },
			{ stat: UnitStat.fromStat(Stat.StatFrostPower) },
			{ stat: UnitStat.fromStat(Stat.StatHolyPower) },
			{ stat: UnitStat.fromStat(Stat.StatNaturePower) },
			{ stat: UnitStat.fromStat(Stat.StatShadowPower) },
			{ stat: UnitStat.fromStat(Stat.StatHealingPower) },
			{ stat: UnitStat.fromStat(Stat.StatSpellHit) },
			{ stat: UnitStat.fromStat(Stat.StatSpellCrit) },
			{ stat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatCastSpeedMultiplier) },
			{ stat: UnitStat.fromStat(Stat.StatSpellPenetration) },
			{ stat: UnitStat.fromStat(Stat.StatMP5) },
		],
	],
	[
		'Defense',
		[
			{ stat: UnitStat.fromStat(Stat.StatArmor) },
			{ stat: UnitStat.fromStat(Stat.StatBonusArmor) },
			{ stat: UnitStat.fromStat(Stat.StatDefense) },
			{ stat: UnitStat.fromStat(Stat.StatDodge) },
			{ stat: UnitStat.fromStat(Stat.StatParry) },
			{ stat: UnitStat.fromStat(Stat.StatBlock) },
			{ stat: UnitStat.fromStat(Stat.StatBlockValue) },
		],
	],
	[
		'Resistance',
		[
			{ stat: UnitStat.fromStat(Stat.StatArcaneResistance) },
			{ stat: UnitStat.fromStat(Stat.StatFireResistance) },
			{ stat: UnitStat.fromStat(Stat.StatFrostResistance) },
			{ stat: UnitStat.fromStat(Stat.StatNatureResistance) },
			{ stat: UnitStat.fromStat(Stat.StatShadowResistance) },
		],
	],
	[
		'Misc',
		[
			{ stat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatThornsDamage), notEditable: true },
			{ stat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatTimewornBonus) },
			{ stat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSanctifiedBonus) },
		],
	],
]);

export class CharacterStats extends Component {
	readonly stats: Array<UnitStat>;
	readonly valueElems: Array<HTMLTableCellElement>;
	readonly meleeCritCapValueElem: HTMLTableCellElement | undefined;

	private readonly player: Player<any>;
	private readonly modifyDisplayStats?: (player: Player<any>) => StatMods;

	constructor(parent: HTMLElement, player: Player<any>, displayStats: Array<UnitStat>, modifyDisplayStats?: (player: Player<any>) => StatMods) {
		super(parent, 'character-stats-root');
		this.stats = [];
		this.player = player;
		this.modifyDisplayStats = modifyDisplayStats;

		const playerLevelRef = ref<HTMLSpanElement>();
		this.player.levelChangeEmitter.on(() => (playerLevelRef.value!.textContent = `Level ${player.getLevel()}`));

		this.rootElem.appendChild(
			<label className="character-stats-label">
				<span>Stats</span>
				<span ref={playerLevelRef} className="ms-auto">
					Level {player.getLevel()}
				</span>
			</label>,
		);

		const table = <table className="character-stats-table"></table>;
		this.rootElem.appendChild(table);

		this.valueElems = [];
		statGroups.forEach((groupedStats, _) => {
			const filteredStats = groupedStats.filter(stat => displayStats.find(displayStat => displayStat.equals(stat.stat)));

			if (!filteredStats.length) return;

			const body = <tbody></tbody>;
			filteredStats.forEach(displayStat => {
				const { stat } = displayStat;
				this.stats.push(stat);

				const statName = stat.getName(player.getClass());

				const row = (
					<tr className="character-stats-table-row">
						<td className="character-stats-table-label">{statName}</td>
						<td className="character-stats-table-value">{this.bonusStatsLink(displayStat)}</td>
					</tr>
				);
				body.appendChild(row);

				const valueElem = row.getElementsByClassName('character-stats-table-value')[0] as HTMLTableCellElement;
				this.valueElems.push(valueElem);
			});

			table.appendChild(body);
		});

		if (this.shouldShowMeleeCritCap(player)) {
			const row = (
				<tr className="character-stats-table-row">
					<td className="character-stats-table-label">Melee Crit Cap</td>
					<td className="character-stats-table-value"></td>
				</tr>
			);

			table.appendChild(row);
			this.meleeCritCapValueElem = row.getElementsByClassName('character-stats-table-value')[0] as HTMLTableCellElement;
		}

		this.updateStats(player);
		TypedEvent.onAny([player.currentStatsEmitter, player.sim.changeEmitter, player.talentsChangeEmitter]).on(() => {
			this.updateStats(player);
		});
	}

	private updateStats(player: Player<any>) {
		const playerStats = player.getCurrentStats();
		const gear = player.getGear();
		const mainHandWeapon = gear.getEquippedItem(ItemSlot.ItemSlotMainHand);
		const offHandItem = gear.getEquippedItem(ItemSlot.ItemSlotOffHand);

		const statMods = this.modifyDisplayStats ? this.modifyDisplayStats(this.player) : {};
		if (!statMods.talents) statMods.talents = new Stats();
		if (!statMods.buffs) statMods.buffs = new Stats();

		const baseStats = Stats.fromProto(playerStats.baseStats);
		const gearStats = Stats.fromProto(playerStats.gearStats);
		const talentsStats = Stats.fromProto(playerStats.talentsStats);
		const buffsStats = Stats.fromProto(playerStats.buffsStats);
		const consumesStats = Stats.fromProto(playerStats.consumesStats);
		const debuffStats = this.getDebuffStats(player);
		const bonusStats = player.getBonusStats();

		const baseDelta = baseStats;
		const gearDelta = gearStats.subtract(baseStats).subtract(bonusStats);
		const talentsDelta = talentsStats.subtract(gearStats).add(statMods.talents);
		const buffsDelta = buffsStats.subtract(talentsStats).add(statMods.buffs);
		const consumesDelta = consumesStats.subtract(buffsStats);

		const finalStats = Stats.fromProto(playerStats.finalStats).add(statMods.talents).add(statMods.buffs).add(debuffStats);

		this.stats.forEach((stat, idx) => {
			const bonusStatValue = bonusStats.getUnitStat(stat);
			let contextualClass: string;
			if (bonusStatValue === 0) {
				contextualClass = 'text-white';
			} else if (bonusStatValue > 0) {
				contextualClass = 'text-success';
			} else {
				contextualClass = 'text-danger';
			}

			const statLinkElemRef = ref<HTMLAnchorElement>();

			const valueElem = (
				<div className="stat-value-link-container">
					<a href="javascript:void(0)" className={`stat-value-link ${contextualClass}`} attributes={{ role: 'button' }} ref={statLinkElemRef}>
						{`${this.statDisplayString(player, finalStats, finalStats, stat)} `}
					</a>
				</div>
			);

			const statLinkElem = statLinkElemRef.value!;

			this.valueElems[idx].querySelector('.stat-value-link-container')?.remove();
			this.valueElems[idx].prepend(valueElem);

			const tooltipContent = (
				<div className="d-flex">
					<div>
						<div className="character-stats-tooltip-row">
							<span>Base:</span>
							<span>{this.statDisplayString(player, baseStats, baseDelta, stat)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Gear:</span>
							<span>{this.statDisplayString(player, gearStats, gearDelta, stat)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Talents:</span>
							<span>{this.statDisplayString(player, talentsStats, talentsDelta, stat)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Buffs:</span>
							<span>{this.statDisplayString(player, buffsStats, buffsDelta, stat)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Consumes:</span>
							<span>{this.statDisplayString(player, consumesStats, consumesDelta, stat)}</span>
						</div>
						{stat.isStat() && debuffStats.getStat(stat.getStat()) != 0 && (
							<div className="character-stats-tooltip-row">
								<span>Debuffs:</span>
								<span>{this.statDisplayString(player, debuffStats, debuffStats, stat)}</span>
							</div>
						)}
						{bonusStatValue != 0 && (
							<div className="character-stats-tooltip-row">
								<span>Bonus:</span>
								<span>{this.statDisplayString(player, bonusStats, bonusStats, stat)}</span>
							</div>
						)}
						<div className="character-stats-tooltip-row">
							<span>Total:</span>
							<span>{this.statDisplayString(player, finalStats, finalStats, stat)}</span>
						</div>
					</div>
				</div>
			);

			if (stat.isStat() && [Stat.StatMeleeHit, Stat.StatExpertise].includes(stat.getStat())) {
				tooltipContent.appendChild(
					<div className="ps-2">
						<div className="character-stats-tooltip-row">
							<span>Axes</span>
							<span>
								{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatAxesSkill)} /{' '}
								{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatTwoHandedAxesSkill)}
							</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Daggers</span>
							<span>{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatDaggersSkill)}</span>
						</div>
						{player.spec === Spec.SpecFeralDruid && (
							<div className="character-stats-tooltip-row">
								<span>Feral Combat</span>
								<span>{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatFeralCombatSkill)}</span>
							</div>
						)}
						<div className="character-stats-tooltip-row">
							<span>Maces</span>
							<span>
								{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatMacesSkill)} /{' '}
								{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatTwoHandedMacesSkill)}
							</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Polearms</span>
							<span>{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatPolearmsSkill)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Staves</span>
							<span>{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatStavesSkill)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Swords</span>
							<span>
								{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatSwordsSkill)} /{' '}
								{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatTwoHandedSwordsSkill)}
							</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Unarmed</span>
							<span>{this.weaponSkillDisplayString(gearStats, PseudoStat.PseudoStatUnarmedSkill)}</span>
						</div>
					</div>,
				);
			} else if (stat.isStat() && stat.getStat() === Stat.StatSpellHit) {
				tooltipContent.appendChild(
					<div className="ps-2">
						<div className="character-stats-tooltip-row">
							<span>Arcane</span>
							<span>{this.spellSchoolHitDisplayString(finalStats, PseudoStat.PseudoStatSchoolHitArcane)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Fire</span>
							<span>{this.spellSchoolHitDisplayString(finalStats, PseudoStat.PseudoStatSchoolHitFire)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Frost</span>
							<span>{this.spellSchoolHitDisplayString(finalStats, PseudoStat.PseudoStatSchoolHitFrost)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Holy</span>
							<span>{this.spellSchoolHitDisplayString(finalStats, PseudoStat.PseudoStatSchoolHitHoly)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Nature</span>
							<span>{this.spellSchoolHitDisplayString(finalStats, PseudoStat.PseudoStatSchoolHitNature)}</span>
						</div>
						<div className="character-stats-tooltip-row">
							<span>Shadow</span>
							<span>{this.spellSchoolHitDisplayString(finalStats, PseudoStat.PseudoStatSchoolHitShadow)}</span>
						</div>
					</div>,
				);
			} else if (stat.isPseudoStat() && stat.getPseudoStat() === PseudoStat.PseudoStatMeleeSpeedMultiplier && (mainHandWeapon || offHandItem)) {
				const speedStat = finalStats.getPseudoStat(PseudoStat.PseudoStatMeleeSpeedMultiplier);
				const offHandWeapon =
					offHandItem &&
					offHandItem.item.weaponType !== WeaponType.WeaponTypeShield &&
					offHandItem.item.weaponType !== WeaponType.WeaponTypeOffHand &&
					offHandItem.item.weaponType !== WeaponType.WeaponTypeUnknown;
				const mainHandLabel = offHandWeapon ? 'Main-hand' : 'Weapon';
				tooltipContent.appendChild(
					<div className="ps-2">
						{mainHandWeapon && (
							<div className="character-stats-tooltip-row">
								<span>{mainHandLabel} Speed</span>
								<span>{(mainHandWeapon.item.weaponSpeed / speedStat).toFixed(2)}s</span>
							</div>
						)}
						{offHandWeapon && (
							<div className="character-stats-tooltip-row">
								<span>Off-hand Speed</span>
								<span>{(offHandItem.item.weaponSpeed / speedStat).toFixed(2)}s</span>
							</div>
						)}
					</div>,
				);
			}

			tippy(statLinkElem, {
				content: tooltipContent,
			});
		});

		if (this.meleeCritCapValueElem) {
			const meleeCritCapInfo = player.getMeleeCritCapInfo();

			const valueElem = (
				<a href="javascript:void(0)" className="stat-value-link" attributes={{ role: 'button' }}>
					{`${this.meleeCritCapDisplayString(player, finalStats)} `}
				</a>
			);

			const capDelta = meleeCritCapInfo.playerCritCapDelta;
			if (capDelta === 0) {
				valueElem.classList.add('text-white');
			} else if (capDelta > 0) {
				valueElem.classList.add('text-danger');
			} else if (capDelta < 0) {
				valueElem.classList.add('text-success');
			}

			this.meleeCritCapValueElem.querySelector('.stat-value-link')?.remove();
			this.meleeCritCapValueElem.prepend(valueElem);

			const tooltipContent = (
				<div>
					<div className="character-stats-tooltip-row">
						<span>Glancing:</span>
						<span>{`${meleeCritCapInfo.glancing.toFixed(2)}%`}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Suppression:</span>
						<span>{`${meleeCritCapInfo.suppression.toFixed(2)}%`}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>To Hit Cap:</span>
						<span>{`${meleeCritCapInfo.remainingMeleeHitCap.toFixed(2)}%`}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>To Exp Cap:</span>
						<span>{`${meleeCritCapInfo.remainingExpertiseCap.toFixed(2)}%`}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Debuffs:</span>
						<span>{`${meleeCritCapInfo.debuffCrit.toFixed(2)}%`}</span>
					</div>
					{meleeCritCapInfo.specSpecificOffset != 0 && (
						<div className="character-stats-tooltip-row">
							<span>Spec Offsets:</span>
							<span>{`${meleeCritCapInfo.specSpecificOffset.toFixed(2)}%`}</span>
						</div>
					)}
					<div className="character-stats-tooltip-row">
						<span>Final Crit Cap:</span>
						<span>{`${meleeCritCapInfo.baseCritCap.toFixed(2)}%`}</span>
					</div>
					<hr />
					<div className="character-stats-tooltip-row">
						<span>Can Raise By:</span>
						<span>{`${(meleeCritCapInfo.remainingExpertiseCap + meleeCritCapInfo.remainingMeleeHitCap).toFixed(2)}%`}</span>
					</div>
				</div>
			);

			tippy(valueElem, {
				content: tooltipContent,
			});
		}
	}

	private statDisplayString(player: Player<any>, stats: Stats, deltaStats: Stats, unitStat: UnitStat): string {
		const rawValue = deltaStats.getUnitStat(unitStat);
		let displayStr: string | undefined;

		if (unitStat.isStat()) {
			const stat = unitStat.getStat();

			switch (stat) {
				case Stat.StatBlockValue:
					const mult = stats.getPseudoStat(PseudoStat.PseudoStatBlockValueMultiplier) || 1;
					const perStr = Math.max(0, stats.getPseudoStat(PseudoStat.PseudoStatBlockValuePerStrength) * deltaStats.getStat(Stat.StatStrength) - 1);
					displayStr = String(Math.round(rawValue * mult + perStr));
					break;
				case Stat.StatMeleeHit:
					displayStr = `${(rawValue / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%`;
					break;
				case Stat.StatSpellHit:
					displayStr = `${(rawValue / Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%`;
					break;
				case Stat.StatSpellDamage:
					const statSpellDamage = Math.round(rawValue);
					const statSpellPower = Math.round(deltaStats.getStat(Stat.StatSpellPower));
					displayStr = statSpellPower + statSpellDamage + ` (+${statSpellDamage})`;
					break;
				case Stat.StatArcanePower:
				case Stat.StatFirePower:
				case Stat.StatFrostPower:
				case Stat.StatHolyPower:
				case Stat.StatNaturePower:
				case Stat.StatShadowPower:
				case Stat.StatHealingPower:
					const schoolDamage = Math.round(rawValue);
					const baseSpellPower = Math.round(deltaStats.getStat(Stat.StatSpellPower) + deltaStats.getStat(Stat.StatSpellDamage));
					displayStr = baseSpellPower + schoolDamage + ` (+${schoolDamage})`;
					break;
				case Stat.StatMeleeCrit:
				case Stat.StatSpellCrit:
					displayStr = `${(rawValue / Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE).toFixed(2)}%`;
					break;
				case Stat.StatArmorPenetration:
					displayStr = `${rawValue} (${(rawValue / Mechanics.ARMOR_PEN_PER_PERCENT_ARMOR).toFixed(2)}%)`;
					break;
				case Stat.StatExpertise:
					// It's just like crit and hit in SoD.
					displayStr = `${rawValue}%`;
					break;
				case Stat.StatDefense:
					displayStr = `${(player.getLevel() * 5 + Math.floor(rawValue / Mechanics.DEFENSE_RATING_PER_DEFENSE)).toFixed(0)}`;
					break;
				case Stat.StatBlock:
					displayStr = `${(rawValue / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE).toFixed(2)}%`;
					break;
				case Stat.StatDodge:
					displayStr = `${(rawValue / Mechanics.DODGE_RATING_PER_DODGE_CHANCE).toFixed(2)}%`;
					break;
				case Stat.StatParry:
					displayStr = `${(rawValue / Mechanics.PARRY_RATING_PER_PARRY_CHANCE).toFixed(2)}%`;
					break;
				case Stat.StatResilience:
					displayStr = `${rawValue} (${(rawValue / Mechanics.RESILIENCE_RATING_PER_CRIT_REDUCTION_CHANCE).toFixed(2)}%)`;
					break;
			}
		} else {
			const pseudoStat = unitStat.getPseudoStat();

			switch (pseudoStat) {
				case PseudoStat.PseudoStatMeleeSpeedMultiplier:
					displayStr = `${(100 * deltaStats.getPseudoStat(PseudoStat.PseudoStatMeleeSpeedMultiplier)).toFixed(2)}%`;
					break;
				case PseudoStat.PseudoStatRangedSpeedMultiplier:
					displayStr = `${(100 * deltaStats.getPseudoStat(PseudoStat.PseudoStatRangedSpeedMultiplier)).toFixed(2)}%`;
					break;
				case PseudoStat.PseudoStatCastSpeedMultiplier:
					displayStr = `${(100 * deltaStats.getPseudoStat(PseudoStat.PseudoStatCastSpeedMultiplier)).toFixed(2)}%`;
					break;
			}
		}

		if (!displayStr) displayStr = String(Math.round(rawValue));

		return displayStr;
	}

	private weaponSkillDisplayString(stats: Stats, pseudoStat: PseudoStat): string {
		return `${300 + stats.getPseudoStat(pseudoStat)}`;
	}

	private spellSchoolHitDisplayString(stats: Stats, pseudoStat: PseudoStat): string {
		return `${(stats.getPseudoStat(pseudoStat) + stats.getStat(Stat.StatSpellHit)).toFixed(2)}%`;
	}

	private raidAvoidanceDebuff(player: Player<any>): Stats {
		const targets = player.sim.encounter.targets;
		let stats = new Stats();
		if (!targets || targets.length === 0) {
			return stats;
		}

		if (player.sim.db) {
			const preset = player.sim.db.getAllPresetEncounters().find(pe => player.sim.encounter.matchesPreset(pe));
			if (preset && preset.path.includes('Scarlet Enclave')) {
				stats = stats.addStat(Stat.StatDodge, -20);
				return stats;
			}
		}

		const targetInputs = targets[0].targetInputs;
		if (!targetInputs || targetInputs.length === 0) {
			return stats;
		}

		const authorityInput = targetInputs.find(x => x.inputType === InputType.Enum && x.label === 'Difficulty Level');

		if (authorityInput && authorityInput.enumValue > 0) {
			stats = stats.addStat(Stat.StatDodge, -4 * authorityInput.enumValue);
		}

		return stats;
	}

	private getDebuffStats(player: Player<any>): Stats {
		const debuffStats = new Stats().add(this.raidAvoidanceDebuff(player));

		return debuffStats;
	}

	private bonusStatsLink(displayStat: DisplayStat): HTMLElement {
		const { stat, notEditable } = displayStat;
		const statName = displayStat.stat.getName(this.player.getClass());
		const linkRef = ref<HTMLAnchorElement>();
		const iconRef = ref<HTMLDivElement>();

		const link = (
			<a
				ref={linkRef}
				href="javascript:void(0)"
				className={clsx('add-bonus-stats text-white ms-2', notEditable && 'invisible')}
				dataset={{ bsToggle: 'popover' }}
				attributes={{ role: 'button' }}>
				<i ref={iconRef} className="fas fa-plus-minus"></i>
			</a>
		);

		tippy(iconRef.value!, { content: `Bonus ${statName}` });
		tippy(linkRef.value!, {
			interactive: true,
			trigger: 'click',
			theme: 'bonus-stats-popover',
			placement: 'right',
			onShow: instance => {
				const picker = new NumberPicker(null, this.player, {
					id: `character-bonus-${stat.isStat() ? 'stat-' + stat.getStat() : 'pseudostat-' + stat.getPseudoStat()}`,
					label: `Bonus ${statName}`,
					extraCssClasses: ['mb-0'],
					changedEvent: (player: Player<any>) => player.bonusStatsChangeEmitter,
					getValue: (player: Player<any>) => player.getBonusStats().getUnitStat(stat),
					setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
						const bonusStats = player.getBonusStats().withUnitStat(stat, newValue);
						player.setBonusStats(eventID, bonusStats);
						instance?.hide();
					},
				});
				instance.setContent(picker.rootElem);
			},
		});

		return link as HTMLElement;
	}

	private shouldShowMeleeCritCap(player: Player<any>): boolean {
		// TODO: Disabled for now while we fix displayed crit cap
		return false;
		// return [Spec.SpecEnhancementShaman, Spec.SpecRetributionPaladin, Spec.SpecRogue, Spec.SpecWarrior, Spec.SpecHunter].includes(player.spec);
	}

	private meleeCritCapDisplayString(player: Player<any>, _: Stats): string {
		const playerCritCapDelta = player.getMeleeCritCap();

		if (playerCritCapDelta === 0.0) {
			return 'Exact';
		}

		const prefix = playerCritCapDelta > 0 ? 'Over by ' : 'Under by ';
		return `${prefix} ${Math.abs(playerCritCapDelta).toFixed(2)}%`;
	}
}
