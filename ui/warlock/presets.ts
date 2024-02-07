import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Debuffs,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
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

import DestructionGear from './gear_sets/destruction.gear.json';
import DestructionAPL from './apls/destruction.apl.json';

export const GearAfflictionPhase1 = PresetUtils.makePresetGear('Affliction', DestructionGear);
export const GearDemonologyPhase1 = PresetUtils.makePresetGear('Demonology', DestructionGear);
export const GearDestructionPhase1 = PresetUtils.makePresetGear('Destruction', DestructionGear);

export const GearPresets = {
  [Phase.Phase1]: [
    GearAfflictionPhase1,
		GearDemonologyPhase1,
		GearDestructionPhase1,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase1][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const RotationAfflictionPhase1 = PresetUtils.makePresetAPLRotation('Affliction', DestructionAPL);
export const RotationDemonologyPhase1 = PresetUtils.makePresetAPLRotation('Demonology', DestructionAPL);
export const RotationDestructionPhase1 = PresetUtils.makePresetAPLRotation('Destruction', DestructionAPL);

export const APLPresets = {
  [Phase.Phase1]: [
    RotationAfflictionPhase1,
		RotationDemonologyPhase1,
		RotationDestructionPhase1,
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

export const DestroTalents = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '-03-0550201',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [
    DestroTalents,
  ],
  [Phase.Phase2]: [
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsAffliction = TalentPresets[Phase.Phase1][0];
export const DefaultTalentsDemonology	= TalentPresets[Phase.Phase1][0];
export const DefaultTalentsDestruction = TalentPresets[Phase.Phase1][0];

export const DefaultTalents = DefaultTalentsDestruction;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Imp,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodSmokedSagefish,
	defaultPotion: Potions.ManaPotion,
	mainHandImbue: WeaponImbue.BlackfathomManaOil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
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
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
