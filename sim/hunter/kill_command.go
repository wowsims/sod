package hunter

// TODO: 2024-06-13 - Rune changed from Kill Command to Kill Shot. Unsure if Kill Command still in the game.
// func (hunter *Hunter) registerKillCommand() {
// 	if hunter.pet == nil || !hunter.HasRune(proto.HunterRune_RuneLegsKillCommand) {
// 		return
// 	}

// 	actionID := core.ActionID{SpellID: 409379}
// 	hasCatlikeReflexes := hunter.HasRune(proto.HunterRune_RuneHelmCatlikeReflexes)

// 	cooldownModifier := 1.0
// 	if hasCatlikeReflexes {
// 		cooldownModifier *= 0.5
// 	}

// 	// For tracking in timeline
// 	hunterAura := hunter.RegisterAura(core.Aura{
// 		Label:     "Kill Command",
// 		ActionID:  actionID,
// 		Duration:  time.Second * 30,
// 		MaxStacks: 3,
// 	})

// 	hunter.pet.killCommandAura = hunter.pet.RegisterAura(core.Aura{
// 		Label:     "Kill Command",
// 		ActionID:  actionID,
// 		Duration:  time.Second * 30,
// 		MaxStacks: 3,
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			// TODO: Make it only work on Claw/Bite after pet abilities refactor
// 			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
// 				aura.RemoveStack(sim)
// 				hunterAura.RemoveStack(sim)
// 			}
// 		},
// 	})

// 	hunter.KillCommand = hunter.RegisterSpell(core.SpellConfig{
// 		ActionID:    actionID,
// 		SpellSchool: core.SpellSchoolPhysical,
// 		Flags:       core.SpellFlagNoOnCastComplete,

// 		ManaCost: core.ManaCostOptions{
// 			BaseCost: 0.015,
// 		},
// 		Cast: core.CastConfig{
// 			CD: core.Cooldown{
// 				Timer:    hunter.NewTimer(),
// 				Duration: time.Second * time.Duration(60*cooldownModifier),
// 			},
// 		},
// 		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
// 			return hunter.pet.IsEnabled()
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			hunter.pet.killCommandAura.Activate(sim)
// 			hunter.pet.killCommandAura.SetStacks(sim, 3)

// 			hunterAura.Activate(sim)
// 			hunterAura.SetStacks(sim, 3)
// 		},
// 	})

// 	hunter.AddMajorCooldown(core.MajorCooldown{
// 		Spell: hunter.KillCommand,
// 		Type:  core.CooldownTypeDPS,
// 	})
// }
