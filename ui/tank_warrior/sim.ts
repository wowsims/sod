import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as WarriorInputs from '../core/components/inputs/warrior_inputs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
// import * as TankWarriorInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecTankWarrior, {
	cssClass: 'tank-warrior-sim-ui',
	cssScheme: 'warrior',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		// Primary
		Stat.StatHealth,
		// Attributes
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatStamina,
		// Physical
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Defense
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatBlock,
		Stat.StatBlockValue,
		Stat.StatParry,
		Stat.StatDodge,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps, PseudoStat.PseudoStatMeleeSpeedMultiplier],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		// Attributes
		Stat.StatStrength,
		Stat.StatAgility,
		// Physical
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Defense
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatBlock,
		Stat.StatBlockValue,
		Stat.StatDodge,
		Stat.StatParry,
		// Resistances
		Stat.StatShadowResistance,
	],
	displayPseudoStats: [PseudoStat.PseudoStatThornsDamage, PseudoStat.PseudoStatMeleeSpeedMultiplier],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatArmor]: 0.174,
				[Stat.StatBonusArmor]: 0.155,
				[Stat.StatStamina]: 2.336,
				[Stat.StatStrength]: 1.555,
				[Stat.StatAgility]: 2.771,
				[Stat.StatAttackPower]: 0.32,
				[Stat.StatMeleeHit]: 1.432,
				[Stat.StatMeleeCrit]: 0.925,
				[Stat.StatBlock]: 1.32,
				[Stat.StatBlockValue]: 1.373,
				[Stat.StatDodge]: 2.606,
				[Stat.StatParry]: 2.649,
				[Stat.StatDefense]: 3.305,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 11.92,
				[PseudoStat.PseudoStatOffHandDps]: 4.69,
				[PseudoStat.PseudoStatMeleeSpeedMultiplier]: 4.69,
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
	playerIconInputs: [WarriorInputs.ShoutPicker<Spec.SpecTankWarrior>(), WarriorInputs.StancePicker<Spec.SpecTankWarrior>()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarriorInputs.StartingRage<Spec.SpecTankWarrior>(),
			WarriorInputs.QueueDelay<Spec.SpecTankWarrior>(),
			WarriorInputs.StanceSnapshot<Spec.SpecTankWarrior>(),
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
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
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
		// Preset builds that the user can quickly select.
		builds: [Presets.PresetBuildTanky, Presets.PresetBuildDamage],
	},

	autoRotation: player => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecTankWarrior,
			tooltip: 'Tank Warrior',
			defaultName: 'Tank',
			iconUrl: getSpecIcon(Class.ClassWarrior, 2),

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {},
				[Faction.Horde]: {},
			},
		},
	],
});

export class TankWarriorSimUI extends IndividualSimUI<Spec.SpecTankWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecTankWarrior>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
