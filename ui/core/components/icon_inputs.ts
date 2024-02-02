import {
	Consumes,
	Debuffs,
	Faction,
	IndividualBuffs,
	RaidBuffs,
} from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';

import { Player } from '../player';
import { Spec } from '../proto/common';
import { Raid } from '../raid';
import { EventID, TypedEvent } from '../typed_event';

import { IconEnumPicker, IconEnumPickerDirection, IconEnumValueConfig } from './icon_enum_picker';
import { IconPicker } from './icon_picker';

import * as InputHelpers from './input_helpers';

// Component Functions

export type IconInputConfig<ModObject, T> = (
	InputHelpers.TypedIconPickerConfig<ModObject, T> |
	InputHelpers.TypedIconEnumPickerConfig<ModObject, T>
);

export const buildIconInput = (parent: HTMLElement, player: Player<Spec>, inputConfig: IconInputConfig<Player<Spec>, any>) => {
	if (inputConfig.type == 'icon') {
		return new IconPicker<Player<Spec>, any>(parent, player, inputConfig);
	} else {
		return new IconEnumPicker<Player<Spec>, any>(parent, player, inputConfig);
	}
};

export function withLabel<ModObject, T>(config: IconInputConfig<ModObject, T>, label: string): IconInputConfig<ModObject, T> {
	config.label = label;
	return config;
}

interface BooleanInputConfig<T> {
	actionId: (player: Player<Spec>) => ActionId | null
	fieldName: keyof T
	value?: number
	showWhen?: (player: Player<Spec>) => boolean
}

export function makeBooleanRaidBuffInput<SpecType extends Spec>(config: BooleanInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => !config.showWhen || config.showWhen(player),
		getValue: (player: Player<SpecType>) => player.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: RaidBuffs) => player.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.buffsChangeEmitter, player.levelChangeEmitter, player.raceChangeEmitter]),
	}, config.actionId, config.fieldName, config.value);
}
// export function makeBooleanPartyBuffInput<SpecType extends Spec>(config: BooleanInputConfig<PartyBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
// 	return InputHelpers.makeBooleanIconInput<any, PartyBuffs, Party>({
// 		getModObject: (player: Player<SpecType>) => player.getParty()!,
// 		getValue: (party: Party) => party.getBuffs(),
// 		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
// 		changeEmitter: (party: Party) => party.buffsChangeEmitter,
// 	}, config.actionId, config.fieldName, config.value);
// }

export function makeBooleanIndividualBuffInput<SpecType extends Spec>(config: BooleanInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => !config.showWhen || config.showWhen(player),
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.levelChangeEmitter, player.raceChangeEmitter]),
	}, config.actionId, config.fieldName, config.value);
}

export function makeBooleanConsumeInput<SpecType extends Spec>(config: BooleanInputConfig<Consumes>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Consumes, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => !config.showWhen || config.showWhen(player),
		getValue: (player: Player<SpecType>) => player.getConsumes(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: Consumes) => player.setConsumes(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter])
	}, config.actionId, config.fieldName, config.value);
}
export function makeBooleanDebuffInput<SpecType extends Spec>(config: BooleanInputConfig<Debuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Debuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => !config.showWhen || config.showWhen(player),
		getValue: (player: Player<SpecType>) => player.getRaid()!.getDebuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: Debuffs) => player.getRaid()!.setDebuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.debuffsChangeEmitter, player.levelChangeEmitter]),
	}, config.actionId, config.fieldName, config.value);
}

interface TristateInputConfig<T, ModObject> {
	actionId: (modObj: ModObject) => ActionId | null
	impId: ActionId
	fieldName: keyof T
	showWhen?: (modObj: ModObject) => boolean
}

export function makeTristateRaidBuffInput<SpecType extends Spec>(config: TristateInputConfig<RaidBuffs, Player<SpecType>>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeTristateIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => (!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: RaidBuffs) => player.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.buffsChangeEmitter, player.levelChangeEmitter, player.raceChangeEmitter]),
	}, config.actionId, config.impId, config.fieldName);
}

export function makeTristateIndividualBuffInput<SpecType extends Spec>(config: TristateInputConfig<IndividualBuffs,  Player<SpecType>>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeTristateIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => (!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.levelChangeEmitter, player.raceChangeEmitter])
	}, config.actionId, config.impId, config.fieldName);
}

export function makeTristateDebuffInput<SpecType extends Spec>(config: TristateInputConfig<Debuffs, Player<SpecType>>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeTristateIconInput<any, Debuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		getValue: (player: Player<SpecType>) => player.getRaid()!.getDebuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: Debuffs) => player.getRaid()!.setDebuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.debuffsChangeEmitter, player.levelChangeEmitter]),
	}, config.actionId, config.impId, config.fieldName);
}

// function makeQuadstateDebuffInput(actionId: ActionId, impId: ActionId, impId2: ActionId, fieldName: keyof Debuffs): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
// 	return InputHelpers.makeQuadstateIconInput<any, Debuffs, Raid>({
// 		getModObject: (player: Player<SpecType>) => player.getRaid()!,
// 		getValue: (raid: Raid) => raid.getDebuffs(),
// 		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
// 		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
// 	}, actionId, impId, impId2, fieldName);
// }

interface MultiStateInputConfig<T, ModObject> {
	actionId: (modObj: ModObject) => ActionId | null
	numStates: number
	fieldName: keyof T
	multiplier?: number
	faction?: Faction
}

export function makeMultistateRaidBuffInput<SpecType extends Spec>(config: MultiStateInputConfig<RaidBuffs, Player<SpecType>>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeMultistateIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => (!config.faction || config.faction == player.getFaction()),
		getValue: (player: Player<SpecType>) => player.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: RaidBuffs) => player.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.buffsChangeEmitter, player.levelChangeEmitter, player.raceChangeEmitter]),
	}, config.actionId, config.numStates, config.fieldName, config.multiplier);
}
// function makeMultistatePartyBuffInput(actionId: ActionId, numStates: number, fieldName: keyof PartyBuffs): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
// 	return InputHelpers.makeMultistateIconInput<any, PartyBuffs, Party>({
// 		getModObject: (player: Player<SpecType>) => player.getParty()!,
// 		getValue: (party: Party) => party.getBuffs(),
// 		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
// 		changeEmitter: (party: Party) => party.buffsChangeEmitter,
// 	}, actionId, numStates, fieldName);
// }
export function makeMultistateIndividualBuffInput<SpecType extends Spec>(config: MultiStateInputConfig<IndividualBuffs, Player<SpecType>>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => (!config.faction || config.faction == player.getFaction()),
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.levelChangeEmitter, player.raceChangeEmitter]),
	}, config.actionId, config.numStates, config.fieldName, config.multiplier);
}
// function makeMultistateMultiplierIndividualBuffInput(actionId: ActionId, numStates: number, multiplier: number, fieldName: keyof IndividualBuffs): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
// 	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<SpecType>>({
// 		getModObject: (player: Player<SpecType>) => player,
// 		getValue: (player: Player<SpecType>) => player.getBuffs(),
// 		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
// 		changeEmitter: (player: Player<SpecType>) => player.buffsChangeEmitter,
// 	}, actionId, numStates, fieldName, multiplier);
// }

export function makeMultistateMultiplierDebuffInput(config: MultiStateInputConfig<Debuffs, Raid>): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeMultistateIconInput<any, Debuffs, Raid>({
		getModObject: (player: Player<any>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getDebuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
	}, config.actionId, config.numStates, config.fieldName, config.multiplier);
}

interface EnumInputConfig<ModObject, Message, T> {
	fieldName: keyof Message
	values: Array<IconEnumValueConfig<ModObject, T>>
	direction?: IconEnumPickerDirection
	numColumns?: number
	faction?: Faction
}

export function makeEnumIndividualBuffInput<SpecType extends Spec>(config: EnumInputConfig<Player<SpecType>, IndividualBuffs, number>): InputHelpers.TypedIconEnumPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeEnumIconInput<any, IndividualBuffs, Player<SpecType>, number>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) => (!config.faction || config.faction == player.getFaction()),
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.levelChangeEmitter, player.raceChangeEmitter]),
	}, config.fieldName, config.values, config.numColumns, config.direction || IconEnumPickerDirection.Vertical)
};
