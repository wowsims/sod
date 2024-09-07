package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"time"
)

func (paladin *Paladin) registerAvengersShield() {
	if !paladin.hasRune(proto.PaladinRune_RuneLegsAvengersShield) {
		return
	}
    
    // Avenger's Shield hits up to 3 targets.
	results := make([]*core.SpellResult, min(3, paladin.Env.GetNumTargets()))

	paladin.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407669},
        SpellCode:   SpellCode_PaladinAvengersShield,
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.26,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 15,
			},
		},
        ExtraCastCondition: // have a shield equipped
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
        BonusCoefficient: 0.091, // for spell damage; we add the AP bonus manually

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			lowDamage := 366 * paladin.baseRuneAbilityDamage() / 100
            highDamage := 448 * paladin.baseRuneAbilityDamage() / 100
            apBonus := 0.091 * spell.MeleeAttackPower()
			for idx := range results {
                baseDamage := sim.Roll(lowDamage, highDamage) + apBonus
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
