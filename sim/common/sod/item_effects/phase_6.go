package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SignetRingDominatorR5   = 234034
	SignetRingDominatorR4   = 234030
	SignetRingDominatorR3   = 234026
	SignetRingDominatorR2   = 234021
	SignetRingDominatorR1   = 234017
	SignetRingInvokerR5     = 234032
	SignetRingInvokerR4     = 234028
	SignetRingInvokerR3     = 234024
	SignetRingInvokerR2     = 234020
	SignetRingInvokerR1     = 234016
	SignetRingFlamekeeperR5 = 234964
	SignetRingFlamekeeperR4 = 234965
	SignetRingFlamekeeperR3 = 234966
	SignetRingFlamekeeperR2 = 234967
	SignetRingFlamekeeperR1 = 234968
	SignetRingPreserverR5   = 234033
	SignetRingPreserverR4   = 234029
	SignetRingPreserverR3   = 234025
	SignetRingPreserverR2   = 234023
	SignetRingPreserverR1   = 234019
	SignetRingProtectorR5   = 234035
	SignetRingProtectorR4   = 234031
	SignetRingProtectorR3   = 234027
	SignetRingProtectorR2   = 234022
	SignetRingProtectorR1   = 234018
	SignetRingSubjugatorR5  = 234436
	SignetRingSubjugatorR4  = 234437
	SignetRingSubjugatorR3  = 234438
	SignetRingSubjugatorR2  = 234439
	SignetRingSubjugatorR1  = 234440
	SignetRingConquerorR5   = 234202
	SignetRingConquerorR4   = 234201
	SignetRingConquerorR3   = 234200
	SignetRingConquerorR2   = 234199
	SignetRingConquerorR1   = 234198
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Rings
	///////////////////////////////////////////////////////////////////////////

	core.NewItemEffect(SignetRingDominatorR5, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingDominatorR4, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingDominatorR3, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingDominatorR2, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingDominatorR1, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingConquerorR5, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingConquerorR4, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingConquerorR3, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingConquerorR2, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingConquerorR1, func(agent core.Agent) { TimeswornStrikeAura(agent.GetCharacter()) })

	core.NewItemEffect(SignetRingProtectorR5, func(agent core.Agent) { TimeswornExpertiseAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingProtectorR4, func(agent core.Agent) { TimeswornExpertiseAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingProtectorR3, func(agent core.Agent) { TimeswornExpertiseAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingProtectorR2, func(agent core.Agent) { TimeswornExpertiseAura(agent.GetCharacter()) })
	core.NewItemEffect(SignetRingProtectorR1, func(agent core.Agent) { TimeswornExpertiseAura(agent.GetCharacter()) })

}
func TimeswornExpertiseAura(character *core.Character) {
	core.MakePermanent(character.GetOrRegisterAura(core.Aura{
		ActionID:   core.ActionID{SpellID: 1214218},
		Label:      "Timeworn Expertise Aura",
		Duration:   core.NeverExpires,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			character.AddStatDynamic(sim, stats.Expertise, aura.Unit.GetStat(stats.Timeworn)*core.CritRatingPerCritChance)
		},
	}))
}

func TimeswornStrikeAura(character *core.Character) {

	chance := character.Unit.GetStat(stats.Timeworn) * 0.01

	timeStrikeSpell := character.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1213381},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:     "Timeworn Strike Aura",
		ActionID: core.ActionID{SpellID: 468782},
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		ProcMask: core.ProcMaskWhiteHit,
		ICD:      time.Millisecond * 100,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			roll := sim.RandomFloat("Timeworn Strike Aura")

			if roll < chance {
				character.Unit.GetStat(stats.Timeworn)
				timeStrikeSpell.Cast(sim, result.Target)
			}
		},
	})
}
