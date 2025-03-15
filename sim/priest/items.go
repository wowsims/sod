package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const (
	// Keep these ordered by ID
	CassandrasTome = 231509
	BandOfFaith    = 236112
	AtieshPriest   = 236399
)

func init() {
	core.AddEffectsToTest = false

	// Keep these ordered by name

	// https://www.wowhead.com/classic/item=236399/atiesh-greatstaff-of-the-guardian
	core.NewItemEffect(AtieshPriest, func(agent core.Agent) {
		character := agent.GetCharacter()
		aura := core.AtieshHealingEffect(&character.Unit)
		character.ItemSwap.RegisterProc(AtieshPriest, aura)
	})

	// https://www.wowhead.com/classic/item=231509/cassandras-tome
	core.NewItemEffect(CassandrasTome, func(agent core.Agent) {
		priest := agent.(PriestAgent).GetPriest()

		actionID := core.ActionID{ItemID: CassandrasTome}
		duration := time.Second * 15

		buffAura := priest.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Cassandra's Tome",
			Duration: duration,
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(ClassSpellMask_PriestAll) && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && !spell.Flags.Matches(core.SpellFlagPureDot|core.SpellFlagChanneled) {
					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt: sim.CurrentTime + core.SpellBatchWindow,
						OnAction: func(sim *core.Simulation) {
							if aura.IsActive() {
								aura.Deactivate(sim)
							}
						},
					})
				}
			},
		}).AttachSpellMod(core.SpellModConfig{
			ClassMask:         ClassSpellMask_PriestAll,
			SpellFlagsExclude: core.SpellFlagPureDot | core.SpellFlagChanneled,
			Kind:              core.SpellMod_BonusCrit_Flat,
			FloatValue:        100 * core.SpellCritRatingPerCritChance,
		})

		spell := priest.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    priest.NewTimer(),
					Duration: time.Minute * 2,
				},
				// Does not seem to share the offensive trinket timer
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		priest.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=236112/band-of-faith
	// Equip: Increases the damage dealt by your damage over time spells by 2%.
	core.NewItemEffect(BandOfFaith, func(agent core.Agent) {
		priest := agent.(PriestAgent).GetPriest()

		core.MakePermanent(priest.RegisterAura(core.Aura{
			Label: "Band of Faith",
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_PeriodicDamageDone_Flat,
			ClassMask: ClassSpellMask_PriestAll,
			IntValue:  2,
		}))
	})

	core.AddEffectsToTest = true
}
