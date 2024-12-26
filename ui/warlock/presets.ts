import { Phase } from '../core/constants/other.js';
import { Player } from '../core/player.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	Alcohol,
	Conjured,
	Consumes,
	Debuffs,
	EnchantedSigil,
	FirePowerBuff,
	Flask,
	Food,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	Race,
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
import DestroAplP1JSON from './apls/p1/destruction.apl.json';
import AfflictionAplP2JSON from './apls/p2/affliction.apl.json';
import DemonologyAplP2JSON from './apls/p2/demonology.apl.json';
import DestroConflagAplP2JSON from './apls/p2/fire.conflag.apl.json';
import DestroMgiAplP2JSON from './apls/p2/fire.imp.apl.json';
import BackdraftAplP3JSON from './apls/p3/backdraft.apl.json';
import NFRuinAplP3JSON from './apls/p3/nf.ruin.apl.json';
import AffAplP4JSON from './apls/p4/affliction.apl.json';
import DestroAplP4JSON from './apls/p4/destruction.apl.json';
import AffAplP5JSON from './apls/p5/affliction.apl.json';
import DemoAplP5JSON from './apls/p5/demonology.apl.json';
import DestroAplP5JSON from './apls/p5/destruction.apl.json';
import AffAplP6JSON from './apls/p6/affliction.apl.json';
import DemoAplP6JSON from './apls/p6/demonology.apl.json';
import DestroAplP6JSON from './apls/p6/destruction.apl.json';
// gear
import DestructionGearJSON from './gear_sets/p1/destruction.gear.json';
import FireImpGearJSON from './gear_sets/p2/fire.imp.gear.json';
import FireSuccubusGearJSON from './gear_sets/p2/fire.succubus.gear.json';
import ShadowGearJSON from './gear_sets/p2/shadow.gear.json';
import BackdraftGearP3JSON from './gear_sets/p3/backdraft.gear.json';
import NFRuinGearP3JSON from './gear_sets/p3/nf.ruin.gear.json';
import AffGearP4JSON from './gear_sets/p4/affliction.gear.json';
import DestroGearP4JSON from './gear_sets/p4/destruction.gear.json';
import AffGearP5JSON from './gear_sets/p5/affliction.gear.json';
import DemoGearP5JSON from './gear_sets/p5/demonology.gear.json';
import DestroGearP5JSON from './gear_sets/p5/destruction.gear.json';
import AffGearP6JSON from './gear_sets/p6/affliction.gear.json';
import DemoGearP6JSON from './gear_sets/p6/demonology.gear.json';
import DestroGearP6JSON from './gear_sets/p6/destruction.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearDestructionPhase1 = PresetUtils.makePresetGear('Destruction', DestructionGearJSON, { customCondition: player => player.getLevel() === 25 });

export const FireImpGearPhase2 = PresetUtils.makePresetGear('P2 Fire Imp', FireImpGearJSON, { customCondition: player => player.getLevel() === 40 });
export const FireSuccubusGearPhase2 = PresetUtils.makePresetGear('P2 Fire Succubus', FireSuccubusGearJSON, { customCondition: player => player.getLevel() === 40 });
export const ShadowGearPhase2 = PresetUtils.makePresetGear('P2 Shadow', ShadowGearJSON, { customCondition: player => player.getLevel() === 40 });

export const BackdraftGearPhase3 = PresetUtils.makePresetGear('P3 Backdraft', BackdraftGearP3JSON, { customCondition: player => player.getLevel() === 50 });
export const NFRuinGearPhase3 = PresetUtils.makePresetGear('P3 NF/Ruin', NFRuinGearP3JSON, { customCondition: player => player.getLevel() === 50 });

export const AffGearPhase4 = PresetUtils.makePresetGear('P4 Aff', AffGearP4JSON, { customCondition: player => player.getLevel() === 60 });
export const DestroGearPhase4 = PresetUtils.makePresetGear('P4 Destro', DestroGearP4JSON, { customCondition: player => player.getLevel() === 60 });

export const AffGearPhase5 = PresetUtils.makePresetGear('P5 Aff', AffGearP5JSON, { customCondition: player => player.getLevel() === 60 });
export const DemoGearPhase5 = PresetUtils.makePresetGear('P5 Demo', DemoGearP5JSON, { customCondition: player => player.getLevel() === 60 });
export const DestroGearPhase5 = PresetUtils.makePresetGear('P5 Destro', DestroGearP5JSON, { customCondition: player => player.getLevel() === 60 });

export const AffGearPhase6 = PresetUtils.makePresetGear('P6 Aff', AffGearP6JSON, { customCondition: player => player.getLevel() === 60 });
export const DemoGearPhase6 = PresetUtils.makePresetGear('P6 Demo', DemoGearP6JSON, { customCondition: player => player.getLevel() === 60 });
export const DestroGearPhase6 = PresetUtils.makePresetGear('P6 Destro', DestroGearP6JSON, { customCondition: player => player.getLevel() === 60 });

export const GearPresets = {
	[Phase.Phase1]: [GearDestructionPhase1],
	[Phase.Phase2]: [FireImpGearPhase2, FireSuccubusGearPhase2, ShadowGearPhase2],
	[Phase.Phase3]: [NFRuinGearPhase3, BackdraftGearPhase3],
	[Phase.Phase4]: [AffGearPhase4, DestroGearPhase4],
	[Phase.Phase5]: [AffGearPhase5, DemoGearPhase5, DestroGearPhase5],
	[Phase.Phase6]: [AffGearPhase6, DemoGearPhase6, DestroGearPhase6],
	[Phase.Phase7]: [],
};

export const DefaultGearAff = GearPresets[Phase.Phase6][0];
export const DefaultGearDemo = GearPresets[Phase.Phase6][1];;
export const DefaultGearDestro = GearPresets[Phase.Phase6][2];;

export const DefaultGear = DefaultGearDestro;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

// P1
export const RotationDestructionPhase1 = PresetUtils.makePresetAPLRotation('Destruction', DestroAplP1JSON, {
	customCondition: player => player.getLevel() === 25,
});

// P2
export const DestroMgiRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Destro Imp', DestroMgiAplP2JSON, {
	customCondition: player => player.getLevel() === 40,
});
export const DestroConflagRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Destro Conflag', DestroConflagAplP2JSON, {
	customCondition: player => player.getLevel() === 40,
});
export const DemonologyRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Demonology', DemonologyAplP2JSON, {
	customCondition: player => player.getLevel() === 40,
});
export const AfflictionRotationPhase2 = PresetUtils.makePresetAPLRotation('P2 Affliction', AfflictionAplP2JSON, {
	customCondition: player => player.getLevel() === 40,
});

// P3
export const BackdraftRotationPhase3 = PresetUtils.makePresetAPLRotation('P3 Backdraft', BackdraftAplP3JSON, {
	customCondition: player => player.getLevel() === 50,
});
export const NFRuinRotationPhase3 = PresetUtils.makePresetAPLRotation('P3 NF/Ruin', NFRuinAplP3JSON, {
	customCondition: player => player.getLevel() === 50,
});

// P4
export const AffRotationPhase4 = PresetUtils.makePresetAPLRotation('P4 Aff', AffAplP4JSON, {
	customCondition: player => player.getLevel() === 60,
});
export const DestroRotationPhase4 = PresetUtils.makePresetAPLRotation('P4 Destro', DestroAplP4JSON, {
	customCondition: player => player.getLevel() === 60,
});

// P5
export const AffRotationPhase5 = PresetUtils.makePresetAPLRotation('P5 Aff', AffAplP5JSON, {
	customCondition: player => player.getLevel() === 60,
});
export const DemoRotationPhase5 = PresetUtils.makePresetAPLRotation('P5 Demo', DemoAplP5JSON, {
	customCondition: player => player.getLevel() === 60,
});
export const DestroRotationPhase5 = PresetUtils.makePresetAPLRotation('P5 Destro', DestroAplP5JSON, {
	customCondition: player => player.getLevel() === 60,
});

// P5
export const AffRotationPhase6 = PresetUtils.makePresetAPLRotation('P6 Aff', AffAplP6JSON, {
	customCondition: player => player.getLevel() === 60,
});
export const DemoRotationPhase6 = PresetUtils.makePresetAPLRotation('P6 Demo', DemoAplP6JSON, {
	customCondition: player => player.getLevel() === 60,
});
export const DestroRotationPhase6 = PresetUtils.makePresetAPLRotation('P6 Destro', DestroAplP6JSON, {
	customCondition: player => player.getLevel() === 60,
});

export const APLPresets = {
	[Phase.Phase1]: [RotationDestructionPhase1],
	[Phase.Phase2]: [DestroMgiRotationPhase2, DestroConflagRotationPhase2, DemonologyRotationPhase2, AfflictionRotationPhase2],
	[Phase.Phase3]: [NFRuinRotationPhase3, BackdraftRotationPhase3],
	[Phase.Phase4]: [AffRotationPhase4, DestroRotationPhase4],
	[Phase.Phase5]: [AffRotationPhase5, DemoRotationPhase5, DestroRotationPhase5],
	[Phase.Phase6]: [AffRotationPhase6, DemoRotationPhase6, DestroRotationPhase6, ],
	[Phase.Phase7]: [],
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
		0: AffRotationPhase6,
		1: DemoRotationPhase6,
		2: DestroRotationPhase6,
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
	enableWhen: (player: Player<any>) => player.getLevel() === 25,
};

// P2
export const DestroMgiTalentsPhase2 = {
	name: 'P2 Destro Imp',
	data: SavedTalents.create({ talentsString: '-01-055020512000415' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 40,
};
export const DestroConflagTalentsPhase2 = {
	name: 'P2 Destro Conflag',
	data: SavedTalents.create({ talentsString: '--0550205120005141' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 40,
};
export const DemonologyTalentsPhase2 = {
	name: 'P2 Demonology',
	data: SavedTalents.create({ talentsString: '-2050033132501051' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 40,
};
export const AfflictionTalentsPhase2 = {
	name: 'P2 Affliction',
	data: SavedTalents.create({ talentsString: '3500253012201105--1' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 40,
};

// P3
export const BackdraftTalentsPhase3 = {
	name: 'P3 Backdraft',
	data: SavedTalents.create({ talentsString: '-032004-5050205102005151' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 50,
};
export const NFRuinTalentsPhase3 = {
	name: 'P3 NF/Ruin',
	data: SavedTalents.create({ talentsString: '25002500102-03-50502051020001' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 50,
};

// P4
export const AffTalentsPhase4_5 = {
	name: 'P4/5 Aff',
	data: SavedTalents.create({ talentsString: '4500253012201005--50502051020001' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 60,
};
export const DestroTalentsPhase4 = {
	name: 'P4 Destro',
	data: SavedTalents.create({ talentsString: '05002-035-5250205122005151' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 60,
};

// P5
export const DemoTalentsPhase5_6 = {
	name: 'P5/6 Demo',
	data: SavedTalents.create({ talentsString: '-230205015250005-52500051020001' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 60,
}
export const DestroTalentsPhase5_6 = {
	name: 'P5/6 Destro',
	data: SavedTalents.create({ talentsString: '05002-23-5550205122005151' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 60,
}

// P6
export const AffTalentsPhase6 = {
	name: 'P6 Aff',
	data: SavedTalents.create({ talentsString: '3500243212201005-2302050152001' }),
	enableWhen: (player: Player<any>) => player.getLevel() === 60,
}

export const TalentPresets = {
	[Phase.Phase1]: [DestroP1Talents],
	[Phase.Phase2]: [DestroMgiTalentsPhase2, DestroConflagTalentsPhase2, DemonologyTalentsPhase2, AfflictionTalentsPhase2],
	[Phase.Phase3]: [NFRuinTalentsPhase3, BackdraftTalentsPhase3],
	[Phase.Phase4]: [AffTalentsPhase4_5, DestroTalentsPhase4],
	[Phase.Phase5]: [AffTalentsPhase4_5, DemoTalentsPhase5_6, DestroTalentsPhase5_6],
	[Phase.Phase6]: [AffTalentsPhase6, DemoTalentsPhase5_6, DestroTalentsPhase5_6],
	[Phase.Phase7]: [],
};

export const DefaultTalentsAff = TalentPresets[Phase.Phase6][0];
export const DefaultTalentsDemo = TalentPresets[Phase.Phase6][1];
export const DefaultTalentsDestro = TalentPresets[Phase.Phase6][2];

export const DefaultTalents = DefaultTalentsDestro;

///////////////////////////////////////////////////////////////////////////
//                                 Builds
///////////////////////////////////////////////////////////////////////////

export const PresetBuildAff = PresetUtils.makePresetBuild('Affliction', {
	gear: DefaultGearAff,
	talents: DefaultTalentsAff,
	rotation: DefaultAPLs[60][0]
});
export const PresetBuildDemo = PresetUtils.makePresetBuild('Demonology', {
	gear: DefaultGearDemo,
	talents: DefaultTalentsDemo,
	rotation: DefaultAPLs[60][1]
});
export const PresetBuildDestro = PresetUtils.makePresetBuild('Destruction', {
	gear: DefaultGearDestro,
	talents: DefaultTalentsDestro,
	rotation: DefaultAPLs[60][2],
});

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
	enchantedSigil: EnchantedSigil.WrathOfTheStormSigil,
	flask: Flask.FlaskOfAncientKnowledge,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	food: Food.FoodDarkclawBisque,
	mainHandImbue: WeaponImbue.EnchantedRepellent,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	mildlyIrradiatedRejuvPot: true,
	spellPowerBuff: SpellPowerBuff.ElixirOfTheMageLord,
	shadowPowerBuff: ShadowPowerBuff.ElixirOfShadowPower,
	zanzaBuff: ZanzaBuff.SpiritOfZanza,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	aspectOfTheLion: true,
	demonicPact: 120,
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
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	mightOfStormwind: true,
	moldarsMoxie: true,
	rallyingCryOfTheDragonslayer: true,
	saygesFortune: SaygesFortune.SaygesDamage,
	slipkiksSavvy: true,
	songflowerSerenade: true,
	spiritOfZandalar: true,
	valorOfAzeroth: true,
	warchiefsBlessing: true,
});

export const DefaultDebuffs = Debuffs.create({
	faerieFire: true,
	homunculi: 100,
	improvedScorch: true,
	improvedShadowBolt: true,
	judgementOfWisdom: true,
	markOfChaos: true,
	occultPoison: true,
	shadowWeaving: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	race: Race.RaceOrc,
	channelClipDelay: 150,
};
