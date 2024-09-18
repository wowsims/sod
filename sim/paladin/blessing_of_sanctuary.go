package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (paladin *Paladin) registerBlessingOfSanctuary() {
	if paladin.Options.PersonalBlessing != proto.Blessings_BlessingOfSanctuary {
		return
	}

	sanctuaryValues := []struct {
		minLevel int32
		maxLevel int32
		spellID  int32
		absorb   float64
		damage   float64
	}{
		{minLevel: 1, maxLevel: 39, spellID: 20911, absorb: 10, damage: 14},
		{minLevel: 40, maxLevel: 49, spellID: 20912, absorb: 14, damage: 21},
		{minLevel: 50, maxLevel: 49, spellID: 20913, absorb: 19, damage: 28},
		{minLevel: 60, maxLevel: 60, spellID: 20914, absorb: 24, damage: 35},
	}

	hasImpSanc := paladin.hasRune(proto.PaladinRune_RuneHeadImprovedSanctuary)
	absorbMult := core.TernaryFloat64(hasImpSanc, 2, 1)
	bonusDamage := core.TernaryFloat64(hasImpSanc, 0.3, 0.0)

	for i, values := range sanctuaryValues {

		if (values.minLevel <= paladin.Level) && (paladin.Level <= values.maxLevel) {

			rank := i + 1
			actionID := core.ActionID{SpellID: values.spellID}
			absorb := values.absorb * absorbMult
			damage := values.damage + bonusDamage*paladin.BlockValue()

			sanctuaryProc := paladin.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolHoly,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskSpellDamage,
				Flags:       core.SpellFlagIgnoreResists,

				Rank: rank,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHit)
				},
			})

			paladin.RegisterAura(core.Aura{
				Label:    "Blessing of Sanctuary Trigger",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for i := range paladin.PseudoStats.BonusDamageTakenBeforeModifiers {
						paladin.PseudoStats.BonusDamageTakenBeforeModifiers[i] -= absorb
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for i := range paladin.PseudoStats.BonusDamageTakenBeforeModifiers {
						paladin.PseudoStats.BonusDamageTakenBeforeModifiers[i] += absorb
					}
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidBlock() {
						sanctuaryProc.Cast(sim, spell.Unit)
					}
				},
			})

		}
	}
}
