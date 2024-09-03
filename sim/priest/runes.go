package priest

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) ApplyRunes() {
	// Head
	priest.registerEyeOfTheVoidCD()
	priest.applyPainAndSuffering()

	// Cloak
	priest.registerVampiricTouchSpell()

	// Chest
	priest.registerVoidPlagueSpell()

	// Bracers
	priest.applySurgeOfLight()
	priest.applyDespair()
	priest.registerVoidZoneSpell()

	// Hands
	priest.registerMindSearSpell()
	priest.RegisterPenanceSpell()
	priest.registerShadowWordDeathSpell()

	// Belt
	priest.registerMindSpikeSpell()

	// Legs
	priest.registerHomunculiSpell()

	// Feet
	priest.registerDispersionSpell()

	// Skill Books
	priest.registerShadowfiendSpell()
}

func (priest *Priest) applyPainAndSuffering() {
	if !priest.HasRune(proto.PriestRune_RuneHelmPainAndSuffering) {
		return
	}

	affectedSpellcodes := []int32{SpellCode_PriestMindBlast, SpellCode_PriestMindFlay, SpellCode_PriestMindSpike}
	priest.RegisterAura(core.Aura{
		Label:    "Pain and Suffering Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if slices.Contains(affectedSpellcodes, spell.SpellCode) && result.Landed() {
				target := result.Target

				var dotToRollover *core.Dot
				dotSpells := core.FilterSlice(
					core.Flatten([][]*core.Spell{priest.ShadowWordPain, {priest.VoidPlague, priest.VampiricTouch}}),
					func(spell *core.Spell) bool {
						return spell != nil && spell.Dot(target).IsActive()
					},
				)

				if len(dotSpells) > 0 {
					dotToRollover = dotSpells[0].Dot(target)
					for _, spell := range dotSpells {
						dot := spell.Dot(target)
						if dot.RemainingDuration(sim) < dotToRollover.RemainingDuration(sim) {
							dotToRollover = dot
						}
					}

					dotToRollover.NumberOfTicks = dotToRollover.OriginalNumberOfTicks
					dotToRollover.Rollover(sim)
				}
			}
		},
	})
}

func (priest *Priest) applySurgeOfLight() {
	if !priest.HasRune(proto.PriestRune_RuneBracersSurgeOfLight) {
		return
	}

	var affectedSpells []*core.Spell

	priest.SurgeOfLightAura = priest.RegisterAura(core.Aura{
		Label:    "Surge of Light Proc",
		ActionID: core.ActionID{SpellID: int32(proto.PriestRune_RuneBracersSurgeOfLight)},
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{priest.Smite, priest.FlashHeal}),
				func(spell *core.Spell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier -= 1
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier += 1
			})
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if spell.SpellCode == SpellCode_PriestSmite {
				aura.Deactivate(sim)
			}
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if spell.SpellCode == SpellCode_PriestFlashHeal {
				aura.Deactivate(sim)
			}
		},
	})

	handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if spell.ProcMask.Matches(core.ProcMaskSpellOrProc) && result.Outcome.Matches(core.OutcomeCrit) {
			priest.SurgeOfLightAura.Activate(sim)
		}
	}

	priest.RegisterAura(core.Aura{
		Label:    "Surge of Light Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: handler,
		OnHealDealt:     handler,
	})
}

func (priest *Priest) applyDespair() {
	if !priest.HasRune(proto.PriestRune_RuneBracersDespair) {
		return
	}

	priest.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagPriest) && !spell.Flags.Matches(core.SpellFlagHelpful) {
			spell.CritDamageBonus += 1
		}
	})
}
