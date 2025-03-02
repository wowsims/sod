import { IndividualSimUI } from '../../../individual_sim_ui';
import { Class, EquipmentSpec, ItemSlot, ItemSpec, Race, Spec } from '../../../proto/common';
import { nameToClass, nameToRace } from '../../../proto_utils/names';
import { buf2hex } from '../../../utils';
import { IndividualImporter } from './individual_importer';

export class IndividualWowheadGearPlannerImporter<SpecType extends Spec> extends IndividualImporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Wowhead Import', allowFileUpload: true });

		this.descriptionElem.appendChild(
			<>
				<p>
					Import settings from{' '}
					<a href="https://www.wowhead.com/classic/gear-planner" target="_blank">
						Wowhead Gear Planner
					</a>
					.
				</p>
				<p>This feature imports gear, race, and (optionally) talents. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.</p>
				<p>To import, paste the gear planner link below and click, 'Import'.</p>
			</>,
		);
	}

	async onImport(url: string) {
		const match = url.match(/www\.wowhead\.com\/classic\/gear-planner\/([a-z\-]+)\/([a-z\-]+)\/([a-zA-Z0-9_\-]+)/);
		if (!match) {
			throw new Error(`Invalid WCL URL ${url}, must look like "https://www.wowhead.com/classic/gear-planner/CLASS/RACE/XXXX"`);
		}

		// Parse all the settings.
		const charClass = nameToClass(match[1].replaceAll('-', ''));
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class: ' + match[1]);
		}

		const race = nameToRace(match[2].replaceAll('-', ''));
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race: ' + match[2]);
		}

		const base64Data = match[3].replaceAll('_', '/').replaceAll('-', '+');
		//console.log('Base64: ' + base64Data);
		const data = Uint8Array.from(atob(base64Data), c => c.charCodeAt(0));
		//console.log('Hex: ' + buf2hex(data));

		// Binary schema
		// Byte 00: ??
		// Byte 01: ?? Seems related to aesthetics (e.g. body type)
		// Byte 02: 8-bit Player Level
		// Byte 03: 8-bit length of talents bytes
		// Next N Bytes: Talents in hex string format

		// Talent hex string looks like '230005232100330150323102505321f03f023203001f'
		// Just like regular wowhead talents string except 'f' instead of '-'.
		const numTalentBytes = data[3];
		const talentBytes = data.subarray(4, 4 + numTalentBytes);
		const talentsHexStr = buf2hex(talentBytes);
		//console.log('Talents hex: ' + talentsHexStr);
		const talentsStr = talentsHexStr.split('f').slice(0, 3).join('-');
		//console.log('Talents: ' + talentsStr);

		// Binary schema for each item:
		// 8-bit slotNumber, high bit = is enchanted
		// 8-bit upper 3 bits for gem count
		// 16-bit item id
		// if enchant bit is set:
		//   8-bit ??, possibly enchant position for multiple enchants?
		//   16-bit enchant id
		const gearBytes = data.subarray(numTalentBytes);
		//console.log(`Remaining ${gearBytes.length} bytes: ${buf2hex(gearBytes)}`);

		const equipmentSpec = EquipmentSpec.create();
		let cur = 0;
		while (cur < gearBytes.length) {
			const itemSpec = ItemSpec.create();
			const slotId = gearBytes[cur] & 0b00111111;
			const isEnchanted = Boolean(gearBytes[cur] & 0b10000000);
			// const randomEnchant = Boolean(gearBytes[cur] & 0b01000000);
			cur++;

			const highid = gearBytes[cur] & 0b00011111;
			cur++;

			itemSpec.id = (highid << 16) + (gearBytes[cur] << 8) + gearBytes[cur + 1];
			cur += 2;
			//console.log(`Slot ID: ${slotId}, isEnchanted: ${isEnchanted}, itemID: ${itemSpec.id}`);

			if (isEnchanted) {
				// Note: this is the enchant SPELL id, not the effect ID.
				const enchantSpellId = (gearBytes[cur] << 16) + (gearBytes[cur + 1] << 8) + gearBytes[cur + 2];
				itemSpec.enchant = this.simUI.sim.db.enchantSpellIdToEffectId(enchantSpellId);
				cur += 3;
				//console.log(`Enchant ID: ${itemSpec.enchant}. Spellid: ${enchantSpellId}`);
			}

			// Ignore tabard / shirt slots
			const itemSlotEntry = Object.entries(IndividualWowheadGearPlannerImporter.slotIDs).find(e => e[1] == slotId);
			if (itemSlotEntry != null) {
				equipmentSpec.items.push(itemSpec);
			}
		}
		this.simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, []);
	}

	static slotIDs: Record<ItemSlot, number> = {
		[ItemSlot.ItemSlotHead]: 1,
		[ItemSlot.ItemSlotNeck]: 2,
		[ItemSlot.ItemSlotShoulder]: 3,
		[ItemSlot.ItemSlotBack]: 15,
		[ItemSlot.ItemSlotChest]: 5,
		[ItemSlot.ItemSlotWrist]: 9,
		[ItemSlot.ItemSlotHands]: 10,
		[ItemSlot.ItemSlotWaist]: 6,
		[ItemSlot.ItemSlotLegs]: 7,
		[ItemSlot.ItemSlotFeet]: 8,
		[ItemSlot.ItemSlotFinger1]: 11,
		[ItemSlot.ItemSlotFinger2]: 12,
		[ItemSlot.ItemSlotTrinket1]: 13,
		[ItemSlot.ItemSlotTrinket2]: 14,
		[ItemSlot.ItemSlotMainHand]: 16,
		[ItemSlot.ItemSlotOffHand]: 17,
		[ItemSlot.ItemSlotRanged]: 18,
	};
}
