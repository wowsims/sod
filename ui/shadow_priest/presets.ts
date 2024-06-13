import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	EnchantedSigil,
	Flask,
	Food,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	SaygesFortune,
	ShadowPowerBuff,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { ShadowPriest_Options as Options } from '../core/proto/priest.js';
import { SavedTalents } from '../core/proto/ui.js';
import Phase1APL from './apls/phase_1.apl.json';
import Phase2APL from './apls/phase_2.apl.json';
import Phase3APL from './apls/phase_3.apl.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';
import Phase3Gear from './gear_sets/phase_3.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);
export const GearPhase3 = PresetUtils.makePresetGear('Phase 3', Phase3Gear);

export const GearPresets = {
	[Phase.Phase1]: [GearPhase1],
	[Phase.Phase2]: [GearPhase2],
	[Phase.Phase3]: [GearPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

// TODO: Add Phase 3 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase3][0];

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
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

// TODO: Add Phase 3 preset and pull from map
export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase3][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsPhase1 = {
	name: 'Phase 1',
	data: SavedTalents.create({
		talentsString: '-20535000001',
	}),
};

export const TalentsPhase2 = {
	name: 'Phase 2',
	data: SavedTalents.create({
		talentsString: '--5022204002501251',
	}),
};

export const TalentsPhase3 = {
	name: 'Phase 3',
	data: SavedTalents.create({
		talentsString: '-0055-5022204002501251',
	}),
};

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2],
	[Phase.Phase3]: [TalentsPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

// TODO: Add Phase 3 preset and pull from map
export const DefaultTalents = TalentPresets[Phase.Phase3][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = Options.create({});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	flask: Flask.FlaskOfRestlessDreams,
	food: Food.FoodNightfinSoup,
	mainHandImbue: WeaponImbue.LesserWizardOil,
	mildlyIrradiatedRejuvPot: true,
	shadowPowerBuff: ShadowPowerBuff.ElixirOfShadowPower,
	spellPowerBuff: SpellPowerBuff.ArcaneElixir,
	zanzaBuff: ZanzaBuff.AtalaiMojoOfForbiddenMagic,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	// aspectOfTheLion: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	fervorOfTheTempleExplorer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	songflowerSerenade: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfShadowNew: TristateEffect.TristateEffectRegular,
	improvedShadowBolt: true,
});

export const OtherDefaults = {
	channelClipDelay: 100,
	profession1: Profession.Alchemy,
	profession2: Profession.Enchanting,
};
