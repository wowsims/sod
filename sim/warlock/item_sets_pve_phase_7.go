package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetPlagueheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Plagueheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasDamage2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasDamage4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasDamage6PBonus()
		},
	},
})

// Increases the damage done by your Incinerate and Corruption abilities by 20%.
func (warlock *Warlock) applyNaxxramasDamage2PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Damage 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := warlock.Corruption
			if warlock.Incinerate != nil {
				affectedSpells = append(affectedSpells, warlock.Incinerate)
			}

			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.20
			}
		},
	})
}

// Your non-periodic critical strikes cause your active Corruption, Immolate, Shadowflame, and Unstable Affliction spells on the target to immediately deal one pulse of their damage to the target.
func (warlock *Warlock) applyNaxxramasDamage4PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Damage 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	affectedDots := []*core.Dot{}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			dotSpells := core.FilterSlice(
				core.Flatten([][]*core.Spell{
					warlock.Corruption,
					warlock.Immolate,
					{warlock.Shadowflame, warlock.UnstableAffliction},
				}),
				func(spell *core.Spell) bool { return spell != nil },
			)

			for _, spell := range dotSpells {
				affectedDots = append(affectedDots, core.FilterSlice(spell.Dots(), func(dot *core.Dot) bool { return dot != nil })...)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Flags.Matches(SpellFlagWarlock) && result.DidCrit() {
				for _, dot := range affectedDots {
					if dot.IsActive() {
						dot.TickOnce(sim)
					}
				}
			}
		},
	}))
}

// Your Curse of Agony does not expire on Undead targets, and continues to grow in power indefinitely.
func (warlock *Warlock) applyNaxxramasDamage6PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Damage 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.CurseOfAgony {
				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					if target.MobType == proto.MobType_MobTypeUndead {
						dot := spell.Dot(target)
						dot.NumberOfTicks = 65_536 // Large enough to be effectively infinite for our purposes
						dot.RecomputeAuraDuration()
					}
					oldApplyEffects(sim, target, spell)
				}
			}
		},
	})
}

var ItemSetPlagueheartStitchings = core.NewItemSet(core.ItemSet{
	Name: "Plagueheart Stitchings",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasTank2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasTank4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasTank6PBonus()
		},
	},
})

// Your Menace ability never misses, and your chance to be Dodged or Parried or for your spells to miss is reduced by 2%.
func (warlock *Warlock) applyNaxxramasTank2PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Tank 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 2 * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(bonusStats)
			} else {
				aura.Unit.AddStatsDynamic(sim, bonusStats)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(bonusStats.Invert())
			} else {
				aura.Unit.AddStatsDynamic(sim, bonusStats.Invert())
			}
		},
	}))
}

// Reduces the cooldown on your Infernal Armor ability by 30 sec and reduces the cooldown on your Demonic Grace ability by 10 sec.
func (warlock *Warlock) applyNaxxramasTank4PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneCloakInfernalArmor) && !warlock.HasRune(proto.WarlockRune_RuneLegsDemonicGrace) {
		return
	}

	label := "S03 - Item - Naxxramas - Warlock - Tank 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if warlock.InfernalArmor != nil {
				warlock.InfernalArmor.CD.FlatModifier -= time.Second * 30
			}

			if warlock.DemonicGrace != nil {
				warlock.DemonicGrace.CD.FlatModifier -= time.Second * 10
			}
		},
	})
}

// When an Undead enemy attempts to attack you, the remaining duration of your active Vengeance is reset to 20 sec.
func (warlock *Warlock) applyNaxxramasTank6PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneHelmVengeance) {
		return
	}

	label := "S03 - Item - Naxxramas - Warlock - Tank 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warlock.VengeanceAura.IsActive() && spell.Unit.MobType == proto.MobType_MobTypeUndead {
				warlock.VengeanceAura.Activate(sim)
			}
		},
	}))
}
