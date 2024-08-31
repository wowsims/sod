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
	shaman.LightningShieldProcs = make([]*core.Spell, LightningShieldRanks+1)

	for rank := 1; rank <= LightningShieldRanks; rank++ {
		level := LightningShieldLevel[rank]

		if level <= int(shaman.Level) {
			shaman.registerNewLightningShieldSpell(rank)
		}
	}
}

func (shaman *Shaman) registerNewLightningShieldSpell(rank int) {
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

	baseCharges := int32(3)
	maxCharges := int32(3)
	if hasRollingThunderRune {
		maxCharges = 9
	} else if hasStaticShockRune {
		baseCharges = 9
		maxCharges = 9
	}

	shaman.LightningShieldProcs[rank] = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: procSpellId},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagShaman | SpellFlagLightning,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	// TODO: Does vanilla have an ICD?
	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: core.Ternary(hasOverchargedRune, time.Second*3, time.Millisecond*3500),
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

			if hasStaticShockRune && spell.ProcMask.Matches(core.ProcMaskMelee) {
				staticShockProcChance := core.TernaryFloat64(shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand, .12, .06)
				if sim.RandomFloat("Static Shock") < staticShockProcChance {
					aura.RemoveStack(sim)
					shaman.LightningShieldProcs[rank].Cast(sim, result.Target)
				}
			}

			if hasRollingThunderRune && spell.SpellCode == SpellCode_ShamanEarthShock && aura.GetStacks() > 3 {
				multiplier := float64(aura.GetStacks() - baseCharges)
				shaman.RollingThunder.DamageMultiplier = multiplier
				shaman.RollingThunder.Cast(sim, result.Target)
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
						shaman.LightningShieldProcs[rank].Cast(sim, aoeTarget)
					}
				}
			} else {
				aura.RemoveStack(sim)
				shaman.LightningShieldProcs[rank].Cast(sim, spell.Unit)
			}
		},
	})

	shaman.LightningShield[rank] = shaman.RegisterSpell(core.SpellConfig{
		ActionID:  core.ActionID{SpellID: spellId},
		SpellCode: SpellCode_ShamanLightningShield,
		ProcMask:  core.ProcMaskEmpty,
		Flags:     core.SpellFlagAPL | SpellFlagShaman | SpellFlagLightning,

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
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if shaman.ActiveShieldAura != nil {
				shaman.ActiveShieldAura.Deactivate(sim)
			}
			shaman.ActiveShield = spell
			shaman.ActiveShieldAura = aura
			shaman.ActiveShieldAura.Activate(sim)
		},
	})
}
