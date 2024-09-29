import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import * as ProtectionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecProtectionPaladin, {
	cssClass: 'protection-paladin-sim-ui',
	cssScheme: 'paladin',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
		`Judgement of the Crusader is currently not implemented; users can manually award themselves the relevant spellpower amount
		for a dps gain that will be slightly inflated given JotC does not benefit from source damage modifiers.`,
		`Be aware that not all item and weapon enchants are currently implemented in the sim, which make some notable Retribution
		weapons like Pendulum of Doom and The Jackhammer undervalued.`,
	],
	warnings: [
		(simUI: IndividualSimUI<Spec.SpecProtectionPaladin>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					if (simUI.player.getSpecOptions().primarySeal == 0) {
						return `Your previously selected seal is no longer available because of a talent or rune change.
							No seal will be cast with this configuration. Please select an available seal in the Settings>Player menu.`;
					} else {
						return '';
					}
				},
			};
		},
	],
	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatHealth,
		Stat.StatMana,
		// Primary Attributes
		Stat.StatStrength,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatIntellect,
		// Melee Offensive
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		// Magic offensive/healing
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellPower,
		Stat.StatHolyPower,
		Stat.StatHealingPower,
		// Defensive
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatBlock,
		Stat.StatBlockValue,
		// Resistances
		Stat.StatFireResistance,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatArcaneResistance,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatMeleeSpeedMultiplier],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatMana,
		// Primary Attributes
		Stat.StatStrength,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatIntellect,
		// Melee Offensive
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		// Magic offensive/healing
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellPower,
		Stat.StatHolyPower,
		Stat.StatHealingPower,
		// Defensive
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDefense,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatBlock,
		Stat.StatBlockValue,
		// Resistances
		Stat.StatFireResistance,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatArcaneResistance,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 2.31,
				[Stat.StatAgility]: 14.7,
				[Stat.StatStamina]: 0.0,
				[Stat.StatIntellect]: 0.02,
				[Stat.StatSpellPower]: 0.4,
				[Stat.StatHolyPower]: 0.22,
				[Stat.StatSpellHit]: 4.92,
				[Stat.StatSpellCrit]: 1.81,
				[Stat.StatAttackPower]: 1.0,
				[Stat.StatMeleeHit]: 0.0,
				[Stat.StatMeleeCrit]: 27.78,
				[Stat.StatMeleeHaste]: 22.6,
				[Stat.StatMana]: 0.0,
				[Stat.StatArmor]: 1.0,
				[Stat.StatDefense]: 25.89,
				[Stat.StatBlock]: 13.64,
				[Stat.StatBlockValue]: 1.93,
				[Stat.StatDodge]: 213.3,
				[Stat.StatParry]: 212.61,
				[Stat.StatHealth]: 0.0,
				[Stat.StatArcaneResistance]: 0.0,
				[Stat.StatFireResistance]: 0.0,
				[Stat.StatFrostResistance]: 0.0,
				[Stat.StatNatureResistance]: 0.0,
				[Stat.StatShadowResistance]: 0.0,
				[Stat.StatBonusArmor]: 0.96,
				[Stat.StatHealingPower]: 0.0,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 3.33,
				[PseudoStat.PseudoStatMeleeSpeedMultiplier]: 3.33,
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
		race: Race.RaceHuman,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		ProtectionPaladinInputs.PrimarySealSelection,
		ProtectionPaladinInputs.RighteousFuryToggle,
		ProtectionPaladinInputs.BlessingSelection,
		ProtectionPaladinInputs.AuraSelection,
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellScorchDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InspirationUptime,
			OtherInputs.InFrontOfTarget,
			ProtectionPaladinInputs.ManualAutoAttacksToggle,
			//OtherInputs.DistanceFromTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		rotations: [...Presets.APLPresets[Phase.Phase4]],
		// Preset talents that the user can quickly select.
		talents: [...Presets.TalentPresets[Phase.Phase5], ...Presets.TalentPresets[Phase.Phase4]],
		// Preset gear configurations that the user can quickly select.
		gear: [...Presets.GearPresets[Phase.Phase4]],
	},

	autoRotation: player => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecProtectionPaladin,
			tooltip: 'Protection Paladin',
			defaultName: 'Protection',
			iconUrl: getSpecIcon(Class.ClassPaladin, 1),

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceUnknown,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase4][0].gear,
				},
				[Faction.Horde]: {},
			},
		},
	],
});

export class ProtectionPaladinSimUI extends IndividualSimUI<Spec.SpecProtectionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
