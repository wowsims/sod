package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const FanOfKnivesSpellID int32 = 409240

func (rogue *Rogue) makeFanOfKnivesWeaponHitSpell(isMH bool) *core.Spell {
	var procMask core.ProcMask
	var weaponMultiplier float64
	var actionID core.ActionID
	activate2PcBonuses := rogue.HasSetBonus(ItemSetNightSlayerBattlearmor, 2)  && rogue.HasAura("Blade Dance") && rogue.HasRune(proto.RogueRune_RuneJustAFleshWound)
	if isMH {
		actionID = core.ActionID{SpellID: FanOfKnivesSpellID}.WithTag(1)
		weaponMultiplier = core.TernaryFloat64(rogue.HasDagger(core.MainHand), 0.75, 0.5)
		procMask = core.ProcMaskMeleeMHSpecial
	} else {
		actionID = core.ActionID{SpellID: FanOfKnivesSpellID}.WithTag(2)
		weaponMultiplier = core.TernaryFloat64(rogue.HasDagger(core.OffHand), 0.75, 0.5)
		weaponMultiplier *= rogue.dwsMultiplier()
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagColdBlooded,

		DamageMultiplier: weaponMultiplier,
		ThreatMultiplier: core.TernaryFloat64(activate2PcBonuses, 2, 1),
	})
}

// TODO: 8 yd range
func (rogue *Rogue) registerFanOfKnives() {
	if !rogue.HasRune(proto.RogueRune_RuneFanOfKnives) {
		return
	}

	activate2PcBonuses := rogue.HasSetBonus(ItemSetNightSlayerBattlearmor, 2)  && rogue.HasAura("Blade Dance") && rogue.HasRune(proto.RogueRune_RuneJustAFleshWound)
	mhSpell := rogue.makeFanOfKnivesWeaponHitSpell(true)
	ohSpell := rogue.makeFanOfKnivesWeaponHitSpell(false)
	results := make([]*core.SpellResult, len(rogue.Env.Encounter.TargetUnits))

	rogue.FanOfKnives = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: FanOfKnivesSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagCarnage,

		EnergyCost: core.EnergyCostOptions{
			Cost: 50 - core.TernaryFloat64(activate2PcBonuses, 20, 0),
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
				results[i] = ohSpell.CalcDamage(sim, aoeTarget, baseDamage, ohSpell.OutcomeMeleeSpecialHit)
			}
			for i := range sim.Encounter.TargetUnits {
				ohSpell.DealDamage(sim, results[i])
			}

			for i, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := mhSpell.Unit.MHWeaponDamage(sim, mhSpell.MeleeAttackPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				results[i] = mhSpell.CalcDamage(sim, aoeTarget, baseDamage, mhSpell.OutcomeMeleeSpecialHit)
			}
			for i := range sim.Encounter.TargetUnits {
				mhSpell.DealDamage(sim, results[i])
			}
		},
	})
}
