package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) applyLavaLash() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaLash) || !shaman.AutoAttacks.IsDualWielding {
		return
	}

	ohWeaponDamageMultiplier := 1.50
	cooldown := time.Second * 6
	manaCost := .01

	damageMultiplier := core.TernaryFloat64(shaman.GetCharacter().Consumes.OffHandImbue == proto.WeaponImbue_FlametongueWeapon, 1.2, 1)

	shaman.LavaLash = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsLavaLash)},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: shaman.ShamanThreatMultiplier(1),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())*ohWeaponDamageMultiplier + spell.BonusWeaponDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
