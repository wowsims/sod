package database

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type WowheadTooltipManager struct {
	TooltipManager
}

func (wtm *WowheadTooltipManager) Read() map[int32]WowheadItemResponse {
	strDB := wtm.TooltipManager.Read()
	return core.MapMap(strDB, func(id int32, tooltip string) (int32, WowheadItemResponse) {
		return id, NewWowheadItemResponse(id, tooltip)
	})
}

func NewWowheadItemTooltipManager(filePath string) *WowheadTooltipManager {
	return &WowheadTooltipManager{
		TooltipManager{
			FilePath:   filePath,
			UrlPattern: "https://nether.wowhead.com/classic/tooltip/item/%s?lvl=60",
		},
	}
}

func NewWowheadSpellTooltipManager(filePath string) *WowheadTooltipManager {
	return &WowheadTooltipManager{
		TooltipManager{
			FilePath:   filePath,
			UrlPattern: "https://nether.wowhead.com/classic/tooltip/spell/%s",
		},
	}
}

type Stats [stats.Len]float64
type WeaponSkills [stats.WeaponSkillLen]float64

type ItemResponse interface {
	GetName() string
	GetQuality() int
	GetIcon() string
	HasBuff() bool
	TooltipWithoutSetBonus() string
	GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string
	GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int
	GetIntValue(pattern *regexp.Regexp) int
	GetStats() Stats
	GetWeaponSkills() []float64
	IsEquippable() bool
	GetItemLevel() int
	GetRequiresLevel() int
	GetRequiredClass() proto.Class
	GetRequiredItemSlot() proto.ItemSlot
	GetSpellRank() int
	GetPhase() int
	GetUnique() bool
	GetItemType() proto.ItemType
	GetArmorType() proto.ArmorType
	GetWeaponType() proto.WeaponType
	GetHandType() proto.HandType
	GetRangedWeaponType() proto.RangedWeaponType
	GetWeaponDamage() (float64, float64)
	GetWeaponSpeed() float64
	GetItemSetID() string
	GetItemSetName() string
	IsHeroic() bool
	GetRequiredProfession() proto.Profession
}

type WowheadItemResponse struct {
	ID      int32
	Name    string `json:"name"`
	Quality int    `json:"quality"`
	Icon    string `json:"icon"`
	Tooltip string `json:"tooltip"`
	Buff    string `json:"buff"`
}

func NewWowheadItemResponse(id int32, tooltip string) WowheadItemResponse {
	response := WowheadItemResponse{}
	err := json.Unmarshal([]byte(tooltip), &response)
	if err != nil {
		fmt.Printf("Failed to decode tooltipBytes: %s\n", tooltip)
		panic(err)
	}
	response.ID = id
	return response
}

func (item WowheadItemResponse) GetName() string {
	return item.Name
}
func (item WowheadItemResponse) GetQuality() int {
	return item.Quality
}
func (item WowheadItemResponse) GetIcon() string {
	return item.Icon
}
func (item WowheadItemResponse) HasBuff() bool {
	return item.Buff != ""
}

func GetRegexStringValue(srcStr string, pattern *regexp.Regexp, matchIdx int) string {
	match := pattern.FindStringSubmatch(srcStr)
	if match == nil {
		return ""
	} else {
		return match[matchIdx]
	}
}
func GetRegexIntValue(srcStr string, pattern *regexp.Regexp, matchIdx int) int {
	matchStr := GetRegexStringValue(srcStr, pattern, matchIdx)

	val, err := strconv.Atoi(matchStr)
	if err != nil {
		return 0
	}

	return val
}
func GetBestRegexIntValue(srcStr string, patterns []*regexp.Regexp, matchIdx int) int {
	best := 0
	for _, pattern := range patterns {
		newVal := GetRegexIntValue(srcStr, pattern, matchIdx)
		if newVal > best {
			best = newVal
		}
	}
	return best
}

func (item WowheadItemResponse) TooltipWithoutSetBonus() string {
	setIdx := strings.Index(item.Tooltip, "Set : ")
	if setIdx == -1 {
		return item.Tooltip
	} else {
		return item.Tooltip[:setIdx]
	}
}

func (item WowheadItemResponse) GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string {
	return GetRegexStringValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WowheadItemResponse) GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int {
	return GetRegexIntValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WowheadItemResponse) GetIntValue(pattern *regexp.Regexp) int {
	return item.GetTooltipRegexValue(pattern, 1)
}

var armorRegex = regexp.MustCompile(`<!--amr-->([0-9]+) Armor`)
var agilityRegex = regexp.MustCompile(`<!--stat3-->([+-][0-9]+) Agility`)
var strengthRegex = regexp.MustCompile(`<!--stat4-->([+-][0-9]+) Strength`)
var intellectRegex = regexp.MustCompile(`<!--stat5-->([+-][0-9]+) Intellect`)
var spiritRegex = regexp.MustCompile(`<!--stat6-->([+-][0-9]+) Spirit`)
var staminaRegex = regexp.MustCompile(`<!--stat7-->([+-][0-9]+) Stamina`)
var spellHealingRegex = regexp.MustCompile(`Increases healing done by spells and effects by up to ([0-9]+)\.`)
var spellPowerRegex = regexp.MustCompile(`Increases damage and healing done by magical spells and effects by up to ([0-9]+)\.`)
var spellPowerRegex2 = regexp.MustCompile(`Increases damage and healing done by magical spells and effects by up to <!--rtg45-->([0-9]+)\.`)
var spellPowerRegex3 = regexp.MustCompile(`Increases healing done by up to [0-9]+ and damage done by up to ([0-9]+) for all magical spells and effects\.`)

// Patch 1.15.4 added the first item with a comgined multi-school power equip effect https://www.wowhead.com/classic/item=231785/echoes-of-betrayal
// so these have to check for other schools before and after the one we're looking for.
var arcaneSpellPowerRegex = regexp.MustCompile(`Increases damage done by (?:\w*? and )*?Arcane(?: and \w*?)* spells and effects by up to ([0-9]+)\.`)
var fireSpellPowerRegex = regexp.MustCompile(`Increases damage done by (?:\w*? and )*?Fire(?: and \w*?)* spells and effects by up to ([0-9]+)\.`)
var frostSpellPowerRegex = regexp.MustCompile(`Increases damage done by (?:\w*? and )*?Frost(?: and \w*?)* spells and effects by up to ([0-9]+)\.`)
var holySpellPowerRegex = regexp.MustCompile(`Increases damage done by (?:\w*? and )*?Holy(?: and \w*?)* spells and effects by up to ([0-9]+)\.`)
var natureSpellPowerRegex = regexp.MustCompile(`Increases damage done by (?:\w*? and )*?Nature(?: and \w*?)* spells and effects by up to ([0-9]+)\.`)
var shadowSpellPowerRegex = regexp.MustCompile(`Increases damage done by (?:\w*? and )*?Shadow(?: and \w*?)* spells and effects by up to ([0-9]+)\.`)

var hitRegex = regexp.MustCompile(`Improves your chance to hit with spells and with melee and ranged attacks by ([0-9]+)%\.`)
var hitRegex2 = regexp.MustCompile(`Improves your chance to hit with all spells and attacks by ([0-9]+)%\.`)
var physicalHitRegex = regexp.MustCompile(`Improves your chance to hit by ([0-9]+)%\.`)
var spellHitRegex = regexp.MustCompile(`Improves your chance to hit with spells by ([0-9]+)%\.`)
var critRegex = regexp.MustCompile(`Improves your chance to get a critical strike with melee and ranged attacks and with spells by ([0-9]+)%\.`)
var critRegex2 = regexp.MustCompile(`Improves your chance to get a critical strike with all spells and attacks by ([0-9]+)%\.`)
var spellCritRegex = regexp.MustCompile(`Improves your chance to get a critical strike with spells by ([0-9]+)%\.`)
var meleeCritRegex = regexp.MustCompile(`Improves your chance to get a critical strike by ([0-9]+)%\.`)
var hasteRegex = regexp.MustCompile(`Improves your haste by ([0-9]+)%\.`)

var spellPenetrationRegex = regexp.MustCompile(`Decreases the magical resistances of your spell targets by ([0-9]+)\.`)
var mp5Regex = regexp.MustCompile(`Restores ([0-9]+) mana per 5 sec\.`)
var attackPowerRegex = regexp.MustCompile(`\+([0-9]+) Attack Power\.`)
var attackPowerRegex2 = regexp.MustCompile(`\+<!--rtg38-->([0-9]+) Attack Power\.`)

var rangedAttackPowerRegex = regexp.MustCompile(`\+([0-9]+) ranged Attack Power\.`)
var rangedAttackPowerRegex2 = regexp.MustCompile(`Increases ranged attack power by ([0-9]+)\.`)
var rangedAttackPowerRegex3 = regexp.MustCompile(`Increases ranged attack power by <!--rtg39-->([0-9]+)\.`)

var feralAttackPowerRegex = regexp.MustCompile(`\+([0-9]+) Attack Power in Cat, Bear, and Dire Bear forms only\.`)

var armorPenetrationRegex = regexp.MustCompile(`Increases armor penetration rating by ([0-9]+)`)
var armorPenetrationRegex2 = regexp.MustCompile(`Increases your armor penetration by <!--rtg44-->([0-9]+)\.`)

var expertiseRegex = regexp.MustCompile(`Reduces the chance for your attacks to be dodged or parried by ([0-9]+)%\.`)
var timewornRegex = regexp.MustCompile(`<span style=\"color: #CD7F32\">Timeworn<\/span>`)
var weaponDamageRegex = regexp.MustCompile(`<!--dmg-->([0-9]+) - ([0-9]+)`)
var weaponDamageRegex2 = regexp.MustCompile(`<!--dmg-->([0-9]+) Damage`)
var weaponDamageBonusSchoolRegex = regexp.MustCompile(`\+ ([0-9]+) - ([0-9]+) [a-zA-Z]+ Damage`)
var weaponSpeedRegex = regexp.MustCompile(`<!--spd-->(([0-9]+).([0-9]+))`)

var axesSkill = regexp.MustCompile(`Increased Axes \+([0-9]+)\.`)
var swordsSkill = regexp.MustCompile(`Increased Swords \+([0-9]+)\.`)
var daggersSkill = regexp.MustCompile(`Increased Daggers \+([0-9]+)\.`)
var unarmedSkill = regexp.MustCompile(`Increased Fist Weapons \+([0-9]+)\.`)
var macesSkill = regexp.MustCompile(`Increased Maces \+([0-9]+)\.`)

var twoHandedAxesSkill = regexp.MustCompile(`Increased Two-handed Axes \+([0-9]+)\.`)
var twoHandedSwordsSkill = regexp.MustCompile(`Increased Two-handed Swords \+([0-9]+)\.`)
var twoHandedMacesSkill = regexp.MustCompile(`Increased Two-handed Maces \+([0-9]+)\.`)
var stavesSkill = regexp.MustCompile(`Increased Staves \+([0-9]+)\.`)
var polearmsSkill = regexp.MustCompile(`Increased Polearms \+([0-9]+)\.`)

var thrownSkill = regexp.MustCompile(`Increased Thrown \+([0-9]+)\.`)
var bowsSkill = regexp.MustCompile(`Increased Bows \+([0-9]+)\.`)
var crossbowsSkill = regexp.MustCompile(`Increased Crossbows \+([0-9]+)\.`)
var gunsSkill = regexp.MustCompile(`Increased Guns \+([0-9]+)\.`)

var feralCombatSkill = regexp.MustCompile(`Increased Feral Combat \+([0-9]+)\.`)

var defenseRegex = regexp.MustCompile(`Increased Defense \+([0-9]+)\.`)
var blockRegex = regexp.MustCompile(`Increases your chance to block attacks with a shield by ([0-9]+)%\.`)
var blockValueRegex = regexp.MustCompile(`Increases the block value of your shield by ([0-9]+)\.`)
var blockValueRegex2 = regexp.MustCompile(`<br>([0-9]+) Block<br>`)
var dodgeRegex = regexp.MustCompile(`Increases your chance to dodge an attack by ([0-9]+)%\.`)
var parryRegex = regexp.MustCompile(`Increases your chance to parry an attack by ([0-9]+)%\.`)
var resilienceRegex = regexp.MustCompile(`Improves your resilience rating by <!--rtg35-->([0-9]+)\.`)
var arcaneResistanceRegex = regexp.MustCompile(`\+([0-9]+) Arcane Resistance`)
var fireResistanceRegex = regexp.MustCompile(`\+([0-9]+) Fire Resistance`)
var frostResistanceRegex = regexp.MustCompile(`\+([0-9]+) Frost Resistance`)
var natureResistanceRegex = regexp.MustCompile(`\+([0-9]+) Nature Resistance`)
var shadowResistanceRegex = regexp.MustCompile(`\+([0-9]+) Shadow Resistance`)
var bonusArmorRegex = regexp.MustCompile(`Has ([0-9]+) bonus armor`)

// Match "Requires <a href=\"/...\" class="...">Profession</a> (level)"
var alchemyRegex = regexp.MustCompile(`Requires <a [ =/\\"\w]*>Alchemy<\/a> \([0-9]+\)`)
var blacksmithingRegex = regexp.MustCompile(`Requires <a [ =/\\"\w]*>Blacksmithing<\/a> \([0-9]+\)`)
var enchantingRegex = regexp.MustCompile(`Requires <a [ =/\\"\w]*>Enchanting<\/a> \([0-9]+\)`)
var engineeringRegex = regexp.MustCompile(`Requires <a [ =/\\"\w]*>Engineering<\/a> \([0-9]+\)`)
var leatherworkingRegex = regexp.MustCompile(`Requires <a [ =/\\"\w]*>Leatherworking<\/a> \([0-9]+\)`)
var tailoringRegex = regexp.MustCompile(`Requires <a [ =/\\"\w]*>Tailoring<\/a> \([0-9]+\)`)

func (item WowheadItemResponse) GetStats() Stats {
	sp := float64(item.GetIntValue(spellPowerRegex)) + float64(item.GetIntValue(spellPowerRegex2)) + float64(item.GetIntValue(spellPowerRegex3))
	baseAP := float64(item.GetIntValue(attackPowerRegex)) + float64(item.GetIntValue(attackPowerRegex2))
	armor, bonusArmor := item.GetArmorValues()
	timeworn := 0
	if item.GetTooltipRegexString(timewornRegex, 0) != "" {
		timeworn = 1
	}
	return Stats{
		proto.Stat_StatArmor:             float64(armor),
		proto.Stat_StatBonusArmor:        float64(bonusArmor),
		proto.Stat_StatStrength:          float64(item.GetIntValue(strengthRegex)),
		proto.Stat_StatAgility:           float64(item.GetIntValue(agilityRegex)),
		proto.Stat_StatStamina:           float64(item.GetIntValue(staminaRegex)),
		proto.Stat_StatIntellect:         float64(item.GetIntValue(intellectRegex)),
		proto.Stat_StatSpirit:            float64(item.GetIntValue(spiritRegex)),
		proto.Stat_StatSpellPower:        sp,
		proto.Stat_StatArcanePower:       float64(item.GetIntValue(arcaneSpellPowerRegex)),
		proto.Stat_StatFirePower:         float64(item.GetIntValue(fireSpellPowerRegex)),
		proto.Stat_StatFrostPower:        float64(item.GetIntValue(frostSpellPowerRegex)),
		proto.Stat_StatHolyPower:         float64(item.GetIntValue(holySpellPowerRegex)),
		proto.Stat_StatNaturePower:       float64(item.GetIntValue(natureSpellPowerRegex)),
		proto.Stat_StatShadowPower:       float64(item.GetIntValue(shadowSpellPowerRegex)),
		proto.Stat_StatSpellHit:          float64(item.GetIntValue(hitRegex) + item.GetIntValue(hitRegex2) + item.GetIntValue(spellHitRegex)),
		proto.Stat_StatMeleeHit:          float64(item.GetIntValue(hitRegex) + item.GetIntValue(hitRegex2) + item.GetIntValue(physicalHitRegex)),
		proto.Stat_StatSpellCrit:         float64(item.GetIntValue(critRegex) + item.GetIntValue(critRegex2) + item.GetIntValue(spellCritRegex)),
		proto.Stat_StatMeleeCrit:         float64(item.GetIntValue(critRegex) + item.GetIntValue(critRegex2) + item.GetIntValue(meleeCritRegex)),
		+proto.Stat_StatSpellHaste:       float64(item.GetIntValue(hasteRegex)),
		+proto.Stat_StatMeleeHaste:       float64(item.GetIntValue(hasteRegex)),
		proto.Stat_StatSpellPenetration:  float64(item.GetIntValue(spellPenetrationRegex)),
		proto.Stat_StatMP5:               float64(item.GetIntValue(mp5Regex)),
		proto.Stat_StatAttackPower:       baseAP,
		proto.Stat_StatRangedAttackPower: baseAP + float64(item.GetIntValue(rangedAttackPowerRegex)) + float64(item.GetIntValue(rangedAttackPowerRegex2)) + float64(item.GetIntValue(rangedAttackPowerRegex3)),
		proto.Stat_StatArmorPenetration:  float64(item.GetIntValue(armorPenetrationRegex) + item.GetIntValue(armorPenetrationRegex2)),
		proto.Stat_StatExpertise:         float64(item.GetIntValue(expertiseRegex)),
		proto.Stat_StatTimeworn:          float64(timeworn),
		proto.Stat_StatDefense:           float64(item.GetIntValue(defenseRegex)),
		proto.Stat_StatBlock:             float64(item.GetIntValue(blockRegex)),
		proto.Stat_StatBlockValue:        float64(item.GetIntValue(blockValueRegex) + item.GetIntValue(blockValueRegex2)),
		proto.Stat_StatDodge:             float64(item.GetIntValue(dodgeRegex)),
		proto.Stat_StatParry:             float64(item.GetIntValue(parryRegex)),
		proto.Stat_StatResilience:        float64(item.GetIntValue(resilienceRegex)),
		proto.Stat_StatArcaneResistance:  float64(item.GetIntValue(arcaneResistanceRegex)),
		proto.Stat_StatFireResistance:    float64(item.GetIntValue(fireResistanceRegex)),
		proto.Stat_StatFrostResistance:   float64(item.GetIntValue(frostResistanceRegex)),
		proto.Stat_StatNatureResistance:  float64(item.GetIntValue(natureResistanceRegex)),
		proto.Stat_StatShadowResistance:  float64(item.GetIntValue(shadowResistanceRegex)),
		proto.Stat_StatHealingPower:      float64(item.GetIntValue(spellHealingRegex)),
		proto.Stat_StatFeralAttackPower:  float64(item.GetIntValue(feralAttackPowerRegex)),
	}
}

func (item WowheadItemResponse) GetWeaponSkills() WeaponSkills {
	return WeaponSkills{
		0,
		float64(item.GetIntValue(axesSkill)),
		float64(item.GetIntValue(swordsSkill)),
		float64(item.GetIntValue(macesSkill)),
		float64(item.GetIntValue(daggersSkill)),
		float64(item.GetIntValue(unarmedSkill)),
		float64(item.GetIntValue(twoHandedAxesSkill)),
		float64(item.GetIntValue(twoHandedSwordsSkill)),
		float64(item.GetIntValue(twoHandedMacesSkill)),
		float64(item.GetIntValue(polearmsSkill)),
		float64(item.GetIntValue(stavesSkill)),
		float64(item.GetIntValue(thrownSkill)),
		float64(item.GetIntValue(bowsSkill)),
		float64(item.GetIntValue(crossbowsSkill)),
		float64(item.GetIntValue(gunsSkill)),
		float64(item.GetIntValue(feralCombatSkill)),
	}
}

var reqLevelRegex = regexp.MustCompile(`<!--rlvl-->([0-9]+)`)      // Items level
var reqLevelRegex2 = regexp.MustCompile(`Requires level ([0-9]+)`) // Spells level

func (item WowheadItemResponse) GetRequiresLevel() int {
	level := item.GetIntValue(reqLevelRegex)

	if level == 0 {
		level = item.GetIntValue(reqLevelRegex2)
	}

	return level
}

var reqClassRegex = regexp.MustCompile(`Requires (Druid|Hunter|Mage|Paladin|Priest|Rogue|Shaman|Warlock|Warrior)`)

func (item WowheadItemResponse) GetRequiredClass() proto.Class {
	class := item.GetTooltipRegexString(reqClassRegex, 1)

	if class == "" {
		return proto.Class_ClassUnknown
	}

	className := "Class" + class
	return proto.Class(proto.Class_value[className])
}

var reqSlotRegex = regexp.MustCompile(`Requires (Back|Belt|Bracer|Chest|Cloak|Boots|Feet|Gloves|Hands|Head|Helm|Legs|Pants|Ring|Shoulder|Trinket|Waist|Wrist)`)

func (item WowheadItemResponse) GetRequiredItemSlot() proto.ItemType {
	slot := item.GetTooltipRegexString(reqSlotRegex, 1)

	switch slot {
	case "Belt":
		fallthrough
	case "Waist":
		return proto.ItemType_ItemTypeWaist
	case "Bracer":
		fallthrough
	case "Wrist":
		return proto.ItemType_ItemTypeWrist
	case "Chest":
		return proto.ItemType_ItemTypeChest
	case "Back":
		fallthrough
	case "Cloak":
		return proto.ItemType_ItemTypeBack
	case "Boots":
		fallthrough
	case "Feet":
		return proto.ItemType_ItemTypeFeet
	case "Gloves":
		fallthrough
	case "Hands":
		return proto.ItemType_ItemTypeHands
	case "Head":
		fallthrough
	case "Helm":
		return proto.ItemType_ItemTypeHead
	case "Legs":
		fallthrough
	case "Pants":
		return proto.ItemType_ItemTypeLegs
	case "Ring":
		return proto.ItemType_ItemTypeFinger
	case "Trinket":
		return proto.ItemType_ItemTypeTrinket
	case "Shoulder":
		return proto.ItemType_ItemTypeShoulder
	default:
		return proto.ItemType_ItemTypeChest
	}
}

var rankRegex = regexp.MustCompile(`Rank ([0-9]+)`) // Spell rank

func (item WowheadItemResponse) GetSpellRank() int {
	return item.GetIntValue(rankRegex)

}

var patternRegexes = []*regexp.Regexp{
	regexp.MustCompile(`Design:`),
	regexp.MustCompile(`Recipe:`),
	regexp.MustCompile(`Pattern:`),
	regexp.MustCompile(`Plans:`),
	regexp.MustCompile(`Schematic:`),
}

func (item WowheadItemResponse) IsPattern() bool {
	for _, pattern := range patternRegexes {
		if pattern.MatchString(item.Tooltip) {
			return true
		}
	}
	return false
}

var randomEnchantRegex = regexp.MustCompile(`Random enchantment`)

func (item WowheadItemResponse) IsRandomEnchant() bool {
	return randomEnchantRegex.MatchString(item.Tooltip)
}

func (item WowheadItemResponse) IsEquippable() bool {
	return item.GetItemType() != proto.ItemType_ItemTypeUnknown &&
		!item.IsPattern()
}

var itemLevelRegex = regexp.MustCompile(`Item Level <!--ilvl-->([0-9]+)<`)

func (item WowheadItemResponse) GetItemLevel() int {
	return item.GetIntValue(itemLevelRegex)
}

var phaseRegex = regexp.MustCompile(`Phase ([0-9])`)

func (item WowheadItemResponse) GetPhase() int {
	phase := item.GetIntValue(phaseRegex)
	if phase != 0 {
		return phase
	}

	// ilvl := item.GetItemLevel()
	// if ilvl <= 164 { // TBC items
	// 	return 0
	// }

	// if ilvl < 200 || ilvl == 200 || ilvl == 213 || ilvl == 226 {
	// 	return 1
	// } else if ilvl == 219 || ilvl == 226 || ilvl == 239 {
	// 	return 2
	// } else if ilvl == 232 || ilvl == 245 || ilvl == 258 {
	// 	return 3
	// } else if ilvl == 251 || ilvl == 258 || ilvl == 259 || ilvl == 264 || ilvl == 268 || ilvl == 270 || ilvl == 271 || ilvl == 272 {
	// 	return 4
	// } else if ilvl == 277 || ilvl == 284 {
	// 	return 5
	// }

	// default to 1
	return 1
}

var uniqueRegex = regexp.MustCompile(`Unique`)

func (item WowheadItemResponse) GetUnique() bool {
	return uniqueRegex.MatchString(item.Tooltip)
}

var itemTypePatterns = map[proto.ItemType]*regexp.Regexp{
	proto.ItemType_ItemTypeHead:     regexp.MustCompile(`<td>Head</td>`),
	proto.ItemType_ItemTypeNeck:     regexp.MustCompile(`<td>Neck</td>`),
	proto.ItemType_ItemTypeShoulder: regexp.MustCompile(`<td>Shoulder</td>`),
	proto.ItemType_ItemTypeBack:     regexp.MustCompile(`<td>Back</td>`),
	proto.ItemType_ItemTypeChest:    regexp.MustCompile(`<td>Chest</td>`),
	proto.ItemType_ItemTypeWrist:    regexp.MustCompile(`<td>Wrist</td>`),
	proto.ItemType_ItemTypeHands:    regexp.MustCompile(`<td>Hands</td>`),
	proto.ItemType_ItemTypeWaist:    regexp.MustCompile(`<td>Waist</td>`),
	proto.ItemType_ItemTypeLegs:     regexp.MustCompile(`<td>Legs</td>`),
	proto.ItemType_ItemTypeFeet:     regexp.MustCompile(`<td>Feet</td>`),
	proto.ItemType_ItemTypeFinger:   regexp.MustCompile(`<td>Finger</td>`),
	proto.ItemType_ItemTypeTrinket:  regexp.MustCompile(`<td>Trinket</td>`),
	proto.ItemType_ItemTypeWeapon:   regexp.MustCompile(`<td>((Main Hand)|(Two-Hand)|(One-Hand)|(Off Hand)|(Held In Off-hand)|(Held In Off-Hand))</td>`),
	proto.ItemType_ItemTypeRanged:   regexp.MustCompile(`<td>(Ranged|Thrown|Relic)</td>`),
}

func (item WowheadItemResponse) GetItemType() proto.ItemType {
	for itemType, pattern := range itemTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return itemType
		}
	}
	return proto.ItemType_ItemTypeUnknown
}

func (item WowheadItemResponse) IsScalableArmorSlot() bool {
	// Special case shields as Base Armor
	if item.GetWeaponType() == proto.WeaponType_WeaponTypeShield {
		return true
	}

	itemType := item.GetItemType()
	switch itemType {
	case
		proto.ItemType_ItemTypeNeck,
		proto.ItemType_ItemTypeFinger,
		proto.ItemType_ItemTypeTrinket,
		proto.ItemType_ItemTypeWeapon:
		return false
	}
	return true
}

func (item WowheadItemResponse) GetArmorValues() (int, int) {
	armorValue := item.GetIntValue(armorRegex)
	bonusArmorValue := item.GetIntValue(bonusArmorRegex)

	if item.IsScalableArmorSlot() {
		armorValue = armorValue - bonusArmorValue
	} else {
		bonusArmorValue = armorValue
		armorValue = 0
	}

	return armorValue, bonusArmorValue
}

var armorTypePatterns = map[proto.ArmorType]*regexp.Regexp{
	proto.ArmorType_ArmorTypeCloth:   regexp.MustCompile(`<span class="q1">(?:<!--asc1-->)?Cloth</span>`),
	proto.ArmorType_ArmorTypeLeather: regexp.MustCompile(`<span class="q1">(?:<!--asc2-->)?Leather</span>`),
	proto.ArmorType_ArmorTypeMail:    regexp.MustCompile(`<span class="q1">(?:<!--asc3-->)?Mail</span>`),
	proto.ArmorType_ArmorTypePlate:   regexp.MustCompile(`<span class="q1">(?:<!--asc4-->)?Plate</span>`),
}

func (item WowheadItemResponse) GetArmorType() proto.ArmorType {
	for armorType, pattern := range armorTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return armorType
		}
	}
	return proto.ArmorType_ArmorTypeUnknown
}

var weaponTypePatterns = map[proto.WeaponType]*regexp.Regexp{
	proto.WeaponType_WeaponTypeAxe:     regexp.MustCompile(`<span class="q1">Axe</span>`),
	proto.WeaponType_WeaponTypeDagger:  regexp.MustCompile(`<span class="q1">Dagger</span>`),
	proto.WeaponType_WeaponTypeFist:    regexp.MustCompile(`<span class="q1">Fist Weapon</span>`),
	proto.WeaponType_WeaponTypeMace:    regexp.MustCompile(`<span class="q1">Mace</span>`),
	proto.WeaponType_WeaponTypeOffHand: regexp.MustCompile(`<td>Held In Off-hand</td>`),
	proto.WeaponType_WeaponTypePolearm: regexp.MustCompile(`<span class="q1">Polearm</span>`),
	proto.WeaponType_WeaponTypeShield:  regexp.MustCompile(`<span class="q1">Shield</span>`),
	proto.WeaponType_WeaponTypeStaff:   regexp.MustCompile(`<span class="q1">Staff</span>`),
	proto.WeaponType_WeaponTypeSword:   regexp.MustCompile(`<span class="q1">Sword</span>`),
}

func (item WowheadItemResponse) GetWeaponType() proto.WeaponType {
	for weaponType, pattern := range weaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return weaponType
		}
	}
	return proto.WeaponType_WeaponTypeUnknown
}

var handTypePatterns = map[proto.HandType]*regexp.Regexp{
	proto.HandType_HandTypeMainHand: regexp.MustCompile(`<td>Main Hand</td>`),
	proto.HandType_HandTypeOneHand:  regexp.MustCompile(`<td>One-Hand</td>`),
	proto.HandType_HandTypeOffHand:  regexp.MustCompile(`<td>((Off Hand)|(Held In Off-hand)|(Held In Off-Hand))</td>`),
	proto.HandType_HandTypeTwoHand:  regexp.MustCompile(`<td>Two-Hand</td>`),
}

func (item WowheadItemResponse) GetHandType() proto.HandType {
	for handType, pattern := range handTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return handType
		}
	}
	return proto.HandType_HandTypeUnknown
}

var rangedWeaponTypePatterns = map[proto.RangedWeaponType]*regexp.Regexp{
	proto.RangedWeaponType_RangedWeaponTypeBow:      regexp.MustCompile(`<span class="q1">Bow</span>`),
	proto.RangedWeaponType_RangedWeaponTypeCrossbow: regexp.MustCompile(`<span class="q1">Crossbow</span>`),
	proto.RangedWeaponType_RangedWeaponTypeGun:      regexp.MustCompile(`<span class="q1">Gun</span>`),
	proto.RangedWeaponType_RangedWeaponTypeIdol:     regexp.MustCompile(`<span class="q1">Idol</span>`),
	proto.RangedWeaponType_RangedWeaponTypeLibram:   regexp.MustCompile(`<span class="q1">Libram</span>`),
	proto.RangedWeaponType_RangedWeaponTypeThrown:   regexp.MustCompile(`<span class="q1">Thrown</span>`),
	proto.RangedWeaponType_RangedWeaponTypeTotem:    regexp.MustCompile(`<span class="q1">Totem</span>`),
	proto.RangedWeaponType_RangedWeaponTypeWand:     regexp.MustCompile(`<span class="q1">Wand</span>`),
	proto.RangedWeaponType_RangedWeaponTypeSigil:    regexp.MustCompile(`<span class="q1">Sigil</span>`),
}

func (item WowheadItemResponse) GetRangedWeaponType() proto.RangedWeaponType {
	for rangedWeaponType, pattern := range rangedWeaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return rangedWeaponType
		}
	}
	return proto.RangedWeaponType_RangedWeaponTypeUnknown
}

// Returns min/max of weapon damage
func (item WowheadItemResponse) GetWeaponDamage() (float64, float64) {
	noCommas := strings.ReplaceAll(item.Tooltip, ",", "")
	min, max := 0.0, 0.0

	// Base damage
	if matches := weaponDamageRegex.FindStringSubmatch(noCommas); len(matches) > 0 {
		minVal, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		maxVal, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		if minVal > maxVal {
			log.Fatalf("Invalid weapon damage for item %s: min = %0.1f, max = %0.1f", item.Name, min, max)
		}
		min, max = minVal, maxVal
	} else if matches := weaponDamageRegex2.FindStringSubmatch(noCommas); len(matches) > 0 {
		val, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		min, max = val, val
	}

	// Add effects like Thunderfury + 16 - 30 Nature Damage that were intended for bonus school damage on top of the base weapon damage
	if matches := weaponDamageBonusSchoolRegex.FindStringSubmatch(noCommas); len(matches) > 0 {
		minVal, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		maxVal, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		if minVal > maxVal {
			log.Fatalf("Invalid bonus school weapon damage for item %s: min = %0.1f, max = %0.1f", item.Name, min, max)
		}

		min += minVal
		max += maxVal
	}

	return min, max
}

func (item WowheadItemResponse) GetWeaponSpeed() float64 {
	if matches := weaponSpeedRegex.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
		speed, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return speed
	}
	return 0
}

func (item WowheadItemResponse) ToItemProto() *proto.UIItem {
	weaponDamageMin, weaponDamageMax := item.GetWeaponDamage()

	itemProto := &proto.UIItem{
		Id:   item.ID,
		Name: item.GetName(),
		Icon: item.GetIcon(),

		Type:             item.GetItemType(),
		ArmorType:        item.GetArmorType(),
		WeaponType:       item.GetWeaponType(),
		HandType:         item.GetHandType(),
		RangedWeaponType: item.GetRangedWeaponType(),

		Stats: toSlice(item.GetStats()),

		WeaponSkills: weaponSkillsToSlice(item.GetWeaponSkills()),

		WeaponDamageMin: weaponDamageMin,
		WeaponDamageMax: weaponDamageMax,
		WeaponSpeed:     item.GetWeaponSpeed(),

		Ilvl:          int32(item.GetItemLevel()),
		Phase:         int32(item.GetPhase()),
		RequiresLevel: int32(item.GetRequiresLevel()),
		Quality:       proto.ItemQuality(item.GetQuality()),
		Unique:        item.GetUnique(),
		Heroic:        item.IsHeroic(),

		RequiredProfession: item.GetRequiredProfession(),
		SetName:            item.GetItemSetName(),
	}

	if item.GetRequiredProfession() != proto.Profession_ProfessionUnknown {
		itemProto.Sources = append(itemProto.Sources, &proto.UIItemSource{
			Source: &proto.UIItemSource_Crafted{
				Crafted: &proto.CraftedSource{
					Profession: item.GetRequiredProfession(),
				},
			},
		})
	}

	return itemProto
}

var itemSetNameRegex = regexp.MustCompile(`<a href="/classic/item-set=-?([0-9]+)/(.*)" class="q">([^<]+)<`)

func (item WowheadItemResponse) GetItemSetID() int {
	idStr := item.GetTooltipRegexString(itemSetNameRegex, 1)
	id, _ := strconv.Atoi(idStr)
	return id

	// // Strip out the 10/25 man prefixes from set names
	// withoutTier := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(original, "Heroes' "), "Valorous "), "Conqueror's "), "Triumphant "), "Sanctified ")
	// if original != withoutTier { // if we found a tier prefix, return now.
	// 	return withoutTier
	// }

	// // Now strip out the season prefix from any pvp set names
	// withoutPvp := strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(original, "Savage Glad", "Glad", 1), "Hateful Glad", "Glad", 1), "Deadly Glad", "Glad", 1), "Furious Glad", "Glad", 1), "Relentless Glad", "Glad", 1), "Wrathful Glad", "Glad", 1)
	// return withoutPvp
}

func (item WowheadItemResponse) GetItemSetName() string {
	return item.GetTooltipRegexString(itemSetNameRegex, 3)

	// // Strip out the 10/25 man prefixes from set names
	// withoutTier := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(original, "Heroes' "), "Valorous "), "Conqueror's "), "Triumphant "), "Sanctified ")
	// if original != withoutTier { // if we found a tier prefix, return now.
	// 	return withoutTier
	// }

	// // Now strip out the season prefix from any pvp set names
	// withoutPvp := strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(original, "Savage Glad", "Glad", 1), "Hateful Glad", "Glad", 1), "Deadly Glad", "Glad", 1), "Furious Glad", "Glad", 1), "Relentless Glad", "Glad", 1), "Wrathful Glad", "Glad", 1)
	// return withoutPvp
}

func (item WowheadItemResponse) IsHeroic() bool {
	return strings.Contains(item.Tooltip, "<span class=\"q2\">Heroic</span>")
}

func (item WowheadItemResponse) GetRequiredProfession() proto.Profession {
	if alchemyRegex.MatchString(item.Tooltip) {
		return proto.Profession_Alchemy
	} else if blacksmithingRegex.MatchString(item.Tooltip) {
		return proto.Profession_Blacksmithing
	} else if enchantingRegex.MatchString(item.Tooltip) {
		return proto.Profession_Enchanting
	} else if engineeringRegex.MatchString(item.Tooltip) {
		return proto.Profession_Engineering
	} else if leatherworkingRegex.MatchString(item.Tooltip) {
		return proto.Profession_Leatherworking
	} else if tailoringRegex.MatchString(item.Tooltip) {
		return proto.Profession_Tailoring
	} else {
		return proto.Profession_ProfessionUnknown
	}
}
