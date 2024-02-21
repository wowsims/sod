package paladin

import (
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

func (paladin *Paladin) newValueCurrentSealRemainingTime(rot *core.APLRotation, config *proto.APLValueCurrentSealRemainingTime) core.APLValue {
	return &APLValueCurrentSealRemainingTime{
		paladin: paladin,
	}
}

func (value *APLValueCurrentSealRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}

func (value *APLValueCurrentSealRemainingTime) GetDuration(sim *core.Simulation) time.Duration {
	paladin := value.paladin
	return max(paladin.CurrentSealExpiration - sim.CurrentTime)
}

func (value *APLValueCurrentSealRemainingTime) String() string {
	return "Current Seal Remaining Time()"
}

// The APLAction for casting the current Seal
func (paladin *Paladin) NewAPLAction(rot *core.APLRotation, config *proto.APLAction) core.APLActionImpl {
	switch config.Action.(type) {
	case *proto.APLAction_CastPaladinPrimarySeal:
		return paladin.newActionPaladinPrimarySealAction(rot, config.GetCastPaladinPrimarySeal())
	default:
		return nil
	}
}

type APLActionCastPaladinPrimarySeal struct {
	paladin    *Paladin
	lastAction time.Duration
}

func (impl *APLActionCastPaladinPrimarySeal) GetInnerActions() []*core.APLAction { return nil }
func (impl *APLActionCastPaladinPrimarySeal) GetAPLValues() []core.APLValue      { return nil }
func (impl *APLActionCastPaladinPrimarySeal) Finalize(*core.APLRotation)         {}
func (impl *APLActionCastPaladinPrimarySeal) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}

func (paladin *Paladin) newActionPaladinPrimarySealAction(_ *core.APLRotation, config *proto.APLActionCastPaladinPrimarySeal) core.APLActionImpl {
	return &APLActionCastPaladinPrimarySeal{
		paladin: paladin,
	}
}

func (action *APLActionCastPaladinPrimarySeal) Execute(sim *core.Simulation) {
	paladin := action.paladin
	paladin.PrimarySealSpell.Cast(sim, paladin.CurrentTarget)
	// paladin.Exorcism[1].Cast(sim, paladin.CurrentTarget)
	action.lastAction = sim.CurrentTime
}

func (action *APLActionCastPaladinPrimarySeal) IsReady(sim *core.Simulation) bool {
	return sim.CurrentTime > action.lastAction
}

func (action *APLActionCastPaladinPrimarySeal) Reset(*core.Simulation) {
	action.lastAction = core.DurationFromSeconds(-100)
}

func (action *APLActionCastPaladinPrimarySeal) String() string {
	return "Cast Primary Seal()"
}
