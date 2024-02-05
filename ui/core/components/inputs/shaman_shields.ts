import { Player } from '../../player.js';
import { Spec } from '../../proto/common.js';
import { ShamanShield } from '../../proto/shaman.js';
import { TypedEvent } from '../../typed_event.js';

import * as InputHelpers from '../../components/input_helpers.js';

type ShamanSpec = Spec.SpecElementalShaman | Spec.SpecEnhancementShaman | Spec.SpecRestorationShaman

export const ShamanShieldInput = <SpecType extends ShamanSpec>() =>
	InputHelpers.makeSpecOptionsEnumIconInput<SpecType, ShamanShield>({
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
		changeEmitter: (player) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.levelChangeEmitter]),
	});
