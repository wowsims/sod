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

///////////////////////////////////////////////////////////////////////////
//                                Encounters
///////////////////////////////////////////////////////////////////////////
export const PresetBuildBackstab = PresetUtils.makePresetBuild('Backstab', {
	gear: DefaultGearBackstab,
	talents: P5TalentBackstabAssassination,
	rotation: DefaultAPLBackstab,
	encounter: PresetUtils.makePresetEncounter(
		'Backstab',
		'https://wowsims.github.io/sod/rogue/#eJztVV1sVEUUvmfu3d3Z09zrdAS5O23tdlVyXSzcu38ppbpLlYhoyEZ52BhjqlIE0mgjiVGfaqEqqUBRo1II6ZsV1OgGX/pgGh4MfSsmNiTSiPGNGG3iT7Ga6My92227LfIA/jx4HiZnzpzznW/OmR9klMSJQwqkC/pIcYhAZ6MNowQ+JlCASYBLAEWyG84BbIURAlwTMwzDxZ4nXux+joXsjng4+QFgHdt/yLInNjl9nwg02OgRCyPs1BHL/miVXBsbtOzzYWf8zZicvC+t31rO9KyanJaTact5S0bVsX4J8bXl/HSyESkbOGQ5Z64onwMB8ldnbDl5OfB57deY9Pls1nRePdAgtckfTGdIagabPWzJ8fIrisDclGl/p5TpOV8x2IlZMxWnOgMH2uT+tpIX1vSBPQT6COjjAF+AdgFC7VMGXiLa//K3Cp+Bf5vCPy2XidZZ3wajoI4eFKEEXdAL22Kum02n0p6bzWbl4LW6qXTaTbW62SJ9Uts7E7V1rMfI3hDq9DczEUYDibsPVyMVUqcXRyARxQiGWr312aqZff+zWTXnFsxzQ/Nmd723j98vtmCHaEeWtKjByUua8OE9F9fy2+et2+at2Vq/jCskK/b7Jm4I0q7xdeJOjItb0SwDcsJCkp6622aiDqOnjHAw9UPefoTfJZLYIpqxiTcglkFyor9kHFITrLzppz38BIhjgIdBDALu4A9jU7KBhjk7C6bvG6EXS/GdjqHcETI+UZ3r/aAFlmXU025AqwyK1jufN2Kaez6qwdnYiqgesnLtggqW75S/mWbe5GcJL2TJuYsyCLXFoyWH8uMg3gV8HcRBwHbeVomqZRteVPoc8gRDqwx1fv4QHSw5YXyQP7AsNo03Jc1KrG9YObS2tzmfnVriRbEd7xWbMc/vXoaewTXJ1bKyOFZtWLU4NZRdMe/Bu8TjeB/vXFxx9YQHHBZTrynf0ib9eLIxdQtGfdhxy4GKOqA47wfRB7hHPI3d/CkfBpbUUzUWrnZcshgd9jOoT6YhGZMUrLEl5Qr8UkiHQZ3Hb55J3YwREaiVeyhj+W6xCx8Tj6LHN6zQ079EvoPftkKxVSuNhVZmqpnuER3oiQ3LSpjxsxhXybJQr+uI909wq1iHLbz5GtdwyfW49sOgHqjzR6H3HFwhMAxaP5F7pzsTZ/Vo8IJuL9QHyqpCbPiYkon8xsAymW/Z5ctUPvWGTi/0m5w+1P18d08859r5uJ78j33pzsHCDcFZ++F14TDxx43hIcXzx973Chv3PDszMLHjy/zmykqhBH8CZHYaeA==',
	),
});

export const PresetBuildBackstabIEA = PresetUtils.makePresetBuild('Backstab IEA', {
	gear: DefaultGearBackstab,
	talents: P5TalentBackstabAssassinationIEA,
	rotation: DefaultAPLBackstab,
	encounter: PresetUtils.makePresetEncounter(
		'Backstab IEA',
		'http://localhost:5173/sod/rogue/#eJztVV1oXEUUvmfu3d3ZE+51MrZ6d5KYzarlujXtvftHmkZ3GxVrlbJoHxYRidrUtgQNFkR9imlDrbFtqqI2lpI3Y/1BlwqSBwl9kOQtFQwBW1R8E9GAP1uroDP3bjbJJrUPrT8PnofhzJlzvvPNOfODjJI4cUiB9MAAKY4Q6G6xYZzABwQKMAMwB/AxIbthCmArjBHgmphnGC72PfJs71MsZHfFw8l3ARvY/sOWPb3ZGfhQoMHGj1oYYaeOWvb7a+TaxLBlnw07k6/E5ORtaf3Gcs5X1OQdOTlvOa/KqAY2KCG+tJyfTjYjZUOHLef0BeVzIED+4rQtJ88HPi/8GpM+n1RM5+CBJqnN/GA6I1IzWOWIJccXDyoCUxXT/k4pF2d9xWAnKmYqTnUGDnTIDW4lz9gDEBsBfQz0SYDPQJuDUOesgV8R7X/5W4XPw79N4Z+Wb4nW3dgB46COHhShBD3QD9tirptNp9Kem81m5eC1u6l02k21u9kifVTbOx+1dWzEyN4Q6vQ3MxFGA4m7D9ciFVKn58YgEcUIhtq9DdmamX3/s1kz5xbNF0cWzO4Gbx+/W9yFXaITWdKiBifPacKH91xcx29asG5bsGbr/TKukKzY75u5IUinxteLWzAubkCzDMgJC0l66m6biQaMnjLCwdQPee0BfqtIYptoxRbehFgGyYn+knFIXbDyph/18RMgjgMeATEMuIPfjy3JJhrm7AyYvm+EnivFdzqGckfI+ER1rg+CFlhWUE+7Aa0yKFqvf9qMae75qAZnE6uiesjK9QsqWL5T/mZaeYufJbyYJecuySDUFo+VHMrfBPEG4EsgDgF28o5qVD3b8JLS55AnGFplaPDzh+hwyQnjvfyeFbFpvCZpVmN9w+qh9b3N+ezUEi+K7XiH2IJ5ftsK9Axen1wrK4sTtYbVilNH2RULHrxHPIx38u6lFVdPeMBhKfW68i1v0o8nm1PXYdSHnbQcqKpDivN+EAOAe8Tj2Msf82FgWT1VY+FSxyWL0VE/g/pkmpIxScGaWFauwC+FdBTUefz6idS1GBGBWr2HMpbvFrvwIfEgenzjKj39S+Sb+Y2rFFu10lhsZaaW6XbRhZ7YuKKEGT+LcYksi/W6gnj/BLeL9djGWy9zDZddj8s/DOqBOnsM+qfgAoFR0AaJ3DvdmTijR4MXdHuhMVDWFGKjx5VM5zcFlpl82y5fZvOpl3U6N2hyel/v07198Zxr5+N68j/2pTuHClcFZ917V4TDxB9Xh4cUzx/73yps2vPk/ND0js/zW6orhRL8CQTAGmY=',
	),
});

export const PresetBuildMutilate = PresetUtils.makePresetBuild('Mutilate', {
	gear: DefaultGearMutilate,
	talents: DefaultTalentsMutilate,
	rotation: DefaultAPLMutilate,
	encounter: PresetUtils.makePresetEncounter(
		'Mutilate',
		'https://wowsims.github.io/sod/rogue/#eJztVVtIFFEYnv/M7O7420zjwWg8Zmxb2bCkzK6OeaPdLEskyqAHHyu0i0hqQlRPIl0UIrfAhyTItyQoSvQ1wid90ygRUjKiF4nypdbMLmdmdr3Wk3Z56H84/Oc7/+X7z+U/SGXiJwaJkgo4Bi2kKkagbJMOPQQeEYjCMMAkQBU5DYMAFdBNgAqsXUFvVf3xi7XnNI9e6vcGewBTtW8xRR8qMdoeM5S0ng4Vfdr9DlV/mM6VuzOKPuK1jTpV/Y1qTPTpfBLnkwnV6OAeqdocn7xSjed8RdZmOlWjbyaDw7Mc5lFfOg5fXZu2zxncZviDYly7nMm1J3HFiHFN0qauqnyM37CTz44q+jtb6Y87iqTdiSthvyxqYEAhr62CXPC1gBQDsRvEpwDPQBgDT/GohJNE+C+/Veg0/G0Kf1qmiFCmOvcOqqCav7ZGqMwyTSsvbIbMfMviQ44ZzsszrbDJxTJD1fIJoe66VxcxDX11HhTlL0rAixISsxk3oMy4Lo93QyAFfejJCeVa87D2/qMyDxcswLOxJGzmhprpAVaOpawYtaAqS5RcEpgTPmxiNt2WRCsX0GV2+SbjrLS5ElrB9q8mkGU67PjDpRIjxQI1WDZmsUzMoBsRe4ETlj/lGwRTesEumrcDO7PcX0+PsMNYTvdiVjBT9lJtABRKNA/65PFqf40h2WYI+U5CkYqtILiItZzCLoeCzLuOsyvFrBALqOvnXeQXQjWwbgkl12K+uJDJkku0hBVhLt2JSiA1ydzuXbYHLOGS0jVflu38VDWANrEG3MfKMEp3/8TDLhd+Ve7igCiEd6DP2a3XZxO3wG6T7u2xDVzt7RijzawJT7FaTKd0ReEWbqdbV6B5uD6oJM42sdPJA19sltzHhQO3kumpn23GdEYT57EotGPB/wOKzKU9chMaB2GGQBcIrYQnlWsCA2KK+7wORdNcJT2a0XXblqFIkYsMR7acdGQ0Er4lymOtCpUP1p6vrfcXmHrELwb/sX5vtEfXJE72g1XF0dj3teHBJeSMjfeiRWcapq8MHX0R2ZNYiVbDDyGt3Iw=',
	),
});

export const PresetBuildMutilateIEA = PresetUtils.makePresetBuild('Mutilate IEA', {
	gear: DefaultGearMutilate,
	talents: P5TalentMutilateSaberslashCTTCIEA,
	rotation: ROTATION_PRESET_MUTILATE_IEA_P5,
	encounter: PresetUtils.makePresetEncounter(
		'Mutilate IEA',
		'https://wowsims.github.io/sod/rogue/#eJztVVtsDGEUnvPP7O70NDOmfyqmf1PWRmSsS2Znd1yqsauCqkvrkugLSlRcGhoNwZM0ri9agkQjUU8aCaHhxYOIJ31r0UaCqHhrBC9s1fWff7aLti6J64PzMDn/+b/zncucOYOGSsLEIilSDXtJZROBUmZCK4FLBFLQAdADsBFuAZRBCwEqsTs6Bitr1+6u2W4EzJJwMNoKmGu8OaGb7TOtQ5cZKkZro44h43yjbl7M58qZPs3sDHLQew56olsPr5j8kOaHh7rVyD1yjbNpzXykW3f5jWr0ndCtK30F3Nzvs94XDu/4gWMOvS7gmI7nmnVwXyHXrqc1q4lritF7QOfP9BEveH+3Zj71lKtpoSjG6bTmhFXZAAum88rKyK7QXlCaQG4B+QbAbZDuQaC4W8EeIv2X3yr0BfztFP609BKpNG86tII3elAJVVANdVBeZNtu3LFjjuO6MTsx2Xbicdt1bC6uHVsnbe4KmTLmYWhzAGX1jRYJooLErseRqDKuqw9aIJKDIQxMjk1xs2bj2Usta576ydzfNGC2p8Tq6Xw2F0tYMRpRXVUo2SMxQR+zcTwdN2Atz1gdezAuYTOelfF2Ji1j84Yjcn6QyLVFdvzLpQojxRKdyCZgmI1GrQ2QEiPACz3O94QWycWc80rQP4rgJ5fTSSyKY9kYLKKFiG3Aq1NfJSwyyNlDq1draTVbjcvZUqygi0UaQSo3gOTdIsRwVHQkTwyvZVkGLnxotgbeIr0NcgV9QL2+0gqwAQ+6lFXgXDoHi6KF3MW4CZqAhdQHVeH1luIzJgSj/Flwd3BXpomuqHwT0gq2GEtpCgujBdxHv/ZlZOEeF+7wBeGIqJZJ2oeIHD0XupAtwFm0xO9oG3hxvL06lCKnWVx6G3i4gm/oFtAdrB6XsUpcQhcNw+C1Ab7Whm8FcIowJN7a462Z8fUWvD/2HE63sE24hq3CcTQy5E0mBpOhS+NDUHHRIeVTh4a4ZYN505jPaOYdfUYhEPy/Jcbw+zPrldF5FOpuQR+BZpAaCA+vro/clHP8RbEklecr+amC5lOetCdn+JaO5NgNQrqTzjFZvdegUXVRzc6a2vBU20yG5eg/9ueyDqd+Cc/4Cz/FY7APvyYPLjHxrDuXmrFp24v97Su6krMzN6kq+AjqqvoQ',
	),
});

export const PresetBuildSaberSlash = PresetUtils.makePresetBuild('Saber Slash', {
	gear: DefaultGearSaber,
	talents: DefaultTalentsSaber,
	rotation: DefaultAPLSaber,
	encounter: PresetUtils.makePresetEncounter(
		'Saber Slash',
		'https://wowsims.github.io/sod/rogue/#eJztVV1oHFUUnnNndnb2xJneXFOc3LSyRlvHMQ13ZrM1f3a3aauhSF3Rh4AIVRJ/StBoQPxBqGva1EogFhUMIhEEg2BpY30JKCVPzVsiGAq2WPGtiPbFpm1EvXNnZ9M0iQrWnwfPw+y53/k/99yzSC2SJR4pkr2wn5RGCXRtcGGCwDECRZgFOAdQIk/CKYBuGCfANF6maJb6H32x7zmacjuzpj8BWEMX33bcmQ7v0HGOBl084mCaXpbQ0TrJTL3huHOmVPpFIt853tkTrjwsyMNZx3tLWtTQ8ojjfuN4cxfr0aIHRhzvxKV6CQ+NKK9fK4NXY51DlyOdLxZsb3ioQXKzP9re5681yLDnDzrye3g4Cj570HG/j5gr83bEGPS9BTvMWjoFD1plbd3khfR+MEZBHwf9JMCXoJ2GVPu8geeI9j/9rcQuwL+dwj9N54nWVdsKExCNHpSgB/bCAOzeKEQ+F+YCEeTzgWjZIsJcTuRDISkvgse0fUMZV8daTO9LoW4t2o0mGkjEIK5Hi0veOjMOjRlMY2pL0JyvwvSHn+wqvHUJvjKawKI5GGT38l3YyduR+o5lMPKSxpX7QOBmdluC7k7Q/LV6LYLLrOjPHczgpF1jd/I7MMtvRnsSkBGakulFr9turMHMx4YZH5XJOw+yw8CHAV9hLyuvOtPLoEUyhBbcxjtxg88lum4ablCuTEvuhmxvrJFLpFOrSUO83d+0JCVVKZ0ktgWqolhRdVDKWAdvQ8GbsYn5iJMgG2RdbPGICgNrJRH9WJ/1swf4/biL7cCNfoNlMjoNtlJOW2d6sr2ekRR1bZkr2tkqVD5y87ES34M7+HYssLuVlrmsPTf566UnnKrmGQuCimr1xkLBEw1VYTNrii9jEqJA0SKNLGBZUpkxJZQrVxmfdDxgj/CHcSfvwiLbtopFVDesVffVDlELb8S0atu3T1fmV8JskD+LT/A+rGNsFfeb2K0repDDdb5dGc5KU5KJXd6qSkeSFt8lqkE/BX4M8EPgHwDuZF1rXM/VtvlkHFadutZqovrvJxoyUbn5JbVQuTbXcq1u4v0ez2Ld/B71YpNq/yjjFYGqI9/E/T/xWKOlMfcmDJyCSwTGQCsTWYvV2zitZ+K9tqdYGzN1xfqxdyOaKbTFyGzhlscVzRfCI7p1umwz676+5/v6s1uFW8jq/n/sj9Z7vXhd/Gz+5C/5ofzX65OHpEB9Bz4qtj31zIUDMw99VdhekRR74DfhTw0j',
	),
});

export const PresetBuildSaberSlashIEA = PresetUtils.makePresetBuild('Saber Slash IEA', {
	gear: DefaultGearSaber,
	talents: P5TalentMutilateSaberslashCTTCIEA,
	rotation: ROTATION_PRESET_SABER_IEA_P5,
	encounter: PresetUtils.makePresetEncounter(
		'Saber Slash IEA',
		'https://wowsims.github.io/sod/rogue/#eJztVl1oHFUUnnNndvbmxJlOri1Obqysa1vHsQ2zs7sxf3TXaDX9TWoF8ySxJGJr0GBB/EGosW1qJRCLFQ1FIooGwVLXlkLAUvLUvCWCoWCLFd+KaF9s2qaod+7sbEx2o4L158H7cPfe8/Odc74597BoUZIgDsmTbthLOocJtHEbxggcJ5CHKYCLAE/BWYB2GCXAFH6yBvXOvide7H3OitmtCd0dA6y25o6Y9mSLc/Azjpo1d9jEuHVNiI4tF4fxN0x7WhdGPwvJd6Zz4YQtLrPicsF03hIe1dbAkGl/YzrTV2qRWvuHTOfE1Voh3jckUb+WDq+GNgevBTanZw1ncF+dOE39aDhfvFYnwl46YIr90GAQfOqAaX8fHK7PGMFBs47OGn6CqhY40CgqaycvxPeCNgzqKKhnAL4E5RzEmmc0vEiU/9ffuthl+LdT+KfXJaK01TTCGAStB53QBd3QD5tWel427fsp38tmU15mnednhMT3fLF7qZ3K7lNoq1iD8d0xVOmckdRRQ+LtwRVIuTjT86OQrMI4xtal6rMlsfXDT0ZJ3DAvvj4cib361B72MN+ArbwZLdekGiMvKVzCpzxcw1ZF0k2RNLvYLuNxkZV1o4VpnDQr7F5+Dyb4HWgUABmxYiK94HUbyWqs+kTTw6t0eXsHOwR8EPAV9rJEVZk6AEqgQ8jget6Kt7tcSJdNwC0SSqdiNiR6Qot0pB2vpPXxbnf1vJaUtFaBGBRkRaGhZFDoWAtvQo/X41rmIhZAEESvZBwiw8BSSQQ/9GQf2847cAN7AFe6dVRn1gQY0jhOz3clehwtKmpxmWV0NnoyHzH5WAffim0sj3VurfAxx6FaQsbo6cecWMRB4A4LAJe5hsigVF+ax4surJs/jjv4duxgW6WfvoDv29wVIgyOlwoPFamiaakFfA/NwsJcIg+2mW/E9aw1/NoFCCoJJnV5klUjUhnM9EpgZ0wH2NN8Fz7CO3Eb21IBISAaliL69wL4t2JcfrVvnyk+H2HEPgL+AeA7wI8ArmLJCgEXwWCWpctYTEv2tXn2M2VuDSxTgfwix1EX3FfOcZTo58CPA34I/H3AB1nbEh31W6xs1MEVH0ojrmZ3laFUKMRnXrFZ5818Ca0vBS2/5XtdDmXt/CE5ZKKR8kcZlwUqvdK13P0T8yWYc9NvQv9ZuEpgBJQBImqhPckJtSocxdvyNeFheb525N1gTeaaQslU7s4n5ZrJ+YdVem7AYHRL7/O9fYkGz84lVPc/9t/AeT1/U3DWfPqXcCz+y83JQ6yU3Ps/zjftevby/slHv8rdX9Tku+BXQNYmBQ==',
	),
});

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
