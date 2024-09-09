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
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, blockBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, -blockBonus)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidBlock() {
				aura.RemoveStack(sim)
			}
		},
	})

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
		ActionID:   procID,
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeCrit,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AutoAttacks.ExtraMHAttack(sim, 1, procID, spell.ActionID)
		},
	})
}

func (paladin *Paladin) getWeaponSpecializationModifier() float64 {
	switch paladin.MainHand().HandType {
	case proto.HandType_HandTypeOneHand:
		return 1 + 0.02*float64(paladin.Talents.OneHandedWeaponSpecialization)
	case proto.HandType_HandTypeTwoHand:
		return 1 + 0.02*float64(paladin.Talents.TwoHandedWeaponSpecialization)
	default:
		return 1
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

	vengeanceMultiplier := []float64{1, 1.03, 1.06, 1.09, 1.12, 1.15}[paladin.Talents.Vengeance]

	procAura := paladin.RegisterAura(core.Aura{
		Label:    "Vengeance Proc",
		ActionID: core.ActionID{SpellID: 20059},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= vengeanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= vengeanceMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= vengeanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= vengeanceMultiplier
		},
	})

	paladin.RegisterAura(core.Aura{
		Label:    "Vengeance",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				procAura.Activate(sim)
			}
		},
	})
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
