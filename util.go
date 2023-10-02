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

func fixups(text string) string {
	replacer := strings.NewReplacer("Z(0)", stopCmd, "Z(PC)", stopCmd,
		"J(0, 0, ", "J(", "J(1, 1, ", "J(",
		"G(0, 0, ", "G(", "G(1, 1, ", "G(",
		"L(0, 0, ", "L(", "L(1, 1, ", "L(",
		"C(0, 0)", noopCmd, "C(1, 1)", noopCmd,
	)
	return replacer.Replace(text)
}
