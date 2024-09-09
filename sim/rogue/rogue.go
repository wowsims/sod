package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SpellFlagBuilder      = core.SpellFlagAgentReserved1
	SpellFlagColdBlooded  = core.SpellFlagAgentReserved2
	SpellFlagDeadlyBrewed = core.SpellFlagAgentReserved3
	SpellFlagCarnage      = core.SpellFlagAgentReserved4 // for Carnage
	SpellFlagRoguePoison  = core.SpellFlagAgentReserved5 // RogueT1
)

const (
	SpellCode_RogueNone int32 = iota

	SpellCode_RogueAmbush
	SpellCode_RogueBackstab
	SpellCode_RogueBetweentheEyes
	SpellCode_RogueBladeFlurry
	SpellCode_RogueCrimsonTempest
	SpellCode_RogueEnvenom
	SpellCode_RogueEviscerate
	SpellCode_RogueGarrote
	SpellCode_RogueGhostlyStrike
	SpellCode_RogueHemorrhage
	SpellCode_RogueMainGauche
	SpellCode_RogueMutilate
	SpellCode_RoguePoisonedKnife
	SpellCode_RogueRupture
	SpellCode_RogueSaberSlash
	SpellCode_RogueSaberSlashDoT
	SpellCode_RogueShadowStrike
	SpellCode_RogueSinisterStrike
)

var TalentTreeSizes = [3]int{15, 19, 17}

const RogueBleedTag = "RogueBleed"

type Rogue struct {
	core.Character

	Talents *proto.RogueTalents
	Options *proto.RogueOptions

	sliceAndDiceDurations [6]time.Duration
	bladeDanceDurations   [6]time.Duration

	Backstab            *core.Spell
	BladeFlurry         *core.Spell
	Feint               *core.Spell
	Garrote             *core.Spell
	Ambush              *core.Spell
	Hemorrhage          *core.Spell
	GhostlyStrike       *core.Spell
	HungerForBlood      *core.Spell
	Mutilate            *core.Spell
	MutilateMH          *core.Spell
	MutilateOH          *core.Spell
	Shiv                *core.Spell
	SinisterStrike      *core.Spell
	SaberSlash          *core.Spell
	saberSlashTick      *core.Spell
	MainGauche          *core.Spell
	Shadowstep          *core.Spell
	Preparation         *core.Spell
	Premeditation       *core.Spell
	ColdBlood           *core.Spell
	Vanish              *core.Spell
	Shadowstrike        *core.Spell
	QuickDraw           *core.Spell
	ShurikenToss        *core.Spell
	BetweenTheEyes      *core.Spell
	PoisonedKnife       *core.Spell
	Blunderbuss         *core.Spell
	FanOfKnives         *core.Spell
	CrimsonTempest      *core.Spell
	CrimsonTempestBleed *core.Spell

	Envenom      *core.Spell
	Eviscerate   *core.Spell
	ExposeArmor  *core.Spell
	Rupture      *core.Spell
	SliceAndDice *core.Spell

	Evasion    *core.Spell
	BladeDance *core.Spell

	DeadlyPoison     [3]*core.Spell
	deadlyPoisonTick *core.Spell
	InstantPoison    [3]*core.Spell
	WoundPoison      [2]*core.Spell
	OccultPoison     *core.Spell
	occultPoisonTick *core.Spell

	instantPoisonProcChanceBonus float64
	additivePoisonBonusChance    float64
	cutthroatBonusChance         float64

	AdrenalineRushAura            *core.Aura
	BladeFlurryAura               *core.Aura
	EnvenomAura                   *core.Aura
	ExposeArmorAuras              core.AuraArray
	EvasionAura                   *core.Aura
	BladeDanceAura                *core.Aura
	SliceAndDiceAura              *core.Aura
	MasterOfSubtletyAura          *core.Aura
	ShadowstepAura                *core.Aura
	ShadowDanceAura               *core.Aura
	StealthAura                   *core.Aura
	WaylayAuras                   core.AuraArray
	RollingWithThePunchesAura     *core.Aura
	RollingWithThePunchesProcAura *core.Aura
	CutthroatProcAura             *core.Aura
	VanishAura                    *core.Aura

	HonorAmongThieves *core.Aura

	woundPoisonDebuffAuras core.AuraArray
}

func (rogue *Rogue) GetCharacter() *core.Character {
	return &rogue.Character
}

func (rogue *Rogue) GetRogue() *Rogue {
	return rogue
}

func (rogue *Rogue) AddRaidBuffs(_ *proto.RaidBuffs)   {}
func (rogue *Rogue) AddPartyBuffs(_ *proto.PartyBuffs) {}

func (rogue *Rogue) finisherFlags() core.SpellFlag {
	return SpellFlagCarnage | core.SpellFlagMeleeMetrics | core.SpellFlagAPL
}

func (rogue *Rogue) builderFlags() core.SpellFlag {
	return SpellFlagBuilder | SpellFlagColdBlooded | SpellFlagCarnage | core.SpellFlagMeleeMetrics | core.SpellFlagAPL
}

func (rogue *Rogue) Initialize() {
	rogue.registerBackstabSpell()
	rogue.registerEviscerate()
	rogue.registerExposeArmorSpell()
	rogue.registerFeintSpell()
	rogue.registerGarrote()
	rogue.registerHemorrhageSpell()
	rogue.registerRupture()
	rogue.registerSinisterStrikeSpell()
	rogue.registerSliceAndDice()
	rogue.registerThistleTeaCD()
	rogue.registerAmbushSpell()

	// Poisons
	rogue.registerInstantPoisonSpell()
	rogue.registerDeadlyPoisonSpell()
	rogue.registerOccultPoisonSpell()
	rogue.registerWoundPoisonSpell()
	rogue.registerSebaciousPoisonSpell()

	// Stealth
	rogue.registerStealthAura()
	rogue.registerVanishSpell()
}

func (rogue *Rogue) ApplyEnergyTickMultiplier(multiplier float64) {
	rogue.EnergyTickMultiplier += multiplier
}

func (rogue *Rogue) Reset(_ *core.Simulation) {
	for _, mcd := range rogue.GetMajorCooldowns() {
		mcd.Disable()
	}
}

func NewRogue(character *core.Character, options *proto.Player, rogueOptions *proto.RogueOptions) *Rogue {
	rogue := &Rogue{
		Character: *character,
		Talents:   &proto.RogueTalents{},
		Options:   rogueOptions,
	}
	core.FillTalentsProto(rogue.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	// Passive rogue threat reduction: https://wotlk.wowhead.com/spell=21184/rogue-passive-dnd
	rogue.PseudoStats.ThreatMultiplier *= 0.71
	// TODO: Be able to Parry based on results
	rogue.PseudoStats.CanParry = true
	maxEnergy := 100.0
	if rogue.Talents.Vigor {
		maxEnergy += 10
	}
	rogue.EnableEnergyBar(maxEnergy)

	rogue.EnableAutoAttacks(rogue, core.AutoAttackOptions{
		MainHand:       rogue.WeaponFromMainHand(),
		OffHand:        rogue.WeaponFromOffHand(),
		Ranged:         rogue.WeaponFromRanged(),
		AutoSwingMelee: true,
	})
	rogue.applyPoisons()

	rogue.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	rogue.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.RangedAttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(rogue.Level)]*core.CritRatingPerCritChance)
	rogue.AddStatDependency(stats.Agility, stats.Dodge, core.DodgePerAgiAtLevel[character.Class][int(rogue.Level)]*core.DodgeRatingPerDodgeChance)
	rogue.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	guardians.ConstructGuardians(&rogue.Character)

	return rogue
}

// Deactivate Stealth if it is active. This must be added to all abilities that cause Stealth to fade.
func (rogue *Rogue) BreakStealth(sim *core.Simulation) {
	if rogue.StealthAura.IsActive() {
		rogue.StealthAura.Deactivate(sim)
		rogue.AutoAttacks.EnableAutoSwing(sim)
	}
}

// Does the rogue have a dagger equipped in the specified hand (main or offhand)?
func (rogue *Rogue) HasDagger(hand core.Hand) bool {
	if hand == core.MainHand {
		return rogue.MainHand().WeaponType == proto.WeaponType_WeaponTypeDagger
	}
	return rogue.OffHand().WeaponType == proto.WeaponType_WeaponTypeDagger
}

// Check if the rogue is considered in "stealth" for the purpose of casting abilities
func (rogue *Rogue) IsStealthed() bool {
	return rogue.StealthAura.IsActive()
}

// Agent is a generic way to access underlying rogue on any of the agents.
type RogueAgent interface {
	GetRogue() *Rogue
}

func (rogue *Rogue) HasRune(rune proto.RogueRune) bool {
	return rogue.HasRuneById(int32(rune))
}

func (rogue *Rogue) baseRuneAbilityDamage() float64 {
	return 5.741530 - 0.255683*float64(rogue.Level) + 0.032656*float64(rogue.Level*rogue.Level)
}

func (rogue *Rogue) baseRuneAbilityDamageCombo() float64 {
	return 8.740728 - 0.415787*float64(rogue.Level) + 0.051973*float64(rogue.Level*rogue.Level)
}

func (rogue *Rogue) getImbueProcMask(imbue proto.WeaponImbue) core.ProcMask {
	var mask core.ProcMask
	if rogue.HasMHWeapon() && rogue.Consumes.MainHandImbue == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if rogue.HasOHWeapon() && rogue.Consumes.OffHandImbue == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}
