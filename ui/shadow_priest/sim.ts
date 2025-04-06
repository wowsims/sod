import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames, SpecOptions } from '../core/proto_utils/utils.js';
import * as ShadowPriestInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecShadowPriest, {
	cssClass: 'shadow-priest-sim-ui',
	cssScheme: 'priest',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['The Homunculi Rune is not currently implemented until more data is available'],

	// All stats for which EP should be calculated.
	epStats: [
		// Attributes
		Stat.StatIntellect,
		Stat.StatSpirit,
		// Spell
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatShadowPower,
		Stat.StatHolyPower,
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
		Stat.StatIntellect,
		Stat.StatSpirit,
		// Spell
		Stat.StatSpellDamage,
		Stat.StatShadowPower,
		Stat.StatHolyPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMP5,
	],
	displayPseudoStats: [PseudoStat.PseudoStatCastSpeedMultiplier],

	modifyDisplayStats: (player: Player<Spec.SpecShadowPriest>) => {
		let stats = new Stats();
		stats = stats.addPseudoStat(PseudoStat.PseudoStatSchoolHitShadow, player.getTalents().shadowFocus * 2 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);

		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultBuild.gear!.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatIntellect]: 0.16,
				[Stat.StatSpirit]: 0.01,
				[Stat.StatSpellPower]: 1,
				[Stat.StatSpellDamage]: 1,
				[Stat.StatShadowPower]: 1,
				[Stat.StatSpellHit]: 5.51,
				[Stat.StatSpellCrit]: 5.99, // Averaged between using and not using Despair for dot crits
				[Stat.StatMP5]: 0.0,
				[Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatCastSpeedMultiplier]: 1.73,
				[PseudoStat.PseudoStatTimewornBonus]: 26.08,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultBuild.settings!.consumes!,
		// Default talents.
		talents: Presets.DefaultBuild.talents!.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultBuild.settings!.specOptions! as SpecOptions<Spec.SpecShadowPriest>,
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
	playerIconInputs: [ShadowPriestInputs.ArmorInput],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.BlessingOfWisdom,
		BuffDebuffInputs.ManaSpringTotem,
		BuffDebuffInputs.StaminaBuff,
		BuffDebuffInputs.SpellWintersChillDebuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.TankAssignment, OtherInputs.ChannelClipDelay, OtherInputs.DistanceFromTarget],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		talents: [
			...Presets.TalentPresets[Phase.Phase8],
			...Presets.TalentPresets[Phase.Phase7],
			...Presets.TalentPresets[Phase.Phase6],
			...Presets.TalentPresets[Phase.Phase5],
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		rotations: [
			...Presets.APLPresets[Phase.Phase8],
			...Presets.APLPresets[Phase.Phase7],
			...Presets.APLPresets[Phase.Phase6],
			...Presets.APLPresets[Phase.Phase5],
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		gear: [
			...Presets.GearPresets[Phase.Phase8],
			...Presets.GearPresets[Phase.Phase7],
			...Presets.GearPresets[Phase.Phase6],
			...Presets.GearPresets[Phase.Phase5],
			...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
		builds: [Presets.PresetBuildPhase8, Presets.PresetBuildPhase7],
	},

	autoRotation: player => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecShadowPriest,
			tooltip: specNames[Spec.SpecShadowPriest],
			defaultName: 'Shadow',
			iconUrl: getSpecIcon(Class.ClassPriest, 2),

			talents: Presets.DefaultBuild.talents!.data,
			specOptions: Presets.DefaultBuild.settings!.specOptions! as SpecOptions<Spec.SpecShadowPriest>,
			consumes: Presets.DefaultBuild.settings!.consumes!,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceDwarf,
				[Faction.Horde]: Race.RaceUndead,
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
});

export class ShadowPriestSimUI extends IndividualSimUI<Spec.SpecShadowPriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecShadowPriest>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
