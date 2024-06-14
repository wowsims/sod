import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../../css_utils';
import { Player } from '../../player';
import { ItemSlot } from '../../proto/common';
import { UIRune } from '../../proto/ui.js';
import { ActionId } from '../../proto_utils/action_id';
import { getEnchantDescription } from '../../proto_utils/enchants';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { slotNames } from '../../proto_utils/names.js';
import { itemTypeToSlotsMap } from '../../proto_utils/utils.js';
import { SimUI } from '../../sim_ui';
import { EventID } from '../../typed_event';
import { Component } from '../component';
import { GearData } from './item_list';
import SelectorModal, { SelectorModalTabs } from './selector_modal';
import { getEmptySlotIconUrl } from './utils';

export default class GearPicker extends Component {
	// ItemSlot is used as the index
	readonly itemPickers: Array<ItemPicker>;
	readonly selectorModal: SelectorModal;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'gear-picker-root');

		const leftSide = document.createElement('div');
		leftSide.classList.add('gear-picker-left', 'tab-panel-col');
		this.rootElem.appendChild(leftSide);

		const rightSide = document.createElement('div');
		rightSide.classList.add('gear-picker-right', 'tab-panel-col');
		this.rootElem.appendChild(rightSide);

		const leftItemPickers = [
			ItemSlot.ItemSlotHead,
			ItemSlot.ItemSlotNeck,
			ItemSlot.ItemSlotShoulder,
			ItemSlot.ItemSlotBack,
			ItemSlot.ItemSlotChest,
			ItemSlot.ItemSlotWrist,
			ItemSlot.ItemSlotMainHand,
			ItemSlot.ItemSlotOffHand,
			ItemSlot.ItemSlotRanged,
		].map(slot => new ItemPicker(leftSide, this, simUI, player, slot));

		const rightItemPickers = [
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
			ItemSlot.ItemSlotFinger1,
			ItemSlot.ItemSlotFinger2,
			ItemSlot.ItemSlotTrinket1,
			ItemSlot.ItemSlotTrinket2,
		].map(slot => new ItemPicker(rightSide, this, simUI, player, slot));

		this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);

		this.selectorModal = new SelectorModal(simUI.rootElem, simUI, player, this);
	}
}

export class ItemRenderer extends Component {
	private readonly player: Player<any>;

	readonly iconElem: HTMLAnchorElement;
	readonly nameElem: HTMLAnchorElement;
	readonly ilvlElem: HTMLSpanElement;
	readonly enchantElem: HTMLAnchorElement;
	readonly runeElem: HTMLAnchorElement;

	constructor(parent: HTMLElement, root: HTMLElement, player: Player<any>) {
		super(parent, 'item-picker-root', root);
		this.player = player;

		const iconElem = ref<HTMLAnchorElement>();
		const nameElem = ref<HTMLAnchorElement>();
		const ilvlElem = ref<HTMLSpanElement>();
		const enchantElem = ref<HTMLAnchorElement>();
		const runeElem = ref<HTMLAnchorElement>();
		const sce = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<>
				<div className="item-picker-icon-wrapper">
					<span className="item-picker-ilvl" ref={ilvlElem} />
					<a ref={iconElem} className="item-picker-icon" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<div ref={sce} className="item-picker-sockets-container"></div>
				</div>
				<div className="item-picker-labels-container">
					<a ref={nameElem} className="item-picker-name" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<a ref={enchantElem} className="item-picker-enchant hide" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<a ref={runeElem} className="item-picker-rune hide" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
				</div>
			</>,
		);

		this.iconElem = iconElem.value!;
		this.nameElem = nameElem.value!;
		this.ilvlElem = ilvlElem.value!;
		this.enchantElem = enchantElem.value!;
		this.runeElem = runeElem.value!;
	}

	clear() {
		this.iconElem.removeAttribute('data-wowhead');
		this.iconElem.removeAttribute('href');
		this.nameElem.removeAttribute('data-wowhead');
		this.nameElem.removeAttribute('href');
		this.enchantElem.removeAttribute('data-wowhead');
		this.enchantElem.removeAttribute('href');
		this.enchantElem.classList.add('hide');
		this.runeElem.removeAttribute('data-wowhead');
		this.runeElem.removeAttribute('href');
		this.runeElem.classList.add('hide');

		this.iconElem.style.backgroundImage = '';

		this.nameElem.innerText = '';
		this.enchantElem.innerText = '';
		this.runeElem.innerText = '';
	}

	update(newItem: EquippedItem) {
		this.nameElem.textContent = newItem.item.name;
		this.ilvlElem.textContent = newItem.item.ilvl.toString();

		if (newItem.randomSuffix) {
			this.nameElem.textContent += ' ' + newItem.randomSuffix.name;
		}

		setItemQualityCssClass(this.nameElem, newItem.item.quality);

		this.player.setWowheadData(newItem, this.iconElem);
		this.player.setWowheadData(newItem, this.nameElem);
		newItem
			.asActionId()
			.fill()
			.then(filledId => {
				filledId.setBackgroundAndHref(this.iconElem);
				filledId.setWowheadHref(this.nameElem);
			});

		if (newItem.enchant) {
			getEnchantDescription(newItem.enchant).then(description => {
				this.enchantElem.textContent = description;
			});
			// Make enchant text hover have a tooltip.
			if (newItem.enchant.spellId) {
				this.enchantElem.href = ActionId.makeSpellUrl(newItem.enchant.spellId);
				ActionId.makeSpellTooltipData(newItem.enchant.spellId).then(url => {
					this.enchantElem.dataset.wowhead = url;
				});
			} else {
				this.enchantElem.href = ActionId.makeItemUrl(newItem.enchant.itemId);
				ActionId.makeItemTooltipData(newItem.enchant.itemId).then(url => {
					this.enchantElem.dataset.wowhead = url;
				});
			}
			this.enchantElem.dataset.whtticon = 'false';
			this.enchantElem.classList.remove('hide');
		} else {
			this.enchantElem.classList.add('hide');
		}

		const isRuneSlot = itemTypeToSlotsMap[newItem._item.type]?.some(slot => this.player.sim.db.hasRuneBySlot(slot, this.player.getClass()));
		if (isRuneSlot) {
			this.iconElem.appendChild(this.createRuneContainer(newItem.rune));

			if (newItem.rune) {
				this.runeElem.classList.remove('hide');
				this.runeElem.textContent = newItem.rune.name;
				this.runeElem.href = ActionId.makeSpellUrl(newItem.rune.id);
				this.runeElem.dataset.wowhead = `domain=classic&spell=${newItem.rune.id}`;
				this.runeElem.dataset.whtticon = 'false';
			}
		}
	}

	private createRuneContainer = (rune: UIRune | null) => {
		const runeIconElem = ref<HTMLImageElement>();
		const runeContainer = (
			<div className="item-picker-rune-container">
				<img ref={runeIconElem} className="item-picker-rune-icon" />
			</div>
		);

		if (rune) {
			ActionId.fromSpellId(rune.id)
				.fill()
				.then(filledId => (runeIconElem.value!.src = filledId.iconUrl));
		} else {
			runeIconElem.value!.src = 'https://wow.zamimg.com/images/wow/icons/medium/inventoryslot_empty.jpg';
		}

		return runeContainer;
	};
}

export class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly onUpdateCallbacks: (() => void)[] = [];

	private readonly itemElem: ItemRenderer;
	private readonly gearPicker: GearPicker;

	// All items and enchants that are eligible for this slot
	private _equippedItem: EquippedItem | null = null;

	constructor(parent: HTMLElement, gearPicker: GearPicker, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');

		this.gearPicker = gearPicker;
		this.simUI = simUI;
		this.player = player;
		this.slot = slot;
		this.itemElem = new ItemRenderer(parent, this.rootElem, player);

		this.item = player.getEquippedItem(slot);

		player.sim.waitForInit().then(() => {
			const openGearSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Items);
			};

			this.itemElem.iconElem.addEventListener('click', openGearSelector);
			this.itemElem.nameElem.addEventListener('click', openGearSelector);
		});

		player.gearChangeEmitter.on(() => {
			this.item = this.player.getEquippedItem(this.slot);
		});

		player.professionChangeEmitter.on(() => {
			if (!!this._equippedItem) {
				this.player.setWowheadData(this._equippedItem, this.itemElem.iconElem);
			}
		});
	}

	createGearData(): GearData {
		return {
			equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
				this.player.equipItem(eventID, this.slot, equippedItem);
			},
			getEquippedItem: () => this.player.getEquippedItem(this.slot),
			changeEvent: this.player.gearChangeEmitter,
		};
	}

	get item(): EquippedItem | null {
		return this._equippedItem;
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.itemElem.clear();
		this.itemElem.nameElem.textContent = slotNames.get(this.slot) ?? '';
		setItemQualityCssClass(this.itemElem.nameElem, null);

		if (!!newItem) {
			this.itemElem.update(newItem);
		} else {
			this.itemElem.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		}

		this._equippedItem = newItem;
		this.onUpdateCallbacks.forEach(callback => callback());
	}

	onUpdate(callback: () => void) {
		this.onUpdateCallbacks.push(callback);
	}

	openSelectorModal(selectedTab: SelectorModalTabs) {
		this.gearPicker.selectorModal.openTab(this.slot, selectedTab, this.createGearData());
	}
}
