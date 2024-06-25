package mage

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (mage *Mage) applyFrostIceArmor() {
	spellID := map[int32]int32{
		25: 7301,
		40: 7320,
		50: 10219,
		60: 10220,
	}[mage.Level]

	armor := map[int32]float64{
		25: 200,
		40: 380,
		50: 470,
		60: 560,
	}[mage.Level]

	frostRes := map[int32]float64{
		25: 0,
		40: 9,
		50: 12,
		60: 15,
	}[mage.Level]

	mage.AddStat(stats.Armor, armor)
	mage.AddStat(stats.FrostResistance, frostRes)

	mage.GetOrRegisterAura(core.Aura{
		Label:    "Ice Armor",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (mage *Mage) applyMageArmor() {
	if mage.Level < 40 {
		return
	}

	spellID := map[int32]int32{
		40: 6117,
		50: 22782,
		60: 22783,
	}[mage.Level]

	spellRes := map[int32]float64{
		40: 5,
		50: 10,
		60: 15,
	}[mage.Level]

	mage.PseudoStats.SpiritRegenRateCasting += .3
	mage.AddStats(stats.Stats{
		stats.ArcaneResistance: spellRes,
		stats.FireResistance:   spellRes,
		stats.FrostResistance:  spellRes,
		stats.NatureResistance: spellRes,
		stats.ShadowResistance: spellRes,
	})

	mage.GetOrRegisterAura(core.Aura{
		Label:    "Mage Armor",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (mage *Mage) applyMoltenArmor() {
	if !mage.HasRune(proto.MageRune_RuneBracersMoltenArmor) {
		return
	}

	mage.AddStat(stats.SpellCrit, 5*core.CritRatingPerCritChance)

	mage.GetOrRegisterAura(core.Aura{
		Label:    "Molten Armor",
		ActionID: core.ActionID{SpellID: int32(proto.MageRune_RuneBracersMoltenArmor)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}
