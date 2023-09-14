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
	DefaultSize     = 100
	DefaultMaxSteps = 10_000
	PC              = "PC"
	PcReg           = 0
	cmdJump         = 203 // J(n, m, t) ; jump to t if m = n
	cmdCopy         = 202 // C(n, m) ; n → m
	cmdSucc         = 201 // S(n) ; n++
	cmdZero         = 200 // Z(n) ; n = 0
	startLabel      = "START"
	stopCmd         = "STOP"
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
	ErrBug = errors.New("EBUG: bug")
)
