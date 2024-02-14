package druid

import (
	"time"

	item_sets "github.com/wowsims/sod/sim/common/sod/items_sets"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (druid *Druid) applyStarsurge() {
	if !druid.HasRune(proto.DruidRune_RuneLegsStarsurge) {
		return
	}

	level := float64(druid.GetCharacter().Level)
	baseCalc := (9.183105 + 0.616405*level + 0.028608*level*level)
	baseLowDamage := baseCalc * 3.81
	baseHighDamage := baseCalc * 4.67

	druid.Starsurge = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 417157},
		SpellCode:   SpellCode_DruidStarsurge,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing,

		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.01 * (1 - 0.03*float64(druid.Talents.Moonglow)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.VengeanceCritMultiplier(),
		BonusCritRating:  core.TernaryFloat64(druid.HasSetBonus(item_sets.ItemSetInsulatedSorcerorLeather, 3), 2, 0) * core.CritRatingPerCritChance,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)*druid.MoonfuryDamageMultiplier() + spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.DidCrit() && druid.NaturesGraceProcAura != nil {
				druid.NaturesGraceProcAura.Activate(sim)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
