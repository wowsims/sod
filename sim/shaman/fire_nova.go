package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

var FireNovaSpellId = [FireNovaTotemRanks + 1]int32{0, 408341, 408342, 408343, 408427, 408345}
var FireNovaSpellCoeff = [FireNovaTotemRanks + 1]float64{0, .214, .214, .214, .214, .214}

func (shaman *Shaman) applyFireNova() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistFireNova) {
		return
	}

	for rank := FireNovaTotemRanks; rank >= 1; rank-- {
		config := shaman.newFireNovaSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.FireNova = shaman.RegisterSpell(config)
			return
		}
	}
}

func (shaman *Shaman) newFireNovaSpellConfig(rank int) core.SpellConfig {
	spellId := FireNovaSpellId[rank]
	baseDamageLow := FireNovaTotemBaseDamage[rank][0]
	baseDamageHigh := FireNovaTotemBaseDamage[rank][1]
	spellCoeff := FireNovaTotemSpellCoeff[rank]
	cooldown := time.Second * 10
	manaCost := .22
	level := FireNovaTotemLevel[rank]

	spell := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellCode:   SpellCode_ShamanFireNova,
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagShaman | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: shaman.callOfFlameMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
				result := spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicCrit)

				spell.DealDamage(sim, result)
			}
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return shaman.ActiveTotems[FireTotem] != nil
		},
	}

	return spell
}
