import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../../css_utils';
import { IndividualSimUI } from '../../individual_sim_ui';
import { BulkSettings, ErrorOutcomeType, ProgressMetrics, TalentLoadout } from '../../proto/api';
import { ItemSpec, SimDatabase, SimEnchant, SimItem } from '../../proto/common';
import { SavedTalents, UIEnchant, UIItem, UIItem_FactionRestriction } from '../../proto/ui';
import { canEquipItem } from '../../proto_utils/utils';
import { RequestTypes } from '../../sim_signal_manager';
import { TypedEvent } from '../../typed_event';
import { cloneChildren } from '../../utils';
import { WorkerProgressCallback } from '../../worker_pool';
import { BooleanPicker } from '../boolean_picker';
import { ContentBlock } from '../content_block';
import SelectorModal, { SelectorModalTabs } from '../gear_picker/selector_modal';
import { ResultsViewer } from '../results_viewer';
import { SimTab } from '../sim_tab';
import Toast from '../toast';
import BulkItemPicker from './bulk/bulk_item_picker';
import BulkSimResultRenderer from './bulk/bulk_sim_result_renderer';
import { BulkGearJsonImporter } from './importers';

export class BulkTab extends SimTab {
	readonly simUI: IndividualSimUI<any>;

	readonly itemsChangedEmitter = new TypedEvent<void>();

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	readonly column1: HTMLElement = this.buildColumn(1, 'raid-settings-col');

	protected items: Array<ItemSpec> = new Array<ItemSpec>();

	private pendingResults: ResultsViewer;
	private pendingDiv: HTMLDivElement;

	// TODO: Make a real options probably
	private doCombos: boolean;
	private fastMode: boolean;
	private simTalents: boolean;
	private autoEnchant: boolean;
	private savedTalents: TalentLoadout[];
	readonly selectorModal: SelectorModal;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<any>) {
		super(parentElem, simUI, { identifier: 'bulk-tab', title: 'Batch' });
		this.simUI = simUI;

		this.leftPanel = (<div className="bulk-tab-left tab-panel-left">{this.column1}</div>) as HTMLDivElement;
		this.rightPanel = (<div className="bulk-tab-right tab-panel-right" />) as HTMLDivElement;

		this.pendingDiv = (<div className="results-pending-overlay d-flex hide" />) as HTMLDivElement;
		this.pendingResults = new ResultsViewer(this.pendingDiv);
		this.pendingResults.hideAll();
		this.selectorModal = new SelectorModal(this.simUI.rootElem, this.simUI, this.simUI.player, undefined, {
			id: 'bulk-selector-modal',
			disabledTabs: [SelectorModalTabs.Items],
		});

		this.contentContainer.appendChild(
			<>
				{this.leftPanel}
				{this.rightPanel}
				{this.pendingDiv}
			</>,
		);

		this.doCombos = true;
		this.fastMode = true;
		this.autoEnchant = true;
		this.savedTalents = [];
		this.simTalents = false;
		this.buildTabContent();

		this.simUI.sim.waitForInit().then(() => {
			this.loadSettings();
		});
	}

	private getSettingsKey(): string {
		return this.simUI.getStorageKey('bulk-settings.v1');
	}

	private loadSettings() {
		const storedSettings = window.localStorage.getItem(this.getSettingsKey());
		if (storedSettings != null) {
			const settings = BulkSettings.fromJsonString(storedSettings, {
				ignoreUnknownFields: true,
			});

			this.doCombos = settings.combinations;
			this.fastMode = settings.fastMode;
			this.autoEnchant = settings.autoEnchant;
			this.savedTalents = settings.talentsToSim;
			this.simTalents = settings.simTalents;
		}
	}

	private storeSettings() {
		const settings = this.createBulkSettings();
		const setStr = BulkSettings.toJsonString(settings, { enumAsInteger: true });
		window.localStorage.setItem(this.getSettingsKey(), setStr);
	}

	protected createBulkSettings(): BulkSettings {
		return BulkSettings.create({
			items: this.items,
			// TODO(Riotdog-GehennasEU): Make all of these configurable.
			// For now, it's always constant iteration combinations mode for "sim my bags".
			combinations: this.doCombos,
			fastMode: this.fastMode,
			autoEnchant: this.autoEnchant,
			simTalents: this.simTalents,
			talentsToSim: this.savedTalents,
			iterationsPerCombo: this.simUI.sim.getIterations(), // TODO(Riotdog-GehennasEU): Define a new UI element for the iteration setting.
		});
	}

	protected createBulkItemsDatabase(): SimDatabase {
		const itemsDb = SimDatabase.create();
		for (const is of this.items) {
			const item = this.simUI.sim.db.lookupItemSpec(is);
			if (!item) {
				throw new Error(`item with ID ${is.id} not found in database`);
			}
			itemsDb.items.push(SimItem.fromJson(UIItem.toJson(item.item), { ignoreUnknownFields: true }));
			if (item.enchant) {
				itemsDb.enchants.push(
					SimEnchant.fromJson(UIEnchant.toJson(item.enchant), {
						ignoreUnknownFields: true,
					}),
				);
			}
		}
		return itemsDb;
	}

	addItem(item: ItemSpec) {
		this.addItems([item]);
	}
	addItems(items: ItemSpec[]) {
		this.items = [...(this.items || []), ...items];
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	setItems(items: ItemSpec[]) {
		this.items = items;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	removeItem(item: ItemSpec) {
		const indexToRemove = this.items.findIndex(i => ItemSpec.equals(i, item));
		if (indexToRemove === -1) return;
		this.items.splice(indexToRemove, 1);
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}
	removeItemByIndex(index: number) {
		if (this.items.length < index) {
			new Toast({
				variant: 'error',
				body: 'Failed to remove item, please report this issue.',
			});
			return;
		}
		this.items.splice(index, 1);
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	clearItems() {
		this.items = new Array<ItemSpec>();
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	hasItem(item: ItemSpec) {
		return this.items.some(i => ItemSpec.equals(i, item));
	}

	getItems(): Array<ItemSpec> {
		const result = new Array<ItemSpec>();
		this.items.forEach(spec => {
			result.push(ItemSpec.clone(spec));
		});
		return result;
	}

	setCombinations(doCombos: boolean) {
		this.doCombos = doCombos;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	setFastMode(fastMode: boolean) {
		this.fastMode = fastMode;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	protected async runBulkSim(onProgress: WorkerProgressCallback) {
		this.pendingResults.setPending();

		try {
			const result = await this.simUI.sim.runBulkSim(this.createBulkSettings(), this.createBulkItemsDatabase(), onProgress);
			if (result.error?.type == ErrorOutcomeType.ErrorOutcomeAborted) {
				new Toast({
					variant: 'info',
					body: 'Bulk sim cancelled.',
				});
			}
		} catch (e) {
			this.simUI.handleCrash(e);
		}
	}

	protected buildTabContent() {
		const itemsBlock = new ContentBlock(this.column1, 'bulk-items', {
			header: { title: 'Items' },
		});
		itemsBlock.bodyElement.classList.add('gear-picker-root', 'gear-picker-root-bulk');

		const itemTextIntro = (
			<div className="bulk-items-text-line">
				<i>
					Notice: This is under very early but active development and experimental. You may also need to update your WoW AddOn if you want to import
					your bags.
				</i>
			</div>
		);

		const itemList = (<div className="tab-panel-col bulk-gear-combo" />) as HTMLElement;

		this.itemsChangedEmitter.on(() => {
			const items = (<></>) as HTMLElement;
			if (!!this.items.length) {
				itemTextIntro.textContent = 'The following items will be simmed together with your equipped gear.';
				for (let i = 0; i < this.items.length; ++i) {
					const spec = this.items[i];
					const item = this.simUI.sim.db.lookupItemSpec(spec);
					new BulkItemPicker(items, this.simUI, this, item!, i);
				}
			}
			itemList.replaceChildren(items);
		});

		itemsBlock.bodyElement.appendChild(
			<>
				{itemTextIntro}
				{itemList}
			</>,
		);

		this.clearItems();

		const resultsBlock = new ContentBlock(this.column1, 'bulk-results', {
			header: {
				title: 'Results',
				extraCssClasses: ['bulk-results-header'],
			},
		});

		resultsBlock.rootElem.hidden = true;
		resultsBlock.bodyElement.classList.add('gear-picker-root', 'gear-picker-root-bulk', 'tab-panel-col');

		this.simUI.sim.bulkSimStartEmitter.on(() => {
			resultsBlock.rootElem.hidden = true;
		});

		this.simUI.sim.bulkSimResultEmitter.on((_, bulkSimResult) => {
			resultsBlock.rootElem.hidden = bulkSimResult.results.length == 0;
			resultsBlock.bodyElement.replaceChildren();

			for (const r of bulkSimResult.results) {
				const resultBlock = new ContentBlock(resultsBlock.bodyElement, 'bulk-result', {
					header: { title: '' },
					bodyClasses: ['bulk-results-body'],
				});
				new BulkSimResultRenderer(resultBlock.bodyElement, this.simUI, r, bulkSimResult.equippedGearResult!);
			}
		});

		const settingsBlock = new ContentBlock(this.rightPanel, 'bulk-settings', {
			header: { title: 'Setup' },
		});

		const bulkSimButton = (<button className="btn btn-primary w-100 bulk-settings-button">Simulate Batch</button>) as HTMLButtonElement;

		let isRunning = false;
		bulkSimButton.addEventListener('click', async () => {
			if (isRunning) return;
			isRunning = true;
			bulkSimButton.disabled = true;

			this.pendingDiv.classList.remove('hide');
			this.leftPanel.classList.add('blurred');
			this.rightPanel.classList.add('blurred');

			const defaultState = cloneChildren(bulkSimButton);
			bulkSimButton.disabled = true;
			bulkSimButton.classList.add('disabled');
			bulkSimButton.replaceChildren(
				<>
					<i className="fa fa-spinner fa-spin" /> Running
				</>,
			);

			let waitAbort = false;
			try {
				await this.simUI.sim.signalManager.abortType(RequestTypes.All);

				this.pendingResults.addAbortButton(async () => {
					if (waitAbort) return;
					try {
						waitAbort = true;
						await this.simUI.sim.signalManager.abortType(RequestTypes.BulkSim);
					} catch (error) {
						console.error('Error on bulk sim abort!');
						console.error(error);
					} finally {
						waitAbort = false;
						if (!isRunning) bulkSimButton.disabled = false;
					}
				});

				let simStart = new Date().getTime();
				let lastTotal = 0;
				let rounds = 0;
				let currentRound = 0;
				let combinations = 0;

				await this.runBulkSim((progressMetrics: ProgressMetrics) => {
					const msSinceStart = new Date().getTime() - simStart;
					const iterPerSecond = progressMetrics.completedIterations / (msSinceStart / 1000);

					if (combinations === 0) {
						combinations = progressMetrics.totalSims;
					}
					if (this.fastMode) {
						if (rounds === 0 && progressMetrics.totalSims > 0) {
							rounds = Math.ceil(Math.log(progressMetrics.totalSims / 20) / Math.log(2)) + 1;
							currentRound = 1;
						}
						if (progressMetrics.totalSims < lastTotal) {
							currentRound += 1;
							simStart = new Date().getTime();
						}
					}

					this.setSimProgress(progressMetrics, iterPerSecond, currentRound, rounds, combinations);
					lastTotal = progressMetrics.totalSims;
				});
			} catch (error) {
				console.error(error);
			} finally {
				isRunning = false;
				if (!waitAbort) bulkSimButton.disabled = false;
				// reset state
				this.pendingDiv.classList.add('hide');
				this.leftPanel.classList.remove('blurred');
				this.rightPanel.classList.remove('blurred');

				this.pendingResults.hideAll();
				bulkSimButton.classList.remove('disabled');
				bulkSimButton.replaceChildren(...defaultState);
			}
		});

		const importButton = (
			<button className="btn btn-secondary w-100 bulk-settings-button">
				<i className="fa fa-download" /> Import From Bags
			</button>
		) as HTMLButtonElement;
		importButton.addEventListener('click', () => new BulkGearJsonImporter(this.simUI.rootElem, this.simUI, this).open());

		const importFavsButton = (
			<button className="btn btn-secondary w-100 bulk-settings-button">
				<i className="fa fa-download" /> Import Favorites
			</button>
		);
		importFavsButton.addEventListener('click', () => {
			const filters = this.simUI.player.sim.getFilters();
			const items = filters.favoriteItems.map(itemID => ItemSpec.create({ id: itemID }));
			this.addItems(items);
		});

		const searchInputRef = ref<HTMLInputElement>();
		const searchResultsRef = ref<HTMLUListElement>();
		const searchWrapper = (
			<div className="search-wrapper hide">
				<input ref={searchInputRef} type="text" placeholder="Search..." className="batch-search-input form-control hide" />
				<ul ref={searchResultsRef} className="batch-search-results hide"></ul>
			</div>
		);

		let allItems = Array<UIItem>();

		searchInputRef.value?.addEventListener('keyup', event => {
			if (event.key == 'Enter') {
				const toAdd = Array<ItemSpec>();
				searchResultsRef.value?.childNodes.forEach(node => {
					const strID = (node as HTMLElement).getAttribute('data-item-id');
					if (strID != null) {
						toAdd.push(ItemSpec.create({ id: Number.parseInt(strID) }));
					}
				});
				this.addItems(toAdd);
			}
		});

		searchInputRef.value?.addEventListener('input', _event => {
			const searchString = searchInputRef.value?.value || '';

			if (!searchString.length) {
				searchResultsRef.value?.replaceChildren();
				searchResultsRef.value?.classList.add('hide');
				return;
			}

			const pieces = searchString.split(' ');
			const items = <></>;

			allItems.forEach(item => {
				let matched = true;
				const lcName = item.name.toLowerCase();
				const lcSetName = item.setName.toLowerCase();

				pieces.forEach(piece => {
					const lcPiece = piece.toLowerCase();
					if (!lcName.includes(lcPiece) && !lcSetName.includes(lcPiece)) {
						matched = false;
						return false;
					}
					return true;
				});

				if (matched) {
					const itemRef = ref<HTMLLIElement>();
					const itemNameRef = ref<HTMLSpanElement>();
					items.appendChild(
						<li ref={itemRef} dataset={{ itemId: item.id.toString() }}>
							<span ref={itemNameRef}>{item.name}</span>
							{item.heroic && <span className="item-quality-uncommon">[H]</span>}
							{item.factionRestriction === UIItem_FactionRestriction.HORDE_ONLY && <span className="faction-horde">(H)</span>}
							{item.factionRestriction === UIItem_FactionRestriction.ALLIANCE_ONLY && <span className="faction-alliance">(A)</span>}
						</li>,
					);
					setItemQualityCssClass(itemNameRef.value!, item.quality);
					itemRef.value?.addEventListener('click', () => this.addItem(ItemSpec.create({ id: item.id })));
				}
			});
			searchResultsRef.value?.replaceChildren(items);
			searchResultsRef.value?.classList.remove('hide');
		});

		const searchButtonContents = (
			<>
				<i className="fa fa-search" /> Add Item
			</>
		);

		const searchButton = <button className="btn btn-secondary w-100 bulk-settings-button">{searchButtonContents.cloneNode(true)}</button>;
		searchButton.addEventListener('click', () => {
			if (searchInputRef.value?.classList.contains('hide')) {
				searchWrapper?.classList.remove('hide');
				searchButton.replaceChildren(<>Close Search Results</>);
				allItems = this.simUI.sim.db.getAllItems().filter(item => canEquipItem(this.simUI.player, item, undefined));
				searchInputRef.value?.classList.remove('hide');
				if (searchInputRef.value?.value) searchResultsRef.value?.classList.remove('hide');
				searchInputRef.value?.focus();
			} else {
				searchButton.replaceChildren(searchButtonContents.cloneNode(true));
				searchWrapper?.classList.add('hide');
				searchInputRef.value?.classList.add('hide');
				searchResultsRef.value?.replaceChildren();
				searchResultsRef.value?.classList.add('hide');
			}
		});

		const clearButton = <button className="btn btn-secondary w-100 bulk-settings-button">Clear all</button>;
		clearButton.addEventListener('click', () => {
			this.clearItems();
			resultsBlock.rootElem.hidden = true;
			resultsBlock.bodyElement.replaceChildren();
		});

		// Talents to sim
		const talentsContainerRef = ref<HTMLDivElement>();
		const talentsToSimDiv = (
			<div className={clsx('talents-picker-container', !this.simTalents && 'hide')}>
				<label className="mb-2">Pick talents to sim (will increase time to sim)</label>
				<div ref={talentsContainerRef} className="talents-container"></div>
			</div>
		);

		const dataStr = window.localStorage.getItem(this.simUI.getSavedTalentsStorageKey());

		let jsonData;
		try {
			if (dataStr !== null) {
				jsonData = JSON.parse(dataStr);
			}
		} catch (e) {
			console.warn('Invalid json for local storage value: ' + dataStr);
		}

		const handleToggle = (frag: HTMLElement, load: TalentLoadout) => {
			const chipDiv = frag.querySelector('.saved-data-set-chip');
			const exists = this.savedTalents.some(talent => talent.name === load.name); // Replace 'id' with your unique identifier

			// console.log('Exists:', exists);
			// console.log('Load Object:', load);
			// console.log('Saved Talents Before Update:', this.savedTalents);

			if (exists) {
				// If the object exists, find its index and remove it
				const indexToRemove = this.savedTalents.findIndex(talent => talent.name === load.name);
				this.savedTalents.splice(indexToRemove, 1);
				chipDiv?.classList.remove('active');
			} else {
				// If the object does not exist, add it
				this.savedTalents.push(load);
				chipDiv?.classList.add('active');
			}

			// console.log('Updated savedTalents:', this.savedTalents);
		};
		for (const name in jsonData) {
			try {
				const savedTalentLoadout = SavedTalents.fromJson(jsonData[name]);
				const loadout = {
					talentsString: savedTalentLoadout.talentsString,
					name: name,
				};

				const index = this.savedTalents.findIndex(talent => JSON.stringify(talent) === JSON.stringify(loadout));
				const talentFragment = (
					<div className={clsx('saved-data-set-chip badge rounded-pill', index !== -1 && 'active')}>
						<a href="javascript:void(0)" className="saved-data-set-name" attributes={{ role: 'button' }}>
							{name}
						</a>
					</div>
				) as HTMLDivElement;

				// console.log('Adding event for loadout', loadout);

				talentsContainerRef.value!.appendChild(talentFragment);
				talentFragment.addEventListener('click', () => handleToggle(talentFragment, loadout));
			} catch (e) {
				console.log(e);
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}

		//////////////////////
		////////////////////////////////////

		settingsBlock.bodyElement.appendChild(
			<>
				{bulkSimButton}
				{importButton}
				{importFavsButton}
				{searchButton}
				{searchWrapper}
				{clearButton}
				{talentsToSimDiv}
			</>,
		);

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-fast-mode',
			label: 'Fast Mode',
			labelTooltip: 'Fast mode reduces accuracy but will run faster.',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.fastMode,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.fastMode = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-combinations',
			label: 'Combinations',
			labelTooltip:
				'When checked bulk simulator will create all possible combinations of the items. When disabled trinkets and rings will still run all combinations becausee they have two slots to fill each.',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.doCombos,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.doCombos = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-auto-enchant',
			label: 'Auto Enchant',
			labelTooltip: 'When checked bulk simulator apply the current enchant for a slot to each replacement item it can.',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.autoEnchant,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.autoEnchant = value;
			},
		});

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-sim-talents',
			label: 'Sim Talents',
			labelTooltip: 'When checked bulk simulator will sim chosen talent setups. Warning, it might cause the bulk sim to run for a lot longer',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.simTalents,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.simTalents = value;
				talentsToSimDiv.classList[value ? 'remove' : 'add']('hide');
			},
		});
	}

	private setSimProgress(progress: ProgressMetrics, iterPerSecond: number, currentRound: number, rounds: number, combinations: number) {
		const secondsRemain = ((progress.totalIterations - progress.completedIterations) / iterPerSecond).toFixed();

		this.pendingResults.setContent(
			<div className="results-sim">
				<div>{combinations} total combinations.</div>
				<div>
					{rounds > 0 && (
						<>
							{currentRound} / {rounds} refining rounds
						</>
					)}
				</div>
				<div>
					{progress.completedSims} / {progress.totalSims}
					<br />
					simulations complete
				</div>
				<div>
					{progress.completedIterations} / {progress.totalIterations}
					<br />
					iterations complete
				</div>
				<div>{secondsRemain} seconds remaining.</div>
			</div>,
		);
	}
}
