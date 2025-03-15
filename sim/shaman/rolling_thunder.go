package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) registerRollingThunder() {
	if !shaman.HasRune(proto.ShamanRune_RuneBracersRollingThunder) {
		return
	}

	actionID := core.ActionID{SpellID: 432129}
	impLightningShieldBonus := []float64{1, 1.05, 1.10, 1.15}[shaman.Talents.ImprovedLightningShield]
	manaMetrics := shaman.NewManaMetrics(actionID)

	shaman.rollingThunderProcChance += 0.50

	// Casts handled in lightning_shield.go
	shaman.RollingThunder = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: ClassSpellMask_ShamanRollingThunder,
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if shaman.ActiveShield == nil || shaman.ActiveShield.ClassSpellMask != ClassSpellMask_ShamanLightningShield {
				return
			}

			rank := shaman.ActiveShield.Rank
			numCharges := float64(shaman.ActiveShieldAura.GetStacks() - 3)

			// TODO: Need a better way to get a spell's base damage directly from a spell
			chargeDamage := LightningShieldBaseDamage[rank]*impLightningShieldBonus + LightningShieldSpellCoef[rank]*shaman.LightningShieldProcs[rank].GetBonusDamage()
			spell.CalcAndDealDamage(sim, target, chargeDamage*numCharges, spell.OutcomeMagicCrit)

			shaman.AddMana(sim, .02*numCharges*shaman.MaxMana(), manaMetrics)
			shaman.ActiveShieldAura.SetStacks(sim, 3)
		},
	})

	affectedSpellClassMasks := ClassSpellMask_ShamanLightningBolt | ClassSpellMask_ShamanChainLightning

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label: "Rolling Thunder Trigger",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, aura := range shaman.LightningShieldAuras {
				if aura != nil {
					aura.MaxStacks = 9
				}
			}
		},
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if shaman.ActiveShield == nil || shaman.ActiveShield.ClassSpellMask != ClassSpellMask_ShamanLightningShield {
				return
			}

			if spell.Matches(ClassSpellMask_ShamanEarthShock) && shaman.ActiveShieldAura.GetStacks() > 3 {
				shaman.RollingThunder.Cast(sim, result.Target)
			} else if spell.Matches(affectedSpellClassMasks) && sim.Proc(shaman.rollingThunderProcChance, "Rolling Thunder") {
				shaman.ActiveShieldAura.AddStack(sim)
			}
		},
	}))
}
