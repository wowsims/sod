package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) applyLavaBurst() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
		return
	}

	shaman.LavaBurst = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(false))

	if shaman.HasRune(proto.ShamanRune_RuneChestOverload) {
		shaman.LavaBurstOverload = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(true))
	}
}

func (shaman *Shaman) newLavaBurstSpellConfig(isOverload bool) core.SpellConfig {
	level := float64(shaman.GetCharacter().Level)
	spellId := core.TernaryInt32(isOverload, 408491, 408490)
	baseCalc := 7.583798 + 0.471881*level + 0.036599*level*level
	baseLowDamage := baseCalc * 4.69
	baseHighDamage := baseCalc * 6.05
	spellCoeff := .571

	flags := SpellFlagFocusable
	if !isOverload {
		flags |= core.SpellFlagAPL
	}

	canOverload := !isOverload && shaman.OverloadAura != nil

	spell := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,

		ManaCost: core.ManaCostOptions{
			BaseCost: .10,
			// Convection does not currently apply to Lava Burst in SoD
			// Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				// Lightning Mastery does not currently apply to Lava Burst in SoD
				// CastTime: time.Second*2 - time.Millisecond*200*time.Duration(shaman.Talents.LightningMastery),
				CastTime: time.Second * 2,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		// Concussion does not currently apply to Lava Burst in SoD
		// DamageMultiplier: 1 + 0.01*float64(shaman.Talents.Concussion)
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if canOverload && result.Landed() && sim.RandomFloat("LvB Overload") < shaman.OverloadChance {
				shaman.LavaBurstOverload.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
