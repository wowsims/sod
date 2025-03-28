package rogue

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetDuskwraithArmor = core.NewItemSet(core.ItemSet{
	Name: "Duskwraith Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
		},
		6: func(agent core.Agent) {
		},
	},
})

var ItemSetDuskwraithLeathers = core.NewItemSet(core.ItemSet{
	Name: "Duskwraith Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveTank2PBonus()
		},
		4: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveTank4PBonus()
		},
		6: func(agent core.Agent) {
			rogue := agent.(RogueAgent).GetRogue()
			rogue.applyScarletEnclaveTank6PBonus()
		},
	},
})

// Your stacks of Rolling with the Punches also increase all damage you deal by 1%.
func (rogue *Rogue) applyScarletEnclaveTank2PBonus() {

	if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Tank 2P Bonus"
	if rogue.HasAura(label) {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldOnStacksChange := rogue.RollingWithThePunchesProcAura.OnStacksChange
			rogue.RollingWithThePunchesProcAura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				oldOnStacksChange(aura, sim, oldStacks, newStacks)
				rogue.PseudoStats.DamageDealtMultiplierAdditive += float64(0.1 * float64(newStacks-oldStacks))
			}
		},
	})
}

// Your Blade Flurry now also strikes a third target and increases your attack speed by an additional 10%. In addition, each combo point you spend reduces the remaining cooldown on your Blade Flurry by 0.5 sec.
func (rogue *Rogue) applyScarletEnclaveTank4PBonus() {
	label := "S03 - Item - Scarlet Enclave - Rogue - Tank 4P Bonus"

	if rogue.HasAura(label) {
		return
	}

}

// Your Rolling with the Punches can now stack up to 10 times, but grants 2% less health per stack. At 10 stacks, each time you Dodge you will gain 15 Energy.
func (rogue *Rogue) applyScarletEnclaveTank6PBonus() {

	if !rogue.HasRune(proto.RogueRune_RuneRollingWithThePunches) {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Tank 6P Bonus"

	if rogue.HasAura(label) {
		return
	}

	statDeps := make([]*stats.StatDependency, 11) // 10 stacks + zero condition
	for i := 1; i < 11; i++ {
		statDeps[i] = rogue.NewDynamicMultiplyStat(stats.Health, 1.0+.04*float64(i)) // 4% health per stack

	}
	metrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1226956})

	energyProc := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1226956},
		SpellSchool: core.SpellSchoolNature,
		ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
			rogue.AddEnergy(sim, 15, metrics)
		},
	})

	energyAura := rogue.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1219291},
		Label:    "S03 - Item - Scarlet Enclave - Rogue - Tank 6P Bonus Energy Gain",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeDodge) {
				energyProc.Cast(sim, result.Target)
			}
		},
	})

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.RollingWithThePunchesProcAura.MaxStacks += 5
			rogue.RollingWithThePunchesProcAura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				if oldStacks != 0 {
					aura.Unit.DisableDynamicStatDep(sim, statDeps[oldStacks])
				}
				if newStacks != 0 {
					aura.Unit.EnableDynamicStatDep(sim, statDeps[newStacks])
				}
				if newStacks == 10 {
					energyAura.Activate(sim)
				}
				if newStacks != 10 && oldStacks == 10 {
					energyAura.Deactivate(sim)
				}

				// repeat the 2p set bonus because we need to override the whole onStackChange because of health scaling changes
				rogue.PseudoStats.DamageDealtMultiplierAdditive += float64(0.1 * float64(newStacks-oldStacks))
			}
		},
	}))
}
