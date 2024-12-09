package version

import (
	"fmt"
)

var (
	Release string
	Commit  string
)

// Version is the specification version that the package types support.
var Version = fmt.Sprintf("%s_%s",
	Release, Commit)
