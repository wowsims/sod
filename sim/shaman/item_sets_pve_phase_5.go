package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetEruptionOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Eruption of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Elemental2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Elemental4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Elemental6PBonus()
		},
	},
})

// Your spell critical strikes now have a 100% chance trigger your Elemental Focus talent.
func (shaman *Shaman) applyT2Elemental2PBonus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	label := "S03 - Item - T2 - Shaman - Elemental 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if shaman.isShamanDamagingSpell(spell) && result.DidCrit() {
				shaman.ClearcastingAura.Activate(sim)
				shaman.ClearcastingAura.SetStacks(sim, shaman.ClearcastingAura.MaxStacks)
			}
		},
	}))
}

// Loyal Beta from your Spirit of the Alpha ability now also increases Fire, Frost, and Nature damage by 5%.
func (shaman *Shaman) applyT2Elemental4PBonus() {
	label := "S03 - Item - T2 - Shaman - Elemental 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if shaman.LoyalBetaAura == nil {
				return
			}

			shaman.LoyalBetaAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1.05
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.05
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.05
			})
		},
	}))
}

// While Clearcasting is active, you deal 15% more non-Physical damage.
func (shaman *Shaman) applyT2Elemental6PBonus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	label := "S03 - Item - T2 - Shaman - Elemental 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ClearcastingAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1.15)
			})

			shaman.ClearcastingAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1 / 1.15)
			})
		},
	}))
}

var ItemSetResolveOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Resolve of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Tank2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Tank4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Tank6PBonus()
		},
	},
})

// Your Flame Shock also grants 30% increased chance to Block for 5 sec or until you Block an attack.
func (shaman *Shaman) applyT2Tank2PBonus() {
	label := "S03 - Item - T2 - Shaman - Tank 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shieldBlockAura := shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467891},
		Label:    "Shield Block",
		Duration: time.Second * 5,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			if result.DidBlock() {
				aura.Deactivate(sim)
			}
		},
	}).AttachStatBuff(stats.Block, 30*core.BlockRatingPerBlockChance)

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           label,
		ClassSpellMask: ClassSpellMask_ShamanFlameShock,
		Callback:       core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shieldBlockAura.Activate(sim)
		},
	})
}

// Each time you Block, your Block amount is increased by 10% of your Spell Damage for 6 sec, stacking up to 3 times.
func (shaman *Shaman) applyT2Tank4PBonus() {
	label := "S03 - Item - T2 - Shaman - Tank 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	spellDamStatDeps := []*stats.StatDependency{
		shaman.NewDynamicStatDependency(stats.SpellDamage, stats.BlockValue, 0),
		shaman.NewDynamicStatDependency(stats.SpellDamage, stats.BlockValue, 0.10),
		shaman.NewDynamicStatDependency(stats.SpellDamage, stats.BlockValue, 0.20),
		shaman.NewDynamicStatDependency(stats.SpellDamage, stats.BlockValue, 0.30),
	}
	spellPowerStatDeps := []*stats.StatDependency{
		shaman.NewDynamicStatDependency(stats.SpellPower, stats.BlockValue, 0),
		shaman.NewDynamicStatDependency(stats.SpellPower, stats.BlockValue, 0.10),
		shaman.NewDynamicStatDependency(stats.SpellPower, stats.BlockValue, 0.20),
		shaman.NewDynamicStatDependency(stats.SpellPower, stats.BlockValue, 0.30),
	}

	// Couldn't find a separate spell for this
	blockAura := shaman.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 467910},
		Label:     "Elemental Shield",
		Duration:  time.Second * 6,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			shaman.DisableDynamicStatDep(sim, spellDamStatDeps[oldStacks])
			shaman.EnableDynamicStatDep(sim, spellDamStatDeps[newStacks])

			shaman.DisableDynamicStatDep(sim, spellPowerStatDeps[oldStacks])
			shaman.EnableDynamicStatDep(sim, spellPowerStatDeps[newStacks])
		},
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidBlock() {
				blockAura.Activate(sim)
				blockAura.AddStack(sim)
			}
		},
	}))
}

// Each time you Block an attack, you have a 50% chance to trigger your Maelstrom Weapon rune.
func (shaman *Shaman) applyT2Tank6PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	label := "S03 - Item - T2 - Shaman - Tank 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidBlock() && sim.Proc(0.50, "T2 6P Proc Maelstrom Weapon") {
				shaman.MaelstromWeaponAura.Activate(sim)
				shaman.MaelstromWeaponAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetImpactOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Impact of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Enhancement2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Enhancement4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Enhancement6PBonus()
		},
	},
})

// Your chance to trigger Static Shock is increased by 12% (6% while dual-wielding)
func (shaman *Shaman) applyT2Enhancement2PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock) {
		return
	}

	label := "S03 - Item - T2 - Shaman - Enhancement 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.staticSHocksProcChance += 0.06
		},
	}))
}

// Main-hand Stormstrike now deals 50% more damage.
func (shaman *Shaman) applyT2Enhancement4PBonus() {
	if !shaman.Talents.Stormstrike {
		return
	}

	label := "S03 - Item - T2 - Shaman - Enhancement 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_ShamanStormstrikeHit,
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMeleeMHSpecial,
		FloatValue: 1.5,
	}))
}

// While Static Shock is engraved, your Lightning Shield now gains a charge each time you hit a target with Lightning Bolt or Chain Lightning, up to a maximum of 9 charges.
// In addition, while Static Shock is engraved, your Lightning Shield can now deal critical damage.
func (shaman *Shaman) applyT2Enhancement6PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock) {
		return
	}

	label := "S03 - Item - T2 - Shaman - Enhancement 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	affectedSpellClassMasks := ClassSpellMask_ShamanLightningBolt | ClassSpellMask_ShamanChainLightning
	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(t26pAura *core.Aura, sim *core.Simulation) {
			for _, aura := range shaman.LightningShieldAuras {
				if aura == nil {
					continue
				}

				oldOnGain := aura.OnGain
				aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
					oldOnGain(aura, sim)
					t26pAura.Activate(sim)
				}

				oldOnExpire := aura.OnExpire
				aura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
					oldOnExpire(aura, sim)
					t26pAura.Deactivate(sim)
				}
			}

			shaman.lightningShieldCanCrit = true
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Tested and it doesn't proc from overloads
			if spell.Matches(affectedSpellClassMasks) && !spell.ProcMask.Matches(core.ProcMaskSpellProc) && result.Landed() {
				shaman.ActiveShieldAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetReliefOfTheTenStorms = core.NewItemSet(core.ItemSet{
	Name: "Relief of the Ten Storms",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Restoration2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Restoration4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyT2Restoration6PBonus()
		},
	},
})

// Your damaging and healing critical strikes now have a 100% chance to trigger your Water Shield, but do not consume a charge or trigger its cooldown.
func (shaman *Shaman) applyT2Restoration2PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsWaterShield) {
		return
	}

	label := "S03 - Item - T2 - Shaman - Restoration 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.DidCrit() {
				shaman.WaterShieldRestore.Cast(sim, aura.Unit)
			}
		},
	}))
}

// Your Chain Lightning now also heals the target of your Earth Shield for 100% of the damage done.
func (shaman *Shaman) applyT2Restoration4PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsWaterShield) {
		return
	}

	label := "S03 - Item - T2 - Shaman - Restoration 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	healthMetrics := shaman.NewHealthMetrics(core.ActionID{SpellID: 467809})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_ShamanChainLightning) {
				shaman.GainHealth(sim, result.Damage, healthMetrics)
			}
		},
	}))
}

// Increases the healing of Chain Heal and the damage of Chain Lightning by 20%.
func (shaman *Shaman) applyT2Restoration6PBonus() {
	label := "S03 - Item - T2 - Shaman - Restoration 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanChainHeal | ClassSpellMask_ShamanChainLightning,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  20,
	}))
}

var ItemSetAugursRegalia = core.NewItemSet(core.ItemSet{
	Name: "Augur's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +7.
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			// No corresponding Soul found so leaving this simple
			shaman.AddStat(stats.Defense, 7)
		},

		3: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyZGTank3PBonus()
		},

		5: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyZGTank5PBonus()
		},
	},
})

// Increases your chance to block attacks with a shield by 10%.
func (shaman *Shaman) applyZGTank3PBonus() {
	label := "S03 - Item - ZG - Shaman - Tank 3P Bonus"
	if shaman.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Block: 10 * core.BlockRatingPerBlockChance}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachStatsBuff(bonusStats))
}

// Increases the chance to trigger your Power Surge rune by an additional 5%.
func (shaman *Shaman) applyZGTank5PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) {
		return
	}

	label := "S03 - Item - ZG - Shaman - Tank 5P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.powerSurgeProcChance += .05
		},
	})
}
