package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}

func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit, level int32) *core.Aura {
	duration := time.Second * 12

	return target.GetOrRegisterAura(core.Aura{
		Label:    fmt.Sprintf("Stormstrike-%s", shaman.Label),
		ActionID: StormstrikeActionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AttackTables[aura.Unit.UnitIndex][proto.CastType_CastTypeMainHand].NatureDamageTakenMultiplier *= 1.2
			shaman.AttackTables[aura.Unit.UnitIndex][proto.CastType_CastTypeOffHand].NatureDamageTakenMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AttackTables[aura.Unit.UnitIndex][proto.CastType_CastTypeMainHand].NatureDamageTakenMultiplier /= 1.2
			shaman.AttackTables[aura.Unit.UnitIndex][proto.CastType_CastTypeOffHand].NatureDamageTakenMultiplier /= 1.2
		},
	})
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) func(*core.Simulation, *core.Unit, *core.Spell) {
	var procMask core.ProcMask
	if isMH {
		procMask = core.ProcMaskMeleeMHSpecial
	} else {
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		var baseDamage float64
		spell.ProcMask = procMask
		if isMH {
			baseDamage = spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()
		} else {
			baseDamage = (spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()) *
				shaman.AutoAttacks.OHConfig().DamageMultiplier
		}

		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
	}
}

func (shaman *Shaman) registerStormstrikeSpell() {
	if !shaman.Talents.Stormstrike {
		return
	}

	manaCost := .063
	cooldown := time.Second * 6

	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    StormstrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},

		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				core.StormstrikeAura(target, shaman.Level).Activate(sim)

				mhHit(sim, target, spell)

				if shaman.HasRune(proto.ShamanRune_RuneChestDualWieldSpec) && shaman.AutoAttacks.IsDualWielding {
					ohHit(sim, target, spell)
				}

				shaman.Stormstrike.SpellMetrics[target.UnitIndex].Hits--
			}
			spell.DealOutcome(sim, result)
		},
	})
}
