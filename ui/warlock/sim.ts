import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import {
	Class,
	Faction,
	ItemSlot,
	PartyBuffs,
	Race,
	Spec,
	Stat,
} from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as WarlockInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWarlock, {
	cssClass: 'warlock-sim-ui',
	cssScheme: 'warlock',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
		"Most abilities and pets are work in progress"
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatIntellect,
		Stat.StatStamina,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatFirePower,
		Stat.StatShadowPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMP5,
		Stat.StatSpellPenetration,
	],
	// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatMana,
		Stat.StatIntellect,
		Stat.StatStamina,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatFirePower,
		Stat.StatShadowPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatSpellPenetration,
		Stat.StatMP5,
	],
	// TODO: Figure out a way to get the stat but right now this comes out wrong
	// due to pet scaling and player getting some dynamic buffs which we cant get here
	// modifyDisplayStats: (player: Player<Spec.SpecWarlock>) => {
	// 	let stats = new Stats();
		
	// 	// Demonic Knowledge rune
	// 	if (player.getEquippedItem(ItemSlot.ItemSlotFeet)?.rune?.id == WarlockRune.RuneBootsDemonicKnowledge) {
	// 		let petIntStaMap = new Map<number, Map<WarlockOptions_Summon, number>>([
	// 			[25, new Map<WarlockOptions_Summon, number>([
	// 				[WarlockOptions_Summon.Imp, 49 + 94],
	// 				[WarlockOptions_Summon.Succubus, 87 + 35],
	// 			])],
	// 			[40, new Map<WarlockOptions_Summon, number>([
	// 				[WarlockOptions_Summon.Imp, 67 + 163],
	// 				[WarlockOptions_Summon.Succubus, 148 + 49],
	// 			])],
	// 			[50, new Map<WarlockOptions_Summon, number>([
	// 				[WarlockOptions_Summon.Imp, 67 + 163],
	// 				[WarlockOptions_Summon.Succubus, 148 + 49],
	// 			])],
	// 			[60, new Map<WarlockOptions_Summon, number>([
	// 				[WarlockOptions_Summon.Imp, 67 + 163],
	// 				[WarlockOptions_Summon.Succubus, 148 + 49],
	// 			])],
	// 		]);

	// 		// Base stats
	// 		let currentTotal = petIntStaMap.get(player.getLevel())!.get(player.getSpecOptions().summon)!;

	// 		// Bonus item stats
	// 		let trinketId = 216509
	// 		if (player.getEquippedItem(ItemSlot.ItemSlotTrinket1)?.id == trinketId || player.getEquippedItem(ItemSlot.ItemSlotTrinket2)?.id == trinketId) {
	// 			currentTotal = currentTotal + 100;
	// 		}

	// 		// Player scaled stats
	// 		let playerStats = Stats.fromProto(player.getCurrentStats().finalStats)
	// 		currentTotal = currentTotal + playerStats.getStat(Stat.StatIntellect) * 0.3 + playerStats.getStat(Stat.StatStamina) * (player.getSpecOptions().summon == WarlockOptions_Summon.Imp ? 0.66 : 0.75)
			
	// 		stats = stats.addStat(Stat.StatSpellPower, currentTotal * 0.1);
	// 	}

	// 	return {
	// 		talents: stats,
	// 	};
	// },

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,

		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 0.18,
			[Stat.StatSpirit]: 0.54,
			[Stat.StatSpellPower]: 1,
			[Stat.StatSpellDamage]: 1,
			[Stat.StatFirePower]: 1,
			[Stat.StatShadowPower]: 1,
			[Stat.StatSpellHit]: 0.93,
			[Stat.StatSpellCrit]: 0.53,
			[Stat.StatSpellHaste]: 0.81,
			[Stat.StatStamina]: 0.01,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,

		// Default buffs and debuffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,

		partyBuffs: PartyBuffs.create({}),

		individualBuffs: Presets.DefaultIndividualBuffs,

		debuffs: Presets.DefaultDebuffs,

		other: Presets.OtherDefaults,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		WarlockInputs.PetInput,
		WarlockInputs.ImpFireboltRank,
		WarlockInputs.ArmorInput,
		WarlockInputs.WeaponImbueInput,
	],

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		// Physical buffs that affect pets
		BuffDebuffInputs.MajorArmorDebuff,
		BuffDebuffInputs.CurseOfRecklessness,
		BuffDebuffInputs.FaerieFire,
		BuffDebuffInputs.PaladinPhysicalBuff,
		BuffDebuffInputs.StrengthBuffHorde,
		BuffDebuffInputs.BattleShoutBuff,
		BuffDebuffInputs.TrueshotAuraBuff,
		BuffDebuffInputs.MeleeCritBuff,
		BuffDebuffInputs.CurseOfVulnerability,
		BuffDebuffInputs.GiftOfArthas,
		BuffDebuffInputs.CrystalYield,
		BuffDebuffInputs.AncientCorrosivePoison,
	],
	excludeBuffDebuffInputs: [
		BuffDebuffInputs.BleedDebuff,
		BuffDebuffInputs.SpellWintersChillDebuff,
		...ConsumablesInputs.FROST_POWER_CONFIG,
	],
	petConsumeInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.DistanceFromTarget,
			OtherInputs.ChannelClipDelay,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			...Presets.TalentPresets[Phase.Phase1],
			...Presets.TalentPresets[CURRENT_PHASE],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase1],
			...Presets.APLPresets[CURRENT_PHASE],
		],

		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase1],
			...Presets.GearPresets[CURRENT_PHASE],
		],
	},

	autoRotation: (player) => {
		return Presets.DefaultAPLs[player.getLevel()][player.getTalentTree()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Affliction DPS',
			defaultName: 'Affliction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 0),

			talents: Presets.DefaultTalentsAffliction.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
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
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Demonology DPS',
			defaultName: 'Demonology',
			iconUrl: getSpecIcon(Class.ClassWarlock, 1),

			talents: Presets.DefaultTalentsDemonology.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
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
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Destruction DPS',
			defaultName: 'Destruction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 2),

			talents: Presets.DefaultTalentsDestruction.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
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
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class WarlockSimUI extends IndividualSimUI<Spec.SpecWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarlock>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
