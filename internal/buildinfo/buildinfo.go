package buildinfo

import "fmt"

var Version = "dev"
var Commit = "unknown"
var BuildDate = "unknown"

func String() string {
	if Version == "dev" && Commit == "unknown" && BuildDate == "unknown" {
		return "DevDuck dev"
	}
	return fmt.Sprintf("DevDuck %s\ncommit: %s\nbuilt: %s", Version, Commit, BuildDate)
}
