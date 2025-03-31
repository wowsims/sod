package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerDivineFavor() {
	if !paladin.Talents.DivineFavor {
		return
	}

	var affectedSpells []*core.Spell
	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_PaladinHolyShock) {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Minute * 2,
	}

	aura := paladin.RegisterAura(core.Aura{
		Label:    "Divine Favor",
		ActionID: core.ActionID{SpellID: 20216},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.BonusCritRating += core.SpellCritRatingPerCritChance * 100
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.BonusCritRating -= core.SpellCritRatingPerCritChance * 100
			})
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(ClassSpellMask_PaladinHolyShock) {
				return
			}
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cd.Set(sim.CurrentTime + cd.Duration)
			paladin.UpdateMajorCooldowns()
		},
	})

	divineFavor := paladin.RegisterSpell(core.SpellConfig{
		ActionID: aura.ActionID,
		Cast: core.CastConfig{
			CD: cd,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},
	})
	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: divineFavor,
		Type:  core.CooldownTypeDPS,
	})
}
