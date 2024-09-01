package core

type ProcMask uint32

// Returns whether there is any overlap between the given masks.
func (pm ProcMask) Matches(other ProcMask) bool {
	return (pm & other) != 0
}

// Actual Blizzard flag values:
// 1  1        = Triggered by script
// 2  2        = Triggers on kill
// 3  4        = Melee auto attack
// 4  8        = On take melee auto attack
// 5  16       = Melee special attack / melee damage
// 6  32       = On take melee special attack
// 7  64       = Ranged auto attack
// 8  128      = On take ranged auto attack
// 9  256      = Ranged special attack / ranged damage
// 10 512      = On take ranged special attack
// 11 1024     = ???? On use combo points? Shapeshift? Change stance? Gain buff? Some rogue stuff
// 12 2048     = ???? Rogue related? Script?
// 13 4096     = ???? Stealth related? Script? On gain/lose stealth? Also possibly on stance change
// 14 8192     = On spell hit on you
// 15 16384    = Cast heal
// 16 32768    = On get healed
// 17 65536    = Deal spell damage
// 18 131072   = On take spell damage
// 19 262144   = Deal periodic damage
// 20 524288   = On take periodic damage
// 21 1048576  = On take any damage
// 22 2097152  = On Apply debuff
// 23 4194304  = ???? On have debuff applied to you? really bizarre mask
// 24 8388608  = On offhand attack
// 25 16777216 = What the fuck?

// Single-bit masks. These don't need to match Blizzard's values.
const (
	// Default value is invalid, to force authors to think about proc masks.
	ProcMaskUnknown ProcMask = 0

	ProcMaskEmpty ProcMask = 1 << iota
	ProcMaskMeleeMHAuto
	ProcMaskMeleeOHAuto
	ProcMaskMeleeMHSpecial
	ProcMaskMeleeOHSpecial
	ProcMaskRangedAuto
	ProcMaskRangedSpecial
	ProcMaskSpellDamage
	ProcMaskSpellHealing
	// Special mask for procs that can trigger things
	ProcMaskProc
	// Mask for FT weapon and rogue poisons, seems to be spell procs from a weapon imbue
	ProcMaskWeaponProc
	// Mask for Fiery Weapon and Blazefury Medalion that trigger melee procs like Art of War Rune or Vengeance Talent
	ProcMaskTriggerInstant
)

const (
	ProcMaskMeleeMH = ProcMaskMeleeMHAuto | ProcMaskMeleeMHSpecial
	ProcMaskMeleeOH = ProcMaskMeleeOHAuto | ProcMaskMeleeOHSpecial
	// Equivalent to in-game mask of 4.
	ProcMaskMeleeWhiteHit = ProcMaskMeleeMHAuto | ProcMaskMeleeOHAuto
	// Equivalent to in-game mask of 68.
	ProcMaskWhiteHit = ProcMaskMeleeMHAuto | ProcMaskMeleeOHAuto | ProcMaskRangedAuto
	// Equivalent to in-game mask of 16.
	ProcMaskMeleeSpecial = ProcMaskMeleeMHSpecial | ProcMaskMeleeOHSpecial
	// Equivalent to in-game mask of 272.
	ProcMaskMeleeOrRangedSpecial = ProcMaskMeleeSpecial | ProcMaskRangedSpecial
	// Equivalent to in-game mask of 20.
	ProcMaskMelee = ProcMaskMeleeWhiteHit | ProcMaskMeleeSpecial
	// Equivalent to in-game mask of 320.
	ProcMaskRanged = ProcMaskRangedAuto | ProcMaskRangedSpecial
	// Equivalent to in-game mask of 340.
	ProcMaskMeleeOrRanged = ProcMaskMelee | ProcMaskRanged

	ProcMaskDirect = ProcMaskMelee | ProcMaskRanged | ProcMaskSpellDamage

	ProcMaskSpecial = ProcMaskMeleeOrRangedSpecial | ProcMaskSpellDamage

	ProcMaskMeleeOrProc = ProcMaskMelee | ProcMaskProc
	ProcMaskSpellOrProc = ProcMaskSpellDamage | ProcMaskProc
)

// Possible outcomes of any hit/damage roll.
type HitOutcome uint16

// Returns whether there is any overlap between the given masks.
func (ho HitOutcome) Matches(other HitOutcome) bool {
	return (ho & other) != 0
}

// Single-bit outcomes.
const (
	OutcomeEmpty HitOutcome = 0

	// These bits are set by the hit roll
	OutcomeMiss HitOutcome = 1 << iota
	OutcomeHit
	OutcomeDodge
	OutcomeGlance
	OutcomeParry
	OutcomeBlock

	// These bits are set by the crit and damage rolls.
	OutcomeCrit
	OutcomeCrush

	OutcomePartial1_4 // 1/4 of the spell was resisted.
	OutcomePartial2_4 // 2/4 of the spell was resisted.
	OutcomePartial3_4 // 3/4 of the spell was resisted.
)

const (
	OutcomePartial = OutcomePartial1_4 | OutcomePartial2_4 | OutcomePartial3_4
	OutcomeLanded  = OutcomeHit | OutcomeCrit | OutcomeCrush | OutcomeGlance | OutcomeBlock
)

func (ho HitOutcome) String() string {
	if ho.Matches(OutcomeMiss) {
		return "Miss"
	} else if ho.Matches(OutcomeDodge) {
		return "Dodge"
	} else if ho.Matches(OutcomeParry) {
		return "Parry"
	} else if ho.Matches(OutcomeGlance) {
		return "Glance" + ho.PartialResistString()
	} else if ho.Matches(OutcomeBlock) {
		if ho.Matches(OutcomeCrit) {
			return "CriticalBlock"
		} else {
			return "Block"
		}
	} else if ho.Matches(OutcomeCrit) {
		return "Crit" + ho.PartialResistString()
	} else if ho.Matches(OutcomeHit) {
		return "Hit" + ho.PartialResistString()
	} else if ho.Matches(OutcomeCrush) {
		return "Crush"
	} else {
		return "Empty"
	}
}

func (ho HitOutcome) PartialResistString() string {
	if ho.Matches(OutcomePartial1_4) {
		return " (25% Resist)"
	} else if ho.Matches(OutcomePartial2_4) {
		return " (50% Resist)"
	} else if ho.Matches(OutcomePartial3_4) {
		return " (75% Resist)"
	} else {
		return ""
	}
}

// Other flags
type SpellFlag uint64

// Returns whether there is any overlap between the given masks.
func (se SpellFlag) Matches(other SpellFlag) bool {
	return (se & other) != 0
}

const (
	SpellFlagNone                    SpellFlag = 0
	SpellFlagIgnoreResists           SpellFlag = 1 << iota // skip spell resist/armor
	SpellFlagIgnoreTargetModifiers                         // skip target damage modifiers
	SpellFlagIgnoreAttackerModifiers                       // skip attacker damage modifiers
	SpellFlagBinary                                        // Does not do partial resists or blocks and could need a different hit roll.
	SpellFlagChanneled                                     // Spell is channeled
	SpellFlagDisease                                       // Spell is categorized as disease
	SpellFlagPoison                                        // Spell is categorized as poison
	SpellFlagHelpful                                       // For healing spells / buffs.
	SpellFlagMeleeMetrics                                  // Marks a spell as a melee ability for metrics.
	SpellFlagNoOnCastComplete                              // Disables the OnCastComplete callback.
	SpellFlagNoMetrics                                     // Disables metrics for a spell.
	SpellFlagNoLogs                                        // Disables logs for a spell.
	SpellFlagAPL                                           // Indicates this spell can be used from an APL rotation.
	SpellFlagMCD                                           // Indicates this spell is a MajorCooldown.
	SpellFlagNoOnDamageDealt                               // Disables OnSpellHitDealt and OnPeriodicDamageDealt aura callbacks for this spell.
	SpellFlagPrepullOnly                                   // Indicates this spell should only be used during prepull. Not enforced, just a signal for the APL UI.
	SpellFlagEncounterOnly                                 // Indicates this spell should only be used during the encounter (not prepull). Not enforced, just a signal for the APL UI.
	SpellFlagPotion                                        // Indicates this spell is a potion spell.
	SpellFlagOffensiveEquipment                            // Indicates this spell is an offensive equippable item activation spell
	SpellFlagDefensiveEquipment                            // Indicates this spell a defensive equippable item activation spell
	SpellFlagResetAttackSwing                              // Indicates this spell resets the melee swing timer.
	SpellFlagCastTimeNoGCD                                 // Indicates this spell is off the GCD (e.g. hunter's Auto Shot)
	SpellFlagCastWhileCasting                              // Indicates this spell can be cast while another spell is being cast (e.g. mage's Fire Blast with Overheat rune)
	SpellFlagPureDot                                       // Indicates this spell is a dot with no initial damage component

	SpellFlagSupressExtraAttack  // Mask for Seal of Righteousness, it does not proc Wild Strikes
	SpellFlagSuppressWeaponProcs // Indicates this spell cannot proc weapon chance on hits or enchants
	SpellFlagSuppressEquipProcs  // Indicates this spell cannot proc Equip procs

	SpellFlagBatchStopAttackMacro // Indicates this spell is being cast in a Macro with a stopattack following it

	// Used to let agents categorize their spells.
	SpellFlagAgentReserved1
	SpellFlagAgentReserved2
	SpellFlagAgentReserved3
	SpellFlagAgentReserved4
	SpellFlagAgentReserved5

	SpellFlagIgnoreModifiers = SpellFlagIgnoreAttackerModifiers | SpellFlagIgnoreTargetModifiers
)

/*
outcome roll hit/miss/crit/glance (assigns Outcome mask) -> If Hit, Crit Roll -> damage (applies metrics) -> trigger proc

So in TBC it looks like they just gave it the cannot miss flag even though they also switched its defense type to physical (??)
the damage type is holy, which ignores armor and as it is magic so it can be partially resisted (due to level resistance).
however it also gains the physical bit mask as I explain in a post above

so there is no hit roll, there is a melee crit roll, a spell damage roll, and melee "on hit"

ok so I did some more testing on this.
Judgement of Blood correctly gets the "always hit" (aka cannot miss flag applied to it) --
its only mitigation events are partial resists at the correct rates
however Judgement of Command is broken. even though it has the "always hit" flag it seems to
be ignored because it is procced by an intermediary dummy spell which does not have the "cannot miss" flag applied to it lmao.
for some god forsaken reason Judgement of Command is ALSO a dummy which then casts the correct Judgement of Command
which deals damage, and this dummy can miss, lmao
I got ~16.4% resists in about almost 96 casts which suggests it uses the spell hit check,
which makes sense because its defensetype is set to 1, Magic

arcane shot - ranged hit, spell dmg, procs special ranged
	OutcomeRollRanged, School Arcane, ProcMask - RangedSpecial

judgement of blood - physical hit/crit, spell damage, "cannot miss", procs special melee and ranged
	Damage is (weapon damage + spell power)*0.7*(bonus holy damage against target)+flat bonus damage
	OtherFlagCannotMiss, OutcomeRollSpecial, School Holy (base damage = weapon damage range), Multiplier 70%

judgement of command - spell hit, melee crit, spell damage, procs special melee and ranged
	OutcomeRollSpell, School Holy


moonfire - spell hit, spell dmg, dot dmg, procs spell hit
stormstrike - melee hit, melee dmg, procs special melee
rupture -

wotlk
shadowflame - requires each 'effect' to have its own school.
*/
