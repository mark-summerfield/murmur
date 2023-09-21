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
	Arity() int
	Lino() int
}

type Value int

func NewValue(v int) Value      { return Value(v) }
func (me Value) Value() int     { return int(me) }
func (me Value) String() string { return strconv.Itoa(int(me)) }

type Label string

func NewLabel(label string) (Label, error) {
	if rx := regexp.MustCompile(`^\pL\w*$`); !rx.MatchString(label) {
		return Label(""), fmt.Errorf("%w: %s", Err102, label)
	}
	return Label(label), nil
}
func (me Label) Value() int     { return -222 }
func (me Label) String() string { return string(me) }

type Address string

func NewAddress(label string) (Address, error) {
	if rx := regexp.MustCompile(`^@\pL\w*$`); !rx.MatchString(label) {
		return Address(""), fmt.Errorf("%w: %s", Err102, label)
	}
	return Address(label), nil
}
func (me Address) Value() int     { return -111 }
func (me Address) String() string { return string(me) }

func NewCommand(lino int, name byte) (Commander, error) {
	switch name {
	case '+':
		return NewAddCommand(lino), nil
	case 'C', 'c', 'T', 't':
		return NewCopyCommand(lino), nil
	case '/':
		return NewDivCommand(lino), nil
	case 'J', 'j':
		return NewJumpCommand(lino), nil
	case 'G', 'g':
		return NewJumpGtCommand(lino), nil
	case 'L', 'l':
		return NewJumpLtCommand(lino), nil
	case '*':
		return NewMulCommand(lino), nil
	case 'P', 'p', 'D', 'd':
		return NewPredCommand(lino), nil
	case '-':
		return NewSubCommand(lino), nil
	case 'S', 's', 'I', 'i':
		return NewSuccCommand(lino), nil
	case 'Z', 'z':
		return NewZeroCommand(lino), nil
	}
	return NewZeroCommand(lino), fmt.Errorf("%w: %c", Err117, name)
}

type CopyCommand int

func NewCopyCommand(lino int) CopyCommand { return CopyCommand(lino) }
func (me CopyCommand) Value() int         { return cmdCopy }
func (me CopyCommand) String() string     { return "C(" }
func (me CopyCommand) Arity() int         { return 2 }
func (me CopyCommand) Lino() int          { return int(me) }

type JumpCommand int

func NewJumpCommand(lino int) JumpCommand { return JumpCommand(lino) }
func (me JumpCommand) Value() int         { return cmdJump }
func (me JumpCommand) String() string     { return "J(" }
func (me JumpCommand) Arity() int         { return 3 }
func (me JumpCommand) Lino() int          { return int(me) }

type JumpGtCommand int

func NewJumpGtCommand(lino int) JumpGtCommand { return JumpGtCommand(lino) }
func (me JumpGtCommand) Value() int           { return cmdJumpGt }
func (me JumpGtCommand) String() string       { return "G(" }
func (me JumpGtCommand) Arity() int           { return 3 }
func (me JumpGtCommand) Lino() int            { return int(me) }

type JumpLtCommand int

func NewJumpLtCommand(lino int) JumpLtCommand { return JumpLtCommand(lino) }
func (me JumpLtCommand) Value() int           { return cmdJumpLt }
func (me JumpLtCommand) String() string       { return "L(" }
func (me JumpLtCommand) Arity() int           { return 3 }
func (me JumpLtCommand) Lino() int            { return int(me) }

type SuccCommand int

func NewSuccCommand(lino int) SuccCommand { return SuccCommand(lino) }
func (me SuccCommand) Value() int         { return cmdSucc }
func (me SuccCommand) String() string     { return "S(" }
func (me SuccCommand) Arity() int         { return 1 }
func (me SuccCommand) Lino() int          { return int(me) }

type ZeroCommand int

func NewZeroCommand(lino int) ZeroCommand { return ZeroCommand(lino) }
func (me ZeroCommand) Value() int         { return cmdZero }
func (me ZeroCommand) String() string     { return "Z(" }
func (me ZeroCommand) Arity() int         { return 1 }
func (me ZeroCommand) Lino() int          { return int(me) }

type PredCommand int

func NewPredCommand(lino int) PredCommand { return PredCommand(lino) }
func (me PredCommand) Value() int         { return cmdPred }
func (me PredCommand) String() string     { return "P(" }
func (me PredCommand) Arity() int         { return 1 }
func (me PredCommand) Lino() int          { return int(me) }

type AddCommand int

func NewAddCommand(lino int) AddCommand { return AddCommand(lino) }
func (me AddCommand) Value() int        { return cmdAdd }
func (me AddCommand) String() string    { return "+(" }
func (me AddCommand) Arity() int        { return 2 }
func (me AddCommand) Lino() int         { return int(me) }

type SubCommand int

func NewSubCommand(lino int) SubCommand { return SubCommand(lino) }
func (me SubCommand) Value() int        { return cmdSub }
func (me SubCommand) String() string    { return "-(" }
func (me SubCommand) Arity() int        { return 2 }
func (me SubCommand) Lino() int         { return int(me) }

type MulCommand int

func NewMulCommand(lino int) MulCommand { return MulCommand(lino) }
func (me MulCommand) Value() int        { return cmdMul }
func (me MulCommand) String() string    { return "*(" }
func (me MulCommand) Arity() int        { return 2 }
func (me MulCommand) Lino() int         { return int(me) }

type DivCommand int

func NewDivCommand(lino int) DivCommand { return DivCommand(lino) }
func (me DivCommand) Value() int        { return cmdDiv }
func (me DivCommand) String() string    { return "/(" }
func (me DivCommand) Arity() int        { return 2 }
func (me DivCommand) Lino() int         { return int(me) }
