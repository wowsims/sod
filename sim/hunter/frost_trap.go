package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getFrostTrapConfig(timer *core.Timer) core.SpellConfig {

	hasLockAndLoad := hunter.HasRune(proto.HunterRune_RuneHelmLockAndLoad)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 13809},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | SpellFlagTrap,
		RequiredLevel: 28,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: 60 * hunter.resourcefulnessManacostModifier(),
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * time.Duration(15 * hunter.resourcefulnessCooldownModifier()),
			},
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= hunter.trapRange()
		},

		BonusHitRating: hunter.trapMastery(),

		DamageMultiplier: 1 + 0.15*float64(hunter.Talents.CleverTraps),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				if hasLockAndLoad {
					hunter.LockAndLoadAura.Activate(sim)
				}
			})
		},
	}
}

func (hunter *Hunter) registerFrostTrapSpell(timer *core.Timer) {
	config := hunter.getFrostTrapConfig(timer)

	if config.RequiredLevel <= int(hunter.Level) {
		hunter.FrostTrap = hunter.GetOrRegisterSpell(config)
	}
}
