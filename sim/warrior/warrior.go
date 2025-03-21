package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SpellFlagOffensive = core.SpellFlagAgentReserved1
)

const (
	ClassSpellMask_WarriorNone uint64 = 0

	ClassSpellMask_WarriorBloodrage uint64 = 1 << iota
	ClassSpellMask_WarriorBloodthirst
	ClassSpellMask_WarriorCleave
	ClassSpellMask_WarriorDeathWish
	ClassSpellMask_WarriorDeepWounds
	ClassSpellMask_WarriorDevastate
	ClassSpellMask_WarriorExecute
	ClassSpellMask_WarriorHeroicStrike
	ClassSpellMask_WarriorMortalStrike
	ClassSpellMask_WarriorOverpower
	ClassSpellMask_WarriorQuickStrike
	ClassSpellMask_WarriorRagingBlow
	ClassSpellMask_WarriorRecklesness
	ClassSpellMask_WarriorRend
	ClassSpellMask_WarriorRevenge
	ClassSpellMask_WarriorShieldSlam
	ClassSpellMask_WarriorShieldWall
	ClassSpellMask_WarriorShockwave
	ClassSpellMask_WarriorSlam
	ClassSpellMask_WarriorSlamMH
	ClassSpellMask_WarriorSlamOH
	ClassSpellMask_WarriorStanceBattle
	ClassSpellMask_WarriorStanceBerserker
	ClassSpellMask_WarriorStanceGladiator
	ClassSpellMask_WarriorStanceDefensive
	ClassSpellMask_WarriorWhirlwind
	ClassSpellMask_WarriorWhirlwindMH
	ClassSpellMask_WarriorWhirlwindOH
)

var TalentTreeSizes = [3]int{18, 17, 17}

type WarriorInputs struct {
	QueueDelay     int32
	StanceSnapshot bool
	Stance         proto.WarriorStance
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
	Stance          Stance
	PreviousStance  Stance // Used for Warrior T1 DPS 4P
	revengeProcAura *core.Aura
	OverpowerAura   *core.Aura

	BloodSurgeAura      *core.Aura
	LastStandAura       *core.Aura
	RampageAura         *core.Aura
	ShieldBlockAura     *core.Aura
	SuddenDeathAura     *core.Aura
	TasteForBloodAura   *core.Aura
	lastMeleeAutoTarget *core.Unit

	// Enrage Auras
	BerserkerRageAura      *core.Aura
	BloodrageAura          *core.Aura
	ConsumedByRageAura     *core.Aura
	EnrageAura             *core.Aura
	FreshMeatEnrageAura    *core.Aura
	WreckingCrewEnrageAura *core.Aura

	// Reaction time values
	reactionTime time.Duration
	LastAMTick   time.Duration

	BattleShout *WarriorSpell

	BattleStanceSpells    []*WarriorSpell
	DefensiveStanceSpells []*WarriorSpell
	BerserkerStanceSpells []*WarriorSpell

	Stances         []*WarriorSpell
	BattleStance    *WarriorSpell
	DefensiveStance *WarriorSpell
	BerserkerStance *WarriorSpell
	GladiatorStance *WarriorSpell

	Bloodrage         *WarriorSpell
	BerserkerRage     *WarriorSpell
	Bloodthirst       *WarriorSpell
	DeathWish         *WarriorSpell
	DemoralizingShout *WarriorSpell
	Execute           *WarriorSpell
	MortalStrike      *WarriorSpell
	Overpower         *WarriorSpell
	Recklessness      *WarriorSpell
	Rend              *WarriorSpell
	Revenge           *WarriorSpell
	ShieldBlock       *WarriorSpell
	ShieldSlam        *WarriorSpell
	ShieldWall        *WarriorSpell
	Slam              *WarriorSpell
	SlamMH            *WarriorSpell
	SlamOH            *WarriorSpell
	SunderArmor       *WarriorSpell
	Devastate         *WarriorSpell
	ThunderClap       *WarriorSpell
	Whirlwind         *WarriorSpell
	WhirlwindMH       *WarriorSpell
	WhirlwindOH       *WarriorSpell
	DeepWounds        *WarriorSpell
	ConcussionBlow    *WarriorSpell
	RagingBlow        *WarriorSpell
	Hamstring         *WarriorSpell
	Rampage           *WarriorSpell
	Shockwave         *WarriorSpell

	HeroicStrike       *WarriorSpell
	HeroicStrikeQueue  *WarriorSpell
	QuickStrike        *WarriorSpell
	Cleave             *WarriorSpell
	CleaveQueue        *WarriorSpell
	curQueueAura       *core.Aura
	curQueuedAutoSpell *WarriorSpell

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura
	GladiatorStanceAura *core.Aura

	defensiveStanceThreatMultiplier   float64
	gladiatorStanceDamageMultiplier   float64
	recklessnessDamageTakenMultiplier float64

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

func (warrior *Warrior) RegisterSpell(stanceMask Stance, config core.SpellConfig) *WarriorSpell {
	ws := &WarriorSpell{
		StanceMask: stanceMask,
	}

	castConditionOld := config.ExtraCastCondition
	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		// Check if we're in a correct stance to cast the spell
		if stance := ws.GetStanceMask(); !ws.stanceOverride && stance != AnyStance && !warrior.StanceMatches(stance) {
			if sim.Log != nil {
				sim.Log("Failed cast to spell %s, wrong stance", ws.ActionID)
			}
			return false
		}
		return castConditionOld == nil || castConditionOld(sim, target)
	}

	ws.Spell = warrior.Unit.RegisterSpell(config)

	if stanceMask.Matches(BattleStance) {
		warrior.BattleStanceSpells = append(warrior.BattleStanceSpells, ws)
	}
	if stanceMask.Matches(DefensiveStance) {
		warrior.DefensiveStanceSpells = append(warrior.DefensiveStanceSpells, ws)
	}
	if stanceMask.Matches(BerserkerStance) {
		warrior.BerserkerStanceSpells = append(warrior.BerserkerStanceSpells, ws)
	}

	return ws
}

func (warrior *Warrior) newStanceOverrideExclusiveEffect(stance Stance, aura *core.Aura) *core.ExclusiveEffect {
	return aura.NewExclusiveEffect("stance-override", false, core.ExclusiveEffect{
		Priority: core.TernaryFloat64(stance == AnyStance, 2, 1),
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			if stance.Matches(BattleStance) {
				for _, spell := range warrior.BattleStanceSpells {
					spell.stanceOverride = true
				}
			}
			if stance.Matches(DefensiveStance) {
				for _, spell := range warrior.DefensiveStanceSpells {
					spell.stanceOverride = true
				}
			}
			if stance.Matches(BerserkerStance) {
				for _, spell := range warrior.BerserkerStanceSpells {
					spell.stanceOverride = true
				}
			}
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			if stance.Matches(BattleStance) {
				for _, spell := range warrior.BattleStanceSpells {
					spell.stanceOverride = false
				}
			}
			if stance.Matches(DefensiveStance) {
				for _, spell := range warrior.DefensiveStanceSpells {
					spell.stanceOverride = false
				}
			}
			if stance.Matches(BerserkerStance) {
				for _, spell := range warrior.BerserkerStanceSpells {
					spell.stanceOverride = false
				}
			}
		},
	})
}

func (warrior *Warrior) Initialize() {
	primaryTimer := warrior.NewTimer()
	overpowerRevengeTimer := warrior.NewTimer()

	warrior.reactionTime = time.Millisecond * 500

	warrior.registerShouts()
	warrior.registerStances()
	warrior.registerBerserkerRageSpell()
	warrior.registerBloodthirstSpell(primaryTimer)
	warrior.registerDemoralizingShoutSpell()
	warrior.registerExecuteSpell()
	warrior.registerMortalStrikeSpell(primaryTimer)
	warrior.registerOverpowerSpell(overpowerRevengeTimer)
	warrior.registerRevengeSpell(overpowerRevengeTimer)
	warrior.registerShieldSlamSpell()
	warrior.registerSlamSpell()
	warrior.registerThunderClapSpell()
	warrior.registerWhirlwindSpell()
	warrior.registerRendSpell()
	warrior.registerHamstringSpell()

	// The sim often re-enables heroic strike in an unrealistic amount of time.
	// This can cause an unrealistic immediate double-hit around wild strikes procs
	queuedRealismICD := &core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Millisecond * time.Duration(warrior.WarriorInputs.QueueDelay),
	}
	warrior.registerHeroicStrikeSpell(queuedRealismICD)
	warrior.registerCleaveSpell(queuedRealismICD)

	warrior.SunderArmor = warrior.registerSunderArmorSpell()

	warrior.registerBloodrageCD()
	warrior.RegisterRecklessnessCD()
}

func (warrior *Warrior) Reset(sim *core.Simulation) {
	warrior.curQueueAura = nil
	warrior.curQueuedAutoSpell = nil

	// Reset Stance
	switch warrior.WarriorInputs.Stance {
	case proto.WarriorStance_WarriorStanceBattle:
		warrior.Stance = BattleStance
		warrior.BattleStanceAura.Activate(sim)
	case proto.WarriorStance_WarriorStanceDefensive:
		warrior.Stance = DefensiveStance
		warrior.DefensiveStanceAura.Activate(sim)
	case proto.WarriorStance_WarriorStanceBerserker:
		warrior.Stance = BerserkerStance
		warrior.BerserkerStanceAura.Activate(sim)
	case proto.WarriorStance_WarriorStanceGladiator:
		warrior.Stance = GladiatorStance
		warrior.GladiatorStanceAura.Activate(sim)
	default:
		// Fallback to checking for Glad Stance rune or checking talent tree
		if warrior.GladiatorStanceAura != nil {
			warrior.Stance = GladiatorStance
			warrior.GladiatorStanceAura.Activate(sim)
		} else if warrior.PrimaryTalentTree == ArmsTree {
			warrior.Stance = BattleStance
			warrior.BattleStanceAura.Activate(sim)
		} else if warrior.PrimaryTalentTree == FuryTree {
			warrior.Stance = BerserkerStance
			warrior.BerserkerStanceAura.Activate(sim)
		} else {
			warrior.Stance = DefensiveStance
			warrior.DefensiveStanceAura.Activate(sim)
		}
	}
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
	warrior.PseudoStats.BlockValuePerStrength = .05 // 20 str = 1 block
	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(warrior.Level)]*core.CritRatingPerCritChance)
	warrior.AddStatDependency(stats.Agility, stats.Dodge, core.DodgePerAgiAtLevel[character.Class][int(warrior.Level)]*core.DodgeRatingPerDodgeChance)
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
	return warrior.BloodrageAura.IsActive() ||
		warrior.BerserkerRageAura.IsActive() ||
		(warrior.EnrageAura != nil && warrior.EnrageAura.IsActive()) ||
		(warrior.ConsumedByRageAura != nil && warrior.ConsumedByRageAura.IsActive()) ||
		(warrior.FreshMeatEnrageAura != nil && warrior.FreshMeatEnrageAura.IsActive()) ||
		(warrior.WreckingCrewEnrageAura != nil && warrior.WreckingCrewEnrageAura.IsActive())
}

type WarriorSpell struct {
	*core.Spell
	StanceMask     Stance
	stanceOverride bool // Allows the override of the StanceMask so that the spell can be used in any stance
}

func (ws *WarriorSpell) IsReady(sim *core.Simulation) bool {
	if ws == nil {
		return false
	}
	return ws.Spell.IsReady(sim)
}

func (ws *WarriorSpell) CanCast(sim *core.Simulation, target *core.Unit) bool {
	if ws == nil {
		return false
	}
	return ws.Spell.CanCast(sim, target)
}

func (ws *WarriorSpell) IsEqual(s *core.Spell) bool {
	if ws == nil || s == nil {
		return false
	}
	return ws.Spell == s
}

// Returns the StanceMask accounting for a possible override
func (ws *WarriorSpell) GetStanceMask() Stance {
	if ws.stanceOverride {
		return AnyStance
	}

	return ws.StanceMask
}
