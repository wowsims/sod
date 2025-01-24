import { Faction, SaygesFortune, Stat } from '../../proto/common';
import { ActionId } from '../../proto_utils/action_id';
import {
	makeBooleanDebuffInput,
	makeBooleanIndividualBuffInput,
	makeBooleanRaidBuffInput,
	makeEnumIndividualBuffInput,
	makeMultistateIndividualBuffInput,
	makeMultistateMultiplierDebuffInput,
	makeMultistateRaidBuffInput,
	makeTristateDebuffInput,
	makeTristateIndividualBuffInput,
	makeTristateRaidBuffInput,
	withLabel,
} from '../icon_inputs';
import { IconPicker, IconPickerDirection } from '../icon_picker';
import * as InputHelpers from '../input_helpers';
import { MultiIconPicker } from '../multi_icon_picker';
import { ItemStatOption, PickerStatOptions } from './stat_options';

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

export const AllStatsBuff = withLabel(
	makeTristateRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 1126, maxLevel: 9 },
				{ id: 5232, minLevel: 10, maxLevel: 19 },
				{ id: 6756, minLevel: 20, maxLevel: 29 },
				{ id: 5234, minLevel: 30, maxLevel: 39 },
				{ id: 8907, minLevel: 40, maxLevel: 49 },
				{ id: 9884, minLevel: 50, maxLevel: 59 },
				{ id: 9885, minLevel: 60 },
			]),
		impId: ActionId.fromSpellId(17055),
		fieldName: 'giftOfTheWild',
	}),
	'Mark of the Wild',
);

// Separate Strength buffs allow us to use a boolean pickers for Horde specifically
export const AllStatsPercentBuff = InputHelpers.makeMultiIconInput({
	values: [
		makeBooleanIndividualBuffInput({
			actionId: () => ActionId.fromSpellId(20217),
			fieldName: 'blessingOfKings',
			showWhen: player => player.getFaction() === Faction.Alliance,
		}),
		makeBooleanRaidBuffInput({
			actionId: () => ActionId.fromSpellId(409580),
			fieldName: 'aspectOfTheLion',
		}),
	],
	label: 'Stats %',
});

export const ArmorBuff = withLabel(
	makeTristateRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 465, maxLevel: 9 },
				{ id: 10290, minLevel: 10, maxLevel: 19 },
				{ id: 643, minLevel: 20, maxLevel: 29 },
				{ id: 10291, minLevel: 30, maxLevel: 39 },
				{ id: 1032, minLevel: 40, maxLevel: 49 },
				{ id: 10292, minLevel: 50, maxLevel: 59 },
				{ id: 10293, minLevel: 60 },
			]),
		impId: ActionId.fromSpellId(20142),
		showWhen: player => player.getFaction() === Faction.Alliance,
		fieldName: 'devotionAura',
	}),
	'Armor',
);

export const PhysDamReductionBuff = withLabel(
	makeTristateRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 8071, minLevel: 4, maxLevel: 13 },
				{ id: 8154, minLevel: 14, maxLevel: 23 },
				{ id: 8155, minLevel: 24, maxLevel: 33 },
				{ id: 10406, minLevel: 34, maxLevel: 43 },
				{ id: 10407, minLevel: 44, maxLevel: 53 },
				{ id: 10408, minLevel: 54 },
			]),
		impId: ActionId.fromSpellId(16293),
		showWhen: player => player.getFaction() === Faction.Horde,
		fieldName: 'stoneskinTotem',
	}),
	'Stoneskin',
);

//export const DamageReductionPercentBuff = withLabel(
//	makeBooleanIndividualBuffInput({
//		actionId: player =>
//			player.getMatchingSpellActionId([
//				{ id: 20911, minLevel: 30, maxLevel: 39 },
//				{ id: 20912, minLevel: 40, maxLevel: 49 },
//				{ id: 20913, minLevel: 50, maxLevel: 59 },
//				{ id: 20914, minLevel: 60 },
//			]),
//		showWhen: player => player.getFaction() === Faction.Alliance,
//		fieldName: 'blessingOfSanctuary',
//	}),
//	'Blessing of Sanctuary',
//);

export const ResistanceBuff = InputHelpers.makeMultiIconInput({
	values: [
		// Shadow
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 976, minLevel: 30, maxLevel: 41 },
					{ id: 10957, minLevel: 42, maxLevel: 55 },
					{ id: 10958, minLevel: 56 },
				]),
			fieldName: 'shadowProtection',
		}),
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 19876, minLevel: 28, maxLevel: 39 },
					{ id: 19895, minLevel: 40, maxLevel: 51 },
					{ id: 19896, minLevel: 52 },
				]),
			fieldName: 'shadowResistanceAura',
		}),
		// Nature
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 10595, minLevel: 30, maxLevel: 43 },
					{ id: 10600, minLevel: 44, maxLevel: 59 },
					{ id: 10601, minLevel: 60 },
				]),
			fieldName: 'natureResistanceTotem',
			showWhen: player => player.getFaction() === Faction.Horde,
		}),
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 20043, minLevel: 46, maxLevel: 55 },
					{ id: 20190, minLevel: 56 },
				]),
			fieldName: 'aspectOfTheWild',
		}),
		// Fire
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 19891, minLevel: 36, maxLevel: 47 },
					{ id: 19899, minLevel: 48, maxLevel: 59 },
					{ id: 19900, minLevel: 60 },
				]),
			fieldName: 'fireResistanceAura',
			showWhen: player => player.getFaction() === Faction.Alliance,
		}),
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 8184, minLevel: 28, maxLevel: 41 },
					{ id: 10537, minLevel: 42, maxLevel: 57 },
					{ id: 10538, minLevel: 58 },
				]),
			fieldName: 'fireResistanceTotem',
			showWhen: player => player.getFaction() === Faction.Horde,
		}),
		// Frost
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 19888, minLevel: 32, maxLevel: 43 },
					{ id: 19897, minLevel: 44, maxLevel: 55 },
					{ id: 19898, minLevel: 56 },
				]),
			fieldName: 'frostResistanceAura',
			showWhen: player => player.getFaction() === Faction.Alliance,
		}),
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 8181, minLevel: 24, maxLevel: 37 },
					{ id: 10478, minLevel: 38, maxLevel: 53 },
					{ id: 10479, minLevel: 54 },
				]),
			fieldName: 'frostResistanceTotem',
			showWhen: player => player.getFaction() === Faction.Horde,
		}),
	],
	label: 'Resistances',
});

export const StaminaBuff = InputHelpers.makeMultiIconInput({
	values: [
		makeTristateRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 1243, maxLevel: 11 },
					{ id: 1244, minLevel: 12, maxLevel: 23 },
					{ id: 1245, minLevel: 24, maxLevel: 35 },
					{ id: 2791, minLevel: 36, maxLevel: 47 },
					{ id: 10937, minLevel: 48, maxLevel: 59 },
					{ id: 10938, minLevel: 60 },
				]),
			impId: ActionId.fromSpellId(14767),
			fieldName: 'powerWordFortitude',
		}),
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingItemActionId([
					{ id: 1180, minLevel: 5, maxLevel: 19 },
					{ id: 1711, minLevel: 20, maxLevel: 34 },
					{ id: 4422, minLevel: 35, maxLevel: 49 },
					{ id: 10307, minLevel: 50 },
				]),
			fieldName: 'scrollOfStamina',
		}),
	],
	label: 'Stamina',
});

export const BloodPactBuff = InputHelpers.makeMultiIconInput({
	values: [
		makeTristateRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 6307, minLevel: 4, maxLevel: 13 },
					{ id: 7804, minLevel: 14, maxLevel: 25 },
					{ id: 7805, minLevel: 26, maxLevel: 37 },
					{ id: 11766, minLevel: 38, maxLevel: 49 },
					{ id: 11767, minLevel: 50 },
				]),
			impId: ActionId.fromSpellId(18696),
			fieldName: 'bloodPact',
		}),
		makeBooleanRaidBuffInput({
			actionId: () => ActionId.fromSpellId(403215),
			fieldName: 'commandingShout',
		}),
	],
	label: 'Blood Pact',
});

export const PaladinPhysicalBuff = InputHelpers.makeMultiIconInput({
	values: [
		makeBooleanRaidBuffInput({
			actionId: () => ActionId.fromSpellId(425600),
			fieldName: 'hornOfLordaeron',
			showWhen: player => player.getFaction() == Faction.Alliance,
		}),
		makeTristateIndividualBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 19740, minLevel: 4, maxLevel: 11 },
					{ id: 19834, minLevel: 12, maxLevel: 21 },
					{ id: 19835, minLevel: 22, maxLevel: 31 },
					{ id: 19836, minLevel: 32, maxLevel: 41 },
					{ id: 19837, minLevel: 42, maxLevel: 51 },
					{ id: 19838, minLevel: 52, maxLevel: 59 },
					{ id: 25291, minLevel: 60 },
				]),
			impId: ActionId.fromSpellId(20048),
			fieldName: 'blessingOfMight',
			showWhen: player => player.getFaction() === Faction.Alliance,
		}),
	],
	label: 'Paladin Physical',
});

export const StrengthBuffHorde = withLabel(
	makeTristateRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 8075, minLevel: 10, maxLevel: 23 },
				{ id: 8160, minLevel: 24, maxLevel: 37 },
				{ id: 8161, minLevel: 38, maxLevel: 51 },
				{ id: 10442, minLevel: 52, maxLevel: 59 },
				{ id: 25361, minLevel: 60 },
			]),
		impId: ActionId.fromSpellId(16295),
		fieldName: 'strengthOfEarthTotem',
		showWhen: player => player.getFaction() === Faction.Horde,
	}),
	'Strength',
);

export const GraceOfAir = withLabel(
	makeTristateRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 8835, minLevel: 42, maxLevel: 55 },
				{ id: 10627, minLevel: 56, maxLevel: 59 },
				{ id: 25359, minLevel: 60 },
			]),
		impId: ActionId.fromSpellId(16295),
		fieldName: 'graceOfAirTotem',
		showWhen: player => player.getFaction() === Faction.Horde,
	}),
	'Agility',
);

export const IntellectBuff = InputHelpers.makeMultiIconInput({
	values: [
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 1459, maxLevel: 13 },
					{ id: 1460, minLevel: 14, maxLevel: 27 },
					{ id: 1461, minLevel: 28, maxLevel: 41 },
					{ id: 10156, minLevel: 42, maxLevel: 55 },
					{ id: 10157, minLevel: 56 },
				]),
			fieldName: 'arcaneBrilliance',
		}),
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingItemActionId([
					{ id: 955, minLevel: 5, maxLevel: 19 },
					{ id: 2290, minLevel: 20, maxLevel: 34 },
					{ id: 4419, minLevel: 35, maxLevel: 49 },
					{ id: 10308, minLevel: 50 },
				]),
			fieldName: 'scrollOfIntellect',
		}),
	],
	label: 'Intellect',
});

export const SpiritBuff = InputHelpers.makeMultiIconInput({
	values: [
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 14752, minLevel: 30, maxLevel: 39 },
					{ id: 14818, minLevel: 40, maxLevel: 49 },
					{ id: 14819, minLevel: 50, maxLevel: 59 },
					{ id: 27841, minLevel: 60 },
				]),
			fieldName: 'divineSpirit',
		}),
		makeBooleanRaidBuffInput({
			actionId: player =>
				player.getMatchingItemActionId([
					{ id: 1181, maxLevel: 14 },
					{ id: 1712, minLevel: 15, maxLevel: 29 },
					{ id: 4424, minLevel: 30, maxLevel: 44 },
					{ id: 10306, minLevel: 45 },
				]),
			fieldName: 'scrollOfSpirit',
		}),
	],
	label: 'Spirit',
});

export const BattleShoutBuff = withLabel(
	makeTristateRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 6673, maxLevel: 11 },
				{ id: 5242, minLevel: 12, maxLevel: 21 },
				{ id: 6192, minLevel: 22, maxLevel: 31 },
				{ id: 11549, minLevel: 32, maxLevel: 41 },
				{ id: 11550, minLevel: 42, maxLevel: 51 },
				{ id: 11551, minLevel: 52, maxLevel: 59 },
				{ id: 25289, minLevel: 60 },
			]),
		impId: ActionId.fromSpellId(12861),
		fieldName: 'battleShout',
	}),
	'Battle Shout',
);

export const TrueshotAuraBuff = withLabel(
	makeBooleanRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 19506, minLevel: 40, maxLevel: 49 },
				{ id: 20905, minLevel: 50, maxLevel: 59 },
				{ id: 20906, minLevel: 60 },
			]),
		fieldName: 'trueshotAura',
	}),
	'Trueshot Aura',
);

export const BlessingOfWisdom = withLabel(
	makeTristateIndividualBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 19742, minLevel: 14, maxLevel: 23 },
				{ id: 19850, minLevel: 24, maxLevel: 33 },
				{ id: 19852, minLevel: 34, maxLevel: 43 },
				{ id: 19853, minLevel: 44, maxLevel: 53 },
				{ id: 19854, minLevel: 54, maxLevel: 59 },
				{ id: 25290, minLevel: 60 },
			]),
		impId: ActionId.fromSpellId(20245),
		fieldName: 'blessingOfWisdom',
		showWhen: player => player.getFaction() === Faction.Alliance,
	}),
	'Blessing of Wisdom',
);
export const ManaSpringTotem = withLabel(
	makeTristateRaidBuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 5675, minLevel: 26, maxLevel: 35 },
				{ id: 10495, minLevel: 36, maxLevel: 45 },
				{ id: 10496, minLevel: 46, maxLevel: 55 },
				{ id: 10497, minLevel: 56 },
			]),
		impId: ActionId.fromSpellId(16208),
		fieldName: 'manaSpringTotem',
		showWhen: player => player.getFaction() === Faction.Horde,
	}),
	'Mana Spring Totem',
);
export const VampiricTouchReplenishment = withLabel(
	makeMultistateRaidBuffInput({ actionId: () => ActionId.fromSpellId(402668), numStates: 21, fieldName: 'vampiricTouch', multiplier: 20 }),
	'Vampiric Touch MP5',
);

export const MeleeCritBuff = withLabel(
	makeBooleanRaidBuffInput({ actionId: player => player.getMatchingSpellActionId([{ id: 24932, minLevel: 40 }]), fieldName: 'leaderOfThePack' }),
	'Leader of the Pack',
);

export const HordeThreatBuff = withLabel(
	makeBooleanRaidBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 408696, minLevel: 40 }]),
		fieldName: 'spiritOfTheAlpha',
		showWhen: player => player.getFaction() === Faction.Horde,
	}),
	'Spirit of The Alpha',
);

export const SpellCritBuff = withLabel(
	makeBooleanRaidBuffInput({ actionId: player => player.getMatchingSpellActionId([{ id: 24907, minLevel: 40 }]), fieldName: 'moonkinAura' }),
	'Moonkin Aura',
);

export const SpellIncreaseBuff = withLabel(
	makeMultistateRaidBuffInput({ actionId: () => ActionId.fromSpellId(425464), numStates: 21, fieldName: 'demonicPact', multiplier: 10 }),
	'Demonic Pact',
);

// Misc Buffs
export const AtieshCastSpeedBuff = makeBooleanRaidBuffInput({
	actionId: () => ActionId.fromSpellId(1219557),
	fieldName: 'atieshCastSpeedBuff',
});
export const AtieshHealingBuff = makeBooleanRaidBuffInput({
	actionId: () => ActionId.fromSpellId(1219553),
	fieldName: 'atieshHealingBuff',
});
export const AtieshSpellCritBuff = makeBooleanRaidBuffInput({
	actionId: () => ActionId.fromSpellId(1219558),
	fieldName: 'atieshSpellCritBuff',
});
export const AtieshSpellPowerBuff = makeBooleanRaidBuffInput({
	actionId: () => ActionId.fromSpellId(1219552),
	fieldName: 'atieshSpellPowerBuff',
});

export const ImprovedStoneskinWindwall = makeBooleanRaidBuffInput({
	actionId: () => ActionId.fromSpellId(457544),
	fieldName: 'improvedStoneskinWindwall',
});

export const RetributionAura = makeTristateRaidBuffInput({
	actionId: player =>
		player.getMatchingSpellActionId([
			{ id: 7294, minLevel: 16, maxLevel: 25 },
			{ id: 10298, minLevel: 26, maxLevel: 35 },
			{ id: 10299, minLevel: 36, maxLevel: 45 },
			{ id: 10300, minLevel: 46, maxLevel: 55 },
			{ id: 10301, minLevel: 56 },
		]),
	impId: ActionId.fromSpellId(20092),
	fieldName: 'retributionAura',
	showWhen: player => player.getFaction() === Faction.Alliance,
});

export const SanctityAura = makeBooleanRaidBuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 20218, minLevel: 30 }]),
	fieldName: 'sanctityAura',
	showWhen: player => player.getFaction() === Faction.Alliance,
});

export const Thorns = makeTristateRaidBuffInput({
	actionId: player =>
		player.getMatchingSpellActionId([
			{ id: 467, minLevel: 6, maxLevel: 13 },
			{ id: 782, minLevel: 14, maxLevel: 23 },
			{ id: 1075, minLevel: 24, maxLevel: 33 },
			{ id: 8914, minLevel: 34, maxLevel: 43 },
			{ id: 9756, minLevel: 44, maxLevel: 53 },
			{ id: 9910, minLevel: 54 },
		]),
	impId: ActionId.fromSpellId(16840),
	fieldName: 'thorns',
});

export const Innervate = makeMultistateIndividualBuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 29166, minLevel: 40 }]),
	numStates: 11,
	fieldName: 'innervates',
});

export const PowerInfusion = makeMultistateIndividualBuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 10060, minLevel: 40 }]),
	numStates: 11,
	fieldName: 'powerInfusions',
});

export const BattleSquawkBuff = makeMultistateRaidBuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 23060, minLevel: 40 }]),
	numStates: 6,
	fieldName: 'battleSquawk',
});

///////////////////////////////////////////////////////////////////////////
//                                 WORLD BUFFS
///////////////////////////////////////////////////////////////////////////

export const RallyingCryOfTheDragonslayer = makeBooleanIndividualBuffInput({
	actionId: () => ActionId.fromSpellId(22888),
	fieldName: 'rallyingCryOfTheDragonslayer',
});
export const ValorOfAzeroth = makeBooleanIndividualBuffInput({
	actionId: () => ActionId.fromSpellId(461475),
	fieldName: 'valorOfAzeroth',
});
export const DragonslayerBuffInput = InputHelpers.makeMultiIconInput({ values: [RallyingCryOfTheDragonslayer, ValorOfAzeroth], label: 'Dragonslayer Buff' });

export const SpiritOfZandalar = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: () => ActionId.fromSpellId(24425),
		fieldName: 'spiritOfZandalar',
	}),
	'Spirit of Zandalar',
);
export const SongflowerSerenade = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: () => ActionId.fromSpellId(15366),
		fieldName: 'songflowerSerenade',
	}),
	'Songflower Serenade',
);
export const WarchiefsBlessing = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: () => ActionId.fromSpellId(16609),
		fieldName: 'warchiefsBlessing',
		showWhen: player => player.getFaction() === Faction.Horde,
	}),
	`Warchief's Blessing`,
);
export const MightOfStormwind = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: () => ActionId.fromSpellId(460940),
		fieldName: 'mightOfStormwind',
		showWhen: player => player.getFaction() === Faction.Alliance,
	}),
	`Might Of Stormwind`,
);

export const SaygesDarkFortune = (inputs: ItemStatOption<SaygesFortune>[]) =>
	makeEnumIndividualBuffInput({
		direction: IconPickerDirection.Horizontal,
		values: [
			{ iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_orb_02.jpg', value: SaygesFortune.SaygesUnknown, text: `Sayge's Dark Fortune` },
			...inputs.map(input => input.config),
		],
		fieldName: 'saygesFortune',
	});

export const SaygesDamage = { actionId: () => ActionId.fromSpellId(23768), value: SaygesFortune.SaygesDamage, text: `Sayge's Damage` };
export const SaygesAgility = { actionId: () => ActionId.fromSpellId(23736), value: SaygesFortune.SaygesAgility, text: `Sayge's Agility` };
export const SaygesIntellect = { actionId: () => ActionId.fromSpellId(23766), value: SaygesFortune.SaygesIntellect, text: `Sayge's Intellect` };
export const SaygesSpirit = { actionId: () => ActionId.fromSpellId(23738), value: SaygesFortune.SaygesSpirit, text: `Sayge's Spirit` };
export const SaygesStamina = { actionId: () => ActionId.fromSpellId(23737), value: SaygesFortune.SaygesStamina, text: `Sayge's Stamina` };

// Dire Maul Buffs
export const FengusFerocity = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 22817, minLevel: 51 }]),
		fieldName: 'fengusFerocity',
	}),
	`Fengus' Ferocity`,
);
export const MoldarsMoxie = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 22818, minLevel: 51 }]),
		fieldName: 'moldarsMoxie',
	}),
	`Moldar's Moxie`,
);
export const SlipKiksSavvy = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 22820, minLevel: 51 }]),
		fieldName: 'slipkiksSavvy',
	}),
	`Slip'kik's Savvy`,
);

// SoD World Buffs
export const BoonOfBlackfathom = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 430947, maxLevel: 39 }]),
		fieldName: 'boonOfBlackfathom',
	}),
	'Boon of Blackfathom',
);

export const AshenvalePvpBuff = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 430352, maxLevel: 39 }]),
		fieldName: 'ashenvalePvpBuff',
	}),
	'Ashenvale PvP Buff',
);

export const SparkOfInspiration = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 438536, maxLevel: 49 }]),
		fieldName: 'sparkOfInspiration',
	}),
	'Spark of Inspiration',
);

export const FervorOfTheTempleExplorer = withLabel(
	makeBooleanIndividualBuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 446695, maxLevel: 59 }]),
		fieldName: 'fervorOfTheTempleExplorer',
	}),
	'Fervor Of The Temple Explorer',
);

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput({
	values: [
		makeBooleanDebuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 7386, minLevel: 10, maxLevel: 21 },
					{ id: 7405, minLevel: 22, maxLevel: 33 },
					{ id: 8380, minLevel: 34, maxLevel: 45 },
					{ id: 11596, minLevel: 46, maxLevel: 57 },
					{ id: 11597, minLevel: 58 },
				]),
			fieldName: 'sunderArmor',
		}),
		makeTristateDebuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 8647, minLevel: 14, maxLevel: 25 },
					{ id: 8649, minLevel: 26, maxLevel: 35 },
					{ id: 8650, minLevel: 36, maxLevel: 45 },
					{ id: 11197, minLevel: 46, maxLevel: 55 },
					{ id: 11198, minLevel: 56 },
				]),
			impId: ActionId.fromSpellId(14169),
			fieldName: 'exposeArmor',
		}),
		makeTristateDebuffInput({
			actionId: player => player.getMatchingSpellActionId([{ id: 439500, minLevel: 60 }]),
			impId: ActionId.fromSpellId(14169),
			fieldName: 'sebaciousPoison',
		}),
		makeMultistateMultiplierDebuffInput({
			actionId: () => ActionId.fromSpellId(402818),
			numStates: 11,
			multiplier: 10,
			reverse: true,
			fieldName: 'homunculi',
		}),
	],
	label: 'Major Armor Penetration',
});

export const CurseOfRecklessness = withLabel(
	makeBooleanDebuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 704, minLevel: 14, maxLevel: 27 },
				{ id: 7658, minLevel: 28, maxLevel: 41 },
				{ id: 7659, minLevel: 42, maxLevel: 55 },
				{ id: 11717, minLevel: 56 },
			]),
		fieldName: 'curseOfRecklessness',
	}),
	'Curse of Recklessness',
);

export const FaerieFire = withLabel(
	makeBooleanDebuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 770, minLevel: 18, maxLevel: 29 },
				{ id: 778, minLevel: 30, maxLevel: 41 },
				{ id: 9749, minLevel: 42, maxLevel: 53 },
				{ id: 9907, minLevel: 54 },
			]),
		fieldName: 'faerieFire',
	}),
	'Faerie Fire',
);

export const curseOfWeaknessDebuff = withLabel(
	makeTristateDebuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 702, minLevel: 4, maxLevel: 11 },
				{ id: 1108, minLevel: 12, maxLevel: 21 },
				{ id: 6205, minLevel: 22, maxLevel: 31 },
				{ id: 7646, minLevel: 32, maxLevel: 41 },
				{ id: 11707, minLevel: 42, maxLevel: 51 },
				{ id: 11708, minLevel: 52 },
			]),
		impId: ActionId.fromSpellId(18181),
		fieldName: 'curseOfWeakness',
	}),
	'Curse of Weakness',
);

export const AttackPowerDebuff = InputHelpers.makeMultiIconInput({
	values: [
		makeTristateDebuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 1160, minLevel: 14, maxLevel: 23 },
					{ id: 6190, minLevel: 24, maxLevel: 33 },
					{ id: 11554, minLevel: 34, maxLevel: 43 },
					{ id: 11555, minLevel: 44, maxLevel: 53 },
					{ id: 11556, minLevel: 54 },
				]),
			impId: ActionId.fromSpellId(12879),
			fieldName: 'demoralizingShout',
		}),
		makeTristateDebuffInput({
			actionId: player =>
				player.getMatchingSpellActionId([
					{ id: 99, minLevel: 10, maxLevel: 19 },
					{ id: 1735, minLevel: 20, maxLevel: 31 },
					{ id: 9490, minLevel: 32, maxLevel: 41 },
					{ id: 9747, minLevel: 42, maxLevel: 51 },
					{ id: 9898, minLevel: 52 },
				]),
			impId: ActionId.fromSpellId(16862),
			fieldName: 'demoralizingRoar',
		}),
		makeMultistateMultiplierDebuffInput({
			actionId: () => ActionId.fromSpellId(402811),
			numStates: 11,
			multiplier: 10,
			reverse: true,
			fieldName: 'homunculi',
		}),
		makeBooleanDebuffInput({
			actionId: () => ActionId.fromSpellId(439473),
			fieldName: 'atrophicPoison',
		}),
	],
	label: 'Attack Power',
});

// TODO: SoD Mangle
export const BleedDebuff = withLabel(makeBooleanDebuffInput({ actionId: () => ActionId.fromSpellId(409828), fieldName: 'mangle' }), 'Bleed');

export const MeleeAttackSpeedDebuff = InputHelpers.makeMultiIconInput({
	values: [
		makeTristateDebuffInput({
			actionId: () => ActionId.fromSpellId(6343),
			impId: ActionId.fromSpellId(403219),
			fieldName: 'thunderClap',
		}),
		makeMultistateMultiplierDebuffInput({
			actionId: () => ActionId.fromSpellId(402808),
			numStates: 11,
			multiplier: 10,
			reverse: true,
			fieldName: 'homunculi',
		}),
		makeBooleanDebuffInput({
			actionId: () => ActionId.fromSpellId(408699),
			fieldName: 'waylay',
		}),
		makeBooleanDebuffInput({
			actionId: () => ActionId.fromSpellId(21992),
			fieldName: 'thunderfury',
		}),
		makeBooleanDebuffInput({
			actionId: () => ActionId.fromSpellId(439472),
			fieldName: 'numbingPoison',
		}),
	],
	label: 'Attack Speed',
});

export const MeleeHitDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 5570, minLevel: 20, maxLevel: 29 },
				{ id: 24974, minLevel: 30, maxLevel: 39 },
				{ id: 24975, minLevel: 40, maxLevel: 49 },
				{ id: 24976, minLevel: 50, maxLevel: 59 },
				{ id: 24977, minLevel: 60 },
			]),
		fieldName: 'insectSwarm',
	}),
	'Insect Swarm',
);

export const SpellISBDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 17803, minLevel: 10 }]),
		fieldName: 'improvedShadowBolt',
	}),
	'Improved Shadow Bolt',
);

export const SpellScorchDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 12873, minLevel: 40 }]),
		fieldName: 'improvedScorch',
	}),
	'Fire Damage',
);

export const SpellWintersChillDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 28595, minLevel: 40 }]),
		fieldName: 'wintersChill',
	}),
	'Frost Damage',
);

export const NatureSpellDamageDebuff = InputHelpers.makeMultiIconInput({
	values: [
		makeBooleanDebuffInput({
			actionId: player => player.getMatchingSpellActionId([{ id: 17364, minLevel: 40 }]),
			fieldName: 'stormstrike',
		}),
		makeBooleanDebuffInput({
			actionId: () => ActionId.fromSpellId(408258),
			fieldName: 'dreamstate',
		}),
	],
	label: 'Nature Damage',
});

export const SpellShadowWeavingDebuff = withLabel(
	makeBooleanDebuffInput({
		actionId: player => player.getMatchingSpellActionId([{ id: 15334, minLevel: 40 }]),
		fieldName: 'shadowWeaving',
	}),
	'Shadow Weaving',
);

export const MarkOfChaos = makeBooleanDebuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 461615 }]),
	fieldName: 'markOfChaos',
});

export const CurseOfElements = makeBooleanDebuffInput({
	actionId: player =>
		player.getMatchingSpellActionId([
			{ id: 1490, minLevel: 32, maxLevel: 45 },
			{ id: 11721, minLevel: 46, maxLevel: 59 },
			{ id: 11722, minLevel: 60 },
		]),
	fieldName: 'curseOfElements',
});

export const CurseOfShadow = makeBooleanDebuffInput({
	actionId: player =>
		player.getMatchingSpellActionId([
			{ id: 17862, minLevel: 44, maxLevel: 59 },
			{ id: 17937, minLevel: 60 },
		]),
	fieldName: 'curseOfShadow',
});

export const WarlockCursesConfig = InputHelpers.makeMultiIconInput({ values: [MarkOfChaos, CurseOfElements, CurseOfShadow], label: 'Warlock Curses' });

export const OccultPoison = withLabel(
	makeBooleanDebuffInput({
		actionId: player => player.getMatchingItemActionId([{ id: 226374, minLevel: 54 }]),
		fieldName: 'occultPoison',
	}),
	'Occult Poison',
);

export const HuntersMark = withLabel(
	makeTristateDebuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 1130, minLevel: 6, maxLevel: 21 },
				{ id: 14323, minLevel: 22, maxLevel: 39 },
				{ id: 14324, minLevel: 40, maxLevel: 57 },
				{ id: 14325, minLevel: 58 },
			]),
		impId: ActionId.fromSpellId(19425),
		fieldName: 'huntersMark',
	}),
	`Hunter's Mark`,
);
export const JudgementOfWisdom = withLabel(
	makeBooleanDebuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 20186, minLevel: 38, maxLevel: 57 },
				{ id: 20355, minLevel: 58 },
			]),
		fieldName: 'judgementOfWisdom',
		showWhen: player => player.getFaction() === Faction.Alliance,
	}),
	'Judgement of Wisdom',
);
export const JudgementOfTheCrusader = withLabel(
	makeTristateDebuffInput({
		actionId: player =>
			player.getMatchingSpellActionId([
				{ id: 20300, minLevel: 22, maxLevel: 31 },
				{ id: 20301, minLevel: 32, maxLevel: 41 },
				{ id: 20302, minLevel: 42, maxLevel: 51 },
				{ id: 20303, minLevel: 52 },
			]),
		impId: ActionId.fromSpellId(20337),
		fieldName: 'judgementOfTheCrusader',
		showWhen: player => player.getFaction() === Faction.Alliance,
	}),
	'Judgement of the Crusader',
);

// Misc Debuffs
export const ImprovedFaerieFire = makeBooleanDebuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 455864, minLevel: 60 }]),
	fieldName: 'improvedFaerieFire',
});
export const MeleeHunter2pcT1Bonus = makeBooleanDebuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 456389, minLevel: 60 }]),
	fieldName: 'meleeHunterDodgeDebuff',
});
export const MekkatorqueFistDebuff = makeBooleanDebuffInput({
	actionId: player => player.getMatchingItemActionId([{ id: 213409, minLevel: 40, maxLevel: 45 }]),
	fieldName: 'mekkatorqueFistDebuff',
});
export const SerpentsStrikerFistDebuff = makeBooleanDebuffInput({
	actionId: player => player.getMatchingItemActionId([{ id: 220589, minLevel: 50, maxLevel: 55 }]),
	fieldName: 'serpentsStrikerFistDebuff',
});
export const JudgementOfLight = makeBooleanDebuffInput({
	actionId: player =>
		player.getMatchingSpellActionId([
			{ id: 20185, minLevel: 30, maxLevel: 39 },
			{ id: 20344, minLevel: 40, maxLevel: 49 },
			{ id: 20345, minLevel: 50, maxLevel: 59 },
			{ id: 20346, minLevel: 60 },
		]),
	fieldName: 'judgementOfLight',
	showWhen: player => player.getFaction() === Faction.Alliance,
});
export const CurseOfVulnerability = makeBooleanDebuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 427143, minLevel: 25 }]),
	fieldName: 'curseOfVulnerability',
});
export const GiftOfArthas = makeBooleanDebuffInput({
	actionId: player =>
		player.getMatchingSpellActionId([
			// SoD Phase 3?
			{ id: 11374, minLevel: 41 },
		]),
	fieldName: 'giftOfArthas',
});
export const CrystalYield = makeBooleanDebuffInput({
	actionId: player => player.getMatchingSpellActionId([{ id: 15235, minLevel: 47 }]),
	fieldName: 'crystalYield',
});
export const AncientCorrosivePoison = makeMultistateMultiplierDebuffInput({
	actionId: () => ActionId.fromItemId(209562),
	numStates: 11,
	multiplier: 10,
	reverse: true,
	fieldName: 'ancientCorrosivePoison',
});

///////////////////////////////////////////////////////////////////////////
//                                 CONFIGS
///////////////////////////////////////////////////////////////////////////

export const RAID_BUFFS_CONFIG = [
	// Core Stat Buffs
	{
		config: AllStatsBuff,
		picker: IconPicker,
		stats: [],
	},
	{
		config: AllStatsPercentBuff,
		picker: MultiIconPicker,
		stats: [],
	},
	{
		config: StaminaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina],
	},
	{
		config: BloodPactBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina],
	},
	{
		config: IntellectBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatIntellect],
	},
	{
		config: SpiritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpirit],
	},

	// Tank-related Buffs
	{
		config: ArmorBuff,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: PhysDamReductionBuff,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	// {
	// 	config: DamageReductionPercentBuff,
	// 	picker: IconPicker,
	// 	stats: [Stat.StatArmor],
	// },
	{
		config: ResistanceBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatNatureResistance, Stat.StatShadowResistance, Stat.StatFireResistance, Stat.StatFrostResistance],
	},

	// Physical Damage Buffs
	{
		config: PaladinPhysicalBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatStrength, Stat.StatAgility],
	},
	{
		config: StrengthBuffHorde,
		picker: IconPicker,
		stats: [Stat.StatStrength],
	},
	{
		config: BattleShoutBuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower],
	},
	{
		config: GraceOfAir,
		picker: IconPicker,
		stats: [Stat.StatAgility],
	},
	{
		config: TrueshotAuraBuff,
		picker: IconPicker,
		stats: [Stat.StatRangedAttackPower],
	},
	{
		config: MeleeCritBuff,
		picker: IconPicker,
		stats: [Stat.StatMeleeCrit],
	},
	// Threat Buffs
	{
		config: HordeThreatBuff,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	// Spell Damage Buffs
	{
		config: SpellIncreaseBuff,
		picker: IconPicker,
		stats: [Stat.StatSpellPower],
	},
	{
		config: SpellCritBuff,
		picker: IconPicker,
		stats: [Stat.StatSpellCrit],
	},
	{
		config: BlessingOfWisdom,
		picker: IconPicker,
		stats: [Stat.StatMP5],
	},
	{
		config: ManaSpringTotem,
		picker: IconPicker,
		stats: [Stat.StatMP5],
	},
	{
		config: VampiricTouchReplenishment,
		picker: IconPicker,
		stats: [Stat.StatMP5],
	},
] as PickerStatOptions[];

export const MISC_BUFFS_CONFIG = [
	{
		config: AtieshSpellPowerBuff,
		picker: IconPicker,
		stats: [],
	},
	{
		config: AtieshSpellCritBuff,
		picker: IconPicker,
		stats: [],
	},
	{
		config: AtieshCastSpeedBuff,
		picker: IconPicker,
		stats: [],
	},
	{
		config: AtieshHealingBuff,
		picker: IconPicker,
		stats: [],
	},
	{
		config: ImprovedStoneskinWindwall,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: Thorns,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: RetributionAura,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: SanctityAura,
		picker: IconPicker,
		stats: [Stat.StatHolyPower],
	},
	{
		config: Innervate,
		picker: IconPicker,
		stats: [Stat.StatMP5],
	},
	{
		config: PowerInfusion,
		picker: IconPicker,
		stats: [],
	},
	{
		config: BattleSquawkBuff,
		picker: IconPicker,
		stats: [Stat.StatMeleeHit],
	},
] as PickerStatOptions[];

export const WORLD_BUFFS_CONFIG = [
	// {
	// 	config: RallyingCryOfTheDragonslayer,
	// 	picker: IconPicker,
	// 	stats: [
	// 		Stat.StatMeleeCrit,
	// 		// TODO: Stat.StatRangedCrit,
	// 		Stat.StatSpellCrit,
	// 		Stat.StatAttackPower,
	// 	],
	// },
	{
		config: DragonslayerBuffInput,
		picker: MultiIconPicker,
		stats: [],
	},
	{
		config: SongflowerSerenade,
		picker: IconPicker,
		stats: [],
	},
	{
		config: SpiritOfZandalar,
		picker: IconPicker,
		stats: [],
	},
	{
		config: WarchiefsBlessing,
		picker: IconPicker,
		stats: [],
	},
	{
		config: MightOfStormwind,
		picker: IconPicker,
		stats: [],
	},
	{
		config: FervorOfTheTempleExplorer,
		picker: IconPicker,
		stats: [],
	},
	{
		config: SparkOfInspiration,
		picker: IconPicker,
		stats: [],
	},
	{
		config: BoonOfBlackfathom,
		picker: IconPicker,
		stats: [
			Stat.StatMeleeCrit,
			// TODO: Stat.StatRangedCrit,
			Stat.StatSpellCrit,
			Stat.StatAttackPower,
		],
	},
	{
		config: AshenvalePvpBuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatSpellPower],
	},
	{
		config: FengusFerocity,
		picker: IconPicker,
		stats: [Stat.StatAttackPower],
	},
	{
		config: MoldarsMoxie,
		picker: IconPicker,
		stats: [Stat.StatStamina],
	},
	{
		config: SlipKiksSavvy,
		picker: IconPicker,
		stats: [Stat.StatSpellCrit],
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
	// Standard Debuffs
	{
		config: MajorArmorDebuff,
		stats: [Stat.StatAttackPower],
		picker: MultiIconPicker,
	},
	{
		config: CurseOfRecklessness,
		picker: IconPicker,
		stats: [Stat.StatAttackPower],
	},
	{
		config: FaerieFire,
		picker: IconPicker,
		stats: [Stat.StatAttackPower],
	},
	{
		config: BleedDebuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},

	// Magic
	{
		config: JudgementOfTheCrusader,
		picker: IconPicker,
		stats: [Stat.StatHolyPower],
	},
	{
		config: SpellISBDebuff,
		picker: IconPicker,
		stats: [Stat.StatShadowPower],
	},
	{
		config: SpellScorchDebuff,
		picker: IconPicker,
		stats: [],
	},
	{
		config: SpellWintersChillDebuff,
		picker: IconPicker,
		stats: [Stat.StatFrostPower],
	},
	{
		config: NatureSpellDamageDebuff,
		picker: MultiIconPicker,
		stats: [],
	},
	{
		config: SpellShadowWeavingDebuff,
		picker: IconPicker,
		stats: [Stat.StatShadowPower],
	},
	{
		config: WarlockCursesConfig,
		picker: MultiIconPicker,
		stats: [],
	},
	{
		config: OccultPoison,
		picker: IconPicker,
		stats: [],
	},

	// Defensive
	{
		config: AttackPowerDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: MeleeAttackSpeedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: curseOfWeaknessDebuff,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: MeleeHitDebuff,
		picker: IconPicker,
		stats: [Stat.StatDodge],
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
	},
] as PickerStatOptions[];

export const MISC_DEBUFFS_CONFIG = [
	// Misc Debuffs
	{
		config: ImprovedFaerieFire,
		picker: IconPicker,
		stats: [],
	},
	{
		config: MeleeHunter2pcT1Bonus,
		picker: IconPicker,
		stats: [Stat.StatMeleeHit],
	},
	{
		config: MekkatorqueFistDebuff,
		picker: IconPicker,
		stats: [],
	},
	{
		config: SerpentsStrikerFistDebuff,
		picker: IconPicker,
		stats: [],
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
	{
		config: JudgementOfLight,
		picker: IconPicker,
		stats: [Stat.StatStamina],
	},
] as PickerStatOptions[];
