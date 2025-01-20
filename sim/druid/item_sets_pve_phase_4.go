package druid

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"time"
)

var ItemSetFeralheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Feralheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
				stats.SpellDamage:       23,
				stats.HealingPower:      44,
			})
		},
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450608}
			manaMetrics := c.NewManaMetrics(actionID)
			energyMetrics := c.NewEnergyMetrics(actionID)
			rageMetrics := c.NewRageMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Mana)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.02,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					c.AddMana(sim, 300, manaMetrics)
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Energy)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				ProcChance: 0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 40, energyMetrics)
					}
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:       "S03 - Druid Energize Trigger - Wildheart Raiment (Rage)",
				Callback:   core.CallbackOnSpellHitTaken,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.03,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasRageBar() {
						c.AddRage(sim, 10, rageMetrics)
					}
				},
			})
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
	},
})

var ItemSetCenarionEclipse = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Eclipse",
	Bonuses: map[int32]core.ApplyEffect{
		// Damage dealt by Thorns increased by 100% and duration increased by 200%.
		2: func(agent core.Agent) {
			// TODO: Thorns
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT1Balance4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT1Balance6PBonus()
		},
	},
})

// Increases your chance to hit with spells and attacks by 3%.
func (druid *Druid) applyT1Balance4PBonus() {
	label := "S03 - Item - T1 - Druid - Balance 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{
		stats.MeleeHit: 3 * core.MeleeHitRatingPerHitChance,
		stats.SpellHit: 3 * core.SpellHitRatingPerHitChance,
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(bonusStats)
			} else {
				aura.Unit.AddStatsDynamic(sim, bonusStats)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if aura.Unit.Env.MeasuringStats && aura.Unit.Env.State != core.Finalized {
				aura.Unit.AddStats(bonusStats.Invert())
			} else {
				aura.Unit.AddStatsDynamic(sim, bonusStats.Invert())
			}
		},
	}))
}

// Reduces the cooldown on Starfall by 50%.
func (druid *Druid) applyT1Balance6PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneCloakStarfall) {
		return
	}

	label := "S03 - Item - T1 - Druid - Balance 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.Starfall.CD.Multiplier -= 50
		},
	})
}

var ItemSetCenarionCunning = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Cunning",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT1Feral2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT1Feral4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT1Feral6PBonus()
		},
	},
})

// Your Faerie Fire and Faerie Fire (Feral) also increase the chance for all attacks to hit that target by 1% for 40 sec.
func (druid *Druid) applyT1Feral2PBonus() {
	label := "S03 - Item - T1 - Druid - Feral 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.ImprovedFaerieFireAuras = druid.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return core.ImprovedFaerieFireAura(target)
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.SpellCode == SpellCode_DruidFaerieFire || spell.SpellCode == SpellCode_DruidFaerieFireFeral) && result.Landed() {
				druid.ImprovedFaerieFireAuras.Get(result.Target).Activate(sim)
			}
		},
	}))
}

// Periodic damage from your Rake and Rip can now be critical strikes.
func (druid *Druid) applyT1Feral4PBonus() {
	label := "S03 - Item - T1 - Druid - Feral 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 455872}, // Tracking in APL
		Label:    label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			druid.AllowRakeRipDoTCrits = true
		},
	}))
}

// Your Rip and Ferocious Bite have a 20% chance per combo point spent to refresh the duration of Savage Roar back to its initial value.
func (druid *Druid) applyT1Feral6PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsSavageRoar) {
		return
	}

	label := "S03 - Item - T1 - Druid - Feral 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 455873}, // Tracking in APL
		Label:    label,
	}))

	druid.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
		if spell != druid.SavageRoar.Spell && druid.SavageRoarAura.IsActive() && sim.Proc(.2*float64(comboPoints), label) {
			druid.SavageRoarAura.Refresh(sim)
		}
	})
}

var ItemSetCenarionRage = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Rage",
	Bonuses: map[int32]core.ApplyEffect{
		// You may cast Rebirth and Innervate while in Bear Form or Dire Bear Form.
		2: func(agent core.Agent) {
			// Nothing to do
		},
		// Reduces the cooldown of Enrage by 30 sec and it no longer reduces your armor.
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.Enrage.CD.FlatModifier -= time.Second * 30
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT1Guardian6PBonus()
		},
	},
})

// Bear Form and Dire Bear Form increase all threat you generate by an additional 20%, and Cower now removes all your threat against the target but has a 20 sec longer cooldown.
func (druid *Druid) applyT1Guardian6PBonus() {
	druid.CenarionRageThreatBonus = .2
}

var ItemSetCenarionBounty = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Bounty",
	Bonuses: map[int32]core.ApplyEffect{
		// When you cast Innervate on another player, it is also cast on you.
		2: func(agent core.Agent) {
			// TODO: Would need to rework innervate to make this work
		},
		// Casting your Healing Touch or Nourish spells gives you a 25% chance to gain Mana equal to 35% of the base cost of the spell.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// Reduces the cooldown on Tranquility by 100% and increases its healing by 100%.
		6: func(agent core.Agent) {
			// Nothing to do
		},
	},
})
