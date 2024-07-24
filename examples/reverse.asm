; reverses array at 60

array 60 0x01 0x02 0x03
;result at N: 0x03 0x02 0x01

var M F0 0x60 ; M
var Total F2 0x03 ; total amount of elements
var N F1 0x80 ; reversed goes here

var PlusCounter F3 0x00 ; PlusCounter
var MinusCounter F4 0x00 ; MinusCounter


var AddrOfElem HALF
var AddrOfElem2 HALF2

.label MAIN:
    load PlusCounter r5 ; load PlusCounter r5
    load MinusCounter r4 ; load MinusCounter r4
    loadi 0x01 r1 ; +1
    loadi 0xFF r3 ; -1
    load M r6 ; load M r6
    load N rA ; result arr

    load Total rF
    loadi 0xFF rE
    add rF rE rB ; -1 from Total
    store rF Total

    add rA rB r8 ; 0x60+Total

    jeq REVERSE r0 ;jeq REVERSE r0

.label END:
    halt ; END

.label REVERSE:
    add r5 r6 r7 ; REVERSE: sum PlusCounter with M

    store r7 AddrOfElem ; store PlusCounter+M to the next instr
    .label HALF:
        load AddrOfElem r2 ; load next value of array into r2

    ; 0x60 + Total (2) - PlusCounter (0)
    ; total is rf
    
    add r8 r4 rD ;

    store rD AddrOfElem2 ; store PlusCounter+M to the next instr
    .label HALF2:
        store r2 AddrOfElem2 ; load next value of array into r2
    
    add r5 r1 r5 ; rise PlusCounter
    add r4 r3 r4 ; decrease MinusCounter
    load Total r0 ; total to r0
    jeq END r5 ; jump to END Total with PlusCounter
    jeq REVERSE r0 ; jump back to start of function
