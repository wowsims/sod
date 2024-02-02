import { Faction, SaygesFortune, Stat } from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";

import { IconEnumPickerDirection } from "../icon_enum_picker";
import {
  makeBooleanDebuffInput,
  makeBooleanIndividualBuffInput,
  makeBooleanRaidBuffInput,
  makeEnumIndividualBuffInput,
  makeMultistateIndividualBuffInput,
  makeMultistateRaidBuffInput,
	makeMultistateMultiplierDebuffInput,
  makeTristateDebuffInput,
  makeTristateIndividualBuffInput,
  makeTristateRaidBuffInput,
  withLabel
} from "../icon_inputs";
import { IconPicker } from "../icon_picker";
import { MultiIconPicker } from "../multi_icon_picker";

import { ItemStatOption, PickerStatOptions } from "./stat_options";

import * as InputHelpers from '../input_helpers';

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

export const AllStatsBuff = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 1126, 							maxLevel: 9  	},
			{ id: 5232, minLevel: 10, maxLevel: 19 	},
			{ id: 6756, minLevel: 20, maxLevel: 29 	},
			{ id: 5234, minLevel: 30, maxLevel: 39 	},
			{ id: 8907, minLevel: 40, maxLevel: 49 	},
			{ id: 9884, minLevel: 50, maxLevel: 59 	},
			{ id: 9885, minLevel: 60 							 	},
		]),
		impId: ActionId.fromSpellId(17055),
		fieldName: 'giftOfTheWild',
	}),
	'Mark of the Wild',
);

// Separate Strength buffs allow us to use a boolean pickers for Horde specifically
export const AllStatsPercentBuffAlliance = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(20217), fieldName: 'blessingOfKings'}),
	makeBooleanRaidBuffInput({actionId: () => ActionId.fromSpellId(409580), fieldName: 'aspectOfTheLion'}),
], 'Stats %');

export const AllStatsPercentBuffHorde = withLabel(
	makeBooleanRaidBuffInput({actionId: () => ActionId.fromSpellId(409580), fieldName: 'aspectOfTheLion'}),
	'Stats %'
);

export const ArmorBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 465, 									maxLevel: 9 	},
			{ id: 10290, 	minLevel: 10, maxLevel: 19 	},
			{ id: 643, 		minLevel: 20, maxLevel: 29 	},
			{ id: 10291, 	minLevel: 30, maxLevel: 39 	},
			{ id: 1032, 	minLevel: 40, maxLevel: 49 	},
			{ id: 10292, 	minLevel: 50, maxLevel: 59 	},
			{ id: 10293, 	minLevel: 60 								},
		]),
		impId: ActionId.fromSpellId(20142),
		fieldName: 'devotionAura'
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			{ id: 3013, 								maxLevel: 14 	},
			{ id: 1478, 	minLevel: 15, maxLevel: 29 	},
			{ id: 4421, 	minLevel: 30, maxLevel: 44 	},
			{ id: 10305, 	minLevel: 45 								},
		]),
		fieldName: 'scrollOfProtection'
	}),
], 'Armor');

export const StaminaBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 1243, 								maxLevel: 11 },
			{ id: 1244, 	minLevel: 12, maxLevel: 23 },
			{ id: 1245, 	minLevel: 24, maxLevel: 35 },
			{ id: 2791, 	minLevel: 36, maxLevel: 47 },
			{ id: 10937, 	minLevel: 48, maxLevel: 60 },
			{ id: 10938, 	minLevel: 60 },
		]),
		impId: ActionId.fromSpellId(14767),
		fieldName: 'powerWordFortitude'
	}),
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 6307, 	minLevel: 4, 	maxLevel: 13 	},
			{ id: 7804, 	minLevel: 14, maxLevel: 25 	},
			{ id: 7805, 	minLevel: 26, maxLevel: 37 	},
			{ id: 11766, 	minLevel: 38, maxLevel: 49 	},
			{ id: 11767, 	minLevel: 50 								},
		]),
		impId: ActionId.fromSpellId(18696),
		fieldName: 'bloodPact'
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			{ id: 1180, 	minLevel: 5, 	maxLevel: 19 	},
			{ id: 1711, 	minLevel: 20, maxLevel: 34 	},
			{ id: 4422, 	minLevel: 35, maxLevel: 49 	},
			{ id: 10307, 	minLevel: 50 								},
		]),
		fieldName: 'scrollOfStamina'
	}),
], 'Stamina');

// Separate Strength buffs allow us to use boolean pickers for each
export const PaladinPhysicalBuff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 19740, minLevel: 4,		maxLevel: 11 	},
			{ id: 19834, minLevel: 12, 	maxLevel: 21 	},
			{ id: 19835, minLevel: 22, 	maxLevel: 31 	},
			{ id: 19836, minLevel: 32, 	maxLevel: 41 	},
			{ id: 19837, minLevel: 42, 	maxLevel: 51 	},
			{ id: 19838, minLevel: 52, 	maxLevel: 59 	},
			{ id: 25291, minLevel: 60 								},
		]),
		impId: ActionId.fromSpellId(20048),
		fieldName: 'blessingOfMight',
	}),
	makeBooleanRaidBuffInput({actionId: () => ActionId.fromSpellId(425600), fieldName: 'hornOfLordaeron'}),
], 'Paladin Physical');

export const StrengthBuffHorde = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 8075, 	minLevel: 10, maxLevel: 23 	},
			{ id: 8160, 	minLevel: 24, maxLevel: 37 	},
			{ id: 8161, 	minLevel: 38, maxLevel: 51 	},
			{ id: 10442, 	minLevel: 52, maxLevel: 59 	},
			{ id: 25361, 	minLevel: 60 								},
		]),
		impId: ActionId.fromSpellId(16295),
		fieldName: 'strengthOfEarthTotem'
	}),
	'Strength',
);;
export const GraceOfAir = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 8835, 	minLevel: 42, maxLevel: 55 	},
			{ id: 10627, 	minLevel: 56, maxLevel: 59 	},
			{ id: 25359, 	minLevel: 60 								},
		]),
		impId: ActionId.fromSpellId(16295),
		fieldName: 'graceOfAirTotem',
	}),
	'Agility',
);

export const IntellectBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 1459, 								maxLevel: 13 	},
			{ id: 1460, 	minLevel: 14, maxLevel: 27 	},
			{ id: 1461, 	minLevel: 28, maxLevel: 41 	},
			{ id: 10156, 	minLevel: 42, maxLevel: 55 	},
			{ id: 10157, 	minLevel: 56 								},
		]),
		fieldName: 'arcaneBrilliance'
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			{ id: 955, 		minLevel: 5, 	maxLevel: 19 	},
			{ id: 2290, 	minLevel: 20, maxLevel: 34 	},
			{ id: 4419, 	minLevel: 35, maxLevel: 49 	},
			{ id: 10308, 	minLevel: 50 								},
		]),
		fieldName: 'scrollOfIntellect'
	}),
], 'Intellect');

export const SpiritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 14752, minLevel: 30, maxLevel: 39 },
			{ id: 14818, minLevel: 40, maxLevel: 49 },
			{ id: 14819, minLevel: 50, maxLevel: 59 },
			{ id: 27841, minLevel: 60 							},
		]),
		fieldName: 'divineSpirit',
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			{ id: 1181, 								maxLevel: 14 	},
			{ id: 1712, 	minLevel: 15, maxLevel: 29 	},
			{ id: 4424, 	minLevel: 30, maxLevel: 44 	},
			{ id: 10306, 	minLevel: 45 								},
		]),
		fieldName: 'scrollOfSpirit',
	}),
], 'Spirit');

export const BattleShoutBuff = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 6673, 								maxLevel: 11 	},
			{ id: 5242, 	minLevel: 12, maxLevel: 21 	},
			{ id: 6192, 	minLevel: 22, maxLevel: 31 	},
			{ id: 11549, 	minLevel: 32, maxLevel: 41 	},
			{ id: 11550, 	minLevel: 42, maxLevel: 51 	},
			{ id: 11551, 	minLevel: 52, maxLevel: 59 	},
			{ id: 25289, 	minLevel: 60 								},
		]),
		impId: ActionId.fromSpellId(12861),
		fieldName: 'battleShout'
	}),
	'Battle Shout',
);

export const TrueshotAuraBuff = withLabel(
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 20906, minLevel: 40 },
		]),
		fieldName: 'trueshotAura'
	}),
	'Trueshot Aura',
);

// export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput([
// ], 'Attack Power %', 1, 40);

export const DamageReductionPercentBuff = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 20911, minLevel: 30, maxLevel: 39 },
			{ id: 20912, minLevel: 40, maxLevel: 49 },
			{ id: 20913, minLevel: 50, maxLevel: 59 },
			{ id: 20914, minLevel: 60 							},
		]),
		fieldName: 'blessingOfSanctuary'
	}),
	'Blessing of Sanctuary',
);

export const ResistanceBuff = InputHelpers.makeMultiIconInput([
	// Shadow
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 976, 		minLevel: 30, maxLevel: 41 	},
			{ id: 10957, 	minLevel: 42, maxLevel: 55 	},
			{ id: 10958, 	minLevel: 56 								},
		]),
		fieldName: 'shadowProtection'
	}),
	// Nature
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 10595, minLevel: 30, maxLevel: 43 },
			{ id: 10600, minLevel: 44, maxLevel: 59 },
			{ id: 10601, minLevel: 60 							},
		]),
		fieldName: 'natureResistanceTotem',
		faction: Faction.Horde
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 20043, minLevel: 46, maxLevel: 55 },
			{ id: 20190, minLevel: 56 							},
		]),
		fieldName: 'aspectOfTheWild'
	}),
	// Fire
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 19891, minLevel: 36, maxLevel: 47 },
			{ id: 19899, minLevel: 48, maxLevel: 59 },
			{ id: 19900, minLevel: 60 							},
		]),
		fieldName: 'fireResistanceAura',
		faction: Faction.Alliance
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 8184, 	minLevel: 28, maxLevel: 41 	},
			{ id: 10537, 	minLevel: 42, maxLevel: 57 	},
			{ id: 10538, 	minLevel: 58 								},
		]),
		fieldName: 'fireResistanceTotem',
		faction: Faction.Horde
	}),
	// Frost
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 19888, minLevel: 32, maxLevel: 43 },
			{ id: 19897, minLevel: 44, maxLevel: 55 },
			{ id: 19898, minLevel: 56 							},
		]),
		fieldName: 'frostResistanceAura',
		faction: Faction.Alliance
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 8181, 	minLevel: 24, maxLevel: 37 	},
			{ id: 10478, 	minLevel: 38, maxLevel: 53 	},
			{ id: 10479, 	minLevel: 54 								},
		]),
		fieldName: 'frostResistanceTotem',
		faction: Faction.Horde
	}),
], 'Resistances');

export const BlessingOfWisdom = withLabel(
	makeTristateIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 19742, minLevel: 14, maxLevel: 23 },
			{ id: 19850, minLevel: 24, maxLevel: 33 },
			{ id: 19852, minLevel: 34, maxLevel: 43 },
			{ id: 19853, minLevel: 44, maxLevel: 53 },
			{ id: 19854, minLevel: 54, maxLevel: 59 },
			{ id: 25290, minLevel: 60 							},
		]),
		impId: ActionId.fromSpellId(20245),
		fieldName: 'blessingOfWisdom',
	}),
	'Blessing of Wisdom',
);

export const ManaSpringTotem = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 5675, 	minLevel: 26, maxLevel: 35 	},
			{ id: 10495, 	minLevel: 36, maxLevel: 45 	},
			{ id: 10496, 	minLevel: 46, maxLevel: 55 	},
			{ id: 10497, 	minLevel: 56 								},
		]),
		impId: ActionId.fromSpellId(16208),
		fieldName: 'manaSpringTotem',
	}),
	'Mana Spring Totem',
);

export const MeleeCritBuff = withLabel(
	makeBooleanRaidBuffInput({actionId: (player) => player.getMatchingSpellActionId([
		{ id: 17007, minLevel: 40 },
	]),
	fieldName: 'leaderOfThePack'}),
	'Leader of the Pack',
);

export const SpellCritBuff = withLabel(
	makeBooleanRaidBuffInput({actionId: (player) => player.getMatchingSpellActionId([
		{ id: 24907, minLevel: 40 },
	]),
	fieldName: 'moonkinAura'}),
	'Moonkin Aura',
);

export const SpellIncreaseBuff = withLabel(
	makeMultistateRaidBuffInput({actionId: () => ActionId.fromSpellId(425464), numStates: 200, fieldName: 'demonicPact', multiplier: 10}),
	'Demonic Pact',
);

// export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput([
// ], 'Defensive CDs');

// Misc Buffs
export const RetributionAura = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 7294, 	minLevel: 16, maxLevel: 25 	},
			{ id: 10298, 	minLevel: 26, maxLevel: 35 	},
			{ id: 10299, 	minLevel: 36, maxLevel: 45 	},
			{ id: 10300, 	minLevel: 46, maxLevel: 55 	},
			{ id: 10301, 	minLevel: 56 								},
		]),
		impId: ActionId.fromSpellId(20092),
		fieldName: 'retributionAura',
	}),
	'Retribution Aura',
);
export const Thorns = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 467, 	minLevel: 6, 	maxLevel: 13 	},
			{ id: 782, 	minLevel: 14, maxLevel: 23 	},
			{ id: 1075, minLevel: 24, maxLevel: 33 	},
			{ id: 8914, minLevel: 34, maxLevel: 43 	},
			{ id: 9756, minLevel: 44, maxLevel: 53 	},
			{ id: 9910, minLevel: 54 								},
		]),
		impId: ActionId.fromSpellId(16840),
		fieldName: 'thorns'
	}),
	'Thorns',
);
export const Innervate = withLabel(
	makeMultistateIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 29166, minLevel: 40 },
		]),
		numStates: 11,
		fieldName: 'innervates',
	}),
	'Innervate',
);
export const PowerInfusion = withLabel(
	makeMultistateIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 10060, minLevel: 40 },
		]),
		numStates: 11,
		fieldName: 'powerInfusions',
	}),
	'Power Infusion',
);

///////////////////////////////////////////////////////////////////////////
//                                 WORLD BUFFS
///////////////////////////////////////////////////////////////////////////

export const RallyingCryOfTheDragonslayer = withLabel(
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(22888), fieldName: 'rallyingCryOfTheDragonslayer'}),
	'Rallying Cry of the Dragonslayer',
);

export const SpiritOfZandalar = withLabel(
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(24425), fieldName: 'spiritOfZandalar'}),
	'Spirit of Zandalar',
);
export const SongflowerSerenade = withLabel(
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(15366), fieldName: 'songflowerSerenade'}),
	'Songflower Serenade',
);
export const WarchiefsBlessing = withLabel(
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(16609), fieldName: 'warchiefsBlessing'}),
	`Warchief's Blessing`,
);

export const SaygesDarkFortune = (inputs: ItemStatOption<SaygesFortune>[]) => makeEnumIndividualBuffInput({
	direction: IconEnumPickerDirection.Horizontal,
	values: [
		{ iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_orb_02.jpg', value: SaygesFortune.SaygesUnknown, text: `Sayge's Dark Fortune` },
		...inputs.map(input => input.config),
	],
	fieldName: 'saygesFortune',
})

export const SaygesDamage = { actionId: () => ActionId.fromSpellId(23768), value: SaygesFortune.SaygesDamage, text: `Sayge's Damage` };
export const SaygesAgility = { actionId: () => ActionId.fromSpellId(23736), value: SaygesFortune.SaygesAgility, text: `Sayge's Agility` };
export const SaygesIntellect = { actionId: () => ActionId.fromSpellId(23766), value: SaygesFortune.SaygesIntellect, text: `Sayge's Intellect` };
export const SaygesSpirit = { actionId: () => ActionId.fromSpellId(23738), value: SaygesFortune.SaygesSpirit, text: `Sayge's Spirit` };
export const SaygesStamina = { actionId: () => ActionId.fromSpellId(23737), value: SaygesFortune.SaygesStamina, text: `Sayge's Stamina` };

// Dire Maul Buffs
export const FengusFerocity = withLabel(
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(22817), fieldName: 'fengusFerocity'}),
	`Fengus' Ferocity`,
);
export const MoldarsMoxie = withLabel(
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(22818), fieldName: 'moldarsMoxie'}),
	`Moldar's Moxie`,
);
export const SlipKiksSavvy = withLabel(
	makeBooleanIndividualBuffInput({actionId: () => ActionId.fromSpellId(22820), fieldName: 'slipkiksSavvy'}),
	`Slip'kik's Savvy`,
);

// SoD World Buffs
export const BoonOfBlackfathom = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 430947, maxLevel: 39 },
		]),
		fieldName: 'boonOfBlackfathom',
	}),
	'Boon of Blackfathom',
);

export const AshenvalePvpBuff = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 430352, maxLevel: 39 },
		]),
		fieldName: 'ashenvalePvpBuff',
	}),
	'Ashenvale PvP Buff',
);

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 7386, 	minLevel: 10, maxLevel: 21 	},
			{ id: 7405, 	minLevel: 22, maxLevel: 33 	},
			{ id: 8380, 	minLevel: 34, maxLevel: 45 	},
			{ id: 11596, 	minLevel: 46, maxLevel: 57 	},
			{ id: 11597, 	minLevel: 58 								},
		]),
		fieldName: 'sunderArmor',
	}),
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 8647, 	minLevel: 14, maxLevel: 25 	},
			{ id: 8649, 	minLevel: 26, maxLevel: 35 	},
			{ id: 8650, 	minLevel: 36, maxLevel: 45 	},
			{ id: 11197, 	minLevel: 46, maxLevel: 55 	},
			{ id: 11198, 	minLevel: 56 								},
		]),
		impId: ActionId.fromSpellId(14169),
		fieldName: 'exposeArmor',
	}),
	makeMultistateMultiplierDebuffInput({
		actionId: () => ActionId.fromSpellId(402818),
		numStates: 101,
		multiplier: 10,
		fieldName: 'homunculi',
	}),
], 'Major Armor Penetration');

export const CurseOfRecklessness = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 704, 		minLevel: 14, maxLevel: 27 	},
			{ id: 7658, 	minLevel: 28, maxLevel: 41 	},
			{ id: 7659, 	minLevel: 42, maxLevel: 55 	},
			{ id: 11717, 	minLevel: 56 								},
		]),
		fieldName: 'curseOfRecklessness'
	}),
	'Curse of Recklessness',
);

export const FaerieFire = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 770, 	minLevel: 18, maxLevel: 29 	},
			{ id: 778, 	minLevel: 30, maxLevel: 41 	},
			{ id: 9749, minLevel: 42, maxLevel: 53 	},
			{ id: 9907, minLevel: 54 								},
		]),
		fieldName: 'faerieFire'
	}),
	'Faerie Fire',
);

// TODO: Classic
// export const MinorArmorDebuff = InputHelpers.makeMultiIconInput([
// 	makeTristateDebuffInput(ActionId.fromSpellId(770), ActionId.fromSpellId(33602), 'faerieFire'),
// 	makeBooleanDebuffInput({actionId: () => ActionId.fromSpellId(50511), fieldName: 'curseOfWeakness'}),
// ], 'Minor ArP');

export const AttackPowerDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 1160, 	minLevel: 14, maxLevel: 23 	},
			{ id: 6190, 	minLevel: 24, maxLevel: 33 	},
			{ id: 11554, 	minLevel: 34, maxLevel: 43 	},
			{ id: 11555, 	minLevel: 44, maxLevel: 53 	},
			{ id: 11556, 	minLevel: 54 								},
		]),
		impId: ActionId.fromSpellId(12879),
		fieldName: 'demoralizingShout'
	}),
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 99, 	minLevel: 10, maxLevel: 19 	},
			{ id: 1735, minLevel: 20, maxLevel: 31 	},
			{ id: 9490, minLevel: 32, maxLevel: 41 	},
			{ id: 9747, minLevel: 42, maxLevel: 51 	},
			{ id: 9898, minLevel: 52 								},
		]),
		impId: ActionId.fromSpellId(16862),
		fieldName: 'demoralizingRoar'
		}),
], 'Attack Power');

// TODO: SoD Mangle
export const BleedDebuff = InputHelpers.makeMultiIconInput([
	// makeBooleanDebuffInput({actionId: () => ActionId.fromSpellId(409828), fieldName: 'mangle'}),
], 'Bleed');

export const MeleeAttackSpeedDebuff = withLabel(
	makeTristateDebuffInput({
		actionId: () => ActionId.fromSpellId(6343),
		impId: ActionId.fromSpellId(12666),
		fieldName: 'thunderClap',
	}),
	'Thunder Clap',
);

export const MeleeHitDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 5570, 	minLevel: 20, maxLevel: 29 	},
			{ id: 24974, 	minLevel: 30, maxLevel: 39 	},
			{ id: 24975, 	minLevel: 40, maxLevel: 49 	},
			{ id: 24976, 	minLevel: 50, maxLevel: 59 	},
			{ id: 24977, 	minLevel: 60 								},
		]),
		fieldName: 'insectSwarm',
	}),
	'Insect Swarm',
);

// TODO: Classic
export const SpellISBDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 17803, minLevel: 10 },
		]),
		fieldName: 'improvedShadowBolt',
	}),
	'Improved Shadow Bolt',
);

export const SpellScorchDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 12873, minLevel: 40 },
		]),
		fieldName: 'improvedScorch',
	}),
	'Improved Scorch',
);

export const SpellWintersChillDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 28595, minLevel: 40 },
		]),
		fieldName: 'wintersChill',
	}),
	'Winters Chill',
);

export const CurseOfElements = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 1490, 	minLevel: 32, maxLevel: 45 	},
			{ id: 11721, 	minLevel: 46, maxLevel: 59 	},
			{ id: 11722, 	minLevel: 60 								},
		]),
		fieldName: 'curseOfElements',
	}),
	'Curse of Elements',
);

export const HuntersMark = withLabel(
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 1130, 	minLevel: 6, 	maxLevel: 21 	},
			{ id: 14323, 	minLevel: 22, maxLevel: 39 	},
			{ id: 14324, 	minLevel: 40, maxLevel: 57 	},
			{ id: 14325, 	minLevel: 58 								},
		]),
		impId: ActionId.fromSpellId(19425),
		fieldName: 'huntersMark',
	}),
	`Hunter's Mark`,
);
export const JudgementOfWisdom = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			{ id: 20186, minLevel: 38, maxLevel: 57 },
			{ id: 20355, minLevel: 58 							},
		]),
		fieldName: 'judgementOfWisdom',
	}),
	'Judgement of Wisdom',
);

// Misc Debuffs
export const JudgementOfLight = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([
		{ id: 20185, minLevel: 30, maxLevel: 39 },
		{ id: 20344, minLevel: 40, maxLevel: 49 },
		{ id: 20345, minLevel: 50, maxLevel: 59 },
		{ id: 20346, minLevel: 60 							},
	]),
	fieldName: 'judgementOfLight',
});
export const CurseOfVulnerability = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([
		{ id: 427143, minLevel: 25 },
	]),
	fieldName: 'curseOfVulnerability',
});
export const GiftOfArthas = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([
		{ id: 11374, minLevel: 38 },
	]),
	fieldName: 'giftOfArthas',
});
export const CrystalYield = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([
		{ id: 15235, minLevel: 47 },
	]),
	fieldName: 'crystalYield',
});
export const AncientCorrosivePoison = makeMultistateMultiplierDebuffInput({
	actionId: () => ActionId.fromSpellId(422996),
	numStates: 101,
	multiplier: 10,
	fieldName: 'ancientCorrosivePoison',
});

///////////////////////////////////////////////////////////////////////////
//                                 CONFIGS
///////////////////////////////////////////////////////////////////////////

export const RAID_BUFFS_CONFIG = [
	// Standard buffs
	{
		config: AllStatsBuff,
		picker: IconPicker,
		stats: []
	},
	{
		config: AllStatsPercentBuffAlliance,
		picker: MultiIconPicker,
		stats: [],
		faction: Faction.Alliance,
	},
	{
		config: AllStatsPercentBuffHorde,
		picker: IconPicker,
		stats: [],
		faction: Faction.Horde,
	},
	{
		config: ArmorBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
	{
		config: StaminaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina]
	},
	{
		config: PaladinPhysicalBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStrength, Stat.StatAgility, Stat.StatAttackPower],
	 faction: Faction.Alliance,
	},
	{
		config: StrengthBuffHorde,
		picker: IconPicker,
		stats: [Stat.StatStrength],
		faction: Faction.Horde,
	},
	{
		config: GraceOfAir,
		picker: IconPicker,
		stats: [Stat.StatAgility],
		faction: Faction.Horde,
	},
	{
		config: IntellectBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatIntellect]
	},
	{
		config: SpiritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpirit]
	},
	{
		config: BattleShoutBuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	{
		config: TrueshotAuraBuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
		},
	// {
	// 	config: AttackPowerPercentBuff,
	// 	picker: MultiIconPicker,
	// 	stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	// },
	{
		config: MeleeCritBuff,
		picker: IconPicker,
		stats: [Stat.StatMeleeCrit]
	},
	{
		config: SpellIncreaseBuff,
		picker: IconPicker,
		stats: [Stat.StatSpellPower]
	},
	{
		config: SpellCritBuff,
		picker: IconPicker,
		stats: [Stat.StatSpellCrit]
	},
	{
		config: DamageReductionPercentBuff,
		picker: IconPicker,
		stats: [Stat.StatArmor],
		faction: Faction.Alliance,
	},
	{
		config: ResistanceBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatNatureResistance, Stat.StatShadowResistance, Stat.StatFireResistance, Stat.StatFrostResistance]
	},
	// {
	// 	config: DefensiveCooldownBuff,
	// 	picker: MultiIconPicker,
	// 	stats: [Stat.StatArmor]
	// },
	{
		config: BlessingOfWisdom,
		picker: IconPicker,
		stats: [Stat.StatMP5],
		faction: Faction.Alliance,
	},
	{
		config: ManaSpringTotem,
		picker: IconPicker,
		stats: [Stat.StatMP5],
		faction: Faction.Horde,
	},

	// // Misc Buffs
	{
		config: Thorns,
		picker: IconPicker,
		stats: [Stat.StatArmor]
	},
	{
		config: RetributionAura,
		picker: IconPicker,
		stats: [Stat.StatArmor],
		faction: Faction.Alliance,
	},
	{
		config: Innervate,
		picker: IconPicker,
		stats: [Stat.StatMP5]
	},
	{
		config: PowerInfusion,
		picker: IconPicker,
		stats: [Stat.StatMP5, Stat.StatSpellPower]
	},
] as PickerStatOptions[]

export const WORLD_BUFFS_CONFIG = [
	{
		config: BoonOfBlackfathom,
		picker: IconPicker,
		stats: [
			Stat.StatMeleeCrit,
			// TODO: Stat.StatRangedCrit,
			Stat.StatSpellCrit,
			Stat.StatAttackPower
		]
	},
	{
		config: AshenvalePvpBuff,
		picker: IconPicker,
		stats: [
			Stat.StatAttackPower,
			Stat.StatSpellPower,
		]
	},
	{
		config: FengusFerocity,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	{
		config: MoldarsMoxie,
		picker: IconPicker,
		stats: [Stat.StatStamina]
	},
	{
		config: RallyingCryOfTheDragonslayer,
		picker: IconPicker,
		stats: [
			Stat.StatMeleeCrit,
			// TODO: Stat.StatRangedCrit,
			Stat.StatSpellCrit,
			Stat.StatAttackPower,
		]
	},
	{
		config: SongflowerSerenade,
		picker: IconPicker,
		stats: []
	},
	{
		config: SpiritOfZandalar,
		picker: IconPicker,
		stats: []
	},
	{
		config: WarchiefsBlessing,
		picker: IconPicker,
		stats: [
			Stat.StatHealth,
			Stat.StatMeleeHaste,
			Stat.StatMP5,
		]
	},
] as PickerStatOptions[];

export const SAYGES_CONFIG = [
	{
		config: SaygesDamage,
		stats: [],
	},
	{
		config: SaygesAgility,
		stats: [Stat.StatAgility],
	},
	{
		config: SaygesIntellect,
		stats: [Stat.StatIntellect],
	},
	{
		config: SaygesSpirit,
		stats: [Stat.StatSpirit, Stat.StatMP5],
	},
	{
		config: SaygesStamina,
		stats: [Stat.StatStamina],
	},
] as ItemStatOption<SaygesFortune>[];

export const DEBUFFS_CONFIG = [
	// // Standard Debuffs
	{ 
		config: MajorArmorDebuff,
		stats: [Stat.StatAttackPower],
		picker: MultiIconPicker,
	},
	{ 
		config: CurseOfRecklessness,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	{ 
		config: FaerieFire,
		picker: IconPicker,
		stats: [Stat.StatAttackPower]
	},
	// // { 
	// // 	config: MinorArmorDebuff,
	// // 	picker: MultiIconPicker,
	// // 	stats: [Stat.StatAttackPower]
	// // },
	{ 
		config: BleedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	},
	{ 
		config: SpellISBDebuff,
		picker: IconPicker,
		stats: [Stat.StatShadowPower]
	},
	{ 
		config: SpellScorchDebuff,
		picker: IconPicker,
		stats: [Stat.StatFirePower]
	},
	{ 
		config: SpellWintersChillDebuff,
		picker: IconPicker,
		stats: [Stat.StatFrostPower]
	},
	{
		config: CurseOfElements,
		picker: IconPicker,
		stats: [Stat.StatFirePower, Stat.StatFrostPower],
	},
	{ 
		config: AttackPowerDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
	{ 
		config: MeleeAttackSpeedDebuff,
		picker: IconPicker,
		stats: [Stat.StatArmor]
	},
	{ 
		config: MeleeHitDebuff,
		picker: IconPicker,
		stats: [Stat.StatDodge]
	},

	// Other Debuffs
	{
		config: HuntersMark,
		picker: IconPicker,
		stats: [Stat.StatRangedAttackPower],
	},
	{
		config: JudgementOfWisdom,
		picker: IconPicker,
		stats: [Stat.StatMP5, Stat.StatIntellect],
		faction: Faction.Alliance,
	},
] as PickerStatOptions[];

export const DEBUFFS_MISC_CONFIG = [
	// // Misc Debuffs
	{
		config: JudgementOfLight,
		picker: IconPicker,
		stats: [Stat.StatStamina]
	},
	{
		config: CurseOfVulnerability,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: GiftOfArthas,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: CrystalYield,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: AncientCorrosivePoison,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
] as PickerStatOptions[];