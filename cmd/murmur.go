// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/mark-summerfield/clip"
	"github.com/mark-summerfield/murmur"
	"github.com/mark-summerfield/ufile"
)

func main() {
	config := getConfig()
	lines := readLines(config.infile)
	out, closer := getOutput(config.outfile)
	defer closer()
	urm := load(out, lines, config)
	if config.step {
		step(out, urm, config.maxSteps, &config.watch)
	} else {
		run(out, urm, config.maxSteps, &config.watch,
			config.chars.IsEmpty())
	}
	if config.dis {
		_, _ = out.WriteString(urm.StringWithRegNums() + "\n")
	}
	if !config.chars.IsEmpty() {
		writeChars(out, urm, config.chars)
	}
}

func getConfig() *Config {
	parser := clip.NewParserVersion(murmur.Version)
	parser.LongDesc = "An Unlimited Register Machine (URM) emulator " +
		"with optional extensions (indirect addressing and some extra " +
		"convenience commands)."
	outfileOpt := parser.Str("outfile",
		"The file to write to [default: stdout].", "")
	maxStepsOpt := parser.IntInRange("maxsteps",
		fmt.Sprintf("The maximum steps to execute [default: %d].",
			int(murmur.DefaultMaxSteps)), 1, math.MaxInt,
		int(murmur.DefaultMaxSteps))
	stepOpt := parser.Flag("step",
		"Run step by step showing registers at every step [default: run "+
			"nonstop showing initial and final registers only].")
	registersOpt := parser.IntInRange("registers",
		"How many registers the URM has; can also be overridden in the "+
			".urm with ^n, e.g., ^3000 [default: 1000 unless overridden].",
		1, math.MaxInt, int(murmur.DefaultSize))
	watchOpt := parser.Str("watch", "A comma-separated list of which "+
		"registers to watch with each item of the form r or r-s or "+
		"label where r and s are integers ≥ 0. Out of range registers "+
		"and missing labels are ignored [default: PC,1-9].", "PC,1-9")
	charsOpt := parser.Str("chars", "The same comma-separated list "+
		"of registers as the watch option, only these are displayed as "+
		"Unicode characters rather than as registers. [no default].", "")
	disOpt := parser.Flag("dis",
		"Show diassembly of all registers at the start and at the end.")
	parser.PositionalCount = clip.OneOrMorePositionals
	parser.PositionalHelp = "ARG1 is the .urm file to run, optionally " +
		"followed by data values for registers 1, 2, … The data may be " +
		"given as space-separated numbers (e.g., 1 2 4 8), or as a " +
		"single argument of quoted text (e.g., \"some text\") in " +
		"which case the registers will be set to the text's Unicode " +
		"code points."
	parser.MustSetPositionalVarName("ARG")
	if err := parser.Parse(); err != nil {
		parser.OnError(err) // doesn't return
		return nil          // never reached
	}
	registers := registersOpt.Value()
	config := Config{maxSteps: maxStepsOpt.Value(),
		step: stepOpt.Value(), registers: registers,
		infile: parser.Positionals[0], args: parser.Positionals[1:],
		outfile: outfileOpt.Value(), watch: parseWatches(watchOpt.Value()),
		chars: parseWatches(charsOpt.Value()), dis: disOpt.Value()}
	return &config
}

type Config struct {
	maxSteps  int
	step      bool
	registers int
	infile    string
	args      []string
	outfile   string
	watch     watches
	chars     watches
	dis       bool
}

// String is for debugging.
func (me Config) String() string {
	outfile := me.outfile
	if outfile == "" {
		outfile = "stdout"
	}
	return fmt.Sprintf(
		"maxsteps=%d step=%t registers=%d infile=%s args=%v outfile=%s "+
			"watches=%v dis=%t", me.maxSteps, me.step, me.registers,
		me.infile, me.args, outfile, me.watch, me.dis)
}

func readLines(filename string) []string {
	raw, err := os.ReadFile(filename)
	if err != nil {
		onError(err)
	}
	return strings.Split(string(raw), "\n")
}

func getOutput(filename string) (*os.File, func()) {
	out := os.Stdout
	closer := func() {}
	if filename != "" {
		out, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE,
			ufile.ModeUserRW)
		if err != nil {
			onError(err)
		}
		closer = func() { _ = out.Close() }
	}
	return out, closer
}

func load(out *os.File, lines []string, config *Config) *murmur.Urm {
	urm, err := murmur.NewX(lines, config.registers)
	if err != nil {
		onError(err)
	}
	if len(config.args) == 1 && len(config.args[0]) > 0 &&
		!unicode.IsDigit(rune(config.args[0][0])) { // text arg
		for i, r := range config.args[0] {
			if err := urm.SetRegToValue(i+1, int(r)); err != nil {
				onError(err)
			}
		}
	} else {
		for i, arg := range config.args {
			if value, err := strconv.Atoi(arg); err != nil {
				onError(err)
			} else {
				if err := urm.SetRegToValue(i+1, value); err != nil {
					onError(err)
				}
			}
		}
	}
	if config.dis {
		_, _ = out.WriteString(urm.StringWithRegNums() + "\n")
	}
	if config.chars.IsEmpty() {
		if err := writeWatched(out, urm, &config.watch); err != nil {
			onError(err)
		}
	}
	return urm
}

func run(out *os.File, urm *murmur.Urm, maxSteps int, watch *watches,
	charsIsEmpty bool) {
	if err := urm.RunX(maxSteps); err != nil {
		dump(out, err, urm)
		onError(err)
	} else if charsIsEmpty {
		if err := writeWatched(out, urm, watch); err != nil {
			onError(err)
		}
	}
}

func step(out *os.File, urm *murmur.Urm, maxSteps int, watch *watches) {
	for {
		err := urm.Step()
		if err == murmur.ErrStop {
			break // normal termination
		} else if err != nil {
			dump(out, err, urm)
			onError(err) // abnormal termination
		} else if urm.Steps() >= maxSteps {
			fmt.Fprintf(os.Stderr, "forced to stop after %d steps\n",
				maxSteps)
			break
		} else {
			if err := writeWatched(out, urm, watch); err != nil {
				onError(err)
			}
		}
	}
}

func writeWatched(out *os.File, urm *murmur.Urm, watch *watches) error {
	regs := []string{}
	for _, label := range watch.labels {
		if s, err := urm.RegForLabelAsString(label); err != nil {
			return err
		} else {
			regs = append(regs, s)
		}
	}
	for _, reg := range watch.regs {
		if reg >= urm.Size() {
			break
		}
		if s, err := urm.RegAsString(reg); err != nil {
			return err
		} else {
			regs = append(regs, s)
		}
	}
	_, err := out.WriteString(fmt.Sprintf("#%d %s\n", urm.Steps(),
		strings.Join(regs, " ")))
	return err

}

func parseWatches(watches string) watches {
	watch := newWatches()
	for _, item := range strings.Split(watches, ",") {
		rs, ss, both := strings.Cut(item, "-")
		r, err := strconv.Atoi(strings.TrimSpace(rs))
		if both { // r-s
			if err != nil {
				onError(fmt.Errorf("expected int; got %s", rs))
				return watch
			}
			s, err := strconv.Atoi(strings.TrimSpace(ss))
			if err != nil {
				onError(fmt.Errorf("expected int; got %s", ss))
				return watch
			}
			for reg := r; reg <= s; reg++ {
				watch.regs = append(watch.regs, reg)
			}
		} else if err == nil { // r
			watch.regs = append(watch.regs, r)
		} else { // label
			if label := strings.TrimSpace(rs); label != "" {
				watch.labels = append(watch.labels, label)
			}
		}
	}
	return watch
}

type watches struct {
	labels []string
	regs   []int
}

func newWatches() watches {
	return watches{labels: []string{}, regs: []int{}}
}

func (me *watches) IsEmpty() bool {
	return (len(me.labels) + len(me.regs)) == 0
}

func dump(out *os.File, err error, urm *murmur.Urm) {
	_, _ = out.WriteString(fmt.Sprintf("error: %s\n%s\n", err,
		urm.StringWithRegNums()))
}

func writeChars(out *os.File, urm *murmur.Urm, chars watches) {
	runes := []rune{}
	for _, label := range chars.labels {
		if value, err := urm.RegValueForLabel(label); err == nil {
			runes = append(runes, charForInt(value))
		}
	}
	for _, reg := range chars.regs {
		if value, err := urm.RegValue(reg); err == nil {
			runes = append(runes, charForInt(value))
		}
	}
	runes = append(runes, '\n')
	_, _ = out.WriteString(string(runes))
}

func charForInt(i int) rune {
	r := rune(i)
	if unicode.IsPrint(r) {
		return r
	}
	return '.'
}

func onError(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(3)
}
