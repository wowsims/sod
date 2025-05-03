package core

import (
	"container/heap"
	"math"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
)

type MoveModifier struct {
	ActionId *ActionID
	Modifier float64
}

type MoveHeap []MoveModifier

func (h MoveHeap) Len() int           { return len(h) }
func (h MoveHeap) Less(i, j int) bool { return h[i].Modifier > h[j].Modifier }
func (h MoveHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MoveHeap) Push(x any) {
	*h = append(*h, x.(MoveModifier))
}

func (h *MoveHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *MoveHeap) activeModifier() float64 {
	heap := *h
	if len(heap) < 1 {
		panic("Movement Heaps should never be missing their base case!!!")
	}
	n := heap[0]
	return n.Modifier
}

func (h *MoveHeap) Find(actionId *ActionID) int {
	for i, mod := range *h {
		if mod.ActionId == actionId {
			return i
		}
	}
	return -1
}

func (move *MovementHandler) addMoveSpeedModifier(moveHeap *MoveHeap, moveMod MoveModifier) {
	heap.Push(moveHeap, moveMod)
	move.updateMoveSpeed()
}

func (move *MovementHandler) removeMoveSpeedModifier(moveHeap *MoveHeap, actionID *ActionID) {
	index := moveHeap.Find(actionID)
	if index == -1 {
		return
	}
	heap.Remove(moveHeap, index)
	move.updateMoveSpeed()
}

func (move *MovementHandler) updateMoveSpeed() {
	move.MoveSpeed = move.baseSpeed * move.getActveModifier(move.moveSpeedBonuses) * (1 - move.getActveModifier(move.moveSpeedPenalties))
}

func (move *MovementHandler) getActveModifier(moveHeap *MoveHeap) float64 {
	return moveHeap.activeModifier()
}

type MovementHandler struct {
	Moving    bool
	MoveSpeed float64

	baseSpeed          float64
	moveAura           *Aura
	moveSpell          *Spell
	moveSpeedBonuses   *MoveHeap
	moveSpeedPenalties *MoveHeap
}

func (unit *Unit) initMovement() {
	unit.MovementHandler = &MovementHandler{
		moveSpeedBonuses: &MoveHeap{
			MoveModifier{
				Modifier: 1,
			},
		},
		moveSpeedPenalties: &MoveHeap{
			MoveModifier{
				Modifier: 0,
			},
		},
		baseSpeed: 7.0,
	}
	unit.MovementHandler.updateMoveSpeed()

	unit.MovementHandler.moveAura = unit.GetOrRegisterAura(Aura{
		Label:     "Movement",
		ActionID:  ActionID{OtherID: proto.OtherAction_OtherActionMove},
		Duration:  NeverExpires,
		MaxStacks: 30,

		OnGain: func(aura *Aura, sim *Simulation) {
			if unit.IsChanneling(sim) {
				unit.ChanneledDot.Cancel(sim)
			}
			unit.AutoAttacks.CancelAutoSwing(sim)
			unit.MovementHandler.Moving = true
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			unit.MovementHandler.Moving = false
			unit.AutoAttacks.EnableAutoSwing(sim)

			// Simulate the delay from starting attack
			unit.AutoAttacks.DelayMeleeBy(sim, time.Millisecond*50)
		},
	})

	unit.MovementHandler.moveSpell = unit.GetOrRegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionMove},
		Flags:    SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			unit.MovementHandler.moveAura.Activate(sim)
			unit.MovementHandler.moveAura.SetStacks(sim, int32(unit.DistanceFromTarget))
		},
	})
}

func (unit *Unit) IsMoving() bool {
	return unit.MovementHandler.Moving
}

func (unit *Unit) MoveTo(moveRange float64, sim *Simulation) {
	if moveRange == unit.DistanceFromTarget || unit.IsMoving() {
		return
	}

	moveDistance := moveRange - unit.DistanceFromTarget
	moveTicks := math.Abs(moveDistance)
	moveInterval := moveDistance / float64(moveTicks)

	unit.MovementHandler.moveSpell.Cast(sim, unit.CurrentTarget)

	sim.AddPendingAction(NewPeriodicAction(sim, PeriodicActionOptions{
		Period:          time.Millisecond * time.Duration(1000/(unit.MovementHandler.MoveSpeed)),
		NumTicks:        int(moveTicks),
		TickImmediately: false,

		OnAction: func(sim *Simulation) {
			unit.DistanceFromTarget += moveInterval
			unit.MovementHandler.moveAura.SetStacks(sim, int32(unit.DistanceFromTarget))

			if unit.DistanceFromTarget == moveRange {
				unit.MovementHandler.moveAura.Deactivate(sim)
			}
		},
	}))
}

// A move speed increase of 30% should be represented as 1.30 and a move speed slow of 70% should be respresented as 0.70
func (unit *Unit) AddMoveSpeedModifier(actionId *ActionID, modifier float64) {
	moveSpeedMod := MoveModifier{
		ActionId: actionId,
		Modifier: modifier,
	}
	if moveSpeedMod.Modifier < 1 {
		unit.MovementHandler.addMoveSpeedModifier(unit.MovementHandler.moveSpeedPenalties, moveSpeedMod)
	} else {
		unit.MovementHandler.addMoveSpeedModifier(unit.MovementHandler.moveSpeedBonuses, moveSpeedMod)
	}

}

func (unit *Unit) RemoveMoveSpeedModifier(actionID *ActionID) {
	unit.MovementHandler.removeMoveSpeedModifier(unit.MovementHandler.moveSpeedPenalties, actionID)
	unit.MovementHandler.removeMoveSpeedModifier(unit.MovementHandler.moveSpeedBonuses, actionID)
}
