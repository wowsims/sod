import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import { CURRENT_PHASE, Phase } from "../core/constants/other";
import { IndividualSimUI, registerSpecConfig } from "../core/individual_sim_ui";
import { Player } from "../core/player";
import { Class, Faction, PartyBuffs, PseudoStat, Race, Spec, Stat } from "../core/proto/common";
import { Stats } from "../core/proto_utils/stats";
import { getSpecIcon } from "../core/proto_utils/utils";
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecTankRogue, {
    cssClass: 'rogue-sim-ui',
    cssScheme: 'rogue',
    // List any known bugs / issues here and they'll be shown on the site.
    knownIssues: [
        'Rotations are not fully optimized, especially for non-standard setups.',
    ],
    warnings: [],
    epStats: [
        Stat.StatAgility,
        Stat.StatStrength,
        Stat.StatAttackPower,
        Stat.StatMeleeHit,
        Stat.StatMeleeCrit,
        Stat.StatMeleeHaste,
        Stat.StatSpellDamage,
        Stat.StatSpellPower,
        Stat.StatNaturePower,
        Stat.StatSpellHit,
        Stat.StatSpellCrit,

        // Tank stats
        Stat.StatArmor,
        Stat.StatBonusArmor,
        Stat.StatStamina,
    ],
    epPseudoStats: [
        PseudoStat.PseudoStatMainHandDps,
        PseudoStat.PseudoStatOffHandDps,
    ],
    // Reference stat against which to caluclate EP.
    epReferenceStat: Stat.StatAttackPower,
    // Which stats to display in the Character Stats seciont, at the bottom of the lef-hand sidebar.
    displayStats: [
        Stat.StatHealth,
        Stat.StatAgility,
        Stat.StatStrength,
        Stat.StatAttackPower,
        Stat.StatSpellPower,
        Stat.StatMeleeHit,
        Stat.StatSpellHit,
        Stat.StatMeleeCrit,
        Stat.StatSpellCrit,
        Stat.StatMeleeHaste,

        // Tank stats
        Stat.StatStamina,
        Stat.StatDodge,
        Stat.StatParry,
        Stat.StatArmor,
        Stat.StatBonusArmor,
    ],

    defaults: {
        // Default equipped gear.
        gear: Presets.DefaultGear.gear,
        // Default EP weights for sorting gear in the gear picker
        epWeights: Stats.fromMap({
            [Stat.StatAgility]: 1.69,
            [Stat.StatStrength]: 1.10,
			[Stat.StatAttackPower]: 1,
			[Stat.StatSpellDamage]: .68,
			[Stat.StatNaturePower]: .68,
			[Stat.StatSpellCrit]: 2.0,
			[Stat.StatSpellHit]: 5.54,
			[Stat.StatMeleeHit]: 14.2,
			[Stat.StatMeleeCrit]: 8.64,
			[Stat.StatMeleeHaste]: 6.51,
            [Stat.StatStamina]: 0.3,
            [Stat.StatArmor]: 0.01,
            [Stat.StatBonusArmor]: 0.01,
        }, {
			[PseudoStat.PseudoStatMainHandDps]: 2.94,
			[PseudoStat.PseudoStatOffHandDps]: 2.45,
		}),

        // Default consumes settings.
        consumes: Presets.DefaultConsumes,
        // Default talents.
        talents: Presets.CombatMutilate40Talents.data,
        specOptions: Presets.DefaultOptions,
        other: Presets.OtherDefaults,
        // Default raid/party buffs settings.
        raidBuffs: Presets.DefaultRaidBuffs,
        partyBuffs: PartyBuffs.create({}),
        individualBuffs: Presets.DefaultIndividualBuffs,
        debuffs: Presets.DefaultDebuffs,
    },

    playerInputs: {
        inputs: []
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
		BuffDebuffInputs.PowerInfusion
    ],
    excludeBuffDebuffInputs: [],
    otherInputs: {
        inputs: [
            OtherInputs.TankAssignment,
            OtherInputs.IncomingHps,
            OtherInputs.HealingCadence,
            OtherInputs.HealingCadenceVariation,
            OtherInputs.BurstWindow,
            OtherInputs.HpPercentForDefensives,
            OtherInputs.InspirationUptime,
        ],
    },
    encounterPicker: {
        // Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
    },

    presets: {
        // Preset talents that the user can quickly select.
		talents: [
			...Presets.TalentPresets[Phase.Phase2],
			...Presets.TalentPresets[Phase.Phase1],
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			...Presets.APLPresets[Phase.Phase2],
			...Presets.APLPresets[Phase.Phase1],
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			...Presets.GearPresets[Phase.Phase2],
			...Presets.GearPresets[Phase.Phase1],
		],
    },

    autoRotation: player => {
        return Presets.DefaultAPLs[player.getLevel()][player.getTalentTree()].rotation.rotation!;
    },

    raidSimPresets: [
        {
            spec: Spec.SpecTankRogue,
            tooltip: 'Combat Tank',
            defaultName: 'Combat',
            iconUrl: getSpecIcon(Class.ClassRogue, 0),

            talents: Presets.DefaultTalentsCombat.data,
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
                    1: Presets.GearPresets[CURRENT_PHASE][0].gear,
                },
                [Faction.Horde]: {
                    1: Presets.GearPresets[CURRENT_PHASE][0].gear,
                },
            },
        },
    ],
});

export class TankRogueSimUI extends IndividualSimUI<Spec.SpecTankRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecTankRogue>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}