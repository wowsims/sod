import { IndividualSimUI } from '../../../individual_sim_ui';
import { PseudoStat, Spec, Stat } from '../../../proto/common';
import { UnitStat } from '../../../proto_utils/stats';
import { specNames } from '../../../proto_utils/utils';
import { IndividualExporter } from './individual_exporter';

export class Individual60UEPExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Sixty Upgrades EP Export', allowDownload: true });
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allUnitStats = UnitStat.getAll();

		const namesToWeights: Record<string, number> = {};
		allUnitStats.forEach(stat => {
			const statName = Individual60UEPExporter.getName(stat);
			const weight = epValues.getUnitStat(stat);
			if (weight == 0 || statName == '') {
				return;
			}

			// Need to add together stats with the same name (e.g. hit/crit/haste).
			if (namesToWeights[statName]) {
				namesToWeights[statName] += weight;
			} else {
				namesToWeights[statName] = weight;
			}
		});

		return (
			`https://sixtyupgrades.com/sod/ep/import?name=${encodeURIComponent(`${specNames[player.spec]} WoWSims Weights`)}` +
			Object.keys(namesToWeights)
				.map(statName => `&${statName}=${namesToWeights[statName].toFixed(3)}`)
				.join('')
		);
	}

	static getName(stat: UnitStat): string {
		if (stat.isStat()) {
			return Individual60UEPExporter.statNames[stat.getStat()];
		} else {
			return Individual60UEPExporter.pseudoStatNames[stat.getPseudoStat()] || '';
		}
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'strength',
		[Stat.StatAgility]: 'agility',
		[Stat.StatStamina]: 'stamina',
		[Stat.StatIntellect]: 'intellect',
		[Stat.StatSpirit]: 'spirit',
		[Stat.StatSpellPower]: 'spellPower',
		[Stat.StatSpellDamage]: 'spellDamage',
		[Stat.StatArcanePower]: 'arcaneDamage',
		[Stat.StatHolyPower]: 'holyDamage',
		[Stat.StatFirePower]: 'fireDamage',
		[Stat.StatFrostPower]: 'frostDamage',
		[Stat.StatNaturePower]: 'natureDamage',
		[Stat.StatShadowPower]: 'shadowDamage',
		[Stat.StatMP5]: 'mp5',
		[Stat.StatSpellHit]: 'spellHit',
		[Stat.StatSpellCrit]: 'spellCrit',
		[Stat.StatSpellHaste]: 'spellHaste',
		[Stat.StatSpellPenetration]: 'spellPen',
		[Stat.StatAttackPower]: 'attackPower',
		[Stat.StatMeleeHit]: 'hit',
		[Stat.StatMeleeCrit]: 'crit',
		[Stat.StatMeleeHaste]: 'haste',
		[Stat.StatArmorPenetration]: 'armorPen',
		[Stat.StatExpertise]: 'expertise',
		[Stat.StatMana]: 'mana',
		[Stat.StatEnergy]: 'energy',
		[Stat.StatRage]: 'rage',
		[Stat.StatArmor]: 'armor',
		[Stat.StatRangedAttackPower]: 'attackPower',
		[Stat.StatDefense]: 'defense',
		[Stat.StatBlock]: 'block',
		[Stat.StatBlockValue]: 'blockValue',
		[Stat.StatDodge]: 'dodge',
		[Stat.StatParry]: 'parry',
		[Stat.StatResilience]: 'resilience',
		[Stat.StatHealth]: 'health',
		[Stat.StatArcaneResistance]: 'arcaneResistance',
		[Stat.StatFireResistance]: 'fireResistance',
		[Stat.StatFrostResistance]: 'frostResistance',
		[Stat.StatNatureResistance]: 'natureResistance',
		[Stat.StatShadowResistance]: 'shadowResistance',
		[Stat.StatBonusArmor]: 'armorBonus',
		[Stat.StatHealingPower]: 'healing',
		[Stat.StatFeralAttackPower]: 'feralAttackPower',
	};
	static pseudoStatNames: Partial<Record<PseudoStat, string>> = {
		[PseudoStat.PseudoStatMainHandDps]: 'dps',
		[PseudoStat.PseudoStatRangedDps]: 'rangedDps',
		// Weapon Skills
		[PseudoStat.PseudoStatUnarmedSkill]: 'unarmedSkill',
		[PseudoStat.PseudoStatDaggersSkill]: 'daggerSkill',
		[PseudoStat.PseudoStatSwordsSkill]: 'swordSkill',
		[PseudoStat.PseudoStatMacesSkill]: 'maceSkill',
		[PseudoStat.PseudoStatAxesSkill]: 'axeSkill',
		[PseudoStat.PseudoStatTwoHandedSwordsSkill]: 'sword2hSkill',
		[PseudoStat.PseudoStatTwoHandedMacesSkill]: 'mace2hSkill',
		[PseudoStat.PseudoStatTwoHandedAxesSkill]: 'axe2hSkill',
		[PseudoStat.PseudoStatPolearmsSkill]: 'polearmSkill',
		[PseudoStat.PseudoStatStavesSkill]: 'staffSkill',
		[PseudoStat.PseudoStatBowsSkill]: 'bowSkill',
		[PseudoStat.PseudoStatCrossbowsSkill]: 'crossbowSkill',
		[PseudoStat.PseudoStatGunsSkill]: 'gunSkill',
		[PseudoStat.PseudoStatThrownSkill]: 'thrownSkill',
	};
}
