package warlock

import (
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

	// Cloak Runes
	warlock.applyDecimation()

	warlock.applyDemonicTactics()
	warlock.applyDemonicPact()
	warlock.applyGrimoireOfSynergy()
	warlock.applyShadowAndFlame()
	warlock.applyDemonicKnowledge()
	warlock.applyDanceOfTheWicked()
}

func (warlock *Warlock) applyVengeance() {
	if !warlock.HasRune(proto.WarlockRune_RuneHelmVengeance) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.WarlockRune_RuneHelmVengeance)}
	healthMetrics := warlock.NewHealthMetrics(actionID)
	var bonusHealth float64

	aura := warlock.RegisterAura(core.Aura{
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
			aura.Activate(sim)
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
	if !warlock.HasRune(proto.WarlockRune_RuneCloakDecimation) {
		return
	}

	affectedSpellCodes := []int32{SpellCode_WarlockShadowBolt, SpellCode_WarlockShadowCleave, SpellCode_WarlockIncinerate, SpellCode_WarlockSoulFire}

	decimationAura := warlock.RegisterAura(core.Aura{
		Label:    "Decimation",
		ActionID: core.ActionID{SpellID: 440873},
		Duration: time.Second * 10,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SoulFire {
				if spell != nil {
					spell.CastTimeMultiplier *= .6
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SoulFire {
				if spell != nil {
					spell.CastTimeMultiplier /= .6
				}
			}
		},
	})

	// Hidden trigger aura
	warlock.RegisterAura(core.Aura{
		Label:    "Decimation Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && sim.IsExecutePhase35() && slices.Contains(affectedSpellCodes, spell.SpellCode) {
				decimationAura.Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) InvocationRefresh(sim *core.Simulation, dot *core.Dot) {
	if dot.RemainingDuration(sim) < time.Second*6 {
		ticksLeft := dot.NumberOfTicks - dot.TickCount
		for i := int32(0); i < ticksLeft; i++ {
			dot.TickOnce(sim)
		}
	}
}

func (warlock *Warlock) EverlastingAfflictionRefresh(sim *core.Simulation, target *core.Unit) {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsEverlastingAffliction) {
		return
	}

	for _, spell := range warlock.Corruption {
		if spell.Dot(target).IsActive() {
			spell.Dot(target).Rollover(sim)
		}
	}
}

func (warlock *Warlock) applyDanceOfTheWicked() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsDanceOfTheWicked) {
		return
	}

	actionId := core.ActionID{SpellID: 412800}
	dodgeModifier := warlock.NewDynamicStatDependency(stats.SpellCrit, stats.Dodge, 1)

	dotwAura := warlock.GetOrRegisterAura(core.Aura{
		Label:    "Dance of the Wicked Proc",
		ActionID: actionId,
		Duration: 15 * time.Second,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.EnableDynamicStatDep(sim, dodgeModifier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.DisableDynamicStatDep(sim, dodgeModifier)
		},
	})

	manaMetric := warlock.NewManaMetrics(actionId)

	var petMetric *core.ResourceMetrics
	if warlock.Pet != nil {
		petMetric = warlock.Pet.NewManaMetrics(actionId)
	}

	handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !spell.ProcMask.Matches(core.ProcMaskDirect) {
			return
		}

		if !result.DidCrit() {
			return
		}

		dotwAura.Activate(sim)

		warlock.AddMana(sim, warlock.MaxMana()*0.02, manaMetric)
		if warlock.Pet != nil && warlock.Pet.IsActive() {
			warlock.Pet.AddMana(sim, warlock.Pet.MaxMana()*0.02, petMetric)
		}
	}

	core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
		Label:                 "Dance of the Wicked",
		OnSpellHitDealt:       handler,
		OnPeriodicDamageDealt: handler,
	}))
}

func (warlock *Warlock) applyDemonicKnowledge() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsDemonicKnowledge) || warlock.Pet == nil {
		return
	}

	wp := warlock.Pet
	oldPetEnable := wp.OnPetEnable
	wp.OnPetEnable = func(sim *core.Simulation) {
		if oldPetEnable != nil {
			oldPetEnable(sim)
		}
		warlock.DemonicKnowledgeAura.Activate(sim)
	}

	oldPetDisable := wp.OnPetDisable
	wp.OnPetDisable = func(sim *core.Simulation) {
		if oldPetDisable != nil {
			oldPetDisable(sim)
		}
		warlock.DemonicKnowledgeAura.Deactivate(sim)
	}

	warlock.DemonicKnowledgeAura = warlock.GetOrRegisterAura(core.Aura{
		Label:    "Demonic Knowledge",
		ActionID: core.ActionID{SpellID: int32(proto.WarlockRune_RuneBootsDemonicKnowledge)},
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.demonicKnowledgeSp = (warlock.Pet.GetStat(stats.Stamina) + warlock.Pet.GetStat(stats.Intellect)) * 0.12
			warlock.AddStatDynamic(sim, stats.SpellPower, warlock.demonicKnowledgeSp)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.AddStatDynamic(sim, stats.SpellPower, -warlock.demonicKnowledgeSp)
		},
	})
}

func (warlock *Warlock) applyGrimoireOfSynergy() {
	if !warlock.HasRune(proto.WarlockRune_RuneBeltGrimoireOfSynergy) || warlock.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 426303}
	dmgMod := 1.25
	procChance := 0.05

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

	warlockProcAura := warlock.GetOrRegisterAura(procAuraConfig)
	petProcAura := warlock.Pet.GetOrRegisterAura(procAuraConfig)

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

	core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
		Label:                 "Grimoire of Synergy",
		OnSpellHitDealt:       handlerFunc(petProcAura),
		OnPeriodicDamageDealt: handlerFunc(petProcAura),
	}))

	core.MakePermanent(warlock.Pet.GetOrRegisterAura(core.Aura{
		Label:                 "Grimoire of Synergy",
		OnSpellHitDealt:       handlerFunc(warlockProcAura),
		OnPeriodicDamageDealt: handlerFunc(warlockProcAura),
	}))
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
// Increases the amount drained by your Drain Life and Drain Soul spells by an additional 6% for each of your Warlock Shadow effects afflicting the target,
// up to a maximum of 18% additional effect. When Drain Soul is cast on a target below 20% health, it instead gains 100% per effect, up to a maximum of 300%.
const SoulSiphonDoTMultiplier = 0.06
const SoulSiphonDoTMultiplierExecute = 0.50
const SoulSiphonDoTMultiplierMax = 0.18
const SoulSiphonDoTMultiplierMaxExecute = 1.50

func (warlock *Warlock) applyDemonicTactics() {
	if !warlock.HasRune(proto.WarlockRune_RuneChestDemonicTactics) {
		return
	}

	warlock.AddStat(stats.MeleeCrit, 10*core.CritRatingPerCritChance)
	warlock.AddStat(stats.SpellCrit, 10*core.SpellCritRatingPerCritChance)

	if warlock.Pet != nil {
		pet := warlock.Pet.GetPet()
		pet.AddStat(stats.MeleeCrit, 10*core.CritRatingPerCritChance)
		pet.AddStat(stats.SpellCrit, 10*core.SpellCritRatingPerCritChance)
	}
}

func (warlock *Warlock) getHighestSP() float64 {
	return warlock.GetStat(stats.SpellPower) + warlock.GetStat(stats.SpellDamage) + max(warlock.GetStat(stats.FirePower), warlock.GetStat(stats.ShadowPower))
}

func (warlock *Warlock) applyDemonicPact() {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsDemonicPact) {
		return
	}

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

	warlock.Pet.RegisterAura(core.Aura{
		Label:    "Demonic Pact Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PreviousTime = 0
			aura.Activate(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() || !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)

			currentSP := warlock.getHighestSP()

			// Remove DP bonus from SP bonus if active
			if demonicPactAuras.Get(&warlock.Unit).IsActive() {
				currentSP -= demonicPactAuras.Get(&warlock.Unit).ExclusiveEffects[0].Priority
			}
			spBonus := max(math.Round(currentSP*0.1), math.Round(float64(warlock.Level)/2))
			for _, dpAura := range demonicPactAuras {
				if dpAura != nil {
					// Force expire/gain because of new sp bonus
					dpAura.Deactivate(sim)

					dpAura.ExclusiveEffects[0].SetPriority(sim, spBonus)
					dpAura.Activate(sim)
				}
			}
		},
	})
}
