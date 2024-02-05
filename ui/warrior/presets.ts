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

import BlankGear from './gear_sets/blank.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase1DWGear from './gear_sets/phase_1_dw.gear.json';

import DefaultAPL from './apls/default.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);

export const GearArmsPhase1 = PresetUtils.makePresetGear('P1 Arms 2H', Phase1Gear, { talentTree: 0 });
export const GearArmsDWPhase1 = PresetUtils.makePresetGear('P1 Arms DW', Phase1DWGear, { talentTree: 0 });
export const GearFuryPhase1 = PresetUtils.makePresetGear('P1 Fury', Phase1Gear, { talentTree: 1 });

export const GearArmsDefault = GearArmsPhase1;
export const GearFuryDefault = GearFuryPhase1;

export const RotationArmsDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);
export const RotationFuryDefault = PresetUtils.makePresetAPLRotation('Default', DefaultAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const Talent25 = {
	name: 'Level 25',
	data: SavedTalents.create({
		talentsString: '303220203-01',
	}),
};

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
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectRegular,
	leaderOfThePack: true,
	devotionAura: TristateEffect.TristateEffectImproved,
	stoneskinTotem: TristateEffect.TristateEffectImproved,
})

export const DefaultDebuffs = Debuffs.create({
	curseOfWeakness: TristateEffect.TristateEffectRegular,
	faerieFire: true,
	mangle: true,
	sunderArmor: true,
})
