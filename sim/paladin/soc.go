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

		DamageMultiplier: paladin.SpellCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(judgeMinDamage, judgeMaxDamage) + socJudgeSpellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHit)
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
		},
	})

	/*
	 * Seal of Command is an Spell/Aura that when active makes the paladin capable of procing
	 * 2 different SpellIDs depending on a paladin's casted spell or melee swing.
	 *
	 * SpellID 20467 (Judgement of Command):
	 *   - Procs off of any "Primary" Judgement (JoL, JoW, JoJ).
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage.
	 *   - Crits off of a melee modifier.
	 *
	 * SpellID 20424 (Seal of Command):
	 *   - Procs off of any melee special ability, or white hit.
	 *   - If the ability is SINGLE TARGET, it hits up to 2 extra targets.
	 *   - Deals hybrid AP/SP damage * current weapon speed.
	 *   - Crits off of a melee modifier.
	 *   - CAN MISS, BE DODGED/PARRIED/BLOCKED.
	 */

	// numHits := min(3, paladin.Env.GetNumTargets()) // primary target + 2 others
	// results := make([]*core.SpellResult, numHits)

	// onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:    core.ActionID{SpellID: 20467}, // Judgement of Command
	// 	SpellSchool: core.SpellSchoolHoly,
	// 	ProcMask:    core.ProcMaskMeleeSpecial,
	// 	Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

	// 	// BonusCritRating: (6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance) +
	// 	// 	(core.TernaryFloat64(paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 4), 5, 0) * core.CritRatingPerCritChance),

	// 	DamageMultiplier: 1 *
	// 		// (1 + paladin.getItemSetLightswornBattlegearBonus4() +
	// 		// 	paladin.getTalentTheArtOfWarBonus()) *
	// 		(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()),
	// 	CritMultiplier:   paladin.MeleeCritMultiplier(),
	// 	ThreatMultiplier: 1,

	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		mhWeaponDamage := 0 +
	// 			spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
	// 			spell.BonusWeaponDamage()
	// 		baseDamage := 0.19*mhWeaponDamage +
	// 			0.08*spell.MeleeAttackPower() +
	// 			0.13*spell.SpellPower()

	// 		// Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
	// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
	// 	},
	// })

	// onSpecialOrSwingActionID := core.ActionID{SpellID: 20424}
	// onSpecialOrSwingProcCleave := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:    onSpecialOrSwingActionID,
	// 	SpellSchool: core.SpellSchoolHoly,
	// 	ProcMask:    core.ProcMaskEmpty,
	// 	Flags:       core.SpellFlagMeleeMetrics,

	// 	DamageMultiplier: 1 *
	// 		// (1 + paladin.getItemSetLightswornBattlegearBonus4()) *
	// 		(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()) *
	// 		0.36, // Only 36% of weapon damage.
	// 	CritMultiplier:   paladin.MeleeCritMultiplier(),
	// 	ThreatMultiplier: 1,

	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		curTarget := target
	// 		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
	// 			baseDamage := 0 +
	// 				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
	// 				spell.BonusWeaponDamage()

	// 			results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
	// 			curTarget = sim.Environment.NextTargetUnit(curTarget)
	// 		}

	// 		curTarget = target
	// 		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
	// 			spell.DealDamage(sim, results[hitIndex])
	// 			curTarget = sim.Environment.NextTargetUnit(curTarget)
	// 		}
	// 	},
	// })

	// onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:    onSpecialOrSwingActionID,
	// 	SpellSchool: core.SpellSchoolHoly,
	// 	ProcMask:    core.ProcMaskEmpty, // unlike SoV, SoC crits don't proc Vengeance
	// 	Flags:       core.SpellFlagMeleeMetrics,

	// 	DamageMultiplier: 1 *
	// 		// (1 + paladin.getItemSetLightswornBattlegearBonus4()) *
	// 		(1 + paladin.getTalentTwoHandedWeaponSpecializationBonus()) *
	// 		0.36, // Only 36% of weapon damage.
	// 	CritMultiplier:   paladin.MeleeCritMultiplier(),
	// 	ThreatMultiplier: 1,

	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := 0 +
	// 			spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
	// 			spell.BonusWeaponDamage()

	// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
	// 	},
	// })

	// // Seal of Command aura.
	// auraActionID := core.ActionID{SpellID: 20375}
	// paladin.SealOfCommandAura = paladin.RegisterAura(core.Aura{
	// 	Label:    "Seal of Command",
	// 	Tag:      "Seal",
	// 	ActionID: auraActionID,
	// 	Duration: SealDuration,

	// 	OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		// Don't proc on misses or our own procs.
	// 		if !result.Landed() || spell == onJudgementProc || spell.SameAction(onSpecialOrSwingActionID) {
	// 			return
	// 		}

	// 		// Differ between judgements and other melee abilities.
	// 		if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
	// 			onJudgementProc.Cast(sim, result.Target)
	// 			// if paladin.Talents.JudgementsOfTheJust > 0 {
	// 			// 	// Special JoJ talent behavior, procs swing seal on judgements
	// 			// 	// For SoC this is a cleave.
	// 			// 	onSpecialOrSwingProcCleave.Cast(sim, result.Target)
	// 			// }
	// 		} else if spell.IsMelee() {
	// 			// Temporary check to avoid AOE double procing.
	// 			// if spell.SpellID == paladin.HammerOfTheRighteous.SpellID || spell.SpellID == paladin.DivineStorm.SpellID {
	// 			// 	onSpecialOrSwingProc.Cast(sim, result.Target)
	// 			// } else {
	// 			// 	onSpecialOrSwingProcCleave.Cast(sim, result.Target)
	// 			// }
	// 		}
	// 	},
	// })

	// aura := paladin.SealOfCommandAura
	// paladin.SealOfCommand = paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:    auraActionID, // Seal of Command self buff.
	// 	SpellSchool: core.SpellSchoolHoly,
	// 	ProcMask:    core.ProcMaskEmpty,
	// 	Flags:       core.SpellFlagAPL,

	// 	ManaCost: core.ManaCostOptions{
	// 		BaseCost:   0.14,
	// 		Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: core.GCDDefault,
	// 		},
	// 	},

	// 	ThreatMultiplier: 1,

	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		if paladin.CurrentSeal != nil {
	// 			paladin.CurrentSeal.Deactivate(sim)
	// 		}
	// 		paladin.CurrentSeal = aura
	// 		paladin.CurrentSeal.Activate(sim)
	// 	},
	// })
}

func (paladin *Paladin) registerSealOfCommandSpellAndAura() {
	paladin.SealOfCommand = make([]*core.Spell, sorRanks+1)
	paladin.SealOfCommandAura = make([]*core.Aura, sorRanks+1)

	for rank := 1; rank <= socRanks; rank++ {
		if int(paladin.Level) >= socLevels[rank] {
			paladin.applySealOfCommandSpellAndAuraBaseConfig(rank)
		}
	}
}
