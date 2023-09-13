// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
	"regexp"
	"strconv"
)

func (me *Urm) setRegsSize(lino int, line string, pc int, sized bool) (bool,
	error) {
	if sized || pc > 0 {
		return sized, fmt.Errorf("line#%d %w", lino, Err106)
	}
	size, err := strconv.Atoi(line[1:]) // skip leading #
	if err != nil {
		return sized, fmt.Errorf("line#%d %w: %w: %s", lino, Err105, err,
			line[1:])
	}
	if size < len(me.regs) {
		me.regs = me.regs[:size]
	} else if size > len(me.regs) {
		me.regs = make([]Reger, size)
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
					return pc, start, fmt.Errorf("line#%d %w", lino, Err108)
				}
				if err = me.setRegToValue(reg, value); err != nil {
					return pc, start, fmt.Errorf("line#%d %w", lino, err)
				}
				return pc, start, nil
			}
			if err = me.addRegLabel(pc, label); err != nil { // label: cmd
				return pc, start, fmt.Errorf("line#%d %w", lino, err)
			} // if err == nil falls through...
		}
		if verr == nil { // label: val; e.g., A: 7
			err = me.setRegToValue(pc, value) // label already handled
			if err != nil {
				err = fmt.Errorf("line#%d %w", lino, err)
			}
		} else if command == stopCmd { // STOP or label: STOP
			pc, err = me.setStopCommand(lino, pc)
		} else { // command or label: command
			pc, start, err = me.readCommand(command, lino, pc, start)
		}
		if err != nil {
			return pc, start, err
		}
	} else {
		return pc, start, fmt.Errorf("line#%d %w: %s", lino, Err107, line)
	}
	return pc, start, nil
}

func (me *Urm) setRegToValue(reg, value int) error {
	if 0 <= reg && reg < len(me.regs) {
		me.regs[reg] = NewValue(value)
		return nil
	}
	return fmt.Errorf("%w: %d", Err112, reg)
}

func (me *Urm) setRegToLabel(reg int, label string) error {
	if 0 <= reg && reg < len(me.regs) {
		if label, err := NewLabel(label); err != nil {
			return err
		} else {
			me.regs[reg] = label
			return nil
		}
	}
	return fmt.Errorf("%w: %d", Err114, reg)
}

func (me *Urm) setRegToCommand(reg int, cmd Commander) error {
	if 0 <= reg && reg < len(me.regs) {
		me.regs[reg] = cmd
		return nil
	}
	return fmt.Errorf("%w: %d", Err114, reg)
}

func (me *Urm) addRegLabel(reg int, label string) error {
	if 0 <= reg && reg < len(me.regs) {
		me.reg_for_label[label] = reg
		me.label_for_reg[reg] = label
		return nil
	}
	return fmt.Errorf("%w: %d", Err111, reg)
}

func (me *Urm) setStopCommand(lino, pc int) (int, error) {
	if err := me.setRegToValue(pc, cmdZero); err != nil {
		return pc, fmt.Errorf("line#%d %w", lino, err)
	}
	pc++
	if pc >= len(me.regs) {
		return pc, fmt.Errorf("line#%d %w", lino, Err109)
	}
	return pc, me.setRegToValue(pc, 0)
}

func (me *Urm) readCommand(command string, lino, pc, start int) (int, int,
	error) {
	var err error
	cmd, err := NewCommand(command[0])
	if err != nil {
		return pc, start, fmt.Errorf("line#%d %w", lino, err)
	}
	rx := regexp.MustCompile(`[\s,]+`)
	ops := make([]string, 0, 3)
	ops = append(ops, rx.Split(command[2:len(command)-1], -1)...)
	// TODO
	return pc, start, fmt.Errorf(
		"readCommand unimplemented: %q #%d %d %d %v %q", command, lino, pc,
		start, ops, cmd) // TODO
}

func (me *Urm) setStart(start int) error {
	return fmt.Errorf("setStart unimplemented: %d", start) // TODO
}
