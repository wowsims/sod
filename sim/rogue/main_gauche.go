package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (rogue *Rogue) registerMainGaucheSpell() {
	if !rogue.HasRune(proto.RogueRune_RuneMainGauche) {
		return
	}
	hasPKRune := rogue.HasRune(proto.RogueRune_RunePoisonedKnife)
	hasQDRune := rogue.HasRune(proto.RogueRune_RuneQuickDraw)

	mainGaucheAura := rogue.RegisterAura(core.Aura{
		Label:    "Main Gauche Buff",
		ActionID: core.ActionID{SpellID: int32(proto.RogueRune_RuneMainGauche)},
		Duration: time.Second * 5,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee|core.ProcMaskRanged) && result.DidParry() {
				aura.Deactivate(sim)
			}
		},
	}).AttachStatBuff(stats.Parry, 100*core.ParryRatingPerParryChance)

	mainGaucheSSAura := rogue.RegisterAura(core.Aura{
		Label:    "Main Gauche Sinister Strike Discount",
		ActionID: core.ActionID{SpellID: 462752},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.SinisterStrike.Cost.FlatModifier -= 20
			rogue.SinisterStrike.ThreatMultiplier *= 1.5
			rogue.SinisterStrike.ApplyMultiplicativeDamageBonus(1.5)
			rogue.Eviscerate.ApplyMultiplicativeDamageBonus(1.5)

			if hasPKRune {
				rogue.PoisonedKnife.Cost.FlatModifier -= 20
				rogue.PoisonedKnife.ThreatMultiplier *= 1.5
				rogue.PoisonedKnife.ApplyMultiplicativeDamageBonus(1.5)
			}

			if hasQDRune {
				rogue.QuickDraw.Cost.FlatModifier -= 20
				rogue.QuickDraw.ThreatMultiplier *= 1.5
				rogue.QuickDraw.ApplyMultiplicativeDamageBonus(1.5)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.SinisterStrike.Cost.FlatModifier += 20
			rogue.SinisterStrike.ThreatMultiplier /= 1.5
			rogue.SinisterStrike.ApplyMultiplicativeDamageBonus(1 / 1.5)
			rogue.Eviscerate.ApplyMultiplicativeDamageBonus(1 / 1.5)

			if hasPKRune {
				rogue.PoisonedKnife.Cost.FlatModifier += 20
				rogue.PoisonedKnife.ThreatMultiplier /= 1.5
				rogue.PoisonedKnife.ApplyMultiplicativeDamageBonus(1 / 1.5)
			}

			if hasQDRune {
				rogue.QuickDraw.Cost.FlatModifier += 20
				rogue.QuickDraw.ThreatMultiplier /= 1.5
				rogue.QuickDraw.ApplyMultiplicativeDamageBonus(1 / 1.5)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellResult *core.SpellResult) {
			if spell.ProcMask.Matches(rogue.SinisterStrike.ProcMask) && spellResult.Landed() {
				rogue.RollingWithThePunchesProcAura.AddStack(sim)
			}
		},
	})

	rogue.MainGauche = rogue.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_RogueMainGauche,
		ActionID:       mainGaucheAura.ActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          rogue.builderFlags(),
		MaxRange:       5,

		EnergyCost: core.EnergyCostOptions{
			Cost:   15,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 15,
			},
			IgnoreHaste: true,
		},

		CritDamageBonus: rogue.lethality(),

		DamageMultiplier: 1,
		ThreatMultiplier: 3,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			// Auras gained regardless of landed hit.
			mainGaucheAura.Activate(sim)
			mainGaucheSSAura.Activate(sim)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, target, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
