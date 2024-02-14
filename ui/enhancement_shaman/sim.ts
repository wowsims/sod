import { ShamanShieldInput } from '../core/components/inputs/shaman_shields.js';
import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { APLRotation } from '../core/proto/apl.js';
import {
	Class,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	Spec,
	Stat,
} from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as ShamanInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecEnhancementShaman, {
	cssClass: 'enhancement-shaman-sim-ui',
	cssScheme: 'shaman',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatIntellect,
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
		Stat.StatSpellPower,
		Stat.StatFirePower,
		Stat.StatNaturePower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatSpellHaste,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatArmorPenetration,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
	],

	defaults: {
		race: Race.RaceTroll,
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: .02,
			[Stat.StatAgility]: 1.12,
			[Stat.StatStrength]: 2.29,
			[Stat.StatSpellPower]: 1.15,
			[Stat.StatFirePower]: 0.63,
			[Stat.StatNaturePower]: 0.48,
			[Stat.StatSpellHit]: 0.03, //default EP assumes cap
			[Stat.StatSpellCrit]: 1.94,
			[Stat.StatSpellHaste]: 2.97,
			[Stat.StatAttackPower]: 1.0,
			[Stat.StatMeleeHit]: 9.62,
			[Stat.StatMeleeCrit]: 14.80,
			[Stat.StatMeleeHaste]: 11.84, //haste is complicated
			[Stat.StatArmorPenetration]: 0.35,
			[Stat.StatExpertise]: 1.92, //default EP assumes cap
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 8.15,
			[PseudoStat.PseudoStatOffHandDps]: 5.81,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({
		}),
		individualBuffs: IndividualBuffs.create({
		}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		ShamanShieldInput<Spec.SpecEnhancementShaman>(),
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.BlessingOfWisdom,
		BuffDebuffInputs.ManaSpringTotem,
		BuffDebuffInputs.SpiritBuff,
	],
	excludeBuffDebuffInputs: [
		BuffDebuffInputs.BleedDebuff,
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			ShamanInputs.SyncTypeInput,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	itemSwapConfig: {
		itemSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
		tooltip: "Swap items will be given max-rank Rockbiter Weapon",
	},
	customSections: [
		// TotemsSection,
	],
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
			Presets.GearBlank,
			...Presets.GearPresets[Phase.Phase1],
			...Presets.GearPresets[CURRENT_PHASE],
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
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase1][0].gear,
				},
				[Faction.Horde]: {
					1: Presets.GearPresets[Phase.Phase1][0].gear,
				},
			},
		},
	],
})

export class EnhancementShamanSimUI extends IndividualSimUI<Spec.SpecEnhancementShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecEnhancementShaman>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
