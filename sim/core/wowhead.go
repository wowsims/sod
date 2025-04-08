package core

import "fmt"

// TODO: Wowhead switched to using a dedicated PTR branch which means we need to use it during PTR cycles
const WowheadBranch = "classic"
const WowheadBranchPTR = "classic-ptr"

func MakeWowheadUrl(path string) string {
	if path[0] != '/' {
		path = "/" + path
	}

	return fmt.Sprintf("https://nether.wowhead.com/%s%s", WowheadBranch, path)
}
