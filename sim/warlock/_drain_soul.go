package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) registerDrainSoulSpell() {
	hasSoulSiphonRune := warlock.HasRune(proto.WarlockRune_RuneChestSoulSiphon)

	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)

	calcSoulSiphonMult := func(target *core.Unit) float64 {
		auras := []*core.Aura{
			warlock.UnstableAffliction.Dot(target).Aura,
			warlock.Corruption.Dot(target).Aura,
			warlock.Seed.Dot(target).Aura,
			warlock.CurseOfAgony.Dot(target).Aura,
			warlock.CurseOfDoom.Dot(target).Aura,
			warlock.CurseOfElementsAuras.Get(target),
			warlock.CurseOfWeaknessAuras.Get(target),
			warlock.CurseOfTonguesAuras.Get(target),
			warlock.ShadowEmbraceDebuffAura(target),
			// missing: death coil
		}
		if warlock.HauntDebuffAuras != nil {
			auras = append(auras, warlock.HauntDebuffAuras.Get(target))
		}
		numActive := 0
		for _, aura := range auras {
			if aura.IsActive() {
				numActive++
			}
		}
		return 1.0 + float64(min(3, numActive))*soulSiphonMultiplier
	}

	warlock.DrainSoul = warlock.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_WarlockDrainSoul,
		ActionID:    core.ActionID{SpellID: 47855},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled | core.SpellFlagHauntSE | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.14,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				// ChannelTime: channelTime,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Drain Soul",
			},
			NumberOfTicks:       5,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := 142 + 0.429*dot.Spell.SpellDamage()
				dot.Snapshot(target, baseDmg*calcSoulSiphonMult(target), isRollover)

				if hasSoulSiphonRune {
					isExecute := sim.IsExecutePhase20()
					perDoTMultiplier := core.TernaryFloat64(isExecute, SoulSiphonDoTMultiplier, SoulSiphonDoTMultiplierExecute)
					maxMultiplier := 1 + core.TernaryFloat64(isExecute, SoulSiphonDoTMultiplierMax, SoulSiphonDoTMultiplierMaxExecute)
					multiplier := 1.0

					hasAura := func(target *core.Unit, label string, rank int) bool {
						for i := 1; i <= rank; i++ {
							if target.HasActiveAura(label + strconv.Itoa(rank)) {
								return true
							}
						}
						return false
					}
					if hasAura(target, "Corruption-"+warlock.Label, 7) {
						multiplier += perDoTMultiplier
					}
					if hasAura(target, "CurseofAgony-"+warlock.Label, 6) {
						multiplier += perDoTMultiplier
					}
					if hasAura(target, "SiphonLife-"+warlock.Label, 3) {
						multiplier += perDoTMultiplier
					}
					if target.HasActiveAura("UnstableAffliction-" + warlock.Label) {
						multiplier += perDoTMultiplier
					}
					if target.HasActiveAura("Haunt-" + warlock.Label) {
						multiplier += perDoTMultiplier
					}

					dot.SnapshotAttackerMultiplier *= max(multiplier, maxMultiplier)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.UpdateExpires(dot.ExpiresAt())
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDmg := (142 + 0.429*spell.SpellDamage()) * calcSoulSiphonMult(target)
				return spell.CalcPeriodicDamage(sim, target, baseDmg, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				mult := (4.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)) / (1 + 0.04*float64(warlock.Talents.DeathsEmbrace))
				warlock.DrainSoul.DamageMultiplier = mult
			}
		})
	})
}
