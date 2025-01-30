package druid

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetEclipseOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Eclipse of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Balance2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Balance4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Balance6PBonus()
		},
	},
})

// Increases the damage done and damage radius of Starfall's stars and Hurricane by 25%.
func (druid *Druid) applyT2Balance2PBonus() {
	label := "S03 - Item - T2 - Druid - Balance 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: "S03 - Item - T2 - Druid - Balance 2P Bonus",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells := core.FilterSlice(
				core.Flatten([][]*DruidSpell{
					druid.Hurricane,
					{druid.Starfall, druid.StarfallTick, druid.StarfallSplash},
				}), func(spell *DruidSpell) bool { return spell != nil },
			)

			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.25
			}
		},
	})
}

// Your Wrath casts have a 10% chance to summon a stand of 3 Treants to attack your target for until cancelled.
func (druid *Druid) applyT2Balance4PBonus() {
	label := "S03 - Item - T2 - Druid - Balance 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	affectedSpellCodes := []int32{SpellCode_DruidWrath, SpellCode_DruidStarsurge}
	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if slices.Contains(affectedSpellCodes, spell.SpellCode) && !druid.t26pcTreants.IsActive() && sim.Proc(0.10, "Summon Treants") {
				druid.t26pcTreants.EnableWithTimeout(sim, druid.t26pcTreants, time.Second*15)
			}
		},
	}))
}

// Your Wrath critical strikes have a 30% chance to make your next Starfire instant cast.
func (druid *Druid) applyT2Balance6PBonus() {
	label := "S03 - Item - T2 - Druid - Balance 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	starfires := []*DruidSpell{}
	buffAura := druid.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 467088},
		Label:     "Astral Power",
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range druid.Starfire {
				if spell != nil {
					starfires = append(starfires, spell)
				}
			}
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			for _, spell := range starfires {
				spell.DamageMultiplierAdditive -= 0.10 * float64(oldStacks)
				spell.DamageMultiplierAdditive += 0.10 * float64(newStacks)
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode == SpellCode_DruidStarfire {
				aura.Deactivate(sim)
			}
		},
	})

	procSpellCodes := []int32{SpellCode_DruidWrath, SpellCode_DruidStarsurge}
	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if slices.Contains(procSpellCodes, spell.SpellCode) && result.DidCrit() && sim.Proc(0.50, "Astral Power") {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetCunningOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Cunning of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Feral2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Feral4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Feral6PBonus()
		},
	},
})

// Increases the duration of Rake by 6 sec and its periodic damage by 50%.
func (druid *Druid) applyT2Feral2PBonus() {
	label := "S03 - Item - T2 - Druid - Feral 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467207},
		Label:    "S03 - Item - T2- Druid - Feral 2P Bonus",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, dot := range druid.Rake.Dots() {
				if dot == nil {
					continue
				}

				dot.NumberOfTicks += int32(6 / dot.TickLength.Seconds())
				dot.RecomputeAuraDuration()
				oldOnSnapshot := dot.OnSnapshot
				dot.OnSnapshot = func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					oldOnSnapshot(sim, target, dot, isRollover)
					dot.SnapshotAttackerMultiplier *= 1.50
				}
			}
		},
	})
}

// Your critical strike chance is increased by 15% while Tiger's Fury is active.
func (druid *Druid) applyT2Feral4PBonus() {
	label := "S03 - Item - T2 - Druid - Feral 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldOnGain := druid.TigersFuryAura.OnGain
			druid.TigersFuryAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
				oldOnGain(aura, sim)
				druid.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: 15 * core.CritRatingPerCritChance})
			}
			oldOnExpire := druid.TigersFuryAura.OnExpire
			druid.TigersFuryAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
				oldOnExpire(aura, sim)
				druid.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: -15 * core.CritRatingPerCritChance})
			}
		},
	})
}

// Your Shred and Mangle(Cat) abilities deal 10% increased damage per your Bleed effect on the target, up to a maximum of 20% increase.
func (druid *Druid) applyT2Feral6PBonus() {
	label := "S03 - Item - T2 - Druid - Feral 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			bleedSpells := []*DruidSpell{druid.Rake, druid.Rip}
			for _, spell := range []*DruidSpell{druid.Shred, druid.MangleCat, druid.FerociousBite} {
				if spell == nil {
					continue
				}

				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					modifier := 0.0
					for _, dotSpell := range bleedSpells {
						if dotSpell.Dot(target).IsActive() {
							modifier += 0.10
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

var ItemSetFuryOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Fury of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Guardian2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Guardian4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Guardian6PBonus()
		},
	},
})

// Swipe(Bear) also causes your Maul to hit 1 additional target for the next 6 sec.
func (druid *Druid) applyT2Guardian2PBonus() {
	label := "S03 - Item - T2 - Druid - Guardian 2P Bonus"
	if druid.Env.GetNumTargets() == 1 || druid.HasAura(label) {
		return
	}

	cleaveAura := druid.RegisterAura(core.Aura{
		Label:    "2P Cleave Buff",
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.FuryOfStormrageMaulCleave = true
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.FuryOfStormrageMaulCleave = false
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.SpellCode == SpellCode_DruidSwipeBear {
				cleaveAura.Activate(sim)
			}
		},
	}))
}

// Your Mangle(Bear), Swipe(Bear), Maul, and Lacerate abilities gain 5% increased critical strike chance against targets afflicted by your Lacerate.
func (druid *Druid) applyT2Guardian4PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsLacerate) {
		return
	}

	label := "S03 - Item - T2 - Druid - Guardian 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range []*DruidSpell{druid.MangleBear, druid.SwipeBear, druid.Maul, druid.Lacerate} {
				if spell == nil {
					continue
				}

				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					bonusCrit := 0.0
					if druid.LacerateBleed.Dot(target).GetStacks() > 0 {
						bonusCrit = 5 * core.CritRatingPerCritChance
					}

					spell.BonusCritRating += bonusCrit
					oldApplyEffects(sim, target, spell)
					spell.BonusCritRating -= bonusCrit
				}
			}
		},
	})
}

// Your Swipe now spreads your Lacerate from your primary target to other targets it strikes.
func (druid *Druid) applyT2Guardian6PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsLacerate) {
		return
	}
	label := "S03 - Item - T2 - Druid - Guardian 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.FuryOfStormrageLacerateSpread = true
	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode == SpellCode_DruidSwipeBear && result.Landed() && result.Target != druid.CurrentTarget {
				currentTargetDoT := druid.LacerateBleed.Dot(druid.CurrentTarget)
				if !currentTargetDoT.IsActive() {
					return
				}

				targetDoT := druid.LacerateBleed.Dot(result.Target)
				targetDoT.Apply(sim)
				targetDoT.SetStacks(sim, currentTargetDoT.GetStacks())
				targetDoT.UpdateExpires(sim, currentTargetDoT.ExpiresAt())
			}
		},
	}))
}

var ItemSetBountyOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Bounty of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		// Your healing spell critical strikes trigger the Dreamstate effect, granting you 50% of your mana regeneration while casting for 8 sec.
		2: func(agent core.Agent) {
		},
		// Your non-periodic spell critical strikes reduce the casting time of your next Healing Touch, Regrowth, or Nourish spell by 0.5 sec.
		4: func(agent core.Agent) {
		},
		// Increases healing from Wild Growth by 10%. In addition, Wild Growth can now be used in Moonkin Form, and its healing is increased by an additional 50% in that form.
		6: func(agent core.Agent) {
		},
	},
})

var ItemSetHaruspexsGarb = core.NewItemSet(core.ItemSet{
	Name: "Haruspex's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 12.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyZGBalance3PBonus()
		},
		5: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyZGBalance5PBonus()
		},
	},
})

// Reduces the cast time and global cooldown of Starfire by 0.5 sec.
func (druid *Druid) applyZGBalance3PBonus() {
	label := "S03 - Item - ZG - Druid - Balance 3P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range druid.Starfire {
				if spell == nil {
					continue
				}

				spell.DefaultCast.CastTime -= time.Millisecond * 500
				spell.DefaultCast.GCD -= time.Millisecond * 500
			}
		},
	})
}

// Increases the critical strike chance of Wrath by 10%.
func (druid *Druid) applyZGBalance5PBonus() {
	label := "S03 - Item - ZG - Druid - Balance 5P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range druid.Wrath {
				if spell == nil {
					continue
				}

				spell.BonusCritRating += 10 * core.SpellCritRatingPerCritChance
			}
		},
	})
}
