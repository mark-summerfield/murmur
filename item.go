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
func (me Label) String() string { return string(me) + ":" }

type CopyCommand string

func (me CopyCommand) Value() int     { return cmdCopy }
func (me CopyCommand) String() string { return "C" }

type JumpCommand string

func (me JumpCommand) Value() int     { return cmdJump }
func (me JumpCommand) String() string { return "J" }

type SuccCommand string

func (me SuccCommand) Value() int     { return cmdSucc }
func (me SuccCommand) String() string { return "S" }

type ZeroCommand string

func (me ZeroCommand) Value() int     { return cmdZero }
func (me ZeroCommand) String() string { return "Z" }
