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
			actionId: () => ActionId.fromSpellId(20154),
			value: PaladinSeal.Righteousness 
		},
		{
			actionId: () => ActionId.fromSpellId(20375),
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
	changeEmitter: (player: Player<Spec.SpecRetributionPaladin>) => TypedEvent.onAny([player.gearChangeEmitter, player.talentsChangeEmitter]),

});

