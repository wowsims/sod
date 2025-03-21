package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const SearingTotemRanks = 6

var SearingTotemSpellId = [SearingTotemRanks + 1]int32{0, 3599, 6363, 6364, 6365, 10437, 10438}
var SearingTotemAttackSpellId = [SearingTotemRanks + 1]int32{0, 3606, 6350, 6351, 6352, 10435, 10436}
var SearingTotemBaseDamage = [SearingTotemRanks + 1][]float64{{0}, {9, 11}, {13, 17}, {19, 25}, {26, 34}, {33, 45}, {40, 54}}
var SearingTotemSpellCoef = [SearingTotemRanks + 1]float64{0, .052, .083, .083, .083, .083, .083}
var SearingTotemManaCost = [SearingTotemRanks + 1]float64{0, 25, 45, 75, 110, 145, 170}
var SearingTotemDuration = [SearingTotemRanks + 1]int{0, 30, 35, 40, 45, 50, 55}
var SearingTotemLevel = [SearingTotemRanks + 1]int{0, 10, 20, 30, 40, 50, 60}

func (shaman *Shaman) registerSearingTotemSpell() {
	shaman.SearingTotem = make([]*core.Spell, SearingTotemRanks+1)

	for rank := 1; rank <= SearingTotemRanks; rank++ {
		config := shaman.newSearingTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.SearingTotem[rank] = shaman.RegisterSpell(config)
		}
	}

	shaman.FireTotems = append(
		shaman.FireTotems,
		core.FilterSlice(shaman.SearingTotem, func(spell *core.Spell) bool { return spell != nil })...,
	)
}

func (shaman *Shaman) newSearingTotemSpellConfig(rank int) core.SpellConfig {
	totemSpellId := SearingTotemSpellId[rank]
	baseDamageLow := SearingTotemBaseDamage[rank][0]
	baseDamageHigh := SearingTotemBaseDamage[rank][1]
	spellCoeff := SearingTotemSpellCoef[rank]
	manaCost := SearingTotemManaCost[rank]
	duration := time.Second * time.Duration(SearingTotemDuration[rank])
	level := SearingTotemLevel[rank]

	attackInterval := time.Millisecond * 2500

	attackSpell := shaman.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanSearingTotemAttack,
		ActionID:       core.ActionID{SpellID: SearingTotemAttackSpellId[rank]},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,

		DamageMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})

	spell := core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanSearingTotem,
		ActionID:       core.ActionID{SpellID: totemSpellId},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: shaman.totemManaMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Searing Totem (Rank %d)", rank),
			},
			// These are the real tick values, but searing totem doesn't start its next
			// cast until the previous missile hits the target. We don't have an option
			// for target distance yet so just pretend the tick rate is lower.
			// https://wotlk.wowhead.com/spell=25530/attack
			//TickLength:           time.Second * 2.2,
			NumberOfTicks: int32(duration / attackInterval),
			TickLength:    attackInterval,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				attackSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if shaman.ActiveTotems[FireTotem] != nil {
				shaman.ActiveTotems[FireTotem].Dot(sim.GetTargetUnit(0)).Cancel(sim)
			}
			spell.Dot(sim.GetTargetUnit(0)).Apply(sim)
			// +1 needed because of rounding issues with totem tick time.
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + duration + 1
			shaman.ActiveTotems[FireTotem] = spell
		},
	}

	return spell
}

const MagmaTotemRanks = 4

var MagmaTotemSpellId = [MagmaTotemRanks + 1]int32{0, 8190, 10585, 10586, 10587}
var MagmaTotemAoeSpellId = [MagmaTotemRanks + 1]int32{0, 8187, 10579, 10580, 10581}
var MagmaTotemBaseDamage = [MagmaTotemRanks + 1]float64{0, 22, 37, 54, 75}
var MagmaTotemSpellCoeff = [MagmaTotemRanks + 1]float64{0, .033, .033, .033, .033}
var MagmaTotemManaCost = [MagmaTotemRanks + 1]float64{0, 230, 360, 500, 650}
var MagmaTotemLevel = [MagmaTotemRanks + 1]int{0, 26, 36, 46, 56}

func (shaman *Shaman) registerMagmaTotemSpell() {
	shaman.MagmaTotem = make([]*core.Spell, MagmaTotemRanks+1)

	for rank := 1; rank <= MagmaTotemRanks; rank++ {
		config := shaman.newMagmaTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.MagmaTotem[rank] = shaman.RegisterSpell(config)
		}
	}

	shaman.FireTotems = append(
		shaman.FireTotems,
		core.FilterSlice(shaman.MagmaTotem, func(spell *core.Spell) bool { return spell != nil })...,
	)
}

func (shaman *Shaman) newMagmaTotemSpellConfig(rank int) core.SpellConfig {
	spellId := MagmaTotemSpellId[rank]
	baseDamage := MagmaTotemBaseDamage[rank]
	spellCoeff := MagmaTotemSpellCoeff[rank]
	manaCost := MagmaTotemManaCost[rank]
	level := MagmaTotemLevel[rank]

	duration := time.Second * 20
	attackInterval := time.Second * 2

	aoeSpell := shaman.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanMagmaTotemAttack,
		ActionID:       core.ActionID{SpellID: MagmaTotemAoeSpellId[rank]},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,

		DamageMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	spell := core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanMagmaTotem,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: shaman.totemManaMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Magma Totem (Rank %d)", rank),
			},
			NumberOfTicks: int32(duration / attackInterval),
			TickLength:    attackInterval,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				aoeSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if shaman.ActiveTotems[FireTotem] != nil {
				shaman.ActiveTotems[FireTotem].Dot(sim.GetTargetUnit(0)).Cancel(sim)
			}
			spell.Dot(sim.GetTargetUnit(0)).Apply(sim)
			// +1 needed because of rounding issues with totem tick time.
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + duration + 1
			shaman.ActiveTotems[FireTotem] = spell
		},
	}

	return spell
}

const FireNovaTotemRanks = 5

var FireNovaTotemSpellId = [FireNovaTotemRanks + 1]int32{0, 1535, 8498, 8499, 11314, 11315}
var FireNovaTotemAoeSpellId = [FireNovaTotemRanks + 1]int32{0, 8349, 8502, 8503, 11306, 11307}
var FireNovaTotemBaseDamage = [FireNovaTotemRanks + 1][]float64{{0, 0}, {53, 62}, {110, 124}, {195, 219}, {295, 331}, {413, 459}}
var FireNovaTotemSpellCoeff = [FireNovaTotemRanks + 1]float64{0, .1, .143, .143, .143, .143}
var FireNovaTotemManaCost = [FireNovaTotemRanks + 1]float64{0, 95, 170, 280, 395, 520}
var FireNovaTotemLevel = [FireNovaTotemRanks + 1]int{0, 12, 22, 32, 42, 52}

func (shaman *Shaman) registerFireNovaTotemSpell() {
	if shaman.HasRune(proto.ShamanRune_RuneWaistFireNova) {
		return
	}

	shaman.FireNovaTotem = make([]*core.Spell, FireNovaTotemRanks+1)

	for rank := 1; rank <= FireNovaTotemRanks; rank++ {
		config := shaman.newFireNovaTotemSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.FireNovaTotem[rank] = shaman.RegisterSpell(config)
		}
	}

	shaman.FireTotems = append(
		shaman.FireTotems,
		core.FilterSlice(shaman.FireNovaTotem, func(spell *core.Spell) bool { return spell != nil })...,
	)
}

func (shaman *Shaman) newFireNovaTotemSpellConfig(rank int) core.SpellConfig {
	spellId := FireNovaTotemSpellId[rank]
	baseDamageLow := FireNovaTotemBaseDamage[rank][0]
	baseDamageHigh := FireNovaTotemBaseDamage[rank][1]
	spellCoeff := FireNovaTotemSpellCoeff[rank]
	cooldown := time.Second * 15
	manaCost := FireNovaTotemManaCost[rank]
	level := FireNovaTotemLevel[rank]

	duration := time.Second * 5
	attackInterval := duration

	novaSpell := shaman.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanFireNovaTotemAttack,
		ActionID:       core.ActionID{SpellID: FireNovaTotemAoeSpellId[rank]},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,

		DamageMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	spell := core.SpellConfig{
		ClassSpellMask: ClassSpellMask_ShamanFireNovaTotem,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: shaman.totemManaMultiplier(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Fire Nova Totem (Rank %d)", rank),
			},
			NumberOfTicks: 1,
			TickLength:    attackInterval,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				novaSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if shaman.ActiveTotems[FireTotem] != nil {
				shaman.ActiveTotems[FireTotem].Dot(sim.GetTargetUnit(0)).Cancel(sim)
			}
			dot := spell.Dot(sim.GetTargetUnit(0))
			if dot.Duration > 0 {
				dot.Apply(sim)
			} else {
				novaSpell.Cast(sim, target)
			}
			// +1 needed because of rounding issues with totem tick time.
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + duration + 1
			shaman.ActiveTotems[FireTotem] = spell
		},
	}

	return spell
}
