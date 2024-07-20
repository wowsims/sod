package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	CraftOfTheShadows     = 227280
	DukesDomain           = 227915
	AccursedChalice       = 228078
	GerminatingPoisonseed = 228081
	GloamingTreeheart     = 228083
	WoodcarvedMoonstalker = 228089
	TheMoltenCore         = 228122
	FistOfTheFiresworn    = 228139
	TreantsBane           = 228486
	BeastmastersBoots	  = 226903
	WhistleOfTheBeast	  = 228432
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=228139/fist-of-the-firesworn
	// Chance on hit: Blasts the enemy for 70 Fire damage.
	// TODO: Proc rate assumed and needs testing
	itemhelpers.CreateWeaponProcDamage(FistOfTheFiresworn, "Fist of the Firesworn", 1.0, 461896, core.SpellSchoolFire, 70, 0, 0.15, core.DefenseTypeMagic)

	// https://www.wowhead.com/classic/item=228486/treants-bane
	// Equip: +75 Attack Power when fighting Elementals.
	core.NewItemEffect(TreantsBane, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeElemental {
			character.PseudoStats.MobTypeAttackPower += 75
		}
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=227915/dukes-domain
	// Use: Expand the Duke's Domain, increasing the Fire Resistance of those who reside within by 50. Lasts for 15 sec. (1 Min, 30 Sec Cooldown)
	// TODO: Raid-wide effect if we ever do raid sim
	core.NewSimpleStatDefensiveTrinketEffect(DukesDomain, stats.Stats{stats.FireResistance: 50}, time.Second*15, time.Second*90)

	// https://www.wowhead.com/classic/item=228078/accursed-chalice
	// Use: Increases your Strength by 80.  Effect lasts for 20 sec. (2 Min Cooldown)
	core.NewSimpleStatOffensiveTrinketEffect(AccursedChalice, stats.Stats{stats.Strength: 80}, time.Second*20, time.Minute*2)

	// https://www.wowhead.com/classic/item=228081/germinating-poisonseed
	// Use: Increases your Nature Damage by up to 115.  Effect lasts for 20 sec. (2 Min Cooldown)
	core.NewSimpleStatOffensiveTrinketEffect(GerminatingPoisonseed, stats.Stats{stats.NaturePower: 115}, time.Second*20, time.Minute*2)

	// https://www.wowhead.com/classic/item=228083/gloaming-treeheart
	// Use: Increases your Nature Resistance by 90.  Effect lasts for 30 sec. (3 Min Cooldown)
	core.NewSimpleStatDefensiveTrinketEffect(GloamingTreeheart, stats.Stats{stats.NatureResistance: 90}, time.Second*30, time.Minute*3)

	// https://www.wowhead.com/classic/item=228089/woodcarved-moonstalker
	// Use: Increases your Strength by 60.  Effect lasts for 15 sec. (1 Min, 30 Sec Cooldown)
	core.NewItemEffect(WoodcarvedMoonstalker, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.PseudoStats.BonusDamage += 4

		aura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: WoodcarvedMoonstalker},
			Label:    "Woodcarved Moonstalker",
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.Strength, 60)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.Strength, -60)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: WoodcarvedMoonstalker},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 90,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	// https://www.wowhead.com/classic/item=228122/the-molten-core
	// Equip: Inflicts 20 Fire damage to nearby enemies every 2 sec.
	core.NewItemEffect(TheMoltenCore, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 461228}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			BonusCoefficient: .045,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 20, spell.OutcomeAlwaysHit)
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "The Molten Core Trigger",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   time.Second * 2,
					Priority: core.ActionPriorityDOT, // High prio
					OnAction: func(sim *core.Simulation) {
						for _, aoeTarget := range sim.Encounter.TargetUnits {
							spell.Cast(sim, aoeTarget)
						}
					},
				})
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Other
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=227280/craft-of-the-shadows
	// Equip: Increases your maximum Energy by 10.
	core.NewItemEffect(CraftOfTheShadows, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.HasEnergyBar() {
			character.EnableEnergyBar(character.MaxEnergy() + 10)
		}
	})


	core.NewItemEffect(BeastmastersBoots, func(agent core.Agent) {
		character := agent.GetCharacter()

		for _, pet := range character.Pets {
			pet.PseudoStats.DamageDealtMultiplier *= 1.03
		}
	})

	core.NewItemEffect(WhistleOfTheBeast, func(agent core.Agent) {
		character := agent.GetCharacter()

		for _, pet := range character.Pets {
			pet.PseudoStats.DamageDealtMultiplier *= 1.03
			pet.MultiplyStat(stats.Health, 1.03)
			pet.MultiplyStat(stats.Armor, 1.10)
			pet.AddStat(stats.MeleeCrit, 2*core.CritRatingPerCritChance)
		}
	})

	core.AddEffectsToTest = true
}
