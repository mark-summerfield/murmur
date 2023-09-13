// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
	"strconv"
	"strings"
)

func (me *Urm) clear(size int) {
	if size == 0 {
		size = DefaultSize
	}
	if size < len(me.regs) {
		me.regs = me.regs[:size]
	} else if size > len(me.regs) {
		me.regs = make([]Reger, size)
	}
	_ = me.setPc(1) // start at register 1 by default
	me.reg_for_label = map[string]int{}
	me.label_for_reg = map[int]string{}
	_ = me.addRegLabel(PcReg, PC)
	me.steps = 0
}

func (me *Urm) pc() int { return me.regs[PcReg].Value() }

func (me *Urm) setPc(reg int) error {
	if 0 <= reg && reg < len(me.regs) {
		me.regs[PcReg] = NewValue(reg)
		return nil
	}
	return fmt.Errorf("%w: %d", Err104, reg)
}

func (me *Urm) asString(withRegNums bool) string {
	texts := []string{fmt.Sprintf("#%d", me.Size())}
	for reg := 1; reg < me.Size(); reg++ {
		reger := me.regs[reg]
		label, hasLabel := me.label_for_reg[reg]
		ops := make([]string, 0, 3)
		switch reger.(type) {
		case CopyCommand:
			reg++
			ops = append(ops, me.operand(reg))
			reg++
			ops = append(ops, me.operand(reg))
		case JumpCommand:
			reg++
			ops = append(ops, me.operand(reg))
			reg++
			ops = append(ops, me.operand(reg))
			reg++
			ops = append(ops, me.operand(reg))
		case SuccCommand:
			reg++
			ops = append(ops, me.operand(reg))
		case ZeroCommand:
			reg++
			ops = append(ops, me.operand(reg))
		}
		text := reger.String()
		if len(ops) > 0 {
			text += strings.Join(ops, ", ") + ")"
		}
		if hasLabel {
			text = fmt.Sprintf("%s:\t%s", label, text)
		} else {
			text = "\t" + text
		}
		if withRegNums {
			text += " ; " + strconv.Itoa(reg)
		}
		texts = append(texts, text)
	}
	return strings.Join(texts, "\n")
}

func (me *Urm) operand(reg int) string {
	if 0 <= reg && reg < len(me.regs) {
		return me.regs[reg].String()
	}
	return strconv.Itoa(reg) // should never be reached
}
