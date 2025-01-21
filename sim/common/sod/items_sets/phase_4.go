package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

var ItemSetThePostmaster = core.NewItemSet(core.ItemSet{
	Name: "The Postmaster",
	Bonuses: map[int32]core.ApplyEffect{
		// +50 Armor.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 50)
		},
		// +10 Fire Resistance.
		// +10 Arcane Resistance.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.ArcaneResistance, 10)
			character.AddStat(stats.FireResistance, 10)
		},
		// Increases damage and healing done by magical spells and effects by up to 12.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 12)
		},
		// Increases run speed by 5%.
		// +10 Intellect.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Intellect, 10)
		},
	},
})

var ItemSetNecropileRaiment = core.NewItemSet(core.ItemSet{
	Name: "Necropile Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// +5 Stamina.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +5 Intellect.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Intellect, 5)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(15)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetIronweaveBattlesuit = core.NewItemSet(core.ItemSet{
	Name: "Ironweave Battlesuit",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases your chance to resist Silence and Interrupt effects by 10%.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// +200 Armor.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 200)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		6: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetTheHighlandersIntent = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Intent",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Spells.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellCrit, 1)
		},
	},
})

var ItemSetTheDefilersIntent = core.NewItemSet(core.ItemSet{
	Name: "The Defiler's Intent",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Spells.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellCrit, 1)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetCadaverousGarb = core.NewItemSet(core.ItemSet{
	Name: "Cadaverous Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +3.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 3)
		},
		// +10 Attack Power.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 10)
			character.AddStat(stats.RangedAttackPower, 10)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(15)
		},
		// Improves your chance to hit by 2%.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 2)
		},
	},
})

var ItemSetTheHighlandersWill = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Will",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Spells.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellCrit, 1)
		},
	},
})

var ItemSetTheHighlandersPurpose = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Purpose",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

var ItemSetTheDefilersPurpose = core.NewItemSet(core.ItemSet{
	Name: "The Defiler's Purpose",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

var ItemSetTheDefilersWill = core.NewItemSet(core.ItemSet{
	Name: "The Defiler's Will",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Spells.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellCrit, 1)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetBloodmailRegalia = core.NewItemSet(core.ItemSet{
	Name: "Bloodmail Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +3.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 3)
		},
		// +10 Attack Power.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 10)
			character.AddStat(stats.RangedAttackPower, 10)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(15)
		},
		// Increases your chance to parry an attack by 1%.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Parry, 1)
		},
	},
})

var ItemSetTheHighlandersFortitude = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Fortitude",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Spells.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellCrit, 1)
		},
	},
})

var ItemSetTheHighlandersDetermination = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Determination",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

var ItemSetTheDefilersDetermination = core.NewItemSet(core.ItemSet{
	Name: "The Defiler's Determination",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

var ItemSetTheDefilersFortitude = core.NewItemSet(core.ItemSet{
	Name: "The Defiler's Fortitude",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

var ItemSetDeathboneGuardian = core.NewItemSet(core.ItemSet{
	Name: "Deathbone Guardian",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +3.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 3)
		},
		// +50 Armor.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 50)
		},
		// +15 All Resistances.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(15)
		},
		// Increases run speed by 5%.
		// +10 Intellect.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Parry, 1)
		},
	},
})

var ItemSetTheHighlandersResolve = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Resolve",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

var ItemSetTheHighlandersResolution = core.NewItemSet(core.ItemSet{
	Name: "The Highlander's Resolution",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

var ItemSetTheDefilersResolution = core.NewItemSet(core.ItemSet{
	Name: "The Defiler's Resolution",
	Bonuses: map[int32]core.ApplyEffect{
		// Increase Stamina +5
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Stamina, 5)
		},
		// +1 Crit with Melee
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Other
///////////////////////////////////////////////////////////////////////////

var ItemSetSpidersKiss = core.NewItemSet(core.ItemSet{
	Name: "Spider's Kiss",
	Bonuses: map[int32]core.ApplyEffect{
		// Chance on Hit: Immobilizes the target and lowers their armor by 100 for 10 sec.
		// Increased Defense +7.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Spider's Kiss", core.ActionID{SpellID: 17333}, stats.Stats{stats.Armor: -100}, time.Second*10)
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:          core.ActionID{SpellID: 446570},
				Name:              "Echoes of the Depraved",
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ProcChance:        0.05,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
			character.AddStat(stats.Defense, 7)
		},
	},
})

var ItemSetDalRendsArms = core.NewItemSet(core.ItemSet{
	Name: "Dal'Rend's Arms",
	Bonuses: map[int32]core.ApplyEffect{
		// +50 Attack Power.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 50)
			character.AddStat(stats.RangedAttackPower, 50)
		},
	},
})

var ItemSetShardOfTheGods = core.NewItemSet(core.ItemSet{
	Name: "Shard of the Gods",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases healing done by spells and effects by up to 55.
		// Increases damage done by magical spells and effects by up to 29.
		// Your spell casts have a chance to summon Servants of the Scale or Flame.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.HealingPower, 55)
			character.AddStat(stats.SpellDamage, 29)
			// TODO: pets [Your spell casts have a chance to summon Servants of the Scale or Flame.]
		},
	},
})

var ItemSetSpiritOfEskhandar = core.NewItemSet(core.ItemSet{
	Name: "Spirit of Eskhandar",
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to hit with all spells and attacks by 1%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 1)
			character.AddStat(stats.SpellHit, 1)
		},
		// Improves your chance to get a critical strike with all spells and attacks by 1%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			character.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		},
		// 1% chance on a melee hit to call forth the spirit of Eskhandar to protect you in battle for 2 min.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:              "Call of Eskhandar Trigger",
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ProcChance:        1,
				ICD:               time.Minute * 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					for _, petAgent := range character.PetAgents {
						if eskhandar, ok := petAgent.(*guardians.Eskhandar); ok {
							eskhandar.EnableWithTimeout(sim, eskhandar, time.Minute*2)
							break
						}
					}
				},
			})
		},
	},
})

var ItemSetCoreHoundsCall = core.NewItemSet(core.ItemSet{
	Name: "Core Hound's Call",
	Bonuses: map[int32]core.ApplyEffect{
		// Small chance on melee hit to call forth a Core Hound for 1 min.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:              "Core Hound's Call Trigger",
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ProcChance:        1,
				ICD:               time.Minute * 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					for _, petAgent := range character.PetAgents {
						if coreHound, ok := petAgent.(*guardians.CoreHound); ok {
							coreHound.EnableWithTimeout(sim, coreHound, time.Minute*1)
							break
						}
					}
				},
			})
		},
		// Small chance on melee hit to call forth the Spirit of Magmadar to assist you in battle. Increasing your attack speed by 10% for 20 sec.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 461270},
				Label:    "Magmadar's Return",
				Duration: time.Second * 20,
			}).AttachMultiplyAttackSpeed(&character.Unit, 1.1)

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:              "Magmadar's Return Trigger",
				Callback:          core.CallbackOnSpellHitDealt,
				Outcome:           core.OutcomeLanded,
				ProcMask:          core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ICD:               time.Minute * 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
	},
})
