package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerThunderClapSpell() {
	hasFuriousThunder := warrior.HasRune(proto.WarriorRune_RuneFuriousThunder)

	info := map[int32]struct {
		spellID    int32
		baseDamage float64
		duration   time.Duration
	}{
		25: {spellID: 8198, baseDamage: 23, duration: time.Second * 14},
		40: {spellID: 8205, baseDamage: 55, duration: time.Second * 22},
		50: {spellID: 11580, baseDamage: 82, duration: time.Second * 26},
		60: {spellID: 11581, baseDamage: 103, duration: time.Second * 30},
	}[warrior.Level]

	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit, Level int32) *core.Aura {
		return core.ThunderClapAura(target, info.spellID, info.duration, core.TernaryInt32(hasFuriousThunder, 16, 10))
	})

	numHits := min(4, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: info.spellID},
		SpellSchool: core.SpellSchoolPhysical,
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
			return hasFuriousThunder || warrior.StanceMatches(BattleStance)
		},

		DamageMultiplier: core.TernaryFloat64(hasFuriousThunder, 2, 1),
		CritMultiplier:   warrior.SpellCritMultiplier(1, 0),
		ThreatMultiplier: core.TernaryFloat64(hasFuriousThunder, 2.5*1.5, 2.5),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				results[hitIndex] = spell.CalcDamage(sim, curTarget, info.baseDamage, spell.OutcomeMagicHitAndCrit)
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
