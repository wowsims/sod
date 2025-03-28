package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerMassRegenerationSpell() {
	if !mage.HasRune(proto.MageRune_RuneLegsMassRegeneration) {
		return
	}

	mage.MassRegeneration = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(proto.MageRune_RuneLegsMassRegeneration)},
		ClassSpellMask: ClassSpellMask_MageMassRegeneration,
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.45,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Mass Regeneration",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.152,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Hot(&mage.Unit).Apply(sim)
		},
	})
}
