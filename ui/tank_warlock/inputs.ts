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
		{ actionId: () => ActionId.fromSpellId(11735), value: Armor.DemonArmor },
	],
});

export const WeaponImbueInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecTankWarlock, WeaponImbue>({
	fieldName: 'weaponImbue',
	values: [
		{ value: WeaponImbue.NoWeaponImbue, tooltip: 'No Weapon Stone' },
		// TODO: Classic warlock weapon stone id based on level
		{ actionId: () => ActionId.fromItemId(13701), value: WeaponImbue.Firestone },
		{ actionId: () => ActionId.fromItemId(13603), value: WeaponImbue.Spellstone },
	],
	showWhen: (player) => player.getEquippedItem(ItemSlot.ItemSlotOffHand) == null,
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
