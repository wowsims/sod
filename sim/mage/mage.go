package mage

import (
	"github.com/wowsims/sod/sim/common/guardians"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

const (
	SpellFlagMage       = core.SpellFlagAgentReserved1
	SpellFlagChillSpell = core.SpellFlagAgentReserved2
)

var TalentTreeSizes = [3]int{16, 16, 17}

func RegisterMage() {
	core.RegisterAgentFactory(
		proto.Player_Mage{},
		proto.Spec_SpecMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Mage)
			if !ok {
				panic("Invalid spec value for Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

const (
	SpellCode_MageNone int32 = iota
	SpellCode_MageArcaneBarrage
	SpellCode_MageArcaneBlast
	SpellCode_MageArcaneExplosion
	SpellCode_MageArcaneMissiles
	SpellCode_MageArcaneSurge
	SpellCode_MageBalefireBolt
	SpellCode_MageFireball
	SpellCode_MageFireBlast
	SpellCode_MageFrostbolt
	SpellCode_MageFrostfireBolt
	SpellCode_MageFrozenOrb
	SpellCode_MageLivingBomb
	SpellCode_MageLivingBombExplosion
	SpellCode_MageLivingFlame
	SpellCode_MageScorch
	SpellCode_MageSpellfrostBolt
)

type Mage struct {
	core.Character

	Talents *proto.MageTalents
	Options *proto.Mage_Options

	frozenOrb *FrozenOrb

	ArcaneBarrage           *core.Spell
	ArcaneBlast             *core.Spell
	ArcaneExplosion         []*core.Spell
	ArcaneMissiles          []*core.Spell
	ArcaneMissilesTickSpell []*core.Spell
	ArcaneSurge             *core.Spell
	BalefireBolt            *core.Spell
	BlastWave               []*core.Spell
	Blizzard                []*core.Spell
	DeepFreeze              *core.Spell
	Fireball                []*core.Spell
	FireBlast               []*core.Spell
	Flamestrike             []*core.Spell
	Frostbolt               []*core.Spell
	FrostfireBolt           *core.Spell
	FrozenOrb               *core.Spell
	IceLance                *core.Spell
	Ignite                  *core.Spell
	LivingBomb              *core.Spell
	LivingFlame             *core.Spell
	ManaGem                 []*core.Spell
	Pyroblast               []*core.Spell
	Scorch                  []*core.Spell
	SpellfrostBolt          *core.Spell

	IcyVeins *core.Spell

	ArcaneBlastAura     *core.Aura
	ArcanePotencyAura   *core.Aura
	ArcanePowerAura     *core.Aura
	ClearcastingAura    *core.Aura
	CombustionAura      *core.Aura
	FingersOfFrostAura  *core.Aura
	HotStreakAura       *core.Aura
	IceArmorAura        *core.Aura
	ImprovedScorchAuras core.AuraArray
	MageArmorAura       *core.Aura
	MissileBarrageAura  *core.Aura
	MoltenArmorAura     *core.Aura

	FingersOfFrostProcChance float64
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func (mage *Mage) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ArcaneBrilliance = true
}
func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) Initialize() {
	mage.registerArcaneMissilesSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFrostboltSpell()
	mage.registerPyroblastSpell()
	mage.registerScorchSpell()

	mage.registerArcaneExplosionSpell()
	mage.registerBlastWaveSpell()
	mage.registerBlizzardSpell()
	mage.registerFlamestrikeSpell()

	mage.registerEvocationCD()
	mage.registerManaGemCD()
}

func (mage *Mage) Reset(sim *core.Simulation) {
}

func NewMage(character *core.Character, options *proto.Player) *Mage {
	mageOptions := options.GetMage()

	mage := &Mage{
		Character: *character,
		Talents:   &proto.MageTalents{},
		Options:   mageOptions.Options,
	}
	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	mage.EnableManaBar()

	mage.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	mage.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[mage.Class][int(mage.Level)]*core.SpellCritRatingPerCritChance)

	switch mage.Options.Armor {
	case proto.Mage_Options_IceArmor:
		mage.applyFrostIceArmor()
	case proto.Mage_Options_MageArmor:
		mage.applyMageArmor()
	case proto.Mage_Options_MoltenArmor:
		mage.applyMoltenArmor()
	}

	// Set mana regen to 12.5 + Spirit/4 each 2s tick
	mage.SpiritManaRegenPerSecond = func() float64 {
		return 6.25 + mage.GetStat(stats.Spirit)/8
	}

	if mage.HasRune(proto.MageRune_RuneCloakFrozenOrb) {
		mage.frozenOrb = mage.NewFrozenOrb()
	}

	guardians.ConstructGuardians(&mage.Character)

	return mage
}

func (mage *Mage) HasRune(rune proto.MageRune) bool {
	return mage.HasRuneById(int32(rune))
}

func (mage *Mage) baseRuneAbilityDamage() float64 {
	return 13.828124 + 0.018012*float64(mage.Level) + 0.044141*float64(mage.Level*mage.Level)
}
