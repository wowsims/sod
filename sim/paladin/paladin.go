package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SealDuration                = time.Second * 30
	SpellFlagSecondaryJudgement = core.SpellFlagAgentReserved1
	SpellFlagPrimaryJudgement   = core.SpellFlagAgentReserved2
)

var TalentTreeSizes = [3]int{14, 15, 15}

const (
	SpellCode_PaladinNone = iota
	SpellCode_PaladinHolyShock
)

type Paladin struct {
	core.Character
	Talents *proto.PaladinTalents

	PaladinAura proto.PaladinAura

	CurrentSeal           *core.Aura
	CurrentSealExpiration time.Duration
	PrimarySealSpell      *core.Spell

	// Variables for max rank seal spells, used in APL Actions.
	MaxRankRighteousness int
	MaxRankCommand       int

	CrusaderStrike *core.Spell
	DivineStorm    *core.Spell
	Consecration   []*core.Spell

	// These cooldowns need reset from elsewhere in the package.
	Exorcism          []*core.Spell
	ExorcismCooldown  *core.Cooldown
	HolyShock         []*core.Spell
	HolyShockCooldown *core.Cooldown
	Judgement         *core.Spell
	DivineFavor       *core.Spell
	// HolyShield            *core.Spell
	// HammerOfTheRighteous  *core.Spell
	// HandOfReckoning       *core.Spell
	// ShieldOfRighteousness *core.Spell
	// AvengersShield        *core.Spell
	// HammerOfWrath         *core.Spell
	// SealOfVengeance       *core.Spell
	SealOfRighteousness []*core.Spell
	SealOfCommand       []*core.Spell
	SealOfMartyrdom     *core.Spell
	// AvengingWrath         *core.Spell
	// DivineProtection      *core.Spell
	// SovDotSpell           *core.Spell
	// SealOfWisdom        *core.Spell
	// SealOfLight         *core.Spell

	// HolyShieldAura          *core.Aura
	// RighteousFuryAura       *core.Aura
	// DivinePleaAura          *core.Aura
	// JudgementOfWisdomAura   *core.Aura
	// JudgementOfLightAura    *core.Aura
	// SealOfVengeanceAura     *core.Aura

	SealOfRighteousnessAura []*core.Aura
	SealOfCommandAura       []*core.Aura
	SealOfMartyrdomAura     *core.Aura
	// AvengingWrathAura       *core.Aura
	// DivineProtectionAura    *core.Aura
	// ForbearanceAura         *core.Aura
	DivineFavorAura *core.Aura
	VengeanceAura   *core.Aura
	// HolyWrath             *core.Spell
	// SealOfWisdomAura        *core.Aura
	// SealOfLightAura         *core.Aura
	// ArtOfWarInstantCast *core.Aura
	// SpiritualAttunementMetrics *core.ResourceMetrics

	DemonAndUndeadTargetCount int32
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

func (paladin *Paladin) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	// Buffs are handled explicitly through APLs now
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	paladin.AutoAttacks.MHConfig().CritMultiplier = paladin.MeleeCritMultiplier()

	// paladin.registerSealOfVengeanceSpellAndAura()
	paladin.registerSealOfRighteousnessSpellAndAura()
	paladin.registerJudgementSpell()
	paladin.registerSealOfCommandSpellAndAura()
	paladin.registerSealOfMartyrdomSpellAndAura()
	// paladin.setupSealOfTheCrusader()
	// paladin.setupSealOfWisdom()
	// paladin.setupSealOfLight()
	// paladin.setupSealOfRighteousness()
	// paladin.setupJudgementRefresh()

	paladin.registerCrusaderStrikeSpell()
	paladin.registerDivineStormSpell()

	paladin.registerConsecrationSpell()
	paladin.registerHolyShockSpell()
	paladin.registerExorcismSpell()
	paladin.registerDivineFavorSpellAndAura()
	// paladin.registerHammerOfWrathSpell()
	// paladin.registerHolyWrathSpell()
	// paladin.registerHolyShieldSpell()
	// paladin.registerHammerOfTheRighteousSpell()
	// paladin.registerHandOfReckoningSpell()
	// paladin.registerShieldOfRighteousnessSpell()
	// paladin.registerAvengersShieldSpell()
	// paladin.registerJudgements()

	// paladin.registerSpiritualAttunement()
	// paladin.registerDivinePleaSpell()
	// paladin.registerDivineProtectionSpell()
	// paladin.registerForbearanceDebuff()

	for i := int32(0); i < paladin.Env.GetNumTargets(); i++ {
		unit := paladin.Env.GetTargetUnit(i)
		if unit.MobType == proto.MobType_MobTypeDemon || unit.MobType == proto.MobType_MobTypeUndead {
			paladin.DemonAndUndeadTargetCount += 1
		}
	}
}

func (paladin *Paladin) Reset(_ *core.Simulation) {
	paladin.CurrentSeal = nil
	// paladin.CurrentJudgement = nil
}

// maybe need to add stat dependencies
func NewPaladin(character *core.Character, talentsStr string) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
	}
	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	// This is used to cache its effect in talents.go
	// paladin.HasTuralyonsOrLiadrinsBattlegear2Pc = paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 2)

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()
	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)
	paladin.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(paladin.Level)]*core.CritRatingPerCritChance)
	paladin.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(paladin.Level)]*core.SpellCritRatingPerCritChance)

	// Paladins get 0.0167 dodge per agi. ~1% per 59.88
	paladin.AddStatDependency(stats.Agility, stats.Dodge, (1.0/59.88)*core.DodgeRatingPerDodgeChance)

	// Paladins get 1 block value per 20 str
	paladin.AddStatDependency(stats.Strength, stats.BlockValue, .05)

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// No diminishing returns in Vanilla.
	paladin.PseudoStats.BaseDodge += 0.034943
	paladin.PseudoStats.BaseParry += 0.05

	return paladin
}

func (paladin *Paladin) HasRune(rune proto.PaladinRune) bool {
	return paladin.HasRuneById(int32(rune))
}

func (paladin *Paladin) Has1hEquipped() bool {
	return paladin.GetMHWeapon().HandType == proto.HandType_HandTypeOneHand
}

func (paladin *Paladin) Has2hEquipped() bool {
	return paladin.GetMHWeapon().HandType == proto.HandType_HandTypeTwoHand
}

func (paladin *Paladin) GetMaxRankSeal(seal proto.PaladinSeal) *core.Spell {
	var returnSpell *core.Spell
	switch seal {
	case proto.PaladinSeal_Righteousness:
		returnSpell = paladin.SealOfRighteousness[paladin.MaxRankRighteousness]
	case proto.PaladinSeal_Command:
		returnSpell = paladin.SealOfCommand[paladin.MaxRankCommand]
	}
	return returnSpell
}
