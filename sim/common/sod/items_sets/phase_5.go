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
	
			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Millisecond * 100,
			}
	
			character.GetOrRegisterAura(core.Aura{
				Label:    "Twin Blades of the Hakkari",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Flags.Matches(core.SpellFlagSuppressEquipProcs) {
						return
					}
					if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) && icd.IsReady(sim) && sim.Proc(0.02, "Twin Blades of the Hakkari") {
						icd.Use(sim)
						aura.Unit.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 468255}, spell)
					}
				},
			})	
		},

	},
})
