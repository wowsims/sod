package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) getAspectOfTheHawkSpellConfig(rank int) core.SpellConfig {
	var impHawkAura *core.Aura
	improvedHawkProcChance := 0.01 * float64(hunter.Talents.ImprovedAspectOfTheHawk)

	spellId := [8]int32{0, 13165, 14318, 14319, 14320, 14321, 14322, 25296}[rank]
	rap := [8]float64{0, 20, 35, 50, 70, 90, 110, 120}[rank]
	//manaCost := [8]float64{0, 20, 35, 50, 70, 90, 110, 120}[rank]
	level := [8]int{0, 10, 18, 28, 38, 48, 58, 60}[rank]

	if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
		improvedHawkBonus := 1.3
		impHawkAura = hunter.GetOrRegisterAura(core.Aura{
			Label:    "Quick Shots",
			ActionID: core.ActionID{SpellID: 6150},
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyRangedSpeed(sim, improvedHawkBonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyRangedSpeed(sim, 1/improvedHawkBonus)
			},
		})
	}

	actionID := core.ActionID{SpellID: spellId}
	aspectOfTheHawkAura := hunter.NewTemporaryStatsAuraWrapped(
		"Aspect of the Hawk"+strconv.Itoa(rank),
		actionID,
		stats.Stats{
			stats.RangedAttackPower: rap,
		},
		core.NeverExpires,
		func(aura *core.Aura) {
			aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell != hunter.AutoAttacks.RangedAuto() {
					return
				}

				if impHawkAura != nil && sim.RandomFloat("Imp Aspect of the Hawk") < improvedHawkProcChance {
					impHawkAura.Activate(sim)
				}
			}
		})

	return core.SpellConfig{
		ActionID:      actionID,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.curAspect != aspectOfTheHawkAura
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if hunter.curAspect != nil {
				hunter.curAspect.Deactivate(sim)
			}

			hunter.curAspect = aspectOfTheHawkAura
			hunter.curAspect.Activate(sim)
		},
	}
}

func (hunter *Hunter) registerAspectOfTheHawkSpell() {
	maxRank := 7

	for i := 1; i <= maxRank; i++ {
		config := hunter.getAspectOfTheHawkSpellConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.GetOrRegisterSpell(config)
		}
	}
}

func (hunter *Hunter) registerAspectOfTheViperSpell() {
	actionID := core.ActionID{SpellID: 415423}
	manaMetrics := hunter.NewManaMetrics(actionID)

	var manaPA *core.PendingAction

	baseManaRegenMultiplier := 0.01

	aspectOfTheViperAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Aspect of the Viper",
		ActionID: actionID,
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.DamageDealtMultiplier *= 0.9

			manaPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 3,
				OnAction: func(s *core.Simulation) {
					hunter.AddMana(sim, hunter.MaxMana()*0.1, manaMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.DamageDealtMultiplier /= 0.9
			manaPA.Cancel(sim)
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == hunter.AutoAttacks.RangedAuto() {
				manaPerRangedHitMultiplier := baseManaRegenMultiplier * hunter.AutoAttacks.Ranged().SwingSpeed
				hunter.AddMana(sim, hunter.MaxMana()*manaPerRangedHitMultiplier, manaMetrics)
			} else if spell == hunter.AutoAttacks.MHAuto() {
				manaPerMHHitMultiplier := baseManaRegenMultiplier * hunter.AutoAttacks.MH().SwingSpeed
				hunter.AddMana(sim, hunter.MaxMana()*manaPerMHHitMultiplier, manaMetrics)
			} else if spell == hunter.AutoAttacks.OHAuto() {
				manaPerOHHitMultiplier := baseManaRegenMultiplier * hunter.AutoAttacks.OH().SwingSpeed
				hunter.AddMana(sim, hunter.MaxMana()*manaPerOHHitMultiplier, manaMetrics)
			}
		},
	})

	hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.curAspect != aspectOfTheViperAura
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if hunter.curAspect != nil {
				hunter.curAspect.Deactivate(sim)
			}

			hunter.curAspect = aspectOfTheViperAura
			hunter.curAspect.Activate(sim)
		},
	})
}
