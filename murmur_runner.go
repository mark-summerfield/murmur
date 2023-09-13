// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
)

func (me *Urm) nextReger() (Reger, error) {
	reg := me.pc()
	if 0 <= reg && reg < len(me.regs) {
		if err := me.setPc(reg + 1); err != nil {
			return nil, err
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
	reger, err := me.nextReger()
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
	op1, err := me.nextReger() // every command has at least one operand
	if err != nil {
		return err
	}
	switch cmd.(type) {
	case CopyCommand:
		return me.doCopy(op1)
	case JumpCommand:
		return me.doJump(op1)
	case SuccCommand:
		return me.doSucc(op1)
	case ZeroCommand:
		return me.doZero(op1)
	}
	return fmt.Errorf("%w: %s", Err110, cmd)
}

func (me *Urm) doCopy(op1 Reger) error {
	op2, err := me.nextReger()
	if err != nil {
		return err
	}
	return fmt.Errorf("doCopy unimplemented: %s → %s", op1, op2) // TODO
}

func (me *Urm) doJump(op1 Reger) error {
	op2, err := me.nextReger()
	if err != nil {
		return err
	}
	op3, err := me.nextReger()
	if err != nil {
		return err
	}
	return fmt.Errorf("doJump unimplemented: %s == %s → %s", op1, op2, op3) // TODO
}

func (me *Urm) doSucc(op1 Reger) error {
	return fmt.Errorf("doSucc unimplemented: %s", op1) // TODO
}

func (me *Urm) doZero(op1 Reger) error {
	return fmt.Errorf("doZero unimplemented: %s", op1) // TODO
}
