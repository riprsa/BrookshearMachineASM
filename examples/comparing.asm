; This program compares two variables A and B;
; If A > B => r0 = 1
; If A < B => r0 = -1
; If A == B => r0 = 0

; Initialize variables
var A F0 0x05  ; Example value for A
var B F1 0x03  ; Example value for B

; Load variables into registers
load A r1      ; Load value of A into r1
load B r2      ; Load value of B into r2

; Compare A and B
xor r1 r2 r3   ; XOR A and B, result in r3
jeq EQUAL r3   ; If r3 is 0, A == B

; A != B, now check if A > B
mov r2 r0      ; Move B into r0 to set up for comparison
jgt GREATER r1 ; If A > B, jump to GREATER

; B > A, set r0 to 0xFF
loadi 0xFF r0  ; Load 0xFF into r0
halt           ; Stop the program

.label EQUAL:
loadi 0x00 r0  ; Load 0x00 into r0
halt           ; Stop the program

.label GREATER:
loadi 0x01 r0  ; Load 0x01 into r0
halt           ; Stop the program
