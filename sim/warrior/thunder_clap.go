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

	results := make([]*core.SpellResult, min(4, warrior.Env.GetNumTargets()))

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: info.spellID},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

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
			return hasFuriousThunder || warrior.StanceMatches(BattleStance) || warrior.StanceMatches(GladiatorStance)
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: core.TernaryFloat64(hasFuriousThunder, 2, 1),
		ThreatMultiplier: core.TernaryFloat64(hasFuriousThunder, 2.5*1.5, 2.5),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, info.baseDamage, spell.OutcomeMagicHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
				if result.Landed() {
					warrior.ThunderClapAuras.Get(result.Target).Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{warrior.ThunderClapAuras},
	})
}
