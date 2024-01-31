import { Player } from "../../player";
import { EarthTotem } from "../../proto/shaman";
import { ActionId } from "../../proto_utils/action_id";
import { ShamanSpecs } from "../../proto_utils/utils";

import { IconEnumValueConfig } from "../icon_enum_picker";

export const StoneskinTotemInputs: IconEnumValueConfig<Player<ShamanSpecs>, EarthTotem>[] = [
	{
		actionId: ActionId.fromSpellId(8071),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 4 && player.getLevel() < 14,
	},
	{
		actionId: ActionId.fromSpellId(8154),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 14 && player.getLevel() < 24,
	},
	{
		actionId: ActionId.fromSpellId(8155),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 24 && player.getLevel() < 34,
	},
	{
		actionId: ActionId.fromSpellId(10406),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 34 && player.getLevel() < 44,
	},
	{
		actionId: ActionId.fromSpellId(10407),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 44 && player.getLevel() < 54,
	},
	{
		actionId: ActionId.fromSpellId(10408),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 54,
	},
];

export const StrengthOfEarthTotemInputs: IconEnumValueConfig<Player<ShamanSpecs>, EarthTotem>[] = [
	{
		actionId: ActionId.fromSpellId(8075),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 10 && player.getLevel() < 24,
	},
	{
		actionId: ActionId.fromSpellId(8160),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 24 && player.getLevel() < 34,
	},
	{
		actionId: ActionId.fromSpellId(8161),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 34 && player.getLevel() < 52,
	},
  {
		actionId: ActionId.fromSpellId(10442),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 52 && player.getLevel() < 60,
	},
  {
		actionId: ActionId.fromSpellId(25361),
		value: EarthTotem.StrengthOfEarthTotem,
		showWhen: (player) => player.getLevel() >= 60,
	},
];

export const TremorTotemInput: IconEnumValueConfig<Player<ShamanSpecs>, EarthTotem> = {
  actionId: ActionId.fromSpellId(8143), value: EarthTotem.TremorTotem
};
