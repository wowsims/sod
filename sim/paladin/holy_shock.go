package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const holyShockRanks = 3

var holyShockLevels = [holyShockRanks + 1]int{0, 40, 48, 56}
var holyShockSpellIds = [holyShockRanks + 1]int32{0, 20473, 20929, 20930}
var holyShockBaseDamages = [holyShockRanks + 1][]float64{{0}, {204, 220}, {279, 301}, {365, 395}}
var holyShockManaCosts = [holyShockRanks + 1]float64{0, 225, 275, 325}

func (paladin *Paladin) getHolyShockBaseConfig(rank int) core.SpellConfig {
	spellId := holyShockSpellIds[rank]
	baseDamageLow := holyShockBaseDamages[rank][0]
	baseDamageHigh := holyShockBaseDamages[rank][1]
	manaCost := holyShockManaCosts[rank]
	level := holyShockLevels[rank]

	spellCoeff := 0.429
	actionID := core.ActionID{SpellID: spellId}
	procChance := 0.2 * float64(paladin.Talents.Illumination) // Chance for Illumination to refund Mana on crit.
	manaMetrics := paladin.NewManaMetrics(paladin.getIlluminationActionID())

	return core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,
		SpellCode:     SpellCode_PaladinHolyShock,
		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
			Multiplier: core.TernaryFloat64(
				paladin.HasRune(proto.PaladinRune_RuneFeetTheArtOfWar),
				0.2,
				1.0,
			),
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
		BonusCritRating:  paladin.getBonusCritChanceFromHolyPower(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if procChance == 1 || sim.RandomFloat("Illumination") < procChance {
				paladin.AddMana(sim, manaCost, manaMetrics)
			}
		},
	}

}

func (paladin *Paladin) registerHolyShockSpell() {
	// If the player has Holy Shock talented, register all holyShockRanks up to their level.
	if !paladin.Talents.HolyShock {
		return
	}

	paladin.HolyShock = make([]*core.Spell, holyShockRanks+1)
	for rank := 1; rank <= holyShockRanks; rank++ {
		if int(paladin.Level) >= holyShockLevels[rank] {
			paladin.HolyShock[rank] = paladin.RegisterSpell(paladin.getHolyShockBaseConfig(rank))
		}
	}
}
