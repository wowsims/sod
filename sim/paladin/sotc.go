package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (paladin *Paladin) registerSealOfTheCrusader() {
	type judge struct {
		spellID int32
		bonus   float64
	}

	var ranks = []struct {
		level      int32
		spellID    int32
		manaCost   float64
		scaleLevel int32
		ap         float64
		scale      float64
		judge      judge
	}{
		{level: 6, spellID: 21082, manaCost: 25, scaleLevel: 12, ap: 31, scale: 0.7, judge: judge{spellID: 21183, bonus: 20}},
		{level: 12, spellID: 20162, manaCost: 40, scaleLevel: 20, ap: 51, scale: 1.1, judge: judge{spellID: 20188, bonus: 30}},
		{level: 22, spellID: 20305, manaCost: 65, scaleLevel: 30, ap: 94, scale: 1.7, judge: judge{spellID: 20300, bonus: 50}},
		{level: 32, spellID: 20306, manaCost: 90, scaleLevel: 40, ap: 145, scale: 2, judge: judge{spellID: 20301, bonus: 80}},
		{level: 42, spellID: 20307, manaCost: 125, scaleLevel: 50, ap: 221, scale: 2.2, judge: judge{spellID: 20302, bonus: 110}},
		{level: 52, spellID: 20308, manaCost: 160, scaleLevel: 60, ap: 306, scale: 2.4, judge: judge{spellID: 20303, bonus: 140}},
	}

	improvedSotC := []float64{1, 1.05, 1.1, 1.15}[paladin.Talents.ImprovedSealOfTheCrusader]

	var libramAp, libramBonus float64
	if paladin.Ranged().ID == LibramOfFervor {
		libramAp = 48
		libramBonus = 33
	}

	for i, rank := range ranks {
		rank := rank
		if paladin.Level < rank.level {
			break
		}

		debuffs := paladin.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
			return core.JudgementOfTheCrusaderAura(&paladin.Unit, target, level, improvedSotC, libramBonus)
		})

		judgeSpell := paladin.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: rank.judge.spellID},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagBatchStartAttackMacro,

			ClassSpellMask: ClassSpellMask_PaladinJudgementOfTheCrusader,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
				debuffs.Get(target).Activate(sim)
			},
		})

		ap := rank.ap + rank.scale*float64(min(paladin.Level, rank.scaleLevel)-rank.level) + libramAp

		aura := paladin.RegisterAura(core.Aura{
			Label:    "Seal of the Crusader" + paladin.Label + strconv.Itoa(i+1),
			ActionID: core.ActionID{SpellID: rank.spellID},
			Duration: time.Second * 30,
			OnGain: func(_ *core.Aura, sim *core.Simulation) {
				paladin.MultiplyMeleeSpeed(sim, 1.4)
				paladin.AutoAttacks.MHAuto().ApplyMultiplicativeBaseDamageBonus(1 / 1.4)
				paladin.AddStatDynamic(sim, stats.AttackPower, ap*improvedSotC)
			},
			OnExpire: func(_ *core.Aura, sim *core.Simulation) {
				paladin.MultiplyMeleeSpeed(sim, 1/1.4)
				paladin.AutoAttacks.MHAuto().ApplyMultiplicativeBaseDamageBonus(1.4)
				paladin.AddStatDynamic(sim, stats.AttackPower, -ap*improvedSotC)
			},
		})

		paladin.aurasSotC = append(paladin.aurasSotC, aura)

		paladin.RegisterSpell(core.SpellConfig{
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

		paladin.spellsJotC = append(paladin.spellsJotC, judgeSpell)
	}
}
