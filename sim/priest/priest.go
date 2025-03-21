package priest

import (
	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{15, 16, 16}

const (
	ClassSpellMask_PriestNone uint64 = 0

	ClassSpellMask_PriestDevouringPlague uint64 = 1 << iota
	ClassSpellMask_PriestFlashHeal
	ClassSpellMask_PriestGreaterHeal
	ClassSpellMask_PriestHeal
	ClassSpellMask_PriestHolyFire
	ClassSpellMask_PriestMindBlast
	ClassSpellMask_PriestMindFlay
	ClassSpellMask_PriestMindSear
	ClassSpellMask_PriestMindSpike
	ClassSpellMask_PriestShadowWordPain
	ClassSpellMask_PriestShadowWordDeath
	ClassSpellMask_PriestSmite
	ClassSpellMask_PriestVampiricTouch
	ClassSpellMask_PriestVoidPlague
	ClassSpellMask_PriestVoidZone
	ClassSpellMask_PriestDispersion
	ClassSpellMask_PriestShadowFiend
	ClassSpellMask_PriestVampiricEmbrace
	ClassSpellMask_PriestInnerFocus

	ClassSpellMask_PriestPenanceDamage
	ClassSpellMask_PriestPenanceHeal
	ClassSpellMask_PriestPenance = ClassSpellMask_PriestPenanceDamage | ClassSpellMask_PriestPenanceHeal

	ClassSpellMask_PriestAll = 1<<iota - 1

	ClassSpellMask_PriestDevouringPlagueDots = ClassSpellMask_PriestShadowWordPain | ClassSpellMask_PriestVoidPlague |
		ClassSpellMask_PriestVampiricTouch

	PriestSpellInstant = ClassSpellMask_PriestDevouringPlague |
		ClassSpellMask_PriestMindFlay |
		ClassSpellMask_PriestMindSear |
		ClassSpellMask_PriestShadowWordPain |
		ClassSpellMask_PriestShadowWordDeath |
		ClassSpellMask_PriestVoidPlague |
		ClassSpellMask_PriestVoidZone |
		ClassSpellMask_PriestDispersion |
		ClassSpellMask_PriestShadowFiend |
		ClassSpellMask_PriestVampiricEmbrace |
		ClassSpellMask_PriestPenance
)

type Priest struct {
	core.Character
	Talents *proto.PriestTalents

	Latency                     float64
	MindBlastCritChanceModifier float64

	CircleOfHealing   *core.Spell
	DevouringPlague   []*core.Spell
	Dispersion        *core.Spell
	EmpoweredRenew    *core.Spell
	EyeOfTheVoid      *core.Spell
	FlashHeal         []*core.Spell
	GreaterHeal       []*core.Spell
	HolyFire          []*core.Spell
	Homunculi         *core.Spell
	InnerFocus        *core.Spell
	MindBlast         []*core.Spell
	MindFlay          [][]*core.Spell // 1 entry for each tick for each rank
	MindSear          []*core.Spell   // 1 entry for each tick
	MindSearTicks     []*core.Spell   // 1 entry for each tick
	MindSpike         *core.Spell
	Penance           *core.Spell
	PenanceHeal       *core.Spell
	PowerWordShield   []*core.Spell
	PrayerOfHealing   []*core.Spell
	PrayerOfMending   *core.Spell
	Renew             []*core.Spell
	Shadowfiend       *core.Spell
	Shadowform        *core.Spell
	ShadowWeavingProc *core.Spell
	ShadowWordDeath   *core.Spell
	ShadowWordPain    []*core.Spell
	Smite             []*core.Spell
	VampiricEmbrace   *core.Spell
	VampiricTouch     *core.Spell
	VoidPlague        *core.Spell
	VoidZone          *core.Spell

	DispersionAura   *core.Aura
	EyeOfTheVoidAura *core.Aura
	HomunculiAura    *core.Aura
	InnerFocusAura   *core.Aura
	ShadowfiendAura  *core.Aura
	ShadowformAura   *core.Aura
	SpiritTapAura    *core.Aura
	SurgeOfLightAura *core.Aura

	MindSpikeAuras       core.AuraArray
	ShadowWeavingAuras   core.AuraArray
	VampiricEmbraceAuras core.AuraArray
	WeakenedSouls        core.AuraArray

	EyeOfTheVoidPet *EyeOfTheVoid
	HomunculiPets   []*Homunculus
	ShadowfiendPet  *Shadowfiend

	PainAndSufferingDoTSpells []*core.Spell

	ProcPrayerOfMending core.ApplySpellResults
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

func (priest *Priest) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ShadowProtection = true
	raidBuffs.DivineSpirit = true
	raidBuffs.PowerWordFortitude = max(
		raidBuffs.PowerWordFortitude,
		core.MakeTristateValue(true, priest.Talents.ImprovedPowerWordFortitude == 2),
	)
}

func (priest *Priest) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {
	priest.registerMindBlast()
	priest.registerMindFlay()
	priest.registerShadowWordPainSpell()
	priest.registerDevouringPlagueSpell()
	priest.RegisterSmiteSpell()
	priest.registerHolyFire()

	priest.registerPowerInfusionCD()
}

func (priest *Priest) RegisterHealingSpells() {
	// priest.registerFlashHealSpell()
	// priest.registerGreaterHealSpell()
	// priest.registerPowerWordShieldSpell()
	// priest.registerPrayerOfHealingSpell()
	// priest.registerRenewSpell()
}

func (priest *Priest) Reset(_ *core.Simulation) {
}

func New(character *core.Character, talents string) *Priest {
	priest := &Priest{
		Character: *character,
		Talents:   &proto.PriestTalents{},
	}
	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents, TalentTreeSizes)

	priest.EnableManaBar()

	priest.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	priest.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[priest.Class][int(priest.Level)]*core.SpellCritRatingPerCritChance)

	// Set mana regen to 12.5 + Spirit/4 each 2s tick
	priest.SpiritManaRegenPerSecond = func() float64 {
		return 6.25 + priest.GetStat(stats.Spirit)/8
	}

	priest.ShadowfiendPet = priest.NewShadowfiend()

	if priest.HasRune(proto.PriestRune_RuneHelmEyeOfTheVoid) {
		priest.EyeOfTheVoidPet = priest.NewEyeOfTheVoid()
	}

	if priest.HasRune(proto.PriestRune_RuneLegsHomunculi) {
		priest.HomunculiPets = make([]*Homunculus, 3)
		priest.HomunculiPets[0] = priest.NewHomunculus(1, 202390)
		priest.HomunculiPets[1] = priest.NewHomunculus(2, 202392)
		priest.HomunculiPets[2] = priest.NewHomunculus(3, 202391)
	}

	guardians.ConstructGuardians(&priest.Character)

	return priest
}

func (priest *Priest) HasRune(rune proto.PriestRune) bool {
	return priest.HasRuneById(int32(rune))
}

func (priest *Priest) baseRuneAbilityDamage() float64 {
	return 9.456667 + 0.635108*float64(priest.Level) + 0.039063*float64(priest.Level*priest.Level)
}

func (priest *Priest) baseRuneAbilityHealing() float64 {
	return 38.258376 + 0.904195*float64(priest.Level) + 0.161311*float64(priest.Level*priest.Level)
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}
