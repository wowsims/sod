import clsx from 'clsx';
import tippy, { Instance as TippyInstance } from 'tippy.js';

import { EventID, TypedEvent } from '../typed_event.js';
import { Input, InputConfig } from './input.js';

export type ListItemAction = 'create' | 'delete' | 'move' | 'copy';

export interface ListPickerActionsConfig {
	create?: {
		// Whether or not to use an icon for the create action button
		// defaults to FALSE
		useIcon?: boolean;
	};
}

export interface ListPickerConfig<ModObject, ItemType> extends Omit<InputConfig<ModObject, Array<ItemType>>, 'id'> {
	itemLabel: string;
	newItem: () => ItemType;
	copyItem: (oldItem: ItemType) => ItemType;
	newItemPicker: (
		parent: HTMLElement,
		listPicker: ListPicker<ModObject, ItemType>,
		index: number,
		config: ListItemPickerConfig<ModObject, ItemType>,
	) => Input<ModObject, ItemType>;
	actions?: ListPickerActionsConfig;
	title?: string;
	titleTooltip?: string;
	inlineMenuBar?: boolean;
	hideUi?: boolean;
	horizontalLayout?: boolean;
	// if set, will remove the border and padding of the list items
	isCompact?: boolean;
	// If set, will disable the delete button if the list is at the minimum.
	minimumItems?: number;
	// If set, only actions included in the list are allowed. Otherwise, all actions are allowed.
	allowedActions?: Array<ListItemAction>;
}

const DEFAULT_CONFIG = {
	actions: {
		create: {
			useIcon: false,
		},
	},
};

export interface ListItemPickerConfig<ModObject, ItemType> extends InputConfig<ModObject, ItemType> {}

interface ItemPickerPair<ItemType> {
	elem: HTMLElement;
	picker: Input<any, ItemType>;
	idx: number;
}

interface ListDragData<ModObject, ItemType> {
	listPicker: ListPicker<ModObject, ItemType>;
	item: ItemPickerPair<ItemType>;
}

let curDragData: ListDragData<any, any> | null = null;

export class ListPicker<ModObject, ItemType> extends Input<ModObject, Array<ItemType>> {
	readonly config: ListPickerConfig<ModObject, ItemType>;
	private readonly itemsDiv: HTMLElement;

	private itemPickerPairs: Array<ItemPickerPair<ItemType>>;

	constructor(parent: HTMLElement, modObject: ModObject, config: ListPickerConfig<ModObject, ItemType>) {
		if (config.isCompact) config.extraCssClasses = [...(config.extraCssClasses || []), 'list-picker-compact'];

		super(parent, 'list-picker-root', modObject, config);
		this.config = { ...DEFAULT_CONFIG, ...config };
		this.itemPickerPairs = [];

		this.rootElem.appendChild(
			<>
				{config.title && <label className="list-picker-title form-label">{config.title}</label>}
				<div className="list-picker-items"></div>
			</>,
		);

		if (this.config.hideUi) {
			this.rootElem.classList.add('d-none');
		}
		if (this.config.horizontalLayout) {
			this.config.inlineMenuBar = true;
			this.rootElem.classList.add('horizontal');
		}

		if (this.config.titleTooltip) {
			const titleTooltip = tippy(this.rootElem.querySelector<HTMLElement>('.list-picker-title')!, {
				content: this.config.titleTooltip,
			});
			this.addOnDisposeCallback(() => titleTooltip?.destroy());
		}

		this.itemsDiv = this.rootElem.getElementsByClassName('list-picker-items')[0] as HTMLElement;

		if (this.actionEnabled('create')) {
			let newItemButton: HTMLElement | null = null;
			let newButtonTooltip: TippyInstance | null = null;
			if (this.config.actions?.create?.useIcon) {
				newItemButton = ListPicker.makeActionElem('link-success', 'fa-plus');
				newButtonTooltip = tippy(newItemButton, {
					allowHTML: false,
					content: `New ${config.itemLabel}`,
				});
				this.addOnDisposeCallback(() => newButtonTooltip?.destroy());
			} else {
				newItemButton = (<button className="btn btn-primary">New {config.itemLabel}</button>) as HTMLButtonElement;
			}
			newItemButton.classList.add('list-picker-new-button');
			newItemButton.addEventListener(
				'click',
				() => {
					const newItem = this.config.newItem();
					const newList = this.config.getValue(this.modObject).concat([newItem]);
					this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
					if (newButtonTooltip) {
						newButtonTooltip.hide();
					}
				},
				{ signal: this.signal },
			);

			this.rootElem.appendChild(newItemButton);
		}

		this.init();
	}

	getInputElem(): HTMLElement {
		return this.rootElem;
	}

	getInputValue(): Array<ItemType> {
		return this.itemPickerPairs.map(pair => pair.picker.getInputValue());
	}

	setInputValue(newValue: Array<ItemType>): void {
		// Add/remove pickers to make the lengths match.
		if (newValue.length < this.itemPickerPairs.length) {
			this.itemPickerPairs.slice(newValue.length).forEach(ipp => ipp.elem.remove());
			this.itemPickerPairs = this.itemPickerPairs.slice(0, newValue.length);
		} else if (newValue.length > this.itemPickerPairs.length) {
			const numToAdd = newValue.length - this.itemPickerPairs.length;
			for (let i = 0; i < numToAdd; i++) {
				this.addNewPicker();
			}
		}

		// Set all the values.
		newValue.forEach((val, i) => this.itemPickerPairs[i].picker.setInputValue(val));
	}

	private actionEnabled(action: ListItemAction): boolean {
		return !this.config.allowedActions || this.config.allowedActions.includes(action);
	}

	private addHoverListeners(button: HTMLButtonElement) {
		button.addEventListener(
			'mouseenter',
			() => {
				button.classList.add('hover');
			},
			{ signal: this.signal },
		);

		button.addEventListener(
			'mouseleave',
			() => {
				button.classList.remove('hover');
			},
			{ signal: this.signal },
		);
	}

	private addNewPicker() {
		const index = this.itemPickerPairs.length;
		const itemContainer = document.createElement('div');
		itemContainer.classList.add('list-picker-item-container');
		if (this.config.inlineMenuBar) {
			itemContainer.classList.add('inline');
		}
		this.itemsDiv.appendChild(itemContainer);

		const itemElem = document.createElement('div');
		itemElem.classList.add('list-picker-item');

		const itemHeader = document.createElement('div');
		itemHeader.classList.add('list-picker-item-header');

		const popover = document.createElement('div');
		popover.classList.add('list-picker-item-popover');
		popover.setAttribute('popover', 'auto');
		itemHeader.appendChild(popover);
		let hasActions = false;

		if (this.config.inlineMenuBar) {
			itemContainer.appendChild(itemElem);
			itemContainer.appendChild(itemHeader);
		} else {
			itemContainer.appendChild(itemHeader);
			itemContainer.appendChild(itemElem);
			if (this.config.itemLabel) {
				const itemLabel = document.createElement('h6');
				itemLabel.classList.add('list-picker-item-title');
				itemLabel.textContent = `${this.config.itemLabel} ${this.itemPickerPairs.length + 1}`;
				itemHeader.appendChild(itemLabel);
			}
		}

		const itemPicker = this.config.newItemPicker(itemElem, this, index, {
			changedEvent: this.config.changedEvent,
			getValue: () => this.getSourceValue()[index],
			setValue: (eventID: EventID, modObj: ModObject, newValue: ItemType) => {
				const newList = this.getSourceValue();
				newList[index] = newValue;
				this.config.setValue(eventID, modObj, newList);
			},
		});

		const item: ItemPickerPair<ItemType> = { elem: itemContainer, picker: itemPicker, idx: index };

		if (this.actionEnabled('delete')) {
			if (!this.config.minimumItems || index + 1 > this.config.minimumItems) {
				hasActions = true;
				const deleteButton = ListPicker.makeActionElem('list-picker-item-delete', 'fa-times');
				deleteButton.classList.add('link-danger');
				popover.appendChild(deleteButton);

				const deleteButtonTooltip = tippy(deleteButton, {
					allowHTML: false,
					content: `Delete ${this.config.itemLabel}`,
				});

				deleteButton.addEventListener(
					'click',
					() => {
						const newList = this.config.getValue(this.modObject);
						newList.splice(index, 1);
						this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
						deleteButtonTooltip.hide();
					},
					{ signal: this.signal },
				);
				this.addOnDisposeCallback(() => deleteButtonTooltip?.destroy());
				this.addHoverListeners(deleteButton);
			}
		}

		if (this.actionEnabled('copy')) {
			hasActions = true;
			const copyButton = ListPicker.makeActionElem('list-picker-item-copy', 'fa-copy');
			popover.appendChild(copyButton);
			const copyButtonTooltip = tippy(copyButton, {
				allowHTML: false,
				content: `Copy to New ${this.config.itemLabel}`,
			});

			copyButton.addEventListener(
				'click',
				() => {
					const newList = this.config.getValue(this.modObject).slice();
					newList.splice(index, 0, this.config.copyItem(newList[index]));
					this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);
					copyButtonTooltip.hide();
				},
				{ signal: this.signal },
			);
			this.addOnDisposeCallback(() => copyButtonTooltip?.destroy());
			this.addHoverListeners(copyButton);
		}

		if (this.actionEnabled('move')) {
			hasActions = true;
			itemContainer.classList.add('draggable');
			if (this.config.itemLabel) {
				itemContainer.classList.add(this.config.itemLabel.toLowerCase().replace(' ', '-'));
			}

			const moveButton = ListPicker.makeActionElem('list-picker-item-move', 'fa-arrows-up-down');
			popover.appendChild(moveButton);

			const moveButtonTooltip = tippy(moveButton, {
				allowHTML: false,
				content: 'Move (Drag+Drop)',
			});

			moveButton.addEventListener(
				'click',
				() => {
					moveButtonTooltip.hide();
				},
				{ signal: this.signal },
			);
			this.addOnDisposeCallback(() => {
				moveButtonTooltip?.destroy();
			});

			this.addHoverListeners(moveButton);

			moveButton.addEventListener(
				'mousedown',
				() => {
					moveButton.setAttribute('draggable', 'true');
					itemContainer.setAttribute('draggable', 'true');
				},
				{ signal: this.signal },
			);

			moveButton.addEventListener(
				'mouseup',
				() => {
					moveButton.removeAttribute('draggable');
					itemContainer.removeAttribute('draggable');
				},
				{ signal: this.signal },
			);

			moveButton.addEventListener(
				'dragstart',
				event => {
					if (event.target == moveButton) {
						const popoverRect = popover.getBoundingClientRect();
						event.dataTransfer!.setDragImage(itemContainer, 0, popoverRect.height / 2);
						event.dataTransfer!.dropEffect = 'move';
						event.dataTransfer!.effectAllowed = 'move';
						itemContainer.classList.add('dragfrom');
						curDragData = {
							listPicker: this,
							item: item,
						};
					}
				},
				{ signal: this.signal },
			);

			const droppingActionOnOtherList = () => curDragData && this.config.itemLabel === 'Action' && curDragData.listPicker !== this;
			const targetIsSelf = () => curDragData && curDragData.listPicker === this && curDragData.item.idx === index;
			const targetIsChild = () => curDragData && curDragData.item.elem.contains(itemContainer);

			const invalidDropTarget = (checkSelf = true) => {
				// Only allow dropping on the same type of list, Value -> Value, Action -> Action
				if (!curDragData || curDragData.listPicker.config.itemLabel !== this.config.itemLabel) {
					return true;
				}

				// Only allow dropping Actions within the same list
				if (droppingActionOnOtherList()) {
					return true;
				}

				// Just skip trying to drop on itself?
				if (checkSelf && targetIsSelf()) {
					return true;
				}

				// Can't drop within itself
				if (checkSelf && targetIsChild()) {
					return true;
				}

				return false;
			};

			let dragEnterCounter = 0;
			itemContainer.addEventListener(
				'dragenter',
				event => {
					if (invalidDropTarget()) {
						return;
					}
					event.stopPropagation();
					dragEnterCounter++;
					itemContainer.classList.add('dragto');
				},
				{ signal: this.signal },
			);

			itemContainer.addEventListener(
				'dragleave',
				event => {
					if (invalidDropTarget()) {
						return;
					}
					event.preventDefault();
					dragEnterCounter--;
					if (dragEnterCounter <= 0) {
						itemContainer.classList.remove('dragto');
					}
				},
				{ signal: this.signal },
			);

			itemContainer.addEventListener(
				'dragover',
				event => {
					if (invalidDropTarget()) {
						if (droppingActionOnOtherList() || targetIsSelf()) {
							event.dataTransfer!.dropEffect = 'none';
						}

						return;
					}
					event.dataTransfer!.dropEffect = 'move';
					event.stopPropagation();
					event.preventDefault();
				},
				{ signal: this.signal },
			);

			const cleanupAfterDrag = () => {
				if (!curDragData) {
					return;
				}
				moveButton.removeAttribute('draggable');
				itemContainer.removeAttribute('draggable');
				curDragData.item.elem.removeAttribute('draggable');
				[...document.querySelectorAll('.dragfrom,.dragto')].forEach(elem => {
					elem.classList.remove('dragfrom');
					elem.classList.remove('dragto');
				});
			};

			itemContainer.addEventListener(
				'dragend',
				event => {
					if (invalidDropTarget(false)) {
						return;
					}
					event.stopPropagation();
					cleanupAfterDrag();
					curDragData = null;
				},
				{ signal: this.signal },
			);

			itemContainer.addEventListener(
				'drop',
				event => {
					if (!curDragData || invalidDropTarget()) {
						if (targetIsSelf()) {
							event.stopPropagation();
							cleanupAfterDrag();
						}
						return;
					}
					event.stopPropagation();
					cleanupAfterDrag();

					const srcIdx = curDragData.item.idx;
					let dstIdx = index;

					const targetRect = itemContainer.getBoundingClientRect();
					if (event.clientY > targetRect.top + targetRect.height / 2) {
						dstIdx++;
					}

					const newList = this.config.getValue(this.modObject);
					let arrElem;

					if (curDragData.listPicker !== this) {
						const oldList = curDragData.listPicker.config.getValue(curDragData.listPicker.modObject);
						arrElem = oldList[srcIdx];
						oldList.splice(srcIdx, 1);
						curDragData.listPicker.config.setValue(TypedEvent.nextEventID(), curDragData.listPicker.modObject, oldList);
					} else {
						arrElem = newList[srcIdx];
						newList.splice(srcIdx, 1);
					}

					newList.splice(dstIdx, 0, arrElem);
					this.config.setValue(TypedEvent.nextEventID(), this.modObject, newList);

					curDragData = null;
				},
				{ signal: this.signal },
			);
		}

		if (hasActions) {
			const actionsButton = ListPicker.makeActionElem('list-picker-item-actions', 'fa-ellipsis');
			itemHeader.appendChild(actionsButton);
			actionsButton.addEventListener(
				'mouseover',
				() => {
					popover.showPopover();
					const actionsButtonRect = actionsButton.getBoundingClientRect();
					const popoverRect = popover.getBoundingClientRect();
					const diff = (popoverRect.height - actionsButtonRect.height) / 2;
					popover.style.top = actionsButtonRect.top - diff + 'px';
					popover.style.left = actionsButtonRect.right - popoverRect.width + 10 + 'px';
					popover.classList.add('hover');
				},
				{ signal: this.signal },
			);
			popover.addEventListener(
				'mouseleave',
				() => {
					popover.classList.remove('hover');
					popover.hidePopover();
				},
				{ signal: this.signal },
			);
		}

		this.itemPickerPairs.push(item);
	}

	static makeActionElem(cssClass: string, iconCssClass: string): HTMLButtonElement {
		return (
			<button type="button" className={clsx('list-picker-item-action', cssClass)}>
				<i className={clsx('fa', 'fa-xl', iconCssClass)} />
			</button>
		) as HTMLButtonElement;
	}

	static getItemHeaderElem(itemPicker: Input<any, any>): HTMLElement {
		const itemElem = itemPicker.rootElem.parentElement!;
		const headerElem = itemElem.nextElementSibling || itemElem.previousElementSibling;
		if (!headerElem?.classList.contains('list-picker-item-header')) {
			throw new Error('Could not find list item header');
		}
		return headerElem as HTMLElement;
	}
}
