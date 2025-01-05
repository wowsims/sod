package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{15, 16, 15}

const (
	SpellFlagShaman    = core.SpellFlagAgentReserved1
	SpellFlagTotem     = core.SpellFlagAgentReserved2
	SpellFlagLightning = core.SpellFlagAgentReserved3
	SpellFlagMaelstrom = core.SpellFlagAgentReserved4
)

const (
	SpellCode_ShamanNone int32 = iota

	SpellCode_ShamanChainHeal
	SpellCode_ShamanChainLightning
	SpellCode_ShamanEarthShock
	SpellCode_ShamanFireNova
	SpellCode_ShamanFireNovaTotem
	SpellCode_ShamanFireNovaTotemAttack
	SpellCode_ShamanFlameShock
	SpellCode_ShamanFrostShock
	SpellCode_ShamanHealingWave
	SpellCode_ShamanLavaLash
	SpellCode_ShamanLesserHealingWave
	SpellCode_ShamanLightningBolt
	SpellCode_ShamanLightningShield
	SpellCode_ShamanLavaBurst
	SpellCode_ShamanMagmaTotem
	SpellCode_ShamanMagmaTotemAttack
	SpellCode_ShamanMoltenBlast
	SpellCode_ShamanSearingTotem
	SpellCode_ShamanSearingTotemAttack
	SpellCode_ShamanStormstrike
	SpellCode_ShamanStormstrikeHit
)

// Indexes into NextTotemDrops for self buffs
const (
	AirTotem int = iota
	EarthTotem
	FireTotem
	WaterTotem
)

func NewShaman(character *core.Character, talents string) *Shaman {
	shaman := &Shaman{
		Character: *character,
		Talents:   &proto.ShamanTalents{},
	}

	core.FillTalentsProto(shaman.Talents.ProtoReflect(), talents, TalentTreeSizes)
	shaman.EnableManaBar()

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(shaman.Level)]*core.CritRatingPerCritChance)
	shaman.AddStatDependency(stats.Agility, stats.Dodge, core.DodgePerAgiAtLevel[character.Class][int(shaman.Level)]*core.DodgeRatingPerDodgeChance)
	shaman.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(shaman.Level)]*core.SpellCritRatingPerCritChance)
	shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	shaman.PseudoStats.BlockValuePerStrength = .05 // 20 str = 1 block

	shaman.ApplyRockbiterImbue(shaman.getImbueProcMask(proto.WeaponImbue_RockbiterWeapon))
	shaman.ApplyFlametongueImbue(shaman.getImbueProcMask(proto.WeaponImbue_FlametongueWeapon))
	shaman.ApplyFrostbrandImbue(shaman.getImbueProcMask(proto.WeaponImbue_FrostbrandWeapon))
	shaman.ApplyWindfuryImbue(shaman.getImbueProcMask(proto.WeaponImbue_WindfuryWeapon))

	if shaman.HasRune(proto.ShamanRune_RuneCloakFeralSpirit) {
		shaman.SpiritWolves = shaman.NewSpiritWolves()
	}

	guardians.ConstructGuardians(&shaman.Character)

	return shaman
}

func (shaman *Shaman) getImbueProcMask(imbue proto.WeaponImbue) core.ProcMask {
	mask := core.ProcMaskUnknown
	if shaman.HasMHWeapon() && shaman.Consumes.MainHandImbue == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if shaman.HasOHWeapon() && shaman.Consumes.OffHandImbue == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	Talents *proto.ShamanTalents

	// Spells
	AncestralAwakening     *core.Spell
	ChainHeal              []*core.Spell
	ChainHealOverload      []*core.Spell
	ChainLightning         []*core.Spell
	ChainLightningOverload []*core.Spell
	EarthShield            *core.Spell
	EarthShock             []*core.Spell
	ElementalMastery       *core.Spell
	FeralSpirit            *core.Spell
	FireNova               *core.Spell
	FireNovaTotem          []*core.Spell
	FlameShock             []*core.Spell
	FrostShock             []*core.Spell
	GraceOfAirTotem        []*core.Spell
	HealingStreamTotem     []*core.Spell
	HealingWave            []*core.Spell
	HealingWaveOverload    []*core.Spell
	LavaBurst              *core.Spell
	LavaBurstOverload      *core.Spell
	LavaLash               *core.Spell
	LesserHealingWave      []*core.Spell
	LightningBolt          []*core.Spell
	LightningBoltOverload  []*core.Spell
	LightningShield        []*core.Spell
	LightningShieldProcs   []*core.Spell // The damage component of lightning shield is a separate spell
	MagmaTotem             []*core.Spell
	ManaSpringTotem        []*core.Spell
	MoltenBlast            *core.Spell
	Riptide                *core.Spell
	RollingThunder         *core.Spell
	SearingTotem           []*core.Spell
	StoneskinTotem         []*core.Spell
	Stormstrike            *core.Spell
	StormstrikeMH          *core.Spell
	StormstrikeOH          *core.Spell
	StrengthOfEarthTotem   []*core.Spell
	TremorTotem            *core.Spell
	WaterShield            *core.Spell
	WaterShieldRestore     *core.Spell
	WindfuryTotem          []*core.Spell
	WindfuryWeaponMH       *core.Spell
	WindfuryWeaponOH       *core.Spell
	WindwallTotem          []*core.Spell

	// Auras
	ClearcastingAura     *core.Aura
	LightningShieldAuras []*core.Aura
	LoyalBetaAura        *core.Aura
	MaelstromWeaponAura  *core.Aura
	PowerSurgeDamageAura *core.Aura
	PowerSurgeHealAura   *core.Aura
	ShamanisticRageAura  *core.Aura
	ShieldMasteryAura    *core.Aura
	SpiritOfTheAlphaAura *core.Aura
	WaterShieldAura      *core.Aura

	// Totems
	ActiveTotems     [4]*core.Spell
	EarthTotems      []*core.Spell
	FireTotems       []*core.Spell
	WaterTotems      []*core.Spell
	AirTotems        []*core.Spell
	Totems           *proto.ShamanTotems
	TotemExpirations [4]time.Duration // The expiration time of each totem (earth, air, fire, water).

	// Shield
	ActiveShield     *core.Spell // Tracks the Shaman's active shield spell
	ActiveShieldAura *core.Aura

	// Pets
	SpiritWolves *SpiritWolves

	// Other data
	ancestralHealingAmount      float64 // Used by Ancestral Awakening
	bonusFlurrySpeed            float64 // Bonus added on top of the normal speed, e.g. Earthfury Impact 6pc
	bonusWindfuryWeaponAP       float64
	elementalFocusProcChance    float64
	lastFlameShockTarget        *core.Unit // Used by Ancestral Guidance rune
	lightningShieldCanCrit      bool
	maelstromWeaponPPMM         *core.PPMManager
	powerSurgeProcChance        float64
	shamanisticRageDRMultiplier float64
	staticSHocksProcChance      float64
	useLavaBurstCritScaling     bool
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

func (shaman *Shaman) AddRaidBuffs(_ *proto.RaidBuffs) {
	// Buffs are handled explicitly through APLs now
}

func (shaman *Shaman) Initialize() {
	// Core abilities
	shaman.registerChainLightningSpell()
	shaman.registerLightningBoltSpell()
	shaman.registerLightningShieldSpell()
	shaman.registerShocks()
	shaman.registerStormstrikeSpell()

	// Imbues
	// In the Initialize due to frost brand adding the aura to the enemy
	shaman.RegisterRockbiterImbue(shaman.getImbueProcMask(proto.WeaponImbue_RockbiterWeapon))
	shaman.RegisterFlametongueImbue(shaman.getImbueProcMask(proto.WeaponImbue_FlametongueWeapon))
	shaman.RegisterWindfuryImbue(shaman.getImbueProcMask(proto.WeaponImbue_WindfuryWeapon))
	shaman.RegisterFrostbrandImbue(shaman.getImbueProcMask(proto.WeaponImbue_FrostbrandWeapon))

	// Totems
	shaman.registerStrengthOfEarthTotemSpell()
	shaman.registerStoneskinTotemSpell()
	shaman.registerTremorTotemSpell()
	shaman.registerSearingTotemSpell()
	shaman.registerMagmaTotemSpell()
	shaman.registerFireNovaTotemSpell()
	shaman.registerHealingStreamTotemSpell()
	shaman.registerManaSpringTotemSpell()
	shaman.registerWindfuryTotemSpell()
	shaman.registerGraceOfAirTotemSpell()
	shaman.registerWindwallTotemSpell()

	// Other Abilities
	shaman.registerShamanisticRageCD()

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

func (shaman *Shaman) baseRuneAbilityDamage() float64 {
	return 7.583798 + 0.471881*float64(shaman.Level) + 0.036599*float64(shaman.Level*shaman.Level)
}

func (shaman *Shaman) Reset(_ *core.Simulation) {
	shaman.ActiveShield = nil
	shaman.ActiveShieldAura = nil

	for i := range shaman.TotemExpirations {
		shaman.TotemExpirations[i] = 0
	}
}
