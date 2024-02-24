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

	cooldown := time.Second * 6
	manaCost := .01

	imbueMultiplier := core.TernaryFloat64(shaman.GetCharacter().Consumes.OffHandImbue == proto.WeaponImbue_FlametongueWeapon, 1.5, 1)

	shaman.LavaLash = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsLavaLash)},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIgnoreResists,

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

		DamageMultiplier: 1.5 * imbueMultiplier,
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: shaman.ShamanThreatMultiplier(1),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()) *
				shaman.AutoAttacks.OHConfig().DamageMultiplier
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
