package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 4 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDeathmistRaiment = core.NewItemSet(core.ItemSet{
	Name: "Deathmist Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 23)
		},
		// Your melee autoattacks and spellcasts have a 6% chance to heal you for 270 to 330 health.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			manaMetrics := c.NewManaMetrics(core.ActionID{SpellID: 450583})

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if c.HasManaBar() {
					c.AddMana(sim, sim.Roll(270, 300), manaMetrics)
				}
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Melee Auto)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskWhiteHit,
				ProcChance: 0.06,
				Handler:    handler,
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Heal Proc on Cast - Dreadmist Raiment (Spell Cast)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.06,
				Handler:    handler,
			})
		},
		// +8 All Resistances.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +200 Armor.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetCorruptedFelheart = core.NewItemSet(core.ItemSet{
	Name: "Corrupted Felheart",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT1Damage2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT1Damage4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT1Damage6PBonus()
		},
	},
})

// Lifetap generates 50% more mana and 100% less threat.
func (warlock *Warlock) applyT1Damage2PBonus() {
	label := "S03 - Item - T1 - Warlock - Damage 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.LifeTap {
				spell.MultiplyMultiplicativeDamageBonus(1.5)
				spell.ThreatMultiplier *= -1
			}
		},
	}))
}

// Increases your critical strike chance with spells and attacks by 2%.
func (warlock *Warlock) applyT1Damage4PBonus() {
	label := "S03 - Item - T1 - Warlock - Damage 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{
		stats.MeleeCrit: 2 * core.CritRatingPerCritChance,
		stats.SpellCrit: 2 * core.SpellCritRatingPerCritChance,
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachStatsBuff(bonusStats))
}

// Your Nightfall talent has a 4% increased chance to trigger.
// Incinerate has a 4% chance to trigger the Warlock’s Decimation.
func (warlock *Warlock) applyT1Damage6PBonus() {
	label := "S03 - Item - T1 - Warlock - Damage 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock6pt1Aura := warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(_ *core.Aura, _ *core.Simulation) {
			warlock.nightfallProcChance += 0.04
		},
	})

	if !warlock.HasRune(proto.WarlockRune_RuneBracerIncinerate) || !warlock.HasRune(proto.WarlockRune_RuneBootsDecimation) {
		return
	}

	core.MakePermanent(warlock6pt1Aura)
	warlock6pt1Aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if spell.Matches(ClassSpellMask_WarlockIncinerate) && result.Landed() && sim.Proc(.04, "T1 6P Incinerate Proc") {
			warlock.DecimationAura.Activate(sim)
		}
	}
}

var ItemSetWickedFelheart = core.NewItemSet(core.ItemSet{
	Name: "Wicked Felheart",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT1Tank2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT1Tank4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyT1Tank6PBonus()
		},
	},
})

// Banish is now instant cast, and can be cast on yourself while you are a Demon. You cannot Banish yourself while you have Forbearance, and doing so will give you Forbearance for 1 min.
func (warlock *Warlock) applyT1Tank2PBonus() {
	label := "S03 - Item - T1 - Warlock - Tank 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	// TODO: Implement if needed
	warlock.RegisterAura(core.Aura{
		Label: label,
	})
}

// Each time you take damage, you and your pet gain mana equal to the damage taken, up to a maximum of 420 mana per event. Can only occur once every few seconds.
func (warlock *Warlock) applyT1Tank4PBonus() {
	label := "S03 - Item - T1 - Warlock - Tank 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	actionID := core.ActionID{SpellID: 457572}
	icd := core.Cooldown{
		Timer:    warlock.NewTimer(),
		Duration: time.Millisecond * 3500,
	}
	manaMetrics := warlock.NewManaMetrics(actionID)
	for _, pet := range warlock.BasePets {
		pet.T1Tank4PManaMetrics = pet.NewManaMetrics(actionID)
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if icd.IsReady(sim) {
				restoreAmount := min(result.Damage, 420)
				warlock.AddMana(sim, restoreAmount, manaMetrics)
				if warlock.ActivePet != nil {
					warlock.ActivePet.AddMana(sim, restoreAmount, warlock.ActivePet.T1Tank4PManaMetrics)
				}
			}
		},
	}))
}

// Your Shadow Cleave hits have a 20% chance to grant you a Soul Shard, reset the cooldown on Soul Fire, and make your next Soul Fire within 10 sec instant.
func (warlock *Warlock) applyT1Tank6PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) {
		return
	}

	label := "S03 - Item - T1 - Warlock - Tank 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	procAura := warlock.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 457643},
		Label:    "Soul Fire!",
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SoulFire {
				spell.CastTimeMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.SoulFire {
				spell.CastTimeMultiplier += 1
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_WarlockSoulFire) {
				aura.Deactivate(sim)
			}
		},
	})

	icd := core.Cooldown{
		Timer:    warlock.NewTimer(),
		Duration: time.Millisecond * 100,
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.Matches(ClassSpellMask_WarlockShadowCleave) && icd.IsReady(sim) && sim.Proc(0.2, "Soul Fire! Proc") {
				procAura.Activate(sim)
				icd.Use(sim)
			}
		},
	}))
}
