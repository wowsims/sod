package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) registerInfernalArmorCD() {
	if !warlock.HasRune(proto.WarlockRune_RuneCloakInfernalArmor) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.WarlockRune_RuneCloakInfernalArmor)}

	// TODO: Unsure if there's a better way to do this
	physResistanceMultiplier := 1.0
	infernalArmorAura := warlock.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Infernal Armor",
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			attackTable := warlock.CurrentTarget.AttackTables[warlock.UnitIndex][proto.CastType_CastTypeMainHand]
			physResistanceMultiplier = 1 - attackTable.GetArmorDamageModifier()

			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= physResistanceMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= physResistanceMultiplier
			warlock.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= physResistanceMultiplier
		},
	})

	spell := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			infernalArmorAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}
