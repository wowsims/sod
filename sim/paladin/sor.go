package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerSealOfRighteousness() {
	type proc struct {
		spellID int32
		value   float64
		scale   float64
		coeff   float64
	}

	type judge struct {
		spellID   int32
		minDamage float64
		maxDamage float64
		scale     float64
		coeff     float64
	}

	// proc coefficients are likely 0.058/0.125/0.185 for ranks 1/2/3, and reach 0.2 for ranks >= 4.
	//  there's a weird SotC interaction, in that both "damage done" and "damage taken" are affected.
	var ranks = []struct {
		level      int32
		spellID    int32
		manaCost   float64
		scaleLevel int32
		proc       proc
		judge      judge
	}{
		{level: 1, spellID: 20154, manaCost: 20, scaleLevel: 7, proc: proc{spellID: 25742, value: 108, scale: 18, coeff: 0.029}, judge: judge{spellID: 20187, minDamage: 15, maxDamage: 15, scale: 1.8, coeff: 0.144}},
		{level: 10, spellID: 20287, manaCost: 40, scaleLevel: 16, proc: proc{spellID: 25740, value: 216, scale: 17, coeff: 0.063}, judge: judge{spellID: 20280, minDamage: 25, maxDamage: 27, scale: 1.9, coeff: 0.312}},
		{level: 18, spellID: 20288, manaCost: 60, scaleLevel: 24, proc: proc{spellID: 25739, value: 352, scale: 23, coeff: 0.184}, judge: judge{spellID: 20281, minDamage: 39, maxDamage: 43, scale: 2.4, coeff: 0.462}},
		{level: 26, spellID: 20289, manaCost: 90, scaleLevel: 32, proc: proc{spellID: 25738, value: 541, scale: 31, coeff: 0.184}, judge: judge{spellID: 20282, minDamage: 57, maxDamage: 63, scale: 2.8, coeff: 0.5}},
		{level: 34, spellID: 20290, manaCost: 120, scaleLevel: 40, proc: proc{spellID: 25737, value: 785, scale: 37, coeff: 0.184}, judge: judge{spellID: 20283, minDamage: 78, maxDamage: 86, scale: 3.1, coeff: 0.5}},
		{level: 42, spellID: 20291, manaCost: 140, scaleLevel: 48, proc: proc{spellID: 25736, value: 1082, scale: 41, coeff: 0.184}, judge: judge{spellID: 20284, minDamage: 102, maxDamage: 112, scale: 3.8, coeff: 0.5}},
		{level: 50, spellID: 20292, manaCost: 170, scaleLevel: 56, proc: proc{spellID: 25735, value: 1407, scale: 47, coeff: 0.184}, judge: judge{spellID: 20285, minDamage: 131, maxDamage: 143, scale: 4.1, coeff: 0.5}},
		{level: 58, spellID: 20293, manaCost: 200, scaleLevel: 60, proc: proc{spellID: 25713, value: 1786, scale: 47, coeff: 0.184}, judge: judge{spellID: 20286, minDamage: 162, maxDamage: 178, scale: 4.1, coeff: 0.5}},
	}

	improvedSoR := paladin.improvedSoR()

	for i, rank := range ranks {
		i, rank := i, rank
		if paladin.Level < rank.level {
			break
		}

		/*
		 * Seal of Righteousness is a Spell/Aura that when active makes the paladin capable of procing
		 * two different SpellIDs depending on a paladin's casted spell or melee swing.
		 *
		 * (Judgement of Righteousness):
		 *   - Deals flat damage that is affected by Improved SoR talent, and
		 *     has a spellpower scaling that is unaffected by that talent.
		 *   - Targets magic defense and rolls to hit and crit.
		 *
		 * (Seal of Righteousness):
		 *   - Procs from white hits.
		 *   - Cannot miss or be dodged/parried/blocked if the underlying white hit lands.
		 *   - Deals damage that is a function of weapon speed, and spellpower.
		 *   - Has 0.85 scale factor on base damage if using 1h, 1.2 if using 2h.
		 *   - CANNOT CRIT.
		 */

		minDamage := rank.judge.minDamage + rank.judge.scale*float64(min(paladin.Level, rank.scaleLevel)-rank.level)
		maxDamage := rank.judge.maxDamage + rank.judge.scale*float64(min(paladin.Level, rank.scaleLevel)-rank.level)

		judgeSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.judge.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagMeleeMetrics,

			BonusCritRating: paladin.holyCrit(), // TODO check whether fanaticism applies

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: rank.judge.coeff,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(minDamage, maxDamage) * improvedSoR
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			},
		})

		value := rank.proc.value + rank.proc.scale*float64(min(paladin.Level, rank.scaleLevel)-rank.level)

		damage := value / 100 * core.TernaryFloat64(paladin.Has2hEquipped(), 1.2, 0.85) * paladin.MainHand().SwingSpeed

		procSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.proc.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagMeleeMetrics,

			DamageMultiplier: 1, // TODO check whether this benefits from weapon specializations
			ThreatMultiplier: 1,

			// Testing seems to show 2h benefits from spellpower about 12% more than 1h weapons.
			BonusCoefficient: rank.proc.coeff * core.TernaryFloat64(paladin.Has2hEquipped(), 1.12, 1.0),

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := damage * improvedSoR
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
			},
		})

		aura := paladin.RegisterAura(core.Aura{
			Label:    "Seal of Righteousness" + paladin.Label + strconv.Itoa(i+1),
			ActionID: core.ActionID{SpellID: rank.spellID},
			Duration: time.Second * 30,

			OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}
				if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		paladin.sealOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
			ActionID:    aura.ActionID,
			SpellSchool: core.SpellSchoolHoly,
			Flags:       core.SpellFlagAPL,

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
				paladin.ApplySeal(aura, judgeSpell, sim)
			},
		})
	}
}
