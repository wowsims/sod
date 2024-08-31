package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerOverpowerSpell(cdTimer *core.Timer) {
	hasTasteForBloodRune := warrior.HasRune(proto.WarriorRune_RuneTasteForBlood)

	bonusDamage := map[int32]float64{
		25: 5,
		40: 15,
		50: 25,
		60: 35,
	}[warrior.Level]

	spellID := map[int32]int32{
		25: 7384,
		40: 7887,
		50: 11584,
		60: 11585,
	}[warrior.Level]

	warrior.RegisterAura(core.Aura{
		Label:    "Overpower Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidDodge() {
				warrior.OverpowerAura.Activate(sim)
			}
		},
	})

	warrior.OverpowerAura = warrior.RegisterAura(core.Aura{
		Label:    "Overpower Aura",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: time.Second * 5,
	})

	warrior.Overpower = warrior.RegisterSpell(BattleStance, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   5 - warrior.FocusedRageDiscount,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second * 5,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.OverpowerAura.IsActive() || (hasTasteForBloodRune && warrior.TasteForBloodAura.IsActive())
		},

		BonusCritRating: 25 * core.CritRatingPerCritChance * float64(warrior.Talents.ImprovedOverpower),

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 0.75,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warrior.OverpowerAura.Deactivate(sim)

			baseDamage := bonusDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
