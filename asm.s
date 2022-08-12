
//func getModuleLoadedOrder(i int) (start uintptr, size uintptr)
TEXT ·getMLO(SB), $0-32
	//All operations push values into AX
	//PEB
	MOVQ 0x60(GS), AX
	BYTE $0x90			//NOP
	//PEB->LDR
	MOVQ 0x18(AX),AX
	BYTE $0x90			//NOP

	//LDR->InMemoryOrderModuleList
	MOVQ 0x20(AX),AX
	BYTE $0x90			//NOP

	//loop things
	XORQ R10,R10
startloop:
	CMPQ R10,i+0(FP)
	BYTE $0x90			//NOP
	JE endloop
	BYTE $0x90			//NOP
	//Flink (get next element)
	MOVQ (AX),AX
	BYTE $0x90			//NOP
	INCQ R10
	JMP startloop
endloop:
	//Flink - 0x10 -> _LDR_DATA_TABLE_ENTRY
	//_LDR_DATA_TABLE_ENTRY->DllBase (offset 0x30)

	MOVQ 0x30(AX),CX
	BYTE $0x90			//NOP
	MOVQ CX, size+16(FP)
	BYTE $0x90			//NOP


	MOVQ 0x20(AX),CX
	BYTE $0x90			//NOP
    MOVQ CX, start+8(FP)
    BYTE $0x90			//NOP


	MOVQ AX,CX
	BYTE $0x90			//NOP
	ADDQ $0x38,CX
	BYTE $0x90			//NOP
	MOVQ CX, modulepath+24(FP)
	//SYSCALL
	RET

