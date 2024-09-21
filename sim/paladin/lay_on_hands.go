package paladin

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (paladin *Paladin) registerLayOnHands() {

	minLevels := []int32{50, 30, 10}
	idx := slices.IndexFunc(minLevels, func(level int32) bool {
		return paladin.Level >= level
	})

	if idx == -1 {
		return
	}

	spellID := []int32{10310, 2800, 633}[idx]
	manaReturn := []float64{550, 250, 0}[idx]

	// Only register the highest available rank of LoH (no benefit to using lower ranks)
	actionID := core.ActionID{SpellID: spellID}
	layOnHandsManaMetrics := paladin.NewManaMetrics(actionID)
	layOnHandsHealthMetrics := paladin.NewHealthMetrics(actionID)
	layOnHands := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagAPL | core.SpellFlagMCD,
		SpellSchool: core.SpellSchoolHoly,
		SpellCode:   SpellCode_PaladinLayOnHands,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * time.Duration(60-10*paladin.Talents.ImprovedLayOnHands),
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			paladin.SpendMana(sim, paladin.CurrentMana(), layOnHandsManaMetrics)
			paladin.GainHealth(sim, paladin.MaxHealth(), layOnHandsHealthMetrics)
			paladin.AddMana(sim, manaReturn, layOnHandsManaMetrics)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell:    layOnHands,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeSurvival,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentHealthPercent() < 0.1 // TODO: better default condition
		},
	})
}
