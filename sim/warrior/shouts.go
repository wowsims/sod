package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const ShoutExpirationThreshold = time.Second * 3

func (warrior *Warrior) newShoutSpellConfig(actionID core.ActionID, allyAuras core.AuraArray) *core.Spell {
	return warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, aura := range allyAuras {
				if aura != nil {
					aura.Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{allyAuras},
	})
}

func (warrior *Warrior) registerBattleShout() {
	rank := core.LevelToBuffRank[core.BattleShout][warrior.Level]
	actionId := core.BattleShoutSpellId[rank]

	warrior.BattleShout = warrior.newShoutSpellConfig(core.ActionID{SpellID: actionId}, warrior.NewPartyAuraArray(func(unit *core.Unit) *core.Aura {
		return core.BattleShoutAura(unit, warrior.Talents.ImprovedBattleShout, warrior.Talents.BoomingVoice)
	}))
}

func (warrior *Warrior) registerShouts() {
	warrior.registerBattleShout()
}
