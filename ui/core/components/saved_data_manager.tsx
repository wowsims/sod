import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { EventID, TypedEvent } from '../typed_event';
import { Component } from './component';
import { ContentBlock, ContentBlockHeaderConfig } from './content_block';

export type SavedDataManagerConfig<ModObject, T> = {
	label: string;
	header?: ContentBlockHeaderConfig;
	presetsOnly?: boolean;
	storageKey: string;
	changeEmitters: Array<TypedEvent<any>>;
	equals: (a: T, b: T) => boolean;
	getData: (modObject: ModObject) => T;
	setData: (eventID: EventID, modObject: ModObject, data: T) => void;
	toJson: (a: T) => any;
	fromJson: (obj: any) => T;
};

export type SavedDataConfig<ModObject, T> = {
	name: string;
	data: T;
	tooltip?: string;
	isPreset?: boolean;

	// If set, will automatically hide the saved data when this evaluates to false.
	enableWhen?: (obj: ModObject) => boolean;
	// Will execute when the saved data is loaded.
	onLoad?: (obj: ModObject) => void;
};

type SavedData<ModObject, T> = {
	name: string;
	data: T;
	elem: HTMLElement;
} & Pick<SavedDataConfig<ModObject, T>, 'enableWhen' | 'onLoad'>;

export class SavedDataManager<ModObject, T> extends Component {
	private readonly modObject: ModObject;
	private readonly config: SavedDataManagerConfig<ModObject, T>;

	private readonly userData: Array<SavedData<ModObject, T>>;
	private readonly presets: Array<SavedData<ModObject, T>>;

	private readonly savedDataDiv: HTMLElement;
	private readonly presetDataDiv: HTMLElement;
	private readonly customDataDiv: HTMLElement;
	private saveInput?: HTMLInputElement;

	private frozen: boolean;

	constructor(parent: HTMLElement, modObject: ModObject, config: SavedDataManagerConfig<ModObject, T>) {
		super(parent, 'saved-data-manager-root');
		this.modObject = modObject;
		this.config = config;

		this.userData = [];
		this.presets = [];
		this.frozen = false;

		const contentBlock = new ContentBlock(this.rootElem, 'saved-data', { header: config.header });

		const savedDataRef = ref<HTMLDivElement>();
		const presetDataRef = ref<HTMLDivElement>();
		const customDataRef = ref<HTMLDivElement>();
		contentBlock.bodyElement.replaceChildren(
			<div ref={savedDataRef} className="saved-data-container hide">
				<div ref={presetDataRef} className="saved-data-presets" />
				<div ref={customDataRef} className="saved-data-custom" />
			</div>,
		);

		this.savedDataDiv = savedDataRef.value!;
		this.presetDataDiv = presetDataRef.value!;
		this.customDataDiv = customDataRef.value!;

		if (!config.presetsOnly) {
			contentBlock.bodyElement.appendChild(this.buildCreateContainer());
		}
	}

	addSavedData(config: SavedDataConfig<ModObject, T>) {
		this.savedDataDiv.classList.remove('hide');

		const newData = this.makeSavedData(config);
		const dataArr = config.isPreset ? this.presets : this.userData;
		const oldIdx = dataArr.findIndex(data => data.name == config.name);

		if (oldIdx == -1) {
			if (config.isPreset) {
				this.presetDataDiv.appendChild(newData.elem);
			} else {
				this.customDataDiv.appendChild(newData.elem);
			}
			dataArr.push(newData);
		} else {
			dataArr[oldIdx].elem.replaceWith(newData.elem);
			dataArr[oldIdx] = newData;
		}
	}

	private makeSavedData(config: SavedDataConfig<ModObject, T>): SavedData<ModObject, T> {
		const deleteButtonRef = ref<HTMLAnchorElement>();
		const dataElem = (
			<div className="saved-data-set-chip badge rounded-pill">
				<a href="javascript:void(0)" className="saved-data-set-name" attributes={{ role: 'button' }}>
					{config.name}
				</a>
				{!config.isPreset && (
					<a ref={deleteButtonRef} href="javascript:void(0)" className="saved-data-set-delete" attributes={{ role: 'button' }}>
						<i className="fa fa-times fa-lg"></i>
					</a>
				)}
			</div>
		) as HTMLElement;

		dataElem.addEventListener('click', () => {
			this.config.setData(TypedEvent.nextEventID(), this.modObject, config.data);

			config.onLoad?.(this.modObject);
			if (this.saveInput) this.saveInput.value = config.name;
		});

		if (!config.isPreset && deleteButtonRef.value) {
			const tooltip = tippy(deleteButtonRef.value, { content: `Delete saved ${this.config.label}` });
			deleteButtonRef.value.addEventListener('click', event => {
				event.stopPropagation();
				const shouldDelete = confirm(`Delete saved ${this.config.label} '${config.name}'?`);
				if (!shouldDelete) return;

				tooltip.destroy();

				const idx = this.userData.findIndex(data => data.name == config.name);
				this.userData[idx].elem.remove();
				this.userData.splice(idx, 1);
				this.saveUserData();
			});
		}

		if (config.tooltip) {
			tippy(dataElem, {
				content: config.tooltip,
				placement: 'bottom',
			});
		}

		const checkActive = () => {
			if (this.config.equals(config.data, this.config.getData(this.modObject))) {
				dataElem.classList.add('active');
				if (this.saveInput) this.saveInput.value = config.name;
			} else {
				dataElem.classList.remove('active');
			}

			if (config.enableWhen && !config.enableWhen(this.modObject)) {
				dataElem.classList.add('disabled');
			} else {
				dataElem.classList.remove('disabled');
			}
		};

		checkActive();
		this.config.changeEmitters.forEach(emitter => emitter.on(checkActive));

		return {
			name: config.name,
			data: config.data,
			elem: dataElem,
			enableWhen: config.enableWhen,
			onLoad: config.onLoad,
		};
	}

	// Save data to window.localStorage.
	private saveUserData() {
		const userData: Record<string, unknown> = {};
		this.userData.forEach(savedData => {
			userData[savedData.name] = this.config.toJson(savedData.data);
		});

		if (this.userData.length == 0 && this.presets.length == 0) this.savedDataDiv.classList.add('hide');

		window.localStorage.setItem(this.config.storageKey, JSON.stringify(userData));
	}

	// Load data from window.localStorage.
	loadUserData() {
		const dataStr = window.localStorage.getItem(this.config.storageKey);
		if (!dataStr) return;

		let jsonData;
		try {
			jsonData = JSON.parse(dataStr);
		} catch (e) {
			console.warn('Invalid json for local storage value: ' + dataStr);
		}

		for (const name in jsonData) {
			try {
				this.addSavedData({
					name: name,
					data: this.config.fromJson(jsonData[name]),
				});
			} catch (e) {
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}
	}

	// Prevent user input from creating / deleting saved data.
	freeze() {
		this.frozen = true;
		this.rootElem.classList.add('frozen');
	}

	private buildCreateContainer(): HTMLElement {
		const saveButtonRef = ref<HTMLButtonElement>();
		const saveInputRef = ref<HTMLInputElement>();
		const savedDataCreateFragment = (
			<div className="saved-data-create-container">
				<label className="form-label">{this.config.label} Name</label>
				<input ref={saveInputRef} className="saved-data-save-input form-control" type="text" placeholder="Name" />
				<button ref={saveButtonRef} className="saved-data-save-button btn btn-primary">
					Save {this.config.label}
				</button>
			</div>
		) as HTMLElement;

		this.saveInput = saveInputRef.value!;
		saveButtonRef.value?.addEventListener('click', () => {
			if (this.frozen) return;

			const newName = this.saveInput?.value;
			if (!newName) {
				alert(`Choose a label for your saved ${this.config.label}!`);
				return;
			}

			if (newName in this.presets) {
				alert(`${this.config.label} with name ${newName} already exists.`);
				return;
			}
			this.addSavedData({
				name: newName,
				data: this.config.getData(this.modObject),
			});
			this.saveUserData();
		});

		return savedDataCreateFragment;
	}
}
