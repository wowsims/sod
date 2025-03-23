package item_effects

import (
	"fmt"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/druid"
	"github.com/wowsims/sod/sim/mage"
	"github.com/wowsims/sod/sim/paladin"
	"github.com/wowsims/sod/sim/priest"
	"github.com/wowsims/sod/sim/shaman"
	"github.com/wowsims/sod/sim/warlock"
)

const (
	InfusionOfSouls = 241039
)

func init() {
	core.AddEffectsToTest = false

	// https://www.wowhead.com/classic/item=241039/infusion-of-souls
	// The Global Cooldown caused by your non-weapon based damaging spells can be reduced by Spell Haste, up to a 0.5 second reduction.
	core.NewItemEffect(InfusionOfSouls, func(agent core.Agent) {
		var classMask uint64
		switch agent.GetCharacter().Class {
		// https://www.wowhead.com/classic/spell=1232094/infusion-of-souls
		case proto.Class_ClassDruid:
			classMask = druid.ClassSpellMask_DruidHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1230948/infusion-of-souls
		case proto.Class_ClassMage:
			classMask = mage.ClassSpellMask_MageHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1232104/infusion-of-souls
		case proto.Class_ClassPaladin:
			classMask = paladin.ClassSpellMask_PaladinExorcism | paladin.ClassSpellMask_PaladinHolyWrath | paladin.ClassSpellMask_PaladinConsecration

		// https://www.wowhead.com/classic/spell=1232095/infusion-of-souls
		case proto.Class_ClassPriest:
			classMask = priest.ClassSpellMask_PriestHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1232096/infusion-of-souls
		case proto.Class_ClassShaman:
			classMask = shaman.ClassSpellMask_ShamanHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1232093/infusion-of-souls
		case proto.Class_ClassWarlock:
			classMask = warlock.ClassSpellMask_WarlockHarmfulGCDSpells
		}

		fmt.Println(classMask)
	})

	core.AddEffectsToTest = true
}
