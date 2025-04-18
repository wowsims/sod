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
	apCoef := 0.07
	attackSpeedReduction := int32(10)
	stanceMask := BattleStance

	if hasFuriousThunder {
		damageMultiplier *= 2
		threatMultiplier *= 1.75
		attackSpeedReduction += 6
		stanceMask = AnyStance
	}

	// Engraving Defense Specialization onto your Ring will grant your Thunder Clap 3% increased chance to hit, to make up for the fact that Weapon Skill does not affect it normally.
	if warrior.Character.HasRuneById(int32(proto.RingRune_RuneRingDefenseSpecialization)) {
		warrior.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusHit_Flat,
			ClassMask:  ClassSpellMask_WarriorThunderClap,
			FloatValue: 3.0,
		})
	}

	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit, Level int32) *core.Aura {
		return core.ThunderClapAura(target, info.spellID, info.duration, attackSpeedReduction)
	})

	results := make([]*core.SpellResult, min(4, warrior.Env.GetNumTargets()))

	warrior.ThunderClap = warrior.RegisterSpell(stanceMask, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: info.spellID},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeRanged,
		ProcMask:    core.ProcMaskRanged,
		Flags:       core.SpellFlagAPL | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost: 20 - []float64{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 4,
			},
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: threatMultiplier,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, info.baseDamage+apCoef*spell.MeleeAttackPower(), spell.OutcomeRangedHitAndCrit)
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
