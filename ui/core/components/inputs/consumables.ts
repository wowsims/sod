import { Player } from "../../player";
import {
	AgilityElixir,
	Conjured,
	Consumes,
	EnchantedSigil,
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
import { FlametongueWeaponImbue, FrostbrandWeaponImbue, RockbiterWeaponImbue, WindfuryWeaponImbue } from "./shaman_imbues";

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
					showWhen: (player: Player<any>) => !option.config.showWhen || option.config.showWhen(player)
				} as IconEnumValueConfig<Player<any>, T>;
			})),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,	
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter, player.gearChangeEmitter, player.professionChangeEmitter]),
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
		{ id: 12662, minLevel: 40 },
	]),
	value: Conjured.ConjuredDemonicRune
}

export const CONJURED_CONFIG: ConsumableStatOption<Conjured>[] = [
	{ config: ConjuredMinorRecombobulator, 			stats: [Stat.StatIntellect] },
	{ config: ConjuredDemonicRune, 					stats: [Stat.StatIntellect] },
]

export const makeConjuredInput = makeConsumeInputFactory({consumesFieldName: 'defaultConjured'});

///////////////////////////////////////////////////////////////////////////
//                             ENCHANTING SIGIL
///////////////////////////////////////////////////////////////////////////

export const EnchantedSigilInnovation: ConsumableInputConfig<EnchantedSigil> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 217308, minLevel: 40 },
	]),
	value: EnchantedSigil.InnovationSigil,
};

export const ENCHANTEDSIGILCONFIG: ConsumableStatOption<EnchantedSigil>[] = [
	{ config: EnchantedSigilInnovation, stats: [] },
];

export const makeEncanthedSigilInput = makeConsumeInputFactory({
	consumesFieldName: 'enchantedSigil',
	showWhen: (player) => player.hasProfession(Profession.Enchanting),
});

///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

export const ExplosiveSolidDynamite: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 10507, minLevel: 40 },
	]),
	showWhen: (player) => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveSolidDynamite,
};

export const ExplosiveDenseDynamite: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 18641, minLevel: 50 },
	]),
	showWhen: (player) => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveDenseDynamite,
};

export const ExplosiveThoriumGrenade: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 15993, minLevel: 50 },
	]),
	showWhen: (player) => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveThoriumGrenade,
};

export const ExplosiveEzThroRadiationBomb: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 215168, minLevel: 40 },
	]),
	value: Explosive.ExplosiveEzThroRadiationBomb,
};

export const ExplosiveHighYealdRadiationBomb: ConsumableInputConfig<Explosive> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 215127, minLevel: 40 },
	]),
	showWhen: (player) => player.hasProfession(Profession.Engineering),
	value: Explosive.ExplosiveHighYieldRadiationBomb,
};

export const EXPLOSIVES_CONFIG: ConsumableStatOption<Explosive>[] = [
	{ config: ExplosiveEzThroRadiationBomb, 	stats: [] },
	{ config: ExplosiveHighYealdRadiationBomb, 	stats: [] },
	{ config: ExplosiveSolidDynamite, 			stats: [] },
	{ config: ExplosiveDenseDynamite,		 	stats: [] },
	{ config: ExplosiveThoriumGrenade, 			stats: [] },
];

export const makeExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'fillerExplosive',
	//showWhen: (player) => !!player.getProfessions().find(p => p == Profession.Engineering),
});

export const Sapper = makeBooleanConsumeInput({
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 10646, minLevel: 50 },
	]),
	fieldName: 'sapper',
	showWhen: (player) => player.hasProfession(Profession.Engineering),
})

export const makeSapperInput = makeConsumeInputFactory({
	consumesFieldName: 'sapper',
	showWhen: (player) => player.hasProfession(Profession.Engineering),
});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS
///////////////////////////////////////////////////////////////////////////

export const FlaskOfTheTitans: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13510, minLevel: 50 },
	]),
	value: Flask.FlaskOfTheTitans,
};
export const FlaskOfDistilledWisdom: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13511, minLevel: 50 },
	]),
	value: Flask.FlaskOfDistilledWisdom,
};
export const FlaskOfSupremePower: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13512, minLevel: 50 },
	]),
	value: Flask.FlaskOfSupremePower,
};
export const FlaskOfChromaticResistance: ConsumableInputConfig<Flask> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13513, minLevel: 50 },
	]),
	value: Flask.FlaskOfChromaticResistance,
};

export const FLASKS_CONFIG: ConsumableStatOption<Flask>[] = [
	{ config: FlaskOfTheTitans, 			stats: [Stat.StatStamina] },
	{ config: FlaskOfDistilledWisdom, 		stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: FlaskOfSupremePower, 			stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: FlaskOfChromaticResistance, 	stats: [Stat.StatStamina] },
];

export const makeFlasksInput = makeConsumeInputFactory({consumesFieldName: 'flask'});

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const DirgesKickChimaerokChops: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 21023, minLevel: 55 },
	]),
	value: Food.FoodDirgesKickChimaerokChops,
};
export const SmokedDesertDumpling: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 20452, minLevel: 45 },
	]),
	value: Food.FoodSmokedDesertDumpling,
};
export const RunnTumTuberSurprise: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 18254, minLevel: 45 },
	]),
	value: Food.FoodRunnTumTuberSurprise,
};
export const BlessSunfruit: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13810, minLevel: 45 },
	]),
	value: Food.FoodBlessSunfruit,
};
export const BlessedSunfruitJuice: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13813, minLevel: 45 },
	]),
	value: Food.FoodBlessedSunfruitJuice,
};
export const TenderWolfSteak: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 18045, minLevel: 40 },
	]),
	value: Food.FoodTenderWolfSteak
};
export const NightfinSoup: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		// Requires Skill 250
		{ id: 13931, minLevel: 41 },
	]),
	value: Food.FoodNightfinSoup,
};
export const GrilledSquid: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13928, minLevel: 50 },
	]),
	value: Food.FoodGrilledSquid,
};
export const HotWolfRibs: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13851, minLevel: 25 },
	]),
	value: Food.FoodHotWolfRibs,
};
export const SmokedSagefish: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 21072, minLevel: 10 },
	]),
	value: Food.FoodSmokedSagefish,
};
export const SagefishDelight: ConsumableInputConfig<Food> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 21217, minLevel: 30 },
	]),
	value: Food.FoodSagefishDelight,
};

export const FOOD_CONFIG: ConsumableStatOption<Food>[] = [
	{ config: DirgesKickChimaerokChops, stats: [Stat.StatStamina] },
	{ config: SmokedDesertDumpling, 	stats: [Stat.StatStrength] },
	{ config: RunnTumTuberSurprise, 	stats: [Stat.StatIntellect] },
	{ config: BlessSunfruit, 			stats: [Stat.StatStrength] },
	{ config: BlessedSunfruitJuice, 	stats: [Stat.StatSpirit] },
	{ config: TenderWolfSteak, 			stats: [Stat.StatStamina, Stat.StatSpirit] },
	{ config: NightfinSoup, 			stats: [Stat.StatMP5, Stat.StatSpellPower] },
	{ config: GrilledSquid, 			stats: [Stat.StatAgility] },
	{ config: HotWolfRibs, 				stats: [Stat.StatSpirit] },
	{ config: SmokedSagefish, 			stats: [Stat.StatMP5] },
	{ config: SagefishDelight, 			stats: [Stat.StatMP5] },
];

export const makeFoodInput = makeConsumeInputFactory({consumesFieldName: 'food'});

///////////////////////////////////////////////////////////////////////////
//                                 PHYSICAL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Agility
export const ElixirOfTheMongoose: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 13452, minLevel: 46 },
	]),
	value: AgilityElixir.ElixirOfTheMongoose,
};
export const ElixirOfGreaterAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 9187, minLevel: 38 },
	]),
	value: AgilityElixir.ElixirOfGreaterAgility,
};
export const ElixirOfLesserAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player: Player<Spec>) => player.getMatchingItemActionId([
		{ id: 3390, minLevel: 18 },
	]),
	value: AgilityElixir.ElixirOfLesserAgility,
};
export const ScrollOfAgility: ConsumableInputConfig<AgilityElixir> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 3012, 	minLevel: 10, maxLevel: 24 	},
		{ id: 1477, 	minLevel: 25, maxLevel: 39 	},
		{ id: 4425, 	minLevel: 40, maxLevel: 54 	},
		{ id: 10309, 	minLevel: 55 								},
	]),
	value: AgilityElixir.ScrollOfAgility,
};

export const AGILITY_CONSUMES_CONFIG: ConsumableStatOption<AgilityElixir>[] = [
	{ config: ElixirOfTheMongoose, 		stats: [Stat.StatAgility] },
	{ config: ElixirOfGreaterAgility, 	stats: [Stat.StatAgility] },
	{ config: ElixirOfLesserAgility, 	stats: [Stat.StatAgility] },
	{ config: ScrollOfAgility, 			stats: [Stat.StatAgility] },
];

export const makeAgilityConsumeInput = makeConsumeInputFactory({consumesFieldName: 'agilityElixir'});

// Strength
export const JujuPower: ConsumableInputConfig<StrengthBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 12451, minLevel: 46 },
	]),
	value: StrengthBuff.JujuPower,
};
export const ElixirOfGiants: ConsumableInputConfig<StrengthBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 9206, minLevel: 46 },
	]),
	value: StrengthBuff.ElixirOfGiants,
};
export const ElixirOfOgresStrength: ConsumableInputConfig<StrengthBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 3391, minLevel: 20 },
	]),
	value: StrengthBuff.ElixirOfOgresStrength,
};
export const ScrollOfStrength: ConsumableInputConfig<StrengthBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 954, 	minLevel: 10, maxLevel: 24 	},
		{ id: 2289, 	minLevel: 25, maxLevel: 39 	},
		{ id: 4426, 	minLevel: 40, maxLevel: 54 	},
		{ id: 10310, 	minLevel: 55 								},
	]),
	value: StrengthBuff.ScrollOfStrength,
};

export const STRENGTH_CONSUMES_CONFIG: ConsumableStatOption<StrengthBuff>[] = [
	{ config: JujuPower,				stats: [Stat.StatStrength] },
	{ config: ElixirOfGiants,			stats: [Stat.StatStrength] },
	{ config: ElixirOfOgresStrength,	stats: [Stat.StatStrength] },
	{ config: ScrollOfStrength, 		stats: [Stat.StatStrength] },
];

export const makeStrengthConsumeInput = makeConsumeInputFactory({consumesFieldName: 'strengthBuff'});

// Other
export const BoglingRootBuff = makeBooleanConsumeInput({actionId: () => ActionId.fromItemId(5206), fieldName: 'boglingRoot'});

///////////////////////////////////////////////////////////////////////////
//                                 PET
///////////////////////////////////////////////////////////////////////////

// export const PetScrollOfAgilityV = makeBooleanConsumeInput({actionId: () => ActionId.fromItemId(27498), fieldName: 'petScrollOfAgility', minLevel: 5});
// export const PetScrollOfStrengthV = makeBooleanConsumeInput({actionId: () => ActionId.fromItemId(27503), fieldName: 'petScrollOfStrength', minLevel: 5});

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
	minLevel: 22,
};
export const GreaterManaPotion: ConsumableInputConfig<Potions> = {
	actionId: () => ActionId.fromItemId(6149),
	value: Potions.GreaterManaPotion,
	minLevel: 31,
};

export const POTIONS_CONFIG: ConsumableStatOption<Potions>[] = [
	{ config: LesserManaPotion,  stats: [Stat.StatIntellect] },
	{ config: ManaPotion, 		 stats: [Stat.StatIntellect] },
	{ config: GreaterManaPotion, stats: [Stat.StatIntellect] },
];

export const makePotionsInput = makeConsumeInputFactory({consumesFieldName: 'defaultPotion'});

///////////////////////////////////////////////////////////////////////////
//                                 SPELL DAMAGE CONSUMES
///////////////////////////////////////////////////////////////////////////

// Arcane
export const GreaterArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 13454, minLevel: 46 },
	]),
	value: SpellPowerBuff.GreaterArcaneElixir,
};
export const ArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		// Requires skill 235
		{ id: 9155, minLevel: 41 },
	]),
	value: SpellPowerBuff.ArcaneElixir,
};
export const LesserArcaneElixir: ConsumableInputConfig<SpellPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 217398, minLevel: 28 },
	]),
	value: SpellPowerBuff.LesserArcaneElixir,
};

export const SPELL_POWER_CONFIG: ConsumableStatOption<SpellPowerBuff>[] = [
	{ config: GreaterArcaneElixir, 	stats: [Stat.StatSpellPower] },
	{ config: ArcaneElixir, 		stats: [Stat.StatSpellPower] },
	{ config: LesserArcaneElixir, 	stats: [Stat.StatSpellPower] },
];

export const makeSpellPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'spellPowerBuff'})

// Fire
export const ElixirOfGreaterFirepower: ConsumableInputConfig<FirePowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		// Requires skill 250
		{ id: 21546, minLevel: 41 },
	]),
	value: FirePowerBuff.ElixirOfGreaterFirepower,
};
export const ElixirOfFirepower: ConsumableInputConfig<FirePowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 6373, minLevel: 18 },
	]),
	value: FirePowerBuff.ElixirOfFirepower,
};

export const FIRE_POWER_CONFIG: ConsumableStatOption<FirePowerBuff>[] = [
	{ config: ElixirOfGreaterFirepower, stats: [Stat.StatFirePower] },
	{ config: ElixirOfFirepower, 		stats: [Stat.StatFirePower] },
];

export const makeFirePowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'firePowerBuff'})

// Frost
export const ElixirOfFrostPower: ConsumableInputConfig<FrostPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 17708, minLevel: 40 },
	]),
	value: FrostPowerBuff.ElixirOfFrostPower,
};

export const FROST_POWER_CONFIG: ConsumableStatOption<FrostPowerBuff>[] = [
	{ config: ElixirOfFrostPower, stats: [Stat.StatFrostPower] },
];

export const makeFrostPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'frostPowerBuff'})

// Shadow
export const ElixirOfShadowPower: ConsumableInputConfig<ShadowPowerBuff> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 9264, minLevel: 40 },
	]),
	value: ShadowPowerBuff.ElixirOfShadowPower,
};

export const SHADOW_POWER_CONFIG: ConsumableStatOption<ShadowPowerBuff>[] = [
	{ config: ElixirOfShadowPower, stats: [Stat.StatShadowPower] },
];

export const makeshadowPowerConsumeInput = makeConsumeInputFactory({consumesFieldName: 'shadowPowerBuff'})

///////////////////////////////////////////////////////////////////////////
//                                 Weapon Imbues
///////////////////////////////////////////////////////////////////////////

// Windfury (Buff)
export const Windfury: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingSpellActionId([
		{ id: 8516, 	minLevel: 32, maxLevel: 41	},
		{ id: 10608, 	minLevel: 42, maxLevel: 51	},
		{ id: 10610, 	minLevel: 52								},
	]),
	value: WeaponImbue.Windfury,
};

// Wild Strikes
export const WildStrikes: ConsumableInputConfig<WeaponImbue> = {
	actionId: () => ActionId.fromSpellId(407975),
	value: WeaponImbue.WildStrikes,
};

// Other Imbues

export const ElementalSharpeningStone: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 18262, minLevel: 50 },
	]),
	value: WeaponImbue.ElementalSharpeningStone,
};
export const BrillianWizardOil: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 20749, minLevel: 45 },
	]),
	value: WeaponImbue.BrillianWizardOil,
};
export const BrilliantManaOil: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 20748, minLevel: 45 },
	]),
	value: WeaponImbue.BrilliantManaOil,
};
export const DenseSharpeningStone: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 12404, minLevel: 35 },
	]),
	value: WeaponImbue.DenseSharpeningStone,
};
export const BlackfathomManaOil: ConsumableInputConfig<WeaponImbue> = {
	actionId: (player) => player.getMatchingItemActionId([
		{ id: 211848, minLevel: 25 },
	]),
	value: WeaponImbue.BlackfathomManaOil,
};
export const BlackfathomSharpeningStone: ConsumableInputConfig<WeaponImbue> = {
	actionId: () => ActionId.fromItemId(211845),
	value: WeaponImbue.BlackfathomSharpeningStone,
};

const SHAMAN_IMBUES: ConsumableStatOption<WeaponImbue>[] = [
	{ config: RockbiterWeaponImbue,		stats: [] },
	{ config: FlametongueWeaponImbue,	stats: [] },
	{ config: FrostbrandWeaponImbue,	stats: [] },
	{ config: WindfuryWeaponImbue,		stats: [] },
]

const CONSUMABLES_IMBUES: ConsumableStatOption<WeaponImbue>[] = [
	{ config: ElementalSharpeningStone, 	stats: [Stat.StatAttackPower] },
	{ config: BrillianWizardOil, 					stats: [Stat.StatSpellPower] },
	{ config: BrilliantManaOil, 					stats: [Stat.StatHealing, Stat.StatSpellPower] },
	{ config: DenseSharpeningStone, 			stats: [Stat.StatAttackPower] },
	{ config: BlackfathomManaOil, 				stats: [Stat.StatSpellPower, Stat.StatMP5] },
	{ config: BlackfathomSharpeningStone, stats: [Stat.StatMeleeHit] },
]

export const WEAPON_IMBUES_OH_CONFIG: ConsumableStatOption<WeaponImbue>[] = [
	...SHAMAN_IMBUES,
	...CONSUMABLES_IMBUES,
];

export const WEAPON_IMBUES_MH_CONFIG: ConsumableStatOption<WeaponImbue>[] = [
	...SHAMAN_IMBUES,
	{ config: Windfury, 		stats: [Stat.StatMeleeHit] },
	{ config: WildStrikes, 	stats: [Stat.StatMeleeHit] },
	...CONSUMABLES_IMBUES,
];

export const makeMainHandImbuesInput = makeConsumeInputFactory({
	consumesFieldName: 'mainHandImbue',
	showWhen: (player) => !!player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand),
});
export const makeOffHandImbuesInput = makeConsumeInputFactory({
	consumesFieldName: 'offHandImbue',
	showWhen: (player) => !!player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand),
});
