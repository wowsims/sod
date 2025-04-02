package core

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
)

type APLValueCurrentHealth struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealth(config *proto.APLValueCurrentHealth) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.ValidationWarning("%s does not use Health", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentHealth{
		unit: unit,
	}
}
func (value *APLValueCurrentHealth) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentHealth) GetFloat(_ *Simulation) float64 {
	return value.unit.Get().CurrentHealth()
}
func (value *APLValueCurrentHealth) String() string {
	return "Current Health"
}

type APLValueCurrentHealthPercent struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealthPercent(config *proto.APLValueCurrentHealthPercent) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.ValidationWarning("%s does not use Health", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentHealthPercent{
		unit: unit,
	}
}
func (value *APLValueCurrentHealthPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentHealthPercent) GetFloat(_ *Simulation) float64 {
	return value.unit.Get().CurrentHealthPercent()
}
func (value *APLValueCurrentHealthPercent) String() string {
	return fmt.Sprintf("Current Health %%")
}

type APLValueCurrentMana struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentMana(config *proto.APLValueCurrentMana) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.ValidationWarning("%s does not use Mana", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentMana{
		unit: unit,
	}
}
func (value *APLValueCurrentMana) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentMana) GetFloat(_ *Simulation) float64 {
	return value.unit.Get().CurrentMana()
}
func (value *APLValueCurrentMana) String() string {
	return "Current Mana"
}

type APLValueCurrentManaPercent struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentManaPercent(config *proto.APLValueCurrentManaPercent) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.ValidationWarning("%s does not use Mana", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentManaPercent{
		unit: unit,
	}
}
func (value *APLValueCurrentManaPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentManaPercent) GetFloat(_ *Simulation) float64 {
	return value.unit.Get().CurrentManaPercent()
}
func (value *APLValueCurrentManaPercent) String() string {
	return fmt.Sprintf("Current Mana %%")
}

type APLValueCurrentRage struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRage(_ *proto.APLValueCurrentRage) APLValue {
	unit := rot.unit
	if !unit.HasRageBar() {
		rot.ValidationWarning("%s does not use Rage", unit.Label)
		return nil
	}
	return &APLValueCurrentRage{
		unit: unit,
	}
}
func (value *APLValueCurrentRage) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentRage) GetFloat(_ *Simulation) float64 {
	return value.unit.CurrentRage()
}
func (value *APLValueCurrentRage) String() string {
	return "Current Rage"
}

type APLValueCurrentEnergy struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentEnergy(_ *proto.APLValueCurrentEnergy) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationWarning("%s does not use Energy", unit.Label)
		return nil
	}
	return &APLValueCurrentEnergy{
		unit: unit,
	}
}
func (value *APLValueCurrentEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentEnergy) GetFloat(_ *Simulation) float64 {
	return value.unit.CurrentEnergy()
}
func (value *APLValueCurrentEnergy) String() string {
	return "Current Energy"
}

type APLValueMaxEnergy struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueMaxEnergy(_ *proto.APLValueMaxEnergy) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationWarning("%s does not use Energy", unit.Label)
		return nil
	}
	return &APLValueMaxEnergy{
		unit: unit,
	}
}
func (value *APLValueMaxEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueMaxEnergy) GetFloat(sim *Simulation) float64 {
	return value.unit.MaxEnergy()
}
func (value *APLValueMaxEnergy) String() string {
	return "Max Energy"
}

type APLValueCurrentComboPoints struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentComboPoints(_ *proto.APLValueCurrentComboPoints) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationWarning("%s does not use Combo Points", unit.Label)
		return nil
	}
	return &APLValueCurrentComboPoints{
		unit: unit,
	}
}
func (value *APLValueCurrentComboPoints) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentComboPoints) GetInt(_ *Simulation) int32 {
	return value.unit.ComboPoints()
}
func (value *APLValueCurrentComboPoints) String() string {
	return "Current Combo Points"
}

type APLValueTimeToEnergyTick struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueTimeToEnergyTick(_ *proto.APLValueTimeToEnergyTick) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationWarning("%s does not use Energy", unit.Label)
		return nil
	}
	return &APLValueTimeToEnergyTick{
		unit: unit,
	}
}
func (value *APLValueTimeToEnergyTick) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueTimeToEnergyTick) GetDuration(sim *Simulation) time.Duration {
	return value.unit.NextEnergyTickAt() - sim.CurrentTime
}
func (value *APLValueTimeToEnergyTick) String() string {
	return "Time to Next Energy Tick"
}

type APLValueEnergyThreshold struct {
	DefaultAPLValueImpl
	unit      *Unit
	threshold float64
}

func (rot *APLRotation) newValueEnergyThreshold(config *proto.APLValueEnergyThreshold) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationWarning("%s does not use Energy", unit.Label)
		return nil
	}
	return &APLValueEnergyThreshold{
		unit:      unit,
		threshold: float64(config.Threshold),
	}
}
func (value *APLValueEnergyThreshold) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueEnergyThreshold) GetBool(_ *Simulation) bool {
	if value.threshold > 0 {
		return value.unit.currentEnergy >= value.threshold
	}
	return value.unit.currentEnergy >= value.unit.maxEnergy+value.threshold
}
func (value *APLValueEnergyThreshold) String() string {
	return "Energy Threshold"
}
