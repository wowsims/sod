package druid

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/common/sod"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Totem Item IDs
const (
	WolfsheadHelm                    = 8345
	IdolOfFerocity                   = 22397
	IdolOfTheMoon                    = 23197
	IdolOfBrutality                  = 23198
	IdolOfCruelty                    = 232424
	IdolMindExpandingMushroom        = 209576
	Catnip                           = 213407
	IdolOfWrath                      = 216490
	BloodBarkCrusher                 = 216499
	IdolOfTheDream                   = 220606
	RitualistsHammer                 = 221446
	Glasir                           = 224281
	Raelar                           = 224282
	IdolOfExsanguinationCat          = 228181
	IdolOfTheSwarm                   = 228180
	IdolOfExsanguinationBear         = 228182
	BloodGuardDragonhideGrips        = 227180
	KnightLieutenantsDragonhideGrips = 227183
	WushoolaysCharmOfNature          = 231280
	PristineEnchantedSouthSeasKelp   = 231316
	IdolOfCelestialFocus             = 232390
	IdolOfFelineFocus                = 232391
	IdolOfUrsinPower                 = 234468
	IdolOfFelineFerocity             = 234469
	IdolOfSiderealWrath              = 234474
	AtieshDruid                      = 236401
	ScarletRotbringer                = 240842
	StaffOfTheGlade                  = 240849
)

func init() {
	core.AddEffectsToTest = false

	// https://www.wowhead.com/classic/item=236401/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshDruid, func(agent core.Agent) {
		character := agent.GetCharacter()
		aura := core.AtieshSpellCritEffect(&character.Unit)
		character.ItemSwap.RegisterProc(AtieshDruid, aura)
	})

	core.NewItemEffect(BloodBarkCrusher, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.newBloodbarkCleaveItem(BloodBarkCrusher)
	})

	// https://www.wowhead.com/classic/item=227180/blood-guards-dragonhide-grips
	// Equip: Reduces the mana cost of your shapeshifts by 150.
	core.NewItemEffect(BloodGuardDragonhideGrips, func(agent core.Agent) {
		registerDragonHideGripsAura(agent.(DruidAgent).GetDruid())
	})

	// https://www.wowhead.com/classic/item=224281/glasir
	// Equip: Critical effects from heals have a chance to heal 3 nearby allies for 200 to 350, and critical spell hits have a chance to damage 3 nearby enemies for 100 to 175 nature damage.
	// (Proc chance: 15%, 15s cooldown)
	core.NewItemEffect(Glasir, func(agent core.Agent) {
		character := agent.GetCharacter()

		numDamageHits := min(3, character.Env.GetNumTargets())
		numHealHits := min(3, len(character.Env.Raid.AllPlayerUnits))
		damageResults := make([]*core.SpellResult, numDamageHits)
		healResults := make([]*core.SpellResult, numHealHits)

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 452424},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.10,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for idx := range damageResults {
					damageResults[idx] = spell.CalcDamage(sim, target, sim.Roll(100, 175), spell.OutcomeMagicHitAndCrit)
					target = sim.Environment.NextTargetUnit(target)
				}

				for _, result := range damageResults {
					spell.DealDamage(sim, result)
				}
			},
		})

		healSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 453009},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.10,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for idx := range healResults {
					healResults[idx] = spell.CalcHealing(sim, target, sim.Roll(200, 350), spell.OutcomeMagicCrit)
					target = sim.Environment.NextTargetUnit(target)
				}

				for _, result := range healResults {
					spell.DealHealing(sim, result)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Gla'sir Damage",
			Callback:   core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeCrit,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.15,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				damageSpell.Cast(sim, result.Target)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Gla'sir Heal",
			Callback:   core.CallbackOnPeriodicHealDealt | core.CallbackOnHealDealt,
			Outcome:    core.OutcomeCrit,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.15,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				healSpell.Cast(sim, result.Target)
			},
		})
	})

	// https://www.wowhead.com/classic/item=23198/idol-of-brutality
	// Equip: Reduces the rage cost of Maul and Swipe by 3.
	core.NewItemEffect(IdolOfBrutality, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Reduced Maul and Swipe Cost",
		}).AttachSpellMod(core.SpellModConfig{
			// For whatever reason also affects Mnagle (Bear)
			ClassMask: ClassSpellMask_DruidMaul | ClassSpellMask_DruidSwipeBear | ClassSpellMask_DruidSwipeCat | ClassSpellMask_DruidMangleBear,
			Kind:      core.SpellMod_PowerCost_Flat,
			IntValue:  -3,
		}))
	})

	// https://www.wowhead.com/classic/item=232390/idol-of-celestial-focus
	// Equip: Increases the damage done by Starfall by 10%, but decreases its radius by 50%.
	core.NewItemEffect(IdolOfCelestialFocus, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of Celestial Focus",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_DruidStarfallTick | ClassSpellMask_DruidStarfallSplash,
			IntValue:  10,
		}))
	})

	// https://www.wowhead.com/classic/item=232391/idol-of-feline-focus
	// Equip: Your Ferocious Bite ability no longer converts additional energy into damage, and refunds 30 energy on a Dodge, Miss, or Parry.
	core.NewItemEffect(IdolOfFelineFocus, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.FerociousBiteExcessEnergyOverride = true

		actionID := core.ActionID{SpellID: 470270}

		energyMetrics := druid.NewEnergyMetrics(actionID)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Idol of Feline Focus",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Matches(ClassSpellMask_DruidFerociousBite) && result.Outcome.Matches(core.OutcomeDodge|core.OutcomeMiss|core.OutcomeParry) {
					druid.AddEnergy(sim, 30, energyMetrics)
				}
			},
		}))
	})

	// https://www.wowhead.com/classic/item=22397/idol-of-ferocity
	// Equip: Reduces the energy cost of Claw and Rake by 3.
	core.NewItemEffect(IdolOfFerocity, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		// TODO: Claw is not implemented
		druid.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidRake | ClassSpellMask_DruidMangleCat) {
				spell.Cost.FlatModifier -= 3
			}
		})
	})

	// https://www.wowhead.com/classic/item=23197/idol-of-the-moon
	// Equip: Increases the damage of your Moonfire spell by up to 33.
	core.NewItemEffect(IdolOfTheMoon, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidMoonfire | ClassSpellMask_DruidSunfire | ClassSpellMask_DruidSunfireCat | ClassSpellMask_DruidStarfallSplash | ClassSpellMask_DruidStarfallTick) {
				spell.BonusDamage += 33
			}
		})
	})

	core.NewItemEffect(IdolMindExpandingMushroom, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.Spirit, 5)
	})

	// https://www.wowhead.com/classic/item=228181/idol-of-exsanguination-cat
	// Equip: The energy cost of your Rake and Rip spells is reduced by 5.
	core.NewItemEffect(IdolOfExsanguinationCat, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidRake | ClassSpellMask_DruidRip) {
				spell.Cost.FlatModifier -= 5
			}
		})
	})

	// https://www.wowhead.com/classic/item=228182/idol-of-exsanguination-bear
	// Equip: Your Lacerate ticks energize you for 3 rage.
	core.NewItemEffect(IdolOfExsanguinationBear, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		rageMetrics := druid.NewRageMetrics(core.ActionID{ItemID: IdolOfExsanguinationBear})

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of Exsanguination (Bear)",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Matches(ClassSpellMask_DruidLacerate) {
					druid.AddRage(sim, 3, rageMetrics)
				}
			},
		}))
	})

	// https://www.wowhead.com/classic/item=234469/idol-of-feline-ferocity
	// Increases the damage of Ferocious Bite and Shred by 3%.
	core.NewItemEffect(IdolOfFelineFerocity, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Improved Shred/Ferocious Bite",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_DruidFerociousBite | ClassSpellMask_DruidShred,
			IntValue:  3,
		}))
	})

	// https://www.wowhead.com/classic/item=234474/idol-of-sidereal-wrath
	// Increases the damage of Moonfire and Wrath by 3%.
	core.NewItemEffect(IdolOfSiderealWrath, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Improved Wrath/Moonfire",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_DruidWrath | ClassSpellMask_DruidMoonfire | ClassSpellMask_DruidSunfire | ClassSpellMask_DruidStarsurge | ClassSpellMask_DruidStarfallSplash | ClassSpellMask_DruidStarfallTick,
			IntValue:  3,
		}))
	})

	// https://www.wowhead.com/classic/item=220606/idol-of-the-dream
	// Equip: Increases the damage of Swipe and Shred by 2%.
	core.NewItemEffect(IdolOfTheDream, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Improved Swipe/Shred",
		})).AttachSpellMod(core.SpellModConfig{
			ClassMask: ClassSpellMask_DruidSwipeBear | ClassSpellMask_DruidSwipeCat | ClassSpellMask_DruidShred,
			Kind:      core.SpellMod_BaseDamageDone_Flat,
			IntValue:  2,
		})
	})

	// https://www.wowhead.com/classic/item=228180/idol-of-the-swarm
	// Equip: The duration of your Insect Swarm spell is increased by 12 sec.
	core.NewItemEffect(IdolOfTheSwarm, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		bonusDuration := time.Second * 12

		druid.GetOrRegisterAura(core.Aura{
			Label: "Idol of the Swarm",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				for _, spell := range druid.InsectSwarm {
					if spell != nil {
						for _, dot := range spell.Dots() {
							if dot != nil {
								dot.NumberOfTicks += 6
								dot.RecomputeAuraDuration()
							}
						}
					}
				}

				for _, aura := range druid.InsectSwarmAuras {
					if aura != nil && !aura.IsPermanent() {
						aura.Duration += bonusDuration
					}
				}
			},
		})
	})

	// https://www.wowhead.com/classic/item=234468/idol-of-ursin-power
	// Increases the damage of Swipe and Mangle by 3%.
	core.NewItemEffect(IdolOfUrsinPower, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Improved Swipe/Mangle",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_DruidSwipeBear | ClassSpellMask_DruidSwipeCat | ClassSpellMask_DruidMangleBear | ClassSpellMask_DruidMangleCat,
			IntValue:  3,
		}))
	})

	// https://www.wowhead.com/classic/item=216490/idol-of-wrath
	// Equip: Increases the damage of your Wrath spell by up to 2%.
	core.NewItemEffect(IdolOfWrath, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Improved Wrath",
		})).AttachSpellMod(core.SpellModConfig{
			ClassMask: ClassSpellMask_DruidWrath,
			Kind:      core.SpellMod_BaseDamageDone_Flat,
			IntValue:  2,
		})
	})

	// https://www.wowhead.com/classic/item=227183/knight-lieutenants-dragonhide-grips
	// Equip: Reduces the mana cost of your shapeshifts by 150.
	core.NewItemEffect(KnightLieutenantsDragonhideGrips, func(agent core.Agent) {
		registerDragonHideGripsAura(agent.(DruidAgent).GetDruid())
	})

	// https://www.wowhead.com/classic/item=231316/pristine-enchanted-south-seas-kelp
	// Increases the critical hit chance of Wrath and Starfire by 2%.
	core.NewItemEffect(PristineEnchantedSouthSeasKelp, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Improved Wrath/Starfire",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Flat,
			ClassMask:  ClassSpellMask_DruidWrath | ClassSpellMask_DruidStarfire,
			FloatValue: 2 * core.SpellCritRatingPerCritChance,
		}))
	})

	// https://www.wowhead.com/classic/item=224282/raelar
	// Equip: Chance on melee auto attack to steal 140 to 220 life from target enemy.
	core.NewItemEffect(Raelar, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 452430}
		healthMetrics := character.NewHealthMetrics(actionID)

		lifesteal := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// TODO: Unsure if this can crit but it we're assuming Gla'sir can
				result := spell.CalcAndDealDamage(sim, target, sim.Roll(140, 220), spell.OutcomeMagicHitAndCrit)
				character.GainHealth(sim, result.Damage, healthMetrics)
			},
		})

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Rae'lar Damage Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskWhiteHit,
			PPM:      1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				lifesteal.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterProc(Raelar, triggerAura)
	})

	// https://www.wowhead.com/classic-ptr/item=240842/scarlet-rotbringer
	// Requires Cat Form, Bear Form, Dire Bear Form
	// Equip: Chance on melee attack to inflict your enemy with a rot, dealing 100 Nature damage every 2 sec to all enemies within an 8 yard radius of the caster for 12 sec. The attack speed of these targets is also slowed by 15%.
	// TODO: PPM not confirmed
	itemhelpers.CreateWeaponProcSpell(ScarletRotbringer, "Scarlet Rotbringer", 1.0, func(character *core.Character) *core.Spell {
		druid := character.Env.GetAgentFromUnit(&character.Unit).(DruidAgent).GetDruid()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1231261},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
			Flags:       core.SpellFlagDisease | core.SpellFlagPureDot | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return druid.InForm(Cat | Bear)
			},

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Rot (Scarlet Rotbringer)",
				},
				NumberOfTicks: 6,
				TickLength:    time.Second * 2,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 100, isRollover)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					for _, aoeTarget := range sim.Encounter.TargetUnits {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
					}
				},
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.Dot(target).Apply(sim)
			},
		})

		for _, dot := range spell.Dots() {
			if dot != nil {
				core.AtkSpeedReductionEffect(dot.Aura, 1.15)
			}
		}

		return spell
	})

	core.NewItemEffect(RitualistsHammer, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.newBloodbarkCleaveItem(RitualistsHammer)
	})

	// https://www.wowhead.com/classic-ptr/item=240849/staff-of-the-glade
	// Equip: Remaining in Cat Form for 5 seconds, causes your Energy Regeneration to increase by 100%, and the damage of your Ferocious Bite to increase by 100%.
	// Equip: You may cast Rebirth and Innervate while in Cat Form.
	core.NewItemEffect(StaffOfTheGlade, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		// https://www.wowhead.com/classic-ptr/spell=1231381/feral-dedication
		auraBuff := druid.RegisterAura(core.Aura{
			ActionID: core.ActionID{
				SpellID: 1231381,
			},
			Duration: core.NeverExpires,
			Label:    "Feral Dedication",
		}).AttachSpellMod(core.SpellModConfig{
			ClassMask:  ClassSpellMask_DruidFerociousBite,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 2.0,
		}).AttachSpellMod(core.SpellModConfig{
			ClassMask: ClassSpellMask_DruidFerociousBite,
			Kind:      core.SpellMod_Custom,
			ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
				druid.EnergyTickMultiplier *= 1.5
			},
			RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
				druid.EnergyTickMultiplier /= 1.5
			},
		})

		// https://www.wowhead.com/classic-ptr/spell=1231380/feral-dedication
		auraTimer := druid.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{
				SpellID: 1231380,
			},
			Label:    "Feral Dedication (Timer)",
			Duration: time.Second * 5,
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if druid.CatFormAura.IsActive() {
					auraBuff.Activate(sim)
				}
			},
		})

		// https://www.wowhead.com/classic-ptr/spell=1231382/feral-dedication
		// Also handle https://www.wowhead.com/classic-ptr/spell=1232896/staff-of-the-glade
		druid.ItemSwap.RegisterProcWithSlots(StaffOfTheGlade, core.MakePermanent(druid.GetOrRegisterAura(core.Aura{
			Label: "Feral Dedication (Passive)",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				// May already be in cat form.
				auraTimer.Activate(sim)
				druid.Innervate.FormMask |= Cat
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				auraTimer.Deactivate(sim)
				auraBuff.Deactivate(sim)
				druid.Innervate.FormMask &= ^Cat
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell == druid.CatForm.Spell {
					auraBuff.Deactivate(sim)
					if druid.CatFormAura.IsActive() {
						auraTimer.Activate(sim)
					}
				}
			},
		})), []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand})
	})

	// https://www.wowhead.com/classic/item=231280/wushoolays-charm-of-nature
	// Use: Aligns the Druid with nature, increasing the damage done by spells by 10%, improving heal effects by 10%, and increasing the critical strike chance of spells by 10% for 20 sec.
	// (2 Min Cooldown)
	core.NewItemEffect(WushoolaysCharmOfNature, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: WushoolaysCharmOfNature}
		duration := time.Second * 20

		// TODO: healing dealt multiplier?
		aura := character.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Aligned with Nature",
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1.10)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1 / 1.10)
			},
			// TODO: healing dealt multiplier?
		}).AttachStatBuff(stats.SpellCrit, 10*core.SpellCritRatingPerCritChance)

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.AddEffectsToTest = true
}

// https://www.wowhead.com/classic/item=213407/catnip
func (druid *Druid) registerCatnipCD() {
	if druid.Consumes.MiscConsumes == nil || !druid.Consumes.MiscConsumes.Catnip {
		return
	}
	sod.RegisterFiftyPercentHasteBuffCD(&druid.Character, core.ActionID{ItemID: Catnip})
}

func (druid *Druid) newBloodbarkCleaveItem(itemID int32) {
	auraActionID := core.ActionID{SpellID: 436482}

	results := make([]*core.SpellResult, min(3, druid.Env.GetNumTargets()))

	damageSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 436481},
		SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMelee, // actually has DefenseTypeNone, but is likely using the greatest CritMultiplier available
		ProcMask:    core.ProcMaskEmpty,

		// TODO: "Causes additional threat" in Tooltip, no clue what the multiplier is.
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, 5, spell.OutcomeMagicCrit)
				target = sim.Environment.NextTargetUnit(target)
			}
			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})

	buffAura := druid.GetOrRegisterAura(core.Aura{
		Label:    "Bloodbark Cleave",
		ActionID: auraActionID,
		Duration: 20 * time.Second,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask&core.ProcMaskMelee != 0 {
				damageSpell.Cast(sim, result.Target)
				return
			}
		},
	})

	mainSpell := druid.GetOrRegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{ItemID: itemID},
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell:    mainSpell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}

func registerDragonHideGripsAura(druid *Druid) {
	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label:    "Dragonhide Grips",
		ActionID: core.ActionID{SpellID: 459594},
		Duration: core.NeverExpires,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: ClassSpellMask_DruidForms,
		IntValue:  -150,
	}))
}
