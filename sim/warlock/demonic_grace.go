package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) registerDemonicGraceSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsDemonicGrace) {
		return
	}

	warlock.DemonicGraceAura = warlock.RegisterAura(core.Aura{
		Label:    "Demonic Grace Aura",
		ActionID: core.ActionID{SpellID: 425463},
		Duration: time.Second * 6,
	}).AttachStatsBuff(stats.Stats{
		stats.Dodge:     20 * core.DodgeRatingPerDodgeChance,
		stats.MeleeCrit: 20 * core.CritRatingPerCritChance,
		stats.SpellCrit: 20 * core.SpellCritRatingPerCritChance,
	})

	warlock.DemonicGrace = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 425463},
		ClassSpellMask: ClassSpellMask_WarlockDemonicGrace,
		Flags:          core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagDemonology,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 20,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.DemonicGraceAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: warlock.DemonicGrace,
		Type:  core.CooldownTypeDPS,
	})
}
