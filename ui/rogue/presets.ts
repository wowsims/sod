import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
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
import { RogueOptions, RogueRune } from '../core/proto/rogue.js';
import { SavedTalents } from '../core/proto/ui.js';
import SinisterApl25 from './apls/basic_strike_25.apl.json';
import MutilateApl40 from './apls/mutilate.apl.json';
import MutilateDPSAPL60 from './apls/Mutilate_60.apl.json';
import MutilateDPSApl50 from './apls/Mutilate_DPS_50.apl.json';
import MutilateIEAApl40 from './apls/mutilate_IEA.apl.json';
import MutilateIEAApl50 from './apls/Mutilate_IEA_50.apl.json';
import P5AssassinationBackstabAPL from './apls/P5_Assassination_Backstab.apl.json';
import P5CombatBackstabAPL from './apls/P5_Combat_Backstab.apl.json';
import P5MutilateAPL from './apls/P5_Mutilate.apl.json';
import P5MutilateIEAAPL from './apls/P5_Mutilate_IEA.apl.json';
import P5SaberAPL from './apls/P5_Saber.apl.json';
import P5SaberIEAAPL from './apls/P5_Saber_IEA.apl.json';
import SaberDPSApl50 from './apls/Saber_DPS_50.apl.json';
import SaberDPSAPL60 from './apls/Saber_DPS_60.apl.json';
import SaberIEAApl50 from './apls/Saber_IEA_50.apl.json';
import SlaughterCutthroatDPSAPL60 from './apls/Slaughter_Cutthroat_60.apl.json';
import BlankGear from './gear_sets/blank.gear.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P1Daggers from './gear_sets/p1_daggers.gear.json';
import P2DaggersGear from './gear_sets/p2_daggers.gear.json';
import P3MutiGear from './gear_sets/p3_muti.gear.json';
import P3MutiHatGear from './gear_sets/p3_muti_hat.gear.json';
import P3SaberGear from './gear_sets/p3_saber.gear.json';
import P4MutiGear from './gear_sets/p4_muti.gear.json';
import P4SaberGear from './gear_sets/p4_saber.gear.json';
import P5BackstabGear from './gear_sets/p5_backstab.gear.json';
import P5MutilateGear from './gear_sets/p5_mutilate.gear.json';
import P5SaberGear from './gear_sets/p5_saber.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const P1GearDaggers = PresetUtils.makePresetGear('P1 Daggers', P1Daggers, { customCondition: player => player.getLevel() === 25 });
export const P1GearSaber = PresetUtils.makePresetGear('P1 Saber', P1CombatGear, { customCondition: player => player.getLevel() === 25 });
export const P2GearDaggers = PresetUtils.makePresetGear('P2 Daggers', P2DaggersGear, { customCondition: player => player.getLevel() === 40 });
export const P3GearMuti = PresetUtils.makePresetGear('P3 Mutilate', P3MutiGear, { customCondition: player => player.getLevel() === 50 });
export const P3GearMutiHat = PresetUtils.makePresetGear('P3 Mutilate (HaT)', P3MutiHatGear, { customCondition: player => player.getLevel() === 50 });
export const P3GearSaber = PresetUtils.makePresetGear('P3 Saber', P3SaberGear, { customCondition: player => player.getLevel() === 50 });
export const P4GearMuti = PresetUtils.makePresetGear('P4 Mutilate', P4MutiGear, { customCondition: player => player.getLevel() === 60 });
export const P4GearSaber = PresetUtils.makePresetGear('P4 Saber', P4SaberGear, { customCondition: player => player.getLevel() === 60 });
export const P5GearBackstab = PresetUtils.makePresetGear('P5 Backstab', P5BackstabGear, { customCondition: player => player.getLevel() === 60 });
export const P5GearMutilate = PresetUtils.makePresetGear('P5 Mutilate', P5MutilateGear, { customCondition: player => player.getLevel() === 60 });
export const P5GearSaber = PresetUtils.makePresetGear('P5 Saber', P5SaberGear, { customCondition: player => player.getLevel() === 60 });

export const GearPresets = {
	[Phase.Phase1]: [P1GearDaggers, P1GearSaber],
	[Phase.Phase2]: [P2GearDaggers],
	[Phase.Phase3]: [P3GearMuti, P3GearMutiHat, P3GearSaber],
	[Phase.Phase4]: [P4GearMuti, P4GearSaber],
	[Phase.Phase5]: [P5GearBackstab, P5GearMutilate, P5GearSaber],
};

export const DefaultGear = GearPresets[Phase.Phase5][0];
export const DefaultGearBackstab = GearPresets[Phase.Phase5][0];
export const DefaultGearMutilate = GearPresets[Phase.Phase5][1];
export const DefaultGearSaber = GearPresets[Phase.Phase5][2];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets[]
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_SINISTER_25 = PresetUtils.makePresetAPLRotation('P1 Sinister', SinisterApl25, {
	customCondition: player => player.getLevel() === 25,
});
export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('P2 Mutilate', MutilateApl40, {
	customCondition: player => player.getLevel() === 40,
});
export const ROTATION_PRESET_MUTILATE_IEA = PresetUtils.makePresetAPLRotation('P2 Mutilate IEA', MutilateIEAApl40, {
	customCondition: player => player.getLevel() === 40,
});
export const ROTATION_PRESET_MUTILATE_DPS_50 = PresetUtils.makePresetAPLRotation('P3 Mutilate DPS', MutilateDPSApl50, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_MUTILATE_IEA_50 = PresetUtils.makePresetAPLRotation('P3 Mutilate IEA', MutilateIEAApl50, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_SABER_SLASH_DPS_50 = PresetUtils.makePresetAPLRotation('P3 Saber Slash DPS', SaberDPSApl50, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_SABER_SLASH_IEA_50 = PresetUtils.makePresetAPLRotation('P3 Saber Slash IEA', SaberIEAApl50, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_SABER_SLASH_DPS_60 = PresetUtils.makePresetAPLRotation('P4 Saber Slash', SaberDPSAPL60, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_MUTILATE_DPS_60 = PresetUtils.makePresetAPLRotation('P4 Mutilate', MutilateDPSAPL60, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SLAUGHTER_CUTTHROAT_DPS_60 = PresetUtils.makePresetAPLRotation('P4 Backstab', SlaughterCutthroatDPSAPL60, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_ASSASSINATION_BACKSTAB_DPS_P5 = PresetUtils.makePresetAPLRotation('P5 Assassination Backstab', P5AssassinationBackstabAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_COMBAT_BACKSTAB_DPS_P5 = PresetUtils.makePresetAPLRotation('P5 Combat Backstab', P5CombatBackstabAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_MUTILATE_DPS_P5 = PresetUtils.makePresetAPLRotation('P5 Mutilate', P5MutilateAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SABER_DPS_P5 = PresetUtils.makePresetAPLRotation('P5 Saber Slash', P5SaberAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_MUTILATE_IEA_P5 = PresetUtils.makePresetAPLRotation('P5 Mutilate IEA', P5MutilateIEAAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SABER_IEA_P5 = PresetUtils.makePresetAPLRotation('P5 Saber Slash IEA', P5SaberIEAAPL, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [ROTATION_PRESET_MUTILATE, ROTATION_PRESET_SINISTER_25],
	[Phase.Phase2]: [ROTATION_PRESET_MUTILATE, ROTATION_PRESET_MUTILATE_IEA],
	[Phase.Phase3]: [ROTATION_PRESET_MUTILATE_DPS_50, ROTATION_PRESET_SABER_SLASH_DPS_50, ROTATION_PRESET_MUTILATE_IEA_50, ROTATION_PRESET_SABER_SLASH_IEA_50],
	[Phase.Phase4]: [ROTATION_PRESET_MUTILATE_DPS_60, ROTATION_PRESET_SLAUGHTER_CUTTHROAT_DPS_60, ROTATION_PRESET_SABER_SLASH_DPS_60],
	[Phase.Phase5]: [
		ROTATION_PRESET_ASSASSINATION_BACKSTAB_DPS_P5,
		ROTATION_PRESET_COMBAT_BACKSTAB_DPS_P5,
		ROTATION_PRESET_MUTILATE_DPS_P5,
		ROTATION_PRESET_SABER_DPS_P5,
		ROTATION_PRESET_MUTILATE_IEA_P5,
		ROTATION_PRESET_SABER_IEA_P5,
	],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {},
	40: {
		[RogueRune.RuneMutilate]: ROTATION_PRESET_MUTILATE,
	},
	50: {
		[RogueRune.RuneMutilate]: ROTATION_PRESET_MUTILATE_DPS_50,
		[RogueRune.RuneSaberSlash]: ROTATION_PRESET_SABER_SLASH_DPS_50,
	},
	60: {
		[RogueRune.RuneMutilate]: ROTATION_PRESET_MUTILATE_DPS_P5,
		[RogueRune.RuneSaberSlash]: ROTATION_PRESET_SABER_DPS_P5,
		[RogueRune.RuneCutthroat]: ROTATION_PRESET_ASSASSINATION_BACKSTAB_DPS_P5,
	},
};

export const DefaultAPLBackstab = APLPresets[Phase.Phase5][0];
export const DefaultAPLMutilate = APLPresets[Phase.Phase5][2];
export const DefaultAPLSaber = APLPresets[Phase.Phase5][3];

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

// Preset name must be unique. Ex: 'Mutilate DPS' cannot be used as a name more than once
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

export const P3TalentsMuti = PresetUtils.makePresetTalents('P3 Mutilate', SavedTalents.create({ talentsString: '00532010555101-3203-05' }), {
	customCondition: player => player.getLevel() === 50,
});

export const P3TalentsMutiHat = PresetUtils.makePresetTalents('P3 Mutilate (HaT)', SavedTalents.create({ talentsString: '005323101551051-3203-01' }), {
	customCondition: player => player.getLevel() === 50,
});

export const P3TalentsSaber = PresetUtils.makePresetTalents('P3 Saber', SavedTalents.create({ talentsString: '005323101551051-320004' }), {
	customCondition: player => player.getLevel() === 50,
});

export const P4TalentsMutiSaber = PresetUtils.makePresetTalents('P4 Mutilate/Saber', SavedTalents.create({ talentsString: '00532310155104-02330520000501' }), {
	customCondition: player => player.getLevel() === 60,
});

export const P4TalentsSlaughter = PresetUtils.makePresetTalents('P4 Backstab', SavedTalents.create({ talentsString: '005323105521051-023305-05' }), {
	customCondition: player => player.getLevel() === 60,
});

export const P5TalentBackstabAssassination = PresetUtils.makePresetTalents(
	'P5 Backstab Assassination',
	SavedTalents.create({ talentsString: '005323105551051-023302-05' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const P5TalentBackstabCombat = PresetUtils.makePresetTalents(
	'P5 Backstab Combat',
	SavedTalents.create({ talentsString: '005023104-0233050020550100221-05' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const P5TalentMutilateSaberslashCarnage = PresetUtils.makePresetTalents(
	'P5 Mutilate/Saber Carnage',
	SavedTalents.create({ talentsString: '00532310155104-02330520000501' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const P5TalentMutilateSaberslashCTTC = PresetUtils.makePresetTalents(
	'P5 Mutilate/Saber CTTC',
	SavedTalents.create({ talentsString: '00532010455104-02330520000501' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const P5TalentBackstabAssassinationIEA = PresetUtils.makePresetTalents(
	'P5 Backstab IEA',
	SavedTalents.create({ talentsString: '005323125501051-023305-05' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const P5TalentMutilateSaberslashCTTCIEA = PresetUtils.makePresetTalents(
	'P5 Mutilate/Saber CTTC IEA',
	SavedTalents.create({ talentsString: '00532012255104-02400520200501' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const TalentPresets = {
	[Phase.Phase1]: [CombatDagger25Talents],
	[Phase.Phase2]: [ColdBloodMutilate40Talents, IEAMutilate40Talents, CombatMutilate40Talents],
	[Phase.Phase3]: [P3TalentsMuti, P3TalentsMutiHat, P3TalentsSaber],
	//	[Phase.Phase4]: [P4TalentsMutiSaber, P4TalentsSlaughter],
	[Phase.Phase5]: [
		P5TalentBackstabAssassination,
		P5TalentBackstabCombat,
		P5TalentMutilateSaberslashCarnage,
		P5TalentMutilateSaberslashCTTC,
		P5TalentBackstabAssassinationIEA,
		P5TalentMutilateSaberslashCTTCIEA,
	],
};

export const DefaultTalentsAssassin = TalentPresets[Phase.Phase5][0];
export const DefaultTalentsCombat = TalentPresets[Phase.Phase5][0];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase5][0];

export const DefaultTalentsBackstab = TalentPresets[Phase.Phase5][0];
export const DefaultTalentsMutilate = TalentPresets[Phase.Phase5][3];
export const DefaultTalentsSaber = TalentPresets[Phase.Phase5][2];

export const DefaultTalents = DefaultTalentsAssassin;

export const PresetBuildBackstab = PresetUtils.makePresetBuild('Backstab', DefaultGearBackstab, P5TalentBackstabAssassination, DefaultAPLBackstab);
export const PresetBuildMutilate = PresetUtils.makePresetBuild('Mutilate', DefaultGearMutilate, DefaultTalentsMutilate, DefaultAPLMutilate);
export const PresetBuildSaber = PresetUtils.makePresetBuild('Saber Slash', DefaultGearSaber, DefaultTalentsSaber, DefaultAPLSaber);
export const PresetBuildMutilateIEA = PresetUtils.makePresetBuild(
	'Mutilate IEA',
	DefaultGearMutilate,
	P5TalentMutilateSaberslashCTTCIEA,
	ROTATION_PRESET_MUTILATE_IEA_P5,
);
export const PresetBuildSaberIEA = PresetUtils.makePresetBuild(
	'Saber Slash IEA',
	DefaultGearSaber,
	P5TalentMutilateSaberslashCTTCIEA,
	ROTATION_PRESET_SABER_IEA_P5,
);

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
});

export const P4Consumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	flask: Flask.FlaskOfSupremePower,
	food: Food.FoodGrilledSquid,
	mainHandImbue: WeaponImbue.WildStrikes,
	miscConsumes: {
		jujuEmber: true,
	},
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const P5Consumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	flask: Flask.FlaskOfSupremePower,
	food: Food.FoodGrilledSquid,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const DefaultConsumes = {
	[Phase.Phase1]: P1Consumes,
	[Phase.Phase2]: P2Consumes,
	[Phase.Phase3]: P3Consumes,
	[Phase.Phase4]: P4Consumes,
	[Phase.Phase5]: P5Consumes,
};

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	demonicPact: 110,
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
	spiritOfZandalar: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
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
	profession2: Profession.Alchemy,
};
