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

	mage.registerFrozenOrbTickSpell()

	hasFoFRune := mage.HasRune(proto.MageRune_RuneChestFingersOfFrost)

	cooldown := time.Minute
	manaCost := 0.11

	// Create a dummy aura for tracking Frozen Orb uptime in the APL
	activeAura := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: int32(proto.MageRune_RuneCloakFrozenOrb)},
		Label:    "Frozen Orb",
		Duration: time.Second * 15,
	})

	mage.FrozenOrb = mage.RegisterSpell(core.SpellConfig{
		ClassSpellMask: ClassSpellMask_MageFrozenOrb,
		ActionID:       core.ActionID{SpellID: int32(proto.MageRune_RuneCloakFrozenOrb)},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,

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

					if hasFoFRune {
						mage.FingersOfFrostAura.Activate(sim)
						mage.FingersOfFrostAura.AddStacks(sim, mage.FingersOfFrostAura.MaxStacks)
					}
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

func (mage *Mage) registerFrozenOrbTickSpell() {
	baseDamage := mage.baseRuneAbilityDamage() * 0.9
	spellCoef := .129

	mage.FrozenOrbTick = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 440809},
		ClassSpellMask: ClassSpellMask_MageFrozenOrbTick,
		SpellSchool:    core.SpellSchoolFrost,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellProc | core.ProcMaskSpellDamageProc,
		Flags:          SpellFlagChillSpell | core.SpellFlagNotAProc | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagBinary,

		BonusCoefficient: spellCoef,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

type FrozenOrb struct {
	core.Pet

	mage *Mage
}

func (mage *Mage) NewFrozenOrbPets() []*FrozenOrb {
	// It's possible to have up to 2 Frozen Orbs active at a time because of Cold Snap
	return []*FrozenOrb{mage.newFrozenOrb(1), mage.newFrozenOrb(2)}
}

func (mage *Mage) newFrozenOrb(idx int32) *FrozenOrb {
	frozenOrb := &FrozenOrb{
		Pet:  core.NewPet(fmt.Sprintf("Frozen Orb %d", idx), &mage.Character, frozenOrbBaseStats, frozenOrbStatInheritance(), false, true),
		mage: mage,
	}

	mage.AddPet(frozenOrb)

	return frozenOrb
}

func (orb *FrozenOrb) GetPet() *core.Pet {
	return &orb.Pet
}

func (orb *FrozenOrb) Initialize() {
}

func (orb *FrozenOrb) Reset(_ *core.Simulation) {
}

func (orb *FrozenOrb) ExecuteCustomRotation(sim *core.Simulation) {
	if success := orb.mage.FrozenOrbTick.Cast(sim, orb.mage.CurrentTarget); success {
		orb.WaitUntil(sim, sim.CurrentTime+time.Second)
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
