package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const socRanks = 5
const socProcSpellCoeff = 0.29
const socJudgeSpellCoeff = 0.423

var socLevels = [socRanks + 1]int{0, 20, 30, 40, 60, 50}
var socManaCosts = [socRanks + 1]float64{0, 65, 110, 140, 180, 210}
var socAuraSpellIDs = [socRanks + 1]int32{0, 20375, 20915, 20918, 20919, 20920}
var socProcSpellIDs = [socRanks + 1]int32{0, 20424, 20944, 20945, 20946, 20947}
var socJudgeSpellIDs = [socRanks + 1]int32{0, 20467, 20963, 20964, 20965, 20966}
var socJudgeBasePoints = [socRanks + 1]float64{0, 92, 145, 203, 260, 338}
var socJudgeRealPointsPerLevel = [socRanks + 1]float64{0, 5.6, 6.1, 5.6, 6.1, 6.1}
var socEffectDieSides = [socRanks + 1]float64{0, 9, 15, 21, 27, 35}
var socLevelMinMaxEffects = [socRanks + 1][]int32{{0}, {20, 28}, {30, 38}, {40, 48}, {50, 58}, {60, 60}}

func (paladin *Paladin) applySealOfCommandSpellAndAuraBaseConfig(rank int) {
	spellIDProc := socProcSpellIDs[rank]
	spellIDAura := socAuraSpellIDs[rank]
	spellIDJudge := socJudgeSpellIDs[rank]
	manaCost := socManaCosts[rank]
	level := socLevels[rank]
	scalingLevelMin := socLevelMinMaxEffects[rank][0]
	scalingLevelMax := socLevelMinMaxEffects[rank][1]
	judgeBasePoints := socJudgeBasePoints[rank]
	judgePointsPerLevel := socJudgeRealPointsPerLevel[rank]
	judgeDieSides := socEffectDieSides[rank]

	levelsToScale := min(paladin.Level, scalingLevelMax) - scalingLevelMin
	judgeMinDamage := judgeBasePoints + float64(levelsToScale)*judgePointsPerLevel
	judgeMaxDamage := judgeMinDamage + judgeDieSides

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellIDJudge}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.SpellCritMultiplier(),
		BonusCritRating:  paladin.getBonusCritChanceFromHolyPower(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(judgeMinDamage, judgeMaxDamage) + socJudgeSpellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: spellIDProc},
		SpellSchool:      core.SpellSchoolHoly,
		ProcMask:         core.ProcMaskEmpty, // This needs figured out properly
		Flags:            core.SpellFlagMeleeMetrics,
		RequiredLevel:    level,
		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   paladin.MeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (-1 + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())) * 0.7
			fullDamage := baseDamage + socProcSpellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, fullDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	ppmm := paladin.AutoAttacks.NewPPMManager(7.0, core.ProcMaskMelee)
	icd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 1,
	}

	auraActionID := core.ActionID{SpellID: spellIDAura}
	paladin.SealOfCommandAura[rank] = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Command" + paladin.Label + strconv.Itoa(rank),
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			// If a white hit, handle seal proc.
			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				if !icd.IsReady(sim) {
					return
				}
				if !ppmm.Proc(sim, spell.ProcMask, "seal of command") {
					return
				}
				// If we get here, SoC has procced, cast it.
				icd.Use(sim)
				onSwingProc.Cast(sim, result.Target)
			}
			// Else handle Judgements.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, result.Target)
			}

		},
	})

	aura := paladin.SealOfCommandAura[rank]
	paladin.SealOfCommand[rank] = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1.0 - (float64(paladin.Talents.Benediction) * 0.03),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
			paladin.CurrentSealExpiration = sim.CurrentTime + SealDuration
		},
	})
}

func (paladin *Paladin) registerSealOfCommandSpellAndAura() {
	paladin.SealOfCommand = make([]*core.Spell, sorRanks+1)
	paladin.SealOfCommandAura = make([]*core.Aura, sorRanks+1)

	for rank := 1; rank <= socRanks; rank++ {
		if int(paladin.Level) >= socLevels[rank] {
			paladin.MaxRankCommand = rank
			paladin.applySealOfCommandSpellAndAuraBaseConfig(rank)
		}
	}
}
