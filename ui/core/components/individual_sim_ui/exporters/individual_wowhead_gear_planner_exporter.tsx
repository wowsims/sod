import { MAX_CHARACTER_LEVEL } from '../../../constants/mechanics';
import { IndividualSimUI } from '../../../individual_sim_ui';
import { Spec } from '../../../proto/common';
import { classNames, raceNames } from '../../../proto_utils/names';
import { IndividualWowheadGearPlannerImporter } from '../importers/individual_wowhead_gear_planner_importer';
import { IndividualExporter } from './individual_exporter';

export class IndividualWowheadGearPlannerExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Wowhead Export', allowDownload: true });
	}

	getData(): string {
		const player = this.simUI.player;

		const classStr = classNames.get(player.getClass())!.replaceAll(' ', '-').toLowerCase();
		const raceStr = raceNames.get(player.getRace())!.replaceAll(' ', '-').toLowerCase();
		const url = `https://www.wowhead.com/classic/gear-planner/${classStr}/${raceStr}/`;

		// See comments on the importer for how the binary formatting is structured.
		let bytes: Array<number> = [];
		bytes.push(6);
		bytes.push(0);
		bytes.push(MAX_CHARACTER_LEVEL);

		let talentsStr = player.getTalentsString().replaceAll('-', 'f') + 'f';
		if (talentsStr.length % 2 == 1) {
			talentsStr += '0';
		}
		//console.log('Talents str: ' + talentsStr);
		bytes.push(talentsStr.length / 2);
		for (let i = 0; i < talentsStr.length; i += 2) {
			bytes.push(parseInt(talentsStr.substring(i, i + 2), 16));
		}

		const to2Bytes = (val: number): Array<number> => {
			return [(val & 0xff00) >> 8, val & 0x00ff];
		};

		const gear = player.getGear();
		gear.getItemSlots()
			.sort((slot1, slot2) => IndividualWowheadGearPlannerImporter.slotIDs[slot1] - IndividualWowheadGearPlannerImporter.slotIDs[slot2])
			.forEach(itemSlot => {
				const item = gear.getEquippedItem(itemSlot);
				if (!item) {
					return;
				}

				let slotId = IndividualWowheadGearPlannerImporter.slotIDs[itemSlot];
				if (item.enchant) {
					slotId = slotId | 0b10000000;
				}
				bytes.push(slotId);
				bytes.push(0 << 5);
				bytes = bytes.concat(to2Bytes(item.item.id));

				if (item.enchant) {
					bytes.push(0);
					bytes = bytes.concat(to2Bytes(item.enchant.spellId));
				}
			});

		//console.log('Hex: ' + buf2hex(new Uint8Array(bytes)));
		const binaryString = String.fromCharCode(...bytes);
		const b64encoded = btoa(binaryString);
		const b64converted = b64encoded.replaceAll('/', '_').replaceAll('+', '-').replace(/=+$/, '');

		return url + b64converted;
	}
}
