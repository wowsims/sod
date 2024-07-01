package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getWyvernStrikeConfig(rank int) core.SpellConfig {
	spellId := [4]int32{0, 458436, 458481, 458482}[rank]
	baseDamage := [4]float64{0, 3, 4, 6}[rank]
	manaCost := [4]float64{0, 55, 75, 100}[rank]
	level := [4]int{0, 1, 50, 60}[rank]

	hasCobraStrikes := hunter.pet != nil && hunter.HasRune(proto.HunterRune_RuneChestCobraStrikes)

	spellConfig := core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost * (1 - 0.02*float64(hunter.Talents.Efficiency)),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5
		},

		DamageMultiplier: 1,
		BonusCoefficient: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "WyvernStrike - Bleed" + hunter.Label + strconv.Itoa(rank),
				Tag:   "WyvernStrike - Bleed",
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 1,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				tickDamage := (baseDamage / 100 * 8 * hunter.WyvernStrike.MeleeAttackPower()) / float64(dot.NumberOfTicks)
				dot.Snapshot(target, tickDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weaponDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, weaponDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			if hasCobraStrikes && result.DidCrit() {
				hunter.CobraStrikesAura.Activate(sim)
				hunter.CobraStrikesAura.SetStacks(sim, 2)
			}
		},
	}

	return spellConfig
}

func (hunter *Hunter) registerWyvernStrikeSpell() {
	if !hunter.Talents.WyvernSting || !hunter.HasRune(proto.HunterRune_RuneBootsWyvernStrike) {
		return
	}

	maxRank := 3
	for i := 1; i <= maxRank; i++ {
		config := hunter.getWyvernStrikeConfig(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.WyvernStrike = hunter.GetOrRegisterSpell(config)
		}
	}
}
