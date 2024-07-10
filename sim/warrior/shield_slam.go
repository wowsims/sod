package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warrior *Warrior) registerShieldSlamSpell() {
	if !warrior.Talents.ShieldSlam {
		return
	}

	rank := map[int32]struct {
		spellID    int32
		damageLow  float64
		damageHigh float64
		threat     float64
	}{
		40: {spellID: 23922, damageLow: 225, damageHigh: 235, threat: 178},
		50: {spellID: 23923, damageLow: 264, damageHigh: 276, threat: 203},
		60: {spellID: 23925, damageLow: 342, damageHigh: 358, threat: 254},
	}[warrior.Level]

	apCoef := 0.15

	defendersResolveAura := core.DefendersResolveAttackPower(warrior.GetCharacter())

	warrior.ShieldSlam = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		SpellCode:   SpellCode_WarriorShieldSlam,
		ActionID:    core.ActionID{SpellID: rank.spellID},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial, // TODO really?
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   20 - warrior.FocusedRageDiscount,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 2,
		FlatThreatBonus:  rank.threat * 2,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(rank.damageLow, rank.damageHigh) + warrior.BlockValue()*2 + apCoef*spell.MeleeAttackPower()
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				if stacks := int32(warrior.GetStat(stats.Defense)); stacks > 0 {
					if !defendersResolveAura.IsActive() {
						defendersResolveAura.Activate(sim)
					}

					if defendersResolveAura.GetStacks() != stacks {
						defendersResolveAura.SetStacks(sim, stacks)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
