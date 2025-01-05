import { TOOLTIP_METRIC_LABELS } from '../../constants/tooltips';
import { ActionMetrics } from '../../proto_utils/sim_result';
import { bucket, formatToCompactNumber, formatToNumber, formatToPercent } from '../../utils';
import { MetricsCombinedTooltipTable } from './metrics_table/metrics_combined_tooltip_table';
import { ColumnSortType, MetricsTable } from './metrics_table/metrics_table';
import { MetricsTotalBar } from './metrics_table/metrics_total_bar';
import { ResultComponentConfig, SimResultData } from './result_component';

export class DamageMetricsTable extends MetricsTable<ActionMetrics> {
	maxDamageAmount: number | null = null;
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'damage-metrics-root';
		config.resultsEmitter.on((_, resultData) => {
			const lastResult = resultData
				? this.getGroupedMetrics(resultData)
						.filter(g => g.length)
						.map(groups => this.mergeMetrics(groups))
				: undefined;
			this.maxDamageAmount = Math.max(...(lastResult || []).map(a => a.damage));
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
				name: 'Damage Done',
				headerCellClass: 'text-center metrics-table-cell--primary-metric',
				columnClass: 'metrics-table-cell--primary-metric',
				getValue: (metric: ActionMetrics) => metric.avgDamage,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<MetricsTotalBar
							spellSchool={metric.spellSchool}
							percentage={metric.totalDamagePercent}
							max={this.maxDamageAmount}
							total={metric.avgDamage}
							value={metric.damage}
						/>,
					);

					const hitValues = metric.damageDone.hit;
					const resistedHitValues = metric.damageDone.resistedHit;
					const critHitValues = metric.damageDone.critHit;
					const resistedCritHitValues = metric.damageDone.resistedCritHit;
					const tickValues = metric.damageDone.tick;
					const resistedTickValues = metric.damageDone.resistedTick;
					const critTickValues = metric.damageDone.critTick;
					const resistedCritTickValues = metric.damageDone.resistedCritTick;
					const glanceValues = metric.damageDone.glance;
					const blockValues = metric.damageDone.block;
					const blockedCritValues = metric.damageDone.blockedCrit;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							headerValues={[, 'Amount']}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.damage,
									totalPercentage: 100,
									data: [
										{
											name: 'Hit',
											...hitValues,
										},
										{
											name: 'Resisted Hit',
											...resistedHitValues,
										},
										{
											name: `Critical Hit`,
											...critHitValues,
										},
										{
											name: `Resisted Critical Hit`,
											...resistedCritHitValues,
										},
										{
											name: 'Tick',
											...tickValues,
										},
										{
											name: 'Resisted Tick',
											...resistedTickValues,
										},
										{
											name: `Critical Tick`,
											...critTickValues,
										},
										{
											name: `Resisted Critical Tick`,
											...resistedCritTickValues,
										},
										{
											name: 'Glancing Blow',
											...glanceValues,
										},
										{
											name: 'Blocked Hit',
											...blockValues,
										},
										{
											name: 'Blocked Critical Hit',
											...blockedCritValues,
										},
									],
								},
							]}
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
				tooltip: TOOLTIP_METRIC_LABELS['Damage Avg Cast'],
				getValue: (metric: ActionMetrics) => {
					if (metric.isPassiveAction) return 0;
					return metric.avgCastHit || metric.avgCastTick;
				},
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<>
							{metric.isPassiveAction ? (
								'-'
							) : (
								<>
									{formatToCompactNumber(metric.avgCastHit || metric.avgCastTick, { fallbackString: '-' })}
									{metric.avgCastHit && metric.avgCastTick ? (
										<> ({formatToCompactNumber(metric.avgCastTick, { fallbackString: '-' })})</>
									) : undefined}
								</>
							)}
						</>,
					);
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
					const relativeBlockedCritPercent = (metric.blockedCrits / metric.landedHits) * 100;

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
											name: 'Blocked Critical Hit',
											value: metric.blockedCrits,
											percentage: relativeBlockedCritPercent,
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
										},
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
				getValue: (metric: ActionMetrics) => metric.avgHit || metric.avgTick,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<>
							{formatToCompactNumber(metric.avgHit || metric.avgTick, { fallbackString: '-' })}
							{metric.avgHit && metric.avgTick ? <> ({formatToCompactNumber(metric.avgTick, { fallbackString: '-' })})</> : undefined}
						</>,
					);
				},
			},
			{
				name: 'Crit %',
				getValue: (metric: ActionMetrics) => metric.critPercent + metric.blockedCritPercent || metric.critTickPercent,
				getDisplayString: (metric: ActionMetrics) =>
					`${formatToPercent(metric.critPercent + metric.blockedCritPercent || metric.critTickPercent, { fallbackString: '-' })}${
						metric.critPercent + metric.blockedCritPercent && metric.critTickPercent
							? ` (${formatToPercent(metric.critTickPercent, { fallbackString: '-' })})`
							: ''
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
				name: 'DPET',
				getValue: (metric: ActionMetrics) => metric.damageThroughput,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.damageThroughput, { fallbackString: '-' }),
			},
			{
				name: 'DPS',
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.dps,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.dps, { minimumFractionDigits: 2, fallbackString: '-' })}</>);
					if (!metric.dps) return;
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

		const actions = player.getDamageActions().map(action => action.forTarget(resultData.filter));
		const actionGroups = ActionMetrics.groupById(actions);
		const petsByName = bucket(player.pets, pet => pet.name);

		const petGroups = Object.values(petsByName).map(pets =>
			ActionMetrics.joinById(
				pets.flatMap(pet => pet.getDamageActions().map(action => action.forTarget(resultData.filter))),
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
