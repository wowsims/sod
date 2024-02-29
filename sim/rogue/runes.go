package rogue

import "github.com/wowsims/sod/sim/core/proto"

func (rogue *Rogue) ApplyRunes() {
	// Apply runes here :)
	if rogue.HasRune(proto.RogueRune_RuneDeadlyBrew) {
		rogue.applyDeadlyBrewInstant()
		rogue.applyDeadlyBrewDeadly()
	}

	rogue.registerWaylayAura()
	rogue.registerMasterOfSubtlety()
	rogue.registerMainGaucheSpell()
	rogue.registerSaberSlashSpell()
	rogue.registerShivSpell()
	rogue.registerShadowstrikeSpell()
	rogue.registerMutilateSpell()
	rogue.registerEnvenom()
}
