package sod

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	WolfsheadTrophy = 7124
	BloodPlague     = 7878
	FrostFever      = 7879
	MarkOfBlood     = 7880
	Obliterate      = 7881
	GrandCrusader   = 7940
	GrandArcanist   = 7941
	GrandSorceror   = 7942
	GrandInquisitor = 7943

	BloodPlagueSpellId = 1219121
	FrostFeverSpellId  = 1219124
	MarkOfBloodSpellId = 1219153
	ObliterateSpellId  = 1219176
)

func init() {
	core.AddEffectsToTest = false

	// Weapon - Dismantle
	core.NewEnchantEffect(7210, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 439164},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(60, 90), spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Enchant Weapon - Dismantle",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// Dismantle only procs on hits that land
				if !result.Landed() {
					return
				}

				// Dismantle only procs on Mechanical units
				if result.Target.MobType != proto.MobType_MobTypeMechanical {
					return
				}

				// Dismantle only procs on direct attacks, not proc effects or DoT ticks
				if !spell.Flags.Matches(core.SpellFlagNotAProc) && spell.ProcMask.Matches(core.ProcMaskProc|core.ProcMaskSpellDamageProc) {
					return
				}

				// TODO: Confirm: Dismantle can not proc itself
				if spell == procSpell {
					return
				}

				// Main-Hand hits only trigger Dismantle if the MH weapon is enchanted with Dismantle
				if core.ProcMaskMeleeMH.Matches(spell.ProcMask) && (character.GetMHWeapon() == nil || character.GetMHWeapon().Enchant.EffectID != 7210) {
					return
				}

				// Off-Hand hits only trigger Dismantle if the MH weapon is enchanted with Dismantle
				if core.ProcMaskMeleeOH.Matches(spell.ProcMask) && (character.GetOHWeapon() == nil || character.GetOHWeapon().Enchant.EffectID != 7210) {
					return
				}

				if spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
					if sim.RandomFloat("Dismantle") < 0.10 {
						// Spells proc both Main-Hand and Off-Hand if both are enchanted
						if character.GetMHWeapon() != nil && character.GetMHWeapon().Enchant.EffectID == 7210 {
							procSpell.Cast(sim, result.Target)
						}
						if character.GetOHWeapon() != nil && character.GetOHWeapon().Enchant.EffectID == 7210 {
							procSpell.Cast(sim, result.Target)
						}
					}
				} else if sim.RandomFloat("Dismantle") < 0.10 {
					// Physical hits only proc on the hand that was hit with
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(7210, aura)
	})

	// Sharpened Chitin Armor Kit
	// Permanently cause an item worn on the chest, legs, hands or feet to cause 20 Nature damage to the attacker when struck in combat.
	// Only usable on items level 45 and above.
	core.NewEnchantEffect(7649, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 233803}

		damage := 20.0
		numEnchants := 0
		for _, item := range character.Equipment {
			if item.Enchant.EffectID == 7649 {
				numEnchants++
			}
		}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHit)
			},
		})

		procAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Thorns +20",
			Callback: core.CallbackOnSpellHitTaken,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				for i := 0; i < numEnchants; i++ {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		}).AttachAdditivePseudoStatBuff(&character.PseudoStats.ThornsDamage, damage*float64(numEnchants))

		character.ItemSwap.RegisterEnchantProc(7649, procAura)
	})

	// Obsidian Scope
	core.AddWeaponEffect(7657, func(agent core.Agent) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 10
		w.BaseDamageMax += 10
	})

	registerDeathKnightDiseaseSpell("Blood Plague", BloodPlague, BloodPlagueSpellId, 120, 5)
	registerDeathKnightDiseaseSpell("Frost Fever", FrostFever, FrostFeverSpellId, 100, 7)

	core.NewEnchantEffect(Obliterate, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: ObliterateSpellId}
		baseDamage := 900.0
		singleDiseaseMultiplier := 1 + 0.6
		doubleDiseaseMultiplier := 1 + 0.6 + 0.6 // Assuming additive until it can be tested by 3 people on the same target dummy

		character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 90,
				},
			},

			DamageMultiplier: 1.0,
			ThreatMultiplier: 1.0,
			BonusCoefficient: 0, // Not affected by things like Gift of Arthas

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return character.HasMHWeapon()
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				multiplier := 1.0
				diseases := target.GetAurasWithTag("Obliterate")

				switch len(diseases) {
				case 2:
					multiplier = doubleDiseaseMultiplier
				case 1:
					multiplier = singleDiseaseMultiplier
				}

				spell.ApplyMultiplicativeDamageBonus(multiplier)
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				spell.ApplyMultiplicativeDamageBonus(1 / multiplier)

				spell.DealDamage(sim, result)
			},
		})

		character.ItemSwap.RegisterEnchantActive(Obliterate, ObliterateSpellId)
	})

	core.NewEnchantEffect(GrandCrusader, func(agent core.Agent) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForEnchant(GrandCrusader, 1.0, 0)

		spellID := int32(1231124)
		strBonus := 120.0
		duration := time.Second * 20
		mhAura := character.NewTemporaryStatsAuraWrapped("Righteous Strength MH", core.ActionID{SpellID: spellID, Tag: 1}, stats.Stats{stats.Strength: strBonus}, duration, func(aura *core.Aura) {
			aura.Tag = "Crusader"
		})
		ohAura := character.NewTemporaryStatsAuraWrapped("Righteous Strength OH", core.ActionID{SpellID: spellID, Tag: 2}, stats.Stats{stats.Strength: strBonus}, duration, func(aura *core.Aura) {
			aura.Tag = "Crusader"
		})
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: spellID})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Grand Crusader",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			DPM:               dpm,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.IsMH() {
					mhAura.Activate(sim)
				} else {
					ohAura.Activate(sim)
				}
				character.GainHealth(sim, sim.RollWithLabel(350, 450, "Righteous Strength"), healthMetrics)
			},
		})

		character.ItemSwap.RegisterEnchantProc(GrandCrusader, aura)
	})

	core.NewEnchantEffect(GrandInquisitor, func(agent core.Agent) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForEnchant(GrandInquisitor, 1.0, 0)

		spellID := int32(1232169)
		righteousInquisitionAura := character.NewTemporaryStatsAuraWrapped("Righteous Inquisition", core.ActionID{SpellID: spellID}, stats.Stats{stats.Strength: 200}, time.Second*20, func(aura *core.Aura) {
			aura.Tag = "Crusader"
		})
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: spellID})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Grand Inquisitor",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			DPM:               dpm,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				righteousInquisitionAura.Activate(sim)
				character.GainHealth(sim, sim.RollWithLabel(350, 450, "Righteous Inquisition"), healthMetrics)
			},
		})

		character.ItemSwap.RegisterEnchantProc(GrandInquisitor, aura)
	})

	core.NewEnchantEffect(GrandSorceror, func(agent core.Agent) {
		character := agent.GetCharacter()

		spellID := int32(1231162)
		righteousBlastingAura := character.NewTemporaryStatsAura("Righteous Blasting", core.ActionID{SpellID: spellID}, stats.Stats{stats.SpellPower: 70}, time.Second*20)
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: spellID})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Grand Sorcerer",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.07,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				righteousBlastingAura.Activate(sim)
				if character.HasManaBar() {
					character.AddMana(sim, sim.RollWithLabel(350, 450, "Righteous Blasting"), manaMetrics)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(GrandSorceror, aura)
	})

	core.NewEnchantEffect(GrandArcanist, func(agent core.Agent) {
		character := agent.GetCharacter()

		righteousFireAura := character.NewTemporaryStatsAura("Righteous Fire", core.ActionID{SpellID: 1231138}, stats.Stats{stats.SpellPower: 140}, time.Second*20)
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 1231138})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Grand Arcanist",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.07,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				righteousFireAura.Activate(sim)
				if character.HasManaBar() {
					character.AddMana(sim, sim.RollWithLabel(350, 450, "Righteous Fire"), manaMetrics)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(GrandArcanist, aura)
	})

	core.AddEffectsToTest = true
}

func registerDeathKnightDiseaseSpell(label string, effectID int32, spellID int32, baseDamage float64, numTicks int32) {
	core.NewEnchantEffect(effectID, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellID}.WithTag(1),
			SpellSchool: core.SpellSchoolShadow, // For some reason, both Frost Fever and Blood Plague are Shadow, verified in logs
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagAPL | core.SpellFlagPureDot | core.SpellFlagDisease,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 1,
				},
			},

			DamageMultiplier: 1.0,
			ThreatMultiplier: 1.0,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: label + "-" + character.Label,
				},

				NumberOfTicks: numTicks,
				TickLength:    time.Second * 3,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, baseDamage, isRollover)
				},

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.Snapshot(target, baseDamage, false)
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
				if !result.Landed() {
					spell.DealOutcome(sim, result)
					return
				}
				spell.Dot(target).Apply(sim)
			},
		})

		character.ItemSwap.RegisterEnchantActive(effectID, spellID)
	})
}
