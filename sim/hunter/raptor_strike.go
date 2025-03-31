package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const RaptorStrikeRanks = 8

var RaptorStrikeSpellId = [RaptorStrikeRanks + 1]int32{0, 2973, 14260, 14261, 14262, 14263, 14264, 14265, 14266}
var RaptorStrikeSpellIdMeleeSpecialist = [RaptorStrikeRanks + 1]int32{0, 415335, 415336, 415337, 415338, 415340, 415341, 415342, 415343}
var RaptorStrikeBaseDamage = [RaptorStrikeRanks + 1]float64{0, 5, 11, 21, 34, 50, 80, 110, 140}
var RaptorStrikeManaCost = [RaptorStrikeRanks + 1]float64{0, 15, 25, 35, 45, 55, 70, 80, 100}
var RaptorStrikeLevel = [RaptorStrikeRanks + 1]int{0, 1, 8, 16, 24, 32, 40, 48, 56}

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if hunter.curQueuedAutoSpell != nil && hunter.curQueuedAutoSpell.CanCast(sim, hunter.CurrentTarget) {
		return hunter.curQueuedAutoSpell
	}
	return mhSwingSpell
}

func (hunter *Hunter) getRaptorStrikeConfig(rank int) core.SpellConfig {
	hasRaptorFury := hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury)
	hasDualWieldSpec := hunter.HasRune(proto.HunterRune_RuneBootsDualWieldSpecialization)
	hasMeleeSpecialist := hunter.HasRune(proto.HunterRune_RuneBeltMeleeSpecialist)

	spellID := core.Ternary(hasMeleeSpecialist, RaptorStrikeSpellIdMeleeSpecialist, RaptorStrikeSpellId)[rank]
	manaCost := RaptorStrikeManaCost[rank]
	level := RaptorStrikeLevel[rank]

	hunter.RaptorStrikeMH = hunter.newRaptorStrikeHitSpell(rank, true)
	if hasDualWieldSpec {
		hunter.RaptorStrikeOH = hunter.newRaptorStrikeHitSpell(rank, false)
	}

	spellConfig := core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterRaptorStrike,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagStrike,

		Rank:          rank,
		RequiredLevel: level,
		MaxRange:      core.MaxMeleeAttackRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hunter.RaptorStrikeMH.Cast(sim, target)
			if hasDualWieldSpec && hunter.AutoAttacks.IsDualWielding {
				hunter.RaptorStrikeOH.Cast(sim, target)
			}

			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
			}

			if hasMeleeSpecialist && sim.Proc(0.3, "Raptor Strike Reset") {
				spell.CD.Reset()
				hunter.MongooseBite.CD.Reset()
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
		spellConfig.Flags |= core.SpellFlagAPL
		spellConfig.Cast.DefaultCast = core.Cast{
			GCD: core.GCDDefault,
		}
	} else {
		spellConfig.ProcMask |= core.ProcMaskMeleeMHAuto
	}

	return spellConfig
}

func (hunter *Hunter) newRaptorStrikeHitSpell(rank int, isMH bool) *core.Spell {
	hasMeleeSpecialist := hunter.HasRune(proto.HunterRune_RuneBeltMeleeSpecialist)
	hasRaptorFury := hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury)
	hasHitAndRun := hunter.HasRune(proto.HunterRune_RuneCloakHitAndRun)

	spellID := core.Ternary(hasMeleeSpecialist, RaptorStrikeSpellIdMeleeSpecialist, RaptorStrikeSpellId)[rank]
	baseDamage := RaptorStrikeBaseDamage[rank]

	castType := proto.CastType_CastTypeMainHand
	procMask := core.ProcMaskMeleeMHSpecial
	damageMultiplier := 1.0
	damageFunc := hunter.MHWeaponDamage

	if !isMH {
		baseDamage /= 2
		castType = proto.CastType_CastTypeOffHand
		procMask = core.ProcMaskMeleeOHSpecial
		damageMultiplier = hunter.AutoAttacks.OHConfig().DamageMultiplier
		damageFunc = hunter.OHNormalizedWeaponDamage // https://www.wowhead.com/classic-ptr/spell=409755/raptor-strike
	}

	return hunter.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterRaptorStrikeHit,
		ActionID:       core.ActionID{SpellID: spellID}.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:    core.SpellSchoolPhysical,
		CastType:       castType,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       procMask,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		BonusCritRating:  float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,
		CritDamageBonus:  hunter.mortalShots(),
		DamageMultiplier: damageMultiplier,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if hasHitAndRun {
				hunter.HitAndRunAura.Activate(sim)
			}

			multiplier := 1.0
			if hasRaptorFury {
				multiplier *= hunter.raptorFuryDamageMultiplier()
			}

			weaponDamage := damageFunc(sim, spell.MeleeAttackPower())
			damage := multiplier * (weaponDamage + baseDamage)
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}

func (hunter *Hunter) makeQueueSpellsAndAura() *core.Spell {
	queueAura := hunter.RegisterAura(core.Aura{
		Label:    "Raptor Strike Queued",
		ActionID: hunter.RaptorStrike.ActionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
			}
			hunter.PseudoStats.DisableDWMissPenalty = true
			hunter.curQueueAura = aura
			hunter.curQueuedAutoSpell = hunter.RaptorStrike
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.DisableDWMissPenalty = false
			hunter.curQueueAura = nil
			hunter.curQueuedAutoSpell = nil
		},
	})

	queueSpell := hunter.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterRaptorStrike,
		ActionID:       hunter.RaptorStrike.WithTag(3),
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		MaxRange: core.MaxMeleeAttackRange,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.curQueueAura != queueAura &&
				hunter.CurrentMana() >= hunter.RaptorStrike.Cost.GetCurrentCost() &&
				!hunter.IsCasting(sim) &&
				hunter.RaptorStrike.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			queueAura.Activate(sim)
		},
	})
	queueSpell.CdSpell = hunter.RaptorStrike

	return queueSpell
}

func (hunter *Hunter) registerRaptorStrikeSpell() {
	rank := map[int32]int{
		25: 4,
		40: 6,
		50: 7,
		60: 8,
	}[hunter.Level]

	hasMeleeSpecialist := hunter.HasRune(proto.HunterRune_RuneBeltMeleeSpecialist)
	config := hunter.getRaptorStrikeConfig(rank)
	hunter.RaptorStrike = hunter.GetOrRegisterSpell(config)

	if !hasMeleeSpecialist {
		hunter.makeQueueSpellsAndAura()
	}
}
