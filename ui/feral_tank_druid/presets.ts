import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
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
	StrengthBuff,
	TristateEffect,
	UnitReference,
	WeaponImbue,
	ZanzaBuff} from '../core/proto/common.js';
import {
	FeralTankDruid_Options as DruidOptions,
	FeralTankDruid_Rotation,
} from '../core/proto/druid.js';
import { SavedTalents } from '../core/proto/ui.js';
import CatWeaveRakeAPL from './apls/catweave_rake.apl.json';
import Phase6APL from './apls/phase_6.apl.json';
import Phase6Gear from './gear_sets/phase_6.gear.json';
import Phase7Gear from './gear_sets/phase_7.gear.json';
import Phase7CatWeave from './gear_sets/phase_7_catweave.gear.json';
import Phase7PreBis from './gear_sets/phase_7_prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const P6Gear = PresetUtils.makePresetGear('Phase 6', Phase6Gear);
export const P7Gear = PresetUtils.makePresetGear('Phase 7', Phase7Gear);
export const P7PreBis = PresetUtils.makePresetGear('Phase 7 Pre-bis', Phase7PreBis);
export const P7CatWeave = PresetUtils.makePresetGear('Phase 7 Cat Weave', Phase7CatWeave);

export const GearPresets = {
  [Phase.Phase1]: [],
  [Phase.Phase2]: [],
  [Phase.Phase3]: [],
  [Phase.Phase4]: [],
  [Phase.Phase5]: [],
  [Phase.Phase6]: [P6Gear],
  [Phase.Phase7]: [P7PreBis,P7Gear,P7CatWeave],
};

export const DefaultGear = GearPresets[Phase.Phase7][1];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPhase6 = PresetUtils.makePresetAPLRotation('Phase 6', Phase6APL, { customCondition: player => player.getLevel() === 60 });
export const APLCatWeaveRake = PresetUtils.makePresetAPLRotation('Cat Weave Rake', CatWeaveRakeAPL, { customCondition: player => player.getLevel() === 60 });

export const APLPresets = {
  [Phase.Phase1]: [],
  [Phase.Phase2]: [],
  [Phase.Phase3]: [],
  [Phase.Phase4]: [],
  [Phase.Phase5]: [],
  [Phase.Phase6]: [APLPhase6],
  [Phase.Phase7]: [APLPhase6, APLCatWeaveRake]
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
  25: APLPresets[Phase.Phase6][0],
  40: APLPresets[Phase.Phase6][0],
  50: APLPresets[Phase.Phase6][0],
  60: APLPresets[Phase.Phase7][0]
};

export const DefaultRotation = FeralTankDruid_Rotation.create({
	maulRageThreshold: 25,
	maintainDemoralizingRoar: true,
	lacerateTime: 8.0,
});

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '014005301-5050021323022151-05',
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
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultPotion: Potions.GreaterStoneshieldPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfTheOldGods,
	food: Food.FoodDirgesKickChimaerokChops,
	mainHandImbue: WeaponImbue.ElementalSharpeningStone,
	miscConsumes: {
		catnip: true,
		jujuEmber: true,
	},
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.SpiritOfZanza,
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
	moldarsMoxie: true,
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