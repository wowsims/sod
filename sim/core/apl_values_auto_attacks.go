package core

import (
	"time"

	"github.com/wowsims/sod/sim/core/proto"
)

type APLValueAutoTimeSinceLast struct {
	DefaultAPLValueImpl
	unit     *Unit
	autoType proto.APLValueAutoTimeSinceLast_AttackType
}

func (rot *APLRotation) newValueAutoTimeSinceLast(config *proto.APLValueAutoTimeSinceLast) APLValue {
	return &APLValueAutoTimeSinceLast{
		unit:     rot.unit,
		autoType: config.AutoType,
	}
}
func (value *APLValueAutoTimeSinceLast) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoTimeSinceLast) GetDuration(sim *Simulation) time.Duration {
	duration := time.Duration(0)
	switch value.autoType {
	case proto.APLValueAutoTimeSinceLast_Melee:
		return max(duration, sim.CurrentTime-value.unit.AutoAttacks.LastAutoAt())
	case proto.APLValueAutoTimeSinceLast_MainHand:
		return max(duration, sim.CurrentTime-value.unit.AutoAttacks.LastMainhandAutoAt())
	case proto.APLValueAutoTimeSinceLast_OffHand:
		return max(duration, sim.CurrentTime-value.unit.AutoAttacks.LastOffhandAutoAt())
	case proto.APLValueAutoTimeSinceLast_Ranged:
		return max(duration, sim.CurrentTime-value.unit.AutoAttacks.LastRangedAutoAt())
	default:
		// defaults to Any
		return max(duration, sim.CurrentTime-value.unit.AutoAttacks.LastAnyAutoAt())
	}
}
func (value *APLValueAutoTimeSinceLast) String() string {
	return "Auto Time Since Last"
}

type APLValueAutoTimeToNext struct {
	DefaultAPLValueImpl
	unit     *Unit
	autoType proto.APLValueAutoTimeToNext_AttackType
}

func (rot *APLRotation) newValueAutoTimeToNext(config *proto.APLValueAutoTimeToNext) APLValue {
	return &APLValueAutoTimeToNext{
		unit:     rot.unit,
		autoType: config.AutoType,
	}
}
func (value *APLValueAutoTimeToNext) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoTimeToNext) GetDuration(sim *Simulation) time.Duration {
	switch value.autoType {
	case proto.APLValueAutoTimeToNext_Melee:
		return max(0, value.unit.AutoAttacks.NextAttackAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_MainHand:
		return max(0, value.unit.AutoAttacks.MainhandSwingAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_OffHand:
		return max(0, value.unit.AutoAttacks.OffhandSwingAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_Ranged:
		return max(0, value.unit.AutoAttacks.NextRangedAttackAt()-sim.CurrentTime)
	}
	// defaults to Any
	return max(0, value.unit.AutoAttacks.NextAnyAttackAt()-sim.CurrentTime)
}
func (value *APLValueAutoTimeToNext) String() string {
	return "Auto Time To Next"
}

type APLValueAutoSwingTime struct {
	DefaultAPLValueImpl
	unit     *Unit
	autoType proto.APLValueAutoSwingTime_SwingType
}

func (rot *APLRotation) newValueAutoSwingTime(config *proto.APLValueAutoSwingTime) APLValue {
	return &APLValueAutoSwingTime{
		unit:     rot.unit,
		autoType: config.AutoType,
	}
}
func (value *APLValueAutoSwingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoSwingTime) GetDuration(sim *Simulation) time.Duration {
	switch value.autoType {
	case proto.APLValueAutoSwingTime_MainHand:
		return max(0, value.unit.AutoAttacks.MainhandSwingSpeed())
	case proto.APLValueAutoSwingTime_OffHand:
		return max(0, value.unit.AutoAttacks.OffhandSwingSpeed())
	case proto.APLValueAutoSwingTime_Ranged:
		return max(0, value.unit.AutoAttacks.RangedSwingSpeed())
	}
	// defaults to 0
	return 0
}
func (value *APLValueAutoSwingTime) String() string {
	return "Auto Swing Time"
}
