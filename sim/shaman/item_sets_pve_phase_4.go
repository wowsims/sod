package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetTheFiveThunders = core.NewItemSet(core.ItemSet{
	Name: "The Five Thunders",
	Bonuses: map[int32]core.ApplyEffect{
		// +40 Attack Power, up to 23 increased damage from spells, and up to 44 increased healing from spells.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.SpellDamage:       23,
				stats.HealingPower:      44,
			})
		},
		// 6% chance on mainhand autoattack and 4% chance on spellcast to increase your damage and healing done by magical spells and effects by up to 95 for 10 sec.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()

			procAura := c.NewTemporaryStatsAura("The Furious Storm", core.ActionID{SpellID: 27775}, stats.Stats{stats.SpellPower: 95}, time.Second*10)
			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "Item - The Furious Storm Proc (MH Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeMHAuto,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "Item - The Furious Storm Proc (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.04,
				Handler:    handler,
			})
		},
		// +8 All Resistances.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +200 Armor.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetEarthfuryEruption = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Eruption",
	Bonuses: map[int32]core.ApplyEffect{
		// The radius of your totems that affect friendly targets is increased to 40 yd.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT1Elemental4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT1Elemental6PBonus()
		},
	},
})

// Your Lightning Bolt critical strikes have a 35% chance to reset the cooldown on Lava Burst and Chain Lightning and make the next Lava Burst, Chain Heal, or Chain Lightning within 10 sec instant.
func (shaman *Shaman) applyT1Elemental4PBonus() {
	label := "S03 - Item - T1 - Shaman - Elemental 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode == SpellCode_ShamanLightningBolt && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() && sim.Proc(.35, "Power Surge") {
				shaman.PowerSurgeDamageAura.Activate(sim)
			}
		},
	}))
}

// Lava Burst now also refreshes the duration of Flame Shock on your target back to 12 sec.
func (shaman *Shaman) applyT1Elemental6PBonus() {
	label := "S03 - Item - T1 - Shaman - Elemental 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	var flameShockSpells []*core.Spell

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			flameShockSpells = core.FilterSlice(shaman.FlameShock, func(spell *core.Spell) bool { return spell != nil })
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode != SpellCode_ShamanLavaBurst {
				return
			}

			for _, spell := range flameShockSpells {
				if dot := spell.Dot(shaman.CurrentTarget); dot.IsActive() {
					dot.NumberOfTicks = dot.OriginalNumberOfTicks
					dot.Rollover(sim)
				}
			}
		},
	}))
}

var ItemSetEarthfuryRelief = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Relief",
	Bonuses: map[int32]core.ApplyEffect{
		// The radius of your totems that affect friendly targets is increased to 40 yd.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// After casting your Healing Wave, Lesser Healing Wave, or Riptide spell, gives you a 25% chance to gain Mana equal to 35% of the base cost of the spell.
		4: func(agent core.Agent) {
			// Not implementing for now
		},
		// Your Healing Wave will now jump to additional nearby targets. Each jump reduces the effectiveness of the heal by 80%, and the spell will jump to up to 2 additional targets.
		6: func(agent core.Agent) {
			// Not implementing for now
		},
	},
})

var ItemSetEarthfuryImpact = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Impact",
	Bonuses: map[int32]core.ApplyEffect{
		// The radius of your totems that affect friendly targets is increased to 40 yd.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT1Enhancement4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT1Enhancement6PBonus()
		},
	},
})

// Increases your critical strike chance with spells and attacks by 2%.
func (shaman *Shaman) applyT1Enhancement4PBonus() {
	label := "S03 - Item - T1 - Shaman - Enhancement 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{
		stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
		stats.SpellCrit: 2 * core.CritRatingPerCritChance,
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
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

// Your Flurry talent grants an additional 10% increase to your attack speed.
func (shaman *Shaman) applyT1Enhancement6PBonus() {
	label := "S03 - Item - T1 - Shaman - Enhancement 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.bonusFlurrySpeed += .10
		},
	})
}

var ItemSetEarthfuryResolve = core.NewItemSet(core.ItemSet{
	Name: "Earthfury Resolve",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT1Tank2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT1Tank4PBonus()
		},

		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT1Tank6PBonus()
		},
	},
})

// Increases your attack speed by 30% for your next 3 swings after you parry, dodge, or block.
func (shaman *Shaman) applyT1Tank2PBonus() {
	label := "S03 - Item - T1 - Shaman - Tank 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	flurryAura := shaman.makeFlurryAura(5)
	// The consumption trigger may not exist if the Shaman doesn't talent into Flurry
	shaman.makeFlurryConsumptionTrigger(flurryAura)

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidParry() || result.DidDodge() || result.DidBlock() {
				flurryAura.Activate(sim)
				flurryAura.SetStacks(sim, 3)
			}
		},
	}))
}

// Your parries and dodges also activate your Shield Mastery rune ability.
func (shaman *Shaman) applyT1Tank4PBonus() {
	label := "S03 - Item - T1 - Shaman - Tank 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidParry() || result.DidDodge() {
				shaman.ShieldMasteryAura.Activate(sim)
				shaman.ShieldMasteryAura.AddStack(sim)
			}
		},
	}))
}

// Your Stoneskin Totem also reduces Physical damage taken by 5% and your Windwall Totem also reduces Magical damage taken by 5%.
func (shaman *Shaman) applyT1Tank6PBonus() {
	label := "S03 - Item - T1 - Shaman - Tank 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	improvedStoneskinAura := core.ImprovedStoneskinTotemAura(&shaman.Unit)
	improvedWindwallAura := core.ImprovedWindwallTotemAura(&shaman.Unit)

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range shaman.StoneskinTotem {
				if spell == nil {
					continue
				}

				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					oldApplyEffects(sim, target, spell)
					improvedStoneskinAura.Activate(sim)
				}
			}

			for _, spell := range shaman.WindwallTotem {
				if spell == nil {
					continue
				}

				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					oldApplyEffects(sim, target, spell)
					improvedWindwallAura.Activate(sim)
				}
			}
		},
	})
}
