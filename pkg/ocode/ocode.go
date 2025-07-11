package ocode

//go:generate go tool enumer -type=OcodeKind -json -text
type OcodeKind int

const (
	OpL        OcodeKind = iota // e.g.) L k / L c1,c2,c3 / などスタックにpushする
	OpOUT                       // OUT命令を追加
	OpDB                        // DB n_1,n_2,n_3,..,n_k
	OpDW                        // DW n_1,n_2,n_3,..,n_k
	OpDD                        // DD n_1,n_2,n_3,..,n_k
	OpRESB                      // RESB n または RESB n-$
	OpALIGNB                    // ALIGNB n (Added)
	OpMOV                       // MOV dst, src
	OpADD                       // ADD
	OpSUB                       // SUB (Added)
	OpAND                       // AND (Added)
	OpOR                        // OR
	OpXOR                       // XOR
	OpNOT                       // NOT
	OpSHR                       // SHR
	OpSHL                       // SHL
	OpSAR                       // SAR
	OpIN                        // IN (Unified)
	OpAAA                       // AAA
	OpAAD                       // AAD
	OpAAM                       // AAM
	OpAAS                       // AAS
	OpADX                       // ADX
	OpALTER                     // ALTER
	OpAMX                       // AMX
	OpCBW                       // CBW
	OpCDQ                       // CDQ
	OpCDQE                      // CDQE
	OpCLC                       // CLC
	OpCLD                       // CLD
	OpCLI                       // CLI
	OpCLTS                      // CLTS
	OpCMC                       // CMC
	OpCPUID                     // CPUID
	OpCQO                       // CQO
	OpCS                        // CS
	OpCWD                       // CWD
	OpCWDE                      // CWDE
	OpDAA                       // DAA
	OpDAS                       // DAS
	OpDIV                       // DIV
	OpDS                        // DS
	OpEMMS                      // EMMS
	OpENTER                     // ENTER
	OpES                        // ES
	OpF2XM1                     // F2XM1
	OpFABS                      // FABS
	OpFADDP                     // FADDP
	OpFCHS                      // FCHS
	OpFCLEX                     // FCLEX
	OpFCOM                      // FCOM
	OpFCOMP                     // FCOMP
	OpFCOMPP                    // FCOMPP
	OpFCOS                      // FCOS
	OpFDECSTP                   // FDECSTP
	OpFDISI                     // FDISI
	OpFDIVP                     // FDIVP
	OpFDIVRP                    // FDIVRP
	OpFENI                      // FENI
	OpFINCSTP                   // FINCSTP
	OpFINIT                     // FINIT
	OpFLD1                      // FLD1
	OpFLDL2E                    // FLDL2E
	OpFLDL2T                    // FLDL2T
	OpFLDLG2                    // FLDLG2
	OpFLDLN2                    // FLDLN2
	OpFLDPI                     // FLDPI
	OpFLDZ                      // FLDZ
	OpFMULP                     // FMULP
	OpFNCLEX                    // FNCLEX
	OpFNDISI                    // FNDISI
	OpFNENI                     // FNENI
	OpFNINIT                    // FNINIT
	OpFNOP                      // FNOP
	OpFNSETPM                   // FNSETPM
	OpFPATAN                    // FPATAN
	OpFPREM                     // FPREM
	OpFPREM1                    // FPREM1
	OpFPTAN                     // FPTAN
	OpFRNDINT                   // FRNDINT
	OpFRSTOR                    // FRSTOR
	OpFS                        // FS
	OpFSCALE                    // FSCALE
	OpFSETPM                    // FSETPM
	OpFSIN                      // FSIN
	OpFSINCOS                   // FSINCOS
	OpFSQRT                     // FSQRT
	OpFSUBP                     // FSUBP
	OpFSUBRP                    // FSUBRP
	OpFTST                      // FTST
	OpFUCOM                     // FUCOM
	OpFUCOMP                    // FUCOMP
	OpFUCOMPP                   // FUCOMPP
	OpFXAM                      // FXAM
	OpFXCH                      // FXCH
	OpFXRSTOR                   // FXRSTOR
	OpFXTRACT                   // FXTRACT
	OpFYL2X                     // FYL2X
	OpFYL2XP1                   // FYL2XP1
	OpGETSEC                    // GETSEC
	OpGS                        // GS
	OpHLT                       // HLT
	OpICEBP                     // ICEBP
	OpIDIV                      // IDIV
	OpIMUL                      // IMUL
	OpINT                       // INT
	OpCMP                       // CMP
	OpCALL                      // CALL
	OpLGDT                      // LGDT
	OpJMP                       // JMP
	OpJMP_FAR                   // JMP FAR (セグメント:オフセット)
	OpJA                        // JA
	OpJAE                       // JAE
	OpJB                        // JB
	OpJBE                       // JBE
	OpJC                        // JC
	OpJE                        // JE
	OpJG                        // JG
	OpJGE                       // JGE
	OpJL                        // JL
	OpJLE                       // JLE
	OpJNA                       // JNA
	OpJNAE                      // JNAE
	OpJNB                       // JNB
	OpJNBE                      // JNBE
	OpJNC                       // JNC
	OpJNE                       // JNE
	OpJNG                       // JNG
	OpJNGE                      // JNGE
	OpJNL                       // JNL
	OpJNLE                      // JNLE
	OpJNO                       // JNO
	OpJNP                       // JNP
	OpJNS                       // JNS
	OpJNZ                       // JNZ
	OpJO                        // JO
	OpJP                        // JP
	OpJPE                       // JPE
	OpJPO                       // JPO
	OpJS                        // JS
	OpJZ                        // JZ
	OpINTO                      // INTO
	OpINVD                      // INVD
	OpIRET                      // IRET
	OpIRETD                     // IRETD
	OpIRETQ                     // IRETQ
	OpJMPE                      // JMPE
	OpLAHF                      // LAHF
	OpLEAVE                     // LEAVE
	OpLFENCE                    // LFENCE
	OpLOADALL                   // LOADALL
	OpLOCK                      // LOCK
	OpMFENCE                    // MFENCE
	OpMONITOR                   // MONITOR
	OpMUL                       // MUL
	OpMWAIT                     // MWAIT
	OpNOP                       // NOP
	OpNTAKEN                    // NTAKEN
	OpPAUSE                     // PAUSE
	OpPOP                       // POP (Added)
	OpPOPA                      // POPA
	OpPOPAD                     // POPAD
	OpPOPF                      // POPF
	OpPOPFD                     // POPFD
	OpPOPFQ                     // POPFQ
	OpPUSH                      // PUSH (Added)
	OpPUSHA                     // PUSHA
	OpPUSHAD                    // PUSHAD
	OpPUSHF                     // PUSHF
	OpPUSHFD                    // PUSHFD
	OpPUSHFQ                    // PUSHFQ
	OpRDMSR                     // RDMSR
	OpRDPMC                     // RDPMC
	OpRDTSC                     // RDTSC
	OpRDTSCP                    // RDTSCP
	OpREP                       // REP
	OpREPE                      // REPE
	OpREPNE                     // REPNE
	OpRET                       // RET (Added)
	OpRETF                      // RETF
	OpRETN                      // RETN
	OpRSM                       // RSM
	OpSAHF                      // SAHF
	OpSETALC                    // SETALC
	OpSFENCE                    // SFENCE
	OpSS                        // SS
	OpSTC                       // STC
	OpSTD                       // STD
	OpSTI                       // STI
	OpSWAPGS                    // SWAPGS
	OpSYSCALL                   // SYSCALL
	OpSYSENTER                  // SYSENTER
	OpSYSEXIT                   // SYSEXIT
	OpSYSRET                    // SYSRET
	OpTAKEN                     // TAKEN
	OpUD2                       // UD2
	OpVMCALL                    // VMCALL
	OpVMLAUNCH                  // VMLAUNCH
	OpVMRESUME                  // VMRESUME
	OpVMXOFF                    // VMXOFF
	OpWAIT                      // WAIT
	OpWBINVD                    // WBINVD
	OpWRMSR                     // WRMSR
	OpXGETBV                    // XGETBV
	OpXRSTOR                    // XRSTOR
	OpXSETBV                    // XSETBV
)

type Ocode struct {
	Kind     OcodeKind
	Operands []string // 数値や変数名など
}
