package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
)

type APLActionChangeTarget struct {
	defaultAPLActionImpl
	unit      *Unit
	newTarget UnitReference
}

func (rot *APLRotation) newActionChangeTarget(config *proto.APLActionChangeTarget) APLActionImpl {
	if config.NewTarget == nil {
		return nil
	}
	newTarget := rot.GetSourceUnit(config.NewTarget)
	if newTarget.Get() == nil {
		return nil
	}
	return &APLActionChangeTarget{
		unit:      rot.unit,
		newTarget: newTarget,
	}
}
func (action *APLActionChangeTarget) IsReady(sim *Simulation) bool {
	return action.unit.CurrentTarget != action.newTarget.Get()
}

func (action *APLActionChangeTarget) IsOffGCDAction() bool {
	return false
}
func (action *APLActionChangeTarget) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.unit.Log(sim, "Changing target to %s", action.newTarget.Get().Label)
	}
	action.unit.CurrentTarget = action.newTarget.Get()
}
func (action *APLActionChangeTarget) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}
func (action *APLActionChangeTarget) String() string {
	return fmt.Sprintf("Change Target(%s)", action.newTarget.Get().Label)
}

type APLActionCancelAura struct {
	defaultAPLActionImpl
	aura *Aura
}

func (rot *APLRotation) newActionCancelAura(config *proto.APLActionCancelAura) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionCancelAura{
		aura: aura.Get(),
	}
}

func (action *APLActionCancelAura) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}

func (action *APLActionCancelAura) IsOffGCDAction() bool {
	return true
}
func (action *APLActionCancelAura) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Cancelling aura %s", action.aura.ActionID)
	}
	action.aura.Deactivate(sim)
}
func (action *APLActionCancelAura) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}
func (action *APLActionCancelAura) String() string {
	return fmt.Sprintf("Cancel Aura(%s)", action.aura.ActionID)
}

type APLActionActivateAura struct {
	defaultAPLActionImpl
	aura *Aura
}

func (rot *APLRotation) newActionActivateAura(config *proto.APLActionActivateAura) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionActivateAura{
		aura: aura.Get(),
	}
}

func (action *APLActionActivateAura) IsReady(sim *Simulation) bool {
	return true
}

func (action *APLActionActivateAura) IsOffGCDAction() bool {
	return false
}
func (action *APLActionActivateAura) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Activating aura %s", action.aura.ActionID)
	}
	action.aura.Activate(sim)
}

func (action *APLActionActivateAura) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionActivateAura) String() string {
	return fmt.Sprintf("Activate Aura(%s)", action.aura.ActionID)
}

type APLActionActivateAuraWithStacks struct {
	defaultAPLActionImpl
	aura      *Aura
	numStacks int32
}

func (rot *APLRotation) newActionActivateAuraWithStacks(config *proto.APLActionActivateAuraWithStacks) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	if aura.Get().MaxStacks == 0 {
		rot.ValidationWarning("%s is not a stackable aura", ProtoToActionID(config.AuraId))
		return nil
	}

	numStacks, err := strconv.Atoi(config.NumStacks)
	if err != nil {
		numStacks = 0
	}

	return &APLActionActivateAuraWithStacks{
		aura:      aura.Get(),
		numStacks: int32(numStacks),
	}
}

func (action *APLActionActivateAuraWithStacks) IsReady(sim *Simulation) bool {
	return true
}

func (action *APLActionActivateAuraWithStacks) IsOffGCDAction() bool {
	return false
}

func (action *APLActionActivateAuraWithStacks) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Activating aura %s (%d stacks)", action.aura.ActionID, action.numStacks)
	}
	action.aura.Activate(sim)
	action.aura.SetStacks(sim, action.numStacks)
}

func (action *APLActionActivateAuraWithStacks) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionActivateAuraWithStacks) String() string {
	return fmt.Sprintf("Activate Aura(%s) Stacks(%d)", action.aura.ActionID, action.numStacks)
}

type APLActionAddComboPoints struct {
	defaultAPLActionImpl
	character *Character
	numPoints int32
	metrics   *ResourceMetrics
}

func (rot *APLRotation) newActionAddComboPoints(config *proto.APLActionAddComboPoints) APLActionImpl {
	character := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).GetCharacter()
	numPoints, err := strconv.Atoi(config.NumPoints)
	metrics := character.NewComboPointMetrics(ActionID{OtherID: proto.OtherAction_OtherActionComboPoints})

	if err != nil {
		numPoints = 0
	}
	return &APLActionAddComboPoints{
		character: character,
		numPoints: int32(numPoints),
		metrics:   metrics,
	}
}

func (action *APLActionAddComboPoints) IsReady(sim *Simulation) bool {
	return true
}

func (action *APLActionAddComboPoints) IsOffGCDAction() bool {
	return false
}

func (action *APLActionAddComboPoints) Execute(sim *Simulation) {
	numPoints := strconv.Itoa(int(action.numPoints))

	if sim.Log != nil {
		action.character.Log(sim, "Adding combo points (%s points)", numPoints)
	}

	action.character.AddComboPoints(sim, action.numPoints, action.metrics)
}

func (action *APLActionAddComboPoints) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}

func (action *APLActionAddComboPoints) String() string {
	numPoints := strconv.Itoa(int(action.numPoints))
	return fmt.Sprintf("Add Combo Points(%s)", numPoints)
}

type APLActionTriggerICD struct {
	defaultAPLActionImpl
	aura *Aura
}

func (rot *APLRotation) newActionTriggerICD(config *proto.APLActionTriggerICD) APLActionImpl {
	aura := rot.GetAPLICDAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionTriggerICD{
		aura: aura.Get(),
	}
}
func (action *APLActionTriggerICD) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}

func (action *APLActionTriggerICD) IsOffGCDAction() bool {
	return false
}
func (action *APLActionTriggerICD) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Triggering ICD %s", action.aura.ActionID)
	}
	action.aura.Icd.Use(sim)
}
func (action *APLActionTriggerICD) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}
func (action *APLActionTriggerICD) String() string {
	return fmt.Sprintf("Trigger ICD(%s)", action.aura.ActionID)
}

type APLActionItemSwap struct {
	defaultAPLActionImpl
	character *Character
	swapSet   proto.APLActionItemSwap_SwapSet
}

func (rot *APLRotation) newActionItemSwap(config *proto.APLActionItemSwap) APLActionImpl {
	if config.SwapSet == proto.APLActionItemSwap_Unknown {
		rot.ValidationWarning("Unknown item swap set")
		return nil
	}

	character := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).GetCharacter()
	if !character.ItemSwap.IsEnabled() {
		if config.SwapSet != proto.APLActionItemSwap_Main {
			rot.ValidationWarning("No swap set configured in Settings.")
		}
		return nil
	}

	return &APLActionItemSwap{
		character: character,
		swapSet:   config.SwapSet,
	}
}
func (action *APLActionItemSwap) IsReady(sim *Simulation) bool {
	return (action.swapSet == proto.APLActionItemSwap_Main) == action.character.ItemSwap.IsSwapped()
}

func (action *APLActionItemSwap) IsOffGCDAction() bool {
	return false
}
func (action *APLActionItemSwap) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.character.Log(sim, "Item Swap to set %s", action.swapSet)
	}

	if action.swapSet == proto.APLActionItemSwap_Main {
		action.character.ItemSwap.reset(sim)
	} else {
		action.character.ItemSwap.SwapItems(sim, action.character.ItemSwap.slots)
	}
}
func (action *APLActionItemSwap) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}
func (action *APLActionItemSwap) String() string {
	return fmt.Sprintf("Item Swap(%s)", action.swapSet)
}

type APLActionMove struct {
	defaultAPLActionImpl
	unit      *Unit
	moveRange APLValue
}

func (rot *APLRotation) newActionMove(config *proto.APLActionMove) APLActionImpl {
	return &APLActionMove{
		unit:      rot.unit,
		moveRange: rot.newAPLValue(config.RangeFromTarget),
	}
}
func (action *APLActionMove) IsReady(sim *Simulation) bool {
	isPrepull := sim.CurrentTime < 0
	return !action.unit.IsMoving() && (action.moveRange.GetFloat(sim) != action.unit.DistanceFromTarget || isPrepull) && !action.unit.IsCasting(sim)
}
func (action *APLActionMove) IsOffGCDAction() bool {
	return false
}
func (action *APLActionMove) Execute(sim *Simulation) {
	moveRange := action.moveRange.GetFloat(sim)
	if sim.Log != nil {
		action.unit.Log(sim, "Moving to %s", moveRange)
	}

	action.unit.MoveTo(moveRange, sim)
}
func (action *APLActionMove) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}
func (action *APLActionMove) String() string {
	return fmt.Sprintf("Move(%s)", action.moveRange)
}

type APLActionCustomRotation struct {
	defaultAPLActionImpl
	unit  *Unit
	agent Agent

	lastExecutedAt time.Duration
}

func (rot *APLRotation) newActionCustomRotation(config *proto.APLActionCustomRotation) APLActionImpl {
	agent := rot.unit.Env.GetAgentFromUnit(rot.unit)
	if agent == nil {
		panic("Agent not found for custom rotation")
	}

	return &APLActionCustomRotation{
		unit:  rot.unit,
		agent: agent,
	}
}
func (action *APLActionCustomRotation) Reset(sim *Simulation) {
	action.lastExecutedAt = -1
}
func (action *APLActionCustomRotation) IsReady(sim *Simulation) bool {
	// Prevent infinite loops by only allowing this action to be performed once at each timestamp.
	return action.lastExecutedAt != sim.CurrentTime
}
func (action *APLActionCustomRotation) IsOffGCDAction() bool {
	return false
}
func (action *APLActionCustomRotation) Execute(sim *Simulation) {
	action.lastExecutedAt = sim.CurrentTime
	action.agent.ExecuteCustomRotation(sim)
}
func (action *APLActionCustomRotation) ExecuteOffGCD(sim *Simulation, time time.Duration) {
	action.Execute(sim) // Default to Execute unless impletented for this APL Action
}
func (action *APLActionCustomRotation) String() string {
	return "Custom Rotation()"
}
