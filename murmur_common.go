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
	me.reg_for_label = map[string]int{}
	me.label_for_reg = map[int]string{}
	_ = me.addRegLabel(PcReg, PC)
	me.steps = 0
	me.labelWidth = minLabelWidth
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
		text := reger.String()
		label, hasLabel := me.label_for_reg[reg]
		if cmd, ok := reger.(Commander); ok {
			ops := make([]string, 0, 3)
			for i := 0; i < cmd.Arity(); i++ {
				reg++
				ops = append(ops, me.operand(reg))
			}
			text += strings.Join(ops, ", ") + ")"
		}
		if hasLabel {
			text = fmt.Sprintf("%*s:%s", me.labelWidth, label, text)
		} else {
			text = fmt.Sprintf("%*s", me.labelWidth+1, text)
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
	panic(fmt.Errorf("%w: operand %d", ErrBug, reg)) // shouldn't happen
}
