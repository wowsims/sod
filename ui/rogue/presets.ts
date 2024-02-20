import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	RogueOptions_PoisonImbue as Poison,
	RogueOptions,
} from '../core/proto/rogue.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';
import P1Daggers from './gear_sets/p1_daggers.gear.json';

import MutilateApl from './apls/mutilate.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearDaggersP1 = PresetUtils.makePresetGear('P1_Daggers', P1Daggers)

export const GearPresets = {
  [Phase.Phase1]: [
    GearDaggersP1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { talentTree: 0 });

export const APLPresets = {
  [Phase.Phase1]: [
    ROTATION_PRESET_MUTILATE,
  ],
  [Phase.Phase2]: [
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
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
		2: APLPresets[Phase.Phase1][0],
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const CombatDagger25Talents = {
	name: 'Combat Dagger',
	data: SavedTalents.create({
		talentsString: '-023305002001',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
		CombatDagger25Talents,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsAssassin = TalentPresets[Phase.Phase1][0];
export const DefaultTalentsCombat 	= TalentPresets[Phase.Phase1][0];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase1][0];

export const DefaultTalents = DefaultTalentsAssassin;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RogueOptions.create({
	mhImbue: Poison.NoPoison,
	ohImbue: Poison.NoPoison,
	applyPoisonsManually: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
