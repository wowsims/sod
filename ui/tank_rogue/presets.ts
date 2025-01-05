import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import { Player } from '../core/proto/api';
import {
	AgilityElixir,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	Flask,
	Food,
	IndividualBuffs,
	Profession,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import { RogueOptions } from '../core/proto/rogue.js';
import { SavedTalents } from '../core/proto/ui.js';
import { playerPresets } from '../raid/presets';
import MutilateApl from './apls/mutilate.apl.json';
import P3MutilateApl from './apls/Mutilate_DPS_50.apl.json';
import MutilateIEAApl from './apls/mutilate_IEA.apl.json';
import P3ExposeMutilateApl from './apls/Mutilate_IEA_50.apl.json';
import P5SaberAPL from './apls/P5_Saber.apl.json';
import P6MGAPL from './apls/P6_MG.apl.json';
import P6SaberAPL from './apls/P6_saber.apl.json';
import P3SaberApl from './apls/Saber_DPS_50.apl.json';
import P3SaberIEAApl from './apls/Saber_IEA_50.apl.json';
import P4SaberWeaveApl from './apls/Saber_Weave_60.apl.json';
import BlankGear from './gear_sets/blank.gear.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P1DaggersGear from './gear_sets/p1_daggers.gear.json';
import P2DaggersGear from './gear_sets/p2_daggers.gear.json';
import P3DaggersGear from './gear_sets/p3_daggers.gear.json';
import P3SaberGear from './gear_sets/p3_saber.gear.json';
import P4SaberGear from './gear_sets/p4_saber.gear.json';
import P5SaberGear from './gear_sets/p5_saber.gear.json';
import P6MGGear from './gear_sets/p6_mg.gear.json';
import P6SaberGear from './gear_sets/p6_saber.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

// TODO: Add gear sets
export const GearBlank = PresetUtils.makePresetGear('Blank', BlankGear);
export const GearDaggersP1 = PresetUtils.makePresetGear('P1 Daggers', P1DaggersGear, { customCondition: player => player.getLevel() == 25 });
export const GearDaggersP2 = PresetUtils.makePresetGear('P2 Daggers', P2DaggersGear, { customCondition: player => player.getLevel() == 40 });
export const GearCombatP1 = PresetUtils.makePresetGear('P1 Combat', P1CombatGear, { customCondition: player => player.getLevel() == 25 });
export const GearDaggersP3 = PresetUtils.makePresetGear('P3 Daggers', P3DaggersGear, { customCondition: player => player.getLevel() == 50 });
export const GearSaberP3 = PresetUtils.makePresetGear('P3 Saber', P3SaberGear, { customCondition: player => player.getLevel() == 50 });
export const GearSaberP4 = PresetUtils.makePresetGear('P4 Saber', P4SaberGear, { customCondition: player => player.getLevel() == 60 });
export const GearSaberP5 = PresetUtils.makePresetGear('P5 Saber', P5SaberGear, { customCondition: player => player.getLevel() == 60 });
export const GearMGP6 = PresetUtils.makePresetGear('P6 Main Gauche', P6MGGear, { customCondition: player => player.getLevel() == 60 });
export const GearSaberP6 = PresetUtils.makePresetGear('P6 Saber', P6SaberGear, { customCondition: player => player.getLevel() == 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearDaggersP1, GearCombatP1],
	[Phase.Phase2]: [GearDaggersP2],
	[Phase.Phase3]: [GearDaggersP3, GearSaberP3],
	[Phase.Phase4]: [GearSaberP4],
	[Phase.Phase5]: [GearSaberP5],
	[Phase.Phase6]: [GearMGP6, GearSaberP6],
};

export const DefaultGear = GearPresets[Phase.Phase6][0];

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl, { customCondition: player => player.getLevel() <= 40 });
export const ROTATION_PRESET_MUTILATE_IEA = PresetUtils.makePresetAPLRotation('Mutilate IEA', MutilateIEAApl, {
	customCondition: player => player.getLevel() <= 40,
});
export const ROTATION_PRESET_MUTILATE_P3 = PresetUtils.makePresetAPLRotation('P3 Mutilate', P3MutilateApl, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_MUTILATE_IEA_P3 = PresetUtils.makePresetAPLRotation('P3 Expose Mutilate', P3ExposeMutilateApl, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_SABER_P3 = PresetUtils.makePresetAPLRotation('P3 Saber', P3SaberApl, { customCondition: player => player.getLevel() === 50 });
export const ROTATION_PRESET_SABER_IEA_P3 = PresetUtils.makePresetAPLRotation('P3 Expose Saber', P3SaberIEAApl, {
	customCondition: player => player.getLevel() === 50,
});
export const ROTATION_PRESET_SABER_WEAVE_P4 = PresetUtils.makePresetAPLRotation('P4 Saber Weave', P4SaberWeaveApl, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SABER_P5 = PresetUtils.makePresetAPLRotation('P5 Saber', P5SaberAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_MG_P6 = PresetUtils.makePresetAPLRotation('P6 Main Gauche', P6MGAPL, {
	customCondition: player => player.getLevel() === 60,
});
export const ROTATION_PRESET_SABER_P6 = PresetUtils.makePresetAPLRotation('P6 Saber', P6SaberAPL, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [ROTATION_PRESET_MUTILATE],
	[Phase.Phase2]: [ROTATION_PRESET_MUTILATE, ROTATION_PRESET_MUTILATE_IEA],
	[Phase.Phase3]: [ROTATION_PRESET_MUTILATE_P3, ROTATION_PRESET_MUTILATE_IEA_P3, ROTATION_PRESET_SABER_P3, ROTATION_PRESET_SABER_IEA_P3],
	[Phase.Phase4]: [ROTATION_PRESET_SABER_WEAVE_P4],
	[Phase.Phase5]: [ROTATION_PRESET_SABER_P5],
	[Phase.Phase6]: [ROTATION_PRESET_MG_P6, ROTATION_PRESET_SABER_P6],
};

export const DefaultAPLs: Record<number, Record<number, PresetUtils.PresetRotation>> = {
	25: {
		0: APLPresets[Phase.Phase1][0],
		1: APLPresets[Phase.Phase1][0],
		2: APLPresets[Phase.Phase1][0],
	},
	40: {
		0: APLPresets[Phase.Phase2][0],
		1: APLPresets[Phase.Phase2][0],
		2: APLPresets[Phase.Phase2][0],
	},
	50: {
		0: APLPresets[Phase.Phase3][0],
		1: APLPresets[Phase.Phase3][0],
		2: APLPresets[Phase.Phase3][0],
	},
	60: {
		0: APLPresets[Phase.Phase6][0],
		1: APLPresets[Phase.Phase6][1],
	},
};

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/classic/talent-calc and copy the numbers in the url.

export const CombatDagger25Talents = PresetUtils.makePresetTalents('P1 Combat Dagger', SavedTalents.create({ talentsString: '-023305002001' }), {
	customCondition: player => player.getLevel() === 25,
});

export const ColdBloodMutilate40Talents = PresetUtils.makePresetTalents('P2 CB Mutilate', SavedTalents.create({ talentsString: '005303103551--05' }), {
	customCondition: player => player.getLevel() === 40,
});

export const IEAMutilate40Talents = PresetUtils.makePresetTalents('P2 CB/IEA Mutilate', SavedTalents.create({ talentsString: '005303121551--05' }), {
	customCondition: player => player.getLevel() === 40,
});

export const CombatMutilate40Talents = PresetUtils.makePresetTalents('P2 AR/BF Mutilate', SavedTalents.create({ talentsString: '-0053052020550100201' }), {
	customCondition: player => player.getLevel() === 40,
});

export const TankMutilate50Talents = PresetUtils.makePresetTalents(
	'P3 HAT/CttC Mutilate',
	SavedTalents.create({ talentsString: '00532012-00532500004501001-02' }),
	{
		customCondition: player => player.getLevel() === 50,
	},
);

export const TankSaber50Talents = PresetUtils.makePresetTalents('P3 Saber Carnage', SavedTalents.create({ talentsString: '0053221-02505501000501031' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TankBladeFlurry50Talents = PresetUtils.makePresetTalents('P3 BF Poison', SavedTalents.create({ talentsString: '0053221205-02330520000501' }), {
	customCondition: player => player.getLevel() === 50,
});

export const TankSaber60Talents = PresetUtils.makePresetTalents('P4 Saber', SavedTalents.create({ talentsString: '305323102-0230550100050150131' }), {
	customCondition: player => player.getLevel() === 60,
});

export const P5TankSaberTalents = PresetUtils.makePresetTalents('P5 Saber', SavedTalents.create({ talentsString: '30532312-0230550100050140231' }), {
	customCondition: player => player.getLevel() === 60,
});

export const P6TankMGTalents = PresetUtils.makePresetTalents('P6 Main Gauche', SavedTalents.create({ talentsString: '305303125-023055010005013013' }), {
	customCondition: player => player.getLevel() === 60,
});

export const P6TankSaberTalents = PresetUtils.makePresetTalents('P6 Saber', SavedTalents.create({ talentsString: '305303105-0230550100050130231' }), {
	customCondition: player => player.getLevel() === 60,
});

export const TalentPresets = {
	[Phase.Phase1]: [CombatDagger25Talents],
	[Phase.Phase2]: [ColdBloodMutilate40Talents, IEAMutilate40Talents, CombatMutilate40Talents],
	[Phase.Phase3]: [TankMutilate50Talents, TankSaber50Talents, TankBladeFlurry50Talents],
	[Phase.Phase4]: [TankSaber60Talents],
	[Phase.Phase5]: [P5TankSaberTalents],
	[Phase.Phase6]: [P6TankMGTalents, P6TankSaberTalents],
};

export const DefaultTalentsAssassin = TalentPresets[Phase.Phase6][0];
export const DefaultTalentsCombat = TalentPresets[Phase.Phase6][0];
export const DefaultTalentsSubtlety = TalentPresets[Phase.Phase6][0];

export const DefaultTalents = P6TankMGTalents;

///////////////////////////////////////////////////////////////////////////
//                                Encounters
///////////////////////////////////////////////////////////////////////////

export const PresetBuildMG = PresetUtils.makePresetBuild('Main Gauche IEA', {
	gear: GearMGP6,
	talents: P6TankMGTalents,
	rotation: ROTATION_PRESET_MG_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Main Gauche',
		'http://localhost:5173/sod/tank_rogue/#eJztVl1oHFUU3nNndnL3xJ1Obrbp7E2sybaGcXGT2d1sTTaRnQa1aayaB8HgH1WI2hJosCDqU6wVUygafdIIJi2xjS3V/gpGGoNPyYNQiz9BUQOCtNKH4EMtVNR7Z2ZnN7vWFysIehnmnnvuOd+c33sH11NikGZiEYdshzcIjJKBcQLHxQNkiUDvOhNm3LUD5wCWAZ6CBYA+yBMDWIgfqkdtYPix54aeNqjZ06wl3wesNc4e0M2FbmvsBEfVOHJAxxrjY8F6Lyb2liZ1cy5mfX6ei8VHgvttt/XKGVMs9ovFhW7rtRNy53Wx2BuzvvoljtS4OKVbp6/EBXvOQ97vKlye0s3vdWtcLKjx9aRunX2xUVCnpnXr+N4m8elLb4F4j01LAybe1c1LkvjtqEQWxL53dPOb+kw7RQMsrVN42KeOE/JszSg0jQGOgzIF6jzAeQgtQXgZtAsAKwD5H1VcJqH/xz862Oh/PsQ/kVBvfSfMgKxMGIBB2C66bwT6m7J2Lmtn05lcys4IOmenbdsW76x4BtRB+nikbuJNORYLO0/VmkqilpU4WI81O8Oo0KvRhIYqEnu3CRhH5BQ1+uGn0KwlIliD4VS6LSe31mKEiyWd1y3wNFI5wa5DbRYUBrmE2EQllbZRdBWXyGeGfYSMQBDf8/Qvd1jE10+zzVwYwupwTTJKFUbyIamIkM5YSLkQorLfXcrtYo8n296jxPHBNvIEmqwBIwdVKfbdYR0jE+CT3J8rpOgHL/f7UpLk/sxa+E3YwGJID6rSgaO7kE6AR3FvYjHOUE/cgHgSpDdXo5aaD7Et/E7s4Xk0kjpVGXk+xF0PRTBa2cYit7/ErZDrsCW88Ws3u5dvw17mIE+aQiA6C8jEySzMO3SRe7HJuMoKU/aAH61cJdwmuxg7dgcF7mCOZzHN2rExGRea+izUCtSwMH/uAStcDgLlsG66pATbwZ/ER9hD2JCMCf3ILPgJ8QSzrmq4TDWL65JrA0mZsMBzCRGugrCLiWU38kZk3ECdBTGWFeMXFEtyC9ezpsoUrBKVyC9Ml+X8iJtzeYn4ORck92eW553Yxm4tSr0q7p1oohYjJ4umu95pZd6lA2P7+F3Yw/LX8nZ3lWrWww4+FSDJ4mlnqXIrylyUblVjBe3EHuYP4n3sHjSTDcIOnF0Viz81pMPnBEWZsStNC9CLXvq9I+/siugopdq7za6qz46gV9O8HW9mG9x29xqirIDKzcu4KjM/A7ubb8UCvx1jjFVXqAtWHeMK+I4ArJ2ncANrKXop/0xkP1dXflEhw0sd/Jctly213Fa+xT0Lrq2nVYVLXRUuGfm3By06sgCjAFcILPq3IHU+86kuR/yWTUBrxF05k06df104t3jzSiGedccPhWbtEyK6k56GxLziKYQeLSq0OfHildDlcc4VWp5wx5eFzLhCl/ZEGd029MzQcPMm2ywk/2U/PdY+57rgtB77WzgG//362CFG2n2PHHa6duxaeWnx/i8Km/0dZxD+AMAtU0s=',
	),
});

export const PresetBuildSaber = PresetUtils.makePresetBuild('Saber Slash', {
	gear: GearSaberP6,
	talents: P6TankSaberTalents,
	rotation: ROTATION_PRESET_SABER_P6,
	encounter: PresetUtils.makePresetEncounter(
		'Saber Slash',
		'http://localhost:5173/sod/tank_rogue/#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',
	),
});

///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = RogueOptions.create({
	honorAmongThievesCritRate: 100,
	pkSwap: false,
});

///////////////////////////////////////////////////////////////////////////
//                         Consumes/Buffs/Debuffs
///////////////////////////////////////////////////////////////////////////

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheHoneyBadger,
	attackPowerBuff: AttackPowerBuff.JujuMight,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfTheTitans,
	food: Food.FoodGrilledSquid,
	mainHandImbue: WeaponImbue.WildStrikes,
	offHandImbue: WeaponImbue.DeadlyPoison,
	strengthBuff: StrengthBuff.JujuPower,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	aspectOfTheLion: true,
	battleShout: TristateEffect.TristateEffectImproved,
	demonicPact: 80,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
	graceOfAirTotem: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectImproved,
	fengusFerocity: true,
	mightOfStormwind: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
	spiritOfZandalar: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	dreamstate: true,
	faerieFire: true,
	homunculi: 100,
	improvedScorch: true,
	mangle: true,
	markOfChaos: true,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
};
