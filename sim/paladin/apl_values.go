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

func (cast *APLActionCastPaladinPrimarySeal) GetInnerActions() []*core.APLAction { return nil }
func (cast *APLActionCastPaladinPrimarySeal) GetAPLValues() []core.APLValue      { return nil }
func (cast *APLActionCastPaladinPrimarySeal) Finalize(*core.APLRotation)         {}
func (cast *APLActionCastPaladinPrimarySeal) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}

func (paladin *Paladin) newActionPaladinPrimarySealAction(_ *core.APLRotation, _ *proto.APLActionCastPaladinPrimarySeal) core.APLActionImpl {
	return &APLActionCastPaladinPrimarySeal{
		paladin: paladin,
	}
}

func (cast *APLActionCastPaladinPrimarySeal) Execute(sim *core.Simulation) {
	cast.lastAction = sim.CurrentTime
	paladin := cast.paladin
	// If the player options are incorrectly configured, then no primary seal will be selected.
	if paladin.PrimarySealSpell == nil {
		return
	}
	paladin.PrimarySealSpell.Cast(sim, paladin.CurrentTarget)
}

func (cast *APLActionCastPaladinPrimarySeal) IsReady(sim *core.Simulation) bool {
	paladin := cast.paladin
	return sim.CurrentTime > cast.lastAction && paladin.PrimarySealSpell.CanCast(sim, paladin.CurrentTarget)
}

func (cast *APLActionCastPaladinPrimarySeal) Reset(*core.Simulation) {
	cast.lastAction = core.DurationFromSeconds(-100)
}

func (cast *APLActionCastPaladinPrimarySeal) String() string {
	return "Cast Primary Seal()"
}
