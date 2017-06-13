package kagami

import (
	"fmt"
)

// Version is the current version of the software
const Version = "0.1.0"

// VersionStability is the status of this build, either stable, rc or dev
const VersionStability = "dev"

// VersionString returns a string in the form version-stability
func VersionString() string {
	if VersionStability != "stable" {
		return fmt.Sprintf("%s-%s", Version, VersionStability)
	}
	return Version
}
