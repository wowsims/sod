import { Player } from '../../player.js';
import { Spec } from '../../proto/common.js';
import { AirTotem, EarthTotem, FireTotem, WaterTotem } from '../../proto/shaman.js';

///////////////////////////////////////////////////////////////////////////
//                                 Earth Totems
///////////////////////////////////////////////////////////////////////////

export const StoneskinTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 8071, minLevel: 4, maxLevel: 13 },
			{ id: 8154, minLevel: 14, maxLevel: 23 },
			{ id: 8155, minLevel: 24, maxLevel: 33 },
			{ id: 10406, minLevel: 34, maxLevel: 43 },
			{ id: 10407, minLevel: 44, maxLevel: 53 },
			{ id: 10408, minLevel: 54 },
		]),
	value: EarthTotem.StoneskinTotem,
};

export const StrengthOfEarthTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 8075, minLevel: 10, maxLevel: 23 },
			{ id: 8160, minLevel: 24, maxLevel: 37 },
			{ id: 8161, minLevel: 38, maxLevel: 51 },
			{ id: 10442, minLevel: 52, maxLevel: 59 },
			{ id: 25361, minLevel: 60 },
		]),
	value: EarthTotem.StrengthOfEarthTotem,
};

export const TremorTotem = {
	actionId: (player: Player<Spec>) => player.getMatchingSpellActionId([{ id: 8143, minLevel: 18 }]),
	value: EarthTotem.TremorTotem,
};

///////////////////////////////////////////////////////////////////////////
//                                 Fire Totems
///////////////////////////////////////////////////////////////////////////

export const SearingTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 3599, minLevel: 10, maxLevel: 19 },
			{ id: 6363, minLevel: 20, maxLevel: 29 },
			{ id: 6364, minLevel: 30, maxLevel: 39 },
			{ id: 6365, minLevel: 40, maxLevel: 49 },
			{ id: 10437, minLevel: 50, maxLevel: 59 },
			{ id: 10438, minLevel: 60 },
		]),
	value: FireTotem.SearingTotem,
};

export const FireNovaTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 1535, minLevel: 12, maxLevel: 21 },
			{ id: 8498, minLevel: 22, maxLevel: 31 },
			{ id: 8499, minLevel: 32, maxLevel: 41 },
			{ id: 11314, minLevel: 42, maxLevel: 51 },
			{ id: 11315, minLevel: 52 },
		]),
	value: FireTotem.FireNovaTotem,
};

export const MagmaTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 8190, minLevel: 26, maxLevel: 35 },
			{ id: 10585, minLevel: 36, maxLevel: 45 },
			{ id: 10586, minLevel: 46, maxLevel: 55 },
			{ id: 10587, minLevel: 56 },
		]),
	value: FireTotem.FireNovaTotem,
};

///////////////////////////////////////////////////////////////////////////
//                                 Water Totems
///////////////////////////////////////////////////////////////////////////

export const HealingStreamTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 5394, minLevel: 20, maxLevel: 29 },
			{ id: 6375, minLevel: 30, maxLevel: 39 },
			{ id: 6377, minLevel: 40, maxLevel: 49 },
			{ id: 10462, minLevel: 50, maxLevel: 59 },
			{ id: 10463, minLevel: 60 },
		]),
	value: WaterTotem.HealingStreamTotem,
};

export const ManaSpringTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 5675, minLevel: 26, maxLevel: 35 },
			{ id: 10495, minLevel: 36, maxLevel: 45 },
			{ id: 10496, minLevel: 46, maxLevel: 55 },
			{ id: 10497, minLevel: 56 },
		]),
	value: WaterTotem.ManaSpringTotem,
};

///////////////////////////////////////////////////////////////////////////
//                                 Air Totems
///////////////////////////////////////////////////////////////////////////

export const WindfuryTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 8512, minLevel: 32, maxLevel: 41 },
			{ id: 10613, minLevel: 42, maxLevel: 51 },
			{ id: 25359, minLevel: 52 },
		]),
	value: AirTotem.WindfuryTotem,
};

export const GraceOfAirTotem = {
	actionId: (player: Player<Spec>) =>
		player.getMatchingSpellActionId([
			{ id: 10627, minLevel: 42, maxLevel: 55 },
			{ id: 10627, minLevel: 56, maxLevel: 59 },
			{ id: 25359, minLevel: 60 },
		]),
	value: AirTotem.GraceOfAirTotem,
};
