package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{17, 17, 16}

const SpellFlagHaunt = core.SpellFlagAgentReserved1

const (
	SpellCode_WarlockNone int32 = iota

	SpellCode_WarlockCorruption
	SpellCode_WarlockDrainLife
	SpellCode_WarlockImmolate
	SpellCode_WarlockIncinerate
	SpellCode_WarlockShadowCleave
	SpellCode_WarlockShadowBolt
	SpellCode_WarlockSoulFire
)

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	Pet *WarlockPet

	ChaosBolt          *core.Spell
	Conflagrate        *core.Spell
	Corruption         []*core.Spell
	DarkPact           *core.Spell
	DrainSoul          *core.Spell
	Haunt              *core.Spell
	Immolate           []*core.Spell
	Incinerate         *core.Spell
	LifeTap            *core.Spell
	SearingPain        *core.Spell
	ShadowBolt         *core.Spell
	ShadowCleave       []*core.Spell
	Shadowburn         *core.Spell
	SoulFire           []*core.Spell
	DemonicGrace       *core.Spell
	DrainLife          *core.Spell
	RainOfFire         *core.Spell
	SiphonLife         *core.Spell
	DeathCoil          *core.Spell
	UnstableAffliction *core.Spell

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
	CurseOfAgony             *core.Spell
	CurseOfDoom              *core.Spell
	AmplifyCurse             *core.Spell
	Shadowflame              *core.Spell
	ShadowflameDot           *core.Spell

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

	// The sum total of demonic pact spell power * seconds.
	DPSPAggregate float64
	PreviousTime  time.Duration

	petStmBonusSP         float64
	demonicKnowledgeSp    float64
	demonicSacrificeAuras []*core.Aura
	zilaGularAura         *core.Aura
	shadowSparkAura       *core.Aura
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) Initialize() {
	warlock.registerChaosBoltSpell()
	warlock.registerCorruptionSpell()
	warlock.registerImmolateSpell()
	warlock.registerIncinerateSpell()
	warlock.registerShadowBoltSpell()
	warlock.registerShadowCleaveSpell()
	warlock.registerLifeTapSpell()
	warlock.registerSoulFireSpell()
	warlock.registerUnstableAfflictionSpell()
	// warlock.registerSeedSpell()
	// warlock.registerDrainSoulSpell()
	warlock.registerConflagrateSpell()
	warlock.registerHauntSpell()
	warlock.registerSiphonLifeSpell()
	warlock.registerMetamorphosisSpell()
	warlock.registerDarkPactSpell()
	warlock.registerShadowBurnSpell()
	warlock.registerSearingPainSpell()
	// warlock.registerInfernoSpell()
	// warlock.registerBlackBook()
	warlock.registerDemonicGraceSpell()
	warlock.registerDrainLifeSpell()
	warlock.registerRainOfFireSpell()
	warlock.registerShadowflameSpell()
	warlock.registerDeathCoilSpell()

	warlock.registerCurseOfElementsSpell()
	warlock.registerCurseOfShadowSpell()
	warlock.registerCurseOfRecklessnessSpell()
	warlock.registerCurseOfAgonySpell()
	warlock.registerAmplifyCurseSpell()
	// warlock.registerCurseOfDoomSpell()
	warlock.registerImmolationAuraSpell()
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = max(raidBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.WarlockOptions_Imp,
		warlock.Talents.ImprovedImp == 3,
	))
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	if sim.CurrentTime == 0 {
		warlock.petStmBonusSP = 0
	}

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

	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	warlock.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	warlock.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[warlock.Class][int(warlock.Level)]*core.SpellCritRatingPerCritChance)
	warlock.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[warlock.Class][int(warlock.Level)]*core.CritRatingPerCritChance)

	if warlock.Options.Armor == proto.WarlockOptions_DemonArmor {
		armor := map[int32]float64{
			25: 210.0,
			40: 390.0,
			50: 480.0,
			60: 570.0,
		}[warlock.GetCharacter().Level]

		shadowRes := map[int32]float64{
			25: 3.0,
			40: 9.0,
			50: 12.0,
			60: 15.0,
		}[warlock.Level]

		warlock.AddStat(stats.Armor, armor)
		warlock.AddStat(stats.ShadowResistance, shadowRes)
	}

	if warlock.Options.Summon != proto.WarlockOptions_NoSummon {
		warlock.Pet = warlock.NewWarlockPet()
	}

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
