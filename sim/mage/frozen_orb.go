package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (mage *Mage) registerFrozenOrbCD() {
	if !mage.HasRune(proto.MageRune_RuneCloakFrozenOrb) {
		return
	}

	cooldown := time.Minute
	manaCost := 0.11

	// Create a dummy aura for tracking Frozen Orb uptime in the APL
	activeAura := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: int32(proto.MageRune_RuneCloakFrozenOrb)},
		Label:    "Frozen Orb",
		Duration: time.Second * 15,
	})

	mage.FrozenOrb = mage.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_MageFrozenOrb,
		ActionID:    core.ActionID{SpellID: int32(proto.MageRune_RuneCloakFrozenOrb)},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, orb := range mage.frozenOrbPets {
				if !orb.IsActive() {
					orb.EnableWithTimeout(sim, orb, time.Second*15)
					activeAura.Activate(sim)
					break
				}
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.FrozenOrb,
		Type:  core.CooldownTypeDPS,
	})
}

type FrozenOrb struct {
	core.Pet

	mage *Mage

	FrozenOrbTick          *core.Spell
	FrozenOrbFingerOfFrost *core.Aura
	TickCount              int64
}

func (mage *Mage) NewFrozenOrbPets() []*FrozenOrb {
	// It's possible to have up to 2 Frozen Orbs active at a time because of Cold Snap
	return []*FrozenOrb{mage.newFrozenOrb(1), mage.newFrozenOrb(2)}
}

func (mage *Mage) newFrozenOrb(idx int32) *FrozenOrb {
	frozenOrb := &FrozenOrb{
		Pet:       core.NewPet(fmt.Sprintf("Frozen Orb %d", idx), &mage.Character, frozenOrbBaseStats, frozenOrbStatInheritance(), false, true),
		mage:      mage,
		TickCount: 0,
	}

	mage.AddPet(frozenOrb)

	return frozenOrb
}

func (orb *FrozenOrb) GetPet() *core.Pet {
	return &orb.Pet
}

func (orb *FrozenOrb) Initialize() {
	orb.registerFrozenOrbTickSpell()

	// Frozen Orb seems to benefit from Frost Specialization
	orb.PseudoStats.SchoolBonusHitChance = orb.mage.PseudoStats.SchoolBonusHitChance
}

func (orb *FrozenOrb) Reset(_ *core.Simulation) {
}

func (orb *FrozenOrb) ExecuteCustomRotation(sim *core.Simulation) {
	if success := orb.FrozenOrbTick.Cast(sim, orb.mage.CurrentTarget); !success {
		orb.Disable(sim)
	}
}

var frozenOrbBaseStats = stats.Stats{}

func frozenOrbStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHit:   ownerStats[stats.SpellHit],
			stats.SpellCrit:  ownerStats[stats.SpellCrit],
			stats.SpellPower: ownerStats[stats.SpellPower],
		}
	}
}

func (orb *FrozenOrb) registerFrozenOrbTickSpell() {
	hasFOFRune := orb.mage.HasRune(proto.MageRune_RuneChestFingersOfFrost)

	baseDamage := orb.mage.baseRuneAbilityDamage() * 0.9
	spellCoef := .129

	orb.FrozenOrbTick = orb.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 440809},
		SpellSchool: core.SpellSchoolFrost | core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | SpellFlagChillSpell,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		BonusCoefficient: spellCoef,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			orb.TickCount += 1
			if orb.TickCount == 15 {
				orb.TickCount = 0
			}
		},
	})

	if hasFOFRune {
		orb.FrozenOrbFingerOfFrost = core.MakePermanent(orb.RegisterAura(core.Aura{
			Label: "Frozen Orb FoF",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == orb.FrozenOrbTick && orb.TickCount == 0 {
					orb.mage.FingersOfFrostAura.Activate(sim)
					orb.mage.FingersOfFrostAura.AddStack(sim)
				}
			},
		}))
	}
}
