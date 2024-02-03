import { IconEnumPicker } from '../components/icon_enum_picker.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { Player } from '../player.js';
import { Spec } from '../proto/common.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
} from '../proto/shaman.js';
import { ActionId } from '../proto_utils/action_id.js';
import { ShamanSpecs } from '../proto_utils/utils.js';
import { EventID, TypedEvent } from '../typed_event.js';

import { ContentBlock } from './content_block.js';
import { Input } from './input.js';

export function TotemsSection(parentElem: HTMLElement, simUI: IndividualSimUI<ShamanSpecs>): ContentBlock {
	let contentBlock = new ContentBlock(parentElem, 'totems-settings', {
		header: { title: 'Totems' }
	});

	let totemDropdownGroup = Input.newGroupContainer();
	totemDropdownGroup.classList.add('totem-dropdowns-container', 'icon-group');

	contentBlock.bodyElement.appendChild(totemDropdownGroup);

	new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'earth-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#ffdfba', value: EarthTotem.NoEarthTotem },
			StrengthOfEarthTotem,
			StoneskinTotem,
			TremorTotem,
		],
		equals: (a: EarthTotem, b: EarthTotem) => a == b,
		zeroValue: EarthTotem.NoEarthTotem,
		changedEvent: (player: Player<ShamanSpecs>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.levelChangeEmitter]),
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.earth || EarthTotem.NoEarthTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.earth = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'fire-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#ffb3ba', value: FireTotem.NoFireTotem },
			SearingTotem,
			// FireNovaTotem,
			// MagmaTotem,
		],
		equals: (a: FireTotem, b: FireTotem) => a == b,
		zeroValue: FireTotem.NoFireTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.fire || FireTotem.NoFireTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.fire = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'water-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#bae1ff', value: WaterTotem.NoWaterTotem },
			{ actionId: () => ActionId.fromSpellId(58774), value: WaterTotem.ManaSpringTotem },
			{ actionId: () => ActionId.fromSpellId(58757), value: WaterTotem.HealingStreamTotem },
		],
		equals: (a: WaterTotem, b: WaterTotem) => a == b,
		zeroValue: WaterTotem.NoWaterTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.water || WaterTotem.NoWaterTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.water = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'air-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#baffc9', value: AirTotem.NoAirTotem },
			{ actionId: () => ActionId.fromSpellId(8512), value: AirTotem.WindfuryTotem },
			{ actionId: () => ActionId.fromSpellId(8835), value: AirTotem.GraceOfAirTotem },
		],
		equals: (a: AirTotem, b: AirTotem) => a == b,
		zeroValue: AirTotem.NoAirTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.air || AirTotem.NoAirTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.air = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	return contentBlock;
}

///////////////////////////////////////////////////////////////////////////
//                                 Earth Totems
///////////////////////////////////////////////////////////////////////////

export const StoneskinTotem = {
	actionId: (player: Player<Spec>) => player.getMatchingSpellActionId([
		{ id: 8071, 	minLevel: 4, 	maxLevel: 13 	},
		{ id: 8154, 	minLevel: 14, maxLevel: 23 	},
		{ id: 8155, 	minLevel: 24, maxLevel: 33 	},
		{ id: 10406, 	minLevel: 34, maxLevel: 43 	},
		{ id: 10407, 	minLevel: 44, maxLevel: 53 	},
		{ id: 10408, 	minLevel: 54 								},
	]),
	value: EarthTotem.StoneskinTotem,
};

export const StrengthOfEarthTotem = {
	actionId: (player: Player<Spec>) => player.getMatchingSpellActionId([
		{ id: 8075, 	minLevel: 10, maxLevel: 23 	},
		{ id: 8160, 	minLevel: 24, maxLevel: 37 	},
		{ id: 8161, 	minLevel: 38, maxLevel: 51 	},
		{ id: 10442, 	minLevel: 52, maxLevel: 59 	},
		{ id: 25361, 	minLevel: 60 								},
	]),
	value: EarthTotem.StrengthOfEarthTotem,
};

export const TremorTotem = {
	actionId: (player: Player<Spec>) => player.getMatchingSpellActionId([
		{ id: 8143, minLevel: 18 },
	]),
	value: EarthTotem.TremorTotem,
};

///////////////////////////////////////////////////////////////////////////
//                                 Fire Totems
///////////////////////////////////////////////////////////////////////////

export const SearingTotem = {
	actionId: (player: Player<Spec>) => player.getMatchingSpellActionId([
		{ id: 3599, minLevel: 10, maxLevel: 19 },
		{ id: 6363, minLevel: 20, maxLevel: 29 },
		{ id: 6364, minLevel: 30, maxLevel: 39 },
		{ id: 6365, minLevel: 40, maxLevel: 49 },
		{ id: 10437, minLevel: 50, maxLevel: 59 },
		{ id: 10438, minLevel: 60 },
	]),
	value: FireTotem.SearingTotem,
};

export const MagmaTotem = {

};

export const FireNovaTotem = {
	
};

///////////////////////////////////////////////////////////////////////////
//                                 Water Totems
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                                 Air Totems
///////////////////////////////////////////////////////////////////////////
