package item_effects

import (
	"math"

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
	AspirantsSealOfTheDawnDamage  = 236354
	InitiatesSealOfTheDawnDamage  = 236355
	SquiresSealOfTheDawnDamage    = 236356
	KnightsSealOfTheDawnDamage    = 236357
	TemplarsSealOfTheDawnDamage   = 236358
	ChampionsSealOfTheDawnDamage  = 236360
	VanguardsSealOfTheDawnDamage  = 236361
	CrusadersSealOfTheDawnDamage  = 236362
	CommandersSealOfTheDawnDamage = 236363
	HighlordsSSealOfTheDawnDamage = 236364

	AspirantsSealOfTheDawnHealing  = 236385
	InitiatesSealOfTheDawnHealing  = 236384
	SquiresSealOfTheDawnHealing    = 236383
	KnightsSealOfTheDawnHealing    = 236382
	TemplarsSealOfTheDawnHealing   = 236380
	ChampionsSealOfTheDawnHealing  = 236379
	VanguardsSealOfTheDawnHealing  = 236378
	CrusadersSealOfTheDawnHealing  = 236376
	CommandersSealOfTheDawnHealing = 236375
	HighlordsSSealOfTheDawnHealing = 236374

	AspirantsSealOfTheDawnTanking  = 236396
	InitiatesSealOfTheDawnTanking  = 236395
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
	core.NewItemEffect(AspirantsSealOfTheDawnDamage, sanctifiedDamageEffect(1219539, 0.83))
	core.NewItemEffect(InitiatesSealOfTheDawnDamage, sanctifiedDamageEffect(1223348, 2.92))
	core.NewItemEffect(SquiresSealOfTheDawnDamage, sanctifiedDamageEffect(1223349, 4.17))
	core.NewItemEffect(KnightsSealOfTheDawnDamage, sanctifiedDamageEffect(1223350, 6.67))
	core.NewItemEffect(TemplarsSealOfTheDawnDamage, sanctifiedDamageEffect(1223351, 8.33))
	core.NewItemEffect(ChampionsSealOfTheDawnDamage, sanctifiedDamageEffect(1223352, 12.08))
	core.NewItemEffect(VanguardsSealOfTheDawnDamage, sanctifiedDamageEffect(1223353, 14.17))
	core.NewItemEffect(CrusadersSealOfTheDawnDamage, sanctifiedDamageEffect(1223354, 18.75))
	core.NewItemEffect(CommandersSealOfTheDawnDamage, sanctifiedDamageEffect(1223355, 21.67))
	core.NewItemEffect(HighlordsSSealOfTheDawnDamage, sanctifiedDamageEffect(1223357, 25.0))

	// https://www.wowhead.com/classic/item=236383/squires-seal-of-the-dawn
	core.NewItemEffect(AspirantsSealOfTheDawnHealing, sanctifiedHealingEffect(1219548, 0.83))
	core.NewItemEffect(InitiatesSealOfTheDawnHealing, sanctifiedHealingEffect(1223379, 2.92))
	core.NewItemEffect(SquiresSealOfTheDawnHealing, sanctifiedHealingEffect(1223380, 4.17))
	core.NewItemEffect(KnightsSealOfTheDawnHealing, sanctifiedHealingEffect(1223381, 6.67))
	core.NewItemEffect(TemplarsSealOfTheDawnHealing, sanctifiedHealingEffect(1223382, 8.33))
	core.NewItemEffect(ChampionsSealOfTheDawnHealing, sanctifiedHealingEffect(1223383, 12.08))
	core.NewItemEffect(VanguardsSealOfTheDawnHealing, sanctifiedHealingEffect(1223384, 14.17))
	core.NewItemEffect(CrusadersSealOfTheDawnHealing, sanctifiedHealingEffect(1223385, 18.75))
	core.NewItemEffect(CommandersSealOfTheDawnHealing, sanctifiedHealingEffect(1223386, 21.67))
	core.NewItemEffect(HighlordsSSealOfTheDawnHealing, sanctifiedHealingEffect(1223387, 25.0))

	// https://www.wowhead.com/classic/item=236394/squires-seal-of-the-dawn
	core.NewItemEffect(AspirantsSealOfTheDawnTanking, sanctifiedTankingEffect(1220514, 2.08, 0.83))
	core.NewItemEffect(InitiatesSealOfTheDawnTanking, sanctifiedTankingEffect(1223367, 2.92, 2.92))
	core.NewItemEffect(SquiresSealOfTheDawnTanking, sanctifiedTankingEffect(1223368, 3.33, 4.17))
	core.NewItemEffect(KnightsSealOfTheDawnTanking, sanctifiedTankingEffect(1223370, 3.75, 6.67))
	core.NewItemEffect(TemplarsSealOfTheDawnTanking, sanctifiedTankingEffect(1223371, 4.17, 8.33))
	core.NewItemEffect(ChampionsSealOfTheDawnTanking, sanctifiedTankingEffect(1223372, 5.0, 12.08))
	core.NewItemEffect(VanguardsSealOfTheDawnTanking, sanctifiedTankingEffect(1223373, 5.42, 14.17))
	core.NewItemEffect(CrusadersSealOfTheDawnTanking, sanctifiedTankingEffect(1223374, 6.25, 18.75))
	core.NewItemEffect(CommandersSealOfTheDawnTanking, sanctifiedTankingEffect(1223375, 6.67, 21.67))
	core.NewItemEffect(HighlordsSSealOfTheDawnTanking, sanctifiedTankingEffect(1223376, 7.08, 25.0))

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
// Increasing your damage by X% and your health by X% for each piece of Sanctified armor equipped.
func sanctifiedDamageEffect(spellID int32, percentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.PseudoStats.SanctifiedBonus == 0 {
			return
		}

		sanctifiedBonus := math.Min(12, float64(character.PseudoStats.SanctifiedBonus))
		multiplier := 1.0 + percentIncrease/100.0*sanctifiedBonus
		healthDep := character.NewDynamicMultiplyStat(stats.Health, multiplier)

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:      "Seal of the Dawn (Damage)",
			ActionID:   core.ActionID{SpellID: spellID},
			BuildPhase: core.CharacterBuildPhaseGear,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDep)
				} else {
					character.EnableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.DamageDealtMultiplier *= multiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDep)
				} else {
					character.DisableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.DamageDealtMultiplier /= multiplier
			},
		}))
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your healing and shielding by X% and your health by X% for each piece of Sanctified armor equipped.
func sanctifiedHealingEffect(spellID int32, percentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.PseudoStats.SanctifiedBonus == 0 {
			return
		}

		sanctifiedBonus := math.Min(12, float64(character.PseudoStats.SanctifiedBonus))
		multiplier := 1.0 + percentIncrease/100.0*sanctifiedBonus
		healthDep := character.NewDynamicMultiplyStat(stats.Health, multiplier)

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:      "Seal of the Dawn (Healing)",
			ActionID:   core.ActionID{SpellID: spellID},
			BuildPhase: core.CharacterBuildPhaseGear,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDep)
				} else {
					character.EnableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.HealingDealtMultiplier *= multiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDep)
				} else {
					character.DisableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.HealingDealtMultiplier /= multiplier
			},
		}))
	}
}

// Equip: Unlocks your potential while inside Naxxramas.
// Increasing your threat caused by X%, your damage by Y%, and your health by Y% for each piece of Sanctified armor equipped.
func sanctifiedTankingEffect(spellID int32, threatPercentIncrease float64, damageHealthPercentIncrease float64) core.ApplyEffect {
	return func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.PseudoStats.SanctifiedBonus == 0 {
			return
		}

		sanctifiedBonus := math.Min(12, float64(character.PseudoStats.SanctifiedBonus))
		damageHealthMultiplier := 1.0 + damageHealthPercentIncrease/100.0*sanctifiedBonus
		threatMultiplier := 1.0 + threatPercentIncrease/100.0*sanctifiedBonus
		healthDep := character.NewDynamicMultiplyStat(stats.Health, damageHealthMultiplier)

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:      "Seal of the Dawn (Tanking)",
			ActionID:   core.ActionID{SpellID: spellID},
			BuildPhase: core.CharacterBuildPhaseGear,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.EnableDynamicStatDep(healthDep)
				} else {
					character.EnableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.ThreatMultiplier *= threatMultiplier
				character.PseudoStats.DamageDealtMultiplier *= damageHealthMultiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
					aura.Unit.StatDependencyManager.DisableDynamicStatDep(healthDep)
				} else {
					character.DisableDynamicStatDep(sim, healthDep)
				}

				character.PseudoStats.ThreatMultiplier /= threatMultiplier
				character.PseudoStats.DamageDealtMultiplier /= damageHealthMultiplier
			},
		}))
	}
}
