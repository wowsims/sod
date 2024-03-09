package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ArcaneMissilesRanks = 8

var ArcaneMissilesSpellId = [ArcaneMissilesRanks + 1]int32{0, 5143, 5144, 5145, 8416, 8417, 10211, 10212, 25345}
var ArcaneMissilesBaseTickDamage = [ArcaneMissilesRanks + 1]float64{0, 26, 38, 57, 86, 115, 153, 196, 230}
var ArcaneMissilesSpellCoeff = [ArcaneMissilesRanks + 1]float64{0, .132, .204, .24, .24, .24, .24, .24, .24}
var ArcaneMissilesCastTime = [ArcaneMissilesRanks + 1]int32{0, 3, 4, 5, 5, 5, 5, 5, 5}
var ArcaneMissilesManaCost = [ArcaneMissilesRanks + 1]float64{0, 85, 140, 235, 320, 410, 500, 595, 655}
var ArcaneMissilesLevel = [ArcaneMissilesRanks + 1]int{0, 8, 16, 24, 32, 40, 48, 56, 56}

func (mage *Mage) registerArcaneMissilesSpell() {
	mage.ArcaneMissiles = make([]*core.Spell, ArcaneMissilesRanks+1)
	mage.ArcaneMissilesTickSpell = make([]*core.Spell, ArcaneMissilesRanks+1)

	for rank := 1; rank < ArcaneMissilesRanks; rank++ {
		config := mage.getArcaneMissilesSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.ArcaneMissiles[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) getArcaneMissilesSpellConfig(rank int) core.SpellConfig {
	spellId := ArcaneMissilesSpellId[rank]
	baseTickDamage := ArcaneMissilesBaseTickDamage[rank]
	spellCoeff := ArcaneMissilesSpellCoeff[rank]
	castTime := ArcaneMissilesCastTime[rank]
	manaCost := ArcaneMissilesManaCost[rank]
	level := ArcaneMissilesLevel[rank]

	numTicks := castTime
	tickLength := time.Second

	hasArcaneBlastRune := mage.HasRune(proto.MageRune_RuneHandsArcaneBlast)
	hasMissileBarrageRune := mage.HasRune(proto.MageRune_RuneBeltMissileBarrage)

	tickSpell := mage.getArcaneMissilesTickSpell(rank)
	mage.ArcaneMissilesTickSpell[rank] = tickSpell

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellCode:   SpellCode_MageArcaneMissiles,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL | core.SpellFlagChanneled | core.SpellFlagNoMetrics,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("ArcaneMissiles-%d-%d", +rank, numTicks),
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					// TODO: This check is necessary to ensure the final tick occurs before
					// Arcane Blast stacks are dropped. To fix this, ticks need to reliably
					// occur before aura expirations.
					dot := mage.ArcaneMissiles[rank].Dot(aura.Unit)
					if dot.TickCount < dot.NumberOfTicks {
						dot.TickCount++
						dot.TickOnce(sim)
					}
					if hasArcaneBlastRune && mage.ArcaneBlastAura.IsActive() {
						mage.ArcaneBlastAura.Deactivate(sim)
					}
				},
			},
			NumberOfTicks: numTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				tickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)

			if hasMissileBarrageRune && mage.MissileBarrageAura.IsActive() {
				mage.MissileBarrageAura.Deactivate(sim)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := baseTickDamage + (spellCoeff * spell.SpellDamage())
			return tickSpell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	}
}

func (mage *Mage) getArcaneMissilesTickSpell(rank int) *core.Spell {
	spellId := ArcaneMissilesSpellId[rank]
	baseTickDamage := ArcaneMissilesBaseTickDamage[rank]
	spellCoeff := ArcaneMissilesSpellCoeff[rank]

	return mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: spellId}.WithTag(1),
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskProc | core.ProcMaskNotInSpellbook,
		Flags:        SpellFlagMage,
		MissileSpeed: 20,

		DamageMultiplier: 1,
		CritMultiplier:   mage.MageCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseTickDamage + (spellCoeff * spell.SpellDamage())
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
