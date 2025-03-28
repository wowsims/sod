package core

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

// Function for applying permanent effects to an Agent.
//
// Passing Character instead of Agent would work for almost all cases,
// but there are occasionally class-specific item effects.
type ApplyEffect func(agent Agent)

var itemEffects = map[int32]ApplyEffect{}
var weaponEffects = map[int32]ApplyEffect{}
var enchantEffects = map[int32]ApplyEffect{}

// IDs of item effects which should be used for tests.
var itemEffectsForTest []int32

// This value can be set before adding item effects, to control whether they are included in tests.
var AddEffectsToTest = true

func HasItemEffect(id int32) bool {
	_, ok := itemEffects[id]
	return ok
}
func HasItemEffectForTest(id int32) bool {
	return slices.Contains(itemEffectsForTest, id)
}

func HasWeaponEffect(id int32) bool {
	_, ok := weaponEffects[id]
	return ok
}

func HasEnchantEffect(id int32) bool {
	_, ok := enchantEffects[id]
	return ok
}

// Registers an ApplyEffect function which will be called before the Sim
// starts, for any Agent that is wearing the item.
func NewItemEffect(id int32, itemEffect ApplyEffect) {
	if WITH_DB {
		if _, hasItem := ItemsByID[id]; !hasItem {
			panic(fmt.Sprintf("No item with ID: %d", id))
		}
	}

	if HasItemEffect(id) {
		panic(fmt.Sprintf("Cannot add multiple effects for one item: %d, %#v", id, itemEffect))
	}

	itemEffects[id] = itemEffect
	if AddEffectsToTest {
		itemEffectsForTest = append(itemEffectsForTest, id)
	}
}

func NewEnchantEffect(id int32, enchantEffect ApplyEffect) {
	if WITH_DB {
		if _, ok := EnchantsByEffectID[id]; !ok {
			panic(fmt.Sprintf("No enchant with ID: %d", id))
		}
	}

	if HasEnchantEffect(id) {
		panic(fmt.Sprintf("Cannot add multiple effects for one enchant: %d, %#v", id, enchantEffect))
	}

	enchantEffects[id] = enchantEffect
}

func AddWeaponEffect(id int32, weaponEffect ApplyEffect) {
	if WITH_DB {
		if _, ok := EnchantsByEffectID[id]; !ok {
			panic(fmt.Sprintf("No enchant with ID: %d", id))
		}
	}
	if HasWeaponEffect(id) {
		panic(fmt.Sprintf("Cannot add multiple effects for one item: %d, %#v", id, weaponEffect))
	}
	weaponEffects[id] = weaponEffect
}

func (equipment *Equipment) applyItemEffects(agent Agent, registeredItemEffects map[int32]bool, registeredItemEnchantEffects map[int32]bool) {
	for _, eq := range equipment {
		if applyItemEffect, ok := itemEffects[eq.ID]; ok && !registeredItemEffects[eq.ID] {
			applyItemEffect(agent)
			registeredItemEffects[eq.ID] = true
		}

		if applyEnchantEffect, ok := enchantEffects[eq.Enchant.EffectID]; ok && !registeredItemEnchantEffects[eq.Enchant.EffectID] {
			applyEnchantEffect(agent)
			registeredItemEnchantEffects[eq.Enchant.EffectID] = true
		}

		if eq.RandomSuffix.ID != 0 && eq.RandomSuffix.Stats.Equals(stats.Stats{}) {
			for _, enchantID := range eq.RandomSuffix.EnchantIDList {
				if applyEnchantEffect, ok := enchantEffects[enchantID]; ok && !registeredItemEnchantEffects[enchantID] {
					applyEnchantEffect(agent)
					registeredItemEnchantEffects[enchantID] = true
				}
			}
		}

		if applyWeaponEffect, ok := weaponEffects[eq.Enchant.EffectID]; ok && !registeredItemEnchantEffects[eq.Enchant.EffectID] {
			applyWeaponEffect(agent)
			registeredItemEnchantEffects[eq.Enchant.EffectID] = true
		}
	}
}

// Helpers for making common types of active item effects.

func NewSimpleStatItemActiveEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration, flags SpellFlag, sharedCDFunc func(*Character) Cooldown, otherEffects ApplyEffect) {
	registerCD := MakeTemporaryStatsOnUseCDRegistration(
		"ItemActive-"+strconv.Itoa(int(itemID)),
		bonus,
		duration,
		SpellConfig{
			ActionID: ActionID{ItemID: itemID},
			Flags:    SpellFlagNoOnCastComplete | flags,
		},
		func(character *Character) Cooldown {
			return Cooldown{
				Timer:    character.NewTimer(),
				Duration: cooldown,
			}
		},
		sharedCDFunc,
	)

	if otherEffects == nil {
		NewItemEffect(itemID, registerCD)
	} else {
		NewItemEffect(itemID, func(agent Agent) {
			registerCD(agent)
			otherEffects(agent)
		})
	}
}

// No shared CD
func NewSimpleStatItemEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration, flags SpellFlag) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, flags, func(character *Character) Cooldown {
		return Cooldown{}
	}, nil)
}

func NewSimpleStatOffensiveTrinketEffectWithOtherEffects(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration, otherEffects ApplyEffect) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, SpellFlagOffensiveEquipment, func(character *Character) Cooldown {
		return Cooldown{
			Timer:    character.GetOffensiveTrinketCD(),
			Duration: duration,
		}
	}, otherEffects)
}
func NewSimpleStatOffensiveTrinketEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatOffensiveTrinketEffectWithOtherEffects(itemID, bonus, duration, cooldown, nil)
}

func NewSimpleStatDefensiveTrinketEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, SpellFlagDefensiveEquipment, func(character *Character) Cooldown {
		return Cooldown{
			Timer:    character.GetDefensiveTrinketCD(),
			Duration: duration,
		}
	}, nil)
}

// TODO: These should ideally be done at the AttackTable level
func NewMobTypeAttackPowerEffect(itemID int32, mobTypes []proto.MobType, bonus float64) {
	NewItemEffect(itemID, func(agent Agent) {
		character := agent.GetCharacter()

		if !slices.Contains(mobTypes, character.CurrentTarget.MobType) {
			return
		}

		aura := MakePermanent(character.GetOrRegisterAura(Aura{
			Label: fmt.Sprintf("Mob type Attack Power Bonus - %s (%d)", character.CurrentTarget.MobType, itemID),
			OnGain: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.MobTypeAttackPower += bonus
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.MobTypeAttackPower -= bonus
			},
		}))

		character.ItemSwap.RegisterProc(itemID, aura)
	})
}

func NewMobTypeSpellPowerEffect(itemID int32, mobTypes []proto.MobType, bonus float64) {
	NewItemEffect(itemID, func(agent Agent) {
		character := agent.GetCharacter()

		if !slices.Contains(mobTypes, character.CurrentTarget.MobType) {
			return
		}

		aura := MakePermanent(character.GetOrRegisterAura(Aura{
			Label: fmt.Sprintf("Mob type Spell Power Bonus - %s (%d)", character.CurrentTarget.MobType, itemID),
			OnGain: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.MobTypeSpellPower += bonus
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.MobTypeSpellPower -= bonus
			},
		}))

		character.ItemSwap.RegisterProc(itemID, aura)
	})
}
