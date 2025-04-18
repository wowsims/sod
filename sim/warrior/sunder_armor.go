package warrior

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerSunderArmorSpell() *WarriorSpell {
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
		warrior.Devastate = warrior.RegisterSpell(DefensiveStance, core.SpellConfig{
			ClassSpellMask: ClassSpellMask_WarriorDevastate,
			ActionID:       core.ActionID{SpellID: int32(proto.WarriorRune_RuneDevastate)},
			SpellSchool:    core.SpellSchoolPhysical,
			DefenseType:    core.DefenseTypeMelee,
			ProcMask:       core.ProcMaskMeleeMHSpecial, // TODO check whether this can actually proc stuff or not
			Flags:          core.SpellFlagMeleeMetrics | SpellFlagOffensive,

			CritDamageBonus:  warrior.impale(),
			DamageMultiplier: 1.5,
			ThreatMultiplier: 1,
			BonusCoefficient: 1,

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return warrior.PseudoStats.CanBlock
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				threatMultiplier := 1.0
				if warrior.Stance == DefensiveStance {
					threatMultiplier = 1.50
				}
				spell.ThreatMultiplier *= threatMultiplier

				weapon := warrior.AutoAttacks.MH()
				baseDamage := weapon.CalculateAverageWeaponDamage(spell.MeleeAttackPower()) / weapon.SwingSpeed
				multiplier := 1 + 0.15*float64(effectiveStacks)
				spell.CalcAndDealDamage(sim, target, baseDamage*multiplier, spell.OutcomeMeleeSpecialCritOnly)
				spell.ThreatMultiplier /= threatMultiplier
			},
		})
	}

	return warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarriorSunderArmor,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost:   15 - float64(warrior.Talents.ImprovedSunderArmor),
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
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeWeaponSpecialNoCrit) // Cannot be blocked

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
