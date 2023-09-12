// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
	"regexp"
	"strconv"
)

func (me *Urm) clear(size int) {
	if size == 0 {
		size = DefaultSize
	}
	if size < len(me.regs) {
		me.regs = me.regs[:size]
	} else if size > len(me.regs) {
		me.regs = make([]RegValue, size)
	}
	_ = me.SetPc(1) // start at register 1 by default
	me.reg_for_label = map[string]int{}
	me.label_for_reg = map[int]string{}
	me.setLabel(PcReg, PC)
	me.steps = 0
}

func (me *Urm) setLabel(reg int, label string) error {
	if 0 <= reg && reg < len(me.regs) {
		rx := regexp.MustCompile(`^\pL\w*$`)
		if !rx.MatchString(label) {
			return fmt.Errorf("setLabel() %w: %d", Err102, reg)
		}
		me.reg_for_label[label] = reg
		me.label_for_reg[reg] = label
	}
	return fmt.Errorf("setLabel() %w: %d", Err104, reg)
}

func (me *Urm) setRegsSize(lino int, line string, pc int, sized bool) (bool,
	error) {
	if sized || pc > 0 {
		return sized, fmt.Errorf("setRegsSize() line#%d: %w", lino, Err106)
	}
	size, err := strconv.Atoi(line[1:]) // skip leading #
	if err != nil {
		return sized, fmt.Errorf("setRegsSize() line#%d %w: %w: %s", lino,
			Err105, err, line[1:])
	}
	if size < len(me.regs) {
		me.regs = me.regs[:size]
	} else if size > len(me.regs) {
		me.regs = make([]RegValue, size)
	}
	return true, nil
}

func (me *Urm) readStatement(lino int, line string, pc, start int) (
	int, int, error) {
	var err error
	pc++
	rx := regexp.MustCompile(`(?:(\w+):)?\s*(\d+|` + stopCmd +
		`|[CJSTZ][(][^)]+[)])`)
	if matches := rx.FindStringSubmatch(line); matches != nil {
		label := matches[1]
		command := matches[2]
		value, verr := strconv.Atoi(command)
		if label != "" {
			if reg, err := strconv.Atoi(label); err == nil { // reg: val
				if verr != nil {
					return pc, start, fmt.Errorf(
						"setRegToValue() line#%d: %w", lino, Err108)
				}
				return pc, start, me.setRegValue(reg, value)
			}
			if err = me.setLabel(pc, label); err != nil { // label: cmd
				return pc, start, err // if err == nil falls through...
			}
		}
		if verr == nil { // label: val; e.g., A: 7
			err = me.setRegValue(pc, value) // label already handled
		} else if command == stopCmd { // STOP or label: STOP
			pc, err = me.setStopCommand(lino, pc)
		} else { // command or label: command
			pc, start, err = me.readCommand(command, lino, pc, start)
		}
		if err != nil {
			return pc, start, err
		}
	} else {
		return pc, start, fmt.Errorf("readStatement() line#%d%w: %s", lino,
			Err107, line)
	}
	return pc, start, nil
}

func (me *Urm) setRegValue(reg, value int) error {
	if 0 <= reg && reg < len(me.regs) {
		me.regs[reg] = Value(value)
		return nil
	}
	return fmt.Errorf("setRegValue() %w: %d", Err104, reg)
}

func (me *Urm) setStopCommand(lino, pc int) (int, error) {
	if err := me.setRegValue(pc, cmdZero); err != nil {
		return pc, err
	}
	pc++
	if pc >= len(me.regs) {
		return pc, fmt.Errorf("setStopCommand() line#%d: %w", lino, Err109)
	}
	return pc, me.setRegValue(pc, 0)
}
