package druid

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/common/sod"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Totem Item IDs
const (
	WolfsheadHelm                    = 8345
	IdolOfFerocity                   = 22397
	IdolOfTheMoon                    = 23197
	IdolOfBrutality                  = 23198
	IdolMindExpandingMushroom        = 209576
	Catnip                           = 213407
	IdolOfWrath                      = 216490
	BloodBarkCrusher                 = 216499
	RitualistsHammer                 = 221446
	IdolOfTheDream                   = 220606
	IdolOfExsanguinationCat          = 228181
	IdolOfTheSwarm                   = 228180
	IdolOfExsanguinationBear         = 228182
	BloodGuardDragonhideGrips        = 227180
	KnightLieutenantsDragonhideGrips = 227183
)

func init() {
	core.AddEffectsToTest = false

	// https://www.wowhead.com/classic/item=22397/idol-of-ferocity
	// Equip: Reduces the energy cost of Claw and Rake by 3.
	core.NewItemEffect(IdolOfFerocity, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		// TODO: Claw is not implemented
		druid.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_DruidRake || spell.SpellCode == SpellCode_DruidMangleCat {
				spell.Cost.FlatModifier -= 3
			}
		})
	})

	// https://www.wowhead.com/classic/item=23197/idol-of-the-moon
	// Equip: Increases the damage of your Moonfire spell by up to 33.
	core.NewItemEffect(IdolOfTheMoon, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		affectedSpellCodes := []int32{SpellCode_DruidMoonfire, SpellCode_DruidSunfire, SpellCode_DruidStarfallSplash, SpellCode_DruidStarfallTick}
		druid.OnSpellRegistered(func(spell *core.Spell) {
			if slices.Contains(affectedSpellCodes, spell.SpellCode) {
				spell.BonusDamage += 33
			}
		})
	})

	// https://www.wowhead.com/classic/item=23198/idol-of-brutality
	// Equip: Reduces the rage cost of Maul and Swipe by 3.
	core.NewItemEffect(IdolOfBrutality, func(agent core.Agent) {
		// Implemented in maul.go and swipe.go
	})

	core.NewItemEffect(IdolMindExpandingMushroom, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.Spirit, 5)
	})

	// https://www.wowhead.com/classic/item=228181/idol-of-exsanguination-cat
	// Equip: The energy cost of your Rake and Rip spells is reduced by 5.
	core.NewItemEffect(IdolOfExsanguinationCat, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellCode == SpellCode_DruidRake || spell.SpellCode == SpellCode_DruidRip {
				spell.Cost.FlatModifier -= 5
			}
		})
	})

	// https://www.wowhead.com/classic/item=228182/idol-of-exsanguination-bear
	// Equip: Your Lacerate ticks energize you for 3 rage.
	core.NewItemEffect(IdolOfExsanguinationBear, func(agent core.Agent) {
		// TODO: Not yet implemented
	})

	// https://www.wowhead.com/classic/item=228180/idol-of-the-swarm
	// Equip: The duration of your Insect Swarm spell is increased by 12 sec.
	core.NewItemEffect(IdolOfTheSwarm, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		bonusDuration := time.Second * 12

		core.MakePermanent(druid.GetOrRegisterAura(core.Aura{
			Label: "Idol of the Swarm",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				for _, spell := range druid.InsectSwarm {
					if spell != nil {
						for _, dot := range spell.Dots() {
							if dot != nil {
								dot.NumberOfTicks += 6
								dot.RecomputeAuraDuration()
							}
						}
					}
				}

				for _, aura := range druid.InsectSwarmAuras {
					if aura != nil && !aura.IsPermanent() {
						aura.Duration += bonusDuration
					}
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				for _, spell := range druid.InsectSwarm {
					if spell != nil {
						for _, dot := range spell.Dots() {
							if dot != nil {
								dot.NumberOfTicks -= 6
								dot.RecomputeAuraDuration()
							}
						}
					}
				}

				for _, aura := range druid.InsectSwarmAuras {
					if aura != nil && !aura.IsPermanent() {
						aura.Duration -= bonusDuration
					}
				}
			},
		}))
	})

	core.NewItemEffect(BloodBarkCrusher, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.newBloodbarkCleaveItem(BloodBarkCrusher)
	})

	core.NewItemEffect(RitualistsHammer, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		druid.newBloodbarkCleaveItem(RitualistsHammer)
	})

	// https://www.wowhead.com/classic/item=227180/blood-guards-dragonhide-grips
	// Equip: Reduces the mana cost of your shapeshifts by 150.
	core.NewItemEffect(BloodGuardDragonhideGrips, func(agent core.Agent) {
		registerDragonHideGripsAura(agent.(DruidAgent).GetDruid())
	})

	// https://www.wowhead.com/classic/item=227183/knight-lieutenants-dragonhide-grips
	// Equip: Reduces the mana cost of your shapeshifts by 150.
	core.NewItemEffect(KnightLieutenantsDragonhideGrips, func(agent core.Agent) {
		registerDragonHideGripsAura(agent.(DruidAgent).GetDruid())
	})

	core.AddEffectsToTest = true
}

// https://www.wowhead.com/classic/item=213407/catnip
func (druid *Druid) registerCatnipCD() {
	if druid.Consumes.MiscConsumes == nil || !druid.Consumes.MiscConsumes.Catnip {
		return
	}
	sod.RegisterFiftyPercentHasteBuffCD(&druid.Character, core.ActionID{ItemID: Catnip})
}

func (druid *Druid) newBloodbarkCleaveItem(itemID int32) {
	auraActionID := core.ActionID{SpellID: 436482}

	results := make([]*core.SpellResult, min(3, druid.Env.GetNumTargets()))

	damageSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 436481},
		SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMelee, // actually has DefenseTypeNone, but is likely using the greatest CritMultiplier available
		ProcMask:    core.ProcMaskEmpty,

		// TODO: "Causes additional threat" in Tooltip, no clue what the multiplier is.
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				results[idx] = spell.CalcDamage(sim, target, 5, spell.OutcomeMagicCrit)
				target = sim.Environment.NextTargetUnit(target)
			}
			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})

	buffAura := druid.GetOrRegisterAura(core.Aura{
		Label:    "Bloodbark Cleave",
		ActionID: auraActionID,
		Duration: 20 * time.Second,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ProcMask&core.ProcMaskMelee != 0 {
				damageSpell.Cast(sim, result.Target)
				return
			}
		},
	})

	mainSpell := druid.GetOrRegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{ItemID: itemID},
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell:    mainSpell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}

func registerDragonHideGripsAura(druid *Druid) {
	const costReduction int32 = 150
	var affectedForms []*DruidSpell

	druid.RegisterAura(core.Aura{
		Label:    "Dragonhide Grips",
		ActionID: core.ActionID{SpellID: 459594},
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedForms = []*DruidSpell{
				druid.CatForm,
				druid.MoonkinForm,
				druid.BearForm,
			}
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedForms {
				if spell != nil {
					spell.Cost.FlatModifier -= costReduction
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedForms {
				if spell != nil {
					spell.Cost.FlatModifier += costReduction
				}
			}
		},
	})
}
