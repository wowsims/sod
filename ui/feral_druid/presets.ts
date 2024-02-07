import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
	Profession,
	Spec,
	Potions,
	Conjured,
	WeaponImbue,
	AgilityElixir,
	StrengthBuff
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';
import Phase1Gear from './gear_sets/p1.gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);

export const GearPresets = {
  [Phase.Phase1]: [
    GearPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1 = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);

export const APLPresets = {
  [Phase.Phase1]: [
    APLPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
  25: APLPresets[Phase.Phase1][0],
  40: APLPresets[Phase.Phase1][0],
};

export const DefaultRotation = FeralDruidRotation.create({
	maintainFaerieFire: false,
	minCombosForRip: 3,
	maxWaitTime: 2.0,
	preroarDuration: 26.0,
	precastTigersFury: false,
	useShredTrick: false,
});

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFeralDruid, DefaultRotation);

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsPhase1 = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '500005001--05',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
    TalentsPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalents = TalentPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = FeralDruidOptions.create({
	latencyMs: 100,
	assumeBleedActive: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodSmokedSagefish,
	defaultPotion: Potions.ManaPotion,
	defaultConjured: Conjured.ConjuredMinorRecombobulator,
	mainHandImbue: WeaponImbue.BlackfathomSharpeningStone,
	agilityElixir: AgilityElixir.ElixirOfLesserAgility,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	boglingRoot: true,
});

export const OtherDefaults = {
	profession2: Profession.Leatherworking,
};
