package warlock

import (
	"github.com/wowsims/sod/sim/core"
)

const LifeTapRanks = 6

func (warlock *Warlock) getLifeTapBaseConfig(rank int) core.SpellConfig {
	spellId := [LifeTapRanks + 1]int32{0, 1454, 1455, 1456, 11687, 11688, 11689}[rank]
	baseDamage := [LifeTapRanks + 1]float64{0, 30, 75, 140, 220, 310, 424}[rank]
	level := [LifeTapRanks + 1]int{0, 6, 16, 26, 36, 46, 56}[rank]

	actionID := core.ActionID{SpellID: spellId}
	var manaPetMetrics *core.ResourceMetrics
	if warlock.Pet != nil {
		manaPetMetrics = warlock.Pet.NewManaMetrics(actionID)
	}
	manaMetrics := warlock.NewManaMetrics(actionID)
	spellCoef := 0.68

	return core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolShadow,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagAffliction,
		RequiredLevel: level,
		Rank:          rank,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCoefficient: spellCoef,

		DamageMultiplier: 1 + 0.1*float64(warlock.Talents.ImprovedLifeTap),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var result *core.SpellResult
			if warlock.IsTanking() {
				result = spell.CalcAndDealDamage(sim, spell.Unit, baseDamage, spell.OutcomeMagicCrit)
			} else {
				result = spell.CalcDamage(sim, spell.Unit, baseDamage, spell.OutcomeMagicCrit)
			}
			restore := result.Damage

			if warlock.MetamorphosisAura != nil && warlock.MetamorphosisAura.IsActive() {
				restore *= 2
			}
			warlock.AddMana(sim, restore, manaMetrics)
			if warlock.Pet != nil && warlock.Pet.IsActive() {
				warlock.Pet.AddMana(sim, restore, manaPetMetrics)
			}
		},
	}
}

func (warlock *Warlock) registerLifeTapSpell() {
	warlock.LifeTap = make([]*core.Spell, 0)
	for i := 1; i <= LifeTapRanks; i++ {
		config := warlock.getLifeTapBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.LifeTap = append(warlock.LifeTap, warlock.GetOrRegisterSpell(config))
		}
	}
}
