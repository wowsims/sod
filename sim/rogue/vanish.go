package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerVanishSpell() {
	has4Pc := rogue.HasSetBonus(ItemSetNightSlayerBattlearmor, 4)
	
	rogue.VanishAura = rogue.RegisterAura(core.Aura{
		Label:    "Vanish",
		ActionID: core.ActionID{SpellID:457437},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(0.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(1/0.5)
		},
	})
	
	rogue.Vanish = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1856},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(300-45*rogue.Talents.Elusiveness),
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if has4Pc {
				rogue.VanishAura.Activate(sim)
				return
			}
			// Pause auto attacks
			rogue.AutoAttacks.CancelAutoSwing(sim)
			// Apply stealth
			rogue.StealthAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Vanish,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDrums,
	})
}
