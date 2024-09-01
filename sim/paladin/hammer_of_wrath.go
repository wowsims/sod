package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core/proto"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerHammerOfWrath() {
	ranks := []struct {
		level     int32
		spellID   int32
		minDamage float64
		maxDamage float64
		manaCost  float64
	}{
		{level: 44, spellID: 24275, manaCost: 295, minDamage: 316, maxDamage: 348},
		{level: 52, spellID: 24274, manaCost: 360, minDamage: 412, maxDamage: 455},
		{level: 60, spellID: 24239, manaCost: 425, minDamage: 504, maxDamage: 566},
	}

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 6,
	}

	hasImprovedHammerOfWrath := paladin.hasRune(proto.PaladinRune_RuneWristImprovedHammerOfWrath)

	for i, rank := range ranks {
		rank := rank
		if paladin.Level < rank.level {
			break
		}

		paladin.GetOrRegisterSpell(core.SpellConfig{
			ActionID:        core.ActionID{SpellID: rank.spellID},
			SpellSchool:     core.SpellSchoolHoly,
			DefenseType:     core.DefenseTypeRanged,
			ProcMask:        core.ProcMaskRangedSpecial, // TODO to be tested
			Flags:           core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
			CastType:        proto.CastType_CastTypeRanged,

			Rank:          i + 1,
			RequiredLevel: int(rank.level),

			ManaCost: core.ManaCostOptions{
				FlatCost: rank.manaCost,
			},
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD:      time.Second,
					CastTime: core.TernaryDuration(hasImprovedHammerOfWrath, 0, time.Second),
				},
				IgnoreHaste: true,
				CD:          cd,
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.429,

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return sim.IsExecutePhase20()
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := sim.Roll(rank.minDamage, rank.maxDamage)
				result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)

				// should be based on target.CurrentHealthPercent(), which is not available
				if hasImprovedHammerOfWrath && result.Landed() && sim.CurrentTime >= time.Duration(0.9*float64(sim.Duration)) {
					cd.Reset()
				}
			},
		})
	}
}
