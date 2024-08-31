package shaman

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Totem Item IDs
const (
	TotemOfRage              = 22395
	TotemOfTheStorm          = 23199
	TotemOfSustaining        = 23200
	TotemCarvedDriftwoodIcon = 209575
	TotemInvigoratingFlame   = 215436
	TotemTormentedAncestry   = 220607
	TotemOfThunder           = 228176
	TotemOfRagingFire        = 228177
	TotemOfEarthenVitality   = 228178
	NaturalAlignmentCrystal  = 230273
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                            Totems
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=23199/totem-of-the-storm
	// Equip: Increases damage done by Chain Lightning and Lightning Bolt by up to 33.
	core.NewItemEffect(TotemOfTheStorm, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		shaman.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_ShamanLightningBolt || spell.SpellCode == SpellCode_ShamanChainLightning {
				spell.BonusDamage += 33
			}
		})
	})

	// https://www.wowhead.com/classic/item=23200/totem-of-sustaining
	// Equip: Increases healing done by Lesser Healing Wave by up to 53.
	core.NewItemEffect(TotemOfSustaining, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		shaman.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_ShamanLesserHealingWave {
				spell.BonusDamage += 53
			}
		})
	})

	core.NewItemEffect(TotemCarvedDriftwoodIcon, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.MP5, 2)
	})

	core.NewItemEffect(TotemInvigoratingFlame, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()

		shaman.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_ShamanFlameShock {
				spell.Cost.FlatModifier -= 10
			}
		})
	})

	// Ancestral Bloodstorm Beacon
	core.NewItemEffect(216615, func(agent core.Agent) {
		character := agent.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436413},
			SpellSchool: core.SpellSchoolNature | core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagAPL | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, 150, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Totem of Tormented Ancestry
	core.NewItemEffect(TotemTormentedAncestry, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Totem of Tormented Ancestry Proc", core.ActionID{SpellID: 446219}, stats.Stats{stats.AttackPower: 15, stats.SpellDamage: 15, stats.HealingPower: 15}, 12*time.Second)

		shaman.RegisterAura(core.Aura{
			Label:    "Totem of Tormented Ancestry",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellCode == SpellCode_ShamanFlameShock {
					procAura.Activate(sim)
				}
			},
		})
	})

	// Totem of Rage
	// Equip: Increases damage done by Earth Shock, Flame Shock, and Frost Shock by up to 30.
	// Acts as extra 30 spellpower for shocks.
	core.NewItemEffect(TotemOfRage, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		affectedSpellCodes := []int32{SpellCode_ShamanEarthShock, SpellCode_ShamanFlameShock, SpellCode_ShamanFrostShock}
		shaman.OnSpellRegistered(func(spell *core.Spell) {
			if slices.Contains(affectedSpellCodes, spell.SpellCode) {
				spell.BonusDamage += 30
			}
		})
	})

	// https://www.wowhead.com/classic/item=228176/totem-of-thunder
	// Equip: The cast time of your Lightning Bolt spell is reduced by -0.1 sec.
	core.NewItemEffect(TotemOfThunder, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		shaman.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_ShamanLightningBolt {
				spell.DefaultCast.CastTime -= time.Millisecond * 100
			}
		})
	})

	// https://www.wowhead.com/classic/item=228177/totem-of-raging-fire
	// Equip: Your Stormstrike spell causes you to gain 50 attack power for 12 sec. (More effective with a two - handed weapon).
	core.NewItemEffect(TotemOfRagingFire, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura1H := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: TotemOfRagingFire}.WithTag(1),
			Label:    "Totem of Raging Fire (1H)",
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, 50)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, -50)
			},
		})
		// TODO: Verify 2H value
		procAura2H := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: TotemOfRagingFire}.WithTag(2),
			Label:    "Totem of Raging Fire (2H)",
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, 200)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, -200)
			},
		})

		shaman.RegisterAura(core.Aura{
			Label:    "Totem of Raging Fire Trigger",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellCode != SpellCode_ShamanStormstrike {
					return
				}

				if shaman.MainHand().HandType == proto.HandType_HandTypeOneHand {
					procAura2H.Deactivate(sim)
					procAura1H.Activate(sim)
				} else if shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand {
					procAura1H.Deactivate(sim)
					procAura2H.Activate(sim)
				}
			},
		})
	})

	// https://www.wowhead.com/classic/item=228178/totem-of-earthen-vitality
	// Equip: While a Shield is equipped, your melee attacks with Rockbiter Weapon restore 2% of your total mana.
	core.NewItemEffect(TotemOfEarthenVitality, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()

		manaMetrics := shaman.NewManaMetrics(core.ActionID{SpellID: 461299})
		core.MakePermanent(shaman.RegisterAura(core.Aura{
			Label:    "Totem of Earthen Vitality Trigger",
			Duration: core.NeverExpires,
			OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto) || !result.Landed() || shaman.OffHand().WeaponType != proto.WeaponType_WeaponTypeShield {
					return
				}
				shaman.AddMana(sim, shaman.MaxMana()*.02, manaMetrics)
			},
		}))
	})

	///////////////////////////////////////////////////////////////////////////
	//                            Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=230273/natural-alignment-crystal
	// Use: Aligns the Shaman with nature, increasing the damage done by spells by 20%, improving heal effects by 20%, and increasing mana cost of spells by 20% for 20 sec.
	// (2 Min Cooldown)
	core.NewItemEffect(NaturalAlignmentCrystal, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()

		duration := time.Second * 20

		aura := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: NaturalAlignmentCrystal},
			Label:    "Nature Aligned",
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.DamageDealtMultiplier *= 1.20
				// shaman.PseudoStats.HealingDealtMultiplier *= 1.20
				shaman.PseudoStats.SchoolCostMultiplier.AddToAllSchools(20)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.DamageDealtMultiplier /= 1.20
				// shaman.PseudoStats.HealingDealtMultiplier /= 1.20
				shaman.PseudoStats.SchoolCostMultiplier.AddToAllSchools(-20)
			},
		})

		spell := shaman.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: NaturalAlignmentCrystal},
			ProcMask: core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    shaman.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    shaman.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		shaman.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.AddEffectsToTest = true
}
