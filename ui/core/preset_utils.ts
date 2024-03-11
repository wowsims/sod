import * as Tooltips from './constants/tooltips.js';
import { Player } from './player';
import {
	APLRotation,
	APLRotation_Type as APLRotationType,
} from './proto/apl';
import {
	EquipmentSpec,
    Faction,
    Spec,
} from './proto/common';
import {
    SavedRotation, SavedTalents,
} from './proto/ui';
import {
    SpecRotation,
	specTypeFunctions,
} from './proto_utils/utils';

export interface PresetGear {
	name: string;
	gear: EquipmentSpec;
	tooltip?: string;
	enableWhen?: (obj: Player<any>) => boolean;
}

export interface PresetRotation {
	name: string;
	rotation: SavedRotation;
	tooltip?: string;
	enableWhen?: (obj: Player<any>) => boolean;
}

export interface PresetGearOptions {
    talentTree?: number,
    talentTrees?: Array<number>,
    faction?: Faction,
    customCondition?: (player: Player<any>) => boolean,

    tooltip?: string,
}

export interface PresetRotationOptions {
    talentTree?: number,
    customCondition?: (player: Player<any>) => boolean,
}

export function makePresetGear(name: string, gearJson: any, options?: PresetGearOptions): PresetGear {
    const gear = EquipmentSpec.fromJson(gearJson);
    return makePresetGearHelper(name, gear, options || {});
}

function makePresetGearHelper(name: string, gear: EquipmentSpec, options: PresetGearOptions): PresetGear {
    const conditions: Array<(player: Player<any>) => boolean> = [];
    if (options.talentTree != undefined) {
        conditions.push((player: Player<any>) => player.getTalentTree() == options.talentTree);
    }
    if (options.talentTrees != undefined) {
        conditions.push((player: Player<any>) => (options.talentTrees || []).includes(player.getTalentTree()));
    }
    if (options.faction != undefined) {
        conditions.push((player: Player<any>) => player.getFaction() == options.faction);
    }
    if (options.customCondition != undefined) {
        conditions.push(options.customCondition);
    }

    return {
        name: name,
        tooltip: options.tooltip || Tooltips.BASIC_BIS_DISCLAIMER,
        gear: gear,
        enableWhen: conditions.length > 0
            ? (player: Player<any>) => conditions.every(cond => cond(player))
            : undefined,
    };
}

export function makePresetAPLRotation(name: string, rotationJson: any, options?: PresetRotationOptions): PresetRotation {
    const rotation = SavedRotation.create({
        rotation: APLRotation.fromJson(rotationJson),
    });
    return makePresetRotationHelper(name, rotation, options);
}

export function makePresetSimpleRotation<SpecType extends Spec>(name: string, spec: SpecType, simpleRotation: SpecRotation<SpecType>, options?: PresetRotationOptions): PresetRotation {
    const rotation = SavedRotation.create({
		rotation: {
			type: APLRotationType.TypeSimple,
			simple: {
				specRotationJson: JSON.stringify(specTypeFunctions[spec].rotationToJson(simpleRotation)),
			},
		},
    });
    return makePresetRotationHelper(name, rotation, options);
}

export interface PresetTalentsOptions {
	customCondition?: (player: Player<any>) => boolean,
}

export function makePresetTalents(name: string, data: SavedTalents, options?: PresetTalentsOptions) {
	const conditions: Array<(player: Player<any>) => boolean> = [];
	if (options && options.customCondition) {
        conditions.push(options.customCondition);
    }

	return {
        name,
        data,
        enableWhen: conditions.length > 0
            ? (player: Player<any>) => conditions.every(cond => cond(player))
            : undefined,
    };
}

function makePresetRotationHelper(name: string, rotation: SavedRotation, options?: PresetRotationOptions): PresetRotation {
    const conditions: Array<(player: Player<any>) => boolean> = [];
    if (options?.talentTree != undefined) {
        conditions.push((player: Player<any>) => player.getTalentTree() == options.talentTree);
    }
    if (options?.customCondition != undefined) {
        conditions.push(options.customCondition);
    }

    return {
        name: name,
        rotation: rotation,
        enableWhen: conditions.length > 0
            ? (player: Player<any>) => conditions.every(cond => cond(player))
            : undefined,
    };
}
