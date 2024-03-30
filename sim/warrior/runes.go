package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) ApplyRunes() {
	if warrior.HasRune(proto.WarriorRune_RuneSingleMindedFury) && warrior.HasMHWeapon() && warrior.HasOHWeapon() {
		warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
	}

	if warrior.HasRune(proto.WarriorRune_RuneFrenziedAssault) && warrior.MainHand().HandType == proto.HandType_HandTypeTwoHand {
		warrior.PseudoStats.MeleeSpeedMultiplier *= 1.2
	}

	if warrior.HasRune(proto.WarriorRune_RuneShieldMastery) && warrior.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield {
		warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
	}

	warrior.FocusedRageDiscount = core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneFocusedRage), 3.0, 0)

	warrior.applyBloodFrenzy()
	warrior.applyFlagellation()
	warrior.applyConsumedByRage()
	warrior.registerQuickStrike()
	warrior.registerRagingBlow()
	warrior.applyBloodSurge()
	warrior.registerRampage()
	warrior.applyTasteForBlood()
	warrior.applyWreckingCrew()
	warrior.applySwordAndBoard()

	// Endless Rage implemented on dps_warrior.go and protection_warrior.go
	// Precise Timing is implemented on slam.go
	// Gladiator implemented on stances.go

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

	warrior.ConsumedByRageAura = warrior.GetOrRegisterAura(core.Aura{
		Label:     "Enrage 10%",
		ActionID:  core.ActionID{SpellID: 425415},
		Duration:  time.Second * 12,
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

	warrior.ConsumedByRageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: 10})

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

			if spell.ProcMask.Matches(core.ProcMaskMelee) {
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
			if warrior.Slam != nil {
				warrior.Slam.DefaultCast.CastTime = 0
				warrior.Slam.CostMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.Slam != nil {
				warrior.Slam.DefaultCast.CastTime = 1500 * time.Millisecond
				warrior.Slam.CostMultiplier += 1
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warrior.Slam != nil && spell == warrior.Slam { // removed even if slam doesn't land
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

func (warrior *Warrior) applyTasteForBlood() {
	if !warrior.HasRune(proto.WarriorRune_RuneTasteForBlood) {
		return
	}

	icd := core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Millisecond * 580,
	}

	warrior.RegisterAura(core.Aura{
		Label:    "Taste for Blood",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != warrior.Rend {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)
			warrior.OverpowerAura.Duration = time.Second * 9
			warrior.OverpowerAura.Activate(sim)
			warrior.OverpowerAura.Duration = time.Second * 5
		},
	})
}

func (warrior *Warrior) applyWreckingCrew() {
	if !warrior.HasRune(proto.WarriorRune_RuneWreckingCrew) {
		return
	}

	warrior.WreckingCrewEnrageAura = warrior.RegisterAura(core.Aura{
		Label:    "Enrage Wrecking Crew",
		ActionID: core.ActionID{SpellID: 427066},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.MeleeCritMultiplier *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.MeleeCritMultiplier /= 1.1
		},
	})

	warrior.WreckingCrewEnrageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: 1})

	warrior.RegisterAura(core.Aura{
		Label:    "Wrecking Crew",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			warrior.WreckingCrewEnrageAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applySwordAndBoard() {
	if !warrior.HasRune(proto.WarriorRune_RuneSwordAndBoard) || !warrior.Talents.ShieldSlam {
		return
	}

	sabAura := warrior.GetOrRegisterAura(core.Aura{
		Label:    "Sword And Board",
		ActionID: core.ActionID{SpellID: int32(proto.WarriorRune_RuneSwordAndBoard)},
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ShieldSlam.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ShieldSlam.CostMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warrior.ShieldSlam {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
		Label: "Sword And Board Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell != warrior.Revenge && spell != warrior.Devastate {
				return
			}

			if sim.RandomFloat("Sword And Board") < 0.3 {
				sabAura.Activate(sim)
				warrior.ShieldSlam.CD.Reset()
			}
		},
	}))
}
