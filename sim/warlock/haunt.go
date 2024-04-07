package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func hauntMultiplier(spell *core.Spell, _ *core.AttackTable) float64 {
	return core.TernaryFloat64(spell.Flags.Matches(SpellFlagHaunt), 1.2, 1)
}

func (warlock *Warlock) registerHauntSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsHaunt) {
		return
	}

	actionID := core.ActionID{SpellID: 403501}

	spellCoeff := 0.714
	baseLowDamage := warlock.baseRuneAbilityDamage() * 2.51
	baseHighDamage := warlock.baseRuneAbilityDamage() * 2.95

	warlock.HauntDebuffAuras = warlock.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Haunt-" + warlock.Label,
			ActionID: actionID,
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.AttackTables[aura.Unit.UnitIndex][proto.CastType_CastTypeMainHand].DamageDoneByCasterMultiplier = hauntMultiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.AttackTables[aura.Unit.UnitIndex][proto.CastType_CastTypeMainHand].DamageDoneByCasterMultiplier = nil
			},
		})
	})

	warlock.Haunt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.ShadowMastery),
		ThreatMultiplier:         1,
		BonusCoefficient:         spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					warlock.HauntDebuffAuras.Get(result.Target).Activate(sim)
					warlock.EverlastingAfflictionRefresh(sim, target)
				}
			})
		},
	})
}
