package paladin

import (
	"fmt"
	"strconv"

	"github.com/wowsims/sod/sim/core"
)

const sorRanks = 8
const sor1hModifier = 0.85
const sor2hModifier = 1.2

var sorLevels = [sorRanks + 1]int{0, 1, 10, 18, 26, 34, 42, 50, 58}
var sorManaCosts = [sorRanks + 1]float64{0, 20, 40, 60, 90, 120, 140, 170, 200}
var sorAuraSpellIds = [sorRanks + 1]int32{0, 20154, 20287, 20288, 20289, 20290, 20291, 20292, 20293}
var sorProcSpellIds = [sorRanks + 1]int32{0, 25742, 25740, 25739, 25738, 25737, 25736, 25735, 25713}
var sorEffectBasePoints = [sorRanks + 1]float64{0, 107, 215, 351, 540, 784, 1081, 1406, 1785}
var sorEffectRealPointsPerLevel = [sorRanks + 1]float64{0, 18, 17, 23, 31, 37, 41, 47, 47}
var sorLevelMinMaxEffects = [sorRanks + 1][]int32{{0}, {1, 7}, {10, 16}, {18, 24}, {26, 32}, {34, 40}, {42, 48}, {50, 56}, {58, 60}}

// SoR Rank 3 has approximately double the seemingly-intended spellpower scaling
var sorEffectBonusCoefficient = [sorRanks + 1]float64{0, 0.029, 0.063, 0.184, 0.1, 0.1, 0.1, 0.1, 0.1}

var jorSpellIDs = [sorRanks + 1]int32{0, 20187, 20280, 20281, 20282, 20283, 20284, 20285, 20286}
var jorEffectBasePoints = [sorRanks + 1]float64{0, 14, 24, 38, 56, 77, 101, 130, 161}
var jorEffectRealPointsPerLevel = [sorRanks + 1]float64{0, 1.8, 1.9, 2.4, 2.8, 3.1, 3.8, 4.1, 4.1}
var jorEffectDieSides = [sorRanks + 1]float64{0, 1, 3, 5, 7, 9, 11, 13, 17}
var jorEffectBonusCoefficient = [sorRanks + 1]float64{0, 0.144, 0.312, 0.462, 0.5, 0.5, 0.5, 0.5, 0.5}

func (paladin *Paladin) applySealOfRighteousnessSpellAndAuraBaseConfig(rank int) {
	spellIdAura := sorAuraSpellIds[rank]
	spellIdProc := sorProcSpellIds[rank]
	basePoints := sorEffectBasePoints[rank]
	pointsPerLevel := sorEffectRealPointsPerLevel[rank]
	scalingLevelMin := sorLevelMinMaxEffects[rank][0]
	scalingLevelMax := sorLevelMinMaxEffects[rank][1]
	effectBonusCoefficient := sorEffectBonusCoefficient[rank]

	jorSpellID := jorSpellIDs[rank]
	jorBasePoints := jorEffectBasePoints[rank]
	jorPointsPerLevel := jorEffectRealPointsPerLevel[rank]
	jorDieSides := jorEffectDieSides[rank]
	jorBonusCoefficient := jorEffectBonusCoefficient[rank]

	manaCost := sorManaCosts[rank]
	level := sorLevels[rank]

	levelsToScale := min(paladin.Level, scalingLevelMax) - scalingLevelMin
	baseCoefficientFinal := basePoints + float64(levelsToScale)*pointsPerLevel

	handednessModifier := sor1hModifier
	if paladin.Has2hEquipped() {
		handednessModifier = sor2hModifier
	}
	weaponSpeed := paladin.GetMHWeapon().SwingSpeed
	impSoRModifier := core.TernaryFloat64(
		paladin.Talents.ImprovedSealOfRighteousness >= 1,
		1+0.03*float64(paladin.Talents.ImprovedSealOfRighteousness),
		1.0)
	baseDamageNoSP := baseCoefficientFinal / 100 * handednessModifier * weaponSpeed * impSoRModifier

	baseJoRMinDamage := jorBasePoints + float64(levelsToScale)*jorPointsPerLevel
	baseJoRMaxDamage := baseJoRMinDamage + jorDieSides
	fmt.Println(baseJoRMinDamage, baseJoRMaxDamage)
	/*
	 * Seal of Righteousness is an Spell/Aura that when active makes the paladin capable of procing
	 * 2 different SpellIDs depending on a paladin's casted spell or melee swing.
	 * NOTE:
	 *   Seal of Righteousness is unique in that it is the only seal that can proc off its own judgements.
	 *
	 * (Judgement of Righteousness):
	 *   - Procs off of the dummy Judgement spell.
	 *   - Cannot miss or be dodged/parried/blocked.
	 *   - Deals a flat damage that is affected by Improved SoR talent, and
	 *     has a spellpower scaling that is unaffacted by that talent.
	 *   - Crits off of a spell modifier.
	 *
	 * (Seal of Righteousness):
	 *   - Procs off of white hits.
	 *   - Cannot miss or be dodged/parried if the underlying white hit lands.
	 *   - Deals damage that is a function of weapon speed, and spellpower.
	 *   - Has 0.85 scale factor on base damage if using 1h, 1.2 if using 2h.
	 *   - CANNOT CRIT.
	 */

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: jorSpellID}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.SpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseJoRMinDamage, baseJoRMaxDamage) + jorBonusCoefficient*spell.SpellPower()
			// Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
		},
	})

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{

		ActionID:      core.ActionID{SpellID: spellIdProc}, // Seal of Righteousness damage bonus.
		SpellSchool:   core.SpellSchoolHoly,
		ProcMask:      core.ProcMaskEmpty,
		Flags:         core.SpellFlagMeleeMetrics,
		RequiredLevel: level,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Testing seems to show 2h benefits from spellpower about 12% more than 1h weapons.
			handednessModifierSP := core.TernaryFloat64(paladin.Has2hEquipped(), 1.12, 1.0)
			baseDamage := baseDamageNoSP + effectBonusCoefficient*spell.SpellPower()*handednessModifierSP

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	// Seal of Righteousness aura.
	auraActionID := core.ActionID{SpellID: spellIdAura}
	paladin.SealOfRighteousnessAura[rank] = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Righteousness" + paladin.Label + strconv.Itoa(rank),
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses or our own procs.
			if !result.Landed() || spell.SpellID == onSwingProc.SpellID {
				return
			}
			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				onSwingProc.Cast(sim, result.Target)
			}
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, result.Target)
			}

		},
	})

	aura := paladin.SealOfRighteousnessAura[rank]
	paladin.SealOfRighteousness[rank] = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Righteousness self buff.
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1.0 - (float64(paladin.Talents.Benediction) * 0.03),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}

func (paladin *Paladin) registerSealOfRighteousnessSpellAndAura() {

	paladin.SealOfRighteousness = make([]*core.Spell, sorRanks+1)
	paladin.SealOfRighteousnessAura = make([]*core.Aura, sorRanks+1)

	for rank := 1; rank <= sorRanks; rank++ {
		if int(paladin.Level) >= sorLevels[rank] {
			paladin.applySealOfRighteousnessSpellAndAuraBaseConfig(rank)
		}
	}

}
