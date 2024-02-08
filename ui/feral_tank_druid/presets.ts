import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
	UnitReference
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	FeralTankDruid_Options as DruidOptions,
	FeralTankDruid_Rotation as DruidRotation,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);

export const GearPresets = {
  [Phase.Phase1]: [
    GearBlank,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const DefaultRotation = DruidRotation.create({
	maulRageThreshold: 25,
	maintainDemoralizingRoar: true,
	lacerateTime: 8.0,
});

export const DefaultAPL = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

export const APLPresets = {
  [Phase.Phase1]: [
    DefaultAPL,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
  25: APLPresets[Phase.Phase1][0],
  40: APLPresets[Phase.Phase1][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-503232132322010353120300313511-20350001',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
    StandardTalents,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalents = TalentPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = DruidOptions.create({
	innervateTarget: UnitReference.create(),
	startingRage: 20,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
