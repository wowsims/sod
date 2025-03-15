package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	InfernalPactEssence         = 216509
	ZIlaGular                   = 223214
	ScytheOfChaos               = 229910
	TheBlackBook                = 230238
	HazzarahsCharmOfDestruction = 231284
	KezansUnstoppableTaint      = 231346
	PlagueheartRing             = 236067
	AtieshWarlock               = 236398
)

func init() {
	// https://www.wowhead.com/classic/item=236398/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshWarlock, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()
		aura := core.AtieshSpellPowerEffect(&warlock.Unit)
		warlock.ItemSwap.RegisterProc(AtieshWarlock, aura)

		for _, pet := range warlock.BasePets {
			petAura := core.AtieshSpellPowerEffect(&pet.Unit)
			warlock.ItemSwap.RegisterProc(AtieshWarlock, petAura)
		}
	})

	// https://www.wowhead.com/classic/item=231284/hazzarahs-charm-of-destruction
	// Increases your critical hit chance by 10%, and increases your pet's attack speed by 50% for 20 sec.
	// This spell does not affect temporary pets or Subjugated Demons.
	core.NewItemEffect(HazzarahsCharmOfDestruction, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		actionID := core.ActionID{ItemID: HazzarahsCharmOfDestruction}
		duration := time.Second * 20
		affectedPet := warlock.ActivePet

		buffAura := warlock.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Massive Destruction",
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				affectedPet = warlock.ActivePet
				if affectedPet != nil {
					affectedPet.MultiplyAttackSpeed(sim, 1.50)
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if affectedPet != nil {
					affectedPet.MultiplyAttackSpeed(sim, 1/1.50)
				}
			},
		}).AttachStatsBuff(stats.Stats{
			stats.SpellCrit: 10 * core.SpellCritRatingPerCritChance,
			stats.MeleeCrit: 10 * core.CritRatingPerCritChance,
		})

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    warlock.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Infernal Pact Essence
	core.NewItemEffect(InfernalPactEssence, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		stats := stats.Stats{
			stats.Stamina:   20,
			stats.Intellect: 80,
		}

		// TODO: Does this affect Infernal or Doomguard?
		warlock.Felhunter.AddStats(stats)
		warlock.Imp.AddStats(stats)
		warlock.Succubus.AddStats(stats)
		warlock.Voidwalker.AddStats(stats)
		if warlock.Felguard != nil {
			warlock.Felguard.AddStats(stats)
		}

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436479},
			SpellSchool: core.SpellSchoolShadow | core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
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

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=229910/scythe-of-chaos
	// Chance on direct damage spell to cause your next pet summoned within 20 sec to be instant cast and not consume a Soul Shard.
	// (Proc chance: 10%, 1m cooldown)
	// Use: Harvest the soul of your summoned demon, granting you an effect that lasts 15 sec.  The effect is canceled if any Demon is summoned. (1 Min Cooldown)
	core.NewItemEffect(ScytheOfChaos, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		summonBuffAura := warlock.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469211},
			Label:    "Scythe of Chaos",
			Duration: time.Second * 20,
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(ClassSpellMask_WarlockSummons) {
					aura.Deactivate(sim)
				}
			},
		}).AttachSpellMod(core.SpellModConfig{
			Kind:       core.SpellMod_CastTime_Pct,
			ClassMask:  ClassSpellMask_WarlockSummons,
			FloatValue: -1,
		})

		core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
			Name:       "Scythe of Chaos Trigger",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.10,
			ICD:        time.Minute,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				summonBuffAura.Activate(sim)
			},
		})

		harvestDemonDuration := time.Second * 15

		felhunterAura := warlock.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469229},
			Label:    "Hunter of Chaos",
			Duration: harvestDemonDuration,
			// Not doing anything for this one
			OnGain:   func(aura *core.Aura, sim *core.Simulation) {},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {},
		})

		impAura := warlock.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469222},
			Label:    "Impish Delight",
			Duration: harvestDemonDuration,
			// Not doing anything for this one
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.15
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.15
			},
		})

		succubusAura := warlock.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469225},
			Label:    "Seduction of the Shadows",
			Duration: harvestDemonDuration,
			// Not doing anything for this one
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.15
			},
		})

		voidwalkerAura := warlock.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469224},
			Label:    "Void Walking",
			Duration: harvestDemonDuration,
			// Not doing anything for this one
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageTakenMultiplier.AddToAllSchools(-100)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageTakenMultiplier.AddToAllSchools(100)
			},
		})

		felguardAura := warlock.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469230},
			Label:    "Fel Invigoration",
			Duration: harvestDemonDuration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= .75
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= .75
			},
		})

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: ScytheOfChaos},
			SpellSchool: core.SpellSchoolShadow,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 1,
				},
			},
			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return warlock.ActivePet != nil
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				switch warlock.ActivePet {
				case warlock.Felhunter:
					felhunterAura.Activate(sim)
				case warlock.Imp:
					impAura.Activate(sim)
				case warlock.Succubus:
					succubusAura.Activate(sim)
				case warlock.Voidwalker:
					voidwalkerAura.Activate(sim)
				case warlock.Felguard:
					felguardAura.Activate(sim)
				}

				warlock.changeActivePet(sim, nil, false)
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				// Must be manually activated
				return false
			},
		})
	})

	// https://www.wowhead.com/classic/item=231346/kezans-unstoppable-taint
	// Reduces the cooldown of your Felguard's Cleave spell by 2 sec.
	core.NewItemEffect(KezansUnstoppableTaint, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		warlock.RegisterAura(core.Aura{
			Label: "Reduced Cleave Cooldown",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				warlock.Succubus.primaryAbility.CD.ApplyFlatCooldownMod(-time.Second * 2)
				if warlock.Felguard != nil {
					warlock.Felguard.primaryAbility.CD.ApplyFlatCooldownMod(-time.Second * 2)
				}
			},
		})
	})

	// https://www.wowhead.com/classic/item=236067/plagueheart-ring
	// Equip: Increases the damage dealt by your damage over time spells by 2%.
	core.NewItemEffect(PlagueheartRing, func(agent core.Agent) {
		priest := agent.(WarlockAgent).GetWarlock()

		priest.OnSpellRegistered(func(spell *core.Spell) {
			// Unlike the Priest ring, the Warlock ring doesn't seem to affect channels https://www.wowhead.com/classic/spell=1222974/damage-over-time-increase
			if spell.Matches(ClassSpellMask_WarlockAll) && !spell.Flags.Matches(core.SpellFlagChanneled) {
				spell.ApplyAdditivePeriodicDamageBonus(2)
			}
		})
	})

	// https://www.wowhead.com/classic/item=230238/the-black-book
	// Empowers your pet, increasing pet damage by 100% and increasing pet armor by 100% for 30 sec.
	// This spell does not affect temporary pets or Subjugated Demons.
	core.NewItemEffect(TheBlackBook, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		actionID := core.ActionID{ItemID: TheBlackBook}
		duration := time.Second * 30

		statDeps := map[string]*stats.StatDependency{}
		for _, pet := range warlock.BasePets {
			statDeps[pet.Name] = pet.NewDynamicMultiplyStat(stats.Armor, 2)
		}

		var affectedPet *WarlockPet

		buffAura := warlock.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Blessing of the Black Book",
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				affectedPet = warlock.ActivePet
				if affectedPet != nil {
					affectedPet.PseudoStats.DamageDealtMultiplier *= 2.0
					affectedPet.EnableDynamicStatDep(sim, statDeps[affectedPet.Name])
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if affectedPet != nil {
					affectedPet.PseudoStats.DamageDealtMultiplier /= 2.0
					affectedPet.DisableDynamicStatDep(sim, statDeps[affectedPet.Name])
				}
			},
		})

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolShadow,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 3,
				},
				SharedCD: core.Cooldown{
					Timer:    warlock.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Zila Gular
	core.NewItemEffect(ZIlaGular, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()

		warlock.zilaGularAura = warlock.GetOrRegisterAura(core.Aura{
			Label:    "Zila Gular",
			ActionID: core.ActionID{SpellID: 448686},
			Duration: time.Second * 20,
		})

		spell := warlock.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 448686},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    warlock.GetOffensiveTrinketCD(),
					Duration: time.Second * 20,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				warlock.zilaGularAura.Activate(sim)
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})
}
