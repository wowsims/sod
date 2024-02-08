import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
	Spec,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	WarriorShout,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	ProtectionWarrior_Options as ProtectionWarriorOptions,
} from '../core/proto/warrior.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

import BlankGear from './gear_sets/blank.gear.json';

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

import DefaultApl from './apls/default.apl.json';

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

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Cooldowns', Spec.SpecProtectionWarrior, ProtectionWarriorRotation.create());

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '2500030023-302-053351225000012521030113321',
	}),
};

export const UATalents = {
	name: 'UA',
	data: SavedTalents.create({
		talentsString: '35023301230051002020120002-2-05035122500000252',
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

export const DefaultOptions = ProtectionWarriorOptions.create({
	shout: WarriorShout.WarriorShoutCommanding,
	startingRage: 0,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});
