import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { APLRotation } from '../core/proto/apl.js';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
// import * as ShamanInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWardenShaman, {
	cssClass: 'warden-shaman-sim-ui',
	cssScheme: 'shaman',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		// Primary
		Stat.StatHealth,
		// Attributes
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatAgility,
		Stat.StatStrength,
		// Physical
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Spell
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatFirePower,
		Stat.StatNaturePower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatMP5,
		// Defense
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatBlock,
		Stat.StatBlockValue,
		Stat.StatParry,
		Stat.StatDodge,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
		PseudoStat.PseudoStatMeleeSpeedMultiplier,
		PseudoStat.PseudoStatCastSpeedMultiplier,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		// Primary
		Stat.StatMana,
		// Attributes
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		// Physical
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Spell
		Stat.StatSpellDamage,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMP5,
		// Defense
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatBlock,
		Stat.StatBlockValue,
		Stat.StatParry,
		Stat.StatDodge,
		// Resistances
		Stat.StatShadowResistance,
	],
	displayPseudoStats: [
		PseudoStat.PseudoStatThornsDamage,
		PseudoStat.PseudoStatMeleeSpeedMultiplier,
		PseudoStat.PseudoStatCastSpeedMultiplier,
	],
	
	defaults: {
		race: Race.RaceTroll,
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatIntellect]: 0.02,
				[Stat.StatAgility]: 1.12,
				[Stat.StatStrength]: 2.29,
				[Stat.StatSpellPower]: 1.15,
				[Stat.StatSpellDamage]: 1.15,
				[Stat.StatFirePower]: 0.63,
				[Stat.StatNaturePower]: 0.48,
				[Stat.StatSpellHit]: 0.03, //default EP assumes cap
				[Stat.StatSpellCrit]: 1.94,
				[Stat.StatMP5]: 0.01,
				[Stat.StatAttackPower]: 1.0,
				[Stat.StatMeleeHit]: 9.62,
				[Stat.StatMeleeCrit]: 14.8,

				// Tank Stats
				[Stat.StatHealth]: 0.15,
				[Stat.StatStamina]: 0.35,
				[Stat.StatArmor]: 1.0,
				[Stat.StatBonusArmor]: 1.0,
				[Stat.StatDefense]: 3.16,
				[Stat.StatBlock]: 26.7,
				[Stat.StatBlockValue]: 3.06,
				[Stat.StatParry]: 26.7,
				[Stat.StatDodge]: 26.7,
				[Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 8.15,
				[PseudoStat.PseudoStatOffHandDps]: 5.81,
				[PseudoStat.PseudoStatMeleeSpeedMultiplier]: 5.81,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
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
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.BlessingOfWisdom, BuffDebuffInputs.ManaSpringTotem, BuffDebuffInputs.SpiritBuff],
	excludeBuffDebuffInputs: [BuffDebuffInputs.BleedDebuff],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InFrontOfTarget,
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InspirationUptime,
		],
	},
	itemSwapConfig: {
		itemSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
	},
	customSections: [],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		talents: [
			...Presets.TalentPresets[Phase.Phase5],
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		rotations: [
			...Presets.APLPresets[Phase.Phase5],
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		gear: [
			...Presets.GearPresets[Phase.Phase5],
			...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
	},

	autoRotation: (player): APLRotation => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecBalanceDruid,
			tooltip: specNames[Spec.SpecBalanceDruid],
			defaultName: 'Balance',
			iconUrl: getSpecIcon(Class.ClassDruid, 0),

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceTauren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {},
				[Faction.Horde]: {},
			},
		},
	],
});

export class WardenShamanSimUI extends IndividualSimUI<Spec.SpecWardenShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWardenShaman>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
