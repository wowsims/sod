package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerDivineFavorSpellAndAura() {
	if !paladin.Talents.DivineFavor {
		return
	}
	var affectedSpells []*core.Spell

	cdTimer := paladin.NewTimer()
	cd := time.Minute * 2

	auraActionID := core.ActionID{SpellID: 20216}

	paladin.DivineFavorAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Favor",
		ActionID: auraActionID,
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{
					paladin.HolyShock,
				}), func(spell *core.Spell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.BonusCritRating += core.CritRatingPerCritChance * 100
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.BonusCritRating -= core.CritRatingPerCritChance * 100
			})
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode != SpellCode_PaladinHolyShock {
				return
			}
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			paladin.UpdateMajorCooldowns()
		},
	})

	paladin.DivineFavor = paladin.RegisterSpell(core.SpellConfig{
		ActionID: auraActionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.DivineFavorAura.Activate(sim)
		},
	})
	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: paladin.DivineFavor,
		Type:  core.CooldownTypeDPS,
	})
}
