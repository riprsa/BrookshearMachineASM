# Brookshear Machine ASM

My simple assembly for Brookshear Machine (also known as VOLE machine), because I hate bytes.

This repo contains a compiler and a translator, both written in Go, and few examples of the ASM. Check out [How to use](#how-to-use) to try it.

## Registers

You can use following registers in your code:
> r0, r1, r2, r3, r4, r5, r6, r7
> r8, r9, rA, rB, rC, rD, rE, rF

## Counter

Brookshear Machine has a counter which points to current instruction. Because an instruction takes two bytes, while memory cell is only 1 byte, the counter is always even.

## Memory

Memory consists of 256 cells, each cell is 1 byte. Each cell has a `value` and an `address`. By definition, `value` and `address` are written as hex, but for `value` you have to write it with prefix `0x`. Example:
>`value`: `0xAF`
>`address`: `4D`

## Labels

Labels are a way to use a place in memory for jumps and variable declarations. Labels do not take space in actual memory. Labels are written in upper cased and should not contain spaces. Declaration looks like this:

`.label <NAME>:`

Prefer to use SNAKE_CASE for labels for readability.

## Variabels

Variable allows you to access a value of an address using its name. Variables change the mentioned space in memory. You have to declare them before any other commands. There are two types of variables. The first one called plain variable:

`var <name> <address> <value>`

The second one called labeled variable:

`var <name> <label>`

Prefer to use CamelCase for variabels for readability.

## Arrays

Arrays are similar to variables, they just edit a specific place in memory with a sequence of values:

`array <address> <value...>`

Note that `address` here should be in hex without `0x` prefix, while values should have this prefix

## Instructions

Load variable's value to a register: (opcode 1)
```assembly
var A F0 0x05
load A r5 ; now r5 contains hex value 5
```

Load immediate value to a register: (opcode 2)
```assembly
loadi 0x03 r5 ; now r5 contains hex value 3
```

Store a register to a variable: (opcode 3)
```assembly
var B F0 0x00
loadi 0xB3 r5
store r5 B ; now B contains hex value B3
```

Move a register to another register: (opcode 4)
```assembly
loadi 0x09 r1
move r1 r2 ; now r2 contains hex value 09
```

Two’s complement add of a register and another one, result to a third register: (opcode 5)
```assembly
var A F0 0x02
var B F1 0x03
load A r1
load B r2
add r1 r2 r0 ; now r0 contains sum of A and B => hex value 05
```

Float add of a register and another register, result to a third register: (opcode 6)
```assembly
var A F0 0x02
var B F1 0x03
load A r1
load B r2
addf r1 r2 r0 ; figure out it yourself, i never tested it lmao
```

Bitwise OR of a register and another register, result to a third register: (opcode 7)
```assembly
loadi 0xAA r1 ; 1010 1010
loadi 0x55 r2 ; 0101 0101
or r1 r2 r0 ; 1010 1010 or 0101 0101 => 1111 1111
```

Bitwise AND of a register and another register, result to a third register: (opcode 8)
```assembly
loadi 0xAA r1 ; 1010 1010
loadi 0x55 r2 ; 0101 0101
and r1 r2 r0 ; 1010 1010 and 0101 0101 => 0000 0000
```

Bitwise XOR of a register and another register, result to a third register: (opcode 9)
```assembly
loadi 0xAF r1 ; 1010 1111
loadi 0x5F r2 ; 0101 1111
xor r1 r2 r0 ; 1010 1111 and 0101 1111 => 1111 0000
```

Rotate the bit pattern in a register one bit to the right `number` times. Each time place the bit that started at the low-order end at the high-order end. `number` is a half-byte (one hex symbol): (opcode A)
```assembly
rotate r1 5 ; here you need to use value without 0x, because i am lazy to fix it
```

Change counter to a label if a register equals r0: (opcode B)
```assembly
.label MAIN:
mov r0 r0 ; beause label itself is not an instruction, it needs a dummy one to simulate action
jeq MAIN r0 ; this jump will loop forever...
halt
```

Halt execution. Ends program: (opcode C)
```
loadi 0x01 r0 ; smth useful
halt ; now stopped
```

Change counter to a label if a register is greater than r0: (opcode D)
```assembly
loadi 0x01 r1
loadi 0x02 r2
mov r1 r0 ; comparing works only with r0
jgt GREATER r2 ; if r2 > r1 ...
loadi 0xFF r0 ; if !(r2 > r1), write -1 to r0
halt ; ensure it is stopped
.label GREATER:
loadi 0x01 r0 ; if r2 > r1, write 1 to r0
halt
```

## Examples

You can find few examples in `./examples` folder. Please go thru them, because you will understand how to work with assembly more deeply. I also want to make insist you read `reverse.asm`, if you need to work with arrays/strings.

## How to use

Due to my laziness, you have to compile the program. Sorry, I won't provide binaries.

Clone repo, then run `go run main.go -asm ./examples/comparing.asm -o a.out`.

To use the compiler run `go run main.go -asm ./examples/comparing.asm -e`.

Change path to use your own `.asm` file, also check variable `DEBUG` in the `translator.go`