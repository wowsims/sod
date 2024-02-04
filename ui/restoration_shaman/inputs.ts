import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';

import {
	ShamanShield,
} from '../core/proto/shaman.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecElementalShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
		{
			actionId: (player: Player<Spec>) => player.getMatchingSpellActionId([
				{ id: 324, 		minLevel: 8, 	maxLevel: 15 	},
				{ id: 325, 		minLevel: 16, maxLevel: 23 	},
				{ id: 905, 		minLevel: 24, maxLevel: 31 	},
				{ id: 945, 		minLevel: 32, maxLevel: 39 	},
				{ id: 8134, 	minLevel: 40, maxLevel: 47	},
				{ id: 10431, 	minLevel: 48, maxLevel: 55 	},
				{ id: 10432, 	minLevel: 56								},
			]),
			value: ShamanShield.LightningShield },
	],
});
