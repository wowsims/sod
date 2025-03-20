package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/common/vanilla"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SanctifiedOrb                        = 20512
	LibramOfHope                         = 22401
	LibramOfFervor                       = 23203
	LibramDiscardedTenetsOfTheSilverHand = 209574
	LibramOfBenediction                  = 215435
	LibramOfDraconicDestruction          = 221457
	LibramOfTheDevoted                   = 228174
	LibramOfAvenging                     = 232421
	HammerOfTheFallenThane               = 227935
	Truthbearer2H                        = 229749
	Truthbearer1H                        = 229806
	HammerOfTheLightbringer              = 230003
	ScrollsOfBlindingLight               = 230272
	GrileksCharmOfValor                  = 231285
	HerosBrand                           = 231328
	ZandalarFreethinkersBreastplate      = 231329
	ZandalarFreethinkersBelt             = 231330
	LibramOfWrath                        = 232420
	LibramOfTheExorcist                  = 234475
	LibramOfSanctity                     = 234476
	LibramOfRighteousness                = 234477
	BandOfRedemption                     = 236130
	ClaymoreOfUnholyMight                = 236299
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(BandOfRedemption, func(agent core.Agent) {
		character := agent.GetCharacter()
		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Band of Redemption Trigger",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.02,
			ICD:        time.Millisecond * 200,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				character.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 1223010}, spell)
			},
		})
		character.ItemSwap.RegisterProc(BandOfRedemption, triggerAura)
	})

	// https://www.wowhead.com/classic/item=236299/claymore-of-unholy-might
	// Chance on hit: Increases the wielder's Strength by 350, but they also take 5% more damage from all sources for 8 sec.
	// TODO: Proc rate assumed and needs testing
	itemhelpers.CreateWeaponProcAura(ClaymoreOfUnholyMight, "Claymore of Unholy Might", 2.0, core.UnholyMightAura)

	core.NewSimpleStatOffensiveTrinketEffect(SanctifiedOrb, stats.Stats{stats.MeleeCrit: 3 * core.CritRatingPerCritChance, stats.SpellCrit: 3 * core.CritRatingPerCritChance}, time.Second*25, time.Minute*3)

	core.NewItemEffect(LibramDiscardedTenetsOfTheSilverHand, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon || character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			aura := core.MakePermanent(character.RegisterAura(core.Aura{
				Label: "Libram Discarded Tenets Of The SilverHand",
			}).AttachAdditivePseudoStatBuff(&character.PseudoStats.MobTypeAttackPower, 15))

			character.ItemSwap.RegisterProc(LibramDiscardedTenetsOfTheSilverHand, aura)
		}
	})

	core.NewItemEffect(LibramOfDraconicDestruction, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			aura := core.MakePermanent(character.RegisterAura(core.Aura{
				Label: "Libram Of Draconic Destruction",
			}).AttachAdditivePseudoStatBuff(&character.PseudoStats.MobTypeAttackPower, 36))

			character.ItemSwap.RegisterProc(LibramOfDraconicDestruction, aura)
		}
	})

	core.NewItemEffect(HerosBrand, func(agent core.Agent) {
		//Increases critical strike chance of holy shock spell by 2%
		paladin := agent.GetCharacter()
		paladin.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(ClassSpellMask_PaladinHolyShock) {
				spell.BonusCritRating += 2.0
			}
		})
	})

	// https://www.wowhead.com/classic/item=230272/scrolls-of-blinding-light
	// Use: Energizes a Paladin with light, increasing melee attack speed by 25%
	// and spell casting speed by 33% for 20 sec. (2 Min Cooldown)
	core.NewItemEffect(ScrollsOfBlindingLight, func(agent core.Agent) {
		character := agent.GetCharacter()

		duration := time.Second * 20

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Blinding Light",
			ActionID: core.ActionID{SpellID: 467522},
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				// Crusader's zeal proc overwrites scrolls regardless of time left
				truthbearerAura := character.GetAuraByID(core.ActionID{SpellID: 465414})
				if truthbearerAura != nil {
					truthbearerAura.Deactivate(sim)
				}
			},
		}).
			AttachMultiplyAttackSpeed(&character.Unit, 1.25).
			AttachMultiplyCastSpeed(&character.Unit, 1.33)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: ScrollsOfBlindingLight},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			ProcMask: core.ProcMaskEmpty,

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

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:     core.CooldownTypeDPS,
			Priority: core.CooldownPriorityDefault,
			Spell:    spell,
		})
	})

	core.NewItemEffect(HammerOfTheLightbringer, func(agent core.Agent) {
		character := agent.GetCharacter()
		vanilla.BlazefuryTriggerAura(character, HammerOfTheLightbringer, 465412, 465411, core.SpellSchoolHoly, 4)
		crusadersZealAura465414(character, HammerOfTheLightbringer)
	})

	// https://www.wowhead.com/classic/item=229806/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer1H, func(agent core.Agent) {
		character := agent.GetCharacter()
		crusadersZealAura465414(character, Truthbearer1H)
	})

	// https://www.wowhead.com/classic/item=229749/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer2H, func(agent core.Agent) {
		character := agent.GetCharacter()
		vanilla.BlazefuryTriggerAura(character, Truthbearer2H, 465412, 465411, core.SpellSchoolHoly, 4)
		crusadersZealAura465414(character, Truthbearer2H)
	})

	core.NewItemEffect(ZandalarFreethinkersBreastplate, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexHoly] += 3
	})

	core.NewItemEffect(ZandalarFreethinkersBelt, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexHoly] += 2
	})

	// https://www.wowhead.com/classic/item=231285/grileks-charm-of-valor
	// Use: Increases the critical hit chance of Holy spells by 10% for 15 sec. If Shock and Awe is engraved, gain an additional 5%. (1 Min, 30 Sec Cooldown)
	core.NewItemEffect(GrileksCharmOfValor, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		duration := time.Second * 15

		aura := paladin.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: GrileksCharmOfValor},
			Label:    "Brilliant Light",
			Duration: duration,
		}).AttachAdditivePseudoStatBuff(
			&paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly],
			core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe), 15.0*core.CritRatingPerCritChance, 10.0*core.CritRatingPerCritChance),
		)

		spell := paladin.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: GrileksCharmOfValor},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    paladin.NewTimer(),
					Duration: time.Second * 90,
				},
				SharedCD: core.Cooldown{
					Timer:    paladin.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		paladin.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	//https://www.wowhead.com/classic/item=232420/libram-of-wrath
	//Equip: Your Holy Shock spell reduces the cast time and mana cost of your next Holy Wrath spell cast within 15 sec by 20%, and increases its damage by 20%. Stacking up to 5 times.
	core.NewItemEffect(LibramOfWrath, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		costMod := paladin.AddDynamicMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Pct,
			ClassMask: ClassSpellMask_PaladinHolyWrath,
		})
		castMod := paladin.AddDynamicMod(core.SpellModConfig{
			Kind:      core.SpellMod_CastTime_Pct,
			ClassMask: ClassSpellMask_PaladinHolyWrath,
		})
		damageMod := paladin.AddDynamicMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_PaladinHolyWrath,
		})

		buffAura := paladin.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 470246},
			Label:     "Libram Of Wrath Buff",
			Duration:  time.Second * 15,
			MaxStacks: 5,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				castMod.UpdateFloatValue(-0.2 * float64(newStacks))
				costMod.UpdateIntValue(-int64(100.0 * (0.2 * float64(newStacks))))
				damageMod.UpdateIntValue(int64(20 * newStacks))
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				castMod.Activate()
				costMod.Activate()
				damageMod.Activate()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				castMod.Deactivate()
				costMod.Deactivate()
				damageMod.Deactivate()
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(ClassSpellMask_PaladinHolyWrath) {
					aura.Deactivate(sim)
				}
			},
		})

		triggerAura := core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
			Name:           "Libram Of Wrath Trigger",
			ClassSpellMask: ClassSpellMask_PaladinHolyShock,
			Callback:       core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			},
		})

		paladin.ItemSwap.RegisterProc(LibramOfWrath, triggerAura)
	})

	core.NewItemEffect(LibramOfTheDevoted, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		actionID := core.ActionID{SpellID: 461309}
		devotedMetrics := paladin.NewManaMetrics(actionID)

		procAura := core.MakePermanent(paladin.RegisterAura(core.Aura{
			Label: "Libram of the Devoted Trigger",
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.DidBlock() {
					amount := paladin.MaxMana() * 0.03
					paladin.AddMana(sim, amount, devotedMetrics)
				}
			},
		}))

		paladin.ItemSwap.RegisterProc(LibramOfTheDevoted, procAura)
	})

	core.NewItemEffect(LibramOfTheExorcist, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		// Increases the damage of Exorcism and Crusader Strike by 3%.
		aura := core.MakePermanent(paladin.RegisterAura(core.Aura{
			Label: "Libram of the Exorcist",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_PaladinCrusaderStrike | ClassSpellMask_PaladinExorcism,
			IntValue:  3,
		}))

		paladin.ItemSwap.RegisterProc(LibramOfTheExorcist, aura)
	})

	core.NewItemEffect(LibramOfSanctity, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		buffAura := core.MakePermanent(paladin.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 1214298},
			Label:    "Libram of Sanctity",
			Duration: time.Minute,
		}))

		core.ExclusiveHolyDamageDealtAura(buffAura, 1.1)

		aura := core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
			Name:           "Improved Holy Shock",
			ClassSpellMask: ClassSpellMask_PaladinHolyShock,
			Callback:       core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		}).AttachSpellMod(core.SpellModConfig{
			// Increases the damage of Holy Shock by 3%, and your Shock and Awe buff now also grants 10% increased Holy Damage. (This effect does not stack with Sanctity Aura).
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_PaladinHolyShock,
			IntValue:  3,
		})

		paladin.ItemSwap.RegisterProc(LibramOfSanctity, aura)
	})

	core.NewItemEffect(LibramOfRighteousness, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		aura := core.MakePermanent(paladin.RegisterAura(core.Aura{
			Label: "Libram of Righteousness",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Flat,
			ClassMask: ClassSpellMask_PaladinHammerOfTheRighteous | ClassSpellMask_PaladinShieldOfRighteousness,
			IntValue:  3,
		}))

		paladin.ItemSwap.RegisterProc(LibramOfRighteousness, aura)
	})

	// https://www.wowhead.com/classic/item=227935/hammer-of-the-fallen-thane
	// Chance on hit: Steals 42 life from target enemy.
	itemhelpers.CreateWeaponProcSpell(HammerOfTheFallenThane, "Hammer of the Fallen Thane", 7.0, func(character *core.Character) *core.Spell {
		actionID := core.ActionID{SpellID: 462221}
		healthMetrics := character.NewHealthMetrics(actionID)

		return character.RegisterSpell(core.SpellConfig{
			ActionID:         actionID,
			SpellSchool:      core.SpellSchoolFire,
			DefenseType:      core.DefenseTypeMagic,
			ProcMask:         core.ProcMaskEmpty,
			BonusCoefficient: 0.2,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealDamage(sim, target, 42, spell.OutcomeAlwaysHit)
				character.GainHealth(sim, result.Damage, healthMetrics)
			},
		})
	})

	core.AddEffectsToTest = true
}

// https://www.wowhead.com/classic/spell=465414/crusaders-zeal
// Used by:
// - https://www.wowhead.com/classic/item=229806/truthbearer and
// - https://www.wowhead.com/classic/item=229749/truthbearer
func crusadersZealAura465414(character *core.Character, itemID int32) *core.Aura {
	procAura := character.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 465414},
		Label:    "Crusader's Zeal",
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.BonusPhysicalDamage += 15
			character.MultiplyAttackSpeed(sim, 1.30)

			// Crusader's zeal proc overwrites scrolls regardless of time left
			scrollsAura := character.GetAuraByID(core.ActionID{SpellID: 467522})
			if scrollsAura != nil {
				scrollsAura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			character.PseudoStats.BonusPhysicalDamage -= 15
			character.MultiplyAttackSpeed(sim, 1/1.30)
		},
	})

	triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:              "Truthbearer (Crusader's Zeal)",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		ProcMask:          core.ProcMaskMelee,
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		PPM:               2.0,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
		},
	})

	character.ItemSwap.RegisterProc(itemID, triggerAura)

	return triggerAura
}
