package core

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type DebuffName int32

const (
	// General Buffs
	DemoralizingShout DebuffName = iota
)

var LevelToDebuffRank = map[DebuffName]map[int32]int32{
	DemoralizingShout: {
		25: 2,
		40: 3,
		50: 4,
		60: 5,
	},
}

func applyDebuffEffects(target *Unit, targetIdx int, debuffs *proto.Debuffs, raid *proto.Raid) {
	level := raid.Parties[0].Players[0].Level
	if debuffs.JudgementOfWisdom && targetIdx == 0 {
		jowAura := JudgementOfWisdomAura(target, level)
		if jowAura != nil {
			MakePermanent(jowAura)
		}
	}

	if debuffs.ImprovedShadowBolt {
		//TODO: Apply periodically
		MakePermanent(ImprovedShadowBoltAura(target, 5))
	}

	if debuffs.ShadowWeaving {
		aura := ShadowWeavingAura(target, 5)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        5,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio so it comes before actual warrior sunders.
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		}, raid)
	}

	if debuffs.MekkatorqueFistDebuff {
		MakePermanent(MekkatorqueFistDebuffAura(target, level))
	}

	if debuffs.CurseOfElements {
		MakePermanent(CurseOfElementsAura(target, level))
	}

	if debuffs.CurseOfShadow {
		MakePermanent(CurseOfShadowAura(target, level))
	}

	if debuffs.ImprovedScorch && targetIdx == 0 {
		aura := ImprovedScorchAura(target)
		ScheduledMajorArmorAura(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        5,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio so it comes before actual mage scorches
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		}, raid)
	}

	if debuffs.WintersChill && targetIdx == 0 {
		MakePermanent(WintersChillAura(target, 5))
	}

	if debuffs.Stormstrike {
		MakePermanent(StormstrikeAura(target))
	} else if debuffs.Dreamstate {
		MakePermanent(DreamstateAura(target))
	}

	if debuffs.GiftOfArthas {
		MakePermanent(GiftOfArthasAura(target))
	}

	if debuffs.CurseOfVulnerability {
		MakePermanent(CurseOfVulnerabilityAura(target))
	}

	if debuffs.Mangle {
		MakePermanent(MangleAura(target, level))
	}

	if debuffs.CrystalYield {
		MakePermanent(CrystalYieldAura(target))
	}

	if debuffs.AncientCorrosivePoison > 0 {
		ApplyFixedUptimeAura(AncientCorrosivePoisonAura(target), float64(debuffs.AncientCorrosivePoison)/100.0, GCDDefault, 1)
	}

	// Major Armor Debuffs
	if targetIdx == 0 {
		if debuffs.ExposeArmor != proto.TristateEffect_TristateEffectMissing {
			// Improved EA
			aura := ExposeArmorAura(target, TernaryInt32(debuffs.ExposeArmor == proto.TristateEffect_TristateEffectRegular, 0, 2), level)
			ScheduledMajorArmorAura(aura, PeriodicActionOptions{
				Period:   time.Second * 3,
				NumTicks: 1,
				OnAction: func(sim *Simulation) {
					aura.Activate(sim)
				},
			}, raid)
		}

		if debuffs.SunderArmor {
			// Sunder Armor
			aura := SunderArmorAura(target, level)
			ScheduledMajorArmorAura(aura, PeriodicActionOptions{
				Period:          time.Millisecond * 1500,
				NumTicks:        5,
				TickImmediately: true,
				Priority:        ActionPriorityDOT, // High prio so it comes before actual warrior sunders.
				OnAction: func(sim *Simulation) {
					aura.Activate(sim)
					if aura.IsActive() {
						aura.AddStack(sim)
					}
				},
			}, raid)
		}

		if debuffs.Homunculi > 0 {
			// Calculate desired downtime based on selected uptimeCount (1 count = 10% uptime, 0%-100%)
			totalDuration := time.Second * 15
			uptimePercent := float64(debuffs.Homunculi) / 100.0
			ApplyFixedUptimeAura(HomunculiArmorAura(target, level), uptimePercent, totalDuration, 1)
		}
	}

	if debuffs.CurseOfRecklessness {
		MakePermanent(CurseOfRecklessnessAura(target, level))
	}

	if debuffs.FaerieFire {
		MakePermanent(FaerieFireAura(target, level))
	}

	if debuffs.CurseOfWeakness != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(CurseOfWeaknessAura(target, GetTristateValueInt32(debuffs.CurseOfWeakness, 1, 2), level))
	}

	if debuffs.DemoralizingRoar != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingRoarAura(target, GetTristateValueInt32(debuffs.DemoralizingRoar, 0, 5), level))
	}
	if debuffs.DemoralizingShout != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingShoutAura(target, 0, GetTristateValueInt32(debuffs.DemoralizingShout, 0, 5), level))
	}
	if debuffs.HuntersMark != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(HuntersMarkAura(target, GetTristateValueInt32(debuffs.HuntersMark, 0, 5), level))
	}

	// Atk spd reduction
	if debuffs.ThunderClap != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(ThunderClapAura(target, GetTristateValueInt32(debuffs.ThunderClap, 0, 3), level))
	}

	// Miss
	if debuffs.InsectSwarm && targetIdx == 0 {
		MakePermanent(InsectSwarmAura(target))
	}
	if debuffs.ScorpidSting && targetIdx == 0 {
		MakePermanent(ScorpidStingAura(target))
	}
}

func StormstrikeAura(unit *Unit) *Aura {
	return exclusiveNatureDamageTakenAura(unit, "Stormstrike", ActionID{SpellID: 17364})
}

func DreamstateAura(unit *Unit) *Aura {
	return exclusiveNatureDamageTakenAura(unit, "Dreamstate", ActionID{SpellID: 408258})
}

func exclusiveNatureDamageTakenAura(unit *Unit, label string, actionID ActionID) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: time.Second * 12,
	})

	aura.NewExclusiveEffect("NatureDamageTaken", true, ExclusiveEffect{
		Priority: 20,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1.2
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= 1.2
		},
	})

	return aura
}

func ImprovedShadowBoltAura(unit *Unit, rank int32) *Aura {
	damageMulti := 1. + 0.04*float64(rank)
	return unit.GetOrRegisterAura(Aura{
		Label:     "Improved Shadow Bolt",
		ActionID:  ActionID{SpellID: 17800},
		Duration:  12 * time.Second,
		MaxStacks: 4,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= damageMulti
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= damageMulti
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.SpellSchool.Matches(SpellSchoolShadow) && result.Landed() {
				aura.RemoveStack(sim)
			}
		},
	})
}

func ShadowWeavingAura(unit *Unit, rank int) *Aura {
	spellId := [6]int32{0, 15257, 15331, 15332, 15333, 15334}[rank]
	return unit.GetOrRegisterAura(Aura{
		Label:     "Shadow Weaving",
		ActionID:  ActionID{SpellID: spellId},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= 1.0 + 0.03*float64(oldStacks)
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1.0 + 0.03*float64(newStacks)
		},
	})
}

func ScheduledMajorArmorAura(aura *Aura, options PeriodicActionOptions, raid *proto.Raid) {
	aura.OnReset = func(aura *Aura, sim *Simulation) {
		aura.Duration = NeverExpires
		StartPeriodicAction(sim, options)
	}
}

var JudgementOfWisdomAuraLabel = "Judgement of Wisdom"

// TODO: Classic verify logic
func JudgementOfWisdomAura(target *Unit, level int32) *Aura {
	spellID := map[int32]int32{
		40: 20186,
		50: 20354,
		60: 20355,
	}[level]
	actionID := ActionID{SpellID: spellID}

	jowMana := 0.0
	if level < 38 {
		return nil
	} else if level < 48 {
		jowMana = 33.0
	} else if level < 58 {
		jowMana = 46.0
	} else {
		jowMana = 59.0
	}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfWisdomAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 10,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			unit := spell.Unit
			if !unit.HasManaBar() {
				return
			}

			if spell.ProcMask.Matches(ProcMaskEmpty | ProcMaskProc | ProcMaskWeaponProc) {
				return // Phantom spells (Romulo's, Lightning Capacitor, etc.) don't proc JoW.
			}

			procChance := 0.5
			if spell.ProcMask.Matches(ProcMaskWhiteHit | ProcMaskRanged) {
				// Apparently ranged/melee can still proc on miss
				if sim.RandomFloat("JoW Proc") > procChance {
					return
				}
			} else { // spell casting
				if !spell.ProcMask.Matches(ProcMaskDirect) {
					return
				}

				if !result.Landed() {
					return
				}

				if sim.RandomFloat("jow") > procChance {
					return
				}
			}

			if unit.JowManaMetrics == nil {
				unit.JowManaMetrics = unit.NewManaMetrics(actionID)
			}
			// JoW returns flat mana
			unit.AddMana(sim, jowMana, unit.JowManaMetrics)
		},
	})
}

var JudgementOfLightAuraLabel = "Judgement of Light"

func JudgementOfLightAura(target *Unit) *Aura {
	actionID := ActionID{SpellID: 20271}

	return target.GetOrRegisterAura(Aura{
		Label:    JudgementOfLightAuraLabel,
		ActionID: actionID,
		Duration: time.Second * 20,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !spell.ProcMask.Matches(ProcMaskMelee) || !result.Landed() {
				return
			}
		},
	})
}

func MekkatorqueFistDebuffAura(target *Unit, playerLevel int32) *Aura {
	if playerLevel < 40 {
		return nil
	}

	spellID := 434841
	resistance := 45.0
	dmgMod := 1.06

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Mekkatorque Debuff",
		ActionID: ActionID{SpellID: int32(spellID)},
		Duration: time.Second * 20,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.FireResistance:   -resistance,
				stats.FrostResistance:  -resistance,
				stats.ArcaneResistance: -resistance,
				stats.NatureResistance: -resistance,
				stats.ShadowResistance: -resistance,
			})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.FireResistance:   resistance,
				stats.FrostResistance:  resistance,
				stats.ArcaneResistance: resistance,
				stats.NatureResistance: resistance,
				stats.ShadowResistance: resistance,
			})
		},
	})

	// 0.01 priority as this overwrites the other spells of this category and does not allow them to be recast
	spellSchoolDamageEffect(aura, stats.SchoolIndexFire, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexFrost, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexArcane, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexNature, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexShadow, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexHoly, dmgMod, 0.01, true)
	return aura
}

func CurseOfElementsAura(target *Unit, playerLevel int32) *Aura {
	if playerLevel < 40 {
		return nil
	}

	spellID := map[int32]int32{
		40: 1490,
		50: 11721,
		60: 11722,
	}[playerLevel]

	resistance := map[int32]float64{
		40: 45,
		50: 60,
		60: 75,
	}[playerLevel]

	dmgMod := map[int32]float64{
		40: 1.06,
		50: 1.08,
		60: 1.10,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Elements",
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Minute * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.FireResistance: -resistance, stats.FrostResistance: -resistance})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.FireResistance: resistance, stats.FrostResistance: resistance})
		},
	})
	spellSchoolDamageEffect(aura, stats.SchoolIndexFire, dmgMod, 0.0, false)
	spellSchoolDamageEffect(aura, stats.SchoolIndexFrost, dmgMod, 0.0, false)
	return aura
}

func CurseOfShadowAura(target *Unit, playerLevel int32) *Aura {
	if playerLevel < 50 {
		return nil
	}

	spellID := map[int32]int32{
		50: 17862,
		60: 17937,
	}[playerLevel]

	resistance := map[int32]float64{
		50: 60,
		60: 75,
	}[playerLevel]

	dmgMod := map[int32]float64{
		50: 1.08,
		60: 1.10,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Shadow",
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Minute * 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: -resistance, stats.ShadowResistance: -resistance})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.ArcaneResistance: resistance, stats.ShadowResistance: resistance})
		},
	})
	spellSchoolDamageEffect(aura, stats.SchoolIndexArcane, dmgMod, 0.0, false)
	spellSchoolDamageEffect(aura, stats.SchoolIndexShadow, dmgMod, 0.0, false)
	return aura
}

func spellSchoolDamageEffect(aura *Aura, school stats.SchoolIndex, multiplier float64, extraPriority float64, exclusive bool) *ExclusiveEffect {
	return aura.NewExclusiveEffect("spellDamage"+strconv.Itoa(int(school)), exclusive, ExclusiveEffect{
		Priority: multiplier + extraPriority,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[school] *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[school] /= multiplier
		},
	})
}

func GiftOfArthasAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Gift of Arthas",
		ActionID: ActionID{SpellID: 11374},
		Duration: time.Minute * 3,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken += 8
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken -= 8
		},
	})
}

func HemorrhageAura(target *Unit, casterLevel int32) *Aura {
	debuffBonusDamage := map[int32]float64{
		40: 3,
		50: 5,
		60: 7,
	}[casterLevel]

	spellID := map[int32]int32{
		40: 16511,
		50: 17347,
		60: 17348,
	}[casterLevel]

	return target.GetOrRegisterAura(Aura{
		Label:     "Hemorrhage",
		ActionID:  ActionID{SpellID: spellID},
		Duration:  time.Second * 8,
		MaxStacks: 30,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken += debuffBonusDamage
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken -= debuffBonusDamage
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.SpellSchool != SpellSchoolPhysical {
				return
			}
			if !result.Landed() || result.Damage == 0 {
				return
			}
			// TODO find out which abilities are actually affected
			aura.RemoveStack(sim)
		},
	})
}

func CurseOfVulnerabilityAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Curse of Vulnerability",
		ActionID: ActionID{SpellID: 427143},
		Duration: time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken += 2
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken -= 2
		},
	})
}

func MangleAura(target *Unit, playerLevel int32) *Aura {
	return bleedDamageAura(target, Aura{
		Label:    "Mangle",
		ActionID: ActionID{SpellID: 409828},
		Duration: time.Minute,
	}, 1.3)
}

// Bleed Damage Multiplier category
const BleedEffectCategory = "BleedDamage"

func bleedDamageAura(target *Unit, config Aura, multiplier float64) *Aura {
	aura := target.GetOrRegisterAura(config)
	aura.NewExclusiveEffect(BleedEffectCategory, true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.PeriodicPhysicalDamageTakenMultiplier /= multiplier
		},
	})
	return aura
}

const SpellFirePowerEffectCategory = "spellFirePowerdebuff"

func ImprovedScorchAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Improved Scorch",
		ActionID:  ActionID{SpellID: 12873},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= 1 + .03*float64(oldStacks)
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1 + .03*float64(newStacks)
		},
	})

	return aura
}

const SpellCritEffectCategory = "spellcritdebuff"

func WintersChillAura(target *Unit, startingStacks int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Winter's Chill",
		ActionID:  ActionID{SpellID: 28593},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, startingStacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks, newStacks int32) {
			aura.Unit.PseudoStats.SchoolCritTakenMultiplier[stats.SchoolIndexFrost] /= 1 + 0.2*float64(oldStacks)
			aura.Unit.PseudoStats.SchoolCritTakenMultiplier[stats.SchoolIndexFrost] *= 1 + 0.2*float64(newStacks)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolCritTakenMultiplier[stats.SchoolIndexFrost] /= 1 + 0.2*float64(aura.stacks)
		},
	})

	// effect = aura.NewExclusiveEffect(SpellCritEffectCategory, true, ExclusiveEffect{
	// 	Priority: 0,
	// 	OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
	// 		ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken += ee.Priority * CritRatingPerCritChance
	// 	},
	// 	OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
	// 		ee.Aura.Unit.PseudoStats.BonusSpellCritRatingTaken -= ee.Priority * CritRatingPerCritChance
	// 	},
	// })
	return aura
}

var majorArmorReductionEffectCategory = "MajorArmorReduction"

func SunderArmorAura(target *Unit, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 7405,
		40: 8380,
		50: 11596,
		60: 11597,
	}[playerLevel]

	arpen := map[int32]float64{
		25: 180,
		40: 270,
		50: 360,
		60: 450,
	}[playerLevel]

	var effect *ExclusiveEffect
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Sunder Armor",
		ActionID:  ActionID{SpellID: spellID},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			effect.SetPriority(sim, arpen*float64(newStacks))
		},
	})

	effect = aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: 0,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, -ee.Priority)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, ee.Priority)
		},
	})

	return aura
}

func ExposeArmorAura(target *Unit, improvedEA int32, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 8647,
		40: 8650,
		50: 11197,
		60: 11198,
	}[playerLevel]

	arpen := map[int32]float64{
		25: 400,
		40: 1050,
		50: 1375,
		60: 1700,
	}[playerLevel]

	arpen *= 1 + 0.25*float64(improvedEA)

	aura := target.GetOrRegisterAura(Aura{
		Label:    "ExposeArmor",
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Second * 30,
	})

	aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: arpen,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, -ee.Priority)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, ee.Priority)
		},
	})

	return aura
}

func HomunculiAttackSpeedAura(target *Unit, playerLevel int32) *Aura {
	multiplier := 1.1

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Cripple (Homunculus)",
		ActionID: ActionID{SpellID: 402808},
		Duration: time.Second * 15,
	})

	aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/multiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, multiplier)
		},
	})

	return aura
}

func HomunculiArmorAura(target *Unit, playerLevel int32) *Aura {
	arpen := float64(185 + 35*(playerLevel-1))

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Degrade (Homunculus)",
		ActionID: ActionID{SpellID: 402818},
		Duration: time.Second * 15,
	})

	aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: arpen,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, -ee.Priority)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, ee.Priority)
		},
	})

	return aura
}

func HomunculiAttackPowerAura(target *Unit, playerLevel int32) *Aura {
	ap := float64(190 + 3*(playerLevel-1))

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Demoralize (Homunculus)",
		ActionID: ActionID{SpellID: 402811},
		Duration: time.Second * 15,
	})

	aura.NewExclusiveEffect("Homonculi AP", true, ExclusiveEffect{
		Priority: ap,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			target.AddStatDynamic(sim, stats.AttackPower, -1*ap)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			target.AddStatDynamic(sim, stats.AttackPower, ap)
		},
	})

	return aura
}

func CurseOfRecklessnessAura(target *Unit, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 704,
		40: 7658,
		50: 7659,
		60: 11717,
	}[playerLevel]

	arpen := map[int32]float64{
		25: 140,
		40: 290,
		50: 465,
		60: 640,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Recklessness",
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Minute * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, -arpen)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, arpen)
		},
	})
	return aura
}

func FaerieFireAura(target *Unit, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 770,
		40: 778,
		50: 9749,
		60: 9907,
	}[playerLevel]

	arpen := map[int32]float64{
		25: 175,
		40: 285,
		50: 395,
		60: 505,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Faerie Fire",
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Second * 40,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, -arpen)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, arpen)
		},
	})
	return aura
}

// TODO: Classic
func CurseOfWeaknessAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Weakness" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 50511},
		Duration: time.Minute * 2,
	})
	return aura
}

const HuntersMarkAuraTag = "HuntersMark"

func HuntersMarkAura(target *Unit, points int32, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 14323,
		40: 14324,
		50: 14324,
		60: 14325,
	}[playerLevel]

	bonus := map[int32]float64{
		25: 45,
		40: 75,
		50: 75,
		60: 110,
	}[playerLevel]

	bonus *= 1 + 0.03*float64(points)

	aura := target.GetOrRegisterAura(Aura{
		Label:    "HuntersMark-" + strconv.Itoa(int(bonus)),
		Tag:      HuntersMarkAuraTag,
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Minute * 2,
	})

	aura.NewExclusiveEffect("HuntersMark", true, ExclusiveEffect{
		Priority: bonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusRangedAttackPowerTaken += bonus
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BonusRangedAttackPowerTaken -= bonus
		},
	})

	return aura
}

// TODO: Classic
func DemoralizingRoarAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingRoar-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 9898},
		Duration: time.Second * 30,
	})
	apReductionEffect(aura, 411*(1+0.08*float64(points)))
	return aura
}

const DemoralizingShoutRanks = 5

var DemoralizingShoutSpellId = [DemoralizingShoutRanks + 1]int32{0, 1160, 6190, 11554, 11555, 11556}
var DemoralizingShoutBaseAP = [DemoralizingShoutRanks + 1]float64{0, 45, 56, 76, 111, 146}
var DemoralizingShoutLevel = [DemoralizingShoutRanks + 1]int{0, 14, 24, 34, 44, 54}

func DemoralizingShoutAura(target *Unit, boomingVoicePts int32, impDemoShoutPts int32, playerLevel int32) *Aura {
	rank := LevelToDebuffRank[DemoralizingShout][target.Level]
	spellId := DemoralizingShoutSpellId[rank]
	baseAPReduction := DemoralizingShoutBaseAP[rank]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingShout-" + strconv.Itoa(int(impDemoShoutPts)),
		ActionID: ActionID{SpellID: spellId},
		Duration: time.Duration(float64(time.Second*30) * (1 + 0.1*float64(boomingVoicePts))),
	})
	apReductionEffect(aura, baseAPReduction*(1+0.08*float64(impDemoShoutPts)))
	return aura
}

// TODO: Classic
func VindicationAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Vindication",
		ActionID: ActionID{SpellID: 26016},
		Duration: time.Second * 10,
	})
	apReductionEffect(aura, 287*float64(points))
	return aura
}

func apReductionEffect(aura *Aura, apReduction float64) *ExclusiveEffect {
	statReduction := stats.Stats{stats.AttackPower: -apReduction}
	return aura.NewExclusiveEffect("APReduction", false, ExclusiveEffect{
		Priority: apReduction,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, statReduction)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, statReduction.Invert())
		},
	})
}

// TODO: Classic
func ThunderClapAura(target *Unit, points int32, playerLevel int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "ThunderClap-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: 47502},
		Duration: time.Second * 30,
	})
	AtkSpeedReductionEffect(aura, []float64{1.1, 1.14, 1.17, 1.2}[points])
	return aura
}

func WaylayAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Waylay",
		ActionID: ActionID{SpellID: 408699},
		Duration: time.Second * 8,
	})
	AtkSpeedReductionEffect(aura, 1.1)
	return aura
}

func AtkSpeedReductionEffect(aura *Aura, speedMultiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("AtkSpdReduction", false, ExclusiveEffect{
		Priority: speedMultiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, 1/speedMultiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(sim, speedMultiplier)
		},
	})
}

func InsectSwarmAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "InsectSwarmMiss",
		ActionID: ActionID{SpellID: 24977},
		Duration: time.Second * 12,
	})
	increasedMissEffect(aura, 0.02)
	return aura
}

func ScorpidStingAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Scorpid Sting",
		ActionID: ActionID{SpellID: 3043},
		Duration: time.Second * 20,
	})
	return aura
}

func increasedMissEffect(aura *Aura, increasedMissChance float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("IncreasedMiss", false, ExclusiveEffect{
		Priority: increasedMissChance,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.IncreasedMissChance += increasedMissChance
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.IncreasedMissChance -= increasedMissChance
		},
	})
}

func CrystalYieldAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Crystal Yield",
		ActionID: ActionID{SpellID: 15235},
		Duration: 2 * time.Minute,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] -= 200
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] += 200
		},
	})
}

func AncientCorrosivePoisonAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Ancient Corrosive Poison",
		ActionID: ActionID{SpellID: 422996},
		Duration: 15 * time.Second,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] -= 150
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] += 150
		},
	})
}
