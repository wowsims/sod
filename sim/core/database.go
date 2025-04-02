package core

import (
	"fmt"
	"slices"
	"sync"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

var WITH_DB = false

var rwMutex = sync.RWMutex{}
var ItemsByID = map[int32]Item{}
var RandomSuffixesByID = map[int32]RandomSuffix{}
var EnchantsByEffectID = map[int32]Enchant{}

func addToDatabase(newDB *proto.SimDatabase) {
	for _, v := range newDB.Items {
		rwMutex.Lock()
		if _, ok := ItemsByID[v.Id]; !ok {
			ItemsByID[v.Id] = ItemFromProto(v)
		}
		rwMutex.Unlock()
	}

	for _, v := range newDB.RandomSuffixes {
		rwMutex.Lock()
		if _, ok := RandomSuffixesByID[v.Id]; !ok {
			RandomSuffixesByID[v.Id] = RandomSuffixFromProto(v)
		}
		rwMutex.Unlock()
	}

	for _, v := range newDB.Enchants {
		rwMutex.Lock()
		if _, ok := EnchantsByEffectID[v.EffectId]; !ok {
			EnchantsByEffectID[v.EffectId] = EnchantFromProto(v)
		}
		rwMutex.Unlock()
	}
}

type Item struct {
	ID             int32
	RequiresLevel  int32
	ClassAllowlist []proto.Class
	Type           proto.ItemType
	ArmorType      proto.ArmorType

	// Weapon Stats
	WeaponType       proto.WeaponType
	HandType         proto.HandType
	RangedWeaponType proto.RangedWeaponType
	WeaponDamageMin  float64
	WeaponDamageMax  float64
	SwingSpeed       float64

	Name                string
	Stats               stats.Stats // Stats applied to wearer
	BonusPhysicalDamage float64
	BonusPeriodicPct    int32

	Quality      proto.ItemQuality
	SetName      string // Empty string if not part of a set.
	SetID        int32  // 0 if not part of a set.
	WeaponSkills stats.WeaponSkills

	Timeworn   bool
	Sanctified bool

	// Modified for each instance of the item.
	RandomSuffix RandomSuffix
	Enchant      Enchant
	Rune         int32

	//Internal use
	TempEnchant int32
}

func ItemFromProto(pData *proto.SimItem) Item {
	return Item{
		ID:                  pData.Id,
		RequiresLevel:       pData.RequiresLevel,
		ClassAllowlist:      pData.ClassAllowlist,
		Name:                pData.Name,
		Type:                pData.Type,
		ArmorType:           pData.ArmorType,
		WeaponType:          pData.WeaponType,
		HandType:            pData.HandType,
		RangedWeaponType:    pData.RangedWeaponType,
		WeaponDamageMin:     pData.WeaponDamageMin,
		WeaponDamageMax:     pData.WeaponDamageMax,
		SwingSpeed:          pData.WeaponSpeed,
		Stats:               stats.FromFloatArray(pData.Stats),
		BonusPhysicalDamage: pData.BonusPhysicalDamage,
		BonusPeriodicPct:    pData.BonusPeriodicPct,
		SetName:             pData.SetName,
		SetID:               pData.SetId,
		WeaponSkills:        stats.WeaponSkillsFloatArray(pData.WeaponSkills),
		Timeworn:            pData.Timeworn,
		Sanctified:          pData.Sanctified,
	}
}

func (item *Item) IsWeapon() bool {
	return !slices.Contains([]proto.WeaponType{proto.WeaponType_WeaponTypeUnknown, proto.WeaponType_WeaponTypeShield, proto.WeaponType_WeaponTypeOffHand}, item.WeaponType)
}

func (item *Item) ToItemSpecProto() *proto.ItemSpec {
	return &proto.ItemSpec{
		Id:           item.ID,
		RandomSuffix: item.RandomSuffix.ID,
		Enchant:      item.Enchant.EffectID,

		Rune: item.Rune,
	}
}

type RandomSuffix struct {
	ID            int32
	Name          string
	Stats         stats.Stats
	EnchantIDList []int32
}

func RandomSuffixFromProto(pData *proto.ItemRandomSuffix) RandomSuffix {
	return RandomSuffix{
		ID:            pData.Id,
		Name:          pData.Name,
		Stats:         stats.FromFloatArray(pData.Stats),
		EnchantIDList: pData.EnchantIdList,
	}
}

type Enchant struct {
	EffectID int32 // Used by UI to apply effect to tooltip
	Stats    stats.Stats
}

func EnchantFromProto(pData *proto.SimEnchant) Enchant {
	return Enchant{
		EffectID: pData.EffectId,
		Stats:    stats.FromFloatArray(pData.Stats),
	}
}

type Rune struct {
	ID int32
}

func RuneFromProto(pData *proto.SimRune) Rune {
	return Rune{
		ID: pData.Id,
	}
}

type ItemSpec struct {
	ID           int32
	RandomSuffix int32
	Enchant      int32
	Rune         int32
}

type Equipment [proto.ItemSlot_ItemSlotRanged + 1]Item

func (equipment *Equipment) MainHand() *Item {
	return &equipment[proto.ItemSlot_ItemSlotMainHand]
}

func (equipment *Equipment) OffHand() *Item {
	return &equipment[proto.ItemSlot_ItemSlotOffHand]
}

func (equipment *Equipment) Ranged() *Item {
	return &equipment[proto.ItemSlot_ItemSlotRanged]
}

func (equipment *Equipment) Head() *Item {
	return &equipment[proto.ItemSlot_ItemSlotHead]
}

func (equipment *Equipment) Shoulders() *Item {
	return &equipment[proto.ItemSlot_ItemSlotShoulder]
}

func (equipment *Equipment) Hands() *Item {
	return &equipment[proto.ItemSlot_ItemSlotHands]
}

func (equipment *Equipment) Neck() *Item {
	return &equipment[proto.ItemSlot_ItemSlotNeck]
}

func (equipment *Equipment) Trinket1() *Item {
	return &equipment[proto.ItemSlot_ItemSlotTrinket1]
}

func (equipment *Equipment) Trinket2() *Item {
	return &equipment[proto.ItemSlot_ItemSlotTrinket2]
}

func (equipment *Equipment) Finger1() *Item {
	return &equipment[proto.ItemSlot_ItemSlotFinger1]
}

func (equipment *Equipment) Finger2() *Item {
	return &equipment[proto.ItemSlot_ItemSlotFinger2]
}

func (equipment *Equipment) EquipItem(item Item) {
	if item.Type == proto.ItemType_ItemTypeFinger {
		if equipment.Finger1().ID == 0 {
			*equipment.Finger1() = item
		} else {
			*equipment.Finger2() = item
		}
	} else if item.Type == proto.ItemType_ItemTypeTrinket {
		if equipment.Trinket1().ID == 0 {
			*equipment.Trinket1() = item
		} else {
			*equipment.Trinket2() = item
		}
	} else if item.Type == proto.ItemType_ItemTypeWeapon {
		if item.WeaponType == proto.WeaponType_WeaponTypeShield && equipment.MainHand().HandType != proto.HandType_HandTypeTwoHand {
			*equipment.OffHand() = item
		} else if item.HandType == proto.HandType_HandTypeMainHand || item.HandType == proto.HandType_HandTypeUnknown {
			*equipment.MainHand() = item
		} else if item.HandType == proto.HandType_HandTypeOffHand {
			*equipment.OffHand() = item
		} else if item.HandType == proto.HandType_HandTypeOneHand || item.HandType == proto.HandType_HandTypeTwoHand {
			if equipment.MainHand().ID == 0 {
				*equipment.MainHand() = item
			} else if equipment.OffHand().ID == 0 {
				*equipment.OffHand() = item
			}
		}
	} else {
		equipment[ItemTypeToSlot(item.Type)] = item
	}
}

func (equipment *Equipment) containsEnchantInSlot(effectID int32, slot proto.ItemSlot) bool {
	equipmentSlot := equipment[slot]
	if equipmentSlot.Enchant.EffectID == effectID || equipmentSlot.TempEnchant == effectID {
		return true
	}

	if len(equipmentSlot.RandomSuffix.EnchantIDList) > 0 {
		for _, enchantID := range equipmentSlot.RandomSuffix.EnchantIDList {
			if enchantID == effectID {
				return true
			}
		}
	}

	return false
}

func (equipment *Equipment) containsEnchantInSlots(effectID int32, possibleSlots []proto.ItemSlot) bool {
	return slices.ContainsFunc(possibleSlots, func(slot proto.ItemSlot) bool {
		return equipment.containsEnchantInSlot(effectID, slot)
	})
}

func (equipment *Equipment) containsItemInSlots(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return slices.ContainsFunc(possibleSlots, func(slot proto.ItemSlot) bool {
		return equipment[slot].ID == itemID
	})
}

func (equipment *Equipment) GetEnchantCount(effectID int32) int32 {
	count := int32(0)

	for itemSlot := proto.ItemSlot(0); itemSlot < NumItemSlots; itemSlot++ {
		if equipment.containsEnchantInSlot(effectID, itemSlot) {
			count++
		}
	}

	return count
}

func (equipment *Equipment) EligibleSlotsForEffect(effectID int32) []proto.ItemSlot {
	var eligibleSlots []proto.ItemSlot

	for itemSlot := proto.ItemSlot(0); itemSlot < NumItemSlots; itemSlot++ {
		if equipment.containsEnchantInSlot(effectID, itemSlot) {
			eligibleSlots = append(eligibleSlots, itemSlot)
		}
	}

	return eligibleSlots
}

func (equipment *Equipment) ToEquipmentSpecProto() *proto.EquipmentSpec {
	return &proto.EquipmentSpec{
		Items: MapSlice(equipment[:], func(item Item) *proto.ItemSpec {
			return item.ToItemSpecProto()
		}),
	}
}

// Structs used for looking up items/enchants
type EquipmentSpec [proto.ItemSlot_ItemSlotRanged + 1]ItemSpec

func ProtoToEquipmentSpec(es *proto.EquipmentSpec) EquipmentSpec {
	var coreEquip EquipmentSpec
	for i, item := range es.Items {
		coreEquip[i] = ItemSpec{
			ID:           item.Id,
			RandomSuffix: item.RandomSuffix,
			Enchant:      item.Enchant,
			Rune:         item.Rune,
		}
	}
	return coreEquip
}

func NewItem(itemSpec ItemSpec) Item {
	item := Item{}
	if foundItem, ok := ItemsByID[itemSpec.ID]; ok {
		item = foundItem
	} else {
		panic(fmt.Sprintf("No item with id: %d", itemSpec.ID))
	}

	if itemSpec.RandomSuffix != 0 {
		if randomSuffix, ok := RandomSuffixesByID[itemSpec.RandomSuffix]; ok {
			item.RandomSuffix = randomSuffix
		} else {
			panic(fmt.Sprintf("No random suffix with id: %d", itemSpec.RandomSuffix))
		}
	}

	if itemSpec.Enchant != 0 {
		if enchant, ok := EnchantsByEffectID[itemSpec.Enchant]; ok {
			item.Enchant = enchant
		}
		// else {
		// 	panic(fmt.Sprintf("No enchant with id: %d", itemSpec.Enchant))
		// }
	}

	if itemSpec.Rune != 0 {
		item.Rune = itemSpec.Rune
		// if rune, ok := RuneBySpellID[itemSpec.Rune]; ok {
		// 	item.Rune = rune.ID
		// }
	}

	return item
}

func NewEquipmentSet(equipSpec EquipmentSpec) Equipment {
	equipment := Equipment{}
	for _, itemSpec := range equipSpec {
		if itemSpec.ID != 0 {
			equipment.EquipItem(NewItem(itemSpec))
		}
	}
	return equipment
}

func ProtoToEquipment(es *proto.EquipmentSpec) Equipment {
	return NewEquipmentSet(ProtoToEquipmentSpec(es))
}

// Like ItemSpec, but uses names for reference instead of ID.
type ItemStringSpec struct {
	Name    string
	Enchant string
}

func EquipmentSpecFromJsonString(jsonString string) *proto.EquipmentSpec {
	es := &proto.EquipmentSpec{}

	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, es); err != nil {
		panic(err)
	}
	return es
}

func ItemSwapFromJsonString(jsonString string) *proto.ItemSwap {
	is := &proto.ItemSwap{}

	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, is); err != nil {
		panic(err)
	}
	return is
}

func (equipment *Equipment) Stats() stats.Stats {
	equipStats := stats.Stats{}
	for _, item := range equipment {
		equipStats = equipStats.Add(ItemEquipmentStats(item, false))
	}
	return equipStats
}

func (equipment *Equipment) BaseStats() stats.Stats {
	equipStats := stats.Stats{}
	for _, item := range equipment {
		equipStats = equipStats.Add(ItemEquipmentStats(item, true))
	}
	return equipStats
}

func ItemEquipmentStats(item Item, includeOnlyBaseStats bool) stats.Stats {
	equipStats := stats.Stats{}

	if item.ID == 0 {
		return equipStats
	}

	equipStats = equipStats.Add(item.Stats)
	equipStats = equipStats.Add(item.RandomSuffix.Stats)
	if !includeOnlyBaseStats {
		equipStats = equipStats.Add(item.Enchant.Stats)
	}

	return equipStats
}

func GetItemByID(id int32) *Item {
	if item, ok := ItemsByID[id]; ok {
		return &item
	}
	return nil
}

func (equipment *Equipment) GetRuneIds() []int32 {
	out := make([]int32, len(equipment))
	for _, v := range equipment {
		if v.Rune != 0 {
			out = append(out, v.Rune)
		}
	}
	return out
}

func ItemTypeToSlot(it proto.ItemType) proto.ItemSlot {
	switch it {
	case proto.ItemType_ItemTypeHead:
		return proto.ItemSlot_ItemSlotHead
	case proto.ItemType_ItemTypeNeck:
		return proto.ItemSlot_ItemSlotNeck
	case proto.ItemType_ItemTypeShoulder:
		return proto.ItemSlot_ItemSlotShoulder
	case proto.ItemType_ItemTypeBack:
		return proto.ItemSlot_ItemSlotBack
	case proto.ItemType_ItemTypeChest:
		return proto.ItemSlot_ItemSlotChest
	case proto.ItemType_ItemTypeWrist:
		return proto.ItemSlot_ItemSlotWrist
	case proto.ItemType_ItemTypeHands:
		return proto.ItemSlot_ItemSlotHands
	case proto.ItemType_ItemTypeWaist:
		return proto.ItemSlot_ItemSlotWaist
	case proto.ItemType_ItemTypeLegs:
		return proto.ItemSlot_ItemSlotLegs
	case proto.ItemType_ItemTypeFeet:
		return proto.ItemSlot_ItemSlotFeet
	case proto.ItemType_ItemTypeFinger:
		return proto.ItemSlot_ItemSlotFinger1
	case proto.ItemType_ItemTypeTrinket:
		return proto.ItemSlot_ItemSlotTrinket1
	case proto.ItemType_ItemTypeWeapon:
		return proto.ItemSlot_ItemSlotMainHand
	case proto.ItemType_ItemTypeRanged:
		return proto.ItemSlot_ItemSlotRanged
	}

	return 255
}

// See getEligibleItemSlots in proto_utils/utils.ts.
var itemTypeToSlotsMap = map[proto.ItemType][]proto.ItemSlot{
	proto.ItemType_ItemTypeHead:     {proto.ItemSlot_ItemSlotHead},
	proto.ItemType_ItemTypeNeck:     {proto.ItemSlot_ItemSlotNeck},
	proto.ItemType_ItemTypeShoulder: {proto.ItemSlot_ItemSlotShoulder},
	proto.ItemType_ItemTypeBack:     {proto.ItemSlot_ItemSlotBack},
	proto.ItemType_ItemTypeChest:    {proto.ItemSlot_ItemSlotChest},
	proto.ItemType_ItemTypeWrist:    {proto.ItemSlot_ItemSlotWrist},
	proto.ItemType_ItemTypeHands:    {proto.ItemSlot_ItemSlotHands},
	proto.ItemType_ItemTypeWaist:    {proto.ItemSlot_ItemSlotWaist},
	proto.ItemType_ItemTypeLegs:     {proto.ItemSlot_ItemSlotLegs},
	proto.ItemType_ItemTypeFeet:     {proto.ItemSlot_ItemSlotFeet},
	proto.ItemType_ItemTypeFinger:   {proto.ItemSlot_ItemSlotFinger1, proto.ItemSlot_ItemSlotFinger2},
	proto.ItemType_ItemTypeTrinket:  {proto.ItemSlot_ItemSlotTrinket1, proto.ItemSlot_ItemSlotTrinket2},
	proto.ItemType_ItemTypeRanged:   {proto.ItemSlot_ItemSlotRanged},
	// ItemType_ItemTypeWeapon is excluded intentionally - the slot cannot be decided based on type alone for weapons.
}

func eligibleSlotsForItem(item *Item) []proto.ItemSlot {
	if item == nil {
		return nil
	}

	if slots, ok := itemTypeToSlotsMap[item.Type]; ok {
		return slots
	}

	if item.Type == proto.ItemType_ItemTypeWeapon {
		switch item.HandType {
		case proto.HandType_HandTypeTwoHand, proto.HandType_HandTypeMainHand:
			return []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}
		case proto.HandType_HandTypeOffHand:
			return []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}
		case proto.HandType_HandTypeOneHand:
			return []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}
		}
	}

	return nil
}
