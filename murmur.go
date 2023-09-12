// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
)

func Hello() string {
	return fmt.Sprintf("Hello murmur v%s", Version)
}
