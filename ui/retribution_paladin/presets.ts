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
	Explosive,
	FirePowerBuff,
	Flask,
	Food,
	HealthElixir,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	RaidBuffs,
	SapperExplosive,
	SaygesFortune,
	SealOfTheDawn,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { PaladinAura, PaladinOptions as RetributionPaladinOptions, PaladinSeal } from '../core/proto/paladin.js';
import { SavedTalents } from '../core/proto/ui.js';
import APLP1RetJson from './apls/p1-ret.apl.json';
import APLP2RetJson from './apls/p2-ret.apl.json';
import APLP3RetJson from './apls/p3-ret.apl.json';
import APLTwistingJson from './apls/p5p6p7-twist.apl.json';
import APLExodinJson from './apls/p7-exodin.apl.json';
import APLExodinFastStackJson from './apls/p7-exodin-fast-stack.apl.json';
import APLSealStackingJson from './apls/p7-seal-stacking.apl.json';
import APLShockadin1HJson from './apls/p7-shockadin-1h.apl.json';
import APLShockadin2HJson from './apls/p7-shockadin-2h.apl.json';
import Phase1RetGearJson from './gear_sets/p1-ret.gear.json';
import Phase2RetSoCGearJson from './gear_sets/p2-retsoc.gear.json';
import Phase2RetSoMGearJson from './gear_sets/p2-retsom.gear.json';
import Phase3RetSoMGearJson from './gear_sets/p3-retsom.gear.json';
import Phase7ExodinFastStackNonNaxxGearJson from './gear_sets/p7-exodin-fast-stack-non-naxx.gear.json';
import Phase7ExodinNaxxGearJson from './gear_sets/p7-exodin-naxx.gear.json';
import Phase7SealStackingNaxxGearJson from './gear_sets/p7-seal-stacking-naxx.gear.json';
import Phase7Shockadin1HNaxxGearJson from './gear_sets/p7-shockadin-1h-naxx.gear.json';
import Phase7Shockadin2HNaxxGearJson from './gear_sets/p7-shockadin-2h-naxx.gear.json';
import Phase7TwistingNaxxGearJson from './gear_sets/p7-twisting-naxx.gear.json';
import Phase7TwistingNonNaxxGearJson from './gear_sets/p7-twisting-non-naxx.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const Phase1RetGear = PresetUtils.makePresetGear('P1', Phase1RetGearJson, {
	customCondition: player => player.getLevel() == 25,
});
export const Phase2RetSoCGear = PresetUtils.makePresetGear('P2 SoC/DS', Phase2RetSoCGearJson, {
	customCondition: player => player.getLevel() == 40,
});
export const Phase2RetSoMGear = PresetUtils.makePresetGear('P2 SoM', Phase2RetSoMGearJson, {
	customCondition: player => player.getLevel() == 40,
});
export const Phase3RetSoMGear = PresetUtils.makePresetGear('P3 SoM', Phase3RetSoMGearJson, {
	customCondition: player => player.getLevel() == 50,
});
export const Phase7Shockadin1HNaxxGear = PresetUtils.makePresetGear('1H Shockadin (Naxx BiS)', Phase7Shockadin1HNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase7Shockadin2HNaxxGear = PresetUtils.makePresetGear('2H Shockadin (Naxx BiS)', Phase7Shockadin2HNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase7ExodinNaxxGear = PresetUtils.makePresetGear('Exodin (Naxx BiS)', Phase7ExodinNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase7ExodinFastStackNonNaxxGear = PresetUtils.makePresetGear('Exodin Fast Stack (Non-Naxx BiS)', Phase7ExodinFastStackNonNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase7SealStackingNaxxGear = PresetUtils.makePresetGear('Seal Stacking (Naxx BiS)', Phase7SealStackingNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase7SealStackingNonNaxxGear = PresetUtils.makePresetGear('Seal Stacking (Non-Naxx BiS)', Phase7SealStackingNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase7TwistingNaxxGear = PresetUtils.makePresetGear('Twisting (Naxx BiS)', Phase7TwistingNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase7TwistingNonNaxxGear = PresetUtils.makePresetGear('Twisting (Non-Naxx BiS)', Phase7TwistingNonNaxxGearJson, {
	customCondition: player => player.getLevel() == 60,
});

export const GearPresets = {
	[Phase.Phase1]: [Phase1RetGear],
	[Phase.Phase2]: [Phase2RetSoCGear, Phase2RetSoMGear],
	[Phase.Phase3]: [Phase3RetSoMGear],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
	[Phase.Phase7]: [
		Phase7TwistingNaxxGear,
		Phase7SealStackingNaxxGear,
		Phase7ExodinNaxxGear,
		Phase7Shockadin1HNaxxGear,
		Phase7Shockadin2HNaxxGear,
		Phase7TwistingNonNaxxGear,
		Phase7SealStackingNonNaxxGear,
		Phase7ExodinFastStackNonNaxxGear,
	],
};

export const DefaultGear = GearPresets[Phase.Phase7][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLP1Ret = PresetUtils.makePresetAPLRotation('P1 Ret', APLP1RetJson, {
	customCondition: player => player.getLevel() === 25,
});
export const APLP2Ret = PresetUtils.makePresetAPLRotation('P2 Ret/Shockadin', APLP2RetJson, {
	customCondition: player => player.getLevel() === 40,
});
export const APLP3Ret = PresetUtils.makePresetAPLRotation('P3 Ret/Shockadin', APLP3RetJson, {
	customCondition: player => player.getLevel() === 50,
});
export const APLSealStacking = PresetUtils.makePresetAPLRotation('Seal Stacking', APLSealStackingJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLShockadin1H = PresetUtils.makePresetAPLRotation('Shockadin 1H', APLShockadin1HJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLShockadin2H = PresetUtils.makePresetAPLRotation('Shockadin 2H', APLShockadin2HJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLExodin = PresetUtils.makePresetAPLRotation('Exodin', APLExodinJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLExodinFastStack = PresetUtils.makePresetAPLRotation('Exodin Fast Stack', APLExodinFastStackJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLTwisting = PresetUtils.makePresetAPLRotation('Twist', APLTwistingJson, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [APLP1Ret],
	[Phase.Phase2]: [APLP2Ret],
	[Phase.Phase3]: [APLP3Ret],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
	[Phase.Phase7]: [APLTwisting, APLExodin, APLSealStacking, APLShockadin1H, APLShockadin2H, APLExodinFastStack],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase7][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const P1RetTalents = PresetUtils.makePresetTalents('P1 Ret', SavedTalents.create({ talentsString: '--05230051' }), {
	customCondition: player => player.getLevel() === 25,
});

export const P2RetTalents = PresetUtils.makePresetTalents('P2 Ret', SavedTalents.create({ talentsString: '--532300512003151' }), {
	customCondition: player => player.getLevel() === 40,
});

export const P2ShockadinTalents = PresetUtils.makePresetTalents('P2 Shockadin', SavedTalents.create({ talentsString: '55050100521151--' }), {
	customCondition: player => player.getLevel() === 40,
});

export const P3RetTalents = PresetUtils.makePresetTalents('P3 Ret', SavedTalents.create({ talentsString: '500501--53230051200315' }), {
	customCondition: player => player.getLevel() === 50,
});

export const RetTalents = PresetUtils.makePresetTalents('Exodin / Stacking / Twisting', SavedTalents.create({ talentsString: '500501-503-52230351200315' }), {
	customCondition: player => player.getLevel() === 60,
});

export const ShockadinTalents = PresetUtils.makePresetTalents('Shockadin', SavedTalents.create({ talentsString: '55053100501051--052303502' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [P1RetTalents],
	[Phase.Phase2]: [P2RetTalents, P2ShockadinTalents],
	[Phase.Phase3]: [P3RetTalents],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
	[Phase.Phase7]: [RetTalents, ShockadinTalents],
};

export const PresetBuildP7Twisting = PresetUtils.makePresetBuild('Twisting', {
	gear: Phase7TwistingNaxxGear,
	talents: RetTalents,
	rotation: APLTwisting,
	options: {
		aura: PaladinAura.SanctityAura,
		isUsingCrusaderStrikeStopAttack: true,
		isUsingDivineStormStopAttack: true,
		isUsingExorcismStopAttack: true,
		isUsingJudgementStopAttack: true,
		isUsingManualStartAttack: false,
		primarySeal: PaladinSeal.Martyrdom,
	},
});
export const PresetBuildP7SealStacking = PresetUtils.makePresetBuild('Seal Stacking', {
	gear: Phase7SealStackingNaxxGear,
	talents: RetTalents,
	rotation: APLSealStacking,
	options: {
		aura: PaladinAura.SanctityAura,
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
		primarySeal: PaladinSeal.Martyrdom,
	},
});
export const PresetBuildP7Shockadin1H = PresetUtils.makePresetBuild('1H Shockadin', {
	gear: Phase7Shockadin1HNaxxGear,
	talents: ShockadinTalents,
	rotation: APLShockadin1H,
	options: {
		aura: PaladinAura.NoPaladinAura,
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
		primarySeal: PaladinSeal.Righteousness,
	},
});
export const PresetBuildP7Shockadin2H = PresetUtils.makePresetBuild('2H Shockadin', {
	gear: Phase7Shockadin2HNaxxGear,
	talents: ShockadinTalents,
	rotation: APLShockadin2H,
	options: {
		aura: PaladinAura.NoPaladinAura,
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
		primarySeal: PaladinSeal.Martyrdom,
	},
});
export const PresetBuildP7Exodin = PresetUtils.makePresetBuild('Exodin', {
	gear: Phase7ExodinNaxxGear,
	talents: RetTalents,
	rotation: APLExodin,
	options: {
		aura: PaladinAura.SanctityAura,
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
		primarySeal: PaladinSeal.Martyrdom,
	},
});

export const DefaultTalents = TalentPresets[Phase.Phase7][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.SanctityAura,
	primarySeal: PaladinSeal.Martyrdom,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	fillerExplosive: Explosive.ExplosiveStratholmeHolyWater,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodSmokedDesertDumpling,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.WildStrikes,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	mildlyIrradiatedRejuvPot: true,
	miscConsumes: {
		boglingRoot: false,
		elixirOfCoalescedRegret: true,
		greaterMarkOfTheDawn: true,
	},
	offHandImbue: WeaponImbue.EnchantedRepellent,
	sapperExplosive: SapperExplosive.SapperFumigator,
	sealOfTheDawn: SealOfTheDawn.SealOfTheDawnDamageR7,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
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

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	commandingShout: true,
	demonicPact: 150,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	hornOfLordaeron: true,
	leaderOfThePack: true,
	moonkinAura: false,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	sanctityAura: false,
	vampiricTouch: 0,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	giftOfArthas: true,
	homunculi: 70, // 70% average uptime default
	improvedScorch: true,
	judgementOfLight: true,
	judgementOfTheCrusader: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	mangle: true,
	markOfChaos: true,
	occultPoison: true,
});

export const OtherDefaults = {
	profession1: Profession.Alchemy,
	profession2: Profession.Engineering,
};
