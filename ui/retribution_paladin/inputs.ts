import { Player } from '../core/player.js';
import { ItemSlot, Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import {
	PaladinAura,
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

export const StartingSealSelection = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRetributionPaladin, PaladinSeal>({
	fieldName: 'seal',
	values: [
		{ actionId: () => ActionId.fromSpellId(20290), value: PaladinSeal.Righteousness },
		{ actionId: () => ActionId.fromSpellId(407798), value: PaladinSeal.Martyrdom},
		{
			actionId: () => ActionId.fromSpellId(20375), value: PaladinSeal.Command,
			showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getTalents().sealOfCommand,
		},
		// the comparison below needs fixed
		// {
		// 	actionId: () => ActionId.fromSpellId(407799), value: PaladinSeal.Martyrdom,
		// 	showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getEquippedItem(ItemSlot.ItemSlotChest)?.rune == PaladinRune.RuneChestSealofMartyrdom,
		// },
	],
	changeEmitter: (player: Player<Spec.SpecRetributionPaladin>) => player.changeEmitter,
});

