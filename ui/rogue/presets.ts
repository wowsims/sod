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
import P6BackstabAPL from './apls/P6_Backstab.apl.json';
import P6BackstabIEAAPL from './apls/P6_Backstab_IEA.apl.json';
import P6MutilateAPL from './apls/P6_Mutilate.apl.json';
import P6MutilateIEAAPL from './apls/P6_Mutilate_IEA.apl.json';
import P6SaberAPL from './apls/P6_Saber.apl.json';
import P6SaberIEAAPL from './apls/P6_Saber_IEA.apl.json';
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
import P6BackstabGear from './gear_sets/p6_backstab.gear.json';
import P6MutilateGear from './gear_sets/p6_mutilate.gear.json';
import P6MutilateIEAGear from './gear_sets/p6_mutilate_IEA.gear.json';
import P6SaberGear from './gear_sets/p6_saber.gear.json';

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
export const P6GearBackstab = PresetUtils.makePresetGear('P6 Backstab', P6BackstabGear, { customCondition: player => player.getLevel() === 60 });
export const P6GearMutilate = PresetUtils.makePresetGear('P6 Mutilate', P6MutilateGear, { customCondition: player => player.getLevel() === 60 });
export const P6GearMutilateIEA = PresetUtils.makePresetGear('P6 Mutilate IEA', P6MutilateIEAGear, { customCondition: player => player.getLevel() === 60 });
export const P6GearSaber = PresetUtils.makePresetGear('P6 Saber', P6SaberGear, { customCondition: player => player.getLevel() === 60 });

export const GearPresets = {
	[Phase.Phase1]: [P1GearDaggers, P1GearSaber],
	[Phase.Phase2]: [P2GearDaggers],
	[Phase.Phase3]: [P3GearMuti, P3GearMutiHat, P3GearSaber],
	[Phase.Phase4]: [P4GearMuti, P4GearSaber],
	[Phase.Phase5]: [P5GearBackstab, P5GearMutilate, P5GearSaber],
	[Phase.Phase6]: [P6GearBackstab, P6GearMutilate, P6GearMutilateIEA, P6GearSaber],
};

export const DefaultGear = GearPresets[Phase.Phase6][0];
export const DefaultGearBackstab = GearPresets[Phase.Phase6][0];
export const DefaultGearMutilate = GearPresets[Phase.Phase6][1];
export const DefaultGearSaber = GearPresets[Phase.Phase6][3];

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
export const ROTATION_PRESET_BACKSTAB_DPS_P6 = PresetUtils.makePresetAPLRotation('P6 Backstab', P6BackstabAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_MUTILATE_DPS_P6 = PresetUtils.makePresetAPLRotation('P6 Mutilate', P6MutilateAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SABER_DPS_P6 = PresetUtils.makePresetAPLRotation('P6 Saber Slash', P6SaberAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_BACKSTAB_IEA_P6 = PresetUtils.makePresetAPLRotation('P6 Backstab IEA', P6BackstabIEAAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_MUTILATE_IEA_P6 = PresetUtils.makePresetAPLRotation('P6 Mutilate IEA', P6MutilateIEAAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SABER_IEA_P6 = PresetUtils.makePresetAPLRotation('P5 Saber Slash IEA', P6SaberIEAAPL, {
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
	[Phase.Phase6]: [
		ROTATION_PRESET_BACKSTAB_DPS_P6,
		ROTATION_PRESET_MUTILATE_DPS_P6,
		ROTATION_PRESET_SABER_DPS_P6,
		ROTATION_PRESET_BACKSTAB_IEA_P6,
		ROTATION_PRESET_MUTILATE_IEA_P6,
		ROTATION_PRESET_SABER_IEA_P6,
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
		[RogueRune.RuneMutilate]: ROTATION_PRESET_MUTILATE_DPS_P6,
		[RogueRune.RuneSaberSlash]: ROTATION_PRESET_SABER_DPS_P6,
		[RogueRune.RuneCutthroat]: ROTATION_PRESET_BACKSTAB_DPS_P6,
	},
};

export const DefaultAPLBackstab = APLPresets[Phase.Phase6][0];
export const DefaultAPLMutilate = APLPresets[Phase.Phase6][1];
export const DefaultAPLSaber = APLPresets[Phase.Phase6][2];

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
	SavedTalents.create({ talentsString: '00532012255104-02330520000501' }),
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
	SavedTalents.create({ talentsString: '00532012255104-02530500000501' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const P6BackstabTalent = PresetUtils.makePresetTalents('P6 Backstab', SavedTalents.create({ talentsString: '005203102-0233-0502530310321005' }), {
	customCondition: player => player.getLevel() === 60,
});
export const P6BackstabIEATalent = PresetUtils.makePresetTalents('P6 Backstab IEA', SavedTalents.create({ talentsString: '00520312-0233-0502530310321005' }), {
	customCondition: player => player.getLevel() === 60,
});
export const P6MutilateSaberTalent = PresetUtils.makePresetTalents(
	'P6 Mutilate/Saber',
	SavedTalents.create({ talentsString: '00530310355104-02330520000501' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);
export const P6MutilateSaberIEATalent = PresetUtils.makePresetTalents(
	'P5 Mutilate/Saber IEA',
	SavedTalents.create({ talentsString: '00530312155104-02330520000501' }),
	{
		customCondition: player => player.getLevel() === 60,
	},
);

export const TalentPresets = {
	[Phase.Phase1]: [CombatDagger25Talents],
	[Phase.Phase2]: [ColdBloodMutilate40Talents, IEAMutilate40Talents, CombatMutilate40Talents],
	[Phase.Phase3]: [P3TalentsMuti, P3TalentsMutiHat, P3TalentsSaber],
	/*	[Phase.Phase4]: [P4TalentsMutiSaber, P4TalentsSlaughter],
	[Phase.Phase5]: [
		P5TalentBackstabAssassination,
		P5TalentBackstabCombat,
		P5TalentMutilateSaberslashCarnage,
		P5TalentMutilateSaberslashCTTC,
		P5TalentBackstabAssassinationIEA,
		P5TalentMutilateSaberslashCTTCIEA,
	],*/ // Cleaner to have only current phase talents available
	[Phase.Phase6]: [P6BackstabTalent, P6BackstabIEATalent, P6MutilateSaberTalent, P6MutilateSaberIEATalent],
};

export const DefaultTalentsAssassin = TalentPresets[Phase.Phase6][2];
export const DefaultTalentsCombat = TalentPresets[Phase.Phase6][2];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase6][0];

export const DefaultTalentsBackstab = TalentPresets[Phase.Phase6][0];
export const DefaultTalentsMutilate = TalentPresets[Phase.Phase6][2];
export const DefaultTalentsSaber = TalentPresets[Phase.Phase6][2];

export const DefaultTalents = DefaultTalentsCombat;

///////////////////////////////////////////////////////////////////////////
//                                Encounters
///////////////////////////////////////////////////////////////////////////
export const PresetBuildBackstab = PresetUtils.makePresetBuild('Backstab', {
	gear: P6GearBackstab,
	talents: P6BackstabTalent,
	rotation: ROTATION_PRESET_BACKSTAB_DPS_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Backstab',
		'https://wowsims.github.io/sod/rogue/#eJztVl1oHFUU3nNndnb2pDOZXFKd3KTpuqJcVjfMzuyum8S6m+hDmorsQ5E8lSiNaCkajBT1xU2a0n+7+IO1FZO0iqGo1NinIBIFpXmrv8RasSCIrT4IiqlW0Duzs91hsw+xVhD0Psycc+45Z77v3HPvHTRUEiOcFMgwlEixTKC/zYRZAicJFOAMwHmAInkQTgMMAA2xcy2oFLff98TIo0bYvD2mJF4HbDJ2HtTNxV5eeouhbMwe0jFinDikm2+2irmlKd38SOELz7YJpSQcv9H5V8uucmBGN7/r5c+JqCbj+2ndnGzlP7/cgapxfEbnpy65PpOVzF+eMoUyLpSvdb7n1zbhc3ZK57sn213pVZ2XhSQby0/r4nl0xgVwTGT/wRW+mPIE2XhpWbO5GjEIV3KC3IBcJuTxSAnayyBNgzwLsADwMRADliDc862M50no//GPDlr6z5f4Ign1t+RAdJ9oSSjCEAzDKAyut6yMbTkpy05atuMkrYxlZxzX4NgpMVdUh+T7Q9veQVNCipFtYZTUy1o8gkJIpsdQ7NN5kChkfJMzhtdhlAlF/SXNSTyKQkzaXZkxXIsqU1BWz02Db04FzMZvZc03W12pMToO7CncwR7DHM2ikdBVmZInQ64nkpTlW6QJCDHxWYS0ZwlT6RnZt6SCkYPVyEx9rrS1MtJ9Gb/30veAvQv4NrCTgPfSzfWhzupgaPEmjM6BS/2FDzuC6k/iFArkrQMpBb70l/N6JRXnFJUZ6QnRArsD09TGGOtEbQ6QEiMsHN0z0Ys7ISu+qsfXIM6Bu4CXNS57ecZf0WmS3YLr6ToPh1JDlrUCGJiKivrpB8Cj9DCw5wH3AdsN2E1v88NqhcpUE10hncWWeHMQ3M5PGA7SgRWhDjYnND/UMzSMrK9p1us9cT3QIrsH72R9mKcbGsC6PrFWVB7nodrC1XLX47WuNDk9Cuww4AFgewF7aG5VZGncQH0OmjzIYXX/EFdwE924GraNQlfSjfhTdJhtwbtof7BB3KutvnXt+s6sW2nHsqs7e0Hn4Iu73C9sYhsxzzZgK6UNuN9Eb1xhTXus5BqrtLc84v6lfSyPaWY32FqdrANZwhRWbT641pX5XA1PPytcZYpaA69j7UipUdkPx6BK23VULzxMb2UJjNFO7EgwQaz5fVjjJVPUswcHY1uDnMQPxuhpuETgCIQmiKCsbo0vSNHKqbyl0FIRugptR150x2K+u2I5k7/hAW98nrfLkro0oVH17pEdI9tjWcvMJ/5lPw18b+Ga5Ln5jb+Vx2B/XBscYqS85+hrhe6HHvlx1+Lmz/J9/kxhCP4ETjAc6g==',
	),
});

export const PresetBuildBackstabIEA = PresetUtils.makePresetBuild('Backstab IEA', {
	gear: P6GearBackstab,
	talents: P6BackstabIEATalent,
	rotation: ROTATION_PRESET_BACKSTAB_IEA_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Backstab IEA',
		'https://wowsims.github.io/sod/rogue/#eJztVk1sG0UU9huv1+uX7mYzSulm4gbX0Ghk5Gi9/sFJKN4EDmmKkA8VyqkKqEFQVRARVAEXnDRVKbTU4keUFpGk5SeqKCqhpxxQQKJqbuVXoYCIhIQocKhURApFKrPrdWzZPgQoEhLMYf3em/nevO/NmzdGXSERwolNhqFA8kUC/RsMmCVwmoAN5wGWAR6AcwADQH3sWgvK+d33Pj7yiB4wbovIsbcAm/S9hzRjsZcX3mEo6bOHNQzqJw9rxtutYm5pSjM+kvnC821CKYiF32r86xVHOTijGd/38hcEqkn/cVozJlv5z6+GUdFPzGj8zBVnzWTJ85dnDKGMC+UbjT/1a5tYc2FK4/sn2x3pdY0XhSTpK89q4ntsxgnguPD+kyN8MeUKkv7KimpxJagTLmcFtQGpSMhj7QUIF8E/DdIswALAx0B0WIJAz3cSLhPf/+MfHbTwn0/xD8TX35IFUX2iJCEPQzAMozDYYZppy0wmrLhpJZNxM21a6aTQzaSVEFN5ZUi6z7drGQ0/UgzuCqBfuapGgyiEeGoMxTWdBz+FtGdKjuENGGJCUX5JcRINoRDjVld6DNejwmSUlK+mwTMnqsz6b0XVM5tdiTHayW5GFjMUiarzgJToAQF94yJjYheEhOtJXF06DuxJ3MMexSzNoB7TBIA84XOmkSRMz+KfAF8JmHItAep/TvIsiWrkYBmZrvWVMuuRzo/+ey99H9h7gO8COw14D91eC02uLQw12oShOXCIvXQ2XK1eFs2qym9NkP6qnf60Xzf1op1RiZEeH7XZ7ZiiFkZYB6pzq3l3WqeLOynJnqpF1yHOgXPQV1UuuX7GX9NonN2CN9KNbhxyJbKMWRUDU1BWPv0QeIgeAfYi4NPA9gN201s9WCVR6bKjVdIZbIk2Vwe39xOGg3SgDprE5pjqQV1DQ2RtTjNuZYlXhObZ3XgH68Mc3dIgrA2x9SLzOA/lUi+nuzZec/Uy0GPAjgAeBHYAsIdm10SWRnXU5qDJDTmgPDPEZdxGt66FbSNoPd2gN0WH2Q68k/ZXF4jzAtaWrlVbmTUnnTStcgdY0Dh44j5nh21sK+bYFmyltAH3zfSmOmvKZSVVWKXc4xHPNO1jOUwxq8HV6mDh+tYhzro0n63E08/sv+iiUsAbWTtSqpfuw3Eo03YWKhcfopx1YjjGBKnmD2Cd60hWLhwajOys5iMa2eg5uELgKPgmiKCr7Iwu+EOlxr3DbikJXXbb0ZedsZjrLlnO5zbd747Pc1bRryxNqFS5a2TPyO5IxjRysX/Z/wp+wL4ufjpP/S0/Ort2feIQI+F+R9+0ux98+NK+xe2f5fq8GXsI/gCC7iJ2',
	),
});

export const PresetBuildMutilate = PresetUtils.makePresetBuild('Mutilate', {
	gear: P6GearMutilate,
	talents: P6MutilateSaberTalent,
	rotation: ROTATION_PRESET_MUTILATE_DPS_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Mutilate',
		'https://wowsims.github.io/sod/rogue/#eJztVU1sG1UQ3nler9eTerN5dYrzSivXQPRkJdHu2tsmcandokJkVcgHDhFCqKAaaBSVtJaqlpMJ/ZVaySr0QC5JRSSqKkFtKD30FPWU3MJfFAKIAhJCFQdOKE2C4L312uQHTg1SpXYOq5lvZr43Mx6/h6ZO4oSTHDkIZVKoENjXEoOrBK4TyMEMwF2AAnkLpgB6gCqsbKBW6H/tZPGYGYztjmvJMcAGc2zYiE1n+LkbDFVz6ZKBIfP+ZSP2SVT45oTvc41Pvt8ijHFh/Gzw72/GhHHhihH7NcM/EFkN5rKIPxXlXwmPbn50xeA3F2TC4mWP+Vsv4U9h/GDwc/dbRMz8sMHPntoqtE9HDV4RmmreO2OI75lRWcA34qjfpLI4G5GKai6MGg7XQybhWqdorketEHIiVAa1AoERUK8CTAJ8AcSEOQh2/6LiXaI8lv9VaPmRH/E9ouxr6gSxfWIloQC9cBAGIL/NstyUlbKtlOvaVrrdclIpy3UsIa5lF/Re9XWl7yctFsAo6kxDVb91Np/QUWjttl3CzRjqC2JAX4rUQKuETajdhgAFNxFC4WxPl3ALhpkw9D/SnCTCKNT2VIdbEqwhJvM/6/dRR6LN/lnfjYAP2ytgc7ES8WGrwy7RF9h+3M260UwaukrJO4oMQmJb2EqfrqF5H3WstXFpS1ZgLmeoyki3QnNsD6apg3G2HSMTgJSYQVGLvD8iiQYMX1M13zQSmxAnQPa1FOGqV9u7owbdxVxsY0nktLXml32vZQtPgGQSl44/AlpgL+JzbC9m6bNelRoNDIIivQhpfCLZrAco3q4TVh22H1pv0bbqo6Y97HlvNp10578wNiYjfqIHuNiUaFxZ5XtfMu93ELcqzbAu7KBt1RlUK5c3pCSFFaQuhofqbckyJg0O9BX28qrGVmc8mWQCabwDm7xzNX3+Yj5+aD0fKs5mf19+POKviIBpiR3FN1kRo5Sua9LFZ+hT69CU17r6T+vp+q6snpA/2tquOG79ULkeUUb9bVpB7UWIB4bKFYjT7V572n+0l/aixSM0MAULBIZAGSSiLP1QYjIQrv5zX801VZWOXMvQh1Kms11VZCa74w1PZrNOJaDPDUaofqB4vNgf32nFssmH7GHh53MbwtM6/kA8JvtrY+oQYnvfgY9zXYff/v309EtfZ/f6nlwv/A3G2evT',
	),
});

export const PresetBuildMutilateIEA = PresetUtils.makePresetBuild('Mutilate IEA', {
	gear: P6GearMutilateIEA,
	talents: P6MutilateSaberIEATalent,
	rotation: ROTATION_PRESET_MUTILATE_IEA_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Mutilate IEA',
		'https://wowsims.github.io/sod/rogue/#eJztVVtsFFUYnv/M7OzsX3Y6PWlleopkXbQ52WTLXHZEWuIOqzGbjZd9UNknI4QSII02NjFKfKgNN29ho/IgL4A2kRgwWtGHPpiGJ/qGRklFjRATQ4gPPhhTykXPmZ2FYYEnMTHR8zDzX7/z/9+5ITVIjnASkio8C5Ok3iRQWW7DEQKfEAjhFMBZgK1wEqAKVGHvd6NeH9v48ugLVspel9MLRwG7rKMHTXt+hO/9lKFmXXrbxLR1cb9pf9wrfAvC95XO597pF8oxofxs8h+P20J587Bpnx/h+0RWl3VZxO/s5d8Ij2F9cNjkxxdlwtL+CPn7KOGKUH4y+d6L/SLmzEGT79k5IKTPpk3eFJJmXdhtiu/uaVnAd2KqX6WwdDorBc1anDY9bqQtwvUHRGtVrUnIS+lJ0JqgHgLtCMAcwNdALFiA1PAvGp4lyv/jHx108j9P8QWiVMxoP0IdGuIQjkPtbscJfMd3PTcIXKdUdDzfdwLPESNw3LrWMDYp23cZtop3ocF01Iwv9tTyBgqp6LoTNmAfZranMWVcynKt7XAmsAf1WVApBHnhRLUYTAiEDJORf5Q4yWdQiMXSkLD3YpqJEOPzsdjqS2tfPN8PhyA2ewmztdTMxmZnyJ2gDhvC+9gqtAqmoVGyQ5FBSDwHuwtZaakpcg6EQP6syyNUY2RYoSF7EEvUwxxbidkZQEqslJj0XXFTZPNdmPlI02PVzC9DnIF2q1ERr06b1GcucjrYDt3XzpwBqV79faCdJduO+6QVFuIaGiTxO5I6irkFRp09jg/TCq4oMEOn3SdgWRSvG2fequU2t7otRYSoVJ2Cdv+dFK1xIprFJUifYI9hhYY4UOgXOeYsdEWIKePLDTzVSvejdLgBUDKsX2fYjxZZptAN7Cms0WqyzeWFPjE7zt7QDoIb4eoJ3FJskbhRoa5zbffQdWwYV9NikvTE+siQm/C869mSuYfo+mRZPfnuTso7AfzORoMI8Y0G1+krbAc26NPJRZRvR9xBm2o/uAV76uugRGsIt1nDwFsZH5Bzz8UnaM7kEB8P8V7RKnsEJSM3gw/Se29PbLusILiGtIKxeH8k+o684rkcPwmLBA6AMkV0qhmbLcjPqZnW3fJM2NMShsL+A+/JMV9e27KcKt+zJRqny15TNRamstR4dPTF0bHc/Y5dLvzLnj7+WnhHcAaP/S0ci/15Z+oQw42+4x+Ga7c9/9uu+Se/La+PPWED/gKA7wPp',
	),
});

export const PresetBuildSaberSlash = PresetUtils.makePresetBuild('Saber Slash', {
	gear: P6GearSaber,
	talents: P6MutilateSaberTalent,
	rotation: ROTATION_PRESET_SABER_DPS_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Saber Slash',
		'https://wowsims.github.io/sod/rogue/#eJztVV1oHFUUnnN3Mjt70rmZ3KY6ubWy3do6LM1yZ380f3a3sa0hlLIPPoQipZXGnxI0uiC2IqRrJLU0sFBfjA9JsMVSWmlD9SEghD41b22hEqvFFkGk+KAPlmiFeufuzGa3qU9WEPQ8DPf8f+dnz6JtkjhxSYHshVFSrBDoa3fgFIFzBApwGeAmQJG8ApcA+oFp/HgLGsXhfQeH3rSbnN64kTwD2GyfmaLOQo975DxH3S4vWRi17x2jzmdtUrcodVcMd/54u2TOSuZ76t644Ejm2Ax1fuxxP5RezXZ5gjpjbe6VO+1o2p/MUPfCku8wNqEif6McDkvmO+oe+c23uT5F3fGxtf7rJHW/fG+tTP3xDJXfo+NUAvhlmjo/+Y+vp9RDt5dO0LRrRm3iGp2yuH69Qsjb0VHQKxCZBv0UwDzAVSA2LEJT9w863iTa//SPEhv9z7f4NtH6WjtBbp9cSSjCIOyFERhYJ0QuIzKeyORynsh2iHQmI3JpISknvKI5qL+oHfg16kSwDU1uoG5+MT6QMFG+OjyvhKsxeqAJI+ZdKxSKEraiMQcRBrlEFKWyI1vCRzDGJWPeybokEUP57MikciUZNcp9/8+HA2nal64Jcn07DYHYqxPbv1esQCxSXok9x7djL+9GO0lNnZFDmm+ExBO4iT0RSgdCae5+u6zwEdh/9DCdk26NFfgWzLI0xvnjaM0CMmI3SSz+/bASzRg7rRsBSxOrEGfBr+uu5eoK2+ETlB0FPg74LntHpYqwSBk0PwdCFrfwXnwsyaW05SKsUrENU56j+P6qRSbUzj1Im8YnkxuXtaSmtWeJZYIqs2qo+id1rId3oeAp3MySIVh/CCoN/BWIYCasyHfhNtanjI164+sTA8uY7q9yRYs7hYIjb62K+Czfinn2jLIyGrrzaHKNjIRzNZhVhReYLk9R1LZJFZhim6vDmQU/kX+6fQ9oABWbVEp55JXzPHWBvcB3N8Bp9FjRo7qy6+Ohll4dLPKt14LdlWJW4m/gy3wI2xh7QPSNbMOKFmSwJWkF+xr0JFzixk4FDQk7/LSoJe3nO9TPIfRbOZl6v3BWyza11Rm5BEsEJkErEwnI3J+Yj8Sqx2RPobX6SBXaJz/yaSHfVZVczq9/SdFX+XQlYi6WLWbuHHpraDj+lHDyyX/Zf537QeGhxNl09m/Fsfm9h4NDkqe+I58Wul59/ef3F56/lt8aaAqD8CcV3P3q',
	),
});

export const PresetBuildSaberSlashIEA = PresetUtils.makePresetBuild('Saber Slash IEA', {
	gear: P6GearSaber,
	talents: P6MutilateSaberIEATalent,
	rotation: ROTATION_PRESET_SABER_IEA_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Saber Slash IEA',
		'https://wowsims.github.io/sod/rogue/#eJztVVtoHFUY3v/s7Ozsn+7p5GRrJycq62rjsDRhZndWc7M7jbcQtO6DYBCRKo2XEjS6IF4Q0jWSWhpZ0BfjQxJssUgrGmofApXQp+atFZRYFSuCiIjoS4lWqOecndnNdqsvVhD0PAznfP///fdzBk2DpIlNfLIbpkipSmB4qwVHCHxAwIczAOcBnoDTACPAIvx1E/XSxCMvjD9rxqyhtJ49CthmHp2n1uqgvf9DjppZWU9i3Lx0kFrvp4RsTcjO6vbKG53icEwcvqX2V8ctcTi4SK3vB+03BavNrMxSazpln73QiYb5ziK1j69LwvSssvyFIuwTh6+pvf9XqXNuntoz011yd5jaJ1/pEq7fXqTie2CGigB+WaDWj3Lz+bzaaOb6IZqzjbhJbL1PpDaiVQl5Pj4FWhWiC6AdAVgB+ASICWsQG/hOw/Mk8v/6Rxeb+s+X+AcSGW7vAzF9YiShBGOwGyZh9DrHKeSdvJtzCwXX8XqcXD7vFHKOWAXHLRlj2qORvT8ZVhRTaHAdNePEzGjGQLHrcd0ydmB8bwyjxsVkCDplbEd9GaIMCpk4CmFPoYzXYIKLg3HBs0kmgWLb4/UKPIVxLvkfTQRoXqJbAl9fLkAA5yTcUYPN36rJwLJbZnfzO3GID6CZpYbGyIsRqYLEdbCb3RSioyFauFzPc6R/8/dBpnEyEGE+34Eey2GaX4/JJUBGzJiIRL4eyUwbJt7T9OBIM5sQl0BmdTFpayqyfYcoOwB8BvBl9pJyFWXRCkSkDwQPd/AhvDbLBbr5FGxStnVDPEbpPTWNfChdvpI0hzdntzWkpC41l0jSAJVmTVFVT8jYIO9Hh/fidpYNg5UtUG7gz4IIOsJKfBfewYaVsr5R+dzsaCOmy7NsKXGfo8IRLy27j9+Lw8zHrmyn4NBlaFMWY8bHD9ixsASSrjcZ3JxNCqSeXl6NkqSoEG/nO7HIbmvhebg1u0W4weV63jWBG6g2xsKpD6eqWC/bXuv2EsjI5Z9AMqApqMScEop/hiKvUBvYQ/zBpnCaGS1F31DHjfYwkusI7sU3TwVXQcCszJ/Bx/k4phi7gvVt7MaWEuRV6bRG6bz6rWiuVFCQsGW3OnWnI/wudb9CXmurN/LC5jd0GrNo8+6/HCUvHJPJ07BOYA4iFaIzzdhjQmYlmqi9Yw/77bVNr98595Zcq8X+GnKmeMNjan1WzFWjxlolyYx7xp8bn0jf4ljF7L/sN2u/5l8VO93H/pYdk1+6OnGI5arv5Lt+/5NP//zq6v2fFncGEn8M/gCWAwy7',
	),
});

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RogueOptions.create({
	honorAmongThievesCritRate: 100,
	pkSwap: false,
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
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodGrilledSquid,
	mainHandImbue: WeaponImbue.WildStrikes,
	miscConsumes: {
		jujuEmber: true,
	},
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const P5Consumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodGrilledSquid,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.ElementalSharpeningStone,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
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
	demonicPact: 120,
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
	improvedScorch: true,
	mangle: true,
	markOfChaos: true,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Alchemy,
};
