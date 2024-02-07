import {
  Consumes,
  Debuffs,
  Flask,
  Food,
  Profession,
  RaidBuffs,
  TristateEffect,
  WeaponImbue,
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
import Phase2Gear from './gear_sets/phase_2.gear.json';

import Phase1APL from './apls/phase_1.apl.json';
import Phase1AGAPL from './apls/phase_1_ag.apl.json';
import Phase2APL from './apls/phase_2.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);

export const DefaultGear = GearPhase1;

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL);
export const APLPhase1AG = PresetUtils.makePresetAPLRotation('Phase 1 (AG)', Phase1AGAPL);
export const APLPhase2 = PresetUtils.makePresetAPLRotation('Phase 2', Phase2APL);

export const DefaultAPL = APLPhase1;

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const TalentsPhase1 = {
  name: 'Phase 1',
  data: SavedTalents.create({
    talentsString: '25003105',
  }),
};

export const TalentsPhase2 = {
  name: 'Phase 2',
  data: SavedTalents.create({
    talentsString: '550031550000151',
  }),
};

export const DefaultTalents = TalentsPhase1;

export const DefaultOptions = ElementalShamanOptions.create({
  shield: ShamanShield.LightningShield,
  totems: ShamanTotems.create({
    earth: EarthTotem.StrengthOfEarthTotem,
    fire: FireTotem.SearingTotem,
    water: WaterTotem.HealingStreamTotem,
    air: AirTotem.WindfuryTotem,
  }),
});

export const DefaultConsumes = Consumes.create({
  flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
  mainHandImbue: WeaponImbue.RockbiterWeapon,
  offHandImbue: WeaponImbue.RockbiterWeapon,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	leaderOfThePack: true,
	moonkinAura: true,
	divineSpirit: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
  improvedScorch: true,
});

export const OtherDefaults = {
  distanceFromTarget: 20,
  profession1: Profession.Engineering,
  profession2: Profession.Tailoring,
  nibelungAverageCasts: 11,
}
