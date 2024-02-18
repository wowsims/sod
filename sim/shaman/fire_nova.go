package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// import (
// 	"time"

// 	"github.com/wowsims/sod/sim/core"
// )

var FireNovaSpellId = [FireNovaTotemRanks + 1]int32{0, 408341, 408342, 408343, 408344, 408345}
var FireNovaSpellCoeff = [FireNovaTotemRanks + 1]float64{0, .214, .214, .214, .214, .214}

func (shaman *Shaman) applyFireNova() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistFireNova) {
		return
	}

	shaman.FireNova = make([]*core.Spell, FireNovaTotemRanks+1)

	for rank := 1; rank <= FireNovaTotemRanks; rank++ {
		config := shaman.newFireNovaSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.FireNova[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newFireNovaSpellConfig(rank int) core.SpellConfig {
	spellId := FireNovaSpellId[rank]
	baseDamageLow := FireNovaTotemBaseDamage[rank][0]
	baseDamageHigh := FireNovaTotemBaseDamage[rank][1]
	// Verify spell coef?
	spellCoeff := FireNovaTotemSpellCoeff[rank]
	cooldown := time.Second * 10
	manaCost := .22
	level := FireNovaTotemLevel[rank]

	spell := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellCode:   SpellCode_ShamanFireNova,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagFocusable | core.SpellFlagAPL,

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

		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: shaman.ShamanThreatMultiplier(1),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.DealDamage(sim, result)
		},
	}

	return spell
}
