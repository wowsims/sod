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

	hasDeadlyBrew := rogue.HasRune(proto.RogueRune_RuneDeadlyBrew)

	baseCost := 20.0
	if ohWeapon := rogue.GetOHWeapon(); ohWeapon != nil {
		baseCost = baseCost + 10*ohWeapon.SwingSpeed
	}

	// Shiv /might/ scale with BonusWeaponDamage, if it's using https://www.wowhead.com/classic/spell=424800/shiv
	rogue.Shiv = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.RogueRune_RuneShiv)},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       rogue.builderFlags(),

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

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.Aggression] * rogue.dwsMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			// TODO check whether Shiv is affected by
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())

				switch rogue.Consumes.OffHandImbue {
				case proto.WeaponImbue_InstantPoison:
					rogue.InstantPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_DeadlyPoison:
					rogue.DeadlyPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_WoundPoison:
					rogue.WoundPoison[ShivProc].Cast(sim, target)
				default:
					if hasDeadlyBrew {
						rogue.InstantPoison[DeadlyBrewProc].Cast(sim, target)
					}
				}
			}

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
