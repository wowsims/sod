import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	Alcohol,
	ArmorElixir,
	AttackPowerBuff,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Food,
	HealthElixir,
	IndividualBuffs,
	Potions,
	Profession,
	Race,
	RaidBuffs,
	SaygesFortune,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Warrior_Options as WarriorOptions, WarriorShout, WarriorStance } from '../core/proto/warrior.js';
import Phase1APLArms from './apls/phase_1_arms.apl.json';
import Phase2APLArms from './apls/phase_2_arms.apl.json';
import Phase2APLFury from './apls/phase_2_fury.apl.json';
import Phase3APLArms from './apls/phase_3_arms.apl.json';
import Phase3APLFury from './apls/phase_3_fury.apl.json';
import Phase3APLGlad from './apls/phase_3_glad.apl.json';
import Phase4APLFury from './apls/phase_4_fury.apl.json';
import Phase4APLGlad from './apls/phase_4_glad.apl.json';
import Phase5APL2H from './apls/phase_5_2h.apl.json';
import Phase5APLDW from './apls/phase_5_dw.apl.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase1DWGear from './gear_sets/phase_1_dw.gear.json';
import Phase22HGear from './gear_sets/phase_2_2h.gear.json';
import Phase2DWGear from './gear_sets/phase_2_dw.gear.json';
import Phase32HGear from './gear_sets/phase_3_2h.gear.json';
import Phase3DWGear from './gear_sets/phase_3_dw.gear.json';
import Phase3GladGear from './gear_sets/phase_3_glad.gear.json';
import Phase42HGear from './gear_sets/phase_4_2h.gear.json';
import Phase4DWGear from './gear_sets/phase_4_dw.gear.json';
import Phase4GladGear from './gear_sets/phase_4_glad.gear.json';
import Phase52HCoreForgedGear from './gear_sets/phase_5_2h_t1.gear.json';
import Phase52HDraconicGear from './gear_sets/phase_5_2h_t2.gear.json';
import Phase5DWCoreForgedGear from './gear_sets/phase_5_dw_t1.gear.json';
import Phase5DWDraconicGear from './gear_sets/phase_5_dw_t2.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArmsPhase1 = PresetUtils.makePresetGear('P1 Arms 2H', Phase1Gear, { customCondition: player => player.getLevel() === 25 });
export const GearArmsDWPhase1 = PresetUtils.makePresetGear('P1 Arms DW', Phase1DWGear, { customCondition: player => player.getLevel() === 25 });
export const GearFuryPhase1 = PresetUtils.makePresetGear('P1 DW Fury', Phase1Gear, { customCondition: player => player.getLevel() === 25 });

export const GearArmsPhase2 = PresetUtils.makePresetGear('P2 2H', Phase22HGear, { customCondition: player => player.getLevel() === 40 });
export const GearFuryPhase2 = PresetUtils.makePresetGear('P2 DW', Phase2DWGear, { customCondition: player => player.getLevel() === 40 });

export const GearArmsPhase3 = PresetUtils.makePresetGear('P3 2H', Phase32HGear, { customCondition: player => player.getLevel() === 50 });
export const GearFuryPhase3 = PresetUtils.makePresetGear('P3 DW', Phase3DWGear, { customCondition: player => player.getLevel() === 50 });
export const GearGladPhase3 = PresetUtils.makePresetGear('P3 Glad', Phase3GladGear, { customCondition: player => player.getLevel() === 50 });

export const Gear2HPhase4 = PresetUtils.makePresetGear('P4 2H', Phase42HGear, { customCondition: player => player.getLevel() === 60 });
export const GearDWPhase4 = PresetUtils.makePresetGear('P4 DW', Phase4DWGear, { customCondition: player => player.getLevel() === 60 });
export const GearGladPhase4 = PresetUtils.makePresetGear('P4 Glad', Phase4GladGear, { customCondition: player => player.getLevel() === 60 });

export const Gear2HCoreForgedPhase5 = PresetUtils.makePresetGear('P5 2H Core Forged', Phase52HCoreForgedGear, {
	customCondition: player => player.getLevel() === 60,
});
export const GearDWCoreForgedPhase5 = PresetUtils.makePresetGear('P5 DW Core Forged', Phase5DWCoreForgedGear, {
	customCondition: player => player.getLevel() === 60,
});
export const Gear2HDraconicPhase5 = PresetUtils.makePresetGear('P5 2H Draconic', Phase52HDraconicGear, {
	customCondition: player => player.getLevel() === 60,
});
export const GearDWDraconicPhase5 = PresetUtils.makePresetGear('P5 DW Draconic', Phase5DWDraconicGear, {
	customCondition: player => player.getLevel() === 60,
});

export const GearPresets = {
	[Phase.Phase1]: [GearArmsPhase1, GearFuryPhase1, GearArmsDWPhase1],
	[Phase.Phase2]: [GearArmsPhase2, GearFuryPhase2],
	[Phase.Phase3]: [GearArmsPhase3, GearFuryPhase3, GearGladPhase3],
	[Phase.Phase4]: [Gear2HPhase4, GearDWPhase4, GearGladPhase4],
	[Phase.Phase5]: [Gear2HCoreForgedPhase5, GearDWCoreForgedPhase5, Gear2HDraconicPhase5, GearDWDraconicPhase5],
};

export const DefaultGear2H = GearPresets[Phase.Phase5][0];
export const DefaultGearDW = GearPresets[Phase.Phase5][1];
// export const DefaultGearGlad = GearPresets[Phase.Phase4][2];

export const DefaultGear = DefaultGear2H;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1Arms = PresetUtils.makePresetAPLRotation('P1 Arms', Phase1APLArms, { customCondition: player => player.getLevel() === 25 });

export const APLPhase2Arms = PresetUtils.makePresetAPLRotation('P2 Arms', Phase2APLArms, { customCondition: player => player.getLevel() === 40 });
export const APLPhase2Fury = PresetUtils.makePresetAPLRotation('P2 Fury', Phase2APLFury, { customCondition: player => player.getLevel() === 40 });

export const APLPhase3Arms = PresetUtils.makePresetAPLRotation('P3 Arms', Phase3APLArms, { customCondition: player => player.getLevel() === 50 });
export const APLPhase3Fury = PresetUtils.makePresetAPLRotation('P3 Fury', Phase3APLFury, { customCondition: player => player.getLevel() === 50 });
export const APLPhase3Glad = PresetUtils.makePresetAPLRotation('P3 Glad', Phase3APLGlad, { customCondition: player => player.getLevel() === 50 });

export const APLPhase4Fury = PresetUtils.makePresetAPLRotation('P4 Fury', Phase4APLFury, { customCondition: player => player.getLevel() === 60 });
export const APLPhase4Glad = PresetUtils.makePresetAPLRotation('P4 Glad', Phase4APLGlad, { customCondition: player => player.getLevel() === 60 });
// No arms rotation right now
export const APLPhase4Arms = APLPhase4Fury;

export const APLPhase52H = PresetUtils.makePresetAPLRotation('Phase 5 2H', Phase5APL2H, {
	customCondition: player => player.getLevel() === 60,
});
export const APLPhase5DW = PresetUtils.makePresetAPLRotation('Phase 5 DW', Phase5APLDW, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [APLPhase1Arms],
	[Phase.Phase2]: [APLPhase2Arms, APLPhase2Fury],
	[Phase.Phase3]: [APLPhase3Arms, APLPhase3Fury, APLPhase3Glad],
	[Phase.Phase4]: [APLPhase4Arms, APLPhase4Fury, APLPhase4Glad],
	[Phase.Phase5]: [APLPhase52H, APLPhase5DW],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
		2: APLPresets[Phase.Phase1][0],
	},
	40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][1],
		2: APLPresets[Phase.Phase2][0],
	},
	50: {
		0: APLPresets[Phase.Phase3][0],
		1: APLPresets[Phase.Phase3][1],
		2: APLPresets[Phase.Phase3][0],
	},
	60: {
		0: APLPresets[Phase.Phase5][0], // 2H
		1: APLPresets[Phase.Phase5][1], // DW
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsPhase1 = PresetUtils.makePresetTalents('Level 25', SavedTalents.create({ talentsString: '303220203-01' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsPhase2Arms = PresetUtils.makePresetTalents('40 Arms', SavedTalents.create({ talentsString: '303050213525100001' }), {
	customCondition: player => player.getLevel() === 40,
});
export const TalentsPhase2Fury = PresetUtils.makePresetTalents('40 Fury', SavedTalents.create({ talentsString: '-05050005405010051' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsPhase3Arms = PresetUtils.makePresetTalents('50 Arms', SavedTalents.create({ talentsString: '303050213520105001-0505' }), {
	customCondition: player => player.getLevel() === 50,
});
export const TalentsPhase3Fury = PresetUtils.makePresetTalents('50 Fury', SavedTalents.create({ talentsString: '303040003-0505000540501003' }), {
	customCondition: player => player.getLevel() === 50,
});
// Glad talents are identical to fury at the moment
export const TalentsPhase3Glad = TalentsPhase3Fury;

export const TalentsPhase4Fury = PresetUtils.makePresetTalents('60 Fury DW', SavedTalents.create({ talentsString: '20305020302-05050005525010051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsPhase4Glad = PresetUtils.makePresetTalents('60 Glad', SavedTalents.create({ talentsString: '30305020302-05050005025012251' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentsPhase52H = PresetUtils.makePresetTalents('60 Fury 2H', SavedTalents.create({ talentsString: '20305020332-55020005025010051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2Arms, TalentsPhase2Fury],
	[Phase.Phase3]: [TalentsPhase3Arms, TalentsPhase3Fury, TalentsPhase3Glad],
	[Phase.Phase4]: [TalentsPhase4Fury, TalentsPhase4Glad],
	[Phase.Phase5]: [TalentsPhase52H, TalentsPhase4Fury],
};

export const DefaultTalents2H = TalentPresets[Phase.Phase5][0];
export const DefaultTalentsDW = TalentPresets[Phase.Phase5][1];
// export const DefaultTalentsGlad = TalentPresets[Phase.Phase4][3];

export const DefaultTalents = DefaultTalents2H;

export const PresetBuild2H = PresetUtils.makePresetBuild('Two-Handed', DefaultGear2H, DefaultTalents2H, DefaultAPLs[60][0]);
export const PresetBuildDW = PresetUtils.makePresetBuild('Dual-Wield', DefaultGearDW, DefaultTalentsDW, DefaultAPLs[60][1]);
// export const PresetBuildGlad = PresetUtils.makePresetBuild('Glad', DefaultGearGlad, DefaultTalentsGlad, DefaultAPLs[60][3]);

///////////////////////////////////////////////////////////////////////////
//                                 Options Presets
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	shout: WarriorShout.WarriorShoutBattle,
	stance: WarriorStance.WarriorStanceBerserker,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	armorElixir: ArmorElixir.ElixirOfSuperiorDefense,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultPotion: Potions.MightyRagePotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	food: Food.FoodSmokedDesertDumpling,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.WildStrikes,
	mildlyIrradiatedRejuvPot: true,
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	sapper: true,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	commandingShout: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	hornOfLordaeron: true,
	leaderOfThePack: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	stoneskinTotem: TristateEffect.TristateEffectRegular,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	mightOfStormwind: true,
	moldarsMoxie: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	songflowerSerenade: true,
	spiritOfZandalar: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	giftOfArthas: true,
	homunculi: 70, // 70% average uptime default
	improvedFaerieFire: true,
	improvedScorch: true,
	mangle: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Alchemy,
	profession2: Profession.Engineering,
	race: Race.RaceHuman,
};
