import { IndividualSimUI } from '../../../individual_sim_ui';
import { Class, EquipmentSpec, Profession, Race, Spec } from '../../../proto/common';
import { nameToClass, nameToProfession, nameToRace } from '../../../proto_utils/names';
import { IndividualImporter } from './individual_importer';

export class IndividualAddonImporter<SpecType extends Spec> extends IndividualImporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Addon Import', allowFileUpload: true });

		this.descriptionElem.appendChild(
			<>
				<p>
					Import settings from the{' '}
					<a href="https://www.curseforge.com/wow/addons/wowsimsexporter" target="_blank">
						WoWSims Exporter In-Game Addon
					</a>
					.
				</p>
				<p>This feature imports gear, race, talents, and professions. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.</p>
				<p>To import, paste the output from the addon below and click, 'Import'.</p>
			</>,
		);
	}

	async onImport(data: string) {
		const importJson = JSON.parse(data);

		// Parse all the settings.
		const charClass = nameToClass((importJson['class'] as string) || '');
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class!');
		}

		const race = nameToRace((importJson['race'] as string) || '');
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race!');
		}

		const professions = (importJson['professions'] as Array<{ name: string; level: number }>).map(profData => nameToProfession(profData.name));
		professions.forEach((prof, i) => {
			if (prof == Profession.ProfessionUnknown) {
				throw new Error(`Could not parse profession '${importJson['professions'][i]}'`);
			}
		});

		const talentsStr = (importJson['talents'] as string) || '';

		const gearJson = importJson['gear'];
		delete gearJson.version;
		gearJson.items = (gearJson.items as Array<any>).filter(item => item != null);
		const equipmentSpec = EquipmentSpec.fromJson(gearJson);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, professions);
	}
}
