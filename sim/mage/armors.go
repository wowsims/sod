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

	stats := stats.Stats{
		stats.Armor:           armor,
		stats.FrostResistance: frostRes,
	}

	mage.IceArmorAura = core.MakePermanent(mage.RegisterAura(core.Aura{
		Label:    "Ice Armor",
		ActionID: core.ActionID{SpellID: spellID},
	}).AttachBuildPhaseStatsBuff(stats))
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

	mage.MageArmorAura = core.MakePermanent(mage.RegisterAura(core.Aura{
		Label:      "Mage Armor",
		ActionID:   core.ActionID{SpellID: spellID},
		BuildPhase: core.CharacterBuildPhaseBuffs,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.SpiritRegenRateCasting += .3
			mage.AddBuildPhaseResistancesDynamic(sim, spellRes)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.PseudoStats.SpiritRegenRateCasting -= .3
			mage.AddBuildPhaseResistancesDynamic(sim, -1*spellRes)
		},
	}))
}

func (mage *Mage) applyMoltenArmor() {
	if !mage.HasRune(proto.MageRune_RuneBracersMoltenArmor) {
		return
	}

	crit := 5.0 * core.SpellCritRatingPerCritChance

	mage.MoltenArmorAura = core.MakePermanent(mage.RegisterAura(core.Aura{
		Label:      "Molten Armor",
		ActionID:   core.ActionID{SpellID: int32(proto.MageRune_RuneBracersMoltenArmor)},
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachBuildPhaseStatBuff(stats.SpellCrit, crit))
}
