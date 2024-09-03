package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (priest *Priest) registerVampiricEmbraceSpell() {
	if !priest.Talents.VampiricEmbrace {
		return
	}

	actionID := core.ActionID{SpellID: 15286}
	manaCost := 40.0
	duration := time.Minute * 1

	partyPlayers := priest.Env.Raid.GetPlayerParty(&priest.Unit).Players
	healthMetrics := priest.NewHealthMetrics(actionID)
	healthReturnedMultuplier := .10 + .05*float64(priest.Talents.ImprovedVampiricEmbrace)

	priest.VampiricEmbraceAuras = priest.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Vampiric Embrace (Health) - " + target.Label,
			Duration: duration,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					healthGained := result.Damage * healthReturnedMultuplier
					for _, player := range partyPlayers {
						player.GetCharacter().GainHealth(sim, healthGained, healthMetrics)
					}
				}
			},
		})
	})

	priest.VampiricEmbrace = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagPriest | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				priest.VampiricEmbraceAuras.Get(target).Activate(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
