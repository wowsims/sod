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
			[1126, 1, 9],
			[5232, 10, 19],
			[6756, 20, 29],
			[5234, 30, 39],
			[8907, 40, 49],
			[9884, 50, 59],
			[9885, 60]
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
			[465, 1, 9],
			[10290, 10, 19],
			[643, 20, 29],
			[10291, 30, 39],
			[1032, 40, 49],
			[10292, 50, 59],
			[10293, 60],
		]),
		impId: ActionId.fromSpellId(20142),
		fieldName: 'devotionAura'
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			[3013, 1, 14],
			[1478, 15, 29],
			[4421, 30, 44],
			[10305, 45],
		]),
		fieldName: 'scrollOfProtection'
	}),
], 'Armor');

export const StaminaBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[1243, 1, 11],
			[1244, 12, 23],
			[1245, 24, 35],
			[2791, 36, 47],
			[10937, 48, 60],
			[10938, 60],
		]),
		impId: ActionId.fromSpellId(14767),
		fieldName: 'powerWordFortitude'
	}),
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[6307, 4, 13],
			[7804, 14, 25],
			[7805, 26, 37],
			[11766, 38, 49],
			[11767, 50]
		]),
		impId: ActionId.fromSpellId(18696),
		fieldName: 'bloodPact'
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			[1180, 5, 19],
			[1711, 20, 34],
			[4422, 35, 49],
			[10307, 50],
		]),
		fieldName: 'scrollOfStamina'
	}),
], 'Stamina');

// Separate Strength buffs allow us to use boolean pickers for each
export const PaladinPhysicalBuff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[19740, 4, 11],
			[19834, 12, 21],
			[19835, 22, 31],
			[19836, 32, 41],
			[19837, 42, 51],
			[19838, 52, 59],
			[25291, 60],
		]),
		impId: ActionId.fromSpellId(20048),
		fieldName: 'blessingOfMight',
	}),
	makeBooleanRaidBuffInput({actionId: () => ActionId.fromSpellId(425600), fieldName: 'hornOfLordaeron'}),
], 'Paladin Physical');

export const StrengthBuffHorde = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[8075, 10, 23],
			[8160, 24, 37],
			[8161, 38, 51],
			[10442, 52, 59],
			[25361, 60]
		]),
		impId: ActionId.fromSpellId(16295),
		fieldName: 'strengthOfEarthTotem'
	}),
	'Strength',
);;

export const AgilityBuffAlliance = withLabel(
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			[3012, 10, 24],
			[1477, 25, 39],
			[4425, 40, 54],
			[10309, 55],
		]),
		fieldName: 'scrollOfAgility'
	}),
	'Agility',
);

export const AgilityBuffHorde = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[8835, 42, 55],
			[10627, 56, 59],
			[25359, 60],
		]),
		impId: ActionId.fromSpellId(16295),
		fieldName: 'graceOfAirTotem',
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			[3012, 10, 24],
			[1477, 25, 39],
			[4425, 40, 54],
			[10309, 55],
		]),
		fieldName: 'scrollOfAgility'
	}),
], 'Agility');

export const IntellectBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[1459, 1, 13],
			[1460, 14, 27],
			[1461, 28, 41],
			[10156, 42, 55],
			[10157, 56],
		]),
		fieldName: 'arcaneBrilliance'
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			[955, 5, 19],
			[2290, 20, 34],
			[4419, 35, 49],
			[10308, 50],
		]),
		fieldName: 'scrollOfIntellect'
	}),
], 'Intellect');

export const SpiritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[14752, 30, 39],
			[14818, 40, 49],
			[14819, 50, 59],
			[27841, 60],
		]),
		fieldName: 'divineSpirit',
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingItemActionId([
			[1181, 1, 14],
			[1712, 15, 29],
			[4424, 30, 44],
			[10306, 45],
		]),
		fieldName: 'scrollOfSpirit',
	}),
], 'Spirit');

export const BattleShoutBuff = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[6673, 1, 11],
			[5242, 12, 21],
			[6192, 22, 31],
			[11549, 32, 41],
			[11550, 42, 51],
			[11551, 52, 59],
			[25289, 60],
		]),
		impId: ActionId.fromSpellId(12861),
		fieldName: 'battleShout'
	}),
	'Battle Shout',
);

export const TrueshotAuraBuff = withLabel(
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([[20906, 40]]),
		fieldName: 'trueshotAura'
	}),
	'Trueshot Aura',
);

// export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput([
// ], 'Attack Power %', 1, 40);

export const DamageReductionPercentBuff = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[20911, 30, 39],
			[20912, 40, 49],
			[20913, 50, 59],
			[20914, 60],
		]),
		fieldName: 'blessingOfSanctuary'
	}),
	'Blessing of Sanctuary',
);

export const ResistanceBuff = InputHelpers.makeMultiIconInput([
	// Shadow
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[976, 30, 41],
			[10957, 42, 55],
			[10958, 56],
		]),
		fieldName: 'shadowProtection'
	}),
	// Nature
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[10595, 30, 43],
			[10600, 44, 59],
			[10601, 60],
		]),
		fieldName: 'natureResistanceTotem',
		faction: Faction.Horde
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[20043, 46, 55],
			[20190, 56],
		]),
		fieldName: 'aspectOfTheWild'
	}),
	// Fire
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[19891, 36, 47],
			[19899, 48, 59],
			[19900, 60],
		]),
		fieldName: 'fireResistanceAura',
		faction: Faction.Alliance
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[8184, 28, 41],
			[10537, 42, 57],
			[10538, 58],
		]),
		fieldName: 'fireResistanceTotem',
		faction: Faction.Horde
	}),
	// Frost
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[19888, 32, 43],
			[19897, 44, 55],
			[19898, 56],
		]),
		fieldName: 'frostResistanceAura',
		faction: Faction.Alliance
	}),
	makeBooleanRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[8181, 24, 37],
			[10478, 38, 53],
			[10479, 54],
		]),
		fieldName: 'frostResistanceTotem',
		faction: Faction.Horde
	}),
], 'Resistances');

export const BlessingOfWisdom = withLabel(
	makeTristateIndividualBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[19742, 14, 23],
			[19850, 24, 33],
			[19852, 19852, 34, 43],
			[19853, 44, 53],
			[19854, 54, 59],
			[25290, 60],
		]),
		impId: ActionId.fromSpellId(20245),
		fieldName: 'blessingOfWisdom',
	}),
	'Blessing of Wisdom',
);

export const ManaSpringTotem = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[5675, 26, 35],
			[10495, 36, 45],
			[10496, 46, 55],
			[10497, 56],
		]),
		impId: ActionId.fromSpellId(16208),
		fieldName: 'manaSpringTotem',
	}),
	'Mana Spring Totem',
);

export const MeleeCritBuff = withLabel(
	makeBooleanRaidBuffInput({actionId: (player) => player.getMatchingSpellActionId([[17007, 40]]), fieldName: 'leaderOfThePack'}),
	'Leader of the Pack',
);

export const SpellCritBuff = withLabel(
	makeBooleanRaidBuffInput({actionId: (player) => player.getMatchingSpellActionId([[24907, 40]]), fieldName: 'moonkinAura'}),
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
			[7294, 16, 25],
			[10298, 26, 35],
			[10299, 36, 45],
			[10300, 46, 55],
			[10301, 56],
		]),
		impId: ActionId.fromSpellId(20092),
		fieldName: 'retributionAura',
	}),
	'Retribution Aura',
);
export const Thorns = withLabel(
	makeTristateRaidBuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[467, 6, 13],
			[782, 14, 23],
			[1075, 24, 33],
			[8914, 34, 43],
			[9756, 44, 53],
			[9910, 54],
		]),
		impId: ActionId.fromSpellId(16840),
		fieldName: 'thorns'
	}),
	'Thorns',
);
export const Innervate = withLabel(
	makeMultistateIndividualBuffInput({actionId: (player) => player.getMatchingSpellActionId([[29166, 40]]), numStates: 11, fieldName: 'innervates'}),
	'Innervate',
);
export const PowerInfusion = withLabel(
	makeMultistateIndividualBuffInput({actionId: (player) => player.getMatchingSpellActionId([[10060, 40]]), numStates: 11, fieldName: 'powerInfusions'}),
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
	makeBooleanIndividualBuffInput({actionId: (player) => player.getMatchingSpellActionId([[430947, 1, 39]]), fieldName: 'boonOfBlackfathom'}),
	'Boon of Blackfathom',
);

export const AshenvalePvpBuff = withLabel(
	makeBooleanIndividualBuffInput({actionId: (player) => player.getMatchingSpellActionId([[430352, 1, 39]]), fieldName: 'ashenvalePvpBuff'}),
	'Ashenvale Pvp Buff',
);

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[7386, 10, 21],
			[7405, 22, 33],
			[8380, 34, 45],
			[11596, 46, 57],
			[11597, 58],
		]),
		fieldName: 'sunderArmor',
	}),
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[8647, 14, 25],
			[8649, 26, 35],
			[8650, 36, 45],
			[11197, 46, 55],
			[11198, 56],
		]),
		impId: ActionId.fromSpellId(14169),
		fieldName: 'exposeArmor',
	}),
	makeMultistateMultiplierDebuffInput({actionId: () => ActionId.fromSpellId(402818), numStates: 101, multiplier: 10, fieldName: 'homunculi'}),
], 'Major Armor Penetration');

export const CurseOfRecklessness = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[704, 14, 27],
			[7658, 28, 41],
			[7659, 42, 55],
			[11717, 56],
		]),
		fieldName: 'curseOfRecklessness'
	}),
	'Curse of Recklessness',
);

export const FaerieFire = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[770, 18, 29],
			[778, 30, 41],
			[9749, 42, 53],
			[9907, 54],
		]),
		fieldName: 'faerieFire'
	}),
	'Faerie Fire',
);

// TODO: Classic
// export const MinorArmorDebuff = InputHelpers.makeMultiIconInput([
// 	makeTristateDebuffInput(ActionId.fromSpellId(770), ActionId.fromSpellId(33602), 'faerieFire'),
// 	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(50511), fieldName: 'curseOfWeakness'}),
// ], 'Minor ArP');

export const AttackPowerDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[1160, 14, 23],
			[6190, 24, 33],
			[11554, 34, 43],
			[11555, 44, 53],
			[11556, 54],
		]),
		impId: ActionId.fromSpellId(12879),
		fieldName: 'demoralizingShout'
	}),
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[99, 10, 19],
			[1735, 20, 31],
			[9490, 32, 41],
			[9747, 42, 51],
			[9898, 52],
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
	makeTristateDebuffInput({actionId: () => ActionId.fromSpellId(6343), impId: ActionId.fromSpellId(12666), fieldName: 'thunderClap'}),
	'Thunder Clap',
);

export const MeleeHitDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[5570, 20, 29],
			[24974, 30, 39],
			[24975, 40, 49],
			[24976, 50, 59],
			[24977, 60],
		]),
		fieldName: 'insectSwarm',
	}),
	'Insect Swarm',
);

// TODO: Classic
export const SpellISBDebuff = withLabel(
	makeBooleanDebuffInput({actionId: (player) => player.getMatchingSpellActionId([[17803, 10]]), fieldName: 'improvedShadowBolt'}),
	'Improved Shadow Bolt',
);

export const SpellScorchDebuff = withLabel(
	makeBooleanDebuffInput({actionId: (player) => player.getMatchingSpellActionId([[12873, 40]]), fieldName: 'improvedScorch'}),
	'Improved Scorch',
);

export const SpellWintersChillDebuff = withLabel(
	makeBooleanDebuffInput({actionId: (player) => player.getMatchingSpellActionId([[28595, 40]]), fieldName: 'wintersChill'}),
	'Winters Chill',
);

export const CurseOfElements = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[1490, 32, 45],
			[11721, 46, 59],
			[11722, 60],
		]),
		fieldName: 'curseOfElements',
	}),
	'Curse of Elements',
);

export const HuntersMark = withLabel(
	makeTristateDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[1130, 6, 21],
			[14323, 22, 39],
			[14324, 40, 57],
			[14325, 58],
		]),
		impId: ActionId.fromSpellId(19425),
		fieldName: 'huntersMark',
	}),
	`Hunter's Mark`,
);
export const JudgementOfWisdom = withLabel(
	makeBooleanDebuffInput({
		actionId: (player) => player.getMatchingSpellActionId([
			[20186, 38, 57],
			[20355, 58],
		]),
		fieldName: 'judgementOfWisdom',
	}),
	'Judgement of Wisdom',
);

// Misc Debuffs
export const JudgementOfLight = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([
		[20185, 30, 39],
		[20344, 40, 49],
		[20345, 50, 59],
		[20346, 60],
	]),
	fieldName: 'judgementOfLight',
});
export const CurseOfVulnerability = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([[427143, 25]]),
	fieldName: 'curseOfVulnerability',
});
export const GiftOfArthas = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([[11374, 38]]),
	fieldName: 'giftOfArthas',
});
export const CrystalYield = makeBooleanDebuffInput({
	actionId: (player) => player.getMatchingSpellActionId([[15235, 47]]),
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
		config: AgilityBuffAlliance,
		picker: IconPicker,
		stats: [Stat.StatAgility],
		faction: Faction.Alliance,
	},
	{
		config: AgilityBuffHorde,
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