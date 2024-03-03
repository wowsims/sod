package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerShivSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneShiv) {
		return
	}

	baseCost := 20.0
	if ohWeapon := rogue.GetOHWeapon(); ohWeapon != nil {
		baseCost = baseCost + 10*ohWeapon.SwingSpeed
	}

	rogue.Shiv = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 424799},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   baseCost - []float64{0, 3, 5}[rogue.Talents.ImprovedSinisterStrike],
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 * rogue.dwsMultiplier(),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			// TODO: Cannot be Parry/Dodge based on testing, can miss
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				// 100% application (except for 1%? It can resist extremely rarely)
				if rogue.getImbueProcMask(proto.WeaponImbue_InstantPoison).Matches(core.ProcMaskMeleeOH) {
					rogue.InstantPoison[ShivProc].Cast(sim, target)
					return
				} else if rogue.getImbueProcMask(proto.WeaponImbue_DeadlyPoison).Matches(core.ProcMaskMeleeOH) {
					rogue.DeadlyPoison[ShivProc].Cast(sim, target)
					return
				} else if rogue.getImbueProcMask(proto.WeaponImbue_WoundPoison).Matches(core.ProcMaskMeleeOH) {
					rogue.WoundPoison[ShivProc].Cast(sim, target)
				}
			}
		},
	})
}
