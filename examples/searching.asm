var N 3C 0x08 ; N is value that the program will look for in array ;
var P 3D 0x00 ; P is place for index of N (answer). 0xFF (-1) if no answer ;
var Cap 3E 0xC0 ; Cap is amount of elements that the array expects ;
var M 3F 0x40 ; M is address of the first element of an array, used in 'array' keyword
array 40 0x0F 0x0E 0x0D 0x0C 0x0B 0x0A 0x09 0x08 0x07 0x06 0x05 0x04 0x03 0x02 0x01 0x00 ; from M to M+Cap, second argument should be equal to M value
var AddrOfElem HALF ; it is possible to edit a part of memory using combination of labels and vars ;


; program ;
.label MAIN:
    jeq SEARCH_FUNC r0

.label END:
    loadi 0x00 rC
    loadi 0x00 rB
    loadi 0x00 r1
    loadi 0x00 r2
    loadi 0x00 r0
    loadi 0x00 rE
    halt

; SEARCH_FUNC(N, M..) -> P ;
.label SEARCH_FUNC:
    loadi 0x00 rC ; for counter
    loadi 0x01 rB

    load Cap r0
    load M r1

    .label LOOP:
        store r1 AddrOfElem ; store r1 to label HALF as second byte
        .label HALF:
        load AddrOfElem r2

        mov r0 rE
        load N r0
        jeq ENDLOOPOK r2 ; if r2 == N -> END ;
        mov rE r0

        add r1 rB r1 ; address of current element +1 ;

        add rC rB rC ; counter of loop +1 ;
        jeq ENDLOOPNOTOK rC
        jeq LOOP r0

    .label ENDLOOPOK:
        add rB rC rC
        loadi 0xFF rB
        add rB rC rC
        store rC P
        jeq END r0

    .label ENDLOOPNOTOK:
        loadi 0xFF rC
        store rC P
        jeq END r0
