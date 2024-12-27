package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const FanOfKnivesSpellID int32 = 409240

func (rogue *Rogue) makeFanOfKnivesWeaponHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: FanOfKnivesSpellID}.WithTag(1)
	procMask := core.ProcMaskMeleeMHSpecial
	flags := core.SpellFlagMeleeMetrics | SpellFlagColdBlooded
	weaponMultiplier := core.TernaryFloat64(rogue.HasDagger(core.MainHand), 0.75, 0.5)

	if !isMH {
		actionID.Tag = 2
		procMask = core.ProcMaskMeleeOHSpecial
		flags |= core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell
		weaponMultiplier = core.TernaryFloat64(rogue.HasDagger(core.OffHand), 0.75, 0.5) * rogue.dwsMultiplier()
	}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    procMask,
		Flags:       flags,

		DamageMultiplier: weaponMultiplier,
		ThreatMultiplier: 1,
	})
}

// TODO: 8 yd range
func (rogue *Rogue) registerFanOfKnives() {
	if !rogue.HasRune(proto.RogueRune_RuneFanOfKnives) {
		return
	}

	mhSpell := rogue.makeFanOfKnivesWeaponHitSpell(true)
	ohSpell := rogue.makeFanOfKnivesWeaponHitSpell(false)
	results := make([]*core.SpellResult, len(rogue.Env.Encounter.TargetUnits))

	rogue.FanOfKnives = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: FanOfKnivesSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagCarnage,

		EnergyCost: core.EnergyCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			// Calc and apply all OH hits first, because MH hits can benefit from an OH felstriker proc.
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := ohSpell.Unit.OHWeaponDamage(sim, ohSpell.MeleeAttackPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				results[i] = ohSpell.CalcDamage(sim, aoeTarget, baseDamage, ohSpell.OutcomeMeleeSpecialHitAndCrit)
			}
			for i := range sim.Encounter.TargetUnits {
				ohSpell.DealDamage(sim, results[i])
			}

			for i, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := mhSpell.Unit.MHWeaponDamage(sim, mhSpell.MeleeAttackPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				results[i] = mhSpell.CalcDamage(sim, aoeTarget, baseDamage, mhSpell.OutcomeMeleeSpecialHitAndCrit)
			}
			for i := range sim.Encounter.TargetUnits {
				mhSpell.DealDamage(sim, results[i])
			}
		},
	})
}
