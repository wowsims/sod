package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) ApplyRunes() {
	if warrior.GetMHWeapon() != nil && warrior.GetOHWeapon() != nil { // This check is to stop memory dereference error if unarmed
		if warrior.GetMHWeapon().HandType == proto.HandType_HandTypeMainHand || warrior.GetMHWeapon().HandType == proto.HandType_HandTypeOneHand &&
			warrior.GetOHWeapon().HandType == proto.HandType_HandTypeOffHand || warrior.GetOHWeapon().HandType == proto.HandType_HandTypeOneHand {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneSingleMindedFury), 1.1, 1)
		}
	}

	if warrior.GetMHWeapon() != nil { // This check is to stop memory dereference error if unarmed
		if warrior.GetMHWeapon().HandType == proto.HandType_HandTypeTwoHand {
			warrior.PseudoStats.MeleeSpeedMultiplier *= core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneFrenziedAssault), 1.2, 1)
		}
	}

	warrior.FocusedRageDiscount = core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneFocusedRage), 3.0, 0)

	warrior.applyBloodFrenzy()
	warrior.applyFlagellation()
	warrior.applyConsumedByRage()
	warrior.registerQuickStrike()
	warrior.registerRagingBlow()
	warrior.applyBloodSurge()

	// Endless Rage implemented on dps_warrior.go and protection_warrior.go
	// Precise Timing is implemented on slam.go

}

func (warrior *Warrior) applyBloodFrenzy() {
	if !warrior.HasRune(proto.WarriorRune_RuneBloodFrenzy) {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 412507})

	warrior.RegisterAura(core.Aura{
		Label:    "Blood Frenzy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagBleed) {
				return
			}

			warrior.AddRage(sim, 3, rageMetrics)
		},
	})
}

func (warrior *Warrior) applyFlagellation() {
	if !warrior.HasRune(proto.WarriorRune_RuneFlagellation) {
		return
	}

	flagellationAura := warrior.RegisterAura(core.Aura{
		Label:    "Flagellation",
		ActionID: core.ActionID{SpellID: 402877},
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.25
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.25
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Flagellation Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != warrior.Bloodrage && spell != warrior.BerserkerRage {
				return
			}

			flagellationAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applyConsumedByRage() {
	if !warrior.HasRune(proto.WarriorRune_RuneConsumedByRage) {
		return
	}

	warrior.ConsumedByRageAura = warrior.RegisterAura(core.Aura{
		Label:     "Consumed By Rage",
		ActionID:  core.ActionID{SpellID: 425418},
		Duration:  time.Second * 10,
		MaxStacks: 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
			warrior.Above80RageCBRActive = true
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Above80RageCBRActive = true
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.1
			warrior.Above80RageCBRActive = false
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Consumed By Rage Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Above80RageCBRActive = false
			aura.Activate(sim)
		},
		OnRageChange: func(aura *core.Aura, sim *core.Simulation, metrics *core.ResourceMetrics) {
			// Refunding rage should not enable CBR
			if warrior.CurrentRage() < 80 || metrics.ActionID.OtherID == proto.OtherAction_OtherActionRefund {
				warrior.Above80RageCBRActive = false
				return
			}

			if warrior.Above80RageCBRActive {
				return
			}

			warrior.ConsumedByRageAura.Activate(sim)
			warrior.ConsumedByRageAura.SetStacks(sim, 12)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !warrior.ConsumedByRageAura.IsActive() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				warrior.ConsumedByRageAura.RemoveStack(sim)
			}
		},
	})
}

func (warrior *Warrior) applyBloodSurge() {
	if !warrior.HasRune(proto.WarriorRune_RuneBloodSurge) {
		return
	}

	warrior.BloodSurgeAura = warrior.RegisterAura(core.Aura{
		Label:    "Blood Surge Proc",
		ActionID: core.ActionID{SpellID: 413399},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Slam.DefaultCast.CastTime = 0
			warrior.Slam.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Slam.DefaultCast.CastTime = 1500 * time.Millisecond
			warrior.Slam.CostMultiplier -= 1
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warrior.Slam { // removed even if slam doesn't land
				aura.Deactivate(sim)
			}
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Blood Surge",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if !spell.Flags.Matches(SpellFlagBloodSurge) {
				return
			}

			if sim.RandomFloat("Blood Surge") > 0.3 {
				return
			}

			warrior.BloodSurgeAura.Activate(sim)
		},
	})
}
