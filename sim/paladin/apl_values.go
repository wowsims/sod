package paladin

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

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
	seal       proto.PaladinSeal
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

// func (paladin *Paladin) executePrimarySeal(sim *core.Simulation) {
// 	if !paladin.GCD.IsReady(sim) {
// 		return
// 	}
// 	// Dummy cast until inputs are linked up.
// 	paladin.Exorcism[1].Cast(sim, paladin.CurrentTarget)
// 	paladin.SealOfMartyrdom.Cast
// }

func (action *APLActionCastPaladinPrimarySeal) Execute(sim *core.Simulation) {

	paladin := action.paladin
	// paladin.SealOfMartyrdom

	paladin.Exorcism[1].Cast(sim, paladin.CurrentTarget)

	// rotAction := &core.PendingAction{
	// 	Priority:     core.ActionPriorityGCD,
	// 	OnAction:     paladin.executePrimarySeal,
	// 	NextActionAt: sim.CurrentTime,
	// }
	// sim.AddPendingAction(rotAction)
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
