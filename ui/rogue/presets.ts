import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import { AgilityElixir, Consumes, Debuffs, IndividualBuffs, Profession, RaidBuffs, StrengthBuff, TristateEffect, WeaponImbue } from '../core/proto/common.js';
import { RogueOptions } from '../core/proto/rogue.js';
import { SavedTalents } from '../core/proto/ui.js';
import SinisterApl25 from './apls/basic_strike_25.apl.json';
import CPCarnageMutilateApl50 from './apls/CP_Carnage_Mutilate.apl.json';
import CPCarnageSaberApl50 from './apls/CP_Carnage_Saber.apl.json';
import CPCttCMutilateApl50 from './apls/CP_CttC_Mutilate.apl.json';
import CPCttCSaberApl50 from './apls/CP_CttC_Saber.apl.json';
import HATCarnageMutilateApl50 from './apls/HAT_Carnage_Mutilate.apl.json';
import HATCarnageSaberApl50 from './apls/HAT_Carnage_Saber.apl.json';
import HATCttCMutilateApl50 from './apls/HAT_CttC_Mutilate.apl.json';
import HATCttCSaberApl50 from './apls/HAT_CttC_Saber.apl.json';
import MutilateApl from './apls/mutilate.apl.json';
import MutilateIEAApl from './apls/mutilate_IEA.apl.json';
import BlankGear from './gear_sets/blank.gear.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P1Daggers from './gear_sets/p1_daggers.gear.json';
import P2DaggersGear from './gear_sets/p2_daggers.gear.json';
import P3DaggersGear from './gear_sets/p3_muti.gear.json';
import P3CombatGear from './gear_sets/p3_saber.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearDaggersP1 = PresetUtils.makePresetGear('P1 Daggers', P1Daggers, { customCondition: player => player.getLevel() == 25 });
export const GearCombatP1 = PresetUtils.makePresetGear('P1 Combat', P1CombatGear, { customCondition: player => player.getLevel() == 25 });
export const GearDaggersP2 = PresetUtils.makePresetGear('P2 Daggers', P2DaggersGear, { customCondition: player => player.getLevel() == 40 });
export const GearDaggersP3 = PresetUtils.makePresetGear('P3 Daggers', P3DaggersGear, { customCondition: player => player.getLevel() >= 50 });
export const GearCombatP3 = PresetUtils.makePresetGear('P3 Combat', P3CombatGear, { customCondition: player => player.getLevel() >= 50 });

export const GearPresets = {
	[Phase.Phase1]: [GearDaggersP1, GearCombatP1],
	[Phase.Phase2]: [GearDaggersP2],
	[Phase.Phase3]: [GearDaggersP3, GearCombatP3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearPresets[Phase.Phase3][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets[]
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { customCondition: player => player.getLevel() == 40 });
export const ROTATION_PRESET_MUTILATE_IEA = PresetUtils.makePresetAPLRotation('Mutilate IEA', MutilateIEAApl, {
	customCondition: player => player.getLevel() == 40,
});
export const ROTATION_PRESET_SINISTER_25 = PresetUtils.makePresetAPLRotation('Sinister', SinisterApl25, { customCondition: player => player.getLevel() == 25 });
export const ROTATION_PRESET_CP_CARNAGE_MUTI_50 = PresetUtils.makePresetAPLRotation('CP/FA Carnage Muti', CPCarnageMutilateApl50, { customCondition: player => player.getLevel() >= 50 })
export const ROTATION_PRESET_CP_CTTC_MUTI_50 = PresetUtils.makePresetAPLRotation('CP/FA CttC Muti', CPCttCMutilateApl50, { customCondition: player => player.getLevel() >= 50 })
export const ROTATION_PRESET_HAT_CARNAGE_MUTI_50 = PresetUtils.makePresetAPLRotation('HAT Carnage Muti', HATCarnageMutilateApl50, { customCondition: player => player.getLevel() >= 50 })
export const ROTATION_PRESET_HAT_CTTC_MUTI_50 = PresetUtils.makePresetAPLRotation('HAT CttC Muti', HATCttCMutilateApl50, { customCondition: player => player.getLevel() >= 50 })
export const ROTATION_PRESET_CP_CARNAGE_SABER_50 = PresetUtils.makePresetAPLRotation('CP/FA Carnage Saber', CPCarnageSaberApl50, { customCondition: player => player.getLevel() >= 50 })
export const ROTATION_PRESET_CP_CTTC_SABER_50 = PresetUtils.makePresetAPLRotation('CP/FA CttC Saber', CPCttCSaberApl50, { customCondition: player => player.getLevel() >= 50 })
export const ROTATION_PRESET_HAT_CARNAGE_SABER_50 = PresetUtils.makePresetAPLRotation('HAT Carnage Saber', HATCarnageSaberApl50, { customCondition: player => player.getLevel() >= 50 })
export const ROTATION_PRESET_HAT_CTTC_SABER_50 = PresetUtils.makePresetAPLRotation('HAT CttC Saber', HATCttCSaberApl50, { customCondition: player => player.getLevel() >= 50 })

export const APLPresets = {
	[Phase.Phase1]: [ROTATION_PRESET_MUTILATE, ROTATION_PRESET_SINISTER_25],
	[Phase.Phase2]: [ROTATION_PRESET_MUTILATE, ROTATION_PRESET_MUTILATE_IEA],
	[Phase.Phase3]: [ROTATION_PRESET_CP_CARNAGE_MUTI_50, ROTATION_PRESET_CP_CARNAGE_SABER_50, 
		ROTATION_PRESET_CP_CTTC_MUTI_50, ROTATION_PRESET_CP_CTTC_SABER_50, 
		ROTATION_PRESET_HAT_CARNAGE_MUTI_50, ROTATION_PRESET_HAT_CARNAGE_SABER_50, 
		ROTATION_PRESET_HAT_CTTC_MUTI_50, ROTATION_PRESET_HAT_CTTC_SABER_50],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
	},
	40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][0],
		2: APLPresets[Phase.Phase2][0],
	},
	50: {
		0: APLPresets[Phase.Phase3][0],
		1: APLPresets[Phase.Phase3][1],
		2: APLPresets[Phase.Phase3][0],
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const CombatDagger25Talents = PresetUtils.makePresetTalents('P1 Combat Dagger', SavedTalents.create({ talentsString: '-023305002001' }), {
	customCondition: player => player.getLevel() == 25,
});

export const ColdBloodMutilate40Talents = PresetUtils.makePresetTalents('P2 CB Mutilate', SavedTalents.create({ talentsString: '005303103551--05' }), {
	customCondition: player => player.getLevel() == 40,
});

export const IEAMutilate40Talents = PresetUtils.makePresetTalents('P2 CB/IEA Mutilate', SavedTalents.create({ talentsString: '005303121551--05' }), {
	customCondition: player => player.getLevel() == 40,
});

export const CombatMutilate40Talents = PresetUtils.makePresetTalents('P2 AR/BF Mutilate', SavedTalents.create({ talentsString: '-0053052020550100201' }), {
	customCondition: player => player.getLevel() == 40,
});

export const ColdBloodDualWieldSpec50Talents = PresetUtils.makePresetTalents('P3 CB/DWS', SavedTalents.create({ talentsString: '005323101501-320015202005' }), {
	customCondition: player => player.getLevel() >= 50,
});

export const ColdBloodSealFate50Talents = PresetUtils.makePresetTalents('P3 Cold Blood', SavedTalents.create({ talentsString: '00532310155105-320005' }), {
	customCondition: player => player.getLevel() >= 50,
});

export const TalentPresets = {
	[Phase.Phase1]: [CombatDagger25Talents],
	[Phase.Phase2]: [ColdBloodMutilate40Talents, IEAMutilate40Talents, CombatMutilate40Talents],
	[Phase.Phase3]: [ColdBloodSealFate50Talents, ColdBloodDualWieldSpec50Talents],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultTalentsAssassin = TalentPresets[Phase.Phase3][0];
export const DefaultTalentsCombat = TalentPresets[Phase.Phase3][1];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase3][0];

export const DefaultTalents = DefaultTalentsAssassin;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RogueOptions.create({
	honorAmongThievesCritRate: 100,
});

///////////////////////////////////////////////////////////////////////////
//                         Consumes/Buffs/Debuffs
///////////////////////////////////////////////////////////////////////////

export const P1Consumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfLesserAgility,
	dragonBreathChili: false,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.BlackfathomSharpeningStone,
});

export const P2Consumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfAgility,
	dragonBreathChili: false,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.ShadowOil,
});

export const P3Consumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfGreaterAgility,
	dragonBreathChili: false,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.ShadowOil,
})

export const DefaultConsumes = {
	[Phase.Phase1]: P1Consumes,
	[Phase.Phase2]: P2Consumes,
	[Phase.Phase3]: P3Consumes,
};

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectRegular,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectRegular,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	faerieFire: true,
	sunderArmor: true,
	mangle: true,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
};
