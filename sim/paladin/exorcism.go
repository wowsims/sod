package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core/stats"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const exorcismRanks = 6

var exorcismLevels = [exorcismRanks + 1]int{0, 20, 28, 36, 44, 52, 60}
var exorcismSpellIDs = [exorcismRanks + 1]int32{0, 415068, 415069, 415070, 415071, 415072, 415073}
var exorcismManaCosts = [exorcismRanks + 1]float64{0, 85, 135, 180, 235, 285, 345}
var exorcismEffectBasePoints = [exorcismRanks + 1]float64{0, 84, 151, 216, 303, 392, 504}
var exorcismEffectDieSides = [exorcismRanks + 1]float64{0, 13, 21, 29, 39, 47, 59}
var exorcismEffectRealPointsPerLevel = [exorcismRanks + 1]float64{0, 1.2, 1.6, 2.0, 2.4, 2.8, 3.2}
var exorcismMinMaxLevels = [exorcismRanks + 1][]int32{{0}, {20, 25}, {28, 33}, {36, 40}, {44, 49}, {52, 72}, {60, 60}}

func (paladin *Paladin) getExorcismBaseConfig(rank int) core.SpellConfig {
	spellId := exorcismSpellIDs[rank]
	manaCost := exorcismManaCosts[rank]
	level := exorcismLevels[rank]
	basePoints := exorcismEffectBasePoints[rank]
	pointsPerLevel := exorcismEffectRealPointsPerLevel[rank]
	dieSides := exorcismEffectDieSides[rank]
	scalingLevelMin := exorcismMinMaxLevels[rank][0]
	scalingLevelMax := exorcismMinMaxLevels[rank][1]

	levelsToScale := min(paladin.Level, scalingLevelMax) - scalingLevelMin
	baseDamageMin := basePoints + float64(levelsToScale)*pointsPerLevel
	baseDamageMax := baseDamageMin + dieSides

	hasExorcist := paladin.HasRune(proto.PaladinRune_RuneLegsExorcist)
	hasWrath := paladin.HasRune(proto.PaladinRune_RuneHeadWrath)

	return core.SpellConfig{
		// The SoD exorcism replacements don't list a rank so the DB doesn't list one either.
		// Using a tag as a replacement so that we can group Exorcism in the UI.
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolHoly,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
			Multiplier: core.TernaryFloat64(
				paladin.HasRune(proto.PaladinRune_RuneFeetTheArtOfWar),
				0.2,
				1.0,
			),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: *paladin.ExorcismCooldown,
		},

		BonusCritRating: paladin.holyPowerCritChance() + paladin.fanaticismCritChance(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.429,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hasExorcist || target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax)

			var bonusCrit float64
			if hasExorcist && (target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead) {
				bonusCrit += 100 * core.CritRatingPerCritChance
			}
			if hasWrath {
				bonusCrit += paladin.GetStat(stats.MeleeCrit)
			}

			spell.BonusCritRating += bonusCrit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.BonusCritRating -= bonusCrit
		},
	}

}

// Exorcism in SoD is by default castable only on demon and undead targets.
// If the paladin has the Exorcist leg rune equipped, they can cast the spell on
// any target and it additionally always crits on demon and undead targets.
func (paladin *Paladin) registerExorcismSpell() {
	paladin.ExorcismCooldown = &core.Cooldown{
		Timer:    paladin.NewTimer(),
		Duration: time.Second * 15,
	}

	if paladin.HasRune(proto.PaladinRune_RuneWristPurifyingPower) {
		paladin.ExorcismCooldown.Duration /= 2
	}

	paladin.Exorcism = make([]*core.Spell, exorcismRanks+1)

	for rank := 1; rank <= exorcismRanks; rank++ {
		config := paladin.getExorcismBaseConfig(rank)
		if config.RequiredLevel <= int(paladin.Level) {
			paladin.Exorcism[rank] = paladin.GetOrRegisterSpell(config)
		}
	}
}
