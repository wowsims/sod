package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            SoD Phase 6 Item Sets
///////////////////////////////////////////////////////////////////////////

var ItemSetDoomcallersCorruption = core.NewItemSet(core.ItemSet{
	Name: "Doomcaller's Corruption",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyTAQDamage2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyTAQDamage4PBonus()
		},
	},
})

// Reduces the cooldown on your Chaos Bolt by 50% and increases Chaos Bolt and Shadow Bolt damage done by 10%.
// In addition, Chaos Bolt can now trigger your Improved Shadow Bolt talent.
func (warlock *Warlock) applyTAQDamage2PBonus() {
	label := "S03 - Item - TAQ - Warlock - Damage 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	aura := warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowBolt {
				spell.DamageMultiplierAdditive += 0.10
			}

			if warlock.ChaosBolt != nil {
				warlock.ChaosBolt.CD.Multiplier *= 0.5
				warlock.ChaosBolt.DamageMultiplierAdditive += 0.10
			}
		},
	})

	if !warlock.HasRune(proto.WarlockRune_RuneHandsChaosBolt) {
		return
	}

	core.MakePermanent(aura)
	aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if result.Landed() && result.DidCrit() && spell.SpellCode == SpellCode_WarlockChaosBolt {
			isbAura := warlock.ImprovedShadowBoltAuras.Get(result.Target)
			isbAura.Activate(sim)
			// This set always uses 30 stacks
			isbAura.SetStacks(sim, core.ISBNumStacksShadowflame)
		}
	}
}

// Each time you hit a target with Conflagrate, you gain 5% increased Fire damage for 20 sec, stacking up to 2 times.
func (warlock *Warlock) applyTAQDamage4PBonus() {
	if !warlock.Talents.Conflagrate {
		return
	}

	label := "S03 - Item - TAQ - Warlock - Damage 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	buffAura := warlock.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1214088},
		Label:     "Infernalist",
		Duration:  time.Second * 20,
		MaxStacks: 2,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1 + 0.05*(float64(newStacks-oldStacks))
		},
	})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode == SpellCode_WarlockConflagrate && result.Landed() {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetDoomcallersMalevolence = core.NewItemSet(core.ItemSet{
	Name: "Doomcaller's Malevolence",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyTAQTank2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyTAQTank4PBonus()
		},
	},
})

// Reduces the cooldown on your Shadow Cleave by 1.5 sec.
func (warlock *Warlock) applyTAQTank2PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) {
		return
	}

	label := "S03 - Item - TAQ - Warlock - Tank 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.ShadowCleave {
				spell.CD.FlatModifier -= time.Millisecond * 1500
			}
		},
	})
}

// The effects of your Demonic Sacrifice now persist while you have a Demon pet active, as long as you do not resummon the sacrificed pet.
// You may have only one Demonic Sacrifice effect active at a time.
func (warlock *Warlock) applyTAQTank4PBonus() {
	label := "S03 - Item - TAQ - Warlock - Tank 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			warlock.maintainBuffsOnSacrifice = true
		},
	}))
}

var ItemSetImplementsOfUnspokenNames = core.NewItemSet(core.ItemSet{
	Name: "Implements of Unspoken Names",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyRAQTank3PBonus()
		},
	},
})

// For 6 sec after using Shadowcleave, your Searing Pain strikes 1 additional target within melee range.
func (warlock *Warlock) applyRAQTank3PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsMetamorphosis) || len(warlock.Env.Encounter.TargetUnits) <= 1 {
		return
	}

	label := "S03 - Item - RAQ - Warlock - Tank 3P Bonus"
	if warlock.HasAura(label) {
		return
	}

	buffAura := warlock.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1214156},
		Label:    "Spreading Pain",
		Duration: time.Second * 6,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode == SpellCode_WarlockSearingPain {
				aura.Deactivate(sim)
				spell.ApplyEffects(sim, warlock.Env.NextTargetUnit(result.Target), spell)
			}
		},
	})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.SpellCode == SpellCode_WarlockShadowCleave {
				buffAura.Activate(sim)
			}
		},
	}))
}
