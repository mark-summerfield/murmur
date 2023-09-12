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
	Jump            = 253 // J(n, m, t) ; jump to t if m = n
	Copy            = 252 // C(n, m) ; n → m
	Succ            = 251 // S(n) ; n++
	Zero            = 250 // Z(n) ; n = 0
	sentinalBase    = 440_000
	startLabel      = "START"
	stopCmd         = "STOP"
)

type mode uint8

const (
	argMode mode = iota
	baseMode
	cmdMode
)

var (
	ErrStop = errors.New(stopCmd) // indicates normal termination
	Err100  = errors.New("E100: stopped after max steps")
	Err101  = errors.New("E101: unrecognized command")
	Err102  = errors.New("E102: invalid label (must match /^\\pL\\w*$/")
	Err103  = errors.New("E103: undefined label")
	Err104  = errors.New("E104: invalid register")
	Err105  = errors.New("E105: invalid size")
	Err106  = errors.New("E106: may only set size before code and data")
	Err107  = errors.New("E107: invalid statement")
	Err108  = errors.New("E108: invalid value")
	Err109  = errors.New("E109: ran out of registers")
	Err110  = errors.New("E110: invalid arity")
	ErrBug  = errors.New("EBUG: bug")
)
