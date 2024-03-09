package sod

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	// Proc effects. Keep these in order by item ID.

	// Fiery War Axe
	itemhelpers.CreateWeaponProcSpell(870, "Fiery War Axe", 1.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 18796},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Fiery War Axe Fireball",
				},
				TickLength:    2 * time.Second,
				NumberOfTicks: 3,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 8
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotCritChance = 0
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(155, 197)
				result := spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(1982, "Nightblade", 1.0, 18211, core.SpellSchoolShadow, 125, 150, 0, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcDamage(2164, "Gut Ripper", 1.0, 18107, core.SpellSchoolPhysical, 95, 26, 0, core.DefenseTypeMelee)

	itemhelpers.CreateWeaponProcAura(7717, "Ravager", 1.0, func(character *core.Character) *core.Aura {
		tickActionID := core.ActionID{SpellID: 9633}
		procActionID := core.ActionID{SpellID: 9632}
		auraActionID := core.ActionID{SpellID: 433801}

		ravegerBladestormTickSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         tickActionID,
			SpellSchool:      core.SpellSchoolPhysical,
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				damage := 5.0 +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			},
		})

		character.GetOrRegisterSpell(core.SpellConfig{
			SpellSchool: core.SpellSchoolPhysical,
			ActionID:    procActionID,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagChanneled,
			Dot: core.DotConfig{
				IsAOE: true,
				Aura: core.Aura{
					Label: "Ravager Whirlwind",
				},
				NumberOfTicks:       3,
				TickLength:          time.Second * 3,
				AffectedByCastSpeed: false,
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					ravegerBladestormTickSpell.Cast(sim, target)
				},
			},
		})

		return character.GetOrRegisterAura(core.Aura{
			Label:    "Ravager Bladestorm",
			ActionID: auraActionID,
			Duration: time.Second * 9,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AutoAttacks.CancelAutoSwing(sim)
				dotSpell := character.GetSpell(procActionID)
				dotSpell.AOEDot().Apply(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AutoAttacks.EnableAutoSwing(sim)
				dotSpell := character.GetSpell(procActionID)
				dotSpell.AOEDot().Cancel(sim)
			},
		})
	})

	itemhelpers.CreateWeaponProcAura(9423, "The Jackhammer", 1.0, func(character *core.Character) *core.Aura {
		return character.GetOrRegisterAura(core.Aura{
			Label:    "The Jackhammer Haste Aura",
			ActionID: core.ActionID{SpellID: 13533},
			Duration: time.Second * 9,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.3)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1/1.3)
			},
		})
	})

	itemhelpers.CreateWeaponProcDamage(9425, "Pendulum of Doom", 0.5, 10373, core.SpellSchoolPhysical, 250, 100, 0, core.DefenseTypeMelee)

	// Pip's Skinner
	core.NewItemEffect(12709, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.AddStat(stats.AttackPower, 45)
		}
	})

	// Fiendish Machete
	core.NewItemEffect(18310, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeElemental {
			character.AddStat(stats.AttackPower, 36)
		}
	})

	//Thunderfury
	core.NewItemEffect(19019, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(19019)
		ppmm := character.AutoAttacks.NewPPMManager(6.0, procMask)

		procActionID := core.ActionID{SpellID: 21992}

		singleTargetSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(1),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 0.5,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)
			},
		})

		makeDebuffAura := func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:    "Thunderfury",
				ActionID: procActionID,
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, -25)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, 25)
				},
			})
		}

		numHits := min(5, character.Env.GetNumTargets())
		debuffAuras := make([]*core.Aura, len(character.Env.Encounter.TargetUnits))
		for i, target := range character.Env.Encounter.TargetUnits {
			debuffAuras[i] = makeDebuffAura(target)
		}

		bounceSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(2),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			ThreatMultiplier: 1,
			FlatThreatBonus:  63,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					result := spell.CalcDamage(sim, curTarget, 0, spell.OutcomeMagicHit)
					if result.Landed() {
						debuffAuras[target.Index].Activate(sim)
					}
					spell.DealDamage(sim, result)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "Thunderfury",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Thunderfury") {
					singleTargetSpell.Cast(sim, result.Target)
					bounceSpell.Cast(sim, result.Target)
				}
			},
		})
	})

	// Mark of the Champion
	core.NewItemEffect(23206, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead || character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			character.AddStat(stats.AttackPower, 150)
		}
	})

	core.AddEffectsToTest = true
}
