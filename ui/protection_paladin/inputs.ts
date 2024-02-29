import { Spec } from '../core/proto/common.js';

import {
	PaladinAura,
	PaladinSeal,
} from '../core/proto/paladin.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const AuraSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin>({
	fieldName: 'aura',
	label: 'Aura',
	values: [
		{ name: 'None', value: PaladinAura.NoPaladinAura },
		{ name: 'Devotion Aura', value: PaladinAura.DevotionAura },
		{ name: 'Retribution Aura', value: PaladinAura.RetributionAura },
	],
});

// export const StartingSealSelection = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecProtectionPaladin>({
// 	fieldName: 'seal',
// 	label: 'Seal',
// 	labelTooltip: 'The seal active before encounter',
// 	values: [
// 		{ name: 'Righteousness', value: PaladinSeal.Righteousness },
// 		{ name: 'Command', value: PaladinSeal.Command },
// 	],
// });

// export const UseAvengingWrath = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecProtectionPaladin>({
// 	fieldName: 'useAvengingWrath',
// 	label: 'Use Avenging Wrath',
// });
