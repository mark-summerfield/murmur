# Murmur

Murmur is an Unlimited Register Machine (URM) emulator.

This library and executable provide a "pure" URM. All instructions and data
are held in the registers. Register 0 holds the program counter (PC). The
urm "assembly" language means that only literal values need be entered as
numbers, everything else can use labels.

## Instructions

Note: The phrase _register x's value_ is shortened to _register x_
throughout.

_PC_ is the program counter (register 0). Every instruction changes the PC.

| Syntax       | PC       | Notes                                 |
| ------------ | -------- | ------------------------------------- |
| `C(n, m)`    | +3       | Copy register n into register m |
| `J(n, m, t)` | +4 or =t | If register n == register m, set PC to t (i.e., jump to t), else PC += 4|
| `J(t)`       | +4 or =t | Set PC to t (i.e., unconditional jump to t; sugar for `J(n, n, t)`)|
| `S(n)`       | +2       | Successor: increment register n (n++) |
| `Z(n)`       | +2       | Zero: set register n to 0 (n = 0) |

For the copy instruction, `T` (“transfer”), may be used instead of `C`.

Any instruction may be prefixed with a label.

In addition to the standard URM instructions, initial data values may be set
using the syntax `label: v`. This sets the “next” register's value to `v`
and sets the register's label to the given label. Or they may be set using
the syntax `reg: v` where `reg` is a register number (with the program
counter being 0), and `v` is the value.

Any register referred to by any instruction may be either a literal register
value or a label. In practice literals will normally only be needed to set
initial data values, with labels used everywhere else.

The exact number of registers to be available may be specified before any
data or instructions using the syntax `*r` where `r` is the number of
registers. If not specified this will default to 100.

## Example

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
*22 ; allocate 22 registers (0..21)
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
by the meta command, `*22`.

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

## License

GPL-3

---
