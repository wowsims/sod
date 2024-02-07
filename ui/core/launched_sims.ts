import { Class, Spec } from './proto/common';
import { specToClass } from './proto_utils/utils';

// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!

export enum LaunchStatus {
	Unlaunched,
	Phase_1,
	Phase_2,
	Phase_3,
	Phase_4,
	Phase_5,
};

export const raidSimStatus: LaunchStatus = LaunchStatus.Unlaunched;

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, LaunchStatus> = {
	[Spec.SpecBalanceDruid]: LaunchStatus.Phase_1,
	[Spec.SpecFeralDruid]: LaunchStatus.Phase_1,
	[Spec.SpecFeralTankDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecRestorationDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecElementalShaman]: LaunchStatus.Phase_1,
	[Spec.SpecEnhancementShaman]: LaunchStatus.Phase_1,
	[Spec.SpecRestorationShaman]: LaunchStatus.Unlaunched,
	[Spec.SpecHunter]: LaunchStatus.Phase_1,
	[Spec.SpecMage]: LaunchStatus.Unlaunched,
	[Spec.SpecRogue]: LaunchStatus.Unlaunched,
	[Spec.SpecHolyPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecProtectionPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecRetributionPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecHealingPriest]: LaunchStatus.Unlaunched,
	[Spec.SpecShadowPriest]: LaunchStatus.Unlaunched,
	[Spec.SpecWarlock]: LaunchStatus.Phase_1,
	[Spec.SpecTankWarlock]: LaunchStatus.Phase_1,
	[Spec.SpecWarrior]: LaunchStatus.Phase_1,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Unlaunched,
};

export function getLaunchedSims(): Array<Spec> {
	return Object.keys(simLaunchStatuses)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => simLaunchStatuses[spec] > LaunchStatus.Unlaunched);
}

export function getLaunchedSimsForClass(klass: Class): Array<Spec> {
	return Object.keys(specToClass)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => specToClass[spec] == klass && isSimLaunched(spec));
}

export function isSimLaunched(specIndex: Spec): boolean {
	return simLaunchStatuses[specIndex] > LaunchStatus.Unlaunched;
}
