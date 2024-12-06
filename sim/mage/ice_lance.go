package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// TODO: Classic review ice lance numbers on live
func (mage *Mage) registerIceLanceSpell() {
	if !mage.HasRune(proto.MageRune_RuneHandsIceLance) {
		return
	}

	hasWintersChillTalent := mage.Talents.WintersChill > 0

	baseDamageLow := mage.baseRuneAbilityDamage() * 0.55
	baseDamageHigh := mage.baseRuneAbilityDamage() * 0.65
	spellCoeff := 0.429
	manaCost := 0.08

	mage.IceLance = mage.RegisterSpell(core.SpellConfig{
		SpellCode:    SpellCode_MageIceLance,
		ActionID:     core.ActionID{SpellID: int32(proto.MageRune_RuneHandsIceLance)},
		SpellSchool:  core.SpellSchoolFrost,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | core.SpellFlagAPL,
		MissileSpeed: 38,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	if !hasWintersChillTalent {
		return
	}

	mage.GlaciateAuras = mage.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
		return unit.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneHandsIceLance), Tag: 1}, // Temp ID
			Label:     "Glaciate",
			Duration:  time.Second * 15,
			MaxStacks: 5,
		})
	})

	oldApplyEffects := mage.IceLance.ApplyEffects
	mage.IceLance.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		var glaciateAura *core.Aura
		modifier := 0.0

		if glaciateAura = mage.GlaciateAuras.Get(target); glaciateAura.IsActive() {
			modifier += 0.20 * float64(glaciateAura.GetStacks())
		}

		spell.DamageMultiplierAdditive += modifier
		oldApplyEffects(sim, target, spell)
		spell.DamageMultiplierAdditive -= modifier

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if glaciateAura != nil {
				glaciateAura.Deactivate(sim)
			}
		})
	}

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Glaciate Trigger",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Only Blizzard ticks proc
			if spell.SpellCode == SpellCode_MageBlizzard && spell.Flags.Matches(SpellFlagChillSpell) {
				glaciateAura := mage.GlaciateAuras.Get(result.Target)
				glaciateAura.Activate(sim)
				glaciateAura.AddStack(sim)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.SpellSchool.Matches(core.SpellSchoolFrost) && spell.Flags.Matches(SpellFlagMage) && spell.SpellCode != SpellCode_MageIceLance {
				glaciateAura := mage.GlaciateAuras.Get(result.Target)
				glaciateAura.Activate(sim)
				glaciateAura.AddStack(sim)
			}
		},
	}))
}
