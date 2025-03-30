package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
)

const IceBarrierRanks = 4

var IceBarrierSpellId = [IceBarrierRanks + 1]int32{0, 11426, 13031, 13032, 13033}
var IceBarrierManaCost = [IceBarrierRanks + 1]float64{0, 305, 360, 420, 480}
var IceBarrierLevel = [IceBarrierRanks + 1]int{0, 40, 46, 52, 58}

func (mage *Mage) registerIceBarrierSpell() {
	mage.IceBarrier = make([]*core.Spell, IceBarrierRanks+1)
	mage.IceBarrierAuras = make([]*core.Aura, IceBarrierRanks+1)

	cdTimer := mage.NewTimer()

	for rank := 1; rank <= IceBarrierRanks; rank++ {
		config := mage.newIceBarrierSpellConfig(rank, cdTimer)

		if config.RequiredLevel <= int(mage.Level) {
			mage.IceBarrier[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newIceBarrierSpellConfig(rank int, cdTimer *core.Timer) core.SpellConfig {
	spellID := IceBarrierSpellId[rank]
	manaCost := IceBarrierManaCost[rank]
	level := IceBarrierLevel[rank]

	cooldown := time.Second * 30

	mage.IceBarrierAuras[rank] = mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: spellID},
		Label:    fmt.Sprintf("Ice Barrier (Rank %d)", rank),
		Duration: time.Minute,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Dummy
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// Dummy
		},
	})

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellID},
		SpellSchool:   core.SpellSchoolFrost,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,
		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cooldown,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Disable an existing barrier
			if mage.activeBarrier != nil && mage.activeBarrier.IsActive() {
				mage.activeBarrier.Deactivate(sim)
				mage.activeBarrier = mage.IceBarrierAuras[rank]
			}
			mage.IceBarrierAuras[rank].Activate(sim)
		},
	}
}
