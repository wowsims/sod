package feral

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/druid"
)

func (cat *FeralDruid) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_CatExcessEnergy:
		return cat.newValueCatExcessEnergy(rot, config.GetCatExcessEnergy())
	case *proto.APLValue_CatNewSavageRoarDuration:
		return cat.newValueCatNewSavageRoarDuration(rot, config.GetCatNewSavageRoarDuration())
	case *proto.APLValue_CatEnergyAfterDuration:
		return cat.newValueCatEnergyAfterDuration(rot, config.GetCatEnergyAfterDuration())
	default:
		return nil
	}
}

type APLValueCatExcessEnergy struct {
	core.DefaultAPLValueImpl
	cat *FeralDruid
}

func (cat *FeralDruid) newValueCatExcessEnergy(rot *core.APLRotation, config *proto.APLValueCatExcessEnergy) core.APLValue {
	return &APLValueCatExcessEnergy{
		cat: cat,
	}
}
func (value *APLValueCatExcessEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCatExcessEnergy) GetFloat(sim *core.Simulation) float64 {
	cat := value.cat
	pendingPool := PoolingActions{}
	pendingPool.create(4)

	/* TODO
	curCp := cat.ComboPoints()
	simTimeRemain := sim.GetRemainingDuration()
	rakeDot := cat.Rake.CurDot()
	ripDot := cat.Rip.CurDot()
	mangleRefreshPending := cat.bleedAura.IsActive() && cat.bleedAura.RemainingDuration(sim) < (simTimeRemain-time.Second)
	endThresh := time.Second * 10

	if ripDot.IsActive() && (ripDot.RemainingDuration(sim) < simTimeRemain-endThresh) && curCp == 5 {
		ripCost := core.Ternary(cat.berserkExpectedAt(sim, ripDot.ExpiresAt()), cat.Rip.DefaultCast.Cost*0.5, cat.Rip.DefaultCast.Cost)
		pendingPool.addAction(ripDot.ExpiresAt(), ripCost)
		cat.ripRefreshPending = true
	}
	if rakeDot.IsActive() && (rakeDot.RemainingDuration(sim) < simTimeRemain-rakeDot.Duration) {
		rakeCost := core.Ternary(cat.berserkExpectedAt(sim, rakeDot.ExpiresAt()), cat.Rake.DefaultCast.Cost*0.5, cat.Rake.DefaultCast.Cost)
		pendingPool.addAction(rakeDot.ExpiresAt(), rakeCost)
	}
	if mangleRefreshPending {
		mangleCost := core.Ternary(cat.berserkExpectedAt(sim, cat.bleedAura.ExpiresAt()), cat.MangleCat.DefaultCast.Cost*0.5, cat.MangleCat.DefaultCast.Cost)
		pendingPool.addAction(cat.bleedAura.ExpiresAt(), mangleCost)
	}
	if cat.SavageRoarAura.IsActive() {
		roarCost := core.Ternary(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.DefaultCast.Cost*0.5, cat.SavageRoar.DefaultCast.Cost)
		pendingPool.addAction(cat.SavageRoarAura.ExpiresAt(), roarCost)
	}
	*/

	pendingPool.sort()

	floatingEnergy := pendingPool.calcFloatingEnergy(cat, sim)
	return cat.CurrentEnergy() - floatingEnergy
}
func (value *APLValueCatExcessEnergy) String() string {
	return "Cat Excess Energy()"
}

type APLValueCatNewSavageRoarDuration struct {
	core.DefaultAPLValueImpl
	cat *FeralDruid
}

func (cat *FeralDruid) newValueCatNewSavageRoarDuration(rot *core.APLRotation, config *proto.APLValueCatNewSavageRoarDuration) core.APLValue {
	return &APLValueCatNewSavageRoarDuration{
		cat: cat,
	}
}
func (value *APLValueCatNewSavageRoarDuration) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueCatNewSavageRoarDuration) GetDuration(sim *core.Simulation) time.Duration {
	cat := value.cat
	return cat.SavageRoarDurationTable[cat.ComboPoints()]
}
func (value *APLValueCatNewSavageRoarDuration) String() string {
	return "New Savage Roar Duration()"
}

func (cat *FeralDruid) NewAPLAction(rot *core.APLRotation, config *proto.APLAction) core.APLActionImpl {
	switch config.Action.(type) {
	case *proto.APLAction_CatOptimalRotationAction:
		return cat.newActionCatOptimalRotationAction(rot, config.GetCatOptimalRotationAction())
	default:
		return nil
	}
}

type APLActionCatOptimalRotationAction struct {
	cat        *FeralDruid
	lastAction time.Duration
}

func (impl *APLActionCatOptimalRotationAction) GetInnerActions() []*core.APLAction { return nil }
func (impl *APLActionCatOptimalRotationAction) GetAPLValues() []core.APLValue      { return nil }
func (impl *APLActionCatOptimalRotationAction) Finalize(*core.APLRotation)         {}
func (impl *APLActionCatOptimalRotationAction) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}
func (impl *APLActionCatOptimalRotationAction) GetSpellFromAction() *core.Spell { return nil }

func (cat *FeralDruid) newActionCatOptimalRotationAction(_ *core.APLRotation, config *proto.APLActionCatOptimalRotationAction) core.APLActionImpl {
	cat.setupRotation(config)

	return &APLActionCatOptimalRotationAction{
		cat: cat,
	}
}

func (action *APLActionCatOptimalRotationAction) IsReady(sim *core.Simulation) bool {
	return sim.CurrentTime > action.lastAction
}

func (action *APLActionCatOptimalRotationAction) Execute(sim *core.Simulation) {
	cat := action.cat

	// If a melee swing resulted in an Omen or Wild Strikes proc, then schedule the
	// next player decision based on latency.
	if (cat.Talents.OmenOfClarity && cat.ClearcastingAura.RemainingDuration(sim) == cat.ClearcastingAura.Duration) || (cat.WildStrikesBuffAura != nil && cat.WildStrikesBuffAura.RemainingDuration(sim) == cat.WildStrikesBuffAura.Duration) {
		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		kickTime := max(cat.NextGCDAt(), sim.CurrentTime+cat.latency)
		cat.NextRotationAction(sim, kickTime)
	}

	if cat.GCD.IsReady(sim) && (cat.rotationAction == nil || sim.CurrentTime >= cat.rotationAction.NextActionAt) {
		cat.OnGCDReady(sim)
	}

	action.lastAction = sim.CurrentTime
}

func (action *APLActionCatOptimalRotationAction) Reset(*core.Simulation) {
	action.lastAction = core.DurationFromSeconds(-100)
}

func (action *APLActionCatOptimalRotationAction) String() string {
	return "Execute Optimal Cat Action()"
}

type APLValueCatEnergyAfterDuration struct {
	core.DefaultAPLValueImpl
	cat                  *FeralDruid
	condition            core.APLValue
	staffOfTheGladeTimer *core.Aura
	staffOfTheGladeBuff  *core.Aura
	has2PieceScarletTier bool
}

func (cat *FeralDruid) newValueCatEnergyAfterDuration(rot *core.APLRotation, config *proto.APLValueCatEnergyAfterDuration) core.APLValue {
	conditionVal := rot.NewAPLValue(config.Condition)
	if conditionVal == nil {
		return nil
	}

	if conditionVal.Type() != proto.APLValueType_ValueTypeDuration {
		rot.ValidationWarning("Value must be a duration type!")
		return nil
	}

	return &APLValueCatEnergyAfterDuration{
		cat:                  cat,
		condition:            conditionVal,
		staffOfTheGladeTimer: cat.GetAuraByID(core.ActionID{SpellID: 1231380}),
		staffOfTheGladeBuff:  cat.GetAuraByID(core.ActionID{SpellID: 1231381}),
		has2PieceScarletTier: cat.HasSetBonus(druid.ItemSetWaywatcherFerocity, 2),
	}
}
func (value *APLValueCatEnergyAfterDuration) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func getDotTicksOverDuration(spell *core.Spell, sim *core.Simulation, duration time.Duration) float64 {
	ticks := 0.0
	for _, target := range sim.Encounter.Targets {
		dot := spell.Dot(&target.Unit)
		if dot.IsActive() {
			effectiveDuration := min(duration, dot.RemainingDuration(sim))
			timeToNextTick := dot.NextTickAt() - sim.CurrentTime
			if timeToNextTick <= effectiveDuration {
				ticks += 1 + math.Floor(float64(effectiveDuration-timeToNextTick)/float64(dot.TickPeriod()))
			}
		}
	}
	return ticks
}
func (value *APLValueCatEnergyAfterDuration) GetFloat(sim *core.Simulation) float64 {
	duration := value.condition.GetDuration(sim)
	cat := value.cat
	energy := cat.CurrentEnergy()

	if duration <= 0 {
		return energy
	}

	timeToNextEnergyTick := cat.NextEnergyTickAt() - sim.CurrentTime
	if timeToNextEnergyTick <= duration {
		energyTicks := 1 + math.Floor(float64(duration-timeToNextEnergyTick)/float64(core.EnergyTickDuration))
		currentEnergyPerTick := cat.EnergyTickMultiplier * core.EnergyPerTick

		// Staff of the Glade equipped, buff not active but activation timer is active.
		// Need to calculate how many ticks will fall inside of the buff, if any.
		if value.staffOfTheGladeBuff != nil && !value.staffOfTheGladeBuff.IsActive() && value.staffOfTheGladeTimer.IsActive() {
			timeToStaffActive := value.staffOfTheGladeTimer.RemainingDuration(sim)
			if timeToStaffActive < timeToNextEnergyTick {
				energy += energyTicks * currentEnergyPerTick * druid.StaffOfTheGladeEnergyMult
				energyTicks = 0
			} else if timeToStaffActive < duration {
				ticksBeforeStaffActive := 1 + math.Floor(float64(timeToStaffActive-timeToNextEnergyTick)/float64(core.EnergyTickDuration))
				energy += (energyTicks - ticksBeforeStaffActive) * currentEnergyPerTick * druid.StaffOfTheGladeEnergyMult
				energyTicks = ticksBeforeStaffActive
			}
		}

		energy += energyTicks * currentEnergyPerTick
	}

	if value.has2PieceScarletTier {
		ticks := getDotTicksOverDuration(cat.Rake.Spell, sim, duration)
		ticks += getDotTicksOverDuration(cat.Rip.Spell, sim, duration)
		energy += ticks * druid.WaywatcherFerocity2pEnergy
	}

	return energy
}
func (value *APLValueCatEnergyAfterDuration) String() string {
	return "Cat Energy After Duration()"
}
