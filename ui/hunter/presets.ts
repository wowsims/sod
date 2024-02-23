import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import {
	AgilityElixir,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Food,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetType as PetType,
	Hunter_Options_QuiverBonus,
} from '../core/proto/hunter.js';

import * as PresetUtils from '../core/preset_utils.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

import Phase1Gear from './gear_sets/phase1.gear.json';
import Phase2GearRanged from './gear_sets/p2_ranged.gear.json';
import Phase2GearMelee from './gear_sets/p2_melee.gear.json';

export const GearBeastMasteryPhase1 = PresetUtils.makePresetGear('P1 Beast Mastery', Phase1Gear, { talentTree: 0 })
export const GearMarksmanPhase1 = PresetUtils.makePresetGear('P1 Marksmanship', Phase1Gear, { talentTree: 1 })
export const GearSurvivalPhase1 = PresetUtils.makePresetGear('P1 Survival', Phase1Gear, { talentTree: 2 })

export const GearRangedPhase2 = PresetUtils.makePresetGear('P2 Ranged', Phase2GearRanged)
export const GearMeleePhase2 = PresetUtils.makePresetGear('P2 Melee', Phase2GearMelee)

export const GearPresets = {
  	[Phase.Phase1]: [
    	GearBeastMasteryPhase1,
			GearMarksmanPhase1,
			GearSurvivalPhase1,
  	],
  	[Phase.Phase2]: [
			GearRangedPhase2,
			GearMeleePhase2,
  	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearMeleePhase2;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import MeleeWeaveP1 from './apls/p1_weave.apl.json';
import MeleeP2 from './apls/p2_melee.apl.json';
import RangedP2 from './apls/p2_ranged.apl.json';

export const APLMeleeWeavePhase1 = PresetUtils.makePresetAPLRotation('Melee Weave P1', MeleeWeaveP1);

export const APLMeleePhase2 = PresetUtils.makePresetAPLRotation('Melee P2', MeleeP2);
export const APLRangedPhase2 = PresetUtils.makePresetAPLRotation('Ranged P2', RangedP2);

export const APLPresets = {
  	[Phase.Phase1]: [
    	APLMeleeWeavePhase1,
  	],
  	[Phase.Phase2]: [
			APLRangedPhase2,
			APLMeleePhase2,
  	],
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
		2: APLPresets[Phase.Phase1][0],
	},
  40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][0],
		2: APLPresets[Phase.Phase2][1],
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsBeastMasteryPhase1 = {
	name: 'P1 Beast Mastery',
	data: SavedTalents.create({
		talentsString: '53000200501',
	}),
};

export const TalentsMarksmanPhase1 = {
	name: 'P1 Marksmanship',
	data: SavedTalents.create({
		talentsString: '-050515',
	}),
};

export const TalentsSurvivalPhase1 = {
	name: 'P1 Survival',
	data: SavedTalents.create({
		talentsString: '--33502001101',
	}),
};

export const TalentsBeastMasteryPhase2 = {
	name: 'P2 Beast Mastery',
	data: SavedTalents.create({
		talentsString: '5300021150501251',
	}),
};

export const TalentsMarksmanPhase2 = {
	name: 'P2 Marksmanship',
	data: SavedTalents.create({
		talentsString: '5-0555100150301',
	}),
};

export const TalentsSurvivalPhase2 = {
	name: 'P2 Survival',
	data: SavedTalents.create({
		talentsString: '--335020051030315',
	}),
};

export const TalentPresets = {
	[Phase.Phase1]: [
		TalentsBeastMasteryPhase1,
		TalentsMarksmanPhase1,
		TalentsSurvivalPhase1,
	],
 	[Phase.Phase2]: [
		TalentsBeastMasteryPhase2,
		TalentsMarksmanPhase2,
		TalentsSurvivalPhase2,
  ],
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalentsBeastMastery = TalentPresets[CURRENT_PHASE][0];
export const DefaultTalentsMarksman 		= TalentPresets[CURRENT_PHASE][1];
export const DefaultTalentsSurvival 		= TalentPresets[CURRENT_PHASE][2];

export const DefaultTalents = TalentsSurvivalPhase2;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.JaggedArrow,
	petType: PetType.WindSerpent,
	petTalents: {},
	petUptime: 1,
	quiverBonus: Hunter_Options_QuiverBonus.Speed14,
	petAttackSpeed: 2.0,
});

export const BMDefaultOptions = HunterOptions.create({
	ammo: Ammo.JaggedArrow,
	petType: PetType.Cat,
	petTalents: {},
	petUptime: 1,
	quiverBonus: Hunter_Options_QuiverBonus.Speed14,
	petAttackSpeed: 2.0,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfAgility,
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.SolidWeightstone,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
});

export const DefaultRaidBuffs = RaidBuffs.create({
  arcaneBrilliance: true,
	bloodPact: TristateEffect.TristateEffectImproved,
	aspectOfTheLion: true,
  battleShout: TristateEffect.TristateEffectImproved,
  divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	trueshotAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
  blessingOfMight: TristateEffect.TristateEffectImproved,
  blessingOfWisdom: TristateEffect.TristateEffectRegular,
  sparkOfInspiration: true,
  saygesFortune: SaygesFortune.SaygesDamage
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	huntersMark: TristateEffect.TristateEffectImproved,
	mangle: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	distanceFromTarget: 15,
  profession1: Profession.Enchanting,
  profession2: Profession.Leatherworking,
}
