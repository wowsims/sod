import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { ResourceType } from '../../proto/api.js';
import { OtherAction } from '../../proto/common.js';
import { ActionId, resourceTypeToIcon } from '../../proto_utils/action_id.js';
import { AuraUptimeLog, CastLog, DamageDealtLog, DpsLog, ResourceChangedLogGroup, SimLog, ThreatLogGroup } from '../../proto_utils/logs_parser.js';
import { resourceNames } from '../../proto_utils/names.js';
import { UnitMetrics } from '../../proto_utils/sim_result.js';
import { orderedResourceTypes } from '../../proto_utils/utils.js';
import { TypedEvent } from '../../typed_event.js';
import { bucket, distinct, maxIndex, stringComparator } from '../../utils.js';
import { actionColors } from './color_settings.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

declare let ApexCharts: any;

type TooltipHandler = (dataPointIndex: number) => Element;

const dpsColor = '#ed5653';
const manaColor = '#2E93fA';
const threatColor = '#b56d07';

export class Timeline extends ResultComponent {
	private readonly dpsResourcesPlotElem: HTMLElement;
	private dpsResourcesPlot: any;

	private readonly rotationPlotElem: HTMLElement;
	private readonly rotationLabels: HTMLElement;
	private readonly rotationTimeline: HTMLElement;
	private readonly rotationHiddenIdsContainer: HTMLElement;
	private readonly chartPicker: HTMLSelectElement;

	private resultData: SimResultData | null;
	private rendered: boolean;

	private hiddenIds: Array<ActionId>;
	private hiddenIdsChangeEmitter;

	private resetCallbacks: (() => void)[] = [];

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'timeline-root';
		super(config);
		this.resultData = null;
		this.rendered = false;
		this.hiddenIds = [];
		this.hiddenIdsChangeEmitter = new TypedEvent<void>();

		this.rootElem.appendChild(
			<div className="timeline-disclaimer">
				<div className="d-flex flex-column">
					<p>
						<i className="warning fa fa-exclamation-triangle fa-xl me-2"></i>
						Timeline data visualizes only 1 sim iteration.
					</p>
					<p>
						Note: You can move the timeline by holding <kbd>Shift</kbd> while scrolling, or by clicking and dragging.
					</p>
				</div>
				<select className="timeline-chart-picker form-select">
					<option className="rotation-option" value="rotation">
						Rotation
					</option>
					<option className="dps-option" value="dps">
						DPS
					</option>
					<option className="threat-option" value="threat">
						Threat
					</option>
				</select>
			</div>,
		);
		this.rootElem.appendChild(
			<div className="timeline-plots-container">
				<div className="timeline-plot dps-resources-plot hide"></div>
				<div className="timeline-plot rotation-plot">
					<div className="rotation-container">
						<div className="rotation-labels"></div>
						<div className="rotation-timeline" draggable={true}></div>
					</div>
					<div className="rotation-hidden-ids"></div>
				</div>
			</div>,
		);

		this.chartPicker = this.rootElem.querySelector('.timeline-chart-picker')!;
		this.chartPicker.addEventListener('change', () => {
			if (this.chartPicker.value == 'rotation') {
				this.dpsResourcesPlotElem.classList.add('hide');
				this.rotationPlotElem.classList.remove('hide');
			} else {
				this.dpsResourcesPlotElem.classList.remove('hide');
				this.rotationPlotElem.classList.add('hide');
			}
			this.updatePlot();
		});

		this.dpsResourcesPlotElem = this.rootElem.querySelector('.dps-resources-plot')!;
		this.dpsResourcesPlot = new ApexCharts(this.dpsResourcesPlotElem, {
			chart: {
				type: 'line',
				foreColor: 'white',
				id: 'dpsResources',
				animations: {
					enabled: false,
				},
			},
			series: [], // Set dynamically
			xaxis: {
				title: {
					text: 'Time (s)',
				},
			},
			noData: {
				text: 'Waiting for data...',
			},
			stroke: {
				width: 2,
				curve: 'straight',
			},
		});

		this.rotationPlotElem = this.rootElem.querySelector('.rotation-plot')!;
		this.rotationLabels = this.rootElem.querySelector('.rotation-labels')!;
		this.rotationTimeline = this.rootElem.querySelector('.rotation-timeline')!;
		this.rotationHiddenIdsContainer = this.rootElem.querySelector('.rotation-hidden-ids')!;

		let isMouseDown = false;
		let startX = 0;
		let scrollLeft = 0;
		this.rotationTimeline.addEventListener('dragstart', event => {
			event.preventDefault();
		});
		this.rotationTimeline.addEventListener('mousedown', event => {
			isMouseDown = true;
			startX = event.pageX - this.rotationTimeline.offsetLeft;
			scrollLeft = this.rotationTimeline.scrollLeft;
		});
		this.rotationTimeline.addEventListener('mouseleave', () => {
			isMouseDown = false;
			this.rotationTimeline.classList.remove('active');
		});
		this.rotationTimeline.addEventListener('mouseup', () => {
			isMouseDown = false;
			this.rotationTimeline.classList.remove('active');
		});
		this.rotationTimeline.addEventListener('mousemove', event => {
			if (!isMouseDown) return;
			event.preventDefault();
			const x = event.pageX - this.rotationTimeline.offsetLeft;
			const walk = (x - startX) * 3; //scroll-fast
			this.rotationTimeline.scrollLeft = scrollLeft - walk;
		});
	}

	onSimResult(resultData: SimResultData) {
		this.resultData = resultData;

		if (this.rendered) {
			this.updatePlot();
		}
	}

	private updatePlot() {
		if (this.resultData == null) {
			return;
		}

		const duration = this.resultData!.result.result.firstIterationDuration || 1;
		const options: any = {
			series: [],
			colors: [],
			xaxis: {
				min: 0,
				max: duration,
				tickAmount: 10,
				decimalsInFloat: 1,
				labels: {
					show: true,
				},
				title: {
					text: 'Time (s)',
				},
			},
			yaxis: [],
			chart: {
				events: {
					beforeResetZoom: () => {
						return {
							xaxis: {
								min: 0,
								max: duration,
							},
						};
					},
				},
				toolbar: {
					show: false,
				},
			},
		};

		let tooltipHandlers: Array<TooltipHandler | null> = [];
		options.tooltip = {
			enabled: true,
			custom: (data: { series: any; seriesIndex: number; dataPointIndex: number; w: any }) => {
				if (tooltipHandlers[data.seriesIndex]) {
					return tooltipHandlers[data.seriesIndex]!(data.dataPointIndex);
				} else {
					throw new Error('No tooltip handler for series ' + data.seriesIndex);
				}
			},
		};

		const players = this.resultData!.result.getRaidIndexedPlayers(this.resultData!.filter);
		if (players.length == 1) {
			const player = players[0];

			const rotationOption = this.rootElem.querySelector('.rotation-option')!;
			rotationOption.classList.remove('hide');
			const threatOption = this.rootElem.querySelector('.threat-option')!;
			threatOption.classList.add('hide');

			try {
				this.updateRotationChart(player, duration);
			} catch (e) {
				console.log('Failed to update rotation chart: ', e);
			}

			const dpsData = this.addDpsSeries(player, options, '');
			this.addDpsYAxis(dpsData.maxDps, options);
			tooltipHandlers.push(dpsData.tooltipHandler);
			tooltipHandlers.push(this.addManaSeries(player, options));
			tooltipHandlers.push(this.addThreatSeries(player, options, ''));
			tooltipHandlers = tooltipHandlers.filter(handler => handler != null);

			this.addMajorCooldownAnnotations(player, options);
		} else {
			if (this.chartPicker.value == 'rotation') {
				this.chartPicker.value = 'dps';
				return;
			}
			const rotationOption = this.rootElem.querySelector('.rotation-option')!;
			rotationOption.classList.add('hide');
			const threatOption = this.rootElem.querySelector('.threat-option')!;
			threatOption.classList.remove('hide');

			this.clearRotationChart();

			if (this.chartPicker.value == 'dps') {
				let maxDps = 0;
				players.forEach(player => {
					const dpsData = this.addDpsSeries(player, options, `var(--bs-${player.classColor}`);
					maxDps = Math.max(maxDps, dpsData.maxDps);
					tooltipHandlers.push(dpsData.tooltipHandler);
				});
				this.addDpsYAxis(maxDps, options);
			} else {
				// threat
				let maxThreat = 0;
				players.forEach(player => {
					tooltipHandlers.push(this.addThreatSeries(player, options, player.classColor));
					maxThreat = Math.max(maxThreat, player.maxThreat);
				});
				this.addThreatYAxis(maxThreat, options);
			}
		}

		this.dpsResourcesPlot.updateOptions(options);
	}

	private addDpsYAxis(maxDps: number, options: any) {
		const dpsAxisMax = Math.ceil(maxDps / 100) * 100;
		options.yaxis.push({
			color: dpsColor,
			seriesName: 'DPS',
			min: 0,
			max: dpsAxisMax,
			tickAmount: 10,
			decimalsInFloat: 0,
			title: {
				text: 'DPS',
				style: {
					color: dpsColor,
				},
			},
			axisBorder: {
				show: true,
				color: dpsColor,
			},
			axisTicks: {
				color: dpsColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [dpsColor],
				},
			},
		});
	}

	private addThreatYAxis(maxThreat: number, options: any) {
		const axisMax = Math.ceil(maxThreat / 10000) * 10000;
		options.yaxis.push({
			color: threatColor,
			seriesName: 'Threat',
			min: 0,
			max: axisMax,
			tickAmount: 10,
			decimalsInFloat: 0,
			title: {
				text: 'Threat',
				style: {
					color: threatColor,
				},
			},
			axisBorder: {
				show: true,
				color: threatColor,
			},
			axisTicks: {
				color: threatColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [threatColor],
				},
			},
		});
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addDpsSeries(unit: UnitMetrics, options: any, colorOverride: string): { maxDps: number; tooltipHandler: TooltipHandler } {
		const dpsLogs = unit.dpsLogs.filter(log => log.timestamp >= 0);

		options.colors.push(colorOverride || dpsColor);
		options.series.push({
			name: 'DPS',
			type: 'line',
			data: dpsLogs.map(log => {
				return {
					x: log.timestamp,
					y: log.dps,
				};
			}),
		});

		return {
			maxDps: dpsLogs[maxIndex(dpsLogs.map(l => l.dps))!].dps,
			tooltipHandler: (dataPointIndex: number) => {
				const log = dpsLogs[dataPointIndex];
				return this.dpsTooltip(log, true, unit, colorOverride);
			},
		};
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addManaSeries(unit: UnitMetrics, options: any): TooltipHandler | null {
		const manaLogs = unit.groupedResourceLogs[ResourceType.ResourceTypeMana].filter(log => log.timestamp >= 0);
		if (manaLogs.length == 0) {
			return null;
		}
		const maxMana = manaLogs[0].valueBefore;

		options.colors.push(manaColor);
		options.series.push({
			name: 'Mana',
			type: 'line',
			data: manaLogs.map(log => {
				return {
					x: log.timestamp,
					y: log.valueAfter,
				};
			}),
		});
		options.yaxis.push({
			seriesName: 'Mana',
			opposite: true, // Appear on right side
			min: 0,
			max: maxMana,
			tickAmount: 10,
			title: {
				text: 'Mana',
				style: {
					color: manaColor,
				},
			},
			axisBorder: {
				show: true,
				color: manaColor,
			},
			axisTicks: {
				color: manaColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [manaColor],
				},
				formatter: (val: string) => {
					const v = parseFloat(val);
					return `${v.toFixed(0)} (${((v / maxMana) * 100).toFixed(0)}%)`;
				},
			},
		} as any);

		return (dataPointIndex: number) => {
			const log = manaLogs[dataPointIndex];
			return this.resourceTooltip(log, maxMana, true);
		};
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addThreatSeries(unit: UnitMetrics, options: any, colorOverride: string): TooltipHandler | null {
		options.colors.push(colorOverride || threatColor);
		options.series.push({
			name: 'Threat',
			type: 'line',
			data: unit.threatLogs
				.filter(log => log.timestamp >= 0)
				.map(log => {
					return {
						x: log.timestamp,
						y: log.threatAfter,
					};
				}),
		});

		return (dataPointIndex: number) => {
			const log = unit.threatLogs[dataPointIndex];
			return this.threatTooltip(log, true, unit, colorOverride);
		};
	}

	private addMajorCooldownAnnotations(unit: UnitMetrics, options: any) {
		const mcdLogs = unit.majorCooldownLogs;
		const mcdAuraLogs = unit.majorCooldownAuraUptimeLogs;

		// Figure out how much to vertically offset cooldown icons, for cooldowns
		// used very close to each other. This is so the icons don't overlap.
		const MAX_ALLOWED_DIST = 10;
		const cooldownIconOffsets = mcdLogs.map(
			(mcdLog, mcdIdx) => mcdLogs.filter((cdLog, cdIdx) => cdIdx < mcdIdx && cdLog.timestamp > mcdLog.timestamp - MAX_ALLOWED_DIST).length,
		);

		const distinctMcdAuras = distinct(mcdAuraLogs, (a, b) => a.actionId!.equalsIgnoringTag(b.actionId!));
		// Sort by name so auras keep their same colors even if timings change.
		distinctMcdAuras.sort((a, b) => stringComparator(a.actionId!.name, b.actionId!.name));
		const mcdAuraColors = mcdAuraLogs.map(
			mcdAuraLog => actionColors[distinctMcdAuras.findIndex(dAura => dAura.actionId!.equalsIgnoringTag(mcdAuraLog.actionId!))],
		);

		options.annotations = {
			position: 'back',
			xaxis: mcdAuraLogs.map((log, i) => {
				return {
					x: log.gainedAt,
					x2: log.fadedAt,
					fillColor: mcdAuraColors[i],
				};
			}),
			points: mcdLogs.map((log, i) => {
				return {
					x: log.timestamp,
					y: 0,
					image: {
						path: log.actionId!.iconUrl,
						width: 20,
						height: 20,
						offsetY: cooldownIconOffsets[i] * -25,
					},
				};
			}),
		};
	}

	private clearRotationChart() {
		this.rotationLabels.innerText = '';
		this.rotationLabels.appendChild(<div className="rotation-label-header"></div>);

		this.rotationTimeline.innerText = '';
		this.rotationTimeline.appendChild(
			<div className="rotation-timeline-header">
				<canvas className="rotation-timeline-canvas"></canvas>
			</div>,
		);
		this.rotationHiddenIdsContainer.innerText = '';
		this.hiddenIdsChangeEmitter = new TypedEvent<void>();
	}

	private updateRotationChart(player: UnitMetrics, duration: number) {
		const targets = this.resultData!.result.getTargets(this.resultData!.filter);
		if (targets.length == 0) {
			return;
		}
		const target = targets[0];

		this.clearRotationChart();

		try {
			this.drawRotationTimeRuler(this.rotationTimeline.querySelector('.rotation-timeline-canvas')!, duration);
		} catch (e) {
			console.log('Failed to draw rotation: ', e);
		}

		orderedResourceTypes.forEach(resourceType => this.addResourceRow(resourceType, player.groupedResourceLogs[resourceType], duration));

		const buffsById = Object.values(bucket(player.auraUptimeLogs, log => log.actionId!.toString()));
		buffsById.sort((a, b) => stringComparator(a[0].actionId!.name, b[0].actionId!.name));
		const debuffsById = Object.values(bucket(target.auraUptimeLogs, log => log.actionId!.toString()));
		debuffsById.sort((a, b) => stringComparator(a[0].actionId!.name, b[0].actionId!.name));
		const buffsAndDebuffsById = buffsById.concat(debuffsById);

		const playerCastsByAbility = this.getSortedCastsByAbility(player);
		playerCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));

		if (player.pets.length > 0) {
			const playerPets = new Map<string, UnitMetrics>();
			player.pets.forEach(petsLog => {
				const petCastsByAbility = this.getSortedCastsByAbility(petsLog);
				if (petCastsByAbility.length > 0) {
					// Because mulle pets can have the same name and we parse cast logs
					// by pet name each individual pet ends up with all the casts of pets
					// with the same name. Because of this we can just grab the first pet
					// of each name and visualize only that.
					if (!playerPets.has(petsLog.name)) {
						playerPets.set(petsLog.name, petsLog);
					}
				}
			});

			playerPets.forEach(pet => {
				this.addSeparatorRow(duration);
				this.addPetRow(pet.name, duration);
				orderedResourceTypes.forEach(resourceType => this.addResourceRow(resourceType, pet.groupedResourceLogs[resourceType], duration));
				const petCastsByAbility = this.getSortedCastsByAbility(pet);
				petCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));
			});
		}

		// Don't add a row for buffs that were already visualized in a cast row.
		const buffsToShow = buffsById.filter(auraUptimeLogs =>
			playerCastsByAbility.findIndex(casts => casts[0].actionId!.equalsIgnoringTag(auraUptimeLogs[0].actionId!)),
		);
		if (buffsToShow.length > 0) {
			this.addSeparatorRow(duration);
			buffsToShow.forEach(auraUptimeLogs => this.addAuraRow(auraUptimeLogs, duration));
		}

		const targetCastsByAbility = this.getSortedCastsByAbility(target);
		if (targetCastsByAbility.length > 0) {
			this.addSeparatorRow(duration);
			targetCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));
		}

		// Add a row for all debuffs, even those which have already been visualized in a cast row.
		const debuffsToShow = debuffsById;
		if (debuffsToShow.length > 0) {
			this.addSeparatorRow(duration);
			debuffsToShow.forEach(auraUptimeLogs => this.addAuraRow(auraUptimeLogs, duration));
		}
	}

	private getSortedCastsByAbility(player: UnitMetrics): Array<Array<CastLog>> {
		const meleeActionIds = player.getMeleeActions().map(action => action.actionId);
		const spellActionIds = player.getSpellActions().map(action => action.actionId);
		const getActionCategory = (actionId: ActionId): number => {
			const fixedCategory = idToCategoryMap[actionId.anyId()];
			if (fixedCategory != null) {
				return fixedCategory;
			} else if (meleeActionIds.find(meleeActionId => meleeActionId.equals(actionId))) {
				return MELEE_ACTION_CATEGORY;
			} else if (spellActionIds.find(spellActionId => spellActionId.equals(actionId))) {
				return SPELL_ACTION_CATEGORY;
			} else {
				return DEFAULT_ACTION_CATEGORY;
			}
		};

		const castsByAbility = Object.values(
			bucket(player.castLogs, log => {
				if (idsToGroupForRotation.includes(log.actionId!.spellId)) {
					return log.actionId!.toStringIgnoringTag();
				} else {
					return log.actionId!.toString();
				}
			}),
		);

		castsByAbility.sort((a, b) => {
			const categoryA = getActionCategory(a[0].actionId!);
			const categoryB = getActionCategory(b[0].actionId!);
			if (categoryA != categoryB) {
				return categoryA - categoryB;
			} else if (a[0].actionId!.anyId() == b[0].actionId!.anyId()) {
				return a[0].actionId!.tag - b[0].actionId!.tag;
			} else {
				return stringComparator(a[0].actionId!.name, b[0].actionId!.name);
			}
		});

		return castsByAbility;
	}

	private makeLabelElem(actionId: ActionId, isHiddenLabel: boolean, isAura?: boolean): JSX.Element {
		const labelText = idsToGroupForRotation.includes(actionId.spellId) ? actionId.baseName : actionId.name;
		const labelIcon = ref<HTMLAnchorElement>();
		const hideElem = ref<HTMLElement>();
		const labelElem = (
			<div className={`rotation-label rotation-row ${isHiddenLabel ? 'rotation-label-hidden' : ''}`}>
				<span ref={hideElem} className={`fas fa-eye${isHiddenLabel ? '' : '-slash'}`}></span>
				<a ref={labelIcon} className="rotation-label-icon"></a>
				<span className="rotation-label-text">{labelText}</span>
			</div>
		);
		const onClickHandler = () => {
			if (isHiddenLabel) {
				const index = this.hiddenIds.findIndex(hiddenId => hiddenId.equals(actionId));
				if (index != -1) {
					this.hiddenIds.splice(index, 1);
				}
			} else {
				this.hiddenIds.push(actionId);
			}
			this.hiddenIdsChangeEmitter.emit(TypedEvent.nextEventID());
		};
		hideElem.value!.addEventListener('click', onClickHandler);
		const tooltip = tippy(hideElem.value!, {
			theme: 'timeline-tooltip',
			placement: 'bottom',
			content: isHiddenLabel ? 'Show Row' : 'Hide Row',
		});

		const updateHidden = () => {
			if (isHiddenLabel == Boolean(this.hiddenIds.find(hiddenId => hiddenId.equals(actionId)))) {
				labelElem.classList.remove('hide');
			} else {
				labelElem.classList.add('hide');
			}
		};
		const event = this.hiddenIdsChangeEmitter.on(updateHidden);
		updateHidden();
		actionId.setBackgroundAndHref(labelIcon.value!);
		actionId.setWowheadDataset(labelIcon.value!, { useBuffAura: isAura });

		this.addOnResetCallback(() => {
			hideElem.value?.removeEventListener('click', onClickHandler);
			tooltip.destroy();
			event.dispose();
		});

		return labelElem;
	}

	private makeRowElem(actionId: ActionId, duration: number): JSX.Element {
		const rowElem = (
			<div
				className="rotation-timeline-row rotation-row"
				style={{
					width: this.timeToPx(duration),
				}}></div>
		);

		const updateHidden = () => {
			if (this.hiddenIds.find(hiddenId => hiddenId.equals(actionId))) {
				rowElem.classList.add('hide');
			} else {
				rowElem.classList.remove('hide');
			}
		};
		const event = this.hiddenIdsChangeEmitter.on(updateHidden);
		updateHidden();
		this.addOnResetCallback(() => event.dispose());
		return rowElem;
	}

	private addPetRow(petName: string, duration: number) {
		const actionId = ActionId.fromPetName(petName);
		const rowElem = this.makeRowElem(actionId, duration);

		const iconElem = document.createElement('div');
		this.rotationLabels.appendChild(iconElem);

		actionId.fill().then(filledActionId => {
			const labelText = idsToGroupForRotation.includes(filledActionId.spellId) ? filledActionId.baseName : filledActionId.name;
			const labelIcon = ref<HTMLAnchorElement>();
			const labelElem = (
				<div className="rotation-label rotation-row">
					<a ref={labelIcon} className="rotation-label-icon"></a>
					<span className="rotation-label-text">{labelText}</span>
				</div>
			);
			filledActionId.setBackgroundAndHref(labelIcon.value!);
			iconElem.appendChild(labelElem);
		});

		this.rotationTimeline.appendChild(rowElem);
	}

	private addSeparatorRow(duration: number) {
		const separatorElem = <div className="rotation-timeline-separator"></div>;
		this.rotationLabels.appendChild(separatorElem.cloneNode());
		separatorElem.style.width = this.timeToPx(duration);
		this.rotationTimeline.appendChild(separatorElem);
	}

	private addResourceRow(resourceType: ResourceType, resourceLogs: Array<ResourceChangedLogGroup>, duration: number) {
		if (resourceLogs.length == 0) {
			return;
		}
		const startValue = resourceLogs[0].valueBefore;

		const labelElem = (
			<div className="rotation-label rotation-row">
				<a
					className="rotation-label-icon"
					style={{
						backgroundImage: `url('${resourceTypeToIcon[resourceType]}')`,
					}}></a>
				<span className="rotation-label-text">{resourceNames.get(resourceType)}</span>
			</div>
		);

		this.rotationLabels.appendChild(labelElem);

		const rowElem = (
			<div
				className="rotation-timeline-row rotation-row"
				style={{
					width: this.timeToPx(duration),
				}}></div>
		);

		resourceLogs.forEach((resourceLogGroup, i) => {
			const cNames = resourceNames.get(resourceType)!.toLowerCase().replaceAll(' ', '-');
			const resourceElem = (
				<div
					className={`rotation-timeline-resource series-color ${cNames}`}
					style={{
						left: this.timeToPx(resourceLogGroup.timestamp),
						width: this.timeToPx((resourceLogs[i + 1]?.timestamp || duration) - resourceLogGroup.timestamp),
					}}></div>
			);

			if (percentageResources.includes(resourceType)) {
				resourceElem.textContent = ((resourceLogGroup.valueAfter / startValue) * 100).toFixed(0) + '%';
			} else {
				if (resourceType == ResourceType.ResourceTypeEnergy) {
					const bgElem = document.createElement('div');
					bgElem.classList.add('rotation-timeline-resource-fill');
					bgElem.classList.add(cNames);
					bgElem.style.height = ((resourceLogGroup.valueAfter / startValue) * 100).toFixed(0) + '%';
					resourceElem.appendChild(bgElem);
				} else {
					resourceElem.textContent = Math.floor(resourceLogGroup.valueAfter).toFixed(0);
				}
			}
			rowElem.appendChild(resourceElem);

			const tooltip = tippy(resourceElem, {
				placement: 'bottom',
				content: this.resourceTooltipElem(resourceLogGroup, startValue, false),
			});
			this.addOnResetCallback(() => tooltip.destroy());
		});
		this.rotationTimeline.appendChild(rowElem);
	}

	private addCastRow(castLogs: Array<CastLog>, aurasById: Array<Array<AuraUptimeLog>>, duration: number) {
		const actionId = castLogs[0].actionId!;

		this.rotationLabels.appendChild(this.makeLabelElem(actionId, false));
		this.rotationHiddenIdsContainer.appendChild(this.makeLabelElem(actionId, true));

		const rowElem = this.makeRowElem(actionId, duration);

		let stackedIconCount = 1;
		let stackedDamageCount = 1;

		castLogs.forEach((castLog, index) => {
			const castElem = (
				<div
					className="rotation-timeline-cast"
					style={{
						left: this.timeToPx(castLog.timestamp),
						minWidth: this.timeToPx(castLog.castTime + castLog.travelTime),
					}}></div>
			);

			rowElem.appendChild(castElem);

			if (castLog.travelTime != 0) {
				const travelTimeElem = (
					<div
						className="rotation-timeline-travel-time"
						style={{
							left: this.timeToPx(castLog.castTime),
							minWidth: this.timeToPx(castLog.travelTime),
						}}
					/>
				);
				castElem.appendChild(travelTimeElem);
			}

			if (castLog.damageDealtLogs.length > 0) {
				const ddl = castLog.damageDealtLogs[0];
				if (ddl.miss || ddl.dodge || ddl.parry) {
					castElem.classList.add('outcome-miss');
				} else if (ddl.glance || ddl.block || ddl.partialResist1_4 || ddl.partialResist2_4 || ddl.partialResist3_4) {
					castElem.classList.add('outcome-partial');
				} else if (ddl.crit) {
					castElem.classList.add('outcome-crit');
				} else {
					castElem.classList.add('outcome-hit');
				}
			}

			const iconElem = (<a className="rotation-timeline-cast-icon" />) as HTMLAnchorElement;
			actionId.setBackground(iconElem);
			castElem.appendChild(iconElem);

			const travelTimeStr = castLog.travelTime == 0 ? '' : ` + ${castLog.travelTime.toFixed(2)}s travel time`;
			const totalDamage = castLog.totalDamage();

			if (index > 0) {
				let timeDelta = null;

				for (let i = index - 1; i >= 0; i--) {
					if (castLog.timestamp != castLogs[i].timestamp) {
						timeDelta = castLog.timestamp - castLogs[i].timestamp;
						break;
					}
				}
				if (timeDelta != null && timeDelta < 0.21) {
					stackedDamageCount = stackedDamageCount + castLog.damageDealtLogs.length;
					stackedIconCount = stackedIconCount + 1;
				} else {
					stackedDamageCount = castLog.damageDealtLogs.length;
					stackedIconCount = 1;
				}
			}

			const startIndex = Math.max(0, 1 + index - stackedIconCount);
			const relevantCastLogs = castLogs.slice(startIndex, index + 1);

			// Reset and initialize aggregatedData for the current iteration
			const aggregatedData = relevantCastLogs.reduce<{
				damageDealtLogs: DamageDealtLog[];
				totalDamage: number;
			}>(
				(acc, log) => {
					// Correct type
					acc.damageDealtLogs.push(...log.damageDealtLogs);
					acc.totalDamage += log.totalDamage();
					return acc;
				},
				{
					damageDealtLogs: [], // Correctly typed as DamageDealtLog[]
					totalDamage: 0,
				},
			);

			const tt = (
				<div className="timeline-tooltip">
					<span>
						{castLog.actionId!.name} from {castLog.timestamp.toFixed(2)}s to {(castLog.timestamp + castLog.castTime).toFixed(2)}s (
						{castLog.castTime > 0 && `${castLog.castTime.toFixed(2)}s, `}
						{castLog.effectiveTime > 0 && `${castLog.effectiveTime.toFixed(2)}s GCD Time`}){travelTimeStr.length > 0 && travelTimeStr}
					</span>
					{aggregatedData.damageDealtLogs.length > 0 && (
						<ul className="rotation-timeline-cast-damage-list">
							{aggregatedData.damageDealtLogs.map((ddl, ddlIndex) => (
								<li>
									<span>
										{ddl.timestamp.toFixed(2)}s - {ddl.result()}
									</span>
									{ddl.source?.isTarget && <span className="threat-metrics"> ({ddl.threat.toFixed(1)} Threat)</span>}
								</li>
							))}
						</ul>
					)}
					{totalDamage > 0 && (
						<span>
							Total: {aggregatedData.totalDamage.toFixed(2)} ({(aggregatedData.totalDamage / (castLog.effectiveTime || 1)).toFixed(2)} DPET)
						</span>
					)}
				</div>
			);

			tippy(castElem, {
				placement: 'bottom',
				content: tt,
			});

			if (stackedIconCount > 1) {
				castElem.appendChild(<div className="stacked-icon-count">{String(stackedDamageCount)}</div>);
			}

			castLog.damageDealtLogs
				.filter(ddl => ddl.tick)
				.forEach(ddl => {
					const tickElem = (
						<div
							className="rotation-timeline-tick"
							style={{
								left: this.timeToPx(ddl.timestamp),
							}}
						/>
					);
					rowElem.appendChild(tickElem);

					const tt = (
						<div className="timeline-tooltip">
							<span>
								{ddl.timestamp.toFixed(2)}s - {ddl.actionId!.name} {ddl.result()}
							</span>
							{ddl.source?.isTarget && <span className="threat-metrics"> ({ddl.threat.toFixed(1)} Threat)</span>}
						</div>
					);

					tippy(tickElem, {
						placement: 'bottom',
						content: tt,
					});
				});
		});

		// If there are any auras that correspond to this cast, visualize them in the same row.
		aurasById
			.filter(auraUptimeLogs => auraUptimeLogs[0].actionId!.equalsIgnoringTag(actionId))
			.forEach(auraUptimeLogs => this.applyAuraUptimeLogsToRow(auraUptimeLogs, rowElem));

		this.rotationTimeline.appendChild(rowElem);
	}

	private addAuraRow(auraUptimeLogs: Array<AuraUptimeLog>, duration: number) {
		const actionId = auraUptimeLogs[0].actionId!;

		const rowElem = this.makeRowElem(actionId, duration);
		this.rotationLabels.appendChild(this.makeLabelElem(actionId, false, true));
		this.rotationHiddenIdsContainer.appendChild(this.makeLabelElem(actionId, true, true));
		this.rotationTimeline.appendChild(rowElem);

		this.applyAuraUptimeLogsToRow(auraUptimeLogs, rowElem);
	}

	private applyAuraUptimeLogsToRow(auraUptimeLogs: Array<AuraUptimeLog>, rowElem: JSX.Element) {
		auraUptimeLogs.forEach(aul => {
			const auraElem = (
				<div
					className="rotation-timeline-aura"
					style={{
						left: this.timeToPx(aul.gainedAt),
						minWidth: this.timeToPx(aul.fadedAt === aul.gainedAt ? 0.001 : aul.fadedAt - aul.gainedAt),
					}}
				/>
			);
			rowElem.appendChild(auraElem);

			const tt = (
				<div className="timeline-tooltip">
					<span>
						{aul.actionId!.name}: {aul.gainedAt.toFixed(2)}s - {aul.fadedAt.toFixed(2)}s
					</span>
				</div>
			);

			const tooltip = tippy(auraElem, {
				placement: 'bottom',
				content: tt,
			});
			this.addOnResetCallback(() => tooltip.destroy());

			aul.stacksChange.forEach((scl, i) => {
				if (scl.timestamp == aul.fadedAt) {
					return;
				}

				const stacksChangeElem = (
					<div
						className="rotation-timeline-stacks-change"
						style={{
							left: this.timeToPx(scl.timestamp - aul.timestamp),
							width: this.timeToPx(aul.stacksChange[i + 1] ? aul.stacksChange[i + 1].timestamp - scl.timestamp : aul.fadedAt - scl.timestamp),
						}}>
						{String(scl.newStacks)}
					</div>
				);
				auraElem.appendChild(stacksChangeElem);
			});
		});
	}

	private timeToPxValue(time: number): number {
		return time * 100;
	}
	private timeToPx(time: number): string {
		return this.timeToPxValue(time) + 'px';
	}

	private drawRotationTimeRuler(canvas: HTMLCanvasElement, duration: number) {
		const height = 30;
		canvas.width = this.timeToPxValue(duration);
		canvas.height = height;

		const ctx = canvas.getContext('2d')!;
		ctx.strokeStyle = 'white';

		ctx.font = 'bold 14px SimDefaultFont';
		ctx.fillStyle = 'white';
		ctx.lineWidth = 2;
		ctx.beginPath();

		// Bottom border line
		ctx.moveTo(0, height);
		ctx.lineTo(canvas.width, height);

		// Tick lines
		const numTicks = 1 + Math.floor(duration * 10);
		for (let i = 0; i <= numTicks; i++) {
			const time = i * 0.1;
			let x = this.timeToPxValue(time);
			if (i == 0) {
				ctx.textAlign = 'left';
				x++;
			} else if (i % 10 == 0 && time + 1 > duration) {
				ctx.textAlign = 'right';
				x--;
			} else {
				ctx.textAlign = 'center';
			}

			let lineHeight = 0;
			if (i % 10 == 0) {
				lineHeight = height * 0.5;
				ctx.fillText(time + 's', x, height - height * 0.6);
			} else if (i % 5 == 0) {
				lineHeight = height * 0.25;
			} else {
				lineHeight = height * 0.125;
			}
			ctx.moveTo(x, height);
			ctx.lineTo(x, height - lineHeight);
		}
		ctx.stroke();
	}

	private dpsTooltip(log: DpsLog, _includeAuras: boolean, player: UnitMetrics, colorOverride: string) {
		const showPlayerLabel = colorOverride != '';
		return (
			<div className="timeline-tooltip dps">
				<div className="timeline-tooltip-header">
					{showPlayerLabel ? (
						<>
							<img className="timeline-tooltip-icon" src="${player.iconUrl}" />
							<span className="" style="color: ${colorOverride}">
								{player.label}
							</span>
							<span> - </span>
						</>
					) : null}
					<span className="bold">{log.timestamp.toFixed(2)}s</span>
				</div>
				<div className="timeline-tooltip-body">
					<ul className="timeline-dps-events">{log.damageLogs.map(damageLog => this.tooltipLogItem(damageLog, damageLog.result())).join('')}</ul>
					<div className="timeline-tooltip-body-row">
						<span className="series-color">DPS: {log.dps.toFixed(2)}</span>
					</div>
				</div>
				{this.tooltipAurasSection(log)}
			</div>
		);
	}

	private threatTooltip(log: ThreatLogGroup, includeAuras: boolean, player: UnitMetrics, colorOverride: string) {
		const showPlayerLabel = colorOverride != '';
		return (
			<div className="timeline-tooltip threat">
				<div className="timeline-tooltip-header">
					{showPlayerLabel ? (
						<>
							<img className="timeline-tooltip-icon" src={player.iconUrl} />
							<span className="" style={{ color: colorOverride }}>
								{player.label}
							</span>
							<span> - </span>
						</>
					) : null}
					<span className="bold">{log.timestamp.toFixed(2)}s</span>
				</div>
				<div className="timeline-tooltip-body">
					<div className="timeline-tooltip-body-row">
						<span className="series-color">Before: {log.threatBefore.toFixed(1)}</span>
					</div>
					<ul className="timeline-threat-events">{log.logs.map(log => this.tooltipLogItem(log, <>{log.threat.toFixed(1)} Threat</>)).join('')}</ul>
					<div className="timeline-tooltip-body-row">
						<span className="series-color">After: {log.threatAfter.toFixed(1)}</span>
					</div>
				</div>
				{includeAuras ? this.tooltipAurasSection(log) : null}
			</div>
		);
	}

	private resourceTooltipElem(log: ResourceChangedLogGroup, maxValue: number, includeAuras: boolean) {
		const valToDisplayString = percentageResources.includes(log.resourceType)
			? (val: number) => `${val.toFixed(1)} (${((val / maxValue) * 100).toFixed(0)}%)`
			: (val: number) => `${val.toFixed(1)}`;

		return (
			<div className={`timeline-tooltip ${resourceNames.get(log.resourceType)!.toLowerCase().replaceAll(' ', '-')}`}>
				<div className="timeline-tooltip-header">
					<span className="bold">{log.timestamp.toFixed(2)}s</span>
				</div>
				<div className="timeline-tooltip-body">
					<div className="timeline-tooltip-body-row">
						<span className="series-color">Before: {valToDisplayString(log.valueBefore)}</span>
					</div>
					<ul className="timeline-mana-events">
						{log.logs.map(manaChangedLog => this.tooltipLogItemElem(manaChangedLog, <>{manaChangedLog.resultString()}</>))}
					</ul>
					<div className="timeline-tooltip-body-row">
						<span className="series-color">After: {valToDisplayString(log.valueAfter)}</span>
					</div>
				</div>
				{includeAuras && this.tooltipAurasSectionElem(log)}
			</div>
		);
	}

	private resourceTooltip(log: ResourceChangedLogGroup, maxValue: number, includeAuras: boolean) {
		return this.resourceTooltipElem(log, maxValue, includeAuras);
	}

	private tooltipLogItem(log: SimLog, value: Element) {
		return this.tooltipLogItemElem(log, value);
	}

	private tooltipLogItemElem(log: SimLog, value: Element): JSX.Element {
		return (
			<li>
				{log.actionId && log.actionId.iconUrl && <img className="timeline-tooltip-icon" src={log.actionId.iconUrl}></img>}
				{log.actionId && <span>{log.actionId.name}</span>}
				<span className="series-color">{value}</span>
			</li>
		);
	}

	private tooltipAurasSection(log: SimLog) {
		if (log.activeAuras.length == 0) {
			return '';
		}
		return this.tooltipAurasSectionElem(log);
	}

	private tooltipAurasSectionElem(log: SimLog): JSX.Element {
		if (log.activeAuras.length == 0) {
			return <></>;
		}

		return (
			<div className="timeline-tooltip-auras">
				<div className="timeline-tooltip-body-row">
					<span className="bold">Active Auras</span>
				</div>
				<ul className="timeline-active-auras">
					{log.activeAuras.map(auraLog => (
						<li>
							{auraLog.actionId!.iconUrl && <img className="timeline-tooltip-icon" src={auraLog.actionId!.iconUrl}></img>}
							<span>{auraLog.actionId!.name}</span>
						</li>
					))}
				</ul>
			</div>
		);
	}

	render() {
		this.reset();
		this.dpsResourcesPlot.render();
		this.rendered = true;
		this.updatePlot();
	}

	addOnResetCallback(callback: () => void) {
		this.resetCallbacks.push(callback);
	}

	reset() {
		this.resetCallbacks.forEach(callback => callback());
		this.resetCallbacks = [];
	}
}

const MELEE_ACTION_CATEGORY = 1;
const SPELL_ACTION_CATEGORY = 2;
const DEFAULT_ACTION_CATEGORY = 3;

// Hard-coded spell categories for controlling rotation ordering.
const idToCategoryMap: Record<number, number> = {
	[OtherAction.OtherActionShoot]: 0,
	[OtherAction.OtherActionMove]: 0.33,
	[OtherAction.OtherActionAttack]: 0.66,

	// Druid
	[48480]: 0.1, // Maul
	[48564]: MELEE_ACTION_CATEGORY + 0.1, // Mangle (Bear)
	[48568]: MELEE_ACTION_CATEGORY + 0.2, // Lacerate
	[48562]: MELEE_ACTION_CATEGORY + 0.3, // Swipe (Bear)

	[48566]: MELEE_ACTION_CATEGORY + 0.1, // Mangle (Cat)
	[48572]: MELEE_ACTION_CATEGORY + 0.2, // Shred
	[49800]: MELEE_ACTION_CATEGORY + 0.51, // Rip
	[52610]: MELEE_ACTION_CATEGORY + 0.52, // Savage Roar
	[48577]: MELEE_ACTION_CATEGORY + 0.53, // Ferocious Bite

	[48465]: SPELL_ACTION_CATEGORY + 0.1, // Starfire
	[48461]: SPELL_ACTION_CATEGORY + 0.2, // Wrath
	[439748]: SPELL_ACTION_CATEGORY + 0.3, // Starfall
	[48468]: SPELL_ACTION_CATEGORY + 0.4, // Insect Swarm
	[48463]: SPELL_ACTION_CATEGORY + 0.5, // Moonfire

	// Hunter
	[48996]: 0.1, // Raptor Strike
	[53217]: 0.6, // Wild Quiver
	[53209]: MELEE_ACTION_CATEGORY + 0.1, // Chimera Shot
	[53353]: MELEE_ACTION_CATEGORY + 0.11, // Chimera Shot Serpent
	[60053]: MELEE_ACTION_CATEGORY + 0.1, // Explosive Shot
	[49050]: MELEE_ACTION_CATEGORY + 0.2, // Aimed Shot
	[49048]: MELEE_ACTION_CATEGORY + 0.21, // Multi Shot
	[49045]: MELEE_ACTION_CATEGORY + 0.22, // Arcane Shot
	[49052]: MELEE_ACTION_CATEGORY + 0.27, // Steady Shot
	[61006]: MELEE_ACTION_CATEGORY + 0.28, // Kill Shot
	[34490]: MELEE_ACTION_CATEGORY + 0.29, // Silencing Shot
	[49001]: MELEE_ACTION_CATEGORY + 0.3, // Serpent Sting
	[3043]: MELEE_ACTION_CATEGORY + 0.3, // Scorpid Sting
	[53238]: MELEE_ACTION_CATEGORY + 0.31, // Piercing Shots
	[63672]: MELEE_ACTION_CATEGORY + 0.32, // Black Arrow
	[49067]: MELEE_ACTION_CATEGORY + 0.33, // Explosive Trap
	[58434]: MELEE_ACTION_CATEGORY + 0.34, // Volley

	// Paladin
	[35395]: MELEE_ACTION_CATEGORY + 0.1, // Crusader Strike
	[53385]: MELEE_ACTION_CATEGORY + 0.2, // Divine Storm
	[42463]: MELEE_ACTION_CATEGORY + 0.3, // Seal of Vengeance
	[61840]: MELEE_ACTION_CATEGORY + 0.4, // Righteous Vengeance
	[61411]: MELEE_ACTION_CATEGORY + 0.5, // Shield of Righteousness
	[53595]: MELEE_ACTION_CATEGORY + 0.51, // Hammer of Righteousness
	[20182]: MELEE_ACTION_CATEGORY + 0.6, // Reckoning
	[48952]: SPELL_ACTION_CATEGORY + 0.1, // Holy Shield
	[31803]: SPELL_ACTION_CATEGORY + 0.2, // Holy Vengeance
	[48801]: SPELL_ACTION_CATEGORY + 0.3, // Exorcism
	[48819]: SPELL_ACTION_CATEGORY + 0.4, // Consecration
	[53408]: SPELL_ACTION_CATEGORY + 0.51, // Judgement of Wisdom
	[20271]: SPELL_ACTION_CATEGORY + 0.52, // Judgement of Light
	[31804]: SPELL_ACTION_CATEGORY + 0.53, // Judgement of Vengeance
	[20467]: SPELL_ACTION_CATEGORY + 0.54, // Judgement of Command
	[20187]: SPELL_ACTION_CATEGORY + 0.55, // Judgement of Righteousness
	[31878]: SPELL_ACTION_CATEGORY + 0.56, // Judgements of the Wise
	[48817]: SPELL_ACTION_CATEGORY + 0.5, // Holy Wrath
	[48806]: SPELL_ACTION_CATEGORY + 0.6, // Hammer of Wrath
	[54428]: SPELL_ACTION_CATEGORY + 0.7, // Divine Plea
	[66233]: SPELL_ACTION_CATEGORY + 0.8, // Ardent Defender

	// Priest
	[48300]: SPELL_ACTION_CATEGORY + 0.11, // Devouring Plague
	[48125]: SPELL_ACTION_CATEGORY + 0.12, // Shadow Word: Pain
	[48160]: SPELL_ACTION_CATEGORY + 0.13, // Vampiric Touch
	[48135]: SPELL_ACTION_CATEGORY + 0.14, // Holy Fire
	[48123]: SPELL_ACTION_CATEGORY + 0.19, // Smite
	[48127]: SPELL_ACTION_CATEGORY + 0.2, // Mind Blast
	[48158]: SPELL_ACTION_CATEGORY + 0.3, // Shadow Word: Death
	[48156]: SPELL_ACTION_CATEGORY + 0.4, // Mind Flay

	// Rogue
	[6774]: MELEE_ACTION_CATEGORY + 0.1, // Slice and Dice
	[8647]: MELEE_ACTION_CATEGORY + 0.2, // Expose Armor
	[48672]: MELEE_ACTION_CATEGORY + 0.3, // Rupture
	[57993]: MELEE_ACTION_CATEGORY + 0.3, // Envenom
	[48668]: MELEE_ACTION_CATEGORY + 0.4, // Eviscerate
	[48666]: MELEE_ACTION_CATEGORY + 0.5, // Mutilate
	[48665]: MELEE_ACTION_CATEGORY + 0.6, // Mutilate (MH)
	[48664]: MELEE_ACTION_CATEGORY + 0.7, // Mutilate (OH)
	[48638]: MELEE_ACTION_CATEGORY + 0.5, // Sinister Strike
	[51723]: MELEE_ACTION_CATEGORY + 0.8, // Fan of Knives
	[57973]: SPELL_ACTION_CATEGORY + 0.1, // Deadly Poison
	[57968]: SPELL_ACTION_CATEGORY + 0.2, // Instant Poison

	// Shaman
	[58804]: 0.11, // Windfury Weapon
	[58790]: 0.12, // Flametongue Weapon
	[58796]: 0.12, // Frostbrand Weapon
	[17364]: MELEE_ACTION_CATEGORY + 0.1, // Stormstrike
	[60103]: MELEE_ACTION_CATEGORY + 0.2, // Lava Lash
	[49233]: SPELL_ACTION_CATEGORY + 0.21, // Flame Shock
	[49231]: SPELL_ACTION_CATEGORY + 0.22, // Earth Shock
	[49236]: SPELL_ACTION_CATEGORY + 0.23, // Frost Shock
	[60043]: SPELL_ACTION_CATEGORY + 0.31, // Lava Burst
	[49238]: SPELL_ACTION_CATEGORY + 0.32, // Lightning Bolt
	[49271]: SPELL_ACTION_CATEGORY + 0.33, // Chain Lightning
	[61657]: SPELL_ACTION_CATEGORY + 0.41, // Fire Nova
	[58734]: SPELL_ACTION_CATEGORY + 0.42, // Magma Totem
	[58704]: SPELL_ACTION_CATEGORY + 0.43, // Searing Totem
	[49281]: SPELL_ACTION_CATEGORY + 0.51, // Lightning Shield
	[49279]: SPELL_ACTION_CATEGORY + 0.52, // Lightning Shield (Proc)
	[2825]: DEFAULT_ACTION_CATEGORY + 0.1, // Bloodlust

	// Warlock
	[47867]: SPELL_ACTION_CATEGORY + 0.01, // Curse of Doom
	[47864]: SPELL_ACTION_CATEGORY + 0.02, // Curse of Agony
	[47813]: SPELL_ACTION_CATEGORY + 0.1, // Corruption
	[59164]: SPELL_ACTION_CATEGORY + 0.2, // Haunt
	[47843]: SPELL_ACTION_CATEGORY + 0.3, // Unstable Affliction
	[47811]: SPELL_ACTION_CATEGORY + 0.31, // Immolate
	[17962]: SPELL_ACTION_CATEGORY + 0.32, // Conflagrate
	[59172]: SPELL_ACTION_CATEGORY + 0.49, // Chaos Bolt
	[47809]: SPELL_ACTION_CATEGORY + 0.5, // Shadow Bolt
	[47838]: SPELL_ACTION_CATEGORY + 0.51, // Incinerate
	[47825]: SPELL_ACTION_CATEGORY + 0.52, // Soul Fire
	[47855]: SPELL_ACTION_CATEGORY + 0.6, // Drain Soul
	[57946]: SPELL_ACTION_CATEGORY + 0.7, // Life Tap
	[47241]: SPELL_ACTION_CATEGORY + 0.8, // Metamorphosis
	[50589]: SPELL_ACTION_CATEGORY + 0.81, // Immolation Aura
	[47193]: SPELL_ACTION_CATEGORY + 0.82, // Demonic Empowerment

	// Mage
	[42842]: SPELL_ACTION_CATEGORY + 0.01, // Frostbolt
	[47610]: SPELL_ACTION_CATEGORY + 0.02, // Frostfire Bolt
	[42897]: SPELL_ACTION_CATEGORY + 0.02, // Arcane Blast
	[42833]: SPELL_ACTION_CATEGORY + 0.02, // Fireball
	[42859]: SPELL_ACTION_CATEGORY + 0.03, // Scorch
	[42891]: SPELL_ACTION_CATEGORY + 0.1, // Pyroblast
	[10212]: SPELL_ACTION_CATEGORY + 0.1, // Arcane Missiles
	[44572]: SPELL_ACTION_CATEGORY + 0.1, // Deep Freeze
	[44781]: SPELL_ACTION_CATEGORY + 0.2, // Arcane Barrage
	[42914]: SPELL_ACTION_CATEGORY + 0.2, // Ice Lance
	[55360]: SPELL_ACTION_CATEGORY + 0.2, // Living Bomb
	[55362]: SPELL_ACTION_CATEGORY + 0.21, // Living Bomb (Explosion)
	[12654]: SPELL_ACTION_CATEGORY + 0.3, // Ignite
	[12472]: SPELL_ACTION_CATEGORY + 0.4, // Icy Veins
	[11129]: SPELL_ACTION_CATEGORY + 0.4, // Combustion
	[12042]: SPELL_ACTION_CATEGORY + 0.4, // Arcane Power
	[11958]: SPELL_ACTION_CATEGORY + 0.41, // Cold Snap
	[12043]: SPELL_ACTION_CATEGORY + 0.41, // Presence of Mind
	[31687]: SPELL_ACTION_CATEGORY + 0.41, // Water Elemental
	[55342]: SPELL_ACTION_CATEGORY + 0.5, // Mirror Image
	[33312]: SPELL_ACTION_CATEGORY + 0.51, // Mana Gems
	[12051]: SPELL_ACTION_CATEGORY + 0.52, // Evocate
	[44401]: SPELL_ACTION_CATEGORY + 0.6, // Missile Barrage
	[44448]: SPELL_ACTION_CATEGORY + 0.6, // Hot Streak
	[44545]: SPELL_ACTION_CATEGORY + 0.6, // Fingers of Frost
	[44549]: SPELL_ACTION_CATEGORY + 0.61, // Brain Freeze
	[12536]: SPELL_ACTION_CATEGORY + 0.61, // Clearcasting

	// Warrior
	[845]: 0.1, // Cleave
	[11608]: 0.1, // Cleave
	[11609]: 0.1, // Cleave
	[20569]: 0.1, // Cleave
	[1608]: 0.1, // Heroic Strike
	[11565]: 0.1, // Heroic Strike
	[11566]: 0.1, // Heroic Strike
	[11567]: 0.1, // Heroic Strike
	[8820]: MELEE_ACTION_CATEGORY + 0.05, // Slam
	[11604]: MELEE_ACTION_CATEGORY + 0.05, // Slam
	[11605]: MELEE_ACTION_CATEGORY + 0.05, // Slam
	[23881]: MELEE_ACTION_CATEGORY + 0.1, // Bloodthirst
	[47486]: MELEE_ACTION_CATEGORY + 0.1, // Mortal Strike
	[30356]: MELEE_ACTION_CATEGORY + 0.1, // Shield Slam
	[47498]: MELEE_ACTION_CATEGORY + 0.21, // Devastate
	[47467]: MELEE_ACTION_CATEGORY + 0.22, // Sunder Armor
	[57823]: MELEE_ACTION_CATEGORY + 0.23, // Revenge
	[1680]: MELEE_ACTION_CATEGORY + 0.24, // Whirlwind
	[7384]: MELEE_ACTION_CATEGORY + 0.25, // Overpower
	[47471]: MELEE_ACTION_CATEGORY + 0.42, // Execute
	[12867]: SPELL_ACTION_CATEGORY + 0.51, // Deep Wounds
	[58874]: SPELL_ACTION_CATEGORY + 0.52, // Damage Shield
	[46924]: SPELL_ACTION_CATEGORY + 0.61, // Bladestorm
	[2565]: SPELL_ACTION_CATEGORY + 0.62, // Shield Block
	[64382]: SPELL_ACTION_CATEGORY + 0.65, // Shattering Throw
	[71]: DEFAULT_ACTION_CATEGORY + 0.1, // Defensive Stance
	[2457]: DEFAULT_ACTION_CATEGORY + 0.1, // Battle Stance
	[2458]: DEFAULT_ACTION_CATEGORY + 0.1, // Berserker Stance

	// Generic
	[53307]: SPELL_ACTION_CATEGORY + 0.931, // Thorns
	[54043]: SPELL_ACTION_CATEGORY + 0.932, // Retribution Aura
	[54758]: SPELL_ACTION_CATEGORY + 0.933, // Hyperspeed Acceleration
	[42641]: SPELL_ACTION_CATEGORY + 0.941, // Sapper
	[40536]: SPELL_ACTION_CATEGORY + 0.942, // Explosive Decoy
	[41119]: SPELL_ACTION_CATEGORY + 0.943, // Saronite Bomb
	[40771]: SPELL_ACTION_CATEGORY + 0.944, // Cobalt Frag Bomb
};

const idsToGroupForRotation: Array<number> = [
	6774, // Slice and Dice
	8647, // Expose Armor
	48668, // Eviscerate
	48672, // Rupture
	51690, // Killing Spree
	57993, // Envenom
];

const percentageResources: Array<ResourceType> = [ResourceType.ResourceTypeHealth, ResourceType.ResourceTypeMana];
