package paladin

import (
	"time"

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
	HerosBrand							 = 231328
	zandalarFreethinkersBreastplate      = 231329
	zandalarFreethinkersBelt             = 231330
	grileksCharmOfValor                  = 231285
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

	
	core.NewItemEffect(zandalarFreethinkersBreastplate, func(agent core.Agent) {
		character := agent.GetCharacter()		
		character.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexHoly] += 3
	})

	
	core.NewItemEffect(zandalarFreethinkersBelt, func(agent core.Agent) {
		character := agent.GetCharacter()		
		character.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexHoly] += 2
	})


	// https://www.wowhead.com/classic/item=231285/grileks-charm-of-valor
	// Use: Increases the critical hit chance of Holy spells by 10% for 15 sec. If Shock and Awe is engraved, gain an additional 5%. (1 Min, 30 Sec Cooldown)
	core.NewItemEffect(grileksCharmOfValor, func(agent core.Agent) {
		character := agent.GetCharacter()
		paladin := agent.(PaladinAgent).GetPaladin()

		character.PseudoStats.BonusDamage += 4

		duration := time.Second * 15

		aura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: grileksCharmOfValor},
			Label:    "Brilliant Light",
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellCrit, core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe), 15.0 * core.CritRatingPerCritChance, 10.0 * core.CritRatingPerCritChance))
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellCrit, core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe), -15.0 * core.CritRatingPerCritChance, -10.0 * core.CritRatingPerCritChance))
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: grileksCharmOfValor},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 90,
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
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})
}
