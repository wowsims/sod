package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (paladin *Paladin) ApplyTalents() {
	paladin.AddStat(stats.MeleeHit, float64(paladin.Talents.Precision)*core.MeleeHitRatingPerHitChance)
	// TODO: paladin.AddStat(stats.RangedHit, float64(paladin.Talents.Precision)*core.MeleeHitRatingPerHitChance)

	paladin.AddStat(stats.MeleeCrit, float64(paladin.Talents.Conviction)*core.CritRatingPerCritChance)
	// TODO: paladin.AddStat(stats.RangedCrit, float64(paladin.Talents.Conviction)*core.CritRatingPerCritChance)

	if paladin.Talents.Toughness > 0 {
		paladin.ApplyEquipScaling(stats.Armor, 1.0+0.02*float64(paladin.Talents.Toughness))
	}

	// These are no-op if untalented.
	paladin.MultiplyStat(stats.Strength, 1.0+0.02*float64(paladin.Talents.DivineStrength))
	paladin.MultiplyStat(stats.Intellect, 1.0+0.02*float64(paladin.Talents.DivineIntellect))
	paladin.AddStat(stats.Defense, 2*float64(paladin.Talents.Anticipation))

	// Shield Specialization bonus is additive. NOTE: Total SBV will be inflated until
	// https://github.com/wowsims/sod/issues/1025 gets resolved.
	paladin.PseudoStats.BlockValueMultiplier += 0.1 * float64(paladin.Talents.ShieldSpecialization)

	paladin.AddStat(stats.Parry, 1*float64(paladin.Talents.Deflection))

	paladin.applyWeaponSpecialization()
	if paladin.Talents.Vengeance > 0 {
		paladin.applyVengeance()
	}
	if paladin.Talents.Vindication > 0 {
		paladin.applyVindication()
	}
	paladin.PseudoStats.SchoolBonusCritChance[stats.SchoolIndexHoly] += core.SpellCritRatingPerCritChance * float64(paladin.Talents.HolyPower)

	paladin.applyRedoubt()
	paladin.applyReckoning()
	paladin.applyImprovedLayOnHands()

	paladin.applyHealingLight()
}

func (paladin *Paladin) improvedSoR() float64 {
	return []float64{1, 1.03, 1.06, 1.09, 1.12, 1.15}[paladin.Talents.ImprovedSealOfRighteousness]
}

func (paladin *Paladin) benediction() int32 {
	return []int32{100, 97, 94, 91, 88, 85}[paladin.Talents.Benediction]
}

func (paladin *Paladin) applyRedoubt() {
	if paladin.Talents.Redoubt == 0 {
		return
	}

	// Redoubt grants 6% block chance per point.
	blockBonus := 6.0 * float64(paladin.Talents.Redoubt) * core.BlockRatingPerBlockChance

	paladin.redoubtAura = paladin.RegisterAura(core.Aura{
		Label:     "Redoubt",
		ActionID:  core.ActionID{SpellID: 20134},
		Duration:  time.Second * 10,
		MaxStacks: 5,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidBlock() {
				aura.RemoveStack(sim)
			}
		},
	}).AttachStatBuff(stats.Block, blockBonus)

	paladin.RegisterAura(core.Aura{
		Label:    "Redoubt Crit Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				paladin.redoubtAura.Activate(sim)
				paladin.redoubtAura.SetStacks(sim, 5)
			}
		},
	})
}

func (paladin *Paladin) applyReckoning() {

	if paladin.Talents.Reckoning == 0 {
		return
	}

	procID := core.ActionID{SpellID: 20178} // Reckoning Proc ID
	procChance := 0.2 * float64(paladin.Talents.Reckoning)

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       "Reckoning Crit Trigger",
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeCrit,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AutoAttacks.ExtraMHAttackProc(sim, 1, procID, spell)
		},
	})
}

func (paladin *Paladin) getWeaponSpecializationModifier() float64 {
	handType := paladin.MainHand().HandType
	if handType == proto.HandType_HandTypeMainHand || handType == proto.HandType_HandTypeOneHand {
		return 1. + 0.02*float64(paladin.Talents.OneHandedWeaponSpecialization)
	} else if handType == proto.HandType_HandTypeTwoHand {
		return 1. + 0.02*float64(paladin.Talents.TwoHandedWeaponSpecialization)
	} else {
		return 1.
	}
}

// Affects all physical damage or spells that can be rolled as physical.
func (paladin *Paladin) applyWeaponSpecialization() {
	paladin.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= paladin.getWeaponSpecializationModifier()
}

func (paladin *Paladin) applyVengeance() {
	if paladin.Talents.Vengeance == 0 {
		return
	}

	vengeanceMultiplier := []float64{0, 0.03, 0.06, 0.09, 0.12, 0.15}[paladin.Talents.Vengeance]
	damageMultiplier := 1.0 + vengeanceMultiplier

	if !paladin.Options.RighteousFury {
		threatMultiplier := 1.0 + (vengeanceMultiplier * 2)
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_Threat_Pct,
			School:     core.SpellSchoolPhysical | core.SpellSchoolHoly,
			FloatValue: 1.0 / threatMultiplier,
		})
	}

	procAura := paladin.RegisterAura(core.Aura{
		Label:    "Vengeance Proc",
		ActionID: core.ActionID{SpellID: 20059},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= damageMultiplier
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= damageMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= damageMultiplier
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= damageMultiplier
		},
	})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label: "Vengeance",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				procAura.Activate(sim)
			}
		},
	}))
}

func (paladin *Paladin) applyVindication() {
	if paladin.Talents.Vindication == 0 {
		return
	}
	//vindicationMultiplier := []float64{1, 1.05, 1.10, 1.15}[paladin.Talents.Vengeance]
	vindicationMultiplier := []*stats.StatDependency{
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.00),
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.05),
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.10),
		paladin.NewDynamicMultiplyStat(stats.AttackPower, 1.15),
	}

	vindicationAura := paladin.RegisterAura(core.Aura{
		Label:    "Vindication Proc",
		ActionID: core.ActionID{SpellID: 26021},
		Duration: time.Second * 30,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, vindicationMultiplier[0])
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.EnableDynamicStatDep(sim, vindicationMultiplier[paladin.Talents.Vindication])
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.DisableDynamicStatDep(sim, vindicationMultiplier[paladin.Talents.Vindication])
		},
	})
	// 	vindicationAuras := paladin.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
	// 		return core.VindicationAura(target, paladin.Talents.Vindication)
	// 	})
	paladin.RegisterAura(core.Aura{
		Label:    "Vindication Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// TODO: Replace with actual proc mask / proc chance
			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
				vindicationAura.Activate(sim)
			}
		},
	})
}

func (paladin *Paladin) applyImprovedLayOnHands() {

	if paladin.Talents.ImprovedLayOnHands > 0 {

		armorMultiplier := []float64{1, 1.15, 1.3}[paladin.Talents.ImprovedLayOnHands]
		auraID := []int32{0, 20233, 20236}[paladin.Talents.ImprovedLayOnHands]

		paladin.RegisterAura(core.Aura{
			Label:    "Lay on Hands",
			ActionID: core.ActionID{SpellID: auraID},
			Duration: time.Minute * 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				paladin.ApplyDynamicEquipScaling(sim, stats.Armor, armorMultiplier)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				paladin.RemoveDynamicEquipScaling(sim, stats.Armor, armorMultiplier)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(ClassSpellMask_PaladinLayOnHands) {
					aura.Activate(sim)
				}
			},
		})
	}
}

func (paladin *Paladin) applyHealingLight() {
	if paladin.Talents.HealingLight > 0 {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  ClassSpellMask_PaladinHolyLight,
			FloatValue: 1 + 0.04*float64(paladin.Talents.HealingLight),
		})
	}
}
