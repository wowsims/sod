package item_sets

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Keep these in alphabetical order.

var ItemSetBlackDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Black Dragon Mail",
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to hit by 1%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 1)
		},
		// Improves your chance to get a critical strike by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 2*core.CritRatingPerCritChance)
		},
		// +10 Fire Resistance.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.FireResistance, 10)
		},
	},
})

var ItemSetBlueDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Blue Dragon Mail",
	Bonuses: map[int32]core.ApplyEffect{
		// +4 All Resistances.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStats(stats.Stats{
				stats.ArcaneResistance: 4,
				stats.FireResistance:   4,
				stats.FrostResistance:  4,
				stats.NatureResistance: 4,
				stats.ShadowResistance: 4,
			})
		},
		// Increases damage and healing done by magical spells and effects by up to 28.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 28)
		},
	},
})

var ItemSetBloodsoulEmbrace = core.NewItemSet(core.ItemSet{
	Name: "Bloodsoul Embrace",
	Bonuses: map[int32]core.ApplyEffect{
		// Restores 12 mana per 5 sec.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MP5, 12)
		},
	},
})

// var ItemSetBloodvineGarb = core.NewItemSet(core.ItemSet{
// 	Name: "Bloodvine Garb",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Improves your chance to get a critical strike with spells by 2%.
// 		2: func(agent core.Agent) {
// 			character := agent.GetCharacter()
// 			character.AddStat(stats.SpellCrit, 2*core.SpellCritRatingPerCritChance)
// 		},
// 	},
// })

// var ItemSetBloodTigerHarness = core.NewItemSet(core.ItemSet{
// 	Name: "Blood Tiger Harness",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Improves your chance to get a critical strike by 1%.
// 		// Improves your chance to get a critical strike with spells by 1%.
// 		2: func(agent core.Agent) {
// 			character := agent.GetCharacter()
// 			character.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
// 			character.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
// 		},
// 	},
// })

var ItemSetDevilsaurArmor = core.NewItemSet(core.ItemSet{
	Name: "Devilsaur Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to hit by 2%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 2)
		},
	},
})

var ItemSetGreenDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Green Dragon Mail",
	Bonuses: map[int32]core.ApplyEffect{
		// Restores 3 mana per 5 sec.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MP5, 3)
		},
		// Allows 15% of your Mana regeneration to continue while casting.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SpiritRegenRateCasting += .15
		},
	},
})

var ItemSetImperialPlate = core.NewItemSet(core.ItemSet{
	Name: "Imperial Plate",
	Bonuses: map[int32]core.ApplyEffect{
		// +100 Armor.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 100)
		},
		// +28 Attack Power.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 28)
			character.AddStat(stats.RangedAttackPower, 28)
		},
		// +18 Stamina.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 18)
		},
	},
})

var ItemSetIronfeatherArmor = core.NewItemSet(core.ItemSet{
	Name: "Ironfeather Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 20.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 20)
		},
	},
})

var ItemSetStormshroudArmor = core.NewItemSet(core.ItemSet{
	Name: "Stormshroud Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// 5% chance of dealing 15 to 25 Nature damage on a successful melee attack.
		2: func(a core.Agent) {
			char := a.GetCharacter()
			proc := char.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 18980},
				SpellSchool: core.SpellSchoolNature,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(15, 25), spell.OutcomeMagicHitAndCrit)
				},
			})
			char.RegisterAura(core.Aura{
				Label:    "Stormshround Armor 2pc",
				ActionID: core.ActionID{SpellID: 18979},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("Stormshroud Armor 2pc") < 0.05 {
						proc.Cast(sim, result.Target)
					}
				},
			})
		},
		// 2% chance on melee attack of restoring 30 energy.
		3: func(a core.Agent) {
			char := a.GetCharacter()
			if !char.HasEnergyBar() {
				return
			}
			metrics := char.NewEnergyMetrics(core.ActionID{SpellID: 23863})
			proc := char.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 23864},
				SpellSchool: core.SpellSchoolNature,
				ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
					char.AddEnergy(sim, 30, metrics)
				},
			})
			char.RegisterAura(core.Aura{
				Label:    "Stormshround Armor 3pc",
				ActionID: core.ActionID{SpellID: 18979},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("Stormshroud Armor 3pc") < 0.02 {
						proc.Cast(sim, result.Target)
					}
				},
			})

		},
		// +14 Attack Power.
		4: func(a core.Agent) {
			a.GetCharacter().AddStat(stats.AttackPower, 14)
			a.GetCharacter().AddStat(stats.RangedAttackPower, 14)
		},
	},
})

var ItemSetTheDarksoul = core.NewItemSet(core.ItemSet{
	Name: "The Darksoul",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +20.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 20)
		},
	},
})

var ItemSetVolcanicArmor = core.NewItemSet(core.ItemSet{
	Name: "Volcanic Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// 5% chance of dealing 15 to 25 Fire damage on a successful melee attack.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			procSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 9057},
				SpellSchool: core.SpellSchoolFire,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(15, 25), spell.OutcomeMagicHitAndCrit)
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Firebolt Trigger (Volcanic Armor)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: .05,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					procSpell.Cast(sim, result.Target)
				},
			})
		},
	},
})
