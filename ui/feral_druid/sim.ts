import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../core/components/other_inputs.js';
import { Phase } from '../core/constants/other.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import { APLRotation_Type as APLRotationType } from '../core/proto/apl.js';
import { Class, Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat, WeaponImbue } from '../core/proto/common.js';
import { Gear } from '../core/proto_utils/gear.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import { TypedEvent } from '../core/typed_event.js';
import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFeralDruid, {
	cssClass: 'feral-druid-sim-ui',
	cssScheme: 'druid',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	// All stats for which EP should be calculated.
	epStats: [
		// Primary
		Stat.StatMana,
		// Attributes
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatSpirit,
		// Physical
		Stat.StatAttackPower,
		Stat.StatFeralAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		// Spell
		Stat.StatMP5,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatBonusPhysicalDamage,
	],
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
		Stat.StatSpirit,
		// Physical
		Stat.StatAttackPower,
		Stat.StatFeralAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		// Spell
		Stat.StatMP5,
	],
	displayPseudoStats: [
		PseudoStat.PseudoStatBonusPhysicalDamage,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 2.38,
				[Stat.StatAgility]: 2.35,
				[Stat.StatAttackPower]: 1,
				[Stat.StatFeralAttackPower]: 1,
				[Stat.StatMeleeHit]: 24.46,
				[Stat.StatMeleeCrit]: 16.67,
				[Stat.StatMana]: 0.04,
				[Stat.StatIntellect]: 0.67,
				[Stat.StatSpirit]: 0.08,
				[Stat.StatMP5]: 0.46,
				[Stat.StatFireResistance]: 0.5,
			},
			{},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default rotation settings.
		rotationType: APLRotationType.TypeSimple,
		simpleRotation: Presets.DefaultRotation,
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

	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.FeralDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.IntellectBuff,
		BuffDebuffInputs.BlessingOfWisdom,
		BuffDebuffInputs.ManaSpringTotem,
		BuffDebuffInputs.JudgementOfWisdom,
	],
	excludeBuffDebuffInputs: [WeaponImbue.ElementalSharpeningStone, WeaponImbue.DenseSharpeningStone, WeaponImbue.WildStrikes],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.ReactionTime,
			// DruidInputs.AssumeBleedActive,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	itemSwapConfig: {
		itemSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged],
	},
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
		rotations: [
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
	},

	autoRotation: player => {
		return Presets.DefaultAPLs[player.getLevel()].rotation.rotation!;
	},

	// simpleRotation: (player: Player<Spec.SpecFeralDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
	// 	const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

	// 	const preroarDuration = Math.min(simple.preroarDuration, 33.0);
	// 	const preRoar = APLPrepullAction.fromJsonString(
	// 		`{"action":{"activateAura":{"auraId":{"spellId":407988}}},"doAtValue":{"const":{"val":"-${(34.0 - preroarDuration).toFixed(2)}s"}}}`,
	// 	);
	// 	const preTF = APLPrepullAction.fromJsonString(`{"action":{"castSpell":{"spellId":{"spellId":5217,"rank":1}}},"doAtValue":{"const":{"val":"-3s"}}}`);
	// 	const doRotation = APLAction.fromJsonString(
	// 		`{"catOptimalRotationAction":{"maxWaitTime":${simple.maxWaitTime.toFixed(2)},"minCombosForRip":${simple.minCombosForRip.toFixed(0)},"maintainFaerieFire":${simple.maintainFaerieFire},"useShredTrick":${simple.useShredTrick}}}`,
	// 	);

	// 	prepullActions.push(...([preroarDuration > 0 ? preRoar : null, simple.precastTigersFury ? preTF : null].filter(a => a) as Array<APLPrepullAction>));

	// 	actions.push(...([doRotation].filter(a => a) as Array<APLAction>));

	// 	return APLRotation.create({
	// 		prepullActions: prepullActions,
	// 		priorityList: actions.map(action =>
	// 			APLListItem.create({
	// 				action: action,
	// 			}),
	// 		),
	// 	});
	// },

	raidSimPresets: [
		{
			spec: Spec.SpecFeralDruid,
			tooltip: specNames[Spec.SpecFeralDruid],
			defaultName: 'Cat',
			iconUrl: getSpecIcon(Class.ClassDruid, 3),

			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
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
});

export class FeralDruidSimUI extends IndividualSimUI<Spec.SpecFeralDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralDruid>) {
		super(parentElem, player, SPEC_CONFIG);
	}

	calcArpTarget(gear: Gear): number {
		let arpTarget = 1399;

		// First handle ArP proc trinkets
		if (gear.hasTrinket(45931)) {
			arpTarget -= 751;
		} else if (gear.hasTrinket(40256)) {
			arpTarget -= 612;
		}

		// Then check for Executioner enchant
		const weapon = gear.getEquippedItem(ItemSlot.ItemSlotMainHand);

		if (weapon != null && weapon!.enchant != null && weapon!.enchant!.effectId == 3225) {
			arpTarget -= 120;
		}

		return arpTarget;
	}

	calcCritCap(gear: Gear): Stats {
		const baseCritCapPercentage = 77.8; // includes 3% Crit debuff
		let agiProcs = 0;

		if (gear.hasRelic(47668)) {
			agiProcs += 200;
		}

		if (gear.hasRelic(50456)) {
			agiProcs += 44 * 5;
		}

		if (gear.hasTrinket(47131) || gear.hasTrinket(47464)) {
			agiProcs += 510;
		}

		if (gear.hasTrinket(47115) || gear.hasTrinket(47303)) {
			agiProcs += 450;
		}

		if (gear.hasTrinket(44253) || gear.hasTrinket(42987)) {
			agiProcs += 300;
		}

		return new Stats().withStat(Stat.StatMeleeCrit, (baseCritCapPercentage - (agiProcs * 1.1 * 1.06 * 1.02) / 83.33) * 45.91);
	}

	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		return Stats.fromProto(this.player.getCurrentStats().finalStats);
	}
}
