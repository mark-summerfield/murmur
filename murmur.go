// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
	"strings"
)

type Urm struct {
	regs          []Item
	reg_for_label map[string]int
	label_for_reg map[int]string
	steps         int
}

// New returns a new Urm with DefaultSize registers. (See NewX.)
func New() *Urm { return NewX(DefaultSize) }

// NewX returns a new Urm with size registers.
func NewX(size int) *Urm {
	urm := Urm{}
	urm.clear(size)
	return &urm
}

// Size returns the number of registers.
func (me *Urm) Size() int { return len(me.regs) }

// Steps returns the number of steps so far.
func (me *Urm) Steps() int { return me.steps }

// Load clears the Urm and repopulates its registers by parsing the given
// slice of strings. The size is how many registers to use; pass 0 to have
// this either read from the lines (#size) or calculated during parsing.
//
// BNF:
//
//	PROGRAM ::= SIZE? (SETREG | COMMAND)+
//	SIZE ::= /#\d+/
//	SETREG ::= /(\d+|\pL\w*)\s*:\s*\d+/ # if LHS is num it is reg else label
//	COMMAND ::= LABEL? /[CJSTZ][(]\d+(:?,\s*\d+)*/
//	LABEL ::= /^\pL\w*:\s*/
//
// Any line may end with a comment which begins with ';'. All comments and
// blank lines are ignored.
func (me *Urm) Load(lines []string, size int) error {
	me.clear(size)
	pc := 0
	start := 1
	sized := false
	for i, line := range lines {
		lino := i + 1
		line = cleanLine(line)
		if line != "" {
			var err error
			if strings.HasPrefix(line, "#") {
				sized, err = me.setRegsSize(lino, line, pc, sized)
			} else {
				pc, start, err = me.readStatement(lino, line, pc, start)
			}
			if err != nil {
				return err
			}
		}
	}
	return me.setStart(start)
}

// Run runs the Urm program for up to DefaultMaxSteps steps.
// (See RunX and Step.)
func (me *Urm) Run() error { return me.RunX(DefaultMaxSteps) }

// Run runs the Urm program for a maximum of maxSteps steps.
// (See Run and Step.)
func (me *Urm) RunX(maxSteps int) error {
	return fmt.Errorf("RunX() unimplemented") // TODO
	/*
		for {
			me.steps++
			if me.steps > maxSteps {
				return fmt.Errorf("RunX() %w: %d", Err100, me.steps)
			}
			err := me.Step()
			if errors.Is(err, ErrStop) {
				return nil // normal termination
			}
			if err != nil {
				return err // error termination
			}
		}
	*/
}

// Steps runs the one step (i.e., executes the next statement) the Urm
// program. (See Run and RunX.)
func (me *Urm) Step() error {
	return fmt.Errorf("Step() unimplemented") // TODO
	/*
		cmd, err := me.nextCommand() // returns cmd & inc PC
		if err != nil {
			return err
		}
		me.steps++
		return me.executeCommand(cmd) // may acquire 1 or more operands
	*/
}

// Pc returns the value of the program counter (register 0).
func (me *Urm) Pc() int { return me.regs[PcReg].Value() }

// SetPc sets the value of the program counter (register 0).
func (me *Urm) SetPc(reg int) error {
	if 0 <= reg && reg < len(me.regs) {
		me.regs[PcReg] = Value(reg)
		return nil
	}
	return fmt.Errorf("SetPc() %w: %d", Err104, reg)
}
