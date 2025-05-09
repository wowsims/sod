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
	Spec,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { PaladinAura, PaladinOptions as RetributionPaladinOptions, PaladinSeal } from '../core/proto/paladin.js';
import { SavedTalents } from '../core/proto/ui.js';
// APLs
import Phase1APLRetJSON from './apls/p1-ret.apl.json';
import Phase2APLRetJSON from './apls/p2-ret.apl.json';
import Phase3APLRetJSON from './apls/p3-ret.apl.json';
import Phase8APLExodinJSON from './apls/p8-exodin.apl.json';
import Phase8APLShockadin1hJSON from './apls/p8-shockadin-1h.apl.json';
import Phase8APLShockadin2hJSON from './apls/p8-shockadin-2h.apl.json';
import Phase8APLSealStackingJSON from './apls/p8-stack.apl.json';
import Phase8APLTwistingJSON from './apls/p8-twist.apl.json';
import Phase8APLWrathLikeJSON from './apls/p8-wrath.apl.json';
// Builds
import Phase8BuildExodinJSON from './builds/p8-exodin.build.json';
import Phase8BuildShockadin1hJSON from './builds/p8-shockadin-1h.build.json';
import Phase8BuildShockadin2hJSON from './builds/p8-shockadin-2h.build.json';
import Phase8BuildSealStackingJSON from './builds/p8-stack.build.json';
import Phase8BuildTwistingJSON from './builds/p8-twist.build.json';
import Phase8BuildWrathLikeJSON from './builds/p8-wrath.build.json';
// Gear
import Phase1GearRetJSON from './gear_sets/p1-ret.gear.json';
import Phase2GearRetSoCJSON from './gear_sets/p2-retsoc.gear.json';
import Phase2GearRetSoMJSON from './gear_sets/p2-retsom.gear.json';
import Phase3GearRetSoMJSON from './gear_sets/p3-retsom.gear.json';
import Phase8GearExodinJSON from './gear_sets/p8-exodin.gear.json';
import Phase8GearShockadin1hJSON from './gear_sets/p8-shockadin-1h.gear.json';
import Phase8GearShockadin2hJSON from './gear_sets/p8-shockadin-2h.gear.json';
import Phase8GearSealStackingJSON from './gear_sets/p8-stack.gear.json';
import Phase8GearTwistingJSON from './gear_sets/p8-twist.gear.json';
import Phase8GearWrathLikeJSON from './gear_sets/p8-wrath.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearRetPhase1 = PresetUtils.makePresetGear('P1', Phase1GearRetJSON, {
	customCondition: player => player.getLevel() == 25,
});
export const GearRetSoCPhase2 = PresetUtils.makePresetGear('P2 SoC/DS', Phase2GearRetSoCJSON, {
	customCondition: player => player.getLevel() == 40,
});
export const GearRetSoMPhase2 = PresetUtils.makePresetGear('P2 SoM', Phase2GearRetSoMJSON, {
	customCondition: player => player.getLevel() == 40,
});
export const GearRetSoMPhase3 = PresetUtils.makePresetGear('P3 SoM', Phase3GearRetSoMJSON, {
	customCondition: player => player.getLevel() == 50,
});
export const GearExodinPhase8 = PresetUtils.makePresetGear('P8 Exodin', Phase8GearExodinJSON, {
	customCondition: player => player.getLevel() == 60,
});
export const GearSealStackingPhase8 = PresetUtils.makePresetGear('P8 Seal Stacking', Phase8GearSealStackingJSON, {
	customCondition: player => player.getLevel() == 60,
});
export const GearShockadin1hPhase8 = PresetUtils.makePresetGear('P8 Shockadin 1h', Phase8GearShockadin1hJSON, {
	customCondition: player => player.getLevel() == 60,
});
export const GearShockadin2hPhase8 = PresetUtils.makePresetGear('P8 Shockadin 2h', Phase8GearShockadin2hJSON, {
	customCondition: player => player.getLevel() == 60,
});
export const GearTwistingPhase8 = PresetUtils.makePresetGear('P8 Twisting', Phase8GearTwistingJSON, {
	customCondition: player => player.getLevel() == 60,
});
export const GearWrathLikePhase8 = PresetUtils.makePresetGear('P8 Wrath-like', Phase8GearWrathLikeJSON, {
	customCondition: player => player.getLevel() == 60,
});

export const GearPresets = {
	[Phase.Phase1]: [GearRetPhase1],
	[Phase.Phase2]: [GearRetSoCPhase2, GearRetSoMPhase2],
	[Phase.Phase3]: [GearRetSoMPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
	[Phase.Phase7]: [],
	[Phase.Phase8]: [GearExodinPhase8, GearSealStackingPhase8, GearShockadin1hPhase8, GearShockadin2hPhase8, GearTwistingPhase8, GearWrathLikePhase8],
};

export const DefaultGear = GearPresets[Phase.Phase8][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLRetPhase1 = PresetUtils.makePresetAPLRotation('P1 Ret', Phase1APLRetJSON, {
	customCondition: player => player.getLevel() === 25,
});
export const APLRetPhase2 = PresetUtils.makePresetAPLRotation('P2 Ret/Shockadin', Phase2APLRetJSON, {
	customCondition: player => player.getLevel() === 40,
});
export const APLRetPhase3 = PresetUtils.makePresetAPLRotation('P3 Ret/Shockadin', Phase3APLRetJSON, {
	customCondition: player => player.getLevel() === 50,
});
export const APLExodinPhase8 = PresetUtils.makePresetAPLRotation('P8 Exodin', Phase8APLExodinJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const APLSealStackingPhase8 = PresetUtils.makePresetAPLRotation('P8 Seal Stacking', Phase8APLSealStackingJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const APLShockadin1hPhase8 = PresetUtils.makePresetAPLRotation('P8 Shockadin 1h', Phase8APLShockadin1hJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const APLShockadin2hPhase8 = PresetUtils.makePresetAPLRotation('P8 Shockadin 2h', Phase8APLShockadin2hJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const APLTwistingPhase8 = PresetUtils.makePresetAPLRotation('P8 Twisting', Phase8APLTwistingJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const APLWrathLikePhase8 = PresetUtils.makePresetAPLRotation('P8 Wrath-like', Phase8APLWrathLikeJSON, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [APLRetPhase1],
	[Phase.Phase2]: [APLRetPhase2],
	[Phase.Phase3]: [APLRetPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
	[Phase.Phase7]: [],
	[Phase.Phase8]: [APLExodinPhase8, APLSealStackingPhase8, APLShockadin1hPhase8, APLShockadin2hPhase8, APLTwistingPhase8, APLWrathLikePhase8],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase8][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsRetPhase1 = PresetUtils.makePresetTalents('P1 Ret', SavedTalents.create({ talentsString: '--05230051' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsRetPhase2 = PresetUtils.makePresetTalents('P2 Ret', SavedTalents.create({ talentsString: '--532300512003151' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsShockadinPhase2 = PresetUtils.makePresetTalents('P2 Shockadin', SavedTalents.create({ talentsString: '55050100521151--' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsRetPhase3 = PresetUtils.makePresetTalents('P3 Ret', SavedTalents.create({ talentsString: '500501--53230051200315' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TalentsRetPhase8 = PresetUtils.makePresetTalents('P8 Ret', SavedTalents.create({ talentsString: '500501-503-52230351200315' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentsShockadin1hPhase8 = PresetUtils.makePresetTalents('P8 Shockadin', SavedTalents.create({ talentsString: '55053100501051--052303502' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsRetPhase1],
	[Phase.Phase2]: [TalentsRetPhase2, TalentsShockadinPhase2],
	[Phase.Phase3]: [TalentsRetPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
	[Phase.Phase7]: [],
	[Phase.Phase8]: [TalentsRetPhase8, TalentsShockadin1hPhase8],
};

export const PresetBuildExodinPhase8 = PresetUtils.makePresetBuildFromJSON('P8 Exodin', Spec.SpecRetributionPaladin, Phase8BuildExodinJSON);
export const PresetBuildSealStackingPhase8 = PresetUtils.makePresetBuildFromJSON('P8 Seal Stacking', Spec.SpecRetributionPaladin, Phase8BuildSealStackingJSON);
export const PresetBuildShockadin1hPhase8 = PresetUtils.makePresetBuildFromJSON('P8 Shockadin 1h', Spec.SpecRetributionPaladin, Phase8BuildShockadin1hJSON);
export const PresetBuildShockadin2hPhase8 = PresetUtils.makePresetBuildFromJSON('P8 Shockadin 2h', Spec.SpecRetributionPaladin, Phase8BuildShockadin2hJSON);
export const PresetBuildTwistingPhase8 = PresetUtils.makePresetBuildFromJSON('P8 Twisting', Spec.SpecRetributionPaladin, Phase8BuildTwistingJSON);
export const PresetBuildWrathLikePhase8 = PresetUtils.makePresetBuildFromJSON('P8 Wrath-like', Spec.SpecRetributionPaladin, Phase8BuildWrathLikeJSON);

export const DefaultTalents = TalentPresets[Phase.Phase8][0];

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
	fillerExplosive: Explosive.ExplosiveUnknown,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodProwlerSteak,
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
	sealOfTheDawn: SealOfTheDawn.SealOfTheDawnUnknown,
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
