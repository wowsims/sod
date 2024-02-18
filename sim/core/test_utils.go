package core

import (
	"log"
	"os"
	"testing"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
	googleProto "google.golang.org/protobuf/proto"
)

var DefaultSimTestOptions = &proto.SimOptions{
	Iterations: 20,
	IsTest:     true,
	Debug:      false,
	RandomSeed: 101,
}
var StatWeightsDefaultSimTestOptions = &proto.SimOptions{
	Iterations: 300,
	IsTest:     true,
	Debug:      false,
	RandomSeed: 101,
}
var AverageDefaultSimTestOptions = &proto.SimOptions{
	Iterations: 2000,
	IsTest:     true,
	Debug:      false,
	RandomSeed: 101,
}

const ShortDuration = 60
const LongDuration = 300

var DefaultTargetProtoLvl25 = &proto.Target{
	Level: 27,
	Stats: stats.Stats{
		stats.Armor:       1104,
		stats.AttackPower: 320,
	}.ToFloatArray(),
	MobType: proto.MobType_MobTypeDemon,

	SwingSpeed:    2,
	MinBaseDamage: 4192.05,
	ParryHaste:    true,
	DamageSpread:  0.3333,
}

var DefaultTargetProtoLvl40 = &proto.Target{
	Level: 42,
	Stats: stats.Stats{
		stats.Armor:       1104,
		stats.AttackPower: 320,
	}.ToFloatArray(),
	MobType: proto.MobType_MobTypeDemon,

	SwingSpeed:    2,
	MinBaseDamage: 4192.05,
	ParryHaste:    true,
	DamageSpread:  0.3333,
}

///////////////////////////////////////////////////////////////////////////
//                                 Raid Buffs
///////////////////////////////////////////////////////////////////////////

var FullRaidBuffsPhase1 = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	AspectOfTheLion:      true,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	BloodPact:            proto.TristateEffect_TristateEffectImproved,
	DevotionAura:         proto.TristateEffect_TristateEffectImproved,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectImproved,
	PowerWordFortitude:   proto.TristateEffect_TristateEffectImproved,
	RetributionAura:      proto.TristateEffect_TristateEffectImproved,
	StrengthOfEarthTotem: proto.TristateEffect_TristateEffectImproved,
	Thorns:               proto.TristateEffect_TristateEffectImproved,
}

var FullRaidBuffsPhase2 = &proto.RaidBuffs{
	ArcaneBrilliance:      true,
	AspectOfTheLion:       true,
	BattleShout:           proto.TristateEffect_TristateEffectImproved,
	BloodPact:             proto.TristateEffect_TristateEffectImproved,
	DevotionAura:          proto.TristateEffect_TristateEffectImproved,
	DivineSpirit:          true,
	FireResistanceAura:    true,
	FireResistanceTotem:   true,
	FrostResistanceAura:   true,
	FrostResistanceTotem:  true,
	GiftOfTheWild:         proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:       true,
	ManaSpringTotem:       proto.TristateEffect_TristateEffectImproved,
	MoonkinAura:           true,
	NatureResistanceTotem: true,
	PowerWordFortitude:    proto.TristateEffect_TristateEffectImproved,
	RetributionAura:       proto.TristateEffect_TristateEffectImproved,
	ShadowProtection:      true,
	StrengthOfEarthTotem:  proto.TristateEffect_TristateEffectImproved,
	Thorns:                proto.TristateEffect_TristateEffectImproved,
	TrueshotAura:          true,
}

///////////////////////////////////////////////////////////////////////////
//                                 Party Buffs
///////////////////////////////////////////////////////////////////////////

var FullPartyBuffs = &proto.PartyBuffs{}

///////////////////////////////////////////////////////////////////////////
//                                 Individual Buffs
///////////////////////////////////////////////////////////////////////////

var FullIndividualBuffsPhase1 = &proto.IndividualBuffs{
	AshenvalePvpBuff:  true,
	BlessingOfKings:   true,
	BlessingOfMight:   proto.TristateEffect_TristateEffectImproved,
	BlessingOfWisdom:  proto.TristateEffect_TristateEffectImproved,
	BoonOfBlackfathom: true,
	SaygesFortune:     proto.SaygesFortune_SaygesDamage,
}

var FullIndividualBuffsPhase2 = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfMight:     proto.TristateEffect_TristateEffectImproved,
	BlessingOfSanctuary: true,
	BlessingOfWisdom:    proto.TristateEffect_TristateEffectImproved,
	SaygesFortune:       proto.SaygesFortune_SaygesDamage,
	SparkOfInspiration:  true,
}

///////////////////////////////////////////////////////////////////////////
//                                 Debuffs
///////////////////////////////////////////////////////////////////////////

var FullDebuffsPhase1 = &proto.Debuffs{
	CurseOfElements:      true,
	CurseOfRecklessness:  true,
	CurseOfVulnerability: true,
	CurseOfWeakness:      proto.TristateEffect_TristateEffectImproved,
	DemoralizingRoar:     proto.TristateEffect_TristateEffectImproved,
	DemoralizingShout:    proto.TristateEffect_TristateEffectImproved,
	Dreamstate:           true,
	ExposeArmor:          proto.TristateEffect_TristateEffectImproved,
	FaerieFire:           true,
	InsectSwarm:          true,
	ImprovedShadowBolt:   true,
	ScorpidSting:         true,
	SunderArmor:          true,
	ThunderClap:          proto.TristateEffect_TristateEffectImproved,
}

var FullDebuffsPhase2 = &proto.Debuffs{
	CurseOfElements:      true,
	CurseOfRecklessness:  true,
	CurseOfVulnerability: true,
	CurseOfWeakness:      proto.TristateEffect_TristateEffectImproved,
	DemoralizingRoar:     proto.TristateEffect_TristateEffectImproved,
	DemoralizingShout:    proto.TristateEffect_TristateEffectImproved,
	Dreamstate:           true,
	ExposeArmor:          proto.TristateEffect_TristateEffectImproved,
	FaerieFire:           true,
	InsectSwarm:          true,
	ImprovedScorch:       true,
	ImprovedShadowBolt:   true,
	JudgementOfLight:     true,
	JudgementOfWisdom:    true,
	ScorpidSting:         true,
	ShadowWeaving:        true,
	Stormstrike:          true,
	SunderArmor:          true,
	ThunderClap:          proto.TristateEffect_TristateEffectImproved,
	WintersChill:         true,
}

///////////////////////////////////////////////////////////////////////////
//                                 Full Buffs
///////////////////////////////////////////////////////////////////////////

var FullBuffsPhase1 = BuffsCombo{
	Label: "Phase 1 Buffs",

	Debuffs: FullDebuffsPhase1,
	Party:   FullPartyBuffs,
	Player:  FullIndividualBuffsPhase1,
	Raid:    FullRaidBuffsPhase1,
}

var FullBuffsPhase2 = BuffsCombo{
	Label: "Phase 2 Buffs",

	Debuffs: FullDebuffsPhase2,
	Party:   FullPartyBuffs,
	Player:  FullIndividualBuffsPhase2,
	Raid:    FullRaidBuffsPhase2,
}

func NewDefaultTarget(playerLevel int32) *proto.Target {
	switch playerLevel {
	case 40:
		return DefaultTargetProtoLvl40
	default:
		return DefaultTargetProtoLvl25
	}
}

func MakeDefaultEncounterCombos(playerLevel int32) []EncounterCombo {
	var DefaultTarget = NewDefaultTarget(playerLevel)

	multipleTargets := make([]*proto.Target, 20)
	for i := range multipleTargets {
		multipleTargets[i] = DefaultTarget
	}

	return []EncounterCombo{
		{
			Label: "ShortSingleTarget",
			Encounter: &proto.Encounter{
				Duration:             ShortDuration,
				ExecuteProportion_20: 0.2,
				ExecuteProportion_25: 0.25,
				ExecuteProportion_35: 0.35,
				Targets: []*proto.Target{
					DefaultTarget,
				},
			},
		},
		{
			Label: "LongSingleTarget",
			Encounter: &proto.Encounter{
				Duration:             LongDuration,
				ExecuteProportion_20: 0.2,
				ExecuteProportion_25: 0.25,
				ExecuteProportion_35: 0.35,
				Targets: []*proto.Target{
					DefaultTarget,
				},
			},
		},
		{
			Label: "LongMultiTarget",
			Encounter: &proto.Encounter{
				Duration:             LongDuration,
				ExecuteProportion_20: 0.2,
				ExecuteProportion_25: 0.25,
				ExecuteProportion_35: 0.35,
				Targets:              multipleTargets,
			},
		},
	}
}

func MakeSingleTargetEncounter(playerLevel int32, variation float64) *proto.Encounter {
	return &proto.Encounter{
		Duration:             LongDuration,
		DurationVariation:    variation,
		ExecuteProportion_20: 0.2,
		ExecuteProportion_25: 0.25,
		ExecuteProportion_35: 0.35,
		Targets: []*proto.Target{
			NewDefaultTarget(playerLevel),
		},
	}
}

func CharacterStatsTest(label string, t *testing.T, raid *proto.Raid, expectedStats stats.Stats) {
	csr := &proto.ComputeStatsRequest{
		Raid: raid,
	}

	result := ComputeStats(csr)
	finalStats := stats.FromFloatArray(result.RaidStats.Parties[0].Players[0].FinalStats.Stats)

	const tolerance = 0.5
	if !finalStats.EqualsWithTolerance(expectedStats, tolerance) {
		t.Fatalf("%s failed: CharacterStats() = %v, expected %v", label, finalStats, expectedStats)
	}
}

func StatWeightsTest(label string, t *testing.T, _swr *proto.StatWeightsRequest, expectedStatWeights stats.Stats) {
	// Make a copy so we can safely change fields.
	swr := googleProto.Clone(_swr).(*proto.StatWeightsRequest)

	swr.Encounter.Duration = LongDuration
	swr.SimOptions.Iterations = 5000

	result := StatWeights(swr)
	resultWeights := stats.FromFloatArray(result.Dps.Weights.Stats)

	const tolerance = 0.05
	if !resultWeights.EqualsWithTolerance(expectedStatWeights, tolerance) {
		t.Fatalf("%s failed: CalcStatWeight() = %v, expected %v", label, resultWeights, expectedStatWeights)
	}
}

func RaidSimTest(label string, t *testing.T, rsr *proto.RaidSimRequest, expectedDps float64) {
	result := RunRaidSim(rsr)
	if result.ErrorResult != "" {
		t.Fatalf("Sim failed with error: %s", result.ErrorResult)
	}
	tolerance := 0.5
	if result.RaidMetrics.Dps.Avg < expectedDps-tolerance || result.RaidMetrics.Dps.Avg > expectedDps+tolerance {
		// Automatically print output if we had debugging enabled.
		if rsr.SimOptions.Debug {
			log.Printf("LOGS:\n%s\n", result.Logs)
		}
		t.Fatalf("%s failed: expected %0f dps from sim but was %0f", label, expectedDps, result.RaidMetrics.Dps.Avg)
	}
}

func RaidBenchmark(b *testing.B, rsr *proto.RaidSimRequest) {
	rsr.Encounter.Duration = LongDuration
	rsr.SimOptions.Iterations = 1

	// Set to false because IsTest adds a lot of computation.
	rsr.SimOptions.IsTest = false

	for i := 0; i < b.N; i++ {
		result := RunRaidSim(rsr)
		if result.ErrorResult != "" {
			b.Fatalf("RaidBenchmark() at iteration %d failed: %v", i, result.ErrorResult)
		}
	}
}

func GetAplRotation(dir string, file string) RotationCombo {
	filePath := dir + "/" + file + ".apl.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to load apl json file: %s, %s", filePath, err)
	}

	return RotationCombo{Label: file, Rotation: APLRotationFromJsonString(string(data))}
}

func GetGearSet(dir string, file string) GearSetCombo {
	filePath := dir + "/" + file + ".gear.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to load gear json file: %s, %s", filePath, err)
	}

	return GearSetCombo{Label: file, GearSet: EquipmentSpecFromJsonString(string(data))}
}
