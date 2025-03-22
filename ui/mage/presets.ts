import { Phase } from '../core/constants/other';
import * as PresetUtils from '../core/preset_utils';
import { Spec } from '../core/proto/common';
import { SavedTalents } from '../core/proto/ui';
// APLs
import Phase1APLArcaneJSON from './apls/p1_arcane.apl.json';
import Phase1APLFireJSON from './apls/p1_fire.apl.json';
import Phase2APLArcaneJSON from './apls/p2_arcane.apl.json';
import Phase2APLFireJSON from './apls/p2_fire.apl.json';
// import Phase3APLArcaneJSON from './apls/p3_arcane.apl.json';
import Phase3APLFireJSON from './apls/p3_fire.apl.json';
import Phase3APLFrostJSON from './apls/p3_frost.apl.json';
import Phase4APLArcaneJSON from './apls/p4_arcane.apl.json';
import Phase4APLFireJSON from './apls/p4_fire.apl.json';
import Phase4APLFrostJSON from './apls/p4_frost.apl.json';
import Phase5APLFireJSON from './apls/p5_fire.apl.json';
import Phase5APLSpellFrostJSON from './apls/p5_spellfrost.apl.json';
// Builds
import Phase6BuildFireJSON from './builds/p6_fire.build.json';
import Phase6BuildFrostJSON from './builds/p6_frost.build.json';
import Phase7BuildFireJSON from './builds/p7_fire.build.json';
import Phase7BuildFrostJSON from './builds/p7_frost.build.json';
// Gear
import Phase1GearFireJSON from './gear_sets/p1_fire.gear.json';
import Phase1GearJSON from './gear_sets/p1_generic.gear.json';
import Phase2GearArcaneJSON from './gear_sets/p2_arcane.gear.json';
import Phase2GearFireJSON from './gear_sets/p2_fire.gear.json';
import Phase2GearFrostJSON from './gear_sets/p2_frost.gear.json';
import Phase3GearFireJSON from './gear_sets/p3_fire.gear.json';
import Phase3GearFrostFFBJSON from './gear_sets/p3_frost_ffb.gear.json';
import Phase4GearArcaneJSON from './gear_sets/p4_arcane.gear.json';
import Phase4GearFireJSON from './gear_sets/p4_fire.gear.json';
import Phase4GearFrostJSON from './gear_sets/p4_frost.gear.json';
import Phase5GearArcaneJSON from './gear_sets/p5_arcane.gear.json';
import Phase5GearFireJSON from './gear_sets/p5_fire.gear.json';
import Phase5GearFrostJSON from './gear_sets/p5_frost.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Build Presets
///////////////////////////////////////////////////////////////////////////

export const PresetBuildFirePhase6 = PresetUtils.makePresetBuildFromJSON('P6 Fire', Spec.SpecMage, Phase6BuildFireJSON);
export const PresetBuildFrostPhase6 = PresetUtils.makePresetBuildFromJSON('P6 Frost', Spec.SpecMage, Phase6BuildFrostJSON);

export const PresetBuildFirePhase7 = PresetUtils.makePresetBuildFromJSON('P7 Fire', Spec.SpecMage, Phase7BuildFireJSON);
export const PresetBuildFrostPhase7 = PresetUtils.makePresetBuildFromJSON('P7 Frost', Spec.SpecMage, Phase7BuildFrostJSON);

export const DefaultBuild = PresetBuildFirePhase7;

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArcanePhase1 = PresetUtils.makePresetGear('P1 Arcane', Phase1GearJSON, {
	customCondition: player => player.getLevel() === 25,
});
export const GearFirePhase1 = PresetUtils.makePresetGear('P1 Fire', Phase1GearFireJSON, {
	customCondition: player => player.getLevel() === 25,
});
export const GearFrostPhase1 = PresetUtils.makePresetGear('P1 Frost', Phase1GearJSON, {
	customCondition: player => player.getLevel() === 25,
});

export const GearArcanePhase2 = PresetUtils.makePresetGear('P2 Arcane', Phase2GearArcaneJSON, {
	customCondition: player => player.getLevel() === 40,
});
export const GearFirePhase2 = PresetUtils.makePresetGear('P2 Fire', Phase2GearFireJSON, {
	customCondition: player => player.getLevel() === 40,
});
export const GearFrostPhase2 = PresetUtils.makePresetGear('P2 Frost', Phase2GearFrostJSON, {
	customCondition: player => player.getLevel() === 40,
});

// No new Phase 3 Arcane presets at the moment
export const GearArcanePhase3 = GearArcanePhase2;
export const GearFirePhase3 = PresetUtils.makePresetGear('P3 Fire', Phase3GearFireJSON, {
	customCondition: player => player.getLevel() === 50,
});
export const GearFrostPhase3 = PresetUtils.makePresetGear('P3 Frost', Phase3GearFrostFFBJSON, {
	customCondition: player => player.getLevel() === 50,
});

// No new Phase 4 Arcane presets at the moment
export const GearArcanePhase4 = PresetUtils.makePresetGear('P4 Arcane', Phase4GearArcaneJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const GearFirePhase4 = PresetUtils.makePresetGear('P4 Fire', Phase4GearFireJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const GearFrostPhase4 = PresetUtils.makePresetGear('P4 Frost', Phase4GearFrostJSON, {
	customCondition: player => player.getLevel() === 60,
});

export const GearArcanePhase5 = PresetUtils.makePresetGear('P5 Arcane', Phase5GearArcaneJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const GearFirePhase5 = PresetUtils.makePresetGear('P5 Fire', Phase5GearFireJSON, {
	customCondition: player => player.getLevel() === 60,
});
export const GearFrostPhase5 = PresetUtils.makePresetGear('P5 Frost', Phase5GearFrostJSON, {
	customCondition: player => player.getLevel() === 60,
});

export const GearPresets = {
	[Phase.Phase1]: [GearArcanePhase1, GearFirePhase1, GearFrostPhase1],
	[Phase.Phase2]: [GearArcanePhase2, GearFirePhase2, GearFrostPhase2],
	[Phase.Phase3]: [GearArcanePhase3, GearFirePhase3, GearFrostPhase3],
	[Phase.Phase4]: [GearArcanePhase4, GearFirePhase4, GearFrostPhase4],
	[Phase.Phase5]: [GearArcanePhase5, GearFirePhase5, GearFrostPhase5],
	[Phase.Phase6]: [PresetBuildFirePhase6.gear!, PresetBuildFrostPhase6.gear!],
	[Phase.Phase7]: [PresetBuildFirePhase7.gear!, PresetBuildFrostPhase7.gear!],
};

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLArcanePhase1 = PresetUtils.makePresetAPLRotation('P1 Arcane', Phase1APLArcaneJSON, {
	customCondition: player => player.getLevel() === 25,
});
export const APLFirePhase1 = PresetUtils.makePresetAPLRotation('P1 Fire', Phase1APLFireJSON, {
	customCondition: player => player.getLevel() === 25,
});

export const APLArcanePhase2 = PresetUtils.makePresetAPLRotation('P2 Arcane', Phase2APLArcaneJSON, {
	customCondition: player => player.getLevel() === 40,
});
export const APLFirePhase2 = PresetUtils.makePresetAPLRotation('P2 Fire', Phase2APLFireJSON, {
	customCondition: player => player.getLevel() === 40,
});

// No new Phase 3 Arcane presets at the moment
export const APLArcanePhase3 = APLArcanePhase2;
export const APLFirePhase3 = PresetUtils.makePresetAPLRotation('P3 Fire', Phase3APLFireJSON, {
	customCondition: player => player.getLevel() === 50,
});
export const APLFrostPhase3 = PresetUtils.makePresetAPLRotation('P3 Frost', Phase3APLFrostJSON, {
	customCondition: player => player.getLevel() === 50,
});

export const APLArcanePhase4 = PresetUtils.makePresetAPLRotation('P4 Arcane', Phase4APLArcaneJSON, {
	customCondition: player => player.getLevel() >= 60,
});
export const APLFirePhase4 = PresetUtils.makePresetAPLRotation('P4 Fire', Phase4APLFireJSON, {
	customCondition: player => player.getLevel() >= 60,
});
export const APLFrostPhase4 = PresetUtils.makePresetAPLRotation('P4 Frost', Phase4APLFrostJSON, {
	customCondition: player => player.getLevel() >= 60,
});

export const APLFirePhase5 = PresetUtils.makePresetAPLRotation('P5 Fire', Phase5APLFireJSON, {
	customCondition: player => player.getLevel() >= 60,
});
export const APLSpellfrostPhase5 = PresetUtils.makePresetAPLRotation('P5 Frost', Phase5APLSpellFrostJSON, {
	customCondition: player => player.getLevel() >= 60,
});

export const APLPresets = {
	[Phase.Phase1]: [APLArcanePhase1, APLFirePhase1, APLFirePhase1],
	[Phase.Phase2]: [APLArcanePhase2, APLFirePhase2, APLFirePhase2],
	[Phase.Phase3]: [APLArcanePhase3, APLFirePhase3, APLFrostPhase3],
	[Phase.Phase4]: [APLArcanePhase4, APLFirePhase4, APLFrostPhase4],
	[Phase.Phase5]: [APLFirePhase5, APLSpellfrostPhase5],
	[Phase.Phase6]: [PresetBuildFirePhase6.rotation!, PresetBuildFrostPhase6.rotation!],
	[Phase.Phase7]: [PresetBuildFirePhase7.rotation!, PresetBuildFrostPhase7.rotation!],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
		2: APLPresets[Phase.Phase1][2],
	},
	40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][1],
		// Normally frost but frost is unfortunately just too bad to warrant including for now
		2: APLPresets[Phase.Phase2][2],
		// Frostfire
		3: APLPresets[Phase.Phase2][2],
	},
	50: {
		0: APLPresets[Phase.Phase3][0],
		1: APLPresets[Phase.Phase3][1],
		2: APLPresets[Phase.Phase3][2],
	},
	60: {
		0: APLPresets[Phase.Phase6][1],
		1: APLPresets[Phase.Phase7][0],
		2: APLPresets[Phase.Phase7][1],
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// P1
export const TalentsArcanePhase1 = PresetUtils.makePresetTalents('25 Arcane', SavedTalents.create({ talentsString: '22500502' }), {
	customCondition: player => player.getLevel() === 25,
});
export const TalentsFirePhase1 = PresetUtils.makePresetTalents('25 Fire', SavedTalents.create({ talentsString: '-5050020121' }), {
	customCondition: player => player.getLevel() === 25,
});

// P2
export const TalentsArcanePhase2 = PresetUtils.makePresetTalents('40 Arcane', SavedTalents.create({ talentsString: '2250050310031531' }), {
	customCondition: player => player.getLevel() === 40,
});
export const TalentsFirePhase2 = PresetUtils.makePresetTalents('40 Fire', SavedTalents.create({ talentsString: '-5050020123033151' }), {
	customCondition: player => player.getLevel() === 40,
});

// P3
// No new Phase 3 Arcane presets at the moment
export const TalentsArcanePhase3 = TalentsArcanePhase2;
export const TalentsFirePhase3 = PresetUtils.makePresetTalents('50 Fire', SavedTalents.create({ talentsString: '-0550020123033151-2035' }), {
	customCondition: player => player.getLevel() === 50,
});
export const TalentsFrostPhase3 = PresetUtils.makePresetTalents('50 Frost', SavedTalents.create({ talentsString: '-055-20350203100351051' }), {
	customCondition: player => player.getLevel() === 50,
});

// P4
export const TalentsArcanePhase4_5 = PresetUtils.makePresetTalents('P4/5 Arcane', SavedTalents.create({ talentsString: '0550050210031531-054-203500001' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsFirePhase4_5 = PresetUtils.makePresetTalents('P4/5 Fire', SavedTalents.create({ talentsString: '21-5052300123033151-203500031' }), {
	customCondition: player => player.getLevel() === 60,
});
export const TalentsFrostfirePhase4 = PresetUtils.makePresetTalents('P4 Frostfire', SavedTalents.create({ talentsString: '-0550320003021-2035020310035105' }), {
	customCondition: player => player.getLevel() === 60,
});

// P5
export const TalentsFrostPhase5 = PresetUtils.makePresetTalents('P5 Spellfrost', SavedTalents.create({ talentsString: '005005001--20353203112351351' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsArcanePhase1, TalentsFirePhase1, TalentsFirePhase1],
	[Phase.Phase2]: [TalentsArcanePhase2, TalentsFirePhase2, TalentsFirePhase2],
	[Phase.Phase3]: [TalentsArcanePhase3, TalentsFirePhase3, TalentsFrostPhase3],
	[Phase.Phase4]: [TalentsArcanePhase4_5, TalentsFirePhase4_5, TalentsFrostfirePhase4],
	[Phase.Phase5]: [TalentsArcanePhase4_5, TalentsFirePhase4_5, TalentsFrostPhase5],
	[Phase.Phase6]: [PresetBuildFirePhase6.talents!, PresetBuildFrostPhase6.talents!],
	[Phase.Phase7]: [PresetBuildFirePhase7.talents!, PresetBuildFrostPhase7.talents!],
};
