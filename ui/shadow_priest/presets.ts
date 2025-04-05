import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import { Spec } from '../core/proto/common';
import {
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Flask,
	Food,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	ShadowPowerBuff,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { ShadowPriest_Options as Options } from '../core/proto/priest.js';
import { SavedTalents } from '../core/proto/ui.js';
// APLs
import Phase1APL from './apls/phase_1.apl.json';
import Phase2APL from './apls/phase_2.apl.json';
import Phase3APL from './apls/phase_3.apl.json';
import Phase4APL from './apls/phase_4.apl.json';
import Phase5APL from './apls/phase_5.apl.json';
import Phase6APL from './apls/phase_6.apl.json';
// Builds
import Phase7BuildJSON from './builds/phase_7.build.json';
import Phase8BuildJSON from './builds/phase_8.build.json';
// Gear
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';
import Phase3Gear from './gear_sets/phase_3.gear.json';
import Phase4Gear from './gear_sets/phase_4.gear.json';
import Phase5CoreForgedGear from './gear_sets/phase_5_t1.gear.json';
import Phase5DraconicGear from './gear_sets/phase_5_t2.gear.json';
import Phase6Gear from './gear_sets/phase_6.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PresetBuildPhase7 = PresetUtils.makePresetBuildFromJSON('Phase 7', Spec.SpecShadowPriest, Phase7BuildJSON);
export const PresetBuildPhase8 = PresetUtils.makePresetBuildFromJSON('Phase 8', Spec.SpecShadowPriest, Phase8BuildJSON);

export const DefaultBuild = PresetBuildPhase8;

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear, { customCondition: player => player.getLevel() === 25 });
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear, { customCondition: player => player.getLevel() === 40 });
export const GearPhase3 = PresetUtils.makePresetGear('Phase 3', Phase3Gear, { customCondition: player => player.getLevel() === 50 });
export const GearPhase4 = PresetUtils.makePresetGear('Phase 4', Phase4Gear, { customCondition: player => player.getLevel() === 60 });
export const GearPhase5Draconic = PresetUtils.makePresetGear('P5 Draconic', Phase5DraconicGear, { customCondition: player => player.getLevel() === 60 });
export const GearPhase5CoreForged = PresetUtils.makePresetGear('P5 Core Forged', Phase5CoreForgedGear, {
	customCondition: player => player.getLevel() === 60,
});
export const GearPhase6 = PresetUtils.makePresetGear('Phase 6', Phase6Gear, { customCondition: player => player.getLevel() === 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearPhase1],
	[Phase.Phase2]: [GearPhase2],
	[Phase.Phase3]: [GearPhase3],
	[Phase.Phase4]: [GearPhase4],
	[Phase.Phase5]: [GearPhase5Draconic, GearPhase5CoreForged],
	[Phase.Phase6]: [GearPhase6],
	[Phase.Phase7]: [PresetBuildPhase7.gear!],
	[Phase.Phase8]: [PresetBuildPhase8.gear!],
};

export const DefaultGear = GearPresets[Phase.Phase8][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL, { customCondition: player => player.getLevel() === 25 });
export const APLPhase2 = PresetUtils.makePresetAPLRotation('Phase 2', Phase2APL, { customCondition: player => player.getLevel() === 40 });
export const APLPhase3 = PresetUtils.makePresetAPLRotation('Phase 3', Phase3APL, { customCondition: player => player.getLevel() === 50 });
export const APLPhase4 = PresetUtils.makePresetAPLRotation('Phase 4', Phase4APL, { customCondition: player => player.getLevel() === 60 });
export const APLPhase5 = PresetUtils.makePresetAPLRotation('Phase 5', Phase5APL, { customCondition: player => player.getLevel() === 60 });
export const APLPhase6 = PresetUtils.makePresetAPLRotation('Phase 6', Phase6APL, { customCondition: player => player.getLevel() === 60 });

export const APLPresets = {
	[Phase.Phase1]: [APLPhase1],
	[Phase.Phase2]: [APLPhase2],
	[Phase.Phase3]: [APLPhase3],
	[Phase.Phase4]: [APLPhase4],
	[Phase.Phase5]: [APLPhase5],
	[Phase.Phase6]: [APLPhase6],
	[Phase.Phase7]: [PresetBuildPhase7.rotation!],
	[Phase.Phase8]: [PresetBuildPhase8.rotation!],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase8][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsPhase1 = PresetUtils.makePresetTalents('Level 25', SavedTalents.create({ talentsString: '-20535000001' }), {
	customCondition: player => player.getLevel() === 25,
});
export const TalentsPhase2 = PresetUtils.makePresetTalents('Level 40', SavedTalents.create({ talentsString: '--5022204002501251' }), {
	customCondition: player => player.getLevel() === 40,
});
export const TalentsPhase3 = PresetUtils.makePresetTalents('Level 50', SavedTalents.create({ talentsString: '-0055-5022204002501251' }), {
	customCondition: player => player.getLevel() === 50,
});
export const TalentsPhase4 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '0512301302--5002504103501251' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2],
	[Phase.Phase3]: [TalentsPhase3],
	[Phase.Phase4]: [TalentsPhase4],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
	[Phase.Phase7]: [PresetBuildPhase7.talents!],
	[Phase.Phase8]: [PresetBuildPhase8.talents!],
};

export const DefaultTalents = TalentPresets[Phase.Phase8][0];

///////////////////////////////////////////////////////////////////////////
//                                 Build Presets
///////////////////////////////////////////////////////////////////////////

export const PresetBuildPhase4 = PresetUtils.makePresetBuild('Phase 4', {
	gear: GearPhase4,
	talents: TalentsPhase4,
	rotation: APLPhase4,
});
export const PresetBuildPhase5Draconic = PresetUtils.makePresetBuild('Phase 5 Draconic', {
	gear: GearPhase5Draconic,
	talents: TalentsPhase4,
	rotation: APLPhase5,
});
export const PresetBuildPhase5CoreForged = PresetUtils.makePresetBuild('Phase 5 Core Forged', {
	gear: GearPhase5CoreForged,
	talents: TalentsPhase4,
	rotation: APLPhase5,
});
export const PresetBuildPhase6 = PresetUtils.makePresetBuild('Phase 6', {
	gear: GearPhase6,
	talents: TalentsPhase4,
	rotation: APLPhase6,
});

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = Options.create({
});

export const DefaultConsumes = Consumes.create({
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodDarkclawBisque,
	mainHandImbue: WeaponImbue.EnchantedRepellent,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	mildlyIrradiatedRejuvPot: true,
	shadowPowerBuff: ShadowPowerBuff.ElixirOfShadowPower,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	zanzaBuff: ZanzaBuff.CerebralCortexCompound,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	demonicPact: 120,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
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
	improvedShadowBolt: true,
	judgementOfWisdom: true,
	occultPoison: true,
	markOfChaos: true,
	wintersChill: true,
});

export const OtherDefaults = {
	channelClipDelay: 100,
	distanceFromTarget: 30,
	profession1: Profession.Alchemy,
	profession2: Profession.Enchanting,
};
