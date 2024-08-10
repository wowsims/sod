package itemhelpers

import (
	"fmt"

	"github.com/wowsims/sod/sim/core"
)

// Create a simple weapon proc that deals damage.
func CreateWeaponProcDamage(itemId int32, itemName string, ppm float64, spellId int32, school core.SpellSchool,
	dmgMin float64, dmgRange float64, bonusCoef float64, defType core.DefenseType) {

	core.NewItemEffect(itemId, func(agent core.Agent) {
		character := agent.GetCharacter()

		sc := core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellId},
			SpellSchool: school,
			DefenseType: defType,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: bonusCoef,
		}

		switch defType {
		case core.DefenseTypeNone:
			sc.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := dmgMin + core.TernaryFloat64(dmgRange > 0, sim.RandomFloat(itemName)*dmgRange, 0)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeAlwaysHit)
			}
		case core.DefenseTypeMagic:
			sc.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := dmgMin + core.TernaryFloat64(dmgRange > 0, sim.RandomFloat(itemName)*dmgRange, 0)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			}
		case core.DefenseTypeMelee:
			sc.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := dmgMin + core.TernaryFloat64(dmgRange > 0, sim.RandomFloat(itemName)*dmgRange, 0)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		case core.DefenseTypeRanged:
			sc.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := dmgMin + core.TernaryFloat64(dmgRange > 0, sim.RandomFloat(itemName)*dmgRange, 0)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeRangedHitAndCrit)
			}
		}

		procSpell := character.RegisterSpell(sc)
		procMask := character.GetProcMaskForItem(itemId)
		ppmm := character.AutoAttacks.NewPPMManager(ppm, procMask)

		character.GetOrRegisterAura(core.Aura{
			Label:    fmt.Sprintf("%s Proc Aura", itemName),
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) {
					return
				}
				if result.Landed() && ppmm.Proc(sim, spell.ProcMask, aura.Label) {
					procSpell.Cast(sim, result.Target)
				}
			},
		})
	})
}

// Create a weapon proc using a custom spell.
func CreateWeaponProcSpell(itemId int32, itemName string, ppm float64, procSpellGenerator func(character *core.Character) *core.Spell) {
	core.NewItemEffect(itemId, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := procSpellGenerator(character)

		procMask := character.GetProcMaskForItem(itemId)
		ppmm := character.AutoAttacks.NewPPMManager(ppm, procMask)
		procLabel := itemName + " Proc"

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: itemName + " Proc Aura",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) {
					return
				}
				if result.Landed() && ppmm.Proc(sim, spell.ProcMask, procLabel) {
					procSpell.Cast(sim, result.Target)
				}
			},
		}))
	})
}

// Create a weapon proc for a custom aura.
func CreateWeaponProcAura(itemId int32, itemName string, ppm float64, procAuraGenerator func(character *core.Character) *core.Aura) {
	core.NewItemEffect(itemId, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := procAuraGenerator(character)

		procMask := character.GetProcMaskForItem(itemId)
		ppmm := character.AutoAttacks.NewPPMManager(ppm, procMask)
		procLabel := itemName + " Proc"

		core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:    itemName + " Proc Aura",
			Duration: core.NeverExpires,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) {
					return
				}
				if result.Landed() && ppmm.Proc(sim, spell.ProcMask, procLabel) {
					procAura.Activate(sim)
				}
			},
		}))
	})
}
