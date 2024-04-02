package core

import (
	"fmt"
	"slices"

	"github.com/wowsims/sod/sim/core/proto"
)

type APLValueRuneIsEquipped struct {
	DefaultAPLValueImpl
	character *Character
	rune      Rune
}

func (rot *APLRotation) newValueRuneIsEquipped(config *proto.APLValueRuneIsEquipped) APLValue {
	character := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).GetCharacter()
	rune := rot.GetAPLRune(config.GetRuneId())

	return &APLValueRuneIsEquipped{
		character: character,
		rune:      rune,
	}
}
func (value *APLValueRuneIsEquipped) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueRuneIsEquipped) GetBool(sim *Simulation) bool {
	return slices.Contains(value.character.Equipment.GetRuneIds(), value.rune.ID)
}
func (value *APLValueRuneIsEquipped) String() string {
	return fmt.Sprintf("Rune Equipped(%s)", value.rune.ID)
}
