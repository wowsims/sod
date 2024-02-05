import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	RaidBuffs,
	TristateEffect
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	EnhancementShaman_Options as EnhancementShamanOptions,
	ShamanImbue,
	ShamanShield,
	ShamanSyncType,
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

export const DefaultGear = Phase1PresetGear

export const Phase1PresetAPL = PresetUtils.makePresetAPLRotation('P1 Preset', Phase1APL);

export const DefaultAPL = Phase1PresetAPL

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'P1 Preset',
	data: SavedTalents.create({
		talentsString: '-5005202101',
	}),
};

export const DefaultOptions = EnhancementShamanOptions.create({
	shield: ShamanShield.LightningShield,
	imbueMh: ShamanImbue.RockbiterWeapon,
	imbueOh: ShamanImbue.RockbiterWeapon,
	syncType: ShamanSyncType.Auto,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	leaderOfThePack: true,
	moonkinAura: true,
	divineSpirit: true,
	battleShout: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	sunderArmor: true,
	curseOfElements: true,
	curseOfRecklessness: true,
	faerieFire: true,
});

export const OtherDefaults = {
};
