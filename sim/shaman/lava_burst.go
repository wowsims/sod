package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) registerLavaBurstSpell() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
		return
	}

	shaman.LavaBurst = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(false))

	if shaman.HasRune(proto.ShamanRune_RuneChestOverload) {
		shaman.LavaBurstOverload = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(true))
	}
}

func (shaman *Shaman) newLavaBurstSpellConfig(isOverload bool) core.SpellConfig {
	hasMaelstromWeaponRune := shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon)

	baseDamageLow := shaman.baseRuneAbilityDamage() * 4.69
	baseDamageHigh := shaman.baseRuneAbilityDamage() * 6.05
	spellCoeff := .571
	castTime := time.Second * 2
	cooldown := time.Second * 8
	manaCost := .10

	var flags core.SpellFlag
	if !isOverload {
		flags |= core.SpellFlagAPL
	}

	spell := core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanLavaBurst,
		ActionID:       core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsLavaBurst)},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          flags,
		MissileSpeed:   20,
		MetricSplits:   MaelstromWeaponSplits,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
			// Convection does not currently apply to Lava Burst in SoD
			// Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				// Lightning Mastery does not currently apply to Lava Burst in SoD
				// CastTime: time.Second*2 - time.Millisecond*200*time.Duration(shaman.Talents.LightningMastery),
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				castTime := shaman.ApplyCastSpeedForSpell(cast.CastTime, spell)
				// The Enhancement 6pT3.5 makes Lava Burst not consume Maelstrom Weapon, so check that here
				if hasMaelstromWeaponRune && shaman.MaelstromWeaponClassMask&ClassSpellMask_ShamanLavaBurst > 0 {
					stacks := shaman.MaelstromWeaponAura.GetStacks()
					spell.SetMetricsSplit(stacks)
					if stacks > 0 {
						return
					}
				}

				if castTime > 0 {
					shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)

			critChanceBonusPct := 100.0
			if shaman.useLavaBurstCritScaling {
				critChanceBonusPct += shaman.GetStat(stats.SpellCrit)
			}

			flameShockActive := false
			for _, spell := range shaman.FlameShock {
				if spell == nil {
					continue
				}

				if spell.Dot(target).IsActive() {
					flameShockActive = true
					break
				}
			}

			if flameShockActive {
				spell.BonusCritRating += 100.0 * core.SpellCritRatingPerCritChance
			}

			spell.ApplyMultiplicativeDamageBonus(critChanceBonusPct / 100)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.ApplyMultiplicativeDamageBonus(100 / critChanceBonusPct)

			if flameShockActive {
				spell.BonusCritRating -= 100.0 * core.SpellCritRatingPerCritChance
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)

				if !isOverload && shaman.procOverload(sim, "Lava Burst Overload", 1) {
					shaman.LavaBurstOverload.Cast(sim, target)
				}
			})
		},
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
