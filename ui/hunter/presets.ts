import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
	Hunter_Options_QuiverBonus,
} from '../core/proto/hunter.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

import Phase1Gear from './gear_sets/phase1.json';
import Phase2Gear from './gear_sets/phase2.json';

export const GearBeastMasteryPhase1 = PresetUtils.makePresetGear('P1 Beast Mastery', Phase1Gear, { talentTree: 0 })
export const GearMarksmanPhase1 = PresetUtils.makePresetGear('P1 Marksmanship', Phase1Gear, { talentTree: 1 })
export const GearSurvivalPhase1 = PresetUtils.makePresetGear('P1 Survival', Phase1Gear, { talentTree: 2 })

export const GearPhase2 = PresetUtils.makePresetGear('P2 Gear', Phase2Gear)

export const GearPresets = {
  	[Phase.Phase1]: [
    	GearBeastMasteryPhase1,
		GearMarksmanPhase1,
		GearSurvivalPhase1,
  	],
  	[Phase.Phase2]: [
		GearPhase2
  	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPhase2;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import MeleeWeaveP1 from './apls/melee.weave.p1.json';
import MeleeP2 from './apls/melee.p2.json';

export const APLMeleeWeavePhase1 = PresetUtils.makePresetAPLRotation('Melee Weave P1', MeleeWeaveP1);
export const APLMeleePhase2 = PresetUtils.makePresetAPLRotation('Melee P2', MeleeP2);

export const APLPresets = {
  	[Phase.Phase1]: [
    	APLMeleeWeavePhase1,
  	],
  	[Phase.Phase2]: [
		APLMeleePhase2
  	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
		2: APLPresets[Phase.Phase1][0],
	},
  40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][0],
		2: APLPresets[Phase.Phase2][0],
	}
};

export const DefaultSimpleRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: StingType.SerpentSting,
	multiDotSerpentSting: true,
});

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsBeastMasteryPhase1 = {
	name: 'P1 Beast Mastery',
	data: SavedTalents.create({
		talentsString: '53000200501',
	}),
};

export const TalentsMarksmanPhase1 = {
	name: 'P1 Marksmanship',
	data: SavedTalents.create({
		talentsString: '-050515',
	}),
};

export const TalentsSurvivalPhase1 = {
	name: 'P1 Survival',
	data: SavedTalents.create({
		talentsString: '--33502001101',
	}),
};

export const TalentsSurvivalPhase2 = {
	name: 'P2 Survival',
	data: SavedTalents.create({
		talentsString: '--335020051030315',
	}),
};

export const TalentPresets = {
  	[Phase.Phase1]: [
    	TalentsBeastMasteryPhase1,
		TalentsMarksmanPhase1,
		TalentsSurvivalPhase1,
  	],
 	[Phase.Phase2]: [
		TalentsSurvivalPhase2
  	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsBeastMastery 	= TalentPresets[Phase.Phase1][0];
export const DefaultTalentsMarksman 		= TalentPresets[Phase.Phase1][1];
export const DefaultTalentsSurvival 		= TalentPresets[Phase.Phase1][2];

export const DefaultTalents = TalentsSurvivalPhase2;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.JaggedArrow,
	petType: PetType.WindSerpent,
	petTalents: {},
	petUptime: 1,
	quiverBonus: Hunter_Options_QuiverBonus.Speed14
});

export const BMDefaultOptions = HunterOptions.create({
	ammo: Ammo.JaggedArrow,
	petType: PetType.Cat,
	petTalents: {},
	petUptime: 1,
	quiverBonus: Hunter_Options_QuiverBonus.Speed14
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
