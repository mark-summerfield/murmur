// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package murmur

import (
	"fmt"
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
	panic("StringWithRegNums() not implemented") // TODO
}
