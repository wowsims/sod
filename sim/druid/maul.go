package druid

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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

func (druid *Druid) registerMaulSpell(realismICD *core.Cooldown) {
	// Add highest available rank for level.
	for rank := len(maulSpells) - 1; rank >= 0; rank-- {
		if druid.Level >= maulSpells[rank].level {
			config := druid.newMaulSpellConfig(maulSpells[rank])
			druid.Maul = druid.RegisterSpell(Bear, config)
			break
		}
	}

	druid.MaulQueue = druid.makeQueueSpellsAndAura(druid.Maul, realismICD)
}

func (druid *Druid) newMaulSpellConfig(maulRank MaulRankInfo) core.SpellConfig {
	flatBaseDamage := maulRank.damage
	rageCost := 15 - float64(druid.Talents.Ferocity)
	hasGore := druid.HasRune(proto.DruidRune_RuneHelmGore)

	switch druid.Ranged().ID {
	case IdolOfBrutality:
		rageCost -= 3
	}
	rageMetrics := druid.NewRageMetrics(core.ActionID{SpellID: maulRank.id})

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: maulRank.id},
		SpellSchool: core.SpellSchoolPhysical,
		SpellCode:   SpellCode_DruidMaul,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:       SpellFlagOmen | core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   rageCost,
			Refund: 0.8,
		},

		DamageMultiplierAdditive: 1 + .1*float64(druid.Talents.SavageFury),
		ThreatMultiplier:         1.75,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Need to specially deactivate CC here in case maul is cast simultaneously with another spell.
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			baseDamage := flatBaseDamage + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			targetCount := core.TernaryInt32(druid.FuryOfStormrageMaulCleave, 2, 1)
			numHits := min(targetCount, druid.Env.GetNumTargets())
			results := make([]*core.SpellResult, numHits)

			for idx := range results {
				results[idx] = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if results[idx].Landed() {
					if hasGore && sim.Proc(0.15, "Gore") {
						druid.AddRage(sim, 10.0, rageMetrics)
						druid.MangleBear.CD.Reset()
					}
				}
				target = sim.Environment.NextTargetUnit(target)
			}

			if !results[0].Landed() {
				spell.IssueRefund(sim)
			}

			if druid.curQueueAura != nil {
				druid.curQueueAura.Deactivate(sim)
			}
		},
	}
}

func (druid *Druid) makeQueueSpellsAndAura(srcSpell *DruidSpell, realismICD *core.Cooldown) *DruidSpell {
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
