import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	Alcohol,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Flask,
	Food,
	HealthElixir,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	Race,
	RaidBuffs,
	SapperExplosive,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import {
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
	Hunter_Options_QuiverBonus,
} from '../core/proto/hunter.js';
import { SavedTalents } from '../core/proto/ui.js';
import MeleeWeaveP1 from './apls/p1_weave.apl.json';
import MeleeP2 from './apls/p2_melee.apl.json';
import RangedBmP2 from './apls/p2_ranged_bm.apl.json';
import RangedMmP2 from './apls/p2_ranged_mm.apl.json';
import MeleeBmP3 from './apls/p3_melee_bm.apl.json';
import RangedMmP3 from './apls/p3_ranged_mm.apl.json';
import RangedP4 from './apls/p4_ranged.apl.json';
import WeaveP4 from './apls/p4_weave.apl.json';
import Phase5AplMeleeBm from './apls/p5_melee_bm.apl.json';
import Phase5AplMeleeSv from './apls/p5_melee_sv.apl.json';
import Phase5AplRanged from './apls/p5_ranged.apl.json';
import Phase5AplWeave from './apls/p5_weave.apl.json';
import Phase6AplMeleeBm from './apls/p6_melee_bm.apl.json';
import Phase6AplMeleeSv from './apls/p6_melee_sv.apl.json';
import Phase6AplRangedDraconic from './apls/p6_ranged_draconic.apl.json';
import Phase6AplRangedKillshot from './apls/p6_ranged_killshot.apl.json';
import Phase6AplWeave from './apls/p6_weave.apl.json';
import Phase2GearMelee from './gear_sets/p2_melee.gear.json';
import Phase2GearRangedBm from './gear_sets/p2_ranged_bm.gear.json';
import Phase2GearRangedMm from './gear_sets/p2_ranged_mm.gear.json';
import Phase3GearMeleeBm from './gear_sets/p3_melee_bm.gear.json';
import Phase3GearRangedMm from './gear_sets/p3_ranged_mm.gear.json';
import Phase4GearRanged from './gear_sets/p4_ranged.gear.json';
import Phase4GearWeave from './gear_sets/p4_weave.gear.json';
import Phase5GearMeleeBm from './gear_sets/p5_melee_bm.gear.json';
import Phase5GearMeleeSv from './gear_sets/p5_melee_sv.gear.json';
import Phase5GearRangedSv from './gear_sets/p5_ranged_sv.gear.json';
import Phase5GearWeave from './gear_sets/p5_weave.gear.json';
import Phase6GearMeleeBm from './gear_sets/p6_melee_bm.gear.json';
import Phase6GearMeleeSv from './gear_sets/p6_melee_sv.gear.json';
import Phase6GearRangedDraconic from './gear_sets/p6_ranged_draconic.gear.json';
import Phase6GearRangedKillshot from './gear_sets/p6_ranged_killshot.gear.json';
import Phase6GearWeave from './gear_sets/p6_weave.gear.json';
import Phase1Gear from './gear_sets/phase1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBeastMasteryPhase1 = PresetUtils.makePresetGear('P1 Beast Mastery', Phase1Gear, { customCondition: player => player.getLevel() === 25 });
export const GearMarksmanPhase1 = PresetUtils.makePresetGear('P1 Marksmanship', Phase1Gear, { customCondition: player => player.getLevel() === 25 });
export const GearSurvivalPhase1 = PresetUtils.makePresetGear('P1 Survival', Phase1Gear, { customCondition: player => player.getLevel() === 25 });

export const GearRangedBmPhase2 = PresetUtils.makePresetGear('P2 Ranged BM', Phase2GearRangedBm, { customCondition: player => player.getLevel() === 40 });
export const GearRangedMmPhase2 = PresetUtils.makePresetGear('P2 Ranged MM', Phase2GearRangedMm, { customCondition: player => player.getLevel() === 40 });
export const GearMeleePhase2 = PresetUtils.makePresetGear('P2 Melee', Phase2GearMelee, { customCondition: player => player.getLevel() === 40 });

export const GearMeleeBmPhase3 = PresetUtils.makePresetGear('P3 Melee BM', Phase3GearMeleeBm, { customCondition: player => player.getLevel() === 50 });
export const GearRangedMmPhase3 = PresetUtils.makePresetGear('P3 Ranged MM', Phase3GearRangedMm, { customCondition: player => player.getLevel() === 50 });

export const GearWeavePhase4 = PresetUtils.makePresetGear('P4 Weave', Phase4GearWeave, { customCondition: player => player.getLevel() === 60 });
export const GearRangedSVPhase4 = PresetUtils.makePresetGear('P4 Ranged', Phase4GearRanged, { customCondition: player => player.getLevel() === 60 });

export const GearWeavePhase5 = PresetUtils.makePresetGear('P5 Weave', Phase5GearWeave, { customCondition: player => player.getLevel() === 60 });
export const GearRangedMMPhase5 = PresetUtils.makePresetGear('P5 Ranged MM', Phase5GearRangedSv, { customCondition: player => player.getLevel() === 60 });
export const GearRangedSVPhase5 = PresetUtils.makePresetGear('P5 Ranged SV', Phase5GearRangedSv, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeBMPhase5 = PresetUtils.makePresetGear('P5 Melee BM', Phase5GearMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeSVPhase5 = PresetUtils.makePresetGear('P5 Melee SV', Phase5GearMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const GearWeavePhase6 = PresetUtils.makePresetGear('P6 Weave', Phase6GearWeave, { customCondition: player => player.getLevel() === 60 });
export const GearRangedDraconicPhase6 = PresetUtils.makePresetGear('P6 Ranged Draconic', Phase6GearRangedDraconic, { customCondition: player => player.getLevel() === 60 });
export const GearRangedKillshotPhase6 = PresetUtils.makePresetGear('P6 Ranged Killshot', Phase6GearRangedKillshot, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeBMPhase6 = PresetUtils.makePresetGear('P6 Melee BM', Phase6GearMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeSVPhase6 = PresetUtils.makePresetGear('P6 Melee SV', Phase6GearMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearBeastMasteryPhase1, GearMarksmanPhase1, GearSurvivalPhase1],
	[Phase.Phase2]: [GearRangedBmPhase2, GearRangedMmPhase2, GearMeleePhase2],
	[Phase.Phase3]: [GearRangedMmPhase3, GearMeleeBmPhase3],
	[Phase.Phase4]: [], //[GearWeavePhase4, GearRangedSVPhase4],
	[Phase.Phase5]: [], //[GearWeavePhase5, GearRangedMMPhase5, GearRangedSVPhase5, GearMeleeBMPhase5, GearMeleeSVPhase5],
	[Phase.Phase6]: [GearWeavePhase6, GearRangedDraconicPhase6, GearRangedKillshotPhase6, GearMeleeBMPhase6, GearMeleeSVPhase6],
};

export const DefaultGearWeave = GearPresets[Phase.Phase6][0];
export const DefaultGearRangedDraconic = GearPresets[Phase.Phase6][1];
export const DefaultGearRangedKillshot = GearPresets[Phase.Phase6][2];
export const DefaultGearMeleeBM = GearPresets[Phase.Phase6][3];
export const DefaultGearMeleeSV = GearPresets[Phase.Phase6][4];

export const DefaultGear = DefaultGearRangedKillshot;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLMeleeWeavePhase1 = PresetUtils.makePresetAPLRotation('P1 Melee Weave', MeleeWeaveP1, { customCondition: player => player.getLevel() === 25 });

export const APLMeleePhase2 = PresetUtils.makePresetAPLRotation('P2 Melee', MeleeP2, { customCondition: player => player.getLevel() === 40 });
export const APLRangedBmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged BM', RangedBmP2, { customCondition: player => player.getLevel() === 40 });
export const APLRangedMmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged MM', RangedMmP2, { customCondition: player => player.getLevel() === 40 });

export const APLMeleeBmPhase3 = PresetUtils.makePresetAPLRotation('P3 Melee BM', MeleeBmP3, { customCondition: player => player.getLevel() === 50 });
export const APLRangedMmPhase3 = PresetUtils.makePresetAPLRotation('P3 Ranged MM', RangedMmP3, { customCondition: player => player.getLevel() === 50 });

export const APLWeavePhase4 = PresetUtils.makePresetAPLRotation('P4 Weave', WeaveP4, { customCondition: player => player.getLevel() === 60 });
export const APLRangedPhase4 = PresetUtils.makePresetAPLRotation('P4 Ranged', RangedP4, { customCondition: player => player.getLevel() === 60 });

export const APLWeavePhase5 = PresetUtils.makePresetAPLRotation('P5 Weave', Phase5AplWeave, { customCondition: player => player.getLevel() === 60 });
export const APLRanged31Phase5 = PresetUtils.makePresetAPLRotation('P5 Ranged 3-1-1', Phase5AplRanged, { customCondition: player => player.getLevel() === 60 });
export const APLRanged22Phase5 = PresetUtils.makePresetAPLRotation('P5 Ranged 2-2', RangedP4, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeBMPhase5 = PresetUtils.makePresetAPLRotation('P5 Melee BM', Phase5AplMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeSVPhase5 = PresetUtils.makePresetAPLRotation('P5 Melee SV', Phase5AplMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const APLWeavePhase6 = PresetUtils.makePresetAPLRotation('P6 Weave', Phase6AplWeave, { customCondition: player => player.getLevel() === 60 });
export const APLRangedDraconicPhase6 = PresetUtils.makePresetAPLRotation('P6 Ranged Draconic', Phase6AplRangedDraconic, { customCondition: player => player.getLevel() === 60 });
export const APLRangedKillshotPhase6 = PresetUtils.makePresetAPLRotation('P6 Ranged Killshot', Phase6AplRangedKillshot, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeBMPhase6 = PresetUtils.makePresetAPLRotation('P6 Melee BM', Phase6AplMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeSVPhase6 = PresetUtils.makePresetAPLRotation('P6 Melee SV', Phase6AplMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const APLPresets = {
	[Phase.Phase1]: [APLMeleeWeavePhase1],
	[Phase.Phase2]: [APLRangedBmPhase2, APLRangedMmPhase2, APLMeleePhase2],
	[Phase.Phase3]: [APLRangedMmPhase3, APLMeleeBmPhase3],
	[Phase.Phase4]: [], //[APLWeavePhase4, APLRangedPhase4],
	[Phase.Phase5]: [], //[APLWeavePhase5, APLRanged31Phase5, APLRanged22Phase5, APLMeleeBMPhase5, APLMeleeSVPhase5],
	[Phase.Phase6]: [APLWeavePhase6, APLRangedDraconicPhase6, APLRangedKillshotPhase6, APLMeleeBMPhase6, APLMeleeSVPhase6],
};

export const DefaultAPLWeave = APLPresets[Phase.Phase6][0];
export const DefaultAPLRangedDraconic = APLPresets[Phase.Phase6][1];
export const DefaultAPLRangedKillshot = APLPresets[Phase.Phase6][2];
export const DefaultAPLMeleeBM = APLPresets[Phase.Phase6][3];
export const DefaultAPLMeleeSV = APLPresets[Phase.Phase6][4];

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsBeastMasteryPhase1 = PresetUtils.makePresetTalents('P1 Beast Mastery', SavedTalents.create({ talentsString: '53000200501' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsMarksmanPhase1 = PresetUtils.makePresetTalents('P1 Marksmanship', SavedTalents.create({ talentsString: '-050515' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsSurvivalPhase1 = PresetUtils.makePresetTalents('P1 Survival', SavedTalents.create({ talentsString: '--33502001101' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsBeastMasteryPhase2 = PresetUtils.makePresetTalents('P2 Beast Mastery', SavedTalents.create({ talentsString: '5300021150501251' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsMarksmanPhase2 = PresetUtils.makePresetTalents('P2 Marksmanship', SavedTalents.create({ talentsString: '-05551001503051' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsSurvivalPhase2 = PresetUtils.makePresetTalents('P2 Survival', SavedTalents.create({ talentsString: '--335020051030315' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsRangedMMPhase3 = PresetUtils.makePresetTalents('P3 Ranged MM', SavedTalents.create({ talentsString: '5-05051404503051-3' }), {
	customCondition: player => player.getLevel() === 50,
});
export const TalentsMeleeBMPhase3 = PresetUtils.makePresetTalents('P3 Melee BM', SavedTalents.create({ talentsString: '0500321150521251--33002' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TalentsWeavePhase4 = PresetUtils.makePresetTalents('60 Weave', SavedTalents.create({ talentsString: '-055500005-3305202202303051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedMMPhase4 = PresetUtils.makePresetTalents('60 Ranged MM', SavedTalents.create({ talentsString: '-05451002503051-33400023023' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedSVPhase4 = PresetUtils.makePresetTalents('60 Ranged SV', SavedTalents.create({ talentsString: '1-054510005-334000250230305' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentsWeavePhase5 = PresetUtils.makePresetTalents('P5 Weave', SavedTalents.create({ talentsString: '-055500005-3305202202303051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedMMPhase5 = PresetUtils.makePresetTalents('P5 Ranged MM', SavedTalents.create({ talentsString: '5-05451005503051-3320202' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedSVPhase5 = PresetUtils.makePresetTalents('P5 Ranged SV', SavedTalents.create({ talentsString: '1-054510005-334000250230305' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeBMPhase5 = PresetUtils.makePresetTalents('P5 Melee BM', SavedTalents.create({ talentsString: '5500020050521251-0505-33002' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeSVPhase5 = PresetUtils.makePresetTalents('P5 Melee SV', SavedTalents.create({ talentsString: '-055500005-3320202412303051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentsWeavePhase6 = PresetUtils.makePresetTalents('P6 Weave', SavedTalents.create({ talentsString: '-054510005-3305202202303051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedMMPhase6 = PresetUtils.makePresetTalents('P6 Ranged MM', SavedTalents.create({ talentsString: '5-05451005503051-3320202' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedSVPhase6 = PresetUtils.makePresetTalents('P6 Ranged SV', SavedTalents.create({ talentsString: '-054510015-334000250230305' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeBMPhase6 = PresetUtils.makePresetTalents('P6 Melee BM', SavedTalents.create({ talentsString: '5500020050501251-0505-33202' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeSVPhase6 = PresetUtils.makePresetTalents('P6 Melee SV', SavedTalents.create({ talentsString: '-055500005-3320202412303051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsBeastMasteryPhase1, TalentsMarksmanPhase1, TalentsSurvivalPhase1],
	[Phase.Phase2]: [TalentsBeastMasteryPhase2, TalentsMarksmanPhase2, TalentsSurvivalPhase2],
	[Phase.Phase3]: [TalentsRangedMMPhase3, TalentsMeleeBMPhase3],
	[Phase.Phase4]: [], //[TalentsWeavePhase4, TalentsRangedMMPhase4, TalentsRangedSVPhase4],
	[Phase.Phase5]: [], //[TalentsWeavePhase5, TalentsRangedMMPhase5, TalentsRangedSVPhase5, TalentsMeleeBMPhase5, TalentsMeleeSVPhase5],
	[Phase.Phase6]: [TalentsWeavePhase6, TalentsRangedMMPhase6, TalentsRangedSVPhase6, TalentsMeleeBMPhase6, TalentsMeleeSVPhase6],
};

export const DefaultTalentsWeave = TalentPresets[Phase.Phase6][0];
export const DefaultTalentsRangedDraconic = TalentPresets[Phase.Phase6][2];
export const DefaultTalentsRangedKillshot = TalentPresets[Phase.Phase6][2];
export const DefaultTalentsMeleeBM = TalentPresets[Phase.Phase6][3];
export const DefaultTalentsMeleeSV = TalentPresets[Phase.Phase6][4];

export const DefaultTalents = DefaultTalentsRangedKillshot;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const OptionsRangedLonewolf = HunterOptions.create({
	ammo: Ammo.ThoriumHeadedArrow,
	quiverBonus: Hunter_Options_QuiverBonus.Speed15,
	petAttackSpeed: 2.0,
	petTalents: {},
	petType: PetType.PetNone,
	petUptime: 1,
	sniperTrainingUptime: 1.0,
});

export const OptionsRangedPet = HunterOptions.create({
	ammo: Ammo.ThoriumHeadedArrow,
	quiverBonus: Hunter_Options_QuiverBonus.Speed15,
	petAttackSpeed: 2.0,
	petTalents: {},
	petType: PetType.WindSerpent,
	petUptime: 1,
	sniperTrainingUptime: 1.0,
});

export const OptionsMeleePet = HunterOptions.create({
	ammo: Ammo.ThoriumHeadedArrow,
	quiverBonus: Hunter_Options_QuiverBonus.Speed15,
	petAttackSpeed: 2.0,
	petTalents: {},
	petType: PetType.Cat,
	petUptime: 1,
	sniperTrainingUptime: 1.0,
});

export const DefaultOptions = OptionsRangedLonewolf;

// Consumable Presets

export const MeleeConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfMadness,
	food: Food.FoodSmokedDesertDumpling,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.WildStrikes,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	miscConsumes: {
		draughtOfTheSands: true,
	},
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	petAttackPowerConsumable: 1,
	petAgilityConsumable: 1,
	petStrengthConsumable: 1,
	sapperExplosive: SapperExplosive.SapperFumigator,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const RangedConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: false,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodGrilledSquid,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.EnchantedRepellent,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	offHandImbue: WeaponImbue.EnchantedRepellent,
	petAttackPowerConsumable: 1,
	petAgilityConsumable: 1,
	petStrengthConsumable: 1,
	sapperExplosive: SapperExplosive.SapperFumigator,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const WeaveConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: false,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodGrilledSquid,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.WildStrikes,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	offHandImbue: WeaponImbue.EnchantedRepellent,
	petAttackPowerConsumable: 1,
	petAgilityConsumable: 1,
	petStrengthConsumable: 1,
	sapperExplosive: SapperExplosive.SapperFumigator,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const DefaultConsumes = RangedConsumes;

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	commandingShout: true,
	demonicPact: 120,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectRegular,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	fervorOfTheTempleExplorer: true,
	mightOfStormwind: true,
	moldarsMoxie: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	spiritOfZandalar: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	huntersMark: TristateEffect.TristateEffectImproved,
	improvedScorch: true,
	judgementOfWisdom: true,
	mangle: true,
	markOfChaos: true,
	occultPoison: true,
	stormstrike: true,
});

export const OtherDefaults = {
	distanceFromTarget: 12,
	profession1: Profession.Enchanting,
	profession2: Profession.Engineering,
	race: Race.RaceTroll,
};

///////////////////////////////////////////////////////////////////////////
//                                 Encounters
///////////////////////////////////////////////////////////////////////////

export const EncounterWeavePhase6 = PresetUtils.makePresetEncounter(
	'Weave',
	'https://wowsims.github.io/sod/hunter/#eJztV2toHFUU3nN3djM5acbJzWv2Jmm321anYzad3c3GpGm7k+CPEqWGWKG/pIZWaikkuCD6r5YWbbUS8kcaUVOtNgQjJaBgEVIjQqtoW6Ulvgui+ECM/SFVaeOZ12Z2G03FCopelt17zzn3O4/73cdiXGYqizOdmdAOFtsKJ9m2MYDdbA8MMTjGoLtJgzGnZ8EMwAWAXnYMtu2AUWAnAXhIHI1htHfXPQ9tv1+NaOvizBgDrFCnn1W0Dzr1iYkYSuqvwwqWqW+Q6OUa0k2OK9rZqH7uR0GDN0n6haIfIMMK9cmDivZ1p37lJXvwNA321ugfXtFQVg8dVPTROY3EJ1zkn53BM2TzuaL/cLmBbD56QdEf32v3Xjys6I9RT1KfOqzQ96Un7AAujira+RoMUXffYUX7pjPdLJeroEed5GGjNMTYg2WPQmQI2ChIVIcTAO+DpEKcco+s/UrCCyz0f/tbG9/9ny/xtyzUXd0OxD8iJfTCFtgKO2AQehqSZrY1mzJNM5vMZMxs2kzTJ2NSN9UrbYH+0M4zVVoYq1EWUZTkQ2cbE2UYwXAym8cYopAxKo89T4xPlCMpkpmWbF4DrEV5p23/ykQsoMA6D+aTUUgg0txkqqWD5JorV38ZqvTkZkurA+TPeO1dCGoK8lcf6SmS12OZoPDkL+tLgDQv2NNvgV4e1PEGEcMq4wY5zFl/SDjZZbOrhIRMq+RxsRRreXWpQRsZ0I/6Uz9vEg2ONhqYbq7qs4ObmojxlSKBmlEnSxx3FvzboSPL5p0U6Dji60UntvBmrDWq5QiPTDMmMztclFItt+Sx2qgicXgYQq7QJOEwOF1kqTSfBfE94AXgn4JjG3UhIq6fVB77RK8TpAfiizt5B4YPQAjrjVrSlY+DF1DBImbUk6Ji3HYmv0exz2PeJe5EYWhUliXjQIssHx+Ih3XZ1afzji6ykC7jwpbOc7PGykQFlk86gRx5pxGX2Gs1Ce6aDoNTX4QspSyD+A7wM+Afg1O3QtJeidr/VWkXciwpQCHlNt5FGa/HFF+DSmIJ4oitkU8P6HJw/PqAHnaTW6hqkmBrQ0HKrvUom864ZCSPvJ38tKLBdac+Ub8+Fy83FBIJxEhbvDAzSTNXY5wvdQ2e8+saHNowblzOzlMNxYvDwW7N+8pGIRZQeqcEv5n8LONNxfQo8drnG6+gsJZhnVFjJzMN3rkh3LL6/iwyWodpbmK1qApgRt6m9YHShALVd7KnxwNfI5K4gi9feP5CE+4gl7dhF88t6pKKFhwen9OC610AXC1uwibeUMSN4sALpnYBFzEtcJ/HBNWWk3KkCMc+uSpFhZhnM6dDVngB8w1/VM8gXUvY4Ye4icrTgxbfUEK30iPCY0mPxxLa90WsEF7FeD/h3Y2bed+fxrOpIwXPCYc6mQV8zLNocR9S0IeDMEUIt4pubOdtv7P3UvlSJl6dq107S2zAVp6+ZpSrEfz959567oW1SvjXGi268HnDm4Uxv+ndHKeKU6Z3g/9+KJDUPWMDoXn3nXct0hvhGk3Vh48o3BD6YqY+3wZPwn5wH0a6dYnBKe99KFtnvV6HNcNgBELT0MgFveijidBqVzObS/mdjSwxHS53R5usKrez0oqNHLLbqVyHKzmTW36v087n0sNheWZPJZdv3/7A9l3xNlPLxcPGP+wfgL7fui44N078JRxVzF2fOKi5SzZ41Oq4b2B236nN53JdnsbaAr8BPTHkGg==',
);

export const EncounterRangedDraconicPhase6 = PresetUtils.makePresetEncounter(
	'Ranged Draconic',
	'https://wowsims.github.io/sod/hunter/#eJztVl1oHFUUnnNnZnf2pBknt5vN5DbGcUllHLth9mdCfpSdRh9i1bKKD3nU0EgrhQYDom8lpGiV1BipYH0wwdqEQiQEFOxLa6CQIGojWKKCFkRRn2ofpCi4nvnZ7RpbI1hB0Puwc+79vvvdc86cc2cxqzGDWcxmLvSCz4bgMVhl+xYADrMJmGawxGDwVhMWQsuHDYBLABW2BPv2wyywVQAuiaUWTFQOPv7s6FNG0rzbYs4CYJPxypRuvj9gv7PYjorxy4yOSeNlWno7Tdjyad1cT9jrPwqavEqrX+v2C0RsMl6iyXcD9uLxAJma083JtP3ZryZqxvSUbs9WTVo+90aoPB9uOE4bvtLt89WA8/lbun1scgdZp+Z0+0WyFOPqMZ1+X58LHLgyq5sX0yiReYS0vx8o3KWlDLATYfgwpIwln+l8HtRpYLOgUBrOAXwCigEWha72f6vgJSb9P/7RwQ//51P8A5MGt/cC1R/VJFRgmNpyP4zBHpFzvZKXd928lysWS67rFjy3UHSLrldRhmFEevKLlCljKyaFirL2TVtWwwQquUJh3ATMIAqaaxfOg52Mkbw7jm3x+sKb1AnZFCZRzRW7vXGklhJE0t59bk8WkSg5t7tUCqRixPhyXt+EZGLk5+nmRqSu9d6H8Aet6PiPya1UI8Y7hEA9uw2xUMevMUO0xblFUzgbkQQ5jbLn7QwCNyYPXA/tidGfRnin6EDuGFqCs35JhKkoUMS1aP8MD2Lmt2kgdlyXEeWdW0S48RGUhPAIU2QwtQzxqbEZHNAv8S5SsDDjpEkhtVLjBPoI+brKw8R6CO/lu7HNadVUnjodMs9WzYCCrDj+e+BUA9CcbaqdTrd0KHnygw7uim7s4tkwe5HzYfby3g13JMgJ4vFdwkGLd9Z21t+Ku7O26+xi+yPxk1fEXryPD24SDd2Va+7S1R+5W6AKcjJBKs7UA4zccqOUnqmavEeU0OH2X1UUsVDd7037GqYUpYjTx6laRazEt4mG1qnPPqL6VHmv6MFd3MFWZzvlX11hTFOjOsgHzRUkSeXyDEiRP/mouiZO6twSnZjmHA1Hj+s3ZFD51nL3gLgf7+EDjS7SctQsy0CJ0a4cshMoHPPaqw8iD5s7312iOqzfBOFxQaUCj4srfLghg2rMYlsy5C0ZypYMdUtGYmwVjkJ0Rdr+VQZr8ZdC89djq8/fYHACpBXo4II+7YmsdGeEXC7na8YQy67IqWi212+JjC6//cRrwVgr90UrF8q3PxGOi+XCjKxtTDRz7cHRp0cPWj2uWbZk51/2X8A+6t8UnTsW/5aOIao3xw8a0Ssbm/f7Dhy6fGTt0U/Lu2PEH4bfAEReM3g=',
);

export const EncounterRangedKillshotPhase6 = PresetUtils.makePresetEncounter(
	'Ranged Killshot',
	'https://wowsims.github.io/sod/hunter/#eJztVE1oVFcUfufOvJn7Tszr82aSvFxjGR8ir1NG3sxkJPGHebqK8W8WLrIQseIUG0QHB0q7k1TxDyUEXKgLE40aAiMSULArHRESS6sRlOiiFYqiXRTbRbEtGM/7yRjrshYK7VnMfPd855x7vnPvu2hxZrAks5kDneCybtgGE2zHKMA+1g8DDC4zWLPQhFEfuTAN8BigyC7Djp0wBGwCQCjyq0aMFXd98mVprxE3VyZZahSwwaid0c0bK+xqtQ2jxp+DOsaNG+S6lCBufEw3p2L21C+SFjfJ+6NuH6HABuPnEd18RlknPObYsG7uT9gPX5nIjVcjuj00Y5L7elD5d39xnxJ+0O1bM17Mo/O6fXz/AkIXhnX7KKGo8fK4Tr+nh70Gfh3SzQcJVAgeoNrPV2Q/5poBdsyXD93RcvyLDw+BOgBsCKI0husA9yBqQJKkq8ufRvExU/63f9TEvv/8iH9iypqmTqD7R3cSitBLn+VOKEPPgkzayXfkM47j5NO5XAf9Z/NONufknHwx2gvblb6qakawBVFyjPG7t8COW4Qwms44FWwN/aPn6M5bGsZRTeeW5iuUwCUFGX8MNFqIFJJ2lnaQ3wwT7lAh7S8M76MMfqXaNtdvgmiXEnVrHmK2nvmmhlgn12JBrMLmVBNXhVpjjKvS7zDjNTI/9QG5I4OgeP0gywSuiGA9iqR2MZJzKn6r/SO62Cq3ID8L5OVPWkURY5vKpd2lvcIXw+ndCBE9FSH6uo6uHuwJkPH9RV0EMr+jFtU52ZZMYksqwWNCq0GYJGk3hLwfc+1bEIulhWaqhVrEvjeC/eaXVeqVemS3L7s11UyR2hiExGygT6jvEB0V1M5CKEGGOhNSIB8PVdf1bZDrcLUovF2I5AaFaLRtqVYiGsagrjSgOut7hNvykW/aRbNsQm0cwmVdB83pHeztPou9+c6Tc6fZIDXv4PidPTYXhtRRvU1nDp4Y47ft5Qk4DMGtt92XDCbDj5+7UyHqcqcZnAKlBu1C0msds5SPAuZFITMLuplVi2jBaqM7PwCL3bZTJz2bLHQFnruFRZ/69qCQHYzw6f5GwdeXPi/tSi5zzEIykvqXPe/2Yfe91FlS/Vt1DDnzfvogC46sfNHt+mzPiwOTm+8XVoeM2wuvAX2b+Tk=',
);

export const EncounterMeleeBMPhase6 = PresetUtils.makePresetEncounter(
	'Melee BM',
	'https://wowsims.github.io/sod/hunter/#eJztVF1IFFEU3nN3dnf2mNN0XXX2arBtaePCyuzYVGqxmxVJDyHUg48lWZqli0tQPa1m9IckQaT1oEJUSP4QPUQQlhQY9GeQ+RBURFERqC/hSz937q76UD1VENQZmPud73z33HPP3LmYLxOVBIhODFgNMbId4iRJWqGDwBCBCqbBZYFiMAFQRYZgZx30ABkFoA42mIHuqr07DtY2qy5tTYCEBgEz1JFeRbtdrvf3+1FSO08p6FHvcGrAx2P3uxVtzK2PTTPu3OLsa0Xv/KJx52y7or0r179c8XPn2hlFa/PpXZ/zUFaHzij6jKDPt4vME8K5zWe/UPRnST/XJHsV/WSbrX7Yreg3D+fxpc/3Kvx94qhdwHSPon20wYPuNDgxkqm9LzdXyh6V6G6xeaiUOgg54EmCdAxcHUB6QLoMMARwHWAG4BbAEyAq74Sr7K2EL4njv/1Ro8l/vsUfiKNCFWcTqqAatkMdxGFznmUZhmEahsWfiGlFwjYKl5SYhlklVUONY887p+bEHMTT4EEXOhGsoIxulMIRI4FZKDOO5a6x/KAIh60E5iIyLpEf3QXdG/QiD4QjxTyQnVY/74Ef0Z86Z2mjOJKgEiNlDlrAliKjGuaEfLKbekfAVvIfl4laTHuQ3+RSi5VgEV2OmcEM9F4VIvsGsGc5v5sVYWmXRtlaLKERzA1l27obczo+IjETP1l2PgHfqcCTHGcwL+PVy8NNujQXsEUmM7CQLkM1pPBVSI0jld4qEIybks1pZkVCzBnu99PFLA8XhRbaUa4XvbWMgtlm041sPZbSVegP5dr19UG6CCa+jWgrC2myi2IfzNY0HxNpbnzV6AIZUnWqLReU+Cgch9Rp8cVmCJwDxwhk0yxZVt0aBB1FA5teHfINTEYrPcFxpzel3BJblALLYv5zXbbdi5ammMfRJbuEjUfNi055ojWT5m/d39hQ2xjYVrsvvrc2sKF5x+6mxob6xkBFUyKhRQPu0F92GU61xH5LnsL+X8qjBlau+y11cIuId/xSrLS+aerIvW1Po7OpY9XwDYq30s4=',
);

export const EncounterMeleeSVPhase6 = PresetUtils.makePresetEncounter(
	'Melee SV',
	'https://wowsims.github.io/sod/hunter/#eJztVF1IFFEU3nP3b/aY03j9m5012ra0aWhjZteJzGI36UF8CCEffLTI0jRdWqJ6szL6wxAfKotIX0oiQ3qI8EXbelCIyCLxISqIoiLQ6kUIrTt3d6wwnzII6gzMnPOd73z33HNnBgMCkUiQqESHDRAndZAgbWRTJ4EBAhVFMvRxLw4TANVkAHY1QA+QEYBKoA7lyRL0VDfvOFy/X3LLm4JE6wPMklK9ovykXO3v96NL+tIlole6z6CbeSz34Iooj3nUpx8VFgwz9JWoXpiVWXC+Q5TflquzN/wsuMyC9jy1eyaAgtTdIarTHL7UwZUneHCXVb8Q1UnOaesV1bPtlvfwiqieYZ5Len5NZPdLvVYDfYz8QUQHc0+nsuV35ZGo4JWI6uHbhkpXJyGHvCfB3QmkB1x9AAMAdwCmAYYBHgOR2ATcG9+48CVx/Lc/arTtnx/xe+KokPiLCdVQC3XQAAmoCoR10zR1ZmY4Go3o7Co1IlE9qptGtasWdjr2nnPLTixA7AIvutGJYIYE9KArbOhJzEVBYb7QPVYU4umwmWRkn8IC4XOjKoR8yNywsY7h+Rnysx74FXz7RFUG1tcZSeoRQCEbHXQley7HAi1P8FBfCiwm+1oV3krEegivC2lE0XGtoi1EQ0kTGc7krJWQGDoTNpUorqGrMTuUhb5bvML6VVgSznkSUSUT0piyGaPUwEIt3+INzvG4ciS5QA/GnEC+kmuvd+xeEYc/zQT4Forpyp/LZ+1yE31dTrvEGu5Qq+qiqKRHN8lkbd/SylJ880nW2tYiJXQVn4eTkp2ZeZjF9oSqMkhpktcM9fvpMiWAOdpSK8v4/JRNvdg+dnpQOYBNtBFlrUBw/zgPzmQnzDNsR9ch0+oCGdb494xfK2QZvA72NuZSfOHBr3JiBE5B+t3Oi08TuAiOFOTTXEGQPDKEHGvSualYpTc07vSlo23xHLvEf7HbstFYWRp5FFuxm9t4LHLVKUwczaZF2w+0NNW3BGvq9yWa64Nb9+/Y09rS1NgSrGhNJuVY0KP9Zb/uqSPxRdEp6f8tHSm4fsui9MHM4PfEtXhZY+vU8dGapzFbOl4L3wABv+hb',
);

///////////////////////////////////////////////////////////////////////////
//                                 Presets
///////////////////////////////////////////////////////////////////////////

export const PresetBuildWeave = PresetUtils.makePresetBuild('Weave', {
	gear: DefaultGearWeave,
	talents: DefaultTalentsWeave,
	rotation: DefaultAPLWeave,
	encounter: EncounterWeavePhase6,
});
export const PresetBuildRangedDraconic = PresetUtils.makePresetBuild('Ranged Draconic', {
	gear: DefaultGearRangedDraconic,
	talents: DefaultTalentsRangedDraconic,
	rotation: DefaultAPLRangedDraconic,
	encounter: EncounterRangedDraconicPhase6,
});
export const PresetBuildRangedKillshot = PresetUtils.makePresetBuild('Ranged Killshot', {
	gear: DefaultGearRangedKillshot,
	talents: DefaultTalentsRangedKillshot,
	rotation: DefaultAPLRangedKillshot,
	encounter: EncounterRangedKillshotPhase6,
});
export const PresetBuildMeleeBM = PresetUtils.makePresetBuild('Melee DW', {
	gear: DefaultGearMeleeBM,
	talents: DefaultTalentsMeleeBM,
	rotation: DefaultAPLMeleeBM,
	encounter: EncounterMeleeBMPhase6,
});
export const PresetBuildMeleeSV = PresetUtils.makePresetBuild('Melee 2H', {
	gear: DefaultGearMeleeSV,
	talents: DefaultTalentsMeleeSV,
	rotation: DefaultAPLMeleeSV,
	encounter: EncounterMeleeSVPhase6,
});