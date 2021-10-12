.data
cSpace:     .asciiz " "
cEndLine:   .asciiz "\n"
iArraySize: .word   10
iArray:     .word   12, 32, 13, 43, 17, 1, -2, -45, 0, 11

.text
main:
    # print integer array
    lw      $t0, iArraySize # load size of iArray
    la      $t1, iArray     # Load base address of iArray
    jal     print

    # Sort increase array
    lw      $t0, iArraySize # load size of iArray
    la      $t1, iArray     # Load base address of iArray
    jal     sortIncrease

    # print integer array
    lw      $t0, iArraySize # load size of iArray
    la      $t1, iArray     # Load base address of iArray
    jal     print

    #stop program
    li      $v0, 10
    syscall

print:                      # print fuction

    add     $t2, $0, $0     # index of iArray
loopPrint:
    beqz    $t0, exitPrint  # Check condition
    li      $v0, 1          # service 1 is print integer
    add     $t3, $t1, $t2   # load desired value into $a0
    lw      $a0, ($t3) 
    syscall

    li      $v0, 4
    la      $a0, cSpace     # print space just like separator
    syscall

    addi    $t0, $t0, -1    # decrease loop count
    addi    $t2, $t2, 4     # increase index
    b       loopPrint
exitPrint:
    li      $v0, 4
    la      $a0, cEndLine   # print end line
    syscall
    jr      $ra             # end of print

# Sort increase function
sortIncrease:
    addi    $t0, $t0, -1        # i counter
    add     $t2, $0, $0         # a[i] address

loopI:
    beqz    $t0, return         # return

    add     $t3, $t1, $t2       # load base address a[i]
    lw      $t4, ($t3)

    add     $t5, $0, $t0        # j counter
    add     $t6, $t2, 4         # j index

loopJ:
    beqz    $t5, increaseI

    add     $t7, $t1, $t6       # load base address a[j]
    lw      $t8, ($t7)

    slt     $t9, $t4, $t8       # compare
    beqz    $t9, swap

    b       increaseJ

increaseI:
    addi    $t0, $t0, -1        # decrease i counter
    addi    $t2, $t2, 4         # increate a[i] index
    b       loopI
increaseJ:
    addi    $t5, $t5, -1        # decrease j counter
    addi    $t6, $t6, 4         # increate a[j] index
    b       loopJ
swap:
    # swap value on register
    add     $t9, $0, $t8
    add     $t8, $0, $t4
    add     $t4, $0, $t9

    # store value from register to memory
    sw      $t4, ($t3)          # set a[j] to a[j]
    sw      $t8, ($t7)          # set a[i] to a[i]
    b       increaseJ

return:
    jr      $ra
