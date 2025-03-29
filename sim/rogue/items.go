package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	BloodSpatteredStilletto          = 216522
	ShadowflameSword                 = 228143
	DreamEater                       = 224122
	VenomousTotem                    = 230250
	RenatakisCharmofTrickery         = 231287
	ZandalarianShadowMasteryTalisman = 231336
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(BloodSpatteredStilletto, func(agent core.Agent) {
		character := agent.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436477},
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					//Confirmed always hits through logs
					spell.CalcAndDealDamage(sim, aoeTarget, 140, spell.OutcomeAlwaysHit)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=228143/shadowflame-sword
	core.NewItemEffect(ShadowflameSword, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()

		if !rogue.Talents.BladeFlurry {
			return
		}

		// TODO: This is treated as a buff, NOT a debuff in-game
		// We don't have the ability to remove resistances for only one agent at a time right now
		procAura := rogue.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 461252},
			Label:    "Shadowflame Fury",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				for _, target := range sim.Encounter.TargetUnits {
					target.AddStatDynamic(sim, stats.Armor, -2000)
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				for _, target := range sim.Encounter.TargetUnits {
					target.AddStatDynamic(sim, stats.Armor, 2000)
				}
			},
		})

		aura := core.MakePermanent(rogue.RegisterAura(core.Aura{
			Label: "Shadowflame Fury Trigger",
			OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(ClassSpellMask_RogueBladeFlurry) {
					procAura.Duration = rogue.BladeFlurryAura.Duration
					procAura.Activate(sim)
				}
			},
		}))

		rogue.ItemSwap.RegisterProc(ShadowflameSword, aura)
	})

	// https://www.wowhead.com/classic/item=224122/dream-eater
	// Damaging finishing moves have a 20% chance per combo point to restore 10 energy.
	core.NewItemEffect(DreamEater, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()

		aura := core.MakePermanent(rogue.RegisterAura(core.Aura{
			Label:    "Dream Eater Dummy",
			ActionID: core.ActionID{ItemID: DreamEater},
		}))

		cpMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 451439})
		rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
			if aura.IsActive() && spell.Matches(ClassSpellMask_RogueBetweentheEyes|ClassSpellMask_RogueCrimsonTempest|ClassSpellMask_RogueEnvenom|ClassSpellMask_RogueEviscerate|ClassSpellMask_RogueRupture) {
				if sim.Proc(0.2*float64(comboPoints), "Dream Eater") {
					rogue.AddEnergy(sim, 10, cpMetrics)
				}
			}
		})

		rogue.ItemSwap.RegisterProc(ShadowflameSword, aura)
	})

	// https://www.wowhead.com/classic/item=231287/renatakis-charm-of-trickery
	// Use: Instantly increases your energy by 60. If Cutthroat is engraved, gain an activation of Cutthroat's Ambush effect. (2 Min Cooldown)
	core.NewItemEffect(RenatakisCharmofTrickery, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()
		cpMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 468458})
		hasCutthroatRune := rogue.HasRune(proto.RogueRune_RuneCutthroat)

		spell := rogue.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: RenatakisCharmofTrickery},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    rogue.NewTimer(),
					Duration: time.Second * 120,
				},
				SharedCD: core.Cooldown{
					Timer:    rogue.GetOffensiveTrinketCD(),
					Duration: time.Second * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				rogue.AddEnergy(sim, 60, cpMetrics)
				if hasCutthroatRune {
					rogue.CutthroatProcAura.Activate(sim)
				}
			},
		})

		rogue.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				// Make sure we have plenty of room so we dont energy cap right after using.
				return rogue.CurrentEnergy() <= 40
			},
		})

	})

	// https://www.wowhead.com/classic/item=230250/venomous-totem
	core.NewItemEffect(VenomousTotem, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()

		aura := rogue.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 467511},
			Label:    "Venomous Totem",
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				rogue.additivePoisonBonusChance += 0.3
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				rogue.additivePoisonBonusChance -= 0.3
			},
		})

		spell := rogue.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: VenomousTotem},
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    rogue.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    rogue.GetOffensiveTrinketCD(),
					Duration: time.Second * 20,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		rogue.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	// https://www.wowhead.com/classic/item=231336/zandalarian-shadow-mastery-talisman
	// Increases the chance Cutthroat's Ambush effect is triggered by 5%.
	core.NewItemEffect(ZandalarianShadowMasteryTalisman, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()
		rogue.cutthroatBonusChance += 0.05
	})

	core.AddEffectsToTest = true
}

// Your Backstab, Mutilate, and Saber Slash critical strikes set the duration of your Rupture on the target to 16 secs
func (rogue *Rogue) ApplyQueensfallRogueEffect(aura *core.Aura) {
	aura.AttachProcTrigger(core.ProcTrigger{
		Name:           "Queensfall Trigger - Rogue",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeCrit,
		ClassSpellMask: ClassSpellMask_RogueBackstab | ClassSpellMask_RogueMutilateHit | ClassSpellMask_RogueSaberSlash,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if dot := rogue.Rupture.Dot(result.Target); dot.IsActive() {
				dot.NumberOfTicks = int32(16 / dot.TickLength.Seconds())
				dot.RecomputeAuraDuration()
				dot.Rollover(sim)
			}
		},
	})
}

// Striking a higher level enemy applies a stack of Coup, increasing their damage taken from your next Envenom by 5% per stack, stacking up to 20 times.
func (rogue *Rogue) ApplyRegicideRogueEffect(aura *core.Aura) {
	// Coup debuff array
	debuffAuras := rogue.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
		aura := unit.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1231417},
			Label:     "Coup",
			MaxStacks: 20,
			Duration:  core.NeverExpires,
		})

		return aura
	})

	envenomDamageMod := rogue.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: ClassSpellMask_RogueEnvenom,
	})

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: "Coup Execute Damage",
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Coup Trigger - Rogue",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: ClassSpellMask_RogueEnvenom,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuffAuras[result.Target.Index].SetStacks(sim, 0)
			}
		},
	}))

	// Apply the Coup debuff to the target hit by melee abilities
	aura.AttachProcTrigger(core.ProcTrigger{
		Name:     "Regicide Trigger - Rogue",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		ProcMask: core.ProcMaskMelee,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			debuffAuras[result.Target.Index].Activate(sim)
			debuffAuras[result.Target.Index].AddStack(sim)
			envenomDamageMod.UpdateFloatValue(1 + float64(debuffAuras[result.Target.Index].GetStacks())*0.05)
			envenomDamageMod.Activate()
		},
	})
}