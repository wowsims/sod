package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic verify numbers such as aoe caps and base damage
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
// https://www.wowhead.com/classic/spell=400614/living-bomb
// https://www.wowhead.com/classic/spell=400613/living-bomb
func (mage *Mage) registerLivingBombSpell() {
	if !mage.HasRune(proto.MageRune_RuneHandsLivingBomb) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.MageRune_RuneHandsLivingBomb)}
	baseDotDamage := mage.baseRuneAbilityDamage() * .85
	baseExplosionDamage := mage.baseRuneAbilityDamage() * 1.71
	dotCoeff := .20
	explosionCoeff := .40
	manaCost := .22

	ticks := int32(4)
	tickLength := time.Second * 3

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(1),
		SpellCode:   SpellCode_MageLivingBomb,
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: explosionCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseExplosionDamage, spell.OutcomeMagicCrit)
			}
		},
	})

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL | core.SpellFlagPureDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Living Bomb (DoT)",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					livingBombExplosionSpell.Cast(sim, aura.Unit)
				},
			},

			NumberOfTicks:    ticks,
			TickLength:       tickLength,
			BonusCoefficient: dotCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.SpellMetrics[target.UnitIndex].Hits -= 1
		},
	})
}
