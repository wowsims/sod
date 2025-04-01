package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var ItemSetCryptstalkerProwess = core.NewItemSet(core.ItemSet{
	Name: "Cryptstalker Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasMelee2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasMelee4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasMelee6PBonus()
		},
	},
})

// Your Wyvern Strike and Mongoose Bite deal 30% more initial damage.
func (hunter *Hunter) applyNaxxramasMelee2PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Melee 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_ImpactDamageDone_Flat,
		ClassMask: ClassSpellMask_HunterWyvernStrike | ClassSpellMask_HunterMongooseBite,
		IntValue:  30,
	}))
}

// Reduces the cooldown on your Wyvern Strike ability by 2 sec, reduces the cooldown on your raptor strike ability by 1 sec, and reduces the cooldown on your Flanking Strike ability by 8sec.
func (hunter *Hunter) applyNaxxramasMelee4PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Melee 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	// Not entirely sure how this will work so taking some liberties
	// Assume that it resets all of them when one crits
	//var spellsToReset []*core.Spell

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_HunterRaptorStrike,
		TimeValue: -time.Second,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_HunterWyvernStrike,
		TimeValue: -time.Second * 2,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_HunterFlankingStrike,
		TimeValue: -time.Second * 8,
	}))
}

// You gain 2% increased damage and critical damage done to Undead for 30 sec each time you hit an Undead enemy with a melee attack, stacking up to 12 times. Stacks are lost upon performing a ranged attack.
func (hunter *Hunter) applyNaxxramasMelee6PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Melee 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	undeadTargets := core.FilterSlice(hunter.Env.Encounter.TargetUnits, func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeUndead })

	units := []*core.Unit{&hunter.Unit}
	if hunter.pet != nil {
		units = append(units, &hunter.pet.Unit)
	}

	buffAuras := []*core.Aura{}
	for _, unit := range units {
		buffAuras = append(buffAuras, unit.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1218587},
			Label:     "Undead Slaying",
			Duration:  time.Second * 30,
			MaxStacks: 12,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				oldMultiplier := 1 + 0.02*float64(oldStacks)
				newMultiplier := 1 + 0.02*float64(newStacks)
				delta := newMultiplier / oldMultiplier

				for _, unit := range undeadTargets {
					for _, at := range aura.Unit.AttackTables[unit.UnitIndex] {
						at.DamageDealtMultiplier *= delta
						at.CritMultiplier *= delta
					}
				}
			},
		}))
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target.MobType != proto.MobType_MobTypeUndead {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeOrMeleeProc) {
				for _, aura := range buffAuras {
					aura.Activate(sim)
					aura.AddStack(sim)
				}
			} else if spell.ProcMask.Matches(core.ProcMaskRangedOrRangedProc) {
				// Ranged attacks actually remove a stack
				for _, aura := range buffAuras {
					if aura.IsActive() {
						aura.Deactivate(sim)
					}
				}
			}
		},
	}))
}

var ItemSetCryptstalkerArmor = core.NewItemSet(core.ItemSet{
	Name: "Cryptstalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasRanged4PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasRanged2PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasRanged6PBonus()
		},
	},
})

// Your Serpent Sting deals 20% more damage.
func (hunter *Hunter) applyNaxxramasRanged2PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Ranged 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterSerpentSting | ClassSpellMask_HunterSoFSerpentSting | ClassSpellMask_HunterChimeraSerpent,
		IntValue:  20,
	}))
}

// Reduces the cooldown on your Chimera Shot, Explosive Shot, and Aimed Shot abilities by 1.5 sec and reduces the cooldown on your Kill Shot ability by 3sec.
func (hunter *Hunter) applyNaxxramasRanged4PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Ranged 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_HunterChimeraShot | ClassSpellMask_HunterExplosiveShot | ClassSpellMask_HunterAimedShot,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Millisecond * 1500,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_HunterKillShot,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 3,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_HunterMultiShot,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 4,
	})
}

// You gain 2% increased damage and critical damage done to Undead for 30 sec each time you hit an Undead enemy with a ranged attack, stacking up to 7 times.
func (hunter *Hunter) applyNaxxramasRanged6PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Ranged 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	undeadTargets := core.FilterSlice(hunter.Env.Encounter.TargetUnits, func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeUndead })

	units := []*core.Unit{&hunter.Unit}
	if hunter.pet != nil {
		units = append(units, &hunter.pet.Unit)
	}

	buffAuras := []*core.Aura{}
	for _, unit := range units {
		buffAuras = append(buffAuras, unit.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 1218587},
			Label:     "Undead Slaying",
			Duration:  time.Second * 30,
			MaxStacks: 7,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				oldMultiplier := 1 + 0.02*float64(oldStacks)
				newMultiplier := 1 + 0.02*float64(newStacks)
				delta := newMultiplier / oldMultiplier

				for _, unit := range undeadTargets {
					for _, at := range aura.Unit.AttackTables[unit.UnitIndex] {
						at.DamageDealtMultiplier *= delta
						at.CritMultiplier *= delta
					}
				}
			},
		}))
	}

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target.MobType == proto.MobType_MobTypeUndead && spell.ProcMask.Matches(core.ProcMaskRangedOrRangedProc) {
				for _, aura := range buffAuras {
					aura.Activate(sim)
					aura.AddStack(sim)
				}
			}
		},
	}))
}
