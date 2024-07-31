import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	Alcohol,
	ArmorElixir,
	AttackPowerBuff,
	Consumes,
	Debuffs,
	EnchantedSigil,
	FirePowerBuff,
	Flask,
	Food,
	HealthElixir,
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
import { EnhancementShaman_Options as EnhancementShamanOptions, ShamanSyncType } from '../core/proto/shaman.js';
import { SavedTalents } from '../core/proto/ui.js';
import Phase4EleTankAPLJSON from './apls/phase_4_ele_tank.apl.json';
import Phase4EnhTankAPLJSON from './apls/phase_4_enh_tank.apl.json';
import Phase4EleTankGearJSON from './gear_sets/phase_4_ele_tank.gear.json';
import Phase4EnhTankGearJSON from './gear_sets/phase_4_enh_tank.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearEnhTankPhase4 = PresetUtils.makePresetGear('P4 Enh', Phase4EnhTankGearJSON);
export const GearEleTankPhase4 = PresetUtils.makePresetGear('P4 Ele', Phase4EleTankGearJSON);

export const GearPresets = {
	[Phase.Phase1]: [],
	[Phase.Phase2]: [],
	[Phase.Phase3]: [],
	[Phase.Phase4]: [GearEnhTankPhase4, GearEleTankPhase4],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearPresets[Phase.Phase4][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLEnhTankPhase4 = PresetUtils.makePresetAPLRotation('P4 Enh', Phase4EnhTankAPLJSON);
export const APLEleTankPhase4 = PresetUtils.makePresetAPLRotation('P4 Ele', Phase4EleTankAPLJSON);

export const APLPresets = {
	[Phase.Phase1]: [],
	[Phase.Phase2]: [],
	[Phase.Phase3]: [],
	[Phase.Phase4]: [APLEnhTankPhase4, APLEleTankPhase4],
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

export const TalentsHardCastPhase4 = PresetUtils.makePresetTalents('60 Ele', SavedTalents.create({ talentsString: '55030155000015-050003-500053' }));
export const TalentsSpellhancePhase4 = PresetUtils.makePresetTalents('60 Spellhance', SavedTalents.create({ talentsString: '5503015000301-05250320550001' }));
export const TalentsDeepEnhPhase4 = PresetUtils.makePresetTalents('60 Enh', SavedTalents.create({ talentsString: '05033150003-0505032015003151' }));

export const TalentPresets = {
	[Phase.Phase1]: [],
	[Phase.Phase2]: [],
	[Phase.Phase3]: [],
	[Phase.Phase4]: [TalentsDeepEnhPhase4, TalentsSpellhancePhase4, TalentsHardCastPhase4],
	[Phase.Phase5]: [],
};

export const DefaultTalents = TalentPresets[Phase.Phase4][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = EnhancementShamanOptions.create({
	syncType: ShamanSyncType.Auto,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	armorElixir: ArmorElixir.ElixirOfSuperiorDefense,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfTheTitans,
	food: Food.FoodBlessSunfruit,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.RockbiterWeapon,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	mildlyIrradiatedRejuvPot: true,
	offHandImbue: WeaponImbue.ConductiveShieldCoating,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	commandingShout: true,
	demonicPact: 80,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	vampiricTouch: 300,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	fengusFerocity: true,
	moldarsMoxie: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	improvedFaerieFire: true,
	improvedScorch: true,
	insectSwarm: true,
	markOfChaos: true,
	occultPoison: true,
	stormstrike: true,
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Alchemy,
	profession2: Profession.Enchanting,
};
