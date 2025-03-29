package item_sets

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/hunter"
	"github.com/wowsims/sod/sim/rogue"
	"github.com/wowsims/sod/sim/warrior"
)

var ItemSetFallenRegality = core.NewItemSet(core.ItemSet{
	Name: "Fallen Regality",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			aura := core.MakePermanent(character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 1232184},
				Label:    "Fallen Regality",
			}))

			switch character.Class {
			// Damaging finishing moves have a 20% chance per combo point to restore 20 energy.
			case proto.Class_ClassRogue:
				roguePlayer := agent.(rogue.RogueAgent).GetRogue()
				energyMetrics := roguePlayer.NewEnergyMetrics(core.ActionID{SpellID: 1232184})
				aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
					roguePlayer.OnComboPointsSpent(func(sim *core.Simulation, spell *core.Spell, comboPoints int32) {
						if spell.ProcMask != core.ProcMaskEmpty && sim.Proc(0.20*float64(comboPoints), "Fallen Regality Proc") {
							roguePlayer.AddEnergy(sim, 20, energyMetrics)
						}
					})
				})

			// Flanking Strike's damage buff is increased by an additional 2% per stack. When striking from behind, your target takes 150% increased damage from Flanking Strike.
			case proto.Class_ClassHunter:
				hunterPlayer := agent.(hunter.HunterAgent).GetHunter()
				if !hunterPlayer.HasRune(proto.HunterRune_RuneLegsFlankingStrike) {
					return
				}

				flankingBuffDamageMod := hunterPlayer.AddDynamicMod(core.SpellModConfig{
					Kind:     core.SpellMod_DamageDone_Flat,
					ProcMask: core.ProcMaskMelee,
				})

				if !character.PseudoStats.InFrontOfTarget {
					hunterPlayer.AddStaticMod(core.SpellModConfig{
						Kind:       core.SpellMod_DamageDone_Pct,
						ClassMask:  hunter.ClassSpellMask_HunterFlankingStrike | hunter.ClassSpellMask_HunterPetFlankingStrike,
						FloatValue: 1.50,
					})
				}

				aura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
					hunterPlayer.FlankingStrike.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
						flankingBuffDamageMod.Activate()
					}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
						flankingBuffDamageMod.Deactivate()
					}).ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
						flankingBuffDamageMod.UpdateIntValue(int64(2 * newStacks))
					})
				})

			// If Cleave hits fewer than its maximum number of targets, it deals 35% more damage for each unused bounce.
			case proto.Class_ClassWarrior:
				warriorPlayer := agent.(warrior.WarriorAgent).GetWarrior()
				targetCount := warriorPlayer.Env.GetNumTargets()

				cleaveDamageMod := warriorPlayer.AddDynamicMod(core.SpellModConfig{
					Kind:       core.SpellMod_DamageDone_Pct,
					ClassMask:  warrior.ClassSpellMask_WarriorCleave,
					FloatValue: 1,
				})

				aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					cleaveDamageMod.Activate()
					// The cleave target count is set during initializing, so set the value here
					cleaveDamageMod.UpdateFloatValue(1 + float64(warriorPlayer.CleaveTargetCount-targetCount)*0.35)
				}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					cleaveDamageMod.Activate()
				})
			}
		},
	},
})

var ItemSetHackAndSmash = core.NewItemSet(core.ItemSet{
	Name: "Hack and Smash",
	Bonuses: map[int32]core.ApplyEffect{
		// The Fire and Nature damage increases from Mercy and Crimson Cleaver are increased by 10%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			fireAura := character.GetAuraByID(core.ActionID{SpellID: 1231498})
			fireAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				School:     core.SpellSchoolFire,
				FloatValue: 1.30 / 1.20, // Revert the 20% and apply 30%
			})

			natureAura := character.GetAuraByID(core.ActionID{SpellID: 1231456})
			natureAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				School:     core.SpellSchoolNature,
				FloatValue: 1.30 / 1.20, // Revert the 20% and apply 30%
			})
		},
	},
})
