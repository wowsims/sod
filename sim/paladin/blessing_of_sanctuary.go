package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	//"github.com/wowsims/sod/sim/core/stats"
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
		{minLevel: 60, maxLevel: 61, spellID: 20914, absorb: 24, damage: 35},
	}

	for i, values := range sanctuaryValues {

		if values.minLevel <= paladin.Level && paladin.Level < values.maxLevel {

			rank := i + 1
			actionID := core.ActionID{SpellID: values.spellID}
			//absorb := values.absorb
			damage := values.damage

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
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidBlock() {
						sanctuaryProc.Cast(sim, spell.Unit)
					}
				},
			})

		}
	}
}
