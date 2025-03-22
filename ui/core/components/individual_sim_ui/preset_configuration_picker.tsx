import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { SpecOptions } from '../../../core/proto_utils/utils';
import { IndividualSimUI } from '../../individual_sim_ui';
import { PresetBuild } from '../../preset_utils';
import { APLRotation, APLRotation_Type } from '../../proto/apl';
import { Consumes, Debuffs, Encounter, EquipmentSpec, HealingModel, IndividualBuffs, ItemSwap, RaidBuffs, Spec } from '../../proto/common';
import { SavedTalents } from '../../proto/ui';
import { Stats } from '../../proto_utils/stats';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
import { ContentBlock } from '../content_block';

export enum PresetConfigurationCategory {
	EPWeights = 'epWeights',
	Gear = 'gear',
	Talents = 'talents',
	Rotation = 'rotation',
	Encounter = 'encounter',
	Settings = 'settings',
}

export class PresetConfigurationPicker extends Component {
	readonly simUI: IndividualSimUI<Spec>;
	readonly builds: Array<PresetBuild>;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>, types?: PresetConfigurationCategory[]) {
		super(parentElem, 'preset-configuration-picker-root');
		this.rootElem.classList.add('saved-data-manager-root');

		this.simUI = simUI;
		this.builds = (this.simUI.individualConfig.presets.builds ?? []).filter(build =>
			Object.keys(build).some(category => types?.includes(category as PresetConfigurationCategory) && !!build[category as PresetConfigurationCategory]),
		);

		if (!this.builds.length) {
			this.rootElem.classList.add('hide');
			return;
		}

		const contentBlock = new ContentBlock(this.rootElem, 'saved-data', {
			header: {
				title: 'Preset Configurations',
				tooltip: 'Preset configurations can apply an optimal combination of gear, talents, rotation and encounter settings.',
			},
		});

		const buildsContainerRef = ref<HTMLDivElement>();

		const container = (
			<div className="saved-data-container">
				<div className="saved-data-presets" ref={buildsContainerRef}></div>
			</div>
		);

		this.simUI.sim.waitForInit().then(() => {
			this.builds.forEach(build => {
				const dataElemRef = ref<HTMLButtonElement>();
				buildsContainerRef.value!.appendChild(
					<button className="saved-data-set-chip badge rounded-pill" ref={dataElemRef}>
						<span className="saved-data-set-name" attributes={{ role: 'button' }} onclick={() => this.applyBuild(build)}>
							{build.name}
						</span>
					</button>,
				);

				let categories = Object.keys(build).filter(c => !['name', 'encounter', 'settings'].includes(c) && build[c as PresetConfigurationCategory]);
				if (build.encounter?.encounter) {
					categories.push('encounter');
				}

				if (build.epWeights) {
					categories.push('stat weights');
				}

				if (build.settings) {
					Object.keys(build.settings).forEach(c => {
						if (['name', 'buffs', 'raidBuffs'].includes(c)) return;

						if (c === 'options') {
							categories.push('Class/Spec Options');
						} else if (c === 'consumes') {
							categories.push('Consumables');
						} else {
							categories.push('Other Settings');
						}
					});
				}

				if (build.settings?.buffs || build.settings?.raidBuffs) {
					categories.push('buffs');
				}

				categories = [...new Set(categories)].sort();

				tippy(dataElemRef.value!, {
					content: (
						<>
							<p className="mb-1">This preset affects the following settings:</p>
							<ul className="mb-0 text-capitalize">
								{categories.map(category => (
									<li>{category}</li>
								))}
							</ul>
						</>
					),
				});

				const checkActive = () => dataElemRef.value!.classList[this.isBuildActive(build) ? 'add' : 'remove']('active');

				checkActive();
				TypedEvent.onAny([
					this.simUI.player.changeEmitter,
					this.simUI.sim.settingsChangeEmitter,
					this.simUI.sim.raid.changeEmitter,
					this.simUI.sim.encounter.changeEmitter,
				]).on(checkActive);
			});
			contentBlock.bodyElement.replaceChildren(container);
		});
	}

	private applyBuild({ gear, rotation, rotationType, talents, epWeights, encounter, settings }: PresetBuild) {
		const eventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			if (gear) this.simUI.player.setGear(eventID, this.simUI.sim.db.lookupEquipmentSpec(gear.gear));
			if (talents) this.simUI.player.setTalentsString(eventID, talents.data.talentsString);
			if (rotationType) {
				this.simUI.player.aplRotation.type = rotationType;
				this.simUI.player.rotationChangeEmitter.emit(eventID);
			} else if (rotation?.rotation.rotation) {
				this.simUI.player.setAplRotation(eventID, rotation.rotation.rotation);
			}
			if (epWeights) this.simUI.player.setEpWeights(eventID, epWeights.epWeights);
			if (encounter) {
				if (encounter.encounter) this.simUI.sim.encounter.fromProto(eventID, encounter.encounter);
				if (encounter.healingModel) this.simUI.player.setHealingModel(eventID, encounter.healingModel);
				if (encounter.tanks) this.simUI.sim.raid.setTanks(eventID, encounter.tanks);
			}
			if (settings) {
				if (settings.level) this.simUI.player.setLevel(eventID, settings.level);
				if (settings.race) this.simUI.player.setRace(eventID, settings.race);
				if (settings.playerOptions?.profession1) this.simUI.player.setProfession1(eventID, settings.playerOptions.profession1);
				if (settings.playerOptions?.profession2) this.simUI.player.setProfession2(eventID, settings.playerOptions.profession2);
				if (settings.playerOptions?.distanceFromTarget) this.simUI.player.setDistanceFromTarget(eventID, settings.playerOptions.distanceFromTarget);
				if (settings.playerOptions?.enableItemSwap !== undefined && settings.playerOptions?.itemSwap) {
					this.simUI.player.itemSwapSettings.setItemSwapSettings(
						eventID,
						settings.playerOptions.enableItemSwap,
						this.simUI.sim.db.lookupItemSwap(settings.playerOptions.itemSwap),
						Stats.fromProto(settings.playerOptions.itemSwap.prepullBonusStats),
					);
				}
				if (settings.specOptions) {
					this.simUI.player.setSpecOptions(eventID, {
						...this.simUI.player.getSpecOptions(),
						...settings.specOptions,
					});
				}
				if (settings.buffs) this.simUI.player.setBuffs(eventID, settings.buffs);
				if (settings.debuffs) this.simUI.sim.raid.setDebuffs(eventID, settings.debuffs);
				if (settings.raidBuffs) this.simUI.sim.raid.setBuffs(eventID, settings.raidBuffs);
				if (settings.consumes) this.simUI.player.setConsumes(eventID, settings.consumes);
			}
		});
	}

	private isBuildActive({ gear, rotation, rotationType, talents, epWeights, encounter, settings }: PresetBuild): boolean {
		const hasGear = gear ? EquipmentSpec.equals(gear.gear, this.simUI.player.getGear().asSpec()) : true;
		const hasTalents = talents
			? SavedTalents.equals(
					talents.data,
					SavedTalents.create({
						talentsString: this.simUI.player.getTalentsString(),
					}),
			  )
			: true;
		let hasRotation = true;
		if (rotationType) {
			hasRotation = rotationType === this.simUI.player.getRotationType();
		} else if (rotation) {
			const activeRotation = this.simUI.player.getResolvedAplRotation();
			// Ensure that the auto rotation can be matched with a preset
			if (activeRotation.type === APLRotation_Type.TypeAuto) activeRotation.type = APLRotation_Type.TypeAPL;
			if (rotation.rotation?.rotation?.type === APLRotation_Type.TypeSimple && rotation.rotation.rotation?.simple?.specRotationJson) {
				hasRotation = this.simUI.player.specTypeFunctions.rotationEquals(
					this.simUI.player.specTypeFunctions.rotationFromJson(JSON.parse(rotation.rotation.rotation.simple.specRotationJson)),
					this.simUI.player.getSimpleRotation(),
				);
			} else {
				hasRotation = APLRotation.equals(rotation.rotation.rotation, activeRotation);
			}
		}
		const hasEpWeights = epWeights ? this.simUI.player.getEpWeights().equals(epWeights.epWeights) : true;
		const hasEncounter = encounter?.encounter ? Encounter.equals(encounter.encounter, this.simUI.sim.encounter.toProto()) : true;
		const hasHealingModel = encounter?.healingModel ? HealingModel.equals(encounter.healingModel, this.simUI.player.getHealingModel()) : true;
		const hasLevel = settings?.level ? this.simUI.player.getLevel() === settings.level : true;
		const hasRace = settings?.race ? this.simUI.player.getRace() === settings.race : true;
		const hasProfession1 = settings?.playerOptions?.profession1 === undefined || this.simUI.player.getProfession1() === settings.playerOptions.profession1;
		const hasProfession2 = settings?.playerOptions?.profession2 === undefined || this.simUI.player.getProfession2() === settings.playerOptions.profession2;
		const hasDistanceFromTarget =
			settings?.playerOptions?.distanceFromTarget === undefined ||
			this.simUI.player.getDistanceFromTarget() === settings.playerOptions.distanceFromTarget;
		const hasEnableItemSwap =
			settings?.playerOptions?.enableItemSwap === undefined ||
			this.simUI.player.itemSwapSettings.getEnableItemSwap() === settings.playerOptions.enableItemSwap;
		const hasItemSwap =
			settings?.playerOptions?.itemSwap === undefined ||
			ItemSwap.equals(this.simUI.player.itemSwapSettings?.toProto(), settings?.playerOptions?.itemSwap);
		const hasSpecOptions = settings?.specOptions ? this.containsAllFields(this.simUI.player.getSpecOptions(), settings.specOptions) : true;
		const hasConsumes = settings?.consumes ? Consumes.equals(this.simUI.player.getConsumes(), settings.consumes) : true;
		const hasRaidBuffs = settings?.raidBuffs ? RaidBuffs.equals(this.simUI.sim.raid.getBuffs(), settings.raidBuffs) : true;
		const hasBuffs = settings?.buffs ? IndividualBuffs.equals(this.simUI.player.getBuffs(), settings.buffs) : true;
		const hasDebuffs = settings?.debuffs ? Debuffs.equals(this.simUI.sim.raid.getDebuffs(), settings.debuffs) : true;

		return (
			hasGear &&
			hasTalents &&
			hasRotation &&
			hasEpWeights &&
			hasEncounter &&
			hasHealingModel &&
			hasLevel &&
			hasRace &&
			hasProfession1 &&
			hasProfession2 &&
			hasDistanceFromTarget &&
			hasEnableItemSwap &&
			hasItemSwap &&
			hasSpecOptions &&
			hasConsumes &&
			hasRaidBuffs &&
			hasBuffs &&
			hasDebuffs
		);
	}

	private containsAllFields<T extends Spec>(full: SpecOptions<T>, partial: Partial<SpecOptions<T>>): boolean {
		return Object.keys(partial).every(key => key in full && full[key as keyof SpecOptions<T>] === partial[key as keyof SpecOptions<T>]);
	}
}
