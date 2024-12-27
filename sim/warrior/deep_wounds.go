package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) applyDeepWounds() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	spellID := map[int32]int32{
		1: 12834,
		2: 12849,
		3: 12867,
	}[warrior.Talents.DeepWounds]

	warrior.DeepWounds = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		SpellCode:   SpellCode_WarriorDeepWounds,
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Deep Wounds",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrRefresh(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Deep Wounds Talent",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskEmpty) || !spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
				return
			}

			// Ravager doesn't proc Deep Wounds
			if spell.ActionID.SpellID == 9633 {
				return
			}

			if result.Outcome.Matches(core.OutcomeCrit) {
				warrior.procDeepWounds(sim, result.Target, spell.IsOH())
			}
		},
	}))
}

func (warrior *Warrior) procDeepWounds(sim *core.Simulation, target *core.Unit, isOh bool) {
	dot := warrior.DeepWounds.Dot(target)

	outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

	var awd float64
	if isOh {
		attackTableOh := warrior.AttackTables[target.UnitIndex][proto.CastType_CastTypeOffHand]
		adm := warrior.AutoAttacks.OHAuto().AttackerDamageMultiplier(attackTableOh, true)
		awd = warrior.AutoAttacks.OH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower()) * 0.5 * adm
	} else { // MH
		attackTableMh := warrior.AttackTables[target.UnitIndex][proto.CastType_CastTypeMainHand]
		adm := warrior.AutoAttacks.MHAuto().AttackerDamageMultiplier(attackTableMh, true)
		awd = warrior.AutoAttacks.MH().CalculateAverageWeaponDamage(dot.Spell.MeleeAttackPower()) * adm
	}

	newDamage := awd * 0.2 * float64(warrior.Talents.DeepWounds)

	dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(dot.NumberOfTicks)
	dot.SnapshotAttackerMultiplier = warrior.DeepWounds.DamageMultiplier * warrior.DeepWounds.DamageMultiplierAdditive

	warrior.DeepWounds.Cast(sim, target)
}
