package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) getShadowCleaveBaseConfig(rank int) core.SpellConfig {
	spellCoeff := [11]float64{0, .047, .1, .187, .286, .286, .286, .286, .286, .286, .286}[rank]
	baseDamage := [11][]float64{{0}, {3, 7}, {7, 12}, {14, 23}, {26, 39}, {42, 64}, {60, 91}, {82, 124}, {105, 158}, {128, 193}, {136, 204}}[rank]
	spellId := [11]int32{0, 403835, 403839, 403840, 403841, 403842, 403843, 403844, 403848, 403851, 403852}[rank]
	manaCost := [11]float64{0, 12, 20, 35, 55, 80, 105, 132, 157, 185, 190}[rank]
	level := [11]int{0, 1, 6, 12, 20, 28, 36, 44, 52, 60, 60}[rank]

	results := make([]*core.SpellResult, min(10, warlock.Env.GetNumTargets()))

	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolShadow,
		ClassSpellMask: ClassSpellMask_WarlockShadowCleave,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDestruction,
		RequiredLevel:  level,
		Rank:           rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.MetamorphosisAura.IsActive()
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 2, // Undocumented 2x multiplier
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				damage := sim.Roll(baseDamage[0], baseDamage[1])
				results[idx] = spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			hasHit := false
			for _, result := range results {
				if result.Landed() {
					hasHit = true
					spell.DealDamage(sim, result)
				}
			}

			if stacks := int32(warlock.GetStat(stats.Defense)); hasHit && stacks > 0 {
				warlock.defendersResolveAura.Activate(sim)
				if warlock.defendersResolveAura.GetStacks() != stacks {
					warlock.defendersResolveAura.SetStacks(sim, stacks)
				}
			}
		},
	}
}

const DefendersResolveSpellDamagePer = 4

func (warlock *Warlock) registerShadowCleaveSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) {
		return
	}

	warlock.defendersResolveAura = core.DefendersResolveSpellDamage(warlock.GetCharacter(), DefendersResolveSpellDamagePer)

	warlock.ShadowCleave = make([]*core.Spell, 0)
	for rank := 1; rank <= ShadowBoltRanks; rank++ {
		config := warlock.getShadowCleaveBaseConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.ShadowCleave = append(warlock.ShadowCleave, warlock.GetOrRegisterSpell(config))
		}
	}
}
