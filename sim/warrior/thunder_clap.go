package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerThunderClapSpell() {
	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit, Level int32) *core.Aura {
		return core.ThunderClapAura(target, warrior.Talents.ImprovedThunderClap, warrior.Level)
	})

	damageMultiplier := 1.0

	if warrior.HasRune(proto.WarriorRune_RuneFuriousThunder) {
		damageMultiplier = 2.0
	}

	baseDamage := map[int32]float64{
		25: 23,
		40: 55,
		50: 82,
		60: 103,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 8198,
		40: 8205,
		50: 11580,
		60: 11581,
	}[warrior.Level]

	numHits := min(4, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 20 - []float64{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap] - warrior.FocusedRageDiscount,
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
			if warrior.HasRune(proto.WarriorRune_RuneFuriousThunder) {
				return true
			}
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		// Cruelty doesn't apply to Thunder Clap
		BonusCritRating:  (0 - float64(warrior.Talents.Cruelty)*1),
		DamageMultiplier: damageMultiplier,
		CritMultiplier:   warrior.critMultiplier(none),
		ThreatMultiplier: 1.85,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				if results[hitIndex].Landed() {
					warrior.ThunderClapAuras.Get(target).Activate(sim)
				}
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},

		RelatedAuras: []core.AuraArray{warrior.ThunderClapAuras},
	})
}
