import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import * as ShadowPriestInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecShadowPriest, {
	cssClass: 'shadow-priest-sim-ui',
	cssScheme: 'priest',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['The Homunculi Rune is not currently implemented until more data is available'],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatShadowPower,
		Stat.StatHolyPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMP5,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatMana,
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatSpellDamage,
		Stat.StatShadowPower,
		Stat.StatHolyPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMP5,
		Stat.StatFireResistance,
	],
	modifyDisplayStats: (player: Player<Spec.SpecShadowPriest>) => {
		let stats = new Stats();
		stats = stats.addPseudoStat(PseudoStat.PseudoStatSchoolHitShadow, player.getTalents().shadowFocus * 2 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);

		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 0.16,
			[Stat.StatSpirit]: 0.01,
			[Stat.StatSpellPower]: 1,
			[Stat.StatSpellDamage]: 1,
			[Stat.StatShadowPower]: 1,
			[Stat.StatSpellHit]: 5.51,
			[Stat.StatSpellCrit]: 5.99, // Averaged between using and not using Despair for dot crits
			[Stat.StatSpellHaste]: 1.65,
			[Stat.StatMP5]: 0.0,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,

		partyBuffs: PartyBuffs.create({}),

		individualBuffs: Presets.DefaultIndividualBuffs,

		debuffs: Presets.DefaultDebuffs,

		other: Presets.OtherDefaults,
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
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		talents: [
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		rotations: [
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		gear: [
			...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
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

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
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
