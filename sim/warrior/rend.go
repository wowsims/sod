package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Blood Frenzy
// Rend can now be used in Berserker stance, Rend's damage is increased by 100%,
// and Rend deals additional damage equal to 3% of your Attack Power each time it deals damage.

func (warrior *Warrior) registerRendSpell() {
	hasBloodFrenzyRune := warrior.HasRune(proto.WarriorRune_RuneBloodFrenzy)

	rend := map[int32]struct {
		ticks   int32
		damage  float64
		spellID int32
	}{
		25: {spellID: 6547, damage: 9, ticks: 5},
		40: {spellID: 11572, damage: 14, ticks: 7},
		50: {spellID: 11573, damage: 18, ticks: 7},
		60: {spellID: 11574, damage: 21, ticks: 7},
	}[warrior.Level]

	damageMultiplier := []float64{1, 1.15, 1.25, 1.35}[warrior.Talents.ImprovedRend]
	if hasBloodFrenzyRune {
		damageMultiplier *= 2
	}

	warrior.Rend = warrior.RegisterSpell(BattleStance|DefensiveStance, core.SpellConfig{
		SpellCode:   SpellCode_WarriorRend,
		ActionID:    core.ActionID{SpellID: rend.spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagAPL | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		RageCost: core.RageCostOptions{
			Cost:   10 - warrior.FocusedRageDiscount,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rend",
				Tag:   "Rend",
			},
			NumberOfTicks: rend.ticks,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				damage := rend.damage
				if hasBloodFrenzyRune {
					damage += .03 * dot.Spell.MeleeAttackPower()
				}

				dot.Snapshot(target, damage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealOutcome(sim, result)
		},
	})

}
