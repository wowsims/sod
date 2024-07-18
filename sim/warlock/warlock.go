package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{17, 17, 16}

const (
	WarlockFlagAffliction  = core.SpellFlagAgentReserved1
	WarlockFlagDemonology  = core.SpellFlagAgentReserved2
	WarlockFlagDestruction = core.SpellFlagAgentReserved3
	WarlockFlagHaunt       = core.SpellFlagAgentReserved4
)

const (
	SpellCode_WarlockNone int32 = iota

	SpellCode_WarlockCorruption
	SpellCode_WarlockCurseOfAgony
	SpellCode_WarlockCurseOfDoom
	SpellCode_WarlockDrainLife
	SpellCode_WarlockDrainSoul
	SpellCode_WarlockHaunt
	SpellCode_WarlockImmolate
	SpellCode_WarlockIncinerate
	SpellCode_WarlockSearingPain
	SpellCode_WarlockShadowflame
	SpellCode_WarlockShadowCleave
	SpellCode_WarlockShadowBolt
	SpellCode_WarlockSoulFire
)

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	BasePets   []*WarlockPet
	ActivePet  *WarlockPet
	Felhunter  *WarlockPet
	Felguard   *WarlockPet
	Imp        *WarlockPet
	Succubus   *WarlockPet
	Voidwalker *WarlockPet

	// Doomguard *DoomguardPet
	// Infernal  *InfernalPet

	ChaosBolt          *core.Spell
	Conflagrate        []*core.Spell
	Corruption         []*core.Spell
	DarkPact           *core.Spell
	DrainSoul          []*core.Spell
	Haunt              *core.Spell
	Immolate           []*core.Spell
	Incinerate         *core.Spell
	LifeTap            []*core.Spell
	SearingPain        []*core.Spell
	ShadowBolt         []*core.Spell
	ShadowCleave       []*core.Spell
	Shadowburn         []*core.Spell
	SoulFire           []*core.Spell
	DemonicGrace       *core.Spell
	DrainLife          []*core.Spell
	RainOfFire         []*core.Spell
	SiphonLife         []*core.Spell
	DeathCoil          []*core.Spell
	Shadowflame        *core.Spell
	UnstableAffliction *core.Spell

	ActiveCurseAura          *core.Aura
	CurseOfElements          *core.Spell
	CurseOfElementsAuras     core.AuraArray
	CurseOfShadow            *core.Spell
	CurseOfShadowAuras       core.AuraArray
	CurseOfRecklessness      *core.Spell
	CurseOfRecklessnessAuras core.AuraArray
	CurseOfWeakness          *core.Spell
	CurseOfWeaknessAuras     core.AuraArray
	CurseOfTongues           *core.Spell
	CurseOfTonguesAuras      core.AuraArray
	CurseOfAgony             []*core.Spell
	CurseOfDoom              *core.Spell
	AmplifyCurse             *core.Spell

	SummonDemonSpells []*core.Spell

	DemonicKnowledgeAura    *core.Aura
	HauntDebuffAuras        core.AuraArray
	ImmolationAura          *core.Spell
	IncinerateAura          *core.Aura
	Metamorphosis           *core.Spell
	MetamorphosisAura       *core.Aura
	NightfallProcAura       *core.Aura
	PyroclasmAura           *core.Aura
	DemonicGraceAura        *core.Aura
	AmplifyCurseAura        *core.Aura
	BackdraftAura           *core.Aura
	ImprovedShadowBoltAuras core.AuraArray
	MarkOfChaosAuras        core.AuraArray
	SoulLinkAura            *core.Aura

	// The sum total of demonic pact spell power * seconds.
	DPSPAggregate float64
	PreviousTime  time.Duration

	demonicKnowledgeSp   float64
	nightfallProcChance  float64
	zilaGularAura        *core.Aura
	shadowSparkAura      *core.Aura
	defendersResolveAura *core.Aura
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) Initialize() {
	warlock.registerCorruptionSpell()
	warlock.registerImmolateSpell()
	warlock.registerShadowBoltSpell()
	warlock.registerLifeTapSpell()
	warlock.registerSoulFireSpell()
	warlock.registerShadowBurnSpell()
	// warlock.registerSeedSpell()
	warlock.registerDrainSoulSpell()
	warlock.registerConflagrateSpell()
	warlock.registerSiphonLifeSpell()
	warlock.registerDarkPactSpell()
	warlock.registerSearingPainSpell()
	// warlock.registerInfernoSpell()
	// warlock.registerBlackBook()
	warlock.registerDrainLifeSpell()
	warlock.registerRainOfFireSpell()
	warlock.registerDeathCoilSpell()

	warlock.registerCurseOfElementsSpell()
	warlock.registerCurseOfShadowSpell()
	warlock.registerCurseOfRecklessnessSpell()
	warlock.registerCurseOfAgonySpell()
	warlock.registerAmplifyCurseSpell()
	warlock.registerCurseOfDoomSpell()
	warlock.registerSummonDemon()

	warlock.registerPetAbilities()
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = max(raidBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.WarlockOptions_Imp,
		warlock.Talents.ImprovedImp == 3,
	))
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	warlock.setDefaultActivePet()

	// warlock.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand,
	// 	proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged}, false)
	// warlock.setupCooldowns(sim)
}

func NewWarlock(character *core.Character, options *proto.Player, warlockOptions *proto.WarlockOptions) *Warlock {
	warlock := &Warlock{
		Character: *character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions,
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	warlock.EnableManaBar()

	warlock.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	warlock.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	warlock.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[warlock.Class][int(warlock.Level)]*core.SpellCritRatingPerCritChance)
	warlock.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[warlock.Class][int(warlock.Level)]*core.CritRatingPerCritChance)

	switch warlock.Options.Armor {
	case proto.WarlockOptions_DemonArmor:
		warlock.applyDemonArmor()
	case proto.WarlockOptions_FelArmor:
		warlock.applyFelArmor()
	}

	warlock.registerPets()
	warlock.setDefaultActivePet()

	guardians.ConstructGuardians(&warlock.Character)

	return warlock
}

func (warlock *Warlock) HasRune(rune proto.WarlockRune) bool {
	return warlock.HasRuneById(int32(rune))
}

func (warlock *Warlock) baseRuneAbilityDamage() float64 {
	return 6.568597 + 0.672028*float64(warlock.Level) + 0.031721*float64(warlock.Level*warlock.Level)
}

func (warlock *Warlock) OnGCDReady(_ *core.Simulation) {
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}

func isWarlockSpell(spell *core.Spell) bool {
	return spell.Flags.Matches(WarlockFlagAffliction) || spell.Flags.Matches(WarlockFlagDemonology) || spell.Flags.Matches(WarlockFlagDestruction)
}
