package hunter

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	// BM Talents
	hunter.applyEnduranceTraining()
	hunter.applyFerocity()
	hunter.applyUnleashedFury()
	hunter.applyBestialDiscipline()
	hunter.applyFrenzy()
	hunter.registerBestialWrathCD()

	// MM Talents
	hunter.applyEfficiency()
	hunter.applyLethalShots()
	hunter.applyMortalShots()
	hunter.applyImprovedArcaneShot()
	hunter.applyImprovedSerpentSting()
	hunter.applyBarrage()
	hunter.applyRangedWeaponSpecialization()

	// Survival Talents
	hunter.applyMonsterSlaying()
	hunter.applyHumanoidSlaying()
	hunter.applySavageStrikes()
	hunter.applyCleverTraps()
	hunter.applySurvivalist()
	hunter.applyTrapMastery()
	hunter.applySurefooted()
	hunter.applyKillerInstinct()
	hunter.applyLightningReflexes()

	// Draught was confirmed to Stack with Monster and Humanoid Slaying talents
	if hunter.Consumes.MiscConsumes != nil && hunter.Consumes.MiscConsumes.DraughtOfTheSands {
		hunter.Env.RegisterPostFinalizeEffect(func() {
			multiplier := 1.03
			for _, t := range hunter.Env.Encounter.Targets {
				for _, at := range hunter.AttackTables[t.UnitIndex] {
					at.DamageDealtMultiplier *= multiplier
					at.CritMultiplier *= multiplier
				}
			}
		})
	}
}

///////////////////////////////////////////////////////////////////////////
//                            Beast Mastery Talents
///////////////////////////////////////////////////////////////////////////

func (hunter *Hunter) applyEnduranceTraining() {
	if hunter.Talents.EnduranceTraining == 0 || hunter.pet == nil {
		return
	}

	hunter.pet.MultiplyStat(stats.Health, 1+(0.03*float64(hunter.Talents.EnduranceTraining)))
}

func (hunter *Hunter) applyFerocity() {
	if hunter.Talents.Ferocity == 0 || hunter.pet == nil {
		return
	}

	hunter.pet.AddStat(stats.MeleeCrit, 3*float64(hunter.Talents.Ferocity)*core.CritRatingPerCritChance)
	hunter.pet.AddStat(stats.SpellCrit, 3*float64(hunter.Talents.Ferocity)*core.SpellCritRatingPerCritChance)
}

func (hunter *Hunter) applyUnleashedFury() {
	if hunter.Talents.UnleashedFury == 0 || hunter.pet == nil {
		return
	}

	hunter.pet.PseudoStats.DamageDealtMultiplierAdditive += 0.04 * float64(hunter.Talents.UnleashedFury)
}

func (hunter *Hunter) applyBestialDiscipline() {
	if hunter.Talents.BestialDiscipline == 0 || hunter.pet == nil {
		return
	}

	hunter.pet.AddFocusRegenMultiplier(0.1 * float64(hunter.Talents.BestialDiscipline))
}

func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 || hunter.pet == nil {
		return
	}

	procChance := 0.2 * float64(hunter.Talents.Frenzy)

	procAura := hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy Proc",
		ActionID: core.ActionID{SpellID: 19625},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.3)
		},
	})

	hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellResult *core.SpellResult) {
			if !spellResult.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if spell.ClassSpellMask == ClassSpellMask_HunterPetFlankingStrike {
				return
			}

			if procChance == 1 || sim.RandomFloat("Frenzy") < procChance {
				procAura.Activate(sim)
			}
		},
	})
}

func (hunter *Hunter) registerBestialWrathCD() {
	if !hunter.Talents.BestialWrath || hunter.pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 19574}

	hunter.BestialWrathPetAura = hunter.pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 18,
	}).AttachMultiplicativePseudoStatBuff(&hunter.pet.PseudoStats.DamageDealtMultiplier, 1.5)

	bwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.BestialWrathPetAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: bwSpell,
		Type:  core.CooldownTypeDPS,
	})
}

///////////////////////////////////////////////////////////////////////////
//                            Marksmanship Talents
///////////////////////////////////////////////////////////////////////////

func (hunter *Hunter) applyEfficiency() {
	if hunter.Talents.Efficiency == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind: core.SpellMod_PowerCost_Pct,
		// Applies to all shots, stings, strikes, and volley
		ClassMask: ClassSpellMask_HunterShots | ClassSpellMask_HunterStrikes | ClassSpellMask_HunterStings | ClassSpellMask_HunterVolley,
		IntValue:  -2 * int64(hunter.Talents.Efficiency),
	})
}

func (hunter *Hunter) applyLethalShots() {
	if hunter.Talents.LethalShots == 0 {
		return
	}

	hunter.AddStat(stats.MeleeCrit, float64(hunter.Talents.LethalShots)*core.CritRatingPerCritChance)
}

func (hunter *Hunter) applyImprovedArcaneShot() {
	if hunter.Talents.ImprovedArcaneShot == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: ClassSpellMask_HunterArcaneShot,
		TimeValue: -time.Millisecond * 200 * time.Duration(hunter.Talents.ImprovedArcaneShot),
	})
}

func (hunter *Hunter) applyImprovedSerpentSting() {
	if hunter.Talents.ImprovedSerpentSting == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterSerpentSting | ClassSpellMask_HunterSoFSerpentSting | ClassSpellMask_HunterChimeraSerpent,
		IntValue:  int64(2 * hunter.Talents.ImprovedSerpentSting),
	})
}

func (hunter *Hunter) applyMortalShots() {
	if hunter.Talents.MortalShots == 0 {
		return
	}

	hunter.AutoAttacks.RangedConfig().CritDamageBonus += 0.06 * float64(hunter.Talents.MortalShots)

	hunter.AddStaticMod(core.SpellModConfig{
		Kind: core.SpellMod_CritDamageBonus_Flat,
		// Applies to all shots, strikes, and volley
		ClassMask:  ClassSpellMask_HunterShots | ClassSpellMask_HunterStrikes | ClassSpellMask_HunterChimeraSerpent | ClassSpellMask_HunterMongooseBite | ClassSpellMask_HunterWingClip | ClassSpellMask_HunterVolley,
		FloatValue: 0.06 * float64(hunter.Talents.MortalShots),
	})
}

func (hunter *Hunter) applyBarrage() {
	if hunter.Talents.Barrage == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterMultiShot | ClassSpellMask_HunterVolley,
		IntValue:  int64(5 * hunter.Talents.Barrage),
	})
}

func (hunter *Hunter) applyRangedWeaponSpecialization() {
	if hunter.Talents.RangedWeaponSpecialization == 0 {
		return
	}

	mult := 1 + 0.01*float64(hunter.Talents.RangedWeaponSpecialization)
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_HunterShots) || spell.ProcMask.Matches(core.ProcMaskRangedAuto) {
			spell.ApplyMultiplicativeDamageBonus(mult)
		}
	})
}

///////////////////////////////////////////////////////////////////////////
//                            Survival Talents
///////////////////////////////////////////////////////////////////////////

func (hunter *Hunter) applyMonsterSlaying() {
	if hunter.Talents.MonsterSlaying == 0 {
		return
	}

	monsterSlayingMobTypes := []proto.MobType{proto.MobType_MobTypeBeast, proto.MobType_MobTypeGiant, proto.MobType_MobTypeDragonkin}
	monsterSlayingTargets := core.FilterSlice(hunter.Env.Encounter.Targets, func(t *core.Target) bool { return slices.Contains(monsterSlayingMobTypes, t.MobType) })
	if len(monsterSlayingTargets) == 0 {
		return
	}

	monsterMultiplier := []float64{1, 1.01, 1.02, 1.03}[hunter.Talents.MonsterSlaying]
	hunter.Env.RegisterPostFinalizeEffect(func() {
		for _, t := range monsterSlayingTargets {
			for _, at := range hunter.AttackTables[t.UnitIndex] {
				at.DamageDealtMultiplier *= monsterMultiplier
				at.CritMultiplier *= monsterMultiplier
			}
		}
	})
}

func (hunter *Hunter) applyHumanoidSlaying() {
	if hunter.Talents.HumanoidSlaying == 0 {
		return
	}

	humanoidSlayingTargets := core.FilterSlice(hunter.Env.Encounter.Targets, func(t *core.Target) bool { return t.MobType == proto.MobType_MobTypeHumanoid })
	if len(humanoidSlayingTargets) == 0 {
		return
	}

	humanoidMultiplier := []float64{1, 1.01, 1.02, 1.03}[hunter.Talents.HumanoidSlaying]
	hunter.Env.RegisterPostFinalizeEffect(func() {
		for _, t := range humanoidSlayingTargets {
			for _, at := range hunter.AttackTables[t.UnitIndex] {
				at.DamageDealtMultiplier *= humanoidMultiplier
				at.CritMultiplier *= humanoidMultiplier
			}
		}
	})
}

func (hunter *Hunter) applySavageStrikes() {
	if hunter.Talents.SavageStrikes == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Flat,
		ClassMask:  ClassSpellMask_HunterRaptorStrikeHit | ClassSpellMask_HunterMongooseBite,
		FloatValue: 10 * float64(hunter.Talents.SavageStrikes) * core.CritRatingPerCritChance,
	})
}

func (hunter *Hunter) applyCleverTraps() {
	if hunter.Talents.CleverTraps == 0 {
		return
	}

	multiplier := 1 + 0.15*float64(hunter.Talents.CleverTraps)

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_HunterTraps) {
			spell.ApplyMultiplicativeDamageBonus(multiplier)
		}
	})
}

func (hunter *Hunter) applySurvivalist() {
	if hunter.Talents.Survivalist == 0 {
		return
	}

	hunter.MultiplyStat(stats.Health, 1.0+0.02*float64(hunter.Talents.Survivalist))
}

func (hunter *Hunter) applyTrapMastery() {
	if hunter.Talents.TrapMastery == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusHit_Flat,
		ClassMask:  ClassSpellMask_HunterTraps,
		FloatValue: float64(5 * hunter.Talents.TrapMastery),
	})
}

func (hunter *Hunter) applySurefooted() {
	if hunter.Talents.Surefooted == 0 {
		return
	}

	hunter.AddStat(stats.MeleeHit, float64(hunter.Talents.Surefooted)*core.MeleeHitRatingPerHitChance)
}

func (hunter *Hunter) applyKillerInstinct() {
	if hunter.Talents.KillerInstinct == 0 {
		return
	}

	hunter.AddStat(stats.MeleeCrit, float64(hunter.Talents.KillerInstinct)*core.CritRatingPerCritChance)
}

func (hunter *Hunter) applyLightningReflexes() {
	if hunter.Talents.LightningReflexes == 0 {
		return
	}

	agiBonus := 0.03 * float64(hunter.Talents.LightningReflexes)
	hunter.MultiplyStat(stats.Agility, 1.0+agiBonus)
}
