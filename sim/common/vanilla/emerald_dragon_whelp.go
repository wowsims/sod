package vanilla

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

type EmeraldDragonWhelp struct {
	core.Pet

	AcidSpit *core.Spell
}

func NewEmeraldDragonWhelp(character *core.Character) *EmeraldDragonWhelp {
	whelpBaseStats := stats.Stats{
		stats.Health:      1500, // https://wowwiki-archive.fandom.com/wiki/Dragon%27s_Call
		stats.Mana:        300,  // TODO: Assumed value. The whelp seems to cast 3 Acid Spits (90 mana) per spawn
		stats.SpellDamage: 220,  // Puts the Acid Spit damage very close to the below log
		// Based on this log but more data needed
		// https://sod.warcraftlogs.com/reports/xTwQVgbjF9cPnd3R#type=damage-done&ability=-13049&view=events&boss=-2&difficulty=0&wipes=2
		stats.MeleeCrit: 4.5 * core.CritRatingPerCritChance,
		stats.SpellCrit: 13 * core.CritRatingPerCritChance,
	}

	whelp := &EmeraldDragonWhelp{
		Pet: core.NewPet("Emerald Dragon Whelp", character, whelpBaseStats, emeraldWhelpingStatInheritance(), false, true),
	}
	whelp.Level = 55

	whelp.EnableManaBar()
	whelp.registerAcidSpitSpell()

	whelp.EnableAutoAttacks(whelp, core.AutoAttackOptions{
		// TODO: Need whelp data
		MainHand: core.Weapon{
			// These stats are a complete guess from looking at the lone log I could find with Dragon's Call below
			// https://vanilla.warcraftlogs.com/reports/tQW9mqDrx3R4AdYZ#type=damage-done&ability=-13049&boss=-2&difficulty=0&wipes=2&source=25
			BaseDamageMin: 80.0,
			BaseDamageMax: 100.0,
			SwingSpeed:    1.5,
			SpellSchool:   core.SpellSchoolPhysical,
		},
		AutoSwingMelee: true,
	})

	return whelp
}

func emeraldWhelpingStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// TODO: Needs more verification
		return stats.Stats{}
	}
}

func (whelp *EmeraldDragonWhelp) Initialize() {}

func (whelp *EmeraldDragonWhelp) ExecuteCustomRotation(sim *core.Simulation) {
	if !whelp.AcidSpit.IsReady(sim) {
		whelp.WaitUntil(sim, whelp.AcidSpit.CD.ReadyAt())
		return
	}

	whelp.AcidSpit.Cast(sim, whelp.CurrentTarget)
}

func (whelp *EmeraldDragonWhelp) Reset(sim *core.Simulation) {
	whelp.Disable(sim)
}

func (whelp *EmeraldDragonWhelp) OnPetDisable(sim *core.Simulation) {
}

func (whelp *EmeraldDragonWhelp) GetPet() *core.Pet {
	return &whelp.Pet
}

func (whelp *EmeraldDragonWhelp) registerAcidSpitSpell() {
	actionID := core.ActionID{SpellID: 9591}

	whelp.AcidSpit = whelp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		// All of the casts and hits in the above log had the same damage so it would seem debuffs are ignored
		Flags: core.SpellFlagIgnoreModifiers,

		// This is causing errors because of nil spell.Unit.Env in mana.go:315 so setting a cooldown instead
		// ManaCost: core.ManaCostOptions{
		// 	FlatCost: 90,
		// },

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
			CD: core.Cooldown{
				Timer:    whelp.NewTimer(),
				Duration: time.Second * 2,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				castTime := whelp.ApplyCastSpeedForSpell(cast.CastTime, spell)

				if castTime > 0 {
					whelp.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(64, 86), spell.OutcomeMagicHitAndCrit)
		},
	})
}

func MakeEmeraldDragonWhelpTriggerAura(agent core.Agent) {
	character := agent.GetCharacter()

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		ActionID: core.ActionID{SpellID: 13049},
		Name:     "Emerald Dragon Whelp Proc",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		PPM:      1.0, // Reported by armaments discord
		ICD:      time.Minute * 1,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			for _, petAgent := range character.PetAgents {
				if whelp, ok := petAgent.(*EmeraldDragonWhelp); ok && !whelp.IsEnabled() {
					whelp.EnableWithTimeout(sim, whelp, time.Second*15)
					break
				}
			}
		},
	})
}

func ConstructEmeralDragonWhelpPets(character *core.Character) {
	// Original could have up to 3 whelps active at a time however the SoD version seems to only summon 1 whelp on a 1 minute cooldown
	character.AddPet(NewEmeraldDragonWhelp(character))
}
