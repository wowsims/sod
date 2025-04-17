package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"

	"github.com/wowsims/sod/sim"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	_ "github.com/wowsims/sod/sim/encounters" // Needed for preset encounters.
	"github.com/wowsims/sod/tools"
	"github.com/wowsims/sod/tools/database"
)

// To do a full re-scrape, delete the previous output file first.
// go run ./tools/database/gen_db -outDir=assets -gen=atlasloot
// go run ./tools/database/gen_db -outDir=assets -gen=wowhead-items
// go run ./tools/database/gen_db -outDir=assets -gen=wowhead-spells -maxid=31000
// go run ./tools/database/gen_db -outDir=assets -gen=wowhead-gearplannerdb
// go run ./tools/database/gen_db -outDir=assets -gen=wago-db2-items
// python3 tools/scrape_runes.py assets/db_inputs/wowhead_rune_tooltips.csv
// python3 tools/scrape_shoulder_runes.py assets/db_inputs/wowhead_shoulder_rune_tooltips.csv

// Lastly run the following to generate db.json (ensure to delete cached versions and/or rebuild for copying of assets during local development)
// Note: This does not make network requests, only regenerates core db binary and json files from existing inputs
// go run ./tools/database/gen_db -outDir=assets -gen=db

var exactId = flag.Int("id", 0, "ID to scan for")
var minId = flag.Int("minid", 1, "Minimum ID to scan for")
var maxId = flag.Int("maxid", 31000, "Maximum ID to scan for")
var outDir = flag.String("outDir", "assets", "Path to output directory for writing generated .go files.")
var genAsset = flag.String("gen", "", "Asset to generate. Valid values are 'db', 'atlasloot', 'wowhead-items', 'wowhead-spells', 'wowhead-itemdb', 'wotlk-items', and 'wago-db2-items'")

func main() {
	flag.Parse()

	if *exactId != 0 {
		minId = exactId
		maxId = exactId
	}

	if *outDir == "" {
		panic("outDir flag is required!")
	}

	dbDir := fmt.Sprintf("%s/database", *outDir)
	inputsDir := fmt.Sprintf("%s/db_inputs", *outDir)

	if *genAsset == "atlasloot" {
		db := database.ReadAtlasLootData(inputsDir)
		db.WriteJson(fmt.Sprintf("%s/atlasloot_db.json", inputsDir))
		return
	} else if *genAsset == "wowhead-items" {
		database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId), database.OtherItemIdsToFetch)
		return
	} else if *genAsset == "wowhead-spells" {
		database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId), []string{})
		return
	} else if *genAsset == "wowhead-gearplannerdb" {
		tools.WriteFile(fmt.Sprintf("%s/wowhead_gearplannerdb.txt", inputsDir), tools.ReadWebRequired(core.MakeWowheadUrl("/data/gear-planner?dv=100")))
		return
	} else if *genAsset == "wago-db2-items" {
		tools.WriteFile(fmt.Sprintf("%s/wago_db2_items.csv", inputsDir), tools.ReadWebRequired("https://wago.tools/db2/ItemSparse/csv?build=1.15.7.60249"))
		return
	} else if *genAsset != "db" {
		panic("Invalid gen value")
	}

	itemTooltips := database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Read()
	spellTooltips := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Read()
	runeTooltips := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_rune_tooltips.csv", inputsDir)).Read()
	shoulderRuneTooltips := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_shoulder_rune_tooltips.csv", inputsDir)).Read()
	wowheadDB := database.ParseWowheadDB(tools.ReadFile(fmt.Sprintf("%s/wowhead_gearplannerdb.txt", inputsDir)))
	atlaslootDB := database.ReadDatabaseFromJson(tools.ReadFile(fmt.Sprintf("%s/atlasloot_db.json", inputsDir)))
	wagoItems := database.ParseWagoDB(tools.ReadFile(fmt.Sprintf("%s/wago_db2_items.csv", inputsDir)))

	db := database.NewWowDatabase()
	db.Encounters = core.PresetEncounters

	// Try to filter out items reworked in SoD. We do this by storing the max ID for each item name in the map.
	// This works in most cases because items typically don't share names, however one example of items where this fails is:
	// https://www.wowhead.com/classic/item=23206/mark-of-the-champion and https://www.wowhead.com/classic/item=23207/mark-of-the-champion
	// In this case, we can check the icon to see if they're the same or not.
	// Ultimately we want to get rid of any item with the same name and icon, but a lower ID than another entry
	itemNameMap := make(map[string]string, len(wowheadDB.Items))
	for id, item := range wowheadDB.Items {
		if _, ok := database.ItemDenyList[item.ID]; ok {
			continue
		}

		otherId, hasEntry := itemNameMap[item.Name]
		if !hasEntry {
			itemNameMap[item.Name] = id
			continue
		}

		idInt, _ := strconv.Atoi(id)
		otherIdInt, _ := strconv.Atoi(otherId)
		if otherIdInt < idInt {
			itemNameMap[item.Name] = id
		}
	}
	filteredWHDBItems := core.FilterMap(wowheadDB.Items, func(_ string, item database.WowheadItem) bool {
		id := itemNameMap[item.Name]

		otherItem := wowheadDB.Items[id]

		if _, ok := database.ItemAllowList[item.ID]; ok {
			return true
		}

		if _, ok := database.ItemDenyList[item.ID]; ok {
			return false
		}

		// Most new items follow this pattern:
		// - Higher item ID (this is a given)
		// - Same icon
		// - If the items have a ClassMask they should match
		// - Ilvl either the same or only slightly modified (use a 3 ilvl diff threshold)
		// - Have a later game version
		if otherItem.ID > item.ID &&
			otherItem.Icon == item.Icon &&
			(item.ClassMask == 0 || (otherItem.ClassMask&item.ClassMask) != 0) &&
			math.Abs(float64(otherItem.Ilvl-item.Ilvl)) < 10 &&
			otherItem.Version != item.Version {
			return false
		}

		return true
	})

	for _, response := range itemTooltips {
		if response.IsEquippable() {
			// Item is not part of an item set OR the item set is not in the deny list
			if itemSetID := response.GetItemSetID(); itemSetID == 0 || !slices.Contains(database.DenyItemSetIds, itemSetID) {
				// Only included items that are in wowheads gearplanner db
				// Wowhead doesn't seem to have a field/flag to signify 'not available / in game' but their gearplanner db has them filtered
				item := response.ToItemProto()
				if _, ok := filteredWHDBItems[strconv.Itoa(int(item.Id))]; ok {
					db.MergeItem(item)
				}
			}
		}
	}
	for _, wowheadItem := range filteredWHDBItems {
		item := wowheadItem.ToProto()
		if _, ok := db.Items[item.Id]; ok {
			db.MergeItem(item)
		}
	}
	for _, item := range atlaslootDB.Items {
		if _, ok := db.Items[item.Id]; ok {
			db.MergeItem(item)
		}
	}

	for id, rune := range runeTooltips {
		db.AddRune(id, rune)
	}

	for id, rune := range shoulderRuneTooltips {
		db.AddShoulderRune(id, rune)
	}

	db.MergeItems(database.ItemOverrides)
	db.MergeEnchants(database.EnchantOverrides)
	db.MergeRunes(database.RuneOverrides)
	ApplyGlobalFilters(db)
	AttachFactionInformation(db, wagoItems)

	leftovers := db.Clone()
	ApplyNonSimmableFilters(leftovers)
	leftovers.WriteBinaryAndJson(fmt.Sprintf("%s/leftover_db.bin", dbDir), fmt.Sprintf("%s/leftover_db.json", dbDir))

	ApplySimmableFilters(db)
	for _, enchant := range db.Enchants {
		if enchant.ItemId != 0 {
			db.AddItemIcon(enchant.ItemId, itemTooltips)
		}
		if enchant.SpellId != 0 {
			db.AddSpellIcon(enchant.SpellId, spellTooltips)
		}
	}

	for _, itemID := range database.ExtraItemIcons {
		db.AddItemIcon(itemID, itemTooltips)
	}

	for _, item := range db.Items {
		for _, source := range item.Sources {
			if crafted := source.GetCrafted(); crafted != nil {
				db.AddSpellIcon(crafted.SpellId, spellTooltips)
			}
		}

		for _, randomSuffixID := range item.RandomSuffixOptions {
			if _, exists := db.RandomSuffixes[randomSuffixID]; !exists {
				db.RandomSuffixes[randomSuffixID] = wowheadDB.RandomSuffixes[strconv.Itoa(int(randomSuffixID))].ToProto()
			}
		}
	}

	for _, spellId := range database.SharedSpellsIcons {
		db.AddSpellIcon(spellId, spellTooltips)
	}

	for _, spellIds := range GetAllTalentSpellIds(&inputsDir) {
		for _, spellId := range spellIds {
			db.AddSpellIcon(spellId, spellTooltips)
		}
	}

	for _, spellIds := range GetAllRotationSpellIds() {
		for _, spellId := range spellIds {
			db.AddSpellIcon(spellId, spellTooltips)
		}
	}

	db.MergeSpellIcons(database.SpellIconoverrides)

	atlasDBProto := atlaslootDB.ToUIProto()
	db.MergeZones(atlasDBProto.Zones)
	db.MergeNpcs(atlasDBProto.Npcs)
	db.MergeFactions(atlasDBProto.Factions)

	db.WriteBinaryAndJson(fmt.Sprintf("%s/db.bin", dbDir), fmt.Sprintf("%s/db.json", dbDir))
}

// Filters out entities which shouldn't be included anywhere.
func ApplyGlobalFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := database.ItemDenyList[item.Id]; ok {
			return false
		}

		for _, pattern := range database.DenyListNameRegexes {
			if _, ok := database.ItemAllowList[item.Id]; !ok && pattern.MatchString(item.Name) {
				return false
			}
		}

		return true
	})

	// There is an 'unavailable' version of every naxx set, e.g. https://www.wowhead.com/classic/item=43728/bonescythe-gauntlets
	// heroesItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
	// 	return strings.HasPrefix(item.Name, "Heroes' ")
	// })
	// db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
	// 	nameToMatch := "Heroes' " + item.Name
	// 	for _, heroItem := range heroesItems {
	// 		if heroItem.Name == nameToMatch {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// })

	db.ItemIcons = core.FilterMap(db.ItemIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
	db.SpellIcons = core.FilterMap(db.SpellIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
}

// // AttachFactionInformation attaches faction information (faction restrictions) to the DB items.
func AttachFactionInformation(db *database.WowDatabase, factionRestrictions map[int32]database.WagoDbItem) {
	for _, item := range db.Items {
		if item.FactionRestriction == proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED {
			item.FactionRestriction = factionRestrictions[item.Id].FactionRestriction
		}
	}
}

// Filters out entities which shouldn't be included in the sim.
func ApplySimmableFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, simmableItemFilter)
}
func ApplyNonSimmableFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(id int32, item *proto.UIItem) bool {
		return !simmableItemFilter(id, item)
	})
}
func simmableItemFilter(_ int32, item *proto.UIItem) bool {
	if _, ok := database.ItemAllowList[item.Id]; ok {
		return true
	}

	if item.Quality < proto.ItemQuality_ItemQualityUncommon && item.Ilvl != 61 { // Fix for Rank 1 and 2 Seal of the Dawn trinkets
		return false
	} else if item.Quality == proto.ItemQuality_ItemQualityArtifact {
		return false
	} else if item.Quality > proto.ItemQuality_ItemQualityHeirloom {
		return false
	} else if item.Quality < proto.ItemQuality_ItemQualityEpic {
		if item.Ilvl < 10 {
			return false
		}
		if item.Ilvl < 10 && item.SetName == "" {
			return false
		}
	} else {
		// Epic and legendary items might come from classic, so use a lower ilvl threshold.
		if item.Quality != proto.ItemQuality_ItemQualityHeirloom && item.Ilvl < 10 {
			return false
		}
	}
	if item.Ilvl == 0 {
		fmt.Printf("Missing ilvl: %s\n", item.Name)
	}

	return true
}

type TalentConfig struct {
	FieldName string `json:"fieldName"`
	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	SpellIds  []int32 `json:"spellIds"`
	MaxPoints int32   `json:"maxPoints"`
}

type TalentTreeConfig struct {
	Name          string         `json:"name"`
	BackgroundUrl string         `json:"backgroundUrl"`
	Talents       []TalentConfig `json:"talents"`
}

func getSpellIdsFromTalentJson(infile *string) []int32 {
	data, err := os.ReadFile(*infile)
	if err != nil {
		log.Fatalf("failed to load talent json file: %s", err)
	}

	var buf bytes.Buffer
	err = json.Compact(&buf, []byte(data))
	if err != nil {
		log.Fatalf("failed to compact json: %s", err)
	}

	var talents []TalentTreeConfig

	err = json.Unmarshal(buf.Bytes(), &talents)
	if err != nil {
		log.Fatalf("failed to parse talent to json %s", err)
	}

	spellIds := make([]int32, 0)

	for _, tree := range talents {
		for _, talent := range tree.Talents {
			spellIds = append(spellIds, talent.SpellIds...)

			// Infer omitted spell IDs.
			if len(talent.SpellIds) < int(talent.MaxPoints) {
				curSpellId := talent.SpellIds[len(talent.SpellIds)-1]
				for i := len(talent.SpellIds); i < int(talent.MaxPoints); i++ {
					curSpellId++
					spellIds = append(spellIds, curSpellId)
				}
			}
		}
	}

	return spellIds
}

func GetAllTalentSpellIds(inputsDir *string) map[string][]int32 {
	talentsDir := fmt.Sprintf("%s/../../ui/core/talents/trees", *inputsDir)
	specFiles := []string{
		"druid.json",
		"hunter.json",
		"mage.json",
		"paladin.json",
		"priest.json",
		"rogue.json",
		"shaman.json",
		"warlock.json",
		"warrior.json",
	}

	ret_db := make(map[string][]int32, 0)

	for _, specFile := range specFiles {
		specPath := fmt.Sprintf("%s/%s", talentsDir, specFile)
		ret_db[specFile[:len(specFile)-5]] = getSpellIdsFromTalentJson(&specPath)
	}

	return ret_db

}

func CreateTempAgent(r *proto.Raid) core.Agent {
	encounter := core.MakeSingleTargetEncounter(60, 0.0)
	env, _, _ := core.NewEnvironment(r, encounter, false)
	return env.Raid.Parties[0].Players[0]
}

type RotContainer struct {
	Name string
	Raid *proto.Raid
}

func GetAllRotationSpellIds() map[string][]int32 {
	sim.RegisterAll()

	rotMapping := []RotContainer{
		{Name: "feral", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassDruid,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_FeralDruid{FeralDruid: &proto.FeralDruid{Options: &proto.FeralDruid_Options{}, Rotation: &proto.FeralDruid_Rotation{}}}), nil, nil, nil)},
		{Name: "balance", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassDruid,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_BalanceDruid{BalanceDruid: &proto.BalanceDruid{Options: &proto.BalanceDruid_Options{}}}), nil, nil, nil)},
		{Name: "druid tank", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassDruid,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_FeralTankDruid{FeralTankDruid: &proto.FeralTankDruid{Options: &proto.FeralTankDruid_Options{}}}), nil, nil, nil)},
		{Name: "elemental", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassShaman,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{}}}), nil, nil, nil)},
		{Name: "enhance", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassShaman,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{}}}), nil, nil, nil)},
		{Name: "hunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassHunter,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{}}}), nil, nil, nil)},
		{Name: "mage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Level:         60,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "2352342212231531-5532323123233121-25221213122351351",
		}, &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{}}}), nil, nil, nil)},
		{Name: "shadow", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassPriest,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_ShadowPriest{ShadowPriest: &proto.ShadowPriest{Options: &proto.ShadowPriest_Options{}}}), nil, nil, nil)},
		{Name: "rogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassRogue,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
			Rotation:  &proto.APLRotation{},
		}, &proto.Player_Rogue{Rogue: &proto.Rogue{Options: &proto.RogueOptions{}}}), nil, nil, nil)},
		{Name: "tank rogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassRogue,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
			Rotation:  &proto.APLRotation{},
		}, &proto.Player_TankRogue{TankRogue: &proto.TankRogue{Options: &proto.RogueOptions{}}}), nil, nil, nil)},
		{Name: "warrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassWarrior,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_Warrior{Warrior: &proto.Warrior{Options: &proto.Warrior_Options{}}}), nil, nil, nil)},
		// TODO: Warrior Tank Sim
		{Name: "warrior tank", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassWarrior,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_TankWarrior{TankWarrior: &proto.TankWarrior{Options: &proto.TankWarrior_Options{}}}), nil, nil, nil)},
		{Name: "ret", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassPaladin,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_RetributionPaladin{RetributionPaladin: &proto.RetributionPaladin{Options: &proto.PaladinOptions{}}}), nil, nil, nil)},
		{Name: "warlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassWarlock,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.WarlockOptions{}}}), nil, nil, nil)},
		{Name: "tank warlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassWarlock,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_TankWarlock{TankWarlock: &proto.TankWarlock{Options: &proto.WarlockOptions{}}}), nil, nil, nil)},
	}

	ret_db := make(map[string][]int32, 0)

	for _, r := range rotMapping {
		f := CreateTempAgent(r.Raid).GetCharacter()

		spells := make([]int32, 0, len(f.Spellbook))

		for _, s := range f.Spellbook {
			if s.SpellID != 0 {
				spells = append(spells, s.SpellID)
			}
		}

		for _, s := range f.GetAuras() {
			if s.ActionID.SpellID != 0 {
				spells = append(spells, s.ActionID.SpellID)
			}
		}

		ret_db[r.Name] = spells
	}

	// What follows is generated from wowhead at https://www.wowhead.com/classic/class=9/warlock#spells
	// After you have opened the page with your class you go to the bottom and click on the "Copy" icon
	// that is to the left of the top navigation "<< first prev 1- 5 etc." and select ID. This will allow you
	// to just paste them in here as a list
	// These need to be updated every phase as new spells are added
	ret_db["druid wowhead"] = []int32{24977, 31018, 21850, 25297, 17402, 9885, 9913, 20748, 9858, 25299, 9896, 25298, 9846, 9863, 17329, 9850, 9853, 18658, 9881, 9835, 9867, 9841, 9876, 22829, 22896, 9889, 9827, 17392, 9907, 9904, 9857, 9830, 9901, 9908, 9910, 9912, 9892, 9898, 9834, 9840, 9894, 24976, 21849, 9888, 17401, 9884, 9880, 9866, 20747, 9875, 9862, 16813, 9849, 9852, 22828, 9856, 9845, 8983, 9821, 22895, 9833, 9823, 9839, 9829, 8905, 22812, 22839, 9758, 1824, 9752, 9754, 9756, 17391, 9747, 9749, 9745, 6787, 9750, 8951, 24975, 9000, 9634, 20719, 22827, 16914, 29166, 8907, 8929, 6783, 20742, 8910, 8918, 16812, 5201, 5196, 8903, 18657, 8992, 8955, 6780, 22842, 9005, 8941, 9493, 6793, 8972, 8928, 1823, 3627, 8950, 769, 8914, 9490, 22568, 6778, 6785, 5225, 17390, 24974, 6798, 778, 5234, 20739, 8940, 6800, 740, 783, 5180, 16811, 5209, 3029, 8998, 5195, 8927, 2091, 9492, 2893, 1850, 5189, 6809, 8949, 1822, 8939, 2782, 780, 1075, 5217, 8926, 2090, 5221, 2908, 5179, 768, 1082, 1735, 5188, 6756, 5215, 20484, 1079, 2912, 16810, 1062, 770, 2637, 6808, 8938, 1066, 8925, 1430, 779, 5211, 8946, 5187, 782, 5178, 5229, 8936, 5487, 99, 6795, 5232, 6807, 8924, 1058, 18960, 339, 5186, 467, 5177, 8921, 774, 417141, 417148, 424760, 410029, 410059, 410060, 424765, 410027, 431468, 416051, 431461, 410061, 416050, 416042, 410021, 416049, 410025, 416044, 410028, 431451, 431447, 431449, 410033, 410023, 416046, 424718, 417123, 5185, 414644, 414647, 408124, 408245, 414683, 407993, 407995, 1126, 408247, 437138, 436895, 407988, 410176, 417157, 414684, 414687, 414689, 408024, 417437, 417045, 408120, 5176, 27764, 9077, 1180, 15590, 198, 227, 199, 16940, 16941, 16857, 17002, 24866, 16858, 16859, 16860, 16861, 16862, 16979, 16947, 16948, 16949, 16950, 16951, 16934, 16935, 16936, 16937, 16938, 17056, 17058, 17059, 17060, 17061, 17104, 24943, 24944, 24945, 24946, 17003, 17004, 17005, 17006, 24894, 17079, 17082, 16918, 16919, 16920, 17069, 17070, 17071, 17072, 17073, 17050, 17051, 17053, 17054, 17055, 16821, 16822, 16823, 16824, 16825, 17245, 17247, 17248, 17249, 17074, 17075, 17076, 17077, 17078, 17111, 17112, 17113, 16966, 16968, 16850, 16923, 16924, 16925, 16926, 16836, 16839, 16840, 17123, 17124, 16814, 16815, 16816, 16817, 16818, 5570, 17007, 16896, 16897, 16899, 16900, 16901, 16845, 16846, 16847, 24858, 16833, 16834, 16835, 16902, 16903, 16904, 16905, 16906, 17063, 17065, 17066, 17067, 17068, 16880, 16689, 16819, 16820, 17116, 16864, 16972, 16974, 16975, 16958, 16961, 17106, 17107, 17108, 16998, 16999, 16942, 16943, 16944, 17118, 17119, 17120, 17121, 17122, 18562, 16929, 16930, 16931, 16932, 16933, 24968, 24969, 24970, 24971, 24972, 16909, 16910, 16911, 16912, 16913, 436956}
	ret_db["hunter wowhead"] = []int32{20904, 20906, 24133, 14287, 25296, 15632, 14311, 409519, 5049, 14927, 13544, 25294, 24632, 25295, 425737, 19801, 14268, 14322, 14325, 14271, 13555, 425736, 14295, 20190, 14305, 409530, 14266, 409755, 415343, 14280, 20910, 14317, 409535, 5048, 14290, 24631, 20903, 14286, 13543, 14277, 20905, 24132, 24510, 15631, 24464, 24478, 14926, 24513, 13554, 425735, 24516, 19879, 14294, 14321, 14273, 4202, 444680, 444682, 24562, 14265, 409754, 415342, 20043, 14304, 409528, 14327, 14279, 20902, 14285, 14316, 409534, 13542, 14270, 20909, 4201, 14289, 24561, 14276, 13553, 425734, 24509, 13159, 15630, 24463, 14310, 409512, 24477, 14925, 14324, 24512, 14264, 409752, 415341, 24515, 19882, 1510, 14320, 14267, 20901, 14284, 4200, 14303, 409526, 3662, 24560, 3034, 14272, 13813, 409532, 13552, 425733, 1543, 14263, 409751, 415340, 14275, 19878, 24508, 13161, 15629, 5384, 24441, 24476, 4199, 14924, 14269, 14288, 24559, 24511, 14326, 24514, 20900, 14283, 14319, 13809, 409520, 3661, 14302, 409524, 3045, 13551, 425732, 19880, 1462, 4198, 24558, 14262, 409750, 415338, 19885, 14323, 3043, 24495, 14282, 5118, 781, 14274, 24440, 1499, 409510, 24475, 14923, 3111, 24494, 24490, 14318, 4197, 2643, 24557, 13550, 425730, 19884, 13795, 409521, 1495, 14261, 409748, 415337, 6197, 1002, 1513, 14281, 20736, 4196, 136, 24556, 2974, 13165, 883, 2641, 6991, 4195, 24547, 982, 13549, 425729, 1515, 19883, 5116, 14260, 409693, 415336, 3044, 1130, 13163, 1978, 425728, 3127, 2480, 7919, 7918, 2764, 415423, 75, 425711, 409433, 409495, 409493, 409494, 410114, 416086, 410109, 410116, 416089, 410118, 431601, 416093, 431611, 425759, 410115, 410122, 410113, 410110, 425758, 410121, 410123, 416083, 416085, 416090, 425762, 410111, 425760, 416091, 409552, 409554, 415320, 409580, 409379, 2973, 409691, 415335, 446374, 1494, 9077, 8737, 264, 5011, 1180, 674, 15590, 266, 196, 201, 200, 227, 2567, 197, 202, 19434, 19461, 19462, 24691, 19590, 19592, 19596, 19574, 19239, 19245, 19306, 19295, 19297, 19298, 19300, 19301, 19263, 19416, 19417, 19418, 19419, 19420, 19583, 19584, 19585, 19586, 19587, 19184, 19387, 19388, 19389, 19390, 19598, 19599, 19600, 19601, 19602, 19621, 19622, 19623, 19624, 19625, 19498, 19499, 19500, 19151, 19152, 19153, 19454, 19455, 19456, 19457, 19458, 19552, 19553, 19554, 19555, 19556, 19549, 19550, 19551, 24386, 24387, 19407, 19412, 19413, 19414, 19415, 19557, 19558, 19286, 19287, 19421, 19422, 19423, 19424, 19425, 19572, 19573, 19575, 24443, 19491, 19493, 19494, 19464, 19465, 19466, 19467, 19468, 19228, 19232, 19233, 19234, 19235, 19577, 19370, 19371, 19373, 19426, 19427, 19429, 19430, 19431, 19168, 19180, 19181, 24296, 24297, 24293, 24294, 24295, 19485, 19487, 19488, 19489, 19490, 19559, 19560, 19507, 19508, 19509, 19510, 19511, 19159, 19160, 19503, 19578, 20895, 19290, 19294, 24283, 19255, 19256, 19257, 19258, 19259, 19609, 19610, 19612, 19376, 19377, 19506, 19616, 19617, 19618, 19619, 19620, 19386, 23989, 458436, 458481, 458482, 409593}
	ret_db["mage wowhead"] = []int32{13021, 18809, 10187, 28612, 10140, 10174, 10225, 10151, 25306, 28609, 25304, 10220, 10193, 12826, 28271, 28272, 28270, 13033, 10161, 10054, 22783, 10207, 23028, 10157, 10212, 25345, 10216, 10181, 12526, 10170, 10202, 10199, 400623, 10150, 10230, 13020, 13032, 10186, 10145, 10177, 10192, 10206, 10160, 10139, 10223, 10180, 10219, 11419, 11420, 12525, 10211, 10053, 10173, 10149, 10215, 13031, 10201, 10197, 400622, 22782, 10205, 13019, 10185, 10179, 10191, 12524, 10169, 10156, 10159, 10144, 10148, 8462, 8417, 10138, 8458, 8423, 6131, 7320, 12825, 11416, 11417, 10059, 11418, 8446, 8439, 3552, 8413, 400621, 8408, 13018, 12523, 8427, 8451, 8402, 8495, 8492, 6117, 8445, 8416, 6129, 8422, 8461, 8407, 12522, 8455, 8438, 6127, 8412, 400620, 8457, 8401, 7302, 3565, 3566, 1461, 6141, 759, 8494, 8444, 120, 865, 8406, 12505, 5145, 2139, 8450, 8400, 2121, 8437, 990, 2138, 400616, 6143, 2948, 1953, 10, 5506, 12051, 543, 7301, 7322, 1463, 12824, 3562, 3567, 3561, 3563, 1008, 3140, 475, 5144, 2855, 2120, 1449, 1460, 2137, 400619, 837, 597, 604, 145, 130, 5505, 7300, 122, 5143, 205, 118, 587, 2136, 400618, 143, 5504, 116, 5019, 400610, 400574, 1459, 425124, 401762, 401763, 415948, 401752, 425189, 415934, 429311, 429309, 429308, 401761, 401759, 415942, 401765, 401767, 401722, 440858, 401764, 401757, 401760, 415936, 401754, 429304, 401749, 429306, 425171, 425170, 401768, 415939, 453695, 453690, 453696, 453697, 453694, 453635, 133, 168, 401502, 440802, 440809, 400640, 425121, 400613, 400614, 401556, 412510, 401417, 401462, 442543, 1180, 201, 227, 5009, 11213, 12574, 12575, 12576, 12577, 11222, 12839, 12840, 12841, 12842, 15058, 15059, 15060, 18462, 18463, 18464, 11232, 12500, 12501, 12502, 12503, 12042, 28574, 11210, 12592, 16757, 16758, 11113, 11083, 12351, 12472, 11129, 11115, 11367, 11368, 29438, 29439, 29440, 11124, 12378, 12398, 12399, 12400, 11100, 12353, 11160, 12518, 12519, 11189, 28332, 11071, 12496, 12497, 11426, 11958, 11207, 12672, 15047, 15052, 15053, 11119, 11120, 12846, 12847, 12848, 11103, 12357, 12358, 12359, 12360, 11242, 12467, 12469, 11237, 12463, 12464, 16769, 16770, 11185, 12487, 12488, 11190, 12489, 12490, 11255, 12598, 11078, 11080, 12342, 11094, 13043, 11069, 12338, 12339, 12340, 12341, 11108, 12349, 12350, 11165, 12475, 11070, 12473, 16763, 16765, 16766, 11252, 12605, 11095, 12872, 12873, 18459, 18460, 29441, 29444, 29445, 29446, 29447, 11247, 12606, 29074, 29075, 29076, 11175, 12569, 12571, 11151, 12952, 12953, 12043, 11366, 11170, 12982, 12983, 12984, 12985, 6057, 6085, 11180, 28592, 28593, 28594, 28595, 436949, 401558}
	ret_db["paladin wowhead"] = []int32{20914, 20924, 20928, 20920, 19979, 25291, 25290, 10293, 10314, 415073, 19900, 25898, 25890, 25916, 25895, 25899, 25918, 24239, 25292, 10318, 20773, 20349, 23214, 19943, 20293, 20357, 20930, 19898, 10301, 20729, 19854, 25894, 10308, 10329, 19838, 10313, 415072, 25782, 24274, 20308, 19896, 10326, 20913, 20923, 20927, 20919, 19978, 10292, 1020, 19942, 2812, 10310, 20348, 20292, 20929, 19899, 20772, 20356, 6940, 10328, 10300, 19853, 10312, 415071, 19897, 24275, 19837, 4987, 19941, 20291, 20307, 20912, 20922, 20918, 19977, 1032, 5589, 20347, 19895, 13819, 10278, 3472, 20166, 5627, 5615, 415070, 19891, 10324, 10299, 19852, 642, 19940, 20290, 19836, 19888, 20306, 20116, 20915, 10291, 19752, 1042, 2800, 20165, 5614, 415069, 19876, 1038, 19939, 10298, 20289, 5599, 19850, 5588, 10322, 2878, 19835, 19746, 1026, 20164, 20305, 643, 879, 415068, 19750, 5502, 1044, 5573, 20288, 7294, 25780, 19742, 647, 19834, 7328, 20162, 1022, 10290, 633, 20287, 853, 1152, 498, 639, 21082, 19740, 20271, 107, 3127, 407669, 407613, 407615, 407676, 465, 407804, 407778, 407784, 426175, 426180, 426178, 416035, 416028, 416031, 410013, 429261, 429255, 425619, 410014, 425618, 410015, 409999, 410002, 410001, 429251, 429247, 429242, 429249, 410008, 410010, 416037, 410011, 425621, 412020, 407631, 635, 425600, 407880, 407803, 425609, 407627, 440677, 412019, 407798, 20154, 21084, 9077, 27762, 8737, 750, 9116, 196, 198, 201, 200, 197, 199, 202, 20096, 20097, 20098, 20099, 20100, 20101, 20102, 20103, 20104, 20105, 20217, 20911, 26573, 20117, 20118, 20119, 20120, 20121, 20060, 20061, 20062, 20063, 20064, 20216, 20257, 20258, 20259, 20260, 20261, 20262, 20263, 20264, 20265, 20266, 9799, 25988, 20174, 20175, 20237, 20238, 20239, 5923, 5924, 5925, 5926, 25829, 20925, 20473, 20210, 20212, 20213, 20214, 20215, 20042, 20045, 20046, 20047, 20048, 20244, 20245, 20254, 20255, 20256, 20138, 20139, 20140, 20141, 20142, 20487, 20488, 20489, 25956, 25957, 20234, 20235, 20091, 20092, 20468, 20469, 20470, 20224, 20225, 20330, 20331, 20332, 20335, 20336, 20337, 20359, 20360, 20361, 20196, 20197, 20198, 20199, 20200, 26022, 26023, 20177, 20179, 20180, 20181, 20182, 20127, 20130, 20135, 20136, 20137, 20066, 20218, 20375, 20148, 20149, 20150, 20205, 20206, 20207, 20208, 20209, 20143, 20144, 20145, 20146, 20147, 20111, 20112, 20113, 9453, 25836, 20049, 20056, 20057, 20058, 20059, 9452, 26016, 26021, 435984, 407799}
	ret_db["priest wowhead"] = []int32{27841, 27801, 27871, 18807, 19280, 19293, 10942, 19275, 25314, 19285, 15261, 10952, 10938, 10901, 21564, 10961, 25316, 27681, 25315, 425277, 10955, 19312, 19266, 19243, 10965, 10947, 10912, 20770, 10894, 19305, 10917, 10876, 27683, 10890, 10929, 425276, 10958, 15267, 10900, 10934, 27800, 17314, 19279, 10964, 10946, 10953, 19311, 14819, 27870, 19242, 19292, 10941, 19274, 10916, 19284, 10951, 10960, 10928, 425275, 10893, 19304, 19265, 15266, 10875, 10937, 10899, 21562, 10963, 10945, 10881, 10933, 27799, 17313, 19278, 10915, 10911, 10909, 10927, 425274, 19310, 19241, 15265, 10898, 10888, 10957, 10892, 19303, 14818, 19291, 9592, 19273, 2060, 19283, 1006, 10874, 8106, 996, 9485, 19264, 9474, 6078, 425273, 6060, 15431, 17312, 19277, 988, 15264, 8192, 2791, 6066, 19309, 19240, 6064, 1706, 8105, 10880, 2767, 19302, 552, 9473, 8131, 6077, 425272, 19289, 9579, 19271, 19282, 15263, 602, 605, 6065, 596, 976, 1004, 19262, 15430, 17311, 19276, 6063, 8104, 8124, 19308, 19238, 9472, 6076, 425271, 992, 19299, 15262, 8129, 1245, 3747, 2055, 8103, 2096, 2010, 984, 2944, 2651, 9578, 6346, 13896, 2061, 19281, 14914, 7128, 453, 6075, 425270, 9484, 18137, 19261, 19236, 527, 600, 970, 19296, 2054, 8102, 528, 8122, 6074, 425269, 598, 588, 1244, 592, 13908, 9035, 2053, 8092, 2006, 594, 10797, 2652, 586, 139, 425268, 17, 591, 2052, 589, 5019, 401946, 425294, 425309, 431663, 425310, 425314, 402855, 425312, 431673, 431669, 431705, 415995, 415997, 425215, 425216, 402859, 415996, 402862, 402849, 431650, 402864, 415991, 402852, 425213, 402848, 402854, 402789, 402799, 2050, 413259, 413260, 402004, 402174, 402284, 402289, 425207, 1243, 401859, 401863, 440247, 401955, 401977, 585, 425204, 1180, 198, 227, 5009, 15268, 15323, 15324, 15325, 15326, 27811, 27815, 27816, 15259, 15307, 15308, 15309, 15310, 18530, 18531, 18533, 18534, 18535, 14752, 18544, 18547, 18548, 18549, 18550, 14913, 15012, 15237, 27789, 27790, 14889, 15008, 15009, 15010, 15011, 15274, 15311, 14912, 15013, 15014, 14747, 14770, 14771, 14750, 14772, 15273, 15312, 15313, 15314, 15316, 14749, 14767, 14748, 14768, 14769, 14911, 15018, 15392, 15448, 14908, 15020, 17191, 15275, 15317, 27839, 27840, 14751, 14892, 15362, 15363, 724, 14531, 14774, 14521, 14776, 14777, 14520, 14780, 14781, 14782, 14783, 18551, 18552, 18553, 18554, 18555, 15407, 10060, 14909, 15017, 15272, 15318, 15320, 15260, 15327, 15328, 17322, 17323, 17325, 15257, 15331, 15332, 15333, 15334, 15473, 15487, 14523, 14784, 14785, 14786, 14787, 27900, 27901, 27902, 27903, 27904, 20711, 15270, 15335, 15336, 15337, 15338, 14901, 15028, 15029, 15030, 15031, 14898, 15349, 15354, 15355, 15356, 14522, 14788, 14789, 14790, 14791, 15286, 14524, 14525, 14526, 14527, 14528, 15019, 436951, 424036, 424035, 424037, 424041}
	ret_db["rogue wowhead"] = []int32{20777, 10414, 29228, 25359, 10463, 25357, 416325, 10468, 10601, 10438, 25361, 16362, 17359, 10538, 16387, 10473, 16356, 10428, 10605, 408484, 16342, 10627, 10396, 416324, 15208, 408477, 10432, 10587, 10497, 15112, 10623, 416246, 10479, 16316, 10408, 408345, 11315, 10448, 10467, 10442, 10614, 10462, 15207, 408475, 10437, 25908, 10486, 17354, 20776, 2860, 408482, 10413, 10526, 16355, 10395, 416323, 10431, 10427, 10622, 416245, 16341, 10472, 10586, 10496, 15111, 10466, 10392, 408474, 10600, 16315, 10407, 11314, 10537, 8835, 10613, 1064, 416244, 930, 408481, 10447, 6377, 8005, 416322, 8134, 6365, 8235, 8170, 8249, 10478, 10456, 10391, 408473, 6392, 8161, 20610, 10412, 16339, 8010, 10585, 10495, 15107, 8058, 16314, 6495, 10406, 421, 408479, 408343, 8499, 959, 416320, 6041, 408472, 945, 8012, 8512, 556, 8177, 6375, 10595, 20608, 6364, 8232, 8184, 8053, 8227, 8038, 8008, 6391, 546, 6196, 8030, 943, 408443, 8190, 5675, 20609, 8046, 8181, 939, 416319, 905, 10399, 8155, 8160, 2870, 408342, 8498, 8166, 131, 8056, 8033, 2645, 5394, 8004, 915, 408442, 6363, 8052, 8027, 913, 416318, 6390, 8143, 526, 325, 8019, 8045, 548, 408441, 8154, 2008, 408341, 1535, 547, 416317, 370, 8050, 8024, 3599, 8075, 8044, 529, 408440, 324, 8018, 5730, 2484, 332, 416316, 8042, 8071, 107, 409324, 409333, 409337, 425874, 408514, 408519, 410093, 410100, 416054, 425883, 425882, 410103, 432241, 410105, 432236, 432238, 410096, 416057, 410094, 410098, 436368, 410095, 410104, 425344, 410097, 416066, 416055, 432234, 410099, 410101, 425343, 410107, 415242, 331, 416247, 408490, 408491, 408507, 403, 408439, 425339, 408521, 8017, 435884, 425336, 408696, 408510, 9077, 8737, 9116, 27763, 1180, 674, 15590, 196, 198, 227, 197, 199, 16176, 16235, 16240, 17485, 17486, 17487, 17488, 17489, 16254, 16271, 16272, 16273, 16274, 16038, 16160, 16161, 16041, 16117, 16118, 16119, 16120, 16035, 16105, 16106, 16107, 16108, 16039, 16109, 16110, 16111, 16112, 16043, 16130, 29179, 29180, 30160, 16164, 16089, 16166, 28996, 28997, 28998, 16266, 29079, 29080, 16259, 16295, 29062, 29064, 29065, 16256, 16281, 16282, 16283, 16284, 16258, 16293, 16181, 16230, 16232, 16233, 16234, 29187, 29189, 29191, 29202, 29205, 29206, 16086, 16544, 16262, 16287, 16182, 16226, 16227, 16228, 16229, 16261, 16290, 16184, 16209, 29192, 29193, 16578, 16579, 16580, 16581, 16582, 16190, 16180, 16196, 16198, 16188, 16268, 16178, 16210, 16211, 16212, 16213, 16187, 16205, 16206, 16207, 16208, 16040, 16113, 16114, 16115, 16116, 16299, 16300, 16301, 28999, 29000, 17364, 16255, 16302, 16303, 16304, 16305, 16179, 16214, 16215, 16216, 16217, 16194, 16218, 16219, 16220, 16221, 16173, 16222, 16223, 16224, 16225, 16252, 16306, 16307, 16308, 16309, 16269, 29082, 29084, 29086, 29087, 29088, 415236, 437009}
	ret_db["shaman wowhead"] = []int32{20777, 10414, 29228, 25359, 10463, 25357, 416325, 10468, 10601, 10438, 25361, 16362, 17359, 10538, 16387, 10473, 16356, 10428, 10605, 408484, 16342, 10627, 10396, 416324, 15208, 408477, 10432, 10587, 10497, 15112, 10623, 416246, 10479, 16316, 10408, 408345, 11315, 10448, 10467, 10442, 10614, 10462, 15207, 408475, 10437, 25908, 10486, 17354, 20776, 2860, 408482, 10413, 10526, 16355, 10395, 416323, 10431, 10427, 10622, 416245, 16341, 10472, 10586, 10496, 15111, 10466, 10392, 408474, 10600, 16315, 10407, 11314, 10537, 8835, 10613, 1064, 416244, 930, 408481, 10447, 6377, 8005, 416322, 8134, 6365, 8235, 8170, 8249, 10478, 10456, 10391, 408473, 6392, 8161, 20610, 10412, 16339, 8010, 10585, 10495, 15107, 8058, 16314, 6495, 10406, 421, 408479, 408343, 8499, 959, 416320, 6041, 408472, 945, 8012, 8512, 556, 8177, 6375, 10595, 20608, 6364, 8232, 8184, 8053, 8227, 8038, 8008, 6391, 546, 6196, 8030, 943, 408443, 8190, 5675, 20609, 8046, 8181, 939, 416319, 905, 10399, 8155, 8160, 2870, 408342, 8498, 8166, 131, 8056, 8033, 2645, 5394, 8004, 915, 408442, 6363, 8052, 8027, 913, 416318, 6390, 8143, 526, 325, 8019, 8045, 548, 408441, 8154, 2008, 408341, 1535, 547, 416317, 370, 8050, 8024, 3599, 8075, 8044, 529, 408440, 324, 8018, 5730, 2484, 332, 416316, 8042, 8071, 107, 409324, 409333, 409337, 425874, 408514, 408519, 410093, 410100, 416054, 425883, 425882, 410103, 432241, 410105, 432236, 432238, 410096, 416057, 410094, 410098, 436368, 410095, 410104, 425344, 410097, 416066, 416055, 432234, 410099, 410101, 425343, 410107, 415242, 331, 416247, 408490, 408491, 408507, 403, 408439, 425339, 408521, 8017, 435884, 425336, 408696, 408510, 9077, 8737, 9116, 27763, 1180, 674, 15590, 196, 198, 227, 197, 199, 16176, 16235, 16240, 17485, 17486, 17487, 17488, 17489, 16254, 16271, 16272, 16273, 16274, 16038, 16160, 16161, 16041, 16117, 16118, 16119, 16120, 16035, 16105, 16106, 16107, 16108, 16039, 16109, 16110, 16111, 16112, 16043, 16130, 29179, 29180, 30160, 16164, 16089, 16166, 28996, 28997, 28998, 16266, 29079, 29080, 16259, 16295, 29062, 29064, 29065, 16256, 16281, 16282, 16283, 16284, 16258, 16293, 16181, 16230, 16232, 16233, 16234, 29187, 29189, 29191, 29202, 29205, 29206, 16086, 16544, 16262, 16287, 16182, 16226, 16227, 16228, 16229, 16261, 16290, 16184, 16209, 29192, 29193, 16578, 16579, 16580, 16581, 16582, 16190, 16180, 16196, 16198, 16188, 16268, 16178, 16210, 16211, 16212, 16213, 16187, 16205, 16206, 16207, 16208, 16040, 16113, 16114, 16115, 16116, 16299, 16300, 16301, 28999, 29000, 17364, 16255, 16302, 16303, 16304, 16305, 16179, 16214, 16215, 16216, 16217, 16194, 16218, 16219, 16220, 16221, 16173, 16222, 16223, 16224, 16225, 16252, 16306, 16307, 16308, 16309, 16269, 29082, 29084, 29086, 29087, 29088, 415236, 437009}
	ret_db["warlock wowhead"] = []int32{18932, 18938, 25311, 20757, 17728, 603, 11722, 11735, 11695, 11668, 25309, 18540, 11661, 25307, 403851, 403852, 28610, 23161, 18881, 11730, 11713, 17926, 11678, 17923, 11726, 18871, 17953, 11717, 17937, 6215, 11689, 17924, 18931, 11672, 11700, 403689, 11704, 11684, 17928, 11708, 11675, 11694, 11660, 403848, 11740, 18937, 20756, 11719, 17925, 11734, 11743, 11667, 1122, 17922, 18930, 18870, 18880, 18647, 17727, 11712, 6353, 17952, 11729, 11721, 11699, 403688, 11688, 11677, 11671, 17862, 11703, 11693, 11659, 403844, 11725, 7659, 11707, 6789, 11683, 17921, 11739, 18869, 20755, 11733, 5484, 11665, 5784, 18879, 11711, 2970, 7651, 403687, 8289, 17951, 2362, 3700, 11687, 7641, 403843, 7648, 5699, 6226, 6219, 17920, 18868, 1490, 7646, 6213, 6229, 20752, 1086, 709, 403686, 1949, 2941, 1098, 691, 710, 6366, 6217, 7658, 3699, 1106, 403842, 1714, 132, 1456, 17919, 18867, 6223, 5138, 8288, 5500, 6202, 6205, 699, 403685, 126, 706, 3698, 1094, 5740, 698, 1088, 403841, 713, 712, 693, 1014, 5676, 1455, 5697, 6222, 704, 689, 403677, 1108, 755, 705, 403840, 6201, 696, 1120, 707, 697, 980, 5782, 1454, 695, 403839, 172, 702, 5019, 403629, 412788, 687, 412789, 426445, 426443, 426452, 416017, 416014, 426467, 431758, 431756, 431747, 416009, 403937, 403932, 403920, 403925, 403919, 403938, 403936, 431745, 431743, 425477, 425476, 416008, 416015, 426301, 403501, 403506, 348, 412758, 426241, 426245, 426246, 426247, 426331, 403828, 403789, 412783, 412784, 437169, 686, 403835, 426320, 426325, 688, 1180, 201, 227, 5009, 18119, 18120, 18121, 18122, 18123, 18288, 17788, 17789, 17790, 17791, 17792, 17778, 17779, 17780, 17781, 17782, 17962, 18223, 18220, 18697, 18698, 18699, 18700, 18701, 18788, 17917, 17918, 18130, 18131, 18132, 18133, 18134, 17954, 17955, 17956, 17957, 17958, 17783, 17784, 17785, 17786, 17787, 18708, 18731, 18743, 18744, 18745, 18746, 18751, 18752, 18218, 18219, 17810, 17811, 17812, 17813, 17814, 18827, 18829, 18830, 18310, 18311, 18312, 18313, 18179, 18180, 18181, 17804, 17805, 17806, 17807, 17808, 17864, 18393, 18213, 18372, 18126, 18127, 18767, 18768, 18703, 18704, 18692, 18693, 17815, 17833, 17834, 17835, 17836, 18694, 18695, 18696, 18182, 18183, 18754, 18755, 18756, 17927, 17929, 17930, 17931, 17932, 17793, 17796, 17801, 17802, 17803, 18774, 18775, 18821, 18823, 18824, 18825, 18705, 18706, 18707, 18135, 18136, 23785, 23822, 23823, 23824, 23825, 18709, 18710, 18094, 18095, 18073, 18096, 17959, 18271, 18272, 18273, 18274, 18275, 17877, 18265, 19028, 18174, 18175, 18176, 18177, 18178, 18769, 18770, 18771, 18772, 18773, 412798, 425463, 437032}
	ret_db["warrior wowhead"] = []int32{23894, 21553, 23925, 25289, 20569, 25286, 11585, 11574, 25288, 6554, 11597, 11581, 20662, 11567, 20560, 23893, 21552, 23924, 11556, 7373, 11601, 11605, 11551, 20617, 1672, 11609, 1719, 11573, 23892, 21551, 23923, 20661, 11566, 11580, 11578, 20559, 11604, 11596, 11555, 11584, 11600, 11550, 20616, 11608, 20660, 11565, 11572, 6552, 8820, 8205, 7402, 1680, 11554, 7379, 8380, 11549, 18499, 20658, 7372, 11564, 1671, 2458, 7369, 20252, 6548, 1464, 7887, 871, 8204, 1161, 6178, 7400, 6190, 5308, 1608, 6574, 6192, 5246, 7405, 845, 6547, 20230, 676, 8198, 285, 694, 2565, 1160, 6572, 5242, 7384, 72, 2687, 71, 6546, 7386, 355, 1715, 284, 6343, 100, 772, 461475, 413479, 107, 3127, 2480, 7919, 7918, 2764, 6673, 2457, 403215, 403196, 416004, 409163, 416005, 403467, 416002, 403472, 426491, 427081, 427082, 427084, 403474, 403480, 425444, 425445, 440492, 440496, 440494, 403475, 425443, 416003, 403470, 403489, 427080, 427076, 427078, 425446, 425447, 403476, 453688, 453690, 459313, 453691, 453689, 453694, 453692, 453635, 402913, 78, 12288, 12707, 12708, 403338, 403228, 429765, 402911, 426490, 440488, 402927, 9077, 8737, 750, 9116, 264, 5011, 1180, 674, 15590, 266, 196, 198, 201, 200, 227, 2567, 197, 199, 202, 402974, 12296, 12297, 12750, 12751, 12752, 12753, 12700, 12781, 12783, 12784, 12785, 16487, 16489, 16492, 23881, 12321, 12835, 12836, 12837, 12838, 12809, 12320, 12852, 12853, 12855, 12856, 12328, 12834, 12849, 12867, 12303, 12788, 12789, 12791, 12792, 16462, 16463, 16464, 16465, 16466, 23584, 23585, 23586, 23587, 23588, 12317, 13045, 13046, 13047, 13048, 12319, 12971, 12972, 12973, 12974, 16493, 16494, 12318, 12857, 12858, 12860, 12861, 20500, 20501, 12301, 12818, 12285, 12697, 12329, 12950, 20496, 12324, 12876, 12877, 12878, 12879, 12313, 12804, 12807, 20502, 20503, 12289, 12668, 23695, 12282, 12663, 12664, 20504, 20505, 12290, 12963, 12286, 12658, 12659, 12797, 12799, 12800, 12311, 12958, 12307, 12944, 12945, 12312, 12803, 12330, 12862, 20497, 20498, 20499, 12308, 12810, 12811, 12302, 12765, 12287, 12665, 12666, 12300, 12959, 12960, 12961, 12962, 12975, 12284, 12701, 12702, 12703, 12704, 12294, 16538, 16539, 16540, 16541, 16542, 12323, 12165, 12830, 12831, 12832, 12833, 23922, 12298, 12724, 12725, 12726, 12727, 12292, 12281, 12812, 12813, 12814, 12815, 12295, 12676, 12677, 12678, 12679, 12299, 12761, 12762, 12763, 12764, 12163, 12711, 12712, 12713, 12714, 12322, 12999, 13000, 13001, 13002}

	return ret_db
}
