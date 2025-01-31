package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const RockbiterWeaponRanks = 7

var RockbiterWeaponEnchantId = [RockbiterWeaponRanks + 1]int32{0, 29, 6, 1, 503, 1663, 683, 1664}
var RockbiterWeaponBonusAP = [RockbiterWeaponRanks + 1]float64{0, 50, 79, 118, 138, 319, 490, 653}
var RockbiterWeaponBonusTPS = [RockbiterWeaponRanks + 1]float64{0, 6, 10, 16, 27, 41, 55, 72}
var RockbiterWeaponLevel = [RockbiterWeaponRanks + 1]int32{0, 1, 8, 16, 24, 34, 44, 54}

var RockbiterWeaponRankByLevel = map[int32]int32{
	25: 4,
	40: 5,
	50: 6,
	60: 7,
}

func (shaman *Shaman) RegisterRockbiterImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	rank := RockbiterWeaponRankByLevel[shaman.Level]
	enchantId := RockbiterWeaponEnchantId[rank]
	bonusThreat := RockbiterWeaponBonusTPS[rank]

	hasMHImbue := procMask.Matches(core.ProcMaskMeleeMH)
	hasOHImbue := procMask.Matches(core.ProcMaskMeleeOH)

	// Nerfed by 90% going into SoD Phase 3, in a... weird way ;)
	bonusAP := RockbiterWeaponBonusAP[rank] * ([]float64{1, 1.07, 1.14, 1.2}[shaman.Talents.ElementalWeapons] - 0.9)

	if hasMHImbue {
		shaman.MainHand().TempEnchant = enchantId
	}
	if hasOHImbue {
		shaman.OffHand().TempEnchant = enchantId
	}

	threatMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusThreat_Flat,
		ProcMask:   procMask,
		FloatValue: bonusThreat,
	})

	aura := shaman.NewDynamicEquipEffectAura(core.DynamicEquipEffectConfig{
		Label:    "Rockbiter Imbue",
		EffectID: enchantId,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			shaman.AddStatDynamic(sim, stats.AttackPower, bonusAP*float64(newStacks-oldStacks))
			threatMod.UpdateFloatValue(bonusThreat * float64(newStacks))
		},
	})

	aura.Duration = time.Minute * 5
	aura.BuildPhase = core.CharacterBuildPhaseGear
	aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		threatMod.Activate()

		if procMask.Matches(core.ProcMaskMeleeMH) {
			shaman.AutoAttacks.MHAuto().FlatThreatBonus += bonusThreat * shaman.AutoAttacks.MH().SwingSpeed
		}
		if procMask.Matches(core.ProcMaskMeleeOH) {
			shaman.AutoAttacks.OHAuto().FlatThreatBonus += bonusThreat * shaman.AutoAttacks.OH().SwingSpeed
		}
	})
	aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		threatMod.Deactivate()

		if procMask.Matches(core.ProcMaskMeleeMH) {
			shaman.AutoAttacks.MHAuto().FlatThreatBonus -= bonusThreat * shaman.AutoAttacks.MH().SwingSpeed
		}
		if procMask.Matches(core.ProcMaskMeleeOH) {
			shaman.AutoAttacks.OHAuto().FlatThreatBonus -= bonusThreat * shaman.AutoAttacks.OH().SwingSpeed
		}
	})

	shaman.RegisterOnItemSwapWithImbue(enchantId, &procMask, aura)
}
