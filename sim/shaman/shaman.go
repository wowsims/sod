package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{25, 29, 26}

// Start looking to refresh 5 minute totems at 4:55.
const TotemRefreshTime5M = time.Second * 295

const (
	SpellFlagShock     = core.SpellFlagAgentReserved1
	SpellFlagElectric  = core.SpellFlagAgentReserved2
	SpellFlagTotem     = core.SpellFlagAgentReserved3
	SpellFlagFocusable = core.SpellFlagAgentReserved4
)

func NewShaman(character *core.Character, talents string, totems *proto.ShamanTotems, selfBuffs SelfBuffs, thunderstormRange bool) *Shaman {
	shaman := &Shaman{
		Character:           *character,
		Talents:             &proto.ShamanTalents{},
		Totems:              totems,
		SelfBuffs:           selfBuffs,
		thunderstormInRange: thunderstormRange,
	}
	shaman.waterShieldManaMetrics = shaman.NewManaMetrics(core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsWaterShield)})

	core.FillTalentsProto(shaman.Talents.ProtoReflect(), talents, TalentTreeSizes)
	shaman.EnableManaBar()

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	shaman.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(shaman.Level)]*core.CritRatingPerCritChance)
	shaman.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(shaman.Level)]*core.SpellCritRatingPerCritChance)
	shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	// Set proper Melee Haste scaling
	shaman.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	if selfBuffs.Shield == proto.ShamanShield_WaterShield {
		shaman.AddStat(stats.MP5, shaman.MaxMana()*.01)
	}

	return shaman
}

// Which buffs this shaman is using.
type SelfBuffs struct {
	Shield  proto.ShamanShield
	ImbueMH proto.ShamanImbue
	ImbueOH proto.ShamanImbue
}

// Indexes into NextTotemDrops for self buffs
const (
	AirTotem int = iota
	EarthTotem
	FireTotem
	WaterTotem
)

type ShamanSpellCode int

const (
	SpellCode_ShamanLightningBolt ShamanSpellCode = iota
	SpellCode_ShamanChainLightning
)

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	thunderstormInRange bool // flag if thunderstorm will be in range.

	Talents   *proto.ShamanTalents
	SelfBuffs SelfBuffs

	Totems *proto.ShamanTotems

	// The expiration time of each totem (earth, air, fire, water).
	TotemExpirations [4]time.Duration

	LightningBolt         []*core.Spell
	LightningBoltOverload []*core.Spell

	ChainLightning         []*core.Spell
	ChainLightningOverload []*core.Spell

	FireNova    *core.Spell
	Stormstrike *core.Spell

	LightningShield     *core.Spell
	LightningShieldAura *core.Aura

	Thunderstorm *core.Spell

	EarthShock     []*core.Spell
	FlameShock     []*core.Spell
	FlameShockDots []*core.Spell
	FrostShock     []*core.Spell

	// Totems
	StoneskinTotem       *core.Spell
	StrengthOfEarthTotem *core.Spell
	TremorTotem          *core.Spell

	SearingTotem  *core.Spell
	MagmaTotem    *core.Spell
	FireNovaTotem *core.Spell

	HealingStreamTotem *core.Spell
	ManaSpringTotem    *core.Spell

	TotemOfWrath    *core.Spell
	WindfuryTotem   *core.Spell
	GraceOfAirTotem *core.Spell

	// Healing Spells
	tidalWaveProc          *core.Aura
	ancestralHealingAmount float64
	AncestralAwakening     *core.Spell
	LesserHealingWave      *core.Spell
	Riptide                *core.Spell

	HealingWave         *core.Spell
	HealingWaveOverload *core.Spell

	ChainHeal         *core.Spell
	ChainHealOverload *core.Spell

	waterShieldManaMetrics *core.ResourceMetrics

	// Runes
	LavaBurst         *core.Spell
	LavaBurstOverload *core.Spell
	MoltenBlast       *core.Spell
	LavaLash          *core.Spell
	EarthShield       *core.Spell

	MaelstromWeaponAura *core.Aura

	// Used by Ancestral Guidance rune
	LastFlameShockTarget *core.Unit
}

// Implemented by each Shaman spec.
type ShamanAgent interface {
	core.Agent

	// The Shaman controlled by this Agent.
	GetShaman() *Shaman
}

func (shaman *Shaman) GetCharacter() *core.Character {
	return &shaman.Character
}

// TODO: Totem buffs are party-wide
func (shaman *Shaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	switch shaman.Totems.Earth {
	case proto.EarthTotem_StrengthOfEarthTotem:
		totem := proto.TristateEffect_TristateEffectRegular
		if shaman.Talents.EnhancingTotems == 3 {
			totem = proto.TristateEffect_TristateEffectImproved
		}
		raidBuffs.StrengthOfEarthTotem = max(raidBuffs.StrengthOfEarthTotem, totem)
	case proto.EarthTotem_StoneskinTotem:
		raidBuffs.StoneskinTotem = max(raidBuffs.StoneskinTotem, core.MakeTristateValue(
			true,
			shaman.Talents.GuardianTotems == 2,
		))
	}

	switch shaman.Totems.Fire {
	}

	switch shaman.Totems.Water {
	case proto.WaterTotem_ManaSpringTotem:
		raidBuffs.ManaSpringTotem = max(raidBuffs.ManaSpringTotem, proto.TristateEffect_TristateEffectRegular)
		if shaman.Talents.RestorativeTotems == 5 {
			raidBuffs.ManaSpringTotem = proto.TristateEffect_TristateEffectImproved
		}
	}

	switch shaman.Totems.Air {
	case proto.AirTotem_GraceOfAirTotem:
		totem := proto.TristateEffect_TristateEffectRegular
		if shaman.Talents.EnhancingTotems == 3 {
			totem = proto.TristateEffect_TristateEffectImproved
		}
		raidBuffs.GraceOfAirTotem = max(raidBuffs.StrengthOfEarthTotem, totem)
	case proto.AirTotem_WindfuryTotem:
		raidBuffs.WindfuryTotem = true
	}
}

func (shaman *Shaman) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if shaman.Talents.ManaTideTotem {
		partyBuffs.ManaTideTotems++
	}
}

func (shaman *Shaman) Initialize() {
	shaman.registerChainLightningSpell()
	// shaman.registerFeralSpirit()
	// shaman.registerFireNovaSpell()
	shaman.registerLightningBoltSpell()
	// shaman.registerLightningShieldSpell()
	// shaman.registerMagmaTotemSpell()
	// shaman.registerManaSpringTotemSpell()
	// shaman.registerHealingStreamTotemSpell()
	// shaman.registerSearingTotemSpell()
	shaman.registerShocks()
	// shaman.registerStormstrikeSpell()
	// shaman.registerStrengthOfEarthTotemSpell()
	// shaman.registerTremorTotemSpell()
	// shaman.registerStoneskinTotemSpell()
	// shaman.registerWindfuryTotemSpell()
	// shaman.registerGraceofAirTotem()

	// // This registration must come after all the totems are registered
	// shaman.registerCallOfTheElements()
}

func (shaman *Shaman) RegisterHealingSpells() {
	// shaman.registerAncestralHealingSpell()
	// shaman.registerLesserHealingWaveSpell()
	// shaman.registerHealingWaveSpell()
	// shaman.registerRiptideSpell()
	// shaman.registerEarthShieldSpell()
	// shaman.registerChainHealSpell()

	// if shaman.Talents.TidalWaves > 0 {
	// 	shaman.tidalWaveProc = shaman.GetOrRegisterAura(core.Aura{
	// 		Label:    "Tidal Wave Proc",
	// 		ActionID: core.ActionID{SpellID: 53390},
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Deactivate(sim)
	// 		},
	// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
	// 			shaman.HealingWave.CastTimeMultiplier *= 0.7
	// 			shaman.LesserHealingWave.BonusCritRating += core.CritRatingPerCritChance * 25
	// 		},
	// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
	// 			shaman.HealingWave.CastTimeMultiplier /= 0.7
	// 			shaman.LesserHealingWave.BonusCritRating -= core.CritRatingPerCritChance * 25
	// 		},
	// 		MaxStacks: 2,
	// 	})
	// }
}

func (shaman *Shaman) HasRune(rune proto.ShamanRune) bool {
	return shaman.HasRuneById(int32(rune))
}

func (shaman *Shaman) Reset(sim *core.Simulation) {
}

func (shaman *Shaman) ElementalCritMultiplier(secondary float64) float64 {
	critBonus := core.TernaryFloat64(shaman.Talents.ElementalFury, 1, 0) + secondary
	return shaman.SpellCritMultiplier(1, critBonus)
}

func (shaman *Shaman) ShamanThreatMultiplier(secondary float64) float64 {
	return core.TernaryFloat64(shaman.HasRune(proto.ShamanRune_RuneLegsWayOfEarth), 1.5, 1) * secondary
}
