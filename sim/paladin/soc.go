package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

// Seal of Command is a spell consisting of:
// - A judgement that has a flat damage roll, and scales with spellpower.
// - A 7ppm on-hit proc with a 1s ICD that deals 70% weapon damage and scales with spellpower.

// Judgement of Command has some unusual behaviour in classic:
// - The judgement operates via a dummy spell, that likely figures out wether or not to apply
//   half damage if the target is not stunned (counter to the tooltip, the base damage only is
//   multiplied by 2 if the target is stunned). These dummy spells are implemented as targetting
//   magic defense type, and have no flags to prevent misses, meaning they roll on spell hit table
//   and can miss. If it succeeds, it calls the "actual" Judgement of Command spell.
// - The actual Judgement of Command has flags to not miss and to avoid block/parry/dodge, but
//   it targets the melee defense type and so crits for double damage.
// - This is accomplished via the use of the SpellFlagPrimaryJudgement spell flag that is used exclusively by
//   the base judgement spell. The Seal of Command aura watches for this spell, and casts the actual
//   Judgement of Command when it successfully is cast.

const socRanks = 5

// Below is the base sp coefficient before it gets reduced by the 70% modifier
// weapon damage % effect to 20% actual.
const socProcSpellCoeff = 0.29
const socJudgeSpellCoeff = 0.429

var socLevels = [socRanks + 1]int{0, 20, 30, 40, 50, 60}
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
	judgeMinDamage := judgeBasePoints + 1 + float64(levelsToScale)*judgePointsPerLevel // 1..judgeDieSides
	judgeMaxDamage := judgeMinDamage + judgeDieSides

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellIDJudge},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,
		SpellCode:   SpellCode_PaladinJudgementOfCommand,

		BonusCritRating: paladin.holyPowerCritChance() + paladin.fanaticismCritChance(),

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1,
		BonusCoefficient: socJudgeSpellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(judgeMinDamage, judgeMaxDamage)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellIDProc},
		SpellSchool:   core.SpellSchoolHoly,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics,
		RequiredLevel: level,

		DamageMultiplier: 0.7 * paladin.getWeaponSpecializationModifier(),
		ThreatMultiplier: 1.0,
		BonusCoefficient: socProcSpellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	ppmm := paladin.AutoAttacks.NewPPMManager(7.0, core.ProcMaskMelee)
	icd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 1,
	}

	auraActionID := core.ActionID{SpellID: spellIDAura}
	aura := paladin.RegisterAura(core.Aura{
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
	paladin.SealOfCommandAura[rank] = aura

	if paladin.Ranged().ID == LibramOfBenediction {
		manaCost -= 10
	}

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
			paladin.ApplySeal(aura, onJudgementProc, sim)
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
