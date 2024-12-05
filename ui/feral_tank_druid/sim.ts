import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { APLRotation_Type as APLRotationType } from '../core/proto/apl.js';
import { Player } from '../core/player.js';
import { Class, Faction, PartyBuffs, PseudoStat, Race, Spec, Stat, WeaponImbue } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFeralTankDruid, {
	cssClass: 'feral-tank-druid-sim-ui',
	cssScheme: 'druid',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatDodge,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		Stat.StatDefense,
		Stat.StatDodge,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatShadowResistance,
	],
	displayPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	
	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatArmor]: 3.5665,
				[Stat.StatBonusArmor]: 0.5187,
				[Stat.StatStamina]: 7.3021,
				[Stat.StatStrength]: 2.3786,
				[Stat.StatAgility]: 4.4974,
				[Stat.StatAttackPower]: 1,
				[Stat.StatMeleeHit]: 2.9282,
				[Stat.StatMeleeCrit]: 1.5143,
				[Stat.StatDefense]: 1.8171,
				[Stat.StatDodge]: 2.0196,
				[Stat.StatHealth]: 0.4465,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 0.0,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default rotation settings.
		rotationType: APLRotationType.TypeSimple,
		simpleRotation: Presets.DefaultRotation,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: Presets.DefaultIndividualBuffs,
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.FeralTankDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellCritBuff, BuffDebuffInputs.SpellISBDebuff],
	excludeBuffDebuffInputs: [WeaponImbue.ElementalSharpeningStone, WeaponImbue.DenseSharpeningStone, WeaponImbue.WildStrikes],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.InspirationUptime,
			OtherInputs.HpPercentForDefensives,
			DruidInputs.StartingRage,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [...Presets.TalentPresets[Phase.Phase2], ...Presets.TalentPresets[Phase.Phase1]],
		// Preset rotations that the user can quickly select.
		rotations: [...Presets.APLPresets[Phase.Phase2], ...Presets.APLPresets[Phase.Phase1]],
		// Preset gear configurations that the user can quickly select.
		gear: [...Presets.GearPresets[Phase.Phase2], ...Presets.GearPresets[Phase.Phase1]],
	},

	autoRotation: player => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFeralTankDruid,
			tooltip: specNames[Spec.SpecFeralTankDruid],
			defaultName: 'Bear',
			iconUrl: getSpecIcon(Class.ClassDruid, 1),

			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceTauren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.DefaultGear.gear,
				},
				[Faction.Horde]: {
					1: Presets.DefaultGear.gear,
				},
			},
		},
	],
});

export class FeralTankDruidSimUI extends IndividualSimUI<Spec.SpecFeralTankDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralTankDruid>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
