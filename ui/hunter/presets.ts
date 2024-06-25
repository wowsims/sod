import { Phase } from '../core/constants/other.js';
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
import Phase2GearMelee from './gear_sets/p2_melee.gear.json';
import Phase2GearRangedBm from './gear_sets/p2_ranged_bm.gear.json';
import Phase2GearRangedMm from './gear_sets/p2_ranged_mm.gear.json';
import Phase3GearMeleeBm from './gear_sets/p3_melee_bm.gear.json';
import Phase3GearRangedMm from './gear_sets/p3_ranged_mm.gear.json';
import Phase1Gear from './gear_sets/phase1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBeastMasteryPhase1 = PresetUtils.makePresetGear('P1 Beast Mastery', Phase1Gear, { customCondition: player => player.getLevel() == 25 });
export const GearMarksmanPhase1 = PresetUtils.makePresetGear('P1 Marksmanship', Phase1Gear, { customCondition: player => player.getLevel() == 25 });
export const GearSurvivalPhase1 = PresetUtils.makePresetGear('P1 Survival', Phase1Gear, { customCondition: player => player.getLevel() == 25 });

export const GearRangedBmPhase2 = PresetUtils.makePresetGear('P2 Ranged BM', Phase2GearRangedBm, { customCondition: player => player.getLevel() == 40 });
export const GearRangedMmPhase2 = PresetUtils.makePresetGear('P2 Ranged MM', Phase2GearRangedMm, { customCondition: player => player.getLevel() == 40 });
export const GearMeleePhase2 = PresetUtils.makePresetGear('P2 Melee', Phase2GearMelee, { customCondition: player => player.getLevel() == 40 });

export const GearMeleeBmPhase3 = PresetUtils.makePresetGear('P3 Melee BM', Phase3GearMeleeBm, { customCondition: player => player.getLevel() == 50 });
export const GearRangedMmPhase3 = PresetUtils.makePresetGear('P3 Ranged MM', Phase3GearRangedMm, { customCondition: player => player.getLevel() == 50 });

export const GearPresets = {
	[Phase.Phase1]: [GearBeastMasteryPhase1, GearMarksmanPhase1, GearSurvivalPhase1],
	[Phase.Phase2]: [GearRangedBmPhase2, GearRangedMmPhase2, GearMeleePhase2],
	[Phase.Phase3]: [GearRangedMmPhase3, GearMeleeBmPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearMeleeBmPhase3;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLMeleeWeavePhase1 = PresetUtils.makePresetAPLRotation('P1 Melee Weave', MeleeWeaveP1, { customCondition: player => player.getLevel() == 25 });

export const APLMeleePhase2 = PresetUtils.makePresetAPLRotation('P2 Melee', MeleeP2, { customCondition: player => player.getLevel() == 40 });
export const APLRangedBmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged BM', RangedBmP2, { customCondition: player => player.getLevel() == 40 });
export const APLRangedMmPhase2 = PresetUtils.makePresetAPLRotation('P2 Ranged MM', RangedMmP2, { customCondition: player => player.getLevel() == 40 });

export const APLMeleeBmPhase3 = PresetUtils.makePresetAPLRotation('P3 Melee BM', MeleeBmP3, { customCondition: player => player.getLevel() == 50 });
export const APLRangedMmPhase3 = PresetUtils.makePresetAPLRotation('P3 Ranged MM', RangedMmP3, { customCondition: player => player.getLevel() == 50 });

export const APLPresets = {
	[Phase.Phase1]: [APLMeleeWeavePhase1],
	[Phase.Phase2]: [APLRangedBmPhase2, APLRangedMmPhase2, APLMeleePhase2],
	[Phase.Phase3]: [APLRangedMmPhase3, APLMeleeBmPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsBeastMasteryPhase1 = PresetUtils.makePresetTalents('P1 Beast Mastery', SavedTalents.create({ talentsString: '53000200501' }), {
	customCondition: player => player.getLevel() == 25,
});

export const TalentsMarksmanPhase1 = PresetUtils.makePresetTalents('P1 Marksmanship', SavedTalents.create({ talentsString: '-050515' }), {
	customCondition: player => player.getLevel() == 25,
});

export const TalentsSurvivalPhase1 = PresetUtils.makePresetTalents('P1 Survival', SavedTalents.create({ talentsString: '--33502001101' }), {
	customCondition: player => player.getLevel() == 25,
});

export const TalentsBeastMasteryPhase2 = PresetUtils.makePresetTalents('P2 Beast Mastery', SavedTalents.create({ talentsString: '5300021150501251' }), {
	customCondition: player => player.getLevel() == 40,
});

export const TalentsMarksmanPhase2 = PresetUtils.makePresetTalents('P2 Marksmanship', SavedTalents.create({ talentsString: '-05551001503051' }), {
	customCondition: player => player.getLevel() == 40,
});

export const TalentsSurvivalPhase2 = PresetUtils.makePresetTalents('P2 Survival', SavedTalents.create({ talentsString: '--335020051030315' }), {
	customCondition: player => player.getLevel() == 40,
});

export const TalentsRangedMMPhase3 = PresetUtils.makePresetTalents('P3 Ranged MM', SavedTalents.create({ talentsString: '5-05051404503051-3' }), {
	customCondition: player => player.getLevel() == 50,
});
export const TalentsMeleeBMPhase3 = PresetUtils.makePresetTalents('P3 Melee BM', SavedTalents.create({ talentsString: '0500321150521251--33002' }), {
	customCondition: player => player.getLevel() == 50,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsBeastMasteryPhase1, TalentsMarksmanPhase1, TalentsSurvivalPhase1],
	[Phase.Phase2]: [TalentsBeastMasteryPhase2, TalentsMarksmanPhase2, TalentsSurvivalPhase2],
	[Phase.Phase3]: [TalentsRangedMMPhase3, TalentsMeleeBMPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultTalents = TalentsMeleeBMPhase3;

export const PresetBuildMeleeBM = PresetUtils.makePresetBuild('Melee BM', GearMeleeBmPhase3, TalentsMeleeBMPhase3, APLMeleeBmPhase3);
export const PresetBuildRangedMM = PresetUtils.makePresetBuild('Ranged MM', GearRangedMmPhase3, TalentsRangedMMPhase3, APLRangedMmPhase3);

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
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.DenseSharpeningStone,
	spellPowerBuff: SpellPowerBuff.ArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	zanzaBuff: ZanzaBuff.AtalaiMojoOfWar,
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
	fervorOfTheTempleExplorer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	songflowerSerenade: true,
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
};
