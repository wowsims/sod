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

	baseDamage := 113 / 100 * hunter.baseRuneAbilityDamage()

	hasCobraStrikes := hunter.pet != nil && hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes)

	// Efficiency talent doesn't apply to this spell even though it has 'shot' in the name
	manaCostMultiplier := 1.0
	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		manaCostMultiplier -= 0.25
	}

	hunter.KillShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.HunterRune_RuneLegsKillShot)},
		SpellSchool:  core.SpellSchoolPhysical,
		DefenseType:  core.DefenseTypeRanged,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		CastType:     proto.CastType_CastTypeRanged,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.03,
			Multiplier: manaCostMultiplier,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.IsExecutePhase20() {
				spell.CD.Reset()
			}
			
			damage := hunter.AutoAttacks.Ranged().CalculateWeaponDamage(sim, spell.RangedAttackPower(target)) + hunter.AmmoDamageBonus + baseDamage
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)

				// For some reason it doesn't count as a 'shot' ability for efficiency talent but it does count for the 'cobra strikes' rune
				if result.Landed() {
					if hasCobraStrikes && result.DidCrit() {
						hunter.CobraStrikesAura.Activate(sim)
						hunter.CobraStrikesAura.SetStacks(sim, 2)
					}
				}
			})
		},
	})
}
