import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import { Spec } from '../core/proto/common';
import {
	AgilityElixir,
	Alcohol,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Flask,
	Food,
	HealthElixir,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	Race,
	RaidBuffs,
	SapperExplosive,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import {
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
	Hunter_Options_QuiverBonus,
} from '../core/proto/hunter.js';
import { SavedTalents } from '../core/proto/ui.js';
import MeleeWeaveP1 from './apls/p1_weave.apl.json';
import MeleeP2 from './apls/p2_melee.apl.json';
import RangedBmP2 from './apls/p2_ranged_bm.apl.json';
import RangedMmP2 from './apls/p2_ranged_mm.apl.json';
import MeleeBmP3 from './apls/p3_melee_bm.apl.json';
import RangedMmP3 from './apls/p3_ranged_mm.apl.json';
import RangedP4 from './apls/p4_ranged.apl.json';
import WeaveP4 from './apls/p4_weave.apl.json';
import Phase5AplMeleeBm from './apls/p5_melee_bm.apl.json';
import Phase5AplMeleeSv from './apls/p5_melee_sv.apl.json';
import Phase5AplRanged from './apls/p5_ranged.apl.json';
import Phase5AplWeave from './apls/p5_weave.apl.json';
import Phase6AplMeleeBm from './apls/p6_melee_bm.apl.json';
import Phase6AplMeleeSv from './apls/p6_melee_sv.apl.json';
import Phase6AplRangedDraconic from './apls/p6_ranged_draconic.apl.json';
import Phase6AplRangedKillshot from './apls/p6_ranged_killshot.apl.json';
import Phase6AplWeave from './apls/p6_weave.apl.json';
import Phase7AplMelee2h from './apls/p7_melee_2h.apl.json';
import Phase7AplMeleeDw from './apls/p7_melee_dw.apl.json';
import Phase7AplRangedKillshot from './apls/p7_ranged_killshot.apl.json';
import Phase7AplRangedLnL from './apls/p7_ranged_lnl.apl.json';
import Phase7AplWeave from './apls/p7_weave.apl.json';
import Phase7NaxxMelee2H from './builds/p7_naxx_melee_2h.build.json';
// Builds
import Phase7NaxxMeleeDW from './builds/p7_naxx_melee_dw.build.json';
import Phase7NaxxRanged from './builds/p7_naxx_ranged.build.json';
import Phase7NaxxWeave from './builds/p7_naxx_weave.build.json';
import Phase8Melee2H from './builds/p8_melee_2h.build.json';
import Phase8MeleeDW from './builds/p8_melee_dw.build.json';
import Phase8Ranged from './builds/p8_ranged.build.json';
import Phase8Weave from './builds/p8_weave.build.json';
import Phase2GearMelee from './gear_sets/p2_melee.gear.json';
import Phase2GearRangedBm from './gear_sets/p2_ranged_bm.gear.json';
import Phase2GearRangedMm from './gear_sets/p2_ranged_mm.gear.json';
import Phase3GearMeleeBm from './gear_sets/p3_melee_bm.gear.json';
import Phase3GearRangedMm from './gear_sets/p3_ranged_mm.gear.json';
import Phase4GearRanged from './gear_sets/p4_ranged.gear.json';
import Phase4GearWeave from './gear_sets/p4_weave.gear.json';
import Phase5GearMeleeBm from './gear_sets/p5_melee_bm.gear.json';
import Phase5GearMeleeSv from './gear_sets/p5_melee_sv.gear.json';
import Phase5GearRangedSv from './gear_sets/p5_ranged_sv.gear.json';
import Phase5GearWeave from './gear_sets/p5_weave.gear.json';
import Phase6GearMeleeBm from './gear_sets/p6_melee_bm.gear.json';
import Phase6GearMeleeSv from './gear_sets/p6_melee_sv.gear.json';
import Phase6GearRangedDraconic from './gear_sets/p6_ranged_draconic.gear.json';
import Phase6GearRangedKillshot from './gear_sets/p6_ranged_killshot.gear.json';
import Phase6GearWeave from './gear_sets/p6_weave.gear.json';
import Phase7GearMelee2h from './gear_sets/p7_melee_2h.gear.json';
import Phase7GearMeleeDw from './gear_sets/p7_melee_dw.gear.json';
import Phase7GearRangedKillshot from './gear_sets/p7_ranged_killshot.gear.json';
import Phase7GearWeave from './gear_sets/p7_weave.gear.json';
import Phase1Gear from './gear_sets/phase1.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Build Presets
///////////////////////////////////////////////////////////////////////////

export const PresetBuildPhase7NaxxMeleeDW = PresetUtils.makePresetBuildFromJSON('P7 Naxx DW', Spec.SpecHunter, Phase7NaxxMeleeDW);
export const PresetBuildPhase7NaxxMelee2H = PresetUtils.makePresetBuildFromJSON('P7 Naxx 2H', Spec.SpecHunter, Phase7NaxxMelee2H);
export const PresetBuildPhase7NaxxRanged = PresetUtils.makePresetBuildFromJSON('P7 Naxx Ranged', Spec.SpecHunter, Phase7NaxxRanged);
export const PresetBuildPhase7NaxxWeave = PresetUtils.makePresetBuildFromJSON('P7 Naxx Weave', Spec.SpecHunter, Phase7NaxxWeave);
export const PresetBuildPhase8NaxxMeleeDW = PresetUtils.makePresetBuildFromJSON('P8 DW', Spec.SpecHunter, Phase8MeleeDW);
export const PresetBuildPhase8NaxxMelee2H = PresetUtils.makePresetBuildFromJSON('P8 2H', Spec.SpecHunter, Phase8Melee2H);
export const PresetBuildPhase8NaxxRanged = PresetUtils.makePresetBuildFromJSON('P8 Ranged', Spec.SpecHunter, Phase8Ranged);
export const PresetBuildPhase8NaxxWeave = PresetUtils.makePresetBuildFromJSON('P8 Weave', Spec.SpecHunter, Phase8Weave);

export const DefaultBuild = PresetBuildPhase7NaxxRanged;
export const DefaultWeights = PresetBuildPhase7NaxxRanged.epWeights!;

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBeastMasteryPhase1 = PresetUtils.makePresetGear('P1 Beast Mastery', Phase1Gear, { customCondition: player => player.getLevel() === 25 });
export const GearMarksmanPhase1 = PresetUtils.makePresetGear('P1 Marksmanship', Phase1Gear, { customCondition: player => player.getLevel() === 25 });
export const GearSurvivalPhase1 = PresetUtils.makePresetGear('P1 Survival', Phase1Gear, { customCondition: player => player.getLevel() === 25 });

export const GearRangedBmPhase2 = PresetUtils.makePresetGear('P2 Ranged BM', Phase2GearRangedBm, { customCondition: player => player.getLevel() === 40 });
export const GearRangedMmPhase2 = PresetUtils.makePresetGear('P2 Ranged MM', Phase2GearRangedMm, { customCondition: player => player.getLevel() === 40 });
export const GearMeleePhase2 = PresetUtils.makePresetGear('P2 Melee', Phase2GearMelee, { customCondition: player => player.getLevel() === 40 });

export const GearMeleeBmPhase3 = PresetUtils.makePresetGear('P3 Melee BM', Phase3GearMeleeBm, { customCondition: player => player.getLevel() === 50 });
export const GearRangedMmPhase3 = PresetUtils.makePresetGear('P3 Ranged MM', Phase3GearRangedMm, { customCondition: player => player.getLevel() === 50 });

export const GearWeavePhase4 = PresetUtils.makePresetGear('P4 Weave', Phase4GearWeave, { customCondition: player => player.getLevel() === 60 });
export const GearRangedSVPhase4 = PresetUtils.makePresetGear('P4 Ranged', Phase4GearRanged, { customCondition: player => player.getLevel() === 60 });

export const GearWeavePhase5 = PresetUtils.makePresetGear('P5 Weave', Phase5GearWeave, { customCondition: player => player.getLevel() === 60 });
export const GearRangedMMPhase5 = PresetUtils.makePresetGear('P5 Ranged MM', Phase5GearRangedSv, { customCondition: player => player.getLevel() === 60 });
export const GearRangedSVPhase5 = PresetUtils.makePresetGear('P5 Ranged SV', Phase5GearRangedSv, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeBMPhase5 = PresetUtils.makePresetGear('P5 Melee BM', Phase5GearMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeSVPhase5 = PresetUtils.makePresetGear('P5 Melee SV', Phase5GearMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const GearWeavePhase6 = PresetUtils.makePresetGear('P6 Weave', Phase6GearWeave, { customCondition: player => player.getLevel() === 60 });
export const GearRangedDraconicPhase6 = PresetUtils.makePresetGear('P6 Ranged Draconic', Phase6GearRangedDraconic, { customCondition: player => player.getLevel() === 60 });
export const GearRangedKillshotPhase6 = PresetUtils.makePresetGear('P6 Ranged Killshot', Phase6GearRangedKillshot, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeBMPhase6 = PresetUtils.makePresetGear('P6 Melee BM', Phase6GearMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeSVPhase6 = PresetUtils.makePresetGear('P6 Melee SV', Phase6GearMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const GearWeavePhase7 = PresetUtils.makePresetGear('P7 Weave', Phase7GearWeave, { customCondition: player => player.getLevel() === 60 });
export const GearRangedKillshotPhase7 = PresetUtils.makePresetGear('P7 Ranged Killshot', Phase7GearRangedKillshot, { customCondition: player => player.getLevel() === 60 });
export const GearMelee2HPhase7 = PresetUtils.makePresetGear('P7 2H', Phase7GearMelee2h, { customCondition: player => player.getLevel() === 60 });
export const GearMeleeDWPhase7 = PresetUtils.makePresetGear('P7 DW', Phase7GearMeleeDw, { customCondition: player => player.getLevel() === 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearBeastMasteryPhase1, GearMarksmanPhase1, GearSurvivalPhase1],
	[Phase.Phase2]: [GearRangedBmPhase2, GearRangedMmPhase2, GearMeleePhase2],
	[Phase.Phase3]: [GearRangedMmPhase3, GearMeleeBmPhase3],
	[Phase.Phase4]: [GearWeavePhase4, GearRangedSVPhase4],
	[Phase.Phase5]: [GearWeavePhase5, GearRangedMMPhase5, GearRangedSVPhase5, GearMeleeBMPhase5, GearMeleeSVPhase5],
	[Phase.Phase6]: [GearWeavePhase6, GearRangedDraconicPhase6, GearRangedKillshotPhase6, GearMeleeBMPhase6, GearMeleeSVPhase6],
	[Phase.Phase7]: [PresetBuildPhase7NaxxWeave.gear!, PresetBuildPhase7NaxxRanged.gear!, PresetBuildPhase7NaxxMelee2H.gear!, PresetBuildPhase7NaxxMeleeDW.gear!],
	[Phase.Phase8]: [PresetBuildPhase8NaxxWeave.gear!, PresetBuildPhase8NaxxRanged.gear!, PresetBuildPhase8NaxxMelee2H.gear!, PresetBuildPhase8NaxxMeleeDW.gear!],
};

export const DefaultGear = DefaultBuild.gear!

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLMeleeWeavePhase1 = PresetUtils.makePresetAPLRotation('P1 Melee Weave', MeleeWeaveP1, { customCondition: player => player.getLevel() === 25 });

export const APLMeleePhase2 = PresetUtils.makePresetAPLRotation('P2 Melee', MeleeP2, { customCondition: player => player.getLevel() === 40 });
export const APLRangedBmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged BM', RangedBmP2, { customCondition: player => player.getLevel() === 40 });
export const APLRangedMmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged MM', RangedMmP2, { customCondition: player => player.getLevel() === 40 });

export const APLMeleeBmPhase3 = PresetUtils.makePresetAPLRotation('P3 Melee BM', MeleeBmP3, { customCondition: player => player.getLevel() === 50 });
export const APLRangedMmPhase3 = PresetUtils.makePresetAPLRotation('P3 Ranged MM', RangedMmP3, { customCondition: player => player.getLevel() === 50 });

export const APLWeavePhase4 = PresetUtils.makePresetAPLRotation('P4 Weave', WeaveP4, { customCondition: player => player.getLevel() === 60 });
export const APLRangedPhase4 = PresetUtils.makePresetAPLRotation('P4 Ranged', RangedP4, { customCondition: player => player.getLevel() === 60 });

export const APLWeavePhase5 = PresetUtils.makePresetAPLRotation('P5 Weave', Phase5AplWeave, { customCondition: player => player.getLevel() === 60 });
export const APLRanged31Phase5 = PresetUtils.makePresetAPLRotation('P5 Ranged 3-1-1', Phase5AplRanged, { customCondition: player => player.getLevel() === 60 });
export const APLRanged22Phase5 = PresetUtils.makePresetAPLRotation('P5 Ranged 2-2', RangedP4, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeBMPhase5 = PresetUtils.makePresetAPLRotation('P5 Melee BM', Phase5AplMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeSVPhase5 = PresetUtils.makePresetAPLRotation('P5 Melee SV', Phase5AplMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const APLWeavePhase6 = PresetUtils.makePresetAPLRotation('P6 Weave', Phase6AplWeave, { customCondition: player => player.getLevel() === 60 });
export const APLRangedDraconicPhase6 = PresetUtils.makePresetAPLRotation('P6 Ranged Draconic', Phase6AplRangedDraconic, { customCondition: player => player.getLevel() === 60 });
export const APLRangedKillshotPhase6 = PresetUtils.makePresetAPLRotation('P6 Ranged Killshot', Phase6AplRangedKillshot, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeBMPhase6 = PresetUtils.makePresetAPLRotation('P6 Melee BM', Phase6AplMeleeBm, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeSVPhase6 = PresetUtils.makePresetAPLRotation('P6 Melee SV', Phase6AplMeleeSv, { customCondition: player => player.getLevel() === 60 });

export const APLWeavePhase7 = PresetUtils.makePresetAPLRotation('P7 Weave', Phase7AplWeave, { customCondition: player => player.getLevel() === 60 });
export const APLRangedKillshotPhase7 = PresetUtils.makePresetAPLRotation('P7 Ranged Killshot', Phase7AplRangedKillshot, { customCondition: player => player.getLevel() === 60 });
export const APLRangedLnLPhase7 = PresetUtils.makePresetAPLRotation('P7 Ranged LnL', Phase7AplRangedLnL, { customCondition: player => player.getLevel() === 60 });
export const APLMeleeDWPhase7 = PresetUtils.makePresetAPLRotation('P7 Melee DW', Phase7AplMeleeDw, { customCondition: player => player.getLevel() === 60 });
export const APLMelee2HPhase7 = PresetUtils.makePresetAPLRotation('P7 Melee 2H', Phase7AplMelee2h, { customCondition: player => player.getLevel() === 60 });

export const APLPresets = {
	[Phase.Phase1]: [APLMeleeWeavePhase1],
	[Phase.Phase2]: [APLRangedBmPhase2, APLRangedMmPhase2, APLMeleePhase2],
	[Phase.Phase3]: [APLRangedMmPhase3, APLMeleeBmPhase3],
	[Phase.Phase4]: [APLWeavePhase4, APLRangedPhase4],
	[Phase.Phase5]: [APLWeavePhase5, APLRanged31Phase5, APLRanged22Phase5, APLMeleeBMPhase5, APLMeleeSVPhase5],
	[Phase.Phase6]: [APLWeavePhase6, APLRangedDraconicPhase6, APLRangedKillshotPhase6, APLMeleeBMPhase6, APLMeleeSVPhase6],
	[Phase.Phase7]: [PresetBuildPhase7NaxxWeave.rotation!, PresetBuildPhase7NaxxRanged.rotation!, PresetBuildPhase7NaxxMelee2H.rotation!, PresetBuildPhase7NaxxMeleeDW.rotation!],
	[Phase.Phase8]: [PresetBuildPhase8NaxxWeave.rotation!, PresetBuildPhase8NaxxRanged.rotation!, PresetBuildPhase8NaxxMelee2H.rotation!, PresetBuildPhase8NaxxMeleeDW.rotation!],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsBeastMasteryPhase1 = PresetUtils.makePresetTalents('P1 Beast Mastery', SavedTalents.create({ talentsString: '53000200501' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsMarksmanPhase1 = PresetUtils.makePresetTalents('P1 Marksmanship', SavedTalents.create({ talentsString: '-050515' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsSurvivalPhase1 = PresetUtils.makePresetTalents('P1 Survival', SavedTalents.create({ talentsString: '--33502001101' }), {
	customCondition: player => player.getLevel() === 25,
});

export const TalentsBeastMasteryPhase2 = PresetUtils.makePresetTalents('P2 Beast Mastery', SavedTalents.create({ talentsString: '5300021150501251' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsMarksmanPhase2 = PresetUtils.makePresetTalents('P2 Marksmanship', SavedTalents.create({ talentsString: '-05551001503051' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsSurvivalPhase2 = PresetUtils.makePresetTalents('P2 Survival', SavedTalents.create({ talentsString: '--335020051030315' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TalentsRangedMMPhase3 = PresetUtils.makePresetTalents('P3 Ranged MM', SavedTalents.create({ talentsString: '5-05051404503051-3' }), {
	customCondition: player => player.getLevel() === 50,
});
export const TalentsMeleeBMPhase3 = PresetUtils.makePresetTalents('P3 Melee BM', SavedTalents.create({ talentsString: '0500321150521251--33002' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TalentsWeavePhase4 = PresetUtils.makePresetTalents('60 Weave', SavedTalents.create({ talentsString: '-055500005-3305202202303051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedMMPhase4 = PresetUtils.makePresetTalents('60 Ranged MM', SavedTalents.create({ talentsString: '-05451002503051-33400023023' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedSVPhase4 = PresetUtils.makePresetTalents('60 Ranged SV', SavedTalents.create({ talentsString: '1-054510005-334000250230305' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentsWeavePhase5 = PresetUtils.makePresetTalents('P5 Weave', SavedTalents.create({ talentsString: '-055500005-3305202202303051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedMMPhase5 = PresetUtils.makePresetTalents('P5 Ranged MM', SavedTalents.create({ talentsString: '5-05451005503051-3320202' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedSVPhase5 = PresetUtils.makePresetTalents('P5 Ranged SV', SavedTalents.create({ talentsString: '1-054510005-334000250230305' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeBMPhase5 = PresetUtils.makePresetTalents('P5 Melee BM', SavedTalents.create({ talentsString: '5500020050521251-0505-33002' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeSVPhase5 = PresetUtils.makePresetTalents('P5 Melee SV', SavedTalents.create({ talentsString: '-055500005-3320202412303051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentsWeavePhase6 = PresetUtils.makePresetTalents('P6 Weave', SavedTalents.create({ talentsString: '-054510005-3305202202303051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedMMPhase6 = PresetUtils.makePresetTalents('P6 Ranged MM', SavedTalents.create({ talentsString: '5-05451005503051-3320202' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedSVPhase6 = PresetUtils.makePresetTalents('P6 Ranged SV', SavedTalents.create({ talentsString: '-054510015-334000250230305' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeBMPhase6 = PresetUtils.makePresetTalents('P6 Melee BM', SavedTalents.create({ talentsString: '5500020050501251-0505-33202' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeSVPhase6 = PresetUtils.makePresetTalents('P6 Melee SV', SavedTalents.create({ talentsString: '-055500005-3320202412303051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentsWeavePhase7 = PresetUtils.makePresetTalents('P7 Weave', SavedTalents.create({ talentsString: '-054510005-3305202202303051' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsRangedMMPhase7 = PresetUtils.makePresetTalents('P7 Ranged', SavedTalents.create({ talentsString: '5-05451005503051-3320202' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMeleeDWPhase7 = PresetUtils.makePresetTalents('P7 DW', SavedTalents.create({ talentsString: '1-052500305-332020241230305' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsMelee2HPhase7 = PresetUtils.makePresetTalents('P7 2H', SavedTalents.create({ talentsString: '-055500005-3320202412303051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsBeastMasteryPhase1, TalentsMarksmanPhase1, TalentsSurvivalPhase1],
	[Phase.Phase2]: [TalentsBeastMasteryPhase2, TalentsMarksmanPhase2, TalentsSurvivalPhase2],
	[Phase.Phase3]: [TalentsRangedMMPhase3, TalentsMeleeBMPhase3],
	[Phase.Phase4]: [TalentsWeavePhase4, TalentsRangedMMPhase4, TalentsRangedSVPhase4],
	[Phase.Phase5]: [TalentsWeavePhase5, TalentsRangedMMPhase5, TalentsRangedSVPhase5, TalentsMeleeBMPhase5, TalentsMeleeSVPhase5],
	[Phase.Phase6]: [TalentsWeavePhase6, TalentsRangedMMPhase6, TalentsRangedSVPhase6, TalentsMeleeBMPhase6, TalentsMeleeSVPhase6],
	[Phase.Phase7]: [PresetBuildPhase7NaxxWeave.talents!, PresetBuildPhase7NaxxRanged.talents!, PresetBuildPhase7NaxxMelee2H.talents!, PresetBuildPhase7NaxxMeleeDW.talents!],
	[Phase.Phase8]: [PresetBuildPhase8NaxxWeave.talents!, PresetBuildPhase8NaxxRanged.talents!, PresetBuildPhase8NaxxMelee2H.talents!, PresetBuildPhase8NaxxMeleeDW.talents!],
};

export const DefaultTalents = DefaultBuild.talents!;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const OptionsRangedLonewolf = HunterOptions.create({
	ammo: Ammo.SearingArrow,
	quiverBonus: Hunter_Options_QuiverBonus.Speed15,
	petAttackSpeed: 2.0,
	petTalents: {},
	petType: PetType.PetNone,
	petUptime: 1,
	sniperTrainingUptime: 1.0,
});

export const DefaultOptions = OptionsRangedLonewolf;

// Consumable Presets
export const RangedConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: false,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodGrilledSquid,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.EnchantedRepellent,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	offHandImbue: WeaponImbue.EnchantedRepellent,
	petAttackPowerConsumable: 1,
	petAgilityConsumable: 1,
	petStrengthConsumable: 1,
	sapperExplosive: SapperExplosive.SapperFumigator,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const DefaultConsumes = RangedConsumes;

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	commandingShout: true,
	demonicPact: 120,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectRegular,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	fervorOfTheTempleExplorer: true,
	mightOfStormwind: true,
	moldarsMoxie: true,
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
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	huntersMark: TristateEffect.TristateEffectImproved,
	improvedScorch: true,
	judgementOfWisdom: true,
	mangle: true,
	markOfChaos: true,
	occultPoison: true,
	stormstrike: true,
});

export const OtherDefaults = {
	distanceFromTarget: 12,
	profession1: Profession.Enchanting,
	profession2: Profession.Engineering,
	race: Race.RaceNightElf,
};
