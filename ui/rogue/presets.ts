import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	RogueOptions,
} from '../core/proto/rogue.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';
import P1Daggers from './gear_sets/p1_daggers.gear.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';

import MutilateApl from './apls/mutilate.apl.json';
import SinisterApl25 from './apls/basic_strike_25.apl.json';
import SinisterApl40 from './apls/basic_strike_40.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearDaggersP1 = PresetUtils.makePresetGear('P1 Daggers', P1Daggers)
export const GearCombatP1 = PresetUtils.makePresetGear("P1 Combat", P1CombatGear)

export const GearPresets = {
  [Phase.Phase1]: [
    GearDaggersP1,
	GearCombatP1,
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
export const ROTATION_PRESET_SINISTER_25 = PresetUtils.makePresetAPLRotation('Sinister', SinisterApl25, { talentTree: 1 });
export const ROTATION_PRESET_SINISTER_40 = PresetUtils.makePresetAPLRotation('Sinister', SinisterApl40, { talentTree: 1 });

export const APLPresets = {
  [Phase.Phase1]: [
    ROTATION_PRESET_MUTILATE,
	ROTATION_PRESET_SINISTER_25,
  ],
  [Phase.Phase2]: [
	ROTATION_PRESET_SINISTER_40,
	ROTATION_PRESET_MUTILATE
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
		2: APLPresets[Phase.Phase1][1],
	},
  40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][0],
		2: APLPresets[Phase.Phase2][0],
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
export const ColdBloodMutilate40Talents = {
	name: 'CB Mutilate',
	data: SavedTalents.create({
		talentsString: '005303103551--05'
	})
};

export const TalentPresets = {
	[Phase.Phase1]: [
		CombatDagger25Talents,
	],
	[Phase.Phase2]: [
		ColdBloodMutilate40Talents,
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

export const DefaultOptions = RogueOptions.create({});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
