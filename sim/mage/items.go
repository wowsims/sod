package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	FireRuby                     = 20036
	StaffOfOrder                 = 229909
	StaffOfInferno               = 229971
	StaffOfRime                  = 229972
	MindQuickeningGem            = 230243
	HazzarahsCharmOfChilledMagic = 231282
	JewelOfKajaro                = 231324
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(FireRuby, func(agent core.Agent) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: FireRuby}
		manaMetrics := character.NewManaMetrics(actionID)

		damageAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Chaos Fire",
			ActionID: core.ActionID{SpellID: 24389},
			Duration: time.Minute * 1,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.FirePower, 100)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.FirePower, -100)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					aura.Deactivate(sim)
				}
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AddMana(sim, sim.Roll(1, 500), manaMetrics)
				damageAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	// https://www.wowhead.com/classic/item=231282/hazzarahs-charm-of-chilled-magic
	// Use: Increases the critical hit chance of your Frostbolt and Frozen Orb spells by 5%, and increases the critical hit damage of your Frostbolt and Frozen Orb spells by 50% for 20 sec.
	// (2 Min Cooldown)
	core.NewItemEffect(HazzarahsCharmOfChilledMagic, func(agent core.Agent) {
		mage := agent.(MageAgent).GetMage()

		duration := time.Second * 20
		affectedSpells := []*core.Spell{}

		aura := mage.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: HazzarahsCharmOfChilledMagic},
			Label:    "Frost Potency",
			Duration: duration,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				affectedSpells = core.FilterSlice(
					core.Flatten([][]*core.Spell{
						mage.Frostbolt,
						{mage.SpellfrostBolt},
						{mage.FrozenOrbTick},
					}),
					func(spell *core.Spell) bool { return spell != nil },
				)
			},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				for _, spell := range affectedSpells {
					spell.BonusCritRating += 5 * core.SpellCritRatingPerCritChance
					spell.CritDamageBonus += 0.50
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				for _, spell := range affectedSpells {
					spell.BonusCritRating -= 5 * core.SpellCritRatingPerCritChance
					spell.CritDamageBonus -= 0.50
				}
			},
		})

		spell := mage.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: HazzarahsCharmOfChilledMagic},
			SpellSchool: core.SpellSchoolArcane,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    mage.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    mage.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		mage.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=231324/jewel-of-kajaro
	// Equip: Reduces the cooldown on your Frozen Orb spell by 10 sec.
	core.NewItemEffect(JewelOfKajaro, func(agent core.Agent) {
		mage := agent.(MageAgent).GetMage()
		if !mage.HasRune(proto.MageRune_RuneCloakFrozenOrb) {
			return
		}

		mage.RegisterAura(core.Aura{
			Label: "Decreased Frozen Orb Cooldown",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				mage.FrozenOrb.CD.FlatModifier -= time.Second * 10
			},
		})
	})

	// https://www.wowhead.com/classic/item=230243/mind-quickening-gem
	// Use: Quickens the mind, increasing the Mage's casting speed of non-channeled spells by 33% for 20 sec. (2 Min Cooldown)
	core.NewItemEffect(MindQuickeningGem, func(agent core.Agent) {
		mage := agent.(MageAgent).GetMage()

		actionID := core.ActionID{ItemID: MindQuickeningGem}
		duration := time.Second * 20

		buffAura := mage.RegisterAura(core.Aura{
			ActionID: actionID,
			Label:    "Mind Quickening",
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				mage.MultiplyCastSpeed(1.33)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				mage.MultiplyCastSpeed(1 / 1.33)
			},
		})

		spell := mage.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    mage.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: core.Cooldown{
					Timer:    mage.GetOffensiveTrinketCD(),
					Duration: duration,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		mage.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityBloodlust,
			Type:     core.CooldownTypeDPS,
		})
	})

	// https://www.wowhead.com/classic/item=229971/staff-of-inferno
	// Equip: When Improved Scorch is talented, targets hit by your Blast Wave will also have 5 stacks of Fire Vulnerability applied to them.
	core.NewItemEffect(StaffOfInferno, func(agent core.Agent) {
		mage := agent.(MageAgent).GetMage()
		if mage.Talents.ImprovedScorch == 0 {
			return
		}

		core.MakePermanent(mage.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469237},
			Label:    "Staff of Inferno",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellCode == SpellCode_MageBlastWave && result.Landed() {
					aura := mage.ImprovedScorchAuras.Get(result.Target)
					aura.Activate(sim)
					aura.SetStacks(sim, 5)
				}
			},
		}))
	})

	core.NewItemEffect(StaffOfOrder, func(agent core.Agent) {
		mage := agent.(MageAgent).GetMage()
		if !mage.Talents.PresenceOfMind {
			return
		}

		core.MakePermanent(mage.RegisterAura(core.Aura{
			Label: "Staff of Order",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellSchool == core.SpellSchoolArcane && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && result.Landed() {
					mage.PresenceOfMind.CD.Set(mage.PresenceOfMind.CD.ReadyAt() - time.Second)
				}
			},
		}))
	})

	core.NewItemEffect(StaffOfRime, func(agent core.Agent) {
		mage := agent.(MageAgent).GetMage()
		if !mage.Talents.IceBarrier {
			return
		}

		statsAura := mage.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 469238},
			Label:    "Staff of Rime",
			Duration: time.Minute,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				mage.AddStatDynamic(sim, stats.FrostPower, 100)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				mage.AddStatDynamic(sim, stats.FrostPower, -100)
			},
		})

		mage.RegisterAura(core.Aura{
			Label: "Staff of Rime Dummy",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				for _, aura := range mage.IceBarrierAuras {
					if aura == nil {
						continue
					}

					oldOnGain := aura.OnGain
					aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
						oldOnGain(aura, sim)
						statsAura.Activate(sim)
					}

					oldOnExpire := aura.OnExpire
					aura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
						oldOnExpire(aura, sim)
						statsAura.Deactivate(sim)
					}
				}
			},
		})
	})

	core.AddEffectsToTest = true
}
