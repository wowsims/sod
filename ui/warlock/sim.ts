import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as WarlockInputs from '../core/components/inputs/warlock_inputs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { WarlockRune } from '../core/proto/warlock';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWarlock, {
	cssClass: 'warlock-sim-ui',
	cssScheme: 'warlock',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		// Primary
		Stat.StatHealth,
		Stat.StatMana,
		// Attributes
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatSpirit,
		// Spell
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatFirePower,
		Stat.StatShadowPower,
		Stat.StatMP5,
	],
	epPseudoStats: [PseudoStat.PseudoStatCastSpeedMultiplier],
	// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		// Primary
		Stat.StatMana,
		// Attributes
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatSpirit,
		// Physical
		Stat.StatAttackPower,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHit,
		// Spell
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatFirePower,
		Stat.StatShadowPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMP5,
	],
	displayPseudoStats: [PseudoStat.PseudoStatCastSpeedMultiplier],

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
		race: Presets.OtherDefaults.race,

		// Default equipped gear.
		gear: Presets.DefaultGear.gear,

		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatIntellect]: 0.32,
				[Stat.StatSpirit]: 0.63,
				[Stat.StatSpellPower]: 1,
				[Stat.StatSpellDamage]: 1,
				[Stat.StatFirePower]: 1,
				[Stat.StatShadowPower]: 1,
				[Stat.StatSpellHit]: 10.69,
				[Stat.StatSpellCrit]: 16.93,
				[Stat.StatStamina]: 0.01,
				[Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatCastSpeedMultiplier]: 3.47,
				[PseudoStat.PseudoStatTimewornBonus]: 21.59,
			},
		),
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
	playerIconInputs: [WarlockInputs.PetInput(), WarlockInputs.ImpFireboltRank(), WarlockInputs.ArmorInput(), WarlockInputs.WeaponImbueInput()],

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		// Physical buffs that affect pets
		BuffDebuffInputs.MajorArmorDebuff,
		BuffDebuffInputs.CurseOfRecklessness,
		BuffDebuffInputs.FaerieFire,
		BuffDebuffInputs.PaladinPhysicalBuff,
		BuffDebuffInputs.StrengthBuffHorde,
		BuffDebuffInputs.BattleShoutBuff,
		BuffDebuffInputs.GraceOfAir,
		BuffDebuffInputs.MeleeCritBuff,
		BuffDebuffInputs.BattleSquawkBuff,
		BuffDebuffInputs.CurseOfVulnerability,
		BuffDebuffInputs.GiftOfArthas,
		BuffDebuffInputs.HolySunder,
		BuffDebuffInputs.SpellWintersChillDebuff,
	],
	excludeBuffDebuffInputs: [BuffDebuffInputs.BleedDebuff],
	petConsumeInputs: [ConsumablesInputs.PetAttackPowerConsumable, ConsumablesInputs.PetAgilityConsumable, ConsumablesInputs.PetStrengthConsumable],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [WarlockInputs.PetPoolManaInput(), OtherInputs.DistanceFromTarget, OtherInputs.ChannelClipDelay],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			...Presets.TalentPresets[Phase.Phase7],
			...Presets.TalentPresets[Phase.Phase6],
			...Presets.TalentPresets[Phase.Phase5],
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase7],
			...Presets.APLPresets[Phase.Phase6],
			...Presets.APLPresets[Phase.Phase5],
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],

		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase7],
			...Presets.GearPresets[Phase.Phase6],
			...Presets.GearPresets[Phase.Phase5],
			...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
		// Preset builds (gear, talents, APL) that the user can quickly select.
		builds: [Presets.PresetBuildAff, Presets.PresetBuildDemo, Presets.PresetBuildDestro],
	},

	autoRotation: player => {
		const level = player.getLevel();
		if (level < 50) {
			return Presets.DefaultAPLs[player.getLevel()][player.getTalentTree()].rotation.rotation!;
		}

		const hasIncinerate = player.getEquippedItem(ItemSlot.ItemSlotWrist)?.rune?.id == WarlockRune.RuneBracerIncinerate;
		const specNumber = hasIncinerate ? 2 : 0;
		return Presets.DefaultAPLs[level][specNumber].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWarlock,
			tooltip: 'Destruction DPS',
			defaultName: 'Destruction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 2),

			talents: Presets.DestroMgiTalentsPhase2.data,
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
					1: Presets.FireImpGearPhase2.gear,
				},
				[Faction.Horde]: {
					1: Presets.FireImpGearPhase2.gear,
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
