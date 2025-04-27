package hunter

import (
	"fmt"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (hunter *Hunter) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_HunterPetIsActive:
		return hunter.newValueHunterPetIsActive(rot)
	case *proto.APLValue_HunterCurrentPetFocus:
		return hunter.newValueHunterCurrentPetFocus(rot, config.GetHunterCurrentPetFocus())
	case *proto.APLValue_HunterCurrentPetFocusPercent:
		return hunter.newValueHunterCurrentPetFocusPercent(rot, config.GetHunterCurrentPetFocusPercent())
	default:
		return nil
	}
}

type APLValueHunterPetIsActive struct {
	core.DefaultAPLValueImpl
	hunter *Hunter
}

func (hunter *Hunter) newValueHunterPetIsActive(_ *core.APLRotation) core.APLValue {
	return &APLValueHunterPetIsActive{
		hunter: hunter,
	}
}

func (value *APLValueHunterPetIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueHunterPetIsActive) GetBool(sim *core.Simulation) bool {
	return value.hunter.pet != nil
}
func (value *APLValueHunterPetIsActive) String() string {
	return fmt.Sprintf("Current Pet Focus %%")
}

type APLValueHunterCurrentPetFocus struct {
	core.DefaultAPLValueImpl
	pet *HunterPet
}

func (hunter *Hunter) newValueHunterCurrentPetFocus(rot *core.APLRotation, config *proto.APLValueHunterCurrentPetFocus) core.APLValue {
	pet := hunter.pet
	if pet == nil {
		return nil
	}
	if !pet.GetPet().HasFocusBar() {
		rot.ValidationWarning("%s does not use Focus", pet.GetPet().Label)
		return nil
	}
	return &APLValueHunterCurrentPetFocus{
		pet: pet,
	}
}
func (value *APLValueHunterCurrentPetFocus) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueHunterCurrentPetFocus) GetFloat(sim *core.Simulation) float64 {
	return value.pet.CurrentFocus()
}
func (value *APLValueHunterCurrentPetFocus) String() string {
	return "Current Pet Focus"
}

type APLValueHunterCurrentPetFocusPercent struct {
	core.DefaultAPLValueImpl
	pet *HunterPet
}

func (hunter *Hunter) newValueHunterCurrentPetFocusPercent(rot *core.APLRotation, config *proto.APLValueHunterCurrentPetFocusPercent) core.APLValue {
	pet := hunter.pet
	if pet == nil {
		return nil
	}
	if !pet.GetPet().HasFocusBar() {
		rot.ValidationWarning("%s does not use Focus", pet.GetPet().Label)
		return nil
	}
	return &APLValueHunterCurrentPetFocusPercent{
		pet: pet,
	}
}
func (value *APLValueHunterCurrentPetFocusPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueHunterCurrentPetFocusPercent) GetFloat(sim *core.Simulation) float64 {
	return value.pet.GetPet().CurrentFocusPercent()
}
func (value *APLValueHunterCurrentPetFocusPercent) String() string {
	return fmt.Sprintf("Current Pet Focus %%")
}