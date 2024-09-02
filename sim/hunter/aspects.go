package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// Utility function to create an Improved Hawk Aura
func (hunter *Hunter) createImprovedHawkAura(auraLabel string, actionID core.ActionID, isMelee bool) *core.Aura {
    bonusMultiplier := 1.3
    return hunter.GetOrRegisterAura(core.Aura{
        Label:    auraLabel,
        ActionID: actionID,
        Duration: time.Second * 12,
        OnGain: func(aura *core.Aura, sim *core.Simulation) {
            if isMelee {
                aura.Unit.MultiplyMeleeSpeed(sim, bonusMultiplier)
            } else {
                aura.Unit.MultiplyRangedSpeed(sim, bonusMultiplier)
            }
        },
        OnExpire: func(aura *core.Aura, sim *core.Simulation) {
            if isMelee {
                aura.Unit.MultiplyMeleeSpeed(sim, 1/bonusMultiplier)
            } else {
                aura.Unit.MultiplyRangedSpeed(sim, 1/bonusMultiplier)
            }
        },
    })
}

// Utility function to create a generic Aspect aura
func (hunter *Hunter) createAspectAura(auraLabel string, actionID core.ActionID, stats stats.Stats, onSpellHitDealt func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult)) *core.Aura {
    aspectAura := hunter.NewTemporaryStatsAuraWrapped(
        auraLabel,
        actionID,
        stats,
        core.NeverExpires,
        func(aura *core.Aura) {
            aura.OnSpellHitDealt = onSpellHitDealt
        },
    )
    aspectAura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})
    return aspectAura
}

// Config for Aspect of the Hawk
func (hunter *Hunter) getAspectOfTheHawkSpellConfig(rank int) core.SpellConfig {
    improvedHawkProcChance := 0.01 * float64(hunter.Talents.ImprovedAspectOfTheHawk)
    spellId := [8]int32{0, 13165, 14318, 14319, 14320, 14321, 14322, 25296}[rank]
    rap := [8]float64{0, 20, 35, 50, 70, 90, 110, 120}[rank]
    level := [8]int{0, 10, 18, 28, 38, 48, 58, 60}[rank]

    var quickShotsAura *core.Aura
    if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
        quickShotsAura = hunter.createImprovedHawkAura("Quick Shots", core.ActionID{SpellID: 6150}, false)
    }

    actionID := core.ActionID{SpellID: spellId}
    aspectAura := hunter.createAspectAura(
        "Aspect of the Hawk"+strconv.Itoa(rank),
        actionID,
        stats.Stats{stats.RangedAttackPower: rap},
        func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
            if spell != hunter.AutoAttacks.RangedAuto() {
                return
            }
            if quickShotsAura != nil && sim.RandomFloat("Imp Aspect of the Hawk") < improvedHawkProcChance {
                quickShotsAura.Activate(sim)
            }
        },
    )

    return core.SpellConfig{
        ActionID:      actionID,
        Flags:         core.SpellFlagAPL,
        Rank:          rank,
        RequiredLevel: level,

        Cast: core.CastConfig{
            DefaultCast: core.Cast{
                GCD: core.GCDDefault,
            },
        },
        ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
            return !aspectAura.IsActive()
        },

        ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
            aspectAura.Activate(sim)
        },
    }
}

// Config for Aspect of the Falcon
func (hunter *Hunter) getAspectOfTheFalconSpellConfig() core.SpellConfig {
    highestHawkRank := hunter.getHighestAspectOfTheHawkRank()
    if highestHawkRank == 0 {
        return core.SpellConfig{}
    }

    hawkConfig := hunter.getAspectOfTheHawkSpellConfig(highestHawkRank)
    improvedHawkProcChance := 0.01 * float64(hunter.Talents.ImprovedAspectOfTheHawk)

    quickStrikesAura := hunter.createImprovedHawkAura("Quick Strikes", core.ActionID{SpellID: 469144}, true)

    return core.SpellConfig{
        ActionID:      core.ActionID{SpellID: 469145},
        Flags:         core.SpellFlagAPL,
        Rank:          highestHawkRank,
        RequiredLevel: hawkConfig.RequiredLevel,

        Cast: core.CastConfig{
            DefaultCast: core.Cast{
                GCD: core.GCDDefault,
            },
        },
        ExtraCastCondition: hawkConfig.ExtraCastCondition,
        ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
            aspectAura := hunter.createAspectAura(
                "Aspect of the Falcon",
                core.ActionID{SpellID: 469145},
                stats.Stats{
                    stats.RangedAttackPower: [8]float64{0, 20, 35, 50, 70, 90, 110, 120}[highestHawkRank],
                    stats.AttackPower:  [8]float64{0, 20, 35, 50, 70, 90, 110, 120}[highestHawkRank],
                },
                func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
                    if sim.RandomFloat("Imp Aspect of the Hawk") < improvedHawkProcChance {
                        if spell == hunter.AutoAttacks.RangedAuto() || spell == hunter.AutoAttacks.MeleeAuto() {
                            if spell == hunter.AutoAttacks.MeleeAuto() {
                                quickStrikesAura.Activate(sim)
                            }
                        }
                    }
                },
            )
            aspectAura.Activate(sim)
            if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
                quickStrikesAura.Activate(sim)
            }
        },
    }
}

// Register Aspect of the Hawk spell configurations
func (hunter *Hunter) registerAspectOfTheHawkSpell() {
    maxRank := 7

    for i := 1; i < maxRank; i++ {
        config := hunter.getAspectOfTheHawkSpellConfig(i)
        if config.RequiredLevel <= int(hunter.Level) {
            hunter.GetOrRegisterSpell(config)
        }
    }
}

// Register Aspect of the Falcon spell configuration
func (hunter *Hunter) registerAspectOfTheFalconSpell() {
    config := hunter.getAspectOfTheFalconSpellConfig()
    if config.ActionID.SpellID != 0 && config.RequiredLevel <= int(hunter.Level) {
        hunter.GetOrRegisterSpell(config)
    }
}





func (hunter *Hunter) registerAspectOfTheViperSpell() {
	actionID := core.ActionID{SpellID: 415423}
	manaMetrics := hunter.NewManaMetrics(actionID)

	var manaPA *core.PendingAction

	baseManaRegenMultiplier := 0.02

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

	aspectOfTheViperAura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})

	hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !aspectOfTheViperAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aspectOfTheViperAura.Activate(sim)
		},
	})
}
