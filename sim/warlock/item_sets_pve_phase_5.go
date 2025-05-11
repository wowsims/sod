package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 5 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetWickedNemesis = core.NewItemSet(core.ItemSet{
	Name: "Wicked Nemesis",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT2Tank2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT2Tank4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT2Tank6PBonus()
		},
	},
})

// While you are targeting an enemy within 30 yards, Life Tap grants you mana at the expense of your target's health but deals 50% reduced damage to them.
// Mana gained remains unchanged.
//
// Bluepost: Tier 2 tank warlock 2 set can no longer crit
func (warlock *Warlock) applyT2Tank2PBonus() {
	label := "S03 - Item - T2 - Warlock - Tank 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	var healthMetrics [LifeTapRanks + 1]*core.ResourceMetrics

	for i, spellId := range LifeTapSpellId {
		healthMetrics[i] = warlock.NewHealthMetrics(core.ActionID{SpellID: spellId})
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarlockLifeTap) {
				warlock.GainHealth(sim, result.Damage, healthMetrics[spell.Rank])
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_WarlockLifeTap) && warlock.CurrentTarget != nil {
				// Enemy hit can partially resist and cannot crit
				spell.Flags &= ^core.SpellFlagBinary
				spell.ApplyMultiplicativeDamageBonus(1 / 2)
				damageResult := spell.CalcDamage(sim, warlock.CurrentTarget, LifeTapBaseDamage[spell.Rank], spell.OutcomeMagicHit)
				spell.DealDamage(sim, damageResult)
				spell.ApplyMultiplicativeDamageBonus(2)
				spell.Flags |= core.SpellFlagBinary
			}
		},
	}))
}

// While Metamorphosis is active, your offensive abilities and Demon summons cost no Soul Shards.
// In addition, you heal for 15% of your maximum health when you damage a target with Shadowburn
func (warlock *Warlock) applyT2Tank4PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) {
		return
	}

	label := "S03 - Item - T2 - Warlock - Tank 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	healingSpell := warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 468062},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 0,
	})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarlockShadowburn) && result.Landed() {
				healAmount := warlock.MaxHealth() * 0.15
				healingSpell.CalcAndDealHealing(sim, healingSpell.Unit, healAmount, healingSpell.OutcomeHealing)
			}
		},
	}))
}

// Any excess healing you deal to yourself is converted into a shield that absorbs damage.
// This shield can absorb uf to 30% of your maximum health, and stacks from multiple heals.
func (warlock *Warlock) applyT2Tank6PBonus() {
	label := "S03 - Item - T2 - Warlock - Tank 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	shieldIncreaseAmount := 0.0
	currentShieldAmount := 0.0

	shieldSpell := warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 470279},
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Shield: core.ShieldConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label:    "Demonic Aegis",
				Duration: 12 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SelfShield().Apply(sim, currentShieldAmount)
		},
	})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			shieldIncreaseAmount = 0.0
			currentShieldAmount = 0.0
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if currentShieldAmount <= 0 || result.Damage <= 0 || spell.Matches(ClassSpellMask_WarlockLifeTap) {
				return
			}

			damageAbsorbed := min(result.Damage, currentShieldAmount)
			currentShieldAmount -= damageAbsorbed

			warlock.GainHealth(sim, damageAbsorbed, shieldSpell.HealthMetrics(result.Target))

			if currentShieldAmount <= 0 {
				currentShieldAmount = 0
				shieldSpell.SelfShield().Deactivate(sim)
			}
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shieldIncreaseAmount = result.Damage + warlock.CurrentHealth() - warlock.MaxHealth()
			if shieldIncreaseAmount > 0 {
				currentShieldAmount += shieldIncreaseAmount
				currentShieldAmount = min(warlock.MaxHealth()*0.30, currentShieldAmount)
				shieldSpell.Cast(sim, result.Target)
			}
		},
	}))
}

var ItemSetCorruptedNemesis = core.NewItemSet(core.ItemSet{
	Name: "Corrupted Nemesis",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT2Damage2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT2Damage4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT2Damage6PBonus()
		},
	},
})

// Increases the damage of your periodic spells and Felguard pet by 10%
func (warlock *Warlock) applyT2Damage2PBonus() {
	label := "S03 - Item - T2 - Warlock - Damage 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.Spellbook {
				if spell.Matches(ClassSpellMask_WarlockAll) && len(spell.Dots()) > 0 {
					spell.ApplyAdditivePeriodicDamageBonus(10)
				}
			}

			if !warlock.HasRune(proto.WarlockRune_RuneBracerSummonFelguard) {
				return
			}

			warlock.Felguard.ApplyOnPetEnable(func(sim *core.Simulation) {
				warlock.Felguard.PseudoStats.DamageDealtMultiplier *= 1.10
			})
			warlock.Felguard.ApplyOnPetDisable(func(sim *core.Simulation, isSacrifice bool) {
				warlock.Felguard.PseudoStats.DamageDealtMultiplier /= 1.10
			})
		},
	})
}

// Periodic damage from your Shadowflame, Unstable Affliction, and Curse of Agony spells and damage done by your Felguard have a 4% chance to grant the Shadow Trance effect.
func (warlock *Warlock) applyT2Damage4PBonus() {
	label := "S03 - Item - T2 - Warlock - Damage 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	procChance := 0.04

	affectedSpellClassMasks := ClassSpellMask_WarlockShadowflame | ClassSpellMask_WarlockCurseOfAgony | ClassSpellMask_WarlockUnstableAffliction
	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(affectedSpellClassMasks) && sim.Proc(procChance, "Proc Shadow Trance") {
				warlock.ShadowTranceAura.Activate(sim)
			}
		},
	}))

	if !warlock.HasRune(proto.WarlockRune_RuneBracerSummonFelguard) {
		return
	}

	core.MakePermanent(warlock.Felguard.RegisterAura(core.Aura{
		Label: label + " - Felguard Bonus",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && sim.Proc(procChance, "Proc Shadow Trance") {
				warlock.ShadowTranceAura.Activate(sim)
			}
		},
	}))
}

// Shadowbolt deals 10% increased damage for each of your effects afflicting the target, up to a maximum of 30%.
func (warlock *Warlock) applyT2Damage6PBonus() {
	label := "S03 - Item - T2 - Warlock - Damage 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	damageMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  ClassSpellMask_WarlockShadowBolt,
		FloatValue: 1.0,
	})

	core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
		Name:           label,
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: ClassSpellMask_WarlockShadowBolt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			damageMod.UpdateFloatValue(1.0 + min(0.30, 0.10*float64(warlock.activeEffects[result.Target.UnitIndex])))
			damageMod.Activate()
		},
	})
}

var ItemSetDemoniacsThreads = core.NewItemSet(core.ItemSet{
	Name: "Demoniac's Threads",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 12.
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.AddStat(stats.SpellPower, 12)
		},
		3: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyZGDemonology3PBonus()
		},
		5: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyZGDemonology5PBonus()
		},
	},
})

// Increases the Attack Power and Spell Damage your Demon pet gains from your attributes by 20%.
func (warlock *Warlock) applyZGDemonology3PBonus() {
	label := "S03 - Item - ZG - Warlock - Demonology 3P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, pet := range warlock.BasePets {
				oldStatInheritance := pet.GetStatInheritance()
				pet.UpdateStatInheritance(
					func(ownerStats stats.Stats) stats.Stats {
						s := oldStatInheritance(ownerStats)
						s[stats.AttackPower] *= 1.20
						s[stats.SpellPower] *= 1.20
						return s
					},
				)
			}
		},
	})
}

// Increases the benefits of your Master Demonologist talent by 50%.
func (warlock *Warlock) applyZGDemonology5PBonus() {
	label := "S03 - Item - ZG - Warlock - Demonology 5P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warlock.masterDemonologistMultiplier += .50
		},
	})
}
