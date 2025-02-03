package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetRaimentsOfFaith = core.NewItemSet(core.ItemSet{
	Name: "Raiments of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasShadow2PBonus()
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasShadow4PBonus()
		},
		6: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasShadow6PBonus()
		},
	},
})

// Your Shadow Word: Pain ability deals 20% more damage.
func (priest *Priest) applyNaxxramasShadow2PBonus() {
	label := "S03 - Item - Naxxramas - Priest - Shadow 2P Bonus"
	if priest.HasAura(label) {
		return
	}

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_PriestShadowWordPain,
		IntValue:  20,
	}))
}

// Reduces the cooldown on your Mind Blast ability by 1.0 sec.
func (priest *Priest) applyNaxxramasShadow4PBonus() {
	label := "S03 - Item - Naxxramas - Priest - Shadow 4P Bonus"
	if priest.HasAura(label) {
		return
	}

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_PriestMindBlast,
		TimeValue: -time.Second,
	}))
}

// Your Mind Flay, Mind Blast, and Mind Spike abilities deal increased damage to Undead targets equal to their critical strike chance.
func (priest *Priest) applyNaxxramasShadow6PBonus() {
	label := "S03 - Item - Naxxramas - Priest - Shadow 6P Bonus"
	if priest.HasAura(label) {
		return
	}

	classSpellMasks := ClassSpellMask_PriestMindBlast | ClassSpellMask_PriestMindSpike
	if priest.HasRune(proto.PriestRune_RuneBracersDespair) {
		classSpellMasks |= ClassSpellMask_PriestMindFlay
	}

	damageMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  classSpellMasks,
		FloatValue: 1,
	})

	core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
		Name:           label,
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: classSpellMasks,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target.MobType != proto.MobType_MobTypeUndead {
				return
			}

			critChanceBonusPct := 100.0
			critChanceBonusPct += priest.GetStat(stats.SpellCrit)

			damageMod.UpdateFloatValue(critChanceBonusPct / 100)
		},
	}).ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		damageMod.Activate()
	}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		damageMod.Deactivate()
	})
}

var ItemSetVestmentsOfFaith = core.NewItemSet(core.ItemSet{
	Name: "Vestments of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.applyNaxxramasHealer2PBonus()
		},
		// Your Penance, Flash Heal Rank 7, and Greater Heal Rank 4 and Rank 5 have a 9% chance to grant the target 10% increased critical strike chance for 15 sec.
		4: func(agent core.Agent) {
		},
		// Your Power Word: Shield has a 50% chance to not deplete when the target is damaged by an Undead enemy.
		6: func(agent core.Agent) {
		},
	},
})

// Reduces the cooldown on your Circle of Healing and Penance abilities by 25%.
func (priest *Priest) applyNaxxramasHealer2PBonus() {
	if !priest.HasRune(proto.PriestRune_RuneHandsPenance) {
		return
	}

	label := "S03 - Item - Naxxramas - Priest - Healer 2P Bonus"
	if priest.HasAura(label) {
		return
	}

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Multi_Flat,
		ClassMask: ClassSpellMask_PriestPenance,
		IntValue:  -25,
	}))
}
