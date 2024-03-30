package warrior

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) newSunderArmorSpell() *core.Spell {
	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(core.SunderArmorAura)
	spellID := map[int32]int32{
		25: 7405,
		40: 8380,
		50: 11596,
		60: 11597,
	}[warrior.Level]

	isDevastate := warrior.HasRune(proto.WarriorRune_RuneDevastate)

	config := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,
		DefenseType: core.DefenseTypeMelee,

		RageCost: core.RageCostOptions{
			Cost:   15 - warrior.FocusedRageDiscount - float64(warrior.Talents.ImprovedSunderArmor),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.CanApplySunderAura(target)
		},

		CritDamageBonus:  warrior.impale(),
		DamageMultiplier: 1,

		ThreatMultiplier: 1,
		// TODO Warrior: set threat according to spell's level
		FlatThreatBonus: 360,

		RelatedAuras: []core.AuraArray{warrior.SunderArmorAuras},
	}

	config.Flags |= core.SpellFlagAPL

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		var result *core.SpellResult

		overrided := false

		if target.GetAura("Degrade (Homunculus)").IsActive() || target.GetAura("ExposeArmor").IsActive() {
			overrided = true
		}

		aura := warrior.SunderArmorAuras.Get(target)
		if isDevastate {
			stacks := core.TernaryFloat64(overrided, 5.0, float64(aura.GetStacks()))
			modifier := 1.5 + 0.1*float64(stacks)
			damage := modifier * warrior.AutoAttacks.MH().AverageDamage() / warrior.SwingSpeed()
			result = spell.CalcDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)
		} else {
			result = spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
		}

		if !result.Landed() {
			spell.IssueRefund(sim)
			return
		}

		aura.Activate(sim)
		if aura.IsActive() && !overrided {
			aura.AddStack(sim)
		}

		spell.DealOutcome(sim, result)
	}
	return warrior.RegisterSpell(config)
}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
