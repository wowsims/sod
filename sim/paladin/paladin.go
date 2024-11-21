package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{14, 15, 15}

const (
	SpellFlag_RV          = core.SpellFlagAgentReserved1
	SpellFlag_Forbearance = core.SpellFlagAgentReserved2
)

const (
	SpellCode_PaladinNone = iota

	SpellCode_PaladinExorcism
	SpellCode_PaladinHolyShock
	SpellCode_PaladinHolyWrath
	SpellCode_PaladinJudgementOfCommand
	SpellCode_PaladinConsecration
	SpellCode_PaladinAvengersShield
	SpellCode_PaladinHolyShield
	SpellCode_PaladinHolyShieldProc
	SpellCode_PaladinLayOnHands
	SpellCode_PaladinHammerOfWrath
	SpellCode_PaladinCrusaderStrike
	SpellCode_PaladinHammerOfTheRighteous
	SpellCode_PaladinShieldOfRighteousness
)

type SealJudgeCode uint8

const (
	SealJudgeCodeNone        SealJudgeCode = 0
	SealJudgeOfMartyrdomCode SealJudgeCode = 1 << iota
	SealJudgeOfRighteousnessCode
	SealJudgeOfCommandCode
	SealJudgeOfTheCrusaderCode
)

type Paladin struct {
	core.Character

	Talents *proto.PaladinTalents
	Options *proto.PaladinOptions

	primarySeal        *core.Spell // the seal configured in options, available via "Cast Primary Seal"
	primaryPaladinAura proto.PaladinAura
	currentPaladinAura *core.Aura

	currentSeal      *core.Aura
	prevSeal         *core.Aura
	allSealAuras     [][]*core.Aura
	aurasSoM         []*core.Aura
	aurasSoR         []*core.Aura
	aurasSoC         []*core.Aura
	aurasSotC        []*core.Aura
	currentJudgement *core.Spell
	prevJudgement    *core.Spell
	allJudgeSpells   [][]*core.Spell
	spellsJoM        []*core.Spell
	spellsJoR        []*core.Spell
	spellsJoC        []*core.Spell
	spellsJotC       []*core.Spell

	// Active abilities and shared cooldowns that are externally manipulated.
	holyShockCooldown *core.Cooldown
	exorcismCooldown  *core.Cooldown
	crusaderStrike    *core.Spell
	divineStorm       *core.Spell
	exorcism          []*core.Spell
	judgement         *core.Spell
	rv                *core.Spell
	holyShieldAura    [3]*core.Aura
	holyShieldProc    [3]*core.Spell
	redoubtAura       *core.Aura
	holyWrath         []*core.Spell

	// highest rank seal spell if available
	sealOfRighteousness *core.Spell
	sealOfCommand       *core.Spell
	sealOfMartyrdom     *core.Spell

	enableMultiJudge    bool
	lingerDuration      time.Duration
	consumeSealsOnJudge bool
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
	paladin.registerRighteousFury()
	// Judgement and Seals
	paladin.registerJudgement()

	paladin.registerSealOfRighteousness()
	paladin.registerSealOfCommand()
	paladin.registerSealOfMartyrdom()
	paladin.registerSealOfTheCrusader()

	paladin.allJudgeSpells = append(paladin.allJudgeSpells, paladin.spellsJoM)
	paladin.allJudgeSpells = append(paladin.allJudgeSpells, paladin.spellsJoR)
	paladin.allJudgeSpells = append(paladin.allJudgeSpells, paladin.spellsJoC)
	paladin.allJudgeSpells = append(paladin.allJudgeSpells, paladin.spellsJotC)

	paladin.allSealAuras = append(paladin.allSealAuras, paladin.aurasSoM)
	paladin.allSealAuras = append(paladin.allSealAuras, paladin.aurasSoR)
	paladin.allSealAuras = append(paladin.allSealAuras, paladin.aurasSoC)
	paladin.allSealAuras = append(paladin.allSealAuras, paladin.aurasSotC)

	// Active abilities
	paladin.registerForbearance()
	paladin.registerCrusaderStrike()
	paladin.registerDivineStorm()
	paladin.registerConsecration()
	paladin.registerHolyShock()
	paladin.registerExorcism()
	paladin.registerDivineFavor()
	paladin.registerHammerOfWrath()
	paladin.registerHolyWrath()
	paladin.registerAvengingWrath()
	paladin.registerAuraMastery()
	paladin.registerHolyShield()
	paladin.registerShieldOfRighteousness()
	paladin.registerBlessingOfSanctuary()
	paladin.registerLayOnHands()

	paladin.enableMultiJudge = false // Was previously true in Phase 4 but disabled in Phase 5
	paladin.lingerDuration = time.Millisecond * 400
	paladin.consumeSealsOnJudge = true

	paladin.registerStopAttackMacros()

	paladin.ResetCurrentPaladinAura()
}

func (paladin *Paladin) Reset(_ *core.Simulation) {
	paladin.ResetCurrentPaladinAura()
	paladin.ResetPrimarySeal(paladin.Options.PrimarySeal)
}

// maybe need to add stat dependencies
func NewPaladin(character *core.Character, options *proto.Player, paladinOptions *proto.PaladinOptions) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
		Options:   paladinOptions,
	}
	core.FillTalentsProto(paladin.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	if paladin.Options.Aura == proto.PaladinAura_SanctityAura {
		paladin.primaryPaladinAura = paladin.Options.Aura
	}

	paladin.PseudoStats.CanParry = true
	paladin.EnableManaBar()
	paladin.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	paladin.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(paladin.Level)]*core.CritRatingPerCritChance)
	paladin.AddStatDependency(stats.Agility, stats.Dodge, core.CritPerAgiAtLevel[character.Class][int(paladin.Level)]*core.CritRatingPerCritChance)
	paladin.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(paladin.Level)]*core.SpellCritRatingPerCritChance)

	// Paladins get 1 block value per 20 str
	paladin.PseudoStats.BlockValuePerStrength = 0.05

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	guardians.ConstructGuardians(&paladin.Character)

	return paladin
}

func (paladin *Paladin) hasRune(rune proto.PaladinRune) bool {
	return paladin.HasRuneById(int32(rune))
}

func (paladin *Paladin) has2hEquipped() bool {
	return paladin.MainHand().HandType == proto.HandType_HandTypeTwoHand
}

func (paladin *Paladin) ResetPrimarySeal(primarySeal proto.PaladinSeal) {
	paladin.currentSeal = nil
	paladin.primarySeal = paladin.getPrimarySealSpell(primarySeal)
}

func (paladin *Paladin) registerStopAttackMacros() {

	if paladin.divineStorm != nil && paladin.Options.IsUsingDivineStormStopAttack {
		paladin.divineStorm.Flags |= core.SpellFlagBatchStopAttackMacro
	}

	if paladin.crusaderStrike != nil && paladin.Options.IsUsingCrusaderStrikeStopAttack {
		paladin.crusaderStrike.Flags |= core.SpellFlagBatchStopAttackMacro
	}

	for _, spellsJoX := range paladin.allJudgeSpells {
		for _, v := range spellsJoX {
			if v != nil && paladin.Options.IsUsingJudgementStopAttack {
				v.Flags |= core.SpellFlagBatchStopAttackMacro
			}
		}
	}
}

func (paladin *Paladin) ResetCurrentPaladinAura() {
	paladin.currentPaladinAura = nil
	if paladin.primaryPaladinAura == proto.PaladinAura_SanctityAura {
		paladin.currentPaladinAura = core.SanctityAuraAura(paladin.GetCharacter())
	}
}

func (paladin *Paladin) getPrimarySealSpell(primarySeal proto.PaladinSeal) *core.Spell {
	// Used in the Cast Primary Seal APLAction to get the max rank spell for the level.
	switch primarySeal {
	case proto.PaladinSeal_Martyrdom:
		return paladin.sealOfMartyrdom
	case proto.PaladinSeal_Command:
		return paladin.sealOfCommand
	case proto.PaladinSeal_Righteousness:
		return paladin.sealOfRighteousness
	default:
		return paladin.sealOfRighteousness
	}
}

func (paladin *Paladin) applySeal(newSeal *core.Aura, judgement *core.Spell, sim *core.Simulation) {
	isSameSealType := false

	if paladin.currentSeal != nil {
		if newSeal.Label[:10] == paladin.currentSeal.Label[:10] {
			isSameSealType = true

			paladin.currentSeal.Deactivate(sim)
			paladin.currentSeal = newSeal
			paladin.currentJudgement = judgement
			paladin.currentSeal.Activate(sim)

			// Set To nil to avoid issues with multi judging during linger window
			if paladin.prevSeal != nil && paladin.prevSeal.IsActive() {
				paladin.prevSeal.Deactivate(sim)
			}

			paladin.prevSeal = nil
			paladin.prevJudgement = nil
		}
	}

	if !isSameSealType {
		if paladin.currentSeal.IsActive() {
			paladin.currentSeal.UpdateExpires(sim, sim.CurrentTime+paladin.lingerDuration) // always update, even if it extends duration
		}

		paladin.prevSeal = paladin.currentSeal
		paladin.currentSeal = newSeal
		paladin.prevJudgement = paladin.currentJudgement // Judgment Spell for the previous Seal (doesn't mean it was cast)
		paladin.currentJudgement = judgement
		paladin.currentSeal.Activate(sim)
	}
}

func (paladin *Paladin) getLibramSealCostReduction() float64 {
	if paladin.Ranged().ID == LibramOfBenediction {
		return 10
	}
	if paladin.Ranged().ID == LibramOfHope {
		return 20
	}
	return 0
}

func (paladin *Paladin) baseRuneAbilityDamage() float64 {
	return 9.046514 + 0.676562*float64(paladin.Level) + 0.019349*float64(paladin.Level*paladin.Level)
}
