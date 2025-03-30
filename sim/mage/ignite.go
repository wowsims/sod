package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

// If two spells proc Ignite at almost exactly the same time, the latter
// overwrites the former.
const IgniteTicks = 2

func (mage *Mage) applyIgnite() {
	if mage.Talents.Ignite == 0 {
		return
	}

	mage.Ignite = mage.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_MageIgnite,
		ActionID:       core.ActionID{SpellID: 12654},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Cast: core.CastConfig{
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Ignite",
			},
			NumberOfTicks: IgniteTicks,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
		},
	})

	igniteMultiplier := 0.08 * float64(mage.Talents.Ignite)

	mage.procIgnite = func(sim *core.Simulation, result *core.SpellResult) {
		dot := mage.Ignite.Dot(result.Target)

		newDamage := result.Damage * igniteMultiplier

		// This was also made to not double dip on Sanctified
		dot.Snapshot(result.Target, (dot.OutstandingDmg()+newDamage)/float64(IgniteTicks), false)

		// Revert double dipping on effects confirmed to not do so
		// - Sanctified
		dot.SnapshotAttackerMultiplier /= mage.PseudoStats.SanctifiedDamageMultiplier

		mage.Ignite.Cast(sim, result.Target)
	}

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Ignite Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if spell.SpellSchool.Matches(core.SpellSchoolFire) && result.DidCrit() {
				mage.procIgnite(sim, result)
			}
		},
	}))
}
