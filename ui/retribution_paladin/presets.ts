import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
	Potions,
	AgilityElixir,
	StrengthBuff,
	WeaponImbue,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	PaladinAura,
	PaladinSeal,
	RetributionPaladin_Options as RetributionPaladinOptions,
} from '../core/proto/paladin.js';

import * as PresetUtils from '../core/preset_utils.js';

import Phase1Gear from './gear_sets/p1gear.json';

import DefaultApl from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

// export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearPhase1 = PresetUtils.makePresetGear('2h Ret', Phase1Gear);

export const GearPresets = {
  [Phase.Phase1]: [
    GearPhase1,
  ],
  [Phase.Phase2]: [
	GearPhase1,
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

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

///////////////////////////////////////////////////////////////////////////
//                                 Talent presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const P1RetTalents = {
	name: 'P1 Ret',
	data: SavedTalents.create({
		talentsString: '--05230051',
	})
};

export const P2DeepRetTalents = {
	name: 'P2 Deep Ret',
	data: SavedTalents.create({
		talentsString: '--532300512003151',
	}),
};

export const P2ShockadinTalents = {
	name: 'P2 Shockadin',
	data: SavedTalents.create({
		talentsString: '55050100521151--',
	}),
}

export const TalentPresets = {
  [Phase.Phase1]: [
    P1RetTalents,
  ],
  [Phase.Phase2]: [
	P2DeepRetTalents,
	P2ShockadinTalents
  ]
};

export const DefaultTalents = TalentPresets[Phase.Phase2][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	seal: PaladinSeal.Righteousness,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodSmokedSagefish,
	defaultPotion: Potions.GreaterManaPotion,
	mainHandImbue: WeaponImbue.WildStrikes,
	agilityElixir: AgilityElixir.ElixirOfLesserAgility,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	boglingRoot: false,
});
