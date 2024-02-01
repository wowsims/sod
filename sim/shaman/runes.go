package shaman

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) ApplyRunes() {
	// Chest
	shaman.applyDualWieldSpec()
	// shaman.applyHealingRain()
	shaman.applyOverload()
	shaman.applyShieldMastery()
	shaman.applyTwoHandedMastery()

	// Hands
	shaman.applyLavaBurst()
	shaman.applyLavaLash()
	shaman.applyMoltenBlast()
	// shaman.applyWaterShield()

	// Waist
	// shaman.applyFireNova()
	shaman.applyMaelstromWeapon()
	shaman.applyPowerSurge()

	// Legs
	// shaman.applyAncestralGuidance()
	// shaman.applyEarthShield()
	// shaman.applyShamanisticRage()
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

	shaman.DualWieldSpecAura = shaman.RegisterAura(core.Aura{
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

	shaman.ShieldMasteryAura = shaman.RegisterAura(core.Aura{
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

	shaman.ShieldMasteryAura = shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneTwoHandedMastery)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyMaelstromWeapon() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	shaman.ShieldMasteryAura = shaman.RegisterAura(core.Aura{
		Label:    "Maelstrom Weapon",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneWaistMaelstromWeapon)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyPowerSurge() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) {
		return
	}

	shaman.ShieldMasteryAura = shaman.RegisterAura(core.Aura{
		Label:    "Power Surge",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneWaistPowerSurge)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyWayOfEarth() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsWayOfEarth) {
		return
	}

	shaman.ShieldMasteryAura = shaman.RegisterAura(core.Aura{
		Label:    "Way of Earth",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsWayOfEarth)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyAncestralAwakening() {
	if !shaman.HasRune(proto.ShamanRune_RuneFeetAncestralAwakening) {
		return
	}

	shaman.ShieldMasteryAura = shaman.RegisterAura(core.Aura{
		Label:    "Ancestral Awakening",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneFeetAncestralAwakening)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}
