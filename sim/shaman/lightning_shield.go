package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const LightningShieldRanks = 7

var LightningShieldSpellId = [LightningShieldRanks + 1]int32{0, 324, 325, 905, 945, 8134, 10431, 10432}
var LightningShieldProcSpellId = [LightningShieldRanks + 1]int32{0, 26364, 26365, 26366, 26367, 26369, 26370, 26363}
var LightningShieldOverchargedProcSpellId = [LightningShieldRanks + 1]int32{0, 432143, 432144, 432145, 432146, 432147, 432148, 432149}
var LightningShieldBaseDamage = [LightningShieldRanks + 1]float64{0, 13, 29, 51, 80, 114, 154, 198}
var LightningShieldSpellCoef = [LightningShieldRanks + 1]float64{0, .147, .227, .267, .267, .267, .267, .267}
var LightningShieldManaCost = [LightningShieldRanks + 1]float64{0, 45, 80, 125, 180, 240, 305}
var LightningShieldLevel = [LightningShieldRanks + 1]int{0, 8, 16, 24, 32, 40, 48, 56}

func (shaman *Shaman) registerLightningShieldSpell() {
	shaman.LightningShield = make([]*core.Spell, LightningShieldRanks+1)

	for rank := 1; rank <= LightningShieldRanks; rank++ {
		config := shaman.newLightningShieldSpellConfig(rank)

		if config.RequiredLevel <= int(shaman.Level) {
			shaman.LightningShield[rank] = shaman.RegisterSpell(config)
		}
	}
}

func (shaman *Shaman) newLightningShieldSpellConfig(rank int) core.SpellConfig {
	hasOverchargedRune := shaman.HasRune(proto.ShamanRune_RuneBracersOvercharged)
	hasRollingThunderRune := shaman.HasRune(proto.ShamanRune_RuneBracersRollingThunder)
	hasStaticShockRune := shaman.HasRune(proto.ShamanRune_RuneBracersStaticShock)

	impLightningShieldBonus := 1 + []float64{0, .05, .10, .15}[shaman.Talents.ImprovedLightningShield]

	spellId := LightningShieldSpellId[rank]
	procSpellId := core.Ternary(hasOverchargedRune, LightningShieldOverchargedProcSpellId, LightningShieldProcSpellId)[rank]
	baseDamage := LightningShieldBaseDamage[rank] * impLightningShieldBonus
	spellCoeff := LightningShieldSpellCoef[rank]
	manaCost := LightningShieldManaCost[rank]
	level := LightningShieldLevel[rank]

	if level > int(shaman.Level) {
		return core.SpellConfig{RequiredLevel: level}
	}

	staticShockProcChance := .06

	baseCharges := int32(3)
	maxCharges := int32(3)
	if hasRollingThunderRune {
		maxCharges = 9
	} else if hasStaticShockRune {
		baseCharges = 9
		maxCharges = 9
	}

	procSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: procSpellId},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	rollingThunder := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 432129},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			chargeDamage := baseDamage + spellCoeff*procSpell.SpellDamage()
			spell.CalcAndDealDamage(sim, target, chargeDamage, spell.OutcomeMagicCrit)
		},
	})

	// TODO: Does vanilla have an ICD?
	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: core.Ternary(hasOverchargedRune, time.Second*1, time.Millisecond*3500),
	}

	manaMetrics := shaman.NewManaMetrics(core.ActionID{SpellID: procSpellId})

	aura := shaman.RegisterAura(core.Aura{
		Label:     fmt.Sprintf("Lightning Shield (Rank %d)", rank),
		ActionID:  core.ActionID{SpellID: spellId},
		Duration:  time.Minute * 10,
		MaxStacks: maxCharges,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, baseCharges)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if hasStaticShockRune && spell.ProcMask.Matches(core.ProcMaskMelee) && sim.RandomFloat("Static Shock") < staticShockProcChance {
				aura.RemoveStack(sim)
				procSpell.Cast(sim, result.Target)
			}

			if hasRollingThunderRune && spell.SpellCode == SpellCode_ShamanEarthShock && aura.GetStacks() > 3 {
				multiplier := float64(aura.GetStacks() - baseCharges)
				rollingThunder.DamageMultiplier = multiplier
				rollingThunder.Cast(sim, result.Target)
				shaman.AddMana(sim, .02*multiplier*shaman.MaxMana(), manaMetrics)
				aura.SetStacks(sim, baseCharges)
			}
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.Landed() {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			icd.Use(sim)

			if hasOverchargedRune {
				// Deals damage to all targets within 8 yards and does not lose stacks
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					if aoeTarget.DistanceFromTarget <= 8 {
						procSpell.Cast(sim, aoeTarget)
					}
				}
			} else {
				aura.RemoveStack(sim)
			}
		},
	})

	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: spellId},
		Flags:    core.SpellFlagAPL,

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
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if shaman.LightningShieldAura != nil {
				shaman.LightningShieldAura.Deactivate(sim)
			}
			shaman.LightningShieldAura = aura
			shaman.LightningShieldAura.Activate(sim)
		},
	}
}

func (shaman *Shaman) rollRollingThunderCharge(sim *core.Simulation) {
	if shaman.LightningShieldAura.IsActive() && sim.RandomFloat("Rolling Thunder") < .30 {
		shaman.LightningShieldAura.AddStack(sim)
	}
}
