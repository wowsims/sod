package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const WindfuryWeaponRanks = 4

var WindfuryWeaponSpellId = [WindfuryWeaponRanks + 1]int32{0, 8232, 8235, 10486, 16362}
var WindfuryWeaponEnchantId = [WindfuryWeaponRanks + 1]int32{0, 283, 284, 525, 1669}
var WindfuryWeaponBonusAP = [WindfuryWeaponRanks + 1]float64{0, 104, 222, 316, 333}
var WindfuryWeaponLevel = [WindfuryWeaponRanks + 1]int32{0, 30, 40, 50, 60}

var WindfuryWeaponRankByLevel = map[int32]int32{
	25: 0,
	40: 2,
	50: 3,
	60: 4,
}

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	level := shaman.GetCharacter().Level
	rank := WindfuryWeaponRankByLevel[level]
	SpellId := WindfuryWeaponSpellId[rank]
	bonusAP := WindfuryWeaponBonusAP[rank] * (1 + .1333*float64(shaman.Talents.ElementalWeapons))

	tag := 1
	procMask := core.ProcMaskMeleeMHSpecial
	weaponDamageFunc := shaman.MHWeaponDamage
	if !isMH {
		tag = 2
		procMask = core.ProcMaskMeleeOHSpecial
		weaponDamageFunc = shaman.OHWeaponDamage
		bonusAP *= 2 // applied after 50% offhand penalty
	}

	spellConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: SpellId, Tag: int32(tag)},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := spell.BonusWeaponDamage()
			mAP := spell.MeleeAttackPower() + bonusAP

			baseDamage1 := constBaseDamage + weaponDamageFunc(sim, mAP)
			baseDamage2 := constBaseDamage + weaponDamageFunc(sim, mAP)
			result1 := spell.CalcDamage(sim, target, baseDamage1, spell.OutcomeMeleeSpecialHitAndCrit)
			result2 := spell.CalcDamage(sim, target, baseDamage2, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DealDamage(sim, result1)
			spell.DealDamage(sim, result2)
		},
	}

	return shaman.RegisterSpell(spellConfig)
}

func (shaman *Shaman) RegisterWindfuryImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	level := shaman.GetCharacter().Level
	rank := WindfuryWeaponRankByLevel[level]
	enchantId := WindfuryWeaponEnchantId[rank]

	icdDuration := time.Millisecond * 1500
	buffDuration := time.Minute * 5

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = enchantId
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = enchantId
	}

	var proc = 0.2
	if procMask == core.ProcMaskMelee {
		proc = 0.36
	}

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: icdDuration,
	}

	mhSpell := shaman.newWindfuryImbueSpell(true)
	ohSpell := shaman.newWindfuryImbueSpell(false)

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Windfury Imbue",
		Duration: buffDuration,
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

			if sim.RandomFloat("Windfury Imbue") < proc {
				icd.Use(sim)

				if spell.IsMH() {
					mhSpell.Cast(sim, result.Target)
				} else {
					ohSpell.Cast(sim, result.Target)
				}
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(enchantId, &procMask, aura)
}

func (shaman *Shaman) ApplyWindfuryImbue(procMask core.ProcMask) {
	if procMask.Matches(core.ProcMaskMeleeMH) && shaman.HasMHWeapon() {
		shaman.ApplyWindfuryImbueToItem(shaman.MainHand())
	}

	if procMask.Matches(core.ProcMaskMeleeOH) && shaman.HasOHWeapon() {
		shaman.ApplyWindfuryImbueToItem(shaman.OffHand())
	}
}

func (shaman *Shaman) ApplyWindfuryImbueToItem(item *core.Item) {
	if item == nil {
		return
	}

	level := shaman.GetCharacter().Level
	rank := WindfuryWeaponRankByLevel[level]
	enchantId := WindfuryWeaponEnchantId[rank]

	item.TempEnchant = enchantId
}
