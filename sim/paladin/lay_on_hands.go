package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerLayOnHands() {

	manaReturn := 0

	layOnHands := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagAPL | core.SpellFlagMCD,
		SpellSchool: core.SpellSchoolHoly,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute*60 - 10*int(paladin.Talents.ImprovedLayOnHands),
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			paladin.SpendMana(sim, paladin.CurrentMana(), layOnHandsManaMetrics)
			target.GainHealth(sim, paladin.MaxHealth(), layOnHandsHealthMetrics)
			target.AddMana(sim, manaReturn, layOnHandsManaMetrics)
		},
	})
}
