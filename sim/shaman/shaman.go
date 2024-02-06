package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{15, 16, 15}

const (
	SpellFlagShock     = core.SpellFlagAgentReserved1
	SpellFlagElectric  = core.SpellFlagAgentReserved2
	SpellFlagTotem     = core.SpellFlagAgentReserved3
	SpellFlagFocusable = core.SpellFlagAgentReserved4
)

func NewShaman(character *core.Character, talents string, totems *proto.ShamanTotems, selfBuffs SelfBuffs) *Shaman {
	shaman := &Shaman{
		Character: *character,
		Talents:   &proto.ShamanTalents{},
		Totems:    totems,
		SelfBuffs: selfBuffs,
	}
	shaman.waterShieldManaMetrics = shaman.NewManaMetrics(core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsWaterShield)})

	core.FillTalentsProto(shaman.Talents.ProtoReflect(), talents, TalentTreeSizes)
	shaman.EnableManaBar()

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(shaman.Level)]*core.CritRatingPerCritChance)
	shaman.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(shaman.Level)]*core.SpellCritRatingPerCritChance)
	shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	// Set proper Melee Haste scaling
	shaman.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	if selfBuffs.Shield == proto.ShamanShield_WaterShield {
		shaman.AddStat(stats.MP5, shaman.MaxMana()*.01)
	}

	shaman.ApplyRockbiterImbue(shaman.getImbueProcMask(proto.ShamanImbue_RockbiterWeapon))
	shaman.ApplyFlametongueImbue(shaman.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))

	if !shaman.HasMHWeapon() {
		shaman.SelfBuffs.ImbueMH = proto.ShamanImbue_NoImbue
	}

	if !shaman.HasOHWeapon() {
		shaman.SelfBuffs.ImbueOH = proto.ShamanImbue_NoImbue
	}

	return shaman
}

func (shaman *Shaman) getImbueProcMask(imbue proto.ShamanImbue) core.ProcMask {
	var mask core.ProcMask
	if shaman.SelfBuffs.ImbueMH == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if shaman.SelfBuffs.ImbueOH == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
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

	SpellCode_HealingWave
	SpellCode_LesserHealingWave
	SpellCode_ChainHeal

	SpellCode_SearingTotem
	SpellCode_MagmaTotem
	SpellCode_FireNovaTotem
)

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	Talents   *proto.ShamanTalents
	SelfBuffs SelfBuffs

	Totems *proto.ShamanTotems

	// The expiration time of each totem (earth, air, fire, water).
	TotemExpirations [4]time.Duration

	LightningBolt         []*core.Spell
	LightningBoltOverload []*core.Spell

	ChainLightning         []*core.Spell
	ChainLightningOverload []*core.Spell

	Stormstrike *core.Spell

	LightningShield     *core.Spell
	LightningShieldAura *core.Aura

	EarthShock     []*core.Spell
	FlameShock     []*core.Spell
	FlameShockDots []*core.Spell
	FrostShock     []*core.Spell

	// Totems
	ActiveTotems [4]*core.Spell

	StrengthOfEarthTotem []*core.Spell
	StoneskinTotem       []*core.Spell
	TremorTotem          *core.Spell

	SearingTotem  []*core.Spell
	MagmaTotem    []*core.Spell
	FireNovaTotem []*core.Spell

	HealingStreamTotem []*core.Spell
	ManaSpringTotem    []*core.Spell

	WindfuryTotem   []*core.Spell
	GraceOfAirTotem []*core.Spell

	// Healing Spells
	tidalWaveProc *core.Aura

	HealingWave         []*core.Spell
	HealingWaveOverload []*core.Spell

	LesserHealingWave []*core.Spell

	ChainHeal         []*core.Spell
	ChainHealOverload []*core.Spell

	waterShieldManaMetrics *core.ResourceMetrics

	// Runes
	LavaBurst         *core.Spell
	LavaBurstOverload *core.Spell
	MoltenBlast       *core.Spell
	LavaLash          *core.Spell
	EarthShield       *core.Spell

	FireNova []*core.Spell

	MaelstromWeaponAura *core.Aura
	PowerSurgeAura      *core.Aura

	// Used by Ancestral Guidance rune
	lastFlameShockTarget *core.Unit

	AncestralAwakening     *core.Spell
	ancestralHealingAmount float64
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
	// Core abilities
	shaman.registerChainLightningSpell()
	shaman.registerLightningBoltSpell()
	// shaman.registerLightningShieldSpell()
	shaman.registerShocks()
	shaman.registerStormstrikeSpell()

	// Imbues
	// In the Initialize due to frost brand adding the aura to the enemy
	shaman.RegisterRockbiterImbue(shaman.getImbueProcMask(proto.ShamanImbue_RockbiterWeapon))
	shaman.RegisterFlametongueImbue(shaman.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))
	shaman.RegisterWindfuryImbue(shaman.getImbueProcMask(proto.ShamanImbue_WindfuryWeapon))
	shaman.RegisterFrostbrandImbue(shaman.getImbueProcMask(proto.ShamanImbue_FrostbrandWeapon))

	// if shaman.ItemSwap.IsEnabled() {
	// 	mh := shaman.ItemSwap.GetItem(proto.ItemSlot_ItemSlotMainHand)
	// 	shaman.ApplyFlametongueImbueToItem(mh, true)
	// 	oh := shaman.ItemSwap.GetItem(proto.ItemSlot_ItemSlotOffHand)
	// 	shaman.ApplyFlametongueImbueToItem(oh, false)
	// 	shaman.RegisterOnItemSwap(func(_ *core.Simulation) {
	// 		shaman.ApplySyncType(proto.ShamanSyncType_Auto)
	// 	})
	// }

	// Totems
	shaman.registerStrengthOfEarthTotemSpell()
	shaman.registerStoneskinTotemSpell()
	shaman.registerTremorTotemSpell()
	shaman.registerSearingTotemSpell()
	shaman.registerMagmaTotemSpell()
	shaman.registerFireNovaTotemSpell()
	shaman.registerHealingStreamTotemSpell()
	shaman.registerManaSpringTotemSpell()
	// shaman.registerWindfuryTotemSpell()
	// shaman.registerGraceofAirTotem()

	// // This registration must come after all the totems are registered
	// shaman.registerCallOfTheElements()

	shaman.RegisterHealingSpells()
}

func (shaman *Shaman) RegisterHealingSpells() {
	shaman.registerLesserHealingWaveSpell()
	shaman.registerHealingWaveSpell()
	shaman.registerChainHealSpell()
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

func (shaman *Shaman) TotemManaMultiplier() float64 {
	return 1 - 0.05*float64(shaman.Talents.TotemicFocus)
}
