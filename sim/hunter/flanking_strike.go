package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerFlankingStrikeSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsFlankingStrike) {
		return
	}

	hasCatlikeReflexes := hunter.HasRune(proto.HunterRune_RuneHelmCatlikeReflexes)

	cooldownModifier := 1.0
	if hasCatlikeReflexes {
		cooldownModifier *= 0.5
	}

	hunter.FlankingStrikeAura = hunter.GetOrRegisterAura(core.Aura{
		Label:     "Flanking Strike Buff",
		ActionID:  core.ActionID{SpellID: 415320},
		MaxStacks: 3,
		Duration:  time.Second * 10,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			hunter.PseudoStats.DamageDealtMultiplier /= 1 + (0.05 * float64(oldStacks))
			hunter.PseudoStats.DamageDealtMultiplier *= 1 + (0.05 * float64(newStacks))
		},
	})

	if hunter.pet != nil {
		hunter.pet.flankingStrike = hunter.pet.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 415320},
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagMeleeMetrics,

			DamageMultiplier: 1,
			BonusCoefficient: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			},
		})

		hunter.pet.RegisterAura(core.Aura{
			Label:    "Flanking Strike Refresh",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial | core.ProcMaskSpellDamage) {
					if sim.RandomFloat("Flanking Strike Refresh") < 0.33 {
						hunter.FlankingStrike.CD.Set(sim.CurrentTime)
					}
				}
			},
		})
	}

	hunter.FlankingStrike = hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 415320},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagStrike,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.015,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * time.Duration(30*cooldownModifier),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= core.MaxMeleeAttackDistance
		},

		CritDamageBonus:  hunter.mortalShots(),
		DamageMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hunter.pet != nil {
				hunter.pet.flankingStrike.Cast(sim, hunter.pet.CurrentTarget)
			}

			hunter.FlankingStrikeAura.Activate(sim)
			hunter.FlankingStrikeAura.AddStack(sim)
		},
	})
}
