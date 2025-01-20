package item_effects

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	BulwarkOfIre = 235868

	// Atiesh
	AtieshSpellPower = 236398
	AtieshHealing    = 236399
	AtieshCastSpeed  = 236400
	AtieshSpellCrit  = 236401

	// Seals of the Dawn
	SquiresSealOfTheDawnDamage    = 236356
	KnightsSealOfTheDawnDamage    = 236357
	TemplarsSealOfTheDawnDamage   = 236358
	ChampionsSealOfTheDawnDamage  = 236360
	VanguardsSealOfTheDawnDamage  = 236361
	CrusadersSealOfTheDawnDamage  = 236362
	CommandersSealOfTheDawnDamage = 236363
	HighlordsSSealOfTheDawnDamage = 236364

	SquiresSealOfTheDawnHealing    = 236383
	KnightsSealOfTheDawnHealing    = 236382
	TemplarsSealOfTheDawnHealing   = 236380
	ChampionsSealOfTheDawnHealing  = 236379
	VanguardsSealOfTheDawnHealing  = 236378
	CrusadersSealOfTheDawnHealing  = 236376
	CommandersSealOfTheDawnHealing = 236375
	HighlordsSSealOfTheDawnHealing = 236374

	SquiresSealOfTheDawnTanking    = 236394
	KnightsSealOfTheDawnTanking    = 236393
	TemplarsSealOfTheDawnTanking   = 236392
	ChampionsSealOfTheDawnTanking  = 236391
	VanguardsSealOfTheDawnTanking  = 236390
	CrusadersSealOfTheDawnTanking  = 236389
	CommandersSealOfTheDawnTanking = 236388
	HighlordsSSealOfTheDawnTanking = 236386
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////
	// https://www.wowhead.com/classic/item=236356/squires-seal-of-the-dawn
	core.NewItemEffect(SquiresSealOfTheDawnDamage, sanctifiedDamageEffect)
	core.NewItemEffect(KnightsSealOfTheDawnDamage, sanctifiedDamageEffect)
	core.NewItemEffect(TemplarsSealOfTheDawnDamage, sanctifiedDamageEffect)
	core.NewItemEffect(ChampionsSealOfTheDawnDamage, sanctifiedDamageEffect)
	core.NewItemEffect(VanguardsSealOfTheDawnDamage, sanctifiedDamageEffect)
	core.NewItemEffect(CrusadersSealOfTheDawnDamage, sanctifiedDamageEffect)
	core.NewItemEffect(CommandersSealOfTheDawnDamage, sanctifiedDamageEffect)
	core.NewItemEffect(HighlordsSSealOfTheDawnDamage, sanctifiedDamageEffect)

	// https://www.wowhead.com/classic/item=236383/squires-seal-of-the-dawn
	core.NewItemEffect(SquiresSealOfTheDawnHealing, sanctifiedHealingEffect)
	core.NewItemEffect(KnightsSealOfTheDawnHealing, sanctifiedHealingEffect)
	core.NewItemEffect(TemplarsSealOfTheDawnHealing, sanctifiedHealingEffect)
	core.NewItemEffect(ChampionsSealOfTheDawnHealing, sanctifiedHealingEffect)
	core.NewItemEffect(VanguardsSealOfTheDawnHealing, sanctifiedHealingEffect)
	core.NewItemEffect(CrusadersSealOfTheDawnHealing, sanctifiedHealingEffect)
	core.NewItemEffect(CommandersSealOfTheDawnHealing, sanctifiedHealingEffect)
	core.NewItemEffect(HighlordsSSealOfTheDawnHealing, sanctifiedHealingEffect)

	// https://www.wowhead.com/classic/item=236394/squires-seal-of-the-dawn
	core.NewItemEffect(SquiresSealOfTheDawnTanking, sanctifiedTankingEffect)
	core.NewItemEffect(KnightsSealOfTheDawnTanking, sanctifiedTankingEffect)
	core.NewItemEffect(TemplarsSealOfTheDawnTanking, sanctifiedTankingEffect)
	core.NewItemEffect(ChampionsSealOfTheDawnTanking, sanctifiedTankingEffect)
	core.NewItemEffect(VanguardsSealOfTheDawnTanking, sanctifiedTankingEffect)
	core.NewItemEffect(CrusadersSealOfTheDawnTanking, sanctifiedTankingEffect)
	core.NewItemEffect(CommandersSealOfTheDawnTanking, sanctifiedTankingEffect)
	core.NewItemEffect(HighlordsSSealOfTheDawnTanking, sanctifiedTankingEffect)

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////
	// https://www.wowhead.com/classic/item=236400/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshCastSpeed, func(agent core.Agent) {
		core.AtieshCastSpeedEffect(&agent.GetCharacter().Unit)
	})
	// https://www.wowhead.com/classic/item=236399/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshHealing, func(agent core.Agent) {
		core.AtieshHealingEffect(&agent.GetCharacter().Unit)
	})
	// https://www.wowhead.com/classic/item=236401/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshSpellCrit, func(agent core.Agent) {
		core.AtieshSpellCritEffect(&agent.GetCharacter().Unit)
	})
	// https://www.wowhead.com/classic/item=236398/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshSpellPower, func(agent core.Agent) {
		core.AtieshSpellPowerEffect(&agent.GetCharacter().Unit)
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Other
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=235868/bulwark-of-ire
	// Deal 100 Shadow damage to melee attackers.
	// Causes twice as much threat as damage dealt.
	core.NewItemEffect(BulwarkOfIre, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThornsDamage += 100

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: BulwarkOfIre},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 2,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 100, spell.OutcomeMagicHit)
			},
		})

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Splintered Shield",
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		}))
	})

	core.AddEffectsToTest = true
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your damage by 2% for each piece of Sanctified armor equipped.
func sanctifiedDamageEffect(agent core.Agent) {
	character := agent.GetCharacter()
	character.PseudoStats.DamageDealtMultiplier *= 1 + 0.02*float64(character.PseudoStats.SanctifiedBonus)
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your healing and shielding by 2% for each piece of Sanctified armor equipped.
func sanctifiedHealingEffect(agent core.Agent) {
	character := agent.GetCharacter()
	character.PseudoStats.HealingDealtMultiplier *= 1 + 0.02*float64(character.PseudoStats.SanctifiedBonus)
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your threat caused by 1% and health by 1 for each piece of Sanctified armor equipped.
func sanctifiedTankingEffect(agent core.Agent) {
	character := agent.GetCharacter()
	character.PseudoStats.ThreatMultiplier *= 1 + 0.01*float64(character.PseudoStats.SanctifiedBonus)
	character.AddStat(stats.Health, 1*float64(character.PseudoStats.SanctifiedBonus))
}
