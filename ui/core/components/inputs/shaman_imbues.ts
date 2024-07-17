import { Class, ItemSlot, WeaponImbue, WeaponType } from '../../proto/common.js';
import { isWeapon } from '../../proto_utils/utils';
import { ConsumableInputConfig } from './consumables';

// Shaman Imbues
export const RockbiterWeaponImbue = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 8017, minLevel: 1, maxLevel: 7 },
				{ id: 8018, minLevel: 8, maxLevel: 15 },
				{ id: 8019, minLevel: 16, maxLevel: 23 },
				{ id: 10399, minLevel: 24, maxLevel: 33 },
				{ id: 16314, minLevel: 34, maxLevel: 43 },
				{ id: 16315, minLevel: 44, maxLevel: 53 },
				{ id: 16316, minLevel: 54 },
			]),
		value: WeaponImbue.RockbiterWeapon,
		showWhen: player => player.isClass(Class.ClassShaman) && isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

export const FlametongueWeaponImbue = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 8024, minLevel: 10, maxLevel: 17 },
				{ id: 8027, minLevel: 18, maxLevel: 25 },
				{ id: 8030, minLevel: 26, maxLevel: 35 },
				{ id: 16339, minLevel: 36, maxLevel: 45 },
				{ id: 16341, minLevel: 46, maxLevel: 55 },
				{ id: 16342, minLevel: 56 },
			]),
		value: WeaponImbue.FlametongueWeapon,
		showWhen: player => player.isClass(Class.ClassShaman) && isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

export const FrostbrandWeaponImbue = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 8033, minLevel: 20, maxLevel: 27 },
				{ id: 8038, minLevel: 28, maxLevel: 37 },
				{ id: 10456, minLevel: 38, maxLevel: 47 },
				{ id: 16355, minLevel: 48, maxLevel: 57 },
				{ id: 16356, minLevel: 58 },
			]),
		value: WeaponImbue.FrostbrandWeapon,
		showWhen: player => player.isClass(Class.ClassShaman) && isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

export const WindfuryWeaponImbue = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 8232, minLevel: 30, maxLevel: 39 },
				{ id: 8235, minLevel: 40, maxLevel: 49 },
				{ id: 10486, minLevel: 50, maxLevel: 59 },
				{ id: 16362, minLevel: 60 },
			]),
		value: WeaponImbue.WindfuryWeapon,
		showWhen: player => player.isClass(Class.ClassShaman) && isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
