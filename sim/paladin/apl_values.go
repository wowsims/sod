package paladin

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// The APLValue for the remaining duration of the primary seal aura.

func (paladin *Paladin) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_CurrentSealRemainingTime:
		return paladin.newValueCurrentSealRemainingTime(rot, config.GetCurrentSealRemainingTime())
	default:
		return nil
	}
}

type APLValueCurrentSealRemainingTime struct {
	core.DefaultAPLValueImpl
	paladin *Paladin
}

func (paladin *Paladin) newValueCurrentSealRemainingTime(_ *core.APLRotation, _ *proto.APLValueCurrentSealRemainingTime) core.APLValue {
	return &APLValueCurrentSealRemainingTime{
		paladin: paladin,
	}
}

func (x *APLValueCurrentSealRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}

func (x *APLValueCurrentSealRemainingTime) GetDuration(sim *core.Simulation) time.Duration {
	if x.paladin.currentSeal.IsActive() {
		return x.paladin.currentSeal.RemainingDuration(sim)
	}
	return 0
}

func (x *APLValueCurrentSealRemainingTime) String() string {
	return "Current Seal Remaining Time()"
}

// The APLAction for casting the current Seal
func (paladin *Paladin) NewAPLAction(rot *core.APLRotation, config *proto.APLAction) core.APLActionImpl {
	switch config.Action.(type) {
	case *proto.APLAction_CastPaladinPrimarySeal:
		return paladin.newActionPaladinPrimarySealAction(rot, config.GetCastPaladinPrimarySeal())
	case *proto.APLAction_PaladinCastWithMacro:
		return paladin.newActionPaladinCastWithMacro(rot, config.GetPaladinCastWithMacro())
	default:
		return nil
	}
}

type APLActionCastPaladinPrimarySeal struct {
	paladin    *Paladin
	lastAction time.Duration
}

func (x *APLActionCastPaladinPrimarySeal) GetInnerActions() []*core.APLAction             { return nil }
func (x *APLActionCastPaladinPrimarySeal) GetAPLValues() []core.APLValue                  { return nil }
func (x *APLActionCastPaladinPrimarySeal) Finalize(*core.APLRotation)                     {}
func (x *APLActionCastPaladinPrimarySeal) GetNextAction(*core.Simulation) *core.APLAction { return nil }
func (x *APLActionCastPaladinPrimarySeal) GetSpellFromAction() *core.Spell {
	return x.paladin.primarySeal
}

func (paladin *Paladin) newActionPaladinPrimarySealAction(_ *core.APLRotation, _ *proto.APLActionCastPaladinPrimarySeal) core.APLActionImpl {
	return &APLActionCastPaladinPrimarySeal{
		paladin: paladin,
	}
}

func (x *APLActionCastPaladinPrimarySeal) Execute(sim *core.Simulation) {
	x.lastAction = sim.CurrentTime
	x.paladin.primarySeal.Cast(sim, x.paladin.CurrentTarget)
}

func (x *APLActionCastPaladinPrimarySeal) IsReady(sim *core.Simulation) bool {
	return sim.CurrentTime > x.lastAction && x.paladin.primarySeal.CanCast(sim, x.paladin.CurrentTarget)
}

func (x *APLActionCastPaladinPrimarySeal) Reset(*core.Simulation) {
	x.lastAction = core.DurationFromSeconds(-100)
}

func (x *APLActionCastPaladinPrimarySeal) String() string {
	return "Cast Primary Seal()"
}

type APLActionPaladinCastWithMacro struct {
	paladin *Paladin
	spell   *core.Spell
	target  core.UnitReference
	macro   proto.APLActionPaladinCastWithMacro_Macro
}

func (x *APLActionPaladinCastWithMacro) GetInnerActions() []*core.APLAction { return nil }
func (x *APLActionPaladinCastWithMacro) GetAPLValues() []core.APLValue      { return nil }
func (x *APLActionPaladinCastWithMacro) Finalize(*core.APLRotation)         {}
func (x *APLActionPaladinCastWithMacro) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}
func (x *APLActionPaladinCastWithMacro) Reset(*core.Simulation) {}
func (x *APLActionPaladinCastWithMacro) GetSpellFromAction() *core.Spell {
	return x.spell
}

func (paladin *Paladin) newActionPaladinCastWithMacro(rot *core.APLRotation, config *proto.APLActionPaladinCastWithMacro) core.APLActionImpl {
	if config.Macro == proto.APLActionPaladinCastWithMacro_Unknown {
		rot.ValidationWarning("Unknown macro")
		return nil
	}

	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	target := rot.GetTargetUnit(config.Target)
	if target.Get() == nil {
		return nil
	}
	return &APLActionPaladinCastWithMacro{
		paladin: paladin,
		spell:   spell,
		target:  target,
		macro:   config.Macro,
	}
}
func (action *APLActionPaladinCastWithMacro) IsReady(sim *core.Simulation) bool {
	return action.spell.CanCast(sim, action.target.Get()) && (!action.spell.Flags.Matches(core.SpellFlagMCD) || action.spell.Unit.GCD.IsReady(sim) || action.spell.DefaultCast.GCD == 0)
}
func (action *APLActionPaladinCastWithMacro) Execute(sim *core.Simulation) {
	if action.macro == proto.APLActionPaladinCastWithMacro_StartAttack {
		action.ExecuteWithStartattack(sim)
	} else if action.macro == proto.APLActionPaladinCastWithMacro_StopAttack {
		action.ExecuteWithStopattack(sim)
	}
}
func (action *APLActionPaladinCastWithMacro) ExecuteWithStartattack(sim *core.Simulation) {
	action.paladin.bypassMacroOptions = true

	actualSpell := action.spell
	if actualSpell == action.paladin.judgement {
		actualSpell = action.paladin.currentJudgement
	}
	oldApplyEffects := actualSpell.ApplyEffects
	oldFlags := actualSpell.Flags
	actualSpell.Flags ^= core.SpellFlagBatchStopAttackMacro
	actualSpell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		oldApplyEffects(sim, target, spell)
		if sim.CurrentTime > 0 {
			action.paladin.AutoAttacks.EnableAutoSwing(sim)
		}
	}
	action.spell.Cast(sim, action.target.Get())

	actualSpell.Flags = oldFlags
	actualSpell.ApplyEffects = oldApplyEffects
	action.paladin.bypassMacroOptions = false
}
func (action *APLActionPaladinCastWithMacro) ExecuteWithStopattack(sim *core.Simulation) {
	action.paladin.bypassMacroOptions = true

	actualSpell := action.spell
	if actualSpell == action.paladin.judgement {
		actualSpell = action.paladin.currentJudgement
	}
	oldApplyEffects := actualSpell.ApplyEffects
	oldFlags := actualSpell.Flags
	actualSpell.Flags |= core.SpellFlagBatchStopAttackMacro
	actualSpell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		oldApplyEffects(sim, target, spell)
		action.paladin.performStopAttack(sim, target, spell)
	}
	action.spell.Cast(sim, action.target.Get())

	actualSpell.Flags = oldFlags
	actualSpell.ApplyEffects = oldApplyEffects
	action.paladin.bypassMacroOptions = false
}
func (action *APLActionPaladinCastWithMacro) String() string {
	return fmt.Sprintf("Cast Spell(%s) With Macro(%s)", action.spell.ActionID, action.macro.String())
}
