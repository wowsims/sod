package druid

import (
	"github.com/wowsims/sod/sim/core"
)

type MaulRankInfo struct {
	id     int32
	level  int32
	damage float64
}

var maulSpells = []MaulRankInfo{
	{
		id:     6807,
		level:  10,
		damage: 18.0,
	},
	{
		id:     6808,
		level:  18,
		damage: 27.0,
	},
	{
		id:     6809,
		level:  26,
		damage: 37.0,
	},
	{
		id:     8972,
		level:  34,
		damage: 49.0,
	},
	{
		id:     9745,
		level:  42,
		damage: 71.0,
	},
	{
		id:     9880,
		level:  50,
		damage: 101.0,
	},
	{
		id:     9881,
		level:  58,
		damage: 128.0,
	},
}

func (druid *Druid) registerMaulSpell() {
	// Add highest available rank for level.
	for rank := len(maulSpells) - 1; rank >= 0; rank-- {
		if druid.Level >= maulSpells[rank].level {
			config := druid.newMaulSpellConfig(maulSpells[rank])
			druid.Maul = druid.RegisterSpell(Bear, config)
			break
		}
	}

	druid.makeQueueSpellsAndAura(druid.Maul)
}

func (druid *Druid) newMaulSpellConfig(maulRank MaulRankInfo) core.SpellConfig {
	flatBaseDamage := maulRank.damage
	actionID := core.ActionID{SpellID: maulRank.id}

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidMaul,
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:          SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1.75,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Need to specially deactivate CC here in case maul is cast simultaneously with another spell.
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			baseDamage := flatBaseDamage + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			targetCount := core.TernaryInt32(druid.FuryOfStormrageMaulCleave, 2, 1)
			numHits := min(targetCount, druid.Env.GetNumTargets())

			for i := int32(0); i < numHits; i++ {
				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if !result.Landed() && i == 0 {
					spell.IssueRefund(sim)
				}

				target = sim.Environment.NextTargetUnit(target)
			}

			if druid.curQueueAura != nil {
				druid.curQueueAura.Deactivate(sim)
			}
		},
	}
}

func (druid *Druid) makeQueueSpellsAndAura(srcSpell *DruidSpell) *DruidSpell {
	realismICD := core.Cooldown{
		Timer:    druid.NewTimer(),
		Duration: core.SpellBatchWindow * 10,
	}

	queueAura := druid.RegisterAura(core.Aura{
		Label:    "Maul Queue Aura-" + srcSpell.ActionID.String(),
		ActionID: srcSpell.ActionID.WithTag(1),
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if druid.curQueueAura != nil {
				druid.curQueueAura.Deactivate(sim)
			}
			druid.curQueueAura = aura
			druid.curQueuedAutoSpell = srcSpell
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.curQueueAura = nil
			druid.curQueuedAutoSpell = nil
		},
	})

	queueSpell := druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID: srcSpell.ActionID.WithTag(1),
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagCastTimeNoGCD,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.curQueueAura == nil &&
				druid.CurrentRage() >= srcSpell.DefaultCast.Cost &&
				!druid.IsCasting(sim) &&
				realismICD.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if realismICD.IsReady(sim) {
				realismICD.Use(sim)
				sim.AddPendingAction(&core.PendingAction{
					NextActionAt: sim.CurrentTime + realismICD.Duration,
					OnAction: func(sim *core.Simulation) {
						queueAura.Activate(sim)
					},
				})
			}
		},
	})

	return queueSpell
}
