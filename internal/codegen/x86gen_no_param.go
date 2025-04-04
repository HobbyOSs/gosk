package codegen

import "github.com/HobbyOSs/gosk/pkg/ocode"

var opcodeMap = map[ocode.OcodeKind]byte{
	ocode.OpAAA:      0x37,
	ocode.OpAAD:      0xD5,
	ocode.OpAAM:      0xD4,
	ocode.OpAAS:      0x3F,
	ocode.OpCBW:      0x98,
	ocode.OpCDQ:      0x99,
	ocode.OpCDQE:     0x98,
	ocode.OpCLC:      0xF8,
	ocode.OpCLD:      0xFC,
	ocode.OpCLI:      0xFA,
	ocode.OpCLTS:     0x06,
	ocode.OpCMC:      0xF5,
	ocode.OpCPUID:    0xA2,
	ocode.OpCQO:      0x99,
	ocode.OpCS:       0x2E,
	ocode.OpCWD:      0x99,
	ocode.OpCWDE:     0x98,
	ocode.OpDAA:      0x27,
	ocode.OpDAS:      0x2F,
	ocode.OpDIV:      0xF6,
	ocode.OpDS:       0x3E,
	ocode.OpEMMS:     0x77,
	ocode.OpENTER:    0xC8,
	ocode.OpES:       0x26,
	ocode.OpF2XM1:    0xD9,
	ocode.OpFABS:     0xD9,
	ocode.OpFADDP:    0xDE,
	ocode.OpFCHS:     0xD9,
	ocode.OpFCLEX:    0xDB,
	ocode.OpFCOM:     0xD8,
	ocode.OpFCOMP:    0xD8,
	ocode.OpFCOMPP:   0xDE,
	ocode.OpFCOS:     0xD9,
	ocode.OpFDECSTP:  0xD9,
	ocode.OpFDISI:    0xDB,
	ocode.OpFDIVP:    0xDE,
	ocode.OpFDIVRP:   0xDE,
	ocode.OpFENI:     0xDB,
	ocode.OpFINCSTP:  0xD9,
	ocode.OpFINIT:    0xDB,
	ocode.OpFLD1:     0xD9,
	ocode.OpFLDL2E:   0xD9,
	ocode.OpFLDL2T:   0xD9,
	ocode.OpFLDLG2:   0xD9,
	ocode.OpFLDLN2:   0xD9,
	ocode.OpFLDPI:    0xD9,
	ocode.OpFLDZ:     0xD9,
	ocode.OpFMULP:    0xDE,
	ocode.OpFNCLEX:   0xDB,
	ocode.OpFNDISI:   0xDB,
	ocode.OpFNENI:    0xDB,
	ocode.OpFNINIT:   0xDB,
	ocode.OpFNOP:     0xD9,
	ocode.OpFNSETPM:  0xDB,
	ocode.OpFPATAN:   0xD9,
	ocode.OpFPREM:    0xD9,
	ocode.OpFPREM1:   0xD9,
	ocode.OpFPTAN:    0xD9,
	ocode.OpFRNDINT:  0xD9,
	ocode.OpFRSTOR:   0xDD,
	ocode.OpFS:       0x64,
	ocode.OpFSCALE:   0xD9,
	ocode.OpFSETPM:   0xDB,
	ocode.OpFSIN:     0xD9,
	ocode.OpFSINCOS:  0xD9,
	ocode.OpFSQRT:    0xD9,
	ocode.OpFSUBP:    0xDE,
	ocode.OpFSUBRP:   0xDE,
	ocode.OpFTST:     0xD9,
	ocode.OpFUCOM:    0xDD,
	ocode.OpFUCOMP:   0xDD,
	ocode.OpFUCOMPP:  0xDA,
	ocode.OpFXAM:     0xD9,
	ocode.OpFXCH:     0xD9,
	ocode.OpFXRSTOR:  0xAE,
	ocode.OpFXTRACT:  0xD9,
	ocode.OpFYL2X:    0xD9,
	ocode.OpFYL2XP1:  0xD9,
	ocode.OpGETSEC:   0x37,
	ocode.OpGS:       0x65,
	ocode.OpHLT:      0xF4,
	ocode.OpICEBP:    0xF1,
	ocode.OpIDIV:     0xF6,
	ocode.OpINTO:     0xCE,
	ocode.OpINVD:     0x08,
	ocode.OpIRET:     0xCF,
	ocode.OpIRETD:    0xCF,
	ocode.OpIRETQ:    0xCF,
	ocode.OpJMPE:     0x00,
	ocode.OpLAHF:     0x9F,
	ocode.OpLEAVE:    0xC9,
	ocode.OpLFENCE:   0xAE,
	ocode.OpLOADALL:  0x05,
	ocode.OpLOCK:     0xF0,
	ocode.OpMFENCE:   0xAE,
	ocode.OpMONITOR:  0x01,
	ocode.OpMUL:      0xF6,
	ocode.OpMWAIT:    0x01,
	ocode.OpNOP:      0x90,
	ocode.OpPAUSE:    0x90,
	ocode.OpPOPA:     0x61,
	ocode.OpPOPAD:    0x61,
	ocode.OpPOPF:     0x9D,
	ocode.OpPOPFD:    0x9D,
	ocode.OpPOPFQ:    0x9D,
	ocode.OpPUSHA:    0x60,
	ocode.OpPUSHAD:   0x60,
	ocode.OpPUSHF:    0x9C,
	ocode.OpPUSHFD:   0x9C,
	ocode.OpPUSHFQ:   0x9C,
	ocode.OpRDMSR:    0x32,
	ocode.OpRDPMC:    0x33,
	ocode.OpRDTSC:    0x31,
	ocode.OpRDTSCP:   0x01,
	ocode.OpREP:      0xF2,
	ocode.OpREPE:     0xF3,
	ocode.OpREPNE:    0xF2,
	ocode.OpRETF:     0xCB,
	ocode.OpRETN:     0xC3,
	ocode.OpRSM:      0xAA,
	ocode.OpSAHF:     0x9E,
	ocode.OpSETALC:   0xD6,
	ocode.OpSFENCE:   0xAE,
	ocode.OpSS:       0x36,
	ocode.OpSTC:      0xF9,
	ocode.OpSTD:      0xFD,
	ocode.OpSTI:      0xFB,
	ocode.OpSWAPGS:   0x01,
	ocode.OpSYSCALL:  0x05,
	ocode.OpSYSENTER: 0x34,
	ocode.OpSYSEXIT:  0x35,
	ocode.OpSYSRET:   0x07,
	ocode.OpTAKEN:    0x3E,
	ocode.OpUD2:      0x0B,
	ocode.OpVMCALL:   0x01,
	ocode.OpVMLAUNCH: 0x01,
	ocode.OpVMRESUME: 0x01,
	ocode.OpVMXOFF:   0x01,
	ocode.OpWAIT:     0x9B,
	ocode.OpWBINVD:   0x09,
	ocode.OpWRMSR:    0x30,
	ocode.OpXGETBV:   0x01,
	ocode.OpXRSTOR:   0xAE,
	ocode.OpXSETBV:   0x01,
}

func GenerateX86NoParam(ocode ocode.Ocode) []byte {
	var binary []byte
	if code, exists := opcodeMap[ocode.Kind]; exists {
		binary = append(binary, code)
	}
	return binary
}
