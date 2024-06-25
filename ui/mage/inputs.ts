import * as InputHelpers from '../core/components/input_helpers.js';
import { ItemSlot, Spec } from '../core/proto/common.js';
import { Mage_Options_ArmorType as ArmorType, MageRune } from '../core/proto/mage.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Armor = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecMage, ArmorType>({
	fieldName: 'armor',
	values: [
		{ value: ArmorType.NoArmor, tooltip: 'No Armor' },
		{
			actionId: () => ActionId.fromSpellId(428741),
			value: ArmorType.MoltenArmor,
			showWhen: player => player.hasRune(ItemSlot.ItemSlotWrist, MageRune.RuneBracersMoltenArmor),
		},
		{
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 6117, minLevel: 34, maxLevel: 45 },
					{ id: 22782, minLevel: 46, maxLevel: 55 },
					{ id: 22783, minLevel: 58 },
				]),
			value: ArmorType.MageArmor,
		},
		{
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 168, minLevel: 1, maxLevel: 9 },
					{ id: 7300, minLevel: 10, maxLevel: 19 },
					{ id: 7301, minLevel: 20, maxLevel: 29 },
					{ id: 7302, minLevel: 30, maxLevel: 39 },
					{ id: 7320, minLevel: 40, maxLevel: 49 },
					{ id: 10219, minLevel: 50, maxLevel: 59 },
					{ id: 10220, minLevel: 60 },
				]),
			value: ArmorType.IceArmor,
		},
	],
	changeEmitter: player => TypedEvent.onAny([player.gearChangeEmitter, player.specOptionsChangeEmitter]),
});

export const MageRotationConfig = {
	inputs: [],
};
