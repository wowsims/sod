import { getWowheadLanguagePrefix } from '../constants/lang';
import { MAX_CHARACTER_LEVEL } from '../constants/mechanics';
import { ResourceType } from '../proto/api';
import { ActionID as ActionIdProto, ItemRandomSuffix, OtherAction } from '../proto/common';
import { IconData, UIItem as Item } from '../proto/ui';
import { buildWowheadTooltipDataset, WowheadTooltipItemParams, WowheadTooltipSpellParams } from '../wowhead';
import { Database } from './database';

// Used to filter action IDs by level
export interface ActionIdConfig {
	id: number;
	minLevel?: number;
	maxLevel?: number;
}

// Uniquely identifies a specific item / spell / thing in WoW. This object is immutable.
export class ActionId {
	readonly itemId: number;
	readonly randomSuffixId: number;
	readonly spellId: number;
	readonly otherId: OtherAction;
	readonly tag: number;
	readonly rank: number;

	readonly baseName: string; // The name without any tag additions.
	readonly name: string;
	readonly iconUrl: string;
	readonly spellIdTooltipOverride: number | null;

	private constructor(
		itemId: number,
		spellId: number,
		otherId: OtherAction,
		tag: number,
		baseName: string,
		name: string,
		iconUrl: string,
		rank: number,
		randomSuffixId?: number,
	) {
		this.itemId = itemId;
		this.randomSuffixId = randomSuffixId || 0;
		this.spellId = spellId;
		this.otherId = otherId;
		(this.rank = rank), (this.tag = tag);

		switch (otherId) {
			case OtherAction.OtherActionNone:
				break;
			case OtherAction.OtherActionWait:
				baseName = 'Wait';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_pocketwatch_01.jpg';
				break;
			case OtherAction.OtherActionManaRegen:
				name = 'Mana Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				if (tag === 1) {
					name += ' (Casting)';
				} else if (tag === 2) {
					name += ' (Not Casting)';
				}
				break;
			case OtherAction.OtherActionEnergyRegen:
				baseName = 'Energy Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeEnergy];
				break;
			case OtherAction.OtherActionComboPoints:
				baseName = 'Combo Point Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeComboPoints];
				break;
			case OtherAction.OtherActionFocusRegen:
				baseName = 'Focus Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeFocus];
				break;
			case OtherAction.OtherActionManaGain:
				baseName = 'Mana Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				break;
			case OtherAction.OtherActionRageGain:
				baseName = 'Rage Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeRage];
				break;
			case OtherAction.OtherActionAttack:
				name = 'Melee';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_04.jpg';
				if (tag === 1) {
					name += ' (Main-Hand)';
				} else if (tag === 2) {
					name += ' (Off-Hand)';
				} else if (tag === 3) {
					name += ' (Extra Attack)';
				}
				break;
			case OtherAction.OtherActionShoot:
				name = 'Shoot';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/ability_marksmanship.jpg';
				if (tag === 3) {
					name += ' (Extra Attack)';
				}
				break;
			case OtherAction.OtherActionMove:
				name = 'Move';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_boots_02.jpg';
				break;
			case OtherAction.OtherActionPet:
				break;
			case OtherAction.OtherActionRefund:
				baseName = 'Refund';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_coin_01.jpg';
				break;
			case OtherAction.OtherActionDamageTaken:
				baseName = 'Damage Taken';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_04.jpg';
				break;
			case OtherAction.OtherActionHealingModel:
				baseName = 'Incoming HPS';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_renew.jpg';
				break;
			case OtherAction.OtherActionPotion:
				baseName = 'Potion';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_alchemy_elixir_04.jpg';
				break;
			case OtherAction.OtherActionExplosives:
				baseName = 'Explosive';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/Inv_misc_bomb_06.jpg';
				break;
			case OtherAction.OtherActionOffensiveEquip:
				baseName = 'Offensive Equipment';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_trinket_naxxramas03.jpg';
				break;
			case OtherAction.OtherActionDefensiveEquip:
				baseName = 'Defensive Equipment';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_trinket_naxxramas05.jpg';
				break;
		}
		this.baseName = baseName;
		this.name = name || baseName;
		this.iconUrl = iconUrl;
		this.spellIdTooltipOverride = this.spellTooltipOverride?.spellId || null;
		if (this.name) this.name += rank ? ` (Rank ${rank})` : '';
	}

	anyId(): number {
		return this.itemId || this.spellId || this.otherId;
	}

	equals(other: ActionId): boolean {
		return this.equalsIgnoringTag(other) && this.tag === other.tag;
	}

	equalsIgnoringTag(other: ActionId): boolean {
		return this.itemId === other.itemId && this.randomSuffixId === other.randomSuffixId && this.spellId === other.spellId && this.otherId === other.otherId;
	}

	setBackground(elem: HTMLElement) {
		if (this.iconUrl) {
			elem.style.backgroundImage = `url('${this.iconUrl}')`;
		}
	}

	static makeItemUrl(id: number, randomSuffixId?: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		const url = new URL(`https://wowhead.com/classic/${langPrefix}item=${id}`);
		url.searchParams.set('level', String(MAX_CHARACTER_LEVEL));
		url.searchParams.set('rand', String(randomSuffixId || 0));
		return url.toString();
	}
	static makeSpellUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		const showBuff = spellIDsToShowBuffs.has(id);

		let url = `https://wowhead.com/classic/${langPrefix}spell=${id}`;
		if (showBuff) url = `${url}?buff=1`;

		return url;
	}
	static async makeItemTooltipData(id: number, params?: Omit<WowheadTooltipItemParams, 'itemId'>) {
		return buildWowheadTooltipDataset({ itemId: id, ...params });
	}
	static async makeSpellTooltipData(id: number, params?: Omit<WowheadTooltipSpellParams, 'spellId'>) {
		return buildWowheadTooltipDataset({ spellId: id, ...params });
	}
	static makeQuestUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		return `https://wowhead.com/classic/${langPrefix}quest=${id}`;
	}
	static makeNpcUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		return `https://wowhead.com/classic/${langPrefix}npc=${id}`;
	}
	static makeZoneUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		return `https://wowhead.com/classic/${langPrefix}zone=${id}`;
	}

	setWowheadHref(elem: HTMLAnchorElement) {
		if (this.itemId) {
			elem.href = ActionId.makeItemUrl(this.itemId, this.randomSuffixId);
		} else if (this.spellId) {
			elem.href = ActionId.makeSpellUrl(this.spellIdTooltipOverride || this.spellId);
		}
	}

	async setWowheadDataset(elem: HTMLElement, params?: Omit<WowheadTooltipItemParams, 'itemId'> | Omit<WowheadTooltipSpellParams, 'spellId'>) {
		(this.itemId
			? ActionId.makeItemTooltipData(this.itemId, params)
			: ActionId.makeSpellTooltipData(this.spellIdTooltipOverride || this.spellId, params)
		).then(url => {
			if (elem) elem.dataset.wowhead = url;
		});
	}

	setBackgroundAndHref(elem: HTMLAnchorElement) {
		this.setBackground(elem);
		this.setWowheadHref(elem);
	}

	async fillAndSet(elem: HTMLAnchorElement, setHref: boolean, setBackground: boolean): Promise<ActionId> {
		const filled = await this.fill();
		if (setHref) {
			filled.setWowheadHref(elem);
		}
		if (setBackground) {
			filled.setBackground(elem);
		}
		return filled;
	}

	// Returns an ActionId with the name and iconUrl fields filled.
	// playerIndex is the optional index of the player to whom this ID corresponds.
	async fill(playerIndex?: number): Promise<ActionId> {
		if (this.name || this.iconUrl) {
			return this;
		}

		if (this.otherId) {
			return this;
		}

		const tooltipData = await ActionId.getTooltipData(this);

		const baseName = tooltipData['name'];
		let name = baseName;
		switch (baseName) {
			case 'Arcane Blast':
				if (this.tag === 1) {
					name += ' (No Stacks)';
				} else if (this.tag === 2) {
					name += ` (1 Stack)`;
				} else if (this.tag > 2) {
					name += ` (${this.tag - 1} Stacks)`;
				}
				break;
			// Arcane Missiles hits are a separate spell and have to use a tag to differentiate from the cast
			case 'Arcane Missiles':
			// Balefire Bolt's aura uses the same spell ID as the cast
			case 'Balefire Bolt':
				break;
			case 'Berserking':
				if (this.tag !== 0) name = `${name} (${this.tag * 5}%)`;
				break;
			case 'Explosive Trap':
				if (this.tag === 1) {
					name += ' (Weaving)';
				}
				break;
			case 'Hot Streak':
				if (this.tag) name = 'Heating Up';
				break;
			// DoT then Explode Spells
			case 'Living Bomb':
			case 'Seed of Corruption':
				if (this.tag === 0) name = `${name} (DoT)`;
				else if (this.tag === 1) name = `${name} (Explosion)`;
				break;
			// Burn Spells
			case 'Fireball':
			case 'Frostfire Bolt':
			case 'Pyroblast':
			case 'Flame Shock':
				if (this.tag === 1) name = `${name} (DoT)`;
				break;
			// Channeled Tick Spells
			case 'Evocation':
			case 'Mind Flay':
			case 'Mind Sear':
				if (this.tag > 0) name = `${name} (${this.tag} Tick)`;
				break;
			case 'Mind Spike':
				if (this.tag === 1) name = `${name} (2pT2.5)`;
				break;
			case 'Shattering Throw':
				if (this.tag === playerIndex) {
					name += ` (self)`;
				}
				break;
			// Combo Point Spenders
			case 'Envenom':
			case 'Eviscerate':
			case 'Expose Armor':
			case 'Rupture':
			case 'Slice and Dice':
				if (this.tag) name += ` (${this.tag} CP)`;
				break;
			case 'Deadly Poison':
			case 'Deadly Poison II':
			case 'Deadly Poison III':
			case 'Deadly Poison IV':
			case 'Deadly Poison V':
			case 'Instant Poison':
			case 'Instant Poison II':
			case 'Instant Poison III':
			case 'Instant Poison IV':
			case 'Instant Poison V':
			case 'Instant Poison VI':
			case 'Wound Poison':
			case 'Occult Poison II':
				if (this.tag === 1) {
					name += ' (Shiv)';
				} else if (this.tag === 2) {
					name += ' (Deadly Brew)';
				} else if (this.tag === 100) {
					name += ' (Tick)';
				}
				break;
			case 'Saber Slash':
				if (this.tag === 100) {
					name += ' (Tick)';
				}
				break;
			// Dual-hit MH/OH spells and weapon imbues
			case 'Mutilate':
			case 'Stormstrike':
			case 'Carve':
			case 'Whirlwind':
			case 'Slam':
			case 'Windfury Weapon':
			case 'Holy Strength': // Weapon - Crusader Enchant
				if (this.tag === 1) {
					name = `${name} (Main-Hand)`;
				} else if (this.tag === 2) {
					name = `${name} (Off-Hand)`;
				}
				break;
			// Shaman Overload + Maelstrom Weapon
			case 'Lightning Bolt':
			case 'Chain Lightning':
			case 'Lava Burst':
			case 'Healing Wave':
			case 'Lesser Healing Wave':
			case 'Chain Heal':
				if (this.tag === 11) {
					name = `${name} OL`;
				} else if (this.tag) {
					name = `${name} (${this.tag} MSW)`;
				}
				break;
			case 'Holy Shield':
				if (this.tag === 1) {
					name += ' (Proc)';
				}
				break;
			case 'Righteous Vengeance':
				if (this.tag === 1) {
					name += ' (Application)';
				} else if (this.tag === 2) {
					name += ' (DoT)';
				}
				break;
			case 'Holy Vengeance':
				if (this.tag === 1) {
					name += ' (Application)';
				} else if (this.tag === 2) {
					name += ' (DoT)';
				}
				break;
			// For targetted buffs, tag is the source player's raid index or -1 if none.
			case 'Bloodlust':
			case 'Ferocious Inspiration':
			case 'Innervate':
			case 'Focus Magic':
			case 'Mana Tide Totem':
			case 'Power Infusion':
				if (this.tag !== -1) {
					if (this.tag === playerIndex || playerIndex === undefined) {
						name += ` (self)`;
					} else {
						name += ` (from #${this.tag + 1})`;
					}
				} else {
					name += ' (raid)';
				}
				break;
			case 'Darkmoon Card: Crusade':
				if (this.tag === 1) {
					name += ' (Melee)';
				} else if (this.tag === 2) {
					name += ' (Spell)';
				}
				break;
			case 'Battle Shout':
				if (this.tag === 1) {
					name += ' (Snapshot)';
				}
				break;
			case 'Heroic Strike':
			case 'Cleave':
			case 'Maul':
				if (this.tag === 1) {
					name += ' (Queue)';
				}
				break;
			// There are many different types of enrages. Try to give clarity to users.
			case 'Enrage':
				if (this.spellId === 13048) name = `${name} (Talent)`;
				else if (this.spellId === 14201) name = `${name} (Fresh Meat)`;
				else if (this.spellId === 425415) name = `${name} (Consumed by Rage)`;
				else if (this.spellId === 427066) name = `${name} (Wrecking Crew)`;
				break;
			case 'Raptor Strike':
				if (this.tag === 1) name = `${name} (Main-Hand)`;
				else if (this.tag === 2) name = `${name} (Off-Hand)`;
				else if (this.tag === 3) name = `${name} (Queue)`;
				break;
			case 'Thunderfury':
				if (this.tag === 1) {
					name += ' (Main)';
				} else if (this.tag === 2) {
					name += ' (Bounce)';
				}
				break;
			case 'Sunfire':
				if (this.spellId === 414689) {
					name = `${name} (Cat)`;
				}
				break;
			case 'Mangle':
				name = this.spellId === 409828 ? `${name} (Cat)` : `${name} (Bear)`;
				break;
			case 'Swipe':
				name = this.spellId === 411128 ? `${name} (Cat)` : `${name} (Bear)`;
				break;
			case 'Starfall':
				if (this.tag === 1) name = `${name} (Tick)`;
				else if (this.tag === 2) name = `${name} (Splash)`;
				break;
			case 'S03 - Item - T1 - Mage - Damage 4P Bonus':
				// Tags correspond to each non-physical spell school
				if (this.tag === 2) name = `${name} (Arcane)`;
				if (this.tag === 3) name = `${name} (Fire)`;
				if (this.tag === 4) name = `${name} (Frost)`;
				if (this.tag === 5) name = `${name} (Holy)`;
				if (this.tag === 6) name = `${name} (Nature)`;
				if (this.tag === 7) name = `${name} (Shadow)`;
				break;
			// Don't do anything for these but avoid adding "(??)"
			case 'S03 - Item - T1 - Shaman - Tank 6P Bonus':
				break;
			case 'Vampiric Touch':
				// Vampiric touch provided to the party
				if (this.tag === 1) name = `${name} (External)`;
				break;
			case 'Totem of Raging Fire':
				if (this.tag === 1) name = `${name} (1H)`;
				else if (this.tag === 2) name = `${name} (2H)`;
				break;
			// Warlock T2 6 Piece Needs Heals to trigger for the player
			case 'Drain Life':
			case 'Death Coil':
				if (this.tag === 1) name += ` (Heal)`;
				break;
			case 'Kill Shot':
				if (this.tag === 1) name = `${name} (Rapid Fire)`;
				break;
			case 'Master Demonologist':
				if (this.tag === 1) name = `${name} (Imp)`;
				else if (this.tag === 2) name = `${name} (Voidwalker)`;
				else if (this.tag === 3) name = `${name} (Succubus)`;
				else if (this.tag === 4) name = `${name} (Felhunter)`;
				else if (this.tag === 5) name = `${name} (Felguard)`;
				break;
			case 'Blood Plague':
				if (this.tag === 2) name = `${name} (Debuff)`;
				break;
			case 'Frost Fever':
				if (this.tag === 2) name = `${name} (Debuff)`;
				break;
			default:
				if (this.tag) {
					name += ' (??)';
				}
				break;
		}

		let iconOverrideId = this.spellIconOverride;

		// Icon Overrides based on tags
		switch (this.spellId) {
			// https://www.wowhead.com/classic/spell=457544/s03-item-t1-shaman-tank-6p-bonus
			case 457544: {
				// Show Stoneskin / Windwall respectively
				if (this.tag === 1) iconOverrideId = ActionId.fromSpellId(10408);
				else if (this.tag === 2) iconOverrideId = ActionId.fromSpellId(15112);
			}
		}

		let iconUrl = ActionId.makeIconUrl(tooltipData['icon']);
		if (iconOverrideId) {
			const overrideTooltipData = await ActionId.getTooltipData(iconOverrideId);
			iconUrl = ActionId.makeIconUrl(overrideTooltipData['icon']);
		}

		return new ActionId(this.itemId, this.spellId, this.otherId, this.tag, baseName, name, iconUrl, this.rank || tooltipData.rank, this.randomSuffixId);
	}

	toString(): string {
		return this.toStringIgnoringTag() + (this.tag ? '-' + this.tag : '');
	}

	toStringIgnoringTag(): string {
		if (this.itemId) {
			return 'item-' + this.itemId;
		} else if (this.spellId) {
			return 'spell-' + this.spellId;
		} else if (this.otherId) {
			return 'other-' + this.otherId;
		} else {
			throw new Error('Empty action id!');
		}
	}

	toProto(): ActionIdProto {
		const protoId = ActionIdProto.create({
			tag: this.tag,
		});

		if (this.itemId) {
			protoId.rawId = {
				oneofKind: 'itemId',
				itemId: this.itemId,
			};
		} else if (this.spellId) {
			protoId.rawId = {
				oneofKind: 'spellId',
				spellId: this.spellId,
			};
			protoId.rank = this.rank;
		} else if (this.otherId) {
			protoId.rawId = {
				oneofKind: 'otherId',
				otherId: this.otherId,
			};
		}

		return protoId;
	}

	toProtoString(): string {
		return ActionIdProto.toJsonString(this.toProto());
	}

	withoutTag(): ActionId {
		return new ActionId(this.itemId, this.spellId, this.otherId, 0, this.baseName, this.baseName, this.iconUrl, this.rank, this.randomSuffixId);
	}

	static fromEmpty(): ActionId {
		return new ActionId(0, 0, OtherAction.OtherActionNone, 0, '', '', '', 0);
	}

	static fromItemId(itemId: number, tag?: number, randomSuffixId?: number): ActionId {
		return new ActionId(itemId, 0, OtherAction.OtherActionNone, tag || 0, '', '', '', 0, randomSuffixId || 0);
	}

	static fromSpellId(spellId: number, rank = 0, tag?: number): ActionId {
		return new ActionId(0, spellId, OtherAction.OtherActionNone, tag || 0, '', '', '', rank);
	}

	static fromOtherId(otherId: OtherAction, tag?: number): ActionId {
		return new ActionId(0, 0, otherId, tag || 0, '', '', '', 0);
	}

	static fromPetName(petName: string): ActionId {
		return petNameToActionId[petName] || new ActionId(0, 0, OtherAction.OtherActionPet, 0, petName, petName, petNameToIcon[petName] || '', 0);
	}

	static fromItem(item: Item): ActionId {
		return ActionId.fromItemId(item.id);
	}

	static fromRandomSuffix(item: Item, randomSuffix: ItemRandomSuffix): ActionId {
		return ActionId.fromItemId(item.id, 0, randomSuffix.id);
	}

	static fromProto(protoId: ActionIdProto): ActionId {
		if (protoId.rawId.oneofKind === 'spellId') {
			return ActionId.fromSpellId(protoId.rawId.spellId, protoId.rank, protoId.tag);
		} else if (protoId.rawId.oneofKind === 'itemId') {
			return ActionId.fromItemId(protoId.rawId.itemId, protoId.tag);
		} else if (protoId.rawId.oneofKind === 'otherId') {
			return ActionId.fromOtherId(protoId.rawId.otherId, protoId.tag);
		} else {
			return ActionId.fromEmpty();
		}
	}

	private static readonly logRegex = /{((SpellID)|(ItemID)|(OtherID)): (\d+)(, Tag: (-?\d+))?}/;
	private static readonly logRegexGlobal = new RegExp(ActionId.logRegex, 'g');
	private static fromMatch(match: RegExpMatchArray): ActionId {
		const idType = match[1];
		const id = parseInt(match[5]);
		return new ActionId(
			idType === 'ItemID' ? id : 0,
			idType === 'SpellID' ? id : 0,
			idType === 'OtherID' ? id : 0,
			match[7] ? parseInt(match[7]) : 0,
			'',
			'',
			'',
			0,
		);
	}
	static fromLogString(str: string): ActionId {
		const match = str.match(ActionId.logRegex);
		if (match) {
			return ActionId.fromMatch(match);
		} else {
			console.warn('Failed to parse action id from log: ' + str);
			return ActionId.fromEmpty();
		}
	}

	static async replaceAllInString(str: string): Promise<string> {
		const matches = [...str.matchAll(ActionId.logRegexGlobal)];

		const replaceData = await Promise.all(
			matches.map(async match => {
				const actionId = ActionId.fromMatch(match);
				const filledId = await actionId.fill();
				return {
					firstIndex: match.index || 0,
					len: match[0].length,
					actionId: filledId,
				};
			}),
		);

		// Loop in reverse order so we can greedily apply the string replacements.
		for (let i = replaceData.length - 1; i >= 0; i--) {
			const data = replaceData[i];
			str = str.substring(0, data.firstIndex) + data.actionId.name + str.substring(data.firstIndex + data.len);
		}

		return str;
	}

	private static makeIconUrl(iconLabel: string): string {
		return `https://wow.zamimg.com/images/wow/icons/large/${iconLabel}.jpg`;
	}

	static async getTooltipData(actionId: ActionId): Promise<IconData> {
		if (actionId.itemId) {
			return await Database.getItemIconData(actionId.itemId);
		} else {
			return await Database.getSpellIconData(actionId.spellId);
		}
	}
	get spellIconOverride(): ActionId | null {
		const override = spellIdIconOverrides.get(JSON.stringify({ spellId: this.spellId }));
		if (!override) return null;
		return override.itemId ? ActionId.fromItemId(override.itemId) : ActionId.fromItemId(override.spellId!);
	}

	get spellTooltipOverride(): ActionId | null {
		const override = spellIdTooltipOverrides.get(JSON.stringify({ spellId: this.spellId, tag: this.tag }));
		if (!override) return null;
		return override.itemId ? ActionId.fromItemId(override.itemId) : ActionId.fromSpellId(override.spellId!);
	}
}

type ActionIdOverride = { itemId?: number; spellId?: number };

// Some items/spells have weird icons, so use this to show a different icon instead.
const spellIdIconOverrides: Map<string, ActionIdOverride> = new Map([
	[JSON.stringify({ spellId: 449288 }), { itemId: 221309 }], // Darkmoon Card: Sandstorm
	[JSON.stringify({ spellId: 455864 }), { spellId: 9907 }], // Tier 1 Balance Druid "Improved Faerie Fire"
	[JSON.stringify({ spellId: 457544 }), { spellId: 10408 }], // Tier 1 Shaman Tank "Improved Stoneskin / Windwall Totem"
]);

const spellIdTooltipOverrides: Map<string, ActionIdOverride> = new Map([]);

const spellIDsToShowBuffs = new Set([
	702, // https://www.wowhead.com/classic/spell=702/curse-of-weakness
	704, // https://www.wowhead.com/classic/spell=704/curse-of-recklessness
	770, // https://www.wowhead.com/classic/spell=770/faerie-fire
	778, // https://www.wowhead.com/classic/spell=778/faerie-fire
	1108, // https://www.wowhead.com/classic/spell=1108/curse-of-weakness
	1490, // https://www.wowhead.com/classic/spell=1490/curse-of-the-elements
	6205, // https://www.wowhead.com/classic/spell=6205/curse-of-weakness
	7646, // https://www.wowhead.com/classic/spell=7646/curse-of-weakness
	7658, // https://www.wowhead.com/classic/spell=7658/curse-of-recklessness
	7659, // https://www.wowhead.com/classic/spell=7659/curse-of-recklessness
	9749, // https://www.wowhead.com/classic/spell=9749/faerie-fire
	9907, // https://www.wowhead.com/classic/spell=9907/faerie-fire
	11707, // https://www.wowhead.com/classic/spell=11707/curse-of-weakness
	11708, // https://www.wowhead.com/classic/spell=11708/curse-of-weakness
	11717, // https://www.wowhead.com/classic/spell=11717/curse-of-recklessness
	11721, // https://www.wowhead.com/classic/spell=11721/curse-of-the-elements
	11722, // https://www.wowhead.com/classic/spell=11722/curse-of-the-elements
	14201, // https://www.wowhead.com/classic/spell=14201/enrage
	16257, // https://www.wowhead.com/classic/spell=16257/flurry
	16277, // https://www.wowhead.com/classic/spell=16277/flurry
	16278, // https://www.wowhead.com/classic/spell=16278/flurry
	16279, // https://www.wowhead.com/classic/spell=16279/flurry
	16280, // https://www.wowhead.com/classic/spell=16280/flurry
	17862, // https://www.wowhead.com/classic/spell=17862/curse-of-shadow
	17937, // https://www.wowhead.com/classic/spell=17937/curse-of-shadow
	18789, // https://www.wowhead.com/classic/spell=18789/burning-wish
	18790, // https://www.wowhead.com/classic/spell=18790/fel-stamina
	18791, // https://www.wowhead.com/classic/spell=18791/touch-of-shadow
	18792, // https://www.wowhead.com/classic/spell=18792/fel-energy
	20185, // https://www.wowhead.com/classic/spell=20185/judgement-of-light
	20186, // https://www.wowhead.com/classic/spell=20186/judgement-of-wisdom
	20300, // https://www.wowhead.com/classic/spell=20300/judgement-of-the-crusader
	20344, // https://www.wowhead.com/classic/spell=20344/judgement-of-light
	20345, // https://www.wowhead.com/classic/spell=20345/judgement-of-light
	20346, // https://www.wowhead.com/classic/spell=20346/judgement-of-light
	20355, // https://www.wowhead.com/classic/spell=20355/judgement-of-wisdom
	20301, // https://www.wowhead.com/classic/spell=20301/judgement-of-the-crusader
	20302, // https://www.wowhead.com/classic/spell=20302/judgement-of-the-crusader
	20303, // https://www.wowhead.com/classic/spell=20303/judgement-of-the-crusader
	23060, // https://www.wowhead.com/classic/spell=23060/battle-squawk
	23736, // https://www.wowhead.com/classic/spell=23736/sayges-dark-fortune-of-agility
	23737, // https://www.wowhead.com/classic/spell=23737/sayges-dark-fortune-of-stamina
	23738, // https://www.wowhead.com/classic/spell=23738/sayges-dark-fortune-of-spirit
	23766, // https://www.wowhead.com/classic/spell=23766/sayges-dark-fortune-of-intelligence
	23768, // https://www.wowhead.com/classic/spell=23768/sayges-dark-fortune-of-damage
	24907, // https://www.wowhead.com/classic/spell=24907/moonkin-aura
	24932, // https://www.wowhead.com/classic/spell=24932/leader-of-the-pack
	402808, // https://www.wowhead.com/classic/spell=402808/cripple
	425415, // https://www.wowhead.com/classic/spell=425415/enrage
	426969, // https://www.wowhead.com/classic/spell=426969/taste-for-blood
	440114, // https://www.wowhead.com/classic/spell=440114/sudden-death
	446393, // https://www.wowhead.com/classic/spell=446393/decay
	457699, // https://www.wowhead.com/classic/spell=457699/echoes-of-defensive-stance
	457706, // https://www.wowhead.com/classic/spell=457706/echoes-of-battle-stance
	457708, // https://www.wowhead.com/classic/spell=457708/echoes-of-berserker-stance
	457814, // https://www.wowhead.com/classic/spell=457814/defensive-forecast
	457816, // https://www.wowhead.com/classic/spell=457816/battle-forecast
	457817, // https://www.wowhead.com/classic/spell=457817/berserker-forecast
	457819, // https://www.wowhead.com/classic/spell=457819/echoes-of-gladiator-stance
	458403, // https://www.wowhead.com/classic/spell=458403/stalker
	461252, // https://www.wowhead.com/classic/spell=461252/shadowflame-fury
	461270, // https://www.wowhead.com/classic/spell=461270/magmadars-return
	461615, // https://www.wowhead.com/classic/spell=461615/mark-of-chaos
	439473, // https://www.wowhead.com/classic/spell=439473/atrophic-poison
	439472, // https://www.wowhead.com/classic/spell=439472/numbing-poison
	1214279, // https://www.wowhead.com/classic/spell=1214279/spell-blasting
	1218345, // https://www.wowhead.com/classic/spell=1218345/glaciate
	1218587, // https://www.wowhead.com/classic/spell=1218587/critical-aim
]);

export const defaultTargetIcon = 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_metamorphosis.jpg';

const petNameToActionId: Record<string, ActionId> = {
	'Eye of the Void': ActionId.fromSpellId(402789),
	'Frozen Orb 1': ActionId.fromSpellId(440802),
	'Frozen Orb 2': ActionId.fromSpellId(440802),
	Homunculi: ActionId.fromSpellId(402799),
	Shadowfiend: ActionId.fromSpellId(401977),
};

// https://wowhead.com/classic/hunter-pets
const petNameToIcon: Record<string, string> = {
	Bat: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_bat.jpg',
	Bear: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_bear.jpg',
	'Bird of Prey': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_owl.jpg',
	Boar: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_boar.jpg',
	'Carrion Bird': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_vulture.jpg',
	Cat: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_cat.jpg',
	Chimaera: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_chimera.jpg',
	'Core Hound': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_corehound.jpg',
	Crab: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_crab.jpg',
	Crocolisk: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_crocolisk.jpg',
	Devilsaur: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_devilsaur.jpg',
	Dragonhawk: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_dragonhawk.jpg',
	'Emerald Dragon Whelp': 'https://wow.zamimg.com/images/wow/icons/medium/inv_misc_head_dragon_green.jpg',
	Eskhandar: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_head_tiger_01.jpg',
	Felguard: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelguard.jpg',
	Felhunter: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelhunter.jpg',
	'Spirit Wolves': 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_feralspirit.jpg',
	Infernal: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summoninfernal.jpg',
	Gorilla: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_gorilla.jpg',
	Hyena: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_hyena.jpg',
	Imp: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonimp.jpg',
	'Mirror Image': 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg',
	Moth: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_moth.jpg',
	'Nether Ray': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_netherray.jpg',
	Owl: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_owl.jpg',
	Raptor: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_raptor.jpg',
	Ravager: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_ravager.jpg',
	Rhino: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_rhino.jpg',
	Scorpid: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_scorpid.jpg',
	Serpent: 'https://wow.zamimg.com/images/wow/icons/medium/spell_nature_guardianward.jpg',
	Silithid: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_silithid.jpg',
	Spider: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_spider.jpg',
	'Spirit Beast': 'https://wow.zamimg.com/images/wow/icons/medium/ability_druid_primalprecision.jpg',
	'Spore Bat': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_sporebat.jpg',
	Succubus: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonsuccubus.jpg',
	Tallstrider: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_tallstrider.jpg',
	Treants: 'https://wow.zamimg.com/images/wow/icons/medium/ability_druid_forceofnature.jpg',
	Turtle: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_turtle.jpg',
	Voidwalker: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonvoidwalker.jpg',
	'Warp Stalker': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_warpstalker.jpg',
	Wasp: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_wasp.jpg',
	'Wind Serpent': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_windserpent.jpg',
	Wolf: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_wolf.jpg',
	Worm: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_worm.jpg',
};

export function getPetIconFromName(name: string): string | ActionId | undefined {
	return petNameToActionId[name] || petNameToIcon[name];
}

export const resourceTypeToIcon: Record<ResourceType, string> = {
	[ResourceType.ResourceTypeNone]: '',
	[ResourceType.ResourceTypeHealth]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_elemental_mote_life01.jpg',
	[ResourceType.ResourceTypeMana]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_elemental_mote_mana.jpg',
	[ResourceType.ResourceTypeEnergy]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_shadow_shadowworddominate.jpg',
	[ResourceType.ResourceTypeRage]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_misc_emotionangry.jpg',
	[ResourceType.ResourceTypeComboPoints]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_mace_2h_pvp410_c_01.jpg',
	[ResourceType.ResourceTypeFocus]: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_focusfire.jpg',
};
