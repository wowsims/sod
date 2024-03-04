package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerExposeArmorSpell() {
	rogue.ExposeArmorAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return core.ExposeArmorAura(target, rogue.Talents.ImprovedExposeArmor, rogue.Level)
	})

	spellID := map[int32]int32{
		25: 8647,
		40: 8650,
		50: 11197,
		60: 11198,
	}[rogue.Level]

	arpenPerCombo := map[int32]float64{
		25: 80,
		40: 210,
		50: 270,
		60: 340,
	}

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: spellID},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | rogue.finisherFlags() | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:   25.0,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0 && rogue.CanApplyExposeAura(target)
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				debuffAura := rogue.ExposeArmorAuras.Get(target)

				if debuffAura.IsActive() {
					// More powerful?
					debuffAura.Deactivate(sim)
				}
				// recalculate and apply armor reduction value
				arpen := arpenPerCombo[rogue.Level]
				arpen *= float64(rogue.ComboPoints())
				// Improved Expose Armor Multiplier
				arpen *= 1 + 0.25*float64(rogue.Talents.ImprovedExposeArmor)
				debuffAura.ExclusiveEffects[0].Priority = arpen
				debuffAura.Activate(sim)

				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{rogue.ExposeArmorAuras},
	})
}

func (rogue *Rogue) CanApplyExposeAura(target *core.Unit) bool {
	return rogue.ExposeArmorAuras.Get(target).IsActive() || !rogue.ExposeArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
