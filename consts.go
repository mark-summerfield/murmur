// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package murmur

import (
	_ "embed"
	"errors"
)

//go:embed Version.dat
var Version string

const (
	DefaultSize     = 200
	DefaultMaxSteps = 100_000
	PC              = "PC"
	PcReg           = 0
	cmdAdd          = 210 // +(r, s) ; r = r + s
	cmdSub          = 209 // -(r, s) ; r = r - s
	cmdMul          = 208 // *(r, s) ; r = r * s
	cmdDiv          = 207 // /(r, s) ; r = r / s (truncating)
	cmdJumpLt       = 206 // L(r, s, t) ; jump to t if r < s
	cmdJumpGt       = 205 // G(r, s, t) ; jump to t if r > s
	cmdPred         = 204 // P(r) or D(r) ; r--
	cmdJump         = 203 // J(r, s, t) ; jump to t if r = s
	cmdCopy         = 202 // C(r, s) or T(r, s) ; r → s
	cmdSucc         = 201 // S(r) or I(r, s) ; r++
	cmdZero         = 200 // Z(r) ; r = 0
	startLabel      = "START"
	stopCmd         = "STOP"
	null            = "NULL"
	minLabelWidth   = 8
	maxLabelWidth   = 32
)

var (
	ErrStop = errors.New(stopCmd) // indicates normal termination
	Err100  = errors.New("E100: stopped after max steps")
	Err101  = errors.New("E101: unrecognized command")
	Err102  = errors.New("E102: invalid label (must match /^\\pL\\w*$/")
	Err103  = errors.New("E103: undefined label")
	Err104  = errors.New("E104: can't set PC to out-of-range register")
	Err105  = errors.New("E105: invalid size")
	Err106  = errors.New(
		"E106: may only set size once, before code and data")
	Err107 = errors.New("E107: invalid statement")
	Err108 = errors.New("E108: can only set numbered reg to numeric value")
	Err109 = errors.New("E109: ran out of registers")
	Err110 = errors.New("E110: invalid command")
	Err111 = errors.New("E111: can't add label for out-of-range register")
	Err112 = errors.New("E112: can't set value for out-of-range register")
	Err113 = errors.New(
		"E113: can't get reg value for out-of-range register")
	Err114 = errors.New("E114: can't set value for out-of-range register")
	Err115 = errors.New("E115: can't get value for out-of-range register")
	Err116 = errors.New("E116: can't get value for invalid label")
	Err117 = errors.New("E117: unrecognized command")
	Err118 = errors.New("E118: wrong number of arguments")
	Err119 = errors.New("E119: ran out of registers")
	Err120 = errors.New("E120: can't get value for out-of-range register")
	Err121 = errors.New(
		"E121: invalid address label (must match /^\\pL\\w*$/")
	Err122 = errors.New("E122: undefined address")
	Err123 = errors.New("E123: invalid data value")
	Err124 = errors.New("E124: undefined label")
	ErrBug = errors.New("EBUG: bug")
)
