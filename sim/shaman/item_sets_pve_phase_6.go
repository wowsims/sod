package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 6 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetStormcallersEruption = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Eruption",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyTAQElemental2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyTAQElemental4PBonus()
		},
	},
})

// You have a 70% chance to avoid interruption caused by damage while casting Lightning Bolt, Chain Lightning, or Lava Burst, and a 10% increased chance to trigger your Elemental Focus talent.
func (shaman *Shaman) applyTAQElemental2PBonus() {
	label := "S03 - Item - TAQ - Shaman - Elemental 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedPushbackSpells := core.FilterSlice(
				core.Flatten(
					[][]*core.Spell{
						shaman.LightningBolt,
						shaman.ChainLightning,
						{shaman.LavaBurst},
					},
				),
				func(spell *core.Spell) bool { return spell != nil },
			)

			for _, spell := range affectedPushbackSpells {
				spell.PushbackReduction += .70
			}

			if shaman.Talents.ElementalFocus {
				shaman.elementalFocusProcChance += .10
			}
		},
	})
}

// Increases the critical strike damage bonus of your Fire, Frost, and Nature spells by 60%.
func (shaman *Shaman) applyTAQElemental4PBonus() {
	label := "S03 - Item - TAQ - Shaman - Elemental 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range shaman.Spellbook {
				if spell.Flags.Matches(SpellFlagShaman) && spell.DefenseType == core.DefenseTypeMagic {
					spell.CritDamageBonus += 0.60
				}
			}
		},
	})
}

var ItemSetStormcallersResolve = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Resolve",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyTAQTank2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyTAQTank4PBonus()
		},
	},
})

// Damaging a target with Stormstrike, Lava Burst, or Molten Blast also reduces all damage you take by 10% for 10 sec.
func (shaman *Shaman) applyTAQTank2PBonus() {
	if !shaman.Talents.Stormstrike && !shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) && !shaman.HasRune(proto.ShamanRune_RuneHandsMoltenBlast) {
		return
	}

	label := "S03 - Item - TAQ - Shaman - Tank 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	affectedSpellClassMasks := ClassSpellMask_ShamanStormstrikeHit | ClassSpellMask_ShamanLavaBurst | ClassSpellMask_ShamanMoltenBlast

	buffAura := shaman.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1213934},
		Label:    "Stormbraced",
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.90
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.90
		},
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(affectedSpellClassMasks) {
				buffAura.Activate(sim)
			}
		},
	}))
}

// Your Spirit of the Alpha also increases your health by 10%, threat by 20%, and damage by 10% when cast on self.
func (shaman *Shaman) applyTAQTank4PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneFeetSpiritOfTheAlpha) || !shaman.IsTanking() {
		return
	}

	label := "S03 - Item - TAQ - Shaman - Tank 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	damageMultiplier := 1.10
	threatMultiplier := 1.20
	healthMultiplier := 1.10
	statDep := shaman.NewDynamicMultiplyStat(stats.Health, healthMultiplier)

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableBuildPhaseStatDep(sim, statDep)

			shaman.PseudoStats.DamageDealtMultiplier *= damageMultiplier
			shaman.PseudoStats.ThreatMultiplier *= threatMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableBuildPhaseStatDep(sim, statDep)

			shaman.PseudoStats.DamageDealtMultiplier /= damageMultiplier
			shaman.PseudoStats.ThreatMultiplier /= threatMultiplier
		},
	}))
}

var ItemSetStormcallersRelief = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Relief",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Riptide increases the amount healed by Chain Heal by an additional 25%.
		2: func(agent core.Agent) {
		},
		// Reduces the cast time of Chain Heal by 0.5 sec.
		4: func(agent core.Agent) {
		},
	},
})

var ItemSetStormcallersImpact = core.NewItemSet(core.ItemSet{
	Name: "Stormcaller's Impact",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyTAQEnhancement2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyTAQEnhancement4PBonus()
		},
	},
})

// Increases Stormstrike and Lava Lash damage by 50%. Stormstrike's damage is increased by an additional 50% when using a Two-handed weapon.
func (shaman *Shaman) applyTAQEnhancement2PBonus() {
	if !shaman.Talents.Stormstrike && !shaman.HasRune(proto.ShamanRune_RuneHandsLavaLash) {
		return
	}

	label := "S03 - Item - TAQ - Shaman - Enhancement 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if shaman.StormstrikeMH != nil {
				shaman.StormstrikeMH.DamageMultiplierAdditive += core.TernaryFloat64(shaman.HasRune(proto.ShamanRune_RuneChestTwoHandedMastery), 1.00, 0.50)
			}

			if shaman.StormstrikeOH != nil {
				shaman.StormstrikeOH.DamageMultiplierAdditive += 0.50
			}

			if shaman.LavaLash != nil {
				shaman.LavaLash.DamageMultiplierAdditive += 0.50
			}
		},
	})
}

// Your Stormstrike, Lava Lash, and Lava Burst critical strikes cause your target to burn for 30% of the damage done over 4 sec.
func (shaman *Shaman) applyTAQEnhancement4PBonus() {
	if !shaman.Talents.Stormstrike && !shaman.HasRune(proto.ShamanRune_RuneHandsLavaLash) && !shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
		return
	}

	label := "S03 - Item - TAQ - Shaman - Enhancement 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	// This is the spell used for the burn proc.
	// https://www.wowhead.com/classic/spell=1213915/burning
	burnSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1213915},
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagIgnoreAttackerModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Burning",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrRefresh(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	affectedSpellClassMasks := ClassSpellMask_ShamanStormstrikeHit | ClassSpellMask_ShamanLavaLash | ClassSpellMask_ShamanLavaBurst

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:             label,
		Callback:         core.CallbackOnSpellHitDealt,
		Outcome:          core.OutcomeCrit,
		CanProcFromProcs: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(affectedSpellClassMasks) {
				return
			}

			dot := burnSpell.Dot(result.Target)
			dotDamage := result.Damage * 0.3
			if dot.IsActive() {
				dotDamage += dot.SnapshotBaseDamage * float64(dot.MaxTicksRemaining())
			}
			dot.SnapshotBaseDamage = dotDamage / float64(dot.NumberOfTicks)
			dot.SnapshotAttackerMultiplier = 1

			burnSpell.Cast(sim, result.Target)
		},
	})
}

var ItemSetGiftOfTheGatheringStorm = core.NewItemSet(core.ItemSet{
	Name: "Gift of the Gathering Storm",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyRAQElemental3PBonus()
		},
	},
})

// Your Lava Burst deals increased damage equal to its critical strike chance.
func (shaman *Shaman) applyRAQElemental3PBonus() {
	label := "S03 - Item - RAQ - Shaman - Elemental 3P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.useLavaBurstCritScaling = true
		},
	})
}
