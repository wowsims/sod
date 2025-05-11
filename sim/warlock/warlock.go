package warlock

import (
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
	ClassSpellMask_WarlockNone uint64 = 0

	ClassSpellMask_WarlockChaosBolt uint64 = 1 << iota
	ClassSpellMask_WarlockConflagrate
	ClassSpellMask_WarlockCorruption
	ClassSpellMask_WarlockCurseOfAgony
	ClassSpellMask_WarlockCurseOfDoom
	ClassSpellMask_WarlockCurseOfElements
	ClassSpellMask_WarlockCurseOfRecklessness
	ClassSpellMask_WarlockCurseOfShadow
	ClassSpellMask_WarlockDeathCoil
	ClassSpellMask_WarlockDemonicSacrifice
	ClassSpellMask_WarlockDrainLife
	ClassSpellMask_WarlockDrainSoul
	ClassSpellMask_WarlockHaunt
	ClassSpellMask_WarlockImmolate
	ClassSpellMask_WarlockIncinerate
	ClassSpellMask_WarlockLifeTap
	ClassSpellMask_WarlockSearingPain
	ClassSpellMask_WarlockShadowflame
	ClassSpellMask_WarlockShadowCleave
	ClassSpellMask_WarlockShadowBolt
	ClassSpellMask_WarlockShadowburn
	ClassSpellMask_WarlockSiphonLife
	ClassSpellMask_WarlockSoulFire
	ClassSpellMask_WarlockUnstableAffliction
	ClassSpellMask_WarlockInfernalArmor
	ClassSpellMask_WarlockDemonicGrace

	ClassSpellMask_WarlockRainOfFire
	ClassSpellMask_WarlockImmolationAura
	ClassSpellMask_WarlockImmolationAuraProc

	ClassSpellMask_WarlockSummonFelguard
	ClassSpellMask_WarlockSummonFelguardCleave
	ClassSpellMask_WarlockSummonFelhunter
	ClassSpellMask_WarlockSummonImp
	ClassSpellMask_WarlockSummonImpFireBolt
	ClassSpellMask_WarlockSummonSuccubus
	ClassSpellMask_WarlockSummonSuccubusLashOfPain
	ClassSpellMask_WarlockSummonVoidwalker

	ClassSpellMask_WarlockAll = 1<<iota - 1

	ClassSpellMask_WarlockCurses = ClassSpellMask_WarlockCurseOfAgony | ClassSpellMask_WarlockCurseOfDoom |
		ClassSpellMask_WarlockCurseOfRecklessness | ClassSpellMask_WarlockCurseOfElements | ClassSpellMask_WarlockCurseOfShadow

	ClassSpellMask_WarlockSummons = ClassSpellMask_WarlockSummonFelguard |
		ClassSpellMask_WarlockSummonFelhunter |
		ClassSpellMask_WarlockSummonImp |
		ClassSpellMask_WarlockSummonSuccubus |
		ClassSpellMask_WarlockSummonVoidwalker

	// TODO: Hellfire
	ClassSpellMask_WarlockHarmfulGCDSpells = ClassSpellMask_WarlockShadowBolt | ClassSpellMask_WarlockSoulFire | ClassSpellMask_WarlockConflagrate |
		ClassSpellMask_WarlockSearingPain | ClassSpellMask_WarlockImmolate | ClassSpellMask_WarlockRainOfFire |
		ClassSpellMask_WarlockCorruption | ClassSpellMask_WarlockCurseOfAgony | ClassSpellMask_WarlockCurseOfDoom |
		ClassSpellMask_WarlockSiphonLife | ClassSpellMask_WarlockDrainSoul | ClassSpellMask_WarlockDrainLife |
		ClassSpellMask_WarlockDeathCoil |
		ClassSpellMask_WarlockChaosBolt | ClassSpellMask_WarlockIncinerate | ClassSpellMask_WarlockShadowflame |
		ClassSpellMask_WarlockHaunt | ClassSpellMask_WarlockUnstableAffliction |
		ClassSpellMask_WarlockSummonFelguardCleave | ClassSpellMask_WarlockSummonImpFireBolt | ClassSpellMask_WarlockSummonSuccubusLashOfPain
)

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	BasePets   []*WarlockPet
	Felhunter  *WarlockPet
	Felguard   *WarlockPet
	Imp        *WarlockPet
	Succubus   *WarlockPet
	Voidwalker *WarlockPet
	// Doomguard *DoomguardPet
	// Infernal  *InfernalPet

	ActivePet     *WarlockPet // The Warlock's current pet
	SacrificedPet *WarlockPet // Stored reference to the Warlock's most recently-sacrified pet

	ChaosBolt          *core.Spell
	Conflagrate        []*core.Spell
	Corruption         []*core.Spell
	DarkPact           *core.Spell
	DrainSoul          []*core.Spell
	Haunt              *core.Spell
	Immolate           []*core.Spell
	Incinerate         *core.Spell
	InfernalArmor      *core.Spell
	InvocationSpellMap map[uint64]*core.Spell
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

	ActiveCurseAura          core.AuraArray
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

	AmplifyCurseAura        *core.Aura
	BackdraftAura           *core.Aura
	defendersResolveAura    *core.Aura
	DecimationAura          *core.Aura
	DemonicGraceAura        *core.Aura
	DemonicKnowledgeAura    *core.Aura
	HauntDebuffAuras        core.AuraArray
	ImmolationAura          *core.Spell
	ImprovedShadowBoltAuras core.AuraArray
	IncinerateAura          *core.Aura
	MarkOfChaosAuras        core.AuraArray
	MasterDemonologistAura  *core.Aura
	Metamorphosis           *core.Spell
	MetamorphosisAura       *core.Aura
	PyroclasmAura           *core.Aura
	ShadowTranceAura        *core.Aura
	SoulLinkAura            *core.Aura
	VengeanceAura           *core.Aura
	zilaGularAura           *core.Aura

	// The sum total of demonic pact spell power * seconds.
	DPSPAggregate float64

	// Extra state and logic variables
	activeEffects                map[int32]int32 // Used by the 6pT2 DPS bonus
	backdraftCastSpeed           float64
	demonicKnowledgeSp           float64
	maintainBuffsOnSacrifice     bool    // Whether to disable the Master Demonologist and Demonic Sacrifice buffs when sacrificing/summoning pets. Used by TAQ 4pc
	masterDemonologistMultiplier float64 // Bonus multiplier applied to the Master Demonologist talent
	nightfallProcChance          float64
	// For effects that buff the damage of shadow bolt for each active Warlock effect on the target, e.g. 2pc DPS 6pc
	shadowBoltActiveEffectMultiplierPer float64
	shadowBoltActiveEffectMultiplierMax float64
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) Initialize() {
	warlock.activeEffects = make(map[int32]int32, len(warlock.Env.AllUnits))

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
	warlock.SacrificedPet = nil
	warlock.ActiveCurseAura = make([]*core.Aura, len(sim.Environment.AllUnits))

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
	warlock.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[warlock.Class][int(warlock.Level)]*core.CritRatingPerCritChance)
	warlock.AddStatDependency(stats.Agility, stats.Dodge, core.DodgePerAgiAtLevel[character.Class][int(warlock.Level)]*core.DodgeRatingPerDodgeChance)
	warlock.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[warlock.Class][int(warlock.Level)]*core.SpellCritRatingPerCritChance)
	warlock.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

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
