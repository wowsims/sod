import { Player } from "../../player";
import {
	AgilityElixir,
	Conjured,
	Consumes,
	Explosive,
	FirePowerBuff,
	Flask,
	Food,
	FrostPowerBuff,
	ItemSlot,
	Potions,
	Profession,
	ShadowPowerBuff,
	Spec,
	SpellPowerBuff,
	Stat,
	StrengthBuff,
	WeaponImbue } from "../../proto/common";
import { ActionId } from "../../proto_utils/action_id";
import { EventID, TypedEvent } from "../../typed_event";

import { IconEnumValueConfig } from "../icon_enum_picker";
import { makeBooleanConsumeInput } from "../icon_inputs";

import { ActionInputConfig, ItemStatOption } from "./stat_options";

import * as InputHelpers from '../input_helpers';

export interface ConsumableInputConfig<T> extends ActionInputConfig<T> {
	value: T,
}

export interface ConsumableStatOption<T> extends ItemStatOption<T> {
	config: ConsumableInputConfig<T>
}

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes,
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventactionId: EventID, player: Player<any>, newValue: T) => void
	showWhen?: (player: Player<any>) => boolean
}

function makeConsumeInputFactory<T extends number>(args: ConsumeInputFactoryArgs<T>): (options: ConsumableStatOption<T>[], tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
	return (options: ConsumableStatOption<T>[], tooltip?: string) => {
		return {
			type: 'iconEnum',
			tooltip: tooltip,
			numColumns: options.length > 5 ? 2 : 1,
			values: [
				{ value: 0 } as unknown as IconEnumValueConfig<Player<any>, T>,
			].concat(options.map(option => {
				return {
					actionId: option.config.actionId,
					value: option.config.value,
					showWhen: (player: Player<any>) =>
						(!option.config.showWhen || option.config.showWhen(player)) &&
						(option.config.faction || player.getFaction()) == player.getFaction(),
				} as IconEnumValueConfig<Player<any>, T>;
			})),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,	
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter, player.gearChangeEmitter]),
			showWhen: (player: Player<any>) => !args.showWhen || args.showWhen(player),
			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();

				if (newConsumes[args.consumesFieldName] === newValue){
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
	showWhen: (player: Player<any>) => player.getGear().hasTrinket(4381)
};
export const ConjuredDemonicRune: ConsumableInputConfig<Conjured> = { 
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		[12662, 40]
	]),
	value: Conjured.ConjuredDemonicRune
};

export const CONJURED_CONFIG: ConsumableStatOption<Conjured>[] = [
	{ config: ConjuredMinorRecombobulator, 	stats: [Stat.StatIntellect] },
	{ config: ConjuredDemonicRune, 					stats: [Stat.StatIntellect] },
]

export const makeConjuredInput = makeConsumeInputFactory({consumesFieldName: 'defaultConjured'});

///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

export const ExplosiveDenseDynamite: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[18641, 50]]),
	value: Explosive.ExplosiveDenseDynamite,
};
export const ExplosiveThoriumGrenade: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[15993, 50]]),
	value: Explosive.ExplosiveThoriumGrenade,
};

// TODO: Add more SoD explosives + implement on back-end
export const EXPLOSIVES_CONFIG: ConsumableStatOption<Explosive>[] = [
	{ config: ExplosiveDenseDynamite, 	stats: [] },
	{ config: ExplosiveThoriumGrenade, 	stats: [] },
];

export const makeExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'fillerExplosive',
	showWhen: (player) => !!player.getProfessions().find(p => p == Profession.Engineering),
});

export const Sapper = makeBooleanConsumeInput({
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[10646, 40]]),
	fieldName: 'sapper',
});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS
///////////////////////////////////////////////////////////////////////////

export const FlaskOfTheTitans: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13510, 50]]),
	value: Flask.FlaskOfTheTitans,
};
export const FlaskOfDistilledWisdom: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13511, 50]]),
	value: Flask.FlaskOfDistilledWisdom,
};
export const FlaskOfSupremePower: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13512, 50]]),
	value: Flask.FlaskOfSupremePower,
};
export const FlaskOfChromaticResistance: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13513, 50]]),
	value: Flask.FlaskOfChromaticResistance,
};

export const FLASKS_CONFIG: ConsumableStatOption<Flask>[] = [
	{ config: FlaskOfTheTitans, 					stats: [Stat.StatStamina] },
	{ config: FlaskOfDistilledWisdom, 		stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: FlaskOfSupremePower, 				stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: FlaskOfChromaticResistance, stats: [Stat.StatStamina] },
];

export const makeFlasksInput = makeConsumeInputFactory({consumesFieldName: 'flask'});

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const DirgesKickChimaerokChops: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[21023, 55]]),
	value: Food.FoodDirgesKickChimaerokChops,
};
export const SmokedDesertDumpling: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[20452, 45]]),
	value: Food.FoodSmokedDesertDumpling,
};
export const RunnTumTuberSurprise: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[18254, 45]]),
	value: Food.FoodRunnTumTuberSurprise,
};
export const BlessSunfruit: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13810, 45]]),
	value: Food.FoodBlessSunfruit,
};
export const BlessedSunfruitJuice: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13813, 45]]),
	value: Food.FoodBlessedSunfruitJuice,
};
export const TenderWolfSteak: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[22480, 40]]),
	value: Food.FoodTenderWolfSteak
};
export const NightfinSoup: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13931, 35]]),
	value: Food.FoodNightfinSoup,
};
export const GrilledSquid: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13928, 35]]),
	value: Food.FoodGrilledSquid,
};
export const HotWolfRibs: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13851, 25]]),
	value: Food.FoodHotWolfRibs,
};
export const SmokedSagefish: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[21072, 10]]),
	value: Food.FoodSmokedSagefish,
};

export const FOOD_CONFIG: ConsumableStatOption<Food>[] = [
	{ config: DirgesKickChimaerokChops, stats: [Stat.StatStamina] },
	{ config: SmokedDesertDumpling, 		stats: [Stat.StatStrength] },
	{ config: RunnTumTuberSurprise, 		stats: [Stat.StatIntellect] },
	{ config: BlessSunfruit, 						stats: [Stat.StatStrength] },
	{ config: BlessedSunfruitJuice, 		stats: [Stat.StatSpirit] },
	{ config: TenderWolfSteak, 					stats: [Stat.StatStamina, Stat.StatSpirit] },
	{ config: NightfinSoup, 						stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: GrilledSquid, 						stats: [Stat.StatAgility] },
	{ config: HotWolfRibs, 							stats: [Stat.StatSpirit] },
	{ config: SmokedSagefish, 					stats: [Stat.StatMP5] },
];

export const makeFoodInput = makeConsumeInputFactory({consumesFieldName: 'food'});

///////////////////////////////////////////////////////////////////////////
//                                 PHYSICAL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Agility
export const ElixirOfTheMongoose: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[13452, 46]]),
	value: AgilityElixir.ElixirOfTheMongoose,
};
export const ElixirOfGreaterAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[9187, 38]]),
	value: AgilityElixir.ElixirOfGreaterAgility,
};
export const ElixirOfLesserAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([[3390, 18]]),
	value: AgilityElixir.ElixirOfLesserAgility,
};
export const ScrollOfAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player) => player.getMatchingItemActionId([
		[3012, 10, 24],
		[1477, 25, 39],
		[4425, 40, 54],
		[10309, 55],
	]),
	value: AgilityElixir.ScrollOfAgility,
};

export const AGILITY_CONSUMES_CONFIG: ConsumableStatOption<AgilityElixir>[] = [
	{ config: ElixirOfTheMongoose, 		stats: [Stat.StatAgility] },
	{ config: ElixirOfGreaterAgility, stats: [Stat.StatAgility] },
	{ config: ElixirOfLesserAgility, 	stats: [Stat.StatAgility] },
	{ config: ScrollOfAgility, 				stats: [Stat.StatAgility] },
];

export const makeAgilityConsumeInput = makeConsumeInputFactory({consumesFieldName: 'agilityElixir'});

// Strength
export const JujuPower: ConsumableInputConfig<StrengthBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[12451, 46]]),
	value: StrengthBuff.JujuPower,
};
export const ElixirOfGiants: ConsumableInputConfig<StrengthBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[9206, 46]]),
	value: StrengthBuff.ElixirOfGiants,
};
export const ElixirOfOgresStrength: ConsumableInputConfig<StrengthBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[3391, 20]]),
	value: StrengthBuff.ElixirOfOgresStrength,
};
export const ScrollOfStrength: ConsumableInputConfig<StrengthBuff> = {
	actionId: () => ActionId.fromItemId(10310),
	value: StrengthBuff.ScrollOfStrength,
};

export const STRENGTH_CONSUMES_CONFIG: ConsumableStatOption<StrengthBuff>[] = [
	{ config: JujuPower, 							stats: [Stat.StatStrength] },
	{ config: ElixirOfGiants, 				stats: [Stat.StatStrength] },
	{ config: ElixirOfOgresStrength,	stats: [Stat.StatStrength] },
	{ config: ScrollOfStrength, 			stats: [Stat.StatStrength] },
];

export const makeStrengthConsumeInput = makeConsumeInputFactory({consumesFieldName: 'strengthBuff'});

// Other
export const BoglingRootBuff = makeBooleanConsumeInput({actionId: () => ActionId.fromItemId(5206), fieldName: 'boglingRoot'});

///////////////////////////////////////////////////////////////////////////
//                                 PET
///////////////////////////////////////////////////////////////////////////

// export const PetScrollOfAgilityV = makeBooleanConsumeInput({actionId: ActionId.fromItemId(27498), fieldName: 'petScrollOfAgility', minLevel: 5});
// export const PetScrollOfStrengthV = makeBooleanConsumeInput({actionId: ActionId.fromItemId(27503), fieldName: 'petScrollOfStrength', minLevel: 5});

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////

export const LesserManaPotion: ConsumableInputConfig<Potions> = {
	actionId: () => ActionId.fromItemId(3385),
	value: Potions.LesserManaPotion,
};
export const ManaPotion: ConsumableInputConfig<Potions> = {
	actionId: () => ActionId.fromItemId(3827),
	value: Potions.ManaPotion,
};

export const POTIONS_CONFIG: ConsumableStatOption<Potions>[] = [
	{ config: LesserManaPotion, stats: [Stat.StatIntellect] },
	{ config: ManaPotion, 			stats: [Stat.StatIntellect] },
];

export const makePotionsInput = makeConsumeInputFactory({consumesFieldName: 'defaultPotion'});

///////////////////////////////////////////////////////////////////////////
//                                 SPELL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Arcane
export const GreaterArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[13454, 46]]),
	value: SpellPowerBuff.GreaterArcaneElixir,
};
export const ArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[9155, 37]]),
	value: SpellPowerBuff.ArcaneElixir,
};

export const SPELL_POWER_CONFIG: ConsumableStatOption<SpellPowerBuff>[] = [
	{ config: GreaterArcaneElixir, stats: [Stat.StatSpellPower] },
	{ config: ArcaneElixir, 			 stats: [Stat.StatSpellPower] },
];

export const makeSpellPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'spellPowerBuff'})

// Fire
export const ElixirOfGreaterFirepower: ConsumableInputConfig<FirePowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[21546, 40]]),
	value: FirePowerBuff.ElixirOfGreaterFirepower,
};
export const ElixirOfFirepower: ConsumableInputConfig<FirePowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[6373, 18]]),
	value: FirePowerBuff.ElixirOfFirepower,
};

export const FIRE_POWER_CONFIG: ConsumableStatOption<FirePowerBuff>[] = [
	{ config: ElixirOfGreaterFirepower, stats: [Stat.StatFirePower] },
	{ config: ElixirOfFirepower, 				stats: [Stat.StatFirePower] },
];

export const makeFirePowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'firePowerBuff'})

// Frost
export const ElixirOfFrostPower: ConsumableInputConfig<FrostPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[17708, 40]]),
	value: FrostPowerBuff.ElixirOfFrostPower,
};

export const FROST_POWER_CONFIG: ConsumableStatOption<FrostPowerBuff>[] = [
	{ config: ElixirOfFrostPower, stats: [Stat.StatFrostPower] },
];

export const makeFrostPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'frostPowerBuff'})

// Shadow
export const ElixirOfShadowPower: ConsumableInputConfig<ShadowPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([[9264, 40]]),
	value: ShadowPowerBuff.ElixirOfShadowPower,
};

export const SHADOW_POWER_CONFIG: ConsumableStatOption<ShadowPowerBuff>[] = [
	{ config: ElixirOfShadowPower, stats: [Stat.StatShadowPower] },
];

export const makeshadowPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'shadowPowerBuff'})

///////////////////////////////////////////////////////////////////////////
//                                 Weapon Imbues
///////////////////////////////////////////////////////////////////////////

export const ElementalSharpeningStone: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([[18262, 50]]),
	value: WeaponImbue.ElementalSharpeningStone,
};
export const BrillianWizardOil: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([[20749, 45]]),
	value: WeaponImbue.BrillianWizardOil,
};
export const BrilliantManaOil: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([[20748, 45]]),
	value: WeaponImbue.BrilliantManaOil,
};
export const DenseSharpeningStone: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([[12404, 35]]),
	value: WeaponImbue.DenseSharpeningStone,
};
export const BlackfathomManaOil: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([[211848, 25]]),
	value: WeaponImbue.BlackfathomManaOil,
};
export const BlackfathomSharpeningStone: ConsumableInputConfig<WeaponImbue> = {
	actionId: () => ActionId.fromItemId(211845),
	value: WeaponImbue.BlackfathomSharpeningStone,
};
export const WildStrikes: ConsumableInputConfig<WeaponImbue> = {
	actionId: () => ActionId.fromSpellId(407975),
	value: WeaponImbue.WildStrikes,
};

export const WEAPON_IMBUES_OH_CONFIG: ConsumableStatOption<WeaponImbue>[] = [
	{ config: ElementalSharpeningStone, stats: [Stat.StatAttackPower] },
	{ config: BrillianWizardOil, stats: [Stat.StatSpellPower] },
	{ config: BrilliantManaOil, stats: [Stat.StatHealing, Stat.StatSpellPower] },
	{ config: DenseSharpeningStone, stats: [Stat.StatAttackPower] },
	{ config: BlackfathomManaOil, stats: [Stat.StatSpellPower, Stat.StatMP5] },
	{ config: BlackfathomSharpeningStone, stats: [Stat.StatMeleeHit] },
];

export const WEAPON_IMBUES_MH_CONFIG: ConsumableStatOption<WeaponImbue>[] = [
	...WEAPON_IMBUES_OH_CONFIG,
	{ config: WildStrikes, stats: [Stat.StatMeleeHit] },
];

export const makeMainHandImbuesInput = makeConsumeInputFactory({
	consumesFieldName: 'mainHandImbue',
	showWhen: (player) => !!player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand),
});
export const makeOffHandImbuesInput = makeConsumeInputFactory({
	consumesFieldName: 'offHandImbue',
	showWhen: (player) => !!player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand),
});
