package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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
			physResistanceMultiplier = attackTable.GetArmorDamageModifier()
			warlock.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(physResistanceMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageTakenMultiplier.MultiplyMagicSchools(1 / physResistanceMultiplier)
		},
	})

	warlock.InfernalArmor = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_WarlockInfernalArmor,
		SpellSchool:    core.SpellSchoolShadow,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			infernalArmorAura.Activate(sim)
		},

		RelatedSelfBuff: infernalArmorAura,
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: warlock.InfernalArmor,
		Type:  core.CooldownTypeSurvival,
	})
}
