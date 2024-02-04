package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getCurseOfAgonyBaseConfig(rank int) core.SpellConfig {
	spellId := [7]int32{0, 980, 1014, 6217, 11711, 11712, 11713}[rank]
	spellCoeff := [7]float64{0, .046, .077, .083, .083, .083, .083}[rank]
	baseDamage := [7]float64{0, 7, 15, 27, 42, 65, 87}[rank]
	manaCost := [7]float64{0, 25, 50, 90, 130, 170, 215}[rank]
	level := [7]int{0, 8, 18, 28, 38, 48, 58}[rank]
	hasInvocationRune := warlock.HasRune(proto.WarlockRune_RuneBeltInvocation)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		Flags:         core.SpellFlagAPL | core.SpellFlagHauntSE | core.SpellFlagResetAttackSwing | core.SpellFlagPureDot,
		ProcMask:      core.ProcMaskSpellDamage,
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

		BonusHitRating: 2 * float64(warlock.Talents.Suppression) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + 0.02*float64(warlock.Talents.ImprovedCurseOfWeakness)) *
			(1 + 0.02*float64(warlock.Talents.ShadowMastery)),
		ThreatMultiplier: 1,
		FlatThreatBonus:  0,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "CurseofAgony-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks: 12,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 0.5 * (baseDamage + spellCoeff*dot.Spell.SpellPower())
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])

				if warlock.AmplifyCurseAura.IsActive() {
					dot.SnapshotAttackerMultiplier *= 1.5
					warlock.AmplifyCurseAura.Deactivate(sim)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				if dot.TickCount%4 == 0 { // CoA ramp up
					dot.SnapshotBaseDamage += 0.5 * (baseDamage + spellCoeff*dot.Spell.SpellPower())
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--

				if hasInvocationRune && spell.Dot(target).IsActive() {
					warlock.InvocationRefresh(sim, spell.Dot(target))
				}

				//warlock.CurseOfDoom.Dot(target).Cancel(sim)
				spell.Dot(target).Apply(sim)
			}
		},
	}
}

func (warlock *Warlock) registerCurseOfAgonySpell() {
	maxRank := 6

	for i := 1; i <= maxRank; i++ {
		config := warlock.getCurseOfAgonyBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.CurseOfAgony = warlock.GetOrRegisterSpell(config)
		}
	}
}

func (warlock *Warlock) registerCurseOfRecklessnessSpell() {
	playerLevel := warlock.Level

	warlock.CurseOfRecklessnessAuras = warlock.NewEnemyAuraArray(core.CurseOfRecklessnessAura)

	spellID := map[int32]int32{
		25: 704,
		40: 7658,
		50: 7659,
		60: 11717,
	}[playerLevel]

	manaCost := map[int32]float64{
		25: 35.0,
		40: 60.0,
		50: 90.0,
		60: 115.0,
	}[playerLevel]

	warlock.CurseOfRecklessness = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating:   float64(warlock.Talents.Suppression) * 2 * core.CritRatingPerCritChance,
		ThreatMultiplier: 1,
		FlatThreatBonus:  156,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfRecklessnessAuras.Get(target).Activate(sim)
			}
		},

		RelatedAuras: []core.AuraArray{warlock.CurseOfRecklessnessAuras},
	})
}

func (warlock *Warlock) registerCurseOfElementsSpell() {
	playerLevel := warlock.Level
	if playerLevel < 40 {
		return
	}

	warlock.CurseOfElementsAuras = warlock.NewEnemyAuraArray(core.CurseOfElementsAura)

	spellID := map[int32]int32{
		40: 1490,
		50: 11721,
		60: 11722,
	}[playerLevel]

	manaCost := map[int32]float64{
		40: 100.0,
		50: 150.0,
		60: 200.0,
	}[playerLevel]

	warlock.CurseOfElements = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating:   float64(warlock.Talents.Suppression) * 2 * core.CritRatingPerCritChance,
		ThreatMultiplier: 1,
		FlatThreatBonus:  156,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfElementsAuras.Get(target).Activate(sim)
			}
		},

		RelatedAuras: []core.AuraArray{warlock.CurseOfElementsAuras},
	})
}

func (warlock *Warlock) registerCurseOfShadowSpell() {
	playerLevel := warlock.Level
	if playerLevel < 50 {
		return
	}

	warlock.CurseOfShadowAuras = warlock.NewEnemyAuraArray(core.CurseOfShadowAura)

	spellID := map[int32]int32{
		50: 17862,
		60: 17937,
	}[playerLevel]

	manaCost := map[int32]float64{
		50: 150.0,
		60: 200.0,
	}[playerLevel]

	warlock.CurseOfShadow = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating:   float64(warlock.Talents.Suppression) * 2 * core.CritRatingPerCritChance,
		ThreatMultiplier: 1,
		FlatThreatBonus:  156,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfShadowAuras.Get(target).Activate(sim)
			}
		},

		RelatedAuras: []core.AuraArray{warlock.CurseOfShadowAuras},
	})
}

func (warlock *Warlock) registerAmplifyCurseSpell() {
	if !warlock.Talents.AmplifyCurse {
		return
	}

	actionID := core.ActionID{SpellID: 18288}

	warlock.AmplifyCurseAura = warlock.GetOrRegisterAura(core.Aura{
		Label:    "Amplify Curse",
		ActionID: actionID,
		Duration: time.Second * 30,
	})

	warlock.AmplifyCurse = warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 3 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.AmplifyCurseAura.Activate(sim)
		},
	})
}
