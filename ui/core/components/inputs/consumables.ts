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
	HealthElixir,
	ItemSlot,
	ManaRegenElixir,
	Potions,
	Profession,
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
import { makeBooleanConsumeInput, makeBooleanMiscConsumeInput, makeEnumConsumeInput } from '../icon_inputs';
import { IconPickerDirection } from '../icon_picker';
import * as InputHelpers from '../input_helpers';
import { MultiIconPicker } from '../multi_icon_picker';
import { DeadlyPoisonWeaponImbue, InstantPoisonWeaponImbue, WoundPoisonWeaponImbue } from './rogue_imbues';
import { FlametongueWeaponImbue, FrostbrandWeaponImbue, RockbiterWeaponImbue, WindfuryWeaponImbue } from './shaman_imbues';
import { ActionInputConfig, ItemStatOption, PickerStatOptions } from './stat_options';

export interface ConsumableInputConfig<T> extends ActionInputConfig<T> {
	value: T;
}

export interface ConsumableStatOption<T> extends ItemStatOption<T> {
	config: ConsumableInputConfig<T>;
}

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes;
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
			numColumns: options.length > 11 ? 4 : options.length > 8 ? 3 : options.length > 5 ? 2 : 1,
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

///////////////////////////////////////////////////////////////////////////
//                                 CONJURED
///////////////////////////////////////////////////////////////////////////

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
	{ config: ConjuredMinorRecombobulator, stats: [Stat.StatIntellect] },
	{ config: ConjuredDemonicRune, stats: [Stat.StatIntellect] },
	{ config: ConjuredRogueThistleTea, stats: [] },
];

export const makeConjuredInput = makeConsumeInputFactory({ consumesFieldName: 'defaultConjured' });

///////////////////////////////////////////////////////////////////////////
//                             ENCHANTING SIGIL
///////////////////////////////////////////////////////////////////////////

export const EnchantedSigilInnovation: ConsumableInputConfig<EnchantedSigil> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 217308, minLevel: 40 }]),
	value: EnchantedSigil.InnovationSigil,
};

export const EnchantedSigilLivingDreams: ConsumableInputConfig<EnchantedSigil> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 221028, minLevel: 50 }]),
	value: EnchantedSigil.LivingDreamsSigil,
};

export const ENCHANTED_SIGIL_CONFIG: ConsumableStatOption<EnchantedSigil>[] = [
	{ config: EnchantedSigilLivingDreams, stats: [] },
	{ config: EnchantedSigilInnovation, stats: [] },
];

export const makeEncanthedSigilInput = makeConsumeInputFactory({
	consumesFieldName: 'enchantedSigil',
	showWhen: player => player.hasProfession(Profession.Enchanting),
});

///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

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
	{ config: ExplosiveEzThroRadiationBomb, stats: [] },
	{ config: ExplosiveHighYieldRadiationBomb, stats: [] },
	{ config: ExplosiveSolidDynamite, stats: [] },
	{ config: ExplosiveDenseDynamite, stats: [] },
	{ config: ExplosiveThoriumGrenade, stats: [] },
	{ config: ExplosiveGoblinLandMine, stats: [] },
];

export const makeExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'fillerExplosive',
	//showWhen: (player) => !!player.getProfessions().find(p => p == Profession.Engineering),
});

export const Sapper = makeBooleanConsumeInput({
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 10646, minLevel: 50 }]),
	fieldName: 'sapper',
	showWhen: player => player.hasProfession(Profession.Engineering),
});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS
///////////////////////////////////////////////////////////////////////////

// Original lvl 50 not obtainable in Phase 3
export const FlaskOfTheTitans: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13510, minLevel: 51 }]),
	value: Flask.FlaskOfTheTitans,
};
// Original lvl 50 not obtainable in Phase 3
export const FlaskOfDistilledWisdom: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13511, minLevel: 51 }]),
	value: Flask.FlaskOfDistilledWisdom,
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
	{ config: FlaskOfDistilledWisdom, stats: [Stat.StatIntellect] },
	{ config: FlaskOfSupremePower, stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: FlaskOfChromaticResistance, stats: [Stat.StatStamina] },
	{ config: FlaskOfRestlessDreams, stats: [Stat.StatSpellPower] },
	{ config: FlaskOfEverlastingNightmares, stats: [Stat.StatAttackPower] },
];

export const makeFlasksInput = makeConsumeInputFactory({ consumesFieldName: 'flask' });

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const DirgesKickChimaerokChops: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 21023, minLevel: 55 }]),
	value: Food.FoodDirgesKickChimaerokChops,
};
export const GrilledSquid: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 13928, minLevel: 50 }]),
	value: Food.FoodGrilledSquid,
};
// Original lvl 50 not obtainable in Phase 3
export const SmokedDesertDumpling: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 20452, minLevel: 51 }]),
	value: Food.FoodSmokedDesertDumpling,
};
// Original lvl 45 not obtainable in Phase 3
export const RunnTumTuberSurprise: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 18254, minLevel: 51 }]),
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
	value: Alcohol.AlcoholRumseyRumLight,
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
	// { config: SpiritOfZanza, stats: [Stat.StatStamina, Stat.StatSpirit] },
	{ config: ROIDS, stats: [Stat.StatStrength] },
	{ config: GroundScorpokAssay, stats: [Stat.StatAgility] },
	{ config: LungJuiceCocktail, stats: [Stat.StatStamina] },
	{ config: CerebralCortexCompound, stats: [Stat.StatIntellect] },
	{ config: GizzardGum, stats: [Stat.StatSpirit] },
	{ config: AtalAiMojoOfWar, stats: [Stat.StatAttackPower] },
	{ config: AtalAiMojoOfForbiddenMagic, stats: [Stat.StatSpellPower] },
	{ config: AtalAiMojoOfLife, stats: [Stat.StatHealingPower] },
];
export const makeZanzaBuffConsumesInput = makeConsumeInputFactory({ consumesFieldName: 'zanzaBuff' });

export const MiscConsumesConfig = InputHelpers.makeMultiIconInput(
	[
		makeBooleanMiscConsumeInput({
			actionId: (player: Player<Spec>) => player.getMatchingItemActionId([{ id: 213407, minLevel: 20 }]),
			fieldName: 'catnip',
			showWhen: player => player.getClass() === Class.ClassDruid,
		}),
		makeBooleanMiscConsumeInput({ actionId: () => ActionId.fromItemId(210708), fieldName: 'elixirOfCoalescedRegret' }),
		makeBooleanMiscConsumeInput({ actionId: () => ActionId.fromItemId(5206), fieldName: 'boglingRoot' }),
	],
	'',
	IconPickerDirection.Vertical,
);

export const MISC_CONSUMES_CONFIG: PickerStatOptions[] = [{ config: MiscConsumesConfig, picker: MultiIconPicker, stats: [] }];

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
		{ actionId: () => ActionId.fromItemId(3012), value: 1, showWhen: player => player.getLevel() >= 10 },
		{ actionId: () => ActionId.fromItemId(1477), value: 2, showWhen: player => player.getLevel() >= 25 },
		{ actionId: () => ActionId.fromItemId(4425), value: 3, showWhen: player => player.getLevel() >= 40 },
		{ actionId: () => ActionId.fromItemId(10309), value: 4, showWhen: player => player.getLevel() >= 55 },
	],
	fieldName: 'petAgilityConsumable',
});

export const PetStrengthConsumable = makeEnumConsumeInput({
	direction: IconPickerDirection.Vertical,
	values: [
		{ value: 0, tooltip: 'None' },
		{ actionId: () => ActionId.fromItemId(954), value: 1, showWhen: player => player.getLevel() >= 10 },
		{ actionId: () => ActionId.fromItemId(2289), value: 2, showWhen: player => player.getLevel() >= 25 },
		{ actionId: () => ActionId.fromItemId(4426), value: 3, showWhen: player => player.getLevel() >= 40 },
		{ actionId: () => ActionId.fromItemId(10310), value: 4, showWhen: player => player.getLevel() >= 55 },
		{ actionId: () => ActionId.fromItemId(12451), value: 5, showWhen: player => player.getLevel() >= 55 },
	],
	fieldName: 'petStrengthConsumable',
});

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////

export const LesserManaPotion: ConsumableInputConfig<Potions> = {
	actionId: () => ActionId.fromItemId(3385),
	value: Potions.LesserManaPotion,
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
export const GreaterStoneshieldPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 13455, minLevel: 46 }]),
	value: Potions.GreaterStoneshieldPotion,
};
export const LesserStoneshieldPotion: ConsumableInputConfig<Potions> = {
	actionId: player => player.getMatchingItemActionId([{ id: 4623, minLevel: 33 }]),
	value: Potions.LesserStoneshieldPotion,
};

export const POTIONS_CONFIG: ConsumableStatOption<Potions>[] = [
	{ config: MajorManaPotion, stats: [Stat.StatIntellect] },
	{ config: SuperiorManaPotion, stats: [Stat.StatIntellect] },
	{ config: GreaterManaPotion, stats: [Stat.StatIntellect] },
	{ config: ManaPotion, stats: [Stat.StatIntellect] },
	{ config: LesserManaPotion, stats: [Stat.StatIntellect] },
	{ config: MightRagePotion, stats: [] },
	{ config: GreatRagePotion, stats: [] },
	{ config: RagePotion, stats: [] },
	{ config: GreaterStoneshieldPotion, stats: [Stat.StatArmor] },
	{ config: LesserStoneshieldPotion, stats: [Stat.StatArmor] },
];

export const makePotionsInput = makeConsumeInputFactory({ consumesFieldName: 'defaultPotion' });

export const MildlyIrradiatedRejuvPotion = makeBooleanConsumeInput({
	actionId: player => player.getMatchingItemActionId([{ id: 215162, minLevel: 35 }]),
	fieldName: 'mildlyIrradiatedRejuvPot',
	showWhen: player => player.hasProfession(Profession.Alchemy),
});

///////////////////////////////////////////////////////////////////////////
//                                 SPELL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Arcane
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
	{ config: GreaterArcaneElixir, stats: [Stat.StatSpellPower] },
	{ config: ArcaneElixir, stats: [Stat.StatSpellPower] },
	{ config: LesserArcaneElixir, stats: [Stat.StatSpellPower] },
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
// Original lvl 40 but not obtainable until ZG
export const MagebloodPotion: ConsumableInputConfig<ManaRegenElixir> = {
	actionId: player => player.getMatchingItemActionId([{ id: 20007, minLevel: 61 }]),
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
export const BrillianWizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20749, minLevel: 45 }]),
		value: WeaponImbue.BrillianWizardOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const WizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20750, minLevel: 40 }]),
		value: WeaponImbue.WizardOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const LesserWizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20746, minLevel: 30 }]),
		value: WeaponImbue.LesserWizardOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const MinorWizardOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20744, minLevel: 5 }]),
		value: WeaponImbue.MinorWizardOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

// Mana Oils
// Original lvl 45 but not obtainable in Phase 3
export const BrilliantManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20748, minLevel: 45 }]),
		value: WeaponImbue.BrilliantManaOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const LesserManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20747, minLevel: 40 }]),
		value: WeaponImbue.LesserManaOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const MinorManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 20745, minLevel: 20 }]),
		value: WeaponImbue.MinorManaOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const BlackfathomManaOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 211848, minLevel: 25 }]),
		value: WeaponImbue.BlackfathomManaOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

// Sharpening Stones
export const ElementalSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 18262, minLevel: 50 }]),
		value: WeaponImbue.ElementalSharpeningStone,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const DenseSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 12404, minLevel: 35 }]),
		value: WeaponImbue.DenseSharpeningStone,
		showWhen: player => isSharpWeaponType(player.getEquippedItem(slot)?.item.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const SolidSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 7964, minLevel: 35 }]),
		value: WeaponImbue.SolidSharpeningStone,
		showWhen: player => isSharpWeaponType(player.getEquippedItem(slot)?.item.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const BlackfathomSharpeningStone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: () => ActionId.fromItemId(211845),
		value: WeaponImbue.BlackfathomSharpeningStone,
		showWhen: player => isSharpWeaponType(player.getEquippedItem(slot)?.item.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

// Weightstones
export const DenseWeightstone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 12643, minLevel: 35 }]),
		value: WeaponImbue.DenseWeightstone,
		showWhen: player => isBluntWeaponType(player.getEquippedItem(slot)?.item.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const SolidWeightstone = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 7965, minLevel: 35 }]),
		value: WeaponImbue.SolidWeightstone,
		showWhen: player => isBluntWeaponType(player.getEquippedItem(slot)?.item.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

// Spell Oils
export const ShadowOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 3824, minLevel: 25 }]),
		value: WeaponImbue.ShadowOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};
export const FrostOil = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 3829, minLevel: 40 }]),
		value: WeaponImbue.FrostOil,
		showWhen: player => isWeapon(player.getEquippedItem(slot)?._item?.weaponType ?? WeaponType.WeaponTypeUnknown),
	};
};

export const ConductiveShieldCoating = (slot: ItemSlot): ConsumableInputConfig<WeaponImbue> => {
	return {
		actionId: player => player.getMatchingItemActionId([{ id: 228980, minLevel: 40 }]),
		value: WeaponImbue.ConductiveShieldCoating,
		showWhen: player =>
			slot === ItemSlot.ItemSlotOffHand && player.getEquippedItem(ItemSlot.ItemSlotOffHand)?._item?.weaponType === WeaponType.WeaponTypeShield,
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
	// These are not yet implemented for rogues
	// { config: OccultPoisonWeaponImbue, stats: [] },
	// { config: SebaciousPoisonWeaponImbue, stats: [] },
];

const CONSUMABLES_IMBUES = (slot: ItemSlot): ConsumableStatOption<WeaponImbue>[] => [
	{ config: ConductiveShieldCoating(slot), stats: [Stat.StatSpellPower] },
	// { config: BrillianWizardOil, stats: [Stat.StatSpellPower] },
	{ config: WizardOil(slot), stats: [Stat.StatSpellPower] },
	{ config: LesserWizardOil(slot), stats: [Stat.StatSpellPower] },
	{ config: MinorWizardOil(slot), stats: [Stat.StatSpellPower] },

	// { config: BrilliantManaOil, stats: [Stat.StatHealingPower, Stat.StatSpellPower] },
	{ config: LesserManaOil(slot), stats: [Stat.StatHealingPower, Stat.StatSpellPower] },
	{ config: MinorManaOil(slot), stats: [Stat.StatHealingPower, Stat.StatSpellPower] },
	{ config: BlackfathomManaOil(slot), stats: [Stat.StatSpellPower, Stat.StatMP5] },

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
