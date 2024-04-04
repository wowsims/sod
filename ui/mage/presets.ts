import { Phase } from '../core/constants/other';
import * as PresetUtils from '../core/preset_utils';
import {
	Consumes,
	Debuffs,
	EnchantedSigil,
	FirePowerBuff,
	Flask,
	Food,
	FrostPowerBuff,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue,
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
import Phase1GearFire from './gear_sets/p1_fire.gear.json';
import Phase1Gear from './gear_sets/p1_generic.gear.json';
import Phase2GearArcane from './gear_sets/p2_arcane.gear.json';
import Phase2GearFire from './gear_sets/p2_fire.gear.json';
import Phase2GearFrost from './gear_sets/p2_frost.gear.json';
import Phase3GearFire from './gear_sets/p3_fire.gear.json';
import Phase3GearFireFFB from './gear_sets/p3_fire_ffb.gear.json';
import Phase3GearFrostFFB from './gear_sets/p3_frost_ffb.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArcanePhase1 = PresetUtils.makePresetGear('P1 Arcane', Phase1Gear, {
	customCondition: player => player.getLevel() == 25,
});
export const GearFirePhase1 = PresetUtils.makePresetGear('P1 Fire', Phase1GearFire, {
	customCondition: player => player.getLevel() == 25,
});
export const GearFrostPhase1 = PresetUtils.makePresetGear('P1 Frost', Phase1Gear, {
	customCondition: player => player.getLevel() == 25,
});

export const GearArcanePhase2 = PresetUtils.makePresetGear('P2 Arcane', Phase2GearArcane, {
	customCondition: player => player.getLevel() == 40,
});
export const GearFirePhase2 = PresetUtils.makePresetGear('P2 Fire', Phase2GearFire, {
	customCondition: player => player.getLevel() == 40,
});
export const GearFrostPhase2 = PresetUtils.makePresetGear('P2 Frost', Phase2GearFrost, {
	customCondition: player => player.getLevel() == 40,
});

// No new Phase 3 Arcane presets at the moment
export const GearArcanePhase3 = GearArcanePhase2;
export const GearFirePhase3 = PresetUtils.makePresetGear('P3 Fire', Phase3GearFire, {
	customCondition: player => player.getLevel() >= 50,
});
export const GearFrostPhase3 = PresetUtils.makePresetGear('P3 Frost', Phase3GearFrostFFB, {
	customCondition: player => player.getLevel() >= 50,
});
export const GearFrostfirePhase3 = PresetUtils.makePresetGear('P3 Fire FFB', Phase3GearFireFFB, {
	customCondition: player => player.getLevel() >= 50,
});

export const GearPresets = {
	[Phase.Phase1]: [GearArcanePhase1, GearFirePhase1, GearFrostPhase1],
	[Phase.Phase2]: [GearArcanePhase2, GearFirePhase2, GearFrostPhase2],
	[Phase.Phase3]: [GearArcanePhase3, GearFirePhase3, GearFrostPhase3, GearFrostfirePhase3],
};

export const DefaultGearArcane = GearPresets[Phase.Phase3][0];
export const DefaultGearFire = GearPresets[Phase.Phase3][1];
export const DefaultGearFrost = GearPresets[Phase.Phase3][2];
export const DefaultGearFrostfire = GearPresets[Phase.Phase3][3];

export const DefaultGear = DefaultGearFire;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLArcanePhase1 = PresetUtils.makePresetAPLRotation('P1 Arcane', Phase1APLArcane, {
	customCondition: player => player.getLevel() == 25,
});
export const APLFirePhase1 = PresetUtils.makePresetAPLRotation('P1 Fire', Phase1APLFire, {
	customCondition: player => player.getLevel() == 25,
});

export const APLArcanePhase2 = PresetUtils.makePresetAPLRotation('P2 Arcane', Phase2APLArcane, {
	customCondition: player => player.getLevel() == 40,
});
export const APLFirePhase2 = PresetUtils.makePresetAPLRotation('P2 Fire', Phase2APLFire, {
	customCondition: player => player.getLevel() == 40,
});

// No new Phase 3 Arcane presets at the moment
export const APLArcanePhase3 = APLArcanePhase2;
export const APLFirePhase3 = PresetUtils.makePresetAPLRotation('P3 Fire', Phase3APLFire, {
	customCondition: player => player.getLevel() >= 50,
});
export const APLFrostPhase3 = PresetUtils.makePresetAPLRotation('P3 Frost', Phase3APLFrost, {
	customCondition: player => player.getLevel() >= 50,
});

export const APLPresets = {
	[Phase.Phase1]: [APLArcanePhase1, APLFirePhase1, APLFirePhase1],
	[Phase.Phase2]: [APLArcanePhase2, APLFirePhase2, APLFirePhase2],
	[Phase.Phase3]: [APLArcanePhase3, APLFirePhase3, APLFrostPhase3],
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
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsArcanePhase1 = PresetUtils.makePresetTalents('P1 Arcane', SavedTalents.create({ talentsString: '22500502' }), {
	customCondition: player => player.getLevel() == 25,
});
export const TalentsFirePhase1 = PresetUtils.makePresetTalents('P1 Fire', SavedTalents.create({ talentsString: '-5050020121' }), {
	customCondition: player => player.getLevel() == 25,
});

export const TalentsArcanePhase2 = PresetUtils.makePresetTalents('P2 Arcane', SavedTalents.create({ talentsString: '2250050310031531' }), {
	customCondition: player => player.getLevel() == 40,
});
export const TalentsFirePhase2 = PresetUtils.makePresetTalents('P2 Fire', SavedTalents.create({ talentsString: '-5050020123033151' }), {
	customCondition: player => player.getLevel() == 40,
});

// No new Phase 3 Arcane presets at the moment
export const TalentsArcanePhase3 = TalentsArcanePhase2;
export const TalentsFirePhase3 = PresetUtils.makePresetTalents('P3 Fire', SavedTalents.create({ talentsString: '-0550020123033151-2035' }), {
	customCondition: player => player.getLevel() >= 50,
});
export const TalentsFrostPhase3 = PresetUtils.makePresetTalents('P3 Frost', SavedTalents.create({ talentsString: '-055-20350203100351051' }), {
	customCondition: player => player.getLevel() >= 50,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsArcanePhase1, TalentsFirePhase1, TalentsFirePhase1],
	[Phase.Phase2]: [TalentsArcanePhase2, TalentsFirePhase2, TalentsFirePhase2],
	[Phase.Phase3]: [TalentsArcanePhase3, TalentsFirePhase3, TalentsFrostPhase3],
};

export const DefaultTalentsArcane = TalentPresets[Phase.Phase3][0];
export const DefaultTalentsFire = TalentPresets[Phase.Phase3][1];
export const DefaultTalentsFrost = TalentPresets[Phase.Phase3][2];
// export const DefaultTalentsFrostfire = TalentPresets[Phase.Phase3][2];

export const DefaultTalents = DefaultTalentsFire;

// export const PresetBuildArcane = PresetUtils.makePresetBuild('Arcane', DefaultGearArcane, DefaultTalentsArcane, DefaultAPLs[50][0]);
export const PresetBuildFire = PresetUtils.makePresetBuild('Fire', DefaultGearFire, DefaultTalentsFire, DefaultAPLs[50][1]);
export const PresetBuildFrost = PresetUtils.makePresetBuild('Frost FFB', DefaultGearFrost, DefaultTalentsFrost, DefaultAPLs[50][2]);
export const PresetBuildFrostfire = PresetUtils.makePresetBuild('Fire FFB', DefaultGearFrostfire, DefaultTalentsFire, DefaultAPLs[50][1]);

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.GreaterManaPotion,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfRestlessDreams,
	food: Food.FoodSagefishDelight,
	frostPowerBuff: FrostPowerBuff.ElixirOfFrostPower,
	mainHandImbue: WeaponImbue.BrillianWizardOil,
	mildlyIrradiatedRejuvPot: true,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	saygesFortune: SaygesFortune.SaygesDamage,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Alchemy,
	profession2: Profession.Enchanting,
};
