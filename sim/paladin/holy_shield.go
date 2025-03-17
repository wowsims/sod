package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var HolyShieldValues = []struct {
	level    int32
	spellID  int32
	procID   int32
	manaCost float64
	damage   float64
}{
	{level: 30, spellID: 20925, procID: 20955, manaCost: 150, damage: 65},
	{level: 50, spellID: 20927, procID: 20956, manaCost: 195, damage: 95},
	{level: 60, spellID: 20928, procID: 20957, manaCost: 240, damage: 130},
}

func (paladin *Paladin) registerHolyShield() {
	if !paladin.Talents.HolyShield {
		return
	}

	numCharges := int32(4)
	defendersResolveSPAura := core.DefendersResolveSpellDamage(paladin.GetCharacter(), 2)
	blockBonus := 30.0 * core.BlockRatingPerBlockChance

	for i, values := range HolyShieldValues {
		rank := i + 1
		level := values.level
		spellID := values.spellID
		procID := values.procID
		manaCost := values.manaCost
		damage := values.damage

		if paladin.Level < level {
			break
		}

		paladin.holyShieldProc[i] = paladin.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: procID},
			ClassSpellMask: ClassSpellMask_PaladinHolyShieldProc,
			SpellSchool:    core.SpellSchoolHoly,
			DefenseType:    core.DefenseTypeMagic,
			ProcMask:       core.ProcMaskSpellDamage,

			RequiredLevel: int(level),
			Rank:          rank,

			DamageMultiplier: 1,
			ThreatMultiplier: 1.2,
			BonusCoefficient: 0.05,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// Spell damage from Holy Shield can crit, but does not miss.
				spell.CalcAndDealDamage(sim, target, paladin.getHolyShieldDamage(sim, damage), spell.OutcomeMagicCrit)
			},
		})

		paladin.holyShieldAura[i] = paladin.RegisterAura(core.Aura{
			Label:     "Holy Shield" + paladin.Label + strconv.Itoa(rank),
			ActionID:  core.ActionID{SpellID: spellID},
			Duration:  time.Second * 10,
			MaxStacks: numCharges,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if aura.MaxStacks > 0 {
					aura.SetStacks(sim, aura.MaxStacks)
				}
				paladin.AddStatDynamic(sim, stats.Block, blockBonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				paladin.AddStatDynamic(sim, stats.Block, -blockBonus)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.DidBlock() {
					paladin.holyShieldProc[i].Cast(sim, spell.Unit)
					if aura.MaxStacks > 0 {
						aura.RemoveStack(sim)
					}
				}
			},
		})

		paladin.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: spellID},
			ClassSpellMask: ClassSpellMask_PaladinHolyShield,
			Flags:          core.SpellFlagAPL,
			RequiredLevel:  int(level),
			Rank:           rank,
			ManaCost: core.ManaCostOptions{
				FlatCost: manaCost,
			},
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				CD: core.Cooldown{
					Timer:    paladin.NewTimer(),
					Duration: time.Second * 10,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				paladin.holyShieldAura[i].Activate(sim)

				if stacks := int32(paladin.GetStat(stats.Defense)); stacks > 0 {
					defendersResolveSPAura.Activate(sim)

					if defendersResolveSPAura.GetStacks() != stacks {
						defendersResolveSPAura.SetStacks(sim, stacks)
					}
				}
			},
		})
	}
}

func (paladin *Paladin) getHolyShieldDamage(sim *core.Simulation, baseDamage float64) float64 {
	damage := baseDamage

	if paladin.holyShieldExtraDamage != nil {
		damage += paladin.holyShieldExtraDamage(sim, paladin)
	}

	return damage
}
