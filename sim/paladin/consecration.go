package paladin

import (
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const consecrationRanks = 5

var consecrationLevels = [consecrationRanks + 1]int{0, 20, 30, 40, 50, 60}
var consecrationSpellIDs = [consecrationRanks + 1]int32{0, 26573, 20116, 20922, 20923, 20924}
var consecrationBaseDamages = [consecrationRanks + 1]float64{0, 64 / 8, 120 / 8, 192 / 8, 280 / 8, 384 / 8}
var consecrationManaCosts = [consecrationRanks + 1]float64{0, 135, 235, 320, 435, 565}

func (paladin *Paladin) getConsecrationBaseConfig(rank int, cd core.Cooldown) core.SpellConfig {
	spellId := consecrationSpellIDs[rank]
	baseDamage := consecrationBaseDamages[rank]
	manaCost := consecrationManaCosts[rank]
	level := consecrationLevels[rank]

	hasWrath := paladin.HasRune(proto.PaladinRune_RuneHeadWrath)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolHoly,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagPureDot | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: cd,
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Consecration-" + paladin.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: false,
			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := paladin.CurrentTarget
				dot.SnapshotBaseDamage = baseDamage + 0.042*dot.Spell.SpellDamage()
				if hasWrath {
					dot.Spell.BonusCritRating += paladin.GetStat(stats.MeleeCrit)
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
					dot.Spell.BonusCritRating -= paladin.GetStat(stats.MeleeCrit)
				}
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				outcomeApplier := core.Ternary(hasWrath, dot.OutcomeMagicHitAndSnapshotCrit, dot.Spell.OutcomeMagicHit)
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, outcomeApplier)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	}
}

func (paladin *Paladin) registerConsecrationSpell() {
	if !paladin.Talents.Consecration {
		return
	}

	cd := core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 8,
	}
	paladin.Consecration = make([]*core.Spell, consecrationRanks+1)
	for rank := 1; rank <= consecrationRanks; rank++ {
		config := paladin.getConsecrationBaseConfig(rank, cd)
		if config.RequiredLevel <= int(paladin.Level) {
			paladin.Consecration[rank] = paladin.GetOrRegisterSpell(config)
		}
	}
}
