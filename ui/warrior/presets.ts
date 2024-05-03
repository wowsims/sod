import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	Alcohol,
	ArmorElixir,
	Consumes,
	Debuffs,
	EnchantedSigil,
	HealthElixir,
	IndividualBuffs,
	Profession,
	RaidBuffs,
	SaygesFortune,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Warrior_Options as WarriorOptions, WarriorShout } from '../core/proto/warrior.js';
import Phase1APLArms from './apls/phase_1_arms.apl.json';
import Phase2APLArms from './apls/phase_2_arms.apl.json';
import Phase2APLFury from './apls/phase_2_fury.apl.json';
import Phase3APLArms from './apls/phase_3_arms.apl.json';
import Phase3APLFury from './apls/phase_3_fury.apl.json';

import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase1DWGear from './gear_sets/phase_1_dw.gear.json';
import Phase22HGear from './gear_sets/phase_2_2h.gear.json';
import Phase2DWGear from './gear_sets/phase_2_dw.gear.json';
import Phase32HGear from './gear_sets/phase_3_2h.gear.json';
import Phase3DWGear from './gear_sets/phase_3_dw.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArmsPhase1 = PresetUtils.makePresetGear('P1 Arms 2H', Phase1Gear, { talentTree: 0 });
export const GearArmsDWPhase1 = PresetUtils.makePresetGear('P1 Arms DW', Phase1DWGear, { talentTree: 0 });
export const GearArmsPhase2 = PresetUtils.makePresetGear('P2 2H', Phase22HGear, { talentTree: 0 });
export const GearFuryPhase1 = PresetUtils.makePresetGear('P1 DW Fury', Phase1Gear, { talentTree: 1 });
export const GearFuryPhase2 = PresetUtils.makePresetGear('P2 DW', Phase2DWGear, { talentTree: 1 });
export const GearArmsPhase3 = PresetUtils.makePresetGear('P3 2H', Phase32HGear, { talentTree: 0 });
export const GearFuryPhase3 = PresetUtils.makePresetGear('P3 DW', Phase3DWGear, { talentTree: 1 });

export const GearPresets = {
	[Phase.Phase1]: [GearArmsPhase1, GearFuryPhase1, GearArmsDWPhase1],
	[Phase.Phase2]: [GearArmsPhase2, GearFuryPhase2],
	[Phase.Phase3]: [GearArmsPhase3, GearFuryPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearPresets[Phase.Phase3][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1Arms = PresetUtils.makePresetAPLRotation('P1 Preset Arms', Phase1APLArms);
export const APLPhase2Arms = PresetUtils.makePresetAPLRotation('P2 Preset Arms', Phase2APLArms);
export const APLPhase2Fury = PresetUtils.makePresetAPLRotation('P2 Preset Fury', Phase2APLFury);
export const APLPhase3Arms = PresetUtils.makePresetAPLRotation('P3 Preset Arms', Phase3APLArms);
export const APLPhase3Fury = PresetUtils.makePresetAPLRotation('P3 Preset Fury', Phase3APLFury);

export const APLPresets = {
	[Phase.Phase1]: [APLPhase1Arms],
	[Phase.Phase2]: [APLPhase2Arms, APLPhase2Fury],
	[Phase.Phase3]: [APLPhase3Arms, APLPhase3Fury],
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
		1: APLPresets[Phase.Phase2][1],
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

export const TalentsPhase1 = {
	name: 'P1',
	data: SavedTalents.create({
		talentsString: '303220203-01',
	}),
};

export const TalentsPhase2Fury = {
	name: 'P2 Fury',
	data: SavedTalents.create({
		talentsString: '-05050005405010051',
	}),
};

export const TalentsPhase2Arms = {
	name: 'P2 Arms',
	data: SavedTalents.create({
		talentsString: '303050213525100001',
	}),
};

export const TalentsPhase3Arms = {
	name: 'P3 Arms',
	data: SavedTalents.create({
		talentsString: '303050213520105001-0505'
	}),
}

export const TalentPhase3Fury = {
	name: 'P3 Fury',
	data: SavedTalents.create({
		talentsString: '303040003-0505000540501003'
	}),
}


export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2Arms, TalentsPhase2Fury],
	[Phase.Phase3]: [TalentsPhase3Arms, TalentPhase3Fury],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultTalentsFury = TalentPresets[Phase.Phase2][1];
export const DefaultTalentsArms = TalentPresets[Phase.Phase3][0];

export const DefaultTalents = DefaultTalentsArms;

///////////////////////////////////////////////////////////////////////////
//                                 Options Presets
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	shout: WarriorShout.WarriorShoutBattle,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfAgility,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.DenseSharpeningStone,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	armorElixir: ArmorElixir.ElixirOfSuperiorDefense,
	healthElixir: HealthElixir.HealthElixirUnknown,
	alcohol: Alcohol.AlcoholUnknown,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	devotionAura: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	stoneskinTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectRegular,
	sparkOfInspiration: true,
	saygesFortune: SaygesFortune.SaygesDamage,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	mangle: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Enchanting,
	profession2: Profession.Leatherworking,
};
