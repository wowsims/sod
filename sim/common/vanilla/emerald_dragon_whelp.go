package vanilla

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

type EmeraldDragonWhelp struct {
	core.Pet

	acidSpit *core.Spell

	disabledAt time.Duration
}

func NewEmeraldDragonWhelp(character *core.Character) *EmeraldDragonWhelp {
	whelpBaseStats := stats.Stats{
		stats.Health:      1500, // https://wowwiki-archive.fandom.com/wiki/Dragon%27s_Call
		stats.Intellect:   20,   // Adding the base 20 intellect to not mess with the base mana function
		stats.Mana:        500,  // TODO: Assumed value. The whelp seems to cast 3 Acid Spits (90 mana) per spawn (Rain: In the log you can see a whelp casting 4 acid spits so i'm increasing this to 500)
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

	whelp.EnableAutoAttacks(whelp, core.AutoAttackOptions{
		// TODO: Need whelp data
		MainHand: core.Weapon{
			// These stats are a complete guess from looking at the lone log I could find with Dragon's Call below
			// https://vanilla.warcraftlogs.com/reports/tQW9mqDrx3R4AdYZ#type=damage-done&ability=-13049&boss=-2&difficulty=0&wipes=2&source=25
			BaseDamageMin: 80.0,
			BaseDamageMax: 100.0,
			SwingSpeed:    2.0,
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

func (whelp *EmeraldDragonWhelp) Initialize() {
	whelp.registerAcidSpitSpell()
}

func (whelp *EmeraldDragonWhelp) ExecuteCustomRotation(sim *core.Simulation) {
	// Run the cast check only on swings or cast completes
	if whelp.AutoAttacks.NextAttackAt() != sim.CurrentTime+whelp.AutoAttacks.MainhandSwingSpeed() && whelp.AutoAttacks.NextAnyAttackAt()-1 > sim.CurrentTime {
		whelp.WaitUntil(sim, whelp.AutoAttacks.NextAttackAt()-1)
		return
	}

	if sim.Proc(0.5, "Acid Spit Cast") {
		// If the whelp will timeout during this cast just dont do it and stop attacks as well
		// If we dont do this the timeline cast time visual for the spell never ends because we
		// dont support hardcast interrupts
		if sim.CurrentTime+whelp.acidSpit.CastTime() >= whelp.disabledAt {
			whelp.AutoAttacks.StopMeleeUntil(sim, whelp.disabledAt, false)
		} else {
			whelp.acidSpit.Cast(sim, whelp.CurrentTarget)
		}
	} else {
		whelp.WaitUntil(sim, whelp.AutoAttacks.NextAttackAt()-1)
	}
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

	whelp.acidSpit = whelp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		// All of the casts and hits in the above log had the same damage so it would seem debuffs are ignored
		Flags: core.SpellFlagIgnoreModifiers | core.SpellFlagResetAttackSwing,

		ManaCost: core.ManaCostOptions{
			FlatCost: 90,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: The one log i was looking at has 0 misses on the spell but it also has only 25 casts
			// so i can't make a good assumption. Right now we leave it with a hit check and we can remove later.
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
				if whelp, ok := petAgent.(*EmeraldDragonWhelp); ok {
					whelp.EnableWithTimeout(sim, whelp, time.Second*15)
					whelp.disabledAt = sim.CurrentTime + time.Second*15
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
