import { Phase } from '../core/constants/other.js';
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

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

import AfflictionGearPhase1 from './gear_sets/p1.affi.tank.gear.json';
import DestructionGearPhase1 from './gear_sets/p1.destro.tank.gear.json';

import DemonologyGearPhase2 from './gear_sets/p2.demo.tank.gear.json';
import DestructionGearPhase2 from './gear_sets/p2.destro.tank.gear.json';

export const GearAfflictionTankPhase1 = PresetUtils.makePresetGear('P1 Affliction', AfflictionGearPhase1);
export const GearDestructionTankPhase1 = PresetUtils.makePresetGear('P1 Destruction', DestructionGearPhase1);

export const GearDemonologyTankPhase2 = PresetUtils.makePresetGear('P2 Demonology', DemonologyGearPhase2);
export const GearDestructionTankPhase2 = PresetUtils.makePresetGear('P2 Destruction', DestructionGearPhase2);

export const GearPresets = {
  	[Phase.Phase1]: [
    	GearAfflictionTankPhase1,
		GearDestructionTankPhase1,
  	],
	[Phase.Phase2]: [
		GearDemonologyTankPhase2,
		GearDestructionTankPhase2,
  	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearDemonologyTankPhase2;


///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import Phase1AfflictionAPL from './apls/p1.affi.tank.apl.json';
import Phase1DestroTankAPL from './apls/p1.destro.tank.apl.json';

import Phase2DemonologyAPL from './apls/p2.demo.tank.apl.json';
import Phase2DestroTankAPL from './apls/p2.destro.tank.apl.json';

export const APLAfflictionTankPhase1 = PresetUtils.makePresetAPLRotation('P1 Affliction', Phase1AfflictionAPL);
export const APLDestructionTankPhase1 = PresetUtils.makePresetAPLRotation('P1 Destruction', Phase1DestroTankAPL);

export const APLDemonologyTankPhase2 = PresetUtils.makePresetAPLRotation('P2 Demonology', Phase2DemonologyAPL);
export const APLDestructionTankPhase2 = PresetUtils.makePresetAPLRotation('P2 Destruction', Phase2DestroTankAPL);


export const APLPresets = {
  	[Phase.Phase1]: [
    	APLAfflictionTankPhase1,
		APLDestructionTankPhase1,
  	],
  	[Phase.Phase2]: [
		APLDemonologyTankPhase2,
		APLDestructionTankPhase2
  	]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
  	25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
	},
  	40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][1],
	}
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const TalentsAfflictionTankPhase1 = {
	name: 'P1 Affliction',
	data: SavedTalents.create({
		talentsString: '050025001-003',
	}),
};

export const TalentsDestructionTankPhase1 = {
	name: 'P1 Destruction',
	data: SavedTalents.create({
		talentsString: '-03-0550201',
	}),
};

export const TalentsDemonologyTankPhase2 = {
	name: 'P2 Demonology',
	data: SavedTalents.create({
		talentsString: '-2050033112501251',
	}),
};

export const TalentsDestructionTankPhase2 = {
	name: 'P2 Destruction',
	data: SavedTalents.create({
		talentsString: '-035-05500050025001',
	}),
};

export const TalentPresets = {
  	[Phase.Phase1]: [
    	TalentsAfflictionTankPhase1,
		TalentsDestructionTankPhase1,
  	],
  	[Phase.Phase2]: [
		TalentsDemonologyTankPhase2,
		TalentsDestructionTankPhase2
  	]
};


export const DefaultTalents = TalentsDemonologyTankPhase2;

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = WarlockOptions.create({
	armor: Armor.DemonArmor,
	summon: Summon.Succubus,
	weaponImbue: WarlockWeaponImbue.NoWeaponImbue,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfLesserAgility,
	defaultPotion: Potions.GreaterManaPotion,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	food: Food.FoodSagefishDelight,
	mainHandImbue: WeaponImbue.BlackfathomManaOil,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	bloodPact: TristateEffect.TristateEffectImproved,
	devotionAura: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	trueshotAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfSanctuary: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	saygesFortune: SaygesFortune.SaygesDamage,
	sparkOfInspiration: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
	curseOfRecklessness: true,
	demoralizingShout: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	improvedScorch: true,
	shadowWeaving: true,
});

export const OtherDefaults = {
	channelClipDelay: 150,
	distanceFromTarget: 5,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
};
