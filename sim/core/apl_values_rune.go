package core

import (
	"fmt"
	"slices"

	"github.com/wowsims/sod/sim/core/proto"
)

type APLRune struct {
	id       int32
	equipped bool
}

type APLValueRuneIsEquipped struct {
	DefaultAPLValueImpl
	character *Character
	rune      APLRune
}

func (rot *APLRotation) newValueRuneIsEquipped(config *proto.APLValueRuneIsEquipped) APLValue {
	character := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).GetCharacter()
	spellId := config.GetRuneId().GetSpellId()
	rune := APLRune{id: spellId, equipped: slices.Contains(character.Equipment.GetRuneIds(), spellId)}

	return &APLValueRuneIsEquipped{
		character: character,
		rune:      rune,
	}
}
func (value *APLValueRuneIsEquipped) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueRuneIsEquipped) GetBool(sim *Simulation) bool {
	return value.rune.equipped
}
func (value *APLValueRuneIsEquipped) String() string {
	return fmt.Sprintf("Rune Equipped(%d)", value.rune.id)
}
