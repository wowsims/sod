import { ref } from 'tsx-vanilla';

import { Player } from '../../player';
import { ItemSlot } from '../../proto/common';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { SimUI } from '../../sim_ui';
import { EventID } from '../../typed_event';
import { Component } from '../component';
import { GearData } from './item_list';
import SelectorModal, { SelectorModalTabs } from './selector_modal';
import { getEmptySlotIconUrl } from './utils';

export default class IconItemSwapPicker extends Component {
	private readonly iconAnchor: HTMLAnchorElement;
	private readonly player: Player<any>;
	private readonly slot: ItemSlot;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'icon-picker-root');
		this.rootElem.classList.add('icon-picker');
		this.player = player;
		this.slot = slot;

		const iconAnchorRef = ref<HTMLAnchorElement>();

		this.rootElem.prepend(<a ref={iconAnchorRef} className="icon-picker-button" href="#" attributes={{ role: 'button' }}></a>);

		this.iconAnchor = iconAnchorRef.value!;

		const selectorModal = new SelectorModal(simUI.rootElem, simUI, this.player);

		player.sim.waitForInit().then(() => {
			this.iconAnchor.addEventListener('click', (event: Event) => {
				event.preventDefault();
				selectorModal.openTab(this.slot, SelectorModalTabs.Items, this.createGearData());
			});
		});

		player.itemSwapSettings.changeEmitter.on(() => {
			this.update(player.itemSwapSettings.getGear().getEquippedItem(slot));
		});
	}

	update(newItem: EquippedItem | null) {
		this.iconAnchor.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconAnchor.removeAttribute('data-wowhead');
		this.iconAnchor.href = '#';

		if (newItem) {
			newItem.asActionId().fillAndSet(this.iconAnchor, true, true);
			this.player.setWowheadData(newItem, this.iconAnchor);
			this.iconAnchor.classList.add('active');
		} else {
			this.iconAnchor.classList.remove('active');
		}
	}

	private createGearData(): GearData {
		return {
			equipItem: (eventID: EventID, newItem: EquippedItem | null) => {
				this.player.itemSwapSettings.equipItem(eventID, this.slot, newItem);
			},
			getEquippedItem: () => this.player.itemSwapSettings.getItem(this.slot),
			changeEvent: this.player.itemSwapSettings.changeEmitter,
		};
	}
}
