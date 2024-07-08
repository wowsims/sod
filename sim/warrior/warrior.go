package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SpellFlagBleed      = core.SpellFlagAgentReserved1
	SpellFlagBloodSurge = core.SpellFlagAgentReserved2
)

var TalentTreeSizes = [3]int{18, 17, 17}

type WarriorInputs struct {
	StanceSnapshot bool
}

const (
	ArmsTree = 0
	FuryTree = 1
	ProtTree = 2
)

type Warrior struct {
	core.Character

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance                 Stance
	revengeProcAura        *core.Aura
	OverpowerAura          *core.Aura
	BerserkerRageAura      *core.Aura
	BloodrageAura          *core.Aura
	ConsumedByRageAura     *core.Aura
	Above80RageCBRActive   bool
	BloodSurgeAura         *core.Aura
	RampageAura            *core.Aura
	WreckingCrewEnrageAura *core.Aura
	EnrageAura             *core.Aura

	// Rune passive
	FocusedRageDiscount float64

	// Reaction time values
	reactionTime       time.Duration
	lastBloodsurgeProc time.Duration
	LastAMTick         time.Duration

	BattleShout *core.Spell

	BattleStance    *core.Spell
	DefensiveStance *core.Spell
	BerserkerStance *core.Spell
	GladiatorStance *core.Spell

	Bloodrage         *core.Spell
	BerserkerRage     *core.Spell
	Bloodthirst       *core.Spell
	DemoralizingShout *core.Spell
	Execute           *core.Spell
	MortalStrike      *core.Spell
	Overpower         *core.Spell
	Rend              *core.Spell
	Revenge           *core.Spell
	ShieldBlock       *core.Spell
	ShieldSlam        *core.Spell
	Slam              *core.Spell
	SunderArmor       *core.Spell
	Devastate         *core.Spell
	ThunderClap       *core.Spell
	Whirlwind         *core.Spell
	WhirlwindMH       *core.Spell
	WhirlwindOH       *core.Spell
	DeepWounds        *core.Spell
	ConcussionBlow    *core.Spell
	RagingBlow        *core.Spell
	Hamstring         *core.Spell
	Rampage           *core.Spell

	HeroicStrike       *core.Spell
	QuickStrike        *core.Spell
	Cleave             *core.Spell
	curQueueAura       *core.Aura
	curQueuedAutoSpell *core.Spell

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura
	GladiatorStanceAura *core.Aura

	ShieldBlockAura *core.Aura

	DemoralizingShoutAuras core.AuraArray
	SunderArmorAuras       core.AuraArray
	ThunderClapAuras       core.AuraArray
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	primaryTimer := warrior.NewTimer()
	overpowerRevengeTimer := warrior.NewTimer()

	warrior.reactionTime = time.Millisecond * 500

	warrior.registerShouts()
	warrior.registerStances()
	warrior.registerBerserkerRageSpell()
	warrior.registerBloodthirstSpell(primaryTimer)
	warrior.registerCleaveSpell()
	warrior.registerDemoralizingShoutSpell()
	warrior.registerExecuteSpell()
	warrior.registerHeroicStrikeSpell()
	warrior.registerMortalStrikeSpell(primaryTimer)
	warrior.registerOverpowerSpell(overpowerRevengeTimer)
	warrior.registerRevengeSpell(overpowerRevengeTimer)
	warrior.registerShieldSlamSpell()
	warrior.registerSlamSpell()
	warrior.registerThunderClapSpell()
	warrior.registerWhirlwindSpell()
	warrior.registerRendSpell()
	warrior.registerHamstringSpell()

	warrior.SunderArmor = warrior.newSunderArmorSpell()

	warrior.registerBloodrageCD()
	warrior.RegisterShieldBlockCD()
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
	warrior.curQueueAura = nil
	warrior.curQueuedAutoSpell = nil
}

func NewWarrior(character *core.Character, talents string, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:     *character,
		Talents:       &proto.WarriorTalents{},
		WarriorInputs: inputs,
	}
	core.FillTalentsProto(warrior.Talents.ProtoReflect(), talents, TalentTreeSizes)

	warrior.PseudoStats.CanParry = true

	warrior.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	warrior.AddStatDependency(stats.Strength, stats.BlockValue, .05) // 20 str = 1 block
	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(warrior.Level)]*core.CritRatingPerCritChance)
	warrior.AddStatDependency(stats.Agility, stats.Dodge, core.DodgePerAgiAtLevel[character.Class][int(warrior.Level)])
	warrior.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	guardians.ConstructGuardians(&warrior.Character)

	return warrior
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}

func (warrior *Warrior) HasRune(rune proto.WarriorRune) bool {
	return warrior.HasRuneById(int32(rune))
}

func (warrior *Warrior) IsEnraged() bool {
	return warrior.ConsumedByRageAura.IsActive() || warrior.BloodrageAura.IsActive() || warrior.BerserkerRageAura.IsActive()
}
