package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (rogue *Rogue) registerPoisonedKnife() {
	if !rogue.HasRune(proto.RogueRune_RunePoisonedKnife) {
		return
	}

	poisonedKnifeMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 425013})
	hasDeadlyBrew := rogue.HasRune(proto.RogueRune_RuneDeadlyBrew)
	hasJustAFleshWound := rogue.HasRune(proto.RogueRune_RuneJustAFleshWound)

	// Poisoned Knife /might/ scale with BonusWeaponDamage, if it's using https://www.wowhead.com/classic/spell=425013/poisoned-knife
	rogue.PoisonedKnife = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RoguePoisonedKnife,
		ActionID:       core.ActionID{SpellID: int32(proto.RogueRune_RunePoisonedKnife)},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          rogue.builderFlags(),
		MaxRange:       25,

		EnergyCost: core.EnergyCostOptions{
			Cost:   25,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 6,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.HasOHWeapon()
		},
		CastType: proto.CastType_CastTypeOffHand,

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1 * rogue.dwsMultiplier(),
		ThreatMultiplier: core.TernaryFloat64(hasJustAFleshWound, 1.5, 1),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			numStacks := 0.0

			// Cannot Miss, Dodge, or Parry as per spell flags
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())

				if rogue.usingOccult {
					numStacks = float64(rogue.occultPoisonTick.Dot(target).GetStacks())
				} else if rogue.usingDeadly {
					numStacks = float64(rogue.deadlyPoisonTick.Dot(target).GetStacks())
				}

				rogue.AddEnergy(sim, numStacks*5, poisonedKnifeMetrics)

				// 100% application of OH poison (except for 1%? It can resist extremely rarely)
				switch rogue.Consumes.OffHandImbue {
				case proto.WeaponImbue_InstantPoison:
					rogue.InstantPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_DeadlyPoison:
					rogue.DeadlyPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_WoundPoison:
					rogue.WoundPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_OccultPoison:
					rogue.OccultPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_SebaciousPoison:
					rogue.SebaciousPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_AtrophicPoison:
					rogue.AtrophicPoison[ShivProc].Cast(sim, target)
				case proto.WeaponImbue_NumbingPoison:
					rogue.NumbingPoison[ShivProc].Cast(sim, target)
				// Add new alternative poisons as they are implemented
				default:
					if hasDeadlyBrew {
						rogue.InstantPoison[DeadlyBrewProc].Cast(sim, target)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
