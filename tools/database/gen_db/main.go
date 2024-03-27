package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
// python3 tools/scrape_runes.py assets/db_inputs/wowhead_rune_tooltips.csv

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
		db := database.ReadAtlasLootData()
		db.WriteJson(fmt.Sprintf("%s/atlasloot_db.json", inputsDir))
		return
	} else if *genAsset == "wowhead-items" {
		database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId), database.OtherItemIdsToFetch)
		return
	} else if *genAsset == "wowhead-spells" {
		database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Fetch(int32(*minId), int32(*maxId), []string{})
		return
	} else if *genAsset == "wowhead-gearplannerdb" {
		tools.WriteFile(fmt.Sprintf("%s/wowhead_gearplannerdb.txt", inputsDir), tools.ReadWebRequired("https://nether.wowhead.com/classic/data/gear-planner?dv=100"))
		return
	} else if *genAsset != "db" {
		panic("Invalid gen value")
	}

	itemTooltips := database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Read()
	spellTooltips := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Read()
	runeTooltips := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_rune_tooltips.csv", inputsDir)).Read()
	wowheadDB := database.ParseWowheadDB(tools.ReadFile(fmt.Sprintf("%s/wowhead_gearplannerdb.txt", inputsDir)))
	atlaslootDB := database.ReadDatabaseFromJson(tools.ReadFile(fmt.Sprintf("%s/atlasloot_db.json", inputsDir)))
	// factionRestrictions := database.ParseItemFactionRestrictionsFromWagoDB(tools.ReadFile(fmt.Sprintf("%s/wago_db2_items.csv", inputsDir)))

	db := database.NewWowDatabase()
	db.Encounters = core.PresetEncounters

	for _, response := range itemTooltips {
		if response.IsEquippable() {
			// Only included items that are in wowheads gearplanner db
			// Wowhead doesn't seem to have a field/flag to signify 'not available / in game' but their gearplanner db has them filtered
			item := response.ToItemProto()
			if _, ok := wowheadDB.Items[strconv.Itoa(int(item.Id))]; ok {
				db.MergeItem(item)
			}
		}
	}
	for _, wowheadItem := range wowheadDB.Items {
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

	db.MergeItems(database.ItemOverrides)
	db.MergeEnchants(database.EnchantOverrides)
	db.MergeRunes(database.RuneOverrides)
	ApplyGlobalFilters(db)
	// AttachFactionInformation(db, factionRestrictions)

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

	atlasDBProto := atlaslootDB.ToUIProto()
	db.MergeZones(atlasDBProto.Zones)
	db.MergeNpcs(atlasDBProto.Npcs)

	db.WriteBinaryAndJson(fmt.Sprintf("%s/db.bin", dbDir), fmt.Sprintf("%s/db.json", dbDir))
}

// Filters out entities which shouldn't be included anywhere.
func ApplyGlobalFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := database.ItemDenyList[item.Id]; ok {
			return false
		}

		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(item.Name) {
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

// AttachFactionInformation attaches faction information (faction restrictions) to the DB items.
// func AttachFactionInformation(db *database.WowDatabase, factionRestrictions map[int32]proto.UIItem_FactionRestriction) {
// 	for _, item := range db.Items {
// 		item.FactionRestriction = factionRestrictions[item.Id]
// 	}
// }

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

	if item.Quality < proto.ItemQuality_ItemQualityUncommon {
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
		// TODO: Druid Tank Sim
		// {Name: "druid tank", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
		// 	Class:     proto.Class_ClassDruid,
		// 	Level:     60,
		// 	Equipment: &proto.EquipmentSpec{},
		// }, &proto.Player_FeralTankDruid{FeralTankDruid: &proto.FeralTankDruid{Options: &proto.FeralTankDruid_Options{}}}), nil, nil, nil)},
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
			Class:     proto.Class_ClassMage,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
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
		}, &proto.Player_Rogue{Rogue: &proto.Rogue{}}), nil, nil, nil)},
		// TODO: Rogue Tank Sim
		// {Name: "tank rogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
		// 	Class:     proto.Class_ClassRogue,
		// 	Level:     60,
		// 	Equipment: &proto.EquipmentSpec{},
		// 	Rotation:  &proto.APLRotation{},
		// }, &proto.Player_TankRogue{TankRogue: &proto.TankRogue{}}), nil, nil, nil)},
		{Name: "warrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassWarrior,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_Warrior{Warrior: &proto.Warrior{Options: &proto.Warrior_Options{}}}), nil, nil, nil)},
		// TODO: Warrior Tank Sim
		// {Name: "warrior tank", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
		// 	Class:     proto.Class_ClassWarrior,
		// 	Level:     60,
		// 	Equipment: &proto.EquipmentSpec{},
		// }, &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{Options: &proto.ProtectionWarrior_Options{}}}), nil, nil, nil)},
		{Name: "ret", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:     proto.Class_ClassPaladin,
			Level:     60,
			Equipment: &proto.EquipmentSpec{},
		}, &proto.Player_RetributionPaladin{RetributionPaladin: &proto.RetributionPaladin{Options: &proto.RetributionPaladin_Options{}}}), nil, nil, nil)},
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
	ret_db["warlock wowhead"] = []int32{18932, 18938, 25311, 20757, 17728, 603, 11722, 11735, 11695, 11668, 25309, 18540, 11661, 25307, 403851, 403852, 28610, 23161, 18881, 11730, 11713, 17926, 11678, 17923, 11726, 18871, 17953, 11717, 17937, 6215, 11689, 17924, 18931, 11672, 11700, 403689, 11704, 11684, 17928, 11708, 11675, 11694, 11660, 403848, 11740, 18937, 20756, 11719, 17925, 11734, 11743, 11667, 1122, 17922, 18930, 18870, 18880, 18647, 17727, 11712, 6353, 17952, 11729, 11721, 11699, 403688, 11688, 11677, 11671, 17862, 11703, 11693, 11659, 403844, 11725, 7659, 11707, 6789, 11683, 17921, 11739, 18869, 20755, 11733, 5484, 11665, 5784, 18879, 11711, 2970, 7651, 403687, 8289, 17951, 2362, 3700, 11687, 7641, 403843, 7648, 5699, 6226, 6219, 17920, 18868, 1490, 7646, 6213, 6229, 20752, 1086, 709, 403686, 1949, 2941, 1098, 691, 710, 6366, 6217, 7658, 3699, 1106, 403842, 1714, 132, 1456, 17919, 18867, 6223, 5138, 8288, 5500, 6202, 6205, 699, 403685, 126, 706, 3698, 1094, 5740, 698, 1088, 403841, 713, 712, 693, 1014, 5676, 1455, 5697, 6222, 704, 689, 403677, 1108, 755, 705, 403840, 6201, 696, 1120, 707, 697, 980, 5782, 1454, 695, 403839, 172, 702, 5019, 403629, 412788, 687, 412789, 426445, 426443, 426452, 416017, 416014, 426467, 431758, 431756, 431747, 416009, 403937, 403932, 403920, 403925, 403919, 403938, 403936, 431745, 431743, 425477, 425476, 416008, 416015, 426301, 403501, 403506, 348, 412758, 426241, 426245, 426246, 426247, 426331, 403828, 403789, 412783, 412784, 437169, 686, 403835, 426320, 426325, 688, 1180, 201, 227, 5009, 18119, 18120, 18121, 18122, 18123, 18288, 17788, 17789, 17790, 17791, 17792, 17778, 17779, 17780, 17781, 17782, 17962, 18223, 18220, 18697, 18698, 18699, 18700, 18701, 18788, 17917, 17918, 18130, 18131, 18132, 18133, 18134, 17954, 17955, 17956, 17957, 17958, 17783, 17784, 17785, 17786, 17787, 18708, 18731, 18743, 18744, 18745, 18746, 18751, 18752, 18218, 18219, 17810, 17811, 17812, 17813, 17814, 18827, 18829, 18830, 18310, 18311, 18312, 18313, 18179, 18180, 18181, 17804, 17805, 17806, 17807, 17808, 17864, 18393, 18213, 18372, 18126, 18127, 18767, 18768, 18703, 18704, 18692, 18693, 17815, 17833, 17834, 17835, 17836, 18694, 18695, 18696, 18182, 18183, 18754, 18755, 18756, 17927, 17929, 17930, 17931, 17932, 17793, 17796, 17801, 17802, 17803, 18774, 18775, 18821, 18823, 18824, 18825, 18705, 18706, 18707, 18135, 18136, 23785, 23822, 23823, 23824, 23825, 18709, 18710, 18094, 18095, 18073, 18096, 17959, 18271, 18272, 18273, 18274, 18275, 17877, 18265, 19028, 18174, 18175, 18176, 18177, 18178, 18769, 18770, 18771, 18772, 18773, 412798, 425463, 437032}

	return ret_db
}
