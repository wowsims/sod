import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Potions,
	RaidBuffs,
	TristateEffect,
	WeaponImbue,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	WarriorShout,
	Warrior_Options as WarriorOptions,
} from '../core/proto/warrior.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

import BlankGear from './gear_sets/blank.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase1DWGear from './gear_sets/phase_1_dw.gear.json';

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);

export const GearArmsPhase1 = PresetUtils.makePresetGear('P1 Arms 2H', Phase1Gear, { talentTree: 0 });
export const GearArmsDWPhase1 = PresetUtils.makePresetGear('P1 Arms DW', Phase1DWGear, { talentTree: 0 });
export const GearFuryPhase1 = PresetUtils.makePresetGear('P1 Fury', Phase1Gear, { talentTree: 1 });

export const GearPresets = {
  [Phase.Phase1]: [
    GearArmsPhase1,
		GearFuryPhase1,
		GearArmsDWPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import Phase1APL from './apls/phase_1.apl.json';

export const APLPhase1 = PresetUtils.makePresetAPLRotation('P1 Preset', Phase1APL);

export const APLPresets = {
  [Phase.Phase1]: [
    APLPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
		2: APLPresets[Phase.Phase1][0],
	},
  40: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
		2: APLPresets[Phase.Phase1][0],
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsPhase1 = {
	name: 'Level 25',
	data: SavedTalents.create({
		talentsString: '303220203-01',
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
export const DefaultTalentsArms = TalentPresets[Phase.Phase1][0];
export const DefaultTalentsFury = TalentPresets[Phase.Phase1][0];

export const DefaultTalents = DefaultTalentsArms;

///////////////////////////////////////////////////////////////////////////
//                                 Options Presets
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	shout: WarriorShout.WarriorShoutBattle,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodGrilledSquid,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.BlackfathomSharpeningStone,
	defaultPotion: Potions.UnknownPotion,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	devotionAura: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	stoneskinTotem: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectRegular,
})

export const DefaultDebuffs = Debuffs.create({
	curseOfWeakness: TristateEffect.TristateEffectRegular,
	faerieFire: true,
	mangle: true,
	sunderArmor: true,
})
