package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) ApplyRunes() {
	if hunter.HasRune(proto.HunterRune_RuneChestHeartOfTheLion) {
		statMultiply := 1.1
		hunter.MultiplyStat(stats.Strength, statMultiply)
		hunter.MultiplyStat(stats.Stamina, statMultiply)
		hunter.MultiplyStat(stats.Agility, statMultiply)
		hunter.MultiplyStat(stats.Intellect, statMultiply)
		hunter.MultiplyStat(stats.Spirit, statMultiply)
	}

	if hunter.HasRune(proto.HunterRune_RuneChestMasterMarksman) {
		hunter.AddStat(stats.MeleeCrit, 5*core.CritRatingPerCritChance)
		hunter.AddStat(stats.SpellCrit, 5*core.SpellCritRatingPerCritChance)
	}

	if hunter.HasRune(proto.HunterRune_RuneChestLoneWolf) && hunter.pet == nil {
		hunter.PseudoStats.DamageDealtMultiplier *= 1.3
	}

	if hunter.HasRune(proto.HunterRune_RuneHandsBeastmastery) && hunter.pet != nil {
		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1.2
	}

	if hunter.HasRune(proto.HunterRune_RuneBootsDualWieldSpecialization) {
		hunter.AutoAttacks.OHConfig().DamageMultiplier *= 1.5
	}

	hunter.applySniperTraining()
	hunter.applyCobraStrikes()
	hunter.applyExposeWeakness()
	hunter.applyInvigoration()
}

func (hunter *Hunter) applyInvigoration() {
	if !hunter.HasRune(proto.HunterRune_RuneBootsInvigoration) || hunter.pet == nil {
		return
	}

	procSpellId := core.ActionID{SpellID: 437999}
	metrics := hunter.NewManaMetrics(procSpellId)
	procSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    procSpellId,
		SpellSchool: core.SpellSchoolNature,
		ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
			hunter.AddMana(sim, hunter.MaxMana()*0.05, metrics)
		},
	})

	core.MakePermanent(hunter.pet.GetOrRegisterAura(core.Aura{
		Label: "Invigoration",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
				return
			}

			if !result.DidCrit() {
				return
			}

			procSpell.Cast(sim, result.Target)
		},
	}))
}

func (hunter *Hunter) applyExposeWeakness() {
	if !hunter.HasRune(proto.HunterRune_RuneBeltExposeWeakness) {
		return
	}

	apBonus := hunter.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 0.4)
	apRangedBonus := hunter.NewDynamicStatDependency(stats.Agility, stats.RangedAttackPower, 0.4)

	procAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Expose Weakness Proc",
		ActionID: core.ActionID{SpellID: 409507},
		Duration: time.Second * 7,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.EnableDynamicStatDep(sim, apBonus)
			hunter.EnableDynamicStatDep(sim, apRangedBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.DisableDynamicStatDep(sim, apBonus)
			hunter.DisableDynamicStatDep(sim, apRangedBonus)
		},
	})

	core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
		Label: "Expose Weakness",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}

			if !result.DidCrit() {
				return
			}

			procAura.Activate(sim)
		},
	}))
}

func (hunter *Hunter) applySniperTraining() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsSniperTraining) {
		return
	}

	hunter.SniperTrainingAura = hunter.GetOrRegisterAura(core.Aura{
		Label:    "Sniper Training",
		ActionID: core.ActionID{SpellID: 415399},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range aura.Unit.Spellbook {
				if spell.ProcMask.Matches(core.ProcMaskRangedSpecial) {
					spell.BonusCritRating += 10 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range aura.Unit.Spellbook {
				if spell.ProcMask.Matches(core.ProcMaskRangedSpecial) {
					spell.BonusCritRating -= 10 * core.CritRatingPerCritChance
				}
			}
		},
	})

	core.ApplyFixedUptimeAura(hunter.SniperTrainingAura, hunter.Options.SniperTrainingUptime, time.Second*6, 0)
}

func (hunter *Hunter) applyCobraStrikes() {
	if !hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes) || hunter.pet == nil {
		return
	}

	hunter.CobraStrikesAura = hunter.pet.GetOrRegisterAura(core.Aura{
		Label:     "Cobra Strikes",
		ActionID:  core.ActionID{SpellID: 425714},
		Duration:  time.Second * 30,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range aura.Unit.Spellbook {
				if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
					spell.BonusCritRating += 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range aura.Unit.Spellbook {
				if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
					spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				aura.RemoveStack(sim)
			}
		},
	})
}
