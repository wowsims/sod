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
	SpellFlagFocusable = core.SpellFlagAgentReserved3
	SpellFlagMaelstrom = core.SpellFlagAgentReserved4
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
	shaman.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(shaman.Level)]*core.SpellCritRatingPerCritChance)
	shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	shaman.ApplyRockbiterImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_RockbiterWeapon))
	shaman.ApplyFlametongueImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_FlametongueWeapon))
	shaman.ApplyFrostbrandImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_FrostbrandWeapon))
	shaman.ApplyWindfuryImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_WindfuryWeapon))

	if shaman.HasRune(proto.ShamanRune_RuneCloakFeralSpirit) {
		shaman.SpiritWolves = &SpiritWolves{
			SpiritWolf1: shaman.NewSpiritWolf(1),
			SpiritWolf2: shaman.NewSpiritWolf(2),
		}
	}

	guardians.ConstructGuardians(&shaman.Character)

	return shaman
}

func (shaman *Shaman) getImbueProcMask(_ *core.Character, imbue proto.WeaponImbue) core.ProcMask {
	var mask core.ProcMask
	if shaman.HasMHWeapon() && shaman.Consumes.MainHandImbue == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if shaman.HasOHWeapon() && shaman.Consumes.OffHandImbue == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}

// Indexes into NextTotemDrops for self buffs
const (
	AirTotem int = iota
	EarthTotem
	FireTotem
	WaterTotem
)

const (
	SpellCode_ShamanNone int32 = iota

	SpellCode_ShamanChainHeal
	SpellCode_ShamanChainLightning
	SpellCode_ShamanEarthShock
	SpellCode_ShamanFireNova
	SpellCode_ShamanFireNovaTotem
	SpellCode_ShamanFlameShock
	SpellCode_ShamanFrostShock
	SpellCode_ShamanHealingWave
	SpellCode_ShamanLesserHealingWave
	SpellCode_ShamanLightningBolt
	SpellCode_ShamanLightningShield
	SpellCode_ShamanLavaBurst
	SpellCode_ShamanMagmaTotem
	SpellCode_ShamanMoltenBlast
	SpellCode_ShamanSearingTotem
	SpellCode_ShamanStormstrike
)

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	Talents *proto.ShamanTalents

	AncestralAwakening     *core.Spell
	ChainHeal              []*core.Spell
	ChainHealOverload      []*core.Spell
	ChainLightning         []*core.Spell
	ChainLightningOverload []*core.Spell
	EarthShield            *core.Spell
	EarthShock             []*core.Spell
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
	RollingThunder         *core.Spell
	SearingTotem           []*core.Spell
	StoneskinTotem         []*core.Spell
	StormstrikeMH          *core.Spell
	StormstrikeOH          *core.Spell
	StrengthOfEarthTotem   []*core.Spell
	TremorTotem            *core.Spell
	WaterShield            *core.Spell
	WindfuryTotem          []*core.Spell
	WindwallTotem          []*core.Spell

	ActiveShield     *core.Spell // Tracks the Shaman's active shield spell
	ActiveShieldAura *core.Aura

	FlurryAura            *core.Aura
	FlurryConsumptionAura *core.Aura // Trigger aura for consuming Flurry stacks on hit
	MaelstromWeaponAura   *core.Aura
	PowerSurgeAura        *core.Aura

	// Totems
	ActiveTotems     [4]*core.Spell
	Totems           *proto.ShamanTotems
	TotemExpirations [4]time.Duration // The expiration time of each totem (earth, air, fire, water).

	SpiritWolves *SpiritWolves

	lastFlameShockTarget *core.Unit // Used by Ancestral Guidance rune

	ancestralHealingAmount float64 // Used by Ancestral Awakening
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
	character := shaman.GetCharacter()

	// Core abilities
	shaman.registerChainLightningSpell()
	shaman.registerLightningBoltSpell()
	shaman.registerLightningShieldSpell()
	shaman.registerShocks()
	shaman.registerStormstrikeSpell()

	// Imbues
	// In the Initialize due to frost brand adding the aura to the enemy
	shaman.RegisterRockbiterImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_RockbiterWeapon))
	shaman.RegisterFlametongueImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_FlametongueWeapon))
	shaman.RegisterWindfuryImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_WindfuryWeapon))
	shaman.RegisterFrostbrandImbue(shaman.getImbueProcMask(character, proto.WeaponImbue_FrostbrandWeapon))

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
	for i := range shaman.TotemExpirations {
		shaman.TotemExpirations[i] = 0
	}
}
