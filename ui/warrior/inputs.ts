import * as InputHelpers from '../core/components/input_helpers.js';
import { ItemSlot, Spec } from '../core/proto/common.js';
import { WarriorRune, WarriorShout, WarriorStance } from '../core/proto/warrior.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { TypedEvent } from '../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const StanceSnapshot = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'stanceSnapshot',
	label: 'Stance Snapshot',
	labelTooltip: 'Ability that is cast at the same time as stance swap will benefit from the bonus of the stance before the swap.',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'shout',
	actionId: () => ActionId.fromSpellId(6673),
	value: WarriorShout.WarriorShoutBattle,
});

export const StancePicker = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarrior, WarriorStance>({
	fieldName: 'stance',
	values: [
		{ value: WarriorStance.WarriorStanceNone },
		{ value: WarriorStance.WarriorStanceBattle, actionId: () => ActionId.fromSpellId(2457) },
		{ value: WarriorStance.WarriorStanceDefensive, actionId: () => ActionId.fromSpellId(71) },
		{ value: WarriorStance.WarriorStanceBerserker, actionId: player => player.getMatchingSpellActionId([{ id: 2458, minLevel: 30 }]) },
		{
			value: WarriorStance.WarriorStanceGladiator,
			actionId: player => player.getMatchingSpellActionId([{ id: 412513, minLevel: 45 }]),
			showWhen: player => player.hasRune(ItemSlot.ItemSlotFeet, WarriorRune.RuneGladiatorStance),
		},
	],
	changeEmitter: player => TypedEvent.onAny([player.specOptionsChangeEmitter, player.gearChangeEmitter]),
	tooltip: 'Starting Stance<br />If blank, automatically chooses a stance based on your talents and runes.',
});
