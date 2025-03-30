package core

import (
	"fmt"
	"math"
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

const (
	BloodPlagueAuraID = 1219121
	FrostFeverAuraID  = 1219124
)

var LevelToDebuffRank = map[DebuffName]map[int32]int32{
	DemoralizingShout: {
		25: 2,
		40: 3,
		50: 4,
		60: 5,
	},
}

func applyDebuffEffects(target *Unit, targetIdx int, debuffs *proto.Debuffs, level int32, units []*Unit) {
	if debuffs.JudgementOfWisdom && targetIdx == 0 {
		jowAura := JudgementOfWisdomAura(target, level)
		if jowAura != nil {
			MakePermanent(jowAura)
		}
	}

	if debuffs.JudgementOfLight && targetIdx == 0 {
		jolAura := JudgementOfLightAura(target, level, units)
		if jolAura != nil {
			MakePermanent(jolAura)
		}
	}

	if targetIdx == 0 {
		if debuffs.JudgementOfTheCrusader == proto.TristateEffect_TristateEffectRegular {
			MakePermanent(JudgementOfTheCrusaderAura(nil, target, level, 1, 0))
		} else if debuffs.JudgementOfTheCrusader == proto.TristateEffect_TristateEffectImproved {
			MakePermanent(JudgementOfTheCrusaderAura(nil, target, level, 1.15, 0))
		}
	}

	if debuffs.ImprovedShadowBolt && targetIdx == 0 {
		ExternalIsbCaster(debuffs, target)
	}

	if debuffs.ShadowWeaving {
		aura := ShadowWeavingAura(target, 5)
		SchedulePeriodicDebuffApplication(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        5,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		})
	}

	if debuffs.OccultPoison {
		aura := OccultPoisonDebuffAura(target, level)
		SchedulePeriodicDebuffApplication(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        5,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		})
	}

	if debuffs.MekkatorqueFistDebuff {
		MakePermanent(MekkatorqueFistDebuffAura(target, level))
	}

	if debuffs.SerpentsStrikerFistDebuff {
		MakePermanent(SerpentsStrikerFistDebuffAura(target, level))
	}

	if debuffs.MarkOfChaos {
		MakePermanent(MarkOfChaosDebuffAura(target))
	} else {
		if debuffs.CurseOfElements {
			MakePermanent(CurseOfElementsAura(target, level))
		}

		if debuffs.CurseOfShadow {
			MakePermanent(CurseOfShadowAura(target, level))
		}
	}

	if debuffs.ImprovedScorch && targetIdx == 0 {
		aura := ImprovedScorchAura(target)
		SchedulePeriodicDebuffApplication(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        5,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		})
	}

	if debuffs.WintersChill && targetIdx == 0 {
		aura := WintersChillAura(target)
		SchedulePeriodicDebuffApplication(aura, PeriodicActionOptions{
			Period:          time.Millisecond * 1500,
			NumTicks:        5,
			TickImmediately: true,
			Priority:        ActionPriorityDOT, // High prio
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			},
		})
	}

	if debuffs.Stormstrike {
		MakePermanent(StormstrikeAura(target))
	} else if debuffs.Dreamstate {
		MakePermanent(DreamstateAura(target))
	}

	if debuffs.GiftOfArthas {
		MakePermanent(GiftOfArthasAura(target))
	}

	if debuffs.HolySunder {
		MakePermanent(HolySunderAura(target))
	}

	if debuffs.CurseOfVulnerability {
		MakePermanent(CurseOfVulnerabilityAura(target))
	}

	if debuffs.Mangle {
		MakePermanent(MangleAura(target, level))
	}

	// Major Armor Debuffs
	if targetIdx == 0 {
		if debuffs.ExposeArmor != proto.TristateEffect_TristateEffectMissing {
			aura := ExposeArmorAura(target, TernaryInt32(debuffs.ExposeArmor == proto.TristateEffect_TristateEffectRegular, 0, 2), level)
			SchedulePeriodicDebuffApplication(aura, PeriodicActionOptions{
				Period:   time.Second * 3,
				NumTicks: 1,
				OnAction: func(sim *Simulation) {
					aura.Activate(sim)
				},
			})
		}

		if debuffs.SebaciousPoison != proto.TristateEffect_TristateEffectMissing {
			aura := SebaciousPoisonAura(target, TernaryInt32(debuffs.SebaciousPoison == proto.TristateEffect_TristateEffectRegular, 0, 2), level)
			SchedulePeriodicDebuffApplication(aura, PeriodicActionOptions{
				Period:   time.Second * 0,
				NumTicks: 1,
				OnAction: func(sim *Simulation) {
					aura.Activate(sim)
				},
			})
		}

		if debuffs.SunderArmor {
			// Sunder Armor
			aura := SunderArmorAura(target, level)
			SchedulePeriodicDebuffApplication(aura, PeriodicActionOptions{
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
			})
		}
	}

	if debuffs.CurseOfRecklessness {
		MakePermanent(CurseOfRecklessnessAura(target, level))
	}

	if debuffs.FaerieFire || debuffs.ImprovedFaerieFire {
		MakePermanent(FaerieFireAura(target, level))
	}

	if debuffs.ImprovedFaerieFire {
		MakePermanent(ImprovedFaerieFireAura(target))
	}

	if debuffs.MeleeHunterDodgeDebuff {
		MakePermanent(MeleeHunterDodgeReductionAura(target, level))
	}

	if debuffs.CurseOfWeakness != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(CurseOfWeaknessAura(target, GetTristateValueInt32(debuffs.CurseOfWeakness, 0, 3), level))
	}

	if debuffs.DemoralizingRoar != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingRoarAura(target, GetTristateValueInt32(debuffs.DemoralizingRoar, 0, 5), level))
	}
	if debuffs.DemoralizingShout != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(DemoralizingShoutAura(target, 0, GetTristateValueInt32(debuffs.DemoralizingShout, 0, 5), level))
	}
	if debuffs.AtrophicPoison {
		MakePermanent(AtrophicPoisonAura(target))
	}

	if debuffs.HuntersMark != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(HuntersMarkAura(target, GetTristateValueInt32(debuffs.HuntersMark, 0, 5), level))
	}

	// Atk spd reduction
	if debuffs.ThunderClap != proto.TristateEffect_TristateEffectMissing {
		// +6% from Furious Thunder rune
		MakePermanent(ThunderClapAura(target, 8205, time.Second*10, GetTristateValueInt32(debuffs.ThunderClap, 10, 16)))
	}
	if debuffs.Waylay {
		MakePermanent(WaylayAura(target))
	}
	if debuffs.Thunderfury {
		MakePermanent(ThunderfuryASAura(target, level))
	}
	if debuffs.NumbingPoison {
		MakePermanent(NumbingPoisonAura(target))
	}

	// Miss
	if debuffs.InsectSwarm && targetIdx == 0 {
		MakePermanent(InsectSwarmAura(target, level))
	}
	if debuffs.ScorpidSting && targetIdx == 0 {
		MakePermanent(ScorpidStingAura(target))
	}

	// Karazhan random suffixes
	if debuffs.FrostFever {
		MakePermanent(FrostFeverAura(target))
	}
	if debuffs.BloodPlague {
		MakePermanent(BloodPlagueAura(target))
	}
}

func StormstrikeAura(unit *Unit) *Aura {
	return exclusiveNatureDamageTakenAura(unit, "Stormstrike", ActionID{SpellID: 17364})
}

func DreamstateAura(unit *Unit) *Aura {
	aura := exclusiveNatureDamageTakenAura(unit, "Dreamstate", ActionID{SpellID: 408258})
	aura.NewExclusiveEffect("ArcaneDamageTaken", false, ExclusiveEffect{
		Priority: 20,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1.2
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= 1.2
		},
	})
	return aura
}

func exclusiveNatureDamageTakenAura(unit *Unit, label string, actionID ActionID) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
		Duration: time.Second * 12,
	})

	aura.NewExclusiveEffect("NatureDamageTaken", false, ExclusiveEffect{
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

func ExternalIsbCaster(_ *proto.Debuffs, target *Unit) {
	isbConfig := target.Env.Raid.Parties[0].Players[0].GetCharacter().IsbConfig
	baseStacks := TernaryInt32(isbConfig.hasShadowflameRune, ISBNumStacksShadowflame, ISBNumStacksBase)
	isbAura := ImprovedShadowBoltAura(target, 5)
	isbCrit := isbConfig.casterCrit / 100.0
	var pa *PendingAction
	MakePermanent(target.GetOrRegisterAura(Aura{
		Label: "Isb External Proc Aura",
		OnGain: func(aura *Aura, sim *Simulation) {
			pa = NewPeriodicAction(sim, PeriodicActionOptions{
				Period: DurationFromSeconds(isbConfig.shadowBoltFrequency),
				OnAction: func(s *Simulation) {
					for i := 0; i < int(isbConfig.isbWarlocks); i++ {
						if sim.Proc(isbCrit, "External Isb Crit") {
							isbAura.Activate(sim)
							isbAura.SetStacks(sim, baseStacks)
						} else if isbAura.IsActive() {
							isbAura.RemoveStack(sim)
						}
					}
				},
			})
			sim.AddPendingAction(pa)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			pa.Cancel(sim)
		},
	}))
}

type IsbConfig struct {
	hasShadowflameRune  bool
	shadowBoltFrequency float64
	casterCrit          float64
	isbWarlocks         int32
	isbShadowPriests    int32
}

func (character *Character) createIsbConfig(player *proto.Player) {
	character.IsbConfig = IsbConfig{
		hasShadowflameRune:  player.IsbUsingShadowflame,
		shadowBoltFrequency: player.IsbSbFrequency,
		casterCrit:          player.IsbCrit,
		isbWarlocks:         player.IsbWarlocks,
		isbShadowPriests:    player.IsbSpriests,
	}
	//Defaults if not configured
	if character.IsbConfig.shadowBoltFrequency == 0.0 {
		character.IsbConfig.shadowBoltFrequency = 3.0
	}
	if character.IsbConfig.casterCrit == 0.0 {
		character.IsbConfig.casterCrit = 25.0
	}
	if character.IsbConfig.isbWarlocks == 0 {
		character.IsbConfig.isbWarlocks = 1
	}
}

const (
	ISBNumStacksBase        = 4
	ISBNumStacksShadowflame = 30
)

func ImprovedShadowBoltAura(unit *Unit, rank int32) *Aura {
	isbLabel := "Improved Shadow Bolt"
	if unit.GetAura(isbLabel) != nil {
		return unit.GetAura(isbLabel)
	}

	isbConfig := unit.Env.Raid.Parties[0].Players[0].GetCharacter().IsbConfig

	priestGcds := []bool{false, true, true, true, true, true}
	priestCurGcd := 0
	externalShadowPriests := isbConfig.isbShadowPriests
	var priestPa *PendingAction

	damageMulti := 1. + 0.04*float64(rank)
	aura := unit.GetOrRegisterAura(Aura{
		Label:     isbLabel,
		ActionID:  ActionID{SpellID: 17800},
		Duration:  12 * time.Second,
		MaxStacks: 30,
		OnReset: func(aura *Aura, sim *Simulation) {
			// External shadow priests simulation
			if externalShadowPriests > 0 {
				priestCurGcd = 0
				priestPa = NewPeriodicAction(sim, PeriodicActionOptions{
					Period: GCDDefault,
					OnAction: func(s *Simulation) {
						if priestGcds[priestCurGcd] {
							for i := 0; i < int(externalShadowPriests); i++ {
								if aura.IsActive() {
									aura.RemoveStack(sim)
								}
							}
						}
						priestCurGcd++
						if priestCurGcd >= len(priestGcds) {
							priestCurGcd = 0
						}
					},
				})
				sim.AddPendingAction(priestPa)
			}
		},
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= damageMulti
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= damageMulti
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.SpellSchool.Matches(SpellSchoolShadow) && result.Landed() && result.Damage > 0 {
				aura.RemoveStack(sim)
			}
		},
	})

	return aura
}

var ShadowWeavingSpellIDs = [6]int32{0, 15257, 15331, 15332, 15333, 15334}

func ShadowWeavingAura(unit *Unit, rank int) *Aura {
	spellId := ShadowWeavingSpellIDs[rank]
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

func SchedulePeriodicDebuffApplication(aura *Aura, options PeriodicActionOptions) {
	aura.OnReset = func(aura *Aura, sim *Simulation) {
		aura.Duration = NeverExpires
		StartPeriodicAction(sim, options)
	}
}

const JudgementAuraTag = "Judgement"

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
		Label:    "Judgement of Wisdom",
		ActionID: actionID,
		Tag:      JudgementAuraTag,
		Duration: time.Second * 10,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			unit := spell.Unit
			if !unit.HasManaBar() {
				return
			}

			if spell.ProcMask.Matches(ProcMaskEmpty|ProcMaskProc|ProcMaskSpellDamageProc) && !spell.Flags.Matches(SpellFlagNotAProc) {
				return // Phantom spells (Romulo's, Lightning Capacitor, etc.) don't proc JoW.
			}

			if !spell.ProcMask.Matches(ProcMaskDirect) {
				return
			}

			// melee auto attacks don't even need to land
			if !result.Landed() && !spell.ProcMask.Matches(ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("jow") < 0.5 {
				if unit.JowManaMetrics == nil {
					unit.JowManaMetrics = unit.NewManaMetrics(actionID)
				}
				// JoW returns flat mana
				unit.AddMana(sim, jowMana, unit.JowManaMetrics)
			}
		},
	})
}

func JudgementOfLightAura(target *Unit, level int32, units []*Unit) *Aura {
	auraActionID := ActionID{SpellID: map[int32]int32{
		30: 20185,
		40: 20344,
		50: 20345,
		60: 20346,
	}[level]}
	healActionID := ActionID{SpellID: map[int32]int32{
		30: 20267,
		40: 20341,
		50: 20342,
		60: 20343,
	}[level]}

	jolHealth := map[int32]float64{
		30: 25.0,
		40: 34.0,
		50: 49.0,
		60: 61.0,
	}[level]

	for _, playerOrPet := range units {
		unit := playerOrPet
		unit.GetOrRegisterSpell(SpellConfig{
			ActionID:         healActionID,
			SpellSchool:      SpellSchoolHoly,
			ProcMask:         ProcMaskEmpty,
			Flags:            SpellFlagPassiveSpell | SpellFlagHelpful | SpellFlagNoOnCastComplete,
			DamageMultiplier: 1.0,
			ThreatMultiplier: 1.0,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				spell.CalcAndDealHealing(sim, unit, jolHealth, spell.OutcomeAlwaysHit)
			},
		})
	}

	return target.GetOrRegisterAura(Aura{
		Label:    "Judgement of Light",
		ActionID: auraActionID,
		Tag:      JudgementAuraTag,
		Duration: time.Second * 10,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			unit := spell.Unit
			if !unit.HasHealthBar() {
				return
			}

			healingSpell := unit.GetSpell(healActionID)
			if healingSpell == nil {
				return
			}

			if !spell.ProcMask.Matches(ProcMaskMelee) || !result.Landed() {
				return
			}

			if sim.RandomFloat("jol") < 0.5 {
				healingSpell.Cast(sim, unit)
			}
		},
	})
}

func JudgementOfTheCrusaderAura(caster *Unit, target *Unit, level int32, mult float64, extraBonus float64) *Aura {
	var spellId int32
	var bonus float64

	switch level {
	case 25:
		spellId = 20300
		bonus = 50
	case 40:
		spellId = 20301
		bonus = 80
	case 50:
		spellId = 20302
		bonus = 110
	default:
		spellId = 20303
		bonus = 140
	}

	bonus *= mult
	bonus += extraBonus

	return target.GetOrRegisterAura(Aura{
		Label:    "Judgement of the Crusader",
		ActionID: ActionID{SpellID: spellId},
		Tag:      JudgementAuraTag,
		Duration: 10 * time.Second,

		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexHoly] += bonus
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexHoly] -= bonus
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.Unit != caster { // caster is nil for permanent auras
				return
			}
			if result.Landed() && spell.ProcMask.Matches(ProcMaskMelee) {
				aura.Refresh(sim)
			}
		},
	})
}

func OccultPoisonDebuffAura(target *Unit, playerLevel int32) *Aura {
	if playerLevel < 54 {
		panic("Occult Poison requires level 54+")
	}

	aura := target.GetOrRegisterAura(Aura{
		Label:     "Occult Poison II",
		ActionID:  ActionID{SpellID: 1214170},
		Duration:  time.Second * 12,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			multiplier := (1 + .04*float64(newStacks)) / (1 + .04*float64(oldStacks))

			// Applies too all except Holy
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= multiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= multiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= multiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= multiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= multiplier
		},
	})

	return aura
}

func MekkatorqueFistDebuffAura(target *Unit, playerLevel int32) *Aura {
	if playerLevel < 40 {
		panic("Mekkatorque's Arcano-Shredder requires level 40+")
	}

	spellID := 434841
	resistance := 45.0
	dmgMod := 1.06

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Mekkatorque Debuff",
		ActionID: ActionID{SpellID: int32(spellID)},
		Duration: time.Second * 20,
	})

	// 0.01 priority as this overwrites the other spells of this category and does not allow them to be recast
	spellSchoolDamageEffect(aura, stats.SchoolIndexArcane, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexFire, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexFrost, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexHoly, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexNature, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexShadow, dmgMod, 0.01, true)

	spellSchoolResistanceEffect(aura, stats.SchoolIndexArcane, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexFire, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexFrost, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexHoly, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexNature, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexShadow, resistance, 0.01, true)

	return aura
}

// Mark of Chaos does not stack with Curse of Shadows and Elements
func MarkOfChaosDebuffAura(target *Unit) *Aura {
	// That's right, 10.01%. Sneaky enough to override lock curses without being much stronger
	dmgMod := 1.1001
	resistance := 75.01

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Mark of Chaos",
		ActionID: ActionID{SpellID: 461615},
		Duration: time.Second, // Duration is set by the applying curse
	})

	// Applies too all except Holy
	// 0.01 priority as this overwrites the other spells of this category and does not allow them to be recast
	spellSchoolDamageEffect(aura, stats.SchoolIndexArcane, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexFire, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexFrost, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexNature, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexShadow, dmgMod, 0.01, true)

	spellSchoolResistanceEffect(aura, stats.SchoolIndexArcane, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexFire, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexFrost, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexNature, resistance, 0.01, true)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexShadow, resistance, 0.01, true)

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
		Label:      "Curse of Elements",
		DispelType: DispelType_Curse,
		ActionID:   ActionID{SpellID: spellID},
		Duration:   time.Minute * 5,
	})
	spellSchoolDamageEffect(aura, stats.SchoolIndexFire, dmgMod, 0.0, false)
	spellSchoolDamageEffect(aura, stats.SchoolIndexFrost, dmgMod, 0.0, false)

	spellSchoolResistanceEffect(aura, stats.SchoolIndexFire, resistance, 0.0, false)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexFrost, resistance, 0.0, false)

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
		Label:      "Curse of Shadow",
		DispelType: DispelType_Curse,
		ActionID:   ActionID{SpellID: spellID},
		Duration:   time.Minute * 5,
	})
	spellSchoolDamageEffect(aura, stats.SchoolIndexArcane, dmgMod, 0.0, false)
	spellSchoolDamageEffect(aura, stats.SchoolIndexShadow, dmgMod, 0.0, false)

	spellSchoolResistanceEffect(aura, stats.SchoolIndexArcane, resistance, 0.0, false)
	spellSchoolResistanceEffect(aura, stats.SchoolIndexShadow, resistance, 0.0, false)

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

func spellSchoolResistanceEffect(aura *Aura, school stats.SchoolIndex, amount float64, extraPriority float64, exclusive bool) *ExclusiveEffect {
	return aura.NewExclusiveEffect("resistance"+strconv.Itoa(int(school)), exclusive, ExclusiveEffect{
		Priority: amount + extraPriority,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddResistancesDynamic(sim, -amount)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddResistancesDynamic(sim, amount)
		},
	})
}

func GiftOfArthasAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Gift of Arthas",
		ActionID: ActionID{SpellID: 11374},
		Duration: time.Minute * 3,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexPhysical] += 8
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexPhysical] -= 8
		},
	})
}

func HolySunderAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Holy Sunder",
		ActionID: ActionID{SpellID: 9176},
		Duration: time.Minute * 1,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] -= 50
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.stats[stats.Armor] += 50
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
			aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexPhysical] += debuffBonusDamage
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolBonusDamageTaken[stats.SchoolIndexPhysical] -= debuffBonusDamage
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
		Label:      "Curse of Vulnerability",
		DispelType: DispelType_Curse,
		ActionID:   ActionID{SpellID: 427143},
		Duration:   time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
				aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] += 2
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			for si := stats.SchoolIndexPhysical; si < stats.SchoolLen; si++ {
				aura.Unit.PseudoStats.SchoolBonusDamageTaken[si] -= 2
			}
		},
	})
}

func MangleAura(target *Unit, _ int32) *Aura {
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
			ee.Aura.Unit.PseudoStats.BleedDamageTakenMultiplier *= multiplier
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.BleedDamageTakenMultiplier /= multiplier
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

func WintersChillAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:     "Winter's Chill",
		ActionID:  ActionID{SpellID: 28593},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks, newStacks int32) {
			aura.Unit.PseudoStats.SchoolCritTakenChance[stats.SchoolIndexFrost] -= 0.02 * float64(oldStacks)
			aura.Unit.PseudoStats.SchoolCritTakenChance[stats.SchoolIndexFrost] += 0.02 * float64(newStacks)
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

	arpen *= []float64{1, 1.25, 1.5}[improvedEA]

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

func SebaciousPoisonAura(target *Unit, improvedEA int32, playerLevel int32) *Aura {
	if playerLevel < 60 {
		return nil
	}

	spellID := map[int32]int32{
		60: 439462,
	}[playerLevel]

	arpen := map[int32]float64{
		60: 1700,
	}[playerLevel]

	arpen *= []float64{1, 1.25, 1.5}[improvedEA]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Sebacious Poison",
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Second * 15,
	})

	aura.NewExclusiveEffect(majorArmorReductionEffectCategory, true, ExclusiveEffect{
		Priority: arpen,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, -ee.Priority)

			// p8 DPS tier bonus tracking
			fmt.Println("Occult Poison activated")
			//rogue.PoisonsActive--
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, ee.Priority)
		},
	})

	return aura
}

func HomunculiAttackSpeedAura(target *Unit, _ int32) *Aura {
	multiplier := 1.1

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Cripple (Homunculus)",
		ActionID: ActionID{SpellID: 402808},
		Duration: time.Second * 15,
	})

	AtkSpeedReductionEffect(aura, multiplier)
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
	ap := float64(190 - 3*(60-playerLevel))

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Demoralize (Homunculus)",
		ActionID: ActionID{SpellID: 402811},
		Duration: time.Second * 15,
	})

	apReductionEffect(aura, ap)

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

	ap := map[int32]float64{
		25: 20,
		40: 45,
		50: 65,
		60: 90,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:      "Curse of Recklessness",
		DispelType: DispelType_Curse,
		ActionID:   ActionID{SpellID: spellID},
		Duration:   time.Minute * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, -arpen)
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, ap)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.Armor, arpen)
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, -ap)
		},
	})
	return aura
}

// Decreases the armor of the target by X for 40 sec.
// Improved: Your Faerie Fire and Faerie Fire (Feral) also increase the chance for all attacks to hit that target by 1% for 40 sec.
func FaerieFireAura(target *Unit, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 770,
		40: 778,
		50: 9749,
		60: 9907,
	}[playerLevel]

	return faerieFireAuraInternal(target, "Faerie Fire", spellID, playerLevel)
}

func FaerieFireFeralAura(target *Unit, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		40: 17390,
		50: 17391,
		60: 17392,
	}[playerLevel]

	return faerieFireAuraInternal(target, "Faerie Fire (Feral)", spellID, playerLevel)
}

func faerieFireAuraInternal(target *Unit, label string, spellID int32, playerLevel int32) *Aura {
	arPen := map[int32]float64{
		25: 175,
		40: 285,
		50: 395,
		60: 505,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Second * 40,
	})

	aura.NewExclusiveEffect("Faerie Fire", true, ExclusiveEffect{
		Priority: arPen,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Armor, -arPen)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Armor, arPen)
		},
	})

	return aura
}

func ImprovedFaerieFireAura(target *Unit) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "S03 - Item - T1 - Druid - Feral 2P Bonus",
		ActionID: ActionID{SpellID: 455864},
		Duration: time.Second * 40,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusMeleeHitRatingTaken += 1 * MeleeHitRatingPerHitChance
			aura.Unit.PseudoStats.BonusSpellHitRatingTaken += 1 * SpellHitRatingPerHitChance
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusMeleeHitRatingTaken -= 1 * MeleeHitRatingPerHitChance
			aura.Unit.PseudoStats.BonusSpellHitRatingTaken -= 1 * SpellHitRatingPerHitChance
		},
	})
}

func MeleeHunterDodgeReductionAura(target *Unit, _ int32) *Aura {
	return target.GetOrRegisterAura(Aura{
		Label:    "Stalked",
		ActionID: ActionID{SpellID: 456389},
		Duration: time.Second * 30,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.DodgeReduction += 0.01
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.DodgeReduction -= 0.01
		},
	})
}

func CurseOfWeaknessAura(target *Unit, points int32, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 6205,
		40: 7646,
		50: 11707,
		60: 11708,
	}[playerLevel]

	modDmgReduction := map[int32]float64{
		25: -10,
		40: -15,
		50: -22,
		60: -31,
	}[playerLevel]

	modDmgReduction *= []float64{1, 1.06, 1.13, 1.20}[points]
	modDmgReduction = math.Floor(modDmgReduction)

	aura := target.GetOrRegisterAura(Aura{
		Label:      "Curse of Weakness" + strconv.Itoa(int(points)),
		DispelType: DispelType_Curse,
		ActionID:   ActionID{SpellID: spellID},
		Duration:   time.Minute * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamage += modDmgReduction
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusPhysicalDamage -= modDmgReduction
		},
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

func DemoralizingRoarAura(target *Unit, points int32, playerLevel int32) *Aura {
	spellID := map[int32]int32{
		25: 1735,
		40: 9490,
		50: 9747,
		60: 9898,
	}[playerLevel]
	baseAPReduction := map[int32]float64{
		25: 55,
		40: 73,
		50: 108,
		60: 138,
	}[playerLevel]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingRoar-" + strconv.Itoa(int(points)),
		ActionID: ActionID{SpellID: spellID},
		Duration: time.Second * 30,
	})
	apReductionEffect(aura, math.Floor(baseAPReduction*(1+0.08*float64(points))))
	return aura
}

const DemoralizingShoutRanks = 5

var DemoralizingShoutSpellId = [DemoralizingShoutRanks + 1]int32{0, 1160, 6190, 11554, 11555, 11556}
var DemoralizingShoutBaseAP = [DemoralizingShoutRanks + 1]float64{0, 45, 56, 76, 111, 146}
var DemoralizingShoutLevel = [DemoralizingShoutRanks + 1]int{0, 14, 24, 34, 44, 54}

func DemoralizingShoutAura(target *Unit, boomingVoicePts int32, impDemoShoutPts int32, playerLevel int32) *Aura {
	rank := LevelToDebuffRank[DemoralizingShout][playerLevel]
	spellId := DemoralizingShoutSpellId[rank]
	baseAPReduction := DemoralizingShoutBaseAP[rank]

	aura := target.GetOrRegisterAura(Aura{
		Label:    "DemoralizingShout-" + strconv.Itoa(int(impDemoShoutPts)),
		ActionID: ActionID{SpellID: spellId},
		Duration: time.Duration(float64(time.Second*30) * (1 + 0.1*float64(boomingVoicePts))),
	})
	apReductionEffect(aura, math.Floor(baseAPReduction*(1+0.08*float64(impDemoShoutPts))))
	return aura
}

func AtrophicPoisonAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Atrophic Poison",
		ActionID: ActionID{SpellID: 439473},
		Duration: time.Second * 15,
	})
	apReductionEffect(aura, 205)
	return aura
}

func VindicationAura(target *Unit, points int32, _ int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Vindication",
		ActionID: ActionID{SpellID: 26016},
		Duration: time.Second * 10,
	})
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

func ThunderClapAura(target *Unit, spellID int32, duration time.Duration, atkSpeedReductionPercent int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "ThunderClap-" + strconv.Itoa(int(atkSpeedReductionPercent)),
		ActionID: ActionID{SpellID: spellID},
		Duration: duration,
	})
	AtkSpeedReductionEffect(aura, 1+0.01*float64(atkSpeedReductionPercent))
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

func ThunderfuryASAura(target *Unit, _ int32) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Thunderfury",
		ActionID: ActionID{SpellID: 21992},
		Duration: time.Second * 12,
	})
	AtkSpeedReductionEffect(aura, 1.2)
	return aura
}

func NumbingPoisonAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Numbing Poison",
		ActionID: ActionID{SpellID: 439472},
		Duration: time.Second * 15,
	})
	AtkSpeedReductionEffect(aura, 1.2)
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

func InsectSwarmAura(target *Unit, level int32) *Aura {
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

func SerpentsStrikerFistDebuffAura(target *Unit, playerLevel int32) *Aura {
	if playerLevel < 50 {
		panic("Serpent's Striker requires level 50+")
	}

	spellID := 447894
	resistance := 60.0
	dmgMod := 1.08

	aura := target.GetOrRegisterAura(Aura{
		Label:    "Serpents Striker Debuff",
		ActionID: ActionID{SpellID: int32(spellID)},
		Duration: time.Second * 20,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.NatureResistance: -resistance,
			})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.NatureResistance: resistance,
			})
		},
	})

	// 0.01 priority as this overwrites the other spells of this category and does not allow them to be recast
	spellSchoolDamageEffect(aura, stats.SchoolIndexNature, dmgMod, 0.01, true)
	spellSchoolDamageEffect(aura, stats.SchoolIndexHoly, dmgMod, 0.01, true)
	return aura
}

func FrostFeverAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Frost Fever",
		ActionID: ActionID{SpellID: FrostFeverAuraID}.WithTag(2),
		Duration: time.Second * 21,
		Tag:      "Obliterate",
	})
	return aura
}

func BloodPlagueAura(target *Unit) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:    "Blood Plague",
		ActionID: ActionID{SpellID: BloodPlagueAuraID}.WithTag(2),
		Duration: time.Second * 15,
		Tag:      "Obliterate",
	})
	return aura
}
