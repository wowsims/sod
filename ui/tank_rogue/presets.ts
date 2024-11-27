import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import { Player } from '../core/proto/api';
import {
	AgilityElixir,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Flask,
	Food,
	IndividualBuffs,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { RogueOptions } from '../core/proto/rogue.js';
import { SavedTalents } from '../core/proto/ui.js';
import { playerPresets } from '../raid/presets';
import MutilateApl from './apls/mutilate.apl.json';
import P3MutilateApl from './apls/Mutilate_DPS_50.apl.json';
import MutilateIEAApl from './apls/mutilate_IEA.apl.json';
import P3ExposeMutilateApl from './apls/Mutilate_IEA_50.apl.json';
import P5SaberAPL from './apls/P5_Saber.apl.json';
import P3SaberApl from './apls/Saber_DPS_50.apl.json';
import P3SaberIEAApl from './apls/Saber_IEA_50.apl.json';
import P4SaberWeaveApl from './apls/Saber_Weave_60.apl.json';
import BlankGear from './gear_sets/blank.gear.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P1DaggersGear from './gear_sets/p1_daggers.gear.json';
import P2DaggersGear from './gear_sets/p2_daggers.gear.json';
import P3DaggersGear from './gear_sets/p3_daggers.gear.json';
import P3SaberGear from './gear_sets/p3_saber.gear.json';
import P4SaberGear from './gear_sets/p4_saber.gear.json';
import P5SaberGear from './gear_sets/p5_saber.gear.json';

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
export const GearDaggersP3 = PresetUtils.makePresetGear('P3 Daggers', P3DaggersGear, { customCondition: player => player.getLevel() == 50 });
export const GearSaberP3 = PresetUtils.makePresetGear('P3 Saber', P3SaberGear, { customCondition: player => player.getLevel() == 50 });
export const GearSaberP4 = PresetUtils.makePresetGear('P4 Saber', P4SaberGear, { customCondition: player => player.getLevel() == 60 });
export const GearSaberP5 = PresetUtils.makePresetGear('P5 Saber', P5SaberGear, { customCondition: player => player.getLevel() == 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearDaggersP1, GearCombatP1],
	[Phase.Phase2]: [GearDaggersP2],
	[Phase.Phase3]: [GearDaggersP3, GearSaberP3],
	[Phase.Phase4]: [GearSaberP4],
	[Phase.Phase5]: [GearSaberP5],
};

export const DefaultGear = GearPresets[Phase.Phase5][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { customCondition: player => player.getLevel() <= 40 });
export const ROTATION_PRESET_MUTILATE_IEA = PresetUtils.makePresetAPLRotation('Mutilate IEA', MutilateIEAApl, {
	customCondition: player => player.getLevel() <= 40,
});
export const ROTATION_PRESET_MUTILATE_P3 = PresetUtils.makePresetAPLRotation('P3 Mutilate', P3MutilateApl, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_MUTILATE_IEA_P3 = PresetUtils.makePresetAPLRotation('P3 Expose Mutilate', P3ExposeMutilateApl, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_SABER_P3 = PresetUtils.makePresetAPLRotation('P3 Saber', P3SaberApl, { customCondition: player => player.getLevel() === 50 });
export const ROTATION_PRESET_SABER_IEA_P3 = PresetUtils.makePresetAPLRotation('P3 Expose Saber', P3SaberIEAApl, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_SABER_WEAVE_P4 = PresetUtils.makePresetAPLRotation('P4 Saber Weave', P4SaberWeaveApl, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SABER_P5 = PresetUtils.makePresetAPLRotation('P5 Saber', P5SaberAPL, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [ROTATION_PRESET_MUTILATE],
	[Phase.Phase2]: [ROTATION_PRESET_MUTILATE, ROTATION_PRESET_MUTILATE_IEA],
	[Phase.Phase3]: [ROTATION_PRESET_MUTILATE_P3, ROTATION_PRESET_MUTILATE_IEA_P3, ROTATION_PRESET_SABER_P3, ROTATION_PRESET_SABER_IEA_P3],
	[Phase.Phase4]: [ROTATION_PRESET_SABER_WEAVE_P4],
	[Phase.Phase5]: [ROTATION_PRESET_SABER_P5],
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
	60: {
		0: APLPresets[Phase.Phase5][0],
		1: APLPresets[Phase.Phase5][0],
		2: APLPresets[Phase.Phase5][0],
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const CombatDagger25Talents = PresetUtils.makePresetTalents('P1 Combat Dagger', SavedTalents.create({ talentsString: '-023305002001' }), {
	customCondition: player => player.getLevel() === 25,
});

export const ColdBloodMutilate40Talents = PresetUtils.makePresetTalents('P2 CB Mutilate', SavedTalents.create({ talentsString: '005303103551--05' }), {
	customCondition: player => player.getLevel() === 40,
});

export const IEAMutilate40Talents = PresetUtils.makePresetTalents('P2 CB/IEA Mutilate', SavedTalents.create({ talentsString: '005303121551--05' }), {
	customCondition: player => player.getLevel() === 40,
});

export const CombatMutilate40Talents = PresetUtils.makePresetTalents('P2 AR/BF Mutilate', SavedTalents.create({ talentsString: '-0053052020550100201' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TankMutilate50Talents = PresetUtils.makePresetTalents(
	'P3 HAT/CttC Mutilate',
	SavedTalents.create({ talentsString: '00532012-00532500004501001-02' }),
	{
		customCondition: player => player.getLevel() === 50,
	},
);

export const TankSaber50Talents = PresetUtils.makePresetTalents('P3 Saber Carnage', SavedTalents.create({ talentsString: '0053221-02505501000501031' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TankBladeFlurry50Talents = PresetUtils.makePresetTalents('P3 BF Poison', SavedTalents.create({ talentsString: '0053221205-02330520000501' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TankSaber60Talents = PresetUtils.makePresetTalents('P4 Saber', SavedTalents.create({ talentsString: '305323102-0230550100050150131' }), {
	customCondition: player => player.getLevel() === 60,
});

export const P5TankSaberTalents = PresetUtils.makePresetTalents('P5 Saber', SavedTalents.create({ talentsString: '30532312-0230550100050140231' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [CombatDagger25Talents],
	[Phase.Phase2]: [ColdBloodMutilate40Talents, IEAMutilate40Talents, CombatMutilate40Talents],
	[Phase.Phase3]: [TankMutilate50Talents, TankSaber50Talents, TankBladeFlurry50Talents],
	[Phase.Phase4]: [TankSaber60Talents],
	[Phase.Phase5]: [P5TankSaberTalents],
};

export const DefaultTalentsAssassin = TalentPresets[Phase.Phase5][0];
export const DefaultTalentsCombat = TalentPresets[Phase.Phase5][0];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase5][0];

export const DefaultTalents = DefaultTalentsCombat;

///////////////////////////////////////////////////////////////////////////
//                                Encounters
///////////////////////////////////////////////////////////////////////////
export const PresetBuildEncounterDefault = PresetUtils.makePresetBuild('Default', {
	encounter: PresetUtils.makePresetEncounter(
		'Default',
		'http://localhost:5173/sod/tank_rogue/?i=ce#eJyTklFgM5LgYBRg1GC0YHRgrGBvYJSYwMi8gJH5BiOr0gFmTgYwiHMQhDD0HCRnzQSBk/aWEJEL9oppYHDN3mgCM8eNJl4hDp/UstQcBTMDCXutB0wMgwlodDpQxRy1tRSZIyD1nzruAAJDMFmw3MEyM/9D68mQq/aOUBmHCEYAjL4g7A==',
	),
});

export const PresetBuildEncounterVael = PresetUtils.makePresetBuild('Vael', {
	encounter: PresetUtils.makePresetEncounter(
		'Vael',
		'http://localhost:5173/sod/tank_rogue/?i=ce#eJyTkldgM5LmYBRg1GC0YHRgrGBvYGSZwMi8gJH5EiPDDUZWpe8snAxgEOcgCGHoOUjOmgkCJ+0tISIX7BXTwOCavdFxFo47qUJqTjmJydnlmXnpCj6JmUUKYYmpOYnFJUWJxVUKJRmpCs75RUWlBSUS9grMWg+YGAYT0Oh0oIo5amspMiejyd+RKu4AAkMw2bHZwVI2KsX6vv9Ve5jRDhGMRVMZORiFFEIyc1MVnEqL8kBx5phSlJqXmJOZl6oQlJqcmlmWmqKV75FfrpCTD5TNzCvJB8diWmZ6Rgk2TZnFCsnA6FbIzwOrK8hJrEwt0lNwyywCCoJlgCqMDIoVNIJTc1KTSxQMFIBG5qWWpRYpFEEs1AQAp5dhnw==',
	),
});

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RogueOptions.create({
	honorAmongThievesCritRate: 100,
	pkSwap: false
});

///////////////////////////////////////////////////////////////////////////
//                         Consumes/Buffs/Debuffs
///////////////////////////////////////////////////////////////////////////

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	flask: Flask.FlaskOfTheTitans,
	food: Food.FoodGrilledSquid,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.DeadlyPoison,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	demonicPact: 80,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	mightOfStormwind: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
	spiritOfZandalar: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	faerieFire: true,
	homunculi: 100,
	improvedFaerieFire: true,
	improvedScorch: true,
	mangle: true,
	markOfChaos: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
};
