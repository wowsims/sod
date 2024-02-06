package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) ApplyRunes() {
	// Chest
	shaman.applyDualWieldSpec()
	// shaman.applyHealingRain()
	// shaman.applyOverload()
	shaman.applyShieldMastery()
	shaman.applyTwoHandedMastery()

	// Hands
	shaman.applyLavaBurst()
	shaman.applyLavaLash()
	shaman.applyMoltenBlast()

	// Waist
	// shaman.applyFireNova()
	shaman.applyMaelstromWeapon()
	shaman.applyPowerSurge()

	// Legs
	shaman.applyAncestralGuidance()
	// shaman.applyEarthShield()
	shaman.applyShamanisticRage()
	shaman.applyWayOfEarth()

	// Feet
	shaman.applyAncestralAwakening()
	// shaman.applyDecoyTotem()
	// shaman.applySpiritOfTheAlpha()
}

func (shaman *Shaman) applyDualWieldSpec() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestDualWieldSpec) {
		return
	}

	if shaman.HasMHWeapon() && shaman.HasOHWeapon() {
		shaman.AddStat(stats.MeleeHit, 10)
		shaman.AddStat(stats.SpellHit, 10)
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Dual Wield Specialization",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestDualWieldSpec)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyShieldMastery() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestShieldMastery) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Shield Mastery",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestShieldMastery)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyTwoHandedMastery() {
	if !shaman.HasRune(proto.ShamanRune_RuneTwoHandedMastery) {
		return
	}

	procAura := shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery Proc",
		ActionID: core.ActionID{SpellID: 436365},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyAttackSpeed(sim, 1/1.3)
		},
	})

	shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneTwoHandedMastery)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				return
			}

			if shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand {
				procAura.Activate(sim)
			} else {
				procAura.Deactivate(sim)
			}
		},
	})
}

func (shaman *Shaman) applyMaelstromWeapon() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	shaman.MaelstromWeaponAura = shaman.RegisterAura(core.Aura{
		Label:    "Maelstrom Weapon",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneWaistMaelstromWeapon)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

const ShamanPowerSurgeProcChance = .05

func (shaman *Shaman) applyPowerSurge() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) {
		return
	}

	intMP5Rate := .15

	shaman.AddStats(
		stats.Stats{
			stats.MP5: shaman.GetStat(stats.Intellect) * intMP5Rate,
		},
	)

	var affectedSpells []*core.Spell

	shaman.PowerSurgeAura = shaman.RegisterAura(core.Aura{
		Label:    "Power Surge Proc",
		ActionID: core.ActionID{SpellID: 440285},
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{
					shaman.ChainLightning,
					shaman.ChainHeal,
					{shaman.LavaBurst},
				}), func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.LavaBurst.CD.Reset()
			for _, spell := range affectedSpells {
				spell.CD.Reset()
			}
			core.Each(affectedSpells, func(spell *core.Spell) { spell.CastTimeMultiplier -= 1 })
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) { spell.CastTimeMultiplier += 1 })
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode != int32(SpellCode_ShamanChainLightning) && spell.SpellCode != int32(SpellCode_ChainHeal) && spell != shaman.LavaBurst {
				return
			}
			aura.Deactivate(sim)
		},
	})
}

func (shaman *Shaman) applyWayOfEarth() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsWayOfEarth) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Way of Earth",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsWayOfEarth)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}
