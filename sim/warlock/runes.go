package warlock

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) ApplyRunes() {
	warlock.applyDemonicTactics()
	warlock.applyDemonicPact()
	warlock.applyGrimoireOfSynergy()
	warlock.applyShadowAndFlame()
	warlock.applyDemonicKnowledge()
	warlock.applyDanceOfTheWicked()
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

	if warlock.Corruption.Dot(target).IsActive() {
		warlock.Corruption.Dot(target).Rollover(sim)
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

	core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
		Label: "Dance of the Wicked",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
		},
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
		ActionID: core.ActionID{SpellID: 412732},
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.demonicKnowledgeSp = (warlock.Pet.GetStat(stats.Stamina) + warlock.Pet.GetStat(stats.Intellect)) * 0.1
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
	dmgMod := 1.05
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

	core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
		Label: "Grimoire of Synergy",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskDirect) {
				return
			}

			if sim.RandomFloat("Grimoire of Synergy") > procChance {
				return
			}

			petProcAura.Activate(sim)
		},
	}))

	core.MakePermanent(warlock.Pet.GetOrRegisterAura(core.Aura{
		Label: "Grimoire of Synergy",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskDirect) {
				return
			}

			if sim.RandomFloat("Grimoire of Synergy") > procChance {
				return
			}

			warlockProcAura.Activate(sim)
		},
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

	core.MakePermanent(warlock.GetOrRegisterAura(core.Aura{
		Label: "Shadow and Flame",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellSchool != core.SpellSchoolFire && spell.SpellSchool != core.SpellSchoolShadow {
				return
			}

			if !result.DidCrit() {
				return
			}

			procAura.Activate(sim)
		},
	}))
}

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

	spellPower := max(warlock.GetStat(stats.SpellPower)*0.1, float64(warlock.Level)/2.0)
	demonicPactAuras := warlock.NewRaidAuraArray(func(u *core.Unit) *core.Aura {
		return core.DemonicPactAura(u, spellPower)
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

			spBonus := max(math.Round(warlock.GetStat(stats.SpellPower)*0.1), math.Round(float64(warlock.Level)/2))
			for _, dpAura := range demonicPactAuras {
				if dpAura != nil {
					dpAura.ExclusiveEffects[0].SetPriority(sim, spBonus)

					// Force expire/gain because of new sp bonus
					dpAura.Deactivate(sim)
					dpAura.Activate(sim)
				}
			}
		},
	})
}
