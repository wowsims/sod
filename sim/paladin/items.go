package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/common/vanilla"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Libram IDs
const (
	SanctifiedOrb                        = 20512
	LibramOfHope                         = 22401
	LibramOfFervor                       = 23203
	LibramDiscardedTenetsOfTheSilverHand = 209574
	LibramOfBenediction                  = 215435
	LibramOfDraconicDestruction          = 221457
	LibramOfTheDevoted                   = 228174
	LibramOfAvenging                     = 232421
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
)

func init() {
	core.NewSimpleStatOffensiveTrinketEffect(SanctifiedOrb, stats.Stats{stats.MeleeCrit: 3 * core.CritRatingPerCritChance, stats.SpellCrit: 3 * core.CritRatingPerCritChance}, time.Second*25, time.Minute*3)

	core.NewItemEffect(LibramDiscardedTenetsOfTheSilverHand, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon || character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.PseudoStats.MobTypeAttackPower += 15
		}
	})

	core.NewItemEffect(LibramOfDraconicDestruction, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDragonkin {
			character.PseudoStats.MobTypeAttackPower += 36
		}
	})

	core.NewItemEffect(HerosBrand, func(agent core.Agent) {
		//Increases critical strike chance of holy shock spell by 2%
		paladin := agent.GetCharacter()
		paladin.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_PaladinHolyShock {
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
				character.MultiplyAttackSpeed(sim, 1.25)
				character.MultiplyCastSpeed(1.33)

				// Crusader's zeal proc overwrites scrolls regardless of time left
				truthbearerAura := character.GetAuraByID(core.ActionID{SpellID: 465414})
				if truthbearerAura != nil {
					truthbearerAura.Deactivate(sim)
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.0/1.25)
				character.MultiplyCastSpeed(1.0 / 1.33)
			},
		})

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
		vanilla.BlazefuryTriggerAura(character, 465412, core.SpellSchoolHoly, 4)
		crusadersZealAura465414(character)
	})

	// https://www.wowhead.com/classic/item=229806/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer1H, func(agent core.Agent) {
		character := agent.GetCharacter()
		crusadersZealAura465414(character)
	})

	// https://www.wowhead.com/classic/item=229749/truthbearer
	// Chance on hit: Increases damage done by 15 and attack speed by 30% for 8 sec.
	// TODO: Proc rate assumed and needs testing
	core.NewItemEffect(Truthbearer2H, func(agent core.Agent) {
		character := agent.GetCharacter()
		vanilla.BlazefuryTriggerAura(character, 465412, core.SpellSchoolHoly, 4)
		crusadersZealAura465414(character)
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
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe), 15.0*core.CritRatingPerCritChance, 10.0*core.CritRatingPerCritChance)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe), -15.0*core.CritRatingPerCritChance, -10.0*core.CritRatingPerCritChance)
			},
		})

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
		holyWrathSpells := []*core.Spell{}

		buffAura := paladin.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 470246},
			Label:     "Libram Of Wrath Buff",
			Duration:  time.Second * 15,
			MaxStacks: 5,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				holyWrathSpells = core.FilterSlice(paladin.holyWrath, func(spell *core.Spell) bool { return spell != nil })
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				core.Each(holyWrathSpells, func(spell *core.Spell) {
					spell.CastTimeMultiplier += (0.2 * float64(oldStacks))
					spell.Cost.Multiplier += int32(100.0 * (0.2 * float64(oldStacks)))
					spell.DamageMultiplier -= (0.2 * float64(oldStacks))

					spell.CastTimeMultiplier -= (0.2 * float64(newStacks))
					spell.Cost.Multiplier -= int32(100.0 * (0.2 * float64(newStacks)))
					spell.DamageMultiplier += (0.2 * float64(newStacks))

				})
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellCode == SpellCode_PaladinHolyWrath {
					aura.Deactivate(sim)
				}
			},
		})

		paladin.RegisterAura(core.Aura{
			Label:    "Libram Of Wrath Trigger",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {

				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellCode == SpellCode_PaladinHolyShock {
					buffAura.Activate(sim)
					buffAura.AddStack(sim)
				}
			},
		})
	})

	core.NewItemEffect(LibramOfTheDevoted, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()

		actionID := core.ActionID{SpellID: 461309}
		devotedMetrics := paladin.NewManaMetrics(actionID)

		core.MakePermanent(paladin.RegisterAura(core.Aura{
			Label: "Libram of the Devoted Trigger",
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.DidBlock() {
					amount := paladin.MaxMana() * 0.03
					paladin.AddMana(sim, amount, devotedMetrics)
				}
			},
		}))
	})

	core.NewItemEffect(LibramOfTheExorcist, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		paladin.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_PaladinCrusaderStrike || spell.SpellCode == SpellCode_PaladinExorcism {
				// Increases the damage of Exorcism and Crusader Strike by 3%.
				spell.DamageMultiplier += 0.03
			}
		})
	})

	core.NewItemEffect(LibramOfSanctity, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		sanctityLibramAura := ApplyLibramOfSanctityAura(&paladin.Unit)

		paladin.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_PaladinHolyShock {
				// Increases the damage of Holy Shock by 3%, and your Shock and Awe buff now also grants 10% increased Holy Damage. (This effect does not stack with Sanctity Aura).
				spell.DamageMultiplier += 0.03
				originalApplyEffects := spell.ApplyEffects

				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					originalApplyEffects(sim, target, spell)
					sanctityLibramAura.Activate(sim)
				}
			}
		})
	})

	core.NewItemEffect(LibramOfRighteousness, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		paladin.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_PaladinHammerOfTheRighteous || spell.SpellCode == SpellCode_PaladinShieldOfRighteousness {
				// Increases the damage of Hammer of the Righteous and Shield of Righteousness by 3%.
				spell.DamageMultiplier += 0.03
			}
		})
	})

}

// https://www.wowhead.com/classic/spell=465414/crusaders-zeal
// Used by:
// - https://www.wowhead.com/classic/item=229806/truthbearer and
// - https://www.wowhead.com/classic/item=229749/truthbearer
func crusadersZealAura465414(character *core.Character) *core.Aura {
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

	return core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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
}

func ApplyLibramOfSanctityAura(unit *core.Unit) *core.Aura {

	sanctityLibramAura := unit.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1214298},
		Label:    "Libram of Sanctity",
		Duration: time.Minute * 1,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			sanctityAura := aura.Unit.GetAuraByID(core.ActionID{SpellID: 20218})

			if sanctityAura == nil || !sanctityAura.IsActive() {
				aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1.1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			sanctityAura := aura.Unit.GetAuraByID(core.ActionID{SpellID: 20218})

			if sanctityAura == nil || !sanctityAura.IsActive() {
				aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= 1.1
			}
		},
	})

	return sanctityLibramAura

}
