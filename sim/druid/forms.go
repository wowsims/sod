package druid

import (
	"github.com/wowsims/sod/sim/common/sod"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type DruidForm uint8

const (
	Humanoid DruidForm = 1 << iota
	Bear
	Cat
	Moonkin
	Tree
	Any = Humanoid | Bear | Cat | Moonkin | Tree
)

func (form DruidForm) Matches(other DruidForm) bool {
	return (form & other) != 0
}

func (druid *Druid) GetForm() DruidForm {
	return druid.form
}

func (druid *Druid) InForm(form DruidForm) bool {
	return druid.form.Matches(form)
}

// TODO: don't hardcode numbers
func (druid *Druid) GetCatWeapon(level int32) core.Weapon {
	// Level 25 values
	claws := core.Weapon{
		BaseDamageMin:        0,
		BaseDamageMax:        0,
		SwingSpeed:           1.0,
		NormalizedSwingSpeed: 1.0,
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
	}

	switch level {
	case 60:
		// Avg: 54.8
		claws.BaseDamageMin = 43.84
		claws.BaseDamageMax = 65.76
	case 50:
		// TODO: Not entirely verified. Value from Balor (Feral mod)
		// Avg: 46.6
		claws.BaseDamageMin = 37.28
		claws.BaseDamageMax = 55.92
	case 40:
		claws.BaseDamageMin = 27.80305996
		claws.BaseDamageMax = 41.70460054
	default: // 25
		claws.BaseDamageMin = 16.3866
		claws.BaseDamageMax = 24.5799
	}

	return claws
}

func (druid *Druid) GetBearWeapon() core.Weapon {
	return core.Weapon{
		BaseDamageMin:        109,
		BaseDamageMax:        165,
		SwingSpeed:           2.5,
		NormalizedSwingSpeed: 2.5,
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
	}
}

// TODO: Class bonus stats for both cat and bear.
func (druid *Druid) GetFormShiftStats() stats.Stats {
	s := stats.Stats{
		stats.AttackPower: float64(druid.Talents.PredatoryStrikes) * 0.5 * float64(druid.Level),
		stats.MeleeCrit:   float64(druid.Talents.SharpenedClaws) * 2 * core.CritRatingPerCritChance,
	}
	/*
		if weapon := druid.GetMHWeapon(); weapon != nil {
			dps := (weapon.WeaponDamageMax+weapon.WeaponDamageMin)/2.0/weapon.SwingSpeed + druid.PseudoStats.BonusMHDps
			weapAp := weapon.Stats[stats.AttackPower] + weapon.Enchant.Stats[stats.AttackPower]
			fap := math.Floor((dps - 54.8) * 14)

			s[stats.AttackPower] += fap
			s[stats.AttackPower] += (fap + weapAp) * ((0.2 / 3) * float64(druid.Talents.PredatoryStrikes))
		}
	*/

	return s
}

func (druid *Druid) GetDynamicPredStrikeStats() stats.Stats {
	// Accounts for ap bonus for 'dynamic' enchants
	// just scourgebane currently, this is a bit hacky but is needed as the bonus varies based on current target
	// so has to be 'cached' differently
	s := stats.Stats{}
	if weapon := druid.GetMHWeapon(); weapon != nil {
		bonusAp := 0.0
		if weapon.Enchant.EffectID == 3247 && druid.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			bonusAp += 140
		}
		s[stats.AttackPower] += bonusAp * ((0.2 / 3) * float64(druid.Talents.PredatoryStrikes))
	}
	return s
}

// TODO: Classic feral and bear
func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}

	srm := druid.getSavageRoarMultiplier()

	statBonus := druid.GetFormShiftStats().Add(stats.Stats{
		stats.AttackPower: float64(druid.Level) * 2,
	})

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 1)
	feralApDep := druid.NewDynamicStatDependency(stats.FeralAttackPower, stats.AttackPower, 1)

	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.Strength, 1.0+0.04*float64(druid.Talents.HeartOfTheWild))
	}

	clawWeapon := druid.GetCatWeapon(druid.Level)

	predBonus := stats.Stats{}

	druid.CatFormAura = druid.RegisterAura(core.Aura{
		Label:      "Cat Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Cat), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.CancelShapeshift(sim)
			}
			druid.form = Cat
			druid.SetCurrentPowerBar(core.EnergyBar)

			druid.AutoAttacks.SetMH(clawWeapon)

			druid.PseudoStats.ThreatMultiplier *= 0.71
			druid.AddStatDynamic(sim, stats.Dodge, 2*float64(druid.Talents.FelineSwiftness))
			druid.SetShapeshift(aura)

			predBonus = druid.GetDynamicPredStrikeStats()
			druid.AddStatsDynamic(sim, predBonus)
			druid.AddStatsDynamic(sim, statBonus)
			druid.EnableDynamicStatDep(sim, agiApDep)
			druid.EnableDynamicStatDep(sim, feralApDep)
			if hotwDep != nil {
				druid.EnableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= srm
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Activate(sim)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.SetCurrentPowerBar(core.ManaBar)

			druid.TigersFuryAura.Deactivate(sim)

			druid.AutoAttacks.SetMH(druid.WeaponFromMainHand())

			druid.PseudoStats.ThreatMultiplier /= 0.71
			druid.AddStatDynamic(sim, stats.Dodge, -2*float64(druid.Talents.FelineSwiftness))
			druid.SetShapeshift(nil)

			druid.AddStatsDynamic(sim, predBonus.Invert())
			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.DisableDynamicStatDep(sim, agiApDep)
			druid.DisableDynamicStatDep(sim, feralApDep)
			if hotwDep != nil {
				druid.DisableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetReplaceMHSwing(nil)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()

				//druid.TigersFuryAura.Deactivate(sim)

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= srm
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Deactivate(sim)
				}
			}
		},
	})

	energyMetrics := druid.NewEnergyMetrics(actionID)

	furorProcChance := 0.2 * float64(druid.Talents.Furor)

	hasWolfheadBonus := false
	if head := druid.Equipment.Head(); head != nil && (head.ID == WolfsheadHelm || head.Enchant.EffectID == sod.WolfsheadTrophy) {
		hasWolfheadBonus = true
	}

	druid.CatForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_DruidCatForm,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.55,
			Multiplier: 100 - 10*druid.Talents.NaturalShapeshifter,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Mimic a "/cast !Cat Form" macro for the purpose of powershifting.
			// To do actions outside of form, e.g. potions or sapper during shift, APL makers
			// should explicitly cancel the aura, do actions, then cast cat form.
			if druid.IsShapeshifted() {
				druid.CancelShapeshift(sim)
			}

			maxShiftEnergy := core.TernaryFloat64(sim.RandomFloat("Furor") < furorProcChance, 40, 0)
			maxShiftEnergy = core.TernaryFloat64(hasWolfheadBonus, maxShiftEnergy+20, maxShiftEnergy)
			energyDelta := maxShiftEnergy - druid.CurrentEnergy()

			if energyDelta > 0 {
				druid.AddEnergy(sim, energyDelta, energyMetrics)
			} else {
				druid.SpendEnergy(sim, -energyDelta, energyMetrics)
			}

			druid.CatFormAura.Activate(sim)
		},
	})
}

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: core.TernaryInt32(druid.Level < 40, 5487, 9634)}
	feralSkillCoefficient := float64((druid.PseudoStats.FeralCombatSkill + 300) / 5)

	statBonus := druid.GetFormShiftStats().Add(stats.Stats{
		stats.AttackPower: 3 * feralSkillCoefficient,
		stats.Health:      core.TernaryFloat64(druid.Level < 40, (18*feralSkillCoefficient)-160, (32*feralSkillCoefficient)-680),
	})

	hasWolfheadBonus := false
	if head := druid.Equipment.Head(); head != nil && (head.ID == WolfsheadHelm || head.Enchant.EffectID == sod.WolfsheadTrophy) {
		hasWolfheadBonus = true
	}

	feralApDep := druid.NewDynamicStatDependency(stats.FeralAttackPower, stats.AttackPower, 1)

	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.Stamina, 1.0+0.04*float64(druid.Talents.HeartOfTheWild))
	}

	sotfdtm := 1.0
	if druid.HasRune(proto.DruidRune_RuneChestSurvivalOfTheFittest) {
		sotfdtm = 0.81
	}

	clawWeapon := druid.GetBearWeapon()
	predBonus := stats.Stats{}

	druid.BearFormThreatMultiplier = 1.3 + 0.03*float64(druid.Talents.FeralInstinct)

	druid.BearFormAura = druid.RegisterAura(core.Aura{
		Label:      "Bear Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Bear), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.CancelShapeshift(sim)
			}
			druid.form = Bear
			druid.SetShapeshift(aura)
			druid.SetCurrentPowerBar(core.RageBar)

			druid.AutoAttacks.SetMH(clawWeapon)

			druid.PseudoStats.ThreatMultiplier *= druid.BearFormThreatMultiplier
			druid.PseudoStats.DamageTakenMultiplier *= sotfdtm

			predBonus = druid.GetDynamicPredStrikeStats()
			druid.AddStatsDynamic(sim, predBonus)
			druid.AddStatsDynamic(sim, statBonus)
			druid.ApplyDynamicEquipScaling(sim, stats.Armor, core.TernaryFloat64(druid.Level < 40, 1.8, 4.6))
			druid.ApplyDynamicEquipScaling(sim, stats.Armor, 1+.02*float64(druid.Talents.ThickHide))

			druid.EnableDynamicStatDep(sim, feralApDep)
			if hotwDep != nil {
				druid.EnableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetReplaceMHSwing(druid.ReplaceBearMHFunc)
				druid.AutoAttacks.EnableAutoSwing(sim)

				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.SetShapeshift(nil)
			druid.SetCurrentPowerBar(core.ManaBar)

			druid.AutoAttacks.SetMH(druid.WeaponFromMainHand())

			druid.PseudoStats.ThreatMultiplier /= druid.BearFormThreatMultiplier
			druid.PseudoStats.DamageTakenMultiplier /= sotfdtm

			druid.AddStatsDynamic(sim, predBonus.Invert())
			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.RemoveDynamicEquipScaling(sim, stats.Armor, core.TernaryFloat64(druid.Level < 40, 1.8, 4.6))
			druid.RemoveDynamicEquipScaling(sim, stats.Armor, 1+.02*float64(druid.Talents.ThickHide))
			druid.DisableDynamicStatDep(sim, feralApDep)

			if hotwDep != nil {
				druid.DisableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.EnableAutoSwing(sim)

				druid.manageCooldownsEnabled()
				druid.UpdateManaRegenRates()
				druid.EnrageAura.Deactivate(sim)
			}
		},
	})

	rageMetrics := druid.NewRageMetrics(actionID)

	furorProcChance := []float64{0, 0.2, 0.4, 0.6, 0.8, 1}[druid.Talents.Furor]

	druid.BearForm = druid.RegisterSpell(Any, core.SpellConfig{
		ClassSpellMask: ClassSpellMask_DruidBearForm,
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.55,
			Multiplier: 100 - 10*druid.Talents.NaturalShapeshifter,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rageDelta := core.TernaryFloat64(hasWolfheadBonus, 5, 0) - druid.CurrentRage()
			if sim.Proc(furorProcChance, "Furor") {
				rageDelta += 10
			}
			if rageDelta > 0 {
				druid.AddRage(sim, rageDelta, rageMetrics)
			} else if rageDelta < 0 {
				druid.SpendRage(sim, -rageDelta, rageMetrics)
			}
			druid.BearFormAura.Activate(sim)
		},
	})
}

func (druid *Druid) manageCooldownsEnabled() {
	// Disable cooldowns not usable in form and/or delay others
	if druid.StartingForm.Matches(Cat | Bear) {
		for _, mcd := range druid.disabledMCDs {
			mcd.Enable()
		}
		druid.disabledMCDs = nil

		if druid.InForm(Humanoid) {
			// Disable cooldown that incurs a gcd, so we dont get stuck out of form when we dont need to (Greater Drums)
			for _, mcd := range druid.GetMajorCooldowns() {
				if mcd.Spell.DefaultCast.GCD > 0 {
					mcd.Disable()
					druid.disabledMCDs = append(druid.disabledMCDs, mcd)
				}
			}
		}
	}
}

// https://www.wowhead.com/classic/spell=24858/moonkin-form
// - Moonfire costs 50% less mana and deals 50% more damage over time
// - Sunfire costs 50% less mana and deals 50% more damage over time
// - Your periodic damage spells can deal critical periodic damage (handled in individual dot snapshots)
// - You gain (2 * Level) spell damage
func (druid *Druid) registerMoonkinFormSpell() {
	if !druid.Talents.MoonkinForm {
		return
	}

	actionID := core.ActionID{SpellID: 24858}

	druid.MoonfireDotMultiplier = 1.0
	druid.SunfireDotMultiplier = 1.0

	druid.MoonkinFormAura = druid.RegisterAura(core.Aura{
		Label:    "Moonkin Form",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.CancelShapeshift(sim)
			}
			druid.form = Moonkin

			druid.AddStatDynamic(sim, stats.SpellDamage, float64(3*druid.Level))

			druid.MoonfireDotMultiplier *= 2.0
			core.Each(druid.Moonfire, func(spell *DruidSpell) {
				if spell != nil {
					spell.Spell.Cost.Multiplier -= 50
				}
			})

			if druid.HasRune(proto.DruidRune_RuneHandsSunfire) {
				druid.Sunfire.Cost.Multiplier -= 50
				druid.SunfireDotMultiplier *= 2.0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.AddStatDynamic(sim, stats.SpellDamage, float64(-3*druid.Level))

			core.Each(druid.Moonfire, func(spell *DruidSpell) {
				if spell != nil {
					spell.Spell.Cost.Multiplier += 50
				}
			})
			druid.MoonfireDotMultiplier /= 2.0

			if druid.HasRune(proto.DruidRune_RuneHandsSunfire) {
				druid.Sunfire.Cost.Multiplier += 50
				druid.SunfireDotMultiplier /= 2.0
			}
		},
	})

	druid.MoonkinForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_DruidMoonkinForm,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.35,
			Multiplier: 100 - 10*druid.Talents.NaturalShapeshifter,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.MoonkinFormAura.Activate(sim)
		},
	})
}
