import { CharacterStats, StatMods } from './components/character_stats';
import { ContentBlock } from './components/content_block';
import { EmbeddedDetailedResults } from './components/detailed_results';
import { EncounterPickerConfig } from './components/encounter_picker';
import * as IconInputs from './components/icon_inputs';
import { BulkTab } from './components/individual_sim_ui/bulk_tab';
import { Individual60UEPExporter, IndividualJsonExporter, IndividualLinkExporter, IndividualPawnEPExporter } from './components/individual_sim_ui/exporters';
import { GearTab } from './components/individual_sim_ui/gear_tab';
import { Individual60UImporter, IndividualAddonImporter, IndividualJsonImporter, IndividualLinkImporter } from './components/individual_sim_ui/importers';
import { RotationTab } from './components/individual_sim_ui/rotation_tab';
import { SettingsTab } from './components/individual_sim_ui/settings_tab';
import { TalentsTab } from './components/individual_sim_ui/talents_tab';
import * as InputHelpers from './components/input_helpers';
import { addRaidSimAction, RaidSimResultsManager } from './components/raid_sim_action';
import { SavedDataConfig } from './components/saved_data_manager';
import { addStatWeightsAction, StatWeightActionSettings } from './components/stat_weights_action';
import { GLOBAL_DISPLAY_PSEUDO_STATS, GLOBAL_DISPLAY_STATS, GLOBAL_EP_PSEUDOSTATS, GLOBAL_EP_STATS, LEVEL_THRESHOLDS } from './constants/other';
import { SimSettingCategories } from './constants/sim_settings';
import * as Tooltips from './constants/tooltips';
import { simLaunchStatuses } from './launched_sims';
import { Player, PlayerConfig, registerSpecConfig as registerPlayerConfig } from './player';
import { PresetBuild, PresetGear, PresetRotation, PresetSettings } from './preset_utils';
import { StatWeightsResult } from './proto/api';
import { APLRotation, APLRotation_Type as APLRotationType } from './proto/apl';
import {
	Consumes,
	Debuffs,
	Encounter as EncounterProto,
	EquipmentSpec,
	Faction,
	HandType,
	IndividualBuffs,
	ItemSlot,
	ItemSwap,
	PartyBuffs,
	Profession,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
} from './proto/common';
import { IndividualSimSettings, SavedTalents } from './proto/ui';
import { professionNames } from './proto_utils/names';
import { Stats, UnitStat } from './proto_utils/stats';
import { getTalentPoints, isHealingSpec, isTankSpec, SpecOptions, SpecRotation, specToEligibleRaces, specToLocalStorageKey } from './proto_utils/utils';
import { SimUI, SimWarning } from './sim_ui';
import { EventID, TypedEvent } from './typed_event';

const SAVED_GEAR_STORAGE_KEY = '__savedGear__';
const SAVED_ROTATION_STORAGE_KEY = '__savedRotation__';
const SAVED_SETTINGS_STORAGE_KEY = '__savedSettings__';
const SAVED_TALENTS_STORAGE_KEY = '__savedTalents__';

export type InputConfig<ModObject> =
	| InputHelpers.TypedBooleanPickerConfig<ModObject>
	| InputHelpers.TypedNumberPickerConfig<ModObject>
	| InputHelpers.TypedEnumPickerConfig<ModObject>;

export interface InputSection {
	tooltip?: string;
	inputs: Array<InputConfig<Player<any>>>;
}

export interface OtherDefaults {
	profession1?: Profession;
	profession2?: Profession;
	distanceFromTarget?: number;
	channelClipDelay?: number;
}

export interface RaidSimPreset<SpecType extends Spec> {
	spec: Spec;
	talents: SavedTalents;
	specOptions: SpecOptions<SpecType>;
	consumes: Consumes;

	defaultName: string;
	defaultFactionRaces: Record<Faction, Race>;
	defaultGear: Record<Faction, Record<number, EquipmentSpec>>;
	otherDefaults?: OtherDefaults;

	tooltip: string;
	iconUrl: string;
}

export interface IndividualSimUIConfig<SpecType extends Spec> extends PlayerConfig<SpecType> {
	// Additional css class to add to the root element.
	cssClass: string;
	// Used to generate schemed components. E.g. 'shaman', 'druid', 'raid'
	cssScheme: string;

	knownIssues?: Array<string>;
	warnings?: Array<(simUI: IndividualSimUI<SpecType>) => SimWarning>;

	epStats: Array<Stat>;
	epPseudoStats?: Array<PseudoStat>;
	epReferenceStat: Stat;
	displayStats: Array<Stat>;
	displayPseudoStats: Array<PseudoStat>;
	modifyDisplayStats?: (player: Player<SpecType>) => StatMods;

	defaults: {
		race?: Race;
		gear: EquipmentSpec;
		epWeights: Stats;
		consumes: Consumes;

		rotationType?: APLRotationType;
		simpleRotation?: SpecRotation<SpecType>;

		talents: SavedTalents;
		specOptions: SpecOptions<SpecType>;

		raidBuffs: RaidBuffs;
		partyBuffs: PartyBuffs;
		individualBuffs: IndividualBuffs;

		debuffs: Debuffs;

		other?: OtherDefaults;

		itemSwap?: ItemSwap;
	};

	playerInputs?: InputSection;
	playerIconInputs: Array<IconInputs.IconInputConfig<Player<SpecType>, any>>;
	petConsumeInputs?: Array<IconInputs.IconInputConfig<Player<SpecType>, any>>;
	rotationInputs?: InputSection;
	rotationIconInputs?: Array<IconInputs.IconInputConfig<Player<any>, any>>;
	includeBuffDebuffInputs: Array<any>;
	excludeBuffDebuffInputs: Array<any>;
	otherInputs: InputSection;
	// Currently, many classes don't support item swapping, and only in certain slots.
	// So enable it only where it is supported.
	itemSwapSlots?: Array<ItemSlot>;

	// For when extra sections are needed (e.g. Shaman totems)
	customSections?: Array<(parentElem: HTMLElement, simUI: IndividualSimUI<SpecType>) => ContentBlock>;

	encounterPicker: EncounterPickerConfig;

	presets: {
		gear: Array<PresetGear>;
		talents: Array<SavedDataConfig<Player<any>, SavedTalents>>;
		rotations: Array<PresetRotation>;
		settings?: Array<PresetSettings>;
		builds?: Array<PresetBuild>;
	};

	raidSimPresets: Array<RaidSimPreset<SpecType>>;
}

export function registerSpecConfig<SpecType extends Spec>(spec: SpecType, config: IndividualSimUIConfig<SpecType>): IndividualSimUIConfig<SpecType> {
	registerPlayerConfig(spec, config);
	return config;
}

export const itemSwapEnabledSpecs: Array<Spec> = [];

export interface Settings {
	raidBuffs: RaidBuffs;
	partyBuffs: PartyBuffs;
	individualBuffs: IndividualBuffs;
	consumes: Consumes;
	race: Race;
	professions?: Array<Profession>;
}

// Extended shared UI for all individual player sims.
export abstract class IndividualSimUI<SpecType extends Spec> extends SimUI {
	readonly player: Player<SpecType>;
	readonly individualConfig: IndividualSimUIConfig<SpecType>;
	private readonly statWeightActionSettings: StatWeightActionSettings;
	
	private raidSimResultsManager: RaidSimResultsManager | null;

	prevEpIterations: number;
	prevEpSimResult: StatWeightsResult | null;
	dpsRefStat?: Stat;
	healRefStat?: Stat;
	tankRefStat?: Stat;

	readonly bt: BulkTab;
	spec: any;

	constructor(parentElem: HTMLElement, player: Player<SpecType>, config: IndividualSimUIConfig<SpecType>) {
		super(parentElem, player.sim, {
			cssClass: config.cssClass,
			cssScheme: config.cssScheme,
			spec: player.spec,
			knownIssues: config.knownIssues,
			simStatus: simLaunchStatuses[player.spec],
			noticeText: (
				<span>
					We've begun working on Phase 8 changes! Because of this, some rune or other spell behavior may not be consistent with live servers due to
					PTR changes.
				</span>
			),
		});
		this.rootElem.classList.add('individual-sim-ui');
		this.player = player;
		this.individualConfig = config;
		this.raidSimResultsManager = null;
		this.prevEpIterations = 0;
		this.prevEpSimResult = null;
		this.statWeightActionSettings = new StatWeightActionSettings(this);

		if ((config.itemSwapSlots || []).length > 0 && !itemSwapEnabledSpecs.includes(player.spec)) {
			itemSwapEnabledSpecs.push(player.spec);
		}

		this.addWarning({
			updateOn: TypedEvent.onAny([this.player.gearChangeEmitter, this.player.professionChangeEmitter]),
			getContent: () => {
				const failedProfReqs = this.player.getGear().getFailedProfessionRequirements(this.player.getProfessions());
				if (failedProfReqs.length == 0) {
					return '';
				}

				return failedProfReqs.map(fpr => `${fpr.name} requires ${professionNames.get(fpr.requiredProfession)!}, but it is not selected.`);
			},
		});
		this.addWarning({
			updateOn: TypedEvent.onAny([this.player.talentsChangeEmitter, this.player.levelChangeEmitter]),
			getContent: () => {
				const talentPoints = getTalentPoints(this.player.getTalentsString());

				if (talentPoints == 0) {
					// Just return here, so we don't show a warning during page load.
					return '';
				} else if (talentPoints < this.player.getLevel() - 9) {
					return Tooltips.UNSPECT_TALENT_POINTS_WARNING;
				} else if (talentPoints > this.player.getLevel() - 9) {
					return Tooltips.TOO_MANY_TALENT_POINTS_WARNING;
				} else {
					return '';
				}
			},
		});
		this.addWarning({
			updateOn: TypedEvent.onAny([this.player.gearChangeEmitter, this.player.talentsChangeEmitter]),
			getContent: () => {
				if (
					!this.player.canDualWield2H() &&
					((this.player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType == HandType.HandTypeTwoHand &&
						this.player.getEquippedItem(ItemSlot.ItemSlotOffHand) != null) ||
						this.player.getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.handType == HandType.HandTypeTwoHand)
				) {
					return Tooltips.TITANS_GRIP_WARNING;
				} else {
					return '';
				}
			},
		});
		// TODO: This is showing the warning even when no items above player level are equipped
		// this.addWarning({
		// 	updateOn: TypedEvent.onAny([this.player.gearChangeEmitter]),
		// 	getContent: () => {
		// 		const playerLevel = player.getLevel();
		// 		if (
		// 			this.player
		// 				.getGear()
		// 				.asArray()
		// 				.filter(item => item != null && item.item.requiresLevel < playerLevel)
		// 		) {
		// 			return Tooltips.GEAR_MIN_LEVEL_WARNING(playerLevel);
		// 		} else {
		// 			return '';
		// 		}
		// 	},
		// });
		(config.warnings || []).forEach(warning => this.addWarning(warning(this)));

		if (!this.isWithinRaidSim) {
			// This needs to go before all the UI components so that gear loading is the
			// first callback invoked from waitForInit().
			this.sim.waitForInit().then(() => {
				this.loadSettings();

				if (isHealingSpec(this.player.spec)) {
					alert(Tooltips.HEALING_SIM_DISCLAIMER);
				}
			});
		}

		this.addSidebarComponents();
		this.addGearTab();
		this.bt = this.addBulkTab();
		this.addSettingsTab();
		this.addTalentsTab();
		this.addRotationTab();

		if (!this.isWithinRaidSim) {
			this.addDetailedResultsTab();
		}

		this.addTopbarComponents();
	}

	private loadSettings() {
		const initEventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			this.applyDefaults(initEventID);
			const savedSettings = window.localStorage.getItem(this.getSettingsStorageKey());
			if (savedSettings != null) {
				try {
					const settings = IndividualSimSettings.fromJsonString(savedSettings);
					this.fromProto(initEventID, settings);
				} catch (e) {
					console.warn('Failed to parse saved settings: ' + e);
				}
			}
			// Loading from link needs to happen after loading saved settings, so that partial link imports
			// (e.g. rotation only) include the previous settings for other categories.
			try {
				const urlParseResults = IndividualLinkImporter.tryParseUrlLocation(window.location);
				if (urlParseResults) {
					this.fromProto(initEventID, urlParseResults.settings, urlParseResults.categories);
				}
			} catch (e) {
				console.warn('Failed to parse link settings: ' + e);
			}
			window.location.hash = '';

			this.player.setName(initEventID, 'Player');

			// This needs to go last so it doesn't re-store things as they are initialized.
			this.changeEmitter.on(_ => {
				const jsonStr = IndividualSimSettings.toJsonString(this.toProto());

				window.localStorage.setItem(this.getSettingsStorageKey(), jsonStr);
			});

			this.statWeightActionSettings.load(initEventID);
		});
	}

	private addSidebarComponents() {
		this.raidSimResultsManager = addRaidSimAction(this);
		addStatWeightsAction(
			this,
			this.individualConfig.epStats.concat(GLOBAL_EP_STATS),
			GLOBAL_EP_PSEUDOSTATS.concat(this.individualConfig.epPseudoStats ?? []),
			this.individualConfig.epReferenceStat,
			this.statWeightActionSettings,
		);

		const displayStats: UnitStat[] = [];

		this.individualConfig.displayStats.forEach(s => displayStats.push(UnitStat.fromStat(s)));
		GLOBAL_DISPLAY_STATS.forEach(s => displayStats.push(UnitStat.fromStat(s)));

		this.individualConfig.displayPseudoStats.forEach(ps => displayStats.push(UnitStat.fromPseudoStat(ps)));
		GLOBAL_DISPLAY_PSEUDO_STATS.forEach(ps => displayStats.push(UnitStat.fromPseudoStat(ps)));

		new CharacterStats(
			this.rootElem.getElementsByClassName('sim-sidebar-stats')[0] as HTMLElement,
			this.player,
			displayStats,
			this.individualConfig.modifyDisplayStats,
		);
	}

	private addGearTab() {
		const gearTab = new GearTab(this.simTabContentsContainer, this);
		gearTab.rootElem.classList.add('active', 'show');
	}

	private addBulkTab(): BulkTab {
		const bulkTab = new BulkTab(this.simTabContentsContainer, this);
		bulkTab.navLink.hidden = !this.sim.getShowExperimental();
		this.sim.showExperimentalChangeEmitter.on(() => {
			bulkTab.navLink.hidden = !this.sim.getShowExperimental();
		});
		return bulkTab;
	}

	private addSettingsTab() {
		new SettingsTab(this.simTabContentsContainer, this);
	}

	private addTalentsTab() {
		new TalentsTab(this.simTabContentsContainer, this);
	}

	private addRotationTab() {
		new RotationTab(this.simTabContentsContainer, this);
	}

	private addDetailedResultsTab() {
		this.addTab(
			'Results',
			'detailed-results-tab',
			`
			<div class="detailed-results">
			</div>
		`,
		);

		new EmbeddedDetailedResults(this.rootElem.getElementsByClassName('detailed-results')[0] as HTMLElement, this, this.raidSimResultsManager!);
	}

	private addTopbarComponents() {
		// TODO: Classic
		this.simHeader.addImportLink('JSON', new IndividualJsonImporter(this.rootElem, this), true);
		this.simHeader.addImportLink('60U SoD', new Individual60UImporter(this.rootElem, this), true);
		//this.simHeader.addImportLink('WoWHead', new IndividualWowheadGearPlannerImporter(this.rootElem, this), false);
		this.simHeader.addImportLink('Addon', new IndividualAddonImporter(this.rootElem, this), true);

		this.simHeader.addExportLink('Link', new IndividualLinkExporter(this.rootElem, this), false);
		this.simHeader.addExportLink('JSON', new IndividualJsonExporter(this.rootElem, this), true);
		//this.simHeader.addExportLink('WoWHead', new IndividualWowheadGearPlannerExporter(this.rootElem, this), false);
		this.simHeader.addExportLink('60U SoD EP', new Individual60UEPExporter(this.rootElem, this), false);
		this.simHeader.addExportLink('Pawn EP', new IndividualPawnEPExporter(this.rootElem, this), false);
		//this.simHeader.addExportLink("CLI", new IndividualCLIExporter(this.rootElem, this), true);
	}

	applyEmptyAplRotation(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			this.player.setAplRotation(
				eventID,
				APLRotation.create({
					type: APLRotationType.TypeAPL,
				}),
			);
		});
	}

	applyDefaults(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			const tankSpec = isTankSpec(this.player.spec);
			const healingSpec = isHealingSpec(this.player.spec);

			//Special case for Totem of Wrath keeps buff and debuff sync'd
			this.player.applySharedDefaults(eventID);
			this.player.setRace(eventID, this.individualConfig.defaults.race ?? specToEligibleRaces[this.player.spec][0]);
			this.player.setLevel(eventID, LEVEL_THRESHOLDS[simLaunchStatuses[this.player.spec].phase]);
			this.player.setGear(eventID, this.sim.db.lookupEquipmentSpec(this.individualConfig.defaults.gear));
			if (this.individualConfig.defaults.itemSwap) {
				this.player.itemSwapSettings.setItemSwapSettings(
					eventID,
					true,
					this.sim.db.lookupItemSwap(this.individualConfig.defaults.itemSwap || ItemSwap.create()),
				);
			}
			this.player.setConsumes(eventID, this.individualConfig.defaults.consumes);
			this.player.setTalentsString(eventID, this.individualConfig.defaults.talents.talentsString);
			this.player.setSpecOptions(eventID, this.individualConfig.defaults.specOptions);
			this.player.setBuffs(eventID, this.individualConfig.defaults.individualBuffs);
			this.player.getParty()!.setBuffs(eventID, this.individualConfig.defaults.partyBuffs);
			this.player.getRaid()!.setBuffs(eventID, this.individualConfig.defaults.raidBuffs);
			this.player.setEpWeights(eventID, this.individualConfig.defaults.epWeights);
			const defaultRatios = this.player.getDefaultEpRatios(tankSpec, healingSpec);
			this.player.setEpRatios(eventID, defaultRatios);
			this.player.setProfession1(eventID, this.individualConfig.defaults.other?.profession1 || Profession.Engineering);
			this.player.setProfession2(eventID, this.individualConfig.defaults.other?.profession2 || Profession.ProfessionUnknown);
			this.player.setDistanceFromTarget(eventID, this.individualConfig.defaults.other?.distanceFromTarget || 0);
			this.player.setChannelClipDelay(eventID, this.individualConfig.defaults.other?.channelClipDelay || 0);

			if (this.isWithinRaidSim) {
				this.sim.raid.setTargetDummies(eventID, 0);
			} else {
				this.sim.raid.setTargetDummies(eventID, healingSpec ? 9 : 0);
				this.sim.encounter.applyDefaults(eventID);
				this.sim.raid.setDebuffs(eventID, this.individualConfig.defaults.debuffs);
				this.sim.applyDefaults(eventID, tankSpec, healingSpec);

				if (tankSpec) {
					this.sim.raid.setTanks(eventID, [this.player.makeUnitReference()]);
				} else {
					this.sim.raid.setTanks(eventID, []);
				}
			}

			this.statWeightActionSettings.applyDefaults(eventID);
		});
	}

	getSavedGearStorageKey(): string {
		return this.getStorageKey(SAVED_GEAR_STORAGE_KEY);
	}

	getSavedRotationStorageKey(): string {
		return this.getStorageKey(SAVED_ROTATION_STORAGE_KEY);
	}

	getSavedSettingsStorageKey(): string {
		return this.getStorageKey(SAVED_SETTINGS_STORAGE_KEY);
	}

	getSavedTalentsStorageKey(): string {
		return this.getStorageKey(SAVED_TALENTS_STORAGE_KEY);
	}

	// Returns the actual key to use for local storage, based on the given key part and the site context.
	getStorageKey(keyPart: string): string {
		// Local storage is shared by all sites under the same domain, so we need to use
		// different keys for each spec site.
		return specToLocalStorageKey[this.player.spec] + keyPart;
	}

	toProto(exportCategories?: Array<SimSettingCategories>): IndividualSimSettings {
		const exportCategory = (cat: SimSettingCategories) => !exportCategories || exportCategories.length == 0 || exportCategories.includes(cat);

		const proto = IndividualSimSettings.create({
			player: this.player.toProto(true, false, exportCategories),
		});

		if (exportCategory(SimSettingCategories.Miscellaneous)) {
			IndividualSimSettings.mergePartial(proto, {
				tanks: this.sim.raid.getTanks(),
			});
		}
		if (exportCategory(SimSettingCategories.Encounter)) {
			IndividualSimSettings.mergePartial(proto, {
				encounter: this.sim.encounter.toProto(),
			});
		}
		if (exportCategory(SimSettingCategories.External)) {
			IndividualSimSettings.mergePartial(proto, {
				partyBuffs: this.player.getParty()?.getBuffs() || PartyBuffs.create(),
				raidBuffs: this.sim.raid.getBuffs(),
				debuffs: this.sim.raid.getDebuffs(),
				targetDummies: this.sim.raid.getTargetDummies(),
			});
		}
		if (exportCategory(SimSettingCategories.UISettings)) {
			IndividualSimSettings.mergePartial(proto, {
				settings: this.sim.toProto(),
				epWeightsStats: this.player.getEpWeights().toProto(),
				epRatios: this.player.getEpRatios(),
				dpsRefStat: this.dpsRefStat,
				healRefStat: this.healRefStat,
				tankRefStat: this.tankRefStat,
			});
		}

		return proto;
	}

	toLink(): string {
		return IndividualLinkExporter.createLink(this);
	}

	fromProto(eventID: EventID, settings: IndividualSimSettings, includeCategories?: Array<SimSettingCategories>) {
		const loadCategory = (cat: SimSettingCategories) => !includeCategories || includeCategories.length == 0 || includeCategories.includes(cat);

		const tankSpec = isTankSpec(this.player.spec);
		const healingSpec = isHealingSpec(this.player.spec);

		TypedEvent.freezeAllAndDo(() => {
			if (!settings.player) {
				return;
			}

			this.player.fromProto(eventID, settings.player, includeCategories);

			if (loadCategory(SimSettingCategories.Miscellaneous)) {
				this.sim.raid.setTanks(eventID, settings.tanks || []);
			}
			if (loadCategory(SimSettingCategories.External)) {
				this.sim.raid.setBuffs(eventID, settings.raidBuffs || RaidBuffs.create());
				this.sim.raid.setDebuffs(eventID, settings.debuffs || Debuffs.create());
				const party = this.player.getParty();
				if (party) {
					party.setBuffs(eventID, settings.partyBuffs || PartyBuffs.create());
				}
				this.sim.raid.setTargetDummies(eventID, settings.targetDummies);
			}
			if (loadCategory(SimSettingCategories.Encounter)) {
				this.sim.encounter.fromProto(eventID, settings.encounter || EncounterProto.create());
			}
			if (loadCategory(SimSettingCategories.UISettings)) {
				if (settings.epWeightsStats) {
					this.player.setEpWeights(eventID, Stats.fromProto(settings.epWeightsStats));
				} else {
					this.player.setEpWeights(eventID, this.individualConfig.defaults.epWeights);
				}

				const defaultRatios = this.player.getDefaultEpRatios(tankSpec, healingSpec);
				if (settings.epRatios) {
					const missingRatios = new Array<number>(defaultRatios.length - settings.epRatios.length).fill(0);
					this.player.setEpRatios(eventID, settings.epRatios.concat(missingRatios));
				} else {
					this.player.setEpRatios(eventID, defaultRatios);
				}

				if (settings.dpsRefStat) {
					this.dpsRefStat = settings.dpsRefStat;
				}
				if (settings.healRefStat) {
					this.healRefStat = settings.healRefStat;
				}
				if (settings.tankRefStat) {
					this.tankRefStat = settings.tankRefStat;
				}

				if (settings.settings) {
					this.sim.fromProto(eventID, settings.settings);
				} else {
					this.sim.applyDefaults(eventID, tankSpec, healingSpec);
				}
			}
		});
	}
}
