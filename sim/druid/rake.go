package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

type RakeRankInfo struct {
	id            int32
	level         int32
	initialDamage float64
	dotTickDamage float64
}

var rakeSpells = []RakeRankInfo{
	{
		id:            1822,
		level:         24,
		initialDamage: 19.0,
		dotTickDamage: 13.0,
	},
	{
		id:            1823,
		level:         34,
		initialDamage: 29.0,
		dotTickDamage: 19.0,
	},
	{

		id:            1824,
		level:         44,
		initialDamage: 43.0,
		dotTickDamage: 25.0,
	},
	{

		id:            9904,
		level:         54,
		initialDamage: 58.0,
		dotTickDamage: 32.0,
	},
}

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
	damageDotTick := rakeRank.dotTickDamage * baseDmgMultiplier

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: rakeRank.id},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | core.SpellFlagAPL,

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
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rake",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				damage := damageDotTick + 0.04*dot.Spell.MeleeAttackPower()
				dot.Snapshot(target, damage, 0, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamageNew(sim, target, 0, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damageInitial + 0.04*spell.MeleeAttackPower()
			if druid.BleedCategories.Get(target).AnyActive() {
				baseDamage *= 1.3
			}

			result := spell.CalcAndDealDamageNew(sim, target, baseDamage, 0, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := damageInitial + 0.04*spell.MeleeAttackPower()
			initial := spell.CalcPeriodicDamageNew(sim, target, baseDamage, 0, spell.OutcomeExpectedMagicAlwaysHit)

			attackTable := spell.Unit.AttackTables[target.UnitIndex][spell.CastType]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := critChance * (spell.CritMultiplier(attackTable) - 1)
			initial.Damage *= 1 + critMod
			return initial
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			tickBase := damageDotTick + 0.04*spell.MeleeAttackPower()
			ticks := spell.CalcPeriodicDamageNew(sim, target, tickBase, 0, spell.OutcomeExpectedMagicAlwaysHit)
			return ticks
		},
	}
}

func (druid *Druid) CurrentRakeCost() float64 {
	return druid.Rake.ApplyCostModifiers(druid.Rake.DefaultCast.Cost)
}
