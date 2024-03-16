import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { Spec } from '../../proto/common';
import { SavedTalents } from '../../proto/ui';
import { classTalentsConfig } from '../../talents/factory';
import { TalentsPicker } from '../../talents/talents_picker';
import { EventID, TypedEvent } from '../../typed_event';
import { SavedDataManager } from '../saved_data_manager';
import { SimTab } from '../sim_tab';

export class TalentsTab extends SimTab {
	protected simUI: IndividualSimUI<Spec>;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, simUI, { identifier: 'talents-tab', title: 'Talents' });
		this.simUI = simUI;

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('talents-tab-left', 'tab-panel-left');
		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('talents-tab-right', 'tab-panel-right', 'within-raid-sim-hide');

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();
	}

	protected buildTabContent() {
		this.buildTalentsPicker(this.leftPanel);
		this.buildSavedTalentsPicker();
	}

	private buildTalentsPicker(parentElem: HTMLElement) {
		new TalentsPicker(parentElem, this.simUI.player, {
			klass: this.simUI.player.getClass(),
			trees: classTalentsConfig[this.simUI.player.getClass()],
			changedEvent: (player: Player<any>) => player.talentsChangeEmitter,
			getValue: (player: Player<any>) => player.getTalentsString(),
			setValue: (eventID: EventID, player: Player<any>, newValue: string) => {
				player.setTalentsString(eventID, newValue);
			},
			pointsPerRow: 5,
		});
	}

	private buildSavedTalentsPicker() {
		const savedTalentsManager = new SavedDataManager<Player<any>, SavedTalents>(this.rightPanel, this.simUI.player, {
			label: 'Talents',
			header: { title: 'Saved Talents' },
			storageKey: this.simUI.getSavedTalentsStorageKey(),
			getData: (player: Player<any>) =>
				SavedTalents.create({
					talentsString: player.getTalentsString(),
				}),
			setData: (eventID: EventID, player: Player<any>, newTalents: SavedTalents) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setTalentsString(eventID, newTalents.talentsString);
				});
			},
			changeEmitters: [this.simUI.player.talentsChangeEmitter, this.simUI.player.levelChangeEmitter],
			equals: (a: SavedTalents, b: SavedTalents) => SavedTalents.equals(a, b),
			toJson: (a: SavedTalents) => SavedTalents.toJson(a),
			fromJson: (obj: any) => SavedTalents.fromJson(obj),
		});

		this.simUI.sim.waitForInit().then(() => {
			savedTalentsManager.loadUserData();
			this.simUI.individualConfig.presets.talents.forEach(config => {
				config.isPreset = true;
				savedTalentsManager.addSavedData({
					name: config.name,
					isPreset: true,
					data: config.data,
					enableWhen: config.enableWhen,
				});
			});
		});
	}
}
