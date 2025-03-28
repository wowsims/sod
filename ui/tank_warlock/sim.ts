import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as WarlockInputs from '../core/components/inputs/warlock_inputs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { WarlockRune } from '../core/proto/warlock.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecTankWarlock, {
	cssClass: 'tank-warlock-sim-ui',
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
		// Defense
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatDodge,
	],
	epPseudoStats: [PseudoStat.PseudoStatMeleeSpeedMultiplier, PseudoStat.PseudoStatCastSpeedMultiplier],
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
		// Defense
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatDodge,
		// Resistances
		Stat.StatShadowResistance,
	],
	displayPseudoStats: [PseudoStat.PseudoStatThornsDamage, PseudoStat.PseudoStatMeleeSpeedMultiplier, PseudoStat.PseudoStatCastSpeedMultiplier],

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
			[Stat.StatStamina]: 0.01,
			[Stat.StatDefense]: 1.5,
			[Stat.StatFireResistance]: 0.5,
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
	playerIconInputs: [WarlockInputs.PetInput(), WarlockInputs.ArmorInput()],

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellWintersChillDebuff],
	excludeBuffDebuffInputs: [BuffDebuffInputs.BleedDebuff],
	petConsumeInputs: [ConsumablesInputs.PetAttackPowerConsumable, ConsumablesInputs.PetAgilityConsumable, ConsumablesInputs.PetStrengthConsumable],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarlockInputs.PetPoolManaInput(),
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InspirationUptime,
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
			...Presets.TalentPresets[Phase.Phase5],
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase5],
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],

		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase5],
			...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
		builds: [Presets.PresetBuildAff, Presets.PresetBuildDemo, Presets.PresetBuildDestro],
	},

	autoRotation: player => {
		let specNumber = player.getTalentTree();
		const level = player.getLevel();

		// Pre-60 had fewer options so it basically came down to master channeler vs no master channeler
		if (level < 60) {
			const hasMasterChanneler = player.getEquippedItem(ItemSlot.ItemSlotChest)?.rune?.id == WarlockRune.RuneChestMasterChanneler;
			specNumber = hasMasterChanneler ? 0 : 1;
			return Presets.DefaultAPLs[level][specNumber].rotation.rotation!;
		}

		return Presets.DefaultAPLs[level][specNumber].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecTankWarlock,
			tooltip: 'Demonology Tank',
			defaultName: 'Demonology',
			iconUrl: getSpecIcon(Class.ClassWarlock, 0),

			talents: Presets.TalentsDemonologyTankPhase2.data,
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
					1: Presets.GearDemonologyTankPhase2.gear,
				},
				[Faction.Horde]: {
					1: Presets.GearDemonologyTankPhase2.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
		{
			spec: Spec.SpecTankWarlock,
			tooltip: 'Destruction Tank',
			defaultName: 'Destruction',
			iconUrl: getSpecIcon(Class.ClassWarlock, 2),

			talents: Presets.TalentsDestructionTankPhase2.data,
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
					1: Presets.GearDestructionTankPhase2.gear,
				},
				[Faction.Horde]: {
					1: Presets.GearDestructionTankPhase2.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class TankWarlockSimUI extends IndividualSimUI<Spec.SpecTankWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecTankWarlock>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
