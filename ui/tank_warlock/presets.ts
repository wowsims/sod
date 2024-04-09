import { Player } from 'ui/core/player.js';

import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
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
	ShadowPowerBuff,
	SpellPowerBuff,
	StrengthBuff,
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
import Phase1AfflictionAPL from './apls/p1.affi.tank.apl.json';
import Phase1DestroTankAPL from './apls/p1.destro.tank.apl.json';
import Phase2DemonologyAPL from './apls/p2.demo.tank.apl.json';
import Phase2DestroTankAPL from './apls/p2.destro.tank.apl.json';
import Phase3TankAPL from './apls/p3.destro.tank.apl.json';
import AfflictionGearPhase1 from './gear_sets/p1.affi.tank.gear.json';
import DestructionGearPhase1 from './gear_sets/p1.destro.tank.gear.json';
import DemonologyGearPhase2 from './gear_sets/p2.demo.tank.gear.json';
import DestructionGearPhase2 from './gear_sets/p2.destro.tank.gear.json';
import TankGearPhase3 from './gear_sets/p3.destro.tank.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearAfflictionTankPhase1 = PresetUtils.makePresetGear('P1 Affliction', AfflictionGearPhase1, { customCondition: player => player.getLevel() == 25 });
export const GearDestructionTankPhase1 = PresetUtils.makePresetGear('P1 Destruction', DestructionGearPhase1, { customCondition: player => player.getLevel() == 25 });

export const GearDemonologyTankPhase2 = PresetUtils.makePresetGear('P2 Demonology', DemonologyGearPhase2, { customCondition: player => player.getLevel() == 40 });
export const GearDestructionTankPhase2 = PresetUtils.makePresetGear('P2 Destruction', DestructionGearPhase2, { customCondition: player => player.getLevel() == 40 });

export const GearTankPhase3 = PresetUtils.makePresetGear('P3 Destruction', TankGearPhase3, { customCondition: player => player.getLevel() == 50 });

export const GearPresets = {
	[Phase.Phase1]: [GearAfflictionTankPhase1, GearDestructionTankPhase1],
	[Phase.Phase2]: [GearDemonologyTankPhase2, GearDestructionTankPhase2],
	[Phase.Phase3]: [GearTankPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

// TODO: Phase 3
export const DefaultGear = GearTankPhase3;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLAfflictionTankPhase1 = PresetUtils.makePresetAPLRotation('P1 Affliction', Phase1AfflictionAPL, { 
	customCondition: player => player.getLevel() == 25 
});
export const APLDestructionTankPhase1 = PresetUtils.makePresetAPLRotation('P1 Destruction', Phase1DestroTankAPL, { 
	customCondition: player => player.getLevel() == 25 
});

export const APLDemonologyTankPhase2 = PresetUtils.makePresetAPLRotation('P2 Demonology', Phase2DemonologyAPL, { 
	customCondition: player => player.getLevel() == 40 
});
export const APLDestructionTankPhase2 = PresetUtils.makePresetAPLRotation('P2 Destruction', Phase2DestroTankAPL, { 
	customCondition: player => player.getLevel() == 40 
});

export const APLTankPhase3 = PresetUtils.makePresetAPLRotation('P3 Destruction', Phase3TankAPL, { 
	customCondition: player => player.getLevel() == 50 
});

export const APLPresets = {
	[Phase.Phase1]: [APLAfflictionTankPhase1, APLDestructionTankPhase1],
	[Phase.Phase2]: [APLDemonologyTankPhase2, APLDestructionTankPhase2],
	[Phase.Phase3]: [APLTankPhase3],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][1],
	},
	40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][1],
	},
	// TODO: Phase 3
	50: {
		0: APLPresets[Phase.Phase3][0],
		1: APLPresets[Phase.Phase3][0],
	},
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
	enableWhen: (player: Player<any>) => player.getLevel() == 25,
};

export const TalentsDestructionTankPhase1 = {
	name: 'P1 Destruction',
	data: SavedTalents.create({
		talentsString: '-03-0550201',
	}),
	enableWhen: (player: Player<any>) => player.getLevel() == 25,
};

export const TalentsDemonologyTankPhase2 = {
	name: 'P2 Demonology',
	data: SavedTalents.create({
		talentsString: '-2050033112501251',
	}),
	enableWhen: (player: Player<any>) => player.getLevel() == 40,
};

export const TalentsDestructionTankPhase2 = {
	name: 'P2 Destruction',
	data: SavedTalents.create({
		talentsString: '-035-05500050025001',
	}),
	enableWhen: (player: Player<any>) => player.getLevel() == 40,
};

export const TalentsTankPhase3 = {
	name: 'P3 Destruction',
	data: SavedTalents.create({
		talentsString: '05-03-505020500050515',
	}),
	enableWhen: (player: Player<any>) => player.getLevel() == 50,
};

export const TalentPresets = {
	[Phase.Phase1]: [TalentsAfflictionTankPhase1, TalentsDestructionTankPhase1],
	[Phase.Phase2]: [TalentsDemonologyTankPhase2, TalentsDestructionTankPhase2],
	[Phase.Phase3]: [TalentsTankPhase3],
	[Phase.Phase4]: [],
	[Phase.Phase5]: [],
};

// TODO: Phase 3
export const DefaultTalents = TalentsTankPhase3;

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
	defaultPotion: Potions.SuperiorManaPotion,
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	firePowerBuff: FirePowerBuff.ElixirOfFirepower,
	shadowPowerBuff: ShadowPowerBuff.ElixirOfShadowPower,
	food: Food.FoodTenderWolfSteak,
	mainHandImbue: WeaponImbue.ShadowOil,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfGiants,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	bloodPact: TristateEffect.TristateEffectImproved,
	devotionAura: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfSanctuary: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	saygesFortune: SaygesFortune.SaygesDamage,
	fervorOfTheTempleExplorer: true,
	songflowerSerenade: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElementsNew: TristateEffect.TristateEffectRegular,
	curseOfShadowNew: TristateEffect.TristateEffectRegular,
	curseOfRecklessness: true,
	demoralizingShout: TristateEffect.TristateEffectImproved,
	faerieFire: true,
	homunculi: 100,
	improvedScorch: true,
	shadowWeaving: true,
});

export const OtherDefaults = {
	channelClipDelay: 150,
	distanceFromTarget: 5,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
};
