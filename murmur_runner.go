// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
)

func (me *Urm) nextReger() (int, Reger, error) {
	reg := me.pc()
	if 0 <= reg && reg < len(me.regs) {
		if err := me.setPc(reg + 1); err != nil {
			return 0, nil, err
		}
		reger := me.regs[reg]
		if label, ok := reger.(Label); ok {
			if reg, ok := me.reg_for_label[label.String()]; ok {
				return reg, NewValue(reg), nil
			} else {
				return 0, nil, fmt.Errorf("%w: %d", Err103, reg)
			}
		}
		return reg, reger, nil
	}
	return 0, nil, fmt.Errorf("%w: %d", Err113, reg)
}

func (me *Urm) nextCommand() (Commander, error) {
	_, reger, err := me.nextReger()
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
	reg1, value1, err := me.nextReger() // every command has > 0 operands
	if err != nil {
		return err
	}
	switch cmd.(type) {
	case CopyCommand:
		return me.doCopy(reg1, value1)
	case JumpCommand:
		return me.doJump(reg1, value1)
	case SuccCommand:
		return me.doSucc(reg1, value1)
	case ZeroCommand:
		return me.doZero(reg1, value1)
	}
	return fmt.Errorf("%w: %s", Err110, cmd)
}

func (me *Urm) doCopy(reg1 int, value1 Reger) error {
	reg2, _, err := me.nextReger()
	if err != nil {
		return err
	}
	fmt.Printf("Copy regs[%d] = regs[%d] ; value=%d\n", reg2, reg1, value1.Value())
	return me.setRegToValue(reg2, value1.Value())
}

func (me *Urm) doJump(reg1 int, value1 Reger) error {
	_, value2, err := me.nextReger()
	if err != nil {
		return err
	}
	_, value3, err := me.nextReger()
	if err != nil {
		return err
	}
	if value1.Value() == value2.Value() {
		fmt.Printf("Jump (PC) regs[0] = %d ; jumped\n", value3.Value())
		return me.setPc(value3.Value())
	}
	fmt.Printf("Jump (PC) regs[0] = %d ; skipped \n", me.pc())
	return nil
}

func (me *Urm) doSucc(reg1 int, value1 Reger) error {
	fmt.Printf("Succ regs[%d] = %d\n", reg1, value1.Value()+1)
	return me.setRegToValue(reg1, value1.Value()+1)
}

func (me *Urm) doZero(reg1 int, value1 Reger) error {
	if reg1 == PcReg {
		return ErrStop
	}
	fmt.Printf("Zero regs[%d] = 0\n", reg1)
	return me.setRegToValue(reg1, 0)
}
