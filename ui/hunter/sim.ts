import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { ItemSlot, PartyBuffs, PseudoStat, Spec, Stat } from '../core/proto/common.js';
import { HunterRune } from '../core/proto/hunter.js';
import { Stats } from '../core/proto_utils/stats.js';
import { SpecOptions } from '../core/proto_utils/utils';
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
		// Attributes
		Stat.StatIntellect,
		Stat.StatStrength,
		Stat.StatAgility,
		// Physical
		Stat.StatAttackPower,
		Stat.StatRangedAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Spell
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatNaturePower,
		Stat.StatArcanePower,
		Stat.StatSpellCrit,
		Stat.StatMP5,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
		PseudoStat.PseudoStatRangedDps,
		PseudoStat.PseudoStatMeleeSpeedMultiplier,
		PseudoStat.PseudoStatRangedSpeedMultiplier,
	],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		// Primary
		Stat.StatMana,
		// Attributes
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		// Physical
		Stat.StatAttackPower,
		Stat.StatRangedAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Spell
		Stat.StatSpellDamage,
		Stat.StatNaturePower,
		Stat.StatArcanePower,
		Stat.StatSpellCrit,
		Stat.StatMP5,
	],
	displayPseudoStats: [PseudoStat.PseudoStatMeleeSpeedMultiplier, PseudoStat.PseudoStatRangedSpeedMultiplier],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultBuild.gear!.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 0.12,
				[Stat.StatAgility]: 1.0,
				[Stat.StatStamina]: 0.0,
				[Stat.StatIntellect]: 0.02,
				[Stat.StatAttackPower]: 1,
				[Stat.StatRangedAttackPower]: 1.0,
				[Stat.StatMeleeHit]: 15.0,
				[Stat.StatMeleeCrit]: 10.22,
				[Stat.StatExpertise]: 4.13,
				[Stat.StatSpellPower]: 0.11,
				[Stat.StatNaturePower]: 0.01,
				[Stat.StatArcanePower]: 0.1,
				[Stat.StatSpellCrit]: 0.02,
				[Stat.StatMP5]: 0.05,
				[Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 3.0,
				[PseudoStat.PseudoStatOffHandDps]: 1.66,
				[PseudoStat.PseudoStatRangedDps]: 2.61,
				[PseudoStat.PseudoStatMeleeSpeedMultiplier]: 3.92,
				[PseudoStat.PseudoStatRangedSpeedMultiplier]: 1.43,
				[PseudoStat.PseudoStatTimewornBonus]: 11.06,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultBuild.settings!.consumes!,
		// Default talents.
		talents: Presets.DefaultBuild.talents!.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultBuild.settings!.specOptions! as SpecOptions<Spec.SpecHunter>,
		other: {
			// Default distance from target.
			distanceFromTarget: Presets.DefaultBuild.settings?.playerOptions?.distanceFromTarget,
			profession1: Presets.DefaultBuild.settings?.playerOptions?.profession1,
			profession2: Presets.DefaultBuild.settings?.playerOptions?.profession2,
		},
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultBuild.settings!.raidBuffs!,
		partyBuffs: Presets.DefaultBuild.settings!.partyBuffs!,
		individualBuffs: Presets.DefaultBuild.settings!.buffs!,
		debuffs: Presets.DefaultBuild.settings!.debuffs!,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [HunterInputs.PetTypeInput, HunterInputs.WeaponAmmo, HunterInputs.QuiverInput],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: HunterInputs.HunterRotationConfig,
	petConsumeInputs: [ConsumablesInputs.PetAttackPowerConsumable, ConsumablesInputs.PetAgilityConsumable, ConsumablesInputs.PetStrengthConsumable],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.StaminaBuff],
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
		talents: [
			...Presets.TalentPresets[Phase.Phase8],
			...Presets.TalentPresets[Phase.Phase7],
			//...Presets.TalentPresets[Phase.Phase6],
			//...Presets.TalentPresets[Phase.Phase5],
			//...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase8],
			...Presets.APLPresets[Phase.Phase7],
			//...Presets.APLPresets[Phase.Phase6],
			//...Presets.APLPresets[Phase.Phase5],
			//...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase8],
			...Presets.GearPresets[Phase.Phase7],
			//...Presets.GearPresets[Phase.Phase6],
			//...Presets.GearPresets[Phase.Phase5],
			//...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
		builds: [Presets.PresetBuildPhase7NaxxMeleeDW, Presets.PresetBuildPhase7NaxxMelee2H, Presets.PresetBuildPhase7NaxxRanged, Presets.PresetBuildPhase7NaxxWeave, 
			Presets.PresetBuildPhase8NaxxMeleeDW, Presets.PresetBuildPhase8NaxxMelee2H, Presets.PresetBuildPhase8NaxxRanged, Presets.PresetBuildPhase8NaxxWeave],
	},

	autoRotation: player => {
		const level = player.getLevel();
		const isMelee =
			player.hasRune(ItemSlot.ItemSlotWaist, HunterRune.RuneBeltMeleeSpecialist) ||
			player.hasRune(ItemSlot.ItemSlotFeet, HunterRune.RuneBootsDualWieldSpecialization) ||
			player.hasRune(ItemSlot.ItemSlotFeet, HunterRune.RuneBootsWyvernStrike);

		if (isMelee) {
			switch (level) {
				case 25:
					return Presets.APLMeleeWeavePhase1.rotation.rotation!;
				case 40:
					return Presets.APLMeleePhase2.rotation.rotation!;
				case 50:
					return Presets.APLMeleeBmPhase3.rotation.rotation!;
				case 60:
					return Presets.APLWeavePhase4.rotation.rotation!;
			}
		} else {
			switch (level) {
				case 25:
					return Presets.APLMeleeWeavePhase1.rotation.rotation!;
				case 40:
					return player.getTalentTree() === 1 ? Presets.APLRangedMmPhase2.rotation.rotation! : Presets.APLRangedBmPhase2.rotation.rotation!;
				case 50:
					return Presets.APLRangedMmPhase3.rotation.rotation!;
				case 60:
					return Presets.APLRangedPhase4.rotation.rotation!;
			}
		}

		throw new Error('Auto rotation not supported for your current configuration.');
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
