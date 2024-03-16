// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment, ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../individual_sim_ui';
import { PresetBuild } from '../../preset_utils';
import { APLRotation } from '../../proto/apl';
import { EquipmentSpec, Spec } from '../../proto/common';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';

export class PresetBuildsPicker extends Component {
	readonly simUI: IndividualSimUI<Spec>;
	readonly builds: Array<PresetBuild>;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, 'preset-builds-picker-root');

		this.simUI = simUI;
		this.builds = this.simUI.individualConfig.presets.builds ?? [];

		if (!this.builds.length) {
			this.rootElem.classList.add('hide');
			return;
		}

		const buildsContainerRef = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<>
				<label className="form-label">Preset Builds</label>
				<div className="presets-container" ref={buildsContainerRef}></div>
			</>,
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

				const checkActive = () => {
					if (this.isBuildActive(build)) {
						dataElemRef.value!.classList.add('active');
					} else {
						dataElemRef.value!.classList.remove('active');
					}
				};

				checkActive();
				TypedEvent.onAny([this.simUI.player.gearChangeEmitter, this.simUI.player.talentsChangeEmitter, this.simUI.player.rotationChangeEmitter]).on(
					checkActive,
				);
			});
		});
	}

	private applyBuild(build: PresetBuild) {
		const eventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			this.simUI.player.setGear(eventID, this.simUI.sim.db.lookupEquipmentSpec(build.gear.gear));
			this.simUI.player.setTalentsString(eventID, build.talents.data.talentsString);
			this.simUI.player.setAplRotation(eventID, build.rotation.rotation.rotation!);
		});
	}

	private isBuildActive(build: PresetBuild): boolean {
		build.rotation.rotation.rotation;
		return (
			EquipmentSpec.equals(build.gear.gear, this.simUI.player.getGear().asSpec()) &&
			build.talents.data.talentsString == this.simUI.player.getTalentsString() &&
			APLRotation.equals(build.rotation.rotation.rotation, this.simUI.player.getResolvedAplRotation())
		);
	}
}
