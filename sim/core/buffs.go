package core

import (
	"math"
	"time"

	googleProto "google.golang.org/protobuf/proto"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type BuffName int32

const (
	// General Buffs
	ArcaneIntellect BuffName = iota
	BattleShout
	BlessingOfMight
	BlessingOfWisdom
	BloodPact
	CommandingShout
	DevotionAura
	DivineSpirit
	GraceOfAir
	ManaSpring
	MarkOfTheWild
	PowerWordFortitude
	StrengthOfEarth
	TrueshotAura
	Windfury
	SanctityAura
	BattleSquawk

	// Resistance
	AspectOfTheWild
	FireResistanceAura
	FireResistanceTotem
	FrostResistanceTotem
	FrostResistanceAura
	NatureResistanceTotem
	ShadowProtection
	ShadowResistanceAura

	// Scrolls
	ScrollOfAgility
	ScrollOfIntellect
	ScrollOfSpirit
	ScrollOfStrength
	ScrollOfStamina
	ScrollOfProtection
)

var LevelToBuffRank = map[BuffName]map[int32]int32{
	BattleShout: {
		25: 3,
		40: 4,
		50: 5,
		60: 7,
	},
	GraceOfAir: {
		50: 1,
		60: 3,
	},
	StrengthOfEarth: {
		25: 2,
		40: 3,
		50: 3,
		60: 5,
	},
	Windfury: {
		40: 1,
		50: 2,
		60: 3,
	},
}

// Stats from buffs pre-tristate buffs
var BuffSpellByLevel = map[BuffName]map[int32]stats.Stats{
	ArcaneIntellect: {
		25: stats.Stats{
			stats.Intellect: 7, // rank 2
		},
		40: stats.Stats{
			stats.Intellect: 15,
		},
		50: stats.Stats{
			stats.Intellect: 22,
		},
		60: stats.Stats{
			stats.Intellect: 31,
		},
	},
	DivineSpirit: {
		25: stats.Stats{
			stats.Spirit: 0,
		},
		40: stats.Stats{
			stats.Spirit: 23,
		},
		50: stats.Stats{
			stats.Spirit: 33,
		},
		60: stats.Stats{
			stats.Spirit: 40,
		},
	},
	AspectOfTheWild: {
		25: stats.Stats{
			stats.NatureResistance: 0,
		},
		40: stats.Stats{
			stats.NatureResistance: 0,
		},
		50: stats.Stats{
			stats.NatureResistance: 45,
		},
		60: stats.Stats{
			stats.NatureResistance: 60,
		},
	},
	BattleShout: {
		25: stats.Stats{
			stats.AttackPower: 57,
		},
		40: stats.Stats{
			stats.AttackPower: 93,
		},
		50: stats.Stats{
			stats.AttackPower: 138,
		},
		60: stats.Stats{
			stats.AttackPower: 193,
		},
	},
	BlessingOfMight: {
		25: stats.Stats{
			stats.AttackPower: 55,
		},
		40: stats.Stats{
			stats.AttackPower: 85,
		},
		50: stats.Stats{
			stats.AttackPower: 115,
		},
		60: stats.Stats{
			stats.AttackPower: 185,
		},
	},
	BlessingOfWisdom: {
		25: stats.Stats{
			stats.MP5: 15,
		},
		40: stats.Stats{
			stats.MP5: 20,
		},
		50: stats.Stats{
			stats.MP5: 25,
		},
		60: stats.Stats{
			stats.MP5: 33,
		},
	},
	BloodPact: {
		25: stats.Stats{
			stats.Stamina: 9,
		},
		40: stats.Stats{
			stats.Stamina: 30,
		},
		50: stats.Stats{
			stats.Stamina: 38,
		},
		60: stats.Stats{
			stats.Stamina: 42,
		},
	},
	CommandingShout: {
		25: stats.Stats{
			stats.Stamina: 18,
		},
		40: stats.Stats{
			stats.Stamina: 28,
		},
		50: stats.Stats{
			stats.Stamina: 35,
		},
		60: stats.Stats{
			stats.Stamina: 42,
		},
	},
	DevotionAura: {
		25: stats.Stats{
			stats.BonusArmor: 275,
		},
		40: stats.Stats{
			stats.BonusArmor: 505,
		},
		50: stats.Stats{
			stats.BonusArmor: 620,
		},
		60: stats.Stats{
			stats.BonusArmor: 735,
		},
	},
	GraceOfAir: {
		25: stats.Stats{
			stats.Agility: 0,
		},
		40: stats.Stats{
			stats.Agility: 0,
		},
		50: stats.Stats{
			stats.Agility: 43,
		},
		60: stats.Stats{
			stats.Agility: 77,
		},
	},
	FireResistanceAura: {
		25: stats.Stats{
			stats.FireResistance: 0,
		},
		40: stats.Stats{
			stats.FireResistance: 30,
		},
		50: stats.Stats{
			stats.FireResistance: 30,
		},
		60: stats.Stats{
			stats.FireResistance: 60,
		},
	},
	FireResistanceTotem: {
		25: stats.Stats{
			stats.FireResistance: 0,
		},
		40: stats.Stats{
			stats.FireResistance: 30,
		},
		50: stats.Stats{
			stats.FireResistance: 45,
		},
		60: stats.Stats{
			stats.FireResistance: 60,
		},
	},
	FrostResistanceAura: {
		25: stats.Stats{
			stats.FrostResistance: 0,
		},
		40: stats.Stats{
			stats.FrostResistance: 30,
		},
		50: stats.Stats{
			stats.FrostResistance: 45,
		},
		60: stats.Stats{
			stats.FrostResistance: 60,
		},
	},
	FrostResistanceTotem: {
		25: stats.Stats{
			stats.FrostResistance: 30,
		},
		40: stats.Stats{
			stats.FrostResistance: 45,
		},
		50: stats.Stats{
			stats.FrostResistance: 45,
		},
		60: stats.Stats{
			stats.FrostResistance: 60,
		},
	},
	ManaSpring: {
		25: stats.Stats{
			stats.MP5: 0,
		},
		40: stats.Stats{
			stats.MP5: 15,
		},
		50: stats.Stats{
			stats.MP5: 20,
		},
		60: stats.Stats{
			stats.MP5: 25,
		},
	},
	MarkOfTheWild: {
		25: stats.Stats{
			stats.BonusArmor:       105,
			stats.Stamina:          4,
			stats.Agility:          4,
			stats.Strength:         4,
			stats.Intellect:        4,
			stats.Spirit:           4,
			stats.ArcaneResistance: 0,
			stats.ShadowResistance: 0,
			stats.NatureResistance: 0,
			stats.FireResistance:   0,
			stats.FrostResistance:  0,
		},
		40: stats.Stats{
			stats.BonusArmor:       195,
			stats.Stamina:          8,
			stats.Agility:          8,
			stats.Strength:         8,
			stats.Intellect:        8,
			stats.Spirit:           8,
			stats.ArcaneResistance: 10,
			stats.ShadowResistance: 10,
			stats.NatureResistance: 10,
			stats.FireResistance:   10,
			stats.FrostResistance:  10,
		},
		50: stats.Stats{
			stats.BonusArmor:       240,
			stats.Stamina:          10,
			stats.Agility:          10,
			stats.Strength:         10,
			stats.Intellect:        10,
			stats.Spirit:           10,
			stats.ArcaneResistance: 15,
			stats.ShadowResistance: 15,
			stats.NatureResistance: 15,
			stats.FireResistance:   15,
			stats.FrostResistance:  15,
		},
		60: stats.Stats{
			stats.BonusArmor:       285,
			stats.Stamina:          12,
			stats.Agility:          12,
			stats.Strength:         12,
			stats.Intellect:        12,
			stats.Spirit:           12,
			stats.ArcaneResistance: 20,
			stats.ShadowResistance: 20,
			stats.NatureResistance: 20,
			stats.FireResistance:   20,
			stats.FrostResistance:  20,
		},
	},
	NatureResistanceTotem: {
		25: stats.Stats{
			stats.NatureResistance: 0,
		},
		40: stats.Stats{
			stats.NatureResistance: 30,
		},
		50: stats.Stats{
			stats.NatureResistance: 45,
		},
		60: stats.Stats{
			stats.NatureResistance: 60,
		},
	},
	PowerWordFortitude: {
		25: stats.Stats{
			stats.Stamina: 20,
		},
		40: stats.Stats{
			stats.Stamina: 32,
		},
		50: stats.Stats{
			stats.Stamina: 43,
		},
		60: stats.Stats{
			stats.Stamina: 54,
		},
	},
	ShadowProtection: {
		25: stats.Stats{
			stats.ShadowResistance: 0,
		},
		40: stats.Stats{
			stats.ShadowResistance: 30,
		},
		50: stats.Stats{
			stats.ShadowResistance: 45,
		},
		60: stats.Stats{
			stats.ShadowResistance: 60,
		},
	},
	ShadowResistanceAura: {
		25: stats.Stats{
			stats.ShadowResistance: 0,
		},
		40: stats.Stats{
			stats.ShadowResistance: 45,
		},
		50: stats.Stats{
			stats.ShadowResistance: 45,
		},
		60: stats.Stats{
			stats.ShadowResistance: 60,
		},
	},
	TrueshotAura: {
		25: stats.Stats{
			stats.RangedAttackPower: 0,
		},
		40: stats.Stats{
			stats.RangedAttackPower: 100,
		},
		50: stats.Stats{
			stats.RangedAttackPower: 150,
		},
		60: stats.Stats{
			stats.RangedAttackPower: 200,
		},
	},
	StrengthOfEarth: {
		25: stats.Stats{
			stats.Strength: 20,
		},
		40: stats.Stats{
			stats.Strength: 36,
		},
		50: stats.Stats{
			stats.Strength: 36,
		},
		60: stats.Stats{
			stats.Strength: 77,
		},
	},
	ScrollOfAgility: {
		25: stats.Stats{
			stats.Agility: 9,
		},
		40: stats.Stats{
			stats.Agility: 13,
		},
		50: stats.Stats{
			stats.Agility: 17,
		},
		60: stats.Stats{
			stats.Agility: 17,
		},
	},
	ScrollOfIntellect: {
		25: stats.Stats{
			stats.Intellect: 8,
		},
		40: stats.Stats{
			stats.Intellect: 12,
		},
		50: stats.Stats{
			stats.Intellect: 16,
		},
		60: stats.Stats{
			stats.Intellect: 16,
		},
	},
	ScrollOfSpirit: {
		25: stats.Stats{
			stats.Spirit: 7,
		},
		40: stats.Stats{
			stats.Spirit: 11,
		},
		50: stats.Stats{
			stats.Spirit: 15,
		},
		60: stats.Stats{
			stats.Spirit: 15,
		},
	},
	ScrollOfStamina: {
		25: stats.Stats{
			stats.Stamina: 8,
		},
		40: stats.Stats{
			stats.Stamina: 12,
		},
		50: stats.Stats{
			stats.Stamina: 16,
		},
		60: stats.Stats{
			stats.Stamina: 16,
		},
	},
	ScrollOfStrength: {
		25: stats.Stats{
			stats.Strength: 9,
		},
		40: stats.Stats{
			stats.Strength: 13,
		},
		50: stats.Stats{
			stats.Strength: 13,
		},
		60: stats.Stats{
			stats.Strength: 17,
		},
	},
	ScrollOfProtection: {
		25: stats.Stats{
			stats.BonusArmor: 120,
		},
		40: stats.Stats{
			stats.BonusArmor: 180,
		},
		50: stats.Stats{
			stats.BonusArmor: 240,
		},
		60: stats.Stats{
			stats.BonusArmor: 240,
		},
	},
}

// Applies buffs that affect individual players.
// TODO: Classic Maximum buff based on character level
func applyBuffEffects(agent Agent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	character := agent.GetCharacter()
	level := character.Level
	bonusResist := float64(0)

	if raidBuffs.ArcaneBrilliance {
		character.AddStats(BuffSpellByLevel[ArcaneIntellect][level])
	} else if raidBuffs.ScrollOfIntellect {
		character.AddStats(BuffSpellByLevel[ScrollOfIntellect][level])
	}

	if raidBuffs.GiftOfTheWild > 0 {
		updateStats := BuffSpellByLevel[MarkOfTheWild][level]
		if raidBuffs.GiftOfTheWild == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.35).Floor()
		}
		character.AddStats(updateStats)
		bonusResist = updateStats[stats.NatureResistance]
	}

	if raidBuffs.NatureResistanceTotem {
		updateStats := BuffSpellByLevel[NatureResistanceTotem][level]
		updateStats[stats.NatureResistance] = updateStats[stats.NatureResistance] - bonusResist
		character.AddStats(updateStats)
	} else if raidBuffs.AspectOfTheWild {
		updateStats := BuffSpellByLevel[AspectOfTheWild][level]
		updateStats[stats.NatureResistance] = updateStats[stats.NatureResistance] - bonusResist
		character.AddStats(updateStats)
	}

	if raidBuffs.FireResistanceAura {
		updateStats := BuffSpellByLevel[FireResistanceAura][level]
		updateStats[stats.FireResistance] = updateStats[stats.FireResistance] - bonusResist
		character.AddStats(updateStats)
	} else if raidBuffs.FireResistanceTotem {
		updateStats := BuffSpellByLevel[FireResistanceTotem][level]
		updateStats[stats.FireResistance] = updateStats[stats.FireResistance] - bonusResist
		character.AddStats(updateStats)
	}

	if raidBuffs.FrostResistanceAura {
		updateStats := BuffSpellByLevel[FrostResistanceAura][level]
		updateStats[stats.FrostResistance] = updateStats[stats.FrostResistance] - bonusResist
		character.AddStats(updateStats)
	} else if raidBuffs.FrostResistanceTotem {
		updateStats := BuffSpellByLevel[FrostResistanceTotem][level]
		updateStats[stats.FrostResistance] = updateStats[stats.FrostResistance] - bonusResist
		character.AddStats(updateStats)
	}

	if raidBuffs.Thorns == proto.TristateEffect_TristateEffectImproved {
		ThornsAura(character, 3)
	} else if raidBuffs.Thorns == proto.TristateEffect_TristateEffectRegular {
		ThornsAura(character, 0)
	}

	if raidBuffs.MoonkinAura {
		character.AddStat(stats.SpellCrit, 3*SpellCritRatingPerCritChance)
	}

	if raidBuffs.LeaderOfThePack {
		character.AddStats(stats.Stats{
			stats.MeleeCrit: 3 * CritRatingPerCritChance,
		})
	}

	if raidBuffs.TrueshotAura {
		character.AddStats(BuffSpellByLevel[TrueshotAura][level])
	}

	if raidBuffs.PowerWordFortitude > 0 {
		updateStats := BuffSpellByLevel[PowerWordFortitude][level]
		if raidBuffs.PowerWordFortitude == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.3).Floor()
		}
		character.AddStats(updateStats)
	} else if raidBuffs.ScrollOfStamina {
		character.AddStats(BuffSpellByLevel[ScrollOfStamina][level])
	}

	if raidBuffs.BloodPact > 0 {
		updateStats := BuffSpellByLevel[BloodPact][level]
		if raidBuffs.BloodPact == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.3).Floor()
		}
		character.AddStats(updateStats)
	} else if raidBuffs.CommandingShout {
		updateStats := BuffSpellByLevel[CommandingShout][level]
		character.AddStats(updateStats)
	}

	if raidBuffs.ShadowResistanceAura {
		updateStats := BuffSpellByLevel[ShadowResistanceAura][level]
		updateStats[stats.ShadowResistance] = updateStats[stats.ShadowResistance] - bonusResist
		character.AddStats(updateStats)
	} else if raidBuffs.ShadowProtection {
		updateStats := BuffSpellByLevel[ShadowProtection][level]
		updateStats[stats.ShadowResistance] = updateStats[stats.ShadowResistance] - bonusResist
		character.AddStats(updateStats)
	}

	if raidBuffs.DivineSpirit {
		character.AddStats(BuffSpellByLevel[DivineSpirit][level])
	} else if raidBuffs.ScrollOfSpirit {
		character.AddStats(BuffSpellByLevel[ScrollOfSpirit][level])
	}

	// Heart of the Lion grants bonus Melee AP as well so give it priority over kings
	if raidBuffs.AspectOfTheLion {
		HeartOfTheLionAura(character)
	} else if individualBuffs.BlessingOfKings {
		MakePermanent(BlessingOfKingsAura(character))
	}

	if raidBuffs.SanctityAura {
		character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1.1
	}

	// TODO: Classic
	/*	if individualBuffs.BlessingOfSanctuary {
			character.PseudoStats.DamageTakenMultiplier *= 0.97
			BlessingOfSanctuaryAura(character)
		}
	*/

	if raidBuffs.DevotionAura == proto.TristateEffect_TristateEffectImproved {
		MakePermanent(DevotionAuraAura(&character.Unit, 2))
	} else if raidBuffs.DevotionAura == proto.TristateEffect_TristateEffectRegular {
		MakePermanent(DevotionAuraAura(&character.Unit, 0))
	}

	if raidBuffs.StoneskinTotem == proto.TristateEffect_TristateEffectImproved {
		MakePermanent(StoneskinTotemAura(&character.Unit, 2))
	} else if raidBuffs.StoneskinTotem == proto.TristateEffect_TristateEffectRegular {
		MakePermanent(StoneskinTotemAura(&character.Unit, 0))
	}

	if raidBuffs.ImprovedStoneskinWindwall {
		MakePermanent(ImprovedStoneskinTotemAura(&character.Unit))
		MakePermanent(ImprovedWindwallTotemAura(&character.Unit))
	}

	if raidBuffs.RetributionAura == proto.TristateEffect_TristateEffectImproved {
		RetributionAura(character, 2)
	} else if raidBuffs.RetributionAura == proto.TristateEffect_TristateEffectRegular {
		RetributionAura(character, 0)
	}

	if raidBuffs.BattleShout != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(BattleShoutAura(&character.Unit, GetTristateValueInt32(raidBuffs.BattleShout, 0, 5), 0))
	}
	if individualBuffs.BlessingOfMight != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(BlessingOfMightAura(&character.Unit, GetTristateValueInt32(individualBuffs.BlessingOfMight, 0, 5), level))
	}

	if raidBuffs.DemonicPact > 0 {
		power := float64(raidBuffs.DemonicPact)
		dpAura := DemonicPactAura(&character.Unit, power, CharacterBuildPhaseBuffs)
		dpAura.ExclusiveEffects[0].Priority = float64(power)
		MakePermanent(dpAura)
	}

	if raidBuffs.StrengthOfEarthTotem != proto.TristateEffect_TristateEffectMissing {
		multiplier := TernaryFloat64(raidBuffs.StrengthOfEarthTotem == proto.TristateEffect_TristateEffectImproved, 1.15, 1)
		MakePermanent(StrengthOfEarthTotemAura(&character.Unit, level, multiplier))
	}

	if raidBuffs.GraceOfAirTotem > 0 {
		multiplier := TernaryFloat64(raidBuffs.GraceOfAirTotem == proto.TristateEffect_TristateEffectImproved, 1.15, 1)
		MakePermanent(GraceOfAirTotemAura(&character.Unit, level, multiplier))
	}

	if individualBuffs.BlessingOfWisdom > 0 {
		updateStats := BuffSpellByLevel[BlessingOfWisdom][level]
		if individualBuffs.BlessingOfWisdom == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.2)
		}
		character.AddStats(updateStats)
	} else if raidBuffs.ManaSpringTotem > 0 {
		updateStats := BuffSpellByLevel[ManaSpring][level]
		if raidBuffs.ManaSpringTotem == proto.TristateEffect_TristateEffectImproved {
			updateStats = updateStats.Multiply(1.25)
		}
		character.AddStats(updateStats)
	}

	if raidBuffs.VampiricTouch > 0 {
		mp5 := float64(raidBuffs.VampiricTouch)
		MakePermanent(VampiricTouchMP5Aura(&character.Unit, mp5))
	}

	if raidBuffs.BattleSquawk > 0 {
		numBattleSquawks := raidBuffs.BattleSquawk
		BattleSquawkAura(&character.Unit, numBattleSquawks)
	}

	// World Buffs
	if individualBuffs.RallyingCryOfTheDragonslayer {
		ApplyRallyingCryOfTheDragonslayer(character)
	}

	if individualBuffs.SpiritOfZandalar {
		ApplySpiritOfZandalar(character)
	}

	if individualBuffs.SongflowerSerenade {
		ApplySongflowerSerenate(character)
	}

	if individualBuffs.WarchiefsBlessing {
		ApplyWarchiefsBlessing(character)
	} else if individualBuffs.MightOfStormwind {
		ApplyMightOfStormwind(character)
	}

	// Dire Maul Buffs
	if individualBuffs.FengusFerocity {
		ApplyFengusFerocity(character)
	}

	if individualBuffs.MoldarsMoxie {
		ApplyMoldarsMoxie(character)
	}

	if individualBuffs.SlipkiksSavvy {
		ApplySlipkiksSavvy(character)
	}

	// Darkmoon Faire Buffs
	if individualBuffs.SaygesFortune != proto.SaygesFortune_SaygesUnknown {
		ApplySaygesFortunes(character, individualBuffs.SaygesFortune)
	}

	// SoD World Buffs
	if individualBuffs.FervorOfTheTempleExplorer {
		ApplyFervorOfTheTempleExplorer(character)
	}

	if individualBuffs.SparkOfInspiration {
		ApplySparkOfInspiration(character)
	}

	if individualBuffs.BoonOfBlackfathom {
		ApplyBoonOfBlackfathom(character)
	}

	if individualBuffs.AshenvalePvpBuff {
		ApplyAshenvaleRallyingCry(character)
	}

	// TODO: Classic provide in APL?
	registerPowerInfusionCD(agent, individualBuffs.PowerInfusions)
	registerManaTideTotemCD(agent, partyBuffs.ManaTideTotems)
	registerInnervateCD(agent, individualBuffs.Innervates)

	character.AddStats(stats.Stats{
		stats.SpellCrit: 2 * SpellCritRatingPerCritChance * float64(partyBuffs.AtieshMage),
	})
	character.AddStats(stats.Stats{
		stats.SpellPower: 33 * float64(partyBuffs.AtieshWarlock),
	})
}

// Applies buffs to pets.
func applyPetBuffEffects(petAgent PetAgent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	// Summoned pets, like Mage Water Elemental, aren't around to receive raid buffs.
	if petAgent.GetPet().IsGuardian() {
		return
	}

	raidBuffs = googleProto.Clone(raidBuffs).(*proto.RaidBuffs)
	partyBuffs = googleProto.Clone(partyBuffs).(*proto.PartyBuffs)
	individualBuffs = googleProto.Clone(individualBuffs).(*proto.IndividualBuffs)

	// We need to modify the buffs a bit because some things are applied to pets by
	// the owner during combat or don't make sense for a pet.
	individualBuffs.Innervates = 0
	individualBuffs.PowerInfusions = 0

	if !petAgent.GetPet().enabledOnStart {
		raidBuffs.ArcaneBrilliance = false
		raidBuffs.DivineSpirit = false
		raidBuffs.GiftOfTheWild = 0
		raidBuffs.PowerWordFortitude = 0
		raidBuffs.Thorns = 0
		raidBuffs.ShadowProtection = false
		raidBuffs.ScrollOfProtection = false
		raidBuffs.ScrollOfStamina = false
		raidBuffs.ScrollOfStrength = false
		raidBuffs.ScrollOfAgility = false
		raidBuffs.ScrollOfIntellect = false
		raidBuffs.ScrollOfSpirit = false
		individualBuffs.BlessingOfKings = false
		individualBuffs.BlessingOfSanctuary = false
		individualBuffs.BlessingOfMight = 0
		individualBuffs.BlessingOfWisdom = 0
	}

	// Pets no longer get world buffs
	individualBuffs.BoonOfBlackfathom = false
	individualBuffs.SparkOfInspiration = false
	individualBuffs.FervorOfTheTempleExplorer = false
	individualBuffs.AshenvalePvpBuff = false
	individualBuffs.SongflowerSerenade = false
	individualBuffs.RallyingCryOfTheDragonslayer = false
	individualBuffs.WarchiefsBlessing = false
	individualBuffs.MightOfStormwind = false
	individualBuffs.SpiritOfZandalar = false
	individualBuffs.FengusFerocity = false
	individualBuffs.MoldarsMoxie = false
	individualBuffs.SlipkiksSavvy = false
	individualBuffs.SaygesFortune = proto.SaygesFortune_SaygesUnknown

	applyBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
}

func BlessingOfKingsAura(character *Character) *Aura {
	character.MultiplyStat(stats.Strength, 1.1)
	character.MultiplyStat(stats.Stamina, 1.1)
	character.MultiplyStat(stats.Agility, 1.1)
	character.MultiplyStat(stats.Intellect, 1.1)
	character.MultiplyStat(stats.Spirit, 1.1)

	return character.RegisterAura(Aura{
		Label:    "Blessing of Kings",
		ActionID: ActionID{SpellID: 20217},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func HeartOfTheLionAura(character *Character) *Aura {
	character.AddStat(stats.AttackPower, float64(40+4*(character.Level-20)))
	character.MultiplyStat(stats.Strength, 1.1)
	character.MultiplyStat(stats.Stamina, 1.1)
	character.MultiplyStat(stats.Agility, 1.1)
	character.MultiplyStat(stats.Intellect, 1.1)
	character.MultiplyStat(stats.Spirit, 1.1)

	return character.RegisterAura(Aura{
		Label:    "Heart of the Lion",
		ActionID: ActionID{SpellID: 409583},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

// TODO: Classic
func InspirationAura(unit *Unit, points int32) *Aura {
	multiplier := 1 - []float64{0, .03, .07, .10}[points]

	return unit.GetOrRegisterAura(Aura{
		Label:    "Inspiration",
		ActionID: ActionID{SpellID: 15363},
		Duration: time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= multiplier
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= multiplier
		},
	})
}

func ApplyInspiration(character *Character, uptime float64) {
	if uptime <= 0 {
		return
	}
	uptime = min(1, uptime)

	inspirationAura := InspirationAura(&character.Unit, 3)

	ApplyFixedUptimeAura(inspirationAura, uptime, time.Millisecond*2500, 1)
}

func DevotionAuraAura(unit *Unit, points int32) *Aura {
	level := unit.Level
	spellID := map[int32]int32{
		25: 643,
		40: 1032,
		50: 10292,
		60: 10293,
	}[level]

	updateStats := BuffSpellByLevel[DevotionAura][level]
	updateStats = updateStats.Multiply(1 + .125*float64(points))

	return unit.RegisterAura(Aura{
		Label:    "Devotion Aura",
		ActionID: ActionID{SpellID: spellID},
		Duration: NeverExpires,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, updateStats)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, updateStats.Multiply(-1))
		},
	})
}

func StoneskinTotemAura(unit *Unit, points int32) *Aura {
	level := unit.Level
	spellID := map[int32]int32{
		25: 8155,
		40: 10406,
		50: 10407,
		60: 10408,
	}[level]
	meleeDamageReduction := map[int32]float64{
		25: -11,
		40: -16,
		50: -22,
		60: -30,
	}[level]
	meleeDamageReduction *= 1 + .1*float64(points)
	meleeDamageReduction = math.Floor(meleeDamageReduction)

	return unit.GetOrRegisterAura(Aura{
		Label:    "Stoneskin",
		ActionID: ActionID{SpellID: spellID},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusDamageTakenAfterModifiers[DefenseTypeMelee] += meleeDamageReduction
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusDamageTakenAfterModifiers[DefenseTypeMelee] += meleeDamageReduction
		},
	})
}

// https://www.wowhead.com/classic/spell=457544/s03-item-t1-shaman-tank-6p-bonus
// Your Stoneskin Totem also reduces Physical damage taken by 5% and your Windwall Totem also reduces Magical damage taken by 5%.
// Restricting to level 60 only
func ImprovedStoneskinTotemAura(unit *Unit) *Aura {
	return unit.GetOrRegisterAura(Aura{
		Label:    "Improved Stoneskin",
		ActionID: ActionID{SpellID: 457544}.WithTag(1),
		Duration: time.Minute * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= .95
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= .95
		},
	})
}
func ImprovedWindwallTotemAura(unit *Unit) *Aura {
	return unit.GetOrRegisterAura(Aura{
		Label:    "Improved Windwall",
		ActionID: ActionID{SpellID: 457544}.WithTag(2),
		Duration: time.Minute * 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= .95
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= .95
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= .95
		},
	})
}

func RetributionAura(character *Character, points int32) *Aura {
	level := character.Level
	spellID := map[int32]int32{
		25: 7294,
		40: 10299,
		50: 10300,
		60: 10301,
	}[level]

	baseDamage := map[int32]int32{
		25: 5,
		40: 12,
		50: 16,
		60: 20,
	}[level]

	actionID := ActionID{SpellID: spellID}

	damage := float64(baseDamage) * (1 + 0.25*float64(points))

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolHoly,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagBinary,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHit)
		},
	})

	return character.RegisterAura(Aura{
		Label:    "Retribution Aura",
		ActionID: actionID,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(ProcMaskMelee) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})
}

func ThornsAura(character *Character, points int32) *Aura {
	level := character.Level
	spellID := map[int32]int32{
		25: 1075,
		40: 8914,
		50: 9756,
		60: 9910,
	}[level]

	baseDamage := map[int32]int32{
		25: 9,
		40: 12,
		50: 15,
		60: 18,
	}[level]

	actionID := ActionID{SpellID: spellID}
	damage := float64(baseDamage) * (1 + 0.25*float64(points))

	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolNature,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagBinary,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHit)
		},
	})

	return MakePermanent(character.RegisterAura(Aura{
		Label:    "Thorns",
		ActionID: actionID,
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Landed() && spell.ProcMask.Matches(ProcMaskMelee) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	}))
}

// func BlessingOfSanctuaryAura(character *Character, level int32) {
// 	if character.Level < 30 {
// 		return
// 	}

// 	spellID := map[int32]int32{
// 		40: 20912,
// 		50: 20913,
// 		60: 20914,
// 	}[level]

// 	physReduction := map[int32]int32{
// 		40: 14,
// 		50: 19,
// 		60: 24,
// 	}[level]

// 	blockDamage := map[int32]int32{
// 		40: 21,
// 		50: 28,
// 		60: 35,
// 	}

// 	actionID := ActionID{SpellID: spellID}

// 	character.RegisterAura(Aura{
// 		Label:    "Blessing of Sanctuary",
// 		ActionID: actionID,
// 		Duration: NeverExpires,
// 		OnReset: func(aura *Aura, sim *Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
// 			if result.Outcome.Matches(OutcomeBlock | OutcomeDodge | OutcomeParry) {
// 			}
// 		},
// 	})
// }

// Used for approximating cooldowns applied by other players to you, such as
// bloodlust, innervate, power infusion, etc. This is specifically for buffs
// which can be consecutively applied multiple times to a single player.
type externalConsecutiveCDApproximation struct {
	ActionID         ActionID
	AuraTag          string
	CooldownPriority int32
	Type             CooldownType
	AuraDuration     time.Duration
	AuraCD           time.Duration

	// Callback for extra activation conditions.
	ShouldActivate CooldownActivationCondition

	// Applies the buff.
	AddAura CooldownActivation
}

// numSources is the number of other players assigned to apply the buff to this player.
// E.g. the number of other shaman in the group using bloodlust.
func registerExternalConsecutiveCDApproximation(agent Agent, config externalConsecutiveCDApproximation, numSources int32) {
	if numSources == 0 {
		panic("Need at least 1 source!")
	}
	character := agent.GetCharacter()

	var nextExternalIndex int

	externalTimers := make([]*Timer, numSources)
	for i := 0; i < int(numSources); i++ {
		externalTimers[i] = character.NewTimer()
	}
	sharedTimer := character.NewTimer()

	spell := character.RegisterSpell(SpellConfig{
		ActionID: config.ActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    sharedTimer,
				Duration: config.AuraDuration, // Assumes that multiple buffs are different sources.
			},
		},
		ExtraCastCondition: func(sim *Simulation, target *Unit) bool {
			if !externalTimers[nextExternalIndex].IsReady(sim) {
				return false
			}

			if character.HasActiveAuraWithTag(config.AuraTag) {
				return false
			}

			return true
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			config.AddAura(sim, character)
			externalTimers[nextExternalIndex].Set(sim.CurrentTime + config.AuraCD)

			nextExternalIndex = (nextExternalIndex + 1) % len(externalTimers)

			if externalTimers[nextExternalIndex].IsReady(sim) {
				sharedTimer.Set(sim.CurrentTime + config.AuraDuration)
			} else {
				sharedTimer.Set(sim.CurrentTime + externalTimers[nextExternalIndex].TimeToReady(sim))
			}
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: config.CooldownPriority,
		Type:     config.Type,

		ShouldActivate: config.ShouldActivate,
	})
}

// var BloodlustActionID = ActionID{SpellID: 2825}

// const SatedAuraLabel = "Sated"
// const BloodlustAuraTag = "Bloodlust"
// const BloodlustDuration = time.Second * 40
// const BloodlustCD = time.Minute * 10

// func registerBloodlustCD(agent Agent) {
// 	character := agent.GetCharacter()
// 	bloodlustAura := BloodlustAura(character, -1)

// 	spell := character.RegisterSpell(SpellConfig{
// 		ActionID: bloodlustAura.ActionID,
// 		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

// 		Cast: CastConfig{
// 			CD: Cooldown{
// 				Timer:    character.NewTimer(),
// 				Duration: BloodlustCD,
// 			},
// 		},

// 		ApplyEffects: func(sim *Simulation, target *Unit, _ *Spell) {
// 			if !target.HasActiveAura(SatedAuraLabel) {
// 				bloodlustAura.Activate(sim)
// 			}
// 		},
// 	})

// 	character.AddMajorCooldown(MajorCooldown{
// 		Spell:    spell,
// 		Priority: CooldownPriorityBloodlust,
// 		Type:     CooldownTypeDPS,
// 		ShouldActivate: func(sim *Simulation, character *Character) bool {
// 			// Haste portion doesn't stack with Power Infusion, so prefer to wait.
// 			return !character.HasActiveAuraWithTag(PowerInfusionAuraTag) && !character.HasActiveAura(SatedAuraLabel)
// 		},
// 	})
// }

// func BloodlustAura(character *Character, actionTag int32) *Aura {
// 	actionID := BloodlustActionID.WithTag(actionTag)

// 	sated := character.GetOrRegisterAura(Aura{
// 		Label:    SatedAuraLabel,
// 		ActionID: ActionID{SpellID: 57724},
// 		Duration: time.Minute * 10,
// 	})

// 	aura := character.GetOrRegisterAura(Aura{
// 		Label:    "Bloodlust-" + actionID.String(),
// 		Tag:      BloodlustAuraTag,
// 		ActionID: actionID,
// 		Duration: BloodlustDuration,
// 		OnGain: func(aura *Aura, sim *Simulation) {
// 			character.MultiplyAttackSpeed(sim, 1.3)
// 			for _, pet := range character.Pets {
// 				if pet.IsEnabled() && !pet.IsGuardian() {
// 					BloodlustAura(&pet.Character, actionTag).Activate(sim)
// 				}
// 			}

// 			if character.HasActiveAura(SatedAuraLabel) {
// 				aura.Deactivate(sim) // immediately remove it person already has sated.
// 				return
// 			}
// 		},
// 		OnExpire: func(aura *Aura, sim *Simulation) {
// 			character.MultiplyAttackSpeed(sim, 1.0/1.3)
// 			sated.Activate(sim)
// 		},
// 	})
// 	multiplyCastSpeedEffect(aura, 1.3)
// 	return aura
// }

var PowerInfusionActionID = ActionID{SpellID: 10060}
var PowerInfusionAuraTag = "PowerInfusion"

const PowerInfusionDuration = time.Second * 15
const PowerInfusionCD = time.Minute * 3

func registerPowerInfusionCD(agent Agent, numPowerInfusions int32) {
	if numPowerInfusions == 0 {
		return
	}

	piAura := PowerInfusionAura(&agent.GetCharacter().Unit, -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         PowerInfusionActionID.WithTag(-1),
			AuraTag:          PowerInfusionAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     PowerInfusionDuration,
			AuraCD:           PowerInfusionCD,
			Type:             CooldownTypeDPS,

			AddAura: func(sim *Simulation, character *Character) { piAura.Activate(sim) },
		},
		numPowerInfusions)
}

func PowerInfusionAura(character *Unit, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 10060, Tag: actionTag}
	aura := character.GetOrRegisterAura(Aura{
		Label:    "PowerInfusion-" + actionID.String(),
		Tag:      PowerInfusionAuraTag,
		ActionID: actionID,
		Duration: PowerInfusionDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.2
			//character.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] /= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] /= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] /= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] /= 1.2
			character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.2
			//character.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})
	return aura
}

// func multiplyCastSpeedEffect(aura *Aura, multiplier float64) *ExclusiveEffect {
// 	return aura.NewExclusiveEffect("MultiplyCastSpeed", false, ExclusiveEffect{
// 		Priority: multiplier,
// 		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
// 			ee.Aura.Unit.MultiplyCastSpeed(multiplier)
// 		},
// 		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
// 			ee.Aura.Unit.MultiplyCastSpeed(1 / multiplier)
// 		},
// 	})
// }

// var TricksOfTheTradeAuraTag = "TricksOfTheTrade"

// const TricksOfTheTradeCD = time.Second * 3600 // CD is 30s from the time buff ends (so 40s with glyph) but that's in order to be able to set the number of TotT you'll have during the fight

// func registerTricksOfTheTradeCD(agent Agent, numTricksOfTheTrades int32) {
// 	if numTricksOfTheTrades == 0 {
// 		return
// 	}

// 	TotTAura := TricksOfTheTradeAura(&agent.GetCharacter().Unit, -1)

// 	registerExternalConsecutiveCDApproximation(
// 		agent,
// 		externalConsecutiveCDApproximation{
// 			ActionID:         ActionID{SpellID: 57933, Tag: -1},
// 			AuraTag:          TricksOfTheTradeAuraTag,
// 			CooldownPriority: CooldownPriorityDefault,
// 			AuraDuration:     TotTAura.Duration,
// 			AuraCD:           TricksOfTheTradeCD,
// 			Type:             CooldownTypeDPS,

// 			ShouldActivate: func(sim *Simulation, character *Character) bool {
// 				return !agent.GetCharacter().GetExclusiveEffectCategory("PercentDamageModifier").AnyActive()
// 			},
// 			AddAura: func(sim *Simulation, character *Character) { TotTAura.Activate(sim) },
// 		},
// 		numTricksOfTheTrades)
// }

// func TricksOfTheTradeAura(character *Unit, actionTag int32) *Aura {
// 	actionID := ActionID{SpellID: 57933, Tag: actionTag}

// 	aura := character.GetOrRegisterAura(Aura{
// 		Label:    "TricksOfTheTrade-" + actionID.String(),
// 		Tag:      TricksOfTheTradeAuraTag,
// 		ActionID: actionID,
// 		Duration: time.Second * 6,
// 		OnGain: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.DamageDealtMultiplier *= 1.15
// 		},
// 		OnExpire: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.DamageDealtMultiplier /= 1.15
// 		},
// 	})

// 	RegisterPercentDamageModifierEffect(aura, 1.15)
// 	return aura
// }

func RegisterPercentDamageModifierEffect(aura *Aura, percentDamageModifier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("PercentDamageModifier", true, ExclusiveEffect{
		Priority: percentDamageModifier,
	})
}

// var DivineGuardianAuraTag = "DivineGuardian"

// const DivineGuardianDuration = time.Second * 6
// const DivineGuardianCD = time.Minute * 2

// var HandOfSacrificeAuraTag = "HandOfSacrifice"

// const HandOfSacrificeDuration = time.Millisecond * 10500 // subtract Divine Shield GCD
// const HandOfSacrificeCD = time.Minute * 5                // use Divine Shield CD here

// func registerHandOfSacrificeCD(agent Agent, numSacs int32) {
// 	if numSacs == 0 {
// 		return
// 	}

// 	hosAura := HandOfSacrificeAura(agent.GetCharacter(), -1)

// 	registerExternalConsecutiveCDApproximation(
// 		agent,
// 		externalConsecutiveCDApproximation{
// 			ActionID:         ActionID{SpellID: 6940, Tag: -1},
// 			AuraTag:          HandOfSacrificeAuraTag,
// 			CooldownPriority: CooldownPriorityLow,
// 			AuraDuration:     HandOfSacrificeDuration,
// 			AuraCD:           HandOfSacrificeCD,
// 			Type:             CooldownTypeSurvival,

// 			ShouldActivate: func(sim *Simulation, character *Character) bool {
// 				return true
// 			},
// 			AddAura: func(sim *Simulation, character *Character) {
// 				hosAura.Activate(sim)
// 			},
// 		},
// 		numSacs)
// }

// func HandOfSacrificeAura(character *Character, actionTag int32) *Aura {
// 	actionID := ActionID{SpellID: 6940, Tag: actionTag}

// 	return character.GetOrRegisterAura(Aura{
// 		Label:    "HandOfSacrifice-" + actionID.String(),
// 		Tag:      HandOfSacrificeAuraTag,
// 		ActionID: actionID,
// 		Duration: HandOfSacrificeDuration,
// 		OnGain: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.DamageTakenMultiplier *= 0.7
// 		},
// 		OnExpire: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.DamageTakenMultiplier /= 0.7
// 		},
// 	})
// }

// var PainSuppressionAuraTag = "PainSuppression"

// const PainSuppressionDuration = time.Second * 8
// const PainSuppressionCD = time.Minute * 3

// func registerPainSuppressionCD(agent Agent, numPainSuppressions int32) {
// 	if numPainSuppressions == 0 {
// 		return
// 	}

// 	psAura := PainSuppressionAura(agent.GetCharacter(), -1)

// 	registerExternalConsecutiveCDApproximation(
// 		agent,
// 		externalConsecutiveCDApproximation{
// 			ActionID:         ActionID{SpellID: 33206, Tag: -1},
// 			AuraTag:          PainSuppressionAuraTag,
// 			CooldownPriority: CooldownPriorityDefault,
// 			AuraDuration:     PainSuppressionDuration,
// 			AuraCD:           PainSuppressionCD,
// 			Type:             CooldownTypeSurvival,

// 			ShouldActivate: func(sim *Simulation, character *Character) bool {
// 				return true
// 			},
// 			AddAura: func(sim *Simulation, character *Character) { psAura.Activate(sim) },
// 		},
// 		numPainSuppressions)
// }

// func PainSuppressionAura(character *Character, actionTag int32) *Aura {
// 	actionID := ActionID{SpellID: 33206, Tag: actionTag}

// 	return character.GetOrRegisterAura(Aura{
// 		Label:    "PainSuppression-" + actionID.String(),
// 		Tag:      PainSuppressionAuraTag,
// 		ActionID: actionID,
// 		Duration: PainSuppressionDuration,
// 		OnGain: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.DamageTakenMultiplier *= 0.6
// 		},
// 		OnExpire: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.DamageTakenMultiplier /= 0.6
// 		},
// 	})
// }

// var GuardianSpiritAuraTag = "GuardianSpirit"

// const GuardianSpiritDuration = time.Second * 10
// const GuardianSpiritCD = time.Minute * 3

// func registerGuardianSpiritCD(agent Agent, numGuardianSpirits int32) {
// 	if numGuardianSpirits == 0 {
// 		return
// 	}

// 	character := agent.GetCharacter()
// 	gsAura := GuardianSpiritAura(character, -1)
// 	healthMetrics := character.NewHealthMetrics(ActionID{SpellID: 47788})

// 	character.AddDynamicDamageTakenModifier(func(sim *Simulation, _ *Spell, result *SpellResult) {
// 		if (result.Damage >= character.CurrentHealth()) && gsAura.IsActive() {
// 			result.Damage = character.CurrentHealth()
// 			character.GainHealth(sim, 0.5*character.MaxHealth(), healthMetrics)
// 			gsAura.Deactivate(sim)
// 		}
// 	})

// 	registerExternalConsecutiveCDApproximation(
// 		agent,
// 		externalConsecutiveCDApproximation{
// 			ActionID:         ActionID{SpellID: 47788, Tag: -1},
// 			AuraTag:          GuardianSpiritAuraTag,
// 			CooldownPriority: CooldownPriorityLow,
// 			AuraDuration:     GuardianSpiritDuration,
// 			AuraCD:           GuardianSpiritCD,
// 			Type:             CooldownTypeSurvival,

// 			ShouldActivate: func(sim *Simulation, character *Character) bool {
// 				return true
// 			},
// 			AddAura: func(sim *Simulation, character *Character) {
// 				gsAura.Activate(sim)
// 			},
// 		},
// 		numGuardianSpirits)
// }

// func GuardianSpiritAura(character *Character, actionTag int32) *Aura {
// 	actionID := ActionID{SpellID: 47788, Tag: actionTag}

// 	return character.GetOrRegisterAura(Aura{
// 		Label:    "GuardianSpirit-" + actionID.String(),
// 		Tag:      GuardianSpiritAuraTag,
// 		ActionID: actionID,
// 		Duration: GuardianSpiritDuration,
// 		OnGain: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.HealingTakenMultiplier *= 1.4
// 		},
// 		OnExpire: func(aura *Aura, sim *Simulation) {
// 			character.PseudoStats.HealingTakenMultiplier /= 1.4
// 		},
// 	})
// }

// func registerRevitalizeHotCD(agent Agent, label string, hotID ActionID, ticks int, tickPeriod time.Duration, uptimeCount int32) {
// 	if uptimeCount == 0 {
// 		return
// 	}

// 	character := agent.GetCharacter()
// 	revActionID := ActionID{SpellID: 48545}

// 	manaMetrics := character.NewManaMetrics(revActionID)
// 	energyMetrics := character.NewEnergyMetrics(revActionID)
// 	rageMetrics := character.NewRageMetrics(revActionID)

// 	// Calculate desired downtime based on selected uptimeCount (1 count = 10% uptime, 0%-100%)
// 	totalDuration := time.Duration(ticks) * tickPeriod
// 	uptimePercent := float64(uptimeCount) / 100.0

// 	aura := character.GetOrRegisterAura(Aura{
// 		Label:    "Revitalize-" + label,
// 		ActionID: hotID,
// 		Duration: totalDuration,
// 		OnGain: func(aura *Aura, sim *Simulation) {
// 			pa := NewPeriodicAction(sim, PeriodicActionOptions{
// 				Period:   tickPeriod,
// 				NumTicks: ticks,
// 				OnAction: func(s *Simulation) {
// 					if s.RandomFloat("Revitalize Proc") < 0.15 {
// 						cpb := aura.Unit.GetCurrentPowerBar()
// 						if cpb == ManaBar {
// 							aura.Unit.AddMana(s, 0.01*aura.Unit.MaxMana(), manaMetrics)
// 						} else if cpb == EnergyBar {
// 							aura.Unit.AddEnergy(s, 8, energyMetrics)
// 						} else if cpb == RageBar {
// 							aura.Unit.AddRage(s, 4, rageMetrics)
// 						}
// 					}
// 				},
// 			})
// 			sim.AddPendingAction(pa)
// 		},
// 	})

// 	ApplyFixedUptimeAura(aura, uptimePercent, totalDuration, 1)
// }

const ShatteringThrowCD = time.Minute * 5

var InnervateAuraTag = "Innervate"

const InnervateDuration = time.Second * 20
const InnervateCD = time.Minute * 6

func InnervateManaThreshold(character *Character) float64 {
	if character.Class == proto.Class_ClassMage {
		// Mages burn mana really fast so they need a higher threshold.
		return character.MaxMana() * 0.7
	} else {
		return 1000
	}
}

func registerInnervateCD(agent Agent, numInnervates int32) {
	if numInnervates == 0 {
		return
	}

	character := agent.GetCharacter()
	innervateThreshold := 0.0
	innervateAura := InnervateAura(character, -1)

	character.Env.RegisterPostFinalizeEffect(func() {
		innervateThreshold = InnervateManaThreshold(character)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 29166, Tag: -1},
			AuraTag:          InnervateAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     InnervateDuration,
			AuraCD:           InnervateCD,
			Type:             CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only cast innervate when very low on mana, to make sure all other mana CDs are prioritized.
				return character.CurrentMana() <= innervateThreshold
			},
			AddAura: func(sim *Simulation, character *Character) {
				innervateAura.Activate(sim)
			},
		},
		numInnervates)
}

func InnervateAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 29166, Tag: actionTag}
	// TODO: Add metrics for increased regen from spirit (either add here and align ticks to mana tick or create mana tick hook?)
	// manaMetrics := character.NewManaMetrics(actionID)
	return character.GetOrRegisterAura(Aura{
		Label:    "Innervate-" + actionID.String(),
		Tag:      InnervateAuraTag,
		ActionID: actionID,
		Duration: InnervateDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SpiritRegenMultiplier += 4
			character.PseudoStats.ForceFullSpiritRegen = true
			character.UpdateManaRegenRates()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SpiritRegenMultiplier -= 4
			character.PseudoStats.ForceFullSpiritRegen = false
			character.UpdateManaRegenRates()
		},
	})
}

var ManaTideTotemActionID = ActionID{SpellID: 16190}
var ManaTideTotemAuraTag = "ManaTideTotem"

const ManaTideTotemDuration = time.Second * 12
const ManaTideTotemCD = time.Minute * 5

func registerManaTideTotemCD(agent Agent, numManaTideTotems int32) {
	if numManaTideTotems == 0 {
		return
	}

	character := agent.GetCharacter()
	initialDelay := time.Duration(0)
	mttAura := ManaTideTotemAura(character, -1)

	character.Env.RegisterPostFinalizeEffect(func() {
		// Use first MTT at 60s, or halfway through the fight, whichever comes first.
		initialDelay = min(character.Env.BaseDuration/2, time.Second*60)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ManaTideTotemActionID.WithTag(-1),
			AuraTag:          ManaTideTotemAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     ManaTideTotemDuration,
			AuraCD:           ManaTideTotemCD,
			Type:             CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// A normal resto shaman would wait to use MTT.
				return sim.CurrentTime >= initialDelay
			},
			AddAura: func(sim *Simulation, character *Character) {
				mttAura.Activate(sim)
			},
		},
		numManaTideTotems)
}

func ManaTideTotemAura(character *Character, actionTag int32) *Aura {
	actionID := ManaTideTotemActionID.WithTag(actionTag)

	metrics := make([]*ResourceMetrics, len(character.Party.Players))
	for i, player := range character.Party.Players {
		char := player.GetCharacter()
		if char.HasManaBar() {
			metrics[i] = char.NewManaMetrics(actionID)
		}
	}

	manaPerTick := map[int32]float64{
		25: 0,
		40: 170, // Rank 1
		50: 230, // Rank 2
		60: 290, // Rank 3
	}[character.Level]

	return character.GetOrRegisterAura(Aura{
		Label:    "ManaTideTotem-" + actionID.String(),
		Tag:      ManaTideTotemAuraTag,
		ActionID: actionID,
		Duration: ManaTideTotemDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period:   ManaTideTotemDuration / 4,
				NumTicks: 4,
				OnAction: func(sim *Simulation) {
					for i, player := range character.Party.Players {
						if metrics[i] != nil {
							char := player.GetCharacter()
							char.AddMana(sim, manaPerTick, metrics[i])
						}
					}
				},
			})
		},
	})
}

const ReplenishmentAuraDuration = time.Second * 15

// Creates the actual replenishment aura for a unit.
// func replenishmentAura(unit *Unit, _ ActionID) *Aura {
// 	if unit.ReplenishmentAura != nil {
// 		return unit.ReplenishmentAura
// 	}

// 	replenishmentDep := unit.NewDynamicStatDependency(stats.Mana, stats.MP5, 0.01)

// 	unit.ReplenishmentAura = unit.RegisterAura(Aura{
// 		Label:    "Replenishment",
// 		ActionID: ActionID{SpellID: 57669},
// 		Duration: ReplenishmentAuraDuration,
// 		OnGain: func(aura *Aura, sim *Simulation) {
// 			aura.Unit.EnableDynamicStatDep(sim, replenishmentDep)
// 		},
// 		OnExpire: func(aura *Aura, sim *Simulation) {
// 			aura.Unit.DisableDynamicStatDep(sim, replenishmentDep)
// 		},
// 	})

// 	return unit.ReplenishmentAura
// }

func DemonicPactAura(unit *Unit, spellpower float64, buildPhase CharacterBuildPhase) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Demonic Pact",
		ActionID:   ActionID{SpellID: 425464},
		Duration:   time.Second * 45,
		BuildPhase: buildPhase,
	})
	spellPowerBonusEffect(aura, spellpower)
	return aura
}

func spellPowerBonusEffect(aura *Aura, spellPowerBonus float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("SpellPowerBonus", false, ExclusiveEffect{
		Priority: spellPowerBonus,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.SpellPower: ee.Priority,
			})
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.SpellPower: -ee.Priority,
			})
		},
	})
}

func StrengthOfEarthTotemAura(unit *Unit, level int32, multiplier float64) *Aura {
	spellId := map[int32]int32{
		25: 8162,
		40: 8163,
		50: 8163,
		60: 25362,
	}[level]
	duration := time.Minute * 2
	updateStats := BuffSpellByLevel[StrengthOfEarth][level].Multiply(multiplier).Floor()

	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Strength of Earth Totem",
		ActionID:   ActionID{SpellID: spellId},
		Duration:   duration,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			unit.AddStatsDynamic(sim, updateStats)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			unit.AddStatsDynamic(sim, updateStats.Multiply(-1))
		},
	})
	return aura
}

func GraceOfAirTotemAura(unit *Unit, level int32, multiplier float64) *Aura {
	spellId := map[int32]int32{
		50: 8835,
		60: 25359,
	}[level]
	duration := time.Minute * 2
	updateStats := BuffSpellByLevel[GraceOfAir][level].Multiply(multiplier).Floor()

	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Grace of Air Totem",
		ActionID:   ActionID{SpellID: spellId},
		Duration:   duration,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			unit.AddStatsDynamic(sim, updateStats)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			unit.AddStatsDynamic(sim, updateStats.Multiply(-1))
		},
	})
	return aura
}

const BattleShoutRanks = 7

var BattleShoutSpellId = [BattleShoutRanks + 1]int32{0, 6673, 5242, 6192, 11549, 11550, 11551, 25289}
var BattleShoutBaseAP = [BattleShoutRanks + 1]float64{0, 20, 40, 57, 93, 138, 193, 232}
var BattleShoutLevel = [BattleShoutRanks + 1]int{0, 1, 12, 22, 32, 42, 52, 60}

func BattleShoutAura(unit *Unit, impBattleShout int32, boomingVoicePts int32) *Aura {
	rank := LevelToBuffRank[BattleShout][unit.Level]
	spellId := BattleShoutSpellId[rank]
	baseAP := BattleShoutBaseAP[rank]

	return unit.GetOrRegisterAura(Aura{
		Label:      "Battle Shout",
		ActionID:   ActionID{SpellID: spellId},
		Duration:   time.Duration(float64(time.Minute*2) * (1 + 0.1*float64(boomingVoicePts))),
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower: math.Floor(baseAP * (1 + 0.05*float64(impBattleShout))),
			})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower: -1 * math.Floor(baseAP*(1+0.05*float64(impBattleShout))),
			})
		},
	})
}

func BlessingOfMightAura(unit *Unit, impBomPts int32, level int32) *Aura {
	spellID := map[int32]int32{
		25: 19835,
		40: 19836,
		50: 19837,
		60: 25291,
	}[level]

	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Blessing of Might",
		ActionID:   ActionID{SpellID: spellID},
		Duration:   NeverExpires,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower: math.Floor(BuffSpellByLevel[BlessingOfMight][level][stats.AttackPower] * (1 + 0.04*float64(impBomPts))),
			})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{
				stats.AttackPower: -1 * math.Floor(BuffSpellByLevel[BlessingOfMight][level][stats.AttackPower]*(1+0.04*float64(impBomPts))),
			})
		},
	})
	return aura
}

// TODO: Are there exclusive AP buffs in SoD?
// func attackPowerBonusEffect(aura *Aura, apBonus float64) *ExclusiveEffect {
// 	return aura.NewExclusiveEffect("AttackPowerBonus", false, ExclusiveEffect{
// 		Priority: apBonus,
// 		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
// 			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
// 				stats.AttackPower:       ee.Priority,
// 				stats.RangedAttackPower: ee.Priority,
// 			})
// 		},
// 		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
// 			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
// 				stats.AttackPower:       -ee.Priority,
// 				stats.RangedAttackPower: -ee.Priority,
// 			})
// 		},
// 	})
// }

func VampiricTouchMP5Aura(unit *Unit, mp5 float64) *Aura {
	actionID := ActionID{SpellID: 402779}.WithTag(1)
	mps := mp5 / 5

	manaMetrics := unit.NewManaMetrics(actionID)
	aura := unit.GetOrRegisterAura(Aura{
		Label:      "Vampiric Touch Replenishment",
		ActionID:   actionID,
		Duration:   NeverExpires,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnReset: func(aura *Aura, sim *Simulation) {
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period:   time.Second * 1,
				Priority: ActionPriorityDOT, // High prio
				OnAction: func(sim *Simulation) {
					unit.AddMana(sim, mps, manaMetrics)
				},
			})
		},
	})

	return aura
}

func BattleSquawkAura(character *Unit, stackcount int32) *Aura {
	aura := character.GetOrRegisterAura(Aura{
		Label:      "Battle Squawk",
		ActionID:   ActionID{SpellID: 23060},
		Duration:   time.Minute * 4,
		MaxStacks:  5,
		BuildPhase: CharacterBuildPhaseBuffs,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, stackcount)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks, newStacks int32) {
			character.MultiplyMeleeSpeed(sim, math.Pow(1.05, float64(newStacks-oldStacks)))
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, 0)
		},
	})
	return aura
}

// func healthBonusEffect(aura *Aura, healthBonus float64) *ExclusiveEffect {
// 	return aura.NewExclusiveEffect("HealthBonus", false, ExclusiveEffect{
// 		Priority: healthBonus,
// 		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
// 			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
// 				stats.Health: ee.Priority,
// 			})
// 		},
// 		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
// 			ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
// 				stats.Health: -ee.Priority,
// 			})
// 		},
// 	})
// }

func ApplyWildStrikes(character *Character) *Aura {
	buffActionID := ActionID{SpellID: 407975}

	var bonusAP float64

	wsBuffAura := character.GetOrRegisterAura(Aura{
		Label:     "Wild Strikes Buff",
		ActionID:  buffActionID,
		Duration:  time.Millisecond * 1500,
		MaxStacks: 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			bonusAP = 0.2 * aura.Unit.GetStat(stats.AttackPower)
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: bonusAP})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: -bonusAP})
		},
	})

	icd := Cooldown{
		Timer:    character.NewTimer(),
		Duration: time.Millisecond * 1500,
	}

	wsBuffAura.Icd = &icd

	MakePermanent(character.GetOrRegisterAura(Aura{
		Label: "Wild Strikes",
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			// charges are removed by every auto or next melee, whether it lands or not
			if wsBuffAura.IsActive() && spell.ProcMask.Matches(ProcMaskMeleeWhiteHit) {
				wsBuffAura.RemoveStack(sim)
			}

			if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMeleeMH) {
				return
			}

			if icd.IsReady(sim) && sim.RandomFloat("Wild Strikes") < 0.2 {
				icd.Use(sim)
				wsBuffAura.Activate(sim)
				// aura is up _after_ the triggering swing lands, so the aura always stays up after the extra attack
				wsBuffAura.SetStacks(sim, 2)
				aura.Unit.AutoAttacks.ExtraMHAttack(sim, 1, buffActionID)
			}
		},
	}))

	return wsBuffAura
}

const WindfuryRanks = 3

var (
	WindfuryBuffLevelToRank = []int32{
		25: 0,
		40: 1,
		50: 2,
		60: 3,
	}
	WindfuryBuffSpellId = [WindfuryRanks + 1]int32{0, 8516, 10608, 10610}
	WindfuryBuffBonusAP = [WindfuryRanks + 1]float64{0, 122, 229, 315}
)

func ApplyWindfury(character *Character) *Aura {
	level := character.Level
	rank := WindfuryBuffLevelToRank[level]
	spellId := WindfuryBuffSpellId[rank]
	bonusAP := WindfuryBuffBonusAP[rank]

	windfuryBuffAura := character.GetOrRegisterAura(Aura{
		Label:     "Windfury Buff",
		ActionID:  ActionID{SpellID: spellId},
		Duration:  time.Millisecond * 1500,
		MaxStacks: 2,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: bonusAP})
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: -bonusAP})
		},
	})

	icd := Cooldown{
		Timer:    character.NewTimer(),
		Duration: time.Millisecond * 1500,
	}

	windfuryBuffAura.Icd = &icd

	MakePermanent(character.GetOrRegisterAura(Aura{
		Label: "Windfury",
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			// charges are removed by every auto or next melee, whether it lands or not
			//  this directly contradicts https://github.com/magey/classic-warrior/wiki/Windfury-Totem#triggered-by-melee-spell-while-an-on-next-swing-attack-is-queued
			//  but can be seen in both "vanilla" and "sod" era logs
			if windfuryBuffAura.IsActive() && spell.ProcMask.Matches(ProcMaskMeleeWhiteHit) {
				windfuryBuffAura.RemoveStack(sim)
			}

			if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMeleeMH) {
				return
			}

			if icd.IsReady(sim) && sim.RandomFloat("Windfury") < 0.2 {
				icd.Use(sim)
				windfuryBuffAura.Activate(sim)
				// aura is up _before_ the triggering swing lands, so if triggered by an auto attack, the aura fades right after the extra attack lands.
				if spell.ProcMask == ProcMaskMeleeMHAuto {
					windfuryBuffAura.SetStacks(sim, 1)
				} else {
					windfuryBuffAura.SetStacks(sim, 2)
				}
				aura.Unit.AutoAttacks.ExtraMHAttack(sim, 1, ActionID{SpellID: 10610})
			}
		},
	}))

	return windfuryBuffAura
}

///////////////////////////////////////////////////////////////////////////
//                            World Buffs
///////////////////////////////////////////////////////////////////////////

func ApplyRallyingCryOfTheDragonslayer(character *Character) *Aura {
	character.AddStat(stats.SpellCrit, 10*SpellCritRatingPerCritChance)
	character.AddStat(stats.MeleeCrit, 5*CritRatingPerCritChance)
	// TODO: character.MultiplyStat(stats.RangedCrit, 1.05)
	character.AddStat(stats.AttackPower, 140)
	character.AddStat(stats.RangedAttackPower, 140)

	return character.RegisterAura(Aura{
		Label:    "Rallying Cry of the Dragonslayer",
		ActionID: ActionID{SpellID: 22888},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplySpiritOfZandalar(character *Character) *Aura {
	character.MultiplyStat(stats.Stamina, 1.15)
	character.MultiplyStat(stats.Agility, 1.15)
	character.MultiplyStat(stats.Strength, 1.15)
	character.MultiplyStat(stats.Intellect, 1.15)
	character.MultiplyStat(stats.Spirit, 1.15)

	return character.RegisterAura(Aura{
		Label:    "Spirit of Zandalar",
		ActionID: ActionID{SpellID: 24425},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplySongflowerSerenate(character *Character) *Aura {
	character.AddStat(stats.MeleeCrit, 5*CritRatingPerCritChance)
	// TODO: character.AddStat(stats.RangedCrit, 1.05)
	character.AddStat(stats.SpellCrit, 5*SpellCritRatingPerCritChance)
	character.AddStat(stats.Stamina, 15)
	character.AddStat(stats.Agility, 15)
	character.AddStat(stats.Strength, 15)
	character.AddStat(stats.Intellect, 15)
	character.AddStat(stats.Spirit, 15)

	return character.RegisterAura(Aura{
		Label:    "Songflower Serenade",
		ActionID: ActionID{SpellID: 15366},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplyWarchiefsBlessing(character *Character) *Aura {
	character.AddStat(stats.Health, 300)
	character.PseudoStats.MeleeSpeedMultiplier *= 1.15
	character.AddStat(stats.MP5, 10)

	return character.RegisterAura(Aura{
		Label:    "Warchief's Blessing",
		ActionID: ActionID{SpellID: 16609},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplyMightOfStormwind(character *Character) *Aura {
	character.AddStat(stats.Health, 300)
	character.PseudoStats.MeleeSpeedMultiplier *= 1.15
	character.AddStat(stats.MP5, 10)

	return character.RegisterAura(Aura{
		Label:    "Might of Stormwind",
		ActionID: ActionID{SpellID: 460940},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplyFengusFerocity(character *Character) *Aura {
	character.AddStat(stats.AttackPower, 200)
	character.AddStat(stats.RangedAttackPower, 200)

	return character.RegisterAura(Aura{
		Label:    "Fengus' Ferocity",
		ActionID: ActionID{SpellID: 22817},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplyMoldarsMoxie(character *Character) *Aura {
	character.MultiplyStat(stats.Stamina, 1.15)

	return character.RegisterAura(Aura{
		Label:    "Moldar's Moxie",
		ActionID: ActionID{SpellID: 22818},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplySlipkiksSavvy(character *Character) *Aura {
	character.AddStat(stats.SpellCrit, 3*SpellCritRatingPerCritChance)

	return character.RegisterAura(Aura{
		Label:    "Slip'kik's Savvy",
		ActionID: ActionID{SpellID: 22820},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplySaygesFortunes(character *Character, fortune proto.SaygesFortune) *Aura {
	var label string
	var spellID int32

	if fortune == proto.SaygesFortune_SaygesDamage {
		character.PseudoStats.DamageDealtMultiplier *= 1.10
		label = "Sayge's Dark Fortune of Damage"
		spellID = 23768
	}

	if fortune == proto.SaygesFortune_SaygesAgility {
		character.MultiplyStat(stats.Agility, 1.10)
		label = "Sayge's Dark Fortune of Agility"
		spellID = 23736
	}

	if fortune == proto.SaygesFortune_SaygesIntellect {
		character.MultiplyStat(stats.Intellect, 1.10)
		label = "Sayge's Dark Fortune of Intellect"
		spellID = 23766
	}

	if fortune == proto.SaygesFortune_SaygesSpirit {
		character.MultiplyStat(stats.Spirit, 1.10)
		label = "Sayge's Dark Fortune of Spirit"
		spellID = 23738
	}

	if fortune == proto.SaygesFortune_SaygesStamina {
		character.MultiplyStat(stats.Stamina, 1.10)
		label = "Sayge's Dark Fortune of Stamina"
		spellID = 23737
	}

	return character.RegisterAura(Aura{
		Label:    label,
		ActionID: ActionID{SpellID: spellID},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplyFervorOfTheTempleExplorer(character *Character) *Aura {
	if character.Level > 59 {
		return nil
	}

	character.AddStat(stats.MeleeCrit, 5*CritRatingPerCritChance)
	// TODO: character.AddStat(stats.RangedCrit, 5 * CritRatingPerCritChance)
	character.AddStat(stats.SpellCrit, 5*CritRatingPerCritChance)
	character.AddStat(stats.SpellDamage, 65) // TODO: confirm if spellpower or spelldamage
	character.MultiplyStat(stats.Strength, 1.08)
	character.MultiplyStat(stats.Stamina, 1.08)
	character.MultiplyStat(stats.Agility, 1.08)
	character.MultiplyStat(stats.Intellect, 1.08)
	character.MultiplyStat(stats.Spirit, 1.08)

	return character.RegisterAura(Aura{
		Label:    "Fervor of the Temple Explorer",
		ActionID: ActionID{SpellID: 446695},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplySparkOfInspiration(character *Character) *Aura {
	if character.Level > 49 {
		return nil
	}

	character.AddStat(stats.SpellCrit, 4*CritRatingPerCritChance)
	character.AddStat(stats.SpellPower, 42)
	character.PseudoStats.MeleeSpeedMultiplier *= 1.1
	character.PseudoStats.RangedSpeedMultiplier *= 1.1

	return character.RegisterAura(Aura{
		Label:    "Spark of Inspiration",
		ActionID: ActionID{SpellID: 438536},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplyBoonOfBlackfathom(character *Character) *Aura {
	if character.Level > 39 {
		return nil
	}

	character.AddStat(stats.MeleeCrit, 2*CritRatingPerCritChance)
	// TODO: character.AddStat(stats.RangedCrit, 2 * CritRatingPerCritChance)
	character.AddStat(stats.SpellHit, 3*SpellHitRatingPerHitChance)
	character.AddStat(stats.AttackPower, 20)
	character.AddStat(stats.RangedAttackPower, 20)
	character.AddStat(stats.SpellPower, 25)

	return character.RegisterAura(Aura{
		Label:    "Boon of Blackfathom",
		ActionID: ActionID{SpellID: 430947},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}

func ApplyAshenvaleRallyingCry(character *Character) *Aura {
	if character.Level > 39 {
		return nil
	}

	character.PseudoStats.DamageDealtMultiplier *= 1.05
	//TODO: healing dealt multiplier?

	return character.RegisterAura(Aura{
		Label:    "Ashenvale Rallying Cry",
		ActionID: ActionID{SpellID: 430352},
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
	})
}
