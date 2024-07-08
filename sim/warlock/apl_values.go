package warlock

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_WarlockShouldRecastDrainSoul:
		return warlock.newValueWarlockShouldRecastDrainSoul(rot, config.GetWarlockShouldRecastDrainSoul())
	case *proto.APLValue_WarlockShouldRefreshCorruption:
		return warlock.newValueWarlockShouldRefreshCorruption(rot, config.GetWarlockShouldRefreshCorruption())
	case *proto.APLValue_WarlockCurrentPetMana:
		return warlock.newValueWarlockCurrentPetMana(rot, config.GetWarlockCurrentPetMana())
	case *proto.APLValue_WarlockCurrentPetManaPercent:
		return warlock.newValueWarlockCurrentPetManaPercent(rot, config.GetWarlockCurrentPetManaPercent())
	default:
		return nil
	}
}

type APLValueWarlockShouldRecastDrainSoul struct {
	core.DefaultAPLValueImpl
	warlock *Warlock
}

func (warlock *Warlock) newValueWarlockShouldRecastDrainSoul(_ *core.APLRotation, _ *proto.APLValueWarlockShouldRecastDrainSoul) core.APLValue {
	return &APLValueWarlockShouldRecastDrainSoul{
		warlock: warlock,
	}
}
func (value *APLValueWarlockShouldRecastDrainSoul) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueWarlockShouldRecastDrainSoul) GetBool(sim *core.Simulation) bool {
	warlock := value.warlock

	// Assert that we're currently channeling Drain Soul.
	if warlock.ChanneledDot == nil {
		return false
	}

	var activeDrainSoul *core.Spell
	for _, spell := range warlock.DrainSoul {
		if spell.CurDot().IsActive() {
			activeDrainSoul = spell
			break
		}
	}
	if activeDrainSoul == nil {
		return false
	}

	curseRefresh := time.Duration(0)
	if warlock.ActiveCurseAura != nil {
		curseRefresh = warlock.ActiveCurseAura.RemainingDuration(sim)
	}

	hauntRefresh := 1000 * time.Second
	if warlock.HauntDebuffAuras != nil {
		hauntRefresh = warlock.HauntDebuffAuras.Get(warlock.CurrentTarget).RemainingDuration(sim) -
			warlock.Haunt.CastTime() -
			warlock.Haunt.TravelTime()
	}

	// the amount of ticks we have left, assuming we continue channeling
	dsDot := warlock.ChanneledDot
	ticksLeft := int(curseRefresh/dsDot.TickPeriod()) + 1
	ticksLeft = min(ticksLeft, int(hauntRefresh/dsDot.TickPeriod()))
	ticksLeft = min(ticksLeft, dsDot.NumTicksRemaining(sim))

	// amount of ticks we'd get assuming we recast drain soul
	recastTicks := int(curseRefresh/warlock.ApplyCastSpeed(dsDot.TickLength)) + 1
	recastTicks = min(recastTicks, int(hauntRefresh/warlock.ApplyCastSpeed(dsDot.TickLength)))
	recastTicks = min(recastTicks, int(dsDot.NumberOfTicks))

	if ticksLeft <= 0 || recastTicks <= 0 {
		return false
	}

	snapshotDmg := activeDrainSoul.ExpectedTickDamageFromCurrentSnapshot(sim, warlock.CurrentTarget) * float64(ticksLeft)
	recastDmg := activeDrainSoul.ExpectedTickDamage(sim, warlock.CurrentTarget) * float64(recastTicks)
	snapshotDPS := snapshotDmg / (time.Duration(ticksLeft) * dsDot.TickPeriod()).Seconds()
	recastDps := recastDmg / (time.Duration(recastTicks)*warlock.ApplyCastSpeed(dsDot.TickLength) + warlock.ChannelClipDelay).Seconds()

	//if sim.Log != nil {
	//	warlock.Log(sim, "Should Recast Drain Soul Calc: %.2f (%d) > %.2f (%d)", recastDps, recastTicks, snapshotDPS, ticksLeft)
	//}
	return recastDps > snapshotDPS
}
func (value *APLValueWarlockShouldRecastDrainSoul) String() string {
	return "Warlock Should Recast Drain Soul()"
}

type APLValueWarlockShouldRefreshCorruption struct {
	core.DefaultAPLValueImpl
	warlock *Warlock
	target  core.UnitReference
}

func (warlock *Warlock) newValueWarlockShouldRefreshCorruption(rot *core.APLRotation, config *proto.APLValueWarlockShouldRefreshCorruption) core.APLValue {
	target := rot.GetTargetUnit(config.TargetUnit)
	if target.Get() == nil {
		return nil
	}

	return &APLValueWarlockShouldRefreshCorruption{
		warlock: warlock,
		target:  target,
	}
}
func (value *APLValueWarlockShouldRefreshCorruption) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueWarlockShouldRefreshCorruption) GetBool(sim *core.Simulation) bool {
	warlock := value.warlock
	target := value.target.Get()

	var dot *core.Dot
	for _, spell := range warlock.Corruption {
		dot = spell.Dot(target)
		if dot.IsActive() {
			break
		}
	}
	if dot == nil || !dot.IsActive() {
		return true
	}

	attackTable := warlock.AttackTables[target.UnitIndex][dot.Spell.CastType]

	// check if reapplying corruption is worthwhile
	snapshotCrit := dot.SnapshotCritChance
	snapshotMult := dot.SnapshotAttackerMultiplier * (snapshotCrit*(dot.Spell.CritMultiplier(attackTable)-1) + 1)

	curCrit := dot.Spell.SpellCritChance(target)
	curDmg := dot.Spell.AttackerDamageMultiplier(attackTable) * (curCrit*(dot.Spell.CritMultiplier(attackTable)-1) + 1)

	relDmgInc := curDmg / snapshotMult

	snapshotDmg := dot.Spell.ExpectedTickDamageFromCurrentSnapshot(sim, target)
	snapshotDmg *= float64(sim.GetRemainingDuration()) / float64(dot.TickPeriod())
	snapshotDmg *= relDmgInc - 1
	snapshotDmg -= dot.Spell.ExpectedTickDamageFromCurrentSnapshot(sim, target)

	//if sim.Log != nil {
	//	warlock.Log(sim, "Relative Corruption Inc: [%.2f], expected dmg gain: [%.2f]", relDmgInc, snapshotDmg)
	//}

	return relDmgInc > 1.15 || snapshotDmg > 10000
}
func (value *APLValueWarlockShouldRefreshCorruption) String() string {
	return "Warlock Should Refresh Corruption()"
}

type APLValueWarlockCurrentPetMana struct {
	core.DefaultAPLValueImpl
	pet *WarlockPet
}

func (warlock *Warlock) newValueWarlockCurrentPetMana(rot *core.APLRotation, config *proto.APLValueWarlockCurrentPetMana) core.APLValue {
	pet := warlock.ActivePet
	if pet == nil {
		return nil
	}
	if !pet.GetPet().HasManaBar() {
		rot.ValidationWarning("%s does not use Mana", pet.GetPet().Label)
		return nil
	}
	return &APLValueWarlockCurrentPetMana{
		pet: pet,
	}
}
func (value *APLValueWarlockCurrentPetMana) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueWarlockCurrentPetMana) GetFloat(sim *core.Simulation) float64 {
	return value.pet.GetPet().CurrentMana()
}
func (value *APLValueWarlockCurrentPetMana) String() string {
	return "Current Pet Mana"
}

type APLValueWarlockCurrentPetManaPercent struct {
	core.DefaultAPLValueImpl
	pet *WarlockPet
}

func (warlock *Warlock) newValueWarlockCurrentPetManaPercent(rot *core.APLRotation, config *proto.APLValueWarlockCurrentPetManaPercent) core.APLValue {
	pet := warlock.ActivePet
	if pet == nil {
		return nil
	}
	if !pet.GetPet().HasManaBar() {
		rot.ValidationWarning("%s does not use Mana", pet.GetPet().Label)
		return nil
	}
	return &APLValueWarlockCurrentPetManaPercent{
		pet: pet,
	}
}
func (value *APLValueWarlockCurrentPetManaPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueWarlockCurrentPetManaPercent) GetFloat(sim *core.Simulation) float64 {
	return value.pet.GetPet().CurrentManaPercent()
}
func (value *APLValueWarlockCurrentPetManaPercent) String() string {
	return fmt.Sprintf("Current Pet Mana %%")
}
