package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic verify numbers / snapshot / travel time
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
// https://www.wowhead.com/classic/spell=401558/living-flame
func (mage *Mage) registerLivingFlameSpell() {
	if !mage.HasRune(proto.MageRune_RuneLegsLivingFlame) {
		return
	}

	level := float64(mage.GetCharacter().Level)
	baseCalc := 13.828124 + 0.018012*level + 0.044141*level*level
	baseDamage := baseCalc * 1
	spellCoeff := .143
	manaCost := .11
	cooldown := time.Minute * 1

	ticks := int32(20)
	tickLength := time.Second * 1

	mage.LivingFlame = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.MageRune_RuneLegsLivingFlame)},
		SpellCode:   SpellCode_MageLivingFlame,
		SpellSchool: core.SpellSchoolSpellfire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL | core.SpellFlagPureDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},

		// Not affected by hit
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Living Flame",
			},

			NumberOfTicks: ticks,
			TickLength:    tickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseDamage + spellCoeff*dot.Spell.SpellDamage()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
					dot.Spell.SpellMetrics[target.UnitIndex].Hits += 1
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			spell.Dot(target).Apply(sim)
			spell.SpellMetrics[target.UnitIndex].Hits -= 1
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.LivingFlame,
		Type:  core.CooldownTypeDPS,
	})
}
