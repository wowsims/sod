package common

// Just import other directories, so importing common from elsewhere is enough.
import (
	_ "github.com/wowsims/sod/sim/common/vanilla"
	_ "github.com/wowsims/sod/sim/common/vanilla/item_sets"

	_ "github.com/wowsims/sod/sim/common/sod"
	_ "github.com/wowsims/sod/sim/common/sod/crafted"
	_ "github.com/wowsims/sod/sim/common/sod/item_effects"
	_ "github.com/wowsims/sod/sim/common/sod/items_sets"
)
