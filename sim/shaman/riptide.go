package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) registerRiptideSpell() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersRiptide) {
		return
	}

	hasPowerSurgeRune := shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge)

	baseHealingMultiplier := 1 + shaman.purificationHealingModifier()
	baseHealingLow := shaman.baseRuneAbilityDamage() * 1.13 * baseHealingMultiplier
	baseHealingHigh := shaman.baseRuneAbilityDamage() * 1.23 * baseHealingMultiplier
	baseHotHealing := shaman.baseRuneAbilityDamage() * 1.20
	spellCoeff := 0.215
	hotCoeff := 0.10

	shaman.Riptide = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.ShamanRune_RuneBracersRiptide)},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL | SpellFlagShaman,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.18,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Riptide",
			},
			NumberOfTicks:    5,
			TickLength:       time.Second * 3,
			BonusCoefficient: hotCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseHotHealing, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTickCounted)

				if hasPowerSurgeRune && sim.Proc(shaman.powerSurgeProcChance, "Power Surge Proc") {
					shaman.PowerSurgeHealAura.Activate(sim)
				}
			},
		},

		BonusCoefficient: spellCoeff,
		BonusCritRating:  float64(shaman.Talents.TidalMastery) * 1 * core.CritRatingPerCritChance,

		DamageMultiplier: 1,
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, spell.Unit, sim.Roll(baseHealingLow, baseHealingHigh), spell.OutcomeHealingCrit)
			spell.Hot(spell.Unit).Apply(sim)
		},
	})
}
