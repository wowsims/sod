package core

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceDwarf:
		character.PseudoStats.ReducedFrostHitTakenChance += 0.02

		// Gun specialization (+1% ranged crit when using a gun).
		if character.Ranged().RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeGun {
			character.AddBonusRangedCritRating(1 * CritRatingPerCritChance)
		}

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
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				stoneFormAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceGnome:
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.02
		character.MultiplyStat(stats.Intellect, 1.05)
	case proto.Race_RaceHuman:
		character.MultiplyStat(stats.Spirit, 1.03)
		character.PseudoStats.MacesSkill += 5
		character.PseudoStats.TwoHandedMacesSkill += 5
		character.PseudoStats.SwordsSkill += 5
		character.PseudoStats.TwoHandedSwordsSkill += 5
	case proto.Race_RaceNightElf:
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.PseudoStats.ReducedPhysicalHitTakenChance += 0.02
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
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.AddStat(stats.Health, character.GetBaseStats()[stats.Health]*0.05)
	case proto.Race_RaceTroll:
		// Bow specialization (+1% ranged crit when using a bow).
		if character.Ranged().RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeBow {
			character.AddBonusRangedCritRating(1 * CritRatingPerCritChance)
		}

		// Beast Slaying (+5% damage to beasts)
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Berserking
		actionID := ActionID{SpellID: 26297}

		var berserkingAura *Aura
		var berserkingPct float64
		if character.HasManaBar() {
			// Mana-using classes gain a flat % reduction in attack and cast speed
			berserkingAura = character.RegisterAura(Aura{
				Label:    "Berserking (Troll)",
				ActionID: actionID,
				Duration: time.Second * 10,
				OnGain: func(aura *Aura, sim *Simulation) {
					healthPctMissing := 1 - character.CurrentHealthPercent()
					// 10% base + 1/3 of missing health percentage up to a max of 30%
					berserkingPct = math.Min(.1+(healthPctMissing/3), .3)

					character.FlatIncreaseCastSpeed(berserkingPct)
					character.FlatIncreaseAttackSpeed(sim, berserkingPct)

					if sim.Log != nil {
						character.Log(sim, "Berserking increased attack speed by %.2f%% (%.2f%% hp)", berserkingPct*100, character.CurrentHealthPercent()*100.0)
					}
				},
				OnExpire: func(aura *Aura, sim *Simulation) {
					character.FlatIncreaseCastSpeed(-1 * berserkingPct)
					character.FlatIncreaseAttackSpeed(sim, -1*berserkingPct)
				},
			})
		} else {
			// Non-mana bar classes gain a flat % reduction in attack and cast speed
			berserkingAura = character.RegisterAura(Aura{
				Label:    "Berserking (Troll)",
				ActionID: actionID,
				Duration: time.Second * 10,
				OnGain: func(aura *Aura, sim *Simulation) {
					healthPctMissing := 1 - character.CurrentHealthPercent()
					// 10% base + 1/3 of missing health percentage up to a max of 30%
					berserkingPct := math.Min(.1+(healthPctMissing/3), .3)

					character.MultiplyCastSpeed(1 + berserkingPct)
					character.MultiplyAttackSpeed(sim, 1+berserkingPct)

					if sim.Log != nil {
						character.Log(sim, "Berserking increased attack speed by %.2f%% (%.2f%% hp)", berserkingPct*100, character.CurrentHealthPercent()*100.0)
					}
				},
				OnExpire: func(aura *Aura, sim *Simulation) {
					character.MultiplyCastSpeed(1 / (1 + berserkingPct))
					character.MultiplyAttackSpeed(sim, 1/(1+berserkingPct))
				},
			})
		}

		berserkingSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				berserkingAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: berserkingSpell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceUndead:
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02
	}
}
