package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerHemorrhageSpell() {
	// Minimum level of 30 to talent Hemorrhage
	if rogue.Level < 30 {
		return
	}

	if !rogue.Talents.Hemorrhage {
		return
	}

	spellID := map[int32]int32{
		40: 16511,
		50: 17347,
		60: 17348,
	}[rogue.Level]

	actionID := core.ActionID{SpellID: spellID}

	var numPlayers int
	for _, u := range rogue.Env.Raid.AllUnits {
		if u.Type == core.PlayerUnit {
			numPlayers++
		}
	}

	var hemoAuras core.AuraArray

	if numPlayers >= 2 {
		hemoAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
			return core.HemorrhageAura(target, rogue.Level)
		})
	}

	rogue.Hemorrhage = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35.0,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				if len(hemoAuras) > 0 {
					hemoAura := hemoAuras.Get(target)
					hemoAura.Activate(sim)
					hemoAura.SetStacks(sim, 30)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
