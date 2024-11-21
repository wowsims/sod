import { PseudoStat, Stat } from "../proto/common";

export enum Phase {
	Phase1 = 1,
	Phase2,
	Phase3,
	Phase4,
	Phase5,
	Phase6,
	Phase7,
}

export const LEVEL_THRESHOLDS: Record<Phase, number> = {
	[Phase.Phase1]: 25,
	[Phase.Phase2]: 40,
	[Phase.Phase3]: 50,
	[Phase.Phase4]: 60,
	[Phase.Phase5]: 60,
	[Phase.Phase6]: 60,
	[Phase.Phase7]: 60,
};

export const CURRENT_PHASE = Phase.Phase3;

export const LEVEL_BRACKETS = [25, 40, 50, 60];

// Github pages serves our site under the /sod directory (because the repo name is wotlk)
export const REPO_NAME = 'sod';

// Get 'elemental_shaman', the pathname part after the repo name
const pathnameParts = window.location.pathname.split('/');
const repoPartIdx = pathnameParts.findIndex(part => part == REPO_NAME);
export const SPEC_DIRECTORY = repoPartIdx == -1 ? '' : pathnameParts[repoPartIdx + 1];

export const GLOBAL_DISPLAY_STATS = [
	Stat.StatHealth,
	Stat.StatFireResistance,
	Stat.StatFrostResistance,
	Stat.StatNatureResistance,
];

export const GLOBAL_DISPLAY_PSEUDO_STATS = [
	PseudoStat.TimewornBonus,
];

export const GLOBAL_EP_STATS = [
	Stat.StatFireResistance,
	Stat.StatFrostResistance,
	Stat.StatNatureResistance,
];

export enum SortDirection {
	ASC,
	DESC,
}
