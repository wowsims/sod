package hunter

import (
	"log"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) getMongooseBiteConfig(rank int) core.SpellConfig {
	spellId := [5]int32{0, 1495, 14269, 14270, 14271}[rank]
	baseDamage := [5]float64{0, 25, 45, 75, 115}[rank]
	manaCost := [5]float64{0, 30, 40, 50, 65}[rank]
	level := [5]int{0, 16, 30, 44, 58}[rank]



	raptorFuryDmgMult := 0.1
	spellConfig := core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 5,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5 && hunter.CobraSlayerAura.IsActive()
		},

		BonusCritRating:  float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hunter.CobraSlayerAura.Deactivate(sim)

			multiplier := 1.0
			if stacks := hunter.RaptorFuryAura.GetStacks(); stacks > 0 {
				multiplier *= 1 + raptorFuryDmgMult*float64(stacks)
			}
			
			damage := multiplier * (baseDamage + (hunter.GetStat(stats.AttackPower) * 0.4))
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	}

	return spellConfig
}

func (hunter *Hunter) registerCobraSlayerAura() {
	hunter.RegisterAura(core.Aura{
		Label:    "Cobra Slayer Trigger",
		Duration: core.NeverExpires,
		MaxStacks: 20,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeDodge) {
				aura.SetStacks(sim, 1);
				hunter.CobraSlayerAura.Activate(sim)
			} else if result.Outcome.Matches(core.OutcomeLanded) {
				if spell.ProcMask == core.ProcMaskMeleeMHAuto || spell.ProcMask == core.ProcMaskMeleeOHAuto {
					if sim.Proc((float64(aura.GetStacks()) * 0.05), "Cobra Slayer") {
						log.Println(">>> AUTO PROC <<<")
						aura.SetStacks(sim, 1);
						hunter.CobraSlayerAura.Activate(sim)
					} else {
						aura.AddStack(sim)
					}
				}
			}
		},
	})

	hunter.CobraSlayerAura = hunter.RegisterAura(core.Aura{
		Label:    "Cobra Slayer Aura",
		Duration: time.Second * 30,
	})
}

func (hunter *Hunter) registerMongooseBiteSpell() {
	if !hunter.HasRune(proto.HunterRune_RuneChestCobraSlayer) {
		return
	}
	
	maxRank := 4
	hunter.registerCobraSlayerAura()
	for i := 1; i <= maxRank; i++ {
		config := hunter.getMongooseBiteConfig(i)
		if config.RequiredLevel <= int(hunter.Level) {
			hunter.MongooseBite = hunter.GetOrRegisterSpell(config)
		}
	}
}
