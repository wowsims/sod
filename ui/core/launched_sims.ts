import { Phase } from './constants/other';
import { Class, Spec } from './proto/common';
import { specToClass } from './proto_utils/utils';

// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!

export enum LaunchStatus {
	Unlaunched,
	Alpha,
	Beta,
};

export type SimStatus = {
	phase: Phase,
	status: LaunchStatus
}

export const raidSimStatus: SimStatus = {
	phase: Phase.Phase1,
	status: LaunchStatus.Unlaunched,
}

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, SimStatus> = {
	[Spec.SpecBalanceDruid]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFeralDruid]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFeralTankDruid]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecRestorationDruid]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecElementalShaman]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Beta,
	},
	[Spec.SpecEnhancementShaman]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Beta,
	},
	[Spec.SpecRestorationShaman]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecHunter]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecMage]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecRogue]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecHolyPaladin]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecProtectionPaladin]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecRetributionPaladin]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecHealingPriest]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecShadowPriest]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecWarlock]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecTankWarlock]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecWarrior]: {
		phase: Phase.Phase2,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecProtectionWarrior]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
};

export function getLaunchedSims(): Array<Spec> {
	return Object.keys(simLaunchStatuses)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => simLaunchStatuses[spec].status > LaunchStatus.Unlaunched);
}

export function getLaunchedSimsForClass(klass: Class): Array<Spec> {
	return Object.keys(specToClass)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => specToClass[spec] == klass && isSimLaunched(spec));
}

export function isSimLaunched(specIndex: Spec): boolean {
	return simLaunchStatuses[specIndex].status > LaunchStatus.Unlaunched;
}
