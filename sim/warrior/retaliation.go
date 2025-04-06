package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerRetaliationCD(sharedTimer *core.Timer) {
	actionID := core.ActionID{SpellID: 20230}

	// The hits will proc in any stance
	attackSpell := warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarriorRetaliation,
		ActionID:       core.ActionID{SpellID: 20240},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMH,
		Flags:          core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warrior.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	buffAura := warrior.RegisterAura(core.Aura{
		ActionID:  actionID,
		Label:     "Retaliation",
		Duration:  time.Second * 15,
		MaxStacks: 30,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) && result.Landed() && result.Damage > 0 {
				attackSpell.Cast(sim, spell.Unit)
			}
		},
	})

	warrior.Retaliation = warrior.RegisterSpell(BattleStance, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_WarriorRecklesness,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 5,
			},
			SharedCD: core.Cooldown{
				Timer:    sharedTimer,
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			buffAura.Activate(sim)
			if warrior.Recklessness != nil {
				warrior.Recklessness.RelatedSelfBuff.Deactivate(sim)
			}
			if warrior.ShieldWall != nil {
				warrior.ShieldWall.RelatedSelfBuff.Deactivate(sim)
			}
		},

		RelatedSelfBuff: buffAura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.Retaliation.Spell,
		Type:  core.CooldownTypeDPS,
		// Require manual CD usage
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return false
		},
	})
}
