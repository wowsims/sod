package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////////////////
//                                 Other
///////////////////////////////////////////////////////////////////////////

var ItemSetTwinBladesofHakkari = core.NewItemSet(core.ItemSet{
	Name: "The Twin Blades of Hakkari",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases Swords +3
		// 2% chance on melee hit to gain 1 extra attack.  (1%, 100ms cooldown)
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SwordsSkill += 3
			if !character.AutoAttacks.AutoSwingMelee {
				return
			}
			
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     "Twin Blades of the Hakkari",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskMelee,
				SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
				ProcChance: 0.02,
				ICD:      time.Millisecond * 100,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					spell.Unit.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 468255}, spell)
				},
			})
		},

	},
})
