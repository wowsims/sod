package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const FlametongueWeaponRanks = 6

var FlametongueWeaponSpellId = [FlametongueWeaponRanks + 1]int32{0, 8024, 8027, 8030, 16339, 16341, 16342}
var FlametongueWeaponEnchantId = [FlametongueWeaponRanks + 1]int32{0, 5, 4, 3, 523, 1665, 1666}
var FlametongueWeaponMaxDamage = [FlametongueWeaponRanks + 1]float64{0, 18, 26, 42, 57, 85, 112}
var FlametongueWeaponLevel = [FlametongueWeaponRanks + 1]int32{0, 10, 18, 26, 36, 46, 56}

var FlametongueWeaponRankByLevel = map[int32]int32{
	25: 2,
	40: 4,
	50: 5,
	60: 6,
}

func (shaman *Shaman) newFlametongueImbueSpell(weapon *core.Item) *core.Spell {
	level := shaman.GetCharacter().Level
	rank := FlametongueWeaponRankByLevel[level]
	spellID := FlametongueWeaponSpellId[rank]
	maxDamage := FlametongueWeaponMaxDamage[rank]

	baseDamage := maxDamage / 4
	spellCoeff := .1

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(spellID)},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if weapon.SwingSpeed != 0 {
				damage := weapon.SwingSpeed * (baseDamage + spellCoeff*spell.SpellDamage())
				damage *= 1 + .05*float64(shaman.Talents.ElementalWeapons)
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (shaman *Shaman) ApplyFlametongueImbueToItem(item *core.Item) {
	if item == nil {
		return
	}

	level := shaman.GetCharacter().Level
	rank := FlametongueWeaponRankByLevel[level]
	enchantId := FlametongueWeaponEnchantId[rank]

	item.TempEnchant = enchantId
}

func (shaman *Shaman) ApplyFlametongueImbue(procMask core.ProcMask) {
	if procMask.Matches(core.ProcMaskMeleeMH) && shaman.HasMHWeapon() {
		shaman.ApplyFlametongueImbueToItem(shaman.MainHand())
	}

	if procMask.Matches(core.ProcMaskMeleeOH) && shaman.HasOHWeapon() {
		shaman.ApplyFlametongueImbueToItem(shaman.OffHand())
	}
}

func (shaman *Shaman) RegisterFlametongueImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown && !shaman.ItemSwap.IsEnabled() {
		return
	}

	level := shaman.GetCharacter().Level
	rank := FlametongueWeaponRankByLevel[level]
	enchantId := FlametongueWeaponEnchantId[rank]

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueImbueSpell(shaman.MainHand())
	ohSpell := shaman.newFlametongueImbueSpell(shaman.OffHand())

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)

			if spell.IsMH() {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(enchantId, &procMask, aura)
}
