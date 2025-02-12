package junkjit

// すべての命令で機械語のサイズを返す
type Assembler interface {
	BufferData() []byte

	DB(x uint8, options ...Option) int
	DW(x uint16, options ...Option) int
	DD(x uint32, options ...Option) int
	DStruct(x any) int

	// no_param
	AAA() int
	AAD() int
	AAM() int
	AAS() int
	ADX() int
	ALTER() int
	AMX() int
	CBW() int
	CDQ() int
	CDQE() int
	CLC() int
	CLD() int
	CLI() int
	CLTS() int
	CMC() int
	CPUID() int
	CQO() int
	CS() int
	CWD() int
	CWDE() int
	DAA() int
	DAS() int
	DIV() int
	DS() int
	EMMS() int
	ENTER() int
	ES() int
	F2XM1() int
	FABS() int
	FADDP() int
	FCHS() int
	FCLEX() int
	FCOM() int
	FCOMP() int
	FCOMPP() int
	FCOS() int
	FDECSTP() int
	FDISI() int
	FDIVP() int
	FDIVRP() int
	FENI() int
	FINCSTP() int
	FINIT() int
	FLD1() int
	FLDL2E() int
	FLDL2T() int
	FLDLG2() int
	FLDLN2() int
	FLDPI() int
	FLDZ() int
	FMULP() int
	FNCLEX() int
	FNDISI() int
	FNENI() int
	FNINIT() int
	FNOP() int
	FNSETPM() int
	FPATAN() int
	FPREM() int
	FPREM1() int
	FPTAN() int
	FRNDINT() int
	FRSTOR() int
	FS() int
	FSCALE() int
	FSETPM() int
	FSIN() int
	FSINCOS() int
	FSQRT() int
	FSUBP() int
	FSUBRP() int
	FTST() int
	FUCOM() int
	FUCOMP() int
	FUCOMPP() int
	FXAM() int
	FXCH() int
	FXRSTOR() int
	FXTRACT() int
	FYL2X() int
	FYL2XP1() int
	GETSEC() int
	GS() int
	HLT() int
	ICEBP() int
	IDIV() int
	IMUL() int
	INTO() int
	INVD() int
	IRET() int
	IRETD() int
	IRETQ() int
	JMPE() int
	LAHF() int
	LEAVE() int
	LFENCE() int
	LOADALL() int
	LOCK() int
	MFENCE() int
	MONITOR() int
	MUL() int
	MWAIT() int
	NOP() int
	NTAKEN() int
	PAUSE() int
	POPA() int
	POPAD() int
	POPF() int
	POPFD() int
	POPFQ() int
	PUSHA() int
	PUSHAD() int
	PUSHF() int
	PUSHFD() int
	PUSHFQ() int
	RDMSR() int
	RDPMC() int
	RDTSC() int
	RDTSCP() int
	REP() int
	REPE() int
	REPNE() int
	RETF() int
	RETN() int
	// REX() int
	// REX.B() int
	// REX.R() int
	// REX.RB() int
	// REX.RX() int
	// REX.RXB() int
	// REX.W() int
	// REX.WB() int
	// REX.WR() int
	// REX.WRB() int
	// REX.WRX() int
	// REX.WRXB() int
	// REX.WX() int
	// REX.WXB() int
	// REX.X() int
	// REX.XB() int
	RSM() int
	SAHF() int
	SETALC() int
	SFENCE() int
	SS() int
	STC() int
	STD() int
	STI() int
	SWAPGS() int
	SYSCALL() int
	SYSENTER() int
	SYSEXIT() int
	SYSRET() int
	TAKEN() int
	UD2() int
	VMCALL() int
	VMLAUNCH() int
	VMRESUME() int
	VMXOFF() int
	WAIT() int
	WBINVD() int
	WRMSR() int
	XGETBV() int
	XRSTOR() int
	XSETBV() int
}

type Options struct {
	Count int
}

type Option func(*Options)

func Count(count int) Option {
	return func(opts *Options) {
		opts.Count = count
	}
}
