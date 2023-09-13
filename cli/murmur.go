// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/mark-summerfield/clip"
	"github.com/mark-summerfield/gong"
	"github.com/mark-summerfield/murmur"
)

func main() {
	config := getConfig()
	lines := readLines(config.infile)
	out, closer := getOutput(config.outfile)
	defer closer()
	urm := load(out, lines, config.registers, &config.watch)
	if config.step {
		step(out, urm, config.maxSteps, &config.watch)
	} else {
		run(out, urm, config.maxSteps, &config.watch)
	}
	if config.dis {
		_, _ = out.WriteString(urm.StringWithRegNums() + "\n")
	}
}

func getConfig() *Config {
	parser := clip.NewParserVersion(murmur.Version)
	parser.LongDesc = "Unlimited Register Machine emulator."
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
		"How many registers the URM has; can also be overridden in "+
			"the .urm with #n, e.g., #200.", 1, math.MaxInt,
		int(murmur.DefaultSize))
	watchOpt := parser.Str("watch", "A comma-separated list of which "+
		"registers to watch with each item of the form n or n-m or "+
		"label where n and m are integers ≥ 0. Out of range registers "+
		"and missing labels are ignored.", "PC,1-9")
	disOpt := parser.Flag("dis",
		"Show diassembly of all registers at the end.")
	parser.PositionalCount = clip.OnePositional
	parser.PositionalHelp = "The .urm file to run."
	if err := parser.Parse(); err != nil {
		parser.OnError(err) // doesn't return
		return nil          // never reached
	}
	registers := registersOpt.Value()
	config := Config{maxSteps: maxStepsOpt.Value(),
		step: stepOpt.Value(), registers: registers,
		infile: parser.Positionals[0], outfile: outfileOpt.Value(),
		watch: parseWatches(watchOpt.Value()), dis: disOpt.Value()}
	return &config
}

type Config struct {
	maxSteps  int
	step      bool
	registers int
	infile    string
	outfile   string
	watch     watches
	dis       bool
}

func (me Config) String() string {
	outfile := me.outfile
	if outfile == "" {
		outfile = "stdout"
	}
	return fmt.Sprintf(
		"maxsteps=%d step=%t registers=%d infile=%s outfile=%s watches=%v",
		me.maxSteps, me.step, me.registers, me.infile, outfile, me.watch)
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
			gong.ModeUserRW)
		if err != nil {
			onError(err)
		}
		closer = func() { _ = out.Close() }
	}
	return out, closer
}

func load(out *os.File, lines []string, registers int,
	watch *watches) *murmur.Urm {
	urm := murmur.New()
	if err := urm.Load(lines, registers); err != nil {
		onError(err)
	}
	if err := writeWatched(out, urm, watch); err != nil {
		onError(err)
	}
	return urm
}

func run(out *os.File, urm *murmur.Urm, maxSteps int, watch *watches) {
	if err := urm.RunX(maxSteps); err != nil {
		onError(err)
	} else {
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
		if s, err := urm.RegAsStringByLabel(label); err != nil {
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
		ns, ms, both := strings.Cut(item, "-")
		n, err := strconv.Atoi(strings.TrimSpace(ns))
		if both { // n-m
			if err != nil {
				onError(fmt.Errorf("expected int; got %s", ns))
				return watch
			}
			m, err := strconv.Atoi(strings.TrimSpace(ms))
			if err != nil {
				onError(fmt.Errorf("expected int; got %s", ms))
				return watch
			}
			for reg := n; reg <= m; reg++ {
				watch.regs = append(watch.regs, reg)
			}
		} else if err == nil { // n
			watch.regs = append(watch.regs, n)
		} else { // label
			watch.labels = append(watch.labels, strings.TrimSpace(ns))
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

func onError(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(3)
}
