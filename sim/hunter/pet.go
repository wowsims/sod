package hunter

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type HunterPet struct {
	core.Pet

	config PetConfig

	hunterOwner *Hunter

	KillCommandAura *core.Aura

	specialAbility *core.Spell
	focusDump      *core.Spell

	FlankingStrike *core.Spell

	uptimePercent    float64
	hasOwnerCooldown bool
}

func (hunter *Hunter) NewHunterPet() *HunterPet {
	if hunter.Options.PetType == proto.Hunter_Options_PetNone {
		return nil
	}
	if hunter.Options.PetUptime <= 0 {
		return nil
	}
	petConfig := PetConfigs[hunter.Options.PetType]

	hunterPetBaseStats := stats.Stats{}

	baseMinDamage := 0.0
	baseMaxDamage := 0.0
	attackSpeed := 2.0

	switch hunter.Level {
	case 25:
		baseMinDamage = 6.5 * attackSpeed
		baseMaxDamage = 12.5 * attackSpeed
		hunterPetBaseStats = stats.Stats{
			stats.Strength:  53,
			stats.Agility:   45,
			stats.Stamina:   120,
			stats.Intellect: 29,
			stats.Spirit:    39,

			stats.AttackPower: -20,

			// Add 1.8% because pets aren't affected by that component of crit suppression.
			stats.MeleeCrit: (3.2 + 1.8) * core.CritRatingPerCritChance,
		}
	case 40:
		baseMinDamage = 9.5 * attackSpeed
		baseMaxDamage = 15.5 * attackSpeed
		hunterPetBaseStats = stats.Stats{
			stats.Strength:  78,
			stats.Agility:   66,
			stats.Stamina:   160,
			stats.Intellect: 37,
			stats.Spirit:    55,

			stats.AttackPower: -20,

			// Add 1.8% because pets aren't affected by that component of crit suppression.
			stats.MeleeCrit: (3.2 + 1.8) * core.CritRatingPerCritChance,
		}
	case 50:
		// TODO:
		baseMinDamage = 9.5 * attackSpeed
		baseMaxDamage = 15.5 * attackSpeed
		hunterPetBaseStats = stats.Stats{
			stats.Strength:  78,
			stats.Agility:   66,
			stats.Stamina:   160,
			stats.Intellect: 37,
			stats.Spirit:    55,

			stats.AttackPower: -20,

			// Add 1.8% because pets aren't affected by that component of crit suppression.
			stats.MeleeCrit: (3.2 + 1.8) * core.CritRatingPerCritChance,
		}
	case 60:
		// TODO:
		baseMinDamage = 9.5 * attackSpeed
		baseMaxDamage = 15.5 * attackSpeed
		hunterPetBaseStats = stats.Stats{
			stats.Strength:  78,
			stats.Agility:   66,
			stats.Stamina:   160,
			stats.Intellect: 37,
			stats.Spirit:    55,

			stats.AttackPower: -20,

			// Add 1.8% because pets aren't affected by that component of crit suppression.
			stats.MeleeCrit: (3.2 + 1.8) * core.CritRatingPerCritChance,
		}
	}

	hp := &HunterPet{
		Pet:         core.NewPet(petConfig.Name, &hunter.Character, hunterPetBaseStats, hunter.makeStatInheritance(), true, false),
		config:      petConfig,
		hunterOwner: hunter,

		hasOwnerCooldown: petConfig.SpecialAbility == FuriousHowl,
	}

	hp.EnableAutoAttacks(hp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  baseMinDamage,
			BaseDamageMax:  baseMaxDamage,
			SwingSpeed:     attackSpeed,
			CritMultiplier: hp.MeleeCritMultiplier(1, 0),
		},
		AutoSwingMelee: true,
	})

	// After checking numerous logs it seems pet auto attacks are hitting for less then what they should if following standard attack formulas
	// TODO: Figure out from where this difference comes
	// TODO: Phase2 this no longer seems to apply
	//hp.AutoAttacks.MHConfig().DamageMultiplier *= 0.45

	// Happiness
	hp.PseudoStats.DamageDealtMultiplier *= 1.25

	// Family scalars
	hp.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= hp.config.Damage
	hp.PseudoStats.ArmorMultiplier *= hp.config.Armor
	hp.MultiplyStat(stats.Health, hp.config.Health)

	hp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	hp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/62.77)

	//core.ApplyPetConsumeEffects(&hp.Character, hunter.Consumes)

	hunter.AddPet(hp)

	return hp
}

func (hp *HunterPet) GetPet() *core.Pet {
	return &hp.Pet
}

func (hp *HunterPet) Talents() *proto.HunterPetTalents {
	if talents := hp.hunterOwner.Options.PetTalents; talents != nil {
		return talents
	}
	return &proto.HunterPetTalents{}
}

func (hp *HunterPet) Initialize() {
	hp.specialAbility = hp.NewPetAbility(hp.config.SpecialAbility, true)
	hp.focusDump = hp.NewPetAbility(hp.config.FocusDump, false)

	focusRegenMultiplier := (1.0 + 0.1*float64(hp.hunterOwner.Talents.BestialDiscipline)) *
		core.TernaryFloat64(hp.hunterOwner.HasRune(proto.HunterRune_RuneHandsBeastmastery), 1.5, 1.0)

	hp.EnableFocusBar(focusRegenMultiplier, func(sim *core.Simulation) {
		if hp.GCD.IsReady(sim) {
			hp.OnGCDReady(sim)
		}
	})

	// TODO: Make this dynamic
	hp.PseudoStats.DamageDealtMultiplier *= hp.Owner.PseudoStats.DamageDealtMultiplier
}

func (hp *HunterPet) Reset(_ *core.Simulation) {
	hp.uptimePercent = min(1, max(0, hp.hunterOwner.Options.PetUptime))
}

func (hp *HunterPet) ExecuteCustomRotation(sim *core.Simulation) {
	percentRemaining := sim.GetRemainingDurationPercent()
	if percentRemaining < 1.0-hp.uptimePercent { // once fight is % completed, disable pet.
		hp.Disable(sim)
		return
	}

	if hp.hasOwnerCooldown && hp.CurrentFocus() < 50 {
		// When a major ability (Furious Howl or Savage Rend) is ready, pool enough
		// energy to use on-demand.
		return
	}

	target := hp.CurrentTarget

	if hp.focusDump == nil {
		if hp.specialAbility.CanCast(sim, target) {
			hp.specialAbility.Cast(sim, target)
		}
		return
	}
	if hp.specialAbility == nil {
		if hp.focusDump.CanCast(sim, target) {
			hp.focusDump.Cast(sim, target)
		}
		return
	}

	if hp.config.RandomSelection {
		if sim.RandomFloat("Hunter Pet Ability") < 0.5 {
			_ = hp.specialAbility.Cast(sim, target) || hp.focusDump.Cast(sim, target)
		} else {
			_ = hp.focusDump.Cast(sim, target) || hp.specialAbility.Cast(sim, target)
		}
	} else {
		if hp.specialAbility.IsReady(sim) {
			_ = hp.specialAbility.Cast(sim, target)
		} else {
			_ = hp.focusDump.Cast(sim, target)
		}
	}
}

func (hp *HunterPet) killCommandMult() float64 {
	if hp.KillCommandAura == nil {
		return 1
	}
	return 1 + 0.2*float64(hp.KillCommandAura.GetStacks())
}

const PetExpertiseScale = 3.25

func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// EJ posts claim this value is passed through math.Floor, but in-game testing
		// shows pets benefit from each point of owner hit rating in WotLK Classic.
		// https://web.archive.org/web/20120112003252/http://elitistjerks.com/f80/t100099-demonology_releasing_demon_you
		ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
		hitRatingFromOwner := ownerHitChance * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.RangedAttackPower] * 0.22,

			stats.MeleeCrit: ownerStats[stats.MeleeCrit],
			stats.SpellCrit: ownerStats[stats.MeleeCrit],

			stats.MeleeHit: hitRatingFromOwner,
			stats.SpellHit: hitRatingFromOwner * 2,
			//stats.Expertise: ownerHitChance * PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,
		}
	}
}

type PetConfig struct {
	Name string

	SpecialAbility PetAbilityType
	FocusDump      PetAbilityType

	Health float64
	Armor  float64
	Damage float64

	// Randomly select between abilities instead of using a prio.
	RandomSelection bool
}

// Abilities reference: https://wotlk.wowhead.com/hunter-pets
// https://wotlk.wowhead.com/guides/hunter-dps-best-pets-taming-loyalty-burning-crusade-classic
var PetConfigs = map[proto.Hunter_Options_PetType]PetConfig{
	proto.Hunter_Options_Cat: {
		Name:           "Cat",
		SpecialAbility: Bite,
		FocusDump:      Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.1,
	},
	proto.Hunter_Options_WindSerpent: {
		Name:           "Wind Serpent",
		SpecialAbility: Unknown,
		FocusDump:      LightningBreath,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.07,
	},
	proto.Hunter_Options_Bat: {
		Name: "Bat",
		//SpecialAbility: SonicBlast,
		FocusDump: Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Bear: {
		Name:           "Bear",
		SpecialAbility: Swipe,
		FocusDump:      Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Boar: {
		Name: "Boar",
		//SpecialAbility: Gore,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_CarrionBird: {
		Name:           "Carrion Bird",
		SpecialAbility: DemoralizingScreech,
		FocusDump:      Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Chimaera: {
		Name: "Chimaera",
		//SpecialAbility: FroststormBreath,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_CoreHound: {
		Name: "Core Hound",
		//SpecialAbility: LavaBreath,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Crab: {
		Name: "Crab",
		//SpecialAbility: Pin,
		FocusDump: Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Crocolisk: {
		Name: "Crocolisk",
		//SpecialAbility: BadAttitude,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Devilsaur: {
		Name: "Devilsaur",
		//SpecialAbility: MonstrousBite,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Dragonhawk: {
		Name: "Dragonhawk",
		//SpecialAbility: FireBreath,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Gorilla: {
		Name: "Gorilla",
		//SpecialAbility: Pummel,
		//FocusDump: Smack,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Hyena: {
		Name: "Hyena",
		//SpecialAbility: TendonRip,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Raptor: {
		Name: "Raptor",
		//SpecialAbility: SavageRend,
		FocusDump: Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Scorpid: {
		Name:           "Scorpid",
		SpecialAbility: ScorpidPoison,
		FocusDump:      Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Serpent: {
		Name: "Serpent",
		//SpecialAbility: PoisonSpit,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Silithid: {
		Name: "Silithid",
		//SpecialAbility: VenomWebSpray,
		FocusDump: Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Spider: {
		Name: "Spider",
		//SpecialAbility:   Web,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_SpiritBeast: {
		Name: "Spirit Beast",
		//SpecialAbility: SpiritStrike,
		FocusDump: Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_SporeBat: {
		Name: "Spore Bat",
		//SpecialAbility: SporeCloud,
		//FocusDump:      Smack,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Tallstrider: {
		Name: "Tallstrider",
		//SpecialAbility:   DustCloud,
		FocusDump: Claw,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Turtle: {
		Name: "Turtle",
		//SpecialAbility: ShellShield,
		FocusDump: Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
	proto.Hunter_Options_Wolf: {
		Name:           "Wolf",
		SpecialAbility: FuriousHowl,
		FocusDump:      Bite,

		Health: 1.0,
		Armor:  1.0,
		Damage: 1.0,
	},
}
