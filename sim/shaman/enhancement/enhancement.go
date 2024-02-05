package enhancement

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/shaman"
)

func RegisterEnhancementShaman() {
	core.RegisterAgentFactory(
		proto.Player_EnhancementShaman{},
		proto.Spec_SpecEnhancementShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewEnhancementShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_EnhancementShaman)
			if !ok {
				panic("Invalid spec value for Enhancement Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewEnhancementShaman(character *core.Character, options *proto.Player) *EnhancementShaman {
	enhOptions := options.GetEnhancementShaman()

	selfBuffs := shaman.SelfBuffs{
		Shield:  enhOptions.Options.Shield,
		ImbueMH: enhOptions.Options.ImbueMh,
		ImbueOH: enhOptions.Options.ImbueOh,
	}

	totems := &proto.ShamanTotems{}
	if enhOptions.Options.Totems != nil {
		totems = enhOptions.Options.Totems
	}

	enh := &EnhancementShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs),
	}

	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultMeleeCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	enh.ApplySyncType(enhOptions.Options.SyncType)

	return enh
}

type EnhancementShaman struct {
	*shaman.Shaman
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)
}

func (enh *EnhancementShaman) AutoSyncWeapons() proto.ShamanSyncType {
	if mh, oh := enh.MainHand(), enh.OffHand(); mh.SwingSpeed != oh.SwingSpeed {
		return proto.ShamanSyncType_NoSync
	}
	return proto.ShamanSyncType_SyncMainhandOffhandSwings
}

func (enh *EnhancementShaman) ApplySyncType(syncType proto.ShamanSyncType) {
	const FlurryICD = time.Millisecond * 500

	if syncType == proto.ShamanSyncType_Auto {
		syncType = enh.AutoSyncWeapons()
	}

	switch syncType {
	case proto.ShamanSyncType_SyncMainhandOffhandSwings:
		enh.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed(); nextMHSwingAt > aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		})
	case proto.ShamanSyncType_DelayOffhandSwings:
		enh.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed() + 100*time.Millisecond; nextMHSwingAt > aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		})
	default:
		enh.AutoAttacks.SetReplaceMHSwing(nil)
	}
}
