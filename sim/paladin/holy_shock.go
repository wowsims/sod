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
	level := holyShockLevels[rank]
	spellId := holyShockSpellIds[rank]
	baseDamageLow := holyShockBaseDamages[rank][0]
	baseDamageHigh := holyShockBaseDamages[rank][1]
	manaCost := holyShockManaCosts[rank]

	// Art of War reduces cost by 80%
	hasAoW := paladin.HasRune(proto.PaladinRune_RuneFeetTheArtOfWar)
	manaCostMultiplier := core.TernaryFloat64(hasAoW, 0.2, 1.0)
	// Infusion of Light increases base damage by 20%
	hasInfusion := paladin.HasRune(proto.PaladinRune_RuneWaistInfusionOfLight)
	damageMultiplier := core.TernaryFloat64(hasInfusion, 1.2, 1.0)
	manaCostActual := core.TernaryFloat64(
		hasAoW,
		manaCost*manaCostMultiplier,
		manaCost,
	)

	spellCoeff := 0.429
	actionID := core.ActionID{SpellID: spellId}
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
			FlatCost:   manaCost,
			Multiplier: manaCostMultiplier,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: *paladin.HolyShockCooldown,
		},

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.SpellCritMultiplier(),
		BonusCritRating:  paladin.getBonusCritChanceFromHolyPower(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			// If we crit, Infusion of Light refunds base mana cost and resets Holy Shock/Exorcism.
			if !result.Outcome.Matches(core.OutcomeCrit) || !hasInfusion {
				return
			}
			// The mana refund in game is bugged and will often refund more mana than intended.
			// For now, refund the base cost if no AoW rune equipped, or the actual cost if it is.
			paladin.AddMana(sim, manaCostActual, manaMetrics)
			paladin.HolyShockCooldown.Reset()
			paladin.ExorcismCooldown.Reset()
		},
	}

}

func (paladin *Paladin) registerHolyShockSpell() {
	// If the player has Holy Shock talented, register all holyShockRanks up to their level.
	paladin.HolyShockCooldown = &core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 15,
	}
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
