import { Class, WeaponImbue } from '../../proto/common.js';
import { ConsumableInputConfig } from './consumables';

// Rogue Imbues
export const InstantPoisonWeaponImbue: ConsumableInputConfig<WeaponImbue> = {
	actionId: player =>
		player.getMatchingItemActionId([
			{ id: 6947, minLevel: 20, maxLevel: 27 },
			{ id: 6949, minLevel: 28, maxLevel: 35 },
			{ id: 6950, minLevel: 36, maxLevel: 43 },
			{ id: 8926, minLevel: 44, maxLevel: 51 },
			{ id: 8927, minLevel: 52, maxLevel: 59 },
			{ id: 8928, minLevel: 60, maxLevel: 60 },
		]),
	value: WeaponImbue.InstantPoison,
	showWhen: player => player.getClass() == Class.ClassRogue,
};

export const DeadlyPoisonWeaponImbue: ConsumableInputConfig<WeaponImbue> = {
	actionId: player =>
		player.getMatchingItemActionId([
			{ id: 2892, minLevel: 30, maxLevel: 37 },
			{ id: 2893, minLevel: 38, maxLevel: 45 },
			{ id: 8984, minLevel: 46, maxLevel: 53 },
			{ id: 8985, minLevel: 54, maxLevel: 59 },
			{ id: 20844, minLevel: 60, maxLevel: 60 },
		]),
	value: WeaponImbue.DeadlyPoison,
	showWhen: player => player.getClass() == Class.ClassRogue,
};

export const WoundPoisonWeaponImbue: ConsumableInputConfig<WeaponImbue> = {
	actionId: player =>
		player.getMatchingItemActionId([
			{ id: 10918, minLevel: 32, maxLevel: 39 },
			{ id: 10920, minLevel: 40, maxLevel: 47 },
			{ id: 10921, minLevel: 48, maxLevel: 55 },
			{ id: 10922, minLevel: 56, maxLevel: 60 },
		]),
	value: WeaponImbue.WoundPoison,
	showWhen: player => player.getClass() == Class.ClassRogue,
};

export const OccultPoisonWeaponImbue: ConsumableInputConfig<WeaponImbue> = {
	actionId: player => player.getMatchingItemActionId([{ id: 226374, minLevel: 56 }]),
	value: WeaponImbue.OccultPoison,
	showWhen: player => player.getClass() == Class.ClassRogue,
};

export const SebaciousPoisonWeaponImbue: ConsumableInputConfig<WeaponImbue> = {
	actionId: player => player.getMatchingItemActionId([{ id: 217345, minLevel: 60 }]),
	value: WeaponImbue.SebaciousPoison,
	showWhen: player => player.getClass() == Class.ClassRogue,
};
