import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
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
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
} from '../core/proto/common.js';
import { EnhancementShaman_Options as EnhancementShamanOptions, ShamanSyncType } from '../core/proto/shaman.js';
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
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

// TODO: Phase 3
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
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
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
		talentsString: '-5005202101',
	}),
};

export const TalentsPhase2 = {
	name: 'Phase 2',
	data: SavedTalents.create({
		talentsString: '-5005202105023051',
	}),
};

export const TalentsPhase3 = {
	name: 'Phase 3',
	data: SavedTalents.create({
		talentsString: '05003-5005222105023051',
	}),
};

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2],
	[Phase.Phase3]: [TalentsPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultTalents = TalentPresets[Phase.Phase3][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = EnhancementShamanOptions.create({
	syncType: ShamanSyncType.Auto,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfAgility,
	defaultPotion: Potions.GreaterManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.WindfuryWeapon,
	offHandImbue: WeaponImbue.WindfuryWeapon,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	saygesFortune: SaygesFortune.SaygesDamage,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElementsNew: TristateEffect.TristateEffectRegular,
	curseOfRecklessness: true,
	dreamstate: true,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	improvedScorch: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Enchanting,
	profession2: Profession.Leatherworking,
};
