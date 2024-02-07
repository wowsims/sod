package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const FrostbrandWeaponRanks = 5

var FrostbrandWeaponSpellId = [FrostbrandWeaponRanks + 1]int32{0, 8033, 8038, 10456, 16355, 16356}
var FrostbrandWeaponEnchantId = [FrostbrandWeaponRanks + 1]int32{0, 2, 12, 524, 1667, 1668}
var FrostbrandWeaponBaseDamage = [FrostbrandWeaponRanks + 1]float64{0, 46, 77, 94, 142, 187}
var FrostbrandWeaponLevel = [FrostbrandWeaponRanks + 1]int32{0, 20, 28, 38, 48, 58}

var FrostbrandWeaponRankByLevel = map[int32]int32{
	25: 1,
	40: 3,
	50: 4,
	60: 5,
}

func (shaman *Shaman) FrostbrandDebuffAura(target *core.Unit, level int32) *core.Aura {
	rank := FrostbrandWeaponRankByLevel[level]
	spellId := FrostbrandWeaponSpellId[rank]

	return target.GetOrRegisterAura(core.Aura{
		Label:    "Frostbrand Attack-" + shaman.Label,
		ActionID: core.ActionID{SpellID: spellId},
		Duration: time.Second * 8,
	})
}

func (shaman *Shaman) newFrostbrandImbueSpell() *core.Spell {
	level := shaman.GetCharacter().Level
	rank := FrostbrandWeaponRankByLevel[level]
	spellId := FrostbrandWeaponSpellId[rank]
	baseDamage := FrostbrandWeaponBaseDamage[rank]

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage + 0.1*spell.SpellDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) RegisterFrostbrandImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	level := shaman.GetCharacter().Level
	rank := FrostbrandWeaponRankByLevel[level]
	enchantId := FrostbrandWeaponEnchantId[rank]

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = enchantId
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = enchantId
	}

	ppmm := shaman.AutoAttacks.NewPPMManager(9.0, procMask)

	mhSpell := shaman.newFrostbrandImbueSpell()
	ohSpell := shaman.newFrostbrandImbueSpell()

	fbDebuffAuras := shaman.NewEnemyAuraArray(shaman.FrostbrandDebuffAura)

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Frostbrand Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if ppmm.Proc(sim, spell.ProcMask, "Frostbrand Weapon") {
				if spell.IsMH() {
					mhSpell.Cast(sim, result.Target)
				} else {
					ohSpell.Cast(sim, result.Target)
				}
				fbDebuffAuras.Get(result.Target).Activate(sim)
			}
		},
	})

	shaman.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3784, 9.0, &ppmm, aura)
}

func (shaman *Shaman) ApplyFrostbrandImbue(procMask core.ProcMask) {
	if procMask.Matches(core.ProcMaskMeleeMH) && shaman.HasMHWeapon() {
		shaman.ApplyFrostbrandImbueToItem(shaman.MainHand())
	}

	if procMask.Matches(core.ProcMaskMeleeOH) && shaman.HasOHWeapon() {
		shaman.ApplyFrostbrandImbueToItem(shaman.OffHand())
	}
}

func (shaman *Shaman) ApplyFrostbrandImbueToItem(item *core.Item) {
	if item == nil {
		return
	}

	level := shaman.GetCharacter().Level
	rank := FrostbrandWeaponRankByLevel[level]
	enchantId := FrostbrandWeaponEnchantId[rank]

	item.TempEnchant = enchantId
}
