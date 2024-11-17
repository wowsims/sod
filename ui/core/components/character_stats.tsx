import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import * as Mechanics from '../constants/mechanics.js';
import { Player } from '../player.js';
import { PseudoStat, Spec, Stat } from '../proto/common.js';
import { getClassStatName, statOrder } from '../proto_utils/names.js';
import { Stats } from '../proto_utils/stats.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { Component } from './component.js';
import { NumberPicker } from './number_picker';

export type StatMods = { talents?: Stats; buffs?: Stats };

const statGroups = new Map<string, Array<Stat>>([
	[
		'Primary',
		[
			Stat.StatHealth,
			Stat.StatMana,
		],
	],
	[
		'Attributes',
		[
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatStamina,
			Stat.StatIntellect,
			Stat.StatSpirit,
		]
	],
	[
		'Physical',
		[
			Stat.StatAttackPower,
			Stat.StatFeralAttackPower,
			Stat.StatRangedAttackPower,
			Stat.StatMeleeHit,
			Stat.StatExpertise,
			Stat.StatMeleeCrit,
			Stat.StatMeleeHaste,
		]
	],
	[
		'Spell',
		[
			Stat.StatSpellPower,
			Stat.StatSpellDamage,
			Stat.StatArcanePower,
			Stat.StatFirePower,
			Stat.StatFrostPower,
			Stat.StatHolyPower,
			Stat.StatNaturePower,
			Stat.StatShadowPower,
			Stat.StatSpellHit,
			Stat.StatSpellCrit,
			Stat.StatSpellHaste,
			Stat.StatSpellPenetration,
			Stat.StatMP5,
		]
	],
	[
		'Defense',
		[
			Stat.StatArmor,
			Stat.StatBonusArmor,
			Stat.StatDefense,
			Stat.StatDodge,
			Stat.StatParry,
			Stat.StatBlock,
			Stat.StatBlockValue,
		]
	],
	[
		'Resistance',
		[
			Stat.StatArcaneResistance,
			Stat.StatFireResistance,
			Stat.StatFrostResistance,
			Stat.StatNatureResistance,
			Stat.StatShadowResistance,
		]
	],
])

export class CharacterStats extends Component {
	readonly stats: Array<Stat>;
	readonly valueElems: Array<HTMLTableCellElement>;
	readonly meleeCritCapValueElem: HTMLTableCellElement | undefined;

	private readonly player: Player<any>;
	private readonly modifyDisplayStats?: (player: Player<any>) => StatMods;

	constructor(parent: HTMLElement, player: Player<any>, stats: Array<Stat>, modifyDisplayStats?: (player: Player<any>) => StatMods) {
		super(parent, 'character-stats-root');
		this.stats = []
		this.player = player;
		this.modifyDisplayStats = modifyDisplayStats;

		const statsSet = new Set(stats);

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
			const filteredStats = groupedStats.filter(stat => statsSet.has(stat))
			if (!filteredStats.length) return

			const body = <tbody></tbody>
			filteredStats.forEach(stat => {
				this.stats.push(stat)

				const statName = getClassStatName(stat, player.getClass());

				const row = (
					<tr className="character-stats-table-row">
						<td className="character-stats-table-label">{statName}</td>
						<td className="character-stats-table-value">{this.bonusStatsLink(stat)}</td>
					</tr>
				);
				body.appendChild(row);
	
				const valueElem = row.getElementsByClassName('character-stats-table-value')[0] as HTMLTableCellElement;
				this.valueElems.push(valueElem);
			})

			table.appendChild(body)
		})

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

		const statMods = this.modifyDisplayStats ? this.modifyDisplayStats(this.player) : {};
		if (!statMods.talents) statMods.talents = new Stats();
		if (!statMods.buffs) statMods.buffs = new Stats();

		const baseStats = Stats.fromProto(playerStats.baseStats);
		const gearStats = Stats.fromProto(playerStats.gearStats);
		const talentsStats = Stats.fromProto(playerStats.talentsStats);
		const buffsStats = Stats.fromProto(playerStats.buffsStats);
		const consumesStats = Stats.fromProto(playerStats.consumesStats);
		const debuffStats = this.getDebuffStats();
		const bonusStats = player.getBonusStats();

		const baseDelta = baseStats;
		const gearDelta = gearStats.subtract(baseStats).subtract(bonusStats);
		const talentsDelta = talentsStats.subtract(gearStats).add(statMods.talents);
		const buffsDelta = buffsStats.subtract(talentsStats).add(statMods.buffs);
		const consumesDelta = consumesStats.subtract(buffsStats);

		const finalStats = Stats.fromProto(playerStats.finalStats).add(statMods.talents).add(statMods.buffs).add(debuffStats);

		this.stats.forEach((stat, idx) => {
			const bonusStatValue = bonusStats.getStat(stat);
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
						{debuffStats.getStat(stat) != 0 && (
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

			if (stat === Stat.StatMeleeHit) {
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
			} else if (stat === Stat.StatSpellHit) {
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

	private statDisplayString(player: Player<any>, stats: Stats, deltaStats: Stats, stat: Stat): string {
		let rawValue = deltaStats.getStat(stat);

		if (stat === Stat.StatBlockValue) {
			rawValue *= stats.getPseudoStat(PseudoStat.PseudoStatBlockValueMultiplier) || 1;
			rawValue += Math.max(0, stats.getPseudoStat(PseudoStat.PseudoStatBlockValuePerStrength) * deltaStats.getStat(Stat.StatStrength) - 1);
		}

		let displayStr = String(Math.round(rawValue));

		if (stat === Stat.StatMeleeHit) {
			displayStr = `${(rawValue / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%`;
		} else if (stat === Stat.StatSpellHit) {
			displayStr = `${(rawValue / Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE).toFixed(2)}%`;
		} else if (stat === Stat.StatSpellDamage) {
			const spDmg = Math.round(rawValue);
			const baseSp = Math.round(deltaStats.getStat(Stat.StatSpellPower));
			displayStr = baseSp + spDmg + ` (+${spDmg})`;
		} else if (
			stat === Stat.StatArcanePower ||
			stat === Stat.StatFirePower ||
			stat === Stat.StatFrostPower ||
			stat === Stat.StatHolyPower ||
			stat === Stat.StatNaturePower ||
			stat === Stat.StatShadowPower
		) {
			const spDmg = Math.round(rawValue);
			const baseSp = Math.round(deltaStats.getStat(Stat.StatSpellPower) + deltaStats.getStat(Stat.StatSpellDamage));
			displayStr = baseSp + spDmg + ` (+${spDmg})`;
		} else if (stat === Stat.StatMeleeCrit || stat === Stat.StatSpellCrit) {
			displayStr = `${(rawValue / Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE).toFixed(2)}%`;
		} else if (stat === Stat.StatMeleeHaste) {
			// Melee Haste doesn't actually exist in vanilla so use the melee speed pseudostat
			displayStr = `${(deltaStats.getPseudoStat(PseudoStat.PseudoStatMeleeSpeedMultiplier) * 100).toFixed(2)}%`;
		} else if (stat === Stat.StatSpellHaste) {
			displayStr = `${(rawValue / Mechanics.HASTE_RATING_PER_HASTE_PERCENT).toFixed(2)}%`;
		} else if (stat === Stat.StatArmorPenetration) {
			displayStr += ` (${(rawValue / Mechanics.ARMOR_PEN_PER_PERCENT_ARMOR).toFixed(2)}%)`;
		} else if (stat === Stat.StatExpertise) {
			// It's just like crit and hit in SoD.
			displayStr += `%`;
		} else if (stat === Stat.StatDefense) {
			displayStr = `${(player.getLevel() * 5 + Math.floor(rawValue / Mechanics.DEFENSE_RATING_PER_DEFENSE)).toFixed(0)}`;
		} else if (stat === Stat.StatBlock) {
			displayStr = `${(rawValue / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE).toFixed(2)}%`;
		} else if (stat === Stat.StatDodge) {
			displayStr = `${(rawValue / Mechanics.DODGE_RATING_PER_DODGE_CHANCE).toFixed(2)}%`;
		} else if (stat === Stat.StatParry) {
			displayStr = `${(rawValue / Mechanics.PARRY_RATING_PER_PARRY_CHANCE).toFixed(2)}%`;
		} else if (stat === Stat.StatResilience) {
			displayStr += ` (${(rawValue / Mechanics.RESILIENCE_RATING_PER_CRIT_REDUCTION_CHANCE).toFixed(2)}%)`;
		}

		return displayStr;
	}

	private weaponSkillDisplayString(stats: Stats, pseudoStat: PseudoStat): string {
		return `${300 + stats.getPseudoStat(pseudoStat)}`;
	}

	private spellSchoolHitDisplayString(stats: Stats, pseudoStat: PseudoStat): string {
		return `${(stats.getPseudoStat(pseudoStat) + stats.getStat(Stat.StatSpellHit)).toFixed(2)}%`;
	}

	private getDebuffStats(): Stats {
		const debuffStats = new Stats();

		// TODO: Classic ui debuffs
		// const debuffs = this.player.sim.raid.getDebuffs();
		// if (debuffs.improvedScorch || debuffs.wintersChill || debuffs.shadowMastery) {
		// 	debuffStats = debuffStats.addStat(Stat.StatSpellCrit, 5 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE);
		// }

		return debuffStats;
	}

	private bonusStatsLink(stat: Stat): HTMLElement {
		const statName = getClassStatName(stat, this.player.getClass());
		const linkRef = ref<HTMLAnchorElement>();
		const iconRef = ref<HTMLDivElement>();

		const link = (
			<a
				ref={linkRef}
				href="javascript:void(0)"
				className="add-bonus-stats text-white ms-2"
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
					id: `character-bonus-stat-${stat}`,
					label: `Bonus ${statName}`,
					extraCssClasses: ['mb-0'],
					changedEvent: (player: Player<any>) => player.bonusStatsChangeEmitter,
					getValue: (player: Player<any>) => player.getBonusStats().getStat(stat),
					setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
						const bonusStats = player.getBonusStats().withStat(stat, newValue);
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
