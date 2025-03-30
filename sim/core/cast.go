package core

import (
	"fmt"
	"time"
)

// A cast corresponds to any action which causes the in-game castbar to be
// shown, and activates the GCD. Note that a cast can also be instant, i.e.
// the effects are applied immediately even though the GCD is still activated.

// Callback for when a cast is finished, i.e. when the in-game castbar reaches full.
type OnCastComplete func(aura *Aura, sim *Simulation, spell *Spell)

type Hardcast struct {
	Expires    time.Duration
	ActionID   ActionID
	OnComplete func(*Simulation, *Unit)
	Target     *Unit
}

// Input for constructing the CastSpell function for a spell.
type CastConfig struct {
	// Default cast values with all static effects applied.
	DefaultCast Cast

	// Dynamic modifications for each cast.
	ModifyCast func(*Simulation, *Spell, *Cast)

	// Ignores haste when calculating the cast time for this cast.
	// Automatically set if GCD and cast times are all 0, e.g. for empty casts.
	IgnoreHaste bool

	CD       Cooldown
	SharedCD Cooldown

	CastTime func(spell *Spell) time.Duration
}

type Cast struct {
	// Amount of resource that will be consumed by this cast.
	Cost float64

	// The length of time the GCD will be on CD as a result of this cast.
	GCD time.Duration

	// The minimum length of time for the GCD. Can be left out to use the default of 1s
	GCDMin time.Duration

	// The amount of time between the call to spell.Cast() and when the spell
	// effects are invoked.
	CastTime time.Duration
}

func (cast *Cast) EffectiveTime() time.Duration {
	gcd := max(0, cast.GCD)
	if cast.GCD > 0 {
		if cast.GCDMin != 0 {
			gcd = max(cast.GCDMin, gcd)
		} else {
			gcd = max(GCDMin, gcd)
		}
	}
	return max(gcd, cast.CastTime)
}

type CastFunc func(*Simulation, *Unit)
type CastSuccessFunc func(*Simulation, *Unit) bool

func (spell *Spell) castFailureHelper(sim *Simulation, message string, vals ...any) bool {
	if sim.CurrentTime < 0 && spell.Unit.Rotation != nil {
		spell.Unit.Rotation.ValidationWarning(fmt.Sprintf(spell.ActionID.String()+" failed to cast: "+message, vals...))
	} else {
		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, fmt.Sprintf(spell.ActionID.String()+" failed to cast: "+message, vals...))
		}
	}
	return false
}

func (unit *Unit) applySpellPushback() {
	unit.RegisterAura(Aura{
		Label:    "Spell Pushback",
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			//No pushback for bosses/NPCs
			if unit.Type == EnemyUnit {
				return
			}

			if !result.Landed() {
				return
			}

			if result.Damage <= 0 {
				return
			}

			if !spell.ProcMask.Matches(ProcMaskDirect) {
				return
			}

			if hc := aura.Unit.Hardcast; aura.Unit.IsCasting(sim) {
				hcSpell := aura.Unit.GetSpell(hc.ActionID)
				// Caster avoided the pushback
				if sim.Roll(0, 1.0) <= hcSpell.PushbackReduction {
					return
				}

				// Do spell pushback
				pushback := DurationFromSeconds(sim.Roll(0.5, 1.0) * unit.PseudoStats.SpellPushbackMultiplier)

				if hcSpell.Flags.Matches(SpellFlagChanneled) {
					newExpires := max(sim.CurrentTime, hc.Expires-pushback)
					if sim.Log != nil {
						aura.Unit.Log(sim, "Unit Hardcast shortened by %s due to spell hit taken, will now occur at %s", pushback, newExpires)
					}

					// Update Dot if present
					if hcDot := hcSpell.CurDot(); hcDot != nil {
						hcDot.UpdateExpires(sim, newExpires)
					}

					aura.Unit.Hardcast.Expires = newExpires
					hcSpell.SpellMetrics[aura.Unit.CurrentTarget.UnitIndex].TotalCastTime -= pushback

				} else {
					if sim.Log != nil {
						aura.Unit.Log(sim, "Unit Hardcast extended by %s due to spell hit taken, will now occur at %s", pushback, hc.Expires+pushback)
					}

					aura.Unit.Hardcast.Expires += pushback
					hcSpell.SpellMetrics[aura.Unit.CurrentTarget.UnitIndex].TotalCastTime += pushback
				}

				// Update GCDTimer
				aura.Unit.SetGCDTimer(sim, aura.Unit.Hardcast.Expires)

				// Update Swing timer
				aura.Unit.AutoAttacks.StopMeleeUntil(sim, aura.Unit.Hardcast.Expires, false)
			}
		},
	})
}

func (spell *Spell) makeCastFunc(config CastConfig) CastSuccessFunc {
	return func(sim *Simulation, target *Unit) bool {
		spell.CurCast = spell.DefaultCast

		if config.ModifyCast != nil {
			config.ModifyCast(sim, spell, &spell.CurCast)
			if spell.CurCast.Cost != spell.DefaultCast.Cost {
				// Costs need to be modified using the unit and spell multipliers, so that
				// their affects are also visible in the spell.CanCast() function, which
				// does not invoke ModifyCast.
				panic("May not modify cost in ModifyCast!")
			}
		}

		if spell.Flags.Matches(SpellFlagSwapped) {
			return spell.castFailureHelper(sim, "spell attached to an un-equipped item")
		}

		if spell.ExtraCastCondition != nil {
			if !spell.ExtraCastCondition(sim, target) {
				return spell.castFailureHelper(sim, "extra spell condition")
			}
		}

		if spell.Cost != nil {
			if !spell.Cost.MeetsRequirement(sim, spell) {
				return spell.castFailureHelper(sim, spell.Cost.CostFailureReason(sim, spell))
			}
		}

		if !config.IgnoreHaste {
			if spell.AllowGCDHasteScaling {
				spell.CurCast.GCD = max(0, spell.Unit.ApplyCastSpeed(spell.CurCast.GCD)).Round(time.Millisecond)
			}

			// Vanilla has no natural GCD reduction besides abilities with 1s GCDs
			// spell.CurCast.GCD = spell.Unit.ApplyFlatCastSpeed(spell.CurCast.GCD)
			spell.CurCast.CastTime = config.CastTime(spell)
		}

		if config.CD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.CD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on cooldown for %s, curTime = %s", spell.CD.TimeToReady(sim), sim.CurrentTime)
			}

			spell.CD.Set(sim.CurrentTime + spell.CurCast.CastTime + spell.CD.GetCurrentDuration())
		}

		if config.SharedCD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.SharedCD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on shared cooldown for %s, curTime = %s", spell.SharedCD.TimeToReady(sim), sim.CurrentTime)
			}

			spell.SharedCD.Set(sim.CurrentTime + spell.CurCast.CastTime + spell.SharedCD.Duration)
		}

		// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
		if spell.CurCast.GCD != 0 && !spell.Unit.GCD.IsReady(sim) {
			return spell.castFailureHelper(sim, "GCD on cooldown for %s, curTime = %s", spell.Unit.GCD.TimeToReady(sim), sim.CurrentTime)
		}

		if hc := spell.Unit.Hardcast; spell.Unit.IsCasting(sim) && !spell.Flags.Matches(SpellFlagCastWhileCasting) {
			return spell.castFailureHelper(sim, "casting %v for %s, curTime = %s", hc.ActionID, hc.Expires-sim.CurrentTime, sim.CurrentTime)
		}

		if dot := spell.Unit.ChanneledDot; spell.Unit.IsChanneling(sim) && !spell.Flags.Matches(SpellFlagCastWhileChanneling) && (spell.Unit.Rotation.interruptChannelIf == nil || !spell.Unit.Rotation.interruptChannelIf.GetBool(sim)) {
			return spell.castFailureHelper(sim, "channeling %v for %s, curTime = %s", dot.ActionID, dot.expires-sim.CurrentTime, sim.CurrentTime)
		}

		if effectiveTime := spell.CurCast.EffectiveTime(); effectiveTime != 0 {
			if spell.Flags.Matches(SpellFlagCastTimeNoGCD) {
				effectiveTime = max(effectiveTime, spell.Unit.GCD.TimeToReady(sim))
			}
			// do not add channeled time here as they have variable cast length
			// cast time for channels is handled in dot.OnExpire
			if !spell.Flags.Matches(SpellFlagChanneled) {
				spell.SpellMetrics[target.UnitIndex].TotalCastTime += effectiveTime
			}
			spell.Unit.SetGCDTimer(sim, sim.CurrentTime+effectiveTime)
		}

		if (spell.CurCast.CastTime > 0) && spell.Unit.IsMoving() {
			return spell.castFailureHelper(sim, "casting/channeling while moving not allowed!")
		}

		// Non melee casts
		if spell.Flags.Matches(SpellFlagResetAttackSwing) && spell.Unit.AutoAttacks.anyEnabled() {
			restartSwingAt := sim.CurrentTime + spell.CurCast.CastTime
			spell.Unit.AutoAttacks.StopMeleeUntil(sim, restartSwingAt, false)
			spell.Unit.AutoAttacks.StopRangedUntil(sim, restartSwingAt)
		}

		// Hardcasts
		if spell.CurCast.CastTime > 0 {
			if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
				spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
					spell.ActionID, max(0, spell.CurCast.Cost), spell.CurCast.CastTime, spell.CurCast.EffectiveTime())
			}

			spell.Unit.Hardcast = Hardcast{
				Expires:  sim.CurrentTime + spell.CurCast.CastTime,
				ActionID: spell.ActionID,
				OnComplete: func(sim *Simulation, target *Unit) {
					spell.LastCastAt = sim.CurrentTime

					if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
						spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
					}

					if spell.Cost != nil {
						if !spell.Cost.MeetsRequirement(sim, spell) {
							spell.castFailureHelper(sim, spell.Cost.CostFailureReason(sim, spell))
							return
						}
						spell.Cost.SpendCost(sim, spell)
					}

					spell.applyEffects(sim, target)

					if !spell.Flags.Matches(SpellFlagNoOnCastComplete | SpellFlagNoLifecycleCallbacks) {
						spell.Unit.OnCastComplete(sim, spell)
					}

					if !sim.Options.Interactive {
						spell.Unit.Rotation.DoNextAction(sim)
					}
				},
				Target: target,
			}

			spell.Unit.newHardcastAction(sim)

			return true
		}

		spell.LastCastAt = sim.CurrentTime

		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
				spell.ActionID, max(0, spell.CurCast.Cost), spell.CurCast.CastTime, spell.CurCast.EffectiveTime())
			spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
		}

		if spell.Cost != nil {
			spell.Cost.SpendCost(sim, spell)
		}

		spell.applyEffects(sim, target)

		if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
			spell.Unit.OnCastComplete(sim, spell)
		}

		return true
	}
}

func (spell *Spell) makeCastFuncSimple() CastSuccessFunc {
	return func(sim *Simulation, target *Unit) bool {
		if spell.Flags.Matches(SpellFlagSwapped) {
			return spell.castFailureHelper(sim, "spell attached to an un-equipped item")
		}

		if spell.ExtraCastCondition != nil {
			if !spell.ExtraCastCondition(sim, target) {
				return spell.castFailureHelper(sim, "extra spell condition")
			}
		}

		if spell.CD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.CD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on cooldown for %s, curTime = %s", spell.CD.TimeToReady(sim), sim.CurrentTime)
			}

			spell.CD.Set(sim.CurrentTime + spell.CD.GetCurrentDuration())
		}

		if spell.SharedCD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.SharedCD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on shared cooldown for %s, curTime = %s", spell.SharedCD.TimeToReady(sim), sim.CurrentTime)
			}

			spell.SharedCD.Set(sim.CurrentTime + spell.SharedCD.Duration)
		}

		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
				spell.ActionID, 0.0, "0s", "0s")
			spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
		}

		spell.applyEffects(sim, target)

		if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
			spell.Unit.OnCastComplete(sim, spell)
		}

		return true
	}
}

func (spell *Spell) makeCastFuncAutosOrProcs() CastSuccessFunc {
	return func(sim *Simulation, target *Unit) bool {
		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
				spell.ActionID, 0.0, "0s", "0s")
			spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
		}

		spell.applyEffects(sim, target)

		if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
			spell.Unit.OnCastComplete(sim, spell)
		}

		return true
	}
}
