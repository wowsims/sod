import { Phase } from '../core/constants/other.js';
import { Player } from '../core/player.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	AtalAi,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Explosive,
	FirePowerBuff,
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
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	WarlockOptions as WarlockOptions,
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';
// apls
// P1
import DestroP1APL from './apls/p1/destruction.apl.json';
// P2
import AfflictionAPL from './apls/p2/affliction.apl.json';
import DemonologyAPL from './apls/p2/demonology.apl.json';
import DestroConflagAPL from './apls/p2/fire.conflag.apl.json';
import DestroMgiAPL from './apls/p2/fire.imp.apl.json';
// P3
import BackdraftAPL from './apls/p3/backdraft.apl.json';
import NFRuinAPL from './apls/p3/nf.ruin.apl.json';
// gear
// P1
import DestructionGear from './gear_sets/p1/destruction.gear.json';
// P2
import FireImpGear from './gear_sets/p2/fire.imp.gear.json';
import FireSuccubusGear from './gear_sets/p2/fire.succubus.gear.json';
import ShadowGear from './gear_sets/p2/shadow.gear.json';
// P3
import BackdraftGear from './gear_sets/p3/backdraft.gear.json';
import NFRuinGear from './gear_sets/p3/nf.ruin.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearDestructionPhase1 = PresetUtils.makePresetGear('Destruction', DestructionGear, { customCondition: player => player.getLevel() == 25 });

export const FireImpGearPhase2 = PresetUtils.makePresetGear('P2 Fire Imp', FireImpGear, { customCondition: player => player.getLevel() == 40 });
export const FireSuccubusGearPhase2 = PresetUtils.makePresetGear('P2 Fire Succubus', FireSuccubusGear, { customCondition: player => player.getLevel() == 40 });
export const ShadowGearPhase2 = PresetUtils.makePresetGear('P2 Shadow', ShadowGear, { customCondition: player => player.getLevel() == 40 });

export const BackdraftGearPhase3 = PresetUtils.makePresetGear('P3 Backdraft', BackdraftGear, { customCondition: player => player.getLevel() == 50 });
export const NFRuinGearPhase3 = PresetUtils.makePresetGear('P3 NF/Ruin', NFRuinGear, { customCondition: player => player.getLevel() == 50 });

export const GearPresets = {
	[Phase.Phase1]: [GearDestructionPhase1],
	[Phase.Phase2]: [FireImpGearPhase2, FireSuccubusGearPhase2, ShadowGearPhase2],
	[Phase.Phase3]: [NFRuinGearPhase3, BackdraftGearPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

// TODO: Phase 3
export const DefaultGear = BackdraftGearPhase3;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

// P1
export const RotationDestructionPhase1 = PresetUtils.makePresetAPLRotation('Destruction', DestroP1APL, { 
	customCondition: player => player.getLevel() == 25 
});

// P2
export const DestroMgiRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Destro Imp', DestroMgiAPL, { 
	customCondition: player => player.getLevel() == 40 
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
export const BackdraftRotationPhase3 = PresetUtils.makePresetAPLRotation('P3 Backdraft', BackdraftAPL, {
	customCondition: player => player.getLevel() == 50,
});
export const NFRuinRotationPhase3 = PresetUtils.makePresetAPLRotation('P3 NF/Ruin', NFRuinAPL, {
	customCondition: player => player.getLevel() == 50,
});

export const APLPresets = {
	[Phase.Phase1]: [RotationDestructionPhase1],
	[Phase.Phase2]: [DestroMgiRotationPhase2, DestroConflagRotationPhase2, DemonologyRotationPhase2, AfflictionRotationPhase2],
	[Phase.Phase3]: [NFRuinRotationPhase3, BackdraftRotationPhase3],
	[Phase.Phase4]: [],
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

export const TalentPresets = {
	[Phase.Phase1]: [DestroP1Talents],
	[Phase.Phase2]: [DestroMgiTalentsPhase2, DestroConflagTalentsPhase2, DemonologyTalentsPhase2, AfflictionTalentsPhase2],
	[Phase.Phase3]: [NFRuinTalentsPhase3, BackdraftTalentsPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

export const DefaultTalents = BackdraftTalentsPhase3;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Imp,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.SuperiorManaPotion,
	defaultAtalAi: AtalAi.AtalAiForbiddenMagic,
	defaultConjured: Conjured.ConjuredDemonicRune,
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.LesserWizardOil,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	shadowPowerBuff: ShadowPowerBuff.ElixirOfShadowPower,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	saygesFortune: SaygesFortune.SaygesDamage,
	fervorOfTheTempleExplorer: true,
	songflowerSerenade: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElementsNew: TristateEffect.TristateEffectImproved,
	curseOfShadowNew: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	homunculi: 100,
	improvedScorch: true,
	shadowWeaving: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
