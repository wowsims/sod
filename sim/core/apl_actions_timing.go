package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
)

type APLActionWait struct {
	defaultAPLActionImpl
	unit     *Unit
	duration APLValue

	curWaitTime time.Duration
}

func (rot *APLRotation) newActionWait(config *proto.APLActionWait) APLActionImpl {
	unit := rot.unit
	durationVal := rot.coerceTo(rot.newAPLValue(config.Duration), proto.APLValueType_ValueTypeDuration)
	if durationVal == nil {
		return nil
	}

	return &APLActionWait{
		unit:     unit,
		duration: durationVal,
	}
}
func (action *APLActionWait) GetAPLValues() []APLValue {
	return []APLValue{action.duration}
}
func (action *APLActionWait) IsReady(sim *Simulation) bool {
	return action.duration.GetDuration(sim) > 0
}

func (action *APLActionWait) IsOffGCDAction() bool {
	return false
}

func (action *APLActionWait) Execute(sim *Simulation) {
	action.unit.Rotation.pushControllingAction(action)
	action.curWaitTime = sim.CurrentTime + action.duration.GetDuration(sim)

	pa := &PendingAction{
		Priority:     ActionPriorityLow,
		OnAction:     action.unit.gcdAction.OnAction,
		NextActionAt: action.curWaitTime,
	}
	sim.AddPendingAction(pa)
}

func (action *APLActionWait) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionWait) GetNextAction(sim *Simulation) *APLAction {
	if sim.CurrentTime >= action.curWaitTime {
		action.unit.Rotation.popControllingAction(action)
		return action.unit.Rotation.getNextAction(sim)
	} else {
		return nil
	}
}

func (action *APLActionWait) String() string {
	return fmt.Sprintf("Wait(%s)", action.duration)
}

type APLActionWaitUntil struct {
	defaultAPLActionImpl
	unit      *Unit
	condition APLValue
}

func (rot *APLRotation) newActionWaitUntil(config *proto.APLActionWaitUntil) APLActionImpl {
	unit := rot.unit
	conditionVal := rot.coerceTo(rot.newAPLValue(config.Condition), proto.APLValueType_ValueTypeBool)
	if conditionVal == nil {
		return nil
	}

	return &APLActionWaitUntil{
		unit:      unit,
		condition: conditionVal,
	}
}
func (action *APLActionWaitUntil) GetAPLValues() []APLValue {
	return []APLValue{action.condition}
}
func (action *APLActionWaitUntil) IsReady(sim *Simulation) bool {
	return !action.condition.GetBool(sim)
}

func (action *APLActionWaitUntil) IsOffGCDAction() bool {
	return false
}

func (action *APLActionWaitUntil) Execute(sim *Simulation) {
	action.unit.Rotation.pushControllingAction(action)
}

func (action *APLActionWaitUntil) GetNextAction(sim *Simulation) *APLAction {
	if action.condition.GetBool(sim) {
		action.unit.Rotation.popControllingAction(action)
		return action.unit.Rotation.getNextAction(sim)
	} else {
		return nil
	}
}

func (action *APLActionWaitUntil) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionWaitUntil) String() string {
	return fmt.Sprintf("WaitUntil(%s)", action.condition)
}

type APLActionRelativeSchedule struct {
	defaultAPLActionImpl
	innerAction *APLAction

	timing   time.Duration
	prevTime time.Duration
}

func (rot *APLRotation) newActionRelativeSchedule(config *proto.APLActionRelativeSchedule) APLActionImpl {
	innerAction := rot.newAPLAction(config.InnerAction)
	if innerAction == nil {
		return nil
	}

	timing := 1 * time.Second
	valid := true

	if durVal, err := time.ParseDuration(strings.TrimSpace(config.Schedule)); err == nil {
		timing = durVal
	} else {
		rot.ValidationWarning("Invalid duration value '%s'", strings.TrimSpace(config.Schedule))
		valid = false
	}

	if !valid {
		return nil
	}

	return &APLActionRelativeSchedule{
		innerAction: innerAction,
		timing:      timing,
		prevTime:    -1 * time.Minute,
	}
}
func (action *APLActionRelativeSchedule) Reset(*Simulation) {
	action.prevTime = -1 * time.Minute
}
func (action *APLActionRelativeSchedule) GetInnerActions() []*APLAction {
	return []*APLAction{action.innerAction}
}
func (action *APLActionRelativeSchedule) IsReady(sim *Simulation) bool {
	isReady := action.innerAction.IsReady(sim) && (action.prevTime != (sim.CurrentTime + action.timing))

	return isReady
}

func (action *APLActionRelativeSchedule) IsOffGCDAction() bool {
	return action.innerAction.impl.IsOffGCDAction()
}

func (action *APLActionRelativeSchedule) Execute(sim *Simulation) {
	action.prevTime = sim.CurrentTime + action.timing
	scheduledTime := action.prevTime

	if action.IsOffGCDAction() {

		//if sim.Log != nil {
		//	sim.Log("APLActionRelativeSchedule Execute Scheduling delayed off GCD action for %f", scheduledTime)
		//}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: scheduledTime,
			OnAction: func(s *Simulation) {
				action.innerAction.ExecuteOffGCD(sim, scheduledTime)
			},
		})
	} else {
		//if sim.Log != nil {
		//	sim.Log("APLActionRelativeSchedule Execute Scheduling delayed on GCD action ")
		//}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: scheduledTime,
			OnAction: func(s *Simulation) {
				if action.innerAction.IsReady(sim) { // Need to check as there is no guarantee it will be ready
					action.innerAction.Execute(sim)
				}
			},
		})
	}
}

func (action *APLActionRelativeSchedule) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionRelativeSchedule) String() string {
	return fmt.Sprintf("RelativeSchedule(%s, %s)", action.timing, action.innerAction)
}

type APLActionSchedule struct {
	defaultAPLActionImpl
	innerAction *APLAction

	timings       []time.Duration
	nextTimingIdx int
}

func (rot *APLRotation) newActionSchedule(config *proto.APLActionSchedule) APLActionImpl {
	innerAction := rot.newAPLAction(config.InnerAction)
	if innerAction == nil {
		return nil
	}

	timingStrs := strings.Split(config.Schedule, ",")
	if len(timingStrs) == 0 {
		return nil
	}

	timings := make([]time.Duration, len(timingStrs))
	valid := true
	for i, timingStr := range timingStrs {
		if durVal, err := time.ParseDuration(strings.TrimSpace(timingStr)); err == nil {
			timings[i] = durVal
		} else {
			rot.ValidationWarning("Invalid duration value '%s'", strings.TrimSpace(timingStr))
			valid = false
		}
	}
	if !valid {
		return nil
	}

	return &APLActionSchedule{
		innerAction: innerAction,
		timings:     timings,
	}
}
func (action *APLActionSchedule) Reset(*Simulation) {
	action.nextTimingIdx = 0
}
func (action *APLActionSchedule) GetInnerActions() []*APLAction {
	return []*APLAction{action.innerAction}
}
func (action *APLActionSchedule) IsReady(sim *Simulation) bool {

	//if sim.Log != nil {
	//	sim.Log("APLActionSchedule IsReady timing index is %d", action.nextTimingIdx)
	//}

	checkA := action.nextTimingIdx < len(action.timings)
	checkB := false

	if checkA {
		if action.IsOffGCDAction() {
			//if sim.Log != nil {
			//	sim.Log("Scheduled action is offGCD B check")
			//}
			checkB = sim.CurrentTime >= action.timings[action.nextTimingIdx]-(time.Millisecond*1500)
		} else {
			//if sim.Log != nil {
			//	sim.Log("Scheduled action is regular B check")
			//}
			checkB = sim.CurrentTime >= action.timings[action.nextTimingIdx]
		}
	}

	checkC := action.innerAction.IsReady(sim)
	isReady := checkA && checkB && checkC

	//if sim.Log != nil && isReady {
	//	sim.Log("Scheduled action is ready")
	//}
	//if sim.Log != nil && !isReady {
	//	sim.Log("Scheduled action is not ready")
	//	if checkA {
	//		sim.Log("Scheduled action is not ready currentTime %f nextactionTime %f ", sim.CurrentTime, action.timings[action.nextTimingIdx])
	//	}
	//	sim.Log("Scheduled action is not ready reason %t %t %t", checkA, checkB, checkC)
	//}
	return isReady
}

func (action *APLActionSchedule) IsOffGCDAction() bool {
	return action.innerAction.impl.IsOffGCDAction()
}

func (action *APLActionSchedule) Execute(sim *Simulation) {

	//if sim.Log != nil {
	//	sim.Log("APLActionSchedule Execute timing index is %d", action.nextTimingIdx)
	//}

	if action.IsOffGCDAction() {
		offGCDTime := action.timings[action.nextTimingIdx]

		//if sim.Log != nil {
		//	sim.Log("APLActionSchedule Execute Scheduling delayed off GCD action for %f", offGCDTime)
		//}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: offGCDTime,
			OnAction: func(s *Simulation) {
				action.innerAction.ExecuteOffGCD(sim, offGCDTime)
			},
		})
		//action.innerAction.ExecuteOffGCD(sim, offGCDTime)
		action.nextTimingIdx++
	} else {
		//if sim.Log != nil {
		//	sim.Log("APLActionSchedule Execute Scheduling non-delayed action ")
		//}
		action.nextTimingIdx++
		action.innerAction.Execute(sim)
	}
}

func (action *APLActionSchedule) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionSchedule) String() string {
	return fmt.Sprintf("Schedule(%s, %s)", action.timings, action.innerAction)
}

type APLActionPeriodicSchedule struct {
	defaultAPLActionImpl
	innerAction *APLAction

	timings       []time.Duration
	period        time.Duration
	nextTimingIdx int
}

func (rot *APLRotation) newActionPeriodicSchedule(config *proto.APLActionPeriodicSchedule) APLActionImpl {
	innerAction := rot.newAPLAction(config.InnerAction)
	if innerAction == nil {
		return nil
	}

	timingStrs := strings.Split(config.Schedule, ",")
	if len(timingStrs) != 2 {
		return nil
	}

	timings := make([]time.Duration, len(timingStrs))
	valid := true
	period := 1 * time.Second
	for i, timingStr := range timingStrs {
		if durVal, err := time.ParseDuration(strings.TrimSpace(timingStr)); err == nil {
			timings[i] = durVal
			if i == 1 {
				period = durVal

				if timings[1] == 0 {
					rot.ValidationWarning("Invalid periodic duration value '%s'", strings.TrimSpace(timingStr))
					valid = false
				}
				timings[1] = timings[0] + durVal
			}
		} else {
			rot.ValidationWarning("Invalid duration value '%s'", strings.TrimSpace(timingStr))
			valid = false
		}
	}
	if !valid {
		return nil
	}

	return &APLActionPeriodicSchedule{
		innerAction: innerAction,
		timings:     timings,
		period:      period,
	}
}
func (action *APLActionPeriodicSchedule) Reset(*Simulation) {
	action.nextTimingIdx = 0
}
func (action *APLActionPeriodicSchedule) GetInnerActions() []*APLAction {
	return []*APLAction{action.innerAction}
}
func (action *APLActionPeriodicSchedule) IsReady(sim *Simulation) bool {

	//if sim.Log != nil {
	//	sim.Log("APLActionPeriodicSchedule IsReady timing index is %d", action.nextTimingIdx)
	//}

	checkA := action.nextTimingIdx < len(action.timings)
	checkB := false

	if checkA {
		if action.IsOffGCDAction() {
			//if sim.Log != nil {
			//	sim.Log("Scheduled action is offGCD B check")
			//}
			checkB = sim.CurrentTime >= action.timings[action.nextTimingIdx]-(time.Millisecond*1500)
		} else {
			//if sim.Log != nil {
			//	sim.Log("Scheduled action is regular B check")
			//}
			checkB = sim.CurrentTime >= action.timings[action.nextTimingIdx]
		}
	}

	checkC := action.innerAction.IsReady(sim)
	isReady := checkA && checkB && checkC

	//if sim.Log != nil && isReady {
	//	sim.Log("Scheduled action is ready")
	//}
	//if sim.Log != nil && !isReady {
	//	sim.Log("Scheduled action is not ready")
	//	if checkA {
	//		sim.Log("Scheduled action is not ready currentTime %f nextactionTime %f ", sim.CurrentTime, action.timings[action.nextTimingIdx])
	//	}
	//	sim.Log("Scheduled action is not ready reason %t %t %t", checkA, checkB, checkC)
	//}
	return isReady
}

func (action *APLActionPeriodicSchedule) IsOffGCDAction() bool {
	return action.innerAction.impl.IsOffGCDAction()
}

func (action *APLActionPeriodicSchedule) Execute(sim *Simulation) {

	//if sim.Log != nil {
	//	sim.Log("APLActionPeriodicSchedule Execute timing index is %d", action.nextTimingIdx)
	//}

	if action.IsOffGCDAction() {
		offGCDTime := action.timings[action.nextTimingIdx]

		//if sim.Log != nil {
		//	sim.Log("APLActionPeriodicSchedule Execute Scheduling delayed off GCD action for %f", offGCDTime)
		//}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: offGCDTime,
			OnAction: func(s *Simulation) {
				action.innerAction.ExecuteOffGCD(sim, offGCDTime)
			},
		})

		action.nextTimingIdx++
		if action.nextTimingIdx >= len(action.timings) {
			action.timings = append(action.timings, action.timings[action.nextTimingIdx-1]+action.period)
		}
	} else {
		//if sim.Log != nil {
		//	sim.Log("APLActionPeriodicSchedule Execute Scheduling non-delayed action ")
		//}
		action.nextTimingIdx++
		if action.nextTimingIdx >= len(action.timings) {
			action.timings = append(action.timings, action.timings[action.nextTimingIdx-1]+action.period)
		}
		action.innerAction.Execute(sim)
	}
}

func (action *APLActionPeriodicSchedule) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionPeriodicSchedule) String() string {
	return fmt.Sprintf("PeriodicSchedule(%s, %s)", action.timings, action.innerAction)
}
