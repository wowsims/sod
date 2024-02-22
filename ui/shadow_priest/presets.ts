import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import {
	Consumes,
	Debuffs,
	EnchantedSigil,
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
import { SavedTalents } from '../core/proto/ui.js';

import {
	ShadowPriest_Options as Options,
} from '../core/proto/priest.js';

import * as PresetUtils from '../core/preset_utils.js';

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
  ],
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[CURRENT_PHASE][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import Phase1APL from './apls/phase_1.apl.json'
import Phase2APL from './apls/phase_2.apl.json'

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL);
export const APLPhase2 = PresetUtils.makePresetAPLRotation('Phase 2', Phase2APL);

export const APLPresets = {
  [Phase.Phase1]: [
		APLPhase1,
  ],
  [Phase.Phase2]: [
		APLPhase2,
  ],
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
  25: APLPresets[Phase.Phase1][0],
  40: APLPresets[CURRENT_PHASE][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsPhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '-20535000001',
	}),
};

export const TalentsPhase2 = {
	name: 'Phase 2',
	data: SavedTalents.create({
		talentsString: '--5022204002501251',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
		TalentsPhase1,
  ],
  [Phase.Phase2]: [
		TalentsPhase2,
  ],
};

export const DefaultTalents = TalentPresets[CURRENT_PHASE][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = Options.create({});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.LesserWizardOil,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	saygesFortune: SaygesFortune.SaygesDamage,
	sparkOfInspiration: true,
});

export const DefaultDebuffs = Debuffs.create({
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	channelClipDelay: 100,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
};
