package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetFrostfireRegalia = core.NewItemSet(core.ItemSet{
	Name: "Frostfire Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyNaxxramasDamage2PBonus()
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyNaxxramasDamage4PBonus()
		},
		6: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.applyNaxxramasDamage6PBonus()
		},
	},
})

// Your Evocation grants you 5% increased damage done every sec you channel it, stacking up to 8 times and lasting 45 sec.
func (mage *Mage) applyNaxxramasDamage2PBonus() {
	label := "S03 - Item - Naxxramas - Mage - Damage 2P Bonus"
	if mage.HasAura(label) {
		return
	}

	buffAura := mage.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1218701},
		Label:     "Evoker",
		Duration:  time.Second * 45,
		MaxStacks: 8,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			mage.PseudoStats.DamageDealtMultiplier /= 1 + 0.05*float64(oldStacks)
			mage.PseudoStats.DamageDealtMultiplier *= 1 + 0.05*float64(newStacks)
		},
	})

	mage.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// var buffAuraPA *core.PendingAction

			hot := mage.Evocation.SelfHot()
			// oldOnGain := hot.OnGain

			hot.OnTick = func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Tick every 4th tick, aka every 1s
				if (hot.NumTicksRemaining(sim)+3)%4 == 0 {
					buffAura.Activate(sim)
					buffAura.AddStack(sim)
				}
			}
		},
	})
}

// Reduces the cooldown on your Evocation ability by 80%.
func (mage *Mage) applyNaxxramasDamage4PBonus() {
	label := "S03 - Item - Naxxramas - Mage - Damage 4P Bonus"
	if mage.HasAura(label) {
		return
	}

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  ClassSpellMask_MageEvocation,
		Kind:       core.SpellMod_Cooldown_Multi_Pct,
		FloatValue: 0.2,
	}))
}

// Your Ignite damage does not decay on Undead targets below 20% health, and Undead targets below 20% health take damage as if they were Frozen.
func (mage *Mage) applyNaxxramasDamage6PBonus() {
	if mage.Talents.Ignite == 0 && !mage.HasRune(proto.MageRune_RuneChestFingersOfFrost) {
		return
	}

	label := "S03 - Item - Naxxramas - Mage - Damage 6P Bonus"
	if mage.HasAura(label) {
		return
	}

	mage.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldProcIgnite := mage.procIgnite
			mage.procIgnite = func(sim *core.Simulation, result *core.SpellResult) {
				dot := mage.Ignite.Dot(result.Target)
				if dot.TickCount > 0 && sim.IsExecutePhase20() && result.Target.MobType == proto.MobType_MobTypeUndead {
					// Effectively negate the decay by reducing the TickCount
					dot.TickCount -= 1
				}

				oldProcIgnite(sim, result)
			}

			oldIsTargetFrozen := mage.isTargetFrozen
			mage.isTargetFrozen = func(target *core.Unit) bool {
				return (sim.IsExecutePhase20() && target.MobType == proto.MobType_MobTypeUndead) || oldIsTargetFrozen(target)
			}
		},
	})
}

var ItemSetFrostfireVestments = core.NewItemSet(core.ItemSet{
	Name: "Frostfire Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		// Allies with your Temporal Beacon heal for [Spell power * 67 / 1000 + (38.258376 + 0.904195 * 60 + 0.161311 * 60 * 60) * 8 / 100] health every 1 sec.
		2: func(agent core.Agent) {
		},
		// Your Regeneration ability grants your target 60% increased movement speed while you are channeling, and each time it heals your target, they have a chance to gain 10% increased attack and casting speed for 15 sec.
		4: func(agent core.Agent) {
		},
		// Damage you deal to Undead causes 25% more chronomantic healing, and you gain mana equal to 5% of the chronomantic healing you generate from damaging Undead.
		6: func(agent core.Agent) {
		},
	},
})
