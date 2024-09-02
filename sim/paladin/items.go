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
	GrileksCharmOfValor                  = 231285
	HerosBrand							 = 231328
	ZandalarFreethinkersBreastplate      = 231329
	ZandalarFreethinkersBelt             = 231330	
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
				paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe), 15.0 * core.CritRatingPerCritChance, 10.0 * core.CritRatingPerCritChance)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += core.TernaryFloat64(paladin.hasRune(proto.PaladinRune_RuneCloakShockAndAwe), -15.0 * core.CritRatingPerCritChance, -10.0 * core.CritRatingPerCritChance)
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
}
