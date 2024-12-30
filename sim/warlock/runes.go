package warlock

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) ApplyRunes() {
	// Helm runes
	warlock.applyVengeance()
	warlock.applyBackdraft()

	// Shoulders
	warlock.applyShoulderRuneEffect()

	// Cloak Runes
	warlock.applyDecimation()
	warlock.registerInfernalArmorCD()

	// Chest Runes
	warlock.applyDemonicTactics()

	// Bracer Runes
	warlock.registerIncinerateSpell()
	warlock.registerUnstableAfflictionSpell()
	warlock.registerImmolationAuraSpell()

	// Glove Runes
	warlock.registerHauntSpell()
	warlock.registerChaosBoltSpell()
	warlock.registerMetamorphosisSpell()
	warlock.registerShadowCleaveSpell()

	// Belt Runes
	warlock.applyGrimoireOfSynergy()
	warlock.applyShadowAndFlame()

	// Pants Runes
	warlock.applyEverlastingAffliction()
	warlock.applyDemonicPact()
	warlock.registerDemonicGraceSpell()

	// Boots Runes
	warlock.applyDemonicKnowledge()
	warlock.applyDanceOfTheWicked()
	warlock.registerShadowflameSpell()
	warlock.applyMarkOfChaos()
}

func (warlock *Warlock) applyShoulderRuneEffect() {
	if warlock.Equipment.Shoulders().Rune == int32(proto.WarlockRune_WarlockRuneNone) {
		return
	}

	switch warlock.Equipment.Shoulders().Rune {
	// Damage
	case int32(proto.WarlockRune_RuneShouldersTransfusionist):
		warlock.applyT1Damage2PBonus()
	case int32(proto.WarlockRune_RuneShouldersRefinedWarlock):
		warlock.applyT1Damage4PBonus()
	case int32(proto.WarlockRune_RuneShouldersDecimator):
		warlock.applyT1Damage6PBonus()
	case int32(proto.WarlockRune_RuneShouldersRotbringer):
		warlock.applyT2Damage2PBonus()
	case int32(proto.WarlockRune_RuneShouldersMalevolent):
		warlock.applyT2Damage4PBonus()
	case int32(proto.WarlockRune_RuneShouldersShadowmancer):
		warlock.applyT2Damage6PBonus()
	case int32(proto.WarlockRune_RuneShouldersInfernalShepherd):
		warlock.applyZGDemonology3PBonus()
	case int32(proto.WarlockRune_RuneShouldersDemonlord):
		warlock.applyZGDemonology5PBonus()
	case int32(proto.WarlockRune_RuneShouldersChaosHarbinger):
		warlock.applyTAQDamage2PBonus()
	case int32(proto.WarlockRune_RuneShouldersArsonist):
		warlock.applyTAQDamage4PBonus()

	// Tank
	case int32(proto.WarlockRune_RuneShouldersDemonicExorcist):
		warlock.applyT1Tank2PBonus()
	case int32(proto.WarlockRune_RuneShouldersPained):
		warlock.applyT1Tank4PBonus()
	case int32(proto.WarlockRune_RuneShouldersFlamewraith):
		warlock.applyT1Tank6PBonus()
	case int32(proto.WarlockRune_RuneShouldersFleshfeaster):
		warlock.applyT2Tank2PBonus()
	case int32(proto.WarlockRune_RuneShouldersAbyssal):
		warlock.applyT2Tank4PBonus()
	case int32(proto.WarlockRune_RuneShouldersVoidborne):
		warlock.applyT2Tank6PBonus()
	case int32(proto.WarlockRune_RuneShouldersUmbralBlade):
		warlock.applyTAQTank2PBonus()
	case int32(proto.WarlockRune_RuneShouldersRitualist):
		warlock.applyTAQTank4PBonus()
	case int32(proto.WarlockRune_RuneShouldersPainSpreader):
		warlock.applyRAQTank3PBonus()
	}
}

func (warlock *Warlock) applyVengeance() {
	if !warlock.HasRune(proto.WarlockRune_RuneHelmVengeance) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.WarlockRune_RuneHelmVengeance)}
	healthMetrics := warlock.NewHealthMetrics(actionID)
	var bonusHealth float64

	warlock.VengeanceAura = warlock.RegisterAura(core.Aura{
		Label:    "Vengeance",
		ActionID: actionID,
		Duration: time.Second * 20,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = warlock.MaxHealth() * 0.30
			warlock.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			warlock.GainHealth(sim, bonusHealth, healthMetrics)

		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
			healthDiff := warlock.CurrentHealth() - warlock.MaxHealth()
			if healthDiff > 0 {
				warlock.RemoveHealth(sim, healthDiff)
			}
		},
	})

	spell := warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.VengeanceAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentHealthPercent() < 0.5
		},
	})
}

func (warlock *Warlock) applyBackdraft() {
	if !warlock.HasRune(proto.WarlockRune_RuneHelmBackdraft) {
		return
	}

	warlock.BackdraftAura = warlock.RegisterAura(core.Aura{
		Label:    "Backdraft",
		ActionID: core.ActionID{SpellID: 427714},
		Duration: time.Second * 15,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.MultiplyCastSpeed(1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.MultiplyCastSpeed(1 / 1.3)
		},
	})
}

func (warlock *Warlock) applyDecimation() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsDecimation) {
		return
	}

	affectedSpellCodes := []int32{SpellCode_WarlockShadowBolt, SpellCode_WarlockShadowCleave, SpellCode_WarlockIncinerate, SpellCode_WarlockSoulFire}

	warlock.DecimationAura = warlock.RegisterAura(core.Aura{
		Label:    "Decimation",
		ActionID: core.ActionID{SpellID: 440873},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SoulFire {
				spell.CastTimeMultiplier *= .6
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SoulFire {
				spell.CastTimeMultiplier /= .6
			}
		},
	})

	// Hidden trigger aura
	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: "Decimation Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && sim.IsExecutePhase35() && slices.Contains(affectedSpellCodes, spell.SpellCode) {
				warlock.DecimationAura.Activate(sim)
			}
		},
	}))
}

func (warlock *Warlock) applyMarkOfChaos() {
	if !warlock.HasRune(proto.WarlockRune_RuneCloakMarkOfChaos) {
		return
	}

	warlock.MarkOfChaosAuras = warlock.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
		return core.MarkOfChaosDebuffAura(target)
	})
}

func (warlock *Warlock) applyMarkOfChaosDebuff(sim *core.Simulation, target *core.Unit, duration time.Duration) {
	aura := warlock.MarkOfChaosAuras.Get(target)
	// Only expire if not set as a permanent raid debuff.
	if !aura.IsPermanent() {
		aura.Duration = duration
		aura.UpdateExpires(sim, sim.CurrentTime+duration)
	}
	aura.Activate(sim)
}

func (warlock *Warlock) InvocationRefresh(sim *core.Simulation, dot *core.Dot) {
	if dot.RemainingDuration(sim) < time.Second*6 {
		ticksLeft := dot.NumberOfTicks - dot.TickCount
		for i := int32(0); i < ticksLeft; i++ {
			dot.TickOnce(sim)
		}
	}
}

func (warlock *Warlock) applyEverlastingAffliction() {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsEverlastingAffliction) {
		return
	}

	affectedSpellCodes := []int32{SpellCode_WarlockDrainLife, SpellCode_WarlockDrainSoul, SpellCode_WarlockShadowBolt, SpellCode_WarlockShadowCleave, SpellCode_WarlockSearingPain, SpellCode_WarlockIncinerate, SpellCode_WarlockHaunt}
	warlock.RegisterAura(core.Aura{
		Label:    "Everlasting Affliction Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && slices.Contains(affectedSpellCodes, spell.SpellCode) {
				for _, spell := range warlock.Corruption {
					if spell.Dot(result.Target).IsActive() {
						spell.Dot(result.Target).Rollover(sim)
					}
				}
			}
		},
	})
}

func (warlock *Warlock) applyDanceOfTheWicked() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsDanceOfTheWicked) {
		return
	}

	actionId := core.ActionID{SpellID: 412800}
	lastCritSnapshot := 0.0

	// DoTW snapshot your current crit each time it procs so we want to add the delta between the last and current snapshot
	dotwAura := warlock.GetOrRegisterAura(core.Aura{
		ActionID: actionId,
		Label:    "Dance of the Wicked Proc",
		Duration: 15 * time.Second,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			lastCritSnapshot = 0
		},
		OnRefresh: func(aura *core.Aura, sim *core.Simulation) {
			newCritSnapshot := warlock.GetStat(stats.SpellCrit)
			warlock.AddStatDynamic(sim, stats.Dodge, newCritSnapshot-lastCritSnapshot)
			fmt.Println(sim.CurrentTime, "Current Crit:", warlock.GetStat(stats.SpellCrit), "Old Crit:", lastCritSnapshot, "Delta:", warlock.GetStat(stats.SpellCrit)-lastCritSnapshot)
			lastCritSnapshot = newCritSnapshot
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.Dodge, -lastCritSnapshot)
		},
	})

	manaMetric := warlock.NewManaMetrics(actionId)
	for _, pet := range warlock.BasePets {
		pet.DanceOfTheWickedManaMetrics = pet.NewManaMetrics(actionId)
	}

	handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if spell.ProcMask.Matches(core.ProcMaskDirect) && result.DidCrit() {
			dotwAura.Activate(sim)
			warlock.AddMana(sim, warlock.MaxMana()*0.02, manaMetric)

			if warlock.ActivePet != nil {
				warlock.ActivePet.AddMana(sim, warlock.ActivePet.MaxMana()*0.02, warlock.ActivePet.DanceOfTheWickedManaMetrics)
			}
		}
	}

	core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
		Label:                 "Dance of the Wicked",
		OnSpellHitDealt:       handler,
		OnPeriodicDamageDealt: handler,
	}))
}

func (warlock *Warlock) applyDemonicKnowledge() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsDemonicKnowledge) {
		return
	}

	for _, pet := range warlock.BasePets {
		oldOnPetEnable := pet.OnPetEnable
		pet.OnPetEnable = func(sim *core.Simulation) {
			oldOnPetEnable(sim)
			warlock.DemonicKnowledgeAura.Activate(sim)
		}

		oldOnPetDisable := pet.OnPetDisable
		pet.OnPetDisable = func(sim *core.Simulation, isSacrifice bool) {
			oldOnPetDisable(sim, isSacrifice)
			warlock.DemonicKnowledgeAura.Deactivate(sim)
		}
	}

	warlock.DemonicKnowledgeAura = warlock.GetOrRegisterAura(core.Aura{
		Label:    "Demonic Knowledge",
		ActionID: core.ActionID{SpellID: int32(proto.WarlockRune_RuneBootsDemonicKnowledge)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			warlock.demonicKnowledgeSp = 0
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.demonicKnowledgeSp = (warlock.ActivePet.GetStat(stats.Stamina) + warlock.ActivePet.GetStat(stats.Intellect)) * .03
			warlock.AddStatDynamic(sim, stats.SpellPower, warlock.demonicKnowledgeSp)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellPower, -warlock.demonicKnowledgeSp)
			warlock.demonicKnowledgeSp = 0
		},
	})
}

func (warlock *Warlock) applyGrimoireOfSynergy() {
	if !warlock.HasRune(proto.WarlockRune_RuneBeltGrimoireOfSynergy) {
		return
	}

	actionID := core.ActionID{SpellID: 426303}
	dmgMod := 1.10
	procChance := 0.10

	procAuraConfig := core.Aura{
		Label:    "Grimoire of Synergy Proc",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= dmgMod
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= dmgMod
		},
	}

	handlerFunc := func(procAura *core.Aura) func(*core.Aura, *core.Simulation, *core.Spell, *core.SpellResult) {
		return func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskDirect) {
				return
			}

			if sim.RandomFloat("Grimoire of Synergy") > procChance {
				return
			}

			procAura.Activate(sim)
		}
	}
	warlockProcAura := warlock.GetOrRegisterAura(procAuraConfig)
	for _, pet := range warlock.BasePets {
		petProcAura := pet.GetOrRegisterAura(procAuraConfig)
		core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
			Label:                 fmt.Sprintf("Grimoire of Synergy Trigger (%s)", pet.Name),
			OnSpellHitDealt:       handlerFunc(petProcAura),
			OnPeriodicDamageDealt: handlerFunc(petProcAura),
		}))
		core.MakePermanent(pet.GetOrRegisterAura(core.Aura{
			Label:                 "Grimoire of Synergy Trigger",
			OnSpellHitDealt:       handlerFunc(warlockProcAura),
			OnPeriodicDamageDealt: handlerFunc(warlockProcAura),
		}))
	}
}

func (warlock *Warlock) applyShadowAndFlame() {
	if !warlock.HasRune(proto.WarlockRune_RuneBeltShadowAndFlame) {
		return
	}

	procAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Shadow and Flame proc",
		ActionID: core.ActionID{SpellID: 426311},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.10
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.10
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.10
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.10
		},
	})

	procHandler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !spell.SpellSchool.Matches(core.SpellSchoolFire | core.SpellSchoolShadow) {
			return
		}

		if !result.DidCrit() {
			return
		}

		procAura.Activate(sim)
	}

	core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
		Label:                 "Shadow and Flame",
		OnSpellHitDealt:       procHandler,
		OnPeriodicDamageDealt: procHandler,
	}))
}

// https://www.wowhead.com/classic/spell=403511/soul-siphon
// Causes your Drain Soul to to deal damage 3 times faster and increases the amount drained by your Drain Life and Drain Soul spells by an additional
// 6% for each of your Warlock Shadow effects afflicting the target, up to a maximum of 18% additional effect.
// When Drain Soul is cast on a target below 20% health, it instead gains 50% per effect, up to a maximum of 150%.
const SoulSiphonDoTMultiplier = 0.06
const SoulSiphonDoTMultiplierExecute = 1.00
const SoulSiphonDoTMultiplierMax = 0.18
const SoulSiphonDoTMultiplierMaxExecute = 3.00

func (warlock *Warlock) calcSoulSiphonMultiplier(target *core.Unit, executeBonus bool) float64 {
	multiplier := 1.0
	perDoTMultiplier := core.TernaryFloat64(executeBonus, SoulSiphonDoTMultiplierExecute, SoulSiphonDoTMultiplier)
	maxMultiplier := 1 + core.TernaryFloat64(executeBonus, SoulSiphonDoTMultiplierMaxExecute, SoulSiphonDoTMultiplierMax)

	for _, spell := range warlock.Corruption {
		if spell.Dot(target).IsActive() {
			multiplier += perDoTMultiplier
			break
		}
	}

	for _, spell := range warlock.CurseOfAgony {
		if spell.Dot(target).IsActive() {
			multiplier += perDoTMultiplier
			break
		}
	}

	if warlock.CurseOfDoom != nil && warlock.CurseOfDoom.Dot(target).IsActive() {
		multiplier += perDoTMultiplier
	}

	for _, spell := range warlock.SiphonLife {
		if spell.Dot(target).IsActive() {
			multiplier += perDoTMultiplier
			break
		}
	}

	if warlock.UnstableAffliction != nil && warlock.UnstableAffliction.Dot(target).IsActive() {
		multiplier += perDoTMultiplier
	}

	if warlock.Shadowflame != nil && warlock.Shadowflame.Dot(target).IsActive() {
		multiplier += perDoTMultiplier
	}

	if warlock.Haunt != nil && warlock.HauntDebuffAuras.Get(target).IsActive() {
		multiplier += perDoTMultiplier
	}

	return min(multiplier, maxMultiplier)
}

// Increases the melee and spell critical strike chance of you and your pet by 10%.
func (warlock *Warlock) applyDemonicTactics() {
	if !warlock.HasRune(proto.WarlockRune_RuneChestDemonicTactics) {
		return
	}

	warlock.AddStat(stats.MeleeCrit, 10*core.CritRatingPerCritChance)
	warlock.AddStat(stats.SpellCrit, 10*core.SpellCritRatingPerCritChance)
}

func (warlock *Warlock) getHighestSP() float64 {
	return warlock.GetStat(stats.SpellPower) + warlock.GetStat(stats.SpellDamage) + max(warlock.GetStat(stats.FirePower), warlock.GetStat(stats.ShadowPower))
}

func (warlock *Warlock) applyDemonicPact() {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsDemonicPact) {
		return
	}

	warlock.PseudoStats.SchoolDamageDealtMultiplier.MultiplyMagicSchools(1.10)

	if warlock.Options.Summon == proto.WarlockOptions_NoSummon {
		return
	}

	icd := core.Cooldown{
		Timer:    warlock.NewTimer(),
		Duration: 1 * time.Second,
	}

	spellPower := max(warlock.getHighestSP()*0.1, float64(warlock.Level)/2.0)
	demonicPactAuras := warlock.NewRaidAuraArray(func(u *core.Unit) *core.Aura {
		return core.DemonicPactAura(u, spellPower, core.CharacterBuildPhaseNone)
	})

	dpTriggerConfig := core.Aura{
		Label:    "Demonic Pact Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() || !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)

			lastBonus := 0.0
			currentSP := warlock.getHighestSP()
			warlockAura := demonicPactAuras.Get(&warlock.Unit)

			// Remove DP bonus from SP bonus if active
			if warlockAura.IsActive() {
				lastBonus = warlockAura.ExclusiveEffects[0].Priority
			}
			newSPBonus := max(math.Round(0.10*(currentSP-lastBonus)), math.Round(float64(warlock.Level)/2))

			if warlockAura.RemainingDuration(sim) < 10*time.Second || newSPBonus >= lastBonus {
				for _, dpAura := range demonicPactAuras {
					if dpAura != nil {
						// Force expire/gain because of new sp bonus
						dpAura.Deactivate(sim)

						dpAura.ExclusiveEffects[0].SetPriority(sim, newSPBonus)
						dpAura.Activate(sim)
						dpAura.SetStacks(sim, int32(newSPBonus))
					}
				}
			}
		},
	}

	for _, pet := range warlock.BasePets {
		pet.RegisterAura(dpTriggerConfig)
	}
}
