package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warrior *Warrior) registerBerserkerRageSpell() {
	if warrior.Level < 32 {
		return
	}

	actionID := core.ActionID{SpellID: 18499}
	rageMetrics := warrior.NewRageMetrics(actionID)
	rageMultiplier := 1.0

	warrior.BerserkerRageAura = warrior.RegisterAura(core.Aura{
		Label:    "Berserker Rage",
		ActionID: actionID,
		Duration: time.Second * 10,

		// Copy from rage.go
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			rageConversionTaken := 0.0091107836*float64(spell.Unit.Level^2) + 3.225598133*float64(spell.Unit.Level) + 4.2652911
			generatedRage := result.Damage * 2.5 / rageConversionTaken
			generatedRage *= rageMultiplier
			warrior.AddRage(sim, generatedRage, rageMetrics)
		},
	})

	warrior.BerserkerRage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			warrior.BerserkerRageAura.Activate(sim)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance) || warrior.StanceMatches(GladiatorStance)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.BerserkerRage,
		Type:  core.CooldownTypeDPS,
	})
}
