// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import "strconv"

type Item interface {
	Value() int
	String() string
}

type Value int

func (me Value) Value() int     { return int(me) }
func (me Value) String() string { return strconv.Itoa(int(me)) }

type Label string

func (me Label) Value() int     { return -222 }
func (me Label) String() string { return string(me) }
