import {
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions_WeaponImbue as WeaponImbue
} from '../core/proto/warlock.js';

import { Player } from '../core/player.js';
import { ItemSlot, Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import * as InputHelpers from '../core/components/input_helpers.js';
import { TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecTankWarlock, Armor>({
	fieldName: 'armor',
	values: [
		{ value: Armor.NoArmor, tooltip: 'No Armor' },
		{ actionId: (player) => player.getMatchingSpellActionId([
			{ id: 706, minLevel: 20, maxLevel: 39 },
			{ id: 11733, minLevel: 40, maxLevel: 49 },
			{ id: 11734, minLevel: 50, maxLevel: 59 },
			{ id: 11735, minLevel: 60 },
		]), value: Armor.DemonArmor },
	],
});

export const WeaponImbueInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecTankWarlock, WeaponImbue>({
	fieldName: 'weaponImbue',
	values: [
		{ value: WeaponImbue.NoWeaponImbue, tooltip: 'No Weapon Stone' },
		{ actionId: (player) => player.getMatchingItemActionId([
			{ id: 1254, minLevel: 28, maxLevel: 35 },
			{ id: 13699, minLevel: 36, maxLevel: 45 },
			{ id: 13700, minLevel: 46, maxLevel: 55 },
			{ id: 13701, minLevel: 56 },
		]), value: WeaponImbue.Firestone },
		{ actionId: (player) => player.getMatchingItemActionId([
			{ id: 5522, minLevel: 36, maxLevel: 47 },
			{ id: 13602, minLevel: 48, maxLevel: 59 },
			{ id: 13603, minLevel: 60 },
		]), value: WeaponImbue.Spellstone },
	],
	showWhen: (player) => player.getEquippedItem(ItemSlot.ItemSlotOffHand) == null && player.getLevel() >= 28,
	changeEmitter: (player: Player<Spec.SpecTankWarlock>) => player.changeEmitter,
});

export const PetInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecTankWarlock, Summon>({
	fieldName: 'summon',
	values: [
		{ value: Summon.NoSummon, tooltip: 'No Pet' },
		{ actionId: () => ActionId.fromSpellId(688), value: Summon.Imp },
		{ actionId: () => ActionId.fromSpellId(697), value: Summon.Voidwalker },
		{ actionId: () => ActionId.fromSpellId(712), value: Summon.Succubus },
		{ actionId: () => ActionId.fromSpellId(691), value: Summon.Felhunter },
	],
	changeEmitter: (player: Player<Spec.SpecTankWarlock>) => player.changeEmitter,
});

export const WarlockRotationConfig = {
	inputs: [],
};
