package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const ranksHolyShock = 3

var holyShockLevels = [ranksHolyShock + 1]int{0, 40, 48, 56}
var holyShockSpellIds = [ranksHolyShock + 1]int32{0, 20473, 20929, 20930}
var holyShockBaseDamages = [ranksHolyShock + 1][]float64{{0}, {204, 220}, {279, 301}, {365, 395}}
var holyShockManaCosts = [ranksHolyShock + 1]float64{0, 225, 275, 325}

func (paladin *Paladin) getHolyShockBaseConfig(rank int) core.SpellConfig {
	spellId := holyShockSpellIds[rank]
	baseDamageLow := holyShockBaseDamages[rank][0]
	baseDamageHigh := holyShockBaseDamages[rank][1]
	manaCost := holyShockManaCosts[rank]
	level := holyShockLevels[rank]

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
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}

}

func (paladin *Paladin) registerHolyShockSpell() {
	// If the player has Holy Shock talented, register all ranksHolyShock up to their level.
	if !paladin.Talents.HolyShock {
		return
	}

	paladin.HolyShock = make([]*core.Spell, ranksHolyShock+1)
	for rank := 1; rank <= ranksHolyShock; rank++ {
		if int(paladin.Level) >= holyShockLevels[rank] {
			paladin.HolyShock[rank] = paladin.RegisterSpell(paladin.getHolyShockBaseConfig(rank))
		}
	}
}
