package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Damaging finishing moves have a 20% chance per combo point to restore 20 energy.
func (rogue *Rogue) ApplyFallenRegalityRogueBonus(aura *core.Aura) {
	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1232184})
	aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
			if spell.ProcMask != core.ProcMaskEmpty && sim.Proc(0.20*float64(comboPoints), "Fallen Regality Proc") {
				rogue.AddEnergy(sim, 20, energyMetrics)
			}
		})
	})
}

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
			rogue.RollingWithThePunchesProcAura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				rogue.PseudoStats.DamageDealtMultiplier /= 1 + 0.01*float64(oldStacks)
				rogue.PseudoStats.DamageDealtMultiplier *= 1 + 0.01*float64(newStacks)
			})
		},
	})
}

// Your Blade Flurry now also strikes a third target and increases your attack speed by an additional 10%. In addition, each combo point you spend reduces the remaining cooldown on your Blade Flurry by 0.5 sec.
func (rogue *Rogue) applyScarletEnclaveTank4PBonus() {

	if !rogue.Talents.BladeFlurry {
		return
	}

	label := "S03 - Item - Scarlet Enclave - Rogue - Tank 4P Bonus"

	if rogue.HasAura(label) {
		return
	}

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
				cdReduction := time.Millisecond * time.Duration(500) * time.Duration(comboPoints)
				rogue.BladeFlurry.CD.ModifyRemainingCooldown(sim, -cdReduction)
			})

			rogue.bladeFlurryAttackSpeedBonus += 0.1
			rogue.bladeFlurryTargetCount += 1

		},
	}))
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

	metrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 1226957})

	energyProc := rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 1226957},
		Flags:    core.SpellFlagNoOnCastComplete,
		ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
			rogue.AddEnergy(sim, 15, metrics)
		},
	})

	energyAura := rogue.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 1226957},
		Label:    "Float Like a Butterfly, Sting Like a Bee",
		Duration: core.NeverExpires,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeDodge) {
				energyProc.Cast(sim, result.Target)
			}
		},
	})

	core.MakePermanent(rogue.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			rogue.rollingWithThePunchesBonusHealthStackMultiplier -= 0.02
			rogue.RollingWithThePunchesProcAura.MaxStacks += 5
			rogue.rollingWithThePunchesMaxStacks += 5
			rogue.RollingWithThePunchesProcAura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				if newStacks == 10 {
					energyAura.Activate(sim)
				} else if newStacks < 10 && oldStacks == 10 {
					energyAura.Deactivate(sim)
				}
			})
		},
	}))
}
