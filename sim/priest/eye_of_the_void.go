package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) registerEyeOfTheVoidCD() {
	if !priest.HasRune(proto.PriestRune_RuneHelmEyeOfTheVoid) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.PriestRune_RuneHelmEyeOfTheVoid)}
	duration := time.Second * 30
	cooldown := time.Minute * 3

	// For timeline only
	priest.EyeOfTheVoidAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Eye of the Void",
		Duration: duration,
	})

	priest.EyeOfTheVoid = priest.RegisterSpell(core.SpellConfig{
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
			priest.EyeOfTheVoidPet.EnableWithTimeout(sim, priest.EyeOfTheVoidPet, duration)
			priest.EyeOfTheVoidAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    priest.EyeOfTheVoid,
		Priority: 1,
		Type:     core.CooldownTypeDPS,
	})
}
