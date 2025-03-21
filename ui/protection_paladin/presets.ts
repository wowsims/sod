import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	Alcohol,
	ArmorElixir,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	Explosive,
	FirePowerBuff,
	Flask,
	Food,
	HealthElixir,
	IndividualBuffs,
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
import { Blessings, PaladinAura, PaladinOptions as ProtectionPaladinOptions,PaladinSeal } from '../core/proto/paladin.js';
import { SavedTalents } from '../core/proto/ui.js';
import APLP4ProtJson from './apls/p4prot.apl.json';
import APLP5ProtJson from './apls/p5prot.apl.json';
import Phase4ProtGearJson from './gear_sets/p4prot.gear.json';
import Phase5ProtGearJson from './gear_sets/p5prot.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const Phase4ProtGear = PresetUtils.makePresetGear('P4 Prot', Phase4ProtGearJson);
export const Phase5ProtGear = PresetUtils.makePresetGear('P5 Prot', Phase5ProtGearJson);

export const GearPresets = {
	[Phase.Phase1]: [],
	[Phase.Phase2]: [],
	[Phase.Phase3]: [],
	[Phase.Phase4]: [Phase4ProtGear, Phase5ProtGear],
	[Phase.Phase5]: [Phase4ProtGear, Phase5ProtGear],
};

export const DefaultGear = GearPresets[Phase.Phase5][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLP5Prot = PresetUtils.makePresetAPLRotation('P5 Prot', APLP5ProtJson);
export const APLP4Prot = PresetUtils.makePresetAPLRotation('P4 Prot', APLP4ProtJson);

export const APLPresets = {
	[Phase.Phase1]: [],
	[Phase.Phase2]: [],
	[Phase.Phase3]: [],
	[Phase.Phase4]: [APLP4Prot, APLP5Prot],
	[Phase.Phase5]: [APLP4Prot, APLP5Prot],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	60: APLPresets[Phase.Phase5][0],
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const P4ProtTalents = {
	name: 'P4 Prot',
	data: SavedTalents.create({
		talentsString: '-053020335001551-0500535',
	}),
};

export const P5ProtTalents = {
	name: 'P5 Prot',
	data: SavedTalents.create({
		talentsString: '-053020335001551-0520335',
	}),
};

export const TalentPresets = {
	[Phase.Phase1]: [],
	[Phase.Phase2]: [],
	[Phase.Phase3]: [],
	[Phase.Phase4]: [P4ProtTalents],
	[Phase.Phase5]: [P5ProtTalents],
};

export const DefaultTalents = TalentPresets[Phase.Phase5][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = ProtectionPaladinOptions.create({
	aura: PaladinAura.SanctityAura,
	primarySeal: PaladinSeal.Martyrdom,
	personalBlessing: Blessings.BlessingOfSanctuary,
	righteousFury: true,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	healthElixir: HealthElixir.ElixirOfFortitude,
	armorElixir: ArmorElixir.ElixirOfSuperiorDefense,
	defaultPotion: Potions.GreaterStoneshieldPotion,
	dragonBreathChili: true,
	food: Food.FoodDarkclawBisque,
	flask: Flask.FlaskOfTheTitans,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	fillerExplosive: Explosive.ExplosiveDenseDynamite,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.MagnificentTrollshine,

	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.ROIDS,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredDemonicRune,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	moldarsMoxie: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	spiritOfZandalar: true,
	valorOfAzeroth: true,
	mightOfStormwind: true,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	hornOfLordaeron: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	battleShout: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	sanctityAura: true,
	demonicPact: 120,
	aspectOfTheLion: true,
	vampiricTouch: 300,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	homunculi: 70, // 70% average uptime default
	faerieFire: true,
	giftOfArthas: true,
	exposeArmor: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	judgementOfTheCrusader: TristateEffect.TristateEffectImproved,
	improvedScorch: true,
	markOfChaos: true,
	occultPoison: true,
	mangle: true,
});

export const OtherDefaults = {
	distanceFromTarget: 5, // Max melee range
	profession1: Profession.Blacksmithing,
	profession2: Profession.Engineering,
};
