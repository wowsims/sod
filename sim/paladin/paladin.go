package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/common/vanilla"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{14, 15, 15}

const (
	SpellCode_PaladinNone = iota
	SpellCode_PaladinHolyShock
	SpellCode_PaladinJudgementOfCommand
)

type Paladin struct {
	core.Character
	Talents *proto.PaladinTalents

	PaladinAura proto.PaladinAura

	CurrentJudgement      *core.Spell
	CurrentSeal           *core.Aura
	CurrentSealExpiration time.Duration
	PrimarySealSpell      *core.Spell

	// Variables for max rank seal spells, used in APL Actions.
	MaxRankRighteousness int
	MaxRankCommand       int

	// Active abilities and shared cooldowns that need externally manipulated.
	CrusaderStrike    *core.Spell
	DivineStorm       *core.Spell
	Consecration      []*core.Spell
	Exorcism          []*core.Spell
	ExorcismCooldown  *core.Cooldown
	HolyShock         []*core.Spell
	HolyShockCooldown *core.Cooldown
	Judgement         *core.Spell
	DivineFavor       *core.Spell
	// HolyShield            *core.Spell
	// HammerOfWrath         []*core.Spell
	// HolyWrath             []*core.Spell

	// Seal spells and their associated auras
	sealOfRighteousness *core.Spell
	SealOfCommand       []*core.Spell
	SealOfMartyrdom     *core.Spell
	SealOfTheCrusader   *core.Spell

	SealOfCommandAura   []*core.Aura
	SealOfMartyrdomAura *core.Aura

	// Auras from talents
	DivineFavorAura *core.Aura
	VengeanceAura   *core.Aura
}

// Implemented by each Paladin spec.
type PaladinAgent interface {
	GetPaladin() *Paladin
}

func (paladin *Paladin) GetCharacter() *core.Character {
	return &paladin.Character
}

func (paladin *Paladin) GetPaladin() *Paladin {
	return paladin
}

func (paladin *Paladin) AddRaidBuffs(_ *proto.RaidBuffs) {
	// Buffs are handled explicitly through APLs now
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {

	// Judgement and Seals
	paladin.registerJudgementSpell()
	paladin.registerSealOfRighteousness()

	paladin.registerSealOfCommand()
	paladin.registerSealOfMartyrdom()
	paladin.registerSealOfTheCrusader()

	// Active abilities
	paladin.registerCrusaderStrikeSpell()
	paladin.registerDivineStormSpell()
	paladin.registerConsecrationSpell()
	paladin.registerHolyShockSpell()
	paladin.registerExorcismSpell()
	paladin.registerDivineFavorSpellAndAura()
	paladin.registerHammerOfWrathSpell()
	paladin.registerHolyWrathSpell()
	// paladin.registerHolyShieldSpell()
}

func (paladin *Paladin) Reset(_ *core.Simulation) {
	paladin.CurrentSeal = nil
}

// maybe need to add stat dependencies
func NewPaladin(character *core.Character, talentsStr string) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
	}
	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	paladin.PseudoStats.CanParry = true
	paladin.EnableManaBar()
	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)
	paladin.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(paladin.Level)]*core.CritRatingPerCritChance)
	paladin.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(paladin.Level)]*core.SpellCritRatingPerCritChance)

	// Paladins get 1 block value per 20 str
	paladin.AddStatDependency(stats.Strength, stats.BlockValue, .05)

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// Dodge per agi at a given level behaves identically in classic to Crit per agi at a given level.
	// paladin.AddStatDependency(stats.Agility, stats.Dodge, core.CritPerAgiAtLevel[character.Class][int(paladin.Level)]*core.DodgeRatingPerDodgeChance)

	// The below requires some verification for the prot paladin sim when it is implemented.
	// Switch these to AddStat as the PsuedoStats are being removed
	// paladin.PseudoStats.BaseDodge += 0.034943
	// paladin.PseudoStats.BaseParry += 0.05

	vanilla.ConstructEmeralDragonWhelpPets(&paladin.Character)
	return paladin
}

func (paladin *Paladin) HasRune(rune proto.PaladinRune) bool {
	return paladin.HasRuneById(int32(rune))
}

func (paladin *Paladin) Has1hEquipped() bool {
	return paladin.MainHand().HandType == proto.HandType_HandTypeOneHand
}

func (paladin *Paladin) Has2hEquipped() bool {
	return paladin.MainHand().HandType == proto.HandType_HandTypeTwoHand
}

func (paladin *Paladin) GetMaxRankSeal(seal proto.PaladinSeal) *core.Spell {
	// Used in the Cast Primary Seal APLAction to get the max rank spell for the level.
	switch seal {
	case proto.PaladinSeal_Martyrdom:
		return paladin.SealOfMartyrdom
	case proto.PaladinSeal_Command:
		return paladin.SealOfCommand[paladin.MaxRankCommand]
	default:
		return paladin.sealOfRighteousness
	}
}

func (paladin *Paladin) ApplySeal(aura *core.Aura, judgement *core.Spell, sim *core.Simulation) {
	if paladin.CurrentSeal != nil {
		paladin.CurrentSeal.Deactivate(sim)
	}
	paladin.CurrentSeal = aura
	paladin.CurrentJudgement = judgement
	paladin.CurrentSeal.Activate(sim)
	paladin.CurrentSealExpiration = sim.CurrentTime + aura.Duration
}

func (paladin *Paladin) GetLibramSealCostReduction() float64 {
	if paladin.Ranged().ID == LibramOfBenediction {
		return 10
	}
	if paladin.Ranged().ID == LibramOfHope {
		return 20
	}
	return 0
}
