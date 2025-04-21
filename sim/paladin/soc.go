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
// - The judgement operates via a dummy spell, that likely figures out whether to apply
//   half damage if the target is not stunned or not (counter to the tooltip, the base damage only is
//   multiplied by 2 if the target is stunned). These dummy spells are implemented as targetting
//   magic defense type, and have no flags to prevent misses, meaning they roll on spell hit table
//   and can miss. If it succeeds, it calls the "actual" Judgement of Command spell.
// - The actual Judgement of Command has flags to not miss and to avoid block/parry/dodge, but
//   it targets the melee defense type and so crits for double damage.
//   The Seal of Command aura watches for the base Judgement spell, and casts the actual
//   Judgement of Command when it successfully is cast.

func (paladin *Paladin) registerSealOfCommand() {
	type judge struct {
		spellID   int32
		minDamage float64
		maxDamage float64
		scale     float64
	}

	type proc struct {
		spellID int32
	}

	ranks := []struct {
		level      int32
		spellID    int32
		manaCost   float64
		scaleLevel int32
		proc       proc
		judge      judge
	}{
		{level: 20, spellID: 20375, manaCost: 65, scaleLevel: 28, proc: proc{spellID: 20424}, judge: judge{spellID: 20467, minDamage: 93, maxDamage: 101, scale: 5.6}},
		{level: 30, spellID: 20915, manaCost: 110, scaleLevel: 38, proc: proc{spellID: 20944}, judge: judge{spellID: 20963, minDamage: 146, maxDamage: 160, scale: 6.1}},
		{level: 40, spellID: 20918, manaCost: 140, scaleLevel: 48, proc: proc{spellID: 20945}, judge: judge{spellID: 20964, minDamage: 204, maxDamage: 224, scale: 5.6}},
		{level: 50, spellID: 20919, manaCost: 180, scaleLevel: 58, proc: proc{spellID: 20946}, judge: judge{spellID: 20965, minDamage: 261, maxDamage: 287, scale: 6.1}},
		{level: 60, spellID: 20920, manaCost: 210, scaleLevel: 60, proc: proc{spellID: 20947}, judge: judge{spellID: 20966, minDamage: 339, maxDamage: 373, scale: 6.1}},
	}

	dpm := paladin.AutoAttacks.NewPPMManager(7, core.ProcMaskMelee)

	icd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 1,
	}

	for i, rank := range ranks {
		rank := rank
		if paladin.Level < rank.level {
			break
		}

		minDamage := rank.judge.minDamage + float64(min(paladin.Level, rank.scaleLevel)-rank.level)*rank.judge.scale
		maxDamage := rank.judge.maxDamage + float64(min(paladin.Level, rank.scaleLevel)-rank.level)*rank.judge.scale

		judgeSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.judge.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | SpellFlag_RV | core.SpellFlagBatchStartAttackMacro,

			ClassSpellMask: ClassSpellMask_PaladinJudgementOfCommand, // used in judgement.go

			DamageMultiplier: paladin.getWeaponSpecializationModifier(),
			ThreatMultiplier: 1,

			BonusCoefficient: 0.429,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

				baseDamage := sim.Roll(minDamage, maxDamage) * 0.5 // unless stunned

				// Seal of Command requires this spell to act as its intermediary dummy,
				// rolling on the spell hit table. If it succeeds, the actual Judgement of Command rolls on the
				// melee special attack crit/hit table, necessitating two discrete spells.
				// All other judgements are cast directly.
				// Used to decide between spell.OutcomeMeleeSpecialCritOnly and spell.OutcomeAlwaysMiss
				dummyJudgeLanded := paladin.judgement.CalcOutcome(sim, target, paladin.judgement.OutcomeMagicHit).Landed()

				outcomeApplier := core.Ternary(dummyJudgeLanded, spell.OutcomeMeleeSpecialCritOnly, spell.OutcomeAlwaysMiss)
				spell.CalcAndDealDamage(sim, target, baseDamage, outcomeApplier)

			},
		})

		procSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: rank.proc.spellID},
			SpellSchool:    core.SpellSchoolHoly,
			DefenseType:    core.DefenseTypeMelee,
			ProcMask:       core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeProc | core.ProcMaskMeleeDamageProc,
			Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNotAProc,
			ClassSpellMask: ClassSpellMask_PaladinSealOfCommand,

			DamageMultiplier: 0.7 * paladin.getWeaponSpecializationModifier(),
			ThreatMultiplier: 1,

			BonusCoefficient: 0.29,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt:     sim.CurrentTime + core.SpellBatchWindow,
					Priority: core.ActionPriorityLow,
					OnAction: func(s *core.Simulation) {
						spell.DealDamage(sim, result)
					},
				})
			},
		})

		aura := paladin.RegisterAura(core.Aura{
			Label:    "Seal of Command" + paladin.Label + strconv.Itoa(i+1),
			ActionID: core.ActionID{SpellID: rank.spellID},
			Duration: time.Second * 30,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
					if icd.IsReady(sim) && dpm.Proc(sim, spell.ProcMask, "seal of command") {
						icd.Use(sim)
						procSpell.Cast(sim, result.Target)
					}
				}
			},
		})

		paladin.aurasSoC = append(paladin.aurasSoC, aura)

		paladin.sealOfCommand = paladin.RegisterSpell(core.SpellConfig{
			ActionID:    aura.ActionID,
			SpellSchool: core.SpellSchoolHoly,
			Flags:       core.SpellFlagAPL | core.SpellFlagBatchStartAttackMacro,

			RequiredLevel: int(rank.level),
			Rank:          i + 1,

			ManaCost: core.ManaCostOptions{
				FlatCost:   rank.manaCost - paladin.getLibramSealCostReduction(),
				Multiplier: paladin.benediction(),
			},
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				paladin.applySeal(aura, judgeSpell, sim)
			},
		})

		paladin.spellsJoC = append(paladin.spellsJoC, judgeSpell)
	}
}
