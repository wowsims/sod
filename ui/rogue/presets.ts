import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import {
	AgilityElixir,
	Consumes,
	Debuffs,
	IndividualBuffs,
	Profession,
	RaidBuffs,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	RogueOptions,
} from '../core/proto/rogue.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';
import P1Daggers from './gear_sets/p1_daggers.gear.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P2DaggersGear from './gear_sets/p2_daggers.gear.json';

import MutilateApl from './apls/mutilate.apl.json';
import MutilateEAApl from './apls/mutilateEA.apl.json';
import SinisterApl25 from './apls/basic_strike_25.apl.json';
import ARBFApl from './apls/ARBF.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearDaggersP1 = PresetUtils.makePresetGear('P1 Daggers', P1Daggers)
export const GearCombatP1 = PresetUtils.makePresetGear("P1 Combat", P1CombatGear)
export const GearDaggersP2 = PresetUtils.makePresetGear("P2 Daggers", P2DaggersGear)

export const GearPresets = {
  [Phase.Phase1]: [
    GearDaggersP1,
	GearCombatP1,
  ],
  [Phase.Phase2]: [
	GearDaggersP2,
  ]
};

export const DefaultGear = GearPresets[CURRENT_PHASE][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { talentTree: 0 });
export const ROTATION_PRESET_MUTILATEEA = PresetUtils.makePresetAPLRotation('Mutilate EA', MutilateEAApl, { talentTree: 0 });
export const ROTATION_PRESET_SINISTER_25 = PresetUtils.makePresetAPLRotation('Sinister', SinisterApl25, { talentTree: 1 });
export const ROTATION_PRESET_ARBF_40 = PresetUtils.makePresetAPLRotation('AR/BF', ARBFApl, { talentTree: 1 });

export const APLPresets = {
  [Phase.Phase1]: [
    ROTATION_PRESET_MUTILATE,
	ROTATION_PRESET_SINISTER_25,
  ],
  [Phase.Phase2]: [
	ROTATION_PRESET_MUTILATE,
	ROTATION_PRESET_MUTILATEEA,
	ROTATION_PRESET_ARBF_40,
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
	},
  40: {
		0: APLPresets[Phase.Phase2][1],
		1: APLPresets[Phase.Phase2][1],
		2: APLPresets[Phase.Phase2][2],
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const CombatDagger25Talents = {
	name: 'P1 Combat Dagger',
	data: SavedTalents.create({
		talentsString: '-023305002001',
	}),
};
export const ColdBloodMutilate40Talents = {
	name: 'P2 CB Mutilate',
	data: SavedTalents.create({
		talentsString: '005303103551--05'
	})
};

export const ARBF40Talents = {
	name: 'P2 AR BF Mutilate',
	data: SavedTalents.create({
		talentsString: '-0053052020550100201',
	}),
};

export const TalentPresets = {
	[Phase.Phase1]: [
		CombatDagger25Talents,
	],
	[Phase.Phase2]: [
		ColdBloodMutilate40Talents,
		ARBF40Talents,
	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsAssassin = TalentPresets[Phase.Phase2][0];
export const DefaultTalentsCombat 	= TalentPresets[Phase.Phase2][1];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase2][0];

export const DefaultTalents = DefaultTalentsAssassin;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RogueOptions.create({});


///////////////////////////////////////////////////////////////////////////
//                         Consumes/Buffs/Debuffs
///////////////////////////////////////////////////////////////////////////


export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfAgility,
	dragonBreathChili: false,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.ShadowOil,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectRegular,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	faerieFire: true,
	sunderArmor: true,
	mangle: true,
});

export const OtherDefaults = {
  	profession1: Profession.Engineering,
  	profession2: Profession.Leatherworking,
}