package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetPlagueheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Plagueheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasDamage2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasDamage4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasDamage6PBonus()
		},
	},
})

// Increases the damage done by your Incinerate and Corruption abilities by 20%.
func (warlock *Warlock) applyNaxxramasDamage2PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Damage 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_WarlockIncinerate | ClassSpellMask_WarlockCorruption,
		IntValue:  20,
	})
}

// Your non-periodic critical strikes cause your active Corruption, Immolate, Shadowflame, and Unstable Affliction spells on the target to immediately deal one pulse of their damage to the target.
func (warlock *Warlock) applyNaxxramasDamage4PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Damage 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	copiedSpellConfig := []struct {
		ClassMask   uint64
		SpellID     int32
		SpellSchool core.SpellSchool
		Flags       core.SpellFlag
	}{
		{
			ClassMask:   ClassSpellMask_WarlockCorruption,
			SpellID:     1219428,
			SpellSchool: core.SpellSchoolShadow,
			Flags:       WarlockFlagAffliction,
		},
		{
			ClassMask:   ClassSpellMask_WarlockImmolate,
			SpellID:     1219425,
			SpellSchool: core.SpellSchoolFire,
			Flags:       WarlockFlagDestruction,
		},
		{
			ClassMask:   ClassSpellMask_WarlockShadowflame,
			SpellID:     1219429,
			SpellSchool: core.SpellSchoolShadow | core.SpellSchoolFire,
			Flags:       WarlockFlagAffliction | WarlockFlagDestruction,
		},
		{
			ClassMask:   ClassSpellMask_WarlockUnstableAffliction,
			SpellID:     1219436,
			SpellSchool: core.SpellSchoolShadow,
			Flags:       WarlockFlagAffliction,
		},
	}

	dotSpellsMap := make(map[uint64]*core.Spell)

	for _, spellConfig := range copiedSpellConfig {
		dotSpellsMap[spellConfig.ClassMask] = warlock.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: spellConfig.SpellID},
			ClassSpellMask: spellConfig.ClassMask,
			SpellSchool:    spellConfig.SpellSchool,
			DefenseType:    core.DefenseTypeMagic,
			ProcMask:       core.ProcMaskSpellDamage,
			Flags:          core.SpellFlagTreatAsPeriodic | core.SpellFlagPureDot | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | WarlockFlagHaunt | spellConfig.Flags,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {},
		})
	}

	var affectedDotSpells []*core.Spell

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedDotSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{
					warlock.Corruption,
					warlock.Immolate,
					{warlock.Shadowflame, warlock.UnstableAffliction},
				}),
				func(spell *core.Spell) bool { return spell != nil },
			)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_WarlockAll) && result.DidCrit() {
				for _, spell := range affectedDotSpells {
					if dot := spell.Dot(result.Target); dot.IsActive() {
						copiedDoTSpell := dotSpellsMap[spell.ClassSpellMask]
						copiedDoTSpell.Cast(sim, result.Target)
						copiedDoTSpell.CalcAndDealDamage(sim, result.Target, dot.SnapshotBaseDamage, copiedDoTSpell.OutcomeMagicCrit)
					}
				}
			}
		},
	}))
}

// Your Curse of Agony does not expire on Undead targets, and continues to grow in power indefinitely.
func (warlock *Warlock) applyNaxxramasDamage6PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Damage 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	warlock.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range warlock.CurseOfAgony {
				oldApplyEffects := spell.ApplyEffects
				spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					if target.MobType == proto.MobType_MobTypeUndead {
						dot := spell.Dot(target)
						dot.NumberOfTicks = 65_536 // Large enough to be effectively infinite for our purposes
						dot.RecomputeAuraDuration()
					}
					oldApplyEffects(sim, target, spell)
				}
			}
		},
	})
}

var ItemSetPlagueheartStitchings = core.NewItemSet(core.ItemSet{
	Name: "Plagueheart Stitchings",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasTank2PBonus()
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasTank4PBonus()
		},
		6: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.applyNaxxramasTank6PBonus()
		},
	},
})

// Your Menace ability never misses, and your chance to be Dodged or Parried or for your spells to miss is reduced by 2%.
func (warlock *Warlock) applyNaxxramasTank2PBonus() {
	label := "S03 - Item - Naxxramas - Warlock - Tank 2P Bonus"
	if warlock.HasAura(label) {
		return
	}

	bonusStats := stats.Stats{stats.Expertise: 2 * core.ExpertiseRatingPerExpertiseChance}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:      label,
		BuildPhase: core.CharacterBuildPhaseBuffs,
	}).AttachStatsBuff(bonusStats))
}

// Reduces the cooldown on your Infernal Armor ability by 10 sec and reduces the cooldown on your Demonic Grace ability by 3 sec.
func (warlock *Warlock) applyNaxxramasTank4PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneCloakInfernalArmor) && !warlock.HasRune(proto.WarlockRune_RuneLegsDemonicGrace) {
		return
	}

	label := "S03 - Item - Naxxramas - Warlock - Tank 4P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_WarlockInfernalArmor,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 10,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: ClassSpellMask_WarlockDemonicGrace,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 3,
	})
}

// When an Undead enemy attempts to attack you, the remaining duration of your active Vengeance is reset to 20 sec.
func (warlock *Warlock) applyNaxxramasTank6PBonus() {
	if !warlock.HasRune(proto.WarlockRune_RuneHelmVengeance) {
		return
	}

	label := "S03 - Item - Naxxramas - Warlock - Tank 6P Bonus"
	if warlock.HasAura(label) {
		return
	}

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warlock.VengeanceAura.IsActive() && spell.Unit.MobType == proto.MobType_MobTypeUndead {
				warlock.VengeanceAura.Activate(sim)
			}
		},
	}))
}
