import { Phase } from '../core/constants/other';
import * as PresetUtils from '../core/preset_utils';
import {
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	FirePowerBuff,
	Flask,
	Food,
	FrostPowerBuff,
	IndividualBuffs,
	MageScroll,
	ManaRegenElixir,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common';
import { Mage_Options as MageOptions, Mage_Options_ArmorType as ArmorType } from '../core/proto/mage';
import { SavedTalents } from '../core/proto/ui';
import Phase1APLArcane from './apls/p1_arcane.apl.json';
import Phase1APLFire from './apls/p1_fire.apl.json';
import Phase2APLArcane from './apls/p2_arcane.apl.json';
import Phase2APLFire from './apls/p2_fire.apl.json';
// import Phase3APLArcane from './apls/p3_arcane.apl.json';
import Phase3APLFire from './apls/p3_fire.apl.json';
import Phase3APLFrost from './apls/p3_frost.apl.json';
import Phase4APLArcane from './apls/p4_arcane.apl.json';
import Phase4APLFire from './apls/p4_fire.apl.json';
import Phase4APLFrost from './apls/p4_frost.apl.json';
import Phase5APLFire from './apls/p5_fire.apl.json';
import Phase5APLSpellfrost from './apls/p5_spellfrost.apl.json';
import Phase1GearFire from './gear_sets/p1_fire.gear.json';
import Phase1Gear from './gear_sets/p1_generic.gear.json';
import Phase2GearArcane from './gear_sets/p2_arcane.gear.json';
import Phase2GearFire from './gear_sets/p2_fire.gear.json';
import Phase2GearFrost from './gear_sets/p2_frost.gear.json';
import Phase3GearFire from './gear_sets/p3_fire.gear.json';
import Phase3GearFrostFFB from './gear_sets/p3_frost_ffb.gear.json';
import Phase4GearArcane from './gear_sets/p4_arcane.gear.json';
import Phase4GearFire from './gear_sets/p4_fire.gear.json';
import Phase4GearFrost from './gear_sets/p4_frost.gear.json';
import Phase5GearFire from './gear_sets/p5_fire.gear.json';
import Phase5GearSpellfrost from './gear_sets/p5_spellfrost.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArcanePhase1 = PresetUtils.makePresetGear('P1 Arcane', Phase1Gear, {
	customCondition: player => player.getLevel() === 25,
});
export const GearFirePhase1 = PresetUtils.makePresetGear('P1 Fire', Phase1GearFire, {
	customCondition: player => player.getLevel() === 25,
});
export const GearFrostPhase1 = PresetUtils.makePresetGear('P1 Frost', Phase1Gear, {
	customCondition: player => player.getLevel() === 25,
});

export const GearArcanePhase2 = PresetUtils.makePresetGear('P2 Arcane', Phase2GearArcane, {
	customCondition: player => player.getLevel() === 40,
});
export const GearFirePhase2 = PresetUtils.makePresetGear('P2 Fire', Phase2GearFire, {
	customCondition: player => player.getLevel() === 40,
});
export const GearFrostPhase2 = PresetUtils.makePresetGear('P2 Frost', Phase2GearFrost, {
	customCondition: player => player.getLevel() === 40,
});

// No new Phase 3 Arcane presets at the moment
export const GearArcanePhase3 = GearArcanePhase2;
export const GearFirePhase3 = PresetUtils.makePresetGear('P3 Fire', Phase3GearFire, {
	customCondition: player => player.getLevel() === 50,
});
export const GearFrostPhase3 = PresetUtils.makePresetGear('P3 Frost', Phase3GearFrostFFB, {
	customCondition: player => player.getLevel() === 50,
});

// No new Phase 4 Arcane presets at the moment
export const GearArcanePhase4 = PresetUtils.makePresetGear('P4 Arcane', Phase4GearArcane, {
	customCondition: player => player.getLevel() === 60,
});
export const GearFirePhase4 = PresetUtils.makePresetGear('P4 Fire', Phase4GearFire, {
	customCondition: player => player.getLevel() === 60,
});
export const GearFrostPhase4 = PresetUtils.makePresetGear('P4 Frost', Phase4GearFrost, {
	customCondition: player => player.getLevel() === 60,
});

export const GearFirePhase5 = PresetUtils.makePresetGear('P5 Fire', Phase5GearFire, {
	customCondition: player => player.getLevel() === 60,
});
export const GearSpellfrostPhase5 = PresetUtils.makePresetGear('P5 Spellfrost', Phase5GearSpellfrost, {
	customCondition: player => player.getLevel() === 60,
});

export const GearPresets = {
	[Phase.Phase1]: [GearArcanePhase1, GearFirePhase1, GearFrostPhase1],
	[Phase.Phase2]: [GearArcanePhase2, GearFirePhase2, GearFrostPhase2],
	[Phase.Phase3]: [GearArcanePhase3, GearFirePhase3, GearFrostPhase3],
	[Phase.Phase4]: [GearArcanePhase4, GearFirePhase4, GearFrostPhase4],
	[Phase.Phase5]: [GearFirePhase5, GearSpellfrostPhase5],
};

export const DefaultGearArcane = GearPresets[Phase.Phase5][1];
export const DefaultGearFire = GearPresets[Phase.Phase5][0];
export const DefaultGearFrost = GearPresets[Phase.Phase5][1];

export const DefaultGear = DefaultGearFire;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLArcanePhase1 = PresetUtils.makePresetAPLRotation('P1 Arcane', Phase1APLArcane, {
	customCondition: player => player.getLevel() === 25,
});
export const APLFirePhase1 = PresetUtils.makePresetAPLRotation('P1 Fire', Phase1APLFire, {
	customCondition: player => player.getLevel() === 25,
});

export const APLArcanePhase2 = PresetUtils.makePresetAPLRotation('P2 Arcane', Phase2APLArcane, {
	customCondition: player => player.getLevel() === 40,
});
export const APLFirePhase2 = PresetUtils.makePresetAPLRotation('P2 Fire', Phase2APLFire, {
	customCondition: player => player.getLevel() === 40,
});

// No new Phase 3 Arcane presets at the moment
export const APLArcanePhase3 = APLArcanePhase2;
export const APLFirePhase3 = PresetUtils.makePresetAPLRotation('P3 Fire', Phase3APLFire, {
	customCondition: player => player.getLevel() === 50,
});
export const APLFrostPhase3 = PresetUtils.makePresetAPLRotation('P3 Frost', Phase3APLFrost, {
	customCondition: player => player.getLevel() === 50,
});

export const APLArcanePhase4 = PresetUtils.makePresetAPLRotation('P4 Arcane', Phase4APLArcane, {
	customCondition: player => player.getLevel() >= 60,
});
export const APLFirePhase4 = PresetUtils.makePresetAPLRotation('P4 Fire', Phase4APLFire, {
	customCondition: player => player.getLevel() >= 60,
});
export const APLFrostPhase4 = PresetUtils.makePresetAPLRotation('P4 Frost', Phase4APLFrost, {
	customCondition: player => player.getLevel() >= 60,
});

export const APLFirePhase5 = PresetUtils.makePresetAPLRotation('P5 Fire', Phase5APLFire, {
	customCondition: player => player.getLevel() >= 60,
});
export const APLSpellfrostPhase5 = PresetUtils.makePresetAPLRotation('P5 Spellfrost', Phase5APLSpellfrost, {
	customCondition: player => player.getLevel() >= 60,
});

export const APLPresets = {
	[Phase.Phase1]: [APLArcanePhase1, APLFirePhase1, APLFirePhase1],
	[Phase.Phase2]: [APLArcanePhase2, APLFirePhase2, APLFirePhase2],
	[Phase.Phase3]: [APLArcanePhase3, APLFirePhase3, APLFrostPhase3],
	[Phase.Phase4]: [APLArcanePhase4, APLFirePhase4, APLFrostPhase4],
	[Phase.Phase5]: [APLFirePhase5, APLSpellfrostPhase5],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
		2: APLPresets[Phase.Phase1][2],
	},
	40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][1],
		// Normally frost but frost is unfortunately just too bad to warrant including for now
		2: APLPresets[Phase.Phase2][2],
		// Frostfire
		3: APLPresets[Phase.Phase2][2],
	},
	50: {
		0: APLPresets[Phase.Phase3][0],
		1: APLPresets[Phase.Phase3][1],
		2: APLPresets[Phase.Phase3][2],
	},
	60: {
		0: APLPresets[Phase.Phase5][1],
		1: APLPresets[Phase.Phase5][0],
		2: APLPresets[Phase.Phase5][1],
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsArcanePhase1 = PresetUtils.makePresetTalents('25 Arcane', SavedTalents.create({ talentsString: '22500502' }), {
	customCondition: player => player.getLevel() === 25,
});
export const TalentsFirePhase1 = PresetUtils.makePresetTalents('25 Fire', SavedTalents.create({ talentsString: '-5050020121' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsArcanePhase2 = PresetUtils.makePresetTalents('40 Arcane', SavedTalents.create({ talentsString: '2250050310031531' }), {
	customCondition: player => player.getLevel() === 40,
});
export const TalentsFirePhase2 = PresetUtils.makePresetTalents('40 Fire', SavedTalents.create({ talentsString: '-5050020123033151' }), {
	customCondition: player => player.getLevel() === 40,
});

// No new Phase 3 Arcane presets at the moment
export const TalentsArcanePhase3 = TalentsArcanePhase2;
export const TalentsFirePhase3 = PresetUtils.makePresetTalents('50 Fire', SavedTalents.create({ talentsString: '-0550020123033151-2035' }), {
	customCondition: player => player.getLevel() === 50,
});
export const TalentsFrostPhase3 = PresetUtils.makePresetTalents('50 Frost', SavedTalents.create({ talentsString: '-055-20350203100351051' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TalentsArcanePhase4 = PresetUtils.makePresetTalents('60 Arcane', SavedTalents.create({ talentsString: '0550050210031531-054-203500001' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsFirePhase4 = PresetUtils.makePresetTalents('60 Fire', SavedTalents.create({ talentsString: '21-5052300123033151-203500031' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsFrostfirePhase4 = PresetUtils.makePresetTalents('60 Frostfire', SavedTalents.create({ talentsString: '-0550320003021-2035020310035105' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsSpellfrostPhase5 = PresetUtils.makePresetTalents(
	'60 Spellfrost',
	SavedTalents.create({ talentsString: '250025001002--05350203100351051' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const TalentPresets = {
	[Phase.Phase1]: [TalentsArcanePhase1, TalentsFirePhase1, TalentsFirePhase1],
	[Phase.Phase2]: [TalentsArcanePhase2, TalentsFirePhase2, TalentsFirePhase2],
	[Phase.Phase3]: [TalentsArcanePhase3, TalentsFirePhase3, TalentsFrostPhase3],
	[Phase.Phase4]: [TalentsArcanePhase4, TalentsFirePhase4, TalentsFrostfirePhase4],
	[Phase.Phase5]: [TalentsSpellfrostPhase5],
};

export const DefaultTalentsArcane = TalentPresets[Phase.Phase4][0];
export const DefaultTalentsFire = TalentPresets[Phase.Phase4][1];
export const DefaultTalentsFrostfire = TalentPresets[Phase.Phase4][2];
export const DefaultTalentsSpellfrost = TalentPresets[Phase.Phase5][0];

export const DefaultTalents = DefaultTalentsFire;

// export const PresetBuildArcane = PresetUtils.makePresetBuild('Arcane', DefaultGearArcane, DefaultTalentsArcane, DefaultAPLs[60][0]);
export const PresetBuildFire = PresetUtils.makePresetBuild('Fire', DefaultGearFire, DefaultTalentsFire, DefaultAPLs[60][1]);
export const PresetBuildSpellfrost = PresetUtils.makePresetBuild('Spellfrost', DefaultGearFrost, DefaultTalentsSpellfrost, DefaultAPLs[60][2]);

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultConsumes = Consumes.create({
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfSupremePower,
	food: Food.FoodRunnTumTuberSurprise,
	frostPowerBuff: FrostPowerBuff.ElixirOfFrostPower,
	mageScroll: MageScroll.MageScrollArcanePower,
	mainHandImbue: WeaponImbue.BrillianWizardOil,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,

	mildlyIrradiatedRejuvPot: true,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	zanzaBuff: ZanzaBuff.CerebralCortexCompound,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	demonicPact: 110,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: true,
	vampiricTouch: 300,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	mightOfStormwind: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	spiritOfZandalar: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	dreamstate: true,
	improvedFaerieFire: true,
	improvedScorch: true,
	judgementOfWisdom: true,
	markOfChaos: true,
	occultPoison: true,
	wintersChill: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Alchemy,
	profession2: Profession.Tailoring,
};
