import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AttackPowerBuff,
	Consumes,
	Debuffs,
	DragonslayerBuff,
	EnchantedSigil,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	ManaRegenElixir,
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
import { ElementalShaman_Options as ElementalShamanOptions } from '../core/proto/shaman.js';
import { SavedTalents } from '../core/proto/ui.js';
import Phase1APL from './apls/phase_1.apl.json';
import Phase2APL from './apls/phase_2.apl.json';
import Phase3APL from './apls/phase_3.apl.json';
import Phase4APL from './apls/phase_4.apl.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';
import Phase3Gear from './gear_sets/phase_3.gear.json';
import Phase4Gear from './gear_sets/phase_4.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);
export const GearPhase3 = PresetUtils.makePresetGear('Phase 3', Phase3Gear);
export const GearPhase4 = PresetUtils.makePresetGear('Phase 4 Burn', Phase4Gear);

export const GearPresets = {
	[Phase.Phase1]: [GearPhase1],
	[Phase.Phase2]: [GearPhase2],
	[Phase.Phase3]: [GearPhase3],
	[Phase.Phase4]: [GearPhase4],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearPresets[Phase.Phase4][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL);
export const APLPhase2 = PresetUtils.makePresetAPLRotation('Phase 2', Phase2APL);
export const APLPhase3 = PresetUtils.makePresetAPLRotation('Phase 3', Phase3APL);
export const APLPhase4 = PresetUtils.makePresetAPLRotation('Phase 4', Phase4APL);

export const APLPresets = {
	[Phase.Phase1]: [APLPhase1],
	[Phase.Phase2]: [APLPhase2],
	[Phase.Phase3]: [APLPhase3],
	[Phase.Phase4]: [APLPhase4],
	[Phase.Phase5]: [],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase4][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

export const TalentsPhase1 = PresetUtils.makePresetTalents('Level 25', SavedTalents.create({ talentsString: '25003105' }));
export const TalentsPhase2 = PresetUtils.makePresetTalents('Level 40', SavedTalents.create({ talentsString: '550031550000151' }));
export const TalentsPhase3 = PresetUtils.makePresetTalents('Level 50', SavedTalents.create({ talentsString: '550031550000151-500203' }));
export const TalentsPhase4 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '550031550000151--50105301005' }));

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2],
	[Phase.Phase3]: [TalentsPhase3],
	[Phase.Phase4]: [TalentsPhase4],
	[Phase.Phase5]: [],
};

export const DefaultTalents = TalentPresets[Phase.Phase4][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = ElementalShamanOptions.create({});

export const DefaultConsumes = Consumes.create({
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultPotion: Potions.MajorManaPotion,
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfSupremePower,
	food: Food.FoodRunnTumTuberSurprise,
	mainHandImbue: WeaponImbue.FlametongueWeapon,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	mildlyIrradiatedRejuvPot: true,
	offHandImbue: WeaponImbue.WizardOil,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.CerebralCortexCompound,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	demonicPact: 80,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: true,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	vampiricTouch: 300,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	dragonslayerBuff: DragonslayerBuff.RallyingCryofTheDragonslayer,
	fengusFerocity: true,
	mightOfStormwind: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	improvedFaerieFire: true,
	improvedScorch: true,
	markOfChaos: true,
	occultPoison: true,
	stormstrike: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	distanceFromTarget: 15,
	profession2: Profession.Alchemy,
	profession1: Profession.Enchanting,
};
