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

	hasImprovedScorchTalent := mage.Talents.ImprovedScorch > 0

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
		SpellCode:   SpellCode_MageLivingBombExplosion,
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: explosionCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseExplosionDamage, spell.OutcomeMagicCrit)

			// Unlike the normal Scorch application, Living Bomb explosions are guaranteed to apply Improved Scorch as long as they have at least 1 point talented.
			if hasImprovedScorchTalent {
				impScorchAura := mage.ImprovedScorchAuras.Get(target)
				impScorchAura.Activate(sim)
				impScorchAura.AddStack(sim)
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

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Living Bomb (DoT)",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, aoeTarget := range sim.Encounter.TargetUnits {
						livingBombExplosionSpell.Cast(sim, aoeTarget)
					}
				},
			},

			NumberOfTicks:    ticks,
			TickLength:       tickLength,
			BonusCoefficient: dotCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

				// Unlike the normal Scorch application, Living Bomb ticks are guaranteed to apply Improved Scorch as long as they have at least 1 point talented.
				if hasImprovedScorchTalent {
					impScorchAura := mage.ImprovedScorchAuras.Get(target)
					impScorchAura.Activate(sim)
					impScorchAura.AddStack(sim)
				}
			},
		},

		BonusCritRating: 2 * float64(mage.Talents.Incinerate),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).ApplyOrReset(sim)
			}
		},
	})
}
