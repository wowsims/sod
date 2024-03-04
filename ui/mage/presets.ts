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
	Mage_Options_ArmorType as ArmorType
} from '../core/proto/mage.js';
import { SavedTalents } from '../core/proto/ui.js';
import APLDefault from './apls/default.apl.json';
import Phase2APLArcane from './apls/p2_arcane.apl.json';
import Phase2APLFire from './apls/p2_fire.apl.json';
import Phase2GearArcane from './gear_sets/p2_arcane.gear.json';
import Phase2GearFire from './gear_sets/p2_fire.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearArcanePhase2 = PresetUtils.makePresetGear('P2 Arcane', Phase2GearArcane, { talentTree: 0 })
export const GearFirePhase2 = PresetUtils.makePresetGear('P2 Fire', Phase2GearFire, { talentTree: 1 })

export const GearPresets = {
	[Phase.Phase1]: [
		GearArcanePhase2,
		GearArcanePhase2,
		GearArcanePhase2,
	],
	[Phase.Phase2]: [
		GearArcanePhase2,
		GearFirePhase2,
	],
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[CURRENT_PHASE][1];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLArcanePhase1 = PresetUtils.makePresetAPLRotation('Default', APLDefault, { talentTree: 0 });
export const APLFirePhase1 = PresetUtils.makePresetAPLRotation('Default', APLDefault, { talentTree: 1 });
export const APLFrostPhase1 = PresetUtils.makePresetAPLRotation('Default', APLDefault, { talentTree: 2 });

export const APLArcanePhase2 = PresetUtils.makePresetAPLRotation('P2 Arcane', Phase2APLArcane, { talentTree: 0 });
export const APLFirePhase2 = PresetUtils.makePresetAPLRotation('P2 Fire', Phase2APLFire, { talentTree: 1 });
export const APLFrostPhase2 = PresetUtils.makePresetAPLRotation('P2 Frost', APLDefault, { talentTree: 2 });

export const APLPresets = {
	[Phase.Phase1]: [
		APLArcanePhase1,
		APLFirePhase1,
		APLFrostPhase1,
	],
	[Phase.Phase2]: [
		APLArcanePhase2,
		APLFirePhase2,
		APLFrostPhase2,
	],
};

// TODO: Add Phase 2 preset and pull from map
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

export const TalentsArcanePhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '50005003021',
	}),
};

export const TalentsFirePhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '50005003021',
	}),
};

export const TalentsFrostPhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '50005003021',
	}),
};

export const TalentsArcanePhase2 = {
	name: 'P2 Arcane',
	data: SavedTalents.create({
		talentsString: '2250050310031531',
	})
};

export const TalentsFirePhase2 = {
	name: 'P2 Fire',
	data: SavedTalents.create({
		talentsString: '-5050020123033151',
	})
};

export const TalentsFrostPhase2 = {
	name: 'P2 Frost',
	data: SavedTalents.create({
		talentsString: '--0535020310025005',
	})
};

export const TalentPresets = {
	[Phase.Phase1]: [
    	TalentsArcanePhase1,
		TalentsFirePhase1,
		TalentsFrostPhase1,
	],
	[Phase.Phase2]: [
		TalentsArcanePhase2,
		TalentsFirePhase2,
		TalentsFrostPhase2,
	],
};

export const DefaultTalentsArcane 	= TalentPresets[CURRENT_PHASE][0];
export const DefaultTalentsFire 	= TalentPresets[CURRENT_PHASE][1];
export const DefaultTalentsFrost 	= TalentPresets[CURRENT_PHASE][2];

export const DefaultTalents = DefaultTalentsArcane;

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
