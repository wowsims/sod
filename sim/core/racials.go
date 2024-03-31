package core

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceDwarf:
		character.AddStat(stats.FrostResistance, 10)
		character.PseudoStats.GunsSkill += 5

		actionID := ActionID{SpellID: 20594}

		statDep := character.NewDynamicMultiplyStat(stats.Armor, 1.1)
		stoneFormAura := character.NewTemporaryStatsAuraWrapped("Stoneform", actionID, stats.Stats{}, time.Second*8, func(aura *Aura) {
			aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			})
			aura.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			})
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				stoneFormAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
		})
	case proto.Race_RaceGnome:
		character.AddStat(stats.ArcaneResistance, 10)
		character.MultiplyStat(stats.Intellect, 1.05)
	case proto.Race_RaceHuman:
		character.MultiplyStat(stats.Spirit, 1.05)
		character.PseudoStats.MacesSkill += 5
		character.PseudoStats.TwoHandedMacesSkill += 5
		character.PseudoStats.SwordsSkill += 5
		character.PseudoStats.TwoHandedSwordsSkill += 5
	case proto.Race_RaceNightElf:
		character.AddStat(stats.NatureResistance, 10)
		character.AddStat(stats.Dodge, 1)
		// TODO: Shadowmeld?
	case proto.Race_RaceOrc:
		// Command (Pet damage +5%)
		for _, pet := range character.Pets {
			pet.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Blood Fury
		actionID := ActionID{SpellID: 20572}
		var bloodFuryAP float64
		bloodFuryAura := character.RegisterAura(Aura{
			Label:    "Blood Fury",
			ActionID: actionID,
			Duration: time.Second * 15,
			// Tooltip is misleading; ap bonus is base AP plus AP from current strength, does not include +attackpower on items/buffs
			OnGain: func(aura *Aura, sim *Simulation) {
				bloodFuryAP = (character.GetBaseStats()[stats.AttackPower] + (character.GetStat(stats.Strength) * 2)) * 0.25
				character.AddStatDynamic(sim, stats.AttackPower, bloodFuryAP)
			},

			OnExpire: func(aura *Aura, sim *Simulation) {
				character.AddStatDynamic(sim, stats.AttackPower, -bloodFuryAP)
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				DefaultCast: Cast{
					GCD: GCDDefault,
				},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				bloodFuryAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})

		// Axe specialization
		character.PseudoStats.AxesSkill += 5
		character.PseudoStats.TwoHandedAxesSkill += 5
	case proto.Race_RaceTauren:
		character.AddStat(stats.NatureResistance, 10)
		character.MultiplyStat(stats.Health, 1.05)
	case proto.Race_RaceTroll:
		character.PseudoStats.BowsSkill += 5
		character.PseudoStats.ThrownSkill += 5

		// Beast Slaying (+5% damage to beasts)
		character.Env.RegisterPostFinalizeEffect(func() {
			for _, t := range character.Env.Encounter.Targets {
				if t.MobType == proto.MobType_MobTypeBeast {
					for _, at := range character.AttackTables[t.UnitIndex] {
						at.DamageDealtMultiplier *= 1.05
						at.CritMultiplier *= 1.05
					}
				}
			}
		})

		// Berserking
		berserkingTimer := character.NewTimer()
		// Baseline cooldown
		makeBerserkingCooldown(character, 0, berserkingTimer)
		// Hard-coded percentage cooldown options
		makeBerserkingCooldown(character, .1, berserkingTimer)
		makeBerserkingCooldown(character, .15, berserkingTimer)
		makeBerserkingCooldown(character, .2, berserkingTimer)
		makeBerserkingCooldown(character, .25, berserkingTimer)
		makeBerserkingCooldown(character, .3, berserkingTimer)
	case proto.Race_RaceUndead:
		character.AddStat(stats.ShadowResistance, 10)
	}
}

// If customPercentage is 0, use the baseline Berserking calculations from health missing
// otherwise create a cooldown hard-coded to the custom percentage.
func makeBerserkingCooldown(character *Character, customPercentage float64, timer *Timer) {
	actionID := ActionID{SpellID: 26297, Tag: int32(customPercentage * 20)}

	label := "Berserking"
	if customPercentage != 0 {
		label = fmt.Sprintf("%s (%d)", label, int(customPercentage*100))
	}

	calcBerserkingPct := func() float64 {
		if customPercentage != 0 {
			return customPercentage
		}
		// from 10% at full health to 30% at 40% or less health
		switch hp := character.CurrentHealthPercent(); {
		case hp >= 1:
			return 0.1
		case hp <= 0.4:
			return 0.3
		default:
			return 0.1 + (1-hp)/3
		}
	}

	var berserkingAura *Aura
	var berserkingHaste float64
	if character.HasManaBar() {
		// Mana-using classes gain a flat % reduction in attack and cast speed
		berserkingAura = character.RegisterAura(Aura{
			Label:    label,
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				berserkingHaste = 1 / (1 - calcBerserkingPct())

				character.MultiplyCastSpeed(berserkingHaste)
				character.MultiplyAttackSpeed(sim, berserkingHaste)

				if sim.Log != nil {
					character.Log(sim, "Berserking increased attack and casting speed by %.2f%% (%.2f%% hp)", berserkingHaste*100-100, character.CurrentHealthPercent()*100)
				}
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1 / berserkingHaste)
				character.MultiplyAttackSpeed(sim, 1/berserkingHaste)
			},
		})
	} else {
		// Non-mana bar classes gain a flat % reduction in attack and cast speed
		berserkingAura = character.RegisterAura(Aura{
			Label:    label,
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				berserkingHaste = 1 + calcBerserkingPct()

				character.MultiplyAttackSpeed(sim, berserkingHaste)

				if sim.Log != nil {
					character.Log(sim, "Berserking increased attack speed by %.2f%% (%.2f%% hp)", berserkingHaste*100-100, character.CurrentHealthPercent()*100)
				}
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyAttackSpeed(sim, 1/berserkingHaste)
			},
		})
	}

	config := SpellConfig{
		ActionID: actionID,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    timer,
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			berserkingAura.Activate(sim)
		},
	}

	switch {
	case character.HasManaBar():
		config.ManaCost = ManaCostOptions{BaseCost: 0.07}
	case character.HasRageBar():
		config.RageCost = RageCostOptions{Cost: 5}
	case character.HasEnergyBar():
		config.EnergyCost = EnergyCostOptions{Cost: 10}
	}

	berserkingSpell := character.RegisterSpell(config)

	character.AddMajorCooldown(MajorCooldown{
		Spell: berserkingSpell,
		Type:  CooldownTypeDPS,
	})
}
