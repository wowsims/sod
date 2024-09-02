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

// Function to get the maximum attack power for Aspect of the Hawk based on character level
func (hunter *Hunter) getMaxAspectOfTheHawkAttackPower(level int) float64 {
    attackPower := [8]float64{0, 20, 35, 50, 70, 90, 110, 120} // Static data for attack power per rank
    levels := [8]int{0, 10, 18, 28, 38, 48, 58, 60} // Levels at which ranks are available

    maxAttackPower := 0.0

    for rank := 1; rank <= 7; rank++ {
        if level >= levels[rank] {
            maxAttackPower = attackPower[rank]
        }
    }
    return maxAttackPower
}

// Configuration for Aspect of the Hawk spell
func (hunter *Hunter) getAspectOfTheHawkSpellConfig(rank int) core.SpellConfig {
    var impHawkAura *core.Aura
    improvedHawkProcChance := 0.01 * float64(hunter.Talents.ImprovedAspectOfTheHawk)

    spellIds := [8]int32{0, 13165, 14318, 14319, 14320, 14321, 14322, 25296}
    levels := [8]int{0, 10, 18, 28, 38, 48, 58, 60}

    spellId := spellIds[rank]
    level := levels[rank]

    if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
        impHawkAura = hunter.createImprovedHawkAura(
            "Quick Shots",
            core.ActionID{SpellID: 6150},
            false, // Ranged
        )
    }

    // Use utility function to get the attack power based on rank
    rap := hunter.getMaxAspectOfTheHawkAttackPower(level)

    actionID := core.ActionID{SpellID: spellId}
    aspectOfTheHawkAura := hunter.NewTemporaryStatsAuraWrapped(
        "Aspect of the Hawk"+strconv.Itoa(rank),
        actionID,
        stats.Stats{stats.RangedAttackPower: rap},
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

    aspectOfTheHawkAura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})

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
            return !aspectOfTheHawkAura.IsActive()
        },

        ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
            aspectOfTheHawkAura.Activate(sim)
        },
    }
}

// Register Aspect of the Hawk spells
func (hunter *Hunter) registerAspectOfTheHawkSpell() {
    for i := 1; i <= 7; i++ {
        config := hunter.getAspectOfTheHawkSpellConfig(i)

        if config.RequiredLevel <= int(hunter.Level) {
            hunter.GetOrRegisterSpell(config)
        }
    }
}

// Configuration for Aspect of the Falcon spell
func (hunter *Hunter) getAspectOfTheFalconSpellConfig(level int) core.SpellConfig {
    // Get the maximum attack power from Aspect of the Hawk for the given level
    maxAttackPower := hunter.getMaxAspectOfTheHawkAttackPower(level)

    actionID := core.ActionID{SpellID: 469145}
    aspectOfTheFalconAura := hunter.NewTemporaryStatsAuraWrapped(
        "Aspect of the Falcon",
        actionID,
        stats.Stats{
            stats.RangedAttackPower: maxAttackPower,
            stats.AttackPower:       maxAttackPower,
        },
        core.NeverExpires,
        nil, // No special proc effects
    )

    var impHawkAura *core.Aura
    if hunter.Talents.ImprovedAspectOfTheHawk > 0 {
        impHawkAura = hunter.createImprovedHawkAura(
            "Quick Strikes",
            core.ActionID{SpellID: 469144},
            true, // Melee
        )
    }

    return core.SpellConfig{
        ActionID:      actionID,
        Flags:         core.SpellFlagAPL,
        Rank:          1, // Single rank
        RequiredLevel: level,

        Cast: core.CastConfig{
            DefaultCast: core.Cast{
                GCD: core.GCDDefault,
            },
        },
        ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
            return !aspectOfTheFalconAura.IsActive()
        },

        ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
            aspectOfTheFalconAura.Activate(sim)

            // Handle Improved Hawk talent proc for melee attacks
            if impHawkAura != nil && sim.RandomFloat("Imp Aspect of the Hawk") < (0.01 * float64(hunter.Talents.ImprovedAspectOfTheHawk)) {
                impHawkAura.Activate(sim)
            }
        },
    }
}

// Register Aspect of the Falcon spell
func (hunter *Hunter) registerAspectOfTheFalconSpell() {
    config := hunter.getAspectOfTheFalconSpellConfig(int(hunter.Level))
    if config.RequiredLevel <= int(hunter.Level) {
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
