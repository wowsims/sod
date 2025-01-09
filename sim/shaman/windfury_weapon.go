package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const WindfuryWeaponRanks = 4

var WindfuryWeaponSpellId = [WindfuryWeaponRanks + 1]int32{0, 8232, 8235, 10486, 16362}
var WindfuryWeaponEnchantId = [WindfuryWeaponRanks + 1]int32{0, 283, 284, 525, 1669}
var WindfuryWeaponBonusAP = [WindfuryWeaponRanks + 1]float64{0, 104, 119, 249, 333}
var WindfuryWeaponLevel = [WindfuryWeaponRanks + 1]int32{0, 30, 40, 50, 60}

var WindfuryWeaponRankByLevel = map[int32]int32{
	25: 0,
	40: 2,
	50: 3,
	60: 4,
}

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	rank := WindfuryWeaponRankByLevel[shaman.Level]

	ewMultiplier := []float64{1, 1.13, 1.27, 1.4}[shaman.Talents.ElementalWeapons]
	bonusAP := WindfuryWeaponBonusAP[rank]

	actionID := core.ActionID{SpellID: WindfuryWeaponSpellId[rank]}.WithTag(core.TernaryInt32(isMH, 1, 2))
	procMask := core.ProcMaskMeleeMHSpecial
	damageMultiplier := 1.0
	weaponDamageFunc := shaman.MHWeaponDamage
	if !isMH {
		procMask = core.ProcMaskMeleeOHSpecial
		damageMultiplier = shaman.AutoAttacks.OHConfig().DamageMultiplier
		weaponDamageFunc = shaman.OHWeaponDamage
	}

	spellConfig := core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: damageMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mAP := spell.MeleeAttackPower() + (bonusAP+shaman.bonusWindfuryWeaponAP)*ewMultiplier*ewMultiplier
			baseDamage := weaponDamageFunc(sim, mAP)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	}

	return shaman.RegisterSpell(spellConfig)
}

func (shaman *Shaman) RegisterWindfuryImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	rank := WindfuryWeaponRankByLevel[shaman.Level]
	enchantId := WindfuryWeaponEnchantId[rank]

	icdDuration := time.Millisecond * 1500

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = enchantId
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = enchantId
	}

	var procChance = 0.2
	if procMask == core.ProcMaskMelee {
		procChance = 0.36
	}

	shaman.WindfuryWeaponMH = shaman.newWindfuryImbueSpell(true)
	shaman.WindfuryWeaponOH = shaman.newWindfuryImbueSpell(false)

	triggerAura := core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:              "Windfury Trigger",
		Callback:          core.CallbackOnSpellHitDealt,
		ProcMask:          procMask,
		Outcome:           core.OutcomeLanded,
		SpellFlagsExclude: core.SpellFlagSuppressEquipProcs,
		ProcChance:        procChance,
		ICD:               icdDuration,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.IsMH() {
				shaman.WindfuryWeaponMH.Cast(sim, result.Target)
				shaman.WindfuryWeaponMH.Cast(sim, result.Target)
			} else {
				shaman.WindfuryWeaponOH.Cast(sim, result.Target)
				shaman.WindfuryWeaponOH.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(enchantId, &procMask, triggerAura)
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

	rank := WindfuryWeaponRankByLevel[shaman.Level]
	enchantId := WindfuryWeaponEnchantId[rank]

	item.TempEnchant = enchantId
}
