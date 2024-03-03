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

	level := float64(mage.GetCharacter().Level)
	baseCalc := (13.828124 + 0.018012*level + 0.044141*level*level)
	baseDotDamage := baseCalc * .85
	baseExplosionDamage := baseCalc * 1.71
	dotCoeff := .20
	explosionCoeff := .40
	manaCost := .22

	ticks := int32(4)
	tickLength := time.Second * 3

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 400614},
		SpellCode:   SpellCode_MageLivingBomb,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage,

		DamageMultiplier: 1,
		CritMultiplier:   mage.MageCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseExplosionDamage + explosionCoeff*spell.SpellDamage()
			// baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeExpectedMagicCrit)
			}
		},
	})

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.MageRune_RuneHandsLivingBomb)},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					livingBombExplosionSpell.Cast(sim, aura.Unit)
				},
			},

			NumberOfTicks: ticks,
			TickLength:    tickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseDotDamage + dotCoeff*dot.Spell.SpellDamage()
				dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
