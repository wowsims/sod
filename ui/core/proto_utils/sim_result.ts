import {
	ActionMetrics as ActionMetricsProto,
	AuraMetrics as AuraMetricsProto,
	DistributionMetrics as DistributionMetricsProto,
	EncounterMetrics as EncounterMetricsProto,
	Party as PartyProto,
	PartyMetrics as PartyMetricsProto,
	Player as PlayerProto,
	Raid as RaidProto,
	RaidMetrics as RaidMetricsProto,
	RaidSimRequest,
	RaidSimResult,
	ResourceMetrics as ResourceMetricsProto,
	ResourceType,
	TargetedActionMetrics as TargetedActionMetricsProto,
	UnitMetrics as UnitMetricsProto,
} from '../proto/api.js';
import { Class, Encounter as EncounterProto, Spec, SpellSchool, Target as TargetProto } from '../proto/common.js';
import { SimRun } from '../proto/ui.js';
import { ActionId, defaultTargetIcon } from '../proto_utils/action_id.js';
import { bucket, formatToNumber, sum } from '../utils.js';
import {
	AuraUptimeLog,
	CastLog,
	DamageDealtLog,
	DpsLog,
	Entity,
	MajorCooldownUsedLog,
	ResourceChangedLogGroup,
	SimLog,
	ThreatLogGroup,
} from './logs_parser.js';
import { cssClassForClass, getTalentTreeIcon, playerToSpec, specToClass } from './utils';

export interface SimResultFilter {
	// Raid index of the player to display, or null for all players.
	player?: number | null;

	// Target index of the target to display, or null for all targets.
	target?: number | null;
}

class SimResultData {
	readonly request: RaidSimRequest;
	readonly result: RaidSimResult;

	constructor(request: RaidSimRequest, result: RaidSimResult) {
		this.request = request;
		this.result = result;
	}

	get iterations() {
		return this.request.simOptions?.iterations || 1;
	}

	get duration() {
		return this.result.avgIterationDuration || 1;
	}

	get firstIterationDuration() {
		return this.result.firstIterationDuration || 1;
	}
}

// Holds all the data from a simulation call, and provides helper functions
// for parsing it.
export class SimResult {
	readonly request: RaidSimRequest;
	readonly result: RaidSimResult;

	readonly raidMetrics: RaidMetrics;
	readonly encounterMetrics: EncounterMetrics;
	readonly logs: Array<SimLog>;

	private players: Array<UnitMetrics>;
	private units: Array<UnitMetrics>;

	private constructor(request: RaidSimRequest, result: RaidSimResult, raidMetrics: RaidMetrics, encounterMetrics: EncounterMetrics, logs: Array<SimLog>) {
		this.request = request;
		this.result = result;
		this.raidMetrics = raidMetrics;
		this.encounterMetrics = encounterMetrics;
		this.logs = logs;

		this.players = raidMetrics.parties.map(party => party.players).flat();
		this.units = this.players.concat(encounterMetrics.targets);
	}

	get iterations() {
		return this.request.simOptions?.iterations || 1;
	}

	get duration() {
		return this.result.avgIterationDuration || 1;
	}

	get firstIterationDuration() {
		return this.result.firstIterationDuration || 1;
	}

	getPlayers(filter?: SimResultFilter): Array<UnitMetrics> {
		if (filter?.player || filter?.player === 0) {
			const player = this.getUnitWithIndex(filter.player);
			return player ? [player] : [];
		} else {
			return this.raidMetrics.parties.map(party => party.players).flat();
		}
	}

	getRaidIndexedPlayers(filter?: SimResultFilter): Array<UnitMetrics> {
		if (filter?.player || filter?.player === 0) {
			const player = this.getPlayerWithRaidIndex(filter.player);
			return player ? [player] : [];
		} else {
			return this.raidMetrics.parties.map(party => party.players).flat();
		}
	}

	// Returns the first player, regardless of which party / raid slot its in.
	getFirstPlayer(): UnitMetrics | null {
		return this.getPlayers()[0] || null;
	}

	getPlayerWithIndex(unitIndex: number): UnitMetrics | null {
		return this.players.find(player => player.unitIndex == unitIndex) || null;
	}
	getPlayerWithRaidIndex(raidIndex: number): UnitMetrics | null {
		return this.players.find(player => player.index == raidIndex) || null;
	}

	getTargets(filter?: SimResultFilter): Array<UnitMetrics> {
		if (filter?.target || filter?.target === 0) {
			const target = this.getUnitWithIndex(filter.target);
			return target ? [target] : [];
		} else {
			return this.encounterMetrics.targets.slice();
		}
	}

	getTargetWithIndex(unitIndex: number): UnitMetrics | null {
		return this.getTargets().find(target => target.unitIndex == unitIndex) || null;
	}
	getTargetWithEncounterIndex(index: number): UnitMetrics | null {
		return this.getTargets().find(target => target.index == index) || null;
	}
	getUnitWithIndex(unitIndex: number): UnitMetrics | null {
		return this.units.find(unit => unit.unitIndex == unitIndex) || null;
	}

	getDamageMetrics(filter: SimResultFilter): DistributionMetricsProto {
		if (filter.player || filter.player === 0) {
			return this.getPlayerWithIndex(filter.player)?.dps || DistributionMetricsProto.create();
		}

		return this.raidMetrics.dps;
	}

	getActionMetrics(filter?: SimResultFilter): Array<ActionMetrics> {
		return ActionMetrics.joinById(
			this.getPlayers(filter)
				.map(player => player.getPlayerAndPetActions().map(action => action.forTarget(filter)))
				.flat(),
		);
	}

	getRaidIndexedActionMetrics(filter?: SimResultFilter): Array<ActionMetrics> {
		return ActionMetrics.joinById(
			this.getRaidIndexedPlayers(filter)
				.map(player => player.getPlayerAndPetActions().map(action => action.forTarget(filter)))
				.flat(),
		);
	}

	getSpellMetrics(filter?: SimResultFilter): Array<ActionMetrics> {
		return this.getActionMetrics(filter).filter(e => e.hitAttempts != 0 && !e.isMeleeAction);
	}

	getMeleeMetrics(filter?: SimResultFilter): Array<ActionMetrics> {
		return this.getActionMetrics(filter).filter(e => e.hitAttempts != 0 && e.isMeleeAction);
	}

	getResourceMetrics(resourceType: ResourceType, filter?: SimResultFilter): Array<ResourceMetrics> {
		return ResourceMetrics.joinById(
			this.getPlayers(filter)
				.map(player => player.resources.filter(resource => resource.type == resourceType))
				.flat(),
		);
	}

	getBuffMetrics(filter?: SimResultFilter): Array<AuraMetrics> {
		return AuraMetrics.joinById(
			this.getPlayers(filter)
				.map(player => player.auras)
				.flat(),
		);
	}

	getDebuffMetrics(filter?: SimResultFilter): Array<AuraMetrics> {
		return AuraMetrics.joinById(
			this.getTargets(filter)
				.map(target => target.auras)
				.flat(),
		).filter(aura => aura.uptimePercent != 0);
	}

	toProto(): SimRun {
		return SimRun.create({
			request: this.request,
			result: this.result,
		});
	}

	static async fromProto(proto: SimRun): Promise<SimResult> {
		return SimResult.makeNew(proto.request || RaidSimRequest.create(), proto.result || RaidSimResult.create());
	}

	static async makeNew(request: RaidSimRequest, result: RaidSimResult): Promise<SimResult> {
		const resultData = new SimResultData(request, result);
		const logs = await SimLog.parseAll(result);
		const raidPromise = RaidMetrics.makeNew(resultData, request.raid!, result.raidMetrics!, logs);
		const encounterPromise = EncounterMetrics.makeNew(resultData, request.encounter!, result.encounterMetrics!, logs);

		const raidMetrics = await raidPromise;
		const encounterMetrics = await encounterPromise;

		return new SimResult(request, result, raidMetrics, encounterMetrics, logs);
	}
}

export class RaidMetrics {
	private readonly raid: RaidProto;
	private readonly metrics: RaidMetricsProto;

	readonly dps: DistributionMetricsProto;
	readonly hps: DistributionMetricsProto;
	readonly parties: Array<PartyMetrics>;

	private constructor(raid: RaidProto, metrics: RaidMetricsProto, parties: Array<PartyMetrics>) {
		this.raid = raid;
		this.metrics = metrics;
		this.dps = this.metrics.dps!;
		this.hps = this.metrics.hps!;
		this.parties = parties;
	}

	static async makeNew(resultData: SimResultData, raid: RaidProto, metrics: RaidMetricsProto, logs: Array<SimLog>): Promise<RaidMetrics> {
		const numParties = Math.min(raid.parties.length, metrics.parties.length);

		const parties = await Promise.all(
			[...new Array(numParties).keys()].map(i => PartyMetrics.makeNew(resultData, raid.parties[i], metrics.parties[i], i, logs)),
		);

		return new RaidMetrics(raid, metrics, parties);
	}
}

export class PartyMetrics {
	private readonly party: PartyProto;
	private readonly metrics: PartyMetricsProto;

	readonly partyIndex: number;
	readonly dps: DistributionMetricsProto;
	readonly hps: DistributionMetricsProto;
	readonly players: Array<UnitMetrics>;

	private constructor(party: PartyProto, metrics: PartyMetricsProto, partyIndex: number, players: Array<UnitMetrics>) {
		this.party = party;
		this.metrics = metrics;
		this.partyIndex = partyIndex;
		this.dps = this.metrics.dps!;
		this.hps = this.metrics.hps!;
		this.players = players;
	}

	static async makeNew(
		resultData: SimResultData,
		party: PartyProto,
		metrics: PartyMetricsProto,
		partyIndex: number,
		logs: Array<SimLog>,
	): Promise<PartyMetrics> {
		const numPlayers = Math.min(party.players.length, metrics.players.length);
		const players = await Promise.all(
			[...new Array(numPlayers).keys()]
				.filter(i => party.players[i].class != Class.ClassUnknown)
				.map(i => UnitMetrics.makeNewPlayer(resultData, party.players[i], metrics.players[i], partyIndex * 5 + i, false, logs)),
		);

		return new PartyMetrics(party, metrics, partyIndex, players);
	}
}

export class UnitMetrics {
	// If this Unit is a pet, player is the owner. If it's a target, player is null.
	private readonly player: PlayerProto | null;
	private readonly target: TargetProto | null;
	private readonly metrics: UnitMetricsProto;

	readonly index: number;
	readonly unitIndex: number;
	readonly name: string;
	readonly spec: Spec;
	readonly petActionId: ActionId | null;
	readonly iconUrl: string;
	readonly classColor: string;
	readonly dps: DistributionMetricsProto;
	readonly dpasp: DistributionMetricsProto;
	readonly hps: DistributionMetricsProto;
	readonly tps: DistributionMetricsProto;
	readonly dtps: DistributionMetricsProto;
	readonly tmi: DistributionMetricsProto;
	readonly tto: DistributionMetricsProto;
	readonly actions: Array<ActionMetrics>;
	readonly auras: Array<AuraMetrics>;
	readonly resources: Array<ResourceMetrics>;
	readonly pets: Array<UnitMetrics>;
	private readonly iterations: number;
	private readonly duration: number;

	readonly logs: Array<SimLog>;
	readonly damageDealtLogs: Array<DamageDealtLog>;
	readonly groupedResourceLogs: Record<ResourceType, Array<ResourceChangedLogGroup>>;
	readonly dpsLogs: Array<DpsLog>;
	readonly auraUptimeLogs: Array<AuraUptimeLog>;
	readonly majorCooldownLogs: Array<MajorCooldownUsedLog>;
	readonly castLogs: Array<CastLog>;
	readonly threatLogs: Array<ThreatLogGroup>;

	// Aura uptime logs, filtered to include only auras that correspond to a
	// major cooldown.
	readonly majorCooldownAuraUptimeLogs: Array<AuraUptimeLog>;

	private constructor(
		player: PlayerProto | null,
		target: TargetProto | null,
		petActionId: ActionId | null,
		metrics: UnitMetricsProto,
		index: number,
		actions: Array<ActionMetrics>,
		auras: Array<AuraMetrics>,
		resources: Array<ResourceMetrics>,
		pets: Array<UnitMetrics>,
		logs: Array<SimLog>,
		resultData: SimResultData,
	) {
		this.player = player;
		this.target = target;
		this.metrics = metrics;

		this.index = index;
		this.unitIndex = metrics.unitIndex;
		this.name = metrics.name;
		this.spec = player ? playerToSpec(player) : 0;
		this.petActionId = petActionId;
		this.iconUrl = this.isPlayer ? getTalentTreeIcon(this.spec, player!.talentsString) : this.isTarget ? defaultTargetIcon : '';
		this.classColor = this.isTarget ? '' : cssClassForClass(specToClass[this.spec]);
		this.dps = this.metrics.dps!;
		this.dpasp = this.metrics.dpasp!;
		this.hps = this.metrics.hps!;
		this.tps = this.metrics.threat!;
		this.dtps = this.metrics.dtps!;
		this.tmi = this.metrics.tmi!;
		this.tto = this.metrics.tto!;
		this.actions = actions;
		this.auras = auras;
		this.resources = resources;
		this.pets = pets;
		this.logs = logs;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;

		this.damageDealtLogs = this.logs.filter((log): log is DamageDealtLog => log.isDamageDealt());
		this.dpsLogs = DpsLog.fromLogs(this.damageDealtLogs);
		this.castLogs = CastLog.fromLogs(this.logs);
		this.threatLogs = ThreatLogGroup.fromLogs(this.logs);

		this.auraUptimeLogs = AuraUptimeLog.fromLogs(
			this.logs,
			new Entity(this.name, '', this.index, this.target != null, this.isPet),
			resultData.firstIterationDuration,
		);
		this.majorCooldownLogs = this.logs.filter((log): log is MajorCooldownUsedLog => log.isMajorCooldownUsed());

		this.groupedResourceLogs = ResourceChangedLogGroup.fromLogs(this.logs);
		AuraUptimeLog.populateActiveAuras(this.dpsLogs, this.auraUptimeLogs);
		AuraUptimeLog.populateActiveAuras(this.groupedResourceLogs[ResourceType.ResourceTypeMana], this.auraUptimeLogs);

		this.majorCooldownAuraUptimeLogs = this.auraUptimeLogs.filter(auraLog =>
			this.majorCooldownLogs.find(mcdLog => mcdLog.actionId!.equals(auraLog.actionId!)),
		);
	}

	get label() {
		if (this.target == null) {
			return `${this.name} (#${this.index + 1})`;
		} else {
			return this.name;
		}
	}

	get isPlayer() {
		return this.player != null;
	}

	get isTarget() {
		return this.target != null;
	}

	get isPet() {
		return this.petActionId != null;
	}

	// Returns the unit index of the target of this unit, as selected by the filter.
	getTargetIndex(filter?: SimResultFilter): number | null {
		if (!filter) {
			return null;
		}

		const index = this.isPlayer ? filter.target : filter.player;
		if (index == null || index == -1) {
			return null;
		}

		return index;
	}

	get inFrontOfTarget(): boolean {
		if (this.isTarget) {
			return true;
		} else if (this.isPlayer) {
			return this.player!.inFrontOfTarget;
		} else {
			return false; // TODO pets
		}
	}

	get chanceOfDeath(): DistributionMetricsProto {
		const p = Math.max(1.0, this.metrics.chanceOfDeath);
		const err = Math.sqrt(Math.abs(p * (1 - p)) / this.iterations);
		return DistributionMetricsProto.create({
			avg: p * 100,
			stdev: err * 100,
		});
	}

	get maxThreat() {
		return this.threatLogs[this.threatLogs.length - 1]?.threatAfter || 0;
	}

	get secondsOomAvg() {
		return this.metrics.secondsOomAvg;
	}

	get totalDamage() {
		return this.dps.avg * this.duration;
	}

	get totalDamageTaken() {
		return this.dtps.avg * this.duration;
	}

	getPlayerAndPetActions(): Array<ActionMetrics> {
		return this.actions.concat(this.pets.map(pet => pet.getPlayerAndPetActions()).flat());
	}

	private getActionsForDisplay(): Array<ActionMetrics> {
		return this.actions.filter(e => e.hitAttempts != 0 || e.tps != 0 || e.dps != 0);
	}

	getMeleeActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay().filter(e => e.isMeleeAction);
	}

	getMeleeDamageActions(): Array<ActionMetrics> {
		return this.getMeleeActions().filter(e => e.dps !== 0 && e.hps === 0);
	}

	getSpellActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay().filter(e => !e.isMeleeAction);
	}

	getSpellDamageActions(): Array<ActionMetrics> {
		return this.getSpellActions().filter(e => e.dps !== 0 && e.hps === 0);
	}

	getDamageActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay().filter(e => e.dps !== 0 && e.hps === 0);
	}

	getThreatActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay().filter(e => e.tps !== 0 && e.hps === 0);
	}

	getHealingActions(): Array<ActionMetrics> {
		return this.getActionsForDisplay();
	}

	getResourceMetrics(resourceType: ResourceType): Array<ResourceMetrics> {
		return this.resources.filter(resource => resource.type == resourceType);
	}

	static async makeNewPlayer(
		resultData: SimResultData,
		player: PlayerProto,
		metrics: UnitMetricsProto,
		raidIndex: number,
		isPet: boolean,
		logs: Array<SimLog>,
	): Promise<UnitMetrics> {
		const playerLogs = logs.filter(
			log => log.source && !log.source.isTarget && isPet == log.source.isPet && (isPet ? log.source.name == metrics.name : log.source.index == raidIndex),
		);
		const petLogs = logs.filter(log => log.source && !log.source.isTarget && log.source.isPet && log.source.index == raidIndex);

		const actionsPromise = Promise.all(metrics.actions.map(actionMetrics => ActionMetrics.makeNew(null, resultData, actionMetrics, raidIndex)));
		const aurasPromise = Promise.all(metrics.auras.map(auraMetrics => AuraMetrics.makeNew(null, resultData, auraMetrics, raidIndex)));
		const resourcesPromise = Promise.all(metrics.resources.map(resourceMetrics => ResourceMetrics.makeNew(null, resultData, resourceMetrics, raidIndex)));
		const petsPromise = Promise.all(metrics.pets.map(petMetrics => UnitMetrics.makeNewPlayer(resultData, player, petMetrics, raidIndex, true, petLogs)));

		let petIdPromise: Promise<ActionId | null> = Promise.resolve(null);
		if (isPet) {
			petIdPromise = ActionId.fromPetName(metrics.name).fill(raidIndex);
		}

		const actions = await actionsPromise;
		const auras = await aurasPromise;
		const resources = await resourcesPromise;
		const pets = await petsPromise;
		const petActionId = await petIdPromise;

		const playerMetrics = new UnitMetrics(player, null, petActionId, metrics, raidIndex, actions, auras, resources, pets, playerLogs, resultData);
		actions.forEach(action => {
			action.unit = playerMetrics;
			action.resources = resources.filter(resourceMetrics => resourceMetrics.actionId.equals(action.actionId));
		});
		auras.forEach(aura => (aura.unit = playerMetrics));
		resources.forEach(resource => (resource.unit = playerMetrics));
		return playerMetrics;
	}

	static async makeNewTarget(
		resultData: SimResultData,
		target: TargetProto,
		metrics: UnitMetricsProto,
		index: number,
		logs: Array<SimLog>,
	): Promise<UnitMetrics> {
		const targetLogs = logs.filter(log => log.source && log.source.isTarget && log.source.index == index);

		const actionsPromise = Promise.all(metrics.actions.map(actionMetrics => ActionMetrics.makeNew(null, resultData, actionMetrics, index)));
		const aurasPromise = Promise.all(metrics.auras.map(auraMetrics => AuraMetrics.makeNew(null, resultData, auraMetrics)));

		const actions = await actionsPromise;
		const auras = await aurasPromise;

		const targetMetrics = new UnitMetrics(null, target, null, metrics, index, actions, auras, [], [], targetLogs, resultData);
		actions.forEach(action => (action.unit = targetMetrics));
		auras.forEach(aura => (aura.unit = targetMetrics));
		return targetMetrics;
	}
}

export class EncounterMetrics {
	private readonly encounter: EncounterProto;
	private readonly metrics: EncounterMetricsProto;

	readonly targets: Array<UnitMetrics>;

	private constructor(encounter: EncounterProto, metrics: EncounterMetricsProto, targets: Array<UnitMetrics>) {
		this.encounter = encounter;
		this.metrics = metrics;
		this.targets = targets;
	}

	static async makeNew(resultData: SimResultData, encounter: EncounterProto, metrics: EncounterMetricsProto, logs: Array<SimLog>): Promise<EncounterMetrics> {
		const numTargets = Math.min(encounter.targets.length, metrics.targets.length);
		const targets = await Promise.all(
			[...new Array(numTargets).keys()].map(i => UnitMetrics.makeNewTarget(resultData, encounter.targets[i], metrics.targets[i], i, logs)),
		);

		return new EncounterMetrics(encounter, metrics, targets);
	}

	get durationSeconds() {
		return this.encounter.duration;
	}
}

export class AuraMetrics {
	unit: UnitMetrics | null;
	readonly actionId: ActionId;
	readonly name: string;
	readonly iconUrl: string;
	private readonly resultData: SimResultData;
	private readonly iterations: number;
	private readonly duration: number;
	private readonly data: AuraMetricsProto;

	private constructor(unit: UnitMetrics | null, actionId: ActionId, data: AuraMetricsProto, resultData: SimResultData) {
		this.unit = unit;
		this.actionId = actionId;
		this.name = actionId.name;
		this.iconUrl = actionId.iconUrl;
		this.data = data;
		this.resultData = resultData;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;
	}

	get uptimePercent() {
		return (this.data.uptimeSecondsAvg / this.duration) * 100;
	}

	get averageProcs() {
		return this.data.procsAvg;
	}

	get ppm() {
		return this.data.procsAvg / (this.duration / 60);
	}

	static async makeNew(unit: UnitMetrics | null, resultData: SimResultData, auraMetrics: AuraMetricsProto, playerIndex?: number): Promise<AuraMetrics> {
		const actionId = await ActionId.fromProto(auraMetrics.id!).fill(playerIndex);
		return new AuraMetrics(unit, actionId, auraMetrics, resultData);
	}

	// Merges an array of metrics into a single metrics.
	static merge(auras: Array<AuraMetrics>, { removeTag, actionIdOverride }: { removeTag?: boolean; actionIdOverride?: ActionId } = {}): AuraMetrics {
		const firstAura = auras[0];
		const unit = auras.every(aura => aura.unit == firstAura.unit) ? firstAura.unit : null;
		let actionId = actionIdOverride || firstAura.actionId;
		if (removeTag) {
			actionId = actionId.withoutTag();
		}
		return new AuraMetrics(
			unit,
			actionId,
			AuraMetricsProto.create({
				uptimeSecondsAvg: Math.max(...auras.map(a => a.data.uptimeSecondsAvg)),
			}),
			firstAura.resultData,
		);
	}

	// Groups similar metrics, i.e. metrics with the same item/spell/other ID but
	// different tags, and returns them as separate arrays.
	static groupById(auras: Array<AuraMetrics>, useTag?: boolean): Array<Array<AuraMetrics>> {
		if (useTag) {
			return Object.values(bucket(auras, aura => aura.actionId.toString()));
		} else {
			return Object.values(bucket(auras, aura => aura.actionId.toStringIgnoringTag()));
		}
	}

	// Merges aura metrics that have the same name/ID, adding their stats together.
	static joinById(auras: Array<AuraMetrics>, useTag?: boolean): Array<AuraMetrics> {
		return AuraMetrics.groupById(auras, useTag).map(aurasToJoin => AuraMetrics.merge(aurasToJoin));
	}
}
export class ResourceMetrics {
	unit: UnitMetrics | null;
	readonly actionId: ActionId;
	readonly name: string;
	readonly iconUrl: string;
	readonly type: ResourceType;
	private readonly resultData: SimResultData;
	private readonly iterations: number;
	private readonly duration: number;
	private readonly data: ResourceMetricsProto;

	private constructor(unit: UnitMetrics | null, actionId: ActionId, data: ResourceMetricsProto, resultData: SimResultData) {
		this.unit = unit;
		this.actionId = actionId;
		this.name = actionId.name;
		this.iconUrl = actionId.iconUrl;
		this.type = data.type;
		this.resultData = resultData;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;
		this.data = data;
	}

	get events() {
		return this.data.events / this.iterations;
	}

	get gain() {
		return this.data.gain / this.iterations;
	}

	get gainPerSecond() {
		return this.data.gain / this.iterations / this.duration;
	}

	get avgGain() {
		return this.data.gain / this.data.events;
	}

	get wastedGain() {
		return (this.data.gain - this.data.actualGain) / this.iterations;
	}

	static async makeNew(
		unit: UnitMetrics | null,
		resultData: SimResultData,
		resourceMetrics: ResourceMetricsProto,
		playerIndex?: number,
	): Promise<ResourceMetrics> {
		const actionId = await ActionId.fromProto(resourceMetrics.id!).fill(playerIndex);
		return new ResourceMetrics(unit, actionId, resourceMetrics, resultData);
	}

	// Merges an array of metrics into a single metrics.
	static merge(
		resources: Array<ResourceMetrics>,
		{ removeTag, actionIdOverride }: { removeTag?: boolean; actionIdOverride?: ActionId } = {},
	): ResourceMetrics {
		const firstResource = resources[0];
		const unit = resources.every(resource => resource.unit == firstResource.unit) ? firstResource.unit : null;
		let actionId = actionIdOverride || firstResource.actionId;
		if (removeTag) {
			actionId = actionId.withoutTag();
		}
		return new ResourceMetrics(
			unit,
			actionId,
			ResourceMetricsProto.create({
				events: sum(resources.map(a => a.data.events)),
				gain: sum(resources.map(a => a.data.gain)),
				actualGain: sum(resources.map(a => a.data.actualGain)),
			}),
			firstResource.resultData,
		);
	}

	// Groups similar metrics, i.e. metrics with the same item/spell/other ID but
	// different tags, and returns them as separate arrays.
	static groupById(resources: Array<ResourceMetrics>, useTag?: boolean): Array<Array<ResourceMetrics>> {
		if (useTag) {
			return Object.values(bucket(resources, resource => resource.actionId.toString()));
		} else {
			return Object.values(bucket(resources, resource => resource.actionId.toStringIgnoringTag()));
		}
	}

	// Merges resource metrics that have the same name/ID, adding their stats together.
	static joinById(resources: Array<ResourceMetrics>, useTag?: boolean): Array<ResourceMetrics> {
		return ResourceMetrics.groupById(resources, useTag).map(resourcesToJoin => ResourceMetrics.merge(resourcesToJoin));
	}
}

// Manages the metrics for a single unit action (e.g. Lightning Bolt).
export class ActionMetrics {
	unit: UnitMetrics | null;
	readonly actionId: ActionId;
	readonly name: string;
	readonly iconUrl: string;
	readonly spellSchool: SpellSchool | null;
	readonly targets: Array<TargetedActionMetrics>;
	private readonly resultData: SimResultData;
	private readonly iterations: number;
	private readonly duration: number;
	private readonly data: ActionMetricsProto;
	private readonly combinedMetrics: TargetedActionMetrics;
	resources: Array<ResourceMetrics>;

	private constructor(unit: UnitMetrics | null, actionId: ActionId, data: ActionMetricsProto, resultData: SimResultData) {
		this.unit = unit;
		this.actionId = actionId;
		this.name = actionId.name;
		this.iconUrl = actionId.iconUrl;
		this.resultData = resultData;
		this.iterations = resultData.iterations;
		this.duration = resultData.duration;
		this.data = data;
		this.spellSchool = data.spellSchool;
		this.targets = data.targets.map(
			tam =>
				new TargetedActionMetrics(tam, {
					iterations: this.iterations,
					duration: this.duration,
				}),
		);
		this.combinedMetrics = TargetedActionMetrics.merge(this.targets);
		this.resources = [];
	}

	get isMeleeAction() {
		return this.data.isMelee;
	}

	get isPassiveAction() {
		return this.data.isPassive;
	}

	get totalDamagePercent() {
		const totalAvgDps = this.resultData.result.raidMetrics?.dps?.avg;
		if (!totalAvgDps) return undefined;

		return (this.avgDamage / (totalAvgDps * this.duration)) * 100;
	}

	get damage() {
		return this.combinedMetrics.damage;
	}

	get avgDamage() {
		return this.combinedMetrics.avgDamage;
	}

	get avgHitDamage() {
		return (
			this.avgDamage -
			this.avgTickDamage -
			this.avgCritDamage +
			this.avgCritTickDamage -
			this.avgGlanceDamage -
			this.avgBlockDamage -
			this.avgBlockedCritDamage -
			this.avgCrushDamage
		);
	}

	get resistedDamage() {
		return this.combinedMetrics.resistedDamage;
	}

	get avgResistedDamage() {
		return this.combinedMetrics.avgResistedDamage;
	}

	get critDamage() {
		return this.combinedMetrics.critDamage;
	}

	get avgCritDamage() {
		return this.combinedMetrics.avgCritDamage;
	}

	get resistedCritDamage() {
		return this.combinedMetrics.resistedCritDamage;
	}

	get avgResistedCritDamage() {
		return this.combinedMetrics.avgResistedCritDamage;
	}

	get tickDamage() {
		return this.combinedMetrics.tickDamage;
	}

	get avgTickDamage() {
		return this.combinedMetrics.avgTickDamage;
	}

	get resistedTickDamage() {
		return this.combinedMetrics.resistedTickDamage;
	}

	get avgResistedTickDamage() {
		return this.combinedMetrics.avgResistedTickDamage;
	}

	get critTickDamage() {
		return this.combinedMetrics.critTickDamage;
	}

	get avgCritTickDamage() {
		return this.combinedMetrics.avgCritTickDamage;
	}
	get resistedCritTickDamage() {
		return this.combinedMetrics.resistedCritTickDamage;
	}

	get avgResistedCritTickDamage() {
		return this.combinedMetrics.avgResistedCritTickDamage;
	}

	get glanceDamage() {
		return this.combinedMetrics.glanceDamage;
	}

	get avgGlanceDamage() {
		return this.combinedMetrics.avgGlanceDamage;
	}

	get blockDamage() {
		return this.combinedMetrics.blockDamage;
	}

	get crushDamage() {
		return this.combinedMetrics.crushDamage;
	}

	get avgCrushDamage() {
		return this.combinedMetrics.avgCrushDamage;
	}

	get avgBlockDamage() {
		return this.combinedMetrics.avgBlockDamage;
	}

	get blockedCritDamage() {
		return this.combinedMetrics.blockedCritDamage;
	}

	get avgBlockedCritDamage() {
		return this.combinedMetrics.avgBlockedCritDamage;
	}

	get dps() {
		return this.combinedMetrics.dps;
	}

	get threat() {
		return this.combinedMetrics.threat;
	}

	get avgThreat() {
		return this.combinedMetrics.avgThreat;
	}

	get totalThreatPercent() {
		const totalAvgTps = this.resultData.result.raidMetrics?.parties[0].players[0].threat?.avg;
		if (!totalAvgTps) return undefined;

		return (this.avgThreat / (totalAvgTps * this.duration)) * 100;
	}

	get tps() {
		return this.combinedMetrics.tps;
	}

	get totalHealingPercent() {
		const totalAvgHps = this.resultData.result.raidMetrics?.hps?.avg;
		if (!totalAvgHps) return undefined;

		return (this.avgHealing / (totalAvgHps * this.duration)) * 100;
	}

	get healing() {
		return this.combinedMetrics.healing;
	}

	get avgHealing() {
		return this.combinedMetrics.healing / this.iterations;
	}

	get critHealing() {
		return this.combinedMetrics.critHealing;
	}

	get avgCritHealing() {
		return this.combinedMetrics.critHealing / this.iterations;
	}

	get hps() {
		return this.combinedMetrics.hps;
	}

	get casts() {
		if (this.isPassiveAction) return 0;
		return this.combinedMetrics.casts;
	}

	get castsPerMinute() {
		if (this.isPassiveAction) return 0;
		return this.combinedMetrics.castsPerMinute;
	}

	get avgCastTimeMs() {
		if (this.isPassiveAction) return 0;
		return this.combinedMetrics.avgCastTimeMs;
	}

	get hpm() {
		const totalHealing = this.combinedMetrics.hps * this.duration;
		const manaMetrics = this.resources.find(r => r.type == ResourceType.ResourceTypeMana);
		if (manaMetrics) {
			return totalHealing / -manaMetrics.gain;
		}

		return 0;
	}

	get damageThroughput() {
		if (this.unit?.isPet && !this.actionId.spellId) return 0;
		return this.combinedMetrics.damageThroughput;
	}

	get threatThroughput() {
		if (this.unit?.isPet && !this.actionId.spellId) return 0;
		return this.combinedMetrics.threatThroughput;
	}

	get healingThroughput() {
		return this.combinedMetrics.healingThroughput;
	}

	get shielding() {
		return this.combinedMetrics.shielding;
	}

	get avgCast() {
		if (this.isPassiveAction) return 0;
		return this.combinedMetrics.avgCast;
	}

	get avgCastHit() {
		if (!this.combinedMetrics.avgCast) return 0;
		return this.combinedMetrics.avgCast - this.avgCastTick;
	}

	get avgCastTick() {
		return this.combinedMetrics.avgCastTick;
	}

	get avgCastHealing() {
		if (this.isPassiveAction) return 0;
		return this.combinedMetrics.avgCastHealing;
	}

	get avgCastThreat() {
		if (this.isPassiveAction) return 0;
		return this.combinedMetrics.avgCastThreat;
	}

	get landedHits() {
		return this.combinedMetrics.landedHits;
	}

	get landedTicks() {
		return this.combinedMetrics.landedTicks;
	}

	get hitAttempts() {
		return this.combinedMetrics.hitAttempts;
	}

	get avgHit() {
		return this.combinedMetrics.avgHit;
	}

	get avgTick() {
		return this.combinedMetrics.avgTick;
	}

	get avgHitHealing() {
		return this.combinedMetrics.avgHitHealing;
	}

	get avgHitThreat() {
		return this.combinedMetrics.avgHitThreat;
	}

	get totalMisses() {
		return this.misses + this.dodges + this.parries;
	}

	get totalMissesPercent() {
		return this.missPercent + this.dodgePercent + this.parryPercent;
	}

	get misses() {
		return this.combinedMetrics.misses;
	}

	get missPercent() {
		return this.combinedMetrics.missPercent;
	}

	get dodges() {
		return this.combinedMetrics.dodges;
	}

	get dodgePercent() {
		return this.combinedMetrics.dodgePercent;
	}

	get parries() {
		return this.combinedMetrics.parries;
	}

	get parryPercent() {
		return this.combinedMetrics.parryPercent;
	}

	get hits() {
		return this.combinedMetrics.hits;
	}

	get resistedHits() {
		return this.combinedMetrics.resistedHits;
	}

	get hitPercent() {
		return this.combinedMetrics.hitPercent;
	}

	get resistedHitPercent() {
		return this.combinedMetrics.resistedHitPercent;
	}

	get ticks() {
		return this.combinedMetrics.ticks;
	}

	get resistedTicks() {
		return this.combinedMetrics.resistedTicks;
	}

	get resistedTickPercent() {
		return this.combinedMetrics.resistedTickPercent;
	}

	get critTicks() {
		return this.combinedMetrics.critTicks;
	}

	get critTickPercent() {
		return this.combinedMetrics.critTickPercent;
	}

	get resistedCritTicks() {
		return this.combinedMetrics.resistedCritTicks;
	}

	get resistedCritTickPercent() {
		return this.combinedMetrics.resistedCritTickPercent;
	}

	get crushes() {
		return this.combinedMetrics.crushes;
	}

	get crushPercent() {
		return this.combinedMetrics.crushPercent;
	}

	get blocks() {
		return this.combinedMetrics.blocks;
	}

	get blockPercent() {
		return this.combinedMetrics.blockPercent;
	}

	get blockedCrits() {
		return this.combinedMetrics.blockedCrits;
	}

	get blockedCritPercent() {
		return this.combinedMetrics.blockedCritPercent;
	}

	get glances() {
		return this.combinedMetrics.glances;
	}

	get glancePercent() {
		return this.combinedMetrics.glancePercent;
	}

	get crits() {
		return this.combinedMetrics.crits;
	}

	get critPercent() {
		return this.combinedMetrics.critPercent;
	}

	get resistedCrits() {
		return this.combinedMetrics.resistedCrits;
	}

	get resistedCritPercent() {
		return this.combinedMetrics.resistedCritPercent;
	}

	get healingPercent() {
		return this.combinedMetrics.healingPercent;
	}

	get healingCritPercent() {
		return this.combinedMetrics.healingCritPercent;
	}

	get damageDone() {
		const normalHitAvgDamage = Number(
			(
				this.avgDamage -
				this.avgResistedDamage +
				this.avgResistedTickDamage +
				this.avgResistedCritDamage -
				this.avgResistedCritTickDamage -
				this.avgCritDamage +
				this.avgCritTickDamage -
				this.avgTickDamage -
				this.avgGlanceDamage -
				this.avgBlockDamage -
				this.avgBlockedCritDamage -
				this.avgCrushDamage
			).toFixed(8),
		);
		const normalResistedHitAvgDamage = Number(
			(this.avgResistedDamage - this.avgResistedCritDamage - this.avgResistedTickDamage + this.avgResistedCritTickDamage).toFixed(8),
		);
		const critHitAvgDamage = Number((this.avgCritDamage - this.avgResistedCritDamage - this.avgCritTickDamage + this.avgResistedCritTickDamage).toFixed(8));
		const resistedCritHitAvgDamage = Number((this.avgResistedCritDamage - this.avgResistedCritTickDamage).toFixed(8));

		const normalTickAvgDamage = Number(
			(this.avgTickDamage - this.avgCritTickDamage - this.avgResistedTickDamage + this.avgResistedCritTickDamage).toFixed(8),
		);
		const normalResistedTickAvgDamage = Number((this.avgResistedTickDamage - this.avgResistedCritTickDamage).toFixed(8));
		const normalCritTickAvgDamage = Number((this.avgCritTickDamage - this.avgResistedCritTickDamage).toFixed(8));

		return {
			hit: {
				value: normalHitAvgDamage,
				percentage: (normalHitAvgDamage / this.avgDamage) * 100,
				average: normalHitAvgDamage / (this.hits - this.resistedHits),
			},
			resistedHit: {
				value: normalResistedHitAvgDamage,
				percentage: (normalResistedHitAvgDamage / this.avgDamage) * 100,
				average: normalResistedHitAvgDamage / this.resistedHits,
			},
			critHit: {
				value: critHitAvgDamage,
				percentage: (critHitAvgDamage / this.avgDamage) * 100,
				average: critHitAvgDamage / (this.crits - this.resistedCrits),
			},
			resistedCritHit: {
				value: resistedCritHitAvgDamage,
				percentage: (resistedCritHitAvgDamage / this.avgDamage) * 100,
				average: resistedCritHitAvgDamage / this.resistedCrits,
			},
			tick: {
				value: normalTickAvgDamage,
				percentage: (normalTickAvgDamage / this.avgDamage) * 100,
				average: normalTickAvgDamage / (this.ticks - this.resistedTicks),
			},
			resistedTick: {
				value: normalResistedTickAvgDamage,
				percentage: (normalResistedTickAvgDamage / this.avgDamage) * 100,
				average: normalResistedTickAvgDamage / this.resistedTicks,
			},
			critTick: {
				value: normalCritTickAvgDamage,
				percentage: (normalCritTickAvgDamage / this.avgDamage) * 100,
				average: normalCritTickAvgDamage / (this.critTicks - this.resistedCritTicks),
			},
			resistedCritTick: {
				value: this.avgResistedCritTickDamage,
				percentage: (this.avgResistedCritTickDamage / this.avgDamage) * 100,
				average: this.avgResistedCritTickDamage / this.resistedCritTicks,
			},
			glance: {
				value: this.avgGlanceDamage,
				percentage: (this.avgGlanceDamage / this.avgDamage) * 100,
				average: this.avgGlanceDamage / this.glances,
			},
			block: {
				value: this.avgBlockDamage,
				percentage: (this.avgBlockDamage / this.avgDamage) * 100,
				average: this.avgBlockDamage / this.blocks,
			},
			blockedCrit: {
				value: this.avgBlockedCritDamage,
				percentage: (this.avgBlockedCritDamage / this.avgDamage) * 100,
				average: this.avgBlockedCritDamage / this.blockedCrits,
			},
			crush: {
				value: this.avgCrushDamage,
				percentage: (this.avgCrushDamage / this.avgDamage) * 100,
				average: this.avgCrushDamage / this.crushes,
			},
		};
	}

	forTarget(filter?: SimResultFilter): ActionMetrics {
		const unitIndex = this.unit!.getTargetIndex(filter);
		if (unitIndex == null) {
			return this;
		} else {
			const target = this.targets.find(target => target.data.unitIndex == unitIndex);
			if (target) {
				const targetData = ActionMetricsProto.clone(this.data);
				targetData.targets = [target.data];
				return new ActionMetrics(this.unit, this.actionId, targetData, this.resultData);
			} else {
				throw new Error('Could not find target with unitIndex ' + unitIndex);
			}
		}
	}

	static async makeNew(unit: UnitMetrics | null, resultData: SimResultData, actionMetrics: ActionMetricsProto, playerIndex?: number): Promise<ActionMetrics> {
		const actionId = await ActionId.fromProto(actionMetrics.id!).fill(playerIndex);
		return new ActionMetrics(unit, actionId, actionMetrics, resultData);
	}

	// Merges an array of metrics into a single metric.
	static merge(actions: Array<ActionMetrics>, { removeTag, actionIdOverride }: { removeTag?: boolean; actionIdOverride?: ActionId } = {}): ActionMetrics {
		const firstAction = actions[0];
		const unit = firstAction.unit;
		let actionId = actionIdOverride || firstAction.actionId;
		if (removeTag) {
			actionId = actionId.withoutTag();
		}

		const maxTargets = Math.max(...actions.map(action => action.targets.length));
		const mergedTargets = [...Array(maxTargets).keys()].map(i => TargetedActionMetrics.merge(actions.map(action => action.targets[i])));
		const isAllPassiveSpells = actions.every(action => action.isPassiveAction);

		return new ActionMetrics(
			unit,
			actionId,
			ActionMetricsProto.create({
				isMelee: firstAction.isMeleeAction,
				isPassive: isAllPassiveSpells,
				targets: mergedTargets.map(t => t.data),
				spellSchool: firstAction.spellSchool || undefined,
			}),
			firstAction.resultData,
		);
	}

	// Groups similar metrics, i.e. metrics with the same item/spell/other ID but
	// different tags, and returns them as separate arrays.
	static groupById(actions: Array<ActionMetrics>, useTag?: boolean): Array<Array<ActionMetrics>> {
		if (useTag) {
			return Object.values(bucket(actions, action => action.actionId.toString()));
		} else {
			return Object.values(bucket(actions, action => action.actionId.toStringIgnoringTag()));
		}
	}

	// Merges action metrics that have the same name/ID, adding their stats together.
	static joinById(actions: Array<ActionMetrics>, useTag?: boolean): Array<ActionMetrics> {
		return ActionMetrics.groupById(actions, useTag).map(actionsToJoin => ActionMetrics.merge(actionsToJoin));
	}
}

type TargetedActionMetricsOptions = {
	iterations: number;
	duration: number;
};

// Manages the metrics for a single action applied to a specific target.
export class TargetedActionMetrics {
	private readonly iterations: number;
	private readonly duration: number;
	readonly data: TargetedActionMetricsProto;

	readonly landedHitsRaw: number;
	readonly landedTicksRaw: number;
	readonly hitAttempts: number;

	constructor(data: TargetedActionMetricsProto, { iterations, duration }: TargetedActionMetricsOptions) {
		this.iterations = iterations;
		this.duration = duration;
		this.data = data;

		this.landedHitsRaw = this.data.hits + this.data.crits + this.data.blocks + this.data.blockedCrits + this.data.glances + this.data.crushes;
		this.landedTicksRaw = this.data.ticks + this.data.critTicks;

		this.hitAttempts =
			this.data.misses +
			this.data.dodges +
			this.data.parries +
			this.data.blocks +
			this.data.blockedCrits +
			this.data.glances +
			this.data.crits +
			this.data.crushes;

		if (this.data.hits != 0) {
			this.hitAttempts += this.data.hits;
		} else if (this.data.hits == 0 && this.data.ticks > 0) {
			this.hitAttempts += this.data.casts;
		}
	}

	get damage() {
		return this.data.damage;
	}

	get avgDamage() {
		return this.data.damage / this.iterations;
	}

	get resistedDamage() {
		return this.data.resistedDamage;
	}

	get avgResistedDamage() {
		return this.data.resistedDamage / this.iterations;
	}

	get critDamage() {
		return this.data.critDamage;
	}

	get avgCritDamage() {
		return this.data.critDamage / this.iterations;
	}

	get resistedCritDamage() {
		return this.data.resistedCritDamage;
	}

	get avgResistedCritDamage() {
		return this.data.resistedCritDamage / this.iterations;
	}

	get tickDamage() {
		return this.data.tickDamage;
	}

	get avgTickDamage() {
		return this.data.tickDamage / this.iterations;
	}

	get resistedTickDamage() {
		return this.data.resistedTickDamage;
	}

	get avgResistedTickDamage() {
		return this.data.resistedTickDamage / this.iterations;
	}

	get critTickDamage() {
		return this.data.critTickDamage;
	}

	get avgCritTickDamage() {
		return this.data.critTickDamage / this.iterations;
	}

	get resistedCritTickDamage() {
		return this.data.resistedCritTickDamage;
	}

	get avgResistedCritTickDamage() {
		return this.data.resistedCritTickDamage / this.iterations;
	}

	get glanceDamage() {
		return this.data.glanceDamage;
	}

	get avgGlanceDamage() {
		return this.data.glanceDamage / this.iterations;
	}

	get crushDamage() {
		return this.data.crushDamage;
	}

	get avgCrushDamage() {
		return this.data.crushDamage / this.iterations;
	}

	get blockDamage() {
		return this.data.blockDamage;
	}

	get avgBlockDamage() {
		return this.data.blockDamage / this.iterations;
	}

	get blockedCritDamage() {
		return this.data.blockedCritDamage;
	}

	get avgBlockedCritDamage() {
		return this.data.blockedCritDamage / this.iterations;
	}

	get dps() {
		return this.data.damage / this.iterations / this.duration;
	}

	get threat() {
		return this.data.threat;
	}

	get avgThreat() {
		return this.data.threat / this.iterations;
	}

	get tps() {
		return this.data.threat / this.iterations / this.duration;
	}

	get healing() {
		return this.data.healing + this.data.shielding;
	}

	get avgHealing() {
		return (this.data.healing + this.data.shielding) / this.iterations;
	}

	get critHealing() {
		return this.data.critHealing;
	}

	get avgCritHealing() {
		return this.data.critHealing / this.iterations;
	}

	get shielding() {
		return this.data.shielding;
	}

	get hps() {
		return (this.data.healing + this.data.shielding) / this.iterations / this.duration;
	}

	get casts() {
		return this.data.casts / this.iterations;
	}

	get castsPerMinute() {
		return this.casts / (this.duration / 60);
	}

	get avgCastTimeMs() {
		return this.data.castTimeMs / this.iterations / this.casts;
	}

	get damageThroughput() {
		if (this.avgCastTimeMs) {
			return this.avgCast / (this.avgCastTimeMs / 1000);
		} else {
			return 0;
		}
	}

	get threatThroughput() {
		if (this.avgCastTimeMs) {
			return Math.max(this.avgHitThreat, this.avgCastThreat) / (this.avgCastTimeMs / 1000);
		} else {
			return 0;
		}
	}

	get healingThroughput() {
		if (this.avgCastTimeMs) {
			return this.hps / (this.avgCastTimeMs / 1000);
		} else {
			return 0;
		}
	}

	get timeSpentCastingMs() {
		return this.data.castTimeMs / this.iterations;
	}

	get avgCast() {
		if (!this.casts) return 0;
		return this.data.damage / this.iterations / (this.casts || 1);
	}

	get avgCastTick() {
		return this.data.tickDamage / this.iterations / (this.casts || 1);
	}

	get avgCastHealing() {
		return (this.data.healing + this.data.shielding) / this.iterations / (this.casts || 1);
	}

	get avgCastThreat() {
		return this.data.threat / this.iterations / (this.casts || 1);
	}

	get landedHits() {
		return this.landedHitsRaw / this.iterations;
	}

	get landedTicks() {
		return this.landedTicksRaw / this.iterations;
	}

	get avgHit() {
		const lhr = this.landedHitsRaw;
		return lhr == 0 ? 0 : (this.data.damage - this.data.tickDamage) / lhr;
	}

	get avgTick() {
		const ltr = this.landedTicksRaw;
		return ltr == 0 ? 0 : this.data.tickDamage / ltr;
	}

	get avgHitHealing() {
		return (this.data.healing + this.data.shielding) / this.iterations / this.landedHits;
	}

	get avgHitThreat() {
		const lhr = this.landedHitsRaw;
		return lhr == 0 ? 0 : this.data.threat / lhr;
	}

	get totalMisses() {
		return this.misses + this.dodges + this.parries;
	}

	get totalMissesPercent() {
		return this.missPercent + this.dodgePercent + this.parryPercent;
	}

	get misses() {
		return this.data.misses / this.iterations;
	}

	get missPercent() {
		return (this.data.misses / this.hitAttempts) * 100;
	}

	get dodges() {
		return this.data.dodges / this.iterations;
	}

	get dodgePercent() {
		return (this.data.dodges / this.hitAttempts) * 100;
	}

	get parries() {
		return this.data.parries / this.iterations;
	}

	get parryPercent() {
		return (this.data.parries / this.hitAttempts) * 100;
	}

	get hits() {
		return this.data.hits / this.iterations;
	}

	get hitPercent() {
		return (this.data.hits / this.hitAttempts) * 100;
	}

	get resistedHits() {
		return this.data.resistedHits / this.iterations;
	}

	get resistedHitPercent() {
		return (this.data.resistedHits / this.hitAttempts) * 100;
	}

	get ticks() {
		return this.data.ticks / this.iterations;
	}

	get resistedTicks() {
		return this.data.resistedTicks / this.iterations;
	}

	get resistedTickPercent() {
		return (this.data.resistedTicks / (this.data.ticks + this.data.critTicks)) * 100;
	}

	get critTicks() {
		return this.data.critTicks / this.iterations;
	}

	get critTickPercent() {
		return (this.data.critTicks / (this.data.ticks + this.data.critTicks)) * 100;
	}

	get resistedCritTicks() {
		return this.data.resistedCritTicks / this.iterations;
	}

	get resistedCritTickPercent() {
		return (this.data.resistedCritTicks / (this.data.ticks + this.data.critTicks)) * 100;
	}

	get blocks() {
		return this.data.blocks / this.iterations;
	}

	get crushes() {
		return this.data.crushes / this.iterations;
	}

	get crushPercent() {
		return (this.data.crushes / this.hitAttempts) * 100;
	}

	get blockPercent() {
		return (this.data.blocks / this.hitAttempts) * 100;
	}

	get blockedCrits() {
		return this.data.blockedCrits / this.iterations;
	}

	get blockedCritPercent() {
		return (this.data.blockedCrits / this.hitAttempts) * 100;
	}

	get glances() {
		return this.data.glances / this.iterations;
	}

	get glancePercent() {
		return (this.data.glances / this.hitAttempts) * 100;
	}

	get crits() {
		return this.data.crits / this.iterations;
	}

	get critPercent() {
		return (this.data.crits / this.hitAttempts) * 100;
	}

	get resistedCrits() {
		return this.data.resistedCrits / this.iterations;
	}

	get resistedCritPercent() {
		return (this.data.resistedCrits / this.hitAttempts) * 100;
	}

	get healingPercent() {
		return ((this.healing - this.critHealing) / this.healing) * 100;
	}

	get healingCritPercent() {
		return (this.data.critHealing / this.healing) * 100;
	}

	// Merges an array of metrics into a single metric.
	static merge(actions: Array<TargetedActionMetrics>): TargetedActionMetrics {
		const { iterations = 1, duration = 1 } = actions[0];

		return new TargetedActionMetrics(
			TargetedActionMetricsProto.create({
				casts: sum(actions.map(a => a.data.casts)),
				hits: sum(actions.map(a => a.data.hits)),
				resistedHits: sum(actions.map(a => a.data.resistedHits)),
				crits: sum(actions.map(a => a.data.crits)),
				resistedCrits: sum(actions.map(a => a.data.resistedCrits)),
				ticks: sum(actions.map(a => a.data.ticks)),
				resistedTicks: sum(actions.map(a => a.data.resistedTicks)),
				critTicks: sum(actions.map(a => a.data.critTicks)),
				resistedCritTicks: sum(actions.map(a => a.data.resistedCritTicks)),
				misses: sum(actions.map(a => a.data.misses)),
				dodges: sum(actions.map(a => a.data.dodges)),
				parries: sum(actions.map(a => a.data.parries)),
				blocks: sum(actions.map(a => a.data.blocks)),
				crushes: sum(actions.map(a => a.data.crushes)),
				blockedCrits: sum(actions.map(a => a.data.blockedCrits)),
				glances: sum(actions.map(a => a.data.glances)),
				damage: sum(actions.map(a => a.data.damage)),
				resistedDamage: sum(actions.map(a => a.data.resistedDamage)),
				critDamage: sum(actions.map(a => a.data.critDamage)),
				resistedCritDamage: sum(actions.map(a => a.data.resistedCritDamage)),
				tickDamage: sum(actions.map(a => a.data.tickDamage)),
				resistedTickDamage: sum(actions.map(a => a.data.resistedTickDamage)),
				critTickDamage: sum(actions.map(a => a.data.critTickDamage)),
				resistedCritTickDamage: sum(actions.map(a => a.data.resistedCritTickDamage)),
				glanceDamage: sum(actions.map(a => a.data.glanceDamage)),
				blockDamage: sum(actions.map(a => a.data.blockDamage)),
				blockedCritDamage: sum(actions.map(a => a.data.blockedCritDamage)),
				crushDamage: sum(actions.map(a => a.data.crushDamage)),
				threat: sum(actions.map(a => a.data.threat)),
				healing: sum(actions.map(a => a.data.healing)),
				critHealing: sum(actions.map(a => a.data.critHealing)),
				shielding: sum(actions.map(a => a.data.shielding)),
				castTimeMs: sum(actions.map(a => a.data.castTimeMs)),
			}),
			{
				iterations,
				duration,
			},
		);
	}
}
