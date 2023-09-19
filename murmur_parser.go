// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
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

func (me *Urm) readStatement(lino int, line string, pc, start int) (int,
	int, error) {
	var err error
	rx := regexp.MustCompile(`(\pL\w*):\s*(\d+(:?\s+\d+)*)`)
	if matches := rx.FindAllStringSubmatch(line, -1); matches != nil &&
		len(matches[0]) > 2 {
		return me.readDataLine(lino, matches, pc, start)
	}
	pc++
	rx = regexp.MustCompile(`(?:(\w+):)?\s*(\d+|` + stopCmd +
		`|[-+*/CDGIJLPSTZcdgijlpstz][(][^)]+[)])`)
	if matches := rx.FindStringSubmatch(line); matches != nil {
		label := matches[1]
		command := matches[2]
		value, verr := strconv.Atoi(command)
		if label != "" {
			if err = me.readLabel(lino, label, value, verr, pc,
				start); err != nil {
				return pc, start, err
			} // if err == nil falls through...
		}
		if verr == nil { // label: val; e.g., A: 7
			err = me.SetRegToValue(pc, value) // label already handled
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

func (me *Urm) readLabel(lino int, label string, value int, verr error, pc,
	start int) error {
	var err error
	if reg, err := strconv.Atoi(label); err == nil { // reg: val
		if verr != nil {
			return fmt.Errorf("line#%d %w", lino, Err108)
		}
		if err = me.SetRegToValue(reg, value); err != nil {
			return fmt.Errorf("line#%d %w", lino, err)
		}
		return nil
	}
	if err = me.addRegLabel(pc, label); err != nil { // label: cmd
		return fmt.Errorf("line#%d %w", lino, err)
	}
	return nil
}

func (me *Urm) readDataLine(lino int, matches [][]string, pc, start int) (
	int, int, error) {
	var err error
	label := matches[0][1]
	if err = me.addRegLabel(pc+1, label); err != nil {
		return pc, start, err
	}
	for _, text := range strings.Fields(matches[0][2]) {
		if value, err := strconv.Atoi(text); err != nil {
			return pc, start, fmt.Errorf("line#%d %w", lino, err)
		} else {
			pc++
			if err := me.SetRegToValue(pc, value); err != nil {
				return pc, start, fmt.Errorf("line#%d %w", lino, err)
			}
		}
	}
	return pc, start, nil
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

func (me *Urm) setRegToAddress(reg int, address string) error {
	if 0 <= reg && reg < len(me.regs) {
		if address, err := NewAddress(address); err != nil {
			return err
		} else {
			me.regs[reg] = address
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
		width := utf8.RuneCountInString(label)
		me.labelWidth = min(max(me.labelWidth, width), maxLabelWidth)
		return nil
	}
	return fmt.Errorf("%w: %d", Err111, reg)
}

func (me *Urm) setStopCommand(lino, pc int) (int, error) {
	if err := me.setRegToCommand(pc, NewZeroCommand(lino)); err != nil {
		return pc, fmt.Errorf("line#%d %w", lino, err)
	}
	pc++
	if pc >= len(me.regs) {
		return pc, fmt.Errorf("line#%d %w", lino, Err109)
	}
	return pc, me.SetRegToValue(pc, 0)
}

func (me *Urm) readCommand(command string, lino, pc, start int) (int, int,
	error) {
	var err error
	if start == 0 {
		start = pc
	}
	cmd, err := NewCommand(lino, command[0])
	if err != nil {
		return pc, start, fmt.Errorf("line#%d %w", lino, err)
	}
	if err = me.setRegToCommand(pc, cmd); err != nil {
		return pc, start, fmt.Errorf("line#%d %w", lino, err)
	}
	ops, err := me.getOps(cmd, command)
	if err != nil {
		return pc, start, fmt.Errorf("line#%d %w", lino, err)
	}
	for _, op := range ops {
		pc++
		if err = me.addOp(pc, op); err != nil {
			return pc, start, fmt.Errorf("line#%d %w", lino, err)
		}
	}
	return pc, start, nil
}

func (me *Urm) getOps(cmd Commander, command string) ([]string, error) {
	rx := regexp.MustCompile(`[\s,]+`)
	ops := make([]string, 0, 3)
	ops = append(ops, rx.Split(command[2:len(command)-1], -1)...)
	arity := cmd.Arity()
	if _, ok := cmd.(JumpCommand); ok && len(ops) == 1 {
		ops = append([]string{"1", "1"}, ops...) // J(LBL) → J(1, 1, LBL)
	}
	if arity != len(ops) {
		return nil, fmt.Errorf("%w: need %d; got %d", Err118, arity,
			len(ops))
	}
	return ops, nil
}

func (me *Urm) addOp(pc int, op string) error {
	if pc >= len(me.regs) {
		return Err119
	}
	if value, err := strconv.Atoi(op); err == nil { // literal reg val
		return me.SetRegToValue(pc, value)
	} else if strings.HasPrefix(op, "@") { // address
		return me.setRegToAddress(pc, op)
	} else { // label
		return me.setRegToLabel(pc, op)
	}
}

func (me *Urm) setStart(start int) error {
	if start == 0 { // All input was numbers, e.g., 1: 203, 2: 0, etc.
		start = 1
	}
	if reg, ok := me.reg_for_label[startLabel]; ok {
		start = reg // START label takes priority over heuristic
	} else { // no START label, so add one
		if err := me.addRegLabel(start, startLabel); err != nil {
			return err // shouldn't happen
		}
	}
	return me.setPc(start)
}
