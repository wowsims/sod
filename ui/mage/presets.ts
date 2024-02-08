import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
	Profession
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Mage_Options_ArmorType as ArmorType,
	Mage_Options as MageOptions
} from '../core/proto/mage.js';

import * as PresetUtils from '../core/preset_utils.js';

import DefaultBlankGear from './gear_sets/blank.gear.json';

import APLDefault from './apls/default.apl.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArcanePhase1 = PresetUtils.makePresetGear('Default', DefaultBlankGear, { talentTree: 0 });
export const GearFirePhase1 = PresetUtils.makePresetGear('Default', DefaultBlankGear, { talentTree: 1 });
export const GearFrostPhase1 = PresetUtils.makePresetGear('Default', DefaultBlankGear, { talentTree: 2 });

export const GearPresets = {
  [Phase.Phase1]: [
    GearArcanePhase1,
		GearFirePhase1,
		GearFrostPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLArcanePhase1 = PresetUtils.makePresetAPLRotation('Default', APLDefault, { talentTree: 0 });
export const APLFirePhase1 = PresetUtils.makePresetAPLRotation('Default', APLDefault, { talentTree: 1 });
export const APLFrostPhase1 = PresetUtils.makePresetAPLRotation('Default', APLDefault, { talentTree: 2 });

export const APLPresets = {
  [Phase.Phase1]: [
    APLArcanePhase1,
		APLFirePhase1,
		APLFrostPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
		2: APLPresets[Phase.Phase1][2],
	},
  40: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
		2: APLPresets[Phase.Phase1][2],
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsArcanePhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '50005003021',
	}),
};

export const TalentsFirePhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '50005003021',
	}),
};

export const TalentsFrostPhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '50005003021',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
    TalentsArcanePhase1,
		TalentsFirePhase1,
		TalentsFrostPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsArcane = TalentPresets[Phase.Phase1][0];
export const DefaultTalentsFire 	= TalentPresets[Phase.Phase1][1];
export const DefaultTalentsFrost 	= TalentPresets[Phase.Phase1][2];

export const DefaultTalents = DefaultTalentsArcane;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
