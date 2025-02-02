import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Explosive,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	RaidBuffs,
	SapperExplosive,
	SaygesFortune,
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
import APLP4ExodinJson from './apls/p4-exodin.apl.json';
import APLP4Exodin6PcT1Json from './apls/p4-exodin-6pcT1.apl.json';
import APLP4Twisting6PcT1Json from './apls/p4-twisting-6pcT1.apl.json';
import APLP5ExodinJson from './apls/p5-exodin-6CF-2DR.apl.json';
import APLP5SealStackingJson from './apls/p5-seal-stacking-6CF-2DR.apl.json';
import APLShockadinJson from './apls/p5-shockadin.apl.json';
import APLTwistingJson from './apls/p5p6p7-twist.apl.json';
import APLP6OneHandJson from './apls/p6-1h-3AQ10.apl.json';
import APLP6ExodinJson from './apls/p6-exodin.apl.json';
import APLTwistSlowJson from './apls/twist-slow-hold-filler.apl.json';
import Phase1RetGearJson from './gear_sets/p1-ret.gear.json';
import Phase2RetSoCGearJson from './gear_sets/p2-retsoc.gear.json';
import Phase2RetSoMGearJson from './gear_sets/p2-retsom.gear.json';
import Phase3RetSoMGearJson from './gear_sets/p3-retsom.gear.json';
import Phase4ExodinGearJson from './gear_sets/p4-exodin.gear.json';
import Phase4Exodin6PcT1GearJson from './gear_sets/p4-exodin-6pcT1.gear.json';
import Phase4TwistingGearJson from './gear_sets/p4-twist.gear.json';
import Phase4Twisting6PcT1GearJson from './gear_sets/p4-twisting-6pcT1.gear.json';
import Phase5ExodinGearJson from './gear_sets/p5-exodin.gear.json';
import Phase5SealStackingGearJson from './gear_sets/p5-seal-stacking.gear.json';
import Phase5ShockadinGearJson from './gear_sets/p5-shockadin.gear.json';
import Phase5TwistingGearJson from './gear_sets/p5-twisting.gear.json';
import Phase6OneHandGearJson from './gear_sets/p6-1h-3AQ10.gear.json';
import Phase6ExodinGearJson from './gear_sets/p6-exodin.gear.json';
import Phase6TwistingGearJson from './gear_sets/p6-twisting.gear.json';

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
export const Phase4TwistGear = PresetUtils.makePresetGear('P4 Twist', Phase4TwistingGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase4Twist6pT1Gear = PresetUtils.makePresetGear('P4 Twist 6pT1', Phase4Twisting6PcT1GearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase4ExodinGear = PresetUtils.makePresetGear('P4 Exodin', Phase4ExodinGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase4Exodin6pT1Gear = PresetUtils.makePresetGear('P4 Exodin 6pT1', Phase4Exodin6PcT1GearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5ExodinGear = PresetUtils.makePresetGear('P5 Exodin', Phase5ExodinGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5SealStackingGear = PresetUtils.makePresetGear('P5 Seal Stacking', Phase5SealStackingGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5ShockadinGear = PresetUtils.makePresetGear('P5 Shockadin', Phase5ShockadinGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5TwistingGear = PresetUtils.makePresetGear('P5 Twisting', Phase5TwistingGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase6ExodinGear = PresetUtils.makePresetGear('P6 Exodin', Phase6ExodinGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase6OneHandGear = PresetUtils.makePresetGear('P6 1H Ret', Phase6OneHandGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase6TwistingGear = PresetUtils.makePresetGear('P6 Twisting', Phase6TwistingGearJson, {
	customCondition: player => player.getLevel() == 60,
});

export const GearPresets = {
	[Phase.Phase1]: [Phase1RetGear],
	[Phase.Phase2]: [Phase2RetSoCGear, Phase2RetSoMGear],
	[Phase.Phase3]: [Phase3RetSoMGear],
	[Phase.Phase4]: [Phase4TwistGear, Phase4Twist6pT1Gear, Phase4ExodinGear, Phase4Exodin6pT1Gear],
	[Phase.Phase5]: [Phase5TwistingGear, Phase5ExodinGear, Phase5SealStackingGear, Phase5ShockadinGear],
	[Phase.Phase6]: [Phase6TwistingGear, Phase6ExodinGear, Phase6OneHandGear],
};

export const DefaultGear = GearPresets[Phase.Phase6][0];

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
export const APLP4RetTwist6pT1 = PresetUtils.makePresetAPLRotation('P4 Twist 6pT1', APLP4Twisting6PcT1Json, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP4RetExodin = PresetUtils.makePresetAPLRotation('P4 Exodin', APLP4ExodinJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP4RetExodin6pT1 = PresetUtils.makePresetAPLRotation('P4 Exodin 6pT1', APLP4Exodin6PcT1Json, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP5Exodin6CF = PresetUtils.makePresetAPLRotation('P5 Exodin 6CF', APLP5ExodinJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP5SealStacking6CF = PresetUtils.makePresetAPLRotation('P5 Seal Stacking 6CF', APLP5SealStackingJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLShockadin = PresetUtils.makePresetAPLRotation('P5 Shockadin', APLShockadinJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP6Exodin = PresetUtils.makePresetAPLRotation('P6 Exodin', APLP6ExodinJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP6OneHand = PresetUtils.makePresetAPLRotation('P6 1H Ret', APLP6OneHandJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLTwisting = PresetUtils.makePresetAPLRotation('Twist', APLTwistingJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLTwistingSlow = PresetUtils.makePresetAPLRotation('Twist Slow', APLTwistSlowJson, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [APLP1Ret],
	[Phase.Phase2]: [APLP2Ret],
	[Phase.Phase3]: [APLP3Ret],
	[Phase.Phase4]: [APLP4RetTwist6pT1, APLP4RetExodin, APLP4RetExodin6pT1],
	[Phase.Phase5]: [APLP5Exodin6CF, APLShockadin, APLP5SealStacking6CF],
	[Phase.Phase6]: [APLTwisting, APLTwistingSlow, APLP6Exodin, APLP6OneHand],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase6][0],
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

export const ShockadinTalents = PresetUtils.makePresetTalents('P5 Shockadin', SavedTalents.create({ talentsString: '55053100501051--052303511' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [P1RetTalents],
	[Phase.Phase2]: [P2RetTalents, P2ShockadinTalents],
	[Phase.Phase3]: [P3RetTalents],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [ShockadinTalents],
	[Phase.Phase6]: [RetTalents],
};

export const PresetBuildTwisting = PresetUtils.makePresetBuild('Twisting', {
	gear: Phase6TwistingGear,
	talents: RetTalents,
	rotation: APLTwisting,
	options: {
		isUsingCrusaderStrikeStopAttack: true,
		isUsingDivineStormStopAttack: true,
		isUsingExorcismStopAttack: true,
		isUsingJudgementStopAttack: true,
		isUsingManualStartAttack: false,
	},
});
export const PresetBuildP5SealStacking = PresetUtils.makePresetBuild('P5 Seal Stacking', {
	gear: Phase5SealStackingGear,
	talents: RetTalents,
	rotation: APLP5SealStacking6CF,
	options: {
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
	},
});
export const PresetBuildP5Exodin = PresetUtils.makePresetBuild('P5 Exodin', {
	gear: Phase5ExodinGear,
	talents: RetTalents,
	rotation: APLP5Exodin6CF,
	options: {
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
	},
});
export const PresetBuildP5Shockadin = PresetUtils.makePresetBuild('P5 Shockadin', {
	gear: Phase5ShockadinGear,
	talents: ShockadinTalents,
	rotation: APLShockadin,
	options: {
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
	},
});
export const PresetBuildP6Exodin = PresetUtils.makePresetBuild('P6 Exodin', {
	gear: Phase6ExodinGear,
	talents: RetTalents,
	rotation: APLP6Exodin,
	options: {
		isUsingCrusaderStrikeStopAttack: false,
		isUsingDivineStormStopAttack: false,
		isUsingExorcismStopAttack: false,
		isUsingJudgementStopAttack: false,
		isUsingManualStartAttack: false,
	},
});

export const DefaultTalents = TalentPresets[Phase.Phase6][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.SanctityAura,
	primarySeal: PaladinSeal.Martyrdom,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	fillerExplosive: Explosive.ExplosiveUnknown,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfMadness,
	food: Food.FoodSmokedDesertDumpling,
	mainHandImbue: WeaponImbue.WildStrikes,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	mildlyIrradiatedRejuvPot: true,
	miscConsumes: {
		boglingRoot: false,
	},
	offHandImbue: WeaponImbue.EnchantedRepellent,
	sapperExplosive: SapperExplosive.SapperFumigator,
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
	demonicPact: 120,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	hornOfLordaeron: true,
	leaderOfThePack: true,
	moonkinAura: true,
	sanctityAura: true,
	vampiricTouch: 300,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	giftOfArthas: true,
	homunculi: 70, // 70% average uptime default
	improvedScorch: true,
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
