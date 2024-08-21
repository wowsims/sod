import { Spec } from '../../proto/common.js';
import { RaidSimResultsManager } from '../raid_sim_action';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component';

export class ToplineResults extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'topline-results-root';
		super(config);

		this.rootElem.classList.add('results-sim');
	}

	onSimResult(resultData: SimResultData) {
		const noManaSpecs = [Spec.SpecFeralTankDruid, Spec.SpecRogue, Spec.SpecWarrior, Spec.SpecTankWarrior];
		const players = resultData.result.getRaidIndexedPlayers(resultData.filter);

		const content = RaidSimResultsManager.makeToplineResultsContent(resultData.result, resultData.filter, {
			showOutOfMana: players.length === 1 && !noManaSpecs.includes(players[0].spec),
		});

		this.rootElem.replaceChildren(content);
	}
}
