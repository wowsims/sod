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
		50: 275,
		60: 340,
	}[rogue.Level]

	arpenPerCombo *= []float64{1, 1.25, 1.5}[rogue.Talents.ImprovedExposeArmor]

	// share ExtraCastCondition() state with ApplyEffects()
	var arpen float64
	var eaAura *core.Aura

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: spellID},
		SpellSchool:  core.SpellSchoolPhysical,
		DefenseType:  core.DefenseTypeMelee,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        rogue.finisherFlags(),
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost: 25,
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
			if rogue.ComboPoints() == 0 {
				return false
			}

			eaAura = rogue.ExposeArmorAuras.Get(target)
			arpen = float64(rogue.ComboPoints()) * arpenPerCombo

			if curActive := eaAura.ExclusiveEffects[0].Category.GetActiveEffect(); curActive != nil {
				return arpen >= curActive.Priority
			}
			return true
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				eaAura.ExclusiveEffects[0].Priority = arpen
				eaAura.Activate(sim)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{rogue.ExposeArmorAuras},
	})
}
