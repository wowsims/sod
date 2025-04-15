package core

import (
	"fmt"
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type OnItemSwap func(*Simulation, proto.ItemSlot, bool)

type ItemSwap struct {
	character       *Character
	onSwapCallbacks [NumItemSlots][]OnItemSwap

	isFuryWarrior        bool
	mhCritMultiplier     float64
	ohCritMultiplier     float64
	rangedCritMultiplier float64

	// Which slots to actually swap.
	slots []proto.ItemSlot

	// Holds the original equip
	originalEquip Equipment
	// Holds the items that are selected for swapping
	swapEquip Equipment
	// Holds items that are currently not equipped
	unEquippedItems Equipment
	swapSet         proto.APLActionItemSwap_SwapSet
	equipmentStats  ItemSwapStats

	initialized bool
}

type ItemSwapStats struct {
	allSlots    stats.Stats
	weaponSlots stats.Stats
}

func (character *Character) enableItemSwap(itemSwap *proto.ItemSwap) {
	var swapItems Equipment
	hasItemSwap := make(map[proto.ItemSlot]bool)

	for idx, itemSpec := range itemSwap.Items {
		itemSlot := proto.ItemSlot(idx)
		hasItemSwap[itemSlot] = itemSpec != nil && itemSpec.Id != 0
		if !slices.Contains(AllWeaponSlots(), itemSlot) && hasItemSwap[itemSlot] {
			panic(fmt.Sprintf("Slot %d is not supported. Currently only Mainhand, Offhand and Ranged are supported.", itemSlot))
		}
		swapItems[itemSlot] = toItem(itemSpec)
	}

	has2HSwap := swapItems[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand
	hasMhEquipped := character.HasMHWeapon()
	hasOhEquipped := character.HasOHWeapon()

	// Handle MH and OH together, because present MH + empty OH --> swap MH and unequip OH
	if hasItemSwap[proto.ItemSlot_ItemSlotOffHand] && hasMhEquipped {
		hasItemSwap[proto.ItemSlot_ItemSlotMainHand] = true
	}

	if has2HSwap && hasOhEquipped {
		hasItemSwap[proto.ItemSlot_ItemSlotOffHand] = true
	}

	slots := SetToSortedSlice(hasItemSwap)

	if len(slots) == 0 {
		return
	}

	var prepullBonusStats stats.Stats
	if itemSwap.PrepullBonusStats != nil {
		prepullBonusStats = stats.FromFloatArray(itemSwap.PrepullBonusStats.Stats)
	}

	equipmentStats := calcItemSwapStatsOffset(character.Equipment, swapItems, prepullBonusStats, slots)

	character.ItemSwap = ItemSwap{
		slots:           slots,
		originalEquip:   character.Equipment,
		swapEquip:       swapItems,
		unEquippedItems: swapItems,
		equipmentStats:  equipmentStats,
		swapSet:         proto.APLActionItemSwap_Main,
		initialized:     false,
	}
}

func (swap *ItemSwap) initialize(character *Character) {
	swap.character = character
}

func (character *Character) RegisterItemSwapCallback(slots []proto.ItemSlot, callback OnItemSwap) {
	if character == nil || !character.ItemSwap.IsEnabled() || len(slots) == 0 {
		return
	}

	if (character.Env != nil) && character.Env.IsFinalized() {
		panic("Tried to add a new item swap callback in a finalized environment!")
	}

	for _, slot := range slots {
		character.ItemSwap.onSwapCallbacks[slot] = append(character.ItemSwap.onSwapCallbacks[slot], callback)
	}
}

// Helper for handling Item Effects that use the itemID to toggle the aura on and off
// This will also get the eligible slots for the item
func (swap *ItemSwap) RegisterProc(itemID int32, aura *Aura) {
	slots := swap.EligibleSlotsForItem(itemID)
	swap.RegisterProcWithSlots(itemID, aura, slots)
}

// Helper for handling Item Effects that use the itemID to toggle the aura on and off
func (swap *ItemSwap) RegisterProcWithSlots(itemID int32, aura *Aura, slots []proto.ItemSlot) {
	swap.registerProcInternal(ItemSwapProcConfig{
		ItemID: itemID,
		Aura:   aura,
		Slots:  slots,
	})
}

// Helper for handling Enchant Effects that use the effectID to toggle the aura on and off
func (swap *ItemSwap) RegisterEnchantProc(effectID int32, aura *Aura) {
	slots := swap.EligibleSlotsForEffect(effectID)
	swap.RegisterEnchantProcWithSlots(effectID, aura, slots)
}
func (swap *ItemSwap) RegisterEnchantProcWithSlots(effectID int32, aura *Aura, slots []proto.ItemSlot) {
	swap.registerProcInternal(ItemSwapProcConfig{
		EnchantId: effectID,
		Aura:      aura,
		Slots:     slots,
	})
}

type ItemSwapProcConfig struct {
	ItemID    int32
	EnchantId int32
	Aura      *Aura
	Slots     []proto.ItemSlot
}

func (swap *ItemSwap) registerProcInternal(config ItemSwapProcConfig) {
	isItemProc := config.ItemID != 0
	isEnchantEffectProc := config.EnchantId != 0

	// Enchant effects such as Weapon/Back do not trigger an ICD
	shouldUpdateIcd := isItemProc && (config.Aura.Icd != nil)

	character := swap.character
	character.RegisterItemSwapCallback(config.Slots, func(sim *Simulation, _ proto.ItemSlot, _ bool) {
		isItemSlotMatch := false

		if isItemProc {
			isItemSlotMatch = character.hasItemEquipped(config.ItemID, config.Slots)
		} else if isEnchantEffectProc {
			isItemSlotMatch = character.hasEnchantEquipped(config.EnchantId, config.Slots)
		}

		if isItemSlotMatch {
			if !config.Aura.IsActive() {
				config.Aura.Activate(sim)
			}
			if swap.initialized && shouldUpdateIcd {
				config.Aura.Icd.Use(sim)
			}
		} else {
			config.Aura.Deactivate(sim)
			if swap.initialized && shouldUpdateIcd {
				// This is a hack to block ActivateAura APL
				// actions from executing for unequipped items.
				config.Aura.Icd.Set(NeverExpires)
			}
		}
	})
}

// Helper for handling Item On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) RegisterActive(itemID int32) {
	swap.registerActiveInternal(ItemSwapActiveConfig{
		ActionID: ActionID{ItemID: itemID},
		ItemID:   itemID,
		Slots:    swap.EligibleSlotsForItem(itemID),
	})
}

// Helper for handling Enchant On Use effects to set a 30s cd on the related spell.
// Currently only used for random suffix "enchants" like the Karazhan gloves.
func (swap *ItemSwap) RegisterEnchantActive(effectID int32, spellID int32) {
	swap.registerActiveInternal(ItemSwapActiveConfig{
		ActionID:  ActionID{SpellID: spellID},
		EnchantId: effectID,
		Slots:     swap.EligibleSlotsForEffect(effectID),
	})
}

type ItemSwapActiveConfig struct {
	ActionID  ActionID
	ItemID    int32
	EnchantId int32
	Slots     []proto.ItemSlot
}

func (swap *ItemSwap) registerActiveInternal(config ItemSwapActiveConfig) {
	isItemActive := config.ItemID != 0
	isEnchantEffectActive := config.EnchantId != 0

	character := swap.character
	character.RegisterItemSwapCallback(config.Slots, func(sim *Simulation, _ proto.ItemSlot, _ bool) {
		spell := character.GetSpell(config.ActionID)
		if spell == nil {
			return
		}

		aura := character.GetAuraByID(spell.ActionID)
		if aura.IsActive() {
			aura.Deactivate(sim)
		}

		var isEquipped bool
		if isItemActive {
			isEquipped = character.hasItemEquipped(config.ItemID, config.Slots)
		} else if isEnchantEffectActive {
			isEquipped = character.hasEnchantEquipped(config.EnchantId, config.Slots)
		}

		if !isEquipped {
			spell.Flags |= SpellFlagSwapped
			return
		}

		spell.Flags &= ^SpellFlagSwapped

		if !swap.initialized {
			return
		}

		spell.CD.Set(sim.CurrentTime + max(spell.CD.TimeToReady(sim), time.Second*30))
	})
}

// // Helper for handling Effects that use PPMManager to toggle the aura on/off
// func (swap *ItemSwap) RegisterOnSwapItemForEffectWithPPMManager(effectID int32, ppm float64, dpm *DynamicProcManager, aura *Aura) {
// 	slots := swap.EligibleSlotsForEffect(effectID)
// 	character := swap.character
// 	character.RegisterItemSwapCallback(slots, func(sim *Simulation, _ proto.ItemSlot) {
// 		procMask := character.GetDynamicProcMaskForWeaponEnchant(effectID)
// 		*dpm = character.AutoAttacks.NewPPMManager(ppm, procMask)

// 		if dpm.Chance(procMask) == 0 {
// 			aura.Deactivate(sim)
// 		} else {
// 			aura.Activate(sim)
// 		}
// 	})

// }

// // Helper for handling Effects that use the effectID to toggle the aura on and off
// func (swap *ItemSwap) RegisterOnSwapItemForEffect(effectID int32, aura *Aura) {
// 	slots := swap.EligibleSlotsForEffect(effectID)
// 	character := swap.character
// 	character.RegisterItemSwapCallback(slots, func(sim *Simulation, _ proto.ItemSlot) {
// 		procMask := character.GetDynamicProcMaskForWeaponEnchant(effectID)

// 		if procMask == ProcMaskUnknown {
// 			aura.Deactivate(sim)
// 		} else {
// 			aura.Activate(sim)
// 		}
// 	})
// }

func (swap *ItemSwap) IsEnabled() bool {
	return swap.character != nil && len(swap.slots) > 0
}

func (swap *ItemSwap) IsValidSwap(swapSet proto.APLActionItemSwap_SwapSet) bool {
	return swap.swapSet != swapSet
}

func (swap *ItemSwap) IsSwapped() bool {
	return swap.swapSet == proto.APLActionItemSwap_Swap1
}

func (character *Character) hasItemEquipped(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return character.Equipment.containsItemInSlots(itemID, possibleSlots)
}

func (character *Character) hasEnchantEquipped(effectID int32, possibleSlots []proto.ItemSlot) bool {
	return character.Equipment.containsEnchantInSlots(effectID, possibleSlots)
}

func (swap *ItemSwap) GetEquippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.character.Equipment[slot]
}

func (swap *ItemSwap) GetUnequippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.unEquippedItems[slot]
}

func (swap *ItemSwap) ItemExistsInMainEquip(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return swap.originalEquip.containsItemInSlots(itemID, possibleSlots)
}

func (swap *ItemSwap) ItemExistsInSwapEquip(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return swap.swapEquip.containsItemInSlots(itemID, possibleSlots)
}

func (swap *ItemSwap) EligibleSlotsForItem(itemID int32) []proto.ItemSlot {
	eligibleSlots := eligibleSlotsForItem(GetItemByID(itemID))

	if len(eligibleSlots) == 0 {
		return []proto.ItemSlot{}
	}

	if !swap.IsEnabled() {
		return eligibleSlots
	} else {
		return FilterSlice(eligibleSlots, func(slot proto.ItemSlot) bool {
			return (swap.originalEquip[slot].ID == itemID) || (swap.swapEquip[slot].ID == itemID)
		})
	}
}

func (swap *ItemSwap) EligibleSlotsForEffect(effectID int32) []proto.ItemSlot {
	if !swap.IsEnabled() {
		return swap.character.Equipment.EligibleSlotsForEffect(effectID)
	}

	var eligibleSlots []proto.ItemSlot
	for itemSlot := proto.ItemSlot(0); itemSlot < NumItemSlots; itemSlot++ {
		if swap.originalEquip.containsEnchantInSlot(effectID, itemSlot) || (swap.IsEnabled() && swap.swapEquip.containsEnchantInSlot(effectID, itemSlot)) {
			eligibleSlots = append(eligibleSlots, itemSlot)
		}
	}

	return eligibleSlots
}

func (swap *ItemSwap) SwapItems(sim *Simulation, swapSet proto.APLActionItemSwap_SwapSet, isReset bool) {
	if !swap.IsEnabled() || (!swap.IsValidSwap(swapSet) && !isReset) {
		return
	}

	character := swap.character
	weaponSlotSwapped := false
	isPrepull := sim.CurrentTime < 0

	for _, slot := range swap.slots {
		if (slot >= proto.ItemSlot_ItemSlotMainHand) && (slot <= proto.ItemSlot_ItemSlotRanged) {
			weaponSlotSwapped = true
		} else if !isReset && !isPrepull {
			continue
		}

		swap.swapItem(sim, slot, isPrepull, isReset)

		for _, onSwap := range swap.onSwapCallbacks[slot] {
			onSwap(sim, slot, isReset)
		}
	}

	if !swap.IsValidSwap(swapSet) {
		return
	}

	statsToSwap := Ternary(isPrepull, swap.equipmentStats.allSlots, swap.equipmentStats.weaponSlots)
	if swap.IsSwapped() {
		statsToSwap = statsToSwap.Invert()
	}

	if sim.Log != nil {
		sim.Log("Item Swap - Stats Change: %v", statsToSwap.FlatString())
	}
	character.AddDynamicEquipStats(sim, statsToSwap)

	if !isPrepull && !isReset && weaponSlotSwapped {
		character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		character.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime)
		character.SetGCDTimer(sim, max(character.NextGCDAt(), sim.CurrentTime+GCDDefault))
	}

	swap.swapSet = swapSet
}

func (swap *ItemSwap) swapItem(sim *Simulation, slot proto.ItemSlot, isPrepull bool, isReset bool) {
	oldItem := *swap.GetEquippedItemBySlot(slot)

	if isReset {
		swap.character.Equipment[slot] = swap.originalEquip[slot]
	} else {
		swap.character.Equipment[slot] = swap.unEquippedItems[slot]
	}

	swap.unEquippedItems[slot] = oldItem

	character := swap.character

	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:
		if character.AutoAttacks.AutoSwingMelee {
			character.AutoAttacks.SetMH(sim, character.WeaponFromMainHand())
		}
	case proto.ItemSlot_ItemSlotOffHand:
		// OH slot handling is more involved because we need to dynamically toggle the OH weapon attack on/off
		// depending on the updated DW status after the swap.
		if character.AutoAttacks.AutoSwingMelee {
			weapon := character.WeaponFromOffHand()
			isCurrentlyDualWielding := character.AutoAttacks.IsDualWielding
			character.AutoAttacks.SetOH(sim, weapon)
			if !isPrepull && !isCurrentlyDualWielding {
				character.AutoAttacks.IsDualWielding = weapon.SwingSpeed != 0
				character.AutoAttacks.EnableMeleeSwing(sim)
			}
			character.PseudoStats.CanBlock = character.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
		}
	case proto.ItemSlot_ItemSlotRanged:
		if character.AutoAttacks.AutoSwingRanged {
			character.AutoAttacks.SetRanged(sim, character.WeaponFromRanged())
		}
	}
}

func (swap *ItemSwap) reset(sim *Simulation) {
	swap.initialized = false
	if !swap.IsEnabled() {
		return
	}

	swap.SwapItems(sim, proto.APLActionItemSwap_Main, true)

	swap.unEquippedItems = swap.swapEquip

	// This is used to set the initial spell flags for unequipped items.
	// Reset is called before the first iteration.
	swap.initialized = true
}

func (swap *ItemSwap) doneIteration(sim *Simulation) {
	swap.reset(sim)
}

func calcItemSwapStatsOffset(originalEquipment Equipment, swapEquipment Equipment, prepullBonusStats stats.Stats, slots []proto.ItemSlot) ItemSwapStats {
	allSlotStats := prepullBonusStats
	weaponSlotStats := stats.Stats{}
	allWeaponSlots := AllWeaponSlots()

	for _, slot := range slots {
		slotStats := ItemEquipmentStats(swapEquipment[slot], false).Subtract(ItemEquipmentStats(originalEquipment[slot], false))
		allSlotStats = allSlotStats.Add(slotStats)

		if slices.Contains(allWeaponSlots, slot) {
			weaponSlotStats = weaponSlotStats.Add(slotStats)
		}
	}

	return ItemSwapStats{
		allSlots:    allSlotStats,
		weaponSlots: weaponSlotStats,
	}
}

func toItem(itemSpec *proto.ItemSpec) Item {
	if itemSpec == nil || itemSpec.Id == 0 {
		return Item{}
	}

	return NewItem(ItemSpec{
		ID:           itemSpec.Id,
		Enchant:      itemSpec.Enchant,
		RandomSuffix: itemSpec.RandomSuffix,
		Rune:         itemSpec.Rune,
	})
}
