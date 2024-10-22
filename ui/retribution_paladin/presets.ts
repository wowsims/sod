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
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { PaladinAura, PaladinOptions as RetributionPaladinOptions, PaladinSeal } from '../core/proto/paladin.js';
import { SavedTalents } from '../core/proto/ui.js';
import APLP1RetJson from './apls/p1ret.apl.json';
import APLP2RetJson from './apls/p2ret.apl.json';
import APLP3RetJson from './apls/p3ret.apl.json';
import APLP4RetJson from './apls/p4ret.apl.json';
import APLP4RetExodinJson from './apls/p4ret-exodin.apl.json';
import APLP4RetExodin6PcT1Json from './apls/p4ret-exodin-6pcT1.apl.json';
import APLP4RetTwisting6PcT1Json from './apls/p4ret-twisting-6pcT1.apl.json';
import APLPP5ExodinJson from './apls/p5ret-exodin-6CF2DR.apl.json';
import APLPP5TwistingSlowJson from './apls/p5ret-twist-4DR-3.5-3.6.apl.json';
import APLPP5TwistingSlowerJson from './apls/p5ret-twist-4DR-3.7-4.0.apl.json';
import APLPP5TwistingCancelAuraJson from './apls/p5ret-twist-4DR-CancelAura.apl.json';
import APLPP5StackingJson from './apls/p5ret-stacking-6CF2DR.apl.json';
import APLPP5ShockadinJson from './apls/p5Shockadin.apl.json';
import Phase1RetGearJson from './gear_sets/p1ret.gear.json';
import Phase2RetSoCGearJson from './gear_sets/p2retsoc.gear.json';
import Phase2RetSoMGearJson from './gear_sets/p2retsom.gear.json';
import Phase3RetSoMGearJson from './gear_sets/p3retsom.gear.json';
import Phase4RetExodinGearJson from './gear_sets/p4ret-exodin.gear.json';
import Phase4RetExodin6PcT1GearJson from './gear_sets/p4ret-exodin-6pcT1.gear.json';
import Phase4RetTwisting6PcT1GearJson from './gear_sets/p4ret-twisting-6pcT1.gear.json';
import Phase4RetGearJson from './gear_sets/p4rettwist.gear.json';
import Phase5ExodinGearJson from './gear_sets/p5exodin.gear.json';
import Phase5ShockadinGearJson from './gear_sets/p5shockadin.gear.json';
import Phase5TwistingGearJson from './gear_sets/p5twisting.gear.json';
import Phase5TwistingSlowerGearJson from './gear_sets/p5twistingSlower.gear.json';
import Phase5StackingGearJson from './gear_sets/p5stacking.gear.json';
import Phase5TwistingHasteGearJson from './gear_sets/p5twistingHaste.gear.json';

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
export const Phase4RetTwistGear = PresetUtils.makePresetGear('P4 Twist', Phase4RetGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase4RetTwist6pT1Gear = PresetUtils.makePresetGear('P4 Ret Twist 6pT1', Phase4RetTwisting6PcT1GearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase4RetExodinGear = PresetUtils.makePresetGear('P4 Ret Exodin', Phase4RetExodinGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase4RetExodin6pT1Gear = PresetUtils.makePresetGear('P4 Ret Exodin 6pT1', Phase4RetExodin6PcT1GearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5TwistingGear = PresetUtils.makePresetGear('P5 Twisting', Phase5TwistingGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5TwistingSlowerGear = PresetUtils.makePresetGear('P5 Twisting Slower', Phase5TwistingSlowerGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5StackingGear = PresetUtils.makePresetGear('P5 Stacking', Phase5StackingGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5TwistingHasteGear = PresetUtils.makePresetGear('P5 Twisting Haste', Phase5TwistingHasteGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5ExodinGear = PresetUtils.makePresetGear('P5 Exodin', Phase5ExodinGearJson, {
	customCondition: player => player.getLevel() == 60,
});
export const Phase5ShockadinGear = PresetUtils.makePresetGear('P5 Shockadin', Phase5ShockadinGearJson, {
	customCondition: player => player.getLevel() == 60,
});

export const GearPresets = {
	[Phase.Phase1]: [Phase1RetGear],
	[Phase.Phase2]: [Phase2RetSoCGear, Phase2RetSoMGear],
	[Phase.Phase3]: [Phase3RetSoMGear],
	[Phase.Phase4]: [Phase4RetTwistGear, Phase4RetTwist6pT1Gear, Phase4RetExodinGear, Phase4RetExodin6pT1Gear],
	[Phase.Phase5]: [Phase5TwistingGear, Phase5TwistingSlowerGear, Phase5StackingGear, Phase5TwistingHasteGear, Phase5ExodinGear, Phase5ShockadinGear],
};

export const DefaultGear = GearPresets[Phase.Phase5][0];

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
export const APLP4RetTwist = PresetUtils.makePresetAPLRotation('P4 Ret Twist', APLP4RetJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP4RetTwist6pT1 = PresetUtils.makePresetAPLRotation('P4 Ret Twist 6pT1', APLP4RetTwisting6PcT1Json, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP4RetExodin = PresetUtils.makePresetAPLRotation('P4 Ret Exodin', APLP4RetExodinJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLP4RetExodin6pT1 = PresetUtils.makePresetAPLRotation('P4 Ret Exodin 6pT1', APLP4RetExodin6PcT1Json, {
	customCondition: player => player.getLevel() === 60,
});
export const APLPP5Twisting4DRSlow = PresetUtils.makePresetAPLRotation('P5 Twist 4DR Slow 3.5-3.6', APLPP5TwistingSlowJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLPP5Twisting4DRSlower = PresetUtils.makePresetAPLRotation('P5 Twist 4DR Slower 3.7+', APLPP5TwistingSlowerJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLPP5Twisting4DRCancelAura = PresetUtils.makePresetAPLRotation('P5 Twist 4DR CancelAura', APLPP5TwistingCancelAuraJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLPP5Stacking6CF2DR = PresetUtils.makePresetAPLRotation('P5 Twist 4DR CancelAura', APLPP5StackingJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLPP5Exodin = PresetUtils.makePresetAPLRotation('P5 Exodin', APLPP5ExodinJson, {
	customCondition: player => player.getLevel() === 60,
});
export const APLPP5Shockadin = PresetUtils.makePresetAPLRotation('P5 Shockadin', APLPP5ShockadinJson, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [APLP1Ret],
	[Phase.Phase2]: [APLP2Ret],
	[Phase.Phase3]: [APLP3Ret],
	[Phase.Phase4]: [APLP4RetTwist, APLP4RetTwist6pT1, APLP4RetExodin, APLP4RetExodin6pT1],
	[Phase.Phase5]: [APLPP5Twisting4DRCancelAura, APLPP5Stacking6CF2DR, APLPP5Twisting4DRSlow, APLPP5Twisting4DRSlower, APLPP5Exodin, APLPP5Shockadin],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase5][0],
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

export const P4RetTalents = PresetUtils.makePresetTalents('P4/P5 Ret', SavedTalents.create({ talentsString: '500501-503-52230351200315' }), {
	customCondition: player => player.getLevel() === 60,
});

export const P5ShockadinTalents = PresetUtils.makePresetTalents('P5 Shockadin', SavedTalents.create({ talentsString: '55053100501051--052303511' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [P1RetTalents],
	[Phase.Phase2]: [P2RetTalents, P2ShockadinTalents],
	[Phase.Phase3]: [P3RetTalents],
	[Phase.Phase4]: [P4RetTalents],
	[Phase.Phase5]: [P4RetTalents, P5ShockadinTalents],
};

export const PresetBuildTwistingCancelAura = PresetUtils.makePresetBuild('TwistingCancelAura', { gear: Phase5TwistingHasteGear, talents: P4RetTalents, rotation: APLPP5Twisting4DRCancelAura });
export const PresetBuildTwistingSlow = PresetUtils.makePresetBuild('TwistingSlow', { gear: Phase5TwistingGear, talents: P4RetTalents, rotation: APLPP5Twisting4DRSlow });
export const PresetBuildTwistingSlower = PresetUtils.makePresetBuild('TwistingSlower', { gear: Phase5TwistingSlowerGear, talents: P4RetTalents, rotation: APLPP5Twisting4DRSlower });
export const PresetBuildSealStacking = PresetUtils.makePresetBuild('SealStacking', { gear: Phase5StackingGear, talents: P4RetTalents, rotation: APLPP5Stacking6CF2DR });
export const PresetBuildExodin = PresetUtils.makePresetBuild('Exodin', { gear: Phase5ExodinGear, talents: P4RetTalents, rotation: APLPP5Exodin });
export const PresetBuildShockadin = PresetUtils.makePresetBuild('Shockadin', { gear: Phase5ShockadinGear, talents: P5ShockadinTalents, rotation: APLPP5Shockadin });


export const DefaultTalents = TalentPresets[Phase.Phase5][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.SanctityAura,
	primarySeal: PaladinSeal.Martyrdom,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	boglingRoot: false,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	fillerExplosive: Explosive.ExplosiveUnknown,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	food: Food.FoodBlessSunfruit,
	flask: Flask.FlaskOfSupremePower,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.MagnificentTrollshine,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfKings: true,
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
	battleShout: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	sanctityAura: true,
	hornOfLordaeron: true,
	leaderOfThePack: true,
	demonicPact: 110,
	aspectOfTheLion: true,
	moonkinAura: true,
	vampiricTouch: 300,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	homunculi: 70, // 70% average uptime default
	faerieFire: true,
	giftOfArthas: true,
	sunderArmor: true,
	judgementOfWisdom: true,
	judgementOfTheCrusader: TristateEffect.TristateEffectImproved,
	improvedFaerieFire: true,
	improvedScorch: true,
	markOfChaos: true,
	occultPoison: true,
	mangle: true,
});

export const OtherDefaults = {
	profession1: Profession.Blacksmithing,
	profession2: Profession.Enchanting,
};
