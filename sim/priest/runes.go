package priest

func (priest *Priest) ApplyRunes() {
	// Chest
	// priest.registerSerendipity() // TODO
	// priest.registerStrengthOfSoul() // TODO
	priest.registerVoidPlagueSpell()
	// priest.registerTwistedFaith() // Nothing to do

	// Hands
	// priest.registerCircleOfHealingSpell() // TODO
	priest.registerMindSearSpell()
	priest.RegisterPenanceSpell()
	priest.registerShadowWordDeathSpell()

	// Belt
	// priest.registerEmpoweredRenew() // TODO
	priest.registerMindSpikeSpell()
	// priest.registerRenewedHope // TODO

	// Legs
	priest.registerHomunculiSpell()
	// priest.registerPowerWordBarrierSpell() // TODO
	// priest.registerPrayerOfMendingSpell() // TODO
	// priest.registerSharedPainSpell() // Nothing to do

	// Feet
	priest.registerDispersionSpell()
	// priest.registerPainSuppressionSpell() // TODO
	// priest.registerSpiritOfTheRedeemerSpell() // TODO

	// Skill Books
	priest.registerShadowfiendSpell()
}
