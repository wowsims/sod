package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) registerMainGaucheSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneMainGauche) {
		return
	}

	actionID := core.ActionID{SpellID: 424919}

	// Aura gained regardless of landed hit.  Need to confirm later with tank sim if parry is being modified correctly
	mainGauchAura := rogue.RegisterAura(core.Aura{
		Label:    "Main Gauche Buff",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AddStatDynamic(sim, stats.Parry, 10*core.ParryRatingPerParryChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AddStatDynamic(sim, stats.Parry, -10*core.ParryRatingPerParryChance)
		},
	})

	rogue.MainGauche = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | core.SpellFlagAPL | SpellFlagColdBlooded,
		EnergyCost: core.EnergyCostOptions{
			Cost:   []float64{15, 12, 10}[rogue.Talents.ImprovedSinisterStrike],
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 20,
			},
			IgnoreHaste: true,
		},

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.Aggression],
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			mainGauchAura.Activate(sim)

			if result.Landed() {
				rogue.AddComboPoints(sim, 3, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
