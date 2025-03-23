package item_effects

import (
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
		character := agent.GetCharacter()

		var classMask uint64
		switch character.Class {
		// https://www.wowhead.com/classic/spell=1232094/infusion-of-souls
		case proto.Class_ClassDruid:
			classMask = druid.ClassSpellMask_DruidHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1230948/infusion-of-souls
		case proto.Class_ClassMage:
			classMask = mage.ClassSpellMask_MageHarmfulGCDSpells

		// https://www.wowhead.com/classic/spell=1232104/infusion-of-souls
		case proto.Class_ClassPaladin:
			// Explicitly lists that it does not work for Holy Shock in the tooltip https://www.wowhead.com/classic-ptr/item=241039/infusion-of-souls?spellModifier=462814
			classMask = paladin.ClassSpellMask_PaladinHarmfulGCDSpells ^ paladin.ClassSpellMask_PaladinHolyShock

		// https://www.wowhead.com/classic/spell=1232095/infusion-of-souls
		case proto.Class_ClassPriest:
			// Explicitly lists that it does not work for Penance in the tooltip https://www.wowhead.com/classic-ptr/item=241039/infusion-of-souls?spellModifier=440247
			classMask = priest.ClassSpellMask_PriestHarmfulGCDSpells ^ priest.ClassSpellMask_PriestPenance

		// https://www.wowhead.com/classic/spell=1232096/infusion-of-souls
		case proto.Class_ClassShaman:
			// Explicitly lists that it does not work while Way of Earth is active
			classMask = core.Ternary(agent.(shaman.ShamanAgent).GetShaman().WayOfEarthActive(), 0, shaman.ClassSpellMask_ShamanHarmfulGCDSpells)

		// https://www.wowhead.com/classic/spell=1232093/infusion-of-souls
		case proto.Class_ClassWarlock:
			classMask = warlock.ClassSpellMask_WarlockHarmfulGCDSpells
		}

		character.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Matches(classMask) {
				spell.AllowGCDHasteScaling = true
			}
		})
	})

	core.AddEffectsToTest = true
}
