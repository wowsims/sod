package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) registerVampiricTouchSpell() {
	if !priest.HasRune(proto.PriestRune_RuneCloakVampiricTouch) {
		return
	}

	hasDespairRune := priest.HasRune(proto.PriestRune_RuneBracersDespair)

	numTicks := int32(5)
	tickLength := time.Second * 3
	baseTickDamage := priest.baseRuneAbilityDamage() * .65
	spellCoef := 0.167
	castTime := time.Millisecond * 1500
	manaCost := 0.16

	partyPlayers := priest.Env.Raid.GetPlayerParty(&priest.Unit).Players
	// https: //www.wowhead.com/classic/spell=402779/vampiric-touch
	manaMetrics := priest.NewManaMetrics(core.ActionID{SpellID: 402779})
	manaReturnedMultiplier := 0.05

	manaGainAura := priest.RegisterAura(core.Aura{
		Label:    "Vampiric Touch (Mana)",
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
				manaGained := result.Damage * manaReturnedMultiplier
				for _, player := range partyPlayers {
					player.GetCharacter().AddMana(sim, manaGained, manaMetrics)
				}
			}
		},
	})

	priest.VampiricTouch = priest.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_PriestVampiricTouch,
		ActionID:    core.ActionID{SpellID: int32(proto.PriestRune_RuneCloakVampiricTouch)},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagPriest | core.SpellFlagAPL | core.SpellFlagPureDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
		},

		CritDamageBonus: priest.periodicCritBonus(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "VampiricTouch",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					manaGainAura.Activate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					manaGainAura.Deactivate(sim)
				},
			},

			NumberOfTicks:    numTicks,
			TickLength:       tickLength,
			BonusCoefficient: spellCoef,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseTickDamage, isRollover)
				if isRollover && priest.VampiricEmbraceAuras != nil {
					priest.VampiricEmbraceAuras.Get(target).Activate(sim)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if hasDespairRune {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
