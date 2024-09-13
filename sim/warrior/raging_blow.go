package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// A ferocious strike that deals 100% weapon damage, but can only be used while Enrage, Berserker Rage, or Bloodrage is active.
// Raging blow cooldown is reduced by 1 second when you use another melee ability while enraged.
func (warrior *Warrior) registerRagingBlow() {
	if !warrior.HasRune(proto.WarriorRune_RuneRagingBlow) {
		return
	}

	warrior.RegisterAura(core.Aura{
		Label:    "Raging Blow CDR",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) && warrior.IsEnraged() && !warrior.RagingBlow.CD.IsReady(sim) && spell.SpellCode != SpellCode_WarriorRagingBlow {
				warrior.RagingBlow.CD.Timer.Set(time.Duration(*warrior.RagingBlow.CD.Timer) - time.Second*1)
			}
		},
	})

	warrior.RagingBlow = warrior.RegisterSpell(AnyStance, core.SpellConfig{
		SpellCode:   SpellCode_WarriorRagingBlow,
		ActionID:    core.ActionID{SpellID: 402911},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagOffensive,

		RageCost: core.RageCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.IsEnraged()
		},

		CritDamageBonus: warrior.impale(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
