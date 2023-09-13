// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
	"regexp"
	"strconv"
)

type Reger interface {
	Value() int
	String() string
}

type Commander interface {
	Reger
	IsCommand()
}

type Integer interface {
	int8 | int16 | int32 | int64 | int |
		uint8 | uint16 | uint32 | uint64 | uint
}

type Value int

func NewValue[T Integer](v T) Value { return Value(v) }
func (me Value) Value() int         { return int(me) }
func (me Value) String() string     { return strconv.Itoa(int(me)) }

type Label string

func NewLabel(label string) (Label, error) {
	if rx := regexp.MustCompile(`^\pL\w*$`); !rx.MatchString(label) {
		return Label(""), fmt.Errorf("%w: %s", Err102, label)
	}
	return Label(label), nil
}
func (me Label) Value() int     { return -222 }
func (me Label) String() string { return string(me) }

func NewCommand(name byte) (Commander, error) {
	switch name {
	case 'C', 'T', 'c', 't':
		return NewCopyCommand(), nil
	case 'J', 'j':
		return NewJumpCommand(), nil
	case 'S', 's':
		return NewSuccCommand(), nil
	case 'Z', 'z':
		return NewZeroCommand(), nil
	}
	return NewZeroCommand(), fmt.Errorf("%w: %c", Err117, name)
}

type CopyCommand struct{}

func NewCopyCommand() CopyCommand     { return CopyCommand(struct{}{}) }
func (me CopyCommand) Value() int     { return cmdCopy }
func (me CopyCommand) String() string { return "C(" }
func (me CopyCommand) IsCommand()     {}

type JumpCommand struct{}

func NewJumpCommand() JumpCommand     { return JumpCommand(struct{}{}) }
func (me JumpCommand) Value() int     { return cmdJump }
func (me JumpCommand) String() string { return "J(" }
func (me JumpCommand) IsCommand()     {}

type SuccCommand struct{}

func NewSuccCommand() SuccCommand     { return SuccCommand(struct{}{}) }
func (me SuccCommand) Value() int     { return cmdSucc }
func (me SuccCommand) String() string { return "S(" }
func (me SuccCommand) IsCommand()     {}

type ZeroCommand struct{}

func NewZeroCommand() ZeroCommand     { return ZeroCommand(struct{}{}) }
func (me ZeroCommand) Value() int     { return cmdZero }
func (me ZeroCommand) String() string { return "Z(" }
func (me ZeroCommand) IsCommand()     {}
