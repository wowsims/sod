import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { BaseModal } from '../core/components/base_modal.jsx';
import { Component } from '../core/components/component.js';
import { EnumPicker } from '../core/components/enum_picker.js';
import { MAX_PARTY_SIZE, Party } from '../core/party.js';
import { Player } from '../core/player.js';
import { Player as PlayerProto } from '../core/proto/api.js';
import { Class, Faction, Profession, Spec } from '../core/proto/common.js';
import { BalanceDruid_Options as BalanceDruidOptions } from '../core/proto/druid.js';
import { cssClassForClass, isTankSpec, newUnitReference, playerToSpec, specToClass } from '../core/proto_utils/utils.js';
import { Raid } from '../core/raid.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { formatDeltaTextElem, getEnumValues } from '../core/utils.js';
import { playerPresets, specSimFactories } from './presets.js';
import { RaidSimUI } from './raid_sim_ui.js';

const NEW_PLAYER = -1;

enum DragType {
	None,
	New,
	Move,
	Swap,
	Copy,
}

export class RaidPicker extends Component {
	readonly raidSimUI: RaidSimUI;
	readonly raid: Raid;
	readonly partyPickers: Array<PartyPicker>;
	readonly newPlayerPicker: NewPlayerPicker;
	readonly playerEditorModal: PlayerEditorModal<Spec>;

	// Hold data about the player being dragged while the drag is happening.
	currentDragPlayer: Player<any> | null = null;
	currentDragPlayerFromIndex: number = NEW_PLAYER;
	currentDragType: DragType = DragType.New;

	// Hold data about the party being dragged while the drag is happening.
	currentDragParty: PartyPicker | null = null;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI) {
		super(parent, 'raid-picker-root');
		this.raidSimUI = raidSimUI;
		this.raid = raidSimUI.sim.raid;

		const raidControls = document.createElement('div');
		raidControls.classList.add('raid-controls');
		this.rootElem.appendChild(raidControls);

		this.newPlayerPicker = new NewPlayerPicker(this.rootElem, this);
		this.playerEditorModal = new PlayerEditorModal();

		new EnumPicker<Raid>(raidControls, this.raidSimUI.sim.raid, {
			id: 'raid-picker-size',
			label: 'Raid Size',
			labelTooltip: 'Number of players participating in the sim.',
			values: [
				{ name: '5', value: 1 },
				{ name: '10', value: 2 },
				{ name: '25', value: 5 },
				{ name: '40', value: 8 },
			],
			changedEvent: (raid: Raid) => raid.numActivePartiesChangeEmitter,
			getValue: (raid: Raid) => raid.getNumActiveParties(),
			setValue: (eventID: EventID, raid: Raid, newValue: number) => {
				raid.setNumActiveParties(eventID, newValue);
			},
		});

		new EnumPicker<NewPlayerPicker>(raidControls, this.newPlayerPicker, {
			id: 'raid-picker-faction',
			label: 'Default Faction',
			labelTooltip: 'Default faction for newly-created players.',
			values: [
				{ name: 'Alliance', value: Faction.Alliance },
				{ name: 'Horde', value: Faction.Horde },
			],
			changedEvent: () => this.raid.sim.factionChangeEmitter,
			getValue: () => this.raid.sim.getFaction(),
			setValue: (eventID: EventID, _picker: NewPlayerPicker, newValue: Faction) => {
				this.raid.sim.setFaction(eventID, newValue);
			},
		});

		const latestPhaseWithAllPresets = Math.min(
			...playerPresets.map(preset => Math.max(...Object.keys(preset.defaultGear[Faction.Alliance]).map(k => parseInt(k)))),
		);
		new EnumPicker<NewPlayerPicker>(raidControls, this.newPlayerPicker, {
			id: 'raid-picker-gear',
			label: 'Default Gear',
			labelTooltip: 'Newly-created players will start with approximate BIS gear from this phase.',
			values: [...Array(latestPhaseWithAllPresets).keys()].map(val => {
				const phase = val + 1;
				return { name: 'Phase ' + phase, value: phase };
			}),
			changedEvent: () => this.raid.sim.phaseChangeEmitter,
			getValue: () => this.raid.sim.getPhase(),
			setValue: (eventID: EventID, _picker: NewPlayerPicker, newValue: number) => {
				this.raid.sim.setPhase(eventID, newValue);
			},
		});

		const partiesContainer = document.createElement('div');
		partiesContainer.classList.add('parties-container');
		this.rootElem.appendChild(partiesContainer);

		this.partyPickers = this.raid.getParties().map((party, i) => new PartyPicker(partiesContainer, party, i, this));

		const updateActiveParties = () => {
			this.partyPickers.forEach(partyPicker => {
				if (partyPicker.index < this.raidSimUI.sim.raid.getNumActiveParties()) {
					partyPicker.rootElem.classList.add('active');
				} else {
					partyPicker.rootElem.classList.remove('active');
				}
			});
		};
		this.raidSimUI.sim.raid.numActivePartiesChangeEmitter.on(updateActiveParties);
		updateActiveParties();

		this.rootElem.ondragend = () => {
			// Uncomment to remove player when dropped 'off' the raid.
			//if (this.currentDragPlayerFromIndex != NEW_PLAYER) {
			//	const playerPicker = this.getPlayerPicker(this.currentDragPlayerFromIndex);
			//	playerPicker.setPlayer(null, null, DragType.None);
			//}

			this.clearDragPlayer();
			this.clearDragParty();
		};
	}

	getCurrentFaction(): Faction {
		return this.raid.sim.getFaction();
	}

	getCurrentPhase(): number {
		return this.raid.sim.getPhase();
	}

	getPlayerPicker(raidIndex: number): PlayerPicker {
		return this.partyPickers[Math.floor(raidIndex / MAX_PARTY_SIZE)].playerPickers[raidIndex % MAX_PARTY_SIZE];
	}

	getPlayerPickers(): Array<PlayerPicker> {
		return [...new Array(25).keys()].map(i => this.getPlayerPicker(i));
	}

	setDragPlayer(player: Player<any>, fromIndex: number, type: DragType) {
		this.clearDragPlayer();

		this.currentDragPlayer = player;
		this.currentDragPlayerFromIndex = fromIndex;
		this.currentDragType = type;

		if (fromIndex != NEW_PLAYER) {
			const playerPicker = this.getPlayerPicker(fromIndex);
			playerPicker.rootElem.classList.add('dragfrom');
		}
	}

	clearDragPlayer() {
		if (this.currentDragPlayerFromIndex != NEW_PLAYER) {
			const playerPicker = this.getPlayerPicker(this.currentDragPlayerFromIndex);
			playerPicker.rootElem.classList.remove('dragfrom');
		}

		this.currentDragPlayer = null;
		this.currentDragPlayerFromIndex = NEW_PLAYER;
		this.currentDragType = DragType.New;
	}

	setDragParty(party: PartyPicker) {
		this.currentDragParty = party;
		party.rootElem.classList.add('dragfrom');
	}
	clearDragParty() {
		if (this.currentDragParty) {
			this.currentDragParty.rootElem.classList.remove('dragfrom');
			this.currentDragParty = null;
		}
	}
}

export class PartyPicker extends Component {
	readonly party: Party;
	readonly index: number;
	readonly raidPicker: RaidPicker;
	readonly playerPickers: Array<PlayerPicker>;

	constructor(parent: HTMLElement, party: Party, index: number, raidPicker: RaidPicker) {
		super(parent, 'party-picker-root');
		this.party = party;
		this.index = index;
		this.raidPicker = raidPicker;

		this.rootElem.setAttribute('draggable', 'true');
		this.rootElem.innerHTML = `
			<div class="party-header">
				<label class="party-label form-label">Group ${index + 1}</label>
				<div class="party-results">
					<span class="party-results-dps"></span>
					<span class="party-results-reference-delta"></span>
				</div>
			</div>
			<div class="players-container">
			</div>
		`;

		const playersContainer = this.rootElem.getElementsByClassName('players-container')[0] as HTMLDivElement;
		this.playerPickers = [...Array(MAX_PARTY_SIZE).keys()].map(i => new PlayerPicker(playersContainer, this, i));

		const dpsResultElem = this.rootElem.getElementsByClassName('party-results-dps')[0] as HTMLElement;
		const referenceDeltaElem = this.rootElem.getElementsByClassName('party-results-reference-delta')[0] as HTMLElement;

		this.raidPicker.raidSimUI.referenceChangeEmitter.on(() => {
			const currentData = this.raidPicker.raidSimUI.getCurrentData();
			const referenceData = this.raidPicker.raidSimUI.getReferenceData();

			const partyDps = currentData?.simResult.raidMetrics.parties[this.index]?.dps.avg || 0;
			const referenceDps = referenceData?.simResult.raidMetrics.parties[this.index]?.dps.avg || 0;

			if (partyDps == 0 && referenceDps == 0) {
				dpsResultElem.textContent = '';
				referenceDeltaElem.textContent = '';
				return;
			}

			dpsResultElem.textContent = `${partyDps.toFixed(1)} DPS`;

			if (!referenceData) {
				referenceDeltaElem.textContent = '';
				return;
			}

			formatDeltaTextElem(referenceDeltaElem, referenceDps, partyDps, 1);
		});

		this.rootElem.ondragstart = event => {
			if (event.target == this.rootElem) {
				event.dataTransfer!.dropEffect = 'move';
				event.dataTransfer!.effectAllowed = 'all';
				this.raidPicker.setDragParty(this);
			}
		};

		let dragEnterCounter = 0;
		this.rootElem.ondragenter = event => {
			event.preventDefault();
			if (!this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter++;
			this.rootElem.classList.add('dragto');
		};
		this.rootElem.ondragleave = event => {
			event.preventDefault();
			if (!this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter--;
			if (dragEnterCounter <= 0) {
				this.rootElem.classList.remove('dragto');
			}
		};
		this.rootElem.ondragover = event => {
			event.preventDefault();
		};
		this.rootElem.ondrop = event => {
			if (!this.raidPicker.currentDragParty) {
				return;
			}

			event.preventDefault();
			dragEnterCounter = 0;
			this.rootElem.classList.remove('dragto');

			const eventID = TypedEvent.nextEventID();
			TypedEvent.freezeAllAndDo(() => {
				const srcPartyPicker = this.raidPicker.currentDragParty!;

				for (let i = 0; i < MAX_PARTY_SIZE; i++) {
					const srcPlayerPicker = srcPartyPicker.playerPickers[i]!;
					const dstPlayerPicker = this.playerPickers[i]!;

					const srcPlayer = srcPlayerPicker.player;
					const dstPlayer = dstPlayerPicker.player;

					srcPlayerPicker.setPlayer(eventID, dstPlayer, DragType.Swap);
					dstPlayerPicker.setPlayer(eventID, srcPlayer, DragType.Swap);
				}
			});

			this.raidPicker.clearDragParty();
		};
	}
}

export class PlayerPicker extends Component {
	// Index of this player within its party (0-4).
	readonly index: number;

	// Index of this player within the whole raid (0-24).
	readonly raidIndex: number;

	player: Player<any> | null;

	readonly partyPicker: PartyPicker;
	readonly raidPicker: RaidPicker;

	private labelElem: HTMLElement | null;
	private iconElem: HTMLImageElement | null;
	private nameElem: HTMLInputElement | null;
	private resultsElem: HTMLElement | null;
	private dpsResultElem: HTMLElement | null;
	private referenceDeltaElem: HTMLElement | null;
	// Can be used to remove any events in addEventListener
	// https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#add_an_abortable_listener
	public abortController: AbortController;
	public signal: AbortSignal;

	constructor(parent: HTMLElement, partyPicker: PartyPicker, index: number) {
		super(parent, 'player-picker-root');
		this.abortController = new AbortController();
		this.signal = this.abortController.signal;
		this.index = index;
		this.raidIndex = partyPicker.index * MAX_PARTY_SIZE + index;
		this.player = null;
		this.partyPicker = partyPicker;
		this.raidPicker = partyPicker.raidPicker;

		this.labelElem = null;
		this.iconElem = null;
		this.nameElem = null;
		this.resultsElem = null;
		this.dpsResultElem = null;
		this.referenceDeltaElem = null;

		this.rootElem.classList.add('player');

		this.partyPicker.party.compChangeEmitter.on(eventID => {
			const newPlayer = this.partyPicker.party.getPlayer(this.index);
			if (newPlayer != this.player) this.setPlayer(eventID, newPlayer, DragType.None);
		});

		this.raidPicker.raidSimUI.referenceChangeEmitter.on(() => {
			const currentData = this.raidPicker.raidSimUI.getCurrentData();
			const referenceData = this.raidPicker.raidSimUI.getReferenceData();

			const playerDps = currentData?.simResult.getPlayerWithRaidIndex(this.raidIndex)?.dps.avg || 0;
			const referenceDps = referenceData?.simResult.getPlayerWithRaidIndex(this.raidIndex)?.dps.avg || 0;

			if (this.player) {
				this.resultsElem?.classList.remove('hide');
				(this.dpsResultElem as HTMLElement).textContent = `${playerDps.toFixed(1)} DPS`;

				if (referenceData) formatDeltaTextElem(this.referenceDeltaElem as HTMLElement, referenceDps, playerDps, 1);
			}
		});

		let dragEnterCounter = 0;
		this.rootElem.ondragenter = event => {
			event.preventDefault();
			if (this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter++;
			this.rootElem.classList.add('dragto');
		};
		this.rootElem.ondragleave = event => {
			event.preventDefault();
			if (this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter--;
			if (dragEnterCounter <= 0) {
				this.rootElem.classList.remove('dragto');
			}
		};
		this.rootElem.ondragover = event => event.preventDefault();
		this.rootElem.ondrop = event => {
			if (this.raidPicker.currentDragParty) {
				return;
			}
			const dropData = event.dataTransfer!.getData('text/plain');

			event.preventDefault();
			dragEnterCounter = 0;
			this.rootElem.classList.remove('dragto');

			const eventID = TypedEvent.nextEventID();
			TypedEvent.freezeAllAndDo(() => {
				if (this.raidPicker.currentDragPlayer == null && dropData.length == 0) {
					return;
				}

				if (this.raidPicker.currentDragPlayerFromIndex == this.raidIndex) {
					this.raidPicker.clearDragPlayer();
					return;
				}

				const dragType = this.raidPicker.currentDragType;

				if (this.raidPicker.currentDragPlayerFromIndex != NEW_PLAYER) {
					const fromPlayerPicker = this.raidPicker.getPlayerPicker(this.raidPicker.currentDragPlayerFromIndex);
					if (dragType == DragType.Swap) {
						fromPlayerPicker.setPlayer(eventID, this.player, dragType);
					} else if (dragType == DragType.Move) {
						fromPlayerPicker.setPlayer(eventID, null, dragType);
					}
				} else if (this.raidPicker.currentDragPlayer == null) {
					// This would be a copy from another window.
					const binary = atob(dropData);
					const bytes = new Uint8Array(binary.length);
					for (let i = 0; i < bytes.length; i++) {
						bytes[i] = binary.charCodeAt(i);
					}
					const playerProto = PlayerProto.fromBinary(bytes);

					const localPlayer = new Player(playerToSpec(playerProto), this.raidPicker.raidSimUI.sim);
					localPlayer.fromProto(eventID, playerProto);
					this.raidPicker.currentDragPlayer = localPlayer;
				}

				if (dragType == DragType.Copy) {
					this.setPlayer(eventID, this.raidPicker.currentDragPlayer!.clone(eventID), dragType);
				} else {
					this.setPlayer(eventID, this.raidPicker.currentDragPlayer, dragType);
				}

				this.raidPicker.clearDragPlayer();
			});
		};

		this.update();
	}

	setPlayer(eventID: EventID, newPlayer: Player<any> | null, dragType: DragType) {
		if (newPlayer == this.player) {
			return;
		}

		TypedEvent.freezeAllAndDo(() => {
			this.player = newPlayer;
			if (newPlayer) {
				this.partyPicker.party.setPlayer(eventID, this.index, newPlayer);

				if (dragType == DragType.New) {
					applyNewPlayerAssignments(eventID, newPlayer, this.raidPicker.raid);
				}
			} else {
				this.partyPicker.party.setPlayer(eventID, this.index, newPlayer);
				this.partyPicker.party.compChangeEmitter.emit(eventID);
			}
		});

		this.update();
	}

	private update() {
		if (this.player == null) {
			this.rootElem.className = 'player-picker-root player';
			this.rootElem.innerHTML = '';

			this.labelElem = null;
			this.iconElem = null;
			this.nameElem = null;
			this.resultsElem = null;
			this.dpsResultElem = null;
			this.referenceDeltaElem = null;
		} else {
			const classCssClass = cssClassForClass(this.player.getClass());

			this.rootElem.className = `player-picker-root player bg-${classCssClass}-dampened`;
			this.rootElem.innerHTML = `
				<div class="player-label">
					<img class="player-icon" src="${this.player.getSpecIcon()}" draggable="true" />
					<div class="player-details">
						<input
							class="player-name text-${classCssClass}"
							type="text"
							value="${this.player.getName()}"
							spellcheck="false"
							maxlength="15"
						/>
						<div class="player-results hide">
							<span class="player-results-dps"></span>
							<span class="player-results-reference-delta"></span>
						</div>
					</div>
				</div>
				<div class="player-options">
					<a
						href="javascript:void(0)"
						class="player-edit"
						role="button"
						data-tippy-content="Click to Edit"
					>
						<i class="fa fa-edit fa-lg"></i>
					</a>
					<a
						href="javascript:void(0)"
						class="player-copy link-warning"
						role="button"
						draggable="true"
						data-tippy-content="Drag to Copy"
					>
						<i class="fa fa-copy fa-lg"></i>
					</a>
					<a
						href="javascript:void(0)"
						class="player-delete link-danger"
						role="button"
						data-tippy-content="Click to Delete"
					>
						<i class="fa fa-times fa-lg"></i>
					</a>
				</div>
			`;

			this.labelElem = this.rootElem.querySelector<HTMLElement>('.player-label')!;
			this.iconElem = this.rootElem.querySelector<HTMLImageElement>('.player-icon')!;
			this.nameElem = this.rootElem.querySelector<HTMLInputElement>('.player-name')!;
			this.resultsElem = this.rootElem.querySelector<HTMLElement>('.player-results')!;
			this.dpsResultElem = this.rootElem.querySelector<HTMLElement>('.player-results-dps')!;
			this.referenceDeltaElem = this.rootElem.querySelector<HTMLElement>('.player-results-reference-delta')!;

			this.bindPlayerEvents();
		}
	}

	private bindPlayerEvents() {
		this.nameElem?.addEventListener('input', _event => {
			this.player?.setName(TypedEvent.nextEventID(), this.nameElem?.value || '');
		});

		this.nameElem?.addEventListener('mousedown', _event => {
			this.partyPicker.rootElem.setAttribute('draggable', 'false');
		});

		this.nameElem?.addEventListener('mouseup', _event => {
			this.partyPicker.rootElem.setAttribute('draggable', 'true');
		});

		const emptyName = 'Unnamed';
		this.nameElem?.addEventListener('focusout', _event => {
			if (this.nameElem && !this.nameElem.value) {
				this.nameElem.value = emptyName;
				this.player?.setName(TypedEvent.nextEventID(), emptyName);
			}
		});

		const dragStart = (event: DragEvent, type: DragType) => {
			if (this.player == null) {
				event.preventDefault();
				return;
			}

			event.dataTransfer!.dropEffect = 'move';
			event.dataTransfer!.effectAllowed = 'all';

			if (this.player) {
				const playerDataProto = this.player.toProto(true);
				event.dataTransfer!.setData('text/plain', btoa(String.fromCharCode(...PlayerProto.toBinary(playerDataProto))));
			}

			this.raidPicker.setDragPlayer(this.player, this.raidIndex, type);
		};

		const editElem = this.rootElem.querySelector<HTMLElement>('.player-edit')!;
		const copyElem = this.rootElem.querySelector<HTMLElement>('.player-copy')!;
		const deleteElem = this.rootElem.querySelector<HTMLElement>('.player-delete')!;

		const editTooltip = tippy(editElem);
		const copyTooltip = tippy(copyElem);
		const deleteTooltip = tippy(deleteElem);

		const onIconDragStartHandler = (event: DragEvent) => {
			event.dataTransfer!.setDragImage(this.rootElem, 20, 20);
			dragStart(event, DragType.Swap);
		};
		this.iconElem?.addEventListener('dragstart', onIconDragStartHandler, { signal: this.signal });

		const onEditClickHandler = () => {
			if (this.player) this.raidPicker.playerEditorModal.openEditor(this.player);
		};
		editElem.addEventListener('click', onEditClickHandler, { signal: this.signal });

		const onCopyDragStartHandler = (event: DragEvent) => {
			event.dataTransfer!.setDragImage(this.rootElem, 20, 20);
			dragStart(event, DragType.Copy);
		};
		copyElem.addEventListener('dragstart', onCopyDragStartHandler, { signal: this.signal });

		const onDeleteClickHandler = () => {
			this.setPlayer(TypedEvent.nextEventID(), null, DragType.None);
			this.dispose();
		};
		deleteElem.addEventListener('click', onDeleteClickHandler, { signal: this.signal });

		this.addOnDisposeCallback(() => {
			editTooltip?.destroy();
			copyTooltip?.destroy();
			deleteTooltip?.destroy();
		});
	}
}

class PlayerEditorModal<SpecType extends Spec> extends BaseModal {
	playerEditorRoot: HTMLDivElement;

	constructor() {
		super(document.body, 'player-editor-modal', {
			closeButton: { fixed: true },
			header: false,
			disposeOnClose: false,
		});

		const playerEditorElemRef = ref<HTMLDivElement>();
		const playerEditorElem = (<div ref={playerEditorElemRef} className="player-editor within-raid-sim"></div>) as HTMLDivElement;

		this.rootElem.id = 'playerEditorModal';
		this.body.appendChild(playerEditorElem);

		this.playerEditorRoot = playerEditorElemRef.value!;
	}

	openEditor(player: Player<SpecType>) {
		this.setData(player);
		super.open();
	}

	setData(player: Player<SpecType>) {
		this.playerEditorRoot.innerHTML = '';
		specSimFactories[player.spec]?.(this.playerEditorRoot!, player);
	}
}

class NewPlayerPicker extends Component {
	readonly raidPicker: RaidPicker;

	constructor(parent: HTMLElement, raidPicker: RaidPicker) {
		super(parent, 'new-player-picker-root');
		this.raidPicker = raidPicker;

		getEnumValues(Class).forEach(wowClass => {
			if (wowClass == Class.ClassUnknown) {
				return;
			}

			const matchingPresets = playerPresets.filter(preset => specToClass[preset.spec] == wowClass);
			if (matchingPresets.length == 0) {
				return;
			}

			const classPresetsContainer = document.createElement('div');
			classPresetsContainer.classList.add('class-presets-container', `bg-${cssClassForClass(wowClass as Class)}-dampened`);
			this.rootElem.appendChild(classPresetsContainer);

			matchingPresets.forEach(matchingPreset => {
				const presetElemFragment = document.createElement('fragment');
				presetElemFragment.innerHTML = `
					<a
						href="javascript:void(0)"
						role="button"
						draggable="true"
						data-tippy-content="${matchingPreset.tooltip}"
					>
						<img class="preset-picker-icon player-icon" src="${matchingPreset.iconUrl}"/>
					</a>
				`;
				const presetElem = presetElemFragment.children[0] as HTMLElement;
				classPresetsContainer.appendChild(presetElem);

				tippy(presetElem);

				presetElem.ondragstart = event => {
					const eventID = TypedEvent.nextEventID();
					TypedEvent.freezeAllAndDo(() => {
						const dragImage = new Image();
						dragImage.src = matchingPreset.iconUrl;
						event.dataTransfer!.setDragImage(dragImage, 30, 30);
						event.dataTransfer!.setData('text/plain', '');
						event.dataTransfer!.dropEffect = 'copy';

						const newPlayer = new Player(matchingPreset.spec, this.raidPicker.raid.sim);
						newPlayer.applySharedDefaults(eventID);
						newPlayer.setRace(eventID, matchingPreset.defaultFactionRaces[this.raidPicker.getCurrentFaction()]);
						newPlayer.setTalentsString(eventID, matchingPreset.talents.talentsString);
						newPlayer.setSpecOptions(eventID, matchingPreset.specOptions);
						newPlayer.setConsumes(eventID, matchingPreset.consumes);
						newPlayer.setName(eventID, matchingPreset.defaultName);
						newPlayer.setProfession1(eventID, matchingPreset.otherDefaults?.profession1 || Profession.Engineering);
						newPlayer.setProfession2(eventID, matchingPreset.otherDefaults?.profession2 || Profession.Enchanting);
						newPlayer.setDistanceFromTarget(eventID, matchingPreset.otherDefaults?.distanceFromTarget || 0);

						// Need to wait because the gear might not be loaded yet.
						this.raidPicker.raid.sim.waitForInit().then(() => {
							newPlayer.setGear(
								eventID,
								this.raidPicker.raid.sim.db.lookupEquipmentSpec(
									matchingPreset.defaultGear[this.raidPicker.getCurrentFaction()][this.raidPicker.getCurrentPhase()],
								),
							);
						});

						this.raidPicker.setDragPlayer(newPlayer, NEW_PLAYER, DragType.New);
					});
				};
			});
		});
	}
}

function applyNewPlayerAssignments(eventID: EventID, newPlayer: Player<any>, raid: Raid) {
	if (isTankSpec(newPlayer.spec)) {
		const tanks = raid.getTanks();
		const emptyIdx = tanks.findIndex(tank => raid.getPlayerFromUnitReference(tank) == null);
		if (emptyIdx == -1) {
			if (tanks.length < 3) {
				raid.setTanks(eventID, tanks.concat([newPlayer.makeUnitReference()]));
			}
		} else {
			tanks[emptyIdx] = newPlayer.makeUnitReference();
			raid.setTanks(eventID, tanks);
		}
	}

	// Spec-specific assignments. For most cases, default to buffing self.
	if (newPlayer.spec == Spec.SpecBalanceDruid) {
		const newOptions = newPlayer.getSpecOptions() as BalanceDruidOptions;
		newOptions.innervateTarget = newUnitReference(newPlayer.getRaidIndex());
		newPlayer.setSpecOptions(eventID, newOptions);
	}
}
