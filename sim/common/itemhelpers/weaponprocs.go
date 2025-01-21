package itemhelpers

import (
	"fmt"

	"github.com/wowsims/sod/sim/core"
)

// Create a simple weapon proc that deals damage.
func CreateWeaponProcDamage(itemID int32, itemName string, ppm float64, spellID int32, school core.SpellSchool,
	dmgMin float64, dmgRange float64, bonusCoef float64, defType core.DefenseType, spellFlagExclude core.SpellFlag) {

	core.NewItemEffect(itemID, func(agent core.Agent) {
		character := agent.GetCharacter()

		sc := core.SpellConfig{
			ActionID:    core.ActionID{SpellID: spellID},
			SpellSchool: school,
			DefenseType: defType,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

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
			// "Phantom Strike Procs"
			// Can proc itself (Only for CoH proc), can't proc equip effects (in SoD at least - Tested), Weapon Enchants (confirmed - procs fiery), can proc imbues (oils),
			// WildStrikes/Windfury (Wound/ Phantom Strike can't proc WF/WS in SoD, Tested for both, Appear to behave like equip affects in SoD)
			sc.ProcMask = core.ProcMaskMeleeSpecial
			sc.Flags = core.SpellFlagSuppressEquipProcs

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
		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(itemID, ppm, 0)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    fmt.Sprintf("%s Proc Aura", itemName),
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(spellFlagExclude) { // This check is needed because it doesn't use aura_helpers ApplyProcTriggerCallback
					return
				}
				if result.Landed() && dpm.Proc(sim, spell.ProcMask, aura.Label) {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterProc(itemID, aura)
	})
}

// Creates a weapon proc for "Chance on Hit" weapon effects
func CreateWeaponCoHProcDamage(itemId int32, itemName string, ppm float64, spellId int32, school core.SpellSchool,
	dmgMin float64, dmgRange float64, bonusCoef float64, defType core.DefenseType) {

	CreateWeaponProcDamage(itemId, itemName, ppm, spellId, school, dmgMin, dmgRange, bonusCoef, defType, core.SpellFlagSuppressWeaponProcs)
}

// Creates a weapon proc for "Equip" weapon effects
func CreateWeaponEquipProcDamage(itemId int32, itemName string, ppm float64, spellId int32, school core.SpellSchool,
	dmgMin float64, dmgRange float64, bonusCoef float64, defType core.DefenseType) {

	CreateWeaponProcDamage(itemId, itemName, ppm, spellId, school, dmgMin, dmgRange, bonusCoef, defType, core.SpellFlagSuppressEquipProcs)
}

// Create a weapon proc using a custom spell.
func CreateWeaponProcSpell(itemId int32, itemName string, ppm float64, procSpellGenerator func(character *core.Character) *core.Spell) {
	core.NewItemEffect(itemId, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := procSpellGenerator(character)
		procSpell.Flags |= core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              itemName + " Proc",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			DPM:               character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(itemId, ppm, 0),
			DPMProcType:       core.DPMProcNoWeaponSpecials,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterProc(itemId, aura)
	})
}

// Create a weapon proc for a custom aura.
func CreateWeaponProcAura(itemID int32, itemName string, ppm float64, procAuraGenerator func(character *core.Character) *core.Aura) {
	core.NewItemEffect(itemID, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := procAuraGenerator(character)

		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(itemID, ppm, 0)
		procLabel := itemName + " Proc"

		aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label:    itemName + " Proc Aura",
			Duration: core.NeverExpires,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(core.SpellFlagSuppressWeaponProcs) {
					return
				}
				if result.Landed() && dpm.Proc(sim, spell.ProcMask, procLabel) {
					procAura.Activate(sim)
				}
			},
		}))

		character.ItemSwap.RegisterProc(itemID, aura)
	})
}
