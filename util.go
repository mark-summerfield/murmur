// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package murmur

import (
	"strings"
)

func cleanLine(line string) string {
	line = strings.TrimSpace(line)
	i := strings.IndexByte(line, ';')
	if i > -1 {
		line = strings.TrimSpace(line[:i])
	}
	return line
}
