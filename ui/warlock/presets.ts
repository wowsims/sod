import { Phase } from '../core/constants/other.js';
import {
	AgilityElixir,
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
	SpellPowerBuff,
	TristateEffect,
	WeaponImbue
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import {
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions as WarlockOptions,
	WarlockOptions_WeaponImbue as WarlockWeaponImbue,
} from '../core/proto/warlock.js';
import * as PresetUtils from '../core/preset_utils.js';
import { Player } from '../core/player.js';

// gear
// P1
import DestructionGear from './gear_sets/p1/destruction.gear.json';
// P2
import FireImpGear from './gear_sets/p2/fire.imp.gear.json';
import FireSuccubusGear from './gear_sets/p2/fire.succubus.gear.json';
import ShadowGear from './gear_sets/p2/shadow.gear.json';

// apls
// P1
import DestroP1APL from './apls/p1/destruction.apl.json';
// P2
import DestroMgiAPL from './apls/p2/fire.imp.apl.json';
import DestroConflagAPL from './apls/p2/fire.conflag.apl.json';
import DemonologyAPL from './apls/p2/demonology.apl.json';
import AfflictionAPL from './apls/p2/affliction.apl.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearDestructionPhase1 = PresetUtils.makePresetGear('Destruction', DestructionGear, { customCondition: (player) => player.getLevel() == 25 });

export const FireImpGearPreset = PresetUtils.makePresetGear('Fire Imp', FireImpGear, { customCondition: (player) => player.getLevel() == 40 });
export const FireSuccubusGearPreset = PresetUtils.makePresetGear('Fire Succubus', FireSuccubusGear, { customCondition: (player) => player.getLevel() == 40 });
export const ShadowGearPreset = PresetUtils.makePresetGear('Shadow', ShadowGear, { customCondition: (player) => player.getLevel() == 40 });

export const GearPresets = {
  	[Phase.Phase1]: [
		GearDestructionPhase1,
	],
	[Phase.Phase2]: [
		FireImpGearPreset,
		FireSuccubusGearPreset,
		ShadowGearPreset,
	]
};

export const DefaultGear = FireImpGearPreset;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const RotationDestructionPhase1 = PresetUtils.makePresetAPLRotation('Destruction', DestroP1APL, { customCondition: (player) => player.getLevel() == 25 });

export const DestroMgiRotationPreset = PresetUtils.makePresetAPLRotation('Destro Imp', DestroMgiAPL, { customCondition: (player) => player.getLevel() == 40 });
export const DestroConflagRotationPreset = PresetUtils.makePresetAPLRotation('Destro Conflag', DestroConflagAPL, { customCondition: (player) => player.getLevel() == 40 });
export const DemonologyRotationPreset = PresetUtils.makePresetAPLRotation('Demonology', DemonologyAPL, { customCondition: (player) => player.getLevel() == 40 });
export const AfflictionRotationPreset = PresetUtils.makePresetAPLRotation('Affliction', AfflictionAPL, { customCondition: (player) => player.getLevel() == 40 });

export const APLPresets = {
  	[Phase.Phase1]: [
		RotationDestructionPhase1,
	],
	[Phase.Phase2]: [
		DestroMgiRotationPreset,
		DestroConflagRotationPreset,
		DemonologyRotationPreset,
		AfflictionRotationPreset
	]
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  	25: {
		0: RotationDestructionPhase1,
		1: RotationDestructionPhase1,
		2: RotationDestructionPhase1,
	},
  	40: {
		0: AfflictionRotationPreset,
		1: DemonologyRotationPreset,
		2: DestroMgiRotationPreset,
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const DestroP1Talents = {name: 'Destruction P1', data: SavedTalents.create({ talentsString: '-03-0550201' }), enableWhen: (player:Player<any>) => player.getLevel() == 25};

export const DestroMgiTalents = {name: 'Destro Imp', data: SavedTalents.create({ talentsString: '-01-055020512000415' }), enableWhen: (player:Player<any>) => player.getLevel() == 40};
export const DestroConflagTalents = {name: 'Destro Conflag', data: SavedTalents.create({ talentsString: '--0550205120005141' }), enableWhen: (player:Player<any>) => player.getLevel() == 40};
export const DemonologyTalents = {name: 'Demonology', data: SavedTalents.create({ talentsString: '-2050033132501051' }), enableWhen: (player:Player<any>) => player.getLevel() == 40};
export const AfflictionTalents = {name: 'Affliction', data: SavedTalents.create({ talentsString: '3500253012201105--1' }), enableWhen: (player:Player<any>) => player.getLevel() == 40};

export const TalentPresets = {
  	[Phase.Phase1]: [
    	DestroP1Talents,
  	],
  	[Phase.Phase2]: [
		DestroMgiTalents,
		DestroConflagTalents,
		DemonologyTalents,
		AfflictionTalents
 	]
};

export const DefaultTalents = DestroMgiTalents;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Imp,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfAgility,
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	fillerExplosive: Explosive.ExplosiveEzThroRadiationBomb,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.LesserWizardOil,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	bloodPact: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	trueshotAura: true,
	thorns: TristateEffect.TristateEffectRegular,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	saygesFortune: SaygesFortune.SaygesDamage,
	sparkOfInspiration: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	improvedScorch: true,
	shadowWeaving: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
