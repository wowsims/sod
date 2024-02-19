import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	RogueOptions_PoisonImbue as Poison,
} from '../core/proto/rogue.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
// TODO: change tooltip based on level and available rank of poison

export const MainHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'mhImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Main Hand Poison' },
		{ actionId: () => ActionId.fromItemId(2823), value: Poison.DeadlyPoison },
		{ actionId: () => ActionId.fromItemId(8679), value: Poison.InstantPoison },
		{ actionId: () => ActionId.fromItemId(11325), value: Poison.WoundPoison },
	],
});

export const OffHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'ohImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Off Hand Poison' },
		{ actionId: () => ActionId.fromItemId(2823), value: Poison.DeadlyPoison },
		{ actionId: () => ActionId.fromItemId(8679), value: Poison.InstantPoison },
		{ actionId: () => ActionId.fromItemId(11325), value: Poison.WoundPoison },
	],
});

export const ApplyPoisonsManually = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRogue>({
	fieldName: 'applyPoisonsManually',
	label: 'Configure poisons manually',
	labelTooltip: 'Prevent automatic poison configuration that is based on equipped weapons.',
});
