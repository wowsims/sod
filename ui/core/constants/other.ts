export enum Phase {
  Phase1 = 1,
  Phase2,
  Phase3,
  Phase4,
  Phase5,
};

export const LEVEL_THRESHOLDS: Record<Phase, number> = {
  [Phase.Phase1]: 25,
  [Phase.Phase2]: 40,
  [Phase.Phase3]: 50,
  [Phase.Phase4]: 60,
  [Phase.Phase5]: 60,
};

export const CURRENT_PHASE = Phase.Phase2;

// Github pages serves our site under the /sod directory (because the repo name is wotlk)
export const REPO_NAME = 'sod';

// Get 'elemental_shaman', the pathname part after the repo name
const pathnameParts = window.location.pathname.split('/');
const repoPartIdx = pathnameParts.findIndex(part => part == REPO_NAME);
export const SPEC_DIRECTORY = repoPartIdx == -1 ? '' : pathnameParts[repoPartIdx + 1];
