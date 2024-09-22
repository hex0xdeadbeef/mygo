package main

/*
	SMALL SIZE OBJECTS
1. What is "small size" objects?
	The "small-sized" object - is the objects that covers less than 4 machine words.
		- For 64-bit system the one machine word is 8 bytes.
		- For 32-bit system the one machine word is 4 bytes
2. Let's consider that we have a small-sized object including 4 ints and there's a single object that differs close to the firs one. We also have a function that sums the first and the last
	numbers in the definition ob the object.

	Small sized object's ASM code:
        TEXT    main.SumA(SB), NOSPLIT|NOFRAME|ABIInternal, $0-32
        FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        FUNCDATA        $5, main.SumA.arginfo1(SB)
        FUNCDATA        $6, main.SumA.argliveinfo(SB)
        PCDATA  $3, $1
        ADDQ    DI, AX
        RET
        TEXT    main.main(SB), NOSPLIT|NOFRAME|ABIInternal, $0-0
        FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        RET

	Big sized object's ASM code:
        TEXT    main.SumB(SB), NOSPLIT|NOFRAME|ABIInternal, $0-40
        FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        FUNCDATA        $5, main.SumB.arginfo1(SB)
        MOVQ    AX, main.b+8(SP)
        MOVQ    BX, main.b+16(SP)
        MOVQ    CX, main.b+24(SP)
        MOVQ    DI, main.b+32(SP)
        MOVQ    SI, main.b+40(SP)
        MOVQ    main.b+8(SP), CX
        LEAQ    (CX)(DI*1), AX
        RET
        TEXT    main.main(SB), NOSPLIT|NOFRAME|ABIInternal, $0-0
        FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        RET

	The key differences are:
		1) The sizes of stacks
			- SMALL: 32 bytes
			- BIG: 40 bytes
		2) Ops with registries:
			- SMALL: the simple operation with registries:
				ADDQ DI, AX
			- BIG: multiple operations that are directed to copy data from registries to the stack and back:
				MOVQ    AX, main.b+8(SP)
				MOVQ    BX, main.b+16(SP)
				MOVQ    CX, main.b+24(SP)
				MOVQ    DI, main.b+32(SP)
				MOVQ    SI, main.b+40(SP)
				MOVQ    main.b+8(SP), CX
		3) Stack usage
			- SMALL: SumA doesn't use the stack
			- BIG: Moves the data between registries and the stack because the big-sized object requires more memory size
		4) Arithmetic operations
			- SMALL: Performs the easier ops:
				ADDQ    DI, AX
			- BIG: Uses LEAQ (Load Effective Address) (more complex calculation)
				LEAQ    (CX)(DI*1), AX
3. The result of differentiating the SumA and SumB
	Considering two funcs we see that the function working with small-sized object consumps less memory than other one because the first uses only the registries not involving the transfer of the data from registries to the stack.

	Small-sized objects allows the compiler to generate more efficient code working with the CPU directively.
4. The exception
	We cannot design the systems working only with small-sized objects always.
*/

type A struct {
	a, b, c, d int
}

type B struct {
	a, b, c, d int
	e          int
}

func SumA(a A) int {
	return a.a + a.d
}

func SumB(b B) int {
	return b.a + b.d
}
