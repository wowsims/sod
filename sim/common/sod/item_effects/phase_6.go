package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	// Brood of Nozdormu Reputations Rings
	SignetRingBronzeDominatorR5   = 234034
	SignetRingBronzeDominatorR4   = 234030
	SignetRingBronzeDominatorR3   = 234026
	SignetRingBronzeDominatorR2   = 234021
	SignetRingBronzeDominatorR1   = 234017
	SignetRingBronzeInvokerR5     = 234032
	SignetRingBronzeInvokerR4     = 234028
	SignetRingBronzeInvokerR3     = 234024
	SignetRingBronzeInvokerR2     = 234020
	SignetRingBronzeInvokerR1     = 234016
	SignetRingBronzeFlamekeeperR5 = 234964
	SignetRingBronzeFlamekeeperR4 = 234965
	SignetRingBronzeFlamekeeperR3 = 234966
	SignetRingBronzeFlamekeeperR2 = 234967
	SignetRingBronzeFlamekeeperR1 = 234968
	SignetRingBronzePreserverR5   = 234033
	SignetRingBronzePreserverR4   = 234029
	SignetRingBronzePreserverR3   = 234025
	SignetRingBronzePreserverR2   = 234023
	SignetRingBronzePreserverR1   = 234019
	SignetRingBronzeProtectorR5   = 234035
	SignetRingBronzeProtectorR4   = 234031
	SignetRingBronzeProtectorR3   = 234027
	SignetRingBronzeProtectorR2   = 234022
	SignetRingBronzeProtectorR1   = 234018
	SignetRingBronzeSubjugatorR5  = 234436
	SignetRingBronzeSubjugatorR4  = 234437
	SignetRingBronzeSubjugatorR3  = 234438
	SignetRingBronzeSubjugatorR2  = 234439
	SignetRingBronzeSubjugatorR1  = 234440
	SignetRingBronzeConquerorR5   = 234202
	SignetRingBronzeConquerorR4   = 234201
	SignetRingBronzeConquerorR3   = 234200
	SignetRingBronzeConquerorR2   = 234199
	SignetRingBronzeConquerorR1   = 234198
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Rings
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=234198/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeConquerorR5, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeConquerorR4, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeConquerorR3, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeConquerorR2, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeConquerorR1, TimeswornStrikeAura)

	// https://www.wowhead.com/classic/item=234034/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeDominatorR5, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeDominatorR4, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeDominatorR3, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeDominatorR2, TimeswornStrikeAura)
	core.NewItemEffect(SignetRingBronzeDominatorR1, TimeswornStrikeAura)

	// https://www.wowhead.com/classic/item=234964/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeFlamekeeperR5, TimeswornPyromancyAura)
	core.NewItemEffect(SignetRingBronzeFlamekeeperR4, TimeswornPyromancyAura)
	core.NewItemEffect(SignetRingBronzeFlamekeeperR3, TimeswornPyromancyAura)
	core.NewItemEffect(SignetRingBronzeFlamekeeperR2, TimeswornPyromancyAura)
	core.NewItemEffect(SignetRingBronzeFlamekeeperR1, TimeswornPyromancyAura)

	// https://www.wowhead.com/classic/item=234032/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeInvokerR5, TimeswornSpellAura)
	core.NewItemEffect(SignetRingBronzeInvokerR4, TimeswornSpellAura)
	core.NewItemEffect(SignetRingBronzeInvokerR3, TimeswornSpellAura)
	core.NewItemEffect(SignetRingBronzeInvokerR2, TimeswornSpellAura)
	core.NewItemEffect(SignetRingBronzeInvokerR1, TimeswornSpellAura)

	// https://www.wowhead.com/classic/item=234032/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzePreserverR5, TimewornHealing)
	core.NewItemEffect(SignetRingBronzePreserverR4, TimewornHealing)
	core.NewItemEffect(SignetRingBronzePreserverR3, TimewornHealing)
	core.NewItemEffect(SignetRingBronzePreserverR2, TimewornHealing)
	core.NewItemEffect(SignetRingBronzePreserverR1, TimewornHealing)

	// https://www.wowhead.com/classic/item=234035/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeProtectorR5, TimeswornExpertiseAura)
	core.NewItemEffect(SignetRingBronzeProtectorR4, TimeswornExpertiseAura)
	core.NewItemEffect(SignetRingBronzeProtectorR3, TimeswornExpertiseAura)
	core.NewItemEffect(SignetRingBronzeProtectorR2, TimeswornExpertiseAura)
	core.NewItemEffect(SignetRingBronzeProtectorR1, TimeswornExpertiseAura)

	// https://www.wowhead.com/classic/item=234436/signet-ring-of-the-bronze-dragonflight
	core.NewItemEffect(SignetRingBronzeSubjugatorR5, TimewornDecayAura)
	core.NewItemEffect(SignetRingBronzeSubjugatorR4, TimewornDecayAura)
	core.NewItemEffect(SignetRingBronzeSubjugatorR3, TimewornDecayAura)
	core.NewItemEffect(SignetRingBronzeSubjugatorR2, TimewornDecayAura)
	core.NewItemEffect(SignetRingBronzeSubjugatorR1, TimewornDecayAura)

	core.AddEffectsToTest = true
}

// https://www.wowhead.com/classic/spell=1214155/timeworn-decay
// Increases the damage dealt by all of your damage over time spells by 2% per piece of Timeworn armor equipped.
func TimewornDecayAura(agent core.Agent) {
	character := agent.GetCharacter()
	multiplier := 1 + 0.02*float64(character.PseudoStats.TimewornBonus)

	character.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode != 0 && len(spell.Dots()) > 0 {
			spell.PeriodicDamageMultiplier *= multiplier
		}
	})
}

// https://www.wowhead.com/classic/spell=1213407/timeworn-expertise
// Reduces the chance for your attacks to be dodged or parried by 1% per piece of Timeworn armor equipped.
func TimeswornExpertiseAura(agent core.Agent) {
	character := agent.GetCharacter()
	stats := stats.Stats{stats.Expertise: float64(character.PseudoStats.TimewornBonus) * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(character.GetOrRegisterAura(core.Aura{
		ActionID:   core.ActionID{SpellID: 1214218},
		Label:      "Timeworn Expertise Aura",
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(stats)
			} else {
				aura.Unit.AddStatsDynamic(sim, stats)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(stats.Multiply(-1))
			} else {
				aura.Unit.AddStatsDynamic(sim, stats.Multiply(-1))
			}
		},
	}))
}

// https://www.wowhead.com/classic/spell=1213405/timeworn-healing
// Increases the effectiveness of your healing and shielding spells by 2% per piece of Timeworn armor equipped.
func TimewornHealing(agent core.Agent) {
	character := agent.GetCharacter()
	healShieldMultiplier := 1 + 0.02*float64(character.PseudoStats.TimewornBonus)

	core.MakePermanent(character.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213405},
		Label:    "Timeworn Healing Aura",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.HealingDealtMultiplier *= healShieldMultiplier
			character.PseudoStats.ShieldDealtMultiplier *= healShieldMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.HealingDealtMultiplier /= healShieldMultiplier
			character.PseudoStats.ShieldDealtMultiplier /= healShieldMultiplier
		},
	}))
}

// https://www.wowhead.com/classic/spell=1215404/timeworn-pyromancy
// Increases the effectiveness of your Fire damage spells by 3% per piece of Timeworn armor equipped.
func TimeswornPyromancyAura(agent core.Agent) {
	character := agent.GetCharacter()
	fireMultiplier := 1 + 0.03*float64(character.PseudoStats.TimewornBonus)

	core.MakePermanent(character.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1215404},
		Label:    "Timeworn Pyromancy Aura",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= fireMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= fireMultiplier
		},
	}))
}

// https://www.wowhead.com/classic/spell=1213398/timeworn-spell
// Increases the casting speed of your spells by 2% per piece of Timeworn armor equipped.
func TimeswornSpellAura(agent core.Agent) {
	character := agent.GetCharacter()
	castSpeedMultiplier := 1 / (1 - 0.02*float64(character.PseudoStats.TimewornBonus))

	core.MakePermanent(character.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213398},
		Label:    "Timeworn Spell Aura",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyCastSpeed(castSpeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.MultiplyCastSpeed(1 / castSpeedMultiplier)
		},
	}))
}

// https://www.wowhead.com/classic/spell=1213390/timeworn-strike
// Gives you 1% chance per piece of Timeworn armor equipped to get an extra attack on regular melee or ranged hit that deals 100% weapon damage.
// (100ms cooldown)
func TimeswornStrikeAura(agent core.Agent) {
	character := agent.GetCharacter()
	procChance := float64(character.PseudoStats.TimewornBonus) * 0.01

	timeStrikeSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1213381},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagNoOnCastComplete,

		BonusCoefficient: 1,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:       "Timeworn Strike Aura",
		ActionID:   core.ActionID{SpellID: 468782},
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskWhiteHit,
		ProcChance: procChance,
		ICD:        time.Millisecond * 100,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			timeStrikeSpell.Cast(sim, result.Target)
		},
	})
}
