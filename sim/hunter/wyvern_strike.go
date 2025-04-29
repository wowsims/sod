package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getWyvernStrikeConfig(rank int) core.SpellConfig {
	spellId := [4]int32{0, 458436, 458481, 458482}[rank]
	manaCost := [4]float64{0, 55, 75, 100}[rank]
	level := [4]int{0, 1, 50, 60}[rank]

	// The spell tooltips list 3/4/6 on the respective ranks, but Zirene confirmed it's actually 10%.
	bleedCoeff := 0.10

	spellConfig := core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterWyvernStrike,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagStrike,

		Rank:          rank,
		RequiredLevel: level,
		MaxRange:      core.MaxMeleeAttackRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
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

		DamageMultiplier: 1,
		BonusCoefficient: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "WyvernStrike - Bleed" + hunter.Label + strconv.Itoa(rank),
				Tag:   "WyvernStrike - Bleed",
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 1,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				tickDamage := bleedCoeff * hunter.WyvernStrike.MeleeAttackPower()
				dot.Snapshot(target, tickDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weaponDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) * 1.40
			result := spell.CalcAndDealDamage(sim, target, weaponDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	}

	return spellConfig
}

func (hunter *Hunter) registerWyvernStrikeSpell() {
	if !hunter.Talents.WyvernSting || !hunter.HasRune(proto.HunterRune_RuneBootsWyvernStrike) {
		return
	}

	rank := map[int32]int{
		1:  1,
		50: 2,
		60: 3,
	}[hunter.Level]

	config := hunter.getWyvernStrikeConfig(rank)
	hunter.WyvernStrike = hunter.GetOrRegisterSpell(config)
}
