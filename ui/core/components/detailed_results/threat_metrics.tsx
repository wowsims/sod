import { TOOLTIP_METRIC_LABELS } from '../../constants/tooltips';
import { ActionMetrics } from '../../proto_utils/sim_result';
import { bucket, formatToCompactNumber, formatToNumber, formatToPercent } from '../../utils';
import { MetricsCombinedTooltipTable } from './metrics_table/metrics_combined_tooltip_table';
import { ColumnSortType, MetricsTable } from './metrics_table/metrics_table';
import { MetricsTotalBar } from './metrics_table/metrics_total_bar';
import { ResultComponentConfig, SimResultData } from './result_component';

export class ThreatMetricsTable extends MetricsTable<ActionMetrics> {
	maxThreatAmount: number | null = null;
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'threat-metrics-root';
		config.resultsEmitter.on((_, resultData) => {
			const lastResult = resultData
				? this.getGroupedMetrics(resultData)
						.filter(g => g.length)
						.map(groups => this.mergeMetrics(groups))
				: undefined;
			this.maxThreatAmount = Math.max(...(lastResult || []).map(a => a.threat));
		});

		super(config, [
			MetricsTable.nameCellConfig((metric: ActionMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
					metricType: metric.constructor?.name,
				};
			}),
			{
				name: 'Threat done',
				headerCellClass: 'text-center metrics-table-cell--primary-metric',
				columnClass: 'metrics-table-cell--primary-metric',
				getValue: (metric: ActionMetrics) => metric.avgDamage,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<MetricsTotalBar
							spellSchool={metric.spellSchool}
							percentage={metric.totalThreatPercent}
							max={this.maxThreatAmount}
							total={metric.avgThreat}
							value={metric.threat}
						/>,
					);
				},
			},
			{
				name: 'Casts',
				getValue: (metric: ActionMetrics) => metric.casts,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.casts, { fallbackString: '-' })}</>);

					if (metric.isPassiveAction || (!metric.landedHits && !metric.totalMisses)) return;
					const relativeHitPercent = ((metric.landedHits || metric.casts) / ((metric.landedHits || metric.casts) + metric.totalMisses)) * 100;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.casts,
									totalPercentage: 100,
									data: [
										{
											name: 'Hits',
											value: metric.landedHits || metric.casts - metric.totalMisses,
											percentage: relativeHitPercent,
										},
										{
											name: 'Miss',
											value: metric.misses,
											percentage: metric.missPercent,
										},
										{
											name: 'Parry',
											value: metric.parries,
											percentage: metric.parryPercent,
										},
										{
											name: 'Dodge',
											value: metric.dodges,
											percentage: metric.dodgePercent,
										},
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'Avg Cast',
				tooltip: TOOLTIP_METRIC_LABELS['Avg Cast Threat'],
				getValue: (metric: ActionMetrics) => {
					if (metric.isPassiveAction) return 0;
					return metric.avgCastThreat;
				},
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{metric.isPassiveAction ? '-' : <>{formatToCompactNumber(metric.avgCastThreat, { fallbackString: '-' })}</>}</>);
				},
			},
			{
				name: 'Hits',
				getValue: (metric: ActionMetrics) => metric.landedHits || metric.landedTicks,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<>
							{formatToNumber(metric.landedHits || metric.landedTicks, { fallbackString: '-' })}
							{metric.landedHits && metric.landedTicks ? <> ({formatToNumber(metric.landedTicks, { fallbackString: '-' })})</> : undefined}
						</>,
					);
					if (!metric.landedHits && !metric.landedTicks) return;

					const relativeHitPercent = ((metric.hits - metric.resistedHits) / metric.landedHits) * 100;
					const relativeResistedHitPercent = (metric.resistedHits / metric.landedHits) * 100;
					const relativeCritPercent = ((metric.crits - metric.resistedCrits) / metric.landedHits) * 100;
					const relativeResistedCritPercent = (metric.resistedCrits / metric.landedHits) * 100;
					const relativeTickPercent = ((metric.ticks - metric.resistedTicks) / metric.landedTicks) * 100;
					const relativeResistedTickPercent = (metric.resistedTicks / metric.landedTicks) * 100;
					const relativeCritTickPercent = ((metric.critTicks - metric.resistedCritTicks) / metric.landedTicks) * 100;
					const relativeResistedCritTickPercent = (metric.resistedCritTicks / metric.landedTicks) * 100;
					const relativeGlancePercent = (metric.glances / metric.landedHits) * 100;
					const relativeBlockPercent = (metric.blocks / metric.landedHits) * 100;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.landedHits,
									totalPercentage: 100,
									data: [
										{
											name: 'Hit',
											value: metric.hits - metric.resistedHits,
											percentage: relativeHitPercent,
										},
										{
											name: 'Resisted Hit',
											value: metric.resistedHits,
											percentage: relativeResistedHitPercent,
										},
										{
											name: `Critical Hit`,
											value: metric.crits - metric.resistedCrits,
											percentage: relativeCritPercent,
										},
										{
											name: `Resisted Critical Hit`,
											value: metric.resistedCrits,
											percentage: relativeResistedCritPercent,
										},
										{
											name: 'Glancing Blow',
											value: metric.glances,
											percentage: relativeGlancePercent,
										},
										{
											name: 'Blocked Hit',
											value: metric.blocks,
											percentage: relativeBlockPercent,
										}
									],
								},
								{
									spellSchool: metric.spellSchool,
									total: metric.landedTicks,
									totalPercentage: 100,
									data: [
										{
											name: 'Tick',
											value: metric.ticks - metric.resistedTicks,
											percentage: relativeTickPercent,
										},
										{
											name: 'Resisted Tick',
											value: metric.resistedTicks,
											percentage: relativeResistedTickPercent,
										},
										{
											name: `Critical Tick`,
											value: metric.critTicks - metric.resistedCritTicks,
											percentage: relativeCritTickPercent,
										},
										{
											name: 'Resisted Critical Tick',
											value: metric.resistedCritTicks,
											percentage: relativeResistedCritTickPercent,
										},
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'Avg Hit',
				tooltip: TOOLTIP_METRIC_LABELS['Avg Hit Threat'],
				getValue: (metric: ActionMetrics) => metric.avgHitThreat,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToCompactNumber(metric.avgHitThreat, { fallbackString: '-' })}</>);
				},
			},
			{
				name: 'Crit %',
				getValue: (metric: ActionMetrics) => metric.critPercent || metric.critTickPercent,
				getDisplayString: (metric: ActionMetrics) =>
					`${formatToPercent(metric.critPercent || metric.critTickPercent, { fallbackString: '-' })}${
						metric.critPercent && metric.critTickPercent ? ` (${formatToPercent(metric.critTickPercent, { fallbackString: '-' })})` : ''
					}`,
			},
			{
				name: 'Miss %',
				getValue: (metric: ActionMetrics) => metric.totalMissesPercent,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToPercent(metric.totalMissesPercent, { fallbackString: '-' })}</>);
					if (!metric.totalMissesPercent) return;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.totalMisses,
									totalPercentage: metric.totalMissesPercent,
									data: [
										{
											name: 'Miss',
											value: metric.misses,
											percentage: metric.missPercent,
										},
										{
											name: 'Parry',
											value: metric.parries,
											percentage: metric.parryPercent,
										},
										{
											name: 'Dodge',
											value: metric.dodges,
											percentage: metric.dodgePercent,
										},
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'TPET',
				getValue: (metric: ActionMetrics) => metric.threatThroughput,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.threatThroughput, { fallbackString: '-' }),
			},
			{
				name: 'TPS',
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.tps,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.tps, { minimumFractionDigits: 2, fallbackString: '-' })}</>);
					if (!metric.tps) return;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							headerValues={[, 'Amount']}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.tps,
									totalPercentage: 100,
									data: [
										{
											name: 'Per Cast',
											value: metric.avgCastThreat,
											percentage: (metric.avgCastThreat / Math.max(metric.avgCastThreat, metric.avgHitThreat)) * 100,
										},
										{
											name: 'Per Hit',
											value: metric.avgHitThreat,
											percentage: (metric.avgHitThreat / Math.max(metric.avgCastThreat, metric.avgHitThreat)) * 100,
										},
									],
								},
							]}
						/>,
					);
				},
			},
		]);
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<ActionMetrics>> {
		const players = resultData.result.getRaidIndexedPlayers(resultData.filter);
		if (players.length != 1) {
			return [];
		}
		const player = players[0];

		const actions = player.getThreatActions().map(action => action.forTarget(resultData.filter));
		const actionGroups = ActionMetrics.groupById(actions);
		const petsByName = bucket(player.pets, pet => pet.name);

		const petGroups = Object.values(petsByName).map(pets =>
			ActionMetrics.joinById(
				pets.flatMap(pet => pet.getThreatActions().map(action => action.forTarget(resultData.filter))),
				true,
			),
		);

		return actionGroups.concat(petGroups);
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		return ActionMetrics.merge(metrics, {
			removeTag: true,
			actionIdOverride: metrics[0]?.unit?.petActionId || undefined,
		});
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
