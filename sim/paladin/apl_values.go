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

func (paladin *Paladin) newActionPaladinPrimarySealAction(_ *core.APLRotation, _ *proto.APLActionCastPaladinPrimarySeal) core.APLActionImpl {
	return &APLActionCastPaladinPrimarySeal{
		paladin: paladin,
	}
}

func (x *APLActionCastPaladinPrimarySeal) Execute(sim *core.Simulation) {
	x.lastAction = sim.CurrentTime
	x.paladin.primarySeal.Cast(sim, x.paladin.CurrentTarget)
}

func (action *APLActionCastPaladinPrimarySeal) ExecuteOffGCD(sim *core.Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (x *APLActionCastPaladinPrimarySeal) IsReady(sim *core.Simulation) bool {
	return sim.CurrentTime > x.lastAction && x.paladin.primarySeal.CanCast(sim, x.paladin.CurrentTarget)
}

func (action *APLActionCastPaladinPrimarySeal) IsOffGCDAction() bool {
	return false
}

func (x *APLActionCastPaladinPrimarySeal) Reset(*core.Simulation) {
	x.lastAction = core.DurationFromSeconds(-100)
}

func (x *APLActionCastPaladinPrimarySeal) String() string {
	return "Cast Primary Seal()"
}
