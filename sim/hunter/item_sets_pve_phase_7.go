package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetCryptstalkerProwess = core.NewItemSet(core.ItemSet{
	Name: "Cryptstalker Prowess",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasMelee2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasMelee4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasMelee6PBonus()
		},
	},
})

// Your Wyvern Strike and Mongoose Bite deal 20% more initial damage.
func (hunter *Hunter) applyNaxxramasMelee2PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Melee 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.WyvernStrike != nil {
				hunter.WyvernStrike.ImpactDamageMultiplierAdditive += 0.20
			}
			hunter.MongooseBite.ImpactDamageMultiplierAdditive += 0.20
		},
	})
}

// Reduces the cooldown on your Wyvern Strike ability by 2 sec, reduces the cooldown on your raptor strike ability by 1 sec, and reduces the cooldown on your Flanking Strike ability by 8sec.
func (hunter *Hunter) applyNaxxramasMelee4PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Melee 4P Bonus"
	if hunter.HasAura(label) {
		return
	}

	// Not entirely sure how this will work so taking some liberties
	// Assume that it resets all of them when one crits
	//var spellsToReset []*core.Spell

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.RaptorStrike != nil {
				hunter.RaptorStrike.CD.FlatModifier -= time.Second
			}
			if hunter.WyvernStrike != nil {
				hunter.WyvernStrike.CD.FlatModifier -= time.Second * 2
			}
			if hunter.FlankingStrike != nil {
				hunter.FlankingStrike.CD.FlatModifier -= time.Second * 8
			}
		},
	}))
}

// You gain 1% increased critical strike chance for 30 sec each time you hit an Undead enemy with a melee attack, stacking up to 35 times.
func (hunter *Hunter) applyNaxxramasMelee6PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Melee 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	buffAura := hunter.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1218587},
		Label:     "Critical Aim",
		Duration:  time.Second * 30,
		MaxStacks: 35,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			hunter.AddStatDynamic(sim, stats.MeleeCrit, float64(newStacks-oldStacks)*core.CritRatingPerCritChance)
			hunter.AddStatDynamic(sim, stats.SpellCrit, float64(newStacks-oldStacks)*core.CritRatingPerCritChance)
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target.MobType == proto.MobType_MobTypeUndead && spell.ProcMask.Matches(core.ProcMaskMelee) {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetCryptstalkerArmor = core.NewItemSet(core.ItemSet{
	Name: "Cryptstalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasRanged2PBonus()
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasRanged4PBonus()
		},
		6: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.applyNaxxramasRanged6PBonus()
		},
	},
})

// Your Serpent Sting deals 20% more damage.
func (hunter *Hunter) applyNaxxramasRanged2PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Ranged 2P Bonus"
	if hunter.HasAura(label) {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hunter.SerpentSting.DamageMultiplierAdditive += 0.20

			if hunter.SerpentStingChimeraShot != nil {
				hunter.SerpentStingChimeraShot.DamageMultiplierAdditive += 0.20
			}
		},
	})
}

// Reduces the cooldown on your Chimera Shot, Explosive Shot, and Aimed Shot abilities by 1.5 sec and reduces the cooldown on your Kill Shot ability by 3sec.
func (hunter *Hunter) applyNaxxramasRanged4PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Ranged 4P Bonus"
	if hunter.HasAura(label) {
		return
	}
	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.ChimeraShot != nil {
				hunter.ChimeraShot.CD.FlatModifier -= time.Millisecond * 1500
			}
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.CD.FlatModifier -= time.Millisecond * 1500 
			}
			if hunter.AimedShot != nil {
				hunter.AimedShot.CD.FlatModifier -= time.Millisecond * 1500 
			}
			if hunter.KillShot != nil {
				hunter.KillShot.CD.FlatModifier -= time.Second * 3 
			}
		},
	}))
}

// You gain 1% increased critical strike chance for 30 sec each time you hit an Undead enemy with a ranged attack, stacking up to 35 times.
func (hunter *Hunter) applyNaxxramasRanged6PBonus() {
	label := "S03 - Item - Naxxramas - Hunter - Ranged 6P Bonus"
	if hunter.HasAura(label) {
		return
	}

	buffAura := hunter.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 1218587},
		Label:     "Critical Aim",
		Duration:  time.Second * 30,
		MaxStacks: 35,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			hunter.AddStatDynamic(sim, stats.MeleeCrit, float64(newStacks-oldStacks)*core.CritRatingPerCritChance)
			hunter.AddStatDynamic(sim, stats.SpellCrit, float64(newStacks-oldStacks)*core.CritRatingPerCritChance)
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target.MobType == proto.MobType_MobTypeUndead && spell.ProcMask.Matches(core.ProcMaskRanged) {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}
