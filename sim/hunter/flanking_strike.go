package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) registerFlankingStrikeSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneLegsFlankingStrike) {
		return
	}

	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:     core.SpellMod_DamageDone_Flat,
		ProcMask: core.ProcMaskMelee,
	})

	buffAura := hunter.GetOrRegisterAura(core.Aura{
		Label:     "Flanking Strike Buff",
		ActionID:  core.ActionID{SpellID: 415320},
		MaxStacks: 3,
		Duration:  time.Second * 10,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damageMod.UpdateIntValue(int64(8 * newStacks))
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
	})

	if hunter.pet != nil {
		hunter.pet.flankingStrike = hunter.pet.GetOrRegisterSpell(core.SpellConfig{
			ClassSpellMask: ClassSpellMask_HunterPetFlankingStrike,
			ActionID:       core.ActionID{SpellID: 415320},
			SpellSchool:    core.SpellSchoolPhysical,
			DefenseType:    core.DefenseTypeMelee,
			ProcMask:       core.ProcMaskMeleeMHSpecial,
			Flags:          core.SpellFlagMeleeMetrics,

			CritDamageBonus:  hunter.mortalShots(),
			DamageMultiplier: 1,
			BonusCoefficient: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			},
		})

		core.MakeProcTriggerAura(&hunter.pet.Unit, core.ProcTrigger{
			Name:           "Flanking Strike Refresh",
			ClassSpellMask: ClassSpellMask_HunterPetSpecials,
			ProcChance:     0.50,
			Callback:       core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				hunter.FlankingStrike.CD.Set(sim.CurrentTime)
			},
		})
	}

	hunter.FlankingStrike = hunter.GetOrRegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterFlankingStrike,
		ActionID:       core.ActionID{SpellID: 415320},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagStrike,

		MaxRange: core.MaxMeleeAttackRange,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.015,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		CritDamageBonus:  hunter.mortalShots(),
		DamageMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hunter.pet != nil {
				hunter.pet.flankingStrike.Cast(sim, hunter.pet.CurrentTarget)
			}

			buffAura.Activate(sim)
			buffAura.AddStack(sim)
		},

		RelatedSelfBuff: buffAura,
	})
}
