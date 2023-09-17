// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"errors"
	"fmt"
	"strconv"
)

type Urm struct {
	regs          []Reger
	reg_for_label map[string]int
	label_for_reg map[int]string
	steps         int
	labelWidth    int
}

// NewX creates a new Urm and populates its registers by parsing the given
// slice of strings. (See NewX.)
func New(lines []string) (*Urm, error) { return NewX(lines, DefaultSize) }

// NewX creates a new Urm and populates its registers by parsing the given
// slice of strings. The size is how many registers to use; pass 0 to have
// this either read from the lines (^size) or calculated during parsing.
//
// BNF:
//
//	PROGRAM ::= SIZE? (SETREG | COMMAND)+
//	SIZE ::= /\^\d+/
//	SETREG ::= /(\d+|\pL\w*)\s*:\s*\d+/ # if LHS is num it is reg else label
//	COMMAND ::= LABEL? /[CJSTZ][(]\d+(:?,\s*\d+)*/
//	LABEL ::= /^\pL\w*:\s*/
//
// Any line may end with a comment which begins with ';'. All comments and
// blank lines are ignored.
func NewX(lines []string, size int) (*Urm, error) {
	urm := Urm{}
	if err := urm.load(lines, size); err != nil {
		return nil, err
	}
	return &urm, nil
}

// Size returns the number of registers.
func (me *Urm) Size() int { return len(me.regs) }

// Steps returns the number of steps so far.
func (me *Urm) Steps() int { return me.steps }

// Run runs the Urm program for up to DefaultMaxSteps steps.
// (See RunX and Step.)
func (me *Urm) Run() error { return me.RunX(DefaultMaxSteps) }

// Run runs the Urm program for a maximum of maxSteps steps.
// (See Run and Step.)
func (me *Urm) RunX(maxSteps int) error {
	for {
		if me.steps >= maxSteps {
			return fmt.Errorf("%w: %d", Err100, me.steps)
		}
		err := me.Step()
		if errors.Is(err, ErrStop) {
			return nil // normal termination
		}
		if err != nil {
			return err // error termination
		}
	}
}

// Step runs the one step (i.e., executes the next statement) of the Urm
// program. (See Run and RunX.)
func (me *Urm) Step() error {
	me.steps++
	cmd, err := me.nextCommand() // returns cmd & inc PC
	if err != nil {
		return err
	}
	return me.executeCommand(cmd) // may acquire one or more operands
}

// RegAsString returns a string representation of a register's value.
func (me *Urm) RegAsString(reg int) (string, error) {
	if 0 <= reg && reg < len(me.regs) {
		reger := me.regs[reg]
		if reger == nil {
			return fmt.Sprintf("%d:%s", reg, null), nil
		}
		if label, ok := reger.(Label); ok {
			return fmt.Sprintf("%d:%s", reg, label), nil
		}
		register := ""
		if label, ok := me.label_for_reg[reg]; ok {
			register = label
		} else {
			register = strconv.Itoa(reg)
		}
		return fmt.Sprintf("%s:%s", register, reger), nil
	}
	return "", fmt.Errorf("%w: %d", Err115, reg)
}

// RegForLabelAsString returns a string representation of a register's value
// where the register is identified by a label.
func (me *Urm) RegForLabelAsString(label string) (string, error) {
	if reg, ok := me.reg_for_label[label]; ok {
		return me.RegAsString(reg)
	} else {
		return "", fmt.Errorf("%w: %q", Err116, label)
	}
}

// SetRegToValue sets the given register to the given value.
func (me *Urm) SetRegToValue(reg, value int) error {
	if 0 <= reg && reg < len(me.regs) {
		me.regs[reg] = NewValue(value)
		return nil
	}
	return fmt.Errorf("%w: %d", Err112, reg)
}

// String returns a string of all the registers.
// Errors are ignored. Mostly for debugging and testing.
func (me *Urm) String() string { return me.asString(false) }

// StringWithRegNums returns a string of all the registers and with register
// numbers in comments.
// Errors are ignored. Mostly for debugging and testing.
func (me *Urm) StringWithRegNums() string { return me.asString(true) }
