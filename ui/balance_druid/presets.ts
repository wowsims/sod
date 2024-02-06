import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	IndividualBuffs,
	PartyBuffs,
	Profession,
	RaidBuffs,
	TristateEffect,
	UnitReference,
	WeaponImbue
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	BalanceDruid_Options as BalanceDruidOptions,
} from '../core/proto/druid.js';

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

export const APLBalancePhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL);

export const DEFAULT_APL = APLBalancePhase1

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.
export const Phase1PresetTalents = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '50005003021',
	}),
};

export const DefaultTalents = Phase1PresetTalents;

export const DefaultOptions = BalanceDruidOptions.create({
	innervateTarget: UnitReference.create(),
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodUnknown,
	mainHandImbue: WeaponImbue.BlackfathomManaOil,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
});

export const DefaultPartyBuffs = PartyBuffs.create({
});

export const DefaultDebuffs = Debuffs.create({
	faerieFire: true,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
