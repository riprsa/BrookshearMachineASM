; N is value that the program will look for in array ;
var N 3C 0x08

; P is place for index of N (answer). 0xFF (-1) if no answer ;
var P 3D 0x00

; Cap is amount of elements that the array expects ;
var Cap 3E 0xC0

; M is address of the first element of an array, used in 'array' keyword
var M 3F 0x40

; from M to M+Cap, second argument should be equal to M value
array 40 0x0F 0x0E 0x0D 0x0C 0x0B 0x0A 0x09 0x08 0x07 0x06 0x05 0x04 0x03 0x02 0x01 0x00

; it is possible to edit a part of memory using combination of labels and vars ;
var AddrOfElem HALF

; program ;
.label MAIN:
    jeq SEARCH_FUNC r0 ;jump right away to SEARCH_FUNC routine

.label END:
    loadi 0x00 rC;END. cleaning the regs after finishing the task...
    loadi 0x00 rB;this task was fun actually
    loadi 0x00 r1;i wrote my own 'assembly' with Go while solving this
    loadi 0x00 r2;but it did not help with the matrixes and the primes..
    loadi 0x00 r0;(╥﹏╥)
    loadi 0x00 rE;cleaning is done = )
    halt;bye-bye

; SEARCH_FUNC(N, M..) -> P ;
.label SEARCH_FUNC:
    ; for counter
    loadi 0x00 rC;SEARCH_FUNC. regC is counter
    loadi 0x01 rB;regB is 1 to add to counter later

    load Cap r0;reg0 contains capacity of array (192)
    load M r1;reg1 contains address of the array, M (64)

    .label LOOP:
        ; store r1 to label HALF as second byte
        store r1 AddrOfElem ;LOOP. store new address to next instruction as immediate value
        .label HALF:
        load AddrOfElem r2 ;load new address to reg2

        ; if r2 == N -> END ;
        mov r0 rE ;capacity to temporary register
        load N r0 ;given var N to reg0
        jeq ENDLOOPOK r2 ;if reg2.value == reg0.value: goto ENDLOOPOK
        mov rE r0 ;move back regE to reg0

        ; address of current element +1 ;
        add r1 rB r1 ;get next address

        ; counter of loop +1 ;
        add rC rB rC ;inc counter
        jeq ENDLOOPNOTOK rC ;if regC.value == reg0.value: goto ENDLOOPNOTOK
        jeq LOOP r0 ;else: goto LOOP

    .label ENDLOOPOK:
        store rC P ;regC to 
        jeq END r0 ;goto END

    .label ENDLOOPNOTOK:
        loadi 0xFF rC ;ENDLOOPNOTOK. -1 to regC
        store rC P ;regC to RAM cell of P
        jeq END r0 ;goto END
