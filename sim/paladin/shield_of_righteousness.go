package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerShieldOfRighteousness() {
	if !paladin.hasRune(proto.PaladinRune_RuneCloakShieldOfRighteousness) {
		return
	}

	// Base damage formula from wowhead tooltip:
	// https://www.wowhead.com/classic/spell=440658/shield-of-righteousness
	// Testing shows there is an additional 20 base damage included.
	damage := (179.0 * paladin.baseRuneAbilityDamage() / 100.0) + 20.0

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.PaladinRune_RuneCloakShieldOfRighteousness)},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		SpellCode:   SpellCode_PaladinShieldOfRighteousness,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Shield of Righteousness has a hidden scaling coefficient of 2.2x SBV (derived from testing)
			baseDamage := damage + paladin.BlockValue()*2.2
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
