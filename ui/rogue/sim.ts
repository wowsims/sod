import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat, Target, WeaponType } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { HonorOfThievesCritRate } from './inputs';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecRogue, {
	cssClass: 'rogue-sim-ui',
	cssScheme: 'rogue',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['Rotations are not fully optimized, especially for non-standard setups.'],
	warnings: [
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.sim.encounter.changeEmitter,
				getContent: () => {
					const hasNoArmor = !!(simUI.sim.encounter.targets ?? []).find((target: Target) => new Stats(target.stats).getStat(Stat.StatArmor) <= 0);
					if (hasNoArmor) {
						return 'One or more targets have no armor. Check advanced encounter settings.';
					} else {
						return '';
					}
				},
			};
		},
		(simUI: IndividualSimUI<Spec.SpecRogue>) => {
			return {
				updateOn: simUI.player.changeEmitter,
				getContent: () => {
					if (simUI.player.getTalents().maceSpecialization) {
						if (
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeMace ||
							simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeMace
						) {
							return '';
						} else {
							return '"Mace Specialization" talent selected, but maces not equipped.';
						}
					} else {
						return '';
					}
				},
			};
		},
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatSpellDamage,
		Stat.StatSpellPower,
		Stat.StatNaturePower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatDefense,
		Stat.StatFireResistance,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps, PseudoStat.PseudoStatMeleeSpeedMultiplier],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatSpellDamage,
		Stat.StatMeleeHit,
		Stat.StatSpellHit,
		Stat.StatMeleeCrit,
		Stat.StatSpellCrit,
		Stat.StatMeleeHaste,
		Stat.StatDefense,
		Stat.StatFireResistance,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatAgility]: 1.23,
				[Stat.StatStrength]: 1.64,
				[Stat.StatAttackPower]: 1.0,
				[Stat.StatSpellDamage]: 1.09,
				[Stat.StatNaturePower]: 0.81,
				[Stat.StatSpellCrit]: 4.93,
				[Stat.StatSpellHit]: 11.94,
				[Stat.StatMeleeHit]: 14.11,
				[Stat.StatMeleeCrit]: 25.97,
				[Stat.StatFireResistance]: 0.5,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 4.05,
				[PseudoStat.PseudoStatOffHandDps]: 2.74,
				[PseudoStat.PseudoStatMeleeSpeedMultiplier]: 2.74,
			},
		),

		// Default consumes settings.
		consumes: Presets.DefaultConsumes[Phase.Phase4],
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

	playerInputs: {
		inputs: [],
	},
	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.SpellCritBuff,
		BuffDebuffInputs.SpellShadowWeavingDebuff,
		BuffDebuffInputs.NatureSpellDamageDebuff,
		BuffDebuffInputs.MekkatorqueFistDebuff,
		BuffDebuffInputs.SpellScorchDebuff,
		BuffDebuffInputs.PowerInfusion,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.TankAssignment, OtherInputs.InFrontOfTarget, HonorOfThievesCritRate],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			...Presets.TalentPresets[Phase.Phase4],
			...Presets.TalentPresets[Phase.Phase3],
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase4],
			...Presets.APLPresets[Phase.Phase3],
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase4],
			...Presets.GearPresets[Phase.Phase3],
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
	},

	autoRotation: player => {
		// Try to find a rotation by hand rune
		const handRuneID = player.getEquippedItem(ItemSlot.ItemSlotHands)?._rune?.id ?? 0;
		const preset = Presets.DefaultAPLs[player.getLevel()][handRuneID];

		if (preset) return preset.rotation.rotation!;

		throw new Error('Auto rotation is not supported for your level / hand rune combination. Please select an APL manually.');
	},

	raidSimPresets: [
		{
			spec: Spec.SpecRogue,
			tooltip: 'Assassination Rogue',
			defaultName: 'Assassination',
			iconUrl: getSpecIcon(Class.ClassRogue, 0),

			talents: Presets.DefaultTalentsAssassin.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes[Phase.Phase1],
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
		},
		{
			spec: Spec.SpecRogue,
			tooltip: 'Combat Rogue',
			defaultName: 'Combat',
			iconUrl: getSpecIcon(Class.ClassRogue, 1),

			talents: Presets.DefaultTalentsCombat.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes[Phase.Phase1],
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
		},
	],
});

export class RogueSimUI extends IndividualSimUI<Spec.SpecRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRogue>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
