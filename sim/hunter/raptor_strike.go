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

	flankingStrikeDmgMult := 0.1
	raptorFuryDmgMult := 0.15 // TODO: Verify value after launch, has been datamined to possibly get changed to 0.1 instead of 0.15 but until further confirmation

	if hasMeleeSpecialist {
		spellId = [9]int32{0, 415335, 415336, 415337, 415338, 415340, 415341, 415342, 415343}[rank]
	}

	hasOHSpell := hasDualWieldSpec && hunter.AutoAttacks.IsDualWielding

	dwSpecMulti := 1.0
	var ohSpell *core.Spell
	if hasOHSpell {
		ohSpell = hunter.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellId}.WithTag(2),
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeOHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

			BonusCritRating: float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,

			DamageMultiplier: 1.5 * dwSpecMulti,
		})
	}

	spellConfig := core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
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
		DamageMultiplier: 1 * dwSpecMulti,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			flankingStrikeStacks := float64(hunter.FlankingStrikeAura.GetStacks())
			raptorFuryStacks := float64(hunter.RaptorFuryAura.GetStacks())

			weaponDamage := 0.0
			if hasMeleeSpecialist {
				weaponDamage = spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				weaponDamage = spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			}

			mhBaseDamage := baseDamage +
				weaponDamage +
				spell.BonusWeaponDamage()

			if hasFlankingStrike && hunter.FlankingStrikeAura.IsActive() {
				mhBaseDamage *= 1.0 + (flankingStrikeDmgMult * flankingStrikeStacks)
			}

			if hasFlankingStrike && sim.RandomFloat("Flanking Strike Refresh") < 0.2 {
				hunter.FlankingStrike.CD.Set(sim.CurrentTime)
			}

			if hasRaptorFury && hunter.RaptorFuryAura.IsActive() {
				mhBaseDamage *= 1.0 + (raptorFuryDmgMult * raptorFuryStacks)
			}

			spell.CalcAndDealDamage(sim, target, mhBaseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if ohSpell != nil {
				ohSpell.Cast(sim, target)

				ohWeaponDamage := 0.0
				if hasMeleeSpecialist {
					ohWeaponDamage = ohSpell.Unit.OHNormalizedWeaponDamage(sim, ohSpell.MeleeAttackPower())
				} else {
					ohWeaponDamage = ohSpell.Unit.OHWeaponDamage(sim, ohSpell.MeleeAttackPower())
				}

				ohBaseDamage := baseDamage*0.5 +
					ohWeaponDamage +
					ohSpell.BonusWeaponDamage()

				if hasFlankingStrike && hunter.FlankingStrikeAura.IsActive() {
					ohBaseDamage *= 1.0 + (flankingStrikeDmgMult * flankingStrikeStacks)
				}

				if hasRaptorFury && hunter.RaptorFuryAura.IsActive() {
					ohBaseDamage *= 1.0 + (raptorFuryDmgMult * raptorFuryStacks)
				}

				ohSpell.CalcAndDealDamage(sim, target, ohBaseDamage, ohSpell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
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
