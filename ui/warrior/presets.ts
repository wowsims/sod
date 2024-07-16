import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	AttackPowerBuff,
	Consumes,
	Debuffs,
	DragonslayerBuff,
	EnchantedSigil,
	Flask,
	Food,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Warrior_Options as WarriorOptions, WarriorShout } from '../core/proto/warrior.js';
import Phase1APLArms from './apls/phase_1_arms.apl.json';
import Phase2APLArms from './apls/phase_2_arms.apl.json';
import Phase2APLFury from './apls/phase_2_fury.apl.json';
import Phase3APLArms from './apls/phase_3_arms.apl.json';
import Phase3APLFury from './apls/phase_3_fury.apl.json';
import Phase3APLGlad from './apls/phase_3_glad.apl.json';
import Phase4APLFury from './apls/phase_4_fury.apl.json';
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

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArmsPhase1 = PresetUtils.makePresetGear('P1 Arms 2H', Phase1Gear, { talentTree: 0 });
export const GearArmsDWPhase1 = PresetUtils.makePresetGear('P1 Arms DW', Phase1DWGear, { talentTree: 0 });
export const GearFuryPhase1 = PresetUtils.makePresetGear('P1 DW Fury', Phase1Gear, { talentTree: 1 });
export const GearArmsPhase2 = PresetUtils.makePresetGear('P2 2H', Phase22HGear, { talentTree: 0 });
export const GearFuryPhase2 = PresetUtils.makePresetGear('P2 DW', Phase2DWGear, { talentTree: 1 });
export const GearArmsPhase3 = PresetUtils.makePresetGear('P3 2H', Phase32HGear, { talentTree: 0 });
export const GearFuryPhase3 = PresetUtils.makePresetGear('P3 DW', Phase3DWGear, { talentTree: 1 });
export const GearGladPhase3 = PresetUtils.makePresetGear('P3 Glad', Phase3GladGear, { talentTree: 1 });
export const Gear2HPhase4 = PresetUtils.makePresetGear('P4 2H', Phase42HGear);
export const GearDWPhase4 = PresetUtils.makePresetGear('P4 DW', Phase4DWGear);
export const GearGladPhase4 = PresetUtils.makePresetGear('P4 Glad', Phase4GladGear);

export const GearPresets = {
	[Phase.Phase1]: [GearArmsPhase1, GearFuryPhase1, GearArmsDWPhase1],
	[Phase.Phase2]: [GearArmsPhase2, GearFuryPhase2],
	[Phase.Phase3]: [GearArmsPhase3, GearFuryPhase3, GearGladPhase3],
	[Phase.Phase4]: [Gear2HPhase4, GearDWPhase4, GearGladPhase4],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearPresets[Phase.Phase4][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1Arms = PresetUtils.makePresetAPLRotation('P1 Arms', Phase1APLArms);
export const APLPhase2Arms = PresetUtils.makePresetAPLRotation('P2 Arms', Phase2APLArms);
export const APLPhase2Fury = PresetUtils.makePresetAPLRotation('P2 Fury', Phase2APLFury);
export const APLPhase3Arms = PresetUtils.makePresetAPLRotation('P3 Arms', Phase3APLArms);
export const APLPhase3Fury = PresetUtils.makePresetAPLRotation('P3 Fury', Phase3APLFury);
export const APLPhase3Glad = PresetUtils.makePresetAPLRotation('P3 Glad', Phase3APLGlad);
export const APLPhase4Fury = PresetUtils.makePresetAPLRotation('P4 Fury', Phase4APLFury);

export const APLPresets = {
	[Phase.Phase1]: [APLPhase1Arms],
	[Phase.Phase2]: [APLPhase2Arms, APLPhase2Fury],
	[Phase.Phase3]: [APLPhase3Arms, APLPhase3Fury, APLPhase3Glad],
	[Phase.Phase4]: [APLPhase4Fury],
	[Phase.Phase5]: [],
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
		1: APLPresets[Phase.Phase4][0],
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

export const TalentsPhase4Arms = PresetUtils.makePresetTalents('60 Arms', SavedTalents.create({ talentsString: '303050213520105001-05050005023' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsPhase4Fury = PresetUtils.makePresetTalents('60 Fury', SavedTalents.create({ talentsString: '20305020302-05050005525010051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsPhase4Prot = PresetUtils.makePresetTalents('60 Prot', SavedTalents.create({ talentsString: '20305020302-03-55200110530201051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2Arms, TalentsPhase2Fury],
	[Phase.Phase3]: [TalentsPhase3Arms, TalentsPhase3Fury, TalentsPhase3Glad],
	[Phase.Phase4]: [TalentsPhase4Arms, TalentsPhase4Fury, TalentsPhase4Prot],
	[Phase.Phase5]: [],
};

export const DefaultTalentsArms = TalentPresets[Phase.Phase4][0];
export const DefaultTalentsFury = TalentPresets[Phase.Phase4][1];
export const DefaultTalentsGlad = TalentPresets[Phase.Phase4][2];

export const DefaultTalents = DefaultTalentsFury;

///////////////////////////////////////////////////////////////////////////
//                                 Options Presets
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	shout: WarriorShout.WarriorShoutBattle,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultPotion: Potions.MightyRagePotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	food: Food.FoodSmokedDesertDumpling,
	flask: Flask.FlaskOfTheTitans,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	commandingShout: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	dragonslayerBuff: DragonslayerBuff.RallyingCryofTheDragonslayer,
	mightOfStormwind: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	songflowerSerenade: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	faerieFire: true,
	giftOfArthas: true,
	homunculi: 70, // 70% average uptime default
	improvedFaerieFire: true,
	mangle: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Blacksmithing,
	profession2: Profession.Enchanting,
};
