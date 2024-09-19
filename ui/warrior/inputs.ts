// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
import * as InputHelpers from '../core/components/input_helpers.js';
import { Spec } from '../core/proto/common.js';

export const RendStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingRendStopAttack',
	label: 'Using Rend StopAttack Macro',
	labelTooltip: '/cast [@target] Rend \n/stopattack',
});

export const BloodthirstStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingBloodthirstStopAttack',
	label: 'Using Bloodthirst StopAttack Macro',
	labelTooltip: '/cast [@target] Bloodthirst \n/stopattack',
});

export const QuickStrikeStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingQuickStrikeStopAttack',
	label: 'Using Quick Strike StopAttack Macro',
	labelTooltip: '/cast [@target] Quick Strike \n/stopattack',
});

export const HamstringStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingHamstringStopAttack',
	label: 'Using Hamstring StopAttack Macro',
	labelTooltip: '/cast [@target] Hamstring \n/stopattack',
});

export const WhirlwindStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingWhirlwindStopAttack',
	label: 'Using Whirlwind StopAttack Macro',
	labelTooltip: '/cast [@target] Whirlwind \n/stopattack',
});

export const ExecuteStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingExecuteStopAttack',
	label: 'Using Execute StopAttack Macro',
	labelTooltip: '/cast [@target] Execute \n/stopattack',
});

export const OverpowerStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingOverpowerStopAttack',
	label: 'Using Overpower StopAttack Macro',
	labelTooltip: '/cast [@target] Overpower \n/stopattack',
});

export const HeroicStrikeStopAttack = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'isUsingHeroicStrikeStopAttack',
	label: 'Using Heroic Strike StopAttack Macro',
	labelTooltip: '/cast [@target] Heroic Strike \n/stopattack',
});
