package paladin

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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

	paladin.layOnHands = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagAPL | core.SpellFlagMCD,
		SpellSchool:    core.SpellSchoolHoly,
		ClassSpellMask: ClassSpellMask_PaladinLayOnHands,
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

			hasNaxxramasHoly2PBonus := paladin.GetAura("S03 - Item - Naxxramas - Paladin - Holy 2P Bonus")
			if hasNaxxramasHoly2PBonus != nil && hasNaxxramasHoly2PBonus.IsActive() {
				paladin.AddMana(sim, paladin.MaxMana()*0.3, layOnHandsManaMetrics)
			} else {
				paladin.AddMana(sim, manaReturn, layOnHandsManaMetrics)
			}
		},
	})
	
	if paladin.Spec == proto.Spec_SpecRetributionPaladin {
		return
	}

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell:    paladin.layOnHands,
		Priority: core.CooldownPriorityBloodlust,
		Type:     core.CooldownTypeSurvival,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentHealthPercent() < 0.1 // TODO: better default condition
		},
	})
}
