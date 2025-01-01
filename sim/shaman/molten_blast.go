package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) applyMoltenBlast() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsMoltenBlast) {
		return
	}

	baseDamageLow := shaman.baseRuneAbilityDamage() * .72
	baseDamageHigh := shaman.baseRuneAbilityDamage() * 1.08
	apCoef := .10
	spCoef := .14
	cooldown := time.Second * 6
	manaCost := .18
	targetCount := int32(10)

	flameShockResetChance := 0.10

	numHits := min(targetCount, shaman.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	shaman.MoltenBlast = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsMoltenBlast)},
		SpellCode:   SpellCode_ShamanMoltenBlast,
		SpellSchool: core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagShaman | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
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

		BonusCoefficient: spCoef,
		DamageMultiplier: 1,
		ThreatMultiplier: 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx := range results {
				// Molten Blast is a magic ability but scales off of Attack Power
				baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + apCoef*spell.MeleeAttackPower()
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Molten Blast Reset Trigger",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(flameShockResetChance, "Molten Blast Reset") {
				shaman.MoltenBlast.CD.Reset()
			}
		},
	}))
}
