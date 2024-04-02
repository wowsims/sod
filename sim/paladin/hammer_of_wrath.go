package paladin

import (
	"github.com/wowsims/sod/sim/core/proto"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerHammerOfWrathSpell() {
	ranks := []struct {
		level      int32
		spellID    int32
		damageLow  float64
		damageHigh float64
		manaCost   float64
	}{
		{level: 44, spellID: 24275, damageLow: 316, damageHigh: 348, manaCost: 295},
		{level: 52, spellID: 24274, damageLow: 412, damageHigh: 455, manaCost: 360},
		{level: 60, spellID: 24239, damageLow: 504, damageHigh: 566, manaCost: 425},
	}

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 6,
	}

	hasImprovedHammerOfWrath := paladin.HasRune(proto.PaladinRune_RuneWristImprovedHammerOfWrath)

	for i, rank := range ranks {
		if paladin.Level < rank.level {
			break
		}

		paladin.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeRanged,
			ProcMask:    core.ProcMaskRangedSpecial, // TODO to be tested
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

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

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return sim.IsExecutePhase20()
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := sim.Roll(rank.damageLow, rank.damageHigh) + 0.429*spell.SpellDamage()
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)

				// should be based on target.CurrentHealthPercent(), which is not available
				if hasImprovedHammerOfWrath && sim.CurrentTime >= time.Duration(0.9*float64(sim.Duration)) {
					cd.Reset()
				}
			},
		})
	}
}
