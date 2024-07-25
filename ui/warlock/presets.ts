import { Phase } from '../core/constants/other.js';
import { Player } from '../core/player.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	Alcohol,
	Conjured,
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
	ShadowPowerBuff,
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	WarlockOptions as WarlockOptions,
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';
// apls
import DestroP1APL from './apls/p1/destruction.apl.json';
import AfflictionAPL from './apls/p2/affliction.apl.json';
import DemonologyAPL from './apls/p2/demonology.apl.json';
import DestroConflagAPL from './apls/p2/fire.conflag.apl.json';
import DestroMgiAPL from './apls/p2/fire.imp.apl.json';
import BackdraftAPLP3 from './apls/p3/backdraft.apl.json';
import NFRuinAPLP3 from './apls/p3/nf.ruin.apl.json';
import AffAPLP4 from './apls/p4/affliction.apl.json';
import DestroAPLP4 from './apls/p4/destruction.apl.json';
// gear
import DestructionGear from './gear_sets/p1/destruction.gear.json';
import FireImpGear from './gear_sets/p2/fire.imp.gear.json';
import FireSuccubusGear from './gear_sets/p2/fire.succubus.gear.json';
import ShadowGear from './gear_sets/p2/shadow.gear.json';
import BackdraftGearP3 from './gear_sets/p3/backdraft.gear.json';
import NFRuinGearP3 from './gear_sets/p3/nf.ruin.gear.json';
import AffGearP4 from './gear_sets/p4/affliction.gear.json';
import DestroGearP4 from './gear_sets/p4/destruction.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearDestructionPhase1 = PresetUtils.makePresetGear('Destruction', DestructionGear, { customCondition: player => player.getLevel() == 25 });

export const FireImpGearPhase2 = PresetUtils.makePresetGear('P2 Fire Imp', FireImpGear, { customCondition: player => player.getLevel() == 40 });
export const FireSuccubusGearPhase2 = PresetUtils.makePresetGear('P2 Fire Succubus', FireSuccubusGear, { customCondition: player => player.getLevel() == 40 });
export const ShadowGearPhase2 = PresetUtils.makePresetGear('P2 Shadow', ShadowGear, { customCondition: player => player.getLevel() == 40 });

export const BackdraftGearPhase3 = PresetUtils.makePresetGear('P3 Backdraft', BackdraftGearP3, { customCondition: player => player.getLevel() == 50 });
export const NFRuinGearPhase3 = PresetUtils.makePresetGear('P3 NF/Ruin', NFRuinGearP3, { customCondition: player => player.getLevel() == 50 });

export const AffGearPhase4 = PresetUtils.makePresetGear('P4 Aff', AffGearP4, { customCondition: player => player.getLevel() == 60 });
export const DestroGearPhase4 = PresetUtils.makePresetGear('P4 Destro', DestroGearP4, { customCondition: player => player.getLevel() == 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearDestructionPhase1],
	[Phase.Phase2]: [FireImpGearPhase2, FireSuccubusGearPhase2, ShadowGearPhase2],
	[Phase.Phase3]: [NFRuinGearPhase3, BackdraftGearPhase3],
	[Phase.Phase4]: [AffGearPhase4, DestroGearPhase4],
	[Phase.Phase5]: [],
};

export const DefaultGearAff = AffGearPhase4;
export const DefaultGearDestro = DestroGearPhase4;

export const DefaultGear = DefaultGearDestro;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

// P1
export const RotationDestructionPhase1 = PresetUtils.makePresetAPLRotation('Destruction', DestroP1APL, {
	customCondition: player => player.getLevel() == 25,
});

// P2
export const DestroMgiRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Destro Imp', DestroMgiAPL, {
	customCondition: player => player.getLevel() == 40,
});
export const DestroConflagRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Destro Conflag', DestroConflagAPL, {
	customCondition: player => player.getLevel() == 40,
});
export const DemonologyRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Demonology', DemonologyAPL, {
	customCondition: player => player.getLevel() == 40,
});
export const AfflictionRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Affliction', AfflictionAPL, {
	customCondition: player => player.getLevel() == 40,
});

// P3
export const BackdraftRotationPhase3 = PresetUtils.makePresetAPLRotation('P3 Backdraft', BackdraftAPLP3, {
	customCondition: player => player.getLevel() == 50,
});
export const NFRuinRotationPhase3 = PresetUtils.makePresetAPLRotation('P3 NF/Ruin', NFRuinAPLP3, {
	customCondition: player => player.getLevel() == 50,
});

// P4
export const DestroRotationPhase4 = PresetUtils.makePresetAPLRotation('P4 Destro', DestroAPLP4, {
	customCondition: player => player.getLevel() == 60,
});
export const AffRotationPhase4 = PresetUtils.makePresetAPLRotation('P4 Aff', AffAPLP4, {
	customCondition: player => player.getLevel() == 60,
});

export const APLPresets = {
	[Phase.Phase1]: [RotationDestructionPhase1],
	[Phase.Phase2]: [DestroMgiRotationPhase2, DestroConflagRotationPhase2, DemonologyRotationPhase2, AfflictionRotationPhase2],
	[Phase.Phase3]: [NFRuinRotationPhase3, BackdraftRotationPhase3],
	[Phase.Phase4]: [AffRotationPhase4, DestroRotationPhase4],
	[Phase.Phase5]: [],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {
		0: RotationDestructionPhase1,
		1: RotationDestructionPhase1,
		2: RotationDestructionPhase1,
	},
	40: {
		0: AfflictionRotationPhase2,
		1: DemonologyRotationPhase2,
		2: DestroMgiRotationPhase2,
	},
	50: {
		0: NFRuinRotationPhase3,
		2: BackdraftRotationPhase3,
	},
	60: {
		0: AffRotationPhase4,
		2: DestroRotationPhase4,
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const DestroP1Talents = {
	name: 'P1 Destruction',
	data: SavedTalents.create({ talentsString: '-03-0550201' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 25,
};

export const DestroMgiTalentsPhase2 = {
	name: 'P2 Destro Imp',
	data: SavedTalents.create({ talentsString: '-01-055020512000415' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 40,
};
export const DestroConflagTalentsPhase2 = {
	name: 'P2 Destro Conflag',
	data: SavedTalents.create({ talentsString: '--0550205120005141' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 40,
};
export const DemonologyTalentsPhase2 = {
	name: 'P2 Demonology',
	data: SavedTalents.create({ talentsString: '-2050033132501051' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 40,
};
export const AfflictionTalentsPhase2 = {
	name: 'P2 Affliction',
	data: SavedTalents.create({ talentsString: '3500253012201105--1' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 40,
};
export const BackdraftTalentsPhase3 = {
	name: 'P3 Backdraft',
	data: SavedTalents.create({ talentsString: '-032004-5050205102005151' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 50,
};
export const NFRuinTalentsPhase3 = {
	name: 'P3 NF/Ruin',
	data: SavedTalents.create({ talentsString: '25002500102-03-50502051020001' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 50,
};
export const AffTalentsPhase3 = {
	name: 'P4 Aff',
	data: SavedTalents.create({ talentsString: '4500253012201005--50502051020001' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 60,
};
export const DestroTalentsPhase3 = {
	name: 'P4 Destro',
	data: SavedTalents.create({ talentsString: '05002-035-5250205122005151' }),
	enableWhen: (player: Player<any>) => player.getLevel() == 60,
};

export const TalentPresets = {
	[Phase.Phase1]: [DestroP1Talents],
	[Phase.Phase2]: [DestroMgiTalentsPhase2, DestroConflagTalentsPhase2, DemonologyTalentsPhase2, AfflictionTalentsPhase2],
	[Phase.Phase3]: [NFRuinTalentsPhase3, BackdraftTalentsPhase3],
	[Phase.Phase4]: [AffTalentsPhase3, DestroTalentsPhase3],
	[Phase.Phase5]: [],
};

export const DefaultTalentsAff = TalentPresets[Phase.Phase4][0];
export const DefaultTalentsDestro = TalentPresets[Phase.Phase4][1];

export const DefaultTalents = DefaultTalentsDestro;

export const PresetBuildAff = PresetUtils.makePresetBuild('Aff', DefaultGearAff, DefaultTalentsAff, DefaultAPLs[60][0]);
export const PresetBuildDestro = PresetUtils.makePresetBuild('Destro', DefaultGearDestro, DefaultTalentsDestro, DefaultAPLs[60][2]);

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.FelArmor,
	summon: Summon.Imp,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	alcohol: Alcohol.AlcoholKreegsStoutBeatdown,
	defaultPotion: Potions.MajorManaPotion,
	defaultConjured: Conjured.ConjuredDemonicRune,
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	flask: Flask.FlaskOfSupremePower,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	food: Food.FoodTenderWolfSteak,
	mainHandImbue: WeaponImbue.WizardOil,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	shadowPowerBuff: ShadowPowerBuff.ElixirOfShadowPower,
	zanzaBuff: ZanzaBuff.GizzardGum,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	demonicPact: 80,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	vampiricTouch: 300,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	dragonslayerBuff: DragonslayerBuff.RallyingCryofTheDragonslayer,
	mightOfStormwind: true,
	moldarsMoxie: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	faerieFire: true,
	homunculi: 100,
	improvedFaerieFire: true,
	improvedScorch: true,
	judgementOfWisdom: true,
	markOfChaos: true,
	occultPoison: true,
	shadowWeaving: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
