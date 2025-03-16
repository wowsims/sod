package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type Stance uint8

const (
	BattleStance Stance = 1 << iota
	DefensiveStance
	BerserkerStance
	GladiatorStance

	AnyStance = BattleStance | DefensiveStance | BerserkerStance | GladiatorStance
)

func (stance Stance) Matches(other Stance) bool {
	return (stance & other) != 0
}

var StanceCodes = ClassSpellMask_WarriorStanceBattle | ClassSpellMask_WarriorStanceDefensive | ClassSpellMask_WarriorStanceBerserker | ClassSpellMask_WarriorStanceGladiator

const stanceEffectCategory = "Stance"

func (warrior *Warrior) StanceMatches(other Stance) bool {
	return warrior.Stance.Matches(other)
}

func (warrior *Warrior) makeStanceSpell(stance Stance, aura *core.Aura, stanceCD *core.Timer) *WarriorSpell {
	SpellClassMask := map[Stance]uint64{
		BattleStance:    ClassSpellMask_WarriorStanceBattle,
		DefensiveStance: ClassSpellMask_WarriorStanceDefensive,
		BerserkerStance: ClassSpellMask_WarriorStanceBerserker,
		GladiatorStance: ClassSpellMask_WarriorStanceGladiator,
	}[stance]
	actionID := aura.ActionID
	maxRetainedRage := 5 * float64(warrior.Talents.TacticalMastery)
	rageMetrics := warrior.NewRageMetrics(actionID)

	stanceSpell := warrior.RegisterSpell(AnyStance, core.SpellConfig{
		ClassSpellMask: SpellClassMask,
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    stanceCD,
				Duration: time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !warrior.StanceMatches(stance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if warrior.CurrentRage() > maxRetainedRage {
				warrior.SpendRage(sim, warrior.CurrentRage()-maxRetainedRage, rageMetrics)
			}

			if warrior.WarriorInputs.StanceSnapshot {
				// Delayed, so same-GCD casts are affected by the current aura.
				//  Alternatively, those casts could just (artificially) happen before the stance change.
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt:     sim.CurrentTime + 10*time.Millisecond,
					OnAction: aura.Activate,
				})
			} else {
				aura.Activate(sim)
			}

			warrior.PreviousStance = warrior.Stance
			warrior.Stance = stance
		},
	})

	warrior.Stances = append(warrior.Stances, stanceSpell)

	return stanceSpell
}

func (warrior *Warrior) registerBattleStanceAura() {
	warrior.BattleStanceAura = warrior.RegisterAura(core.Aura{
		Label:    "Battle Stance",
		ActionID: core.ActionID{SpellID: 2457},
		Duration: core.NeverExpires,
	})
	warrior.BattleStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.ThreatMultiplier *= 0.8
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.ThreatMultiplier /= 0.8
		},
	})
}

func (warrior *Warrior) registerDefensiveStanceAura() {
	warrior.defensiveStanceThreatMultiplier = 1.3 * []float64{1, 1.03, 1.06, 1.09, 1.12, 1.15}[warrior.Talents.Defiance]

	warrior.DefensiveStanceAura = warrior.RegisterAura(core.Aura{
		Label:    "Defensive Stance",
		ActionID: core.ActionID{SpellID: 71},
		Duration: core.NeverExpires,
	})
	warrior.DefensiveStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.ThreatMultiplier *= warrior.defensiveStanceThreatMultiplier
			ee.Aura.Unit.PseudoStats.DamageDealtMultiplier *= 0.9
			ee.Aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.9
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.ThreatMultiplier /= warrior.defensiveStanceThreatMultiplier
			ee.Aura.Unit.PseudoStats.DamageDealtMultiplier /= 0.9
			ee.Aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.9
		},
	})
}

func (warrior *Warrior) registerBerserkerStanceAura() {
	warrior.BerserkerStanceAura = warrior.RegisterAura(core.Aura{
		Label:    "Berserker Stance",
		ActionID: core.ActionID{SpellID: 2458},
		Duration: core.NeverExpires,
	})
	warrior.BerserkerStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.ThreatMultiplier *= 0.8
			ee.Aura.Unit.PseudoStats.DamageTakenMultiplier *= 1.1
			ee.Aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, core.CritRatingPerCritChance*3)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.PseudoStats.ThreatMultiplier /= 0.8
			ee.Aura.Unit.PseudoStats.DamageTakenMultiplier /= 1.1
			ee.Aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, -core.CritRatingPerCritChance*3)
		},
	})
}

// An aggressive stance that increases damage while you are wearing a shield by 10% and increases block chance by 10%, but reduces armor by 30% and threat generated by 30%.
// In addition, you gain 50% increased Rage when your auto-attack damages an enemy not targeting you.
// While wearing a shield in Gladiator Stance, you may use all abilities that are restricted to other stances.
func (warrior *Warrior) registerGladiatorStanceAura() {
	if !warrior.HasRune(proto.WarriorRune_RuneGladiatorStance) {
		return
	}

	warrior.gladiatorStanceDamageMultiplier = 1.1
	isTanking := warrior.IsTanking()

	gladStanceDamageAura := warrior.RegisterAura(core.Aura{
		Label:    "Gladiator Stance Damage Bonus",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= warrior.gladiatorStanceDamageMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= warrior.gladiatorStanceDamageMultiplier
		},
	})

	// Use a periodic action aura to verify that the player is still using a shield for the damage bonus and stance override.
	// This is needed if Item Swapping.
	var gladStanceValidationPA *core.PendingAction
	var gladStanceStanceOverrideEE *core.ExclusiveEffect
	gladStanceValidationAura := warrior.RegisterAura(core.Aura{
		Label:    "Gladiator Stance Shield Validation",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			gladStanceValidationPA = nil
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			gladStanceValidationPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Second * 2,
				TickImmediately: true,
				OnAction: func(sim *core.Simulation) {
					if warrior.GladiatorStanceAura.IsActive() && warrior.PseudoStats.CanBlock {
						if !gladStanceDamageAura.IsActive() {
							gladStanceDamageAura.Activate(sim)
						}
						if !gladStanceStanceOverrideEE.IsActive() {
							gladStanceStanceOverrideEE.Activate(sim)
						}
					} else {
						gladStanceDamageAura.Deactivate(sim)
						gladStanceStanceOverrideEE.Deactivate(sim)
						return
					}
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if gladStanceValidationPA != nil {
				gladStanceValidationPA.Cancel(sim)
			}
			if gladStanceDamageAura.IsActive() {
				gladStanceDamageAura.Deactivate(sim)
			}
			if !gladStanceStanceOverrideEE.IsActive() {
				gladStanceStanceOverrideEE.Activate(sim)
			}
		},
	})

	warrior.GladiatorStanceAura = warrior.RegisterAura(core.Aura{
		Label:    "Gladiator Stance",
		ActionID: core.ActionID{SpellID: int32(proto.WarriorRune_RuneGladiatorStance)},
		Duration: core.NeverExpires,
	})
	warrior.GladiatorStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Block, 10*core.BlockRatingPerBlockChance)
			ee.Aura.Unit.PseudoStats.ArmorMultiplier *= 0.7
			ee.Aura.Unit.PseudoStats.ThreatMultiplier *= 0.7
			if !isTanking {
				warrior.AddDamageDealtRageMultiplier(1.5)
			}

			gladStanceValidationAura.Activate(sim)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Block, -10*core.BlockRatingPerBlockChance)
			ee.Aura.Unit.PseudoStats.ArmorMultiplier /= 0.7
			ee.Aura.Unit.PseudoStats.ThreatMultiplier /= 0.7
			if !isTanking {
				warrior.AddDamageDealtRageMultiplier(1 / 1.5)
			}

			gladStanceValidationAura.Deactivate(sim)
		},
	})
	gladStanceStanceOverrideEE = warrior.newStanceOverrideExclusiveEffect(AnyStance, warrior.GladiatorStanceAura)
}

func (warrior *Warrior) registerStances() {
	warrior.Stances = make([]*WarriorSpell, 0)
	stanceCD := warrior.NewTimer()
	warrior.registerBattleStanceAura()
	warrior.registerDefensiveStanceAura()
	warrior.registerBerserkerStanceAura()
	warrior.registerGladiatorStanceAura()
	warrior.BattleStance = warrior.makeStanceSpell(BattleStance, warrior.BattleStanceAura, stanceCD)
	warrior.DefensiveStance = warrior.makeStanceSpell(DefensiveStance, warrior.DefensiveStanceAura, stanceCD)
	warrior.BerserkerStance = warrior.makeStanceSpell(BerserkerStance, warrior.BerserkerStanceAura, stanceCD)
	if warrior.HasRune(proto.WarriorRune_RuneGladiatorStance) {
		warrior.GladiatorStance = warrior.makeStanceSpell(GladiatorStance, warrior.GladiatorStanceAura, stanceCD)
	}
}
