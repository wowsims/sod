package item_effects

import (
	"time"

	"github.com/wowsims/sod/sim/common/itemhelpers"
	"github.com/wowsims/sod/sim/common/sod"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const (
	AutomaticCrowdPummeler              = 210741
	ElectromagneticHyperfluxReactivator = 213281
	ElectrocutionersNeedle              = 213286
	ToxicRevengerTwo                    = 213291
	SuperchargedHeadchopper             = 213296
	MachinistsGloves                    = 213319
	MiniaturizedCombustionChamber       = 213347
	Shawarmageddon                      = 213105
	MekkatorquesArcanoShredder          = 213409
	GyromaticExperiment420b             = 213348
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	// Electromagnetic Hyperflux Reactivator
	core.NewItemEffect(ElectromagneticHyperfluxReactivator, func(agent core.Agent) {
		character := agent.GetCharacter()

		forkedLightning := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 11828},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					baseDamage := sim.Roll(153, 173)
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		dmgShieldActionID := core.ActionID{SpellID: 11841}

		dmgShieldProc := character.RegisterSpell(core.SpellConfig{
			ActionID:    dmgShieldActionID,
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 5, spell.OutcomeMagicHit)
			},
		})

		dmgShieldAura := character.RegisterAura(core.Aura{
			Label:    "Static Barrier",
			ActionID: dmgShieldActionID,
			Duration: time.Minute * 10,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					dmgShieldProc.Cast(sim, spell.Unit)
				}
			},
		})

		activationSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 11826},
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 30,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmgShieldAura.Activate(sim)
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + 3*time.Second,
					OnAction: func(s *core.Simulation) {
						forkedLightning.Cast(sim, target)
					},
				})
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
			ShouldActivate: func(_ *core.Simulation, _ *core.Character) bool {
				// Only castable with manual APL Action
				return false
			},
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	// https://www.wowhead.com/classic/item=213319/machinists-gloves
	// Equip: +30 Attack Power when fighting Mechanical units.
	core.NewMobTypeAttackPowerEffect(MachinistsGloves, []proto.MobType{proto.MobType_MobTypeMechanical}, 30)

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Plate
	///////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// Miniaturized Combustion Chamber
	core.NewItemEffect(MiniaturizedCombustionChamber, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 435167}
		manaMetrics := character.NewManaMetrics(actionID)
		healthMetrics := character.NewHealthMetrics(actionID)

		manaRoll := 0.0
		dmgRoll := 0.0

		regChannel := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagChanneled,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 30,
				},
			},

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Miniaturized Combustion Chamber",
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						character.AutoAttacks.CancelAutoSwing(sim)
						manaRoll = sim.RollWithLabel(1, 150, "Miniaturized Combustion Chamber")
						dmgRoll = sim.RollWithLabel(1, 150, "Miniaturized Combustion Chamber")
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						character.AutoAttacks.EnableAutoSwing(sim)
					},
				},
				SelfOnly:      true,
				NumberOfTicks: 10,
				TickLength:    time.Second,

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					character.AddMana(sim, manaRoll, manaMetrics)
					character.RemoveHealth(sim, dmgRoll, healthMetrics)
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				spell.SelfHot().Apply(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    regChannel,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeMana,
			ShouldActivate: func(_ *core.Simulation, _ *core.Character) bool {
				// Only castable with manual APL Action
				return false
			},
		})
	})

	core.NewItemEffect(GyromaticExperiment420b, func(agent core.Agent) {
		character := agent.GetCharacter()

		hasteAura := character.RegisterAura(core.Aura{
			Label:    "Gyromatic Experiment 420b",
			ActionID: core.ActionID{SpellID: 435899},
			Duration: time.Second * 20,
		}).AttachMultiplyAttackSpeed(&character.Unit, 1.05)

		chickenAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Cluck Cluck??",
			ActionID: core.ActionID{SpellID: 435896},
			Duration: time.Second * 5,
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 435899},
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 30,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				if sim.RandomFloat("Gyromatic Experiment 420b") > 0.95 {
					chickenAura.Activate(sim)
					character.WaitUntil(sim, chickenAura.ExpiresAt())
					character.AutoAttacks.DelayMeleeBy(sim, time.Second*5)
				} else {
					hasteAura.Activate(sim)
				}
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Weapons
	///////////////////////////////////////////////////////////////////////////

	// Automatic Crowd Pummeler
	core.NewItemEffect(AutomaticCrowdPummeler, func(agent core.Agent) {
		character := agent.GetCharacter()
		sod.RegisterFiftyPercentHasteBuffCD(character, core.ActionID{ItemID: AutomaticCrowdPummeler})
	})

	itemhelpers.CreateWeaponCoHProcDamage(ElectrocutionersNeedle, "Electrocutioner's Needle", 6.5, 434839, core.SpellSchoolNature, 25, 10, 0.05, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponCoHProcDamage(SuperchargedHeadchopper, "Supercharged Headchopper", 1.5, 434842, core.SpellSchoolNature, 80, 20, 0.1, core.DefenseTypeMagic)

	itemhelpers.CreateWeaponProcSpell(ToxicRevengerTwo, "Toxic Revenger II", 3.0, func(character *core.Character) *core.Spell {
		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 435169},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Toxic Revenger II Poison Cloud",
				},
				TickLength:    5 * time.Second,
				NumberOfTicks: 3,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 30, isRollover)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					result := spell.CalcOutcome(sim, aoeTarget, spell.OutcomeMagicHit)
					if result.Landed() {
						spell.Dot(aoeTarget).Apply(sim)
					}
				}
			},
		})
	})

	core.NewItemEffect(Shawarmageddon, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{SpellID: 434488}

		fireStrike := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 434488},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 7.0, spell.OutcomeMagicHitAndCrit)
			},
		})

		spicyAura := character.RegisterAura(core.Aura{
			Label:    "Spicy!",
			ActionID: actionID,
			Duration: time.Second * 30,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if result.Landed() {
					fireStrike.Cast(sim, spell.Unit.CurrentTarget)
				}
			},
		}).AttachMultiplyAttackSpeed(&character.Unit, 1.04)

		spicy := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Cast: core.CastConfig{
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				spicyAura.Activate(sim)
			},
		})

		character.AddMajorEquipmentCooldown(core.MajorCooldown{
			Spell: spicy,
			Type:  core.CooldownTypeDPS,
		})
	})

	// Mekkatorque's Arcano-Shredder
	itemhelpers.CreateWeaponProcSpell(MekkatorquesArcanoShredder, "Mekkatorque", 5.0, func(character *core.Character) *core.Spell {
		procAuras := character.NewEnemyAuraArray(core.MekkatorqueFistDebuffAura)

		return character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 434841},
			SpellSchool: core.SpellSchoolArcane,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.05,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				if target.Level <= 45 {
					spell.CalcAndDealDamage(sim, target, 30, spell.OutcomeMagicHitAndCrit)
					procAuras.Get(target).Activate(sim)
				}
			},
		})
	})

	core.AddEffectsToTest = true
}
