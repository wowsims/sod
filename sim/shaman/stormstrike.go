package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) registerStormstrikeSpell() {
	if !shaman.Talents.Stormstrike {
		return
	}

	hasDualWieldSpecRune := shaman.HasRune(proto.ShamanRune_RuneChestDualWieldSpec)

	shaman.StormstrikeMH = shaman.newStormstrikeHitSpell(true)
	if hasDualWieldSpecRune {
		shaman.StormstrikeOH = shaman.newStormstrikeHitSpell(false)
	}

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanStormstrike,
		ActionID:       core.ActionID{SpellID: 17364},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: .063,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// offhand always swings first
			if shaman.AutoAttacks.IsDualWielding && shaman.StormstrikeOH != nil {
				shaman.StormstrikeOH.Cast(sim, target)
			}
			shaman.StormstrikeMH.Cast(sim, target)
		},
	})
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) *core.Spell {
	procMask := core.ProcMaskMeleeMHSpecial
	flags := core.SpellFlagMeleeMetrics
	damageMultiplier := 1.0
	damageFunc := shaman.MHWeaponDamage
	if !isMH {
		// Only the main-hand hit triggers procs / the debuff
		procMask = core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc | core.ProcMaskMeleeDamageProc
		flags |= core.SpellFlagNoOnCastComplete
		damageMultiplier = shaman.AutoAttacks.OHConfig().DamageMultiplier
		damageFunc = shaman.OHWeaponDamage
	}

	stormStrikeAuras := shaman.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
		return core.StormstrikeAura(target)
	})

	return shaman.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanStormstrikeHit,
		ActionID:       core.ActionID{SpellID: 17364}.WithTag(int32(core.Ternary(isMH, 1, 2))),
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       procMask,
		Flags:          flags,

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damageFunc(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if isMH && result.Landed() {
				stormStrikeAuras.Get(target).Activate(sim)
			}
		},
	})
}
