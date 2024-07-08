package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) ApplyRunes() {
	if warrior.HasRune(proto.WarriorRune_RuneFrenziedAssault) && warrior.MainHand().HandType == proto.HandType_HandTypeTwoHand {
		warrior.PseudoStats.MeleeSpeedMultiplier *= 1.2
	}

	if warrior.HasRune(proto.WarriorRune_RuneShieldMastery) && warrior.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield {
		warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
	}

	// Chest
	warrior.applyFlagellation()
	warrior.registerRagingBlow()
	warrior.applyBloodFrenzy()

	// Bracers
	warrior.registerRampage()
	warrior.applySwordAndBoard()
	warrior.applyWreckingCrew()

	// Gloves
	warrior.applySingleMindedFury()
	warrior.registerQuickStrike()

	// Belt
	warrior.applyFocusedRage()
	// Precise Timing is implemented on slam.go
	warrior.applyBloodSurge()

	// Pants
	warrior.applyConsumedByRage()

	warrior.applyTasteForBlood()

	// Endless Rage implemented on dps_warrior.go and protection_warrior.go

	// Gladiator implemented on stances.go
}

// You gain Rage from Physical damage taken as if you were wearing no armor.
func (warrior *Warrior) applyFlagellation() {
	if !warrior.HasRune(proto.WarriorRune_RuneFlagellation) {
		return
	}

	// TODO: Rage gain from hits
}

// Rend can now be used in Berserker stance, Rend's damage is increased by 100%, and Rend deals additional damage equal to 3% of your Attack Power each time it deals damage.
func (warrior *Warrior) applyBloodFrenzy() {
	if !warrior.HasRune(proto.WarriorRune_RuneBloodFrenzy) {
		return
	}

	originalRendCastCondition := warrior.Rend.ExtraCastCondition
	warrior.Rend.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return originalRendCastCondition(sim, target) || warrior.StanceMatches(BerserkerStance)
	}

	// Rend AP scaling implemented in rend.go
}

// Enrages you (activating abilities which require being Enraged) for 12 sec  after you exceed 60 Rage.
// In addition, Whirlwind also strikes with off-hand melee weapons while you are Enraged
func (warrior *Warrior) applyConsumedByRage() {
	if !warrior.HasRune(proto.WarriorRune_RuneConsumedByRage) {
		return
	}

	warrior.ConsumedByRageAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Enrage (Consumed by Rage)",
		ActionID: core.ActionID{SpellID: 425415},
		Duration: time.Second * 12,
	})

	warrior.ConsumedByRageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: 0})

	warrior.RegisterAura(core.Aura{
		Label:    "Consumed By Rage Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnRageChange: func(aura *core.Aura, sim *core.Simulation, metrics *core.ResourceMetrics) {
			// Refunding rage should not enable CBR
			if warrior.CurrentRage() < 60 || metrics.ActionID.OtherID == proto.OtherAction_OtherActionRefund {
				return
			}

			warrior.ConsumedByRageAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applyFocusedRage() {
	warrior.FocusedRageDiscount = core.TernaryFloat64(warrior.HasRune(proto.WarriorRune_RuneFocusedRage), 3.0, 0)
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

// Your melee critical hits Enrage you (activating abilities which require being Enraged), and increase Mortal Strike, Bloodthirst, and Shield Slam damage by 10% for 12 sec.
func (warrior *Warrior) applyWreckingCrew() {
	if !warrior.HasRune(proto.WarriorRune_RuneWreckingCrew) {
		return
	}

	affectedSpells := core.FilterSlice(
		[]*core.Spell{warrior.MortalStrike, warrior.Bloodthirst, warrior.ShieldSlam},
		func(spell *core.Spell) bool { return spell != nil },
	)

	warrior.WreckingCrewEnrageAura = warrior.RegisterAura(core.Aura{
		Label:    "Enrage (Wrecking Crew)",
		ActionID: core.ActionID{SpellID: 427066},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplier *= 1.1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.MeleeCritMultiplier /= 1.1
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Wrecking Crew Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.Outcome.Matches(core.OutcomeCrit) {
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

func (warrior *Warrior) applySingleMindedFury() {
	if !warrior.HasRune(proto.WarriorRune_RuneSingleMindedFury) {
		return
	}

	if warrior.HasRune(proto.WarriorRune_RuneSingleMindedFury) && warrior.HasMHWeapon() && warrior.HasOHWeapon() {
		warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
	}
}
