package database

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/tailscale/hujson"
	"github.com/wowsims/sod/sim/core/proto"
)

// Example db input file: https://nether.wowhead.com/classic/data/gear-planner?dv=100

func ParseWowheadDB(dbContents string) WowheadDatabase {
	var wowheadDB WowheadDatabase

	// Each part looks like 'WH.setPageData("wow.gearPlanner.some.name", {......});'
	parts := strings.Split(dbContents, "WH.setPageData(")

	for _, dbPart := range parts {
		// fmt.Printf("Part len: %d\n", len(dbPart))
		if len(dbPart) < 10 {
			continue
		}
		dbPart = strings.TrimSpace(dbPart)
		dbPart = strings.TrimRight(dbPart, ");")

		if dbPart[0] != '"' {
			continue
		}
		secondQuoteIdx := strings.Index(dbPart[1:], "\"")
		if secondQuoteIdx == -1 {
			continue
		}
		dbName := dbPart[1 : secondQuoteIdx+1]
		// fmt.Printf("DB name: %s\n", dbName)

		commaIdx := strings.Index(dbPart, ",")
		dbContents := dbPart[commaIdx+1:]
		if dbName == "wow.gearPlanner.classic.item" {
			standardized, err := hujson.Standardize([]byte(dbContents)) // Removes invalid JSON, such as trailing commas
			if err != nil {
				log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, dbContents[0:30], dbContents[len(dbContents)-30:])
			}

			err = json.Unmarshal(standardized, &wowheadDB.Items)
			if err != nil {
				log.Fatalf("failed to parse wowhead item db to json %s\n\n%s", err, dbContents[0:30])
			}
		}

		if dbName == "wow.gearPlanner.classic.randomEnchant" {
			standardized, err := hujson.Standardize([]byte(dbContents)) // Removes invalid JSON, such as trailing commas
			if err != nil {
				log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, dbContents[0:30], dbContents[len(dbContents)-30:])
			}

			err = json.Unmarshal(standardized, &wowheadDB.RandomSuffixes)
			if err != nil {
				log.Fatalf("failed to parse wowhead random suffix db to json %s\n\n%s", err, dbContents[0:30])
			}
		}
	}

	fmt.Printf("\n--\nWowhead DB items loaded: %d\n--\n", len(wowheadDB.Items))
	fmt.Printf("\n--\nWowhead DB random suffixes loaded: %d\n--\n", len(wowheadDB.RandomSuffixes))

	return wowheadDB
}

type WowheadDatabase struct {
	Items          map[string]WowheadItem
	RandomSuffixes map[string]WowheadRandomSuffix
}

type WowheadRandomSuffix struct {
	ID    int32                    `json:"id"`
	Name  string                   `json:"name"`
	Stats WowheadRandomSuffixStats `json:"stats"`
}

type WowheadRandomSuffixStats struct {
	Armor             int32 `json:"armor"`
	Strength          int32 `json:"str"`
	Agility           int32 `json:"agi"`
	Stamina           int32 `json:"sta"`
	Intellect         int32 `json:"int"`
	Spirit            int32 `json:"spi"`
	SpellPower        int32 `json:"spldmg"`
	ArcanePower       int32 `json:"arcsplpwr"`
	FirePower         int32 `json:"firsplpwr"`
	FrostPower        int32 `json:"frosplpwr"`
	HolyPower         int32 `json:"holsplpwr"`
	NaturePower       int32 `json:"natsplpwr"`
	ShadowPower       int32 `json:"shasplpwr"`
	MeleeCrit         int32 `json:"mlecritstrkpct"`
	MP5               int32 `json:"manargn"`
	AttackPower       int32 `json:"mleatkpwr"`
	RangedAttackPower int32 `json:"rgdatkpwr"`
	Defense           int32 `json:"def"`
	Block             int32 `json:"blockpct"`
	Dodge             int32 `json:"dodgepct"`
	ArcaneResistance  int32 `json:"arcres"`
	FireResistance    int32 `json:"firres"`
	FrostResistance   int32 `json:"frores"`
	NatureResistance  int32 `json:"natres"`
	ShadowResistance  int32 `json:"shares"`
	Healing           int32 `json:"splheal"`
}

func (wrs WowheadRandomSuffix) ToProto() *proto.ItemRandomSuffix {
	stats := Stats{
		proto.Stat_StatArmor:             float64(wrs.Stats.Armor),
		proto.Stat_StatStrength:          float64(wrs.Stats.Strength),
		proto.Stat_StatAgility:           float64(wrs.Stats.Agility),
		proto.Stat_StatStamina:           float64(wrs.Stats.Stamina),
		proto.Stat_StatIntellect:         float64(wrs.Stats.Intellect),
		proto.Stat_StatSpirit:            float64(wrs.Stats.Spirit),
		proto.Stat_StatSpellPower:        float64(wrs.Stats.SpellPower),
		proto.Stat_StatArcanePower:       float64(wrs.Stats.ArcanePower),
		proto.Stat_StatFirePower:         float64(wrs.Stats.FirePower),
		proto.Stat_StatFrostPower:        float64(wrs.Stats.FrostPower),
		proto.Stat_StatHolyPower:         float64(wrs.Stats.HolyPower),
		proto.Stat_StatNaturePower:       float64(wrs.Stats.NaturePower),
		proto.Stat_StatShadowPower:       float64(wrs.Stats.ShadowPower),
		proto.Stat_StatMeleeCrit:         float64(wrs.Stats.MeleeCrit),
		proto.Stat_StatMP5:               float64(wrs.Stats.MP5),
		proto.Stat_StatAttackPower:       float64(wrs.Stats.AttackPower),
		proto.Stat_StatRangedAttackPower: float64(wrs.Stats.RangedAttackPower),
		proto.Stat_StatDefense:           float64(wrs.Stats.Defense),
		proto.Stat_StatBlock:             float64(wrs.Stats.Block),
		proto.Stat_StatDodge:             float64(wrs.Stats.Dodge),
		proto.Stat_StatArcaneResistance:  float64(wrs.Stats.ArcaneResistance),
		proto.Stat_StatFireResistance:    float64(wrs.Stats.FireResistance),
		proto.Stat_StatFrostResistance:   float64(wrs.Stats.FrostResistance),
		proto.Stat_StatNatureResistance:  float64(wrs.Stats.NatureResistance),
		proto.Stat_StatShadowResistance:  float64(wrs.Stats.ShadowResistance),
		proto.Stat_StatHealingPower:      float64(wrs.Stats.Healing),
	}

	return &proto.ItemRandomSuffix{
		Id:    wrs.ID,
		Name:  wrs.Name,
		Stats: toSlice(stats),
	}
}

type WowheadItem struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Version int32  `json:"versionNum"`

	Quality       int32  `json:"quality"`
	Ilvl          int32  `json:"itemLevel"`
	Phase         int32  `json:"contentPhase"`
	RequiresLevel int32  `json:"requiredLevel"`
	RaceMask      uint16 `json:"raceMask"`
	ClassMask     uint16 `json:"classMask"`

	Stats               WowheadItemStats `json:"stats"`
	RandomSuffixOptions []int32          `json:"randomEnchants"`

	SourceTypes   []int32             `json:"source"` // 1 = Crafted, 2 = Dropped by, 3 = sold by zone vendor? barely used, 4 = Quest, 5 = Sold by
	SourceDetails []WowheadItemSource `json:"sourcemore"`
}
type WowheadItemStats struct {
	Armor int32 `json:"armor"`
}
type WowheadItemSource struct {
	C        int32  `json:"c"`
	Name     string `json:"n"`    // Name of crafting spell
	Icon     string `json:"icon"` // Icon corresponding to the named entity
	SkillID  int32  `json:"s"`    // Skill ID
	EntityID int32  `json:"ti"`   // Crafting Spell ID / NPC ID / ?? / Quest ID
	ZoneID   int32  `json:"z"`    // Only for drop / sold by sources
}

func (wi WowheadItem) ToProto() *proto.UIItem {
	var sources []*proto.UIItemSource
	for i, details := range wi.SourceDetails {
		switch wi.SourceTypes[i] {
		case 1: // Crafted
			sources = append(sources, &proto.UIItemSource{
				Source: &proto.UIItemSource_Crafted{
					Crafted: &proto.CraftedSource{
						Profession: WowheadProfessionIDs[details.SkillID],
						SpellId:    details.EntityID,
					},
				},
			})
		case 2: // Dropped by
			sources = append(sources, &proto.UIItemSource{
				Source: &proto.UIItemSource_Drop{
					Drop: &proto.DropSource{
						NpcId:  details.EntityID,
						ZoneId: details.ZoneID,
					},
				},
			})
		case 3: // Sold by zone vendor? barely used
		case 4: // Quest
			if details.EntityID != 0 {
				sources = append(sources, &proto.UIItemSource{
					Source: &proto.UIItemSource_Quest{
						Quest: &proto.QuestSource{
							Id:   details.EntityID,
							Name: details.Name,
						},
					},
				})
			}
		case 5: // Sold by
			sources = append(sources, &proto.UIItemSource{
				Source: &proto.UIItemSource_SoldBy{
					SoldBy: &proto.SoldBySource{
						NpcId:   details.EntityID,
						NpcName: details.Name,
						ZoneId:  details.ZoneID,
					},
				},
			})
		}
	}

	return &proto.UIItem{
		Id:                  wi.ID,
		Name:                wi.Name,
		Icon:                wi.Icon,
		Ilvl:                wi.Ilvl,
		Phase:               wi.getPhase(),
		RequiresLevel:       wi.RequiresLevel,
		FactionRestriction:  wi.getFactionRstriction(),
		ClassAllowlist:      wi.getClassRestriction(),
		Sources:             sources,
		RandomSuffixOptions: wi.RandomSuffixOptions,
	}
}

var SoDVersionRegex = regexp.MustCompile(`115[0-9]+`)

// Get the SoD phase corresponding to the item's version number
// 11500 (1.15.0) = phase 1
// 11501 (1.15.1) = phase 2
// 11502 (1.15.2) = phase 3
// 11503 (1.15.3) = phase 4
// etc.
// Anything else we'll fall back to phase 1
func (wi WowheadItem) getPhase() int32 {
	versionNumStr := strconv.Itoa(int(wi.Version))
	if SoDVersionRegex.MatchString(versionNumStr) && wi.Phase != 0 {
		return wi.Phase
	}

	if wi.Version >= 11500 && wi.Version < 11600 {
		return wi.Version - 11500 + 1
	}

	return 1
}

func (wi WowheadItem) getFactionRstriction() proto.UIItem_FactionRestriction {
	if wi.RaceMask == 77 {
		return proto.UIItem_FACTION_RESTRICTION_ALLIANCE_ONLY
	} else if wi.RaceMask == 178 {
		return proto.UIItem_FACTION_RESTRICTION_HORDE_ONLY
	} else {
		return proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED
	}
}

type ClassMask uint16

const (
	ClassMaskWarrior     ClassMask = 1 << iota
	ClassMaskPaladin               // 2
	ClassMaskHunter                // 4
	ClassMaskRogue                 // 8
	ClassMaskPriest                // 16
	ClassMaskDeathKnight           // 32
	ClassMaskShaman                // 64
	ClassMaskMage                  // 128
	ClassMaskWarlock               // 256
	ClassMaskUnknown               // 512 seemingly unused?
	ClassMaskDruid                 // 1024
)

func (wi WowheadItem) getClassRestriction() []proto.Class {
	classAllowlist := []proto.Class{}
	if wi.ClassMask&uint16(ClassMaskWarrior) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassWarrior)
	}
	if wi.ClassMask&uint16(ClassMaskPaladin) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassPaladin)
	}
	if wi.ClassMask&uint16(ClassMaskHunter) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassHunter)
	}
	if wi.ClassMask&uint16(ClassMaskRogue) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassRogue)
	}
	if wi.ClassMask&uint16(ClassMaskPriest) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassPriest)
	}
	if wi.ClassMask&uint16(ClassMaskDruid) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassDruid)
	}
	if wi.ClassMask&uint16(ClassMaskShaman) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassShaman)
	}
	if wi.ClassMask&uint16(ClassMaskMage) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassMage)
	}
	if wi.ClassMask&uint16(ClassMaskWarlock) != 0 {
		classAllowlist = append(classAllowlist, proto.Class_ClassWarlock)
	}

	return classAllowlist
}

var WowheadProfessionIDs = map[int32]proto.Profession{
	//129: proto.Profession_FirstAid,
	164: proto.Profession_Blacksmithing,
	165: proto.Profession_Leatherworking,
	171: proto.Profession_Alchemy,
	//185: proto.Profession_Cooking,
	186: proto.Profession_Mining,
	197: proto.Profession_Tailoring,
	202: proto.Profession_Engineering,
	333: proto.Profession_Enchanting,
}
