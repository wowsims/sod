import { Phase } from '../core/constants/other.js';
import {
	AgilityElixir,
	Consumes,
	Debuffs,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	StrengthBuff,
	TristateEffect,
	WeaponImbue
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions as WarlockOptions,
	WarlockOptions_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';

import * as PresetUtils from '../core/preset_utils.js';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

import AfflictionTankGear from './gear_sets/affi.tank.gear.json';
import DestroTankGear from './gear_sets/destro.tank.gear.json';

export const GearAfflictionTankPhase1 = PresetUtils.makePresetGear('Affliction Tank', AfflictionTankGear);
export const GearDemologyTankPhase1 = PresetUtils.makePresetGear('Demonology Tank', AfflictionTankGear);
export const GearDestructionTankPhase1 = PresetUtils.makePresetGear('Destruction Tank', DestroTankGear);

export const GearPresets = {
  [Phase.Phase1]: [
    GearAfflictionTankPhase1,
		GearDemologyTankPhase1,
		GearDestructionTankPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];


///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import AfflictionTankAPL from './apls/affi.tank.apl.json';
import DestroTankAPL from './apls/destro.tank.apl.json';

export const APLAfflictionTankPhase1 = PresetUtils.makePresetAPLRotation('Affliction Tank', AfflictionTankAPL);
export const APLDemonologyTankPhase1 = PresetUtils.makePresetAPLRotation('Demonology Tank', DestroTankAPL);
export const APLDestructionTankPhase1 = PresetUtils.makePresetAPLRotation('Destruction Tank', DestroTankAPL);

export const APLPresets = {
  [Phase.Phase1]: [
    APLAfflictionTankPhase1,
		APLDemonologyTankPhase1,
		APLDestructionTankPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
		2: APLPresets[Phase.Phase1][2],
	},
  40: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
		2: APLPresets[Phase.Phase1][2],
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsAfflictionTankPhase1 = {
	name: 'Affliction Tank',
	data: SavedTalents.create({
		talentsString: '050025001-003',
	}),
};

export const TalentsDestructionTankPhase1 = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '-03-0550201',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
    TalentsAfflictionTankPhase1,
		TalentsDestructionTankPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsAffliction 	= TalentPresets[Phase.Phase1][0]
// export const DefaultTalentsDemonology 	= TalentPresets[Phase.Phase1][1];
export const DefaultTalentsDestruction 	= TalentPresets[Phase.Phase1][1];

export const DefaultTalents = DefaultTalentsAffliction;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Succubus,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodSmokedSagefish,
	defaultPotion: Potions.ManaPotion,
	mainHandImbue: WeaponImbue.BlackfathomManaOil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	agilityElixir: AgilityElixir.ElixirOfLesserAgility,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	divineSpirit: true,
	aspectOfTheLion: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	blessingOfMight: TristateEffect.TristateEffectImproved,
});

export const DefaultDebuffs = Debuffs.create({
	homunculi: 70, // 70% average uptime default
	faerieFire: true,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
