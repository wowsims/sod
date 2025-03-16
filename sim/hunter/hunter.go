package hunter

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{16, 14, 16}

const (
	SpellFlagStrike = core.SpellFlagAgentReserved1
	SpellFlagSting  = core.SpellFlagAgentReserved2
)

const (
	ClassSpellMask_HunterNone uint64 = 0

	// Shots
	ClassSpellMask_HunterAimedShot uint64 = 1 << iota
	ClassSpellMask_HunterArcaneShot
	ClassSpellMask_HunterChimeraShot
	ClassSpellMask_HunterExplosiveShot
	ClassSpellMask_HunterKillShot
	ClassSpellMask_HunterMultiShot
	ClassSpellMask_HunterSteadyShot

	// Strikes
	ClassSpellMask_HunterFlankingStrike
	ClassSpellMask_HunterRaptorStrike
	ClassSpellMask_HunterRaptorStrikeHit
	ClassSpellMask_HunterWyvernStrike

	// Stings
	ClassSpellMask_HunterSerpentSting

	// Traps
	ClassSpellMask_HunterExplosiveTrap
	ClassSpellMask_HunterFreezingTrap
	ClassSpellMask_HunterImmolationTrap

	// Other
	ClassSpellMask_HunterCarve
	ClassSpellMask_HunterCarveHit
	ClassSpellMask_HunterMongooseBite
	ClassSpellMask_HunterWingClip
	ClassSpellMask_HunterVolley
	ClassSpellMask_HunterChimeraSerpent
	ClassSpellMask_HunterHuntersMark

	// Pet Spells
	ClassSpellMask_HunterPetFlankingStrike
	ClassSpellMask_HunterPetClaw
	ClassSpellMask_HunterPetBite
	ClassSpellMask_HunterPetLightningBreath
	ClassSpellMask_HunterPetLavaBreath
	ClassSpellMask_HunterPetScreech
	ClassSpellMask_HunterPetScorpidPoison

	ClassSpellMask_HunterAll = 1<<iota - 1

	ClassSpellMask_HunterTraps   = ClassSpellMask_HunterExplosiveTrap | ClassSpellMask_HunterFreezingTrap | ClassSpellMask_HunterImmolationTrap
	ClassSpellMask_HunterShots   = ClassSpellMask_HunterAimedShot | ClassSpellMask_HunterArcaneShot | ClassSpellMask_HunterChimeraShot | ClassSpellMask_HunterExplosiveShot | ClassSpellMask_HunterKillShot | ClassSpellMask_HunterMultiShot | ClassSpellMask_HunterSteadyShot
	ClassSpellMask_HunterStrikes = ClassSpellMask_HunterFlankingStrike | ClassSpellMask_HunterRaptorStrike | ClassSpellMask_HunterRaptorStrikeHit | ClassSpellMask_HunterWyvernStrike
)

func RegisterHunter() {
	core.RegisterAgentFactory(
		proto.Player_Hunter{},
		proto.Spec_SpecHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Hunter)
			if !ok {
				panic("Invalid spec value for Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

type Hunter struct {
	core.Character

	Talents *proto.HunterTalents
	Options *proto.Hunter_Options

	pet *HunterPet

	AmmoDPS                   float64
	AmmoDamageBonus           float64
	NormalizedAmmoDamageBonus float64

	// Miscellaneous set bonuses that require extra logic inside of spells
	SerpentStingAPCoeff float64

	curQueueAura       *core.Aura
	curQueuedAutoSpell *core.Spell

	AimedShot               *core.Spell
	ArcaneShot              *core.Spell
	ChimeraShot             *core.Spell
	ExplosiveShot           *core.Spell
	ExplosiveTrap           *core.Spell
	ImmolationTrap          *core.Spell
	FreezingTrap            *core.Spell
	KillCommand             *core.Spell
	KillShot                *core.Spell
	MultiShot               *core.Spell
	FocusFire               *core.Spell
	RapidFire               *core.Spell
	RaptorStrike            *core.Spell
	RaptorStrikeMH          *core.Spell
	RaptorStrikeOH          *core.Spell
	FlankingStrike          *core.Spell
	WyvernStrike            *core.Spell
	MongooseBite            *core.Spell
	ScorpidSting            *core.Spell
	SerpentSting            *core.Spell
	SerpentStingChimeraShot *core.Spell
	SilencingShot           *core.Spell
	SteadyShot              *core.Spell
	Volley                  *core.Spell
	CarveMH                 *core.Spell
	CarveOH                 *core.Spell
	WingClip                *core.Spell
	HuntersMark             *core.Spell

	Shots       []*core.Spell
	Strikes     []*core.Spell
	MeleeSpells []*core.Spell
	LastShot    *core.Spell

	FlankingStrikeAura *core.Aura
	RaptorFuryAura     *core.Aura
	SniperTrainingAura *core.Aura
	CobraStrikesAura   *core.Aura
	HitAndRunAura      *core.Aura

	// The aura that allows you to cast Mongoose Bite
	DefensiveState *core.Aura

	ImprovedSteadyShotAura *core.Aura
	LockAndLoadAura        *core.Aura
	RapidFireAura          *core.Aura
	BestialWrathPetAura    *core.Aura

	HuntersMarkAuras core.AuraArray
}

func (hunter *Hunter) GetCharacter() *core.Character {
	return &hunter.Character
}

func (hunter *Hunter) GetHunter() *Hunter {
	return hunter
}

func (hunter *Hunter) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if raidBuffs.TrueshotAura && hunter.Talents.TrueshotAura {
		hunter.AddStat(stats.RangedAttackPower, map[int32]float64{
			25: 0,
			40: 100,
			50: 150,
			60: 200,
		}[hunter.Level])
	}

	raidBuffs.AspectOfTheLion = true
	// Hunter gains an additional 10% stats from Aspect of the Lion
	statMultiply := 1.1
	hunter.MultiplyStat(stats.Strength, statMultiply)
	hunter.MultiplyStat(stats.Stamina, statMultiply)
	hunter.MultiplyStat(stats.Agility, statMultiply)
	hunter.MultiplyStat(stats.Intellect, statMultiply)
	hunter.MultiplyStat(stats.Spirit, statMultiply)
}
func (hunter *Hunter) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (hunter *Hunter) Initialize() {
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(ClassSpellMask_HunterShots) {
			hunter.Shots = append(hunter.Shots, spell)
		}
	})
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagStrike) {
			hunter.Strikes = append(hunter.Strikes, spell)
		}
	})
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) {
			hunter.MeleeSpells = append(hunter.MeleeSpells, spell)
		}
	})

	hunter.registerAspectOfTheHawkSpell()
	hunter.registerAspectOfTheFalconSpell()
	hunter.registerAspectOfTheViperSpell()
	hunter.registerHuntersMark()

	multiShotTimer := hunter.NewTimer()
	arcaneShotTimer := hunter.NewTimer()

	hunter.registerSerpentStingSpell()

	hunter.registerArcaneShotSpell(arcaneShotTimer)
	hunter.registerAimedShotSpell(arcaneShotTimer)
	hunter.registerMultiShotSpell(multiShotTimer)
	hunter.registerExplosiveShotSpell()
	hunter.registerChimeraShotSpell()
	hunter.registerSteadyShotSpell()
	hunter.registerKillShotSpell()

	hunter.registerRaptorStrikeSpell()
	hunter.registerFlankingStrikeSpell()
	hunter.registerWyvernStrikeSpell()
	hunter.registerMongooseBiteSpell()
	hunter.registerCarveSpell()
	hunter.registerWingClipSpell()
	hunter.registerVolleySpell()

	// Trap Launcher rune also splits the cooldowns between frost traps and fire traps, without the rune all traps share a cd
	if hunter.HasRune(proto.HunterRune_RuneBootsTrapLauncher) {
		fireTraps := hunter.NewTimer()
		frostTraps := hunter.NewTimer()

		hunter.registerExplosiveTrapSpell(fireTraps)
		hunter.registerImmolationTrapSpell(fireTraps)
		hunter.registerFreezingTrapSpell(frostTraps)
	} else {
		traps := hunter.NewTimer()

		hunter.registerExplosiveTrapSpell(traps)
		hunter.registerImmolationTrapSpell(traps)
		hunter.registerFreezingTrapSpell(traps)
	}

	// hunter.registerKillCommand()
	hunter.registerRapidFire()
	hunter.registerFocusFireSpell()
}

func (hunter *Hunter) Reset(sim *core.Simulation) {
}

func NewHunter(character *core.Character, options *proto.Player) *Hunter {
	hunterOptions := options.GetHunter()

	hunter := &Hunter{
		Character: *character,
		Talents:   &proto.HunterTalents{},
		Options:   hunterOptions.Options,
	}
	core.FillTalentsProto(hunter.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	hunter.EnableManaBar()

	hunter.PseudoStats.CanParry = true

	rangedWeapon := hunter.WeaponFromRanged()

	if hunter.HasRangedWeapon() {
		// Ammo
		switch hunter.Options.Ammo {
		case proto.Hunter_Options_RazorArrow:
			hunter.AmmoDPS = 7.5
		case proto.Hunter_Options_SolidShot:
			hunter.AmmoDPS = 7.5
		case proto.Hunter_Options_JaggedArrow:
			hunter.AmmoDPS = 13
		case proto.Hunter_Options_AccurateSlugs:
			hunter.AmmoDPS = 13
		case proto.Hunter_Options_MithrilGyroShot:
			hunter.AmmoDPS = 15
		case proto.Hunter_Options_IceThreadedArrow:
			hunter.AmmoDPS = 16.5
		case proto.Hunter_Options_IceThreadedBullet:
			hunter.AmmoDPS = 16.5
		case proto.Hunter_Options_ThoriumHeadedArrow:
			hunter.AmmoDPS = 17.5
		case proto.Hunter_Options_ThoriumShells:
			hunter.AmmoDPS = 17.5
		case proto.Hunter_Options_RockshardPellets:
			hunter.AmmoDPS = 18
		case proto.Hunter_Options_Doomshot:
			hunter.AmmoDPS = 20
		case proto.Hunter_Options_MiniatureCannonBalls:
			hunter.AmmoDPS = 20.5
		}
		hunter.AmmoDamageBonus = hunter.AmmoDPS * rangedWeapon.SwingSpeed
		hunter.NormalizedAmmoDamageBonus = hunter.AmmoDPS * 2.8

		// Quiver
		switch hunter.Options.QuiverBonus {
		case proto.Hunter_Options_Speed10:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.1
		case proto.Hunter_Options_Speed11:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.11
		case proto.Hunter_Options_Speed12:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.12
		case proto.Hunter_Options_Speed13:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.13
		case proto.Hunter_Options_Speed14:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.14
		case proto.Hunter_Options_Speed15:
			hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
		}
	}

	hunter.EnableAutoAttacks(hunter, core.AutoAttackOptions{
		MainHand:        hunter.WeaponFromMainHand(),
		OffHand:         hunter.WeaponFromOffHand(),
		Ranged:          rangedWeapon,
		ReplaceMHSwing:  hunter.TryRaptorStrike,
		AutoSwingRanged: true,
		AutoSwingMelee:  true,
	})

	hunter.AutoAttacks.RangedConfig().Flags |= core.SpellFlagCastTimeNoGCD
	hunter.AutoAttacks.RangedConfig().Cast = core.CastConfig{
		DefaultCast: core.Cast{
			CastTime: time.Millisecond * 500,
		},
		ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
			cast.CastTime = spell.CastTime()
		},
		CastTime: func(spell *core.Spell) time.Duration {
			return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed() * math.Max(spell.CastTimeMultiplier, 0))
		},
	}
	hunter.AutoAttacks.RangedConfig().ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return !hunter.IsCasting(sim)
	}
	hunter.AutoAttacks.RangedConfig().CritDamageBonus = hunter.mortalShots()
	hunter.AutoAttacks.RangedConfig().BonusCoefficient = 1
	hunter.AutoAttacks.RangedConfig().ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := hunter.RangedWeaponDamage(sim, spell.RangedAttackPower(target, false)) +
			hunter.AmmoDamageBonus
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

		splitIdx := spell.GetMetricsSplitIdx()

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			newSplitIdx := spell.GetMetricsSplitIdx()

			// We have to dynamically update the split metrics to ensure extra attacks' damage are categorized correctly
			spell.SetMetricsSplit(splitIdx)
			spell.DealDamage(sim, result)
			spell.SetMetricsSplit(newSplitIdx)
		})
	}

	hunter.pet = hunter.NewHunterPet()

	hunter.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	hunter.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	hunter.AddStatDependency(stats.Agility, stats.RangedAttackPower, 2)
	hunter.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class][int(character.Level)]*core.CritRatingPerCritChance)
	hunter.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class][int(character.Level)]*core.SpellCritRatingPerCritChance)

	guardians.ConstructGuardians(&hunter.Character)

	return hunter
}

func (hunter *Hunter) HasRune(rune proto.HunterRune) bool {
	return hunter.HasRuneById(int32(rune))
}

func (hunter *Hunter) baseRuneAbilityDamage() float64 {
	return 2.976264 + 0.641066*float64(hunter.Level) + 0.022519*float64(hunter.Level*hunter.Level)
}

func (hunter *Hunter) OnGCDReady(_ *core.Simulation) {
}

// Agent is a generic way to access underlying hunter on any of the agents.
type HunterAgent interface {
	GetHunter() *Hunter
}
