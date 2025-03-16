package core

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type OnApplyEffects func(aura *Aura, sim *Simulation, target *Unit, spell *Spell)
type ApplySpellResults func(sim *Simulation, target *Unit, spell *Spell)
type ExpectedDamageCalculator func(sim *Simulation, target *Unit, spell *Spell, useSnapshot bool) *SpellResult
type CanCastCondition func(sim *Simulation, target *Unit) bool

type SpellConfig struct {
	// See definition of Spell (below) for comments on these.
	ActionID
	// Used to identify spells with multiple ranks that need to be referenced
	ClassSpellMask uint64
	SpellSchool    SpellSchool
	DefenseType    DefenseType
	ProcMask       ProcMask
	Flags          SpellFlag
	CastType       proto.CastType
	MissileSpeed   float64
	BaseCost       float64
	MetricSplits   int
	Rank           int
	RequiredLevel  int

	ManaCost   ManaCostOptions
	EnergyCost EnergyCostOptions
	RageCost   RageCostOptions
	FocusCost  FocusCostOptions

	Cast               CastConfig
	ExtraCastCondition CanCastCondition

	BonusHitRating  float64
	BonusCritRating float64

	CritDamageBonus float64

	BaseDamageMultiplierAdditivePct     int64
	DamageMultiplier                    float64
	DamageMultiplierAdditivePct         int64
	ImpactDamageMultiplierAdditivePct   int64
	PeriodicDamageMultiplierAdditivePct int64

	BonusDamage      float64 // Bonus scaling power e.g. Idol of the Moon "Increases the damage of X spell by N" https://www.wowhead.com/classic/item=23197/idol-of-the-moon
	BonusCoefficient float64 // EffectBonusCoefficient in SpellEffect client DB table, "SP mod" on Wowhead (not necessarily shown there even if > 0)

	ThreatMultiplier float64

	FlatThreatBonus float64

	// Chance to avoid interruption caused by damage while casting spell
	// Apply Aura: Modifies Pushback Reduction (9)
	PushbackReduction float64

	// Performs the actions of this spell.
	ApplyEffects ApplySpellResults

	// Optional field. Calculates expected average damage.
	ExpectedInitialDamage ExpectedDamageCalculator
	ExpectedTickDamage    ExpectedDamageCalculator

	Dot    DotConfig
	Hot    DotConfig
	Shield ShieldConfig

	RelatedAuras []AuraArray
}

type Spell struct {
	// ID for this spell.
	ActionID

	// Used to identify spells with multiple ranks that need to be referenced
	// The specific class spell id should be a unique bit
	ClassSpellMask uint64

	// The unit who will perform this spell.
	Unit *Unit

	SpellSchool       SpellSchool         // Schoolmask of all schools this spell uses. Do not change this! Whatever you try to do is a hack and probably wrong.
	SchoolIndex       stats.SchoolIndex   // Do not change this! Whatever you try to do is a hack and probably wrong.
	SchoolBaseIndices []stats.SchoolIndex // Base school indices for multi schools. Do not change this! Whatever you try to do is a hack and probably wrong.
	DefenseType       DefenseType

	// Controls which effects can proc from this spell.
	ProcMask ProcMask

	// Flags
	Flags SpellFlag

	// From which slot this spell cast. Usually from Mainhand
	CastType proto.CastType

	// Speed in yards/second. Spell missile speeds can be found in the game data.
	// Example: https://wow.tools/dbc/?dbc=spellmisc&build=3.4.0.44996
	MissileSpeed float64

	Rank          int
	RequiredLevel int

	ResourceMetrics *ResourceMetrics
	healthMetrics   []*ResourceMetrics

	Cost *SpellCost // Cost for the spell.

	DefaultCast        Cast // Default cast parameters with all static effects applied.
	CD                 *SpellCooldown
	SharedCD           *SpellCooldown
	ExtraCastCondition CanCastCondition

	castTimeFn func(spell *Spell) time.Duration // allows to override CastTime()

	castFn CastSuccessFunc // Performs a cast of this spell.

	SpellMetrics []SpellMetrics

	splitIdx          int32
	splitSpellMetrics [][]SpellMetrics // Used to split metrics by some condition, via SetMetricsSplit
	splitTags         []int32          // Tags for each splitSpellMetrics used in doneIteration, defaults to the metrics splitIdx.

	casts int // Sum of casts on all targets, for efficient CPM calculation

	// Performs the actions of this spell.
	ApplyEffects ApplySpellResults

	// Optional field. Calculates expected average damage.
	expectedInitialDamageInternal ExpectedDamageCalculator
	expectedTickDamageInternal    ExpectedDamageCalculator

	// The current or most recent cast data.
	CurCast    Cast
	LastCastAt time.Duration

	BonusHitRating     float64
	BonusCritRating    float64
	CastTimeMultiplier float64

	baseDamageMultiplierAdditivePct     int64 // Stores an integer representation of the Spell's Base Damage Multiplier
	damageMultiplierAdditivePct         int64 // Stores an integer representation of the Spell's Additive Damage Multiplier before Imapct or Periodic-only bonuses
	impactDamageMultiplierAdditivePct   int64 // Stores an integer representation of the Spell's Additive Impact Damage Multiplier
	periodicDamageMultiplierAdditivePct int64 // Stores an integer representation of the Spell's Additive Periodic Damage Multiplier

	baseDamageMultiplier     float64 // Stores the Spell's calculated Base Damage Multiplier
	damageMultiplier         float64 // Stores the Spell's calculated Damage Multiplier before Imapct or Periodic-only bonuses
	impactDamageMultiplier   float64 // Stores the Spell's calculated Impact Damage Multiplier
	periodicDamageMultiplier float64 // Stores the Spell's calculated Damage Multiplier

	BonusDamage      float64 // Bonus scaling power e.g. Idol of the Moon "Increases the damage of X spell by N" https://www.wowhead.com/classic/item=23197/idol-of-the-moon
	BonusCoefficient float64 // EffectBonusCoefficient in SpellEffect client DB table, "SP mod" on Wowhead (not necessarily shown there even if > 0)

	CritDamageBonus float64

	// Multiplier for all threat generated by this effect.
	ThreatMultiplier float64

	// Adds a fixed amount of threat to this spell, before multipliers.
	FlatThreatBonus float64

	// Chance to avoid interruption caused by damage while casting spell
	// Apply Aura: Modifies Pushback Reduction (9)
	PushbackReduction float64

	resultCache SpellResult

	dots   DotArray
	aoeDot *Dot

	shields    ShieldArray
	selfShield *Shield

	// Per-target auras that are related to this spell, usually buffs or debuffs applied by the spell.
	RelatedAuras []AuraArray

	// Reference to a spell to be considered as the CD
	// Defaults to this spell (Used for Next Melee spells)
	CdSpell *Spell
}

func (unit *Unit) OnSpellRegistered(handler SpellRegisteredHandler) {
	for _, spell := range unit.Spellbook {
		handler(spell)
	}
	unit.spellRegistrationHandlers = append(unit.spellRegistrationHandlers, handler)
}

// Registers a new spell to the unit. Returns the newly created spell.
func (unit *Unit) RegisterSpell(config SpellConfig) *Spell {
	if len(unit.Spellbook) > 200 {
		panic(fmt.Sprintf("Over 200 registered spells when registering %s! There is probably a spell being registered every iteration.", config.ActionID))
	}

	// Default the other damage multiplier to 1 if only one or the other is set.
	if config.DamageMultiplierAdditivePct != 0 && config.DamageMultiplier == 0 {
		config.DamageMultiplier = 1
	}

	// Default CastSlot to mainhand
	if config.CastType == proto.CastType_CastTypeUnknown {
		config.CastType = proto.CastType_CastTypeMainHand
	}

	if (config.DamageMultiplier != 0 || config.ThreatMultiplier != 0) && config.ProcMask == ProcMaskUnknown {
		panic("ProcMask for spell " + config.ActionID.String() + " not set")
	}

	if (config.DamageMultiplier != 0 || config.ThreatMultiplier != 0) && config.SpellSchool == SpellSchoolNone {
		panic("SpellSchool for spell " + config.ActionID.String() + " not set")
	}

	if config.Cast.CD.Timer != nil && config.Cast.CD.Duration == 0 {
		panic("Cast.CD w/o Duration specified for spell " + config.ActionID.String())
	}

	if config.Cast.SharedCD.Timer != nil && config.Cast.SharedCD.Duration == 0 {
		panic("Cast.SharedCD w/o Duration specified for spell " + config.ActionID.String())
	}

	if config.Cast.CastTime == nil {
		config.Cast.CastTime = func(spell *Spell) time.Duration {
			return spell.Unit.ApplyCastSpeedForSpell(spell.DefaultCast.CastTime, spell)
		}
	}

	spell := &Spell{
		ActionID:       config.ActionID,
		ClassSpellMask: config.ClassSpellMask,
		DefenseType:    config.DefenseType,
		Unit:           unit,
		ProcMask:       config.ProcMask,
		Flags:          config.Flags,
		CastType:       config.CastType,
		MissileSpeed:   config.MissileSpeed,

		SpellSchool:       config.SpellSchool,
		SchoolIndex:       config.SpellSchool.GetSchoolIndex(),
		SchoolBaseIndices: config.SpellSchool.GetBaseIndices(),

		DefaultCast:        config.Cast.DefaultCast,
		CD:                 newSpellCooldown(config.Cast.CD),
		SharedCD:           newSpellCooldown(config.Cast.SharedCD),
		ExtraCastCondition: config.ExtraCastCondition,

		castTimeFn: config.Cast.CastTime,

		ApplyEffects: config.ApplyEffects,

		expectedInitialDamageInternal: config.ExpectedInitialDamage,
		expectedTickDamageInternal:    config.ExpectedTickDamage,

		BonusDamage:        config.BonusDamage,
		BonusHitRating:     config.BonusHitRating,
		BonusCritRating:    config.BonusCritRating,
		CastTimeMultiplier: 1,

		CritDamageBonus: 1 + config.CritDamageBonus,

		baseDamageMultiplierAdditivePct:     config.BaseDamageMultiplierAdditivePct,
		damageMultiplier:                    config.DamageMultiplier,
		damageMultiplierAdditivePct:         config.DamageMultiplierAdditivePct,
		impactDamageMultiplierAdditivePct:   config.ImpactDamageMultiplierAdditivePct,
		periodicDamageMultiplierAdditivePct: config.PeriodicDamageMultiplierAdditivePct,

		BonusCoefficient: config.BonusCoefficient,

		ThreatMultiplier: config.ThreatMultiplier,
		FlatThreatBonus:  config.FlatThreatBonus,

		PushbackReduction: config.PushbackReduction,

		splitSpellMetrics: make([][]SpellMetrics, max(1, config.MetricSplits)),
		splitTags:         make([]int32, max(1, config.MetricSplits)),

		RelatedAuras: config.RelatedAuras,
	}

	spell.updateBaseDamageMultiplier()
	spell.updateImpactDamageMultiplier()
	spell.updatePeriodicDamageMultiplier()

	spell.Rank = config.Rank
	spell.RequiredLevel = config.RequiredLevel

	spell.CdSpell = spell

	if config.ManaCost.BaseCost != 0 || config.ManaCost.FlatCost != 0 {
		spell.Cost = newManaCost(spell, config.ManaCost)
	} else if config.EnergyCost.Cost != 0 {
		spell.Cost = newEnergyCost(spell, config.EnergyCost)
	} else if config.RageCost.Cost != 0 {
		spell.Cost = newRageCost(spell, config.RageCost)
	} else if config.FocusCost.Cost != 0 {
		spell.Cost = newFocusCost(spell, config.FocusCost)
	}

	if spell.Cost != nil {
		spell.DefaultCast.Cost = spell.Cost.BaseCost
	}

	spell.createDots(config.Dot, false)
	spell.createDots(config.Hot, true)
	spell.createShields(config.Shield)

	var emptyCast Cast

	if spell.DefaultCast == emptyCast && spell.Cost != nil {
		panic("Empty DefaultCast with a cost for spell " + config.ActionID.String())
	}

	if spell.DefaultCast.GCD == 0 && spell.DefaultCast.CastTime == 0 {
		config.Cast.IgnoreHaste = true
	}

	if spell.DefaultCast == emptyCast {
		if config.ExtraCastCondition == nil && config.Cast.CD.Timer == nil && config.Cast.SharedCD.Timer == nil {
			spell.castFn = spell.makeCastFuncAutosOrProcs()
		} else {
			spell.castFn = spell.makeCastFuncSimple()
		}
	} else {
		spell.castFn = spell.makeCastFunc(config.Cast)
	}

	if spell.ApplyEffects == nil {
		spell.ApplyEffects = func(*Simulation, *Unit, *Spell) {}
	}

	unit.Spellbook = append(unit.Spellbook, spell)

	for _, handler := range unit.spellRegistrationHandlers {
		handler(spell)
	}

	if unit.Env != nil && unit.Env.IsFinalized() {
		spell.finalize()
	}

	return spell
}

// Returns the first registered spell with the given ID, or nil if there are none.
func (unit *Unit) GetSpell(actionID ActionID) *Spell {
	for _, spell := range unit.Spellbook {
		if spell.ActionID.SameAction(actionID) {
			return spell
		}
	}
	return nil
}

// Retrieves an existing spell with the same ID as the config uses, or registers it if there is none.
func (unit *Unit) GetOrRegisterSpell(config SpellConfig) *Spell {
	registered := unit.GetSpell(config.ActionID)
	if registered == nil {
		return unit.RegisterSpell(config)
	} else {
		return registered
	}
}

func (spell *Spell) Dots() []*Dot {
	return spell.dots
}
func (spell *Spell) Dot(target *Unit) *Dot {
	return spell.dots.Get(target)
}
func (spell *Spell) CurDot() *Dot {
	return spell.dots.Get(spell.Unit.CurrentTarget)
}
func (spell *Spell) AOEDot() *Dot {
	return spell.aoeDot
}
func (spell *Spell) DotOrAOEDot(target *Unit) *Dot {
	if spell.aoeDot != nil {
		return spell.aoeDot
	}
	return spell.dots.Get(target)
}
func (spell *Spell) Hot(target *Unit) *Dot {
	return spell.dots.Get(target)
}
func (spell *Spell) CurHot() *Dot {
	return spell.dots.Get(spell.Unit.CurrentTarget)
}
func (spell *Spell) AOEHot() *Dot {
	return spell.aoeDot
}
func (spell *Spell) SelfHot() *Dot {
	return spell.aoeDot
}
func (spell *Spell) Shield(target *Unit) *Shield {
	return spell.shields.Get(target)
}
func (spell *Spell) SelfShield() *Shield {
	return spell.selfShield
}

// Metrics for the current iteration
func (spell *Spell) CurDamagePerCast() float64 {
	if spell.SpellMetrics[0].Casts == 0 {
		return 0
	} else {
		casts := int32(0)
		damage := 0.0
		for _, opponent := range spell.Unit.GetOpponents() {
			casts += spell.SpellMetrics[opponent.UnitIndex].Casts
			damage += spell.SpellMetrics[opponent.UnitIndex].TotalDamage
		}
		return damage / float64(casts)
	}
}

// Current casts per minute
func (spell *Spell) CurCPM(sim *Simulation) float64 {
	if sim.CurrentTime <= 0 {
		return 0
	}
	casts := float64(spell.casts)
	minutes := float64(sim.CurrentTime) / float64(time.Minute)
	return casts / minutes
}

func (spell *Spell) finalize() {
	spell.splitTags[0] = spell.Tag
	for i := range spell.splitSpellMetrics {
		spell.splitSpellMetrics[i] = make([]SpellMetrics, len(spell.Unit.Env.AllUnits))
		if spell.splitTags[i] == 0 {
			spell.splitTags[i] = int32(i) // default to tag = splitIndex
		}
	}

	spell.SpellMetrics = spell.splitSpellMetrics[0]

	// Set the "static" "default" cost here
	if spell.Cost != nil {
		spell.DefaultCast.Cost = spell.Cost.GetCurrentCost()
	}
}

func (spell *Spell) reset(_ *Simulation) {
	for i := range spell.splitSpellMetrics {
		for j := range spell.SpellMetrics {
			spell.splitSpellMetrics[i][j] = SpellMetrics{}
		}
	}
	spell.casts = 0
	spell.LastCastAt = 0
}

func (spell *Spell) GetMetricsSplitIdx() int32 {
	return spell.splitIdx
}

func (spell *Spell) SetMetricsSplit(splitIdx int32) {
	spell.splitIdx = splitIdx
	spell.SpellMetrics = spell.splitSpellMetrics[splitIdx]
	spell.Tag = spell.splitTags[splitIdx]
}

// TagSplitMetric allows to override the default tag for a given splitIdx. Use after spell registration.
func (spell *Spell) TagSplitMetric(splitIdx int32, tag int32) {
	spell.splitTags[splitIdx] = tag
}

func (spell *Spell) doneIteration() {
	if spell.Flags.Matches(SpellFlagNoMetrics) {
		return
	}

	for i, spellMetrics := range spell.splitSpellMetrics {
		spell.Unit.Metrics.addSpellMetrics(spell, spell.ActionID.WithTag(spell.splitTags[i]), spellMetrics)
	}
}

func (spell *Spell) HealthMetrics(target *Unit) *ResourceMetrics {
	if spell.healthMetrics == nil {
		spell.healthMetrics = make([]*ResourceMetrics, len(spell.Unit.AttackTables))
	}
	if spell.healthMetrics[target.UnitIndex] == nil {
		spell.healthMetrics[target.UnitIndex] = target.NewHealthMetrics(spell.ActionID)
	}
	return spell.healthMetrics[target.UnitIndex]
}

func (spell *Spell) ReadyAt() time.Duration {
	return BothTimersReadyAt(spell.CdSpell.CD.Timer, spell.CdSpell.SharedCD.Timer)
}

func (spell *Spell) IsReady(sim *Simulation) bool {
	if spell == nil {
		return false
	}

	return BothTimersReady(spell.CdSpell.CD.Timer, spell.CdSpell.SharedCD.Timer, sim)
}

func (spell *Spell) TimeToReady(sim *Simulation) time.Duration {
	return MaxTimeToReady(spell.CdSpell.CD.Timer, spell.CdSpell.SharedCD.Timer, sim)
}

// Returns whether a call to Cast() would be successful, without actually doing a cast.
func (spell *Spell) CanCast(sim *Simulation, target *Unit) bool {
	if spell == nil {
		return false
	}

	if spell.Flags.Matches(SpellFlagSwapped) {
		//if sim.Log != nil {
		//	sim.Log("Cant cast because of item swap")
		//}
		return false
	}

	if spell.ExtraCastCondition != nil && !spell.ExtraCastCondition(sim, target) {
		//if sim.Log != nil {
		//	sim.Log("Cant cast because of extra condition")
		//}
		return false
	}

	// While moving only instant casts are possible
	if spell.DefaultCast.CastTime > 0 && spell.Unit.IsMoving() {
		//if sim.Log != nil {
		//	sim.Log("Cant cast because moving")
		//}
		return false
	}

	// While casting no other action is possible except rare cast-while-casting spells
	if spell.Unit.IsCasting(sim) && !spell.Flags.Matches(SpellFlagCastWhileCasting) {
		//if sim.Log != nil {
		//	sim.Log("Cant cast because already casting")
		//}
		return false
	}

	// While channeling no other action is possible except rare cast-while-channeling spells
	if spell.Unit.IsChanneling(sim) && !spell.Flags.Matches(SpellFlagCastWhileChanneling) && (spell.Unit.Rotation.interruptChannelIf == nil || !spell.Unit.Rotation.interruptChannelIf.GetBool(sim)) {
		//if sim.Log != nil {
		//	sim.Log("Cant cast because already channeling")
		//}
		return false
	}

	if spell.DefaultCast.GCD > 0 && !spell.Unit.GCD.IsReady(sim) {
		//if sim.Log != nil {
		//	sim.Log("Cant cast because of GCD")
		//}
		return false
	}

	if !BothTimersReady(spell.CD.Timer, spell.SharedCD.Timer, sim) {
		//if sim.Log != nil {
		//	sim.Log("Cant cast because of CDs")
		//}
		return false
	}

	if spell.Cost != nil {
		if !spell.Cost.MeetsRequirement(sim, spell) {
			//if sim.Log != nil {
			//	sim.Log("Cant cast because of resource cost")
			//}
			return false
		}
	}

	return true
}

func (spell *Spell) Cast(sim *Simulation, target *Unit) bool {
	if target == nil {
		target = spell.Unit.CurrentTarget
	}
	return spell.castFn(sim, target)
}

func (spell *Spell) applyEffects(sim *Simulation, target *Unit) {
	spell.SpellMetrics[target.UnitIndex].Casts++
	spell.casts++

	// Not sure if we want to split this flag into its own?
	// Both are used to optimize away unneccesery calls and 99%
	// of the time are gonna be used together. For now just in one
	if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
		spell.Unit.OnApplyEffects(sim, target, spell)
	}

	spell.ApplyEffects(sim, target, spell)
}

func (spell *Spell) ApplyAOEThreatIgnoreMultipliers(threatAmount float64) {
	numTargets := spell.Unit.Env.GetNumTargets()
	for i := int32(0); i < numTargets; i++ {
		spell.SpellMetrics[i].TotalThreat += threatAmount
	}
}
func (spell *Spell) ApplyAOEThreat(threatAmount float64) {
	spell.ApplyAOEThreatIgnoreMultipliers(threatAmount * spell.Unit.PseudoStats.ThreatMultiplier)
}

func (spell *Spell) finalizeExpectedDamage(result *SpellResult) {
	if !spell.SpellSchool.Matches(SpellSchoolPhysical) {
		result.Damage /= result.ResistanceMultiplier
		averagePartialResistMultiplier := 1.0 - AverageMagicPartialResistPerLevelMultiplier*float64(result.Target.Level-spell.Unit.Level)
		result.Damage *= averagePartialResistMultiplier
		result.ResistanceMultiplier = averagePartialResistMultiplier
	}
	result.inUse = false
}
func (spell *Spell) ExpectedInitialDamage(sim *Simulation, target *Unit) float64 {
	result := spell.expectedInitialDamageInternal(sim, target, spell, false)
	spell.finalizeExpectedDamage(result)
	return result.Damage
}
func (spell *Spell) ExpectedTickDamage(sim *Simulation, target *Unit) float64 {
	result := spell.expectedTickDamageInternal(sim, target, spell, false)
	spell.finalizeExpectedDamage(result)
	return result.Damage
}
func (spell *Spell) ExpectedTickDamageFromCurrentSnapshot(sim *Simulation, target *Unit) float64 {
	result := spell.expectedTickDamageInternal(sim, target, spell, true)
	spell.finalizeExpectedDamage(result)
	return result.Damage
}

// Time until either the cast is finished or GCD is ready again, whichever is longer
func (spell *Spell) EffectiveCastTime() time.Duration {
	// TODO: this is wrong for spells like shadowfury, that have a GCD of less than 1s
	return max(spell.Unit.SpellGCD(),
		spell.Unit.ApplyCastSpeedForSpell(spell.DefaultCast.EffectiveTime(), spell))
}

// Time until the cast is finished (ignoring GCD)
func (spell *Spell) CastTime() time.Duration {
	return spell.castTimeFn(spell)
}

func (spell *Spell) TravelTime() time.Duration {
	if spell.MissileSpeed == 0 {
		return 0
	} else {
		return time.Duration(float64(time.Second) * spell.Unit.DistanceFromTarget / spell.MissileSpeed)
	}
}

// Returns true if the given mask matches the spell mask
func (spell *Spell) Matches(mask uint64) bool {
	return spell.ClassSpellMask&mask > 0
}

// Applies an additive multiplier to spell base damage. Equivalent to Modifies Spell Effectiveness (8).
func (spell *Spell) ApplyAdditiveBaseDamageBonus(percent int64) {
	spell.baseDamageMultiplierAdditivePct += percent
	spell.updateBaseDamageMultiplier()
}

func (spell *Spell) SetMultiplicativeDamageBonus(multiplier float64) {
	spell.damageMultiplier = multiplier
	spell.updateImpactDamageMultiplier()
	spell.updatePeriodicDamageMultiplier()
}

// Applies a multiplicative multiplier to full Direct and Periodic spell damage. Equivalent to Mod Damage Done % or similar effects.
func (spell *Spell) ApplyMultiplicativeDamageBonus(multiplier float64) {
	spell.damageMultiplier *= multiplier
	spell.updateImpactDamageMultiplier()
	spell.updatePeriodicDamageMultiplier()
}

func (spell *Spell) SetAdditiveDamageBonus(percent int64) {
	spell.damageMultiplierAdditivePct = percent
	spell.updateImpactDamageMultiplier()
	spell.updatePeriodicDamageMultiplier()
}

// Applies an additive multiplier to full Direct and Periodic spell damage. Equivalent to Modifies Damage/Healing Done + Modifies Periodic Damage/Healing Done (22).
func (spell *Spell) ApplyAdditiveDamageBonus(percent int64) {
	spell.damageMultiplierAdditivePct += percent
	spell.updateImpactDamageMultiplier()
	spell.updatePeriodicDamageMultiplier()
}

// Applies an additive multiplier to just Direct spell damage. Equivalent to Modifies Damage/Healing Done.
func (spell *Spell) ApplyAdditiveImpactDamageBonus(percent int64) {
	spell.impactDamageMultiplierAdditivePct += percent
	spell.updateImpactDamageMultiplier()
}

// Applies an additive multiplier to just Periodic spell dammage. Equivalent to Modifies Periodic Damage/Healing Done (22).
func (spell *Spell) ApplyAdditivePeriodicDamageBonus(percent int64) {
	spell.periodicDamageMultiplierAdditivePct += percent
	spell.updatePeriodicDamageMultiplier()
}

func (spell *Spell) updateBaseDamageMultiplier() {
	spell.baseDamageMultiplier = float64(100+spell.baseDamageMultiplierAdditivePct) / 100.0
}

func (spell *Spell) updateImpactDamageMultiplier() {
	spell.impactDamageMultiplier = spell.damageMultiplier * (float64(100+spell.damageMultiplierAdditivePct+spell.impactDamageMultiplierAdditivePct) / 100.0)
}

func (spell *Spell) updatePeriodicDamageMultiplier() {
	spell.periodicDamageMultiplier = spell.damageMultiplier * (float64(100+spell.damageMultiplierAdditivePct+spell.periodicDamageMultiplierAdditivePct) / 100.0)
}

func (spell *Spell) GetBaseDamageMultiplierAdditive() int64 {
	return spell.baseDamageMultiplierAdditivePct
}

func (spell *Spell) GetDamageMultiplier() float64 {
	return spell.damageMultiplier
}

func (spell *Spell) GetDamageMultiplierAdditive() int64 {
	return spell.damageMultiplierAdditivePct
}

func (spell *Spell) GetImpactDamageMultiplierAdditive() int64 {
	return spell.impactDamageMultiplierAdditivePct
}

func (spell *Spell) GetPeriodicDamageMultiplierAdditive() int64 {
	return spell.periodicDamageMultiplierAdditivePct
}

type CostType uint8

const (
	CostTypeUnknown CostType = iota

	CostTypeMana
	CostTypeEnergy
	// CostTypeComboPoints
	CostTypeRage
	CostTypeFocus
)

// Handles computing the cost of spells and checking whether the Unit
// meets them.
type SpellCostFunctions interface {
	// Get the type of resource used to cast the spell
	CostType() CostType

	// Whether the Unit associated with the spell meets the resource cost
	// requirements to cast the spell.
	MeetsRequirement(*Simulation, *Spell) bool

	// Returns a message for when the cast fails due to lack of resources.
	CostFailureReason(*Simulation, *Spell) string

	// Subtracts the resources used from a cast from the Unit.
	SpendCost(*Simulation, *Spell)

	// Space for handling refund mechanics. Not all spells provide refunds.
	IssueRefund(*Simulation, *Spell)
}

type SpellCost struct {
	BaseCost     float64 // The base power cost before all modifiers.
	FlatModifier int32   // Flat value added to base cost before pct mods
	Multiplier   int32   // Multiplier for cost, stored as an int, e.g. 0.5 is stored as 50
	spell        *Spell
	SpellCostFunctions
}

func (sc *SpellCost) ApplyCostModifiers(cost float64) float64 {
	spell := sc.spell
	cost = max(0, cost+float64(sc.FlatModifier))
	cost = max(0, cost*float64(spell.Unit.GetSchoolCostModifier(spell))/100)
	return max(0, cost*float64(sc.Multiplier)/100)
}

// Get power cost after all modifiers.
func (sc *SpellCost) GetCurrentCost() float64 {
	return sc.ApplyCostModifiers(sc.BaseCost)
}

func (spell *Spell) IssueRefund(sim *Simulation) {
	spell.Cost.IssueRefund(sim, spell)
}

type SpellCooldown struct {
	*Cooldown

	flatModifier time.Duration // Flat value added to base cooldown before pct mods

	// An int representation of the cooldown's multiplier percentage
	// Starts at 0 and any value added are an offset from 100%.
	// For example, for -50%, add -50
	multiplierPct int64
}

func newSpellCooldown(cd Cooldown) *SpellCooldown {
	return &SpellCooldown{
		Cooldown:      &cd,
		flatModifier:  0,
		multiplierPct: 0,
	}
}

func (cd *SpellCooldown) ApplyFlatCooldownMod(duration time.Duration) {
	cd.flatModifier += duration
}

func (cd *SpellCooldown) ApplyFlatPercentCooldownMod(percent int64) {
	cd.multiplierPct += percent
}

func (cd *SpellCooldown) applyCooldownModifiers(duration time.Duration) time.Duration {
	duration = max(0, duration+cd.flatModifier)
	return max(0, time.Duration(float64(duration)*float64(100+cd.multiplierPct)/100))
}

// Get cooldown after all modifiers.
func (cd *SpellCooldown) GetCurrentDuration() time.Duration {
	return cd.applyCooldownModifiers(cd.Duration)
}
