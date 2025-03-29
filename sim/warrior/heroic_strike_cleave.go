package warrior

import (
	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerHeroicStrikeSpell(realismICD *core.Cooldown) {
	flatDamageBonus := map[int32]float64{
		25: 44,
		40: 80,
		50: 111,
		60: core.TernaryFloat64(core.IncludeAQ, 157, 138),
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 1608,
		40: 11565,
		50: 11566,
		60: core.TernaryInt32(core.IncludeAQ, 25286, 11567),
	}[warrior.Level]

	// No known equation
	threat := map[int32]float64{
		25: 68,  //guess
		40: 103, //guess
		50: 120,
		60: core.TernaryFloat64(core.IncludeAQ, 173, 145),
	}[warrior.Level]

	warrior.HeroicStrike = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_WarriorHeroicStrike,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.ImprovedHeroicStrike),
			Refund: 0.8,
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  threat,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
	warrior.HeroicStrikeQueue = warrior.makeQueueSpellsAndAura(warrior.HeroicStrike, realismICD)
}

func (warrior *Warrior) registerCleaveSpell(realismICD *core.Cooldown) {
	flatDamageBonus := map[int32]float64{
		25: 5,
		40: 18,
		50: 32,
		60: 50,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 845,
		40: 11608,
		50: 11609,
		60: 20569,
	}[warrior.Level]

	threat := map[int32]float64{
		25: 20, //guess
		40: 60, //guess
		50: 80,
		60: 100,
	}[warrior.Level]

	flatDamageBonus *= []float64{1, 1.4, 1.8, 2.2}[warrior.Talents.ImprovedCleave]

	warrior.CleaveTargetCount += 2
	targetCount := warrior.Env.GetNumTargets()

	warrior.Cleave = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		ClassSpellMask: ClassSpellMask_WarriorCleave,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost: 20,
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  threat,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for i := int32(0); i < min(targetCount, warrior.CleaveTargetCount); i++ {
				baseDamage := flatDamageBonus + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(target)
				if curTarget == target {
					break
				}
			}

			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
	warrior.CleaveQueue = warrior.makeQueueSpellsAndAura(warrior.Cleave, realismICD)
}

func (warrior *Warrior) makeQueueSpellsAndAura(srcSpell *WarriorSpell, realismICD *core.Cooldown) *WarriorSpell {
	isQueueQueued := false

	queueAura := warrior.RegisterAura(core.Aura{
		Label:    "HS/Cleave Queue Aura-" + srcSpell.ActionID.String(),
		ActionID: srcSpell.ActionID.WithTag(1),
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			isQueueQueued = false
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
			warrior.PseudoStats.DisableDWMissPenalty = true
			warrior.curQueueAura = aura
			warrior.curQueuedAutoSpell = srcSpell
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DisableDWMissPenalty = false
			warrior.curQueueAura = nil
			warrior.curQueuedAutoSpell = nil
		},
	})

	queueSpell := warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID: srcSpell.ActionID.WithTag(1),
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagCastTimeNoGCD,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.curQueueAura == nil &&
				!isQueueQueued &&
				warrior.CurrentRage() >= srcSpell.DefaultCast.Cost &&
				!warrior.IsCasting(sim) &&
				realismICD.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if realismICD.IsReady(sim) {
				isQueueQueued = true
				realismICD.Use(sim)
				sim.AddPendingAction(&core.PendingAction{
					NextActionAt: sim.CurrentTime + realismICD.Duration,
					OnAction: func(sim *core.Simulation) {
						queueAura.Activate(sim)
						isQueueQueued = false
					},
				})
			}
		},
	})

	return queueSpell
}

func (warrior *Warrior) TryHSOrCleave(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !warrior.curQueueAura.IsActive() {
		return mhSwingSpell
	}

	if !warrior.curQueuedAutoSpell.CanCast(sim, warrior.CurrentTarget) {
		warrior.curQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	return warrior.curQueuedAutoSpell.Spell
}
