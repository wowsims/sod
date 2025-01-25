package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) ApplyRunes() {
	// Head
	priest.registerEyeOfTheVoidCD()
	priest.applyPainAndSuffering()

	// Shoulders
	priest.applyShoulderRuneEffect()

	// Cloak
	priest.registerVampiricTouchSpell()

	// Chest
	priest.registerVoidPlagueSpell()
	priest.applyTwistedFaith()

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

func (priest *Priest) applyShoulderRuneEffect() {
	if priest.Equipment.Shoulders().Rune == int32(proto.PriestRune_PriestRuneNone) {
		return
	}

	switch priest.Equipment.Shoulders().Rune {
	// Shadow
	case int32(proto.PriestRune_RuneShouldersRefinedPriest):
		priest.applyT1Shadow4PBonus()
	case int32(proto.PriestRune_RuneShouldersMindBreaker):
		priest.applyT1Shadow6PBonus()
	case int32(proto.PriestRune_RuneShouldersDeathdealer):
		priest.applyT2Shadow2PBonus()
	case int32(proto.PriestRune_RuneShouldersSpiritFont):
		priest.applyT2Shadow4PBonus()
	case int32(proto.PriestRune_RuneShouldersZealot):
		priest.applyT2Shadow6PBonus()
	case int32(proto.PriestRune_RuneShouldersUnwaveringDefiler):
		priest.applyTAQShadow2PBonus()
	case int32(proto.PriestRune_RuneShouldersContemnor):
		priest.applyTAQShadow4PBonus()
	case int32(proto.PriestRune_RuneShouldersPlaguebringer):
		priest.applyRAQShadow3PBonus()

	// Healer
	case int32(proto.PriestRune_RuneShouldersFaithful):
		priest.applyT2Healer2PBonus()
	case int32(proto.PriestRune_RuneShouldersSerendipitous):
		priest.applyT2Healer4PBonus()
	case int32(proto.PriestRune_RuneShouldersPenitent):
		priest.applyZGDiscipline3PBonus()
	}
}

func (priest *Priest) applyTwistedFaith() {
	if !priest.HasRune(proto.PriestRune_RuneChestTwistedFaith) {
		return
	}

	priest.RegisterAura(core.Aura{
		Label: "Twisted Faith",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := core.Flatten(priest.MindFlay)
			affectedSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{
					affectedSpells,
					priest.MindBlast,
				}), func(spell *core.Spell) bool { return spell != nil },
			)

			swpSpells := core.FilterSlice(priest.ShadowWordPain, func(spell *core.Spell) bool { return spell != nil })

			for _, spell := range affectedSpells {
				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					modifier := 0.0

					for _, spell := range swpSpells {
						if spell.Dot(target).IsActive() {
							modifier += 0.50
							break
						}
					}

					spell.DamageMultiplierAdditive += modifier
					oldApplyEffects(sim, target, spell)
					spell.DamageMultiplierAdditive -= modifier
				}
			}
		},
	})
}

func (priest *Priest) applyPainAndSuffering() {
	if !priest.HasRune(proto.PriestRune_RuneHelmPainAndSuffering) {
		return
	}

	priest.PainAndSufferingDoTSpells = []*core.Spell{}
	affectedSpellClassMasks := ClassSpellMask_PriestMindBlast | ClassSpellMask_PriestMindFlay | ClassSpellMask_PriestMindSpike
	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: "Pain and Suffering Trigger",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			priest.PainAndSufferingDoTSpells = append(
				priest.PainAndSufferingDoTSpells,
				core.FilterSlice(
					core.Flatten(
						[][]*core.Spell{
							priest.ShadowWordPain,
							{priest.VoidPlague, priest.VampiricTouch},
						},
					),
					func(spell *core.Spell) bool { return spell != nil },
				)...,
			)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(affectedSpellClassMasks) && result.Landed() {
				target := result.Target

				var dotToRollover *core.Dot
				dotSpells := core.FilterSlice(priest.PainAndSufferingDoTSpells, func(spell *core.Spell) bool { return spell.Dot(target).IsActive() })

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
	}))
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
			if spell.Matches(ClassSpellMask_PriestSmite) {
				aura.Deactivate(sim)
			}
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if spell.Matches(ClassSpellMask_PriestFlashHeal) {
				aura.Deactivate(sim)
			}
		},
	})

	handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if spell.ProcMask.Matches(core.ProcMaskSpellOrSpellProc) && result.Outcome.Matches(core.OutcomeCrit) {
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
		if spell.Matches(ClassSpellMask_PriestAll) && !spell.Flags.Matches(core.SpellFlagHelpful) {
			spell.CritDamageBonus += 1
		}
	})
}
