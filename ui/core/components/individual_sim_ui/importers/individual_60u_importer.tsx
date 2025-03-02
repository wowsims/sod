import { IndividualSimUI } from '../../../individual_sim_ui';
import { Class, EquipmentSpec, ItemSpec, Race, Spec } from '../../../proto/common';
import { nameToClass, nameToRace } from '../../../proto_utils/names';
import { talentSpellIdsToTalentString } from '../../../talents/factory';
import { IndividualImporter } from './individual_importer';

export class Individual60UImporter<SpecType extends Spec> extends IndividualImporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Sixty Upgrades SoD Import', allowFileUpload: true });

		this.descriptionElem.innerHTML = `
            <p>
                Import settings from <a href="https://sixtyupgrades.com/sod/" target="_blank">Sixty Upgrades Season of Discovery</a>.
            </p>
            <p>
                This feature imports gear, race, and (optionally) talents. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.
            </p>
            <p>
                To import, paste the output from the site's export option below and click, 'Import'.
            </p>
        `;
	}

	onImport(data: string) {
		const importJson = JSON.parse(data);

		// Parse all the settings.
		const charClass = nameToClass((importJson?.character?.gameClass as string) || '');
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class!');
		}

		const race = nameToRace((importJson?.character?.race as string) || '');
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race!');
		}

		let talentsStr = '';
		if (importJson?.talents?.length > 0) {
			const talentIds = (importJson.talents as Array<any>).map(talentJson => talentJson.spellId);
			talentsStr = talentSpellIdsToTalentString(charClass, talentIds);
		}

		const equipmentSpec = EquipmentSpec.create();
		(importJson.items as Array<any>).forEach(itemJson => {
			const itemSpec = ItemSpec.create();
			itemSpec.id = itemJson.id;
			if (itemJson.enchant?.id) {
				itemSpec.enchant = itemJson.enchant.id;
			}
			equipmentSpec.items.push(itemSpec);
		});

		this.simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, []);
	}
}
