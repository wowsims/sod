// Only include this file in the build when we specify the 'with_db' tag.
// Without the tag, the database will start out completely empty.
//go:build with_db

package core

import (
	"github.com/wowsims/sod/assets/database"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	db := database.Load()
	WITH_DB = true

	simDB := &proto.SimDatabase{
		Items:          make([]*proto.SimItem, len(db.Items)),
		Enchants:       make([]*proto.SimEnchant, len(db.Enchants)),
		RandomSuffixes: make([]*proto.ItemRandomSuffix, len(db.RandomSuffixes)),
	}

	for i, item := range db.Items {
		simDB.Items[i] = &proto.SimItem{
			Id:               item.Id,
			Name:             item.Name,
			Type:             item.Type,
			ArmorType:        item.ArmorType,
			WeaponType:       item.WeaponType,
			HandType:         item.HandType,
			RangedWeaponType: item.RangedWeaponType,
			Stats:            item.Stats,
			WeaponDamageMin:  item.WeaponDamageMin,
			WeaponDamageMax:  item.WeaponDamageMax,
			WeaponSpeed:      item.WeaponSpeed,
			SetName:          item.SetName,
			WeaponSkills:     item.WeaponSkills,
		}
	}

	for i, enchant := range db.Enchants {
		simDB.Enchants[i] = &proto.SimEnchant{
			EffectId: enchant.EffectId,
			Stats:    enchant.Stats,
		}
	}

	for i, suffix := range db.RandomSuffixes {
		simDB.RandomSuffixes[i] = &proto.ItemRandomSuffix{
			Id:    suffix.Id,
			Name:  suffix.Name,
			Stats: suffix.Stats,
		}
	}

	addToDatabase(simDB)
}
