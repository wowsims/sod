import { IndividualSimUI } from '../../../individual_sim_ui';
import { PseudoStat, Spec, Stat } from '../../../proto/common';
import { classNames } from '../../../proto_utils/names';
import { UnitStat } from '../../../proto_utils/stats';
import { specNames } from '../../../proto_utils/utils';
import { IndividualExporter } from './individual_exporter';

export class IndividualPawnEPExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Pawn EP Export', allowDownload: true });
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allUnitStats = UnitStat.getAll();

		const namesToWeights: Record<string, number> = {};
		allUnitStats.forEach(stat => {
			const statName = IndividualPawnEPExporter.getName(stat);
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
			`( Pawn: v1: "${specNames[player.spec]} WoWSims Weights": Class=${classNames.get(player.getClass())},` +
			Object.keys(namesToWeights)
				.map(statName => `${statName}=${namesToWeights[statName].toFixed(3)}`)
				.join(',') +
			' )'
		);
	}

	static getName(stat: UnitStat): string {
		if (stat.isStat()) {
			return IndividualPawnEPExporter.statNames[stat.getStat()];
		} else {
			return IndividualPawnEPExporter.pseudoStatNames[stat.getPseudoStat()] || '';
		}
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'Strength',
		[Stat.StatAgility]: 'Agility',
		[Stat.StatStamina]: 'Stamina',
		[Stat.StatIntellect]: 'Intellect',
		[Stat.StatSpirit]: 'Spirit',
		[Stat.StatSpellPower]: 'SpellDamage',
		[Stat.StatSpellDamage]: 'SpellDamage',
		[Stat.StatArcanePower]: 'ArcaneSpellDamage',
		[Stat.StatFirePower]: 'FireSpellDamage',
		[Stat.StatFrostPower]: 'FrostSpellDamage',
		[Stat.StatHolyPower]: 'HolySpellDamage',
		[Stat.StatNaturePower]: 'NatureSpellDamage',
		[Stat.StatShadowPower]: 'ShadowSpellDamage',
		[Stat.StatMP5]: 'Mp5',
		[Stat.StatSpellHit]: 'SpellHitRating',
		[Stat.StatSpellCrit]: 'SpellCritRating',
		[Stat.StatSpellHaste]: 'SpellHasteRating',
		[Stat.StatSpellPenetration]: 'SpellPen',
		[Stat.StatAttackPower]: 'Ap',
		[Stat.StatMeleeHit]: 'HitRating',
		[Stat.StatMeleeCrit]: 'CritRating',
		[Stat.StatMeleeHaste]: 'HasteRating',
		[Stat.StatArmorPenetration]: 'ArmorPenetration',
		[Stat.StatExpertise]: 'ExpertiseRating',
		[Stat.StatMana]: 'Mana',
		[Stat.StatEnergy]: 'Energy',
		[Stat.StatRage]: 'Rage',
		[Stat.StatArmor]: 'Armor',
		[Stat.StatRangedAttackPower]: 'Ap',
		[Stat.StatDefense]: 'DefenseRating',
		[Stat.StatBlock]: 'BlockRating',
		[Stat.StatBlockValue]: 'BlockValue',
		[Stat.StatDodge]: 'DodgeRating',
		[Stat.StatParry]: 'ParryRating',
		[Stat.StatResilience]: 'ResilienceRating',
		[Stat.StatHealth]: 'Health',
		[Stat.StatArcaneResistance]: 'ArcaneResistance',
		[Stat.StatFireResistance]: 'FireResistance',
		[Stat.StatFrostResistance]: 'FrostResistance',
		[Stat.StatNatureResistance]: 'NatureResistance',
		[Stat.StatShadowResistance]: 'ShadowResistance',
		[Stat.StatBonusArmor]: 'Armor2',
		[Stat.StatHealingPower]: 'Healing',
		[Stat.StatFeralAttackPower]: 'FeralAttackPower',
	};
	static pseudoStatNames: Partial<Record<PseudoStat, string>> = {
		[PseudoStat.PseudoStatMainHandDps]: 'MeleeDps',
		[PseudoStat.PseudoStatRangedDps]: 'RangedDps',
	};
}
