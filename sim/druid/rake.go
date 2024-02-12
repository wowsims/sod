package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

type RakeRankInfo struct {
	id            int32
	level         int32
	initialDamage float64
	dotDamage     float64
}

var rakeSpells = []RakeRankInfo{
	{
		id:            1822,
		level:         24,
		initialDamage: 19.0,
		dotDamage:     39.0,
	},
	{
		id:            1823,
		level:         34,
		initialDamage: 29.0,
		dotDamage:     57.0,
	},
	{

		id:            1824,
		level:         44,
		initialDamage: 43.0,
		dotDamage:     75.0,
	},
	{

		id:            9904,
		level:         54,
		initialDamage: 58.0,
		dotDamage:     96.0,
	},
}

var rakeTicks = 3.0

// SoD balance passive EFFECT1 and EFFECT2 mod for Rake
// See https://www.wowhead.com/classic/spell=436895/s03-tuning-and-overrides-passive-druid
// Mod Eff# should be base value only.
var baseDmgMultiplier = 1.5

func (druid *Druid) registerRakeSpell() {
	// Add highest available rake rank for level.
	for rank := len(rakeSpells) - 1; rank >= 0; rank-- {
		if druid.Level >= rakeSpells[rank].level {
			config := druid.newRakeSpellConfig(rakeSpells[rank])
			druid.Rake = druid.RegisterSpell(Cat, config)
			return
		}
	}
}

func (druid *Druid) newRakeSpellConfig(rakeRank RakeRankInfo) core.SpellConfig {
	damageInitial := rakeRank.initialDamage * baseDmgMultiplier
	damageDotTick := (rakeRank.dotDamage / rakeTicks) * baseDmgMultiplier

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: rakeRank.id},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial | core.ProcMaskSuppressedExtraAttackAura,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   40 - float64(druid.Talents.Ferocity),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
		CritMultiplier:   druid.MeleeCritMultiplier(1, 0),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			NumberOfTicks: 3,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = damageDotTick + 0.04*dot.Spell.MeleeAttackPower()
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = 0
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damageInitial + 0.04*spell.MeleeAttackPower()
			if druid.BleedCategories.Get(target).AnyActive() {
				baseDamage *= 1.3
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := damageInitial + 0.04*spell.MeleeAttackPower()
			initial := spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)

			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := (critChance * (spell.CritMultiplier - 1))
			initial.Damage *= 1 + critMod
			return initial
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			tickBase := damageDotTick + 0.04*spell.MeleeAttackPower()
			ticks := spell.CalcPeriodicDamage(sim, target, tickBase, spell.OutcomeExpectedMagicAlwaysHit)
			return ticks
		},
	}
}

func (druid *Druid) CurrentRakeCost() float64 {
	return druid.Rake.ApplyCostModifiers(druid.Rake.DefaultCast.Cost)
}
