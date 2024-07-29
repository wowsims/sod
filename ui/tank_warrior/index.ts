import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import { Sim } from '../core/sim.js';
import { TypedEvent } from '../core/typed_event.js';
import { TankWarriorSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecTankWarrior>(Spec.SpecTankWarrior, sim);

sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new TankWarriorSimUI(document.body, player);
