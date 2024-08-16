package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerSweepingStrikesCD() {
	if !warrior.Talents.SweepingStrikes {
		return
	}

	// Procs from auto attacks and most abilities https://www.wowhead.com/classic/spell=12723/sweeping-strikes
	var curDmg float64
	hitSchoolDamagWithValue := warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 12723},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
		},
	})

	// Procs from WW, also Execute? https://www.wowhead.com/classic/spell=26654/sweeping-strikes
	hitSpecialNormalized := warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26654},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	actionID := core.ActionID{SpellID: 12292}

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

			var spellToUse *WarriorSpell

			if spell.SpellCode == SpellCode_WarriorWhirlwindMH || (spell.SpellCode == SpellCode_WarriorExecute && !sim.IsExecutePhase20()) {
				spellToUse = hitSpecialNormalized
			} else {
				curDmg = result.Damage
				curDmg /= result.ResistanceMultiplier // Undo armor reduction to get the raw damage value.
				spellToUse = hitSchoolDamagWithValue
			}

			spellToUse.Cast(sim, warrior.Env.NextTargetUnit(result.Target))
			spellToUse.SpellMetrics[result.Target.UnitIndex].Casts--
			if aura.GetStacks() > 0 {
				aura.RemoveStack(sim)
			}
		},
	})

	SweepingStrikes := warrior.RegisterSpell(BattleStance, core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagHelpful,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			ssAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: SweepingStrikes.Spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sim.GetNumTargets() >= 2
		},
	})
}
