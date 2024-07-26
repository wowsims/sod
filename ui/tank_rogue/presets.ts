import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import { AgilityElixir, Consumes, Debuffs, IndividualBuffs, Profession, RaidBuffs, StrengthBuff, TristateEffect, WeaponImbue } from '../core/proto/common.js';
import { RogueOptions } from '../core/proto/rogue.js';
import { SavedTalents } from '../core/proto/ui.js';
import MutilateApl from './apls/mutilate.apl.json';
import P3MutilateApl from './apls/Mutilate_DPS_50.apl.json';
import MutilateIEAApl from './apls/mutilate_IEA.apl.json';
import P3ExposeMutilateApl from './apls/Mutilate_IEA_50.apl.json';
import P3SaberApl from './apls/Saber_DPS_50.apl.json';
import P3SaberIEAApl from './apls/Saber_IEA_50.apl.json';
import BlankGear from './gear_sets/blank.gear.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P1DaggersGear from './gear_sets/p1_daggers.gear.json';
import P2DaggersGear from './gear_sets/p2_daggers.gear.json';
import P3DaggersGear from './gear_sets/p3_daggers.gear.json';
import P3SaberGear from './gear_sets/p3_saber.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

// TODO: Add gear sets
export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearDaggersP1 = PresetUtils.makePresetGear('P1 Daggers', P1DaggersGear, { customCondition: player => player.getLevel() == 25 });
export const GearDaggersP2 = PresetUtils.makePresetGear('P2 Daggers', P2DaggersGear, { customCondition: player => player.getLevel() == 40 });
export const GearCombatP1 = PresetUtils.makePresetGear('P1 Combat', P1CombatGear, { customCondition: player => player.getLevel() == 25 });
export const GearDaggersP3 = PresetUtils.makePresetGear('P3 Daggers', P3DaggersGear, { customCondition: player => player.getLevel() >= 50 });
export const GearSaberP3 = PresetUtils.makePresetGear('P3 Saber', P3SaberGear, { customCondition: player => player.getLevel() >= 50 });

export const GearPresets = {
	[Phase.Phase1]: [GearDaggersP1, GearCombatP1],
	[Phase.Phase2]: [GearDaggersP2],
	[Phase.Phase3]: [GearDaggersP3, GearSaberP3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearPresets[Phase.Phase3][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { customCondition: player => player.getLevel() <= 40 });
export const ROTATION_PRESET_MUTILATE_IEA = PresetUtils.makePresetAPLRotation('Mutilate IEA', MutilateIEAApl, {
	customCondition: player => player.getLevel() <= 40,
});
export const ROTATION_PRESET_MUTILATE_P3 = PresetUtils.makePresetAPLRotation('P3 Mutilate', P3MutilateApl, {
	customCondition: player => player.getLevel() >= 50,
});
export const ROTATION_PRESET_MUTILATE_IEA_P3 = PresetUtils.makePresetAPLRotation('P3 Expose Mutilate', P3ExposeMutilateApl, {
	customCondition: player => player.getLevel() >= 50,
});
export const ROTATION_PRESET_SABER_P3 = PresetUtils.makePresetAPLRotation('P3 Saber', P3SaberApl, { customCondition: player => player.getLevel() >= 50 });
export const ROTATION_PRESET_SABER_IEA_P3 = PresetUtils.makePresetAPLRotation('P3 Expose Saber', P3SaberIEAApl, {
	customCondition: player => player.getLevel() >= 50,
});

export const APLPresets = {
	[Phase.Phase1]: [ROTATION_PRESET_MUTILATE],
	[Phase.Phase2]: [ROTATION_PRESET_MUTILATE, ROTATION_PRESET_MUTILATE_IEA],
	[Phase.Phase3]: [ROTATION_PRESET_MUTILATE_P3, ROTATION_PRESET_MUTILATE_IEA_P3, ROTATION_PRESET_SABER_P3, ROTATION_PRESET_SABER_IEA_P3],
	[Phase.Phase4]: [],
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
		1: APLPresets[Phase.Phase2][0],
		2: APLPresets[Phase.Phase2][0],
	},
	50: {
		0: APLPresets[Phase.Phase3][0],
		1: APLPresets[Phase.Phase3][0],
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

export const TankMutilate50Talents = PresetUtils.makePresetTalents(
	'P3 HAT/CttC Mutilate',
	SavedTalents.create({ talentsString: '00532012-00532500004501001-02' }),
	{
		customCondition: player => player.getLevel() >= 50,
	},
);

export const TankSaber50Talents = PresetUtils.makePresetTalents('P3 Saber Carnage', SavedTalents.create({ talentsString: '0053221-02505501000501031' }), {
	customCondition: player => player.getLevel() >= 50,
});

export const TankBladeFlurry50Talents = PresetUtils.makePresetTalents('P3 BF Poison', SavedTalents.create({ talentsString: '0053221205-02330520000501' }), {
	customCondition: player => player.getLevel() >= 50,
});

export const TalentPresets = {
	[Phase.Phase1]: [CombatDagger25Talents],
	[Phase.Phase2]: [ColdBloodMutilate40Talents, IEAMutilate40Talents, CombatMutilate40Talents],
	[Phase.Phase3]: [TankMutilate50Talents, TankSaber50Talents, TankBladeFlurry50Talents],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultTalentsAssassin = TalentPresets[Phase.Phase3][0];
export const DefaultTalentsCombat = TalentPresets[Phase.Phase3][0];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase3][0];

export const DefaultTalents = DefaultTalentsCombat;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RogueOptions.create({
	honorAmongThievesCritRate: 100,
});

///////////////////////////////////////////////////////////////////////////
//                         Consumes/Buffs/Debuffs
///////////////////////////////////////////////////////////////////////////

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfAgility,
	dragonBreathChili: false,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.DeadlyPoison,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectRegular,
	fireResistanceAura: true,
	fireResistanceTotem: true,
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
