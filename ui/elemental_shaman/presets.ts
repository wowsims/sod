import {
  Consumes,
  Flask,
  Food,
  Profession,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
  AirTotem,
  EarthTotem,
  ElementalShaman_Options as ElementalShamanOptions,
  FireTotem,
  ShamanShield,
  ShamanTotems,
  WaterTotem,
} from '../core/proto/shaman.js';

import * as PresetUtils from '../core/preset_utils.js';

import BlankGear from './gear_sets/blank.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';

import Phase1APL from './apls/phase_1.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BlankPresetGear = PresetUtils.makePresetGear('Blank', BlankGear);
export const Phase1PresetGear = PresetUtils.makePresetGear('Phase 1', Phase1Gear);

export const DefaultGear = Phase1PresetGear;

export const Phase1PresetAPL = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL);

export const DefaultAPL = Phase1PresetAPL

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
  name: 'Phase 1',
  data: SavedTalents.create({
    talentsString: '25003105',
  }),
};

export const DefaultOptions = ElementalShamanOptions.create({
  shield: ShamanShield.WaterShield,
  totems: ShamanTotems.create({
    earth: EarthTotem.StrengthOfEarthTotem,
    air: AirTotem.WindfuryTotem,
    fire: FireTotem.TotemOfWrath,
    water: WaterTotem.ManaSpringTotem,
    useFireElemental: true,
  }),
});

export const OtherDefaults = {
    distanceFromTarget: 20,
    profession1: Profession.Engineering,
    profession2: Profession.Tailoring,
    nibelungAverageCasts: 11,
}

export const DefaultConsumes = Consumes.create({
  flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});