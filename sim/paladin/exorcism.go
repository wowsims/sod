package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) getExorcismBaseConfig(rank int, guaranteed_crit bool) core.SpellConfig {
	spellId := [4]int32{0, 879, 5614, 5615}[rank]
	baseDamageLow := [4]float64{0, 90, 160, 225}[rank]
	baseDamageHigh := [4]float64{0, 102, 180, 253}[rank]
	manaCost := [4]float64{0, 85, 135, 180}[rank]
	level := [4]int{0, 20, 28, 36}[rank]

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
				Duration: time.Second * 15,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.SpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellPower()

			bonusCrit := core.TernaryFloat64(
				guaranteed_crit,
				100*core.CritRatingPerCritChance,
				0)

			spell.BonusCritRating += bonusCrit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.BonusCritRating -= bonusCrit
		},
	}

}

// Exorcism in SoD is by default castable only on demon and undead targets.
// If the paladin has the Exorcist leg rune equipped, they can cast the spell on
// any target and it additonally always crits on demon and undead targets.
func (paladin *Paladin) registerExorcismSpell() {

	guaranteed_crit := false
	target := paladin.CurrentTarget
	target_is_demon_or_undead := (target.MobType == proto.MobType_MobTypeDemon) ||
		(target.MobType == proto.MobType_MobTypeUndead)

	if !paladin.HasRune(proto.PaladinRune_RuneLegsExorcist) {
		if !target_is_demon_or_undead {
			return
		}
	} else {
		if target_is_demon_or_undead {
			guaranteed_crit = true
		}
	}

	maxRank := 3
	for i := 1; i <= maxRank; i++ {
		config := paladin.getExorcismBaseConfig(i, guaranteed_crit)
		if config.RequiredLevel <= int(paladin.Level) {
			paladin.Exorcism = paladin.GetOrRegisterSpell(config)
		}
	}
}
