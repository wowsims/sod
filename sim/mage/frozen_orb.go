package mage

import (
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
			mage.frozenOrb.EnableWithTimeout(sim, mage.frozenOrb, time.Second*15)
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

func (mage *Mage) NewFrozenOrb() *FrozenOrb {
	frozenOrb := &FrozenOrb{
		Pet:       core.NewPet("Frozen Orb", &mage.Character, frozenOrbBaseStats, createFrozenOrbInheritance(), false, true),
		mage:      mage,
		TickCount: 0,
	}

	mage.AddPet(frozenOrb)

	return frozenOrb
}

func (ffo *FrozenOrb) GetPet() *core.Pet {
	return &ffo.Pet
}

func (ffo *FrozenOrb) Initialize() {
	ffo.registerFrozenOrbTickSpell()
}

func (ffo *FrozenOrb) Reset(_ *core.Simulation) {
}

func (ffo *FrozenOrb) ExecuteCustomRotation(sim *core.Simulation) {
	if success := ffo.FrozenOrbTick.Cast(sim, ffo.mage.CurrentTarget); !success {
		ffo.Disable(sim)
	}
}

var frozenOrbBaseStats = stats.Stats{}

var createFrozenOrbInheritance = func() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHit:   ownerStats[stats.SpellHit],
			stats.SpellCrit:  ownerStats[stats.SpellCrit],
			stats.SpellPower: ownerStats[stats.SpellPower],
		}
	}
}

func (ffo *FrozenOrb) registerFrozenOrbTickSpell() {
	hasFOFRune := ffo.mage.HasRune(proto.MageRune_RuneChestFingersOfFrost)
	baseDamage := ffo.mage.baseRuneAbilityDamage() * 0.9
	spellCoef := .129

	ffo.FrozenOrbTick = ffo.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 440809},
		SpellSchool: core.SpellSchoolFrost | core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,

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
			ffo.TickCount += 1
			if ffo.TickCount == 15 {
				ffo.TickCount = 0
			}
		},
	})

	if hasFOFRune {
		ffo.FrozenOrbFingerOfFrost = core.MakePermanent(ffo.RegisterAura(core.Aura{
			Label: "Frozen Orb FoF",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == ffo.FrozenOrbTick && ffo.TickCount == 0 {
					ffo.mage.FingersOfFrostAura.Activate(sim)
					ffo.mage.FingersOfFrostAura.AddStack(sim)
				}
			},
		}))
	}
}
