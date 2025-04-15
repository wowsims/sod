package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetTheEarthshatterersStorm = core.NewItemSet(core.ItemSet{
	Name: "The Earthshatterer's Storm",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasElemental2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasElemental4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasElemental6PBonus()
		},
	},
})

const EleTier32pFlameShockDamageBonus = 40

// Increases periodic damage done by your Flame Shock ability by 40%.
func (shaman *Shaman) applyNaxxramasElemental2PBonus() {
	label := "S03 - Item - Naxxramas - Shaman - Elemental 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanFlameShock,
		Kind:      core.SpellMod_PeriodicDamageDone_Flat,
		IntValue:  EleTier32pFlameShockDamageBonus,
	}))
}

// Reduces the cooldown on your Lava Burst ability by 2 sec.
func (shaman *Shaman) applyNaxxramasElemental4PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
		return
	}

	label := "S03 - Item - Naxxramas - Shaman - Elemental 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_ShamanLavaBurst,
		TimeValue: -time.Second * 2,
	}))
}

// You gain 6% increased damage done to Undead for 30 sec for each time your Overload triggers, stacking up to 7 times.
func (shaman *Shaman) applyNaxxramasElemental6PBonus() {
	label := "S03 - Item - Naxxramas - Shaman - Elemental 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	undeadTargets := core.FilterSlice(shaman.Env.Encounter.TargetUnits, func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeUndead })

	buffAura := shaman.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1219370},
		Label:     "Undead Slaying",
		Duration:  time.Second * 30,
		MaxStacks: 6,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			oldMultiplier := 1 + 0.03*float64(oldStacks)
			newMultiplier := 1 + 0.03*float64(newStacks)
			delta := newMultiplier / oldMultiplier

			for _, unit := range undeadTargets {
				for _, at := range aura.Unit.AttackTables[unit.UnitIndex] {
					at.DamageDealtMultiplier *= delta
					at.CritMultiplier *= delta
				}
			}
		},
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ActionID.Tag == CastTagOverload {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetTheEarthshatterersRage = core.NewItemSet(core.ItemSet{
	Name: "The Earthshatterer's Rage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasEnhancement2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasEnhancement4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasEnhancement6PBonus()
		},
	},
})

// Increases damage done by your Lightning Shield by 100%.
func (shaman *Shaman) applyNaxxramasEnhancement2PBonus() {
	// Hotfix April 14, 2025: The T3 Enhancement Shaman 2-piece bonus now requires the player to have the  Static Shock rune engraved. The tooltip will be updated at a later date.
	if !shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock) {
		return
	}

	label := "S03 - Item - Naxxramas - Shaman - Enhancement 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_ShamanLightningShieldProc | ClassSpellMask_ShamanRollingThunder,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  100,
	}))
}

// Reduces the cooldown on your Lava Lash and Stormstrike abilities by 1.5 sec.
func (shaman *Shaman) applyNaxxramasEnhancement4PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaLash) && !shaman.Talents.Stormstrike {
		return
	}

	label := "S03 - Item - Naxxramas - Shaman - Enhancement 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_ShamanLavaLash | ClassSpellMask_ShamanStormstrike,
		TimeValue: -time.Millisecond * 1500,
	}))
}

// You gain 2% increased damage done to Undead for 30 sec for each charge of Maelstrom Weapon you earn, stacking up to 20 times.
func (shaman *Shaman) applyNaxxramasEnhancement6PBonus() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	label := "S03 - Item - Naxxramas - Shaman - Enhancement 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	undeadTargets := core.FilterSlice(shaman.Env.Encounter.TargetUnits, func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeUndead })

	buffAura := shaman.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1219370},
		Label:     "Undead Slaying",
		Duration:  time.Second * 30,
		MaxStacks: 7,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			oldMultiplier := 1 + 0.03*float64(oldStacks)
			newMultiplier := 1 + 0.03*float64(newStacks)
			delta := newMultiplier / oldMultiplier

			for _, unit := range undeadTargets {
				for _, at := range aura.Unit.AttackTables[unit.UnitIndex] {
					at.DamageDealtMultiplier *= delta
					at.CritMultiplier *= delta
				}
			}
		},
	})

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldOnRefresh := shaman.MaelstromWeaponAura.OnRefresh
			shaman.MaelstromWeaponAura.OnRefresh = func(aura *core.Aura, sim *core.Simulation) {
				if oldOnRefresh != nil {
					oldOnRefresh(aura, sim)
				}

				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	})
}

var ItemSetTheEarthshatterersResolve = core.NewItemSet(core.ItemSet{
	Name: "The Earthshatterer's Resolve",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasTank2PBonus()
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasTank4PBonus()
		},
		6: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.applyNaxxramasTank6PBonus()
		},
	},
})

// Your Earth Shock ability never misses when used as a taunt, and your chance to be Dodged or Parried is reduced by 2%.
func (shaman *Shaman) applyNaxxramasTank2PBonus() {
	label := "S03 - Item - Naxxramas - Shaman - Tank 2P Bonus"
	if shaman.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 2 * core.ExpertiseRatingPerExpertiseChance}

	hasRune := shaman.HasRune(proto.ShamanRune_RuneLegsWayOfEarth)
	hitMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusHit_Flat,
		ClassMask:  ClassSpellMask_ShamanEarthShock,
		FloatValue: 100 * core.SpellHitRatingPerHitChance,
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hasRune {
				hitMod.Activate()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hitMod.Deactivate()
		},
	}).AttachStatsBuff(bonusStats))
}

// Increases the damage taken reduction from your Shamanistic Rage ability by an additional 15% and during Shamanistic Rage your attack speed and spellcasting speed are increased by 30%.
func (shaman *Shaman) applyNaxxramasTank4PBonus() {
	label := "S03 - Item - Naxxramas - Shaman - Tank 4P Bonus"
	if shaman.HasAura(label) {
		return
	}

	shamRageDRBonus := 0.15
	attackCastSpeedBonus := 1.30

	shaman.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			shaman.shamanisticRageDRMultiplier += shamRageDRBonus
			shaman.ShamanisticRageAura.AttachMultiplyAttackSpeed(&shaman.Unit, attackCastSpeedBonus)
			shaman.ShamanisticRageAura.AttachMultiplyCastSpeed(&shaman.Unit, attackCastSpeedBonus)
		},
	})
}

// You take 20% reduced damage from Undead enemies.
func (shaman *Shaman) applyNaxxramasTank6PBonus() {
	label := "S03 - Item - Naxxramas - Shaman - Tank 6P Bonus"
	if shaman.HasAura(label) {
		return
	}

	damageTakenMultiplier := 0.80

	undeadTargets := core.FilterSlice(shaman.Env.Encounter.TargetUnits, func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeUndead })

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, target := range undeadTargets {
				for _, at := range target.AttackTables[shaman.UnitIndex] {
					at.DamageDealtMultiplier *= damageTakenMultiplier
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, target := range undeadTargets {
				for _, at := range target.AttackTables[shaman.UnitIndex] {
					at.DamageDealtMultiplier /= damageTakenMultiplier
				}
			}
		},
	}))
}

var ItemSetTheEarthshatteres = core.NewItemSet(core.ItemSet{
	Name: "The Earthshatterer",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Earth Shield ability no longer loses charges.
		2: func(agent core.Agent) {
		},
		// Your Healing Wave Rank 9 and Rank 10 and Lesser Healing Wave Rank 6 spells have a 10% chance to imbue your target with Totemic Power.
		4: func(agent core.Agent) {
		},
		// The target of your Spirit of the Alpha ability takes 20% reduced damage from Undead enemies.
		6: func(agent core.Agent) {
		},
	},
})
