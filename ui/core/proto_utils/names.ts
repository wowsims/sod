import { ResourceType } from '../proto/api.js';
import { ArmorType, Class, ItemSlot, ItemType, Profession, PseudoStat, Race, RangedWeaponType, Stat, WeaponType } from '../proto/common.js';
import { RepLevel, SourceFilterOption } from '../proto/ui.js';

export const armorTypeNames: Map<ArmorType, string> = new Map([
	[ArmorType.ArmorTypeUnknown, 'Unknown'],
	[ArmorType.ArmorTypeCloth, 'Cloth'],
	[ArmorType.ArmorTypeLeather, 'Leather'],
	[ArmorType.ArmorTypeMail, 'Mail'],
	[ArmorType.ArmorTypePlate, 'Plate'],
]);

export const weaponTypeNames: Map<WeaponType, string> = new Map([
	[WeaponType.WeaponTypeUnknown, 'Unknown'],
	[WeaponType.WeaponTypeAxe, 'Axe'],
	[WeaponType.WeaponTypeDagger, 'Dagger'],
	[WeaponType.WeaponTypeFist, 'Fist'],
	[WeaponType.WeaponTypeMace, 'Mace'],
	[WeaponType.WeaponTypeOffHand, 'Misc'],
	[WeaponType.WeaponTypePolearm, 'Polearm'],
	[WeaponType.WeaponTypeShield, 'Shield'],
	[WeaponType.WeaponTypeStaff, 'Staff'],
	[WeaponType.WeaponTypeSword, 'Sword'],
]);

export const rangedWeaponTypeNames: Map<RangedWeaponType, string> = new Map([
	[RangedWeaponType.RangedWeaponTypeUnknown, 'Unknown'],
	[RangedWeaponType.RangedWeaponTypeBow, 'Bow'],
	[RangedWeaponType.RangedWeaponTypeCrossbow, 'Crossbow'],
	[RangedWeaponType.RangedWeaponTypeGun, 'Gun'],
	[RangedWeaponType.RangedWeaponTypeIdol, 'Idol'],
	[RangedWeaponType.RangedWeaponTypeLibram, 'Libram'],
	[RangedWeaponType.RangedWeaponTypeSigil, 'Sigil'],
	[RangedWeaponType.RangedWeaponTypeThrown, 'Thrown'],
	[RangedWeaponType.RangedWeaponTypeTotem, 'Totem'],
	[RangedWeaponType.RangedWeaponTypeWand, 'Wand'],
]);

export const raceNames: Map<Race, string> = new Map([
	[Race.RaceUnknown, 'None'],
	[Race.RaceDwarf, 'Dwarf'],
	[Race.RaceGnome, 'Gnome'],
	[Race.RaceHuman, 'Human'],
	[Race.RaceNightElf, 'Night Elf'],
	[Race.RaceOrc, 'Orc'],
	[Race.RaceTauren, 'Tauren'],
	[Race.RaceTroll, 'Troll'],
	[Race.RaceUndead, 'Undead'],
]);

export function nameToRace(name: string): Race {
	const normalized = name.toLowerCase().replaceAll(' ', '');
	for (const [key, value] of raceNames) {
		if (value.toLowerCase().replaceAll(' ', '') == normalized) {
			return key;
		}
	}
	return Race.RaceUnknown;
}

export const classNames: Map<Class, string> = new Map([
	[Class.ClassUnknown, 'None'],
	[Class.ClassDruid, 'Druid'],
	[Class.ClassHunter, 'Hunter'],
	[Class.ClassMage, 'Mage'],
	[Class.ClassPaladin, 'Paladin'],
	[Class.ClassPriest, 'Priest'],
	[Class.ClassRogue, 'Rogue'],
	[Class.ClassShaman, 'Shaman'],
	[Class.ClassWarlock, 'Warlock'],
	[Class.ClassWarrior, 'Warrior'],
]);

export function nameToClass(name: string): Class {
	const lower = name.toLowerCase();
	for (const [key, value] of classNames) {
		if (value.toLowerCase().replace(/\s+/g, '') == lower) {
			return key;
		}
	}
	return Class.ClassUnknown;
}

export const professionNames: Map<Profession, string> = new Map([
	[Profession.ProfessionUnknown, 'None'],
	[Profession.Alchemy, 'Alchemy'],
	[Profession.Blacksmithing, 'Blacksmithing'],
	[Profession.Enchanting, 'Enchanting'],
	[Profession.Engineering, 'Engineering'],
	[Profession.Herbalism, 'Herbalism'],
	[Profession.Leatherworking, 'Leatherworking'],
	[Profession.Mining, 'Mining'],
	[Profession.Skinning, 'Skinning'],
	[Profession.Tailoring, 'Tailoring'],
]);

export function nameToProfession(name: string): Profession {
	const lower = name.toLowerCase();
	for (const [key, value] of professionNames) {
		if (value.toLowerCase() == lower) {
			return key;
		}
	}
	return Profession.ProfessionUnknown;
}

export const statOrder: Array<Stat> = [
	Stat.StatHealth,
	Stat.StatMana,
	Stat.StatArmor,
	Stat.StatBonusArmor,
	Stat.StatStamina,
	Stat.StatStrength,
	Stat.StatAgility,
	Stat.StatIntellect,
	Stat.StatSpirit,
	Stat.StatSpellPower,
	Stat.StatSpellDamage,
	Stat.StatArcanePower,
	Stat.StatFirePower,
	Stat.StatFrostPower,
	Stat.StatHolyPower,
	Stat.StatNaturePower,
	Stat.StatShadowPower,
	Stat.StatSpellHit,
	Stat.StatSpellCrit,
	Stat.StatSpellHaste,
	Stat.StatSpellPenetration,
	Stat.StatMP5,
	Stat.StatAttackPower,
	Stat.StatRangedAttackPower,
	Stat.StatFeralAttackPower,
	Stat.StatMeleeHit,
	Stat.StatMeleeCrit,
	Stat.StatMeleeHaste,
	Stat.StatArmorPenetration,
	Stat.StatExpertise,
	Stat.StatTimeworn,
	Stat.StatEnergy,
	Stat.StatRage,
	Stat.StatDefense,
	Stat.StatBlock,
	Stat.StatBlockValue,
	Stat.StatDodge,
	Stat.StatParry,
	Stat.StatResilience,
	Stat.StatArcaneResistance,
	Stat.StatFireResistance,
	Stat.StatFrostResistance,
	Stat.StatNatureResistance,
	Stat.StatShadowResistance,
];

export const statNames: Map<Stat, string> = new Map([
	[Stat.StatStrength, 'Strength'],
	[Stat.StatAgility, 'Agility'],
	[Stat.StatStamina, 'Stamina'],
	[Stat.StatIntellect, 'Intellect'],
	[Stat.StatSpirit, 'Spirit'],
	[Stat.StatSpellPower, 'Spell Power'],
	[Stat.StatSpellDamage, 'Spell Damage'],
	[Stat.StatArcanePower, 'Arcane Damage'],
	[Stat.StatFirePower, 'Fire Damage'],
	[Stat.StatFrostPower, 'Frost Damage'],
	[Stat.StatHolyPower, 'Holy Damage'],
	[Stat.StatNaturePower, 'Nature Damage'],
	[Stat.StatShadowPower, 'Shadow Damage'],
	[Stat.StatMP5, 'MP5'],
	[Stat.StatSpellHit, 'Spell Hit'],
	[Stat.StatSpellCrit, 'Spell Crit'],
	[Stat.StatSpellHaste, 'Spell Haste'],
	[Stat.StatSpellPenetration, 'Spell Pen'],
	[Stat.StatAttackPower, 'Attack Power'],
	[Stat.StatFeralAttackPower, 'Feral AP'],
	[Stat.StatMeleeHit, 'Melee Hit'],
	[Stat.StatMeleeCrit, 'Melee Crit'],
	[Stat.StatMeleeHaste, 'Melee Speed'],
	[Stat.StatArmorPenetration, 'Armor Pen'],
	[Stat.StatExpertise, 'Expertise'],
	[Stat.StatTimeworn, 'Timeworn'],
	[Stat.StatMana, 'Mana'],
	[Stat.StatEnergy, 'Energy'],
	[Stat.StatRage, 'Rage'],
	[Stat.StatArmor, 'Armor'],
	[Stat.StatRangedAttackPower, 'Ranged AP'],
	[Stat.StatDefense, 'Defense'],
	[Stat.StatBlock, 'Block'],
	[Stat.StatBlockValue, 'Block Value'],
	[Stat.StatDodge, 'Dodge'],
	[Stat.StatParry, 'Parry'],
	[Stat.StatResilience, 'Resilience'],
	[Stat.StatHealth, 'Health'],
	[Stat.StatArcaneResistance, 'Arcane Resistance'],
	[Stat.StatFireResistance, 'Fire Resistance'],
	[Stat.StatFrostResistance, 'Frost Resistance'],
	[Stat.StatNatureResistance, 'Nature Resistance'],
	[Stat.StatShadowResistance, 'Shadow Resistance'],
	[Stat.StatBonusArmor, 'Bonus Armor'],
]);

export const pseudoStatOrder: Array<PseudoStat> = [
	PseudoStat.PseudoStatMainHandDps,
	PseudoStat.PseudoStatOffHandDps,
	PseudoStat.PseudoStatRangedDps,
	PseudoStat.PseudoStatBlockValueMultiplier,
];
export const pseudoStatNames: Map<PseudoStat, string> = new Map([
	[PseudoStat.PseudoStatMainHandDps, 'Main Hand DPS'],
	[PseudoStat.PseudoStatOffHandDps, 'Off Hand DPS'],
	[PseudoStat.PseudoStatRangedDps, 'Ranged DPS'],
	[PseudoStat.PseudoStatBlockValueMultiplier, 'Block Value Multiplier'],
]);

export function getClassStatName(stat: Stat, playerClass: Class): string {
	const statName = statNames.get(stat);
	if (!statName) return 'UnknownStat';
	if (playerClass == Class.ClassHunter) {
		return statName.replace('Melee', 'Physical');
	} else {
		return statName;
	}
}

// TODO: Make sure BE exports the spell schools properly
export enum SpellSchool {
	None = 0,
	Physical = 1 << 1,
	Arcane = 1 << 2,
	Fire = 1 << 3,
	Frost = 1 << 4,
	Holy = 1 << 5,
	Nature = 1 << 6,
	Shadow = 1 << 7,
}

export const spellSchoolNames: Map<number, string> = new Map([
	[SpellSchool.Physical, 'Physical'],
	[SpellSchool.Arcane, 'Arcane'],
	[SpellSchool.Fire, 'Fire'],
	[SpellSchool.Frost, 'Frost'],
	[SpellSchool.Holy, 'Holy'],
	[SpellSchool.Nature, 'Nature'],
	[SpellSchool.Shadow, 'Shadow'],
	[SpellSchool.Nature + SpellSchool.Arcane, 'Astral'],
	[SpellSchool.Shadow + SpellSchool.Fire, 'Shadowflame'],
	[SpellSchool.Fire + SpellSchool.Arcane, 'Spellfire'],
	[SpellSchool.Arcane + SpellSchool.Frost, 'Spellfrost'],
	[SpellSchool.Frost + SpellSchool.Fire, 'Frostfire'],
	[SpellSchool.Shadow + SpellSchool.Frost, 'Shadowfrost'],
	[SpellSchool.Arcane + SpellSchool.Fire + SpellSchool.Frost, 'Chimeric'],
]);

export const itemTypeNames: Map<ItemType, string> = new Map([
	[ItemType.ItemTypeHead, 'Helm'],
	[ItemType.ItemTypeNeck, 'Neck'],
	[ItemType.ItemTypeShoulder, 'Shoulders'],
	[ItemType.ItemTypeBack, 'Cloak'],
	[ItemType.ItemTypeChest, 'Chest'],
	[ItemType.ItemTypeWrist, 'Bracers'],
	[ItemType.ItemTypeHands, 'Gloves'],
	[ItemType.ItemTypeWaist, 'Belt'],
	[ItemType.ItemTypeLegs, 'Pants'],
	[ItemType.ItemTypeFeet, 'Boots'],
	[ItemType.ItemTypeFinger, 'Ring'],
	[ItemType.ItemTypeTrinket, 'Trinket'],
	[ItemType.ItemTypeWeapon, 'Weapon'],
	[ItemType.ItemTypeRanged, 'Ranged'],
]);

export const slotNames: Map<ItemSlot, string> = new Map([
	[ItemSlot.ItemSlotHead, 'Head'],
	[ItemSlot.ItemSlotNeck, 'Neck'],
	[ItemSlot.ItemSlotShoulder, 'Shoulders'],
	[ItemSlot.ItemSlotBack, 'Back'],
	[ItemSlot.ItemSlotChest, 'Chest'],
	[ItemSlot.ItemSlotWrist, 'Wrist'],
	[ItemSlot.ItemSlotHands, 'Hands'],
	[ItemSlot.ItemSlotWaist, 'Waist'],
	[ItemSlot.ItemSlotLegs, 'Legs'],
	[ItemSlot.ItemSlotFeet, 'Feet'],
	[ItemSlot.ItemSlotFinger1, 'Finger 1'],
	[ItemSlot.ItemSlotFinger2, 'Finger 2'],
	[ItemSlot.ItemSlotTrinket1, 'Trinket 1'],
	[ItemSlot.ItemSlotTrinket2, 'Trinket 2'],
	[ItemSlot.ItemSlotMainHand, 'Main Hand'],
	[ItemSlot.ItemSlotOffHand, 'Off Hand'],
	[ItemSlot.ItemSlotRanged, 'Ranged'],
]);

export const resourceNames: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, 'None'],
	[ResourceType.ResourceTypeHealth, 'Health'],
	[ResourceType.ResourceTypeMana, 'Mana'],
	[ResourceType.ResourceTypeEnergy, 'Energy'],
	[ResourceType.ResourceTypeRage, 'Rage'],
	[ResourceType.ResourceTypeComboPoints, 'Combo Points'],
	[ResourceType.ResourceTypeFocus, 'Focus'],
]);

export const resourceColors: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, '#ffffff'],
	[ResourceType.ResourceTypeHealth, '#22ba00'],
	[ResourceType.ResourceTypeMana, '#2e93fa'],
	[ResourceType.ResourceTypeEnergy, '#ffd700'],
	[ResourceType.ResourceTypeRage, '#ff0000'],
	[ResourceType.ResourceTypeComboPoints, '#ffa07a'],
	[ResourceType.ResourceTypeFocus, '#cd853f'],
]);

export function stringToResourceType(str: string): ResourceType {
	for (const [key, val] of resourceNames) {
		if (val.toLowerCase() == str.toLowerCase()) {
			return key;
		}
	}
	return ResourceType.ResourceTypeNone;
}

export const sourceNames: Map<SourceFilterOption, string> = new Map([
	[SourceFilterOption.SourceUnknown, 'Unknown'],
	[SourceFilterOption.SourceCrafting, 'Crafting'],
	[SourceFilterOption.SourceQuest, 'Quest'],
	[SourceFilterOption.SourceReputation, 'Reputation'],
	[SourceFilterOption.SourceDungeon, 'Dungeon'],
	[SourceFilterOption.SourceRaid, 'Raid'],
	// [SourceFilterOption.SourceWorldBoss, 'World Boss'],
	[SourceFilterOption.SourceWorldBOE, 'World Drops'],
]);

// export const raidNames: Map<RaidFilterOption, string> = new Map([
// 	[RaidFilterOption.RaidUnknown, 'Unknown'],
// 	[RaidFilterOption.RaidVanilla, 'Vanilla'],
// 	[RaidFilterOption.RaidTbc, 'TBC'],
// 	[RaidFilterOption.RaidNaxxramas, 'Naxxramas'],
// 	[RaidFilterOption.RaidEyeOfEternity, 'Eye of Eternity'],
// 	[RaidFilterOption.RaidObsidianSanctum, 'Obsidian Sanctum'],
// 	[RaidFilterOption.RaidVaultOfArchavon, 'Vault of Archavon'],
// 	[RaidFilterOption.RaidUlduar, 'Ulduar'],
// 	[RaidFilterOption.RaidTrialOfTheCrusader, 'Trial of the Crusader'],
// 	[RaidFilterOption.RaidOnyxiasLair, 'Onyxia\'s Lair'],
// 	[RaidFilterOption.RaidIcecrownCitadel, 'Icecrown Citadel'],
// 	[RaidFilterOption.RaidRubySanctum, 'Ruby Sanctum'],
// ]);

export const REP_LEVEL_NAMES: Record<RepLevel, string> = {
	[RepLevel.RepLevelUnknown]: 'Unknown',
	[RepLevel.RepLevelHated]: 'Hated',
	[RepLevel.RepLevelHostile]: 'Hostile',
	[RepLevel.RepLevelUnfriendly]: 'Unfriendly',
	[RepLevel.RepLevelNeutral]: 'Neutral',
	[RepLevel.RepLevelFriendly]: 'Friendly',
	[RepLevel.RepLevelHonored]: 'Honored',
	[RepLevel.RepLevelRevered]: 'Revered',
	[RepLevel.RepLevelExalted]: 'Exalted',
};
