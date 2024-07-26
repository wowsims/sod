package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Thunder Clap now increases the time between attacks by an additional 6%, can be used in any stance, deals 100% increased damage, and deals 50% increased threat.
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

	damageMultiplier := 1.0
	threatMultiplier := 2.5
	apCoef := 0.05
	attackSpeedReduction := int32(10)
	stanceMask := BattleStance

	if hasFuriousThunder {
		damageMultiplier *= 2
		threatMultiplier *= 1.5
		apCoef *= 2
		attackSpeedReduction += 6
		stanceMask = AnyStance
	}

	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit, Level int32) *core.Aura {
		return core.ThunderClapAura(target, info.spellID, info.duration, attackSpeedReduction)
	})

	results := make([]*core.SpellResult, min(4, warrior.Env.GetNumTargets()))

	warrior.ThunderClap = warrior.RegisterSpell(stanceMask, core.SpellConfig{
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

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: threatMultiplier,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, info.baseDamage+apCoef*spell.MeleeAttackPower(), spell.OutcomeMagicHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
				if result.Landed() {
					// TODO: Thunder Clap now increases the time between attacks by an additional 6%
					warrior.ThunderClapAuras.Get(result.Target).Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{warrior.ThunderClapAuras},
	})
}
