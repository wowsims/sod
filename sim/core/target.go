package core

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type Encounter struct {
	Duration          time.Duration
	DurationVariation time.Duration
	Targets           []*Target
	TargetUnits       []*Unit

	ExecuteProportion_20 float64
	ExecuteProportion_25 float64
	ExecuteProportion_35 float64

	EndFightAtHealth float64
	// DamageTaken is used to track health fights instead of duration fights.
	//  Once primary target has taken its health worth of damage, fight ends.
	DamageTaken float64
	// In health fight: set to true until we get something to base on
	DurationIsEstimate bool

	// Value to multiply by, for damage spells which are subject to the aoe cap.
	aoeCapMultiplier float64
}

func NewEncounter(options *proto.Encounter) Encounter {
	options.ExecuteProportion_25 = max(options.ExecuteProportion_25, options.ExecuteProportion_20)
	options.ExecuteProportion_35 = max(options.ExecuteProportion_35, options.ExecuteProportion_25)

	encounter := Encounter{
		Duration:             DurationFromSeconds(options.Duration),
		DurationVariation:    DurationFromSeconds(options.DurationVariation),
		ExecuteProportion_20: max(options.ExecuteProportion_20, 0),
		ExecuteProportion_25: max(options.ExecuteProportion_25, 0),
		ExecuteProportion_35: max(options.ExecuteProportion_35, 0),
		Targets:              []*Target{},
	}
	// If UseHealth is set, we use the sum of targets health.
	if options.UseHealth {
		for _, t := range options.Targets {
			encounter.EndFightAtHealth += t.Stats[stats.Health]
		}
		if encounter.EndFightAtHealth == 0 {
			encounter.EndFightAtHealth = 1 // default to something so we don't instantly end without anything.
		}
	}

	for targetIndex, targetOptions := range options.Targets {
		target := NewTarget(targetOptions, int32(targetIndex))
		encounter.Targets = append(encounter.Targets, target)
		encounter.TargetUnits = append(encounter.TargetUnits, &target.Unit)
	}
	if len(encounter.Targets) == 0 {
		// Add a dummy target. The only case where targets aren't specified is when
		// computing character stats, and targets won't matter there.
		target := NewTarget(&proto.Target{}, 0)
		encounter.Targets = append(encounter.Targets, target)
		encounter.TargetUnits = append(encounter.TargetUnits, &target.Unit)
	}

	if encounter.EndFightAtHealth > 0 {
		// Until we pre-sim set duration to 10m
		encounter.Duration = time.Minute * 10
		encounter.DurationIsEstimate = true
	}

	encounter.updateAOECapMultiplier()

	return encounter
}

func (encounter *Encounter) AOECapMultiplier() float64 {
	return encounter.aoeCapMultiplier
}
func (encounter *Encounter) updateAOECapMultiplier() {
	encounter.aoeCapMultiplier = min(10/float64(len(encounter.Targets)), 1)
}

func (encounter *Encounter) doneIteration(sim *Simulation) {
	for i := range encounter.Targets {
		target := encounter.Targets[i]
		target.doneIteration(sim)
	}
}

func (encounter *Encounter) GetMetricsProto() *proto.EncounterMetrics {
	metrics := &proto.EncounterMetrics{
		Targets: make([]*proto.UnitMetrics, len(encounter.Targets)),
	}

	i := 0
	for _, target := range encounter.Targets {
		metrics.Targets[i] = target.GetMetricsProto()
		i++
	}

	return metrics
}

// Target is an enemy/boss that can be the target of player attacks/spells.
type Target struct {
	Unit

	AI TargetAI
}

func NewTarget(options *proto.Target, targetIndex int32) *Target {
	unitStats := stats.Stats{}
	if options.Stats != nil {
		copy(unitStats[:], options.Stats)
	}

	target := &Target{
		Unit: Unit{
			Type:        EnemyUnit,
			Index:       targetIndex,
			Label:       "Target " + strconv.Itoa(int(targetIndex)+1),
			Level:       options.Level,
			MobType:     options.MobType,
			auraTracker: newAuraTracker(),
			stats:       unitStats,
			PseudoStats: stats.NewPseudoStats(),
			Metrics:     NewUnitMetrics(),

			StatDependencyManager: stats.NewStatDependencyManager(),
		},
	}
	defaultRaidBossLevel := int32(CharacterMaxLevel + 3)
	target.GCD = target.NewTimer()
	if target.Level == 0 {
		target.Level = defaultRaidBossLevel
	}

	target.AddStatDependency(stats.Defense, stats.Dodge, MissDodgeParryBlockCritChancePerDefense)
	target.AddStatDependency(stats.Defense, stats.Parry, MissDodgeParryBlockCritChancePerDefense)
	target.AddStatDependency(stats.Defense, stats.Block, MissDodgeParryBlockCritChancePerDefense)

	target.PseudoStats.CanBlock = true
	target.PseudoStats.CanParry = true
	target.PseudoStats.ParryHaste = options.ParryHaste
	target.PseudoStats.InFrontOfTarget = true
	target.PseudoStats.DamageSpread = options.DamageSpread

	preset := GetPresetTargetWithID(options.Id)
	if preset != nil && preset.AI != nil {
		target.AI = preset.AI()
	}

	return target
}

func (target *Target) Reset(sim *Simulation) {
	target.Unit.reset(sim, nil)
	target.SetGCDTimer(sim, 0)
	if target.AI != nil {
		target.AI.Reset(sim)
	}
}

func (target *Target) NextTarget() *Target {
	nextIndex := target.Index + 1
	if nextIndex >= target.Env.GetNumTargets() {
		nextIndex = 0
	}
	return target.Env.GetTarget(nextIndex)
}

func (target *Target) GetMetricsProto() *proto.UnitMetrics {
	metrics := target.Metrics.ToProto()
	metrics.Name = target.Label
	metrics.UnitIndex = target.UnitIndex
	metrics.Auras = target.auraTracker.GetMetricsProto()
	return metrics
}

func (character *Character) IsTanking() bool {
	for _, target := range character.Env.Encounter.TargetUnits {
		if target.CurrentTarget == &character.Unit {
			return true
		}
	}
	return false
}

func GetWeaponSkill(unit *Unit, weapon *Item) float64 {
	if weapon == nil {
		return 0
	}

	if unit.PseudoStats.FeralCombatEnabled && unit.PseudoStats.FeralCombatSkill != 0 {
		return unit.PseudoStats.FeralCombatSkill
	}

	if weapon.HandType == proto.HandType_HandTypeTwoHand {
		switch weapon.WeaponType {
		case proto.WeaponType_WeaponTypeAxe:
			return unit.PseudoStats.TwoHandedAxesSkill
		case proto.WeaponType_WeaponTypeMace:
			return unit.PseudoStats.TwoHandedMacesSkill
		case proto.WeaponType_WeaponTypeSword:
			return unit.PseudoStats.TwoHandedSwordsSkill
		case proto.WeaponType_WeaponTypePolearm:
			return unit.PseudoStats.PolearmsSkill
		case proto.WeaponType_WeaponTypeStaff:
			return unit.PseudoStats.StavesSkill
		default:
			return 0
		}
	} else if weapon.RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown {
		switch weapon.RangedWeaponType {
		case proto.RangedWeaponType_RangedWeaponTypeBow:
			return unit.PseudoStats.BowsSkill
		case proto.RangedWeaponType_RangedWeaponTypeCrossbow:
			return unit.PseudoStats.CrossbowsSkill
		case proto.RangedWeaponType_RangedWeaponTypeGun:
			return unit.PseudoStats.GunsSkill
		case proto.RangedWeaponType_RangedWeaponTypeThrown:
			return unit.PseudoStats.ThrownSkill
		default:
			return 0
		}
	} else if weapon.HandType == proto.HandType_HandTypeMainHand || weapon.HandType == proto.HandType_HandTypeOneHand || weapon.HandType == proto.HandType_HandTypeOffHand {
		switch weapon.WeaponType {
		case proto.WeaponType_WeaponTypeAxe:
			return unit.PseudoStats.AxesSkill
		case proto.WeaponType_WeaponTypeFist:
			return unit.PseudoStats.UnarmedSkill
		case proto.WeaponType_WeaponTypeMace:
			return unit.PseudoStats.MacesSkill
		case proto.WeaponType_WeaponTypeSword:
			return unit.PseudoStats.SwordsSkill
		case proto.WeaponType_WeaponTypeDagger:
			return unit.PseudoStats.DaggersSkill
		default:
			return 0
		}
	} else if weapon.HandType == proto.HandType_HandTypeUnknown && weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeUnknown {
		// Needed for Paladin Hammer of Wrath to use 300 skill
		return float64(unit.Level) * 5.0
	} else {
		return 0
	}
}

// Holds cached values for outcome/damage calculations, for a specific attacker+defender pair.
// These are updated dynamically when attacker or defender stats change.
type AttackTable struct {
	Attacker *Unit
	Defender *Unit

	Weapon *Item

	BaseMissChance      float64
	HitSuppression      float64
	BaseSpellMissChance float64
	BaseBlockChance     float64
	BaseDodgeChance     float64
	BaseParryChance     float64
	BaseGlanceChance    float64
	BaseCritChance      float64

	GlanceMultiplierMin  float64
	GlanceMultiplierMax  float64
	MeleeCritSuppression float64
	SpellCritSuppression float64

	// All "Apply Aura: Mod All Damage Done Against Creature" effects in Vanilla and TBC also increase the CritMultiplier.
	//  Explicitly for hunters' "Monster Slaying" and "Humanoid Slaying", but likewise for rogues' "Murder", or trolls' "Beastslaying".
	CritMultiplier float64

	DamageDealtMultiplier  float64 // attacker buff, applied in applyAttackerModifiers()
	DamageTakenMultiplier  float64 // defender debuff, applied in applyTargetModifiers()
	HealingDealtMultiplier float64

	// This is for "Apply Aura: Mod Damage Done By Caster" effects.
	// If set, the damage taken multiplier is multiplied by the callbacks result.
	DamageDoneByCasterMultiplier func(spell *Spell, attackTable *AttackTable) float64
}

func NewAttackTable(attacker *Unit, defender *Unit, weapon *Item) *AttackTable {
	// Source: https://github.com/magey/classic-warrior/wiki/Attack-table
	table := &AttackTable{
		Attacker: attacker,
		Defender: defender,
		Weapon:   weapon,

		CritMultiplier: 1,

		DamageDealtMultiplier:  1,
		DamageTakenMultiplier:  1,
		HealingDealtMultiplier: 1,
	}

	if defender.Type == EnemyUnit {
		baseWeaponSkill := float64(attacker.Level * 5)
		weaponSkill := baseWeaponSkill + GetWeaponSkill(attacker, weapon)
		targetDefense := float64(defender.Level * 5)

		if targetDefense-weaponSkill > 10 {
			table.HitSuppression = (targetDefense - weaponSkill - 10) * 0.002
			table.BaseMissChance = 0.05 + (targetDefense-weaponSkill)*0.002
		} else {
			table.HitSuppression = 0
			table.BaseMissChance = 0.05 + (targetDefense-weaponSkill)*0.001
		}

		if targetDefense-baseWeaponSkill > 10 {
			table.BaseParryChance = 0.05 + (targetDefense-baseWeaponSkill)*0.006 // = 14
		} else {
			table.BaseParryChance = 0.05 + (targetDefense-baseWeaponSkill)*0.001 // = 5 / 5.5 / 6
		}

		table.BaseSpellMissChance = UnitLevelFloat64(defender.Level-attacker.Level, 0.04, 0.05, 0.06, 0.17)
		table.BaseBlockChance = 0.05
		table.BaseDodgeChance = 0.05 + (targetDefense-weaponSkill)*0.001
		table.BaseGlanceChance = 0.1 + (targetDefense-baseWeaponSkill)*0.02

		table.GlanceMultiplierMin = max(min(1.3-0.05*(targetDefense-weaponSkill), 0.91), 0.01)
		table.GlanceMultiplierMax = max(min(1.2-0.03*(targetDefense-weaponSkill), 0.99), 0.2)

		if targetDefense > baseWeaponSkill {
			table.MeleeCritSuppression = (targetDefense - baseWeaponSkill) * 0.002
		} else {
			table.MeleeCritSuppression = (targetDefense - baseWeaponSkill) * 0.0004
		}

		// TODO (maybe): This is technically not correct, but it shouldn't matter outside of edge cases.
		// These 1.8% should only be subtracted from crit chance gained from auras,
		// i.e. talents, gear and buffs, NOT base crit and crit from agility!
		// See https://github.com/magey/classic-warrior/wiki/Attack-table#critical-strike
		// That means if a character with <2% crit from auras attacks a +3 level target the sim will be wrong.
		// The chance of that being the case once bosses are +3 in SoD should be very small though.
		// Most (all?) affected specs have crit in their talents to begin with.
		if (defender.Level - attacker.Level) >= 3 {
			table.MeleeCritSuppression += 0.018
		}
		table.SpellCritSuppression = UnitLevelFloat64(defender.Level-attacker.Level, 0, 0, 0.003, 0.021)
	} else {

		levelDelta := 0.0004 * 5 * float64(defender.Level-attacker.Level)

		table.BaseSpellMissChance = 0.05

		// Apply base Parry
		if defender.PseudoStats.CanParry {
			table.BaseParryChance = levelDelta // + 0.05 applied as stats in character.go
		} else {
			table.BaseParryChance = 0
		}
		// Apply base Block
		if defender.PseudoStats.CanBlock {
			table.BaseBlockChance = levelDelta // + 0.05 applied as stats in character.go
		} else {
			table.BaseBlockChance = 0
		}

		table.BaseMissChance = 0.05 + levelDelta
		table.BaseDodgeChance = levelDelta // base dodge applied with class base stats
		table.BaseCritChance = 0.05 - levelDelta
	}

	return table
}

func ModNonMeleeAttackTable(table *AttackTable, attacker *Unit, defender *Unit, weapon *Item) {
	weaponSkill := float64(attacker.Level*5) + float64(GetWeaponSkill(attacker, weapon))

	table.BaseGlanceChance = 0.8 // min((float64(attacker.Level)-10)*0.03, 0.6)

	table.GlanceMultiplierMin = max(min(1.3-0.05*(float64(defender.Level*5)-weaponSkill)-0.7, 0.6), 0.01)
	table.GlanceMultiplierMax = max(min(1.2-0.03*(float64(defender.Level*5)-weaponSkill)-0.3, 0.99), 0.2)
}
