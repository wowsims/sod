package item_effects

import (
	"time"

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
	AbandonedExperiment = 241037
	InfusionOfSouls     = 241039
	HandOfRebornJustice = 242310
)

func init() {
	core.AddEffectsToTest = false

	/* ! Please keep items ordered alphabetically ! */

	// https://www.wowhead.com/classic-ptr/item=241037/abandoned-experiment
	// Use: After drinking the experiment, ranged or melee attacks increase your attack speed by 2% for 30 sec.
	// This effect stacks up to 15 times. (2 Min Cooldown)
	core.NewItemEffect(AbandonedExperiment, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: AbandonedExperiment}
		duration := time.Second * 30

		buffAura := character.RegisterAura(core.Aura{
			ActionID:  actionID,
			Label:     "Abandoned Experiment",
			MaxStacks: 15,
			Duration:  duration,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				character.MultiplyAttackSpeed(sim, 1/(1+0.02*float64(oldStacks)))
				character.MultiplyAttackSpeed(sim, 1+0.02*float64(newStacks))
			},
		})

		buffAura.MakeDependentProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Abandoned Experiment Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMelee | core.ProcMaskRanged,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !buffAura.IsActive() {
					// Should not refresh
					buffAura.Activate(sim)
				}
				buffAura.AddStack(sim)
			},
		})

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})
		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: cdSpell,
		})
	})

	// https://www.wowhead.com/classic-ptr/item=242310/hand-of-reborn-justice
	// Equip: 2% chance on melee or ranged hit to gain 1 extra attack. (Proc chance: 2%, 2s cooldown)
	core.NewItemEffect(HandOfRebornJustice, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Hand of Reborn Justice Trigger (Melee)",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskMelee | core.ProcMaskRanged,
			SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
			ProcChance:        0.02,
			ICD:               time.Second * 2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMelee) {
					spell.Unit.AutoAttacks.ExtraMHAttackProc(sim, 1, core.ActionID{SpellID: 1232044}, spell)
				} else {
					character.AutoAttacks.StoreExtraRangedAttack(sim, 1, core.ActionID{SpellID: 1213381}, spell.ActionID)
				}
			},
		})

		character.ItemSwap.RegisterProc(HandOfInjustice, triggerAura)
	})

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
