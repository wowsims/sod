package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
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

const (
	// This could be value or bitflag if we ended up needing multiple flags at the same time.
	//1 to 5 are used by MaelstromWeapon Stacks
	CastTagLightningOverload int32 = 6
)

// Shared precomputation logic for LB and CL.
func (shaman *Shaman) newElectricSpellConfig(actionID core.ActionID, baseCost float64, baseCastTime time.Duration, isOverload bool) core.SpellConfig {
	flags := SpellFlagElectric | SpellFlagFocusable | SpellFlagMaelstrom
	if !isOverload {
		flags |= core.SpellFlagAPL
	}

	spell := core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        flags,
		MetricSplits: 6,

		ManaCost: core.ManaCostOptions{
			FlatCost:   baseCost,
			Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: baseCastTime - time.Millisecond*200*time.Duration(shaman.Talents.LightningMastery),
				GCD:      core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(shaman.MaelstromWeaponAura.GetStacks())
				castTime := shaman.ApplyCastSpeedForSpell(cast.CastTime, spell)

				if castTime > 0 {
					shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
				}
			},
		},

		BonusCritRating: 0 +
			float64(shaman.Talents.TidalMastery)*core.CritRatingPerCritChance +
			float64(shaman.Talents.CallOfThunder)*core.CritRatingPerCritChance,

		CritDamageBonus: shaman.elementalFury(),

		DamageMultiplier: shaman.concussionMultiplier(),
		ThreatMultiplier: 1,
	}

	return spell
}
