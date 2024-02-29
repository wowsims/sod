import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
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
import Phase2GearMelee from './gear_sets/p2_melee.gear.json';
import Phase2GearRangedBm from './gear_sets/p2_ranged_bm.gear.json';
import Phase2GearRangedMm from './gear_sets/p2_ranged_mm.gear.json';
import Phase1Gear from './gear_sets/phase1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBeastMasteryPhase1 = PresetUtils.makePresetGear('P1 Beast Mastery', Phase1Gear, { talentTree: 0 })
export const GearMarksmanPhase1 = PresetUtils.makePresetGear('P1 Marksmanship', Phase1Gear, { talentTree: 1 })
export const GearSurvivalPhase1 = PresetUtils.makePresetGear('P1 Survival', Phase1Gear, { talentTree: 2 })

export const GearRangedBmPhase2 = PresetUtils.makePresetGear('P2 Ranged BM', Phase2GearRangedBm)
export const GearRangedMmPhase2 = PresetUtils.makePresetGear('P2 Ranged MM', Phase2GearRangedMm)
export const GearMeleePhase2 = PresetUtils.makePresetGear('P2 Melee', Phase2GearMelee)

export const GearPresets = {
  	[Phase.Phase1]: [
    	GearBeastMasteryPhase1,
			GearMarksmanPhase1,
			GearSurvivalPhase1,
  	],
  	[Phase.Phase2]: [
			GearRangedBmPhase2,
			GearRangedMmPhase2,
			GearMeleePhase2,
  	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearMeleePhase2;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLMeleeWeavePhase1 = PresetUtils.makePresetAPLRotation('P1 Melee Weave', MeleeWeaveP1);

export const APLMeleePhase2 = PresetUtils.makePresetAPLRotation('P2 Melee', MeleeP2);
export const APLRangedBmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged BM', RangedBmP2);
export const APLRangedMmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged MM', RangedMmP2);

export const APLPresets = {
  	[Phase.Phase1]: [
    	APLMeleeWeavePhase1,
  	],
  	[Phase.Phase2]: [
		APLRangedBmPhase2,
		APLRangedMmPhase2,
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
		1: APLPresets[Phase.Phase2][1],
		2: APLPresets[Phase.Phase2][2],
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
		talentsString: '-05551001503051',
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
export const DefaultTalentsMarksman 	= TalentPresets[CURRENT_PHASE][1];
export const DefaultTalentsSurvival 	= TalentPresets[CURRENT_PHASE][2];

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
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.SolidWeightstone,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectRegular,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectRegular,
  	blessingOfWisdom: TristateEffect.TristateEffectRegular,
  	sparkOfInspiration: true,
  	saygesFortune: SaygesFortune.SaygesDamage
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	faerieFire: true,
	homunculi: 100, // 70% average uptime default
	mangle: true,
});

export const OtherDefaults = {
	distanceFromTarget: 15,
  	profession1: Profession.Enchanting,
  	profession2: Profession.Leatherworking,
}
