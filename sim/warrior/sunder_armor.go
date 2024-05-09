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

	spell_level := map[int32]int32{
		25: 22,
		40: 34,
		50: 46,
		60: 58,
	}[warrior.Level]

	var effectiveStacks int32
	var canApplySunder bool

	if warrior.HasRune(proto.WarriorRune_RuneDevastate) {
		warrior.Devastate = warrior.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: int32(proto.WarriorRune_RuneDevastate)},
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeMHSpecial, // TODO check whether this can actually proc stuff or not
			Flags:       core.SpellFlagMeleeMetrics,

			CritDamageBonus:  warrior.impale(),
			DamageMultiplier: 1.5,

			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				weapon := warrior.AutoAttacks.MH()
				baseDamage := weapon.CalculateAverageWeaponDamage(spell.MeleeAttackPower()) / weapon.SwingSpeed
				multiplier := 1 + 0.1*float64(effectiveStacks)
				spell.CalcAndDealDamage(sim, target, baseDamage*multiplier, spell.OutcomeMeleeSpecialCritOnly)
			},
		})
	}

	return warrior.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

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
			sa := warrior.SunderArmorAuras.Get(target)
			if sa.IsActive() {
				effectiveStacks = sa.GetStacks()
				canApplySunder = true
			} else if sa.ExclusiveEffects[0].Category.AnyActive() {
				effectiveStacks = sa.MaxStacks
				canApplySunder = false
			} else {
				effectiveStacks = 0
				canApplySunder = true
			}
			return canApplySunder || warrior.Devastate != nil
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  2.25 * 2 * float64(spell_level),

		RelatedAuras: []core.AuraArray{warrior.SunderArmorAuras},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeWeaponSpecialNoCrit) // completely stopped by blocks

			if !result.Landed() {
				spell.IssueRefund(sim)
				return
			}

			if warrior.Devastate != nil {
				warrior.Devastate.Cast(sim, target)
			}

			if canApplySunder {
				sa := warrior.SunderArmorAuras.Get(target)
				sa.Activate(sim)
				sa.AddStack(sim)
			}
		},
	})
}
