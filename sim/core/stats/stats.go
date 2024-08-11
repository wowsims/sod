package stats

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
)

type Stats [Len]float64

type Stat byte

// Use internal representation instead of proto.Stat so we can add functions
// and use 'byte' as the data type.
//
// This needs to stay synced with proto.Stat.
const (
	Strength Stat = iota
	Agility
	Stamina
	Intellect
	Spirit
	SpellPower
	ArcanePower
	FirePower
	FrostPower
	HolyPower
	NaturePower
	ShadowPower
	MP5
	SpellHit
	SpellCrit
	SpellHaste
	SpellPenetration
	AttackPower
	MeleeHit
	MeleeCrit
	MeleeHaste
	ArmorPenetration
	Expertise
	Mana
	Energy
	Rage
	Armor
	RangedAttackPower
	Defense
	Block
	BlockValue
	Dodge
	Parry
	Resilience
	Health
	ArcaneResistance
	FireResistance
	FrostResistance
	NatureResistance
	ShadowResistance
	BonusArmor
	HealingPower
	SpellDamage
	FeralAttackPower

	// DO NOT add new stats here without discussing it first; new stats come with
	// a performance penalty.

	Len
)

type WeaponSkills [WeaponSkillLen]float64

type WeaponSkill byte

const (
	WeaponSkillUnknown WeaponSkill = iota
	WeaponSkillAxes
	WeaponSkillSwords
	WeaponSkillMaces
	WeaponSkillDaggers
	WeaponSkillUnarmed
	WeaponSkillTwoHandedAxes
	WeaponSkillTwoHandedSwords
	WeaponSkillTwoHandedMaces
	WeaponSkillPolearms
	WeaponSkillStaves
	WeaponSkillThrown
	WeaponSkillBows
	WeaponSkillCrossbows
	WeaponSkillGuns
	WeaponSkillFeralCombat

	WeaponSkillLen
)

var PseudoStatsLen = len(proto.PseudoStat_name)
var UnitStatsLen = int(Len) + PseudoStatsLen

type SchoolIndex byte

// If you add a new multi-school you also need to update
// core/spell_school.go accordingly!

const (
	SchoolIndexNone     SchoolIndex = 0
	SchoolIndexPhysical SchoolIndex = iota
	SchoolIndexArcane
	SchoolIndexFire
	SchoolIndexFrost
	SchoolIndexHoly
	SchoolIndexNature
	SchoolIndexShadow

	SchoolLen

	// School is composed of multiple base schools.
	SchoolIndexMultischool SchoolIndex = iota - 1 // This is deliberately set this way to be a continuous sequence.
)

// Check if school index is a multi-school.
func (schoolIndex SchoolIndex) IsMultiSchool() bool {
	return schoolIndex == SchoolIndexMultischool
}

type SchoolValueArrayValues interface{ int32 | float64 }
type SchoolValueArray[T SchoolValueArrayValues] [SchoolLen]T

func NewSchoolValueArray[T SchoolValueArrayValues](defaultVal T) SchoolValueArray[T] {
	d := defaultVal
	return SchoolValueArray[T]{
		d, d, d, d, d, d, d, d,
	}
}

func (sfa *SchoolValueArray[T]) AddToAllSchools(change T) {
	for idx := SchoolIndexPhysical; idx < SchoolLen; idx++ {
		sfa[idx] += change
	}
}

func (sfa *SchoolValueArray[T]) AddToMagicSchools(change T) {
	for idx := SchoolIndexArcane; idx < SchoolLen; idx++ {
		sfa[idx] += change
	}
}

func (sfa *SchoolValueArray[T]) MultiplyAllSchools(mult T) {
	for idx := SchoolIndexPhysical; idx < SchoolLen; idx++ {
		sfa[idx] *= mult
	}
}

func (sfa *SchoolValueArray[T]) MultiplyMagicSchools(mult T) {
	for idx := SchoolIndexArcane; idx < SchoolLen; idx++ {
		sfa[idx] *= mult
	}
}

// If you add a new defense type you also need to update
// core/constants.go accordingly!

type DefenseTypeIndex byte

const (
	DefenseTypeIndexNone  DefenseTypeIndex = 0
	DefenseTypeIndexMagic DefenseTypeIndex = iota
	DefenseTypeIndexMelee
	DefenseTypeIndexRanged

	DefenseTypeLen
)

func ProtoArrayToStatsList(protoStats []proto.Stat) []Stat {
	stats := make([]Stat, len(protoStats))
	for i, v := range protoStats {
		stats[i] = Stat(v)
	}
	return stats
}

func (s Stat) StatName() string {
	switch s {
	case Strength:
		return "Strength"
	case Agility:
		return "Agility"
	case Stamina:
		return "Stamina"
	case Intellect:
		return "Intellect"
	case Spirit:
		return "Spirit"
	case SpellCrit:
		return "SpellCrit"
	case SpellHit:
		return "SpellHit"
	case SpellPower:
		return "SpellPower"
	case ArcanePower:
		return "ArcanePower"
	case FirePower:
		return "FirePower"
	case FrostPower:
		return "FrostPower"
	case HolyPower:
		return "HolyPower"
	case NaturePower:
		return "NaturePower"
	case ShadowPower:
		return "ShadowPower"
	case SpellDamage:
		return "SpellDamage"
	case HealingPower:
		return "HealingPower"
	case SpellHaste:
		return "SpellHaste"
	case MP5:
		return "MP5"
	case SpellPenetration:
		return "SpellPenetration"
	case AttackPower:
		return "AttackPower"
	case MeleeHit:
		return "MeleeHit"
	case MeleeHaste:
		return "MeleeHaste"
	case MeleeCrit:
		return "MeleeCrit"
	case Expertise:
		return "Expertise"
	case ArmorPenetration:
		return "ArmorPenetration"
	case Mana:
		return "Mana"
	case Energy:
		return "Energy"
	case Rage:
		return "Rage"
	case Armor:
		return "Armor"
	case BonusArmor:
		return "BonusArmor"
	case RangedAttackPower:
		return "RangedAttackPower"
	case Defense:
		return "Defense"
	case Block:
		return "Block"
	case BlockValue:
		return "BlockValue"
	case Dodge:
		return "Dodge"
	case Parry:
		return "Parry"
	case Resilience:
		return "Resilience"
	case Health:
		return "Health"
	case FireResistance:
		return "FireResistance"
	case NatureResistance:
		return "NatureResistance"
	case FrostResistance:
		return "FrostResistance"
	case ShadowResistance:
		return "ShadowResistance"
	case ArcaneResistance:
		return "ArcaneResistance"
	case FeralAttackPower:
		return "FeralAttackPower"
	}

	return "none"
}

func FromFloatArray(values []float64) Stats {
	var stats Stats
	copy(stats[:], values)
	return stats
}

func WeaponSkillsFloatArray(values []float64) WeaponSkills {
	var stats WeaponSkills
	copy(stats[:], values)
	return stats
}

// Adds two Stats together, returning the new Stats.
func (stats Stats) Add(other Stats) Stats {
	for k := range stats {
		stats[k] += other[k]
	}
	return stats
}

// Adds another to Stats to this, in-place. For performance, only.
func (stats *Stats) AddInplace(other *Stats) {
	for k := range stats {
		stats[k] += other[k]
	}
}

// Subtracts another Stats from this one, returning the new Stats.
func (stats Stats) Subtract(other Stats) Stats {
	for k := range stats {
		stats[k] -= other[k]
	}
	return stats
}

func (stats Stats) Invert() Stats {
	for k, v := range stats {
		stats[k] = -v
	}
	return stats
}

func (stats Stats) Multiply(multiplier float64) Stats {
	for k := range stats {
		stats[k] *= multiplier
	}
	return stats
}

func (stats Stats) Floor() Stats {
	for k := range stats {
		stats[k] = math.Floor(stats[k])
	}
	return stats
}

// Multiplies two Stats together by multiplying the values of corresponding
// stats, like a dot product operation.
func (stats Stats) DotProduct(other Stats) Stats {
	for k := range stats {
		stats[k] *= other[k]
	}
	return stats
}

func (stats Stats) Equals(other Stats) bool {
	return stats == other
}

func (stats Stats) EqualsWithTolerance(other Stats, tolerance float64) bool {
	for k, v := range stats {
		if v < other[k]-tolerance || v > other[k]+tolerance {
			return false
		}
	}
	return true
}

func (stats Stats) String() string {
	var sb strings.Builder
	sb.WriteString("\n{\n")

	for statIdx, statValue := range stats {
		if statValue == 0 {
			continue
		}
		if name := Stat(statIdx).StatName(); name != "none" {
			_, _ = fmt.Fprintf(&sb, "\t%s: %0.3f,\n", name, statValue)
		}
	}

	sb.WriteString("\n}")
	return sb.String()
}

// Like String() but without the newlines.
func (stats Stats) FlatString() string {
	var sb strings.Builder
	sb.WriteString("{")

	for statIdx, statValue := range stats {
		if statValue == 0 {
			continue
		}
		if name := Stat(statIdx).StatName(); name != "none" {
			_, _ = fmt.Fprintf(&sb, "\"%s\": %0.3f,", name, statValue)
		}
	}

	sb.WriteString("}")
	return sb.String()
}

func (stats Stats) ToFloatArray() []float64 {
	return stats[:]
}

type PseudoStats struct {
	///////////////////////////////////////////////////
	// Effects that apply when this unit is the attacker.
	///////////////////////////////////////////////////

	SchoolCostMultiplier SchoolValueArray[int32] // Multipliers for spell costs stored as ints

	CastSpeedMultiplier   float64
	MeleeSpeedMultiplier  float64
	RangedSpeedMultiplier float64

	MeleeCritMultiplier float64

	FiveSecondRuleRefreshTime time.Duration // last time a spell was cast
	SpiritRegenRateCasting    float64       // percentage of spirit regen allowed during casting. Spell effect MOD_MANA_REGEN_INTERRUPT (134)

	// Both of these are currently only used for innervate.
	ForceFullSpiritRegen  bool    // If set, automatically uses full spirit regen regardless of FSR refresh time.
	SpiritRegenMultiplier float64 // Multiplier on spirit portion of mana regen. For spell effect MOD_POWER_REGEN_PERCENT (110)

	// If true, allows block/parry.
	InFrontOfTarget bool

	// "Apply Aura: Mod Damage Done (Physical)", applies to abilities with EffectSpellCoefficient > 0.
	//  This includes almost all "(Normalized) Weapon Damage", but also some "School Damage (Physical)" abilities.
	BonusDamage float64 // Comes from '+X Weapon Damage' effects

	BonusMHDps     float64
	BonusOHDps     float64
	BonusRangedDps float64

	DisableDWMissPenalty bool    // Used by Heroic Strike and Cleave
	IncreasedMissChance  float64 // Insect Swarm and Scorpid Sting
	DodgeReduction       float64 // Used by Warrior talent 'Weapon Mastery' and SWP boss auras.


	MobTypeAttackPower float64 // Bonus AP against mobs of the current type.
	MobTypeSpellPower  float64 // Bonus SP against mobs of the current type.

	ThreatMultiplier float64 // Modulates the threat generated. Affected by things like salv.

	DamageDealtMultiplier       float64                   // All damage
	SchoolDamageDealtMultiplier SchoolValueArray[float64] // For specific spell schools. DO NOT use with multi school idices! See helper functions on Unit!

	// Important when unit is attacker or target
	BlockValueMultiplier float64

	// Only used for NPCs, governs variance in enemy auto-attack damage
	DamageSpread float64

	// Weapon Skills
	UnarmedSkill         float64
	DaggersSkill         float64
	SwordsSkill          float64
	MacesSkill           float64
	AxesSkill            float64
	TwoHandedSwordsSkill float64
	TwoHandedMacesSkill  float64
	TwoHandedAxesSkill   float64
	PolearmsSkill        float64
	StavesSkill          float64

	// Ranged Skills
	BowsSkill      float64
	CrossbowsSkill float64
	GunsSkill      float64
	ThrownSkill    float64

	// Special Feral Weapon Skill
	FeralCombatEnabled bool
	FeralCombatSkill   float64

	///////////////////////////////////////////////////
	// Effects that apply when this unit is the target.
	///////////////////////////////////////////////////

	CanBlock bool
	CanParry bool
	Stunned  bool // prevents blocks, dodges, and parries

	ParryHaste bool

	ReducedCritTakenChance float64 // Reduces chance to be crit.

	BonusRangedAttackPowerTaken float64 // Hunters mark
	BonusMeleeHitRatingTaken    float64 // Formerly Imp FF and SW Radiance;
	BonusSpellHitRatingTaken    float64 // Imp FF

	BonusHealingTaken float64 // Talisman of Troll Divinity

	BonusDamageTakenBeforeModifiers [DefenseTypeLen]float64 // Flat damage reduction values BEFORE Modifiers like Blessing of Sanctuary
	BonusDamageTakenAfterModifiers  [DefenseTypeLen]float64 // Flat damage reduction values AFTER Modifiers like Stoneskin Totem, Windwall Totem, etc.

	DamageTakenMultiplier       float64                   // All damage
	SchoolDamageTakenMultiplier SchoolValueArray[float64] // For specific spell schools. DO NOT use with multi school index! See helper functions on Unit!
	SchoolCritTakenChance       SchoolValueArray[float64] // For spell school crit. DO NOT use with multi school index! See helper functions on Unit!
	SchoolBonusDamageTaken      SchoolValueArray[float64] // For spell school bonus damage taken. DO NOT use with multi school index! See helper functions on Unit!
	SchoolBonusHitChance        SchoolValueArray[float64] // Spell school-specific hit bonuses such as ring runes

	BleedDamageTakenMultiplier  float64 // Modifies damage taken from bleed effects
	PoisonDamageTakenMultiplier float64 // Modifies damage taken from poison effects

	ArmorMultiplier float64 // Major/minor/special multiplicative armor modifiers

	HealingTakenMultiplier float64
}

func NewPseudoStats() PseudoStats {
	return PseudoStats{
		SchoolCostMultiplier: NewSchoolValueArray(int32(100)),

		CastSpeedMultiplier:   1,
		MeleeSpeedMultiplier:  1,
		RangedSpeedMultiplier: 1,
		SpiritRegenMultiplier: 1,

		MeleeCritMultiplier: 1,

		ThreatMultiplier: 1,

		DamageDealtMultiplier:       1,
		SchoolDamageDealtMultiplier: NewSchoolValueArray(1.0),

		BlockValueMultiplier: 1,

		DamageSpread: 0.3333,

		// Target effects.
		DamageTakenMultiplier:       1,
		SchoolDamageTakenMultiplier: NewSchoolValueArray(1.0),
		SchoolCritTakenChance:       NewSchoolValueArray(0.0),
		SchoolBonusDamageTaken:      NewSchoolValueArray(0.0),

		BleedDamageTakenMultiplier:  1,
		PoisonDamageTakenMultiplier: 1,

		ArmorMultiplier: 1,

		HealingTakenMultiplier: 1,
		UnarmedSkill:           0,
		DaggersSkill:           0,
		SwordsSkill:            0,
		MacesSkill:             0,
		AxesSkill:              0,
		TwoHandedSwordsSkill:   0,
		TwoHandedMacesSkill:    0,
		TwoHandedAxesSkill:     0,
		PolearmsSkill:          0,
		StavesSkill:            0,

		BowsSkill:      0,
		GunsSkill:      0,
		CrossbowsSkill: 0,
		ThrownSkill:    0,

		FeralCombatEnabled: false,
		FeralCombatSkill:   0,
	}
}

type UnitStat int

func (s UnitStat) IsStat() bool                                 { return int(s) < int(Len) }
func (s UnitStat) IsPseudoStat() bool                           { return !s.IsStat() }
func (s UnitStat) EqualsStat(other Stat) bool                   { return int(s) == int(other) }
func (s UnitStat) EqualsPseudoStat(other proto.PseudoStat) bool { return int(s) == int(other) }
func (s UnitStat) StatIdx() int {
	if !s.IsStat() {
		panic("Is a pseudo stat")
	}
	return int(s)
}
func (s UnitStat) PseudoStatIdx() int {
	if s.IsStat() {
		panic("Is a regular stat")
	}
	return int(s) - int(Len)
}
func (s UnitStat) AddToStatsProto(p *proto.UnitStats, value float64) {
	if s.IsStat() {
		p.Stats[s.StatIdx()] += value
	} else {
		p.PseudoStats[s.PseudoStatIdx()] += value
	}
}

func UnitStatFromIdx(s int) UnitStat                     { return UnitStat(s) }
func UnitStatFromStat(s Stat) UnitStat                   { return UnitStat(s) }
func UnitStatFromPseudoStat(s proto.PseudoStat) UnitStat { return UnitStat(int(s) + int(Len)) }
