// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	_ "embed"
	"fmt"
)

//go:embed Version.dat
var Version string

func Hello() string {
	return fmt.Sprintf("Hello murmur v%s", Version)
}
