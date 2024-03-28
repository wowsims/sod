import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	EnchantedSigil,
	FirePowerBuff,
	Food,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue,
} from '../core/proto/common.js';
import { ElementalShaman_Options as ElementalShamanOptions } from '../core/proto/shaman.js';
import { SavedTalents } from '../core/proto/ui.js';
import Phase1APL from './apls/phase_1.apl.json';
import Phase2APL from './apls/phase_2.apl.json';
import Phase3APL from './apls/phase_3.apl.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);

export const GearPresets = {
	[Phase.Phase1]: [GearPhase1],
	[Phase.Phase2]: [GearPhase2],
	[Phase.Phase3]: [],
};

export const DefaultGear = GearPresets[Phase.Phase2][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL);
export const APLPhase2 = PresetUtils.makePresetAPLRotation('Phase 2', Phase2APL);
export const APLPhase3 = PresetUtils.makePresetAPLRotation('Phase 3', Phase3APL);

export const APLPresets = {
	[Phase.Phase1]: [APLPhase1],
	[Phase.Phase2]: [APLPhase2],
	[Phase.Phase3]: [APLPhase3],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsPhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '25003105',
	}),
};

export const TalentsPhase2 = {
	name: 'Phase 2',
	data: SavedTalents.create({
		talentsString: '550031550000151',
	}),
};

export const TalentsPhase3 = {
	name: 'Phase 3',
	data: SavedTalents.create({
		talentsString: '550031450030151-500003',
	}),
};

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2],
	[Phase.Phase3]: [TalentsPhase3],
};

export const DefaultTalents = TalentPresets[Phase.Phase3][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = ElementalShamanOptions.create({});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.LesserWizardOil,
	offHandImbue: WeaponImbue.LesserWizardOil,
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
	saygesFortune: SaygesFortune.SaygesDamage,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
	improvedScorch: true,
	stormstrike: true,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
	profession1: Profession.Enchanting,
	profession2: Profession.Leatherworking,
};
