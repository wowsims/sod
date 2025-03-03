package hunter

import (
	"github.com/wowsims/sod/sim/core"
)

func (hunter *Hunter) getHuntersMark(rank int) core.SpellConfig {
	spellId := [5]int32{0, 1130, 14323, 14324, 14325}[rank]
	manaCost := [5]float64{0, 15, 30, 45, 60}[rank]
	level := [5]int{0, 6, 22, 40, 58}[rank]

	hunter.HuntersMarkAuras = hunter.NewEnemyAuraArray(hunter.myHuntersMarkAura)

	return core.SpellConfig{
		ClassSpellMask: ClassSpellMask_HunterHuntersMark,
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolArcane,
		DefenseType:    core.DefenseTypeNone,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		Rank:           rank,
		RequiredLevel:  level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura := hunter.HuntersMarkAuras.Get(target)
			// Can only have 1 hunter's mark up at a time
			for _, auras := range hunter.HuntersMarkAuras {
				if auras.IsActive() {
					auras.Deactivate(sim)
				}
			}
			aura.Activate(sim)
		},

		RelatedAuras: []core.AuraArray{hunter.HuntersMarkAuras},
	}
}

func (hunter *Hunter) myHuntersMarkAura(target *core.Unit, playerLevel int32) *core.Aura {
	return core.HuntersMarkAura(target, hunter.Talents.ImprovedHuntersMark, hunter.Level)
}

func (hunter *Hunter) registerHuntersMark() {
	maxRank := 4
	for i := 1; i <= maxRank; i++ {
		config := hunter.getHuntersMark(i)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.HuntersMark = append(hunter.HuntersMark, hunter.GetOrRegisterSpell(config))
		}
	}
	
}
