import clsx from 'clsx';

import { getLaunchedSimsForClass, LaunchStatus, raidSimStatus, simLaunchStatuses } from '../launched_sims.js';
import { Class, Spec } from '../proto/common.js';
import {
	classIcons,
	classNames,
	getSpecSiteUrl,
	naturalClassOrder,
	raidSimIcon,
	raidSimLabel,
	raidSimSiteUrl,
	specNames,
	specToClass,
	textCssClassForClass,
	textCssClassForSpec,
	titleIcons,
} from '../proto_utils/utils.js';
import { Component } from './component.js';

interface ClassOptions {
	type: 'Class';
	index: Class;
}

interface SpecOptions {
	type: 'Spec';
	index: Spec;
}

interface RaidOptions {
	type: 'Raid';
}

type SimTitleDropdownConfig = {
	noDropdown?: boolean;
};

// Dropdown menu for selecting a player.
export class SimTitleDropdown extends Component {
	private readonly dropdownMenu: HTMLElement | undefined;

	private readonly specLabels: Record<Spec, string> = {
		[Spec.SpecBalanceDruid]: 'Balance',
		[Spec.SpecFeralDruid]: 'Feral DPS',
		[Spec.SpecFeralTankDruid]: 'Feral Tank',
		[Spec.SpecRestorationDruid]: 'Restoration',
		[Spec.SpecElementalShaman]: 'Elemental',
		[Spec.SpecEnhancementShaman]: 'Enhancement',
		[Spec.SpecRestorationShaman]: 'Restoration',
		[Spec.SpecHunter]: 'Hunter',
		[Spec.SpecMage]: 'Mage',
		[Spec.SpecRogue]: 'DPS',
		[Spec.SpecTankRogue]: 'Tank',
		[Spec.SpecHolyPaladin]: 'Holy',
		[Spec.SpecProtectionPaladin]: 'Protection',
		[Spec.SpecRetributionPaladin]: 'Retribution',
		[Spec.SpecHealingPriest]: 'Healing',
		[Spec.SpecShadowPriest]: 'Shadow',
		[Spec.SpecWarlock]: 'DPS',
		[Spec.SpecTankWarlock]: 'Tank',
		[Spec.SpecWarrior]: 'DPS',
		[Spec.SpecTankWarrior]: 'Tank',
	};

	constructor(parent: HTMLElement, currentSpecIndex: Spec | null, config: SimTitleDropdownConfig = {}) {
		super(parent, 'sim-title-dropdown-root');

		const rootLinkArgs: SpecOptions | RaidOptions = currentSpecIndex === null ? { type: 'Raid' } : { type: 'Spec', index: currentSpecIndex };
		const rootLink = this.buildRootSimLink(rootLinkArgs);

		if (config.noDropdown) {
			this.rootElem.innerHTML = rootLink.outerHTML;
			return;
		}

		this.rootElem.innerHTML = `
      <div class="dropdown sim-link-dropdown">
        ${rootLink.outerHTML}
        <ul class="dropdown-menu"></ul>
      </div>
    `;

		this.dropdownMenu = this.rootElem.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
		this.buildDropdown();

		// Prevent Bootstrap from closing the menu instead of opening class menus
		this.dropdownMenu.addEventListener('click', event => {
			const target = event.target as HTMLElement;
			const link = target.closest('a:not([href="javascript:void(0)"]');

			if (!link) {
				event.stopPropagation();
				event.preventDefault();
			}
		});
	}

	private buildDropdown() {
		// TODO Classic
		// if (raidSimStatus >= LaunchStatus.Alpha) {
		// 	// Add the raid sim to the top of the dropdown
		// 	let raidListItem = document.createElement('li');
		// 	raidListItem.appendChild(this.buildRaidLink());
		// 	this.dropdownMenu?.appendChild(raidListItem);
		// }

		naturalClassOrder.forEach(classIndex => {
			const listItem = document.createElement('li');
			const sims = getLaunchedSimsForClass(classIndex);

			if (sims.length == 1) {
				// The class only has one listed sim so make a direct link to the sim
				listItem.appendChild(this.buildClassLink(classIndex));
				this.dropdownMenu?.appendChild(listItem);
			} else if (sims.length > 1) {
				// Add the class to the dropdown with an additional spec dropdown
				listItem.appendChild(this.buildClassDropdown(classIndex));
				this.dropdownMenu?.appendChild(listItem);
			}
		});
	}

	private buildClassDropdown(classIndex: Class) {
		const sims = getLaunchedSimsForClass(classIndex);
		const dropdownFragment = document.createElement('fragment');
		const dropdownMenu = document.createElement('ul');
		dropdownMenu.classList.add('dropdown-menu');

		// Generate the class link to act as a dropdown toggle for the spec dropdown
		const classLink = this.buildClassLink(classIndex);

		// Generate links for a class's specs
		sims.forEach(specIndex => {
			const listItem = document.createElement('li');
			const link = this.buildSpecLink(specIndex);

			listItem.appendChild(link);
			dropdownMenu.appendChild(listItem);
		});

		dropdownFragment.innerHTML = `
			<div class="dropend sim-link-dropdown">
				${classLink.outerHTML}
				${dropdownMenu.outerHTML}
			</div>
    	`;

		return dropdownFragment.children[0] as HTMLElement;
	}

	private buildRootSimLink(data: SpecOptions | RaidOptions): Element {
		let label;

		if (data.type == 'Raid') label = raidSimLabel;
		else {
			const classIndex = specToClass[data.index];
			if (getLaunchedSimsForClass(classIndex).length > 1)
				// If the class has multiple sims, use the spec name
				label = specNames[data.index];
			// If the class has only 1 sim, use the class name
			else label = classNames[classIndex];
		}

		return (
			<a href="javascript:void(0)" className={clsx('sim-link', this.getContextualKlass(data))} dataset={{ bsToggle: 'dropdown', bsTrigger: 'click' }}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath(data)} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-label text-white">WoWSims - Season of Discovery</span>
						<span className="sim-link-title">{label}</span>
						{this.launchStatusLabel(data)}
					</div>
				</div>
			</a>
		);
	}

	private buildRaidLink(): HTMLElement {
		const textKlass = this.getContextualKlass({ type: 'Raid' });
		const iconPath = this.getSimIconPath({ type: 'Raid' });
		const label = raidSimLabel;

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
      <a href="${raidSimSiteUrl}" class="sim-link ${textKlass}">
        <div class="sim-link-content">
          <img src="${iconPath}" class="sim-link-icon">
          <div class="d-flex flex-column">
            <span class="sim-link-title">${label}</span>
            ${this.launchStatusLabel({ type: 'Raid' })}
          </div>
        </div>
      </a>
    `;

		return fragment.children[0] as HTMLElement;
	}

	private buildClassLink(classIndex: Class): Element {
		const specIndexes = getLaunchedSimsForClass(classIndex);
		const hasSpecSims = specIndexes.length > 1;
		const href = hasSpecSims ? 'javascript:void(0)' : getSpecSiteUrl(specIndexes[0]);

		return (
			<a
				href={href}
				className={clsx('sim-link', this.getContextualKlass({ type: 'Class', index: classIndex }))}
				dataset={hasSpecSims ? { bsToggle: 'dropdown' } : {}}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Class', index: classIndex })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-title">{classNames[classIndex]}</span>
						{!hasSpecSims && this.launchStatusLabel({ type: 'Spec', index: specIndexes[0] })}
					</div>
				</div>
			</a>
		);
	}

	private buildSpecLink(specIndex: Spec): Element {
		const href = getSpecSiteUrl(specIndex);

		return (
			<a href={href} className={clsx('sim-link', this.getContextualKlass({ type: 'Spec', index: specIndex }))}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Spec', index: specIndex })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-label">{classNames[specToClass[specIndex]]}</span>
						<span className="sim-link-title">{this.specLabels[specIndex]}</span>
						{this.launchStatusLabel({ type: 'Spec', index: specIndex })}
					</div>
				</div>
			</a>
		);
	}

	private launchStatusLabel(data: SpecOptions | RaidOptions): Element {
		const status = data.type == 'Raid' ? raidSimStatus.status : simLaunchStatuses[data.index].status;
		const phase = data.type == 'Raid' ? raidSimStatus.phase : simLaunchStatuses[data.index].phase;

		return (
			<span className="launch-status-label text-brand">
				{status === LaunchStatus.Unlaunched ? (
					<>Not Yet Supported</>
				) : (
					<>
						Phase {phase}
						{status != LaunchStatus.Launched && <> - {LaunchStatus[status]}</>}
					</>
				)}
			</span>
		);
	}

	private getSimIconPath(data: ClassOptions | SpecOptions | RaidOptions): string {
		let iconPath: string;

		if (data.type == 'Raid') {
			iconPath = raidSimIcon;
		} else if (data.type == 'Class') {
			iconPath = classIcons[data.index];
		} else {
			iconPath = titleIcons[data.index];
		}

		return iconPath;
	}

	private getContextualKlass(data: ClassOptions | SpecOptions | RaidOptions): string {
		if (data.type == 'Raid')
			// Raid link
			return 'text-white';
		else if (data.type == 'Class')
			// Class links
			return textCssClassForClass(data.index);
		else return textCssClassForSpec(data.index);
	}
}
