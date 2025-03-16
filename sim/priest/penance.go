package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) RegisterPenanceSpell() {
	if !priest.HasRune(proto.PriestRune_RuneHandsPenance) {
		return
	}
	priest.Penance = priest.makePenanceSpell(false)
	priest.PenanceHeal = priest.makePenanceSpell(true)
}

// https://www.wowhead.com/classic/spell=402284/penance
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044
func (priest *Priest) makePenanceSpell(isHeal bool) *core.Spell {
	baseDamage := priest.baseRuneAbilityDamage() * 1.28
	baseHealing := priest.baseRuneAbilityHealing() * .85
	spellCoeff := 0.285
	manaCost := .16
	cooldown := time.Second * 12

	var classSpellMask uint64
	var procMask core.ProcMask
	flags := core.SpellFlagChanneled | core.SpellFlagAPL
	if isHeal {
		classSpellMask = ClassSpellMask_PriestPenanceHeal
		flags |= core.SpellFlagHelpful
		procMask = core.ProcMaskSpellHealing
	} else {
		classSpellMask = ClassSpellMask_PriestPenanceDamage
		procMask = core.ProcMaskSpellDamage
	}

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: 402284},
		SpellSchool:   core.SpellSchoolHoly,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      procMask,
		Flags:         flags,
		RequiredLevel: 1,

		ClassSpellMask: ClassSpellMask_PriestPenance | classSpellMask,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 0,

		Dot: core.Ternary(!isHeal, core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},
			NumberOfTicks:       2,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    spellCoeff,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.OutcomeTick)
			},
		}, core.DotConfig{}),
		Hot: core.Ternary(isHeal, core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},
			NumberOfTicks:       2,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    spellCoeff,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicHealing(sim, target, baseHealing, dot.Spell.OutcomeHealingCrit)
			},
		}, core.DotConfig{}),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if isHeal {
				hot := spell.Hot(target)
				hot.Apply(sim)
				// Do immediate tick
				hot.TickOnce(sim)
			} else {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					dot := spell.Dot(target)
					dot.Apply(sim)
					// Do immediate tick
					dot.TickOnce(sim)
				}
				spell.DealOutcome(sim, result)
			}
		},
	})
}
