package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Totem Item IDs
const (
	StormfuryTotem           = 31031
	TotemOfAncestralGuidance = 32330
	TotemOfStorms            = 23199
	TotemOfTheVoid           = 28248
	TotemOfHex               = 40267
	VentureCoLightningRod    = 38361
	ThunderfallTotem         = 45255
)

// Shared precomputation logic for LB and CL.
func (shaman *Shaman) newElectricSpellConfig(actionID core.ActionID, baseCost float64, baseCastTime time.Duration, isOverload bool) core.SpellConfig {
	hasMaelstromWeaponRune := shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon)

	flags := core.SpellFlagNone
	if !isOverload {
		flags = core.SpellFlagAPL
	}

	spell := core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        flags,
		MetricSplits: MaelstromWeaponSplits,

		ManaCost: core.ManaCostOptions{
			FlatCost:   baseCost,
			Multiplier: 100 - 2*shaman.Talents.Convection,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: baseCastTime - time.Millisecond*200*time.Duration(shaman.Talents.LightningMastery),
				GCD:      core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				castTime := shaman.ApplyCastSpeedForSpell(cast.CastTime, spell)
				if hasMaelstromWeaponRune {
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

		BonusCritRating: []float64{0, 1, 2, 3, 4, 6}[shaman.Talents.CallOfThunder] * core.CritRatingPerCritChance,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
	}

	return spell
}
