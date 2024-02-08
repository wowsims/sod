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

import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase1DWGear from './gear_sets/phase_1_dw.gear.json';
import Phase2DWGear from './gear_sets/phase_2_dw.gear.json';

export const GearArmsPhase1 = PresetUtils.makePresetGear('P1 Arms 2H', Phase1Gear, { talentTree: 0 });
export const GearArmsDWPhase1 = PresetUtils.makePresetGear('P1 Arms DW', Phase1DWGear, { talentTree: 0 });
export const GearFuryPhase1 = PresetUtils.makePresetGear('P1 Fury', Phase1Gear, { talentTree: 1 });
export const GearFuryPhase2 = PresetUtils.makePresetGear('P2 Fury', Phase2DWGear, { talentTree: 1 });

export const GearPresets = {
  [Phase.Phase1]: [
		GearArmsPhase1,
		GearFuryPhase1,
		GearArmsDWPhase1,
  ],
  [Phase.Phase2]: [
		GearFuryPhase2
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase2][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import Phase1APL from './apls/phase_1.apl.json';
import Phase2APL from './apls/phase_2.apl.json'

export const APLPhase1 = PresetUtils.makePresetAPLRotation('P1 Preset', Phase1APL);
export const APLPhase2 = PresetUtils.makePresetAPLRotation('P2 Preset', Phase2APL);


export const APLPresets = {
  [Phase.Phase1]: [
    APLPhase1,
  ],
  [Phase.Phase2]: [
	APLPhase2
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
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][1],
		2: APLPresets[Phase.Phase2][2],
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

export const TalentsPhase2Fury = {
	name: 'Level 40',
	data: SavedTalents.create({
		talentsString: '-05050005405010051',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
    TalentsPhase1,
  ],
  [Phase.Phase2]: [
	TalentsPhase2Fury,
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsArms = TalentPresets[Phase.Phase1][0];
export const DefaultTalentsFury = TalentPresets[Phase.Phase2][0];

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
	offHandImbue: WeaponImbue.DenseSharpeningStone,
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
