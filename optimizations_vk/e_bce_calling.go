package main

/*
	BCE CALLING
1. Go is safe language that prevents us from unsafe situations. A participant of this action the BCE (Bounds Check Elimination). This method impacts the efficiency dramatically.
2. The code that performs the BCE method at each clobbering the specific index:
	func toUint64A(b []byte) uint64 {
		return uint64(b[0]) |
			uint64(b[1])<<8 |
			uint64(b[2])<<16 |
			uint64(b[3])<<24 |
			uint64(b[4])<<32 |
			uint64(b[5])<<40 |
			uint64(b[6])<<48 |
			uint64(b[7])<<56
	}

	The Assembler code of this function:
3.
	TEXT    main.main(SB), NOSPLIT|NOFRAME|ABIInternal, $0-0
			FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
			FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
			RET
			TEXT    main.toUint64(SB), NOSPLIT|ABIInternal, $24-24
			PUSHQ   BP
			MOVQ    SP, BP
			SUBQ    $16, SP
			MOVQ    AX, main.b+32(FP)
			FUNCDATA        $0, gclocals·wgcWObbY2HYnK2SU/U22lA==(SB)
			FUNCDATA        $1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
			FUNCDATA        $5, main.toUint64.arginfo1(SB)
			FUNCDATA        $6, main.toUint64.argliveinfo(SB)
			PCDATA  $3, $1
			TESTQ   BX, BX
			JLS     main_toUint64_pc185
			NOP
			CMPQ    BX, $1
			JLS     main_toUint64_pc172
			CMPQ    BX, $2
			JLS     main_toUint64_pc159
			CMPQ    BX, $3
			JLS     main_toUint64_pc146
			CMPQ    BX, $4
			JLS     main_toUint64_pc133
			NOP
			CMPQ    BX, $5
			JLS     main_toUint64_pc117
			CMPQ    BX, $6
			JLS     main_toUint64_pc104
			CMPQ    BX, $7
			JLS     main_toUint64_pc91
			MOVQ    (AX), AX
			ADDQ    $16, SP
			POPQ    BP
			RET
	main_toUint64_pc91:
			MOVL    $7, AX
			MOVQ    AX, CX
			PCDATA  $1, $1
			CALL    runtime.panicIndex(SB)
	main_toUint64_pc104:
			MOVL    $6, AX
			MOVQ    AX, CX
			CALL    runtime.panicIndex(SB)
	main_toUint64_pc117:
			MOVL    $5, AX
			MOVQ    AX, CX
			NOP
			CALL    runtime.panicIndex(SB)
	main_toUint64_pc133:
			MOVL    $4, AX
			MOVQ    AX, CX
			CALL    runtime.panicIndex(SB)
	main_toUint64_pc146:
			MOVL    $3, AX
			MOVQ    AX, CX
			CALL    runtime.panicIndex(SB)
	main_toUint64_pc159:
			MOVL    $2, AX
			MOVQ    AX, CX
			CALL    runtime.panicIndex(SB)
	main_toUint64_pc172:
			MOVL    $1, AX
			MOVQ    AX, CX
			CALL    runtime.panicIndex(SB)
	main_toUint64_pc185:
			XORL    AX, AX
			MOVQ    AX, CX
			NOP
			CALL    runtime.panicIndex(SB)
			XCHGL   AX, AX
4. To suppres the panic-checks we can use the trick with referencing to the last index of the slice.

	The Assembler code of the function ToUint64B:
		TEXT    main.main(SB), NOSPLIT|NOFRAME|ABIInternal, $0-0
		FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
		FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
		RET
		TEXT    main.ToUint64B(SB), NOSPLIT|ABIInternal, $24-24
		PUSHQ   BP
		MOVQ    SP, BP
		SUBQ    $16, SP
		MOVQ    AX, main.b+32(FP)
		FUNCDATA        $0, gclocals·wgcWObbY2HYnK2SU/U22lA==(SB)
		FUNCDATA        $1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
		FUNCDATA        $5, main.ToUint64B.arginfo1(SB)
		FUNCDATA        $6, main.ToUint64B.argliveinfo(SB)
		PCDATA  $3, $1
		CMPQ    BX, $7
		JLS     main_ToUint64B_pc28
		MOVQ    (AX), AX
		ADDQ    $16, SP
		POPQ    BP
		RET
		main_ToUint64B_pc28:
		MOVL    $7, AX
		MOVQ    BX, CX
		PCDATA  $1, $1
		CALL    runtime.panicIndex(SB)
		XCHGL   AX, AX
5. As we can see, by clobbering the last index we descend the number of Assembler instructions.
*/

// This function unwinds into the many BCE checks, it results in performance falling.
func ToUint64A(b []byte) uint64 {
	return uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56
}

// Clobbering the last index
func ToUint64B(b []byte) uint64 {
	const (
		lastIndex = 7
	)
	_ = b[lastIndex]
	return uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56
}
