import { IndividualSimUI } from "../../individual_sim_ui";
import { Player } from "../../player";
import { Spec, Stat } from "../../proto/common";
import { TypedEvent } from "../../typed_event";

import { Component } from "../component";
import { IconEnumPicker } from "../icon_enum_picker";
import { buildIconInput } from "../icon_inputs.js";
import { IconPicker } from "../icon_picker";
import { relevantStatOptions } from "../inputs/stat_options";

import * as ConsumablesInputs from '../inputs/consumables';

export class ConsumesPicker extends Component {
	protected simUI: IndividualSimUI<Spec>;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, 'consumes-picker-root');
		this.simUI = simUI;

		this.buildPotionsPicker();
		this.buildFlaskPicker();
		this.buildWeaponImbuePicker();
		this.buildFoodPicker();
		this.buildPhysicalBuffPicker();
		this.buildSpellPowerBuffPicker();
		this.buildEngPicker();
		this.buildEnchPicker();
		this.buildPetPicker();
	}

	private buildPotionsPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Potions</label>
        <div class="consumes-row-inputs consumes-potions"></div>
      </div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const potionsElem = this.rootElem.querySelector('.consumes-potions') as HTMLElement;

		const potionsOptions = ConsumablesInputs.makePotionsInput(
			relevantStatOptions(ConsumablesInputs.POTIONS_CONFIG, this.simUI),
			'Combat Potion',
		)
		const potionsPicker = buildIconInput(potionsElem, this.simUI.player, potionsOptions);

		const conjuredOptions = ConsumablesInputs.makeConjuredInput(
			relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI),
		);
		const conjuredPicker = buildIconInput(potionsElem, this.simUI.player, conjuredOptions);
		const irradiatedRejuvPicker = buildIconInput(potionsElem, this.simUI.player, ConsumablesInputs.MildlyIrradiatedRejuvPotion);

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter]).on(() => {
			this.updateRow(row, [potionsPicker, conjuredPicker, irradiatedRejuvPicker]);
		});
	}

	private buildFlaskPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Elixirs</label>
        <div class="consumes-row-inputs consumes-flasks"></div>
      </div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const flasksElem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;

		const flasksOptions = ConsumablesInputs.makeFlasksInput(
			relevantStatOptions(ConsumablesInputs.FLASKS_CONFIG, this.simUI)
		);
		const flasksPicker = buildIconInput(flasksElem, this.simUI.player, flasksOptions);
		
		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, [flasksPicker]));
		this.updateRow(row, [flasksPicker]);
	}

	private buildWeaponImbuePicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
    	<div class="consumes-row input-root input-inline">
        <label class="form-label">Weapon Imbues</label>
        <div class="consumes-row-inputs consumes-weapon-imbues"></div>
    	</div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const imbuesElem = this.rootElem.querySelector('.consumes-weapon-imbues') as HTMLElement;

		const mhImbueOptions = ConsumablesInputs.makeMainHandImbuesInput(
			relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_MH_CONFIG, this.simUI),
			'Main-Hand',
		);
		const mhPicker = buildIconInput(imbuesElem, this.simUI.player, mhImbueOptions);

		const ohImbueOptions = ConsumablesInputs.makeOffHandImbuesInput(
			relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_OH_CONFIG, this.simUI),
			'Off-Hand',
		);
		const ohPicker = buildIconInput(imbuesElem, this.simUI.player, ohImbueOptions);

		TypedEvent.onAny([this.simUI.player.gearChangeEmitter]).on(() => this.updateRow(row, [mhPicker, ohPicker]));
		this.updateRow(row, [mhPicker, ohPicker]);
	}

	private buildFoodPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Food</label>
        <div class="consumes-row-inputs consumes-food"></div>
      </div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const foodsElem = this.rootElem.querySelector('.consumes-food') as HTMLElement;

		const foodOptions = ConsumablesInputs.makeFoodInput(
			relevantStatOptions(ConsumablesInputs.FOOD_CONFIG, this.simUI),
		);
		const foodPicker = buildIconInput(foodsElem, this.simUI.player, foodOptions);

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, [foodPicker]));
		this.updateRow(row, [foodPicker]);
	}

	private buildPhysicalBuffPicker() {
		const includeAgi = this.simUI.individualConfig.epStats.includes(Stat.StatAgility)
		const includeStr = this.simUI.individualConfig.epStats.includes(Stat.StatStrength)

		if (!includeAgi && !includeStr) return;

		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Physical</label>
        <div class="consumes-row-inputs consumes-physical"></div>
      </div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const physicalConsumesElem = this.rootElem.querySelector('.consumes-physical') as HTMLElement;

		const agiBuffOptions = ConsumablesInputs.makeAgilityConsumeInput(
			relevantStatOptions(ConsumablesInputs.AGILITY_CONSUMES_CONFIG, this.simUI),
			'Agility',
		);
		const agiBuffPicker = buildIconInput(physicalConsumesElem, this.simUI.player, agiBuffOptions)

		const strBuffOptions = ConsumablesInputs.makeStrengthConsumeInput(
			relevantStatOptions(ConsumablesInputs.STRENGTH_CONSUMES_CONFIG, this.simUI),
			'Strength',
		);
		const strBuffPicker = buildIconInput(physicalConsumesElem, this.simUI.player, strBuffOptions);

		buildIconInput(physicalConsumesElem, this.simUI.player, ConsumablesInputs.BoglingRootBuff);

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, [strBuffPicker, agiBuffPicker]));
		this.updateRow(row, [strBuffPicker, agiBuffPicker]);
	}

	private buildSpellPowerBuffPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Spells</label>
        <div class="consumes-row-inputs consumes-spells"></div>
      </div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const spellsCnsumesElem = this.rootElem.querySelector('.consumes-spells') as HTMLElement;

		const spBuffOptions = ConsumablesInputs.makeSpellPowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.SPELL_POWER_CONFIG, this.simUI),
			'Arcane',
		);
		const spBuffPicker =buildIconInput(spellsCnsumesElem, this.simUI.player, spBuffOptions);

		const fireBuffOptions = ConsumablesInputs.makeFirePowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.FIRE_POWER_CONFIG, this.simUI),
			'Fire',
		);
		const fireBuffPicker =buildIconInput(spellsCnsumesElem, this.simUI.player, fireBuffOptions);

		const frostBuffOptions = ConsumablesInputs.makeFrostPowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.FROST_POWER_CONFIG, this.simUI),
			'Frost',
		);
		const frostBuffPicker =buildIconInput(spellsCnsumesElem, this.simUI.player, frostBuffOptions);

		const shadowBuffOptions = ConsumablesInputs.makeshadowPowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.SHADOW_POWER_CONFIG, this.simUI),
			'Shadow',
		);
		const shadowBuffPicker =buildIconInput(spellsCnsumesElem, this.simUI.player, shadowBuffOptions);

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(
			() => this.updateRow(row, [spBuffPicker, fireBuffPicker, frostBuffPicker, shadowBuffPicker])
		);
		this.updateRow(row, [spBuffPicker, fireBuffPicker, frostBuffPicker, shadowBuffPicker]);
	}

	private buildEngPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Engineering</label>
        <div class="consumes-row-inputs consumes-engi">
				</div>
      </div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const engiConsumesElem = this.rootElem.querySelector('.consumes-engi') as HTMLElement;
		
		const sapperPicker = buildIconInput(engiConsumesElem, this.simUI.player, ConsumablesInputs.Sapper);

		const explosiveOptions = ConsumablesInputs.makeExplosivesInput(
			relevantStatOptions(ConsumablesInputs.EXPLOSIVES_CONFIG, this.simUI),
			'Explosives',
		);
		const explosivePicker = buildIconInput(engiConsumesElem, this.simUI.player, explosiveOptions)

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter]).on(
			() => this.updateRow(row, [sapperPicker, explosivePicker])
		);
		this.updateRow(row, [sapperPicker , explosivePicker]);
	}

	private buildEnchPicker() {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <div class="consumes-row input-root input-inline">
        <label class="form-label">Enchanting</label>
        <div class="consumes-row-inputs consumes-ench">
				</div>
      </div>
    `;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const enchConsumesElem = this.rootElem.querySelector('.consumes-ench') as HTMLElement;

		const enchantedSigilOptions = ConsumablesInputs.makeEncanthedSigilInput(
			relevantStatOptions(ConsumablesInputs.ENCHANTEDSIGILCONFIG, this.simUI),
			'Enchanted Sigils',
		);
		const enchantedSigilpicker = buildIconInput(enchConsumesElem, this.simUI.player, enchantedSigilOptions)

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter]).on(
			() => this.updateRow(row, [enchantedSigilpicker])
		);
		this.updateRow(row, [enchantedSigilpicker]);
	}

	private buildPetPicker() {
		if (!this.simUI.individualConfig.petConsumeInputs?.length) return

		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Pet</label>
				<div class="consumes-row-inputs consumes-pet"></div>
			</div>
		`;

		this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const petConsumesElem = this.rootElem.querySelector('.consumes-pet') as HTMLElement;

		this.simUI.individualConfig.petConsumeInputs.map(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput));
	}

	private updateRow(rowElem: HTMLElement, pickers: (IconPicker<Player<Spec>, any> | IconEnumPicker<Player<Spec>, any>)[]) {
		if (!!pickers.find(p => p?.showWhen())) {
			rowElem.classList.remove('hide');
	 } else {
		 rowElem.classList.add('hide');
	 }
	}
}
