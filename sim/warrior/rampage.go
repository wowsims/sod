package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Go on a rampage, increasing your attack power by 10% for 30 sec.  This ability can only be used while Enraged.
func (warrior *Warrior) registerRampage() {
	if !warrior.HasRune(proto.WarriorRune_RuneRampage) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.WarriorRune_RuneRampage)}
	statDep := warrior.NewDynamicMultiplyStat(stats.AttackPower, 1.10)

	warrior.RampageAura = warrior.RegisterAura(core.Aura{
		Label:    "Rampage",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.EnableDynamicStatDep(sim, statDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.DisableDynamicStatDep(sim, statDep)
		},
	})

	warrior.Rampage = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL | core.SpellFlagCastTimeNoGCD,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.IsEnraged()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.RampageAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Type:  core.CooldownTypeDPS,
		Spell: warrior.Rampage.Spell,
	})
}
