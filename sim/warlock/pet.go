package warlock

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type OnPetDisable func(sim *core.Simulation, isSacrifice bool)

type WarlockPet struct {
	core.Pet

	OnPetDisable OnPetDisable

	owner *Warlock

	primaryAbility   *core.Spell
	secondaryAbility *core.Spell

	SoulLinkAura           *core.Aura
	DemonicEmpowermentAura *core.Aura

	DanceOfTheWickedManaMetrics *core.ResourceMetrics
	LifeTapManaMetrics          *core.ResourceMetrics
	T1Tank4PManaMetrics         *core.ResourceMetrics // https://www.wowhead.com/classic/spell=457572/s03-item-t1-warlock-tank-4p-bonus

	manaPooling bool
}

type PetConfig struct {
	Name          string
	PowerModifier float64 // GetUnitPowerModifier("pet")
	Stats         stats.Stats
	AutoAttacks   core.AutoAttackOptions
}

func (warlock *Warlock) setDefaultActivePet() {
	switch warlock.Options.Summon {
	case proto.WarlockOptions_Imp:
		warlock.ActivePet = warlock.Imp
	case proto.WarlockOptions_Felguard:
		warlock.ActivePet = warlock.Felguard
	case proto.WarlockOptions_Felhunter:
		warlock.ActivePet = warlock.Felhunter
	case proto.WarlockOptions_Succubus:
		warlock.ActivePet = warlock.Succubus
	case proto.WarlockOptions_Voidwalker:
		warlock.ActivePet = warlock.Voidwalker
	}
}

func (warlock *Warlock) changeActivePet(sim *core.Simulation, newPet *WarlockPet, isSacrifice bool) {
	if warlock.ActivePet != nil {
		warlock.ActivePet.Disable(sim, isSacrifice)

		if isSacrifice {
			warlock.SacrificedPet = warlock.ActivePet

			// Sacrificed pets lose all buffs, but don't remove trigger auras
			for _, aura := range warlock.ActivePet.GetAuras() {
				if aura.Duration == core.NeverExpires {
					continue
				}

				aura.Deactivate(sim)
			}
		}
	}

	warlock.ActivePet = newPet

	if newPet != nil {
		newPet.Enable(sim, newPet)
	}
}

func (warlock *Warlock) registerPets() {
	warlock.Felhunter = warlock.makeFelhunter()
	warlock.Imp = warlock.makeImp()
	warlock.Succubus = warlock.makeSuccubus()
	warlock.Voidwalker = warlock.makeVoidwalker()

	warlock.BasePets = []*WarlockPet{warlock.Felhunter, warlock.Imp, warlock.Succubus, warlock.Voidwalker}

	if warlock.HasRune(proto.WarlockRune_RuneBracerSummonFelguard) {
		warlock.Felguard = warlock.makeFelguard()
		warlock.BasePets = append(warlock.BasePets, warlock.Felguard)
	}
}

func (warlock *Warlock) makePet(cfg PetConfig, enabledOnStart bool) *WarlockPet {
	wp := &WarlockPet{
		Pet:          core.NewPet(cfg.Name, &warlock.Character, cfg.Stats, warlock.makeStatInheritance(), enabledOnStart, false),
		owner:        warlock,
		OnPetDisable: func(sim *core.Simulation, isSacrifice bool) {},
	}

	wp.EnableManaBarWithModifier(cfg.PowerModifier)

	if cfg.Name == "Imp" {
		// Imp gets 1mp/5 non casting regen per spirit
		wp.PseudoStats.SpiritRegenMultiplier = 1
		wp.PseudoStats.SpiritRegenRateCasting = 0
		wp.SpiritManaRegenPerSecond = func() float64 {
			// 1mp5 per spirit
			return wp.GetStat(stats.Spirit) / 5
		}

		// Mage spell crit scaling for imp
		wp.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[proto.Class_ClassMage][int(wp.Level)]*core.SpellCritRatingPerCritChance)
	} else {
		// Warrior scaling for all other pets
		wp.AddStat(stats.AttackPower, -20)
		wp.AddStatDependency(stats.Strength, stats.AttackPower, 2)

		// Warrior crit scaling
		wp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[proto.Class_ClassWarrior][int(wp.Level)]*core.CritRatingPerCritChance)
		wp.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[proto.Class_ClassWarrior][int(wp.Level)]*core.SpellCritRatingPerCritChance)

		// Imps generally don't melee
		wp.EnableAutoAttacks(wp, cfg.AutoAttacks)
		wp.AutoAttacks.MHConfig().DamageMultiplier *= 1 + 0.04*float64(warlock.Talents.UnholyPower)
	}

	wp.ApplyOnPetEnable(func(sim *core.Simulation) {
		// Warlock pets only inherit the owner's cast speed
		wp.EnableDynamicCastSpeedInheritance(sim)
	})

	// Having either the Fire or Shadow ring rune grants the pet 6% spell hit, 4.8% physical hit, and 0.5% expertise
	hasRingRune := warlock.HasRuneById(int32(proto.RingRune_RuneRingFireSpecialization)) || warlock.HasRuneById(int32(proto.RingRune_RuneRingShadowSpecialization))
	if hasRingRune {
		wp.AddStat(stats.MeleeHit, 4.8*core.MeleeHitRatingPerHitChance)
		wp.AddStat(stats.SpellHit, 6*core.SpellHitRatingPerHitChance)
		wp.AddStat(stats.Expertise, 0.5*core.ExpertiseRatingPerExpertiseChance)
	}

	core.ApplyPetConsumeEffects(&wp.Character, warlock.Consumes)

	warlock.AddPet(wp)

	return wp
}

func (warlock *Warlock) registerPetAbilities() {
	if warlock.Felguard != nil {
		warlock.Felguard.registerFelguardCleaveSpell()
		warlock.Felguard.registerFelguardDemonicFrenzyAura()
	}
	warlock.Imp.registerImpFireboltSpell()
	warlock.Succubus.registerSuccubusLashOfPainSpell()
}

func (wp *WarlockPet) GetPet() *core.Pet {
	return &wp.Pet
}

func (wp *WarlockPet) Initialize() {
}

func (wp *WarlockPet) Reset(_ *core.Simulation) {
	wp.manaPooling = false
}

func (wp *WarlockPet) Disable(sim *core.Simulation, isSacrifice bool) {
	wp.Pet.Disable(sim)

	if wp.OnPetDisable != nil {
		wp.OnPetDisable(sim, isSacrifice)
	}
}

func (wp *WarlockPet) ApplyOnPetDisable(newOnPetDisable OnPetDisable) {
	oldOnPetDisable := wp.OnPetDisable
	if oldOnPetDisable == nil {
		wp.OnPetDisable = oldOnPetDisable
	} else {
		wp.OnPetDisable = func(sim *core.Simulation, isSacrifice bool) {
			oldOnPetDisable(sim, isSacrifice)
			newOnPetDisable(sim, isSacrifice)
		}
	}
}

func (wp *WarlockPet) ExecuteCustomRotation(sim *core.Simulation) {
	if !wp.IsEnabled() || wp.primaryAbility == nil {
		return
	}

	if wp.manaPooling {
		maxPossibleCasts := sim.GetRemainingDuration().Seconds() / wp.primaryAbility.CurCast.CastTime.Seconds()

		if wp.CurrentMana() > (maxPossibleCasts*wp.primaryAbility.CurCast.Cost)*0.75 {
			wp.manaPooling = false
			wp.WaitUntil(sim, sim.CurrentTime+10*time.Millisecond)
			return
		}

		if wp.CurrentMana() >= wp.MaxMana()*0.94 {
			wp.manaPooling = false
			wp.WaitUntil(sim, sim.CurrentTime+10*time.Millisecond)
			return
		}

		if wp.manaPooling {
			return
		}
	}

	if !wp.primaryAbility.IsReady(sim) {
		wp.WaitUntil(sim, wp.primaryAbility.CD.ReadyAt())
		return
	}

	if wp.Unit.CurrentMana() >= wp.primaryAbility.CurCast.Cost {
		wp.primaryAbility.Cast(sim, wp.CurrentTarget)
	} else if !wp.owner.Options.PetPoolMana {
		wp.manaPooling = true
	}
}

func (warlock *Warlock) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// based on testing for WotLK Classic the following is true:
		// - pets are meele hit capped if and only if the warlock has 210 (8%) spell hit rating or more
		//   - this is unaffected by suppression and by magic hit debuffs like FF
		// - pets gain expertise from 0% to 6.5% relative to the owners hit, reaching cap at 17% spell hit
		//   - this is also unaffected by suppression and by magic hit debuffs like FF
		//   - this is continious, i.e. not restricted to 0.25 intervals
		// - pets gain spell hit from 0% to 17% relative to the owners hit, reaching cap at 12% spell hit
		// spell hit rating is floor'd
		//   - affected by suppression and ff, but in weird ways:
		// 3/3 suppression => 262 hit  (9.99%) results in misses, 263 (10.03%) no misses
		// 2/3 suppression => 278 hit (10.60%) results in misses, 279 (10.64%) no misses
		// 1/3 suppression => 288 hit (10.98%) results in misses, 289 (11.02%) no misses
		// 0/3 suppression => 314 hit (11.97%) results in misses, 315 (12.01%) no misses
		// 3/3 suppression + FF => 209 hit (7.97%) results in misses, 210 (8.01%) no misses
		// 2/3 suppression + FF => 222 hit (8.46%) results in misses, 223 (8.50%) no misses
		//
		// the best approximation of this behaviour is that we scale the warlock's spell hit by `1/12*17` floor
		// the result and then add the hit percent from suppression/ff

		// does correctly not include ff/misery
		ownerHitChance := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance
		highestSchoolPower := ownerStats[stats.SpellPower] + ownerStats[stats.SpellDamage] + max(ownerStats[stats.FirePower], ownerStats[stats.ShadowPower])

		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      highestSchoolPower * 0.565,
			stats.MP5:              ownerStats[stats.Intellect] * 0.315,
			stats.SpellDamage:      highestSchoolPower * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:         math.Floor(ownerHitChance / 12.0 * 17.0),
			stats.MeleeCrit:        ownerStats[stats.MeleeCrit],
			stats.SpellCrit:        ownerStats[stats.SpellCrit],
			stats.Dodge:            ownerStats[stats.MeleeCrit],
		}
	}
}
