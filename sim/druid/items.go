package druid

import (
	"time"

	"github.com/wowsims/sod/sim/common/sod"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Totem Item IDs
const (
	WolfsheadHelm             = 8345
	IdolMindExpandingMushroom = 209576
	Catnip                    = 213407
	IdolOfWrath               = 216490
	BloodBarkCrusher          = 216499
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(IdolMindExpandingMushroom, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.Spirit, 5)
	})

	core.NewItemEffect(BloodBarkCrusher, func(agent core.Agent) {
		character := agent.GetCharacter()
		auraActionID := core.ActionID{SpellID: 436482}
		numHits := min(3, character.Env.GetNumTargets())

		triggeredDmgSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 436481},
			SpellSchool: core.SpellSchoolStormstrike,
			DefenseType: core.DefenseTypeMelee, // actually has DefenseTypeNone, but is likely using the greatest CritMultiplier available
			ProcMask:    core.ProcMaskEmpty,

			// TODO: "Causes additional threat" in Tooltip, no clue what the multiplier is.
			ThreatMultiplier: 1,
			DamageMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.CalcAndDealDamage(sim, curTarget, 5, spell.OutcomeMagicHitAndCrit)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		mainAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Bloodbark Cleave",
			ActionID: auraActionID,
			Duration: 20 * time.Second,

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask&core.ProcMaskMelee != 0 {
					triggeredDmgSpell.Cast(sim, result.Target)
					return
				}
			},
		})

		mainSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: auraActionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				mainAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    mainSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
		})
	})

	core.AddEffectsToTest = true
}

// https://www.wowhead.com/classic/item=213407/catnip
func (druid *Druid) registerCatnipCD() {
	if druid.Consumes.DefaultConjured != proto.Conjured_ConjuredDruidCatnip {
		return
	}
	sod.RegisterFiftyPercentHasteBuffCD(&druid.Character, core.ActionID{ItemID: Catnip})
}
