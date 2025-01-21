import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs';
import { Phase } from '../core/constants/other';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui';
import { Player } from '../core/player';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common';
import { Stats } from '../core/proto_utils/stats';
import { getSpecIcon, specNames } from '../core/proto_utils/utils';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecElementalShaman, {
	cssClass: 'elemental-shaman-sim-ui',
	cssScheme: 'shaman',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	// All stats for which EP should be calculated.
	epStats: [
		// Attributes
		Stat.StatIntellect,
		Stat.StatSpirit,
		// Spell
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatFirePower,
		Stat.StatNaturePower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMP5,
	],
	epPseudoStats: [PseudoStat.PseudoStatCastSpeedMultiplier],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		// Primary
		Stat.StatMana,
		// Attributes
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatSpirit,
		// Spell
		Stat.StatSpellDamage,
		Stat.StatNaturePower,
		Stat.StatFirePower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMP5,
	],
	displayPseudoStats: [PseudoStat.PseudoStatCastSpeedMultiplier],

	defaults: {
		race: Race.RaceTroll,
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatIntellect]: 2.3,
				[Stat.StatSpellPower]: 1,
				[Stat.StatSpellDamage]: 1,
				[Stat.StatFirePower]: 0.3,
				[Stat.StatNaturePower]: 0.7,
				[Stat.StatSpellHit]: 20.08,
				[Stat.StatSpellCrit]: 17.91,
				[Stat.StatMP5]: 0.02,
				[Stat.StatStrength]: 0.01,
				[Stat.StatAttackPower]: 0.01,
				[Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatCastSpeedMultiplier]: 19.64,
				[PseudoStat.PseudoStatTimewornBonus]: 35.74,
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
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [BuffDebuffInputs.BleedDebuff],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.DistanceFromTarget],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
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
			...Presets.TalentPresets[Phase.Phase6],
			...Presets.TalentPresets[Phase.Phase5],
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase6],
			...Presets.APLPresets[Phase.Phase5],
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase6],
			...Presets.GearPresets[Phase.Phase5],
			...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
		// Preset build configurations that the user can quickly select
		builds: [Presets.PresetBuildPhase6, Presets.PresetBuildPhase5, Presets.PresetBuildPhase4],
	},

	autoRotation: player => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecElementalShaman,
			tooltip: specNames[Spec.SpecElementalShaman],
			defaultName: 'Elemental',
			iconUrl: getSpecIcon(Class.ClassShaman, 0),

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceUnknown,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase1][0].gear,
					2: Presets.GearPresets[Phase.Phase2][0].gear,
				},
				[Faction.Horde]: {
					1: Presets.GearPresets[Phase.Phase1][0].gear,
					2: Presets.GearPresets[Phase.Phase2][0].gear,
				},
			},
		},
	],
});

export class ElementalShamanSimUI extends IndividualSimUI<Spec.SpecElementalShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecElementalShaman>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
