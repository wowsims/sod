package hunter

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	if hunter.pet != nil {
		hunter.applyFrenzy()
		hunter.registerBestialWrathCD()

		hunter.pet.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*3*float64(hunter.Talents.Ferocity))
		hunter.pet.AddStat(stats.SpellCrit, core.SpellCritRatingPerCritChance*3*float64(hunter.Talents.Ferocity))

		hunter.pet.PseudoStats.DamageDealtMultiplierAdditive += 0.04 * float64(hunter.Talents.UnleashedFury)

		if hunter.Talents.EnduranceTraining > 0 {
			hunter.pet.MultiplyStat(stats.Health, 1+(0.03*float64(hunter.Talents.EnduranceTraining)))
		}
	}

	monsterSlayingMobTypes := []proto.MobType{proto.MobType_MobTypeBeast, proto.MobType_MobTypeGiant, proto.MobType_MobTypeDragonkin}
	monsterSlayingTargets := core.FilterSlice(hunter.Env.Encounter.Targets, func(t *core.Target) bool { return slices.Contains(monsterSlayingMobTypes, t.MobType) })
	humanoidSlayingTargets := core.FilterSlice(hunter.Env.Encounter.Targets, func(t *core.Target) bool { return t.MobType == proto.MobType_MobTypeHumanoid })

	if hunter.Talents.MonsterSlaying > 0 && len(monsterSlayingTargets) > 0 {
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

	if hunter.Talents.HumanoidSlaying > 0 && len(humanoidSlayingTargets) > 0 {
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

	// Draught was confirmed to Stack with talents
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

	if hunter.Talents.BestialDiscipline > 0 {
		core.MakePermanent(hunter.RegisterAura(core.Aura{
			Label: "Bestial Discipline",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				if hunter.pet != nil {
					hunter.pet.AddFocusRegenMultiplier(0.1 * float64(hunter.Talents.BestialDiscipline))
				}
			},
		}))
	}

	hunter.AddStat(stats.MeleeHit, float64(hunter.Talents.Surefooted)*1*core.MeleeHitRatingPerHitChance)
	hunter.AddStat(stats.SpellHit, float64(hunter.Talents.Surefooted)*1*core.SpellHitRatingPerHitChance)

	hunter.AddStat(stats.MeleeCrit, float64(hunter.Talents.KillerInstinct)*1*core.CritRatingPerCritChance)

	if hunter.Talents.LethalShots > 0 {
		hunter.AddStat(stats.MeleeCrit, 1*float64(hunter.Talents.LethalShots)*core.CritRatingPerCritChance)
	}

	if hunter.Talents.RangedWeaponSpecialization > 0 {
		mult := 1 + 0.01*float64(hunter.Talents.RangedWeaponSpecialization)
		hunter.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(ClassSpellMask_HunterShots) || spell.ProcMask.Matches(core.ProcMaskRangedAuto) {
				spell.ApplyMultiplicativeDamageBonus(mult)
			}
		})
	}

	if hunter.Talents.Survivalist > 0 {
		hunter.MultiplyStat(stats.Health, 1.0+0.02*float64(hunter.Talents.Survivalist))
	}

	if hunter.Talents.LightningReflexes > 0 {
		agiBonus := 0.03 * float64(hunter.Talents.LightningReflexes)
		hunter.MultiplyStat(stats.Agility, 1.0+agiBonus)
	}

	hunter.applyEfficiency()
	hunter.applyTrapMastery()
	hunter.applyCleverTraps()
	hunter.applyImprovedSerpentSting()
}

func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 {
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
	if !hunter.Talents.BestialWrath {
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

func (hunter *Hunter) mortalShots() float64 {
	return 0.06 * float64(hunter.Talents.MortalShots)
}

func (hunter *Hunter) applyTrapMastery() {
	if hunter.Talents.TrapMastery == 0 {
		return
	}

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_HunterTraps) {
			spell.BonusHitRating += 5 * float64(hunter.Talents.TrapMastery)
		}
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

func (hunter *Hunter) applyEfficiency() {
	if hunter.Talents.Efficiency == 0 {
		return
	}

	costModifier := 2 * hunter.Talents.Efficiency

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		// applies to Stings, Shots, Strikes and Volley
		if spell.Cost != nil && (spell.Flags.Matches(SpellFlagSting|SpellFlagStrike) || spell.Matches(ClassSpellMask_HunterShots|ClassSpellMask_HunterVolley)) {
			spell.Cost.FlatModifier -= costModifier
		}
	})
}

func (hunter *Hunter) applyImprovedSerpentSting() {
	if hunter.Talents.ImprovedSerpentSting == 0 {
		return
	}

	damageBonus := int64(2 * hunter.Talents.ImprovedSerpentSting)

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Improved Serpent Sting",
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_HunterSerpentSting | ClassSpellMask_HunterSoFSerpentSting | ClassSpellMask_HunterChimeraSerpent,
		IntValue:  damageBonus,
	}))
}
