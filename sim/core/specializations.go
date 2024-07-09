package core

///////////////////////////////////////////////////////////////////////////
//                            Weapon Specialization Auras
///////////////////////////////////////////////////////////////////////////

func (character *Character) SwordSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Sword Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SwordsSkill += 5
			character.PseudoStats.TwoHandedSwordsSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.SwordsSkill -= 5
			character.PseudoStats.TwoHandedSwordsSkill -= 5
		},
	})
}

func (character *Character) AxeSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Axe Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.AxesSkill += 5
			character.PseudoStats.TwoHandedAxesSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.AxesSkill -= 5
			character.PseudoStats.TwoHandedAxesSkill -= 5
		},
	})
}

func (character *Character) MaceSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Mace Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.MacesSkill += 5
			character.PseudoStats.TwoHandedMacesSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.MacesSkill -= 5
			character.PseudoStats.TwoHandedMacesSkill -= 5
		},
	})
}

func (character *Character) DaggerSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Dagger Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DaggersSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DaggersSkill -= 5
		},
	})
}

func (character *Character) FistWeaponSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Fist Weapon Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.UnarmedSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.UnarmedSkill -= 5
		},
	})
}

func (character *Character) PoleWeaponSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Pole Weapon Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.StavesSkill += 5
			character.PseudoStats.PolearmsSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.StavesSkill -= 5
			character.PseudoStats.PolearmsSkill -= 5
		},
	})
}

func (character *Character) GunSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Gun Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.GunsSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.GunsSkill -= 5
		},
	})
}

func (character *Character) BowSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Bow Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.BowsSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.BowsSkill -= 5
		},
	})
}

func (character *Character) CrossbowSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Crossbow Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.CrossbowsSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.CrossbowsSkill -= 5
		},
	})
}

func (character *Character) ThrownSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Thrown Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.ThrownSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.ThrownSkill -= 5
		},
	})
}

func (character *Character) FeralCombatSpecializationAura() *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:      "Feral Combat Skill Specialization",
		BuildPhase: CharacterBuildPhaseBuffs,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.FeralCombatSkill += 5
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.FeralCombatSkill -= 5
		},
	})
}
