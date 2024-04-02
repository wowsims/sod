package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) ApplyRunes() {
	// Apply runes here :)
	if rogue.HasRune(proto.RogueRune_RuneDeadlyBrew) {
		rogue.applyDeadlyBrewInstant()
		rogue.applyDeadlyBrewDeadly()
	}

	rogue.registerWaylayAura()
	rogue.registerMasterOfSubtlety()
	rogue.registerMainGaucheSpell()
	rogue.registerSaberSlashSpell()
	rogue.registerShivSpell()
	rogue.registerShadowstrikeSpell()
	rogue.registerMutilateSpell()
	rogue.registerEnvenom()
	rogue.registerShadowstep()
	rogue.registerShurikenTossSpell()
	rogue.registerQuickDrawSpell()
	rogue.registerBetweenTheEyes()
	rogue.registerPoisonedKnife()
	rogue.registerHonorAmongThieves()
	rogue.applyCombatPotency()
	rogue.applyFocusedAttacks()
	rogue.applyCarnage()
	rogue.applyUnfairAdvantage()
}

func (rogue *Rogue) applyCombatPotency() {
	if !rogue.HasRune(proto.RogueRune_RuneCombatPotency) {
		return
	}

	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 432292})

	rogue.RegisterAura(core.Aura{
		Label:    "Combat Potency",
		ActionID: energyMetrics.ActionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
				return
			}

			if sim.RandomFloat("Combat Potency") < 0.2 {
				rogue.AddEnergy(sim, 15, energyMetrics)
			}
		},
	})
}

func (rogue *Rogue) applyFocusedAttacks() {
	if !rogue.HasRune(proto.RogueRune_RuneFocusedAttacks) {
		return
	}

	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: int32(proto.RogueRune_RuneFocusedAttacks)})

	rogue.RegisterAura(core.Aura{
		Label:    "Focused Attacks",
		ActionID: energyMetrics.ActionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) || !result.DidCrit() {
				return
			}
			// TODO Check whether certain spells don't trigger this
			rogue.AddEnergy(sim, 2, energyMetrics)
		},
	})
}

func (rogue *Rogue) registerHonorAmongThieves() {
	if !rogue.HasRune(proto.RogueRune_RuneHonorAmongThieves) {
		return
	}

	comboMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: int32(proto.RogueRune_RuneHonorAmongThieves)})

	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: time.Second,
	}

	rogue.HonorAmongThieves = rogue.RegisterAura(core.Aura{
		Label:    "Honor Among Thieves",
		ActionID: comboMetrics.ActionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			// In an ideal party, you'd probably get up to 6 ability crits/s (Rate = 600).
			//  Survival Hunters, Enhancement Shamans, and Assassination Rogues are particularly good.
			if rogue.Options.HonorAmongThievesCritRate <= 0 {
				return
			}

			if rogue.Options.HonorAmongThievesCritRate > 2000 {
				rogue.Options.HonorAmongThievesCritRate = 2000 // limited, so performance doesn't suffer
			}

			rateToDuration := float64(time.Second) * 100 / float64(rogue.Options.HonorAmongThievesCritRate)

			pa := &core.PendingAction{}
			pa.OnAction = func(sim *core.Simulation) {
				rogue.tryHonorAmongThievesProc(sim, icd, comboMetrics)
				pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
				sim.AddPendingAction(pa)
			}
			pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
			sim.AddPendingAction(pa)
		},
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && !spell.ProcMask.Matches(core.ProcMaskWhiteHit) {
				rogue.tryHonorAmongThievesProc(sim, icd, comboMetrics)
			}
		},
		OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				rogue.tryHonorAmongThievesProc(sim, icd, comboMetrics)
			}
		},
	})
}

func (rogue *Rogue) tryHonorAmongThievesProc(sim *core.Simulation, icd core.Cooldown, metrics *core.ResourceMetrics) {
	if icd.IsReady(sim) {
		rogue.AddComboPoints(sim, 1, metrics)
		icd.Use(sim)
	}
}

// Apply the effects of the Cut to the Chase talent
func (rogue *Rogue) ApplyCutToTheChase(sim *core.Simulation) {
	// Rune check is done in envenom.go and eviscerate.go
	if rogue.SliceAndDiceAura.IsActive() {
		rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
		rogue.SliceAndDiceAura.Activate(sim)
	}
}
