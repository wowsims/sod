package shaman

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	// Keep these ordered by ID
	TotemOfRage               = 22395
	TotemOfTheStorm           = 23199
	TotemOfSustaining         = 23200
	TotemCarvedDriftwoodIcon  = 209575
	TotemInvigoratingFlame    = 215436
	TotemTormentedAncestry    = 220607
	TerrestrisTank            = 224279
	TotemOfThunder            = 228176
	TotemOfRagingFire         = 228177
	TotemOfEarthenVitality    = 228178
	NaturalAlignmentCrystal   = 230273
	WushoolaysCharmOfSpirits  = 231281
	TerrestrisEle             = 231890
	TotemOfRelentlessThunder  = 232392
	TotemOfTheElements        = 232409
	TotemOfAstralFlow         = 232416
	TotemOfConductiveCurrents = 232419
)

func init() {
	core.AddEffectsToTest = false

	// Keep these ordered by name

	// https://www.wowhead.com/classic/item=231281/wushoolays-charm-of-spirits
	// Use: Increases the damage dealt by your Lightning Shield spell by 100% for 20 sec. (2 Min Cooldown)
	core.NewItemEffect(WushoolaysCharmOfSpirits, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()

		duration := time.Second * 20
		actionID := core.ActionID{ItemID: WushoolaysCharmOfSpirits}

		aura := shaman.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Wushoolay's Charm of Spirits",
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				for _, spell := range shaman.LightningShieldProcs {
					if spell != nil {
						spell.DamageMultiplier *= 2
					}
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				for _, spell := range shaman.LightningShieldProcs {
					if spell != nil {
						spell.DamageMultiplier /= 2
					}
				}
			},
		})

		spell := shaman.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
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
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

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
				shaman.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1.20)
				// shaman.PseudoStats.HealingDealtMultiplier *= 1.20
				shaman.PseudoStats.SchoolCostMultiplier.AddToAllSchools(20)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1 / 1.20)
				// shaman.PseudoStats.HealingDealtMultiplier /= 1.20
				shaman.PseudoStats.SchoolCostMultiplier.AddToAllSchools(-20)
			},
		})

		spell := shaman.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: NaturalAlignmentCrystal},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
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

	// https://www.wowhead.com/classic/item=231890/terrestris
	core.NewItemEffect(TerrestrisEle, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()

		boonOfEarth := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469208},
			Label:    "Boon of Earth",
			Duration: time.Minute * 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.MeleeCrit, 1*core.CritRatingPerCritChance)
				shaman.AddStatDynamic(sim, stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.MeleeCrit, -1*core.CritRatingPerCritChance)
				shaman.AddStatDynamic(sim, stats.SpellCrit, -1*core.SpellCritRatingPerCritChance)
			},
		})

		boonOfFire := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469209},
			Label:    "Boon of Fire",
			Duration: time.Minute * 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.SpellDamage, 16)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.SpellDamage, -16)
			},
		})

		boonOfWater := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469210},
			Label:    "Boon of Water",
			Duration: time.Minute * 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.HealingPower, 31)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.HealingPower, 31)
			},
		})

		boonOfAir := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 452456},
			Label:    "Boon of Air",
			Duration: time.Minute * 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.MoveSpeed *= 1.15
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.MoveSpeed /= 1.15
			},
		})

		core.MakePermanent(shaman.RegisterAura(core.Aura{
			Label: "Terrestris Boon Trigger",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				core.Each(shaman.EarthTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfEarth.Activate(sim)
					}
				})
				core.Each(shaman.FireTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfFire.Activate(sim)
					}
				})
				core.Each(shaman.WaterTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfWater.Activate(sim)
					}
				})
				core.Each(shaman.AirTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfAir.Activate(sim)
					}
				})
			},
		}))
	})

	// https://www.wowhead.com/classic/item=224279/terrestris
	core.NewItemEffect(TerrestrisTank, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()

		boonOfEarth := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 452464},
			Label:    "Boon of Earth",
			Duration: time.Minute * 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.Block, 3*core.BlockRatingPerBlockChance)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.Block, 31*core.BlockRatingPerBlockChance)
			},
		})

		fireExplosion := shaman.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 453085},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, 8, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		boonOfFire := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 452460},
			Label:    "Boon of Fire",
			Duration: time.Minute * 2,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.DidBlock() {
					fireExplosion.Cast(sim, result.Target)
				}
			},
		})

		waterHealActionID := core.ActionID{SpellID: 453081}
		waterHealthMetrics := shaman.NewHealthMetrics(waterHealActionID)
		waterHeal := shaman.RegisterSpell(core.SpellConfig{
			ActionID:    waterHealActionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				shaman.GainHealth(sim, 20, waterHealthMetrics)
			},
		})

		boonOfWater := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 452454},
			Label:    "Boon of Water",
			Duration: time.Minute * 2,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.DidBlock() {
					waterHeal.Cast(sim, aura.Unit)
				}
			},
		})

		boonOfAir := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 452456},
			Label:    "Boon of Air",
			Duration: time.Minute * 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.MoveSpeed *= 1.15
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.MoveSpeed /= 1.15
			},
		})

		core.MakePermanent(shaman.RegisterAura(core.Aura{
			Label: "Terrestris Boon Trigger",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				core.Each(shaman.EarthTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfEarth.Activate(sim)
					}
				})
				core.Each(shaman.FireTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfFire.Activate(sim)
					}
				})
				core.Each(shaman.WaterTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfWater.Activate(sim)
					}
				})
				core.Each(shaman.AirTotems, func(spell *core.Spell) {
					oldApplyEffects := spell.ApplyEffects
					spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						oldApplyEffects(sim, target, spell)
						boonOfAir.Activate(sim)
					}
				})
			},
		}))
	})

	core.NewItemEffect(TotemCarvedDriftwoodIcon, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.MP5, 2)
	})

	// https://www.wowhead.com/classic/item=232416/totem-of-astral-flow
	// Increases the attack power bonus on Windfury Weapon attacks by 68.
	core.NewItemEffect(TotemOfAstralFlow, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		shaman.bonusWindfuryWeaponAP += 68
	})

	// https://www.wowhead.com/classic/item=232419/totem-of-conductive-currents
	// While Frostbrand Weapon is active, your Water Shield triggers reduce the cast time of your next Chain Lightning spell within 15 sec by 20%, and increases its damage by 20%.
	// Stacking up to 5 times.
	core.NewItemEffect(TotemOfConductiveCurrents, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		if shaman.getImbueProcMask(proto.WeaponImbue_FrostbrandWeapon) == core.ProcMaskUnknown || !shaman.HasRune(proto.ShamanRune_RuneHandsWaterShield) {
			return
		}

		affectedSpells := []*core.Spell{}

		buffAura := shaman.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 470272},
			Label:     "Totem of Conductive Currents",
			Duration:  time.Second * 15,
			MaxStacks: 5,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				affectedSpells = core.FilterSlice(shaman.ChainLightning, func(spell *core.Spell) bool { return spell != nil })
				affectedSpells = core.FilterSlice(shaman.ChainLightningOverload, func(spell *core.Spell) bool { return spell != nil })
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				oldStackValue := .20 * float64(oldStacks)
				newStackValue := .20 * float64(newStacks)

				for _, spell := range affectedSpells {
					spell.DamageMultiplier /= 1 + oldStackValue
					spell.DamageMultiplier *= 1 + newStackValue

					spell.CastTimeMultiplier += oldStackValue
					spell.CastTimeMultiplier -= newStackValue
				}
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellCode == SpellCode_ShamanChainLightning {
					aura.Deactivate(sim)
				}
			},
		})

		core.MakePermanent(shaman.RegisterAura(core.Aura{
			Label: "Totem of Conductive Currents Trigger",
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == shaman.WaterShieldRestore {
					buffAura.Activate(sim)
					buffAura.AddStack(sim)
				}
			},
		}))
	})

	// https://www.wowhead.com/classic/item=232392/totem-of-relentless-thunder
	// While a Shield is equipped, your melee attacks with Rockbiter Weapon trigger your Maelstrom Weapon rune 100% more often.
	core.NewItemEffect(TotemOfRelentlessThunder, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
			return
		}

		core.MakePermanent(shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 470081},
			Label:    "Totem of Raging Storms",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				oldPPMM := shaman.maelstromWeaponPPMM
				newPPMM := shaman.AutoAttacks.NewPPMManager(oldPPMM.GetPPM()*2, core.ProcMaskMelee)
				shaman.maelstromWeaponPPMM = &newPPMM

				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:          time.Second * 2,
					TickImmediately: true,
					OnAction: func(sim *core.Simulation) {
						if shaman.OffHand().WeaponType != proto.WeaponType_WeaponTypeShield {
							shaman.maelstromWeaponPPMM = oldPPMM
							aura.Deactivate(sim)
							return
						}

						if !aura.IsActive() {
							shaman.maelstromWeaponPPMM = &newPPMM
							aura.Activate(sim)
						}
					},
				})
			},
		}))
	})

	core.NewItemEffect(TotemOfTheElements, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		shaman.RegisterAura(core.Aura{
			Label: "Totem of the Elements",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				shaman.ClearcastingAura.MaxStacks = 2
			},
		})
	})

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
		procAura := shaman.NewTemporaryStatsAura(
			"Totem of Tormented Ancestry Proc",
			core.ActionID{SpellID: 446219},
			stats.Stats{stats.AttackPower: 10, stats.SpellDamage: 10, stats.HealingPower: 10},
			12*time.Second,
		)

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
	// Equip: Your Stormstrike spell causes you to gain 24 attack power for 12 sec. (More effective with a two - handed weapon).
	core.NewItemEffect(TotemOfRagingFire, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura1H := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: TotemOfRagingFire}.WithTag(1),
			Label:    "Totem of Raging Fire (1H)",
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, 24)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, -24)
			},
		})
		// TODO: Verify 2H value
		procAura2H := shaman.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: TotemOfRagingFire}.WithTag(2),
			Label:    "Totem of Raging Fire (2H)",
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, 48)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.AttackPower, -48)
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

	core.AddEffectsToTest = true
}
