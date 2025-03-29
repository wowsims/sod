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

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              fmt.Sprintf("%s Proc Aura", itemName),
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			DPM:               character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(itemID, ppm, 0),
			DPMProcCheck:      core.DPMProc,
			SpellFlagsExclude: spellFlagExclude,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterProc(itemID, aura)
	})
}

// Creates a weapon proc for "Chance on Hit" weapon effects
func CreateWeaponCoHProcDamage(itemID int32, itemName string, ppm float64, spellId int32, school core.SpellSchool,
	dmgMin float64, dmgRange float64, bonusCoef float64, defType core.DefenseType) {

	CreateWeaponProcDamage(itemID, itemName, ppm, spellId, school, dmgMin, dmgRange, bonusCoef, defType, core.SpellFlagSuppressWeaponProcs)
}

// Creates a weapon proc for "Equip" weapon effects
func CreateWeaponEquipProcDamage(itemID int32, itemName string, ppm float64, spellId int32, school core.SpellSchool,
	dmgMin float64, dmgRange float64, bonusCoef float64, defType core.DefenseType) {

	CreateWeaponProcDamage(itemID, itemName, ppm, spellId, school, dmgMin, dmgRange, bonusCoef, defType, core.SpellFlagSuppressEquipProcs)
}

// Create a weapon proc using a custom spell.
func CreateWeaponProcSpell(itemID int32, itemName string, ppm float64, procSpellGenerator func(character *core.Character) *core.Spell) {
	createWeaponProcSpell(itemID, itemName, ppm, core.DPMProc, procSpellGenerator)
}

func CreateFeralWeaponProcSpell(itemID int32, itemName string, ppm float64, procSpellGenerator func(character *core.Character) *core.Spell) {
	createWeaponProcSpell(itemID, itemName, ppm, core.DPMProcWithWeaponSpecials, procSpellGenerator)
}

func createWeaponProcSpell(itemID int32, itemName string, ppm float64, dpmProcCheck core.DPMProcCheck, procSpellGenerator func(character *core.Character) *core.Spell) {
	core.NewItemEffect(itemID, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := procSpellGenerator(character)
		procSpell.Flags |= core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              itemName + " Proc",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			DPM:               character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(itemID, ppm, 0),
			DPMProcCheck:      dpmProcCheck,
			SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterProc(itemID, aura)
	})
}

// Create a weapon proc for a custom aura.
func CreateWeaponProcAura(itemID int32, itemName string, ppm float64, procAuraGenerator func(character *core.Character) *core.Aura) {
	core.NewItemEffect(itemID, func(agent core.Agent) {
		AddWeaponProcAura(agent.GetCharacter(), itemID, itemName, ppm, procAuraGenerator)
	})
}

// Create a weapon proc for a custom aura and add it to an existing item effect.
func AddWeaponProcAura(character *core.Character, itemID int32, itemName string, ppm float64, procAuraGenerator func(character *core.Character) *core.Aura) {
	procAura := procAuraGenerator(character)
	dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(itemID, ppm, 0)

	aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:              itemName + " Proc",
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		SpellFlagsExclude: core.SpellFlagSuppressWeaponProcs,
		DPM:               dpm,
		DPMProcCheck:      core.DPMProc,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
			if procAura.MaxStacks > 0 {
				procAura.SetStacks(sim, procAura.MaxStacks)
			}
		},
	})

	character.ItemSwap.RegisterProc(itemID, aura)
}
