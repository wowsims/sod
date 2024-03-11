package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ShurikenTossSpellID int32 = int32(proto.RogueRune_RuneShurikenToss)

func (rogue *Rogue) makeShurikenTossHitSpell() *core.Spell {
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: ShurikenTossSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1,
		CritMultiplier:   rogue.RangedCritMultiplier(true),
		ThreatMultiplier: 1,
		CastType:         proto.CastType_CastTypeRanged,
	})
}

func (rogue *Rogue) registerShurikenTossSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneShurikenToss) {
		return
	}
	results := make([]*core.SpellResult, 5)

	//hit := rogue.makeShurikenTossHitSpell()

	rogue.ShurikenToss = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: ShurikenTossSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		CastType: proto.CastType_CastTypeRanged,

		DamageMultiplier: 1,
		CritMultiplier:   rogue.RangedCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			// Deal damage to main target
			baseDamage := spell.MeleeAttackPower() * 0.15
			results[0] = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			// TODO: Make bounces work without triggering "Infinte loop detected"
			// calc and apply for target and up to 4 other targets
			/**hits := 0
			maxHits := 4
			currentTarget := target.Index
			// Find any additional targets up to the bounce limit
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				if hits >= maxHits {
					break
				}
				if aoeTarget.Index != currentTarget {
					results[hits+1] = hit.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeRangedHitAndCrit)
					hits++
				}
			}*/
			// only deal up to number of hits dealt
			/**for i := 0; i < hits; i++ {
				hit.DealDamage(sim, results[i+1])
			}*/
			if results[0].Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			}
		},
	})
}
