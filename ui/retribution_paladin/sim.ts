import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, HandType, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import * as RetributionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecRetributionPaladin, {
	cssClass: 'retribution-paladin-sim-ui',
	cssScheme: 'paladin',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [
		(simUI: IndividualSimUI<Spec.SpecRetributionPaladin>) => {
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
		// Attributes
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		// Physical
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Spell
		Stat.StatSpellPower,
		Stat.StatSpellDamage,
		Stat.StatHolyPower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatSpellHaste,
		Stat.StatMP5,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatMeleeSpeedMultiplier, PseudoStat.PseudoStatTimewornBonus],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		// Primary
		Stat.StatMana,
		// Attributes
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		// Physical
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		// Spell
		Stat.StatSpellHaste,
		Stat.StatSpellDamage,
		Stat.StatHolyPower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatMP5,
	],
	displayPseudoStats: [],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 2.53,
				[Stat.StatAgility]: 1.13,
				[Stat.StatIntellect]: 0.15,
				[Stat.StatSpellPower]: 0.32,
				[Stat.StatSpellHit]: 0.41,
				[Stat.StatSpellCrit]: 0.01,
				[Stat.StatSpellHaste]: 0.12,
				[Stat.StatMP5]: 0.05,
				[Stat.StatAttackPower]: 1,
				[Stat.StatMeleeHit]: 1.96,
				[Stat.StatMeleeCrit]: 1.16,
				[Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 7.33,
				[PseudoStat.PseudoStatMeleeSpeedMultiplier]: 7.33,
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
	playerIconInputs: [RetributionPaladinInputs.PrimarySealSelection, RetributionPaladinInputs.AuraSelection],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
			RetributionPaladinInputs.CrusaderStrikeStopAttack,
			RetributionPaladinInputs.JudgementStopAttack,
			RetributionPaladinInputs.DivineStormStopAttack,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		rotations: [
			...Presets.APLPresets[Phase.Phase6],
			...Presets.APLPresets[Phase.Phase5],
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		// Preset talents that the user can quickly select.
		talents: [
			...Presets.TalentPresets[Phase.Phase6],
			...Presets.TalentPresets[Phase.Phase5],
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
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
		builds: [
			Presets.PresetBuildTwisting,
			Presets.PresetBuildP5SealStacking,
			Presets.PresetBuildP5Exodin,
			Presets.PresetBuildP5Shockadin,
			Presets.PresetBuildP6Exodin,
		],
	},

	autoRotation: player => {
		const level = player.getLevel();
		if (level < 60) {
			return Presets.DefaultAPLs[level].rotation.rotation!;
		}

		if (player.getTalents().holyShock) {
			return Presets.APLShockadin.rotation.rotation!;
		}

		const gear = player.getGear();
		const itemSlots = gear.getItemSlots();
		let coreForgedCount = 0;
		for (let i = 0; i < itemSlots.length; i++) {
			const item = gear.getEquippedItem(itemSlots[i]);
			if (!item) continue;
			if (item.item.setName === 'Lawbringer Radiance') {
				coreForgedCount++;
			}
		}

		const mainHand = gear.getEquippedItem(ItemSlot.ItemSlotMainHand);

		if (mainHand?.item) {
			if (mainHand.item.handType == HandType.HandTypeOneHand) {
				return Presets.APLP6OneHand.rotation.rotation!;
			} else if (mainHand.item.weaponSpeed < 3) {
				if (coreForgedCount >= 6) {
					return Presets.APLP5Exodin6CF.rotation.rotation!;
				} else {
					return Presets.APLP6Exodin.rotation.rotation!;
				}
			}
		}

		if (coreForgedCount >= 6) {
			return Presets.APLP5SealStacking6CF.rotation.rotation!;
		}

		return Presets.APLTwisting.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecRetributionPaladin,
			tooltip: 'Retribution Paladin',
			defaultName: 'Retribution',
			iconUrl: getSpecIcon(Class.ClassPaladin, 2),

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Horde]: Race.RaceUnknown,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.GearPresets[Phase.Phase1][0].gear,
					2: Presets.GearPresets[Phase.Phase2][0].gear,
					3: Presets.GearPresets[Phase.Phase3][0].gear,
					4: Presets.GearPresets[Phase.Phase5][0].gear,
					5: Presets.GearPresets[Phase.Phase5][0].gear,
					6: Presets.GearPresets[Phase.Phase6][0].gear,
				},
				[Faction.Horde]: {},
			},
		},
	],
});

export class RetributionPaladinSimUI extends IndividualSimUI<Spec.SpecRetributionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRetributionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
