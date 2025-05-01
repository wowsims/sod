package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerKillShotSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsKillShot) {
		return
	}

	hunter.KillShot = hunter.RegisterSpell(hunter.newKillShotConfig())
}

func (hunter *Hunter) newKillShotConfig() core.SpellConfig {
	baseDamage := 113 / 100 * hunter.baseRuneAbilityDamage() * 5.21

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterKillShot,
		ActionID:       core.ActionID{SpellID: int32(proto.HunterRune_RuneLegsKillShot)},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeRanged,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		CastType:       proto.CastType_CastTypeRanged,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.IsExecutePhase20() && spell.CD.Duration != 0 {
				spell.CD.Reset()
			}

			damage := hunter.AutoAttacks.Ranged().CalculateWeaponDamage(sim, spell.RangedAttackPower(target, false)) + hunter.AmmoDamageBonus + baseDamage
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	}
}
