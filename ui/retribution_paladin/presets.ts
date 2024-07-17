import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Explosive,
	Food,
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
import { PaladinAura, PaladinSeal, RetributionPaladin_Options as RetributionPaladinOptions } from '../core/proto/paladin.js';
import { SavedTalents } from '../core/proto/ui.js';
import APLP1RetJson from './apls/p1ret.apl.json';
import APLP2RetJson from './apls/p2ret.apl.json';
import APLP3RetJson from './apls/p3ret.apl.json';
import APLP4RetJson from './apls/p3ret.apl.json'; // TODO: Phase 4
import Phase1RetGearJson from './gear_sets/p1ret.gear.json';
import Phase2RetSoCGearJson from './gear_sets/p2retsoc.gear.json';
import Phase2RetSoMGearJson from './gear_sets/p2retsom.gear.json';
import Phase3RetSoMGearJson from './gear_sets/p3retsom.gear.json';
import Phase4RetGearJson from './gear_sets/p4rettwist.gear.json'; // TODO: Phase 4

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const Phase1RetGear = PresetUtils.makePresetGear('P1', Phase1RetGearJson);
export const Phase2RetSoCGear = PresetUtils.makePresetGear('P2 SoC/DS', Phase2RetSoCGearJson);
export const Phase2RetSoMGear = PresetUtils.makePresetGear('P2 SoM', Phase2RetSoMGearJson);
export const Phase3RetSoMGear = PresetUtils.makePresetGear('P3 SoM', Phase3RetSoMGearJson);
export const Phase4RetGear = PresetUtils.makePresetGear('P4 Placeholder', Phase4RetGearJson);

export const GearPresets = {
	[Phase.Phase1]: [Phase1RetGear],
	[Phase.Phase2]: [Phase2RetSoCGear, Phase2RetSoMGear],
	[Phase.Phase3]: [Phase3RetSoMGear],
	[Phase.Phase4]: [Phase4RetGear],
	[Phase.Phase5]: [],
};

export const DefaultGear = GearPresets[Phase.Phase4][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLP1Ret = PresetUtils.makePresetAPLRotation('P1 Ret', APLP1RetJson);
export const APLP2Ret = PresetUtils.makePresetAPLRotation('P2 Ret/Shockadin', APLP2RetJson);
export const APLP3Ret = PresetUtils.makePresetAPLRotation('P3 Ret/Shockadin', APLP3RetJson);
export const APLP4Ret = PresetUtils.makePresetAPLRotation('P4 Ret/Placeholder', APLP4RetJson); // TODO: Phase 4

export const APLPresets = {
	[Phase.Phase1]: [APLP1Ret],
	[Phase.Phase2]: [APLP2Ret],
	[Phase.Phase3]: [APLP3Ret],
	[Phase.Phase4]: [APLP4Ret],
	[Phase.Phase5]: [],
};

export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
	25: APLPresets[Phase.Phase1][0],
	40: APLPresets[Phase.Phase2][0],
	50: APLPresets[Phase.Phase3][0],
	60: APLPresets[Phase.Phase4][0], // TODO: Phase 4
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const P1RetTalents = {
	name: 'P1 Ret',
	data: SavedTalents.create({
		talentsString: '--05230051',
	}),
};

export const P2RetTalents = {
	name: 'P2 Ret',
	data: SavedTalents.create({
		talentsString: '--532300512003151',
	}),
};

export const P2ShockadinTalents = {
	name: 'P2 Shockadin',
	data: SavedTalents.create({
		talentsString: '55050100521151--',
	}),
};

export const P3RetTalents = {
	name: 'P3 Ret',
	data: SavedTalents.create({
		talentsString: '500501--53230051200315',
	}),
};

export const P4RetTalents = {
	name: 'P4 Ret PlaceHolder',
	data: SavedTalents.create({
		talentsString: '500501--53230051200315',
	}),
};

export const TalentPresets = {
	[Phase.Phase1]: [P1RetTalents],
	[Phase.Phase2]: [P2RetTalents, P2ShockadinTalents],
	[Phase.Phase3]: [P3RetTalents],
	[Phase.Phase4]: [P4RetTalents],
	[Phase.Phase5]: [],
};

// TODO: Phase 3
export const DefaultTalents = TalentPresets[Phase.Phase3][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	primarySeal: PaladinSeal.Martyrdom,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	boglingRoot: false,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	enchantedSigil: EnchantedSigil.LivingDreamsSigil,
	fillerExplosive: Explosive.ExplosiveEzThroRadiationBomb,
	food: Food.FoodBlessSunfruit,
	mainHandImbue: WeaponImbue.WildStrikes,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfGiants,
	zanzaBuff: ZanzaBuff.AtalaiMojoOfWar,
	attackPowerBuff: AttackPowerBuff.WinterfallFirewater,
	defaultConjured: Conjured.ConjuredDemonicRune,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	saygesFortune: SaygesFortune.SaygesDamage,
	fervorOfTheTempleExplorer: true,
	songflowerSerenade: true,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	battleShout: TristateEffect.TristateEffectImproved,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	sanctityAura: true,
	leaderOfThePack: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	sunderArmor: true,
	judgementOfWisdom: true,
});

export const OtherDefaults = {
	profession1: Profession.Blacksmithing,
	profession2: Profession.Enchanting,
};
