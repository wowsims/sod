package database

import (
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Note: EffectId AND SpellId are required for all enchants, because they are
// used by various importers/exporters. ItemId is optional.

var EnchantOverrides = []*proto.UIEnchant{
	// Armor Kits
	{EffectId: 15, ItemId: 2304, SpellId: 2831, Name: "Light Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 16, ItemId: 2313, SpellId: 2832, Name: "Medium Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 17, ItemId: 4265, SpellId: 2833, Name: "Heavy Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 18, ItemId: 8173, SpellId: 10344, Name: "Thick Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 32}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 1843, ItemId: 15564, SpellId: 19057, Name: "Rugged Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	// {EffectId: 2503, ItemId: 18251, SpellId: 22725, Name: "Core Armor Kit", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Defense: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},

	// TODO: Arcanums
	// {EffectId: 2544, ItemId: 18330, SpellId: 22844, Name: "Arcanum of Focus", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.SpellPower: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	// {EffectId: 2543, ItemId: 18329, SpellId: 22840, Name: "Arcanum of Rapidity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHaste: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	// {EffectId: 1506, ItemId: 11645, SpellId: 15397, Name: "Arcanum of Voracity", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},
	// {EffectId: 1483, ItemId: 11622, SpellId: 15340, Name: "Arcanum of Rumination", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 150}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},

	// {EffectId: 2590, ItemId: 19789, SpellId: 24167, Name: "Prophetic Aura", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MP5: 4, stats.Stamina: 10, stats.Healing: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs}, EnchantType: proto.EnchantType_EnchantTypeKit},

	// Head
	// {EffectId: 3795, ItemId: 44069, SpellId: 59777, Name: "Arcanum of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},

	// Shoulder
	// {EffectId: 2604, ItemId: 20078, SpellId: 24420, Name: "Zandalar Signet of Serenity", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Healing: 33}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	// {EffectId: 2605, ItemId: 20076, SpellId: 24421, Name: "Zandalar Signet of Mojo", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18, stats.Healing: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	// {EffectId: 2606, ItemId: 20077, SpellId: 24422, Name: "Zandalar Signet of Might", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	// {EffectId: 2715, ItemId: 23547, SpellId: 29475, Name: "Resilience of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Healing: 31, stats.MP5: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	// {EffectId: 2717, ItemId: 23548, SpellId: 29483, Name: "Might of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 26, stats.MeleeCrit: 0.01}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	// {EffectId: 2716, ItemId: 23549, SpellId: 29480, Name: "Fortitude of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 16, stats.Armor: 100}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	// {EffectId: 2721, ItemId: 23545, SpellId: 29467, Name: "Power of the Scourge", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 15, stats.Healing: 15, stats.SpellCrit: 0.01}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},

	// Back
	{EffectId: 2, SpellId: 7454, Name: "Enchant Cloak - Minor Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ArcaneResistance: 1, stats.FrostResistance: 1, stats.FireResistance: 1, stats.NatureResistance: 1, stats.ShadowResistance: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 783, SpellId: 7771, Name: "Enchant Cloak - Minor Protection", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 247, SpellId: 13419, Name: "Enchant Cloak - Minor Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 744, SpellId: 13421, Name: "Enchant Cloak - Lesser Protection", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 256, SpellId: 7861, Name: "Enchant Cloak - Lesser Fire Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FireResistance: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 804, SpellId: 13522, Name: "Enchant Cloak - Lesser Shadow Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ShadowResistance: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 848, SpellId: 13635, Name: "Enchant Cloak - Defense", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 2463, SpellId: 13657, Name: "Enchant Cloak - Fire Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FireResistance: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 884, SpellId: 13746, Name: "Enchant Cloak - Greater Defense", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 903, SpellId: 13794, Name: "Enchant Cloak - Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ArcaneResistance: 3, stats.FireResistance: 3, stats.FrostResistance: 3, stats.NatureResistance: 3, stats.ShadowResistance: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 849, SpellId: 13882, Name: "Enchant Cloak - Lesser Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	// {EffectId: 1888, SpellId: 20014, Name: "Enchant Cloak - Greater Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ArcaneResistance: 5, stats.FireResistance: 5, stats.FrostResistance: 5, stats.NatureResistance: 5, stats.ShadowResistance: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	// {EffectId: 1889, SpellId: 20015, Name: "Enchant Cloak - Superior Defense", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	// {EffectId: 2622, SpellId: 25086, Name: "Enchant Cloak - Dodge", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Dodge: 1 * core.DodgeRatingPerDodgeChance}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	// {EffectId: 2619, SpellId: 25081, Name: "Enchant Cloak - Greater Fire Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FireResistance: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	// {EffectId: 2620, SpellId: 25082, Name: "Enchant Cloak - Greater Nature Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.NatureResistance: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	// {EffectId: 910, SpellId: 25083, Name: "Enchant Cloak - Stealth", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	// {EffectId: 2621, SpellId: 25084, Name: "Enchant Cloak - Subtlety", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},

	// Chest
	{EffectId: 41, SpellId: 7420, Name: "Enchant Chest - Minor Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 24, SpellId: 7443, Name: "Enchant Chest - Minor Mana", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mana: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 44, SpellId: 7426, Name: "Enchant Chest - Minor Absorption", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 242, SpellId: 7748, Name: "Enchant Chest - Lesser Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 246, SpellId: 7776, Name: "Enchant Chest - Lesser Mana", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 254, SpellId: 7857, Name: "Enchant Chest - Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 63, SpellId: 13538, Name: "Enchant Chest - Lesser Absorption", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 843, SpellId: 13607, Name: "Enchant Chest - Mana", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 847, SpellId: 13626, Name: "Enchant Chest - Minor Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 1, stats.Agility: 1, stats.Strength: 1, stats.Intellect: 1, stats.Spirit: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 850, SpellId: 13640, Name: "Enchant Chest - Greater Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 857, SpellId: 13663, Name: "Enchant Chest - Greater Mana", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 7223, SpellId: 435903, Name: "Enchant Chest - Retricutioner", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 866, SpellId: 13700, Name: "Enchant Chest - Lesser Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 2, stats.Agility: 2, stats.Strength: 2, stats.Intellect: 2, stats.Spirit: 2}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 908, SpellId: 13858, Name: "Enchant Chest - Superior Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	// {EffectId: 913, SpellId: 13917, Name: "Enchant Chest - Superior Mana", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	// {EffectId: 928, SpellId: 13941, Name: "Enchant Chest - Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 3, stats.Agility: 3, stats.Strength: 3, stats.Intellect: 3, stats.Spirit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	// {EffectId: 1892, SpellId: 20026, Name: "Enchant Chest - Major Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 100}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	// {EffectId: 1893, SpellId: 20028, Name: "Enchant Chest - Major Mana", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mana: 100}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	// {EffectId: 1891, SpellId: 20025, Name: "Enchant Chest - Greater Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 4, stats.Agility: 4, stats.Strength: 4, stats.Intellect: 4, stats.Spirit: 4}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},

	// Wrist
	{EffectId: 66, SpellId: 7457, Name: "Enchant Bracer - Minor Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 243, SpellId: 7766, Name: "Enchant Bracer - Minor Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 41, SpellId: 7418, Name: "Enchant Bracer - Minor Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 924, SpellId: 7428, Name: "Enchant Bracer - Minor Deflect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Defense: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 247, SpellId: 7779, Name: "Enchant Bracer - Minor Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 248, SpellId: 7782, Name: "Enchant Bracer - Minor Strength", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 255, SpellId: 7859, Name: "Enchant Bracer - Lesser Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 724, SpellId: 13501, Name: "Enchant Bracer - Lesser Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 823, SpellId: 13536, Name: "Enchant Bracer - Lesser Strength", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 723, SpellId: 13622, Name: "Enchant Bracer - Lesser Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 851, SpellId: 13642, Name: "Enchant Bracer - Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 925, SpellId: 13646, Name: "Enchant Bracer - Lesser Deflection", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Defense: 2}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 852, SpellId: 13648, Name: "Enchant Bracer - Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 856, SpellId: 13661, Name: "Enchant Bracer - Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 905, SpellId: 13822, Name: "Enchant Bracer - Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 907, SpellId: 13846, Name: "Enchant Bracer - Greater Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 923, SpellId: 13931, Name: "Enchant Bracer - Deflection", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Defense: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 927, SpellId: 13939, Name: "Enchant Bracer - Greater Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 929, SpellId: 13945, Name: "Enchant Bracer - Greater Stamina", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 1883, SpellId: 20008, Name: "Enchant Bracer - Greater Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Intellect: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 1884, SpellId: 20009, Name: "Enchant Bracer - Superior Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 2565, SpellId: 23801, Name: "Enchant Bracer - Mana Regeneration", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MP5: 4}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 1885, SpellId: 20010, Name: "Enchant Bracer - Superior Strength", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 2566, SpellId: 23802, Name: "Enchant Bracer - Healing Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Healing: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	// {EffectId: 1886, SpellId: 20011, Name: "Enchant Bracer - Superior Stamina", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},

	// Hands
	// {EffectId: 2614, SpellId: 25073, Name: "Shadow Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 846, SpellId: 13620, Name: "Enchant Gloves - Fishing", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 845, SpellId: 13617, Name: "Enchant Gloves - Herbalism", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 844, SpellId: 13612, Name: "Enchant Gloves - Mining", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 865, SpellId: 13698, Name: "Enchant Gloves - Skinning", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 904, SpellId: 13815, Name: "Enchant Gloves - Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 906, SpellId: 13841, Name: "Enchant Gloves - Advanced Mining", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 909, SpellId: 13868, Name: "Enchant Gloves - Advanced Herbalism", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 856, SpellId: 13887, Name: "Enchant Gloves - Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 931, SpellId: 13948, Name: "[NYI] Enchant Gloves - Minor Haste", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 930, SpellId: 13947, Name: "Enchant Gloves - Riding Skill", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 1887, SpellId: 20012, Name: "Enchant Gloves - Greater Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 927, SpellId: 20013, Name: "Enchant Gloves - Greater Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 2616, SpellId: 25078, Name: "Enchant Gloves - Fire Power", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.FirePower: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 2615, SpellId: 25074, Name: "Enchant Gloves - Frost Power", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.FrostPower: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 2617, SpellId: 25079, Name: "Enchant Gloves - Healing Power", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Healing: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 2614, SpellId: 25073, Name: "Enchant Gloves - Shadow Power", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.ShadowPower: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 2564, SpellId: 25080, Name: "Enchant Gloves - Superior Agility", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Agility: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// {EffectId: 2613, SpellId: 25072, Name: "Enchant Gloves - Threat", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},

	// Feet TODO: Classic
	{EffectId: 247, SpellId: 7867, Name: "Enchant Boots - Minor Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 66, SpellId: 7863, Name: "Enchant Boots - Minor Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 849, SpellId: 13637, Name: "Enchant Boots - Lesser Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 724, SpellId: 13644, Name: "Enchant Boots - Lesser Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 255, SpellId: 13687, Name: "Enchant Boots - Lesser Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 852, SpellId: 13836, Name: "Enchant Boots - Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 911, SpellId: 13890, Name: "Enchant Boots - Minor Speed", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	// {EffectId: 904, SpellId: 13935, Name: "Enchant Boots - Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	// {EffectId: 929, SpellId: 20020, Name: "Enchant Boots - Greater Stamina", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	// {EffectId: 851, SpellId: 20024, Name: "Enchant Boots - Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	// {EffectId: 1887, SpellId: 20023, Name: "Enchant Boots - Greater Agility", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},

	// Weapon
	{EffectId: 249, SpellId: 7786, Name: "Enchant Weapon - Minor Beastslayer", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 250, SpellId: 7788, Name: "Enchant Weapon - Minor Striking", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 241, SpellId: 13503, Name: "Enchant Weapon - Lesser Striking", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 853, SpellId: 13653, Name: "Enchant Weapon - Lesser Beastslayer", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 854, SpellId: 13655, Name: "Enchant Weapon - Lesser Elemental Slayer", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 2443, SpellId: 21931, Name: "Enchant Weapon - Winter's Might", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FrostPower: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 943, SpellId: 13693, Name: "Enchant Weapon - Striking", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 7210, SpellId: 435481, Name: "Enchant Weapon - Dismantle", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 912, SpellId: 13915, Name: "Enchant Weapon - Demonslaying", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 805, SpellId: 13943, Name: "Enchant Weapon - Greater Striking", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 803, SpellId: 13898, Name: "Enchant Weapon - Fiery Weapon", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 1894, SpellId: 20029, Name: "Enchant Weapon - Icy Chill", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 2564, SpellId: 23800, Name: "Enchant Weapon - Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 2563, SpellId: 23799, Name: "Enchant Weapon - Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 1899, SpellId: 20033, Name: "Enchant Weapon - Unholy Weapon", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 1900, SpellId: 20034, Name: "Enchant Weapon - Crusader", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 2505, SpellId: 22750, Name: "Enchant Weapon - Healing Power", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Healing: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 1898, SpellId: 20032, Name: "Enchant Weapon - Lifestealing", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 2568, SpellId: 23804, Name: "Enchant Weapon - Mighty Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 2567, SpellId: 23803, Name: "Enchant Weapon - Mighty Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 2504, SpellId: 22749, Name: "Enchant Weapon - Spell Power", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellDamage: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	// {EffectId: 1897, SpellId: 20031, Name: "Enchant Weapon - Superior Striking", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 36, SpellId: 6296, Name: "Fiery Blaze Enchantment", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},

	// 2H Weapon
	{EffectId: 723, SpellId: 7793, Name: "Enchant 2H Weapon - Lesser Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 241, SpellId: 7745, Name: "Enchant 2H Weapon - Minor Impact", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 255, SpellId: 13380, Name: "Enchant 2H Weapon - Lesser Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 943, SpellId: 13529, Name: "Enchant 2H Weapon - Lesser Impact", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 1897, SpellId: 13695, Name: "Enchant 2H Weapon - Impact", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	// {EffectId: 963, SpellId: 13937, Name: "Enchant 2H Weapon - Greater Impact", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	// {EffectId: 2646, SpellId: 27837, Name: "Enchant 2H Weapon - Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	// {EffectId: 1896, SpellId: 20030, Name: "Enchant 2H Weapon - Superior Impact", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	// {EffectId: 1904, SpellId: 20036, Name: "Enchant 2H Weapon - Major Intellect", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Intellect: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	// {EffectId: 1903, SpellId: 20035, Name: "Enchant 2H Weapon - Major Spirit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 9}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 34, SpellId: 7218, Name: "Iron Counterweight", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHaste: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},

	// Ranged Scopes
	{EffectId: 30, ItemId: 4405, SpellId: 3974, Name: "Crude Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 32, ItemId: 4406, SpellId: 3975, Name: "Standard Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 33, ItemId: 4407, SpellId: 3976, Name: "Accurate Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 663, ItemId: 10546, SpellId: 12459, Name: "Deadly Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 664, ItemId: 10548, SpellId: 12460, Name: "Sniper Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},

	// TODO: Add ranged hit stat
	//{EffectId: 2523, ItemId: 18283, SpellId: 22779, Name: "Biznicks 247x128 Accurascope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MeleeHit: 3}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
}
