import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
  Consumes,
  Debuffs,
  EnchantedSigil,
  FirePowerBuff,
  Food,
  IndividualBuffs,
  Potions,
  Profession,
  RaidBuffs,
  SaygesFortune,
  SpellPowerBuff,
  TristateEffect,
  WeaponImbue,
} from '../core/proto/common.js';
import {
  AirTotem,
  EarthTotem,
  ElementalShaman_Options as ElementalShamanOptions,
  FireTotem,
  ShamanShield,
  ShamanTotems,
  WaterTotem,
} from '../core/proto/shaman.js';
import { SavedTalents } from '../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';

export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);

export const GearPresets = {
  [Phase.Phase1]: [
    GearPhase1,
  ],
  [Phase.Phase2]: [
    GearPhase2,
  ]
};

export const DefaultGear = GearPresets[CURRENT_PHASE][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import Phase1APL from './apls/phase_1.apl.json';
import Phase1AGAPL from './apls/phase_1_ag.apl.json';
import Phase2APL from './apls/phase_2.apl.json';

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL);
export const APLPhase1AG = PresetUtils.makePresetAPLRotation('Phase 1 (AG)', Phase1AGAPL);
export const APLPhase2 = PresetUtils.makePresetAPLRotation('Phase 2', Phase2APL);

export const APLPresets = {
	[Phase.Phase1]: [
		APLPhase1,
		APLPhase1AG,
	],
	[Phase.Phase2]: [
		APLPhase2,
	]
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

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

export const TalentPresets = {
	[Phase.Phase1]: [
		TalentsPhase1,
	],
	[Phase.Phase2]: [
		TalentsPhase2,
	]
};

export const DefaultTalents = TalentPresets[CURRENT_PHASE][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

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
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.LesserWizardOil,
	offHandImbue: WeaponImbue.LesserWizardOil,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
  	sparkOfInspiration: true,
  	saygesFortune: SaygesFortune.SaygesDamage
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
	dreamstate: true,
	improvedScorch: true,
});

export const OtherDefaults = {
  	distanceFromTarget: 20,
  	profession1: Profession.Enchanting,
  	profession2: Profession.Leatherworking,
}
