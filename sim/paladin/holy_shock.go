package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) getHolyShockBaseConfig(rank int) core.SpellConfig {
	spellId := [4]int32{0, 20473, 20929, 20930}[rank]
	baseDamageLow := [4]float64{0, 204, 365}[rank]
	baseDamageHigh := [4]float64{0, 220, 301, 395}[rank]
	manaCost := [4]float64{0, 225, 275, 325}[rank]
	level := [4]int{0, 40, 48, 56}[rank]

	spellCoeff := 0.429
	actionID := core.ActionID{SpellID: spellId}

	return core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.SpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellPower()

			// bonusCrit := core.TernaryFloat64(
			// 	guaranteed_crit,
			// 	100*core.CritRatingPerCritChance,
			// 	0)

			// spell.BonusCritRating += bonusCrit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			// spell.BonusCritRating -= bonusCrit
		},
	}

}

// Exorcism in SoD is by default castable only on demon and undead targets.
// If the paladin has the Exorcist leg rune equipped, they can cast the spell on
// any target and it additonally always crits on demon and undead targets.
func (paladin *Paladin) registerHolyShockSpell() {

	if !paladin.Talents.HolyShock {
		return
	}

	maxRank := 3
	for i := 1; i <= maxRank; i++ {
		config := paladin.getHolyShockBaseConfig(i)
		if config.RequiredLevel <= int(paladin.Level) {
			paladin.HolyShock = paladin.GetOrRegisterSpell(config)
		}
	}
}
