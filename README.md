# Murmur

Murmur is an Unlimited Register Machine (URM) emulator.

- [Introduction](#introduction)
- [Instructions](#instructions)
    - [Indirect Addressing Extension](#indirect-addressing-extension)
    - [Comparison Jumps Extensions](#comparison-jumps-extensions)
    - [Simple Arithmetic Extensions](#simple-arithmetic-extensions)
- [Examples](#examples)
    - [Addition](#addition)
    - [Other Examples](#other-examples)
    - [Sorting](#sorting)
- [Gvim](#gvim)
- [License](#license)

## Introduction

This library and command line executable provide a “pure” URM. All
instructions and data are held in the registers. Register 0 holds the
program counter (PC). The urm “assembly” language means that in most cases
only literal register values need be entered as numbers, everything else can
use labels.

The command line executable is built using `go build` as usual, although a
couple of build scripts are included. Run `./murmur -h` to see the command
line usage.

A pre-built binary (`murmur.exe` 2.9Mb;
MD5: 5a318ddb9ee9ab9d13dfa40afd7d2d08)
is provided as a convenience for Windows users.

## Instructions

Note: The phrase _register x's value_ is shortened to _register x_
throughout.

_PC_ is the program counter (register 0). Every instruction changes the PC.

| Syntax       | PC       | Notes                                 |
| ------------ | -------- | ------------------------------------- |
| `C(r,s)`    | +3       | Copy register r into register s       |
| `J(r,s,t)` | +4 or =t | If register r == register s, set PC to t (i.e., jump to t), else PC += 4|
| `J(t)`       | +4 or =t | Set PC to t (i.e., unconditional jump to t; sugar for `J(r,r,t)`)|
| `S(r)`       | +2       | Successor: increment register r (r++) |
| `Z(r)`       | +2       | Zero: set register r to 0 (r = 0) |

Instructions are _case-insensitive_. For the copy instruction, `T`
(“transfer”), may be used instead of `C`, and for the successor instruction,
`I` (“increment”), may be used instead of `S`.

Any instruction may be prefixed with a _case-sensitive_ label.

In addition to the standard URM instructions, initial data values may be
given. Two syntaxes are supported. First, `reg: v` where `reg` is a register
number (with the program counter being 0), and `v` is the value. The second
is more versatile, with the synax `label: v`. This sets the “next”
register's value to `v` and sets the register's label to the given label.
This syntax can actually accept any number of values, e.g., `label: v1 v2 v3
…`, setting the “next” register's value to `v1`, and the register after that
to `v2`, and so on. Furthermore, it is possible to set a label's value to
`HERE`, in which case the value stored for that label is the current
register number. This is useful to mark the end of data items. (See for
example, `eg/index.urm` and `eg/index2.urm`.)

Any register referred to by any instruction may be either a literal register
value or a label. In practice literals should normally only be needed to set
initial data values, with labels used everywhere else.

The exact number of registers to be available may be specified before any
data or instructions using the syntax `^n` where `n` is the number of
registers. If not specified this will default to 100.

All values must be `0` or a positive integer.

_All bets are off if a resulting value would be `< 0`._

Some extensions are supported as described below. Simply any of them (or all
of them) that aren't wanted.

### Indirect Addressing Extension

A significant inconvenience of the basic URM is that all addressing is
direct. This means that it isn't really possible to loop over a range of
registers using some kind of index counter. To address this issue, `murmur`
supports a URM extension: indirect addressing. Using this extension makes it
possible to access a register whose register number is stored in another
register.

For example:

    A: 0      ; sets register 1 to 0 and labels it A; not used ; 1
    B: 6      ; sets register 2 to 6 and labels it B; 2
    INDEX: 1  ; sets register 3 to 1 and labels it INDEX; 3
    S(2)      ; B++ ∴ B = 7; increments register 2, i.e., B ; 4
    S(B)      ; B++ ∴ B = 8; increments B, i.e., register 2 ; 6
    S(INDEX)  ; I++ ∴ I = 2; increments INDEX, i.e., register 3 ; 8
    S(@INDEX) ; B++ ∴ B = 9; increments register 2, i.e., B via INDEX; 10

Here we can see register 2—labelled `B`—incremented three times in three
different ways. First by register number, then by label, and at the end,
indirectly, where the `S` command is applied to `@INDEX`. The `@` label
qualifier says to use the register whose number is held by the register
whose label follows the `@`. So, in this case, `@INDEX` means look at the
value held in the `INDEX` register (register 3). The value is 2 (it was
originally set to 1 but the `S(INDEX)` command incremented it). So actually
apply the command to register 2 (i.e., `B`).

| Syntax       | Notes                                            |
| ------------ | ------------------------------------------------ |
| `@LABEL`     | Use the register whose number is held in `LABEL` |

For an example of use see `eg/index.urm` which uses indirect indexing to
iterate over an “array” of five registers, adding 2 to each one. (The
example, `eg/index2.urm` is the same except it uses multiple data values on
a single labelled line.)

### Comparison Jumps Extensions

| Syntax       | PC       | Notes                                 |
| ------------ | -------- | ------------------------------------- |
| `G(r,s,t)` | +4 or =t | If register r > register s, set PC to t (i.e., jump to t), else PC += 4|
| `L(r,s,t)` | +4 or =t | If register r < register s, set PC to t (i.e., jump to t), else PC += 4|

### Simple Arithmetic Extensions

| Syntax       | PC       | Notes                                    |
| ------------ | -------- | ---------------------------------------- |
| `P(r)`       | +2       | Prececessor: decrement register r (r-\-) |
| `+(r,s)`     | +2       | Add: add register s to r (r = r + s)     |
| `-(r,s)`     | +2       | Subtract: subtract s from r (r = r - s)  |
| `*(r,s)`     | +2       | Multiply: multiply r by s (r = r * s)    |
| `/(r,s)`     | +2       | Divide: divide r by s (r = r / s)        |

For the precedessor instruction, `D` (“decrement”), may be used instead of
`P`.

Division is truncating.

## Examples

The most interesting examples are the [Sorting](#sorting) ones, but they all
depend on using the extensions.

### Addition

Here is a program (`eg/26.urm`) to add register 1 (labelled `A` and
containing 19) to register 2 (labelled `B`, containing 7), using register 3
(labelled `C` as a scratch value), and storing their sum back into register
1 (`A`). Furthermore, at the end, registers 2 and 3 (`B` and `C`) are
cleared.

Note that the only literal values needed are for setting the initial data:
all other registers are referred to by label.

```
; Addition: A = A + B
^22 ; allocate 22 registers (0..21)
A:     19 ; 1 — at the end this will be 26
B:     7  ; 2 — at the end this will be 0
C:     0  ; 0 — at the end this will be 0
START: J(B, C, END) ; 4
       S(A) ; 8
       S(C) ; 10
       J(START) ; 12
END:   Z(B) ; 16
       Z(C) ; 18
       STOP ; 20
```

Here, 22 registers will be available (the exact number needed) as specified
by the meta command, `^22`.

The register number is shown in a comment at the end of each line.

The first three lines label registers 1, 2, and 3 as `A`, `B`, and `C` and
set their values to 19, 7, and 0.

The line labelled `START` is the first line with a command. The urm will
begin execution from the register labelled `START` or if there isn't one,
from the first register containing a command. So in this case, with or
without the `START` label, execution would begin from register 4.

Program execution stops when the program counter (PC) is set to 0. The
`STOP` meta-command is syntactic sugar for `Z(PC)`, itself syntactic sugar
for `Z(0)`. (This is why we need 22 registers, register 20 for the `Z`
command and register 21 for the `0` value.)

It should be easy to follow how the program works. Or use the command line
`murmur` program wih the options `-s` (step) and `-wPC,A,B,C` to see all the
steps and how the program counter and first three registers change at each
step. Or use `-d` (dis-assemble) to see the registers before and after the
run.

If the extensions are used the above can be achieved like this:

```
; Add A = A + B
A: 19
B: 7
   +(A, B)
   STOP
```

### Other Examples

Most of the other numbered examples are commented. There are also named
examples including, `add.urm`, `sub.urm`, `mul.urm`, `lt.urm`, `lte.urm`,
and `max.urm`.

For example, to add two numbers, run, say, `./murmur -wA eg/add.urm 23 91`.
The `-wA` option says to watch register `A`. The numbers given as arguments
to the `.urm` program set registers 1 and 2 (`A` and `B`) and store the
result in `A`. This will output:

```
#0 A:23
#368 A:114
```

Or to find the maximum of two numbers, run, say,
`./murmur -wA eg/max.urm 0 65 82`.
The numbers set registers 1, 2, and 3 (`A`, `B`, and `C`) and store the
maximum of `B` and `C` into `A`. This will output:

```
#0 A:0
#331 A:82
```

The algorithm for the `lt.urm` (“less than”) example is almost the same as
`max.urm` only instead of copying the max of the two numbers into `A`, it
sets `A` to 99 (“invalid”) at the start, and at the end sets `A` to 1 if `B
< C`, otherwise to 0 (i.e., if `B ≥ C`).

_Note that the examples that match `*x.urm` (except for `max.urm`) use one
or more of the extensions listed above._

### Sorting

Using the indirect addressing extension alone is sufficient to implement
sort algorithms with a reasonable number of instructions. And using the
other extensions can make the resulting implementations even more compact.

See the `eg/*sort*.urm` examples.

## Gvim

Gvim users might wish to copy the file `vim.urm` to their `vim/syntax`
folder and add a line like the following to their `.gvimrc` file:

    au BufRead,BufNewFile,BufEnter *.urm set ft=urm

## License

GPL-3

---
