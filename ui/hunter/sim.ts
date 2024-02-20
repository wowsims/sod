import { Phase } from '../core/constants/other.js';
import {
	Class,
	Debuffs,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	Race,
	RaidBuffs,
	RangedWeaponType,
	Spec,
	Stat, PseudoStat,
	TristateEffect,
} from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import * as HunterInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecHunter, {
	cssClass: 'hunter-sim-ui',
	cssScheme: 'hunter',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],
	warnings: [
	],

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
		Stat.StatArmorPenetration,
		Stat.StatMP5,
		Stat.StatSpellPower,
		Stat.StatNaturePower,
		Stat.StatArcanePower,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatRangedDps,
	],
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
		Stat.StatArmorPenetration,
		Stat.StatMP5,
		Stat.StatSpellPower,
		Stat.StatNaturePower,
	],
	modifyDisplayStats: (player: Player<Spec.SpecHunter>) => {
		let stats = new Stats();
		stats = stats.addStat(Stat.StatMeleeCrit, player.getTalents().lethalShots * 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);

		const rangedWeapon = player.getEquippedItem(ItemSlot.ItemSlotRanged);
		if (rangedWeapon?.enchant?.effectId == 3608) {
			stats = stats.addStat(Stat.StatMeleeCrit, 40);
		}
		if (player.getRace() == Race.RaceDwarf && rangedWeapon?.item.rangedWeaponType == RangedWeaponType.RangedWeaponTypeGun) {
			stats = stats.addStat(Stat.StatMeleeCrit, 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
		}
		if (player.getRace() == Race.RaceTroll && rangedWeapon?.item.rangedWeaponType == RangedWeaponType.RangedWeaponTypeBow) {
			stats = stats.addStat(Stat.StatMeleeCrit, 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
		}

		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatStamina]: 0.5,
			[Stat.StatAgility]: 2.65,
			[Stat.StatIntellect]: 1.1,
			[Stat.StatAttackPower]: 1.0,
			[Stat.StatRangedAttackPower]: 1.0,
			[Stat.StatMeleeHit]: 2,
			[Stat.StatMeleeCrit]: 1.5,
			[Stat.StatMeleeHaste]: 1.39,
			[Stat.StatArmorPenetration]: 1.32,
			[Stat.StatSpellPower]: 0.1,
			[Stat.StatNaturePower]: 0.1,
		}, {
			[PseudoStat.PseudoStatRangedDps]: 6.32,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default rotation settings.
		simpleRotation: Presets.DefaultSimpleRotation,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			aspectOfTheLion: true,
			arcaneBrilliance: true,
			powerWordFortitude: TristateEffect.TristateEffectImproved,
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			battleShout: TristateEffect.TristateEffectImproved,
		}),
		partyBuffs: PartyBuffs.create({
		}),
		individualBuffs: IndividualBuffs.create({
			blessingOfWisdom: TristateEffect.TristateEffectImproved,
			blessingOfMight: TristateEffect.TristateEffectImproved,
		}),
		debuffs: Debuffs.create({
			homunculi: 70, // 70% average uptime default
			faerieFire: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		HunterInputs.PetTypeInput,
		HunterInputs.WeaponAmmo,
		HunterInputs.QuiverInput,
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: HunterInputs.HunterRotationConfig,
	petConsumeInputs: [
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.StaminaBuff,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
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
			...Presets.TalentPresets[Phase.Phase1],
			...Presets.TalentPresets[Phase.Phase2],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase1],
			...Presets.APLPresets[Phase.Phase2],
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase1],
			...Presets.GearPresets[Phase.Phase2],
		],
	},

	autoRotation: (player) => {
		return Presets.DefaultAPLs[player.getLevel()][player.getTalentTree()].rotation.rotation!;
	},
	
	raidSimPresets: [
		{
			spec: Spec.SpecHunter,
			tooltip: 'Beast Mastery Hunter',
			defaultName: 'Beast Mastery',
			iconUrl: getSpecIcon(Class.ClassHunter, 0),

			talents: Presets.DefaultTalentsBeastMastery.data,
			specOptions: Presets.BMDefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase1][0].gear,
				},
				[Faction.Horde]: {
					1: Presets.GearPresets[Phase.Phase1][0].gear,
				},
			},
		},
		{
			spec: Spec.SpecHunter,
			tooltip: 'Marksmanship Hunter',
			defaultName: 'Marksmanship',
			iconUrl: getSpecIcon(Class.ClassHunter, 1),
			talents: Presets.DefaultTalentsMarksman.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase1][1].gear,
				},
				[Faction.Horde]: {
					1: Presets.GearPresets[Phase.Phase1][1].gear,
				},
			},
		},
		{
			spec: Spec.SpecHunter,
			tooltip: 'Survival Hunter',
			defaultName: 'Survival',
			iconUrl: getSpecIcon(Class.ClassHunter, 2),

			talents: Presets.DefaultTalentsSurvival.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase1][2].gear,
				},
				[Faction.Horde]: {
					1: Presets.GearPresets[Phase.Phase1][2].gear,
				},
			},
		},
	],
});

export class HunterSimUI extends IndividualSimUI<Spec.SpecHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHunter>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
