import { Phase } from '../core/constants/other.js';
import {
	AgilityElixir,
	AttackPowerBuff,
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
	Spec,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
	UnitReference
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	FeralTankDruid_Options as DruidOptions,
	FeralTankDruid_Rotation as DruidRotation,
	FeralTankDruid_Rotation,
} from '../core/proto/druid.js';

import * as PresetUtils from '../core/preset_utils.js';

import Phase5Gear from './gear_sets/phase_5.gear.json';

import Phase5APL from './apls/phase_5.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const P5Gear = PresetUtils.makePresetGear('Phase 5', Phase5Gear);

export const GearPresets = {
  [Phase.Phase1]: [P5Gear],
  [Phase.Phase2]: [P5Gear],
  [Phase.Phase3]: [P5Gear],
  [Phase.Phase4]: [P5Gear],
  [Phase.Phase5]: [P5Gear]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase5][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase5 = PresetUtils.makePresetAPLRotation('Phase 5', Phase5APL, { customCondition: player => player.getLevel() === 60 });

export const APLPresets = {
  [Phase.Phase1]: [APLPhase5],
  [Phase.Phase2]: [APLPhase5],
  [Phase.Phase3]: [APLPhase5],
  [Phase.Phase4]: [APLPhase5],
  [Phase.Phase5]: [APLPhase5]
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
  25: APLPresets[Phase.Phase1][0],
  40: APLPresets[Phase.Phase2][0],
  50: APLPresets[Phase.Phase3][0],
  60: APLPresets[Phase.Phase5][0]
};

export const DefaultRotation = FeralTankDruid_Rotation.create({
	maulRageThreshold: 25,
	maintainDemoralizingRoar: true,
	lacerateTime: 8.0,
});

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('Default', Spec.SpecFeralTankDruid, DefaultRotation);

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '500005001-5050321303022151-05002',
	}),
};

export const TalentPresets = {
  [Phase.Phase1]: [StandardTalents],
  [Phase.Phase2]: [StandardTalents],
  [Phase.Phase3]: [StandardTalents],
  [Phase.Phase4]: [StandardTalents],
  [Phase.Phase5]: [StandardTalents],
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultTalents = TalentPresets[Phase.Phase5][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = DruidOptions.create({
	innervateTarget: UnitReference.create(),
	startingRage: 20,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultPotion: Potions.GreaterStoneshieldPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.FlowingWatersSigil,
	flask: Flask.FlaskOfTheTitans,
	food: Food.FoodDirgesKickChimaerokChops,
	mainHandImbue: WeaponImbue.ElementalSharpeningStone,
	miscConsumes: {
		catnip: true,
		jujuEmber: true,
	},
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	mightOfStormwind: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
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
	sunderArmor: true,
});

export const OtherDefaults = {
	profession1: Profession.Enchanting,
	profession2: Profession.Alchemy,
};