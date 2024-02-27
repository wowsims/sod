import { Phase } from '../core/constants/other.js';
import {
	Consumes,
	Flask,
	Food,
	Potions,
	AgilityElixir,
	RaidBuffs,
	EnchantedSigil,
	IndividualBuffs,
	Debuffs,
	StrengthBuff,
	WeaponImbue,
	TristateEffect,
	SaygesFortune,
	Profession,
	SpellPowerBuff,
	Explosive,
} from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';

import {
	PaladinAura,
	PaladinSeal,
	RetributionPaladin_Options as RetributionPaladinOptions,
} from '../core/proto/paladin.js';

import * as PresetUtils from '../core/preset_utils.js';


// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////
import Phase1RetGearJson from './gear_sets/p1ret.gear.json';
import Phase2RetSoCGearJson from './gear_sets/p2retsoc.gear.json'
import Phase2RetSoMGearJson from './gear_sets/p2retsom.gear.json'

// export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const Phase1RetGear = PresetUtils.makePresetGear('P1', Phase1RetGearJson);
export const Phase2RetSoCGear = PresetUtils.makePresetGear('P2 SoC/DS', Phase2RetSoCGearJson);
export const Phase2RetSoMGear = PresetUtils.makePresetGear('P2 SoM', Phase2RetSoMGearJson)


export const GearPresets = {
  [Phase.Phase1]: [
    Phase1RetGear,
  ],
  [Phase.Phase2]: [
	Phase2RetSoCGear,
	Phase2RetSoMGear,
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultGear = GearPresets[Phase.Phase2][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

import APLP1RetJson from './apls/p1ret.apl.json';
import APLP2RetJson from './apls/p2ret.apl.json';

export const APLP1Ret = PresetUtils.makePresetAPLRotation('P1 Ret', APLP1RetJson);
export const APLP2Ret = PresetUtils.makePresetAPLRotation('P2 Ret/Shockadin', APLP2RetJson);

export const APLPresets = {
  [Phase.Phase1]: [
    APLP1Ret,
  ],
  [Phase.Phase2]: [
	APLP2Ret,
  ]
};

// TODO: Add Phase 2 preset and pull from map
export const DefaultAPLs: Record<number, PresetUtils.PresetRotation> = {
  25: APLPresets[Phase.Phase1][0],
  40: APLPresets[Phase.Phase2][0],
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
	})
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
}

export const TalentPresets = {
  [Phase.Phase1]: [
    P1RetTalents,
  ],
  [Phase.Phase2]: [
	P2RetTalents,
	P2ShockadinTalents
  ]
};

export const DefaultTalents = TalentPresets[Phase.Phase2][0];

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RetributionPaladinOptions.create({
	aura: PaladinAura.RetributionAura,
	primarySeal: PaladinSeal.Command,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskUnknown,
	food: Food.FoodDragonbreathChili,
	defaultPotion: Potions.GreaterManaPotion,
	mainHandImbue: WeaponImbue.WildStrikes,
	agilityElixir: AgilityElixir.ElixirOfLesserAgility,
	strengthBuff: StrengthBuff.ElixirOfOgresStrength,
	spellPowerBuff: SpellPowerBuff.LesserArcaneElixir,
	fillerExplosive: Explosive.ExplosiveEzThroRadiationBomb,
	enchantedSigil: EnchantedSigil.InnovationSigil,
	boglingRoot: false,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfMight: TristateEffect.TristateEffectImproved,
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	sparkOfInspiration: true,
	saygesFortune: SaygesFortune.SaygesDamage
  });

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	battleShout: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	trueshotAura: true,
	sanctityAura: true,
})

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	faerieFire: true,
	homunculi: 70, // 70% average uptime default
	sunderArmor: true,
	judgementOfWisdom: true
})

export const OtherDefaults = {
	profession1: Profession.Blacksmithing,
	profession2: Profession.Enchanting,
  };
  