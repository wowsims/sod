import { Spec } from '../core/proto/common.js';
import { ShamanSyncType } from '../core/proto/shaman.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SyncTypeInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecEnhancementShaman>({
	fieldName: 'syncType',
	label: 'Sync/Stagger Setting',
	labelTooltip:
		`Choose your sync or stagger option Perfect
		<ul>
			<li><div>Auto: Will auto pick sync options based on your weapons attack speeds</div></li>
			<li><div>None: No Sync or Staggering, used for mismatched weapon speeds</div></li>
			<li><div>Perfect Sync: Makes your weapons always attack at the same time, for match weapon speeds</div></li>
			<li><div>Delayed Offhand: Adds a slight delay to the offhand attacks while staying within the 0.5s flurry ICD window</div></li>
		</ul>`,
	values: [
		{ name: "Automatic", value: ShamanSyncType.Auto },
		{ name: 'None', value: ShamanSyncType.NoSync },
		{ name: 'Perfect Sync', value: ShamanSyncType.SyncMainhandOffhandSwings },
		{ name: 'Delayed Offhand', value: ShamanSyncType.DelayOffhandSwings },
	],
});
