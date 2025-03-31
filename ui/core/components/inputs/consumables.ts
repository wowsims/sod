import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import {
	AgilityElixir,
	Alcohol,
	ArmorElixir,
	AttackPowerBuff,
	Class,
	Conjured,
	Consumes,
	EnchantedSigil,
	Explosive,
	FirePowerBuff,
	Flask,
	Food,
	FrostPowerBuff,
	HandType,
	HealthElixir,
	ItemSlot,
	MageScroll,
	ManaRegenElixir,
	Potions,
	Profession,
	SapperExplosive,
	SealOfTheDawn,
	ShadowPowerBuff,
	Spec,
	SpellPowerBuff,
	Stat,
	StrengthBuff,
	WeaponImbue,
	WeaponType,
	ZanzaBuff,
} from '../../proto/common';
import { ActionId } from '../../proto_utils/action_id';
import { isBluntWeaponType, isSharpWeaponType, isWeapon } from '../../proto_utils/utils';
import { EventID, TypedEvent } from '../../typed_event';
import { IconEnumValueConfig } from '../icon_enum_picker';
import { makeBooleanConsumeInput, makeBooleanMiscConsumeInput, makeBooleanPetMiscConsumeInput, makeEnumConsumeInput } from '../icon_inputs';
import { IconPicker, IconPickerDirection } from '../icon_picker';
import * as InputHelpers from '../input_helpers';
import { MultiIconPicker, MultiIconPickerConfig, MultiIconPickerItemConfig } from '../multi_icon_picker';
import {
	AtrophicPoisonWeaponImbue,
	DeadlyPoisonWeaponImbue,
	InstantPoisonWeaponImbue,
	NumbingPoisonWeaponImbue,
	OccultPoisonWeaponImbue,
	SebaciousPoisonWeaponImbue,
	WoundPoisonWeaponImbue,
} from './rogue_imbues';
import { FlametongueWeaponImbue, FrostbrandWeaponImbue, RockbiterWeaponImbue, WindfuryWeaponImbue } from './shaman_imbues';
import { ActionInputConfig, ItemStatOption, PickerStatOptions, StatOptions } from './stat_options';

export interface ConsumableInputConfig<T> extends ActionInputConfig<T> {
	value: T;
}

export interface ConsumableStatOption<T> extends ItemStatOption<T> {
	config: ConsumableInputConfig<T>;
}

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes;
	numColumns?: number;
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventactionId: EventID, player: Player<any>, newValue: T) => void;
	showWhen?: (player: Player<any>) => boolean;
}

function makeConsumeInputFactory<T extends number>(
	args: ConsumeInputFactoryArgs<T>,
): (options: ConsumableStatOption<T>[], tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
	return (options: ConsumableStatOption<T>[], tooltip?: string) => {
		return {
			type: 'iconEnum',
			tooltip: tooltip,
			numColumns: args.numColumns ? args.numColumns : options.length > 11 ? 4 : options.length > 8 ? 3 : options.length > 5 ? 2 : 1,
			values: [{ value: 0 } as unknown as IconEnumValueConfig<Player<any>, T>].concat(
				options.map(option => {
					return {
						actionId: option.config.actionId,
						value: option.config.value,
						showWhen: (player: Player<any>) => !option.config.showWhen || option.config.showWhen(player),
					} as IconEnumValueConfig<Player<any>, T>;
				}),
			),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,
			changedEvent: (player: Player<any>) =>
				TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter, player.gearChangeEmitter, player.professionChangeEmitter]),
			showWhen: (player: Player<any>) => !args.showWhen || args.showWhen(player),
			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();

				if (newConsumes[args.consumesFieldName] === newValue) {
					return;
				}

				(newConsumes[args.consumesFieldName] as number) = newValue;
				TypedEvent.freezeAllAndDo(() => {
					player.setConsumes(eventID, newConsumes);
					if (args.onSet) {
						args.onSet(eventID, player, newValue as T);
					}
				});
			},
		};
	};
}

type MultiIconConsumeInputFactoryArg<ModObject> = Omit<MultiIconPickerConfig<ModObject>, 'values'>;

export const makeMultiIconConsumesInputFactory = <ModObject>(
	config: MultiIconConsumeInputFactoryArg<ModObject>,
): ((parent: HTMLElement, modObj: ModObject, simUI: IndividualSimUI<Spec>, options: StatOptions<any, any>) => MultiIconPicker<any>) => {
	return (parent: HTMLElement, modObj: ModObject, simUI: IndividualSimUI<Spec>, options: StatOptions<any, any>) => {
		const pickerConfig = {
			...config,
			values: options.map(option => option.config) as Array<MultiIconPickerItemConfig<ModObject>>,
		};
		return new MultiIconPicker(parent, modObj, pickerConfig, simUI);
	};
};

///////////////////////////////////////////////////////////////////////////
//                                 CONJURED
///////////////////////////////////////////////////////////////////////////

export const ConjuredHealthstone: ConsumableInputConfig<Conjured> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 5509, minLevel: 24 }]),
	value: Conjured.ConjuredHealthstone,
};
export const ConjuredGreaterHealthstone: ConsumableInputConfig<Conjured> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 5510, minLevel: 36 }]),
	value: Conjured.ConjuredGreaterHealthstone,
};
export const ConjuredMajorHealthstone: ConsumableInputConfig<Conjured> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 9421, minLevel: 48 }]),
	value: Conjured.ConjuredMajorHealthstone,
};

export const ConjuredMinorRecombobulator: ConsumableInputConfig<Conjured> = {
	actionId: () => ActionId.fromItemId(4381),
	value: Conjured.ConjuredMinorRecombobulator,
	showWhen: (player: Player<any>) => player.getGear().hasTrinket(4381),
};
export const ConjuredDemonicRune: ConsumableInputConfig<Conjured> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 12662, minLevel: 40 }]),
	value: Conjured.ConjuredDemonicRune,
};
export const ConjuredRogueThistleTea: ConsumableInputConfig<Conjured> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 7676, minLevel: 25 }]),
	value: Conjured.ConjuredRogueThistleTea,
	showWhen: player => player.getClass() == Class.ClassRogue,
};

export const CONJURED_CONFIG: ConsumableStatOption<Conjured>[] = [
	{ config: ConjuredMajorHealthstone, stats: [Stat.StatArmor] },
	{ config: ConjuredGreaterHealthstone, stats: [Stat.StatArmor] },
	{ config: ConjuredHealthstone, stats: [Stat.StatArmor] },

	{ config: ConjuredDemonicRune, stats: [Stat.StatIntellect] },
	{ config: ConjuredMinorRecombobulator, stats: [Stat.StatIntellect] },

	{ config: ConjuredRogueThistleTea, stats: [] },
];

export const makeConjuredInput = makeConsumeInputFactory({ consumesFieldName: 'defaultConjured' });

///////////////////////////////////////////////////////////////////////////
//                             ENCHANTING SIGIL
///////////////////////////////////////////////////////////////////////////

export const EnchantedSigilInnovation: ConsumableInputConfig<EnchantedSigil> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 217308, minLevel: 40 }]),
	value: EnchantedSigil.InnovationSigil,
	showWhen: player => player.hasProfession(Profession.Enchanting),
};
export const EnchantedSigilLivingDreams: ConsumableInputConfig<EnchantedSigil> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 221028, minLevel: 50 }]),
	value: EnchantedSigil.LivingDreamsSigil,
	showWhen: player => player.hasProfession(Profession.Enchanting),
};
export const EnchantedSigilFlowingWaters: ConsumableInputConfig<EnchantedSigil> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 228978, minLevel: 50 }]),
	value: EnchantedSigil.FlowingWatersSigil,
};
export const EnchantedSigilWrathOfTheStorm: ConsumableInputConfig<EnchantedSigil> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233995, minLevel: 50 }]),
	value: EnchantedSigil.WrathOfTheStormSigil,
};

export const ENCHANTED_SIGIL_CONFIG: ConsumableStatOption<EnchantedSigil>[] = [
	{ config: EnchantedSigilWrathOfTheStorm, stats: [] },
	{ config: EnchantedSigilFlowingWaters, stats: [] },
	{ config: EnchantedSigilLivingDreams, stats: [] },
	{ config: EnchantedSigilInnovation, stats: [] },
];

export const makeEncanthedSigilInput = makeConsumeInputFactory({ consumesFieldName: 'enchantedSigil' });

///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

export const SapperGoblinSapper: ConsumableInputConfig<SapperExplosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 10646, minLevel: 50 }]),
	showWhen: player => player.hasProfession(Profession.Engineering),
	value: SapperExplosive.SapperGoblinSapper,
};

export const SapperFumigator: ConsumableInputConfig<SapperExplosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233985, minLevel: 60 }]),
	value: SapperExplosive.SapperFumigator,
};

export const ExplosiveStratholmeHolyWater: ConsumableInputConfig<Explosive> = {
	actionId: () => ActionId.fromItemId(13180),
	value: Explosive.ExplosiveStratholmeHolyWater,
};

export const ExplosiveObsidianBomb: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233986, minLevel: 60 }]),
	value: Explosive.ExplosiveObsidianBomb,
};

export const ExplosiveSolidDynamite: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 10507, minLevel: 40 }]),
	showWhen: player => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveSolidDynamite,
};

export const ExplosiveGoblinLandMine: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 4395, minLevel: 40 }]),
	showWhen: player => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveGoblinLandMine,
};

export const ExplosiveDenseDynamite: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 18641, minLevel: 50 }]),
	showWhen: player => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveDenseDynamite,
};

export const ExplosiveThoriumGrenade: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 15993, minLevel: 50 }]),
	showWhen: player => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveThoriumGrenade,
};

export const ExplosiveEzThroRadiationBomb: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 215168, minLevel: 40 }]),
	value: Explosive.ExplosiveEzThroRadiationBomb,
};

export const ExplosiveHighYieldRadiationBomb: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 215127, minLevel: 40 }]),
	showWhen: player => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveHighYieldRadiationBomb,
};

export const EXPLOSIVES_CONFIG: ConsumableStatOption<Explosive>[] = [
	{ config: ExplosiveStratholmeHolyWater, stats: [] },
	{ config: ExplosiveObsidianBomb, stats: [] },
	{ config: ExplosiveEzThroRadiationBomb, stats: [] },
	{ config: ExplosiveHighYieldRadiationBomb, stats: [] },
	{ config: ExplosiveSolidDynamite, stats: [] },
	{ config: ExplosiveDenseDynamite, stats: [] },
	{ config: ExplosiveThoriumGrenade, stats: [] },
	{ config: ExplosiveGoblinLandMine, stats: [] },
];

export const SAPPER_CONFIG: ConsumableStatOption<SapperExplosive>[] = [
	{ config: SapperFumigator, stats: [] },
	{ config: SapperGoblinSapper, stats: [] },
];

export const makeExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'fillerExplosive',
	//showWhen: (player) => !!player.getProfessions().find(p => p == Profession.Engineering),
});

export const makeSappersInput = makeConsumeInputFactory({
	consumesFieldName: 'sapperExplosive',
});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS
///////////////////////////////////////////////////////////////////////////

// Original lvl 50 not obtainable in Phase 3
export const FlaskOfTheTitans: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13510, minLevel: 51 }]),
	value: Flask.FlaskOfTheTitans,
};
export const FlaskOfTheOldGods: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233965, minLevel: 60 }]),
	value: Flask.FlaskOfTheOldGods,
};
export const FlaskOfUnyieldingSorrow: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233966, minLevel: 60 }]),
	value: Flask.FlaskOfUnyieldingSorrow,
};
export const FlaskOfAncientKnowledge: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233964, minLevel: 60 }]),
	value: Flask.FlaskOfAncientKnowledge,
};
// Original lvl 50 not obtainable in Phase 3
export const FlaskOfDistilledWisdom: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13511, minLevel: 51 }]),
	value: Flask.FlaskOfDistilledWisdom,
};
export const FlaskOfMadness: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233962, minLevel: 60 }]),
	value: Flask.FlaskOfMadness,
};
// Original lvl 50 not obtainable in Phase 3
export const FlaskOfSupremePower: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13512, minLevel: 51 }]),
	value: Flask.FlaskOfSupremePower,
};
// Original lvl 50 not obtainable in Phase 3
export const FlaskOfChromaticResistance: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13513, minLevel: 51 }]),
	value: Flask.FlaskOfChromaticResistance,
};
export const FlaskOfRestlessDreams: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 222952, minLevel: 50, maxLevel: 59 }]),
	value: Flask.FlaskOfRestlessDreams,
	showWhen: player => player.hasProfession(Profession.Alchemy),
};
export const FlaskOfEverlastingNightmares: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 221024, minLevel: 50, maxLevel: 59 }]),
	value: Flask.FlaskOfEverlastingNightmares,
	showWhen: player => player.hasProfession(Profession.Alchemy),
};

export const FLASKS_CONFIG: ConsumableStatOption<Flask>[] = [
	{ config: FlaskOfTheTitans, stats: [Stat.StatStamina] },
	{ config: FlaskOfTheOldGods, stats: [Stat.StatStamina] },
	{ config: FlaskOfMadness, stats: [Stat.StatAttackPower] },
	{ config: FlaskOfDistilledWisdom, stats: [Stat.StatIntellect] },
	{ config: FlaskOfUnyieldingSorrow, stats: [Stat.StatIntellect] },
	{ config: FlaskOfSupremePower, stats: [Stat.StatMP5, Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: FlaskOfAncientKnowledge, stats: [Stat.StatMP5, Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: FlaskOfChromaticResistance, stats: [] },
	{ config: FlaskOfRestlessDreams, stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: FlaskOfEverlastingNightmares, stats: [Stat.StatAttackPower] },
];

export const makeFlasksInput = makeConsumeInputFactory({ consumesFieldName: 'flask' });

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const ProwlerSteak: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 238637, minLevel: 55 }]),
	value: Food.FoodProwlerSteak,
};
export const FiletOFlank: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 238638, minLevel: 55 }]),
	value: Food.FoodFiletOFlank,
};
export const SunriseOmelette: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 238639, minLevel: 55 }]),
	value: Food.FoodSunriseOmelette,
};
export const SpecklefinFeast: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 238641, minLevel: 55 }]),
	value: Food.FoodSpecklefinFeast,
};
export const GrandLobsterBanquet: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 238642, minLevel: 55 }]),
	value: Food.FoodGrandLobsterBanquet,
};
export const DarkclawBisque: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 232436, minLevel: 45 }]),
	value: Food.FoodDarkclawBisque,
};
export const SmokedRedgill: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 232438, minLevel: 45 }]),
	value: Food.FoodSmokedRedgill,
};
export const DirgesKickChimaerokChops: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 21023, minLevel: 55 }]),
	value: Food.FoodDirgesKickChimaerokChops,
};
export const GrilledSquid: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13928, minLevel: 50 }]),
	value: Food.FoodGrilledSquid,
};
export const SmokedDesertDumpling: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 20452, minLevel: 50 }]),
	value: Food.FoodSmokedDesertDumpling,
};
export const RunnTumTuberSurprise: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 18254, minLevel: 45 }]),
	value: Food.FoodRunnTumTuberSurprise,
};
export const BlessSunfruit: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13810, minLevel: 45 }]),
	value: Food.FoodBlessSunfruit,
};
export const BlessedSunfruitJuice: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13813, minLevel: 45 }]),
	value: Food.FoodBlessedSunfruitJuice,
};
export const NightfinSoup: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13931, minLevel: 35 }]),
	value: Food.FoodNightfinSoup,
};
export const TenderWolfSteak: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 18045, minLevel: 40 }]),
	value: Food.FoodTenderWolfSteak,
};
export const SagefishDelight: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 21217, minLevel: 30 }]),
	value: Food.FoodSagefishDelight,
};
export const HotWolfRibs: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13851, minLevel: 25 }]),
	value: Food.FoodHotWolfRibs,
};
export const SmokedSagefish: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 21072, minLevel: 10 }]),
	value: Food.FoodSmokedSagefish,
};

// Ordered by level
export const FOOD_CONFIG: ConsumableStatOption<Food>[] = [
	{ config: SunriseOmelette, stats: [Stat.StatSpellPower, Stat.StatHealingPower, Stat.StatStamina] },
	{ config: ProwlerSteak, stats: [Stat.StatStrength, Stat.StatStamina] },
	{ config: FiletOFlank, stats: [Stat.StatAgility, Stat.StatStamina] },
	{ config: SpecklefinFeast, stats: [Stat.StatAttackPower, Stat.StatSpellPower, Stat.StatHealingPower, Stat.StatStamina] },
	{ config: GrandLobsterBanquet, stats: [Stat.StatAttackPower, Stat.StatSpellPower, Stat.StatHealingPower, Stat.StatStamina] },
	{ config: DarkclawBisque, stats: [Stat.StatSpellDamage] },
	{ config: SmokedRedgill, stats: [Stat.StatHealingPower] },
	{ config: DirgesKickChimaerokChops, stats: [Stat.StatStamina] },
	{ config: GrilledSquid, stats: [Stat.StatAgility] },
	{ config: SmokedDesertDumpling, stats: [Stat.StatStrength] },
	{ config: RunnTumTuberSurprise, stats: [Stat.StatIntellect] },
	{ config: BlessSunfruit, stats: [Stat.StatStrength] },
	{ config: BlessedSunfruitJuice, stats: [Stat.StatSpirit] },
	{ config: NightfinSoup, stats: [Stat.StatMP5] },
	{ config: TenderWolfSteak, stats: [Stat.StatStamina, Stat.StatSpirit] },
	{ config: SagefishDelight, stats: [Stat.StatMP5] },
	{ config: HotWolfRibs, stats: [Stat.StatSpirit] },
	{ config: SmokedSagefish, stats: [Stat.StatMP5] },
];

export const makeFoodInput = makeConsumeInputFactory({ consumesFieldName: 'food' });

export const DragonBreathChili = makeBooleanConsumeInput({
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 12217, minLevel: 35 }]),
	fieldName: 'dragonBreathChili',
});

export const RumseyRumBlackLabel: ConsumableInputConfig<Alcohol> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 21151, minLevel: 1 }]),
	value: Alcohol.AlcoholRumseyRumBlackLabel,
};
export const GordokGreenGrog: ConsumableInputConfig<Alcohol> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 18269, minLevel: 56 }]),
	value: Alcohol.AlcoholGordokGreenGrog,
};
export const RumseyRumDark: ConsumableInputConfig<Alcohol> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 21114, minLevel: 1 }]),
	value: Alcohol.AlcoholRumseyRumDark,
};
export const RumseyRumLight: ConsumableInputConfig<Alcohol> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 20709, minLevel: 1 }]),
	value: Alcohol.AlcoholRumseyRumLight,
};
export const KreegsStoutBeatdown: ConsumableInputConfig<Alcohol> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 18284, minLevel: 56 }]),
	value: Alcohol.AlcoholKreegsStoutBeatdown,
};

export const ALCOHOL_CONFIG: ConsumableStatOption<Alcohol>[] = [
	{ config: RumseyRumBlackLabel, stats: [Stat.StatStamina] },
	{ config: GordokGreenGrog, stats: [Stat.StatStamina] },
	{ config: RumseyRumDark, stats: [Stat.StatStamina] },
	{ config: RumseyRumLight, stats: [Stat.StatStamina] },
	{ config: KreegsStoutBeatdown, stats: [Stat.StatSpirit] },
];

export const makeAlcoholInput = makeConsumeInputFactory({ consumesFieldName: 'alcohol' });

///////////////////////////////////////////////////////////////////////////
//                                 DEFENSIVE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Armor
export const ElixirOfTheIronside: ConsumableInputConfig<ArmorElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233969, minLevel: 60 }]),
	value: ArmorElixir.ElixirOfTheIronside,
};
export const ElixirOfSuperiorDefense: ConsumableInputConfig<ArmorElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13445, minLevel: 43 }]),
	value: ArmorElixir.ElixirOfSuperiorDefense,
};
export const ElixirOfGreaterDefense: ConsumableInputConfig<ArmorElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 8951, minLevel: 29 }]),
	value: ArmorElixir.ElixirOfGreaterDefense,
};
export const ElixirOfDefense: ConsumableInputConfig<ArmorElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 3389, minLevel: 16 }]),
	value: ArmorElixir.ElixirOfDefense,
};
export const ElixirOfMinorDefense: ConsumableInputConfig<ArmorElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 5997, minLevel: 1 }]),
	value: ArmorElixir.ElixirOfMinorDefense,
};
export const ScrollOfProtection: ConsumableInputConfig<ArmorElixir> = {
	actionId: player =>
		player.getMatchingItemActionId([
			{ id: 3013, minLevel: 1, maxLevel: 14 },
			{ id: 1478, minLevel: 15, maxLevel: 29 },
			{ id: 4421, minLevel: 30, maxLevel: 44 },
			{ id: 10305, minLevel: 45 },
		]),
	value: ArmorElixir.ScrollOfProtection,
};
export const ARMOR_CONSUMES_CONFIG: ConsumableStatOption<ArmorElixir>[] = [
	{ config: ElixirOfTheIronside, stats: [Stat.StatArmor] },
	{ config: ElixirOfSuperiorDefense, stats: [Stat.StatArmor] },
	{ config: ElixirOfGreaterDefense, stats: [Stat.StatArmor] },
	{ config: ElixirOfDefense, stats: [Stat.StatArmor] },
	{ config: ElixirOfMinorDefense, stats: [Stat.StatArmor] },
	{ config: ScrollOfProtection, stats: [Stat.StatArmor] },
];

export const makeArmorConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'armorElixir' });

// Health
export const ElixirOfFortitude: ConsumableInputConfig<HealthElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 3825, minLevel: 25 }]),
	value: HealthElixir.ElixirOfFortitude,
};
export const ElixirOfMinorFortitude: ConsumableInputConfig<HealthElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 2458, minLevel: 2 }]),
	value: HealthElixir.ElixirOfMinorFortitude,
};
export const HEALTH_CONSUMES_CONFIG: ConsumableStatOption<HealthElixir>[] = [
	{ config: ElixirOfFortitude, stats: [Stat.StatStamina] },
	{ config: ElixirOfMinorFortitude, stats: [Stat.StatStamina] },
];

export const makeHealthConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'healthElixir' });

///////////////////////////////////////////////////////////////////////////
//                                 PHYSICAL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Attack Power
export const JujuMight: ConsumableInputConfig<AttackPowerBuff> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 12460, minLevel: 55 }]),
	value: AttackPowerBuff.JujuMight,
};
export const WinterfallFirewater: ConsumableInputConfig<AttackPowerBuff> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 12820, minLevel: 45 }]),
	value: AttackPowerBuff.WinterfallFirewater,
};

export const ATTACK_POWER_CONSUMES_CONFIG: ConsumableStatOption<AttackPowerBuff>[] = [
	{ config: JujuMight, stats: [Stat.StatAttackPower] },
	{ config: WinterfallFirewater, stats: [Stat.StatAttackPower] },
];

export const makeAttackPowerConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'attackPowerBuff' });

// Agility
export const ElixirOfTheHoneyBadger: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 233967, minLevel: 60 }]),
	value: AgilityElixir.ElixirOfTheHoneyBadger,
};
export const ElixirOfTheMongoose: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13452, minLevel: 46 }]),
	value: AgilityElixir.ElixirOfTheMongoose,
};
export const ElixirOfGreaterAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 9187, minLevel: 38 }]),
	value: AgilityElixir.ElixirOfGreaterAgility,
};
export const ElixirOfAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 8949, minLevel: 27 }]),
	value: AgilityElixir.ElixirOfAgility,
};
export const ElixirOfLesserAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 3390, minLevel: 18 }]),
	value: AgilityElixir.ElixirOfLesserAgility,
};
export const ScrollOfAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: player =>
		player.getMatchingItemActionId([
			{ id: 3012, minLevel: 10, maxLevel: 24 },
			{ id: 1477, minLevel: 25, maxLevel: 39 },
			{ id: 4425, minLevel: 40, maxLevel: 54 },
			{ id: 10309, minLevel: 55 },
		]),
	value: AgilityElixir.ScrollOfAgility,
};

export const AGILITY_CONSUMES_CONFIG: ConsumableStatOption<AgilityElixir>[] = [
	{ config: ElixirOfTheHoneyBadger, stats: [Stat.StatAgility, Stat.StatMeleeCrit] },
	{ config: ElixirOfTheMongoose, stats: [Stat.StatAgility, Stat.StatMeleeCrit] },
	{ config: ElixirOfGreaterAgility, stats: [Stat.StatAgility] },
	{ config: ElixirOfAgility, stats: [Stat.StatAgility] },
	{ config: ElixirOfLesserAgility, stats: [Stat.StatAgility] },
	{ config: ScrollOfAgility, stats: [Stat.StatAgility] },
];

export const makeAgilityConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'agilityElixir' });

// Strength
export const JujuPower: ConsumableInputConfig<StrengthBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 12451, minLevel: 55 }]),
	value: StrengthBuff.JujuPower,
};
export const ElixirOfGiants: ConsumableInputConfig<StrengthBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 9206, minLevel: 46 }]),
	value: StrengthBuff.ElixirOfGiants,
};
export const ElixirOfOgresStrength: ConsumableInputConfig<StrengthBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 3391, minLevel: 20 }]),
	value: StrengthBuff.ElixirOfOgresStrength,
};
export const ScrollOfStrength: ConsumableInputConfig<StrengthBuff> = {
	actionId: player =>
		player.getMatchingItemActionId([
			{ id: 954, minLevel: 10, maxLevel: 24 },
			{ id: 2289, minLevel: 25, maxLevel: 39 },
			{ id: 4426, minLevel: 40, maxLevel: 54 },
			{ id: 10310, minLevel: 55 },
		]),
	value: StrengthBuff.ScrollOfStrength,
};

export const STRENGTH_CONSUMES_CONFIG: ConsumableStatOption<StrengthBuff>[] = [
	{ config: JujuPower, stats: [Stat.StatStrength] },
	{ config: ElixirOfGiants, stats: [Stat.StatStrength] },
	{ config: ElixirOfOgresStrength, stats: [Stat.StatStrength] },
	{ config: ScrollOfStrength, stats: [Stat.StatStrength] },
];

export const makeStrengthConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'strengthBuff' });

///////////////////////////////////////////////////////////////////////////
//                                 Misc Throughput Consumes
///////////////////////////////////////////////////////////////////////////

// Seal of the Dawn Consumes
export const SEAL_OF_THE_DAWN_CONSUMES_CONFIG: ConsumableStatOption<SealOfTheDawn>[] = [
	{ config: { actionId: () => ActionId.fromItemId(236364), value: SealOfTheDawn.SealOfTheDawnDamageR10 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236363), value: SealOfTheDawn.SealOfTheDawnDamageR9 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236362), value: SealOfTheDawn.SealOfTheDawnDamageR8 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236361), value: SealOfTheDawn.SealOfTheDawnDamageR7 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236360), value: SealOfTheDawn.SealOfTheDawnDamageR6 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236358), value: SealOfTheDawn.SealOfTheDawnDamageR5 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236357), value: SealOfTheDawn.SealOfTheDawnDamageR4 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236356), value: SealOfTheDawn.SealOfTheDawnDamageR3 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236355), value: SealOfTheDawn.SealOfTheDawnDamageR2 }, stats: [] },
	{ config: { actionId: () => ActionId.fromItemId(236354), value: SealOfTheDawn.SealOfTheDawnDamageR1 }, stats: [] },

	{ config: { actionId: () => ActionId.fromItemId(236386), value: SealOfTheDawn.SealOfTheDawnTankR10 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236388), value: SealOfTheDawn.SealOfTheDawnTankR9 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236389), value: SealOfTheDawn.SealOfTheDawnTankR8 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236390), value: SealOfTheDawn.SealOfTheDawnTankR7 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236391), value: SealOfTheDawn.SealOfTheDawnTankR6 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236392), value: SealOfTheDawn.SealOfTheDawnTankR5 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236393), value: SealOfTheDawn.SealOfTheDawnTankR4 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236394), value: SealOfTheDawn.SealOfTheDawnTankR3 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236395), value: SealOfTheDawn.SealOfTheDawnTankR2 }, stats: [Stat.StatDefense] },
	{ config: { actionId: () => ActionId.fromItemId(236396), value: SealOfTheDawn.SealOfTheDawnTankR1 }, stats: [Stat.StatDefense] },
];
export const makeSealOfTheDawnConsumesInput = makeConsumeInputFactory({
	consumesFieldName: 'sealOfTheDawn',
	numColumns: 11,
	showWhen: player => player.getLevel() === 60,
});

// Blasted Lands Consumes
export const ROIDS: ConsumableInputConfig<ZanzaBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 8410, minLevel: 45 }]),
	value: ZanzaBuff.ROIDS,
};
export const GroundScorpokAssay: ConsumableInputConfig<ZanzaBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 8412, minLevel: 45 }]),
	value: ZanzaBuff.GroundScorpokAssay,
};
export const LungJuiceCocktail: ConsumableInputConfig<ZanzaBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 8411, minLevel: 45 }]),
	value: ZanzaBuff.LungJuiceCocktail,
};
export const CerebralCortexCompound: ConsumableInputConfig<ZanzaBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 8423, minLevel: 45 }]),
	value: ZanzaBuff.CerebralCortexCompound,
};
export const GizzardGum: ConsumableInputConfig<ZanzaBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 8424, minLevel: 45 }]),
	value: ZanzaBuff.GizzardGum,
};

// Zanza Potions
export const SpiritOfZanza: ConsumableInputConfig<ZanzaBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 20079, minLevel: 55 }]),
	value: ZanzaBuff.SpiritOfZanza,
};

// Atal'ai Potions
export const AtalAiMojoOfWar: ConsumableInputConfig<ZanzaBuff> = {
	actionId: () => ActionId.fromItemId(221196),
	value: ZanzaBuff.AtalaiMojoOfWar,
	showWhen: (player: Player<any>) => player.getLevel() == 50,
};

export const AtalAiMojoOfForbiddenMagic: ConsumableInputConfig<ZanzaBuff> = {
	actionId: () => ActionId.fromItemId(221030),
	value: ZanzaBuff.AtalaiMojoOfForbiddenMagic,
	showWhen: (player: Player<any>) => player.getLevel() == 50,
};

export const AtalAiMojoOfLife: ConsumableInputConfig<ZanzaBuff> = {
	actionId: () => ActionId.fromItemId(221311),
	value: ZanzaBuff.AtalaiMojoOfLife,
	showWhen: (player: Player<any>) => player.getLevel() == 50,
};

export const ZANZA_BUFF_CONSUMES_CONFIG: ConsumableStatOption<ZanzaBuff>[] = [
	{ config: SpiritOfZanza, stats: [Stat.StatStamina, Stat.StatSpirit] },
	{ config: ROIDS, stats: [Stat.StatStrength] },
	{ config: GroundScorpokAssay, stats: [Stat.StatAgility] },
	{ config: LungJuiceCocktail, stats: [Stat.StatStamina] },
	{ config: CerebralCortexCompound, stats: [Stat.StatIntellect] },
	{ config: GizzardGum, stats: [Stat.StatSpirit] },
	{ config: AtalAiMojoOfWar, stats: [Stat.StatAttackPower] },
	{ config: AtalAiMojoOfForbiddenMagic, stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: AtalAiMojoOfLife, stats: [Stat.StatHealingPower] },
];
export const makeZanzaBuffConsumesInput = makeConsumeInputFactory({ consumesFieldName: 'zanzaBuff' });

export const DraughtOfTheSands = makeBooleanMiscConsumeInput({
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: player.getClass() === Class.ClassHunter ? 235497 : 235825, minLevel: 55 }]),
	fieldName: 'draughtOfTheSands',
	showWhen: player => player.getClass() === Class.ClassHunter || player.getClass() === Class.ClassRogue,
});

export const Catnip = makeBooleanMiscConsumeInput({
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 213407, minLevel: 20 }]),
	fieldName: 'catnip',
	showWhen: player => player.getClass() === Class.ClassDruid,
});
export const JujuFlurry = makeBooleanMiscConsumeInput({
	actionId: () => ActionId.fromItemId(12450),
	fieldName: 'jujuFlurry',
	showWhen: player => player.getLevel() >= 55,
});
export const RaptorPunch = makeBooleanMiscConsumeInput({
	actionId: () => ActionId.fromItemId(5342),
	fieldName: 'raptorPunch',
});
export const elixirOfCoalescedRegret = makeBooleanMiscConsumeInput({ actionId: () => ActionId.fromItemId(210708), fieldName: 'elixirOfCoalescedRegret' });
export const BoglingRoot = makeBooleanMiscConsumeInput({ actionId: () => ActionId.fromItemId(5206), fieldName: 'boglingRoot' });

export const MISC_OFFENSIVE_CONSUMES_CONFIG: PickerStatOptions[] = [
	{ config: DraughtOfTheSands, picker: IconPicker, stats: [] },
	{ config: Catnip, picker: IconPicker, stats: [] },
	{ config: JujuFlurry, picker: IconPicker, stats: [Stat.StatAttackPower] },
	{ config: RaptorPunch, picker: IconPicker, stats: [Stat.StatIntellect] },
	{ config: elixirOfCoalescedRegret, picker: IconPicker, stats: [] },
	{ config: BoglingRoot, picker: IconPicker, stats: [Stat.StatAttackPower] },
];

export const makeMiscOffensiveConsumesInput = makeMultiIconConsumesInputFactory({
	direction: IconPickerDirection.Vertical,
	tooltip: 'Misc Offensive',
});

export const GreaterMarkOfTheDawn = makeBooleanMiscConsumeInput({
	actionId: () => ActionId.fromItemId(23196),
	fieldName: 'greaterMarkOfTheDawn',
	showWhen: player => player.getLevel() >= 60,
});
export const JujuEmber = makeBooleanMiscConsumeInput({
	actionId: () => ActionId.fromItemId(12455),
	fieldName: 'jujuEmber',
	showWhen: player => player.getLevel() >= 55,
});
export const JujuChill = makeBooleanMiscConsumeInput({
	actionId: () => ActionId.fromItemId(12457),
	fieldName: 'jujuChill',
	showWhen: player => player.getLevel() >= 55,
});
export const JujuEscape = makeBooleanMiscConsumeInput({
	actionId: () => ActionId.fromItemId(12459),
	fieldName: 'jujuEscape',
	showWhen: player => player.getLevel() >= 55,
});

export const MISC_DEFENSIVE_CONSUMES_CONFIG: PickerStatOptions[] = [
	{ config: GreaterMarkOfTheDawn, picker: IconPicker, stats: [] },
	{ config: JujuEmber, picker: IconPicker, stats: [] },
	{ config: JujuChill, picker: IconPicker, stats: [] },
	{ config: JujuEscape, picker: IconPicker, stats: [Stat.StatDodge] },
];

export const makeMiscDefensiveConsumesInput = makeMultiIconConsumesInputFactory({
	direction: IconPickerDirection.Vertical,
	tooltip: 'Misc Defensive',
});

export const MageScrollArcaneRecovery: ConsumableInputConfig<MageScroll> = {
	actionId: () => ActionId.fromItemId(211953),
	value: MageScroll.MageScrollArcaneRecovery,
	showWhen: player => player.isClass(Class.ClassMage),
};
export const MageScrollArcaneAccuracy: ConsumableInputConfig<MageScroll> = {
	actionId: () => ActionId.fromItemId(211954),
	value: MageScroll.MageScrollArcaneAccuracy,
	showWhen: player => player.isClass(Class.ClassMage),
};
export const MageScrollArcanePower: ConsumableInputConfig<MageScroll> = {
	actionId: () => ActionId.fromItemId(211957),
	value: MageScroll.MageScrollArcanePower,
	showWhen: player => player.isClass(Class.ClassMage),
};
export const MageScrollFireProtection: ConsumableInputConfig<MageScroll> = {
	actionId: () => ActionId.fromItemId(211955),
	value: MageScroll.MageScrollFireProtection,
	showWhen: player => player.isClass(Class.ClassMage),
};
export const MageScrollFrostProtection: ConsumableInputConfig<MageScroll> = {
	actionId: () => ActionId.fromItemId(211956),
	value: MageScroll.MageScrollFrostProtection,
	showWhen: player => player.isClass(Class.ClassMage),
};

export const MAGE_SCROLL_CONSUMES_CONFIG: ConsumableStatOption<MageScroll>[] = [
	{ config: MageScrollArcaneRecovery, stats: [] },
	{ config: MageScrollArcaneAccuracy, stats: [] },
	{ config: MageScrollArcanePower, stats: [] },
	{ config: MageScrollFireProtection, stats: [] },
	{ config: MageScrollFrostProtection, stats: [] },
];
export const makeMageScrollsInput = makeConsumeInputFactory({ consumesFieldName: 'mageScroll' });

///////////////////////////////////////////////////////////////////////////
//                                 PET
///////////////////////////////////////////////////////////////////////////

export const PetAttackPowerConsumable = makeEnumConsumeInput({
	direction: IconPickerDirection.Vertical,
	values: [
		{ value: 0, tooltip: 'None' },
		{ actionId: () => ActionId.fromItemId(12460), value: 1, showWhen: player => player.getLevel() >= 55 },
	],
	fieldName: 'petAttackPowerConsumable',
});

export const PetAgilityConsumable = makeEnumConsumeInput({
	direction: IconPickerDirection.Vertical,
	values: [
		{ value: 0, tooltip: 'None' },
		{ actionId: () => ActionId.fromItemId(10309), value: 1, showWhen: player => player.getLevel() >= 55 },
		{ actionId: () => ActionId.fromItemId(4425), value: 2, showWhen: player => player.getLevel() >= 40 },
		{ actionId: () => ActionId.fromItemId(1477), value: 3, showWhen: player => player.getLevel() >= 25 },
		{ actionId: () => ActionId.fromItemId(3012), value: 4, showWhen: player => player.getLevel() >= 10 },
	],
	fieldName: 'petAgilityConsumable',
});

export const PetStrengthConsumable = makeEnumConsumeInput({
	direction: IconPickerDirection.Vertical,
	values: [
		{ value: 0, tooltip: 'None' },
		{ actionId: () => ActionId.fromItemId(12451), value: 1, showWhen: player => player.getLevel() >= 55 },
		{ actionId: () => ActionId.fromItemId(10310), value: 2, showWhen: player => player.getLevel() >= 55 },
		{ actionId: () => ActionId.fromItemId(4426), value: 3, showWhen: player => player.getLevel() >= 40 },
		{ actionId: () => ActionId.fromItemId(2289), value: 4, showWhen: player => player.getLevel() >= 25 },
		{ actionId: () => ActionId.fromItemId(954), value: 5, showWhen: player => player.getLevel() >= 10 },
	],
	fieldName: 'petStrengthConsumable',
});

export const JujuFlurryPet = makeBooleanPetMiscConsumeInput({
	actionId: () => ActionId.fromItemId(12450),
	fieldName: 'jujuFlurry',
	showWhen: player => player.getLevel() >= 55,
});

export const MISC_PET_CONSUMES: PickerStatOptions[] = [{ config: JujuFlurryPet, picker: IconPicker, stats: [] }];

export const makeMiscPetConsumesInput = makeMultiIconConsumesInputFactory({
	direction: IconPickerDirection.Vertical,
	tooltip: 'Misc Pet Consumes',
});

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////

export const GreaterHealingPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 1710, minLevel: 21 }]),
	value: Potions.GreaterHealingPotion,
};
export const SuperiorHealingPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 3928, minLevel: 35 }]),
	value: Potions.SuperiorHealingPotion,
};
export const MajorHealingPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 13446, minLevel: 45 }]),
	value: Potions.MajorHealingPotion,
};

export const ManaPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 3827, minLevel: 22 }]),
	value: Potions.ManaPotion,
};
export const GreaterManaPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 6149, minLevel: 31 }]),
	value: Potions.GreaterManaPotion,
};
export const SuperiorManaPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 13443, minLevel: 41 }]),
	value: Potions.SuperiorManaPotion,
};
export const MajorManaPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 13444, minLevel: 49 }]),
	value: Potions.MajorManaPotion,
};

export const MightRagePotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 13442, minLevel: 46 }]),
	value: Potions.MightyRagePotion,
	showWhen: player => player.getClass() == Class.ClassWarrior,
};
export const GreatRagePotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 5633, minLevel: 25 }]),
	value: Potions.GreatRagePotion,
	showWhen: player => player.getClass() == Class.ClassWarrior,
};
export const RagePotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 5631, minLevel: 4 }]),
	value: Potions.RagePotion,
	showWhen: player => player.getClass() == Class.ClassWarrior,
};

export const MagicResistancePotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 9036, minLevel: 32 }]),
	value: Potions.MagicResistancePotion,
};
// TODO: Not yet implemented in the back-end. Missing school shields and shields don't actually absorb damage right now
// export const GreaterArcaneProtectionPotion: ConsumableInputConfig<Potions> = {
// 	actionId: player => player.getMatchingItemActionId([{ id: 13461, minLevel: 48 }]),
// 	value: Potions.GreaterArcaneProtectionPotion,
// };
// export const GreaterFireProtectionPotion: ConsumableInputConfig<Potions> = {
// 	actionId: player => player.getMatchingItemActionId([{ id: 13457, minLevel: 48 }]),
// 	value: Potions.GreaterFireProtectionPotion,
// };
// export const GreaterFrostProtectionPotion: ConsumableInputConfig<Potions> = {
// 	actionId: player => player.getMatchingItemActionId([{ id: 13456, minLevel: 48 }]),
// 	value: Potions.GreaterFrostProtectionPotion,
// };
// export const GreaterHolyProtectionPotion: ConsumableInputConfig<Potions> = {
// 	actionId: player => player.getMatchingItemActionId([{ id: 13460, minLevel: 48 }]),
// 	value: Potions.GreaterHolyProtectionPotion,
// };
// export const GreaterNatureProtectionPotion: ConsumableInputConfig<Potions> = {
// 	actionId: player => player.getMatchingItemActionId([{ id: 13458, minLevel: 48 }]),
// 	value: Potions.GreaterNatureProtectionPotion,
// };
// export const GreaterShadowProtectionPotion: ConsumableInputConfig<Potions> = {
// 	actionId: player => player.getMatchingItemActionId([{ id: 13459, minLevel: 48 }]),
// 	value: Potions.GreaterShadowProtectionPotion,
// };

export const GreaterStoneshieldPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 13455, minLevel: 46 }]),
	value: Potions.GreaterStoneshieldPotion,
};
export const LesserStoneshieldPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 4623, minLevel: 33 }]),
	value: Potions.LesserStoneshieldPotion,
};

export const POTIONS_CONFIG: ConsumableStatOption<Potions>[] = [
	{ config: MajorHealingPotion, stats: [Stat.StatArmor] },
	{ config: SuperiorHealingPotion, stats: [Stat.StatArmor] },
	{ config: GreaterHealingPotion, stats: [Stat.StatArmor] },

	{ config: MajorManaPotion, stats: [Stat.StatIntellect] },
	{ config: SuperiorManaPotion, stats: [Stat.StatIntellect] },
	{ config: GreaterManaPotion, stats: [Stat.StatIntellect] },
	{ config: ManaPotion, stats: [Stat.StatIntellect] },

	{ config: MightRagePotion, stats: [] },
	{ config: GreatRagePotion, stats: [] },
	{ config: RagePotion, stats: [] },

	// { config: MagicResistancePotion, stats: [] },
	{ config: GreaterStoneshieldPotion, stats: [Stat.StatArmor] },
	{ config: LesserStoneshieldPotion, stats: [Stat.StatArmor] },
];

export const makePotionsInput = makeConsumeInputFactory({ consumesFieldName: 'defaultPotion' });

export const MildlyIrradiatedRejuvPotion = makeBooleanConsumeInput({
	actionId: player => player.getMatchingItemActionId([{ id: 215162, minLevel: 35 }]),
	fieldName: 'mildlyIrradiatedRejuvPot',
});

///////////////////////////////////////////////////////////////////////////
//                                 SPELL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Arcane
export const ElixirOfTheMageLord: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 233968, minLevel: 60 }]),
	value: SpellPowerBuff.ElixirOfTheMageLord,
};
export const GreaterArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 13454, minLevel: 46 }]),
	value: SpellPowerBuff.GreaterArcaneElixir,
};
export const ArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 9155, minLevel: 37 }]),
	value: SpellPowerBuff.ArcaneElixir,
};
export const LesserArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 217398, minLevel: 28 }]),
	value: SpellPowerBuff.LesserArcaneElixir,
};

export const SPELL_POWER_CONFIG: ConsumableStatOption<SpellPowerBuff>[] = [
	{ config: ElixirOfTheMageLord, stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: GreaterArcaneElixir, stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: ArcaneElixir, stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: LesserArcaneElixir, stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
];

export const makeSpellPowerConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'spellPowerBuff' });

// Fire
// Original lvl 40 not obtainable in Phase 3
export const ElixirOfGreaterFirepower: ConsumableInputConfig<FirePowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 21546, minLevel: 51 }]),
	value: FirePowerBuff.ElixirOfGreaterFirepower,
};
export const ElixirOfFirepower: ConsumableInputConfig<FirePowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 6373, minLevel: 18 }]),
	value: FirePowerBuff.ElixirOfFirepower,
};

export const FIRE_POWER_CONFIG: ConsumableStatOption<FirePowerBuff>[] = [
	{ config: ElixirOfGreaterFirepower, stats: [Stat.StatFirePower] },
	{ config: ElixirOfFirepower, stats: [Stat.StatFirePower] },
];

export const makeFirePowerConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'firePowerBuff' });

// Frost
export const ElixirOfFrostPower: ConsumableInputConfig<FrostPowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 17708, minLevel: 40 }]),
	value: FrostPowerBuff.ElixirOfFrostPower,
};

export const FROST_POWER_CONFIG: ConsumableStatOption<FrostPowerBuff>[] = [{ config: ElixirOfFrostPower, stats: [Stat.StatFrostPower] }];

export const makeFrostPowerConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'frostPowerBuff' });

// Shadow
export const ElixirOfShadowPower: ConsumableInputConfig<ShadowPowerBuff> = {
	actionId: player => player.getMatchingItemActionId([{ id: 9264, minLevel: 40 }]),
	value: ShadowPowerBuff.ElixirOfShadowPower,
};

export const SHADOW_POWER_CONFIG: ConsumableStatOption<ShadowPowerBuff>[] = [{ config: ElixirOfShadowPower, stats: [Stat.StatShadowPower] }];

export const makeShadowPowerConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'shadowPowerBuff' });

// MP5
export const MagebloodPotion: ConsumableInputConfig<ManaRegenElixir> = {
	actionId: player => player.getMatchingItemActionId([{ id: 20007, minLevel: 40 }]),
	value: ManaRegenElixir.MagebloodPotion,
};

export const MP5_CONFIG: ConsumableStatOption<ManaRegenElixir>[] = [{ config: MagebloodPotion, stats: [Stat.StatMP5] }];

export const makeMp5ConsumeInput = makeConsumeInputFactory({ consumesFieldName: 'manaRegenElixir' });

///////////////////////////////////////////////////////////////////////////
//                                 Weapon Imbues
///////////////////////////////////////////////////////////////////////////

// Windfury (Buff)
export const Windfury: ConsumableInputConfig<WeaponImbue> = {
	actionId: player =>
		player.getMatchingSpellActionId([
			{ id: 8512, minLevel: 32, maxLevel: 41 },
			{ id: 10613, minLevel: 42, maxLevel: 51 },
			{ id: 10614, minLevel: 52 },
		]),
	value: WeaponImbue.Windfury,
};

// Wild Strikes
export const WildStrikes: ConsumableInputConfig<WeaponImbue> = {
	actionId: () => ActionId.fromSpellId(407975),
	value: WeaponImbue.WildStrikes,
};

// Other Imbues

// Wizard Oils
export const BlessedWizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 238234, minLevel: 50 }]),
		value: WeaponImbue.BlessedWizardOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || weapon.item.weaponType !== WeaponType.WeaponTypeOffHand;
		},
	};
};
export const EnchantedRepellent = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 233996, minLevel: 60 }]),
		value: WeaponImbue.EnchantedRepellent,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || weapon.item.weaponType !== WeaponType.WeaponTypeOffHand;
		},
	};
};
export const BrilliantWizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20749, minLevel: 45 }]),
		value: WeaponImbue.BrilliantWizardOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const WizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20750, minLevel: 40 }]),
		value: WeaponImbue.WizardOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const LesserWizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20746, minLevel: 30 }]),
		value: WeaponImbue.LesserWizardOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const MinorWizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20744, minLevel: 5 }]),
		value: WeaponImbue.MinorWizardOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};

// Mana Oils
// Original lvl 45 but not obtainable in Phase 3
export const BrilliantManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20748, minLevel: 45 }]),
		value: WeaponImbue.BrilliantManaOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const LesserManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20747, minLevel: 40 }]),
		value: WeaponImbue.LesserManaOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const MinorManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20745, minLevel: 20 }]),
		value: WeaponImbue.MinorManaOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const BlackfathomManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 211848, minLevel: 25 }]),
		value: WeaponImbue.BlackfathomManaOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};

// Sharpening Stones
export const ConsecratedSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 238241, minLevel: 50 }]),
		value: WeaponImbue.ConsecratedSharpeningStone,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || weapon.item.weaponType !== WeaponType.WeaponTypeOffHand;
		},
	};
};
export const WeightedConsecratedSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 237810, minLevel: 50 }]),
		value: WeaponImbue.WeightedConsecratedSharpeningStone,
		showWhen: player => {
			return player.getEquippedItem(slot)?.item.handType === HandType.HandTypeTwoHand;
		},
	};
};
export const ElementalSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 18262, minLevel: 50 }]),
		value: WeaponImbue.ElementalSharpeningStone,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const DenseSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 12404, minLevel: 35 }]),
		value: WeaponImbue.DenseSharpeningStone,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isSharpWeaponType(weapon.item.weaponType);
		},
	};
};
export const SolidSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 7964, minLevel: 35 }]),
		value: WeaponImbue.SolidSharpeningStone,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isSharpWeaponType(weapon.item.weaponType);
		},
	};
};
export const BlackfathomSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: () => ActionId.fromItemId(211845),
		value: WeaponImbue.BlackfathomSharpeningStone,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isSharpWeaponType(weapon.item.weaponType);
		},
	};
};

// Weightstones
export const DenseWeightstone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 12643, minLevel: 35 }]),
		value: WeaponImbue.DenseWeightstone,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isBluntWeaponType(weapon.item.weaponType);
		},
	};
};
export const SolidWeightstone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 7965, minLevel: 35 }]),
		value: WeaponImbue.SolidWeightstone,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isBluntWeaponType(weapon.item.weaponType);
		},
	};
};

// Spell Oils
export const ShadowOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 3824, minLevel: 25 }]),
		value: WeaponImbue.ShadowOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};
export const FrostOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 3829, minLevel: 40 }]),
		value: WeaponImbue.FrostOil,
		showWhen: player => {
			const weapon = player.getEquippedItem(slot);
			return !weapon || isWeapon(weapon.item.weaponType);
		},
	};
};

export const ConductiveShieldCoating = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 228980, minLevel: 40 }]),
		value: WeaponImbue.ConductiveShieldCoating,
		showWhen: player =>
			slot === ItemSlot.ItemSlotOffHand && player.getEquippedItem(ItemSlot.ItemSlotOffHand)?.item?.weaponType === WeaponType.WeaponTypeShield,
	};
};
export const MagnificentTrollshine = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 232611, minLevel: 45 }]),
		value: WeaponImbue.MagnificentTrollshine,
		showWhen: player =>
			slot === ItemSlot.ItemSlotOffHand && player.getEquippedItem(ItemSlot.ItemSlotOffHand)?.item?.weaponType === WeaponType.WeaponTypeShield,
	};
};

const SHAMAN_IMBUES = (slot: ItemSlot): ConsumableStatOption<WeaponImbue>[] => [
	{ config: RockbiterWeaponImbue(slot), stats: [] },
	{ config: FlametongueWeaponImbue(slot), stats: [] },
	{ config: FrostbrandWeaponImbue(slot), stats: [] },
	{ config: WindfuryWeaponImbue(slot), stats: [] },
];

const ROGUE_IMBUES: ConsumableStatOption<WeaponImbue>[] = [
	{ config: InstantPoisonWeaponImbue, stats: [] },
	{ config: DeadlyPoisonWeaponImbue, stats: [] },
	{ config: WoundPoisonWeaponImbue, stats: [] },
	{ config: OccultPoisonWeaponImbue, stats: [] },
	{ config: SebaciousPoisonWeaponImbue, stats: [] },
	{ config: AtrophicPoisonWeaponImbue, stats: [] },
	{ config: NumbingPoisonWeaponImbue, stats: [] },
];

const CONSUMABLES_IMBUES = (slot: ItemSlot): ConsumableStatOption<WeaponImbue>[] => [
	{ config: BlessedWizardOil(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: EnchantedRepellent(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: MagnificentTrollshine(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: ConductiveShieldCoating(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: BrilliantWizardOil(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: WizardOil(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: LesserWizardOil(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: MinorWizardOil(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage] },

	{ config: BrilliantManaOil(slot), stats: [Stat.StatHealingPower, Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: LesserManaOil(slot), stats: [Stat.StatHealingPower, Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: MinorManaOil(slot), stats: [Stat.StatHealingPower, Stat.StatSpellPower, Stat.StatSpellDamage] },
	{ config: BlackfathomManaOil(slot), stats: [Stat.StatSpellPower, Stat.StatSpellDamage, Stat.StatMP5] },

	{ config: WeightedConsecratedSharpeningStone(slot), stats: [Stat.StatAttackPower] },
	{ config: ConsecratedSharpeningStone(slot), stats: [Stat.StatAttackPower] },
	{ config: ElementalSharpeningStone(slot), stats: [Stat.StatAttackPower] },
	{ config: DenseSharpeningStone(slot), stats: [Stat.StatAttackPower] },
	{ config: SolidSharpeningStone(slot), stats: [Stat.StatAttackPower] },
	{ config: BlackfathomSharpeningStone(slot), stats: [Stat.StatMeleeHit] },

	{ config: DenseWeightstone(slot), stats: [Stat.StatAttackPower] },
	{ config: SolidWeightstone(slot), stats: [Stat.StatAttackPower] },

	{ config: ShadowOil(slot), stats: [Stat.StatAttackPower] },
	{ config: FrostOil(slot), stats: [Stat.StatAttackPower] },
];

export const WEAPON_IMBUES_OH_CONFIG: ConsumableStatOption<WeaponImbue>[] = [
	...ROGUE_IMBUES,
	...SHAMAN_IMBUES(ItemSlot.ItemSlotOffHand),
	...CONSUMABLES_IMBUES(ItemSlot.ItemSlotOffHand),
];

export const WEAPON_IMBUES_MH_CONFIG: ConsumableStatOption<WeaponImbue>[] = [
	...ROGUE_IMBUES,
	...SHAMAN_IMBUES(ItemSlot.ItemSlotMainHand),
	{ config: Windfury, stats: [Stat.StatMeleeHit] },
	{ config: WildStrikes, stats: [Stat.StatMeleeHit] },
	...CONSUMABLES_IMBUES(ItemSlot.ItemSlotMainHand),
];

export const makeMainHandImbuesInput = makeConsumeInputFactory({
	consumesFieldName: 'mainHandImbue',
	showWhen: player => !!player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand),
});
export const makeOffHandImbuesInput = makeConsumeInputFactory({
	consumesFieldName: 'offHandImbue',
	showWhen: player => !!player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand),
});
