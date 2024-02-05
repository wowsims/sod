package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

const RockbiterWeaponRanks = 7

var RockbiterWeaponEnchantId = [RockbiterWeaponRanks + 1]int32{0, 29, 6, 1, 503, 1663, 683, 1664}
var RockbiterWeaponBonusAP = [RockbiterWeaponRanks + 1]float64{0, 50, 79, 118, 194, 355, 522, 653}
var RockbiterWeaponLevel = [RockbiterWeaponRanks + 1]int32{0, 1, 8, 16, 24, 34, 44, 54}

var RockbiterRankByLevel = map[int32]int32{
	25: 4,
	40: 5,
	50: 6,
	60: 7,
}

func (shaman *Shaman) RegisterRockbiterImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	level := shaman.GetCharacter().Level
	rank := RockbiterRankByLevel[level]
	enchantId := RockbiterWeaponEnchantId[rank]
	bonusAP := RockbiterWeaponBonusAP[rank] * (1 + .07*float64(shaman.Talents.ElementalWeapons))

	duration := time.Minute * 5

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = enchantId
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = enchantId
	}

	// TODO: Rockbiter +threat

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Rockbiter Imbue",
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower: bonusAP,
			})
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower: bonusAP,
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower: -1 * bonusAP,
			})
		},
	})

	shaman.RegisterOnItemSwapWithImbue(enchantId, &procMask, aura)
}

func (shaman *Shaman) ApplyRockbiterImbue(procMask core.ProcMask) {
	if procMask.Matches(core.ProcMaskMeleeMH) && shaman.HasMHWeapon() {
		shaman.ApplyRockbiterImbueToItem(shaman.MainHand())
	}

	if procMask.Matches(core.ProcMaskMeleeOH) && shaman.HasOHWeapon() {
		shaman.ApplyRockbiterImbueToItem(shaman.OffHand())
	}
}

func (shaman *Shaman) ApplyRockbiterImbueToItem(item *core.Item) {
	if item == nil {
		return
	}

	level := shaman.GetCharacter().Level
	rank := RockbiterRankByLevel[level]
	enchantId := RockbiterWeaponEnchantId[rank]
	bonusAP := RockbiterWeaponBonusAP[rank]

	newStats := stats.Stats{stats.AttackPower: bonusAP}

	item.Stats = item.Stats.Add(newStats)
	item.TempEnchant = enchantId
}
