package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) getMongooseBiteConfig(rank int) core.SpellConfig {
	spellId := [5]int32{0, 1495, 14269, 14270, 14271}[rank]
	baseDamage := [5]float64{0, 25, 45, 75, 115}[rank]
	manaCost := [5]float64{0, 30, 40, 50, 65}[rank]
	level := [5]int{0, 16, 30, 44, 58}[rank]

	hasCobraSlayer := hunter.HasRune(proto.HunterRune_RuneHandsCobraSlayer)
	hasRaptorFury := hunter.HasRune(proto.HunterRune_RuneBracersRaptorFury)
	hasMeleeSpecialist := hunter.HasRune(proto.HunterRune_RuneBeltMeleeSpecialist)

	spellConfig := core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterMongooseBite,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

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
				Duration: time.Second * 5,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DefensiveState.IsActive()
		},

		BonusCritRating:  float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,
		CritDamageBonus:  hunter.mortalShots(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hunter.DefensiveState.Deactivate(sim)

			if hasMeleeSpecialist && sim.Proc(0.3, "Raptor Strike Reset") {
				hunter.RaptorStrike.CD.Reset()
				spell.CD.Reset()
			}

			multiplier := 1.0
			if hasRaptorFury {
				multiplier *= hunter.raptorFuryDamageMultiplier()
			}

			damage := baseDamage
			if hasCobraSlayer {
				damage += spell.MeleeAttackPower() * 0.45
			}
			damage *= multiplier

			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	}

	return spellConfig
}

func (hunter *Hunter) registerMongooseBiteSpell() {
	// Aura is only used as a pre-requisite for Mongoose Bite
	hunter.DefensiveState = hunter.RegisterAura(core.Aura{
		Label:    "Defensive State",
		ActionID: core.ActionID{SpellID: 5302},
		Duration: time.Second * 5,

		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidDodge() {
				aura.Activate(sim)
			}
		},
	})

	rank := map[int32]int{
		25: 1,
		40: 2,
		50: 3,
		60: 4,
	}[hunter.Level]

	config := hunter.getMongooseBiteConfig(rank)
	hunter.MongooseBite = hunter.GetOrRegisterSpell(config)
}
