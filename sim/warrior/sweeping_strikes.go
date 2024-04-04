package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerSweepingStrikesCD() {
	if !warrior.Talents.SweepingStrikes {
		return
	}

	actionID := core.ActionID{SpellID: 12723}

	var curDmg float64
	ssHit := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
		},
	})

	ssAura := warrior.RegisterAura(core.Aura{
		Label:     "Sweeping Strikes",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, 5)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if aura.GetStacks() == 0 || result.Damage <= 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			// TODO BDR: Is this correct?
			// Auto attacks, Cleave, Slam etc. trigger https://www.wowhead.com/classic/spell=12723/sweeping-strikes
			// -> School dmg, can't crit, EffectBonusCoefficient=0
			// WW procs https://www.wowhead.com/classic/spell=26654/sweeping-strikes
			// -> Normalized dmg, CAN crit, EffectBonusCoefficient=0
			if spell == warrior.Execute && !sim.IsExecutePhase20() {
				curDmg = spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else if spell == warrior.Whirlwind {
				curDmg = spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				curDmg = result.Damage
			}

			// Undo armor reduction to get the raw damage value.
			curDmg /= result.ResistanceMultiplier

			ssHit.Cast(sim, warrior.Env.NextTargetUnit(result.Target))
			ssHit.SpellMetrics[result.Target.UnitIndex].Casts--
			if aura.GetStacks() > 0 {
				aura.RemoveStack(sim)
			}
		},
	})

	SweepingStrikes := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance) || warrior.StanceMatches(GladiatorStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			ssAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: SweepingStrikes,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sim.GetNumTargets() >= 2
		},
	})
}
