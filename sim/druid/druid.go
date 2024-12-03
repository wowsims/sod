package druid

import (
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SpellFlagOmen    = core.SpellFlagAgentReserved1
	SpellFlagBuilder = core.SpellFlagAgentReserved2
)

var TalentTreeSizes = [3]int{16, 16, 15}

const (
	SpellCode_DruidNone int32 = iota

	SpellCode_DruidFaerieFire
	SpellCode_DruidFaerieFireFeral
	SpellCode_DruidFerociousBite
	SpellCode_DruidInsectSwarm
	SpellCode_DruidMangleCat
	SpellCode_DruidMangleBear
	SpellCode_DruidSwipeCat
	SpellCode_DruidLacerateDirect
	SpellCode_DruidLacerateDot
	SpellCode_DruidSwipeBear
	SpellCode_DruidMaul
	SpellCode_DruidMoonfire
	SpellCode_DruidRake
	SpellCode_DruidRip
	SpellCode_DruidShred
	SpellCode_DruidStarfire
	SpellCode_DruidStarsurge
	SpellCode_DruidWrath
	SpellCode_DruidStarfall
	SpellCode_DruidStarfallTick
	SpellCode_DruidStarfallSplash
	SpellCode_DruidSunfire
)

type Druid struct {
	core.Character
	SelfBuffs

	Talents *proto.DruidTalents

	DruidSpells []*DruidSpell

	StartingForm DruidForm

	RebirthTiming     float64
	BleedsActive      int
	AssumeBleedActive bool

	ReplaceBearMHFunc core.ReplaceMHSwing

	Barkskin             *DruidSpell
	DemoralizingRoar     *DruidSpell
	Enrage               *DruidSpell
	FaerieFire           *DruidSpell
	FerociousBite        *DruidSpell
	ForceOfNature        *DruidSpell
	FrenziedRegeneration *DruidSpell
	GiftOfTheWild        *DruidSpell
	Hurricane            []*DruidSpell
	Innervate            *DruidSpell
	InsectSwarm          []*DruidSpell
	Lacerate             *DruidSpell
	LacerateBleed        *DruidSpell
	Languish             *DruidSpell
	MangleBear           *DruidSpell
	MangleCat            *DruidSpell
	Berserk              *DruidSpell
	Maul                 *DruidSpell
	MaulQueueSpell       *DruidSpell
	Moonfire             []*DruidSpell
	Rebirth              *DruidSpell
	Rake                 *DruidSpell
	Rip                  *DruidSpell
	SavageRoar           *DruidSpell
	Shred                *DruidSpell
	Starfire             []*DruidSpell
	Starfall             *DruidSpell
	StarfallTick         *DruidSpell
	StarfallSplash       *DruidSpell
	Starsurge            *DruidSpell
	Sunfire              *DruidSpell
	SunfireCat           *DruidSpell
	SurvivalInstincts    *DruidSpell
	SwipeBear            *DruidSpell
	SwipeCat             *DruidSpell
	TigersFury           *DruidSpell
	Typhoon              *DruidSpell
	curQueuedAutoSpell   *DruidSpell
	MaulQueue            *DruidSpell
	Wrath                []*DruidSpell

	BearForm    *DruidSpell
	CatForm     *DruidSpell
	MoonkinForm *DruidSpell

	BarkskinAura             *core.Aura
	BearFormAura             *core.Aura
	BerserkAura              *core.Aura
	CatFormAura              *core.Aura
	curQueueAura             *core.Aura
	ClearcastingAura         *core.Aura
	DemoralizingRoarAuras    core.AuraArray
	DreamstateManaRegenAura  *core.Aura
	EnrageAura               *core.Aura
	EclipseAura              *core.Aura
	FaerieFireAuras          core.AuraArray
	ImprovedFaerieFireAuras  core.AuraArray
	FrenziedRegenerationAura *core.Aura
	FurorAura                *core.Aura
	FuryOfStormrageAura      *core.Aura
	InsectSwarmAuras         core.AuraArray
	MaulQueueAura            *core.Aura
	MoonkinFormAura          *core.Aura
	NaturesGraceProcAura     *core.Aura
	PredatoryInstinctsAura   *core.Aura
	SurvivalInstinctsAura    *core.Aura
	TigersFuryAura           *core.Aura
	SavageRoarAura           *core.Aura
	SolarEclipseProcAura     *core.Aura
	LunarEclipseProcAura     *core.Aura
	WildStrikesBuffAura      *core.Aura
	Tank2PieceAqAura         *core.Aura
	Tank2PieceAqProcAura     *core.Aura

	BleedCategories         core.ExclusiveCategoryArray
	SavageRoarDurationTable [6]time.Duration

	// Extra data used for various calculations and overrides
	FerociousBiteExcessEnergyOverride bool // When true, disables the excess energy consumption of Ferocious bite
	// Sunfire/Moonfire modifiers applied when in Moonkin form
	MoonfireDotMultiplier float64
	ShredPositionOverride bool
	SunfireDotMultiplier  float64
	t26pcTreants          *T2Treants

	form         DruidForm
	disabledMCDs []*core.MajorCooldown
}

type SelfBuffs struct {
	InnervateTarget *proto.UnitReference
}

func (druid *Druid) GetCharacter() *core.Character {
	return &druid.Character
}

func (druid *Druid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if (raidBuffs.GiftOfTheWild == proto.TristateEffect_TristateEffectRegular) && (druid.Talents.ImprovedMarkOfTheWild > 0) {
		druid.AddStats(core.BuffSpellByLevel[core.MarkOfTheWild][druid.Level].Multiply(0.07 * float64(druid.Talents.ImprovedMarkOfTheWild)))
	}

	// TODO: These should really be aura attached to the actual forms
	if druid.InForm(Moonkin) {
		raidBuffs.MoonkinAura = true
	}

	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
		raidBuffs.LeaderOfThePack = true
	}
}

func (druid *Druid) TryMaul(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !druid.curQueueAura.IsActive() {
		return mhSwingSpell
	}

	if !druid.curQueuedAutoSpell.CanCast(sim, druid.CurrentTarget) {
		druid.curQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	return druid.curQueuedAutoSpell.Spell
}

func (druid *Druid) RegisterSpell(formMask DruidForm, config core.SpellConfig) *DruidSpell {
	prev := config.ExtraCastCondition
	prevModify := config.Cast.ModifyCast

	ds := &DruidSpell{FormMask: formMask}
	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		// Check if we're in allowed form to cast
		// Allow 'humanoid' auto unshift casts
		if (ds.FormMask != Any && !druid.InForm(ds.FormMask)) && !ds.FormMask.Matches(Humanoid) {
			if sim.Log != nil {
				sim.Log("Failed cast to spell %s, wrong form", ds.ActionID)
			}
			return false
		}
		return prev == nil || prev(sim, target)
	}
	config.Cast.ModifyCast = func(sim *core.Simulation, s *core.Spell, c *core.Cast) {
		if !druid.InForm(ds.FormMask) && ds.FormMask.Matches(Humanoid) {
			druid.CancelShapeshift(sim)
		}
		if prevModify != nil {
			prevModify(sim, s, c)
		}
	}

	ds.Spell = druid.Unit.RegisterSpell(config)
	druid.DruidSpells = append(druid.DruidSpells, ds)

	return ds
}

func (druid *Druid) Initialize() {
	druid.BleedCategories = druid.GetEnemyExclusiveCategories(core.BleedEffectCategory)

	druid.registerFaerieFireSpell()
	druid.registerInnervateCD()
	druid.registerCatnipCD()
}

func (druid *Druid) RegisterBalanceSpells() {
	druid.registerHurricaneSpell()
	druid.registerInsectSwarmSpell()
	druid.registerMoonfireSpell()
	druid.registerStarfireSpell()
	druid.registerWrathSpell()
}

// TODO: Classic feral
func (druid *Druid) RegisterFeralCatSpells() {
	druid.registerCatFormSpell()
	// druid.registerBearFormSpell()
	// druid.registerEnrageSpell()
	druid.registerFerociousBiteSpell()
	// druid.registerMangleBearSpell()
	// druid.registerMaulSpell()
	// druid.registerLacerateSpell()
	druid.registerRakeSpell()
	druid.registerRipSpell()
	druid.registerShredSpell()
	// druid.registerSwipeBearSpell()
	druid.registerTigersFurySpell()
}

// TODO: Classic feral tank
func (druid *Druid) RegisterFeralTankSpells(realismICD *core.Cooldown) {
	// druid.registerBarkskinCD()
	// druid.registerBerserkCD()
	druid.registerBearFormSpell()
	// druid.registerDemoralizingRoarSpell()
	druid.registerEnrageSpell()
	druid.registerFrenziedRegenerationCD()
	druid.registerMangleBearSpell()
	druid.registerRakeSpell()
	druid.registerRipSpell()
	druid.registerMaulSpell(realismICD)
	druid.registerLacerateBleedSpell()
	druid.registerLacerateSpell()
	// druid.registerSurvivalInstinctsCD()
	druid.registerSwipeBearSpell()
}

func (druid *Druid) Reset(_ *core.Simulation) {
	druid.BleedsActive = 0
	druid.form = druid.StartingForm
	druid.disabledMCDs = []*core.MajorCooldown{}
}

func New(character *core.Character, form DruidForm, selfBuffs SelfBuffs, talents string) *Druid {
	druid := &Druid{
		Character:    *character,
		SelfBuffs:    selfBuffs,
		Talents:      &proto.DruidTalents{},
		StartingForm: form,
		form:         form,
	}
	core.FillTalentsProto(druid.Talents.ProtoReflect(), talents, TalentTreeSizes)
	druid.EnableManaBar()

	druid.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	druid.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(druid.Level)]*core.CritRatingPerCritChance)
	druid.AddStatDependency(stats.Agility, stats.Dodge, core.DodgePerAgiAtLevel[character.Class][int(druid.Level)]*core.DodgeRatingPerDodgeChance)
	druid.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(druid.Level)]*core.SpellCritRatingPerCritChance)
	druid.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// Druids get extra melee haste
	// druid.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	guardians.ConstructGuardians(&druid.Character)
	druid.t26pcTreants = druid.NewT2Treants()

	return druid
}

type DruidSpell struct {
	*core.Spell
	FormMask DruidForm
}

func (ds *DruidSpell) IsReady(sim *core.Simulation) bool {
	if ds == nil {
		return false
	}
	return ds.Spell.IsReady(sim)
}

func (ds *DruidSpell) CanCast(sim *core.Simulation, target *core.Unit) bool {
	if ds == nil {
		return false
	}
	return ds.Spell.CanCast(sim, target)
}

func (ds *DruidSpell) IsEqual(s *core.Spell) bool {
	if ds == nil || s == nil {
		return false
	}
	return ds.Spell == s
}

func (druid *Druid) HasRune(rune proto.DruidRune) bool {
	return druid.HasRuneById(int32(rune))
}

func (druid *Druid) baseRuneAbilityDamage() float64 {
	return 9.183105 + 0.616405*float64(druid.Level) + 0.028608*float64(druid.Level*druid.Level)
}

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type DruidAgent interface {
	GetDruid() *Druid
}
