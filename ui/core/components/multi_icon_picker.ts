import tippy from 'tippy.js';

import { Player } from '../player.js';
import { ActionId } from '../proto_utils/action_id.js';
import { SimUI } from '../sim_ui.js';
import { TypedEvent } from '../typed_event.js';
import { isRightClick } from '../utils.js';
import { Component } from './component.js';
import { IconPicker, IconPickerConfig, IconPickerDirection } from './icon_picker.js';

export interface MultiIconPickerItemConfig<ModObject> extends IconPickerConfig<ModObject, any> {}

export interface MultiIconPickerConfig<ModObject> {
	values: Array<MultiIconPickerItemConfig<ModObject>>;
	categoryId?: ActionId;
	// The direction the menu will open in relative to the root element
	direction?: IconPickerDirection;
	label?: string;
	// Hover tooltip.
	tooltip?: string;
	showWhen?: (obj: Player<any>) => boolean;
}

// Icon-based UI for a dropdown with multiple icon pickers.
// ModObject is the object being modified (Sim, Player, or Target).
export class MultiIconPicker<ModObject> extends Component {
	protected simUI: SimUI;

	private readonly config: MultiIconPickerConfig<ModObject>;
	private readonly pickers: Array<IconPicker<ModObject, any>>;
	private currentValue: ActionId | null;

	private readonly buttonElem: HTMLAnchorElement;
	private readonly dropdownMenu: HTMLElement;

	constructor(parent: HTMLElement, modObj: ModObject, config: MultiIconPickerConfig<ModObject>, simUI: SimUI) {
		super(parent, 'multi-icon-picker-root');
		this.rootElem.classList.add('icon-picker', (config.direction ?? IconPickerDirection.Horizontal) === 'vertical' ? 'dropdown' : 'dropend');

		this.simUI = simUI;
		this.config = config;
		this.currentValue = null;

		this.rootElem.innerHTML = `
			<a
				class="icon-picker-button"
				role="button"
				data-bs-toggle="dropdown"
				aria-expanded="false"
				data-disable-wowhead-touch-tooltip="true"
				data-whtticon="false"
			></a>
			<ul class="dropdown-menu"></ul>
			<label class="multi-icon-picker-label form-label"></label>
	  `;

		if (config.tooltip) {
			const tooltip = tippy(this.rootElem, {
				content: config.tooltip,
			});
			this.addOnDisposeCallback(() => tooltip.destroy());
		}

		const labelElem = this.rootElem.querySelector<HTMLElement>('.multi-icon-picker-label')!;
		if (config.label) {
			labelElem.textContent = config.label;
		} else {
			labelElem.remove();
		}

		this.buttonElem = this.rootElem.querySelector<HTMLAnchorElement>('.icon-picker-button')!;
		this.dropdownMenu = this.rootElem.querySelector<HTMLElement>('.dropdown-menu')!;

		if (this.config.direction == IconPickerDirection.Horizontal) {
			this.dropdownMenu.style.gridAutoFlow = 'column';
		}

		this.buttonElem.addEventListener('hide.bs.dropdown', event => {
			if (event.hasOwnProperty('clickEvent')) {
				event.preventDefault();
			}
		});

		this.buttonElem.addEventListener('contextmenu', event => event.preventDefault());

		this.buttonElem.addEventListener('mousedown', event => {
			if (isRightClick(event)) {
				this.clearPicker();
			}
		});

		this.buildBlankOption();

		this.pickers = this.config.values.map((pickerConfig, _) => {
			const optionContainer = document.createElement('li');
			optionContainer.classList.add('icon-picker-option', 'dropdown-option');
			this.dropdownMenu.appendChild(optionContainer);

			return new IconPicker(optionContainer, modObj, pickerConfig);
		});
		simUI.sim.waitForInit().then(() => this.updateButtonImage());
		simUI.changeEmitter.on(() => {
			this.updateButtonImage();

			if (this.showWhen()) {
				this.rootElem.classList.remove('hide');
			} else {
				this.rootElem.classList.add('hide');
			}
		});
	}

	showWhen(): boolean {
		return !this.config.showWhen || (this.config.showWhen(this.simUI.sim.raid.getPlayer(0)!) && !!this.pickers.find(p => p.showWhen()));
	}

	private buildBlankOption() {
		const listItem = document.createElement('li');
		this.dropdownMenu.appendChild(listItem);

		const option = document.createElement('a');
		option.classList.add('icon-dropdown-option', 'dropdown-option');
		listItem.appendChild(option);

		option.addEventListener('click', () => this.clearPicker());
	}

	private clearPicker() {
		TypedEvent.freezeAllAndDo(() => {
			this.pickers.forEach(picker => {
				picker.setInputValue(null);
				picker.inputChanged(TypedEvent.nextEventID());
			});
			this.updateButtonImage();
		});
	}

	private updateButtonImage() {
		this.currentValue = this.getMaxValue();

		if (this.currentValue) {
			this.buttonElem.classList.add('active');
			if (this.config.categoryId != null) {
				this.config.categoryId.fillAndSet(this.buttonElem, false, true);
			} else {
				this.currentValue.fillAndSet(this.buttonElem, false, true);
			}
		} else {
			this.buttonElem.classList.remove('active');
			if (this.config.categoryId != null) {
				this.config.categoryId.fillAndSet(this.buttonElem, false, true);
			} else {
				this.buttonElem.style.backgroundImage = '';
			}
			this.buttonElem.removeAttribute('href');
		}
	}

	private getMaxValue(): ActionId | null {
		return (
			this.pickers
				.filter(picker => picker.showWhen())
				.map(picker => picker.getActionId())
				.filter(id => id != null)[0] || null
		);
	}
}
