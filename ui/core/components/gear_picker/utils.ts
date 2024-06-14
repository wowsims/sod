import { ItemSlot } from '../../proto/common';

const emptySlotIcons: Record<ItemSlot, string> = {
	[ItemSlot.ItemSlotHead]: '/sod/assets/item_slots/head.jpg',
	[ItemSlot.ItemSlotNeck]: '/sod/assets/item_slots/neck.jpg',
	[ItemSlot.ItemSlotShoulder]: '/sod/assets/item_slots/shoulders.jpg',
	[ItemSlot.ItemSlotBack]: '/sod/assets/item_slots/shirt.jpg',
	[ItemSlot.ItemSlotChest]: '/sod/assets/item_slots/chest.jpg',
	[ItemSlot.ItemSlotWrist]: '/sod/assets/item_slots/wrists.jpg',
	[ItemSlot.ItemSlotHands]: '/sod/assets/item_slots/hands.jpg',
	[ItemSlot.ItemSlotWaist]: '/sod/assets/item_slots/waist.jpg',
	[ItemSlot.ItemSlotLegs]: '/sod/assets/item_slots/legs.jpg',
	[ItemSlot.ItemSlotFeet]: '/sod/assets/item_slots/feet.jpg',
	[ItemSlot.ItemSlotFinger1]: '/sod/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotFinger2]: '/sod/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotTrinket1]: '/sod/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotTrinket2]: '/sod/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotMainHand]: '/sod/assets/item_slots/mainhand.jpg',
	[ItemSlot.ItemSlotOffHand]: '/sod/assets/item_slots/offhand.jpg',
	[ItemSlot.ItemSlotRanged]: '/sod/assets/item_slots/ranged.jpg',
};
export function getEmptySlotIconUrl(slot: ItemSlot): string {
	return emptySlotIcons[slot];
}
