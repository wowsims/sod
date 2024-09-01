package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
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
	// rogue.registerShivSpell()
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
	rogue.registerBladeDance()
	rogue.applyJustAFleshWound()
	rogue.applyRollingWithThePunches()
	rogue.registerCutthroat()
	rogue.registerBlunderbussSpell()
	rogue.registerFanOfKnives()
	rogue.registerCrimsonTempestSpell()
}

const SlaughterFromTheShadowsDamageMultiplier = 1.60
const SlaughterFromTheShadowsCostReduction = 30.0

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
			isFoKOH := false

			if spell.ActionID.SpellID == 409240 && spell.ActionID.Tag == 2 {
				isFoKOH = true
			}

			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOH) || isFoKOH {
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
			isFoKOH := false

			if spell.ActionID.SpellID == 409240 && spell.ActionID.Tag == 2 {
				isFoKOH = true
			}

			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged|core.ProcMaskTriggerInstant) || !result.DidCrit() || isFoKOH {
				return
			}
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
	refreshSlice := rogue.SliceAndDiceAura.IsActive()
	refreshBladeDance := rogue.BladeDanceAura.IsActive()
	// Refresh the lowest duration of SnD or Blade Dance
	if refreshBladeDance && refreshSlice {
		if rogue.SliceAndDiceAura.RemainingDuration(sim) > rogue.BladeDanceAura.RemainingDuration(sim) {
			refreshSlice = false
		} else {
			refreshBladeDance = false
		}
	}
	if refreshSlice {
		rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
		rogue.SliceAndDiceAura.Activate(sim)
	} else if refreshBladeDance {
		rogue.BladeDanceAura.Duration = rogue.bladeDanceDurations[5]
		rogue.BladeDanceAura.Activate(sim)
	}
}

func (rogue *Rogue) registerBladeDance() {
	if !rogue.HasRune(proto.RogueRune_RuneBladeDance) {
		return
	}

	justAFleshWound := rogue.HasRune(proto.RogueRune_RuneJustAFleshWound)

	rogue.bladeDanceDurations = [6]time.Duration{
		0,
		time.Duration(time.Second * 14),
		time.Duration(time.Second * 18),
		time.Duration(time.Second * 22),
		time.Duration(time.Second * 26),
		time.Duration(time.Second * 30),
	}

	cachedBonusAP := 0.0
	apProcAura := rogue.RegisterAura(core.Aura{
		Label:    "Defender's Resolve",
		ActionID: core.ActionID{SpellID: 462230},
		Duration: rogue.bladeDanceDurations[5],
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			cachedBonusAP = 4 * max(rogue.GetStat(stats.Defense), 0)
			rogue.AddStatDynamic(sim, stats.AttackPower, cachedBonusAP)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AddStatDynamic(sim, stats.AttackPower, -cachedBonusAP)
		},
	})

	rogue.BladeDanceAura = rogue.RegisterAura(core.Aura{
		Label:    "Blade Dance",
		ActionID: core.ActionID{SpellID: int32(proto.RogueRune_RuneBladeDance)},
		Duration: rogue.bladeDanceDurations[5],
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AddStatDynamic(sim, stats.Parry, 10*core.ParryRatingPerParryChance)
			if justAFleshWound {
				rogue.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= 0.8
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AddStatDynamic(sim, stats.Parry, -10*core.ParryRatingPerParryChance)
			if justAFleshWound {
				rogue.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= 0.8
			}
		},
	})

	rogue.BladeDance = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.RogueRune_RuneBladeDance)},
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:   25,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BladeDanceAura.Duration = rogue.bladeDanceDurations[rogue.ComboPoints()]
			rogue.BladeDanceAura.Activate(sim)
			apProcAura.Activate(sim)
			rogue.SpendComboPoints(sim, spell)
		},
	})
}

func (rogue *Rogue) applyJustAFleshWound() {
	if !rogue.HasRune(proto.RogueRune_RuneJustAFleshWound) {
		return
	}
	// Mod threat
	// TODO: Confirm threat mod
	rogue.PseudoStats.ThreatMultiplier *= 2.112

	// Blade Dance 20% Physical DR - Added in registerBladeDance()

	// -6% to be critically hit
	rogue.PseudoStats.ReducedCritTakenChance += 6

	// Replace Feint with Tease
	// TODO: Warrior sim from wrath did not implement it. May implement later

	// Shuriken Toss and Poisoned Knife gain 50% threat mod
	// Implemented in the relevant files
}

func (rogue *Rogue) applyRollingWithThePunches() {
	if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
		return
	}

	statDeps := make([]*stats.StatDependency, 11) // 10 stacks + zero condition
	for i := 1; i < 6; i++ {
		statDeps[i] = rogue.NewDynamicMultiplyStat(stats.Health, 1.0+.06*float64(i))
	}

	rogue.RollingWithThePunchesProcAura = rogue.RegisterAura(core.Aura{
		Label:     "Rolling with the Punches Proc",
		ActionID:  core.ActionID{SpellID: 400015},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if oldStacks != 0 {
				aura.Unit.DisableDynamicStatDep(sim, statDeps[oldStacks])
			}
			if newStacks != 0 {
				aura.Unit.EnableDynamicStatDep(sim, statDeps[newStacks])
			}
		},
	})

	rogue.RollingWithThePunchesAura = rogue.RegisterAura(core.Aura{
		Label:           "Rolling with the Punches",
		ActionID:        core.ActionID{SpellID: int32(proto.RogueRune_RuneRollingWithThePunches)},
		ActionIDForProc: core.ActionID{SpellID: int32(proto.RogueRune_RuneRollingWithThePunches) - 1},
		Duration:        core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee|core.ProcMaskRanged) && result.Outcome.Matches(core.OutcomeDodge|core.OutcomeParry) {
				if rogue.RollingWithThePunchesProcAura.IsActive() {
					rogue.RollingWithThePunchesProcAura.AddStack(sim)
				} else {
					rogue.RollingWithThePunchesProcAura.Activate(sim)
				}
			}
		},
	})
}

func (rogue *Rogue) rollCutthroat(sim *core.Simulation) {
	if sim.RandomFloat("Cutthroat") < 0.15 {
		rogue.CutthroatProcAura.Activate(sim)
	}
}

func (rogue *Rogue) registerCutthroat() {
	if !rogue.HasRune(proto.RogueRune_RuneCutthroat) {
		return
	}

	rogue.CutthroatProcAura = rogue.RegisterAura(core.Aura{
		Label:    "Cutthroat",
		ActionID: core.ActionID{SpellID: 462707},
		Duration: time.Second * 10,
	})
}
