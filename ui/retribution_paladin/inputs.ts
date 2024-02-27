import { Player } from '../core/player.js';
import { ItemSlot, Spec } from '../core/proto/common.js';
import { TypedEvent } from '../core/typed_event.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	PaladinSeal,
	PaladinRune,
} from '../core/proto/paladin.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

// export const AuraSelection = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRetributionPaladin, PaladinAura>({
// 	fieldName: 'aura',
// 	values: [
// 		{ value: PaladinAura.NoPaladinAura, tooltip: 'No Aura' },
// 		{ actionId: () => ActionId.fromSpellId(10299), value: PaladinAura.RetributionAura },
// 	],
// });

// The below is used in the custom APL action "Cast Primary Seal".
// Only shows SoC if it's talented, only shows SoM if the relevant rune is equipped.
export const PrimarySealSelection = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRetributionPaladin, PaladinSeal>({
	fieldName: 'primarySeal',
	values: [
		{
			actionId: (player) => player.getMatchingSpellActionId([
				{id: 20154,	maxLevel: 9},
				{id: 20287, minLevel: 10, maxLevel: 17},
				{id: 20288, minLevel: 18, maxLevel: 25},
				{id: 20289, minLevel: 26, maxLevel: 33},
				{id: 20290, minLevel: 34, maxLevel: 41},
				{id: 20291, minLevel: 42, maxLevel: 49},
				{id: 20292, minLevel: 50, maxLevel: 57},
				{id: 20293, minLevel: 58},
			]),
			value: PaladinSeal.Righteousness 
		},
		{
			actionId: (player) => player.getMatchingSpellActionId([
				{id: 20375, maxLevel: 29},
				{id: 20915, minLevel: 30, maxLevel: 39},
				{id: 20918, minLevel: 40, maxLevel: 49},
				{id: 20919, minLevel: 50, maxLevel: 59},
				{id: 20920, minLevel: 60},
			]),
			value: PaladinSeal.Command,
			showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getTalents().sealOfCommand,
		},
		{
			actionId: () => ActionId.fromSpellId(407798),
			value: PaladinSeal.Martyrdom,
			showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getEquippedItem(ItemSlot.ItemSlotChest)?.rune?.id == PaladinRune.RuneChestSealofMartyrdom,
		},
	],
	// changeEmitter: (player: Player<Spec.SpecRetributionPaladin>) => player.changeEmitter,
	changeEmitter: (player: Player<Spec.SpecRetributionPaladin>) => TypedEvent.onAny([
		player.gearChangeEmitter, player.talentsChangeEmitter, player.specOptionsChangeEmitter, player.levelChangeEmitter
	]),
});

