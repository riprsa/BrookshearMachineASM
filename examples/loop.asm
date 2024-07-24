; example of a sum function and a loop function

; vars: C0-FF ;
var A C0 0x05
var B C1 0x02
var Sum C2 0x00

; main ;
.label Main
    load A r1
    load B r2
    jeq SumFunc r0
.label SumFuncReturn:
    loadi 0x00 rA
    loadi 0x01 rF
    loadi 0x05 r0
    jeq LoopFunc r0
.label LoopFuncReturn:
    halt
; main end ;

; SumFunc(r1, r2) -> r3 ;
.label SumFunc:
    add r1 r2 r3
    store r3 Sum
    jeq SumFuncReturn r0
; SumFunc end ;

; LoopFunc(factor: rF, counter: rA, until: r0) -> rA ;
.label LoopFunc:
    .label Loop:
        add rF rA rA
        jeq EndLoop rA
        jeq Loop r0
    .label EndLoop:
        jeq LoopFuncReturn r0
; LoopFunc end ;