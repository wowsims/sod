package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const RVTicks = 4

func (paladin *Paladin) registerRV() {
	if !paladin.hasRune(proto.PaladinRune_RuneCloakRighteousVengeance) {
		return
	}

	paladin.RegisterAura(core.Aura{
		Label:    "Righteous Vengeance",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {

			if spell.Flags.Matches(SpellFlag_RV) && result.DidCrit() {
				paladin.procRV(sim, result)
			}
		},
	})

	paladin.rv = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 440675},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagPureDot | core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagNoOnCastComplete,

		// SpellFlagIgnoreTargetModifiers was thought to be used based on wowhead flags
		// WCL parses show that this is not the case

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Cast: core.CastConfig{
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Righteous Vengeance",
			},
			NumberOfTicks: RVTicks,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
		},
	})
}

func (paladin *Paladin) procRV(sim *core.Simulation, result *core.SpellResult) {
	dot := paladin.rv.Dot(result.Target)

	newDamage := result.Damage * 0.3
	outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

	dot.Snapshot(result.Target, (outstandingDamage+newDamage)/float64(RVTicks), false)

	paladin.rv.Cast(sim, result.Target)
}
