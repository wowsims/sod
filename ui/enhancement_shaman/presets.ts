import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	RaidBuffs,
	SapperExplosive,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { EnhancementShaman_Options as EnhancementShamanOptions, ShamanSyncType } from '../core/proto/shaman.js';
import { SavedTalents } from '../core/proto/ui.js';
import Phase1APL from './apls/phase_1.apl.json';
import Phase2APL from './apls/phase_2.apl.json';
import Phase3APL from './apls/phase_3.apl.json';
import Phase4APL from './apls/phase_4.apl.json';
import Phase5APL from './apls/phase_5.apl.json';
import Phase6APL from './apls/phase_6.apl.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';
import Phase3Gear from './gear_sets/phase_3.gear.json';
import Phase4Gear2H from './gear_sets/phase_4_2h.gear.json';
import Phase4GearDW from './gear_sets/phase_4_dw.gear.json';
import Phase5Gear2H from './gear_sets/phase_5_2h.gear.json';
import Phase5GearDW from './gear_sets/phase_5_dw.gear.json';
import Phase6Gear2H from './gear_sets/phase_6_2h.gear.json';
import Phase6GearDW from './gear_sets/phase_6_dw.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearPhase1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear, { customCondition: player => player.getLevel() === 25 });
export const GearPhase2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear, { customCondition: player => player.getLevel() === 40 });
export const GearPhase3 = PresetUtils.makePresetGear('Phase 3', Phase3Gear, { customCondition: player => player.getLevel() === 50 });
export const GearDWPhase4 = PresetUtils.makePresetGear('Phase 4 DW', Phase4GearDW, { customCondition: player => player.getLevel() === 60 });
export const Gear2HPhase4 = PresetUtils.makePresetGear('Phase 4 2H', Phase4Gear2H, { customCondition: player => player.getLevel() === 60 });
export const GearDWPhase5 = PresetUtils.makePresetGear('Phase 5 DW', Phase5GearDW, { customCondition: player => player.getLevel() === 60 });
export const Gear2HPhase5 = PresetUtils.makePresetGear('Phase 5 2H', Phase5Gear2H, { customCondition: player => player.getLevel() === 60 });
export const GearDWPhase6 = PresetUtils.makePresetGear('Phase 6 DW', Phase6GearDW, { customCondition: player => player.getLevel() === 60 });
export const Gear2HPhase6 = PresetUtils.makePresetGear('Phase 6 2H', Phase6Gear2H, { customCondition: player => player.getLevel() === 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearPhase1],
	[Phase.Phase2]: [GearPhase2],
	[Phase.Phase3]: [GearPhase3],
	[Phase.Phase4]: [GearDWPhase4, Gear2HPhase4],
	[Phase.Phase5]: [GearDWPhase5, Gear2HPhase5],
	[Phase.Phase6]: [GearDWPhase6, Gear2HPhase6],
};

export const DefaultGear = GearPresets[Phase.Phase6][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase1 = PresetUtils.makePresetAPLRotation('Phase 1', Phase1APL, { customCondition: player => player.getLevel() === 25 });
export const APLPhase2 = PresetUtils.makePresetAPLRotation('Phase 2', Phase2APL, { customCondition: player => player.getLevel() === 40 });
export const APLPhase3 = PresetUtils.makePresetAPLRotation('Phase 3', Phase3APL, { customCondition: player => player.getLevel() === 50 });
export const APLPhase4 = PresetUtils.makePresetAPLRotation('Phase 4', Phase4APL, { customCondition: player => player.getLevel() === 60 });
export const APLPhase5 = PresetUtils.makePresetAPLRotation('Phase 5', Phase5APL, { customCondition: player => player.getLevel() === 60 });
export const APLPhase6 = PresetUtils.makePresetAPLRotation('Phase 6', Phase6APL, { customCondition: player => player.getLevel() === 60 });

export const APLPresets = {
	[Phase.Phase1]: [APLPhase1],
	[Phase.Phase2]: [APLPhase2],
	[Phase.Phase3]: [APLPhase3],
	[Phase.Phase4]: [APLPhase4],
	[Phase.Phase5]: [APLPhase5],
	[Phase.Phase6]: [APLPhase6],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase6][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

export const TalentsPhase1 = PresetUtils.makePresetTalents('Level 25', SavedTalents.create({ talentsString: '-5005202101' }), {
	customCondition: player => player.getLevel() === 25,
});
export const TalentsPhase2 = PresetUtils.makePresetTalents('Level 40', SavedTalents.create({ talentsString: '-5005202105023051' }), {
	customCondition: player => player.getLevel() === 40,
});
export const TalentsPhase3 = PresetUtils.makePresetTalents('Level 50', SavedTalents.create({ talentsString: '05003-5005132105023051' }), {
	customCondition: player => player.getLevel() === 50,
});
export const TalentsPhase4 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '05023105003-5005032105023051' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [TalentsPhase1],
	[Phase.Phase2]: [TalentsPhase2],
	[Phase.Phase3]: [TalentsPhase3],
	[Phase.Phase4]: [TalentsPhase4],
	[Phase.Phase5]: [],
	[Phase.Phase6]: [],
};

export const DefaultTalents = TalentPresets[Phase.Phase4][0];

///////////////////////////////////////////////////////////////////////////
//                                 Build Presets
///////////////////////////////////////////////////////////////////////////

export const PresetBuildDW = PresetUtils.makePresetBuild('Dual Wield', {
	gear: GearDWPhase6,
	talents: TalentsPhase4,
	rotation: APLPhase6,
});
export const PresetBuild2H = PresetUtils.makePresetBuild('2-Handed', {
	gear: Gear2HPhase6,
	talents: TalentsPhase4,
	rotation: APLPhase6,
});

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = EnhancementShamanOptions.create({
	syncType: ShamanSyncType.Auto,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	flask: Flask.FlaskOfAncientKnowledge,
	food: Food.FoodSmokedDesertDumpling,
	mainHandImbue: WeaponImbue.WindfuryWeapon,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	mildlyIrradiatedRejuvPot: true,
	offHandImbue: WeaponImbue.FlametongueWeapon,
	sapperExplosive: SapperExplosive.SapperFumigator,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	demonicPact: 120,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	vampiricTouch: 300,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	fengusFerocity: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	spiritOfZandalar: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	improvedScorch: true,
	markOfChaos: true,
	occultPoison: true,
	stormstrike: true,
});

export const OtherDefaults = {
	profession1: Profession.Alchemy,
	profession2: Profession.Engineering,
};
