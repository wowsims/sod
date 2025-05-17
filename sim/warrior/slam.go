package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerSlamSpell() {
	if warrior.Level < 30 {
		return
	}

	hasBloodSurgeRune := warrior.HasRune(proto.WarriorRune_RuneBloodSurge)

	var castTime time.Duration
	var cooldown core.Cooldown
	if warrior.HasRune(proto.WarriorRune_RunePreciseTiming) {
		castTime = 0
		cooldown = core.Cooldown{
			Timer:    warrior.NewTimer(),
			Duration: 6 * time.Second,
		}
	} else {
		castTime = time.Millisecond*1500 - time.Millisecond*100*time.Duration(warrior.Talents.ImprovedSlam)
	}

	requiredLevel := map[int32]int{
		40: 38,
		50: 46,
		60: 54,
	}[warrior.Level]

	spellID := map[int32]int32{
		40: 8820,
		50: 11604,
		60: 11605,
	}[warrior.Level]

	warrior.SlamMH = warrior.newSlamHitSpell(true)
	canHitOffhand := hasBloodSurgeRune && warrior.AutoAttacks.IsDualWielding
	if canHitOffhand {
		warrior.SlamOH = warrior.newSlamHitSpell(false)
	}

	warrior.Slam = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_WarriorSlam,
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagOffensive,

		RequiredLevel: requiredLevel,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
			CD: cooldown,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if spell.CastTime() > 0 {
					warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, true)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if canHitOffhand && warrior.BloodSurgeAura.IsActive() {
				warrior.SlamOH.Cast(sim, target)
			}
			warrior.SlamMH.Cast(sim, target)
		},
	})
}

func (warrior *Warrior) newSlamHitSpell(isMH bool) *WarriorSpell {
	spellID := map[int32]int32{
		40: 8820,
		50: 11604,
		60: 11605,
	}[warrior.Level]

	requiredLevel := map[int32]float64{
		40: 38,
		50: 46,
		60: 54,
	}[warrior.Level]

	flatDamageBonus := map[int32]float64{
		40: 43,
		50: 68,
		60: 87,
	}[warrior.Level]

	flags := core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete
	damageFunc := warrior.MHWeaponDamage
	if !isMH {
		flatDamageBonus /= 2
		flags |= core.SpellFlagPassiveSpell
		damageFunc = warrior.OHWeaponDamage
	}

	return warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ClassSpellMask: core.Ternary(isMH, ClassSpellMask_WarriorSlamMH, ClassSpellMask_WarriorSlamOH),
		ActionID:       core.ActionID{SpellID: spellID}.WithTag(int32(core.Ternary(isMH, 1, 2))),
		SpellSchool:    core.SpellSchoolPhysical,
		CastType:       core.Ternary(isMH, proto.CastType_CastTypeMainHand, proto.CastType_CastTypeOffHand),
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.Ternary(isMH, core.ProcMaskMeleeMHSpecial, core.ProcMaskMeleeOHSpecial),
		Flags:          flags,

		CritDamageBonus: warrior.impale(),
		FlatThreatBonus: 1 * requiredLevel,

		DamageMultiplier: core.Ternary(isMH, 1.0, warrior.AutoAttacks.OHConfig().DamageMultiplier),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus + damageFunc(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if isMH && !result.Landed() {
				warrior.Slam.IssueRefund(sim)
			}
		},
	})
}
