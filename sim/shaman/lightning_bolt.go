package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) registerLightningBoltSpell() {
	maxRank := 10
	overloadRuneEquipped := shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	for i := 1; i <= maxRank; i++ {
		config := shaman.newLightningBoltSpellConfig(i, false)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.LightningBolt = shaman.RegisterSpell(config)

			if overloadRuneEquipped {
				shaman.LightningBoltOverload = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(i, true))
			}
		}
	}
}

func (shaman *Shaman) newLightningBoltSpellConfig(rank int, isOverload bool) core.SpellConfig {
	// TODO: Get updated spell ranks + separate out Lightning Bolt overload ids + damage
	spellCoeff := [9]float64{0, 0.123, 0.231, 0.443, 0.571, 0.571, 0.571, 0.571, 0.571}[rank]
	baseDamage := [9][]float64{{0}, {13, 16}, {28, 33}, {48, 57}, {69, 79}, {108, 123}, {148, 167}, {198, 221}, {248, 277}}[rank]
	spellId := [9]int32{0, 5176, 5177, 5178, 5179, 5180, 6780, 8905, 9912}[rank]
	manaCost := [9]float64{0, 20, 35, 55, 70, 100, 125, 155, 180}[rank]
	level := [9]int{0, 1, 6, 14, 22, 30, 38, 46, 54}[rank]
	castTime := [9]int{0, 1500, 1700, 2000, 2000, 2000, 2000, 2000, 2000}[rank]

	spell := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: spellId},
		manaCost,
		time.Millisecond*time.Duration(castTime),
		isOverload,
	)

	dmgBonus := shaman.electricSpellBonusDamage(0.7143)

	canOverload := !isOverload && shaman.OverloadAura != nil

	spell.RequiredLevel = level
	spell.Rank = rank
	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := dmgBonus + sim.Roll(719, 819) + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		if canOverload && result.Landed() && sim.RandomFloat("LB Overload") < shaman.OverloadChance {
			shaman.LightningBoltOverload.Cast(sim, target)
		}

		spell.DealDamage(sim, result)
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell, CastTagLightningBoltOverload)
	}

	return spell
}

func (shaman *Shaman) getLightningBoltBaseConfig(
	actionId core.ActionID,
	spellCoeff float64,
	baseDamageLow float64,
	baseDamageHigh float64,
) core.SpellConfig {

}
