import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { ItemSlot, PartyBuffs, PseudoStat, Spec, Stat } from '../core/proto/common.js';
import { HunterRune } from '../core/proto/hunter.js';
import { Stats } from '../core/proto_utils/stats.js';
import * as HunterInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecHunter, {
	cssClass: 'hunter-sim-ui',
	cssScheme: 'hunter',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatRangedAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMP5,
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatNaturePower,
		Stat.StatArcanePower,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps, PseudoStat.PseudoStatRangedDps],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatRangedAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatAttackPower,
		Stat.StatRangedAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMP5,
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatNaturePower,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 0.3,
				[Stat.StatAgility]: 0.64,
				[Stat.StatStamina]: 0.0,
				[Stat.StatIntellect]: 0.02,
				[Stat.StatAttackPower]: 1,
				[Stat.StatRangedAttackPower]: 1.0,
				[Stat.StatMeleeHit]: 3.29,
				[Stat.StatMeleeCrit]: 4.45,
				[Stat.StatMeleeHaste]: 1.08,
				[Stat.StatSpellPower]: 0.03,
				[Stat.StatNaturePower]: 0.01,
				[Stat.StatArcanePower]: 0.01,
				[Stat.StatMP5]: 0.05,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 2.11,
				[PseudoStat.PseudoStatOffHandDps]: 1.39,
				[PseudoStat.PseudoStatRangedDps]: 6.32,
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
	playerIconInputs: [HunterInputs.PetTypeInput, HunterInputs.WeaponAmmo, HunterInputs.QuiverInput],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: HunterInputs.HunterRotationConfig,
	petConsumeInputs: [ConsumablesInputs.PetScrollOfAgility, ConsumablesInputs.PetScrollOfStrength],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellScorchDebuff, BuffDebuffInputs.StaminaBuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			//HunterInputs.NewRaptorStrike,
			HunterInputs.PetAttackSpeedInput,
			HunterInputs.PetUptime,
			HunterInputs.SniperTrainingUptime,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [...Presets.TalentPresets[Phase.Phase3], ...Presets.TalentPresets[Phase.Phase2], ...Presets.TalentPresets[Phase.Phase1]],
		// Preset rotations that the user can quickly select.
		rotations: [...Presets.APLPresets[Phase.Phase3], ...Presets.APLPresets[Phase.Phase2], ...Presets.APLPresets[Phase.Phase1]],
		// Preset gear configurations that the user can quickly select.
		gear: [...Presets.GearPresets[Phase.Phase3], ...Presets.GearPresets[Phase.Phase2], ...Presets.GearPresets[Phase.Phase1]],
		builds: [Presets.PresetBuildMeleeBM, Presets.PresetBuildRangedMM],
	},

	autoRotation: player => {
		const isMelee =
			player.getEquippedItem(ItemSlot.ItemSlotWaist)?.rune?.id == HunterRune.RuneBeltMeleeSpecialist &&
			player.getEquippedItem(ItemSlot.ItemSlotFeet)?.rune?.id == HunterRune.RuneBootsDualWieldSpecialization;

		if (isMelee) {
			return player.getLevel() == 50
				? Presets.APLMeleeBmPhase3.rotation.rotation!
				: player.getLevel() == 40
				? Presets.APLMeleePhase2.rotation.rotation!
				: Presets.APLMeleeWeavePhase1.rotation.rotation!;
		} else {
			if (player.getLevel() == 50) {
				return Presets.APLRangedMmPhase3.rotation.rotation!;
			} else if (player.getLevel() == 40) {
				if (player.getTalentTree() == 1) {
					return Presets.APLRangedMmPhase2.rotation.rotation!;
				} else {
					return Presets.APLRangedBmPhase2.rotation.rotation!;
				}
			} else {
				return Presets.APLMeleeWeavePhase1.rotation.rotation!;
			}
		}
	},

	raidSimPresets: [
		// Raid sim presets dont work very well with SoD specs between phases
		// and we dont support raid sim atm so just comment this out
		// {
		// 	spec: Spec.SpecHunter,
		// 	tooltip: 'Beast Mastery Hunter',
		// 	defaultName: 'Beast Mastery',
		// 	iconUrl: getSpecIcon(Class.ClassHunter, 0),
		// 	talents: Presets.DefaultTalentsBeastMastery.data,
		// 	specOptions: Presets.BMDefaultOptions,
		// 	consumes: Presets.DefaultConsumes,
		// 	defaultFactionRaces: {
		// 		[Faction.Unknown]: Race.RaceUnknown,
		// 		[Faction.Alliance]: Race.RaceNightElf,
		// 		[Faction.Horde]: Race.RaceOrc,
		// 	},
		// 	defaultGear: {
		// 		[Faction.Unknown]: {},
		// 		[Faction.Alliance]: {
		// 			1: Presets.GearPresets[Phase.Phase1][0].gear,
		// 		},
		// 		[Faction.Horde]: {
		// 			1: Presets.GearPresets[Phase.Phase1][0].gear,
		// 		},
		// 	},
		// },
		// {
		// 	spec: Spec.SpecHunter,
		// 	tooltip: 'Marksmanship Hunter',
		// 	defaultName: 'Marksmanship',
		// 	iconUrl: getSpecIcon(Class.ClassHunter, 1),
		// 	talents: Presets.DefaultTalentsMarksman.data,
		// 	specOptions: Presets.DefaultOptions,
		// 	consumes: Presets.DefaultConsumes,
		// 	defaultFactionRaces: {
		// 		[Faction.Unknown]: Race.RaceUnknown,
		// 		[Faction.Alliance]: Race.RaceNightElf,
		// 		[Faction.Horde]: Race.RaceOrc,
		// 	},
		// 	defaultGear: {
		// 		[Faction.Unknown]: {},
		// 		[Faction.Alliance]: {
		// 			1: Presets.GearPresets[Phase.Phase1][1].gear,
		// 		},
		// 		[Faction.Horde]: {
		// 			1: Presets.GearPresets[Phase.Phase1][1].gear,
		// 		},
		// 	},
		// },
		// {
		// 	spec: Spec.SpecHunter,
		// 	tooltip: 'Survival Hunter',
		// 	defaultName: 'Survival',
		// 	iconUrl: getSpecIcon(Class.ClassHunter, 2),
		// 	talents: Presets.DefaultTalentsSurvival.data,
		// 	specOptions: Presets.DefaultOptions,
		// 	consumes: Presets.DefaultConsumes,
		// 	defaultFactionRaces: {
		// 		[Faction.Unknown]: Race.RaceUnknown,
		// 		[Faction.Alliance]: Race.RaceNightElf,
		// 		[Faction.Horde]: Race.RaceOrc,
		// 	},
		// 	defaultGear: {
		// 		[Faction.Unknown]: {},
		// 		[Faction.Alliance]: {
		// 			1: Presets.GearPresets[Phase.Phase1][2].gear,
		// 		},
		// 		[Faction.Horde]: {
		// 			1: Presets.GearPresets[Phase.Phase1][2].gear,
		// 		},
		// 	},
		// },
	],
});

export class HunterSimUI extends IndividualSimUI<Spec.SpecHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHunter>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
