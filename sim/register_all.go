package sim

import (
	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/druid/balance"
	"github.com/wowsims/sod/sim/paladin/retribution"
	dpsrogue "github.com/wowsims/sod/sim/rogue/dps_rogue"
	tankrogue "github.com/wowsims/sod/sim/rogue/tank_rogue"
	"github.com/wowsims/sod/sim/shaman/elemental"
	"github.com/wowsims/sod/sim/shaman/enhancement"
	"github.com/wowsims/sod/sim/shaman/warden"

	"github.com/wowsims/sod/sim/druid/feral"
	// restoDruid "github.com/wowsims/sod/sim/druid/restoration"
	// feralTank "github.com/wowsims/sod/sim/druid/tank"
	_ "github.com/wowsims/sod/sim/encounters"
	"github.com/wowsims/sod/sim/hunter"
	"github.com/wowsims/sod/sim/mage"

	// holyPaladin "github.com/wowsims/sod/sim/paladin/holy"
	"github.com/wowsims/sod/sim/paladin/protection"
	// "github.com/wowsims/sod/sim/paladin/retribution"
	// healingPriest "github.com/wowsims/sod/sim/priest/healing"
	"github.com/wowsims/sod/sim/priest/shadow"

	// restoShaman "github.com/wowsims/sod/sim/shaman/restoration"
	dpsWarlock "github.com/wowsims/sod/sim/warlock/dps"
	tankWarlock "github.com/wowsims/sod/sim/warlock/tank"
	dpsWarrior "github.com/wowsims/sod/sim/warrior/dps_warrior"
	tankWarrior "github.com/wowsims/sod/sim/warrior/tank_warrior"
)

var registered = false

func RegisterAll() {
	if registered {
		return
	}
	registered = true

	balance.RegisterBalanceDruid()
	feral.RegisterFeralDruid()
	// feralTank.RegisterFeralTankDruid()
	// restoDruid.RegisterRestorationDruid()
	elemental.RegisterElementalShaman()
	enhancement.RegisterEnhancementShaman()
	warden.RegisterWardenShaman()
	// restoShaman.RegisterRestorationShaman()
	hunter.RegisterHunter()
	mage.RegisterMage()
	// healingPriest.RegisterHealingPriest()
	shadow.RegisterShadowPriest()
	dpsrogue.RegisterDpsRogue()
	tankrogue.RegisterTankRogue()
	dpsWarrior.RegisterDpsWarrior()
	tankWarrior.RegisterTankWarrior()
	// holyPaladin.RegisterHolyPaladin()
	protection.RegisterProtectionPaladin()
	retribution.RegisterRetributionPaladin()
	dpsWarlock.RegisterDpsWarlock()
	tankWarlock.RegisterTankWarlock()
}
