package shaman

import "github.com/wowsims/sod/sim/core/proto"

func (shaman *Shaman) registerEarthShieldSpell() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsEarthShield) {
		return
	}

	shaman.PseudoStats.SpellPushbackMultiplier *= 0.70

	// actionID := core.ActionID{SpellID: 49284}
	// spCoeff := 0.286

	// icd := core.Cooldown{
	// 	Timer:    shaman.NewTimer(),
	// 	Duration: time.Millisecond * 3500,
	// }

	// shaman.EarthShield = shaman.RegisterSpell(core.SpellConfig{
	// 	ActionID:    actionID,
	// 	SpellSchool: core.SpellSchoolNature,
	// 	DefenseType: core.DefenseTypeMagic,
	// 	ProcMask:    core.ProcMaskEmpty,
	// 	Flags:       core.SpellFlagHelpful | core.SpellFlagAPL | SpellFlagShaman,

	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: core.GCDDefault,
	// 		},
	// 	},

	// 	Hot: core.DotConfig{
	// 		Aura: core.Aura{
	// 			Label:    "Earth Shield",
	// 			ActionID: core.ActionID{SpellID: 379},
	// 			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 				if !result.Landed() {
	// 					return
	// 				}
	// 				if !icd.IsReady(sim) {
	// 					return
	// 				}
	// 				icd.Use(sim)
	// 				shaman.EarthShield.Hot(result.Target).ManualTick(sim)
	// 			},
	// 			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
	// 			},
	// 		},
	// 		NumberOfTicks: 6 + shaman.Talents.ImprovedEarthShield,
	// 		TickLength:    time.Minute*10 + 1, // tick length longer than expire time.
	// 		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
	// 			dot.SnapshotBaseDamage = 377 + dot.Spell.HealingPower(target)*spCoeff
	// 			dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
	// 		},
	// 	},

	// 	DamageMultiplier: 1
	// 	ThreatMultiplier: 1,

	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.Hot(target).Apply(sim)
	// 	},
	// })
}
