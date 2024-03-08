import { CURRENT_PHASE, Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	EnchantedSigil,
	FirePowerBuff,
	Food,
	FrostPowerBuff,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue
} from '../core/proto/common.js';
import {
	Mage_Options as MageOptions,
	Mage_Options_ArmorType as ArmorType,
} from '../core/proto/mage.js';
import { SavedTalents } from '../core/proto/ui.js';
import Phase1APLArcane from './apls/p1_arcane.apl.json';
import Phase1APLFire from './apls/p1_fire.apl.json';
import Phase2APLArcane from './apls/p2_arcane.apl.json';
import Phase2APLFire from './apls/p2_fire.apl.json';
import Phase2APLFrostfire from './apls/p2_frostfire.apl.json';
import Phase1GearFire from './gear_sets/p1_fire.gear.json';
import Phase1Gear from './gear_sets/p1_generic.gear.json';
import Phase2GearArcane from './gear_sets/p2_arcane.gear.json';
import Phase2GearFire from './gear_sets/p2_fire.gear.json';
import Phase2GearFrostfire from './gear_sets/p2_frostfire.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArcanePhase1 = PresetUtils.makePresetGear('P1 Arcane', Phase1Gear, { customCondition: player => player.getLevel() == 25 });
export const GearFirePhase1 = PresetUtils.makePresetGear('P1 Fire', Phase1GearFire, { customCondition: player => player.getLevel() == 25 });
export const GearFrostPhase1 = PresetUtils.makePresetGear('P1 Frost', Phase1Gear, { customCondition: player => player.getLevel() == 25 });

export const GearArcanePhase2 = PresetUtils.makePresetGear('P2 Arcane', Phase2GearArcane, { customCondition: player => player.getLevel() == 40 });
export const GearFirePhase2 = PresetUtils.makePresetGear('P2 Fire', Phase2GearFire, { customCondition: player => player.getLevel() == 40 });
export const GearFrosthase2 = PresetUtils.makePresetGear('P2 Frostfire', Phase2GearFrostfire, { customCondition: player => player.getLevel() == 40 });

export const GearPresets = {
	[Phase.Phase1]: [
		GearArcanePhase1,
		GearFirePhase1,
		GearFrostPhase1,
	],
	[Phase.Phase2]: [
		GearArcanePhase2,
		GearFirePhase2,
		GearFrosthase2,
	],
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[CURRENT_PHASE][1];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLArcanePhase1 = PresetUtils.makePresetAPLRotation('P1 Arcane', Phase1APLArcane, { customCondition: player => player.getLevel() == 25 });
export const APLFirePhase1 = PresetUtils.makePresetAPLRotation('P1 Fire', Phase1APLFire, { customCondition: player => player.getLevel() == 25 });

export const APLArcanePhase2 = PresetUtils.makePresetAPLRotation('P2 Arcane', Phase2APLArcane, { customCondition: player => player.getLevel() == 40 });
export const APLFirePhase2 = PresetUtils.makePresetAPLRotation('P2 Fire', Phase2APLFire, { customCondition: player => player.getLevel() == 40 });
export const APLFrostfirePhase2 = PresetUtils.makePresetAPLRotation('P2 Frostfire', Phase2APLFrostfire, { customCondition: player => player.getLevel() == 40 });

export const APLPresets = {
	[Phase.Phase1]: [
		APLArcanePhase1,
		APLFirePhase1,
		APLFirePhase1,
	],
	[Phase.Phase2]: [
		APLArcanePhase2,
		APLFirePhase2,
		APLFrostfirePhase2,
	],
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
		2: APLPresets[Phase.Phase2][2],
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsArcanePhase1 = PresetUtils.makePresetTalents('P1 Arcane', SavedTalents.create({talentsString: '22500502'}), {
	customCondition: player => player.getLevel() == 25,
});

export const TalentsFirePhase1 = PresetUtils.makePresetTalents('P1 Fire', SavedTalents.create({talentsString: '-5050020121'}), {
	customCondition: player => player.getLevel() == 25,
});

export const TalentsArcanePhase2 = PresetUtils.makePresetTalents('P2 Arcane', SavedTalents.create({talentsString: '2250050310031531'}), {
	customCondition: player => player.getLevel() == 40,
});

export const TalentsFirePhase2 = PresetUtils.makePresetTalents('P2 Fire', SavedTalents.create({talentsString: '-5050020123033151'}), {
	customCondition: player => player.getLevel() == 40,
});

export const TalentsFrostPhase2 = PresetUtils.makePresetTalents('P2 Frostfire', SavedTalents.create({talentsString: '-055002012002-203500031'}), {
	customCondition: player => player.getLevel() == 40,
});

export const TalentPresets = {
	[Phase.Phase1]: [
    	TalentsArcanePhase1,
		TalentsFirePhase1,
		TalentsFirePhase1,
	],
	[Phase.Phase2]: [
		TalentsArcanePhase2,
		TalentsFirePhase2,
		TalentsFrostPhase2,
	],
};

export const DefaultTalentsArcane = TalentPresets[CURRENT_PHASE][0];
export const DefaultTalentsFire	= TalentPresets[CURRENT_PHASE][1];
export const DefaultTalentsFrost = TalentPresets[CURRENT_PHASE][2];

export const DefaultTalents = DefaultTalentsFire;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodSagefishDelight,
	frostPowerBuff: FrostPowerBuff.ElixirOfFrostPower,
	mainHandImbue: WeaponImbue.LesserWizardOil,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
  	sparkOfInspiration: true,
  	saygesFortune: SaygesFortune.SaygesDamage
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
});

export const OtherDefaults = {
  	distanceFromTarget: 20,
  	profession1: Profession.Enchanting,
  	profession2: Profession.Tailoring,
};
