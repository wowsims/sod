package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) RegisterOnItemSwapWithImbue(effectID int32, procMask *core.ProcMask, aura *core.Aura) {
	shaman.RegisterItemSwapCallback(core.MeleeWeaponSlots(), func(sim *core.Simulation, _ proto.ItemSlot, _ bool) {
		mask := core.ProcMaskUnknown
		if shaman.MainHand().TempEnchant == effectID {
			mask |= core.ProcMaskMeleeMH
		}
		if shaman.OffHand().TempEnchant == effectID {
			mask |= core.ProcMaskMeleeOH
		}
		*procMask = mask

		if mask == core.ProcMaskUnknown {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}
