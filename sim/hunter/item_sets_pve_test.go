package hunter

import (
	"testing"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func TestRangedItemSets(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 4, 60, "../../ui/hunter/builds/instances/molten_core", "ranged", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 5, 60, "../../ui/hunter/builds/instances/blackwing_lair", "ranged", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 6, 60, "../../ui/hunter/builds/instances/temple_of_ahnqiraj", "ranged", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 7, 60, "../../ui/hunter/builds/instances/naxxramas", "ranged", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}

func Test2hItemSets(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 4, 60, "../../ui/hunter/builds/instances/molten_core", "2hmelee", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 5, 60, "../../ui/hunter/builds/instances/blackwing_lair", "2hmelee", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 6, 60, "../../ui/hunter/builds/instances/temple_of_ahnqiraj", "2hmelee", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 7, 60, "../../ui/hunter/builds/instances/naxxramas", "2hmelee", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}

func TestDwItemSets(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 4, 60, "../../ui/hunter/builds/instances/molten_core", "dwbmmelee", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 5, 60, "../../ui/hunter/builds/instances/zulgurub", "dwbmmelee", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 6, 60, "../../ui/hunter/builds/instances/ruins_of_ahnqiraj", "dwbmmelee", ItemFilters, proto.Stat_StatAgility, Stats),
		core.GetTestBuildFromJSON(proto.Class_ClassHunter, 7, 60, "../../ui/hunter/builds/instances/naxxramas", "dwbmmelee", ItemFilters, proto.Stat_StatAgility, Stats),
	}))
}