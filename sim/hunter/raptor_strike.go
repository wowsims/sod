package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if hunter.curQueuedAutoSpell != nil && hunter.curQueuedAutoSpell.CanCast(sim, hunter.CurrentTarget) {
		return hunter.curQueuedAutoSpell
	}
	return mhSwingSpell
}

func (hunter *Hunter) getRaptorStrikeConfig(rank int) core.SpellConfig {
	spellId := [9]int32{0, 2973, 14260, 14261, 14262, 14263, 14264, 14265, 14266}[rank]
	baseDamage := [9]float64{0, 5, 11, 21, 34, 50, 80, 110, 140}[rank]
	manaCost := [9]float64{0, 15, 25, 35, 45, 55, 70, 80, 100}[rank]
	level := [9]int{0, 1, 8, 16, 24, 32, 40, 48, 56}[rank]
	hasFlankingStrike := hunter.HasRune(proto.HunterRune_RuneLegsFlankingStrike)
	hasRaptorFury := hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury)
	hasDualWieldSpec := hunter.HasRune(proto.HunterRune_RuneBootsDualWieldSpecialization)
	hasMeleeSpecialist := hunter.HasRune(proto.HunterRune_RuneBeltMeleeSpecialist)

	// https://www.wowhead.com/classic/news/class-tuning-incoming-hunter-shaman-warlock-season-of-discovery-339072?webhook
	raptorFuryDmgMult := 0.1

	if hasMeleeSpecialist {
		spellId = [9]int32{0, 415335, 415336, 415337, 415338, 415340, 415341, 415342, 415343}[rank]
	}

	hasOHSpell := hasDualWieldSpec && hunter.AutoAttacks.IsDualWielding

	var ohSpell *core.Spell
	if hasOHSpell {
		ohSpell = hunter.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellId}.WithTag(2),
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeOHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

			BonusCritRating: float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,

			DamageMultiplier: 1.5,
			BonusCoefficient: 1,
		})
	}

	spellConfig := core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:         core.SpellFlagMeleeMetrics,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * time.Duration(core.TernaryInt(hasMeleeSpecialist, 3, 6)),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5
		},

		BonusCritRating:  float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			multiplier := 1.0
			if stacks := hunter.RaptorFuryAura.GetStacks(); stacks > 0 {
				multiplier *= 1 + raptorFuryDmgMult*float64(stacks)
			}

			var weaponDamage float64
			if hasMeleeSpecialist {
				weaponDamage = spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				weaponDamage = spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			}

			damage := multiplier * (weaponDamage + baseDamage)
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if ohSpell != nil {
				ohSpell.Cast(sim, target)

				var weaponDamage float64
				if hasMeleeSpecialist {
					weaponDamage = ohSpell.Unit.OHNormalizedWeaponDamage(sim, ohSpell.MeleeAttackPower())
				} else {
					weaponDamage = ohSpell.Unit.OHWeaponDamage(sim, ohSpell.MeleeAttackPower())
				}

				damage := multiplier * (weaponDamage + baseDamage*0.5)
				ohSpell.CalcAndDealDamage(sim, target, damage, ohSpell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
			}

			if hasFlankingStrike && sim.RandomFloat("Flanking Strike Refresh") < 0.2 {
				hunter.FlankingStrike.CD.Set(sim.CurrentTime)
			}

			if hasMeleeSpecialist && sim.RandomFloat("Raptor Strike Reset") < 0.3 {
				spell.CD.Reset()
			}

			if hasRaptorFury {
				if !hunter.RaptorFuryAura.IsActive() {
					hunter.RaptorFuryAura.Activate(sim)
				}
				hunter.RaptorFuryAura.AddStack(sim)
			}
		},
	}

	if hasMeleeSpecialist {
		spellConfig.ProcMask ^= core.ProcMaskMeleeMHAuto
		spellConfig.Flags |= core.SpellFlagAPL
		spellConfig.Cast.DefaultCast = core.Cast{
			GCD: core.GCDDefault,
		}
	}

	return spellConfig
}

func (hunter *Hunter) makeQueueSpellsAndAura(srcSpell *core.Spell) *core.Spell {
	queueAura := hunter.RegisterAura(core.Aura{
		Label:    "RaptorStrikeQueue" + srcSpell.ActionID.String(),
		ActionID: srcSpell.ActionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
			}
			hunter.PseudoStats.DisableDWMissPenalty = true
			hunter.curQueueAura = aura
			hunter.curQueuedAutoSpell = srcSpell
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.DisableDWMissPenalty = false
			hunter.curQueueAura = nil
			hunter.curQueuedAutoSpell = nil
		},
	})

	queueSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: srcSpell.WithTag(1),
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.curQueueAura != queueAura &&
				hunter.CurrentMana() >= srcSpell.DefaultCast.Cost &&
				sim.CurrentTime >= hunter.Hardcast.Expires &&
				hunter.DistanceFromTarget <= 5 &&
				srcSpell.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			queueAura.Activate(sim)
		},
	})
	queueSpell.CdSpell = srcSpell

	return queueSpell
}

func (hunter *Hunter) registerRaptorStrikeSpell() {
	maxRank := 8

	hasMeleeSpecialist := hunter.HasRune(proto.HunterRune_RuneBeltMeleeSpecialist)
	for i := 1; i <= maxRank; i++ {
		config := hunter.getRaptorStrikeConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.RaptorStrike = hunter.GetOrRegisterSpell(config)
			if !hasMeleeSpecialist {
				hunter.makeQueueSpellsAndAura(hunter.RaptorStrike)
			}
		}
	}
}
