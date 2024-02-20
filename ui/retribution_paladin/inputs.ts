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

export const PrimarySealSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecRetributionPaladin>({
	fieldName: 'primarySeal',
	label: 'Primary Seal',
	labelTooltip: 'The primary Seal to be used in the paladin\'s rotation.',
	values: [
		{ name: "Seal of Righteousness (Rank 1)", value: PaladinSeal.RighteousnessRank1},
		{ name: "Seal of Righteousness (Rank 2)", value: PaladinSeal.RighteousnessRank2},
		// { name: "Seal of Righteousness (Rank 3)", value: PaladinSeal.RighteousnessRank3},
		// { 
		// 	name: "Seal of Righteousness (Rank 4)",
		// 	value: PaladinSeal.RighteousnessRank4,
		// 	// showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getLevel() >= 25,
		// },
		// { name: "Seal of Righteousness (Rank 5)", value: PaladinSeal.RighteousnessRank5},


		// { actionId: () => ActionId.fromSpellId(20290), value: PaladinSeal.Righteousness },
		// { actionId: () => ActionId.fromSpellId(407798), value: PaladinSeal.Martyrdom},
		// {
		// 	actionId: () => ActionId.fromSpellId(20375), value: PaladinSeal.Command,
		// 	showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getTalents().sealOfCommand,
		// },
		// the comparison below needs fixed
		// {
		// 	actionId: () => ActionId.fromSpellId(407799), value: PaladinSeal.Martyrdom,
		// 	showWhen: (player: Player<Spec.SpecRetributionPaladin>) => player.getEquippedItem(ItemSlot.ItemSlotChest)?.rune == PaladinRune.RuneChestSealofMartyrdom,
		// },
	],
	changeEmitter: (player: Player<Spec.SpecRetributionPaladin>) => player.changeEmitter,
});

