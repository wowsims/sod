import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { Spec } from '../../proto/common';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
import { IconEnumPicker } from '../icon_enum_picker';
import { buildIconInput } from '../icon_inputs.js';
import { IconPicker } from '../icon_picker';
import * as ConsumablesInputs from '../inputs/consumables';
import { relevantStatOptions } from '../inputs/stat_options';
import { MultiIconPicker } from '../multi_icon_picker';

export class ConsumesPicker extends Component {
	protected simUI: IndividualSimUI<Spec>;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, 'consumes-picker-root');
		this.simUI = simUI;

		this.simUI.sim.waitForInit().then(() => {
			this.buildPotionsPicker();
			this.buildFlaskPicker();
			this.buildWeaponImbuePicker();
			this.buildFoodPicker();
			this.buildPhysicalBuffPickers();
			this.buildDefensiveBuffPickers();
			this.buildSpellPowerBuffPickers();
			this.buildMiscConsumesPickers();
			this.buildEngPickers();
			this.buildEnchPicker();
			this.buildPetPicker();
		});
	}

	private buildPotionsPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Potions</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-potions"></div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const potionsElem = this.rootElem.querySelector<HTMLElement>('.consumes-potions')!;

		const potionsOptions = ConsumablesInputs.makePotionsInput(relevantStatOptions(ConsumablesInputs.POTIONS_CONFIG, this.simUI), 'Potions');
		const conjuredOptions = ConsumablesInputs.makeConjuredInput(relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI));

		const pickers = [
			buildIconInput(potionsElem, this.simUI.player, potionsOptions),
			buildIconInput(potionsElem, this.simUI.player, conjuredOptions),
			buildIconInput(potionsElem, this.simUI.player, ConsumablesInputs.MildlyIrradiatedRejuvPotion),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildFlaskPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Flasks</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-flasks"></div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const flasksElem = this.rootElem.querySelector<HTMLElement>('.consumes-flasks')!;

		const flasksOptions = ConsumablesInputs.makeFlasksInput(relevantStatOptions(ConsumablesInputs.FLASKS_CONFIG, this.simUI));

		const pickers = [buildIconInput(flasksElem, this.simUI.player, flasksOptions)];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildWeaponImbuePicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Weapon Imbues</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-weapon-imbues"></div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const imbuesElem = this.rootElem.querySelector<HTMLElement>('.consumes-weapon-imbues')!;

		const mhImbueOptions = ConsumablesInputs.makeMainHandImbuesInput(
			relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_MH_CONFIG, this.simUI),
			'Main-Hand',
		);
		const ohImbueOptions = ConsumablesInputs.makeOffHandImbuesInput(relevantStatOptions(ConsumablesInputs.WEAPON_IMBUES_OH_CONFIG, this.simUI), 'Off-Hand');

		const pickers = [buildIconInput(imbuesElem, this.simUI.player, mhImbueOptions), buildIconInput(imbuesElem, this.simUI.player, ohImbueOptions)];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.gearChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildFoodPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Food</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-food"></div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const foodsElem = this.rootElem.querySelector<HTMLElement>('.consumes-food')!;

		const foodOptions = ConsumablesInputs.makeFoodInput(relevantStatOptions(ConsumablesInputs.FOOD_CONFIG, this.simUI));
		const alcoholOptions = ConsumablesInputs.makeAlcoholInput(relevantStatOptions(ConsumablesInputs.ALCOHOL_CONFIG, this.simUI));

		const pickers = [
			buildIconInput(foodsElem, this.simUI.player, foodOptions),
			buildIconInput(foodsElem, this.simUI.player, ConsumablesInputs.DragonBreathChili),
			buildIconInput(foodsElem, this.simUI.player, alcoholOptions),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildPhysicalBuffPickers() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Physical</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-physical"></div>
			</div>
		`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const physicalConsumesElem = this.rootElem.querySelector<HTMLElement>('.consumes-physical')!;

		const apBuffOptions = ConsumablesInputs.makeAttackPowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.ATTACK_POWER_CONSUMES_CONFIG, this.simUI),
			'Attack Power',
		);
		const agiBuffOptions = ConsumablesInputs.makeAgilityConsumeInput(relevantStatOptions(ConsumablesInputs.AGILITY_CONSUMES_CONFIG, this.simUI), 'Agility');
		const strBuffOptions = ConsumablesInputs.makeStrengthConsumeInput(
			relevantStatOptions(ConsumablesInputs.STRENGTH_CONSUMES_CONFIG, this.simUI),
			'Strength',
		);

		const pickers = [
			buildIconInput(physicalConsumesElem, this.simUI.player, apBuffOptions),
			buildIconInput(physicalConsumesElem, this.simUI.player, agiBuffOptions),
			buildIconInput(physicalConsumesElem, this.simUI.player, strBuffOptions),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildDefensiveBuffPickers() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Defensive</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-defensive"></div>
			</div>
		`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const defensiveConsumesElem = this.rootElem.querySelector<HTMLElement>('.consumes-defensive')!;

		const healthBuffOptions = ConsumablesInputs.makeHealthConsumeInput(relevantStatOptions(ConsumablesInputs.HEALTH_CONSUMES_CONFIG, this.simUI), 'Health');

		const armorBuffOptions = ConsumablesInputs.makeArmorConsumeInput(relevantStatOptions(ConsumablesInputs.ARMOR_CONSUMES_CONFIG, this.simUI), 'Armor');

		const pickers = [
			buildIconInput(defensiveConsumesElem, this.simUI.player, healthBuffOptions),
			buildIconInput(defensiveConsumesElem, this.simUI.player, armorBuffOptions),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildSpellPowerBuffPickers() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Spells</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-spells"></div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const spellsCnsumesElem = this.rootElem.querySelector<HTMLElement>('.consumes-spells')!;

		const spBuffOptions = ConsumablesInputs.makeSpellPowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.SPELL_POWER_CONFIG, this.simUI),
			'Spell Damage',
		);
		const fireBuffOptions = ConsumablesInputs.makeFirePowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.FIRE_POWER_CONFIG, this.simUI),
			'Fire Damage',
		);
		const frostBuffOptions = ConsumablesInputs.makeFrostPowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.FROST_POWER_CONFIG, this.simUI),
			'Frost Damage',
		);
		const shadowBuffOptions = ConsumablesInputs.makeShadowPowerConsumeInput(
			relevantStatOptions(ConsumablesInputs.SHADOW_POWER_CONFIG, this.simUI),
			'Shadow Damage',
		);
		const mp5BuffOptions = ConsumablesInputs.makeMp5ConsumeInput(relevantStatOptions(ConsumablesInputs.MP5_CONFIG, this.simUI), 'Mana Regen');

		const pickers = [
			buildIconInput(spellsCnsumesElem, this.simUI.player, spBuffOptions),
			buildIconInput(spellsCnsumesElem, this.simUI.player, fireBuffOptions),
			buildIconInput(spellsCnsumesElem, this.simUI.player, frostBuffOptions),
			buildIconInput(spellsCnsumesElem, this.simUI.player, shadowBuffOptions),
			buildIconInput(spellsCnsumesElem, this.simUI.player, mp5BuffOptions),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildMiscConsumesPickers() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Miscellaneous</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-misc"></div>
			</div>
		`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const miscConsumesElem = this.rootElem.querySelector<HTMLElement>('.consumes-misc')!;

		const sealOfTheDawnOptions = ConsumablesInputs.makeSealOfTheDawnConsumesInput(
			relevantStatOptions(ConsumablesInputs.SEAL_OF_THE_DAWN_CONSUMES_CONFIG, this.simUI),
			'Seal of the Dawn Buffs',
		);

		const zanzaBuffOptions = ConsumablesInputs.makeZanzaBuffConsumesInput(
			relevantStatOptions(ConsumablesInputs.ZANZA_BUFF_CONSUMES_CONFIG, this.simUI),
			'Zanza Buffs',
		);
		const mageScrollOptions = ConsumablesInputs.makeMageScrollsInput(ConsumablesInputs.MAGE_SCROLL_CONSUMES_CONFIG, 'Mage Scrolls');
		const miscOffensiveConsumesOptions = relevantStatOptions(ConsumablesInputs.MISC_OFFENSIVE_CONSUMES_CONFIG, this.simUI);
		const miscDefensiveConsumesOptions = relevantStatOptions(ConsumablesInputs.MISC_DEFENSIVE_CONSUMES_CONFIG, this.simUI);

		const pickers = [
			buildIconInput(miscConsumesElem, this.simUI.player, sealOfTheDawnOptions),
			buildIconInput(miscConsumesElem, this.simUI.player, zanzaBuffOptions),
			buildIconInput(miscConsumesElem, this.simUI.player, mageScrollOptions),
			ConsumablesInputs.makeMiscOffensiveConsumesInput(miscConsumesElem, this.simUI.player, this.simUI, miscOffensiveConsumesOptions),
			ConsumablesInputs.makeMiscDefensiveConsumesInput(miscConsumesElem, this.simUI.player, this.simUI, miscDefensiveConsumesOptions),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildEngPickers() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Engineering</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-engi"></div>
			</div>
		`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const engiConsumesElem = this.rootElem.querySelector<HTMLElement>('.consumes-engi')!;

		const explosiveOptions = ConsumablesInputs.makeExplosivesInput(relevantStatOptions(ConsumablesInputs.EXPLOSIVES_CONFIG, this.simUI), 'Explosives');
		const sapperOptions = ConsumablesInputs.makeSappersInput(relevantStatOptions(ConsumablesInputs.SAPPER_CONFIG, this.simUI), 'Sappers');

		const pickers = [
			buildIconInput(engiConsumesElem, this.simUI.player, sapperOptions),
			buildIconInput(engiConsumesElem, this.simUI.player, explosiveOptions),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildEnchPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Enchanting</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-ench">
						</div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const enchConsumesElem = this.rootElem.querySelector<HTMLElement>('.consumes-ench')!;

		const enchantedSigilOptions = ConsumablesInputs.makeEncanthedSigilInput(
			relevantStatOptions(ConsumablesInputs.ENCHANTED_SIGIL_CONFIG, this.simUI),
			'Enchanted Sigils',
		);

		const pickers = [buildIconInput(enchConsumesElem, this.simUI.player, enchantedSigilOptions)];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter, this.simUI.player.professionChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private buildPetPicker() {
		if (!this.simUI.individualConfig.petConsumeInputs?.length) return;

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Pet</label>
				<div class="picker-group icon-group consumes-row-inputs consumes-pet"></div>
			</div>
		`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const petConsumesElem = this.rootElem.querySelector<HTMLElement>('.consumes-pet')!;

		// const miscPetConsumesOptions = relevantStatOptions(ConsumablesInputs.MISC_PET_CONSUMES, this.simUI);

		const pickers = [
			...this.simUI.individualConfig.petConsumeInputs.map(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput)),
			// ConsumablesInputs.makeMiscPetConsumesInput(petConsumesElem, this.simUI.player, this.simUI, miscPetConsumesOptions),
		];

		TypedEvent.onAny([this.simUI.player.levelChangeEmitter]).on(() => this.updateRow(row, pickers));
		this.updateRow(row, pickers);
	}

	private updateRow(rowElem: HTMLElement, pickers: (IconPicker<Player<Spec>, any> | IconEnumPicker<Player<Spec>, any> | MultiIconPicker<Player<Spec>>)[]) {
		if (!!pickers.find(p => p?.showWhen())) {
			rowElem.classList.remove('hide');
		} else {
			rowElem.classList.add('hide');
		}
	}
}
