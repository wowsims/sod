package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const SunfireTicks = int32(4)

// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
func (druid *Druid) registerSunfireSpell() {
	if !druid.HasRune(proto.DruidRune_RuneHandsSunfire) {
		return
	}

	moonfuryMultiplier := druid.MoonfuryDamageMultiplier()
	impMoonfireMultiplier := druid.ImprovedMoonfireDamageMultiplier()

	baseLowDamage := druid.baseRuneAbilityDamage() * 1.3 * moonfuryMultiplier * impMoonfireMultiplier
	baseHighDamage := druid.baseRuneAbilityDamage() * 1.52 * moonfuryMultiplier * impMoonfireMultiplier
	baseDotDamage := druid.baseRuneAbilityDamage() * 0.65 * moonfuryMultiplier * impMoonfireMultiplier
	spellCoeff := .15
	spellDotCoeff := .13

	druid.SunfireDotMultiplier = 1

	druid.Sunfire = druid.RegisterSpell(Humanoid|Bear|Cat|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 414684},
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | SpellFlagOmen | core.SpellFlagResetAttackSwing,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.21,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Sunfire",
				ActionID: core.ActionID{SpellID: 414684},
			},
			NumberOfTicks:    SunfireTicks,
			TickLength:       time.Second * 3,
			BonusCoefficient: spellDotCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, baseDotDamage, false)
				dot.SnapshotAttackerMultiplier *= druid.SunfireDotMultiplier
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		BonusCritRating: druid.ImprovedMoonfireCritBonus() * core.SpellCritRatingPerCritChance,

		CritDamageBonus: druid.vengeance(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				dot := spell.Dot(target)
				dot.NumberOfTicks = SunfireTicks
				dot.RecomputeAuraDuration()
				dot.Apply(sim)
			}
		},
	})
}
