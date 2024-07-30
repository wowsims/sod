import { Spec } from '../../proto/common.js';
import { RaidSimResultsManager } from '../raid_sim_action.jsx';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class ToplineResults extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'topline-results-root';
		super(config);

		this.rootElem.classList.add('results-sim');
	}

	onSimResult(resultData: SimResultData) {
		const content = RaidSimResultsManager.makeToplineResultsContent(resultData.result, resultData.filter);

		const noManaSpecs = [Spec.SpecFeralTankDruid, Spec.SpecRogue, Spec.SpecWarrior, Spec.SpecTankWarrior];

		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length == 1 && !noManaSpecs.includes(players[0].spec)) {
			const player = players[0];
			const secondsOOM = player.secondsOomAvg;
			const percentOOM = secondsOOM / resultData.result.encounterMetrics.durationSeconds;
			const dangerLevel = percentOOM < 0.01 ? 'safe' : percentOOM < 0.05 ? 'warning' : 'danger';

			content.appendChild(
				<div className={`results-sim-percent-oom ${dangerLevel} damage-metrics`}>
					<span className="topline-result-avg">{secondsOOM.toFixed(1)}s</span>
				</div>,
			);
		}

		this.rootElem.innerHTML = '';
		this.rootElem.appendChild(content);
	}
}
