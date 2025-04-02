package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// https://www.wowhead.com/classic/spell=408507/lava-lash
func (shaman *Shaman) applyLavaLash() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaLash) || !shaman.AutoAttacks.IsDualWielding {
		return
	}

	cooldown := time.Second * 6
	manaCost := .01

	damageMultiplier := 1.0
	// When off-hand is imbued with flametongue weapon increases damage by 125%
	if shaman.GetCharacter().Consumes.OffHandImbue == proto.WeaponImbue_FlametongueWeapon {
		damageMultiplier += 1.50
	}

	shaman.LavaLash = shaman.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanLavaLash,
		ActionID:       core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsLavaLash)},
		SpellSchool:    core.SpellSchoolFire,
		CastType:       proto.CastType_CastTypeOffHand,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
			// Refund: 0.8, -- Not implemented for ManaCostOption
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

		DamageMultiplier: damageMultiplier * (1 + (.02 * float64(shaman.Talents.WeaponMastery))),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())) *
				shaman.AutoAttacks.OHConfig().DamageMultiplier
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
