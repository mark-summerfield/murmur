// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"errors"
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
		if address, ok := reger.(Address); ok {
			if reg, ok := me.reg_for_label[address.String()[1:]]; ok {
				reg := me.regs[reg] // find addressed reg
				return NewValue(reg.Value()), nil
			} else {
				return nil, fmt.Errorf("%w: %d", Err122, reg)
			}
		} else if label, ok := reger.(Label); ok {
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
	if !ok { // not a command set by C J S etc; but maybe by literal value?
		value, ok := reger.(Value)
		if ok {
			switch value.Value() {
			case cmdAdd:
				cmd = NewAddCommand(0)
			case cmdCopy:
				cmd = NewCopyCommand(0)
			case cmdDiv:
				cmd = NewDivCommand(0)
			case cmdJump:
				cmd = NewJumpCommand(0)
			case cmdJumpGt:
				cmd = NewJumpGtCommand(0)
			case cmdJumpLt:
				cmd = NewJumpLtCommand(0)
			case cmdMul:
				cmd = NewMulCommand(0)
			case cmdPred:
				cmd = NewPredCommand(0)
			case cmdSub:
				cmd = NewSubCommand(0)
			case cmdSucc:
				cmd = NewSuccCommand(0)
			case cmdZero:
				cmd = NewZeroCommand(0)
			default:
				ok = false
			}
			if ok { // replace the literal command value with actual command
				me.regs[me.pc()-1] = cmd
			}

		}
		if !ok { // Not a Commander and not a value in {cmdCopy ...}
			reg := null
			if reger != nil {
				reg = reger.String()
			}
			return nil, fmt.Errorf("%w: %s", Err101, reg)
		}
	}
	return cmd, nil
}

func (me *Urm) executeCommand(cmd Commander) error {
	var err error
	switch cmd.(type) {
	case AddCommand, DivCommand, MulCommand, SubCommand:
		err = me.doMath(cmd.Value())
	case CopyCommand:
		err = me.doCopy()
	case JumpCommand, JumpGtCommand, JumpLtCommand:
		err = me.doJump(cmd.Value())
	case PredCommand:
		err = me.doIncOrDec(-1)
	case SuccCommand:
		err = me.doIncOrDec(1)
	case ZeroCommand:
		err = me.doZero()
	default:
		return fmt.Errorf("%w: %s", Err110, cmd)
	}
	if err != nil && !errors.Is(err, ErrStop) && cmd.Lino() != 0 {
		err = fmt.Errorf("line#%d %w", cmd.Lino(), err)
	}
	return err
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
	value, err := me.RegValue(op1.Value())
	if err != nil {
		return err
	}
	return me.SetRegToValue(op2.Value(), value)
}

func (me *Urm) doJump(cmd int) error {
	op1, err := me.nextReger(true)
	if err != nil {
		return err
	}
	value1, err := me.RegValue(op1.Value())
	if err != nil {
		return err
	}
	op2, err := me.nextReger(true)
	if err != nil {
		return err
	}
	value2, err := me.RegValue(op2.Value())
	if err != nil {
		return err
	}
	op3, err := me.nextReger(true)
	if err != nil {
		return err
	}
	met := false
	switch cmd {
	case cmdJump:
		if value1 == value2 {
			met = true
		}
	case cmdJumpLt:
		if value1 < value2 {
			met = true
		}
	case cmdJumpGt:
		if value1 > value2 {
			met = true
		}
	}
	if met {
		return me.setPc(op3.Value())
	}
	return nil
}

func (me *Urm) doIncOrDec(amount int) error {
	op1, err := me.nextReger(true)
	if err != nil {
		return err
	}
	reg := op1.Value()
	if value, err := me.RegValue(reg); err != nil {
		return err
	} else {
		return me.SetRegToValue(reg, value+amount)
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
	return me.SetRegToValue(reg, 0)
}

func (me *Urm) doMath(cmd int) error {
	op1, err := me.nextReger(true)
	if err != nil {
		return err
	}
	target := op1.Value()
	value1, err := me.RegValue(target)
	if err != nil {
		return err
	}
	op2, err := me.nextReger(true)
	if err != nil {
		return err
	}
	value2, err := me.RegValue(op2.Value())
	if err != nil {
		return err
	}
	value := 0
	switch cmd {
	case cmdAdd:
		value = value1 + value2
	case cmdDiv:
		value = value1 / value2
	case cmdMul:
		value = value1 * value2
	case cmdSub:
		value = value1 - value2
	}
	return me.SetRegToValue(target, value)
}
