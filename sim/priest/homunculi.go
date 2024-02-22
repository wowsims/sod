package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) registerHomunculiSpell() {
	if !priest.HasRune(proto.PriestRune_RuneLegsHomunculi) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.PriestRune_RuneLegsHomunculi)}
	duration := time.Minute * 2
	cooldown := time.Minute * 2

	// For timeline only
	priest.HomunculiAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Homunculi",
		Duration: duration,
	})

	priest.Homunculi = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			core.Each(priest.HomunculiPets, func(homunculus *Homunculus) {
				homunculus.EnableWithTimeout(sim, homunculus, duration)
			})
			priest.HomunculiAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    priest.Homunculi,
		Priority: 1,
		Type:     core.CooldownTypeDPS,
	})
}
