package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) registerVoidZoneSpell() {
	if !priest.HasRune(proto.PriestRune_RuneBracersVoidZone) {
		return
	}

	ticks := int32(10)
	tickLength := time.Second * 1

	baseTickDamage := priest.baseRuneAbilityDamage() * .51
	spellCoeff := 0.084
	manaCost := .21
	cooldown := time.Second * 30

	hasDespairRune := priest.HasRune(proto.PriestRune_RuneBracersDespair)

	priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.PriestRune_RuneBracersVoidZone)},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagPureDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: cooldown,
			},
		},

		BonusCritRating: priest.forceOfWillCritRating(),
		BonusHitRating:  priest.shadowHitModifier(),

		CritDamageBonus: core.TernaryFloat64(hasDespairRune, 1, 0),

		DamageMultiplier: priest.forceOfWillDamageModifier() * priest.darknessDamageModifier(),
		ThreatMultiplier: priest.shadowThreatModifier(),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "VoidZone-" + strconv.Itoa(1),
			},

			NumberOfTicks:    ticks,
			TickLength:       tickLength,
			BonusCoefficient: spellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseTickDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					if hasDespairRune {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTickSnapshotCritCounted)
					} else {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTickCounted)
					}
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseTickDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})
}
