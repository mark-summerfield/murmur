// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
)

func (me *Urm) nextReger(inc bool) (Reger, error) {
	reg := me.pc()
	if 0 <= reg && reg < len(me.regs) {
		if inc {
			if err := me.setPc(reg + 1); err != nil {
				return nil, err
			}
		}
		reger := me.regs[reg]
		if label, ok := reger.(Label); ok {
			if reg, ok := me.reg_for_label[label.String()]; ok {
				return NewValue(reg), nil
			} else {
				return nil, fmt.Errorf("%w: %d", Err103, reg)
			}
		}
		return reger, nil
	}
	return nil, fmt.Errorf("%w: %d", Err113, reg)
}

func (me *Urm) nextCommand() (Commander, error) {
	reger, err := me.nextReger(true)
	if err != nil {
		return nil, err
	}
	cmd, ok := reger.(Commander)
	if !ok { // not a command set by C J S or Z; but maybe by literal value?
		value, ok := reger.(Value)
		if ok {
			switch value.Value() {
			case cmdCopy:
				cmd = NewCopyCommand()
			case cmdJump:
				cmd = NewJumpCommand()
			case cmdSucc:
				cmd = NewSuccCommand()
			case cmdZero:
				cmd = NewZeroCommand()
			default:
				ok = false
			}
			if ok { // replace the literal command value with actual command
				me.regs[me.pc()-1] = cmd
			}

		}
		if !ok { // Not a Commander and not a value in {cmdCopy ...}
			return nil, fmt.Errorf("%w: %s", Err101, reger)
		}
	}
	return cmd, nil
}

func (me *Urm) executeCommand(cmd Commander) error {
	switch cmd.(type) {
	case CopyCommand:
		return me.doCopy()
	case JumpCommand:
		return me.doJump()
	case SuccCommand:
		return me.doSucc()
	case ZeroCommand:
		return me.doZero()
	}
	return fmt.Errorf("%w: %s", Err110, cmd)
}

func (me *Urm) regValue(reg int) (int, error) {
	if 0 <= reg && reg <= len(me.regs) {
		return me.regs[reg].Value(), nil
	}
	return 0, fmt.Errorf("%w: %d", Err120, reg)
}

func (me *Urm) doCopy() error {
	op1, err := me.nextReger(true)
	if err != nil {
		return err
	}
	op2, err := me.nextReger(true)
	if err != nil {
		return err
	}
	value, err := me.regValue(op1.Value())
	if err != nil {
		return err
	}
	return me.setRegToValue(op2.Value(), value)
}

func (me *Urm) doJump() error {
	op1, err := me.nextReger(true)
	if err != nil {
		return err
	}
	value1, err := me.regValue(op1.Value())
	if err != nil {
		return err
	}
	op2, err := me.nextReger(true)
	if err != nil {
		return err
	}
	value2, err := me.regValue(op2.Value())
	if err != nil {
		return err
	}
	op3, err := me.nextReger(true)
	if err != nil {
		return err
	}
	if value1 == value2 {
		return me.setPc(op3.Value())
	}
	return nil
}

func (me *Urm) doSucc() error {
	op1, err := me.nextReger(true)
	if err != nil {
		return err
	}
	reg := op1.Value()
	if value, err := me.regValue(reg); err != nil {
		return err
	} else {
		return me.setRegToValue(reg, value+1)
	}
}

func (me *Urm) doZero() error {
	op1, err := me.nextReger(false)
	if err != nil {
		return err
	}
	reg := op1.Value()
	if reg == PcReg {
		return ErrStop
	}
	if err := me.setPc(me.pc() + 1); err != nil { // manually inc PC
		return err
	}
	return me.setRegToValue(reg, 0)
}
