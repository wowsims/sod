package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// The APLValue for the remaining duration of the primary seal aura.

func (paladin *Paladin) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_PrimarySealRemainingTime:
		return paladin.newValuePrimarySealRemainingTime(rot, config.GetPrimarySealRemainingTime())
	default:
		return nil
	}
}

type APLValuePrimarySealRemainingTime struct {
	core.DefaultAPLValueImpl
	paladin *Paladin
}

func (paladin *Paladin) newValuePrimarySealRemainingTime(rot *core.APLRotation, config *proto.APLValuePrimarySealRemainingTime) core.APLValue {
	return &APLValuePrimarySealRemainingTime{
		paladin: paladin,
	}
}

func (value *APLValuePrimarySealRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}

func (value *APLValuePrimarySealRemainingTime) GetDuration(sim *core.Simulation) time.Duration {
	paladin := value.paladin
	// if paladin.CurrentSeal == nil {
	// 	return 0
	// }
	return max(paladin.CurrentSealExpiration - sim.CurrentTime)
}

func (value *APLValuePrimarySealRemainingTime) String() string {
	return "Primary Seal Remaining Time()"
}

// The APLAction for casting the primary Seal
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
