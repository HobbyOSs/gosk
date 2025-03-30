// Code generated by "enumer -type=OcodeKind -json -text"; DO NOT EDIT.

package ocode

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _OcodeKindName = "OpLOpOUTOpDBOpDWOpDDOpRESBOpMOVOpADDOpSUBOpANDOpOROpXOROpNOTOpSHROpSHLOpSAROpINOpAAAOpAADOpAAMOpAASOpADXOpALTEROpAMXOpCBWOpCDQOpCDQEOpCLCOpCLDOpCLIOpCLTSOpCMCOpCPUIDOpCQOOpCSOpCWDOpCWDEOpDAAOpDASOpDIVOpDSOpEMMSOpENTEROpESOpF2XM1OpFABSOpFADDPOpFCHSOpFCLEXOpFCOMOpFCOMPOpFCOMPPOpFCOSOpFDECSTPOpFDISIOpFDIVPOpFDIVRPOpFENIOpFINCSTPOpFINITOpFLD1OpFLDL2EOpFLDL2TOpFLDLG2OpFLDLN2OpFLDPIOpFLDZOpFMULPOpFNCLEXOpFNDISIOpFNENIOpFNINITOpFNOPOpFNSETPMOpFPATANOpFPREMOpFPREM1OpFPTANOpFRNDINTOpFRSTOROpFSOpFSCALEOpFSETPMOpFSINOpFSINCOSOpFSQRTOpFSUBPOpFSUBRPOpFTSTOpFUCOMOpFUCOMPOpFUCOMPPOpFXAMOpFXCHOpFXRSTOROpFXTRACTOpFYL2XOpFYL2XP1OpGETSECOpGSOpHLTOpICEBPOpIDIVOpIMULOpINTOpCMPOpCALLOpLGDTOpJMPOpJAOpJAEOpJBOpJBEOpJCOpJEOpJGOpJGEOpJLOpJLEOpJNAOpJNAEOpJNBOpJNBEOpJNCOpJNEOpJNGOpJNGEOpJNLOpJNLEOpJNOOpJNPOpJNSOpJNZOpJOOpJPOpJPEOpJPOOpJSOpJZOpINTOOpINVDOpIRETOpIRETDOpIRETQOpJMPEOpLAHFOpLEAVEOpLFENCEOpLOADALLOpLOCKOpMFENCEOpMONITOROpMULOpMWAITOpNOPOpNTAKENOpPAUSEOpPOPAOpPOPADOpPOPFOpPOPFDOpPOPFQOpPUSHAOpPUSHADOpPUSHFOpPUSHFDOpPUSHFQOpRDMSROpRDPMCOpRDTSCOpRDTSCPOpREPOpREPEOpREPNEOpRETOpRETFOpRETNOpRSMOpSAHFOpSETALCOpSFENCEOpSSOpSTCOpSTDOpSTIOpSWAPGSOpSYSCALLOpSYSENTEROpSYSEXITOpSYSRETOpTAKENOpUD2OpVMCALLOpVMLAUNCHOpVMRESUMEOpVMXOFFOpWAITOpWBINVDOpWRMSROpXGETBVOpXRSTOROpXSETBV"

var _OcodeKindIndex = [...]uint16{0, 3, 8, 12, 16, 20, 26, 31, 36, 41, 46, 50, 55, 60, 65, 70, 75, 79, 84, 89, 94, 99, 104, 111, 116, 121, 126, 132, 137, 142, 147, 153, 158, 165, 170, 174, 179, 185, 190, 195, 200, 204, 210, 217, 221, 228, 234, 241, 247, 254, 260, 267, 275, 281, 290, 297, 304, 312, 318, 327, 334, 340, 348, 356, 364, 372, 379, 385, 392, 400, 408, 415, 423, 429, 438, 446, 453, 461, 468, 477, 485, 489, 497, 505, 511, 520, 527, 534, 542, 548, 555, 563, 572, 578, 584, 593, 602, 609, 618, 626, 630, 635, 642, 648, 654, 659, 664, 670, 676, 681, 685, 690, 694, 699, 703, 707, 711, 716, 720, 725, 730, 736, 741, 747, 752, 757, 762, 768, 773, 779, 784, 789, 794, 799, 803, 807, 812, 817, 821, 825, 831, 837, 843, 850, 857, 863, 869, 876, 884, 893, 899, 907, 916, 921, 928, 933, 941, 948, 954, 961, 967, 974, 981, 988, 996, 1003, 1011, 1019, 1026, 1033, 1040, 1048, 1053, 1059, 1066, 1071, 1077, 1083, 1088, 1094, 1102, 1110, 1114, 1119, 1124, 1129, 1137, 1146, 1156, 1165, 1173, 1180, 1185, 1193, 1203, 1213, 1221, 1227, 1235, 1242, 1250, 1258, 1266}

const _OcodeKindLowerName = "oplopoutopdbopdwopddopresbopmovopaddopsubopandoporopxoropnotopshropshlopsaropinopaaaopaadopaamopaasopadxopalteropamxopcbwopcdqopcdqeopclcopcldopcliopcltsopcmcopcpuidopcqoopcsopcwdopcwdeopdaaopdasopdivopdsopemmsopenteropesopf2xm1opfabsopfaddpopfchsopfclexopfcomopfcompopfcomppopfcosopfdecstpopfdisiopfdivpopfdivrpopfeniopfincstpopfinitopfld1opfldl2eopfldl2topfldlg2opfldln2opfldpiopfldzopfmulpopfnclexopfndisiopfneniopfninitopfnopopfnsetpmopfpatanopfpremopfprem1opfptanopfrndintopfrstoropfsopfscaleopfsetpmopfsinopfsincosopfsqrtopfsubpopfsubrpopftstopfucomopfucompopfucomppopfxamopfxchopfxrstoropfxtractopfyl2xopfyl2xp1opgetsecopgsophltopicebpopidivopimulopintopcmpopcalloplgdtopjmpopjaopjaeopjbopjbeopjcopjeopjgopjgeopjlopjleopjnaopjnaeopjnbopjnbeopjncopjneopjngopjngeopjnlopjnleopjnoopjnpopjnsopjnzopjoopjpopjpeopjpoopjsopjzopintoopinvdopiretopiretdopiretqopjmpeoplahfopleaveoplfenceoploadalloplockopmfenceopmonitoropmulopmwaitopnopopntakenoppauseoppopaoppopadoppopfoppopfdoppopfqoppushaoppushadoppushfoppushfdoppushfqoprdmsroprdpmcoprdtscoprdtscpoprepoprepeoprepneopretopretfopretnoprsmopsahfopsetalcopsfenceopssopstcopstdopstiopswapgsopsyscallopsysenteropsysexitopsysretoptakenopud2opvmcallopvmlaunchopvmresumeopvmxoffopwaitopwbinvdopwrmsropxgetbvopxrstoropxsetbv"

func (i OcodeKind) String() string {
	if i < 0 || i >= OcodeKind(len(_OcodeKindIndex)-1) {
		return fmt.Sprintf("OcodeKind(%d)", i)
	}
	return _OcodeKindName[_OcodeKindIndex[i]:_OcodeKindIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _OcodeKindNoOp() {
	var x [1]struct{}
	_ = x[OpL-(0)]
	_ = x[OpOUT-(1)]
	_ = x[OpDB-(2)]
	_ = x[OpDW-(3)]
	_ = x[OpDD-(4)]
	_ = x[OpRESB-(5)]
	_ = x[OpMOV-(6)]
	_ = x[OpADD-(7)]
	_ = x[OpSUB-(8)]
	_ = x[OpAND-(9)]
	_ = x[OpOR-(10)]
	_ = x[OpXOR-(11)]
	_ = x[OpNOT-(12)]
	_ = x[OpSHR-(13)]
	_ = x[OpSHL-(14)]
	_ = x[OpSAR-(15)]
	_ = x[OpIN-(16)]
	_ = x[OpAAA-(17)]
	_ = x[OpAAD-(18)]
	_ = x[OpAAM-(19)]
	_ = x[OpAAS-(20)]
	_ = x[OpADX-(21)]
	_ = x[OpALTER-(22)]
	_ = x[OpAMX-(23)]
	_ = x[OpCBW-(24)]
	_ = x[OpCDQ-(25)]
	_ = x[OpCDQE-(26)]
	_ = x[OpCLC-(27)]
	_ = x[OpCLD-(28)]
	_ = x[OpCLI-(29)]
	_ = x[OpCLTS-(30)]
	_ = x[OpCMC-(31)]
	_ = x[OpCPUID-(32)]
	_ = x[OpCQO-(33)]
	_ = x[OpCS-(34)]
	_ = x[OpCWD-(35)]
	_ = x[OpCWDE-(36)]
	_ = x[OpDAA-(37)]
	_ = x[OpDAS-(38)]
	_ = x[OpDIV-(39)]
	_ = x[OpDS-(40)]
	_ = x[OpEMMS-(41)]
	_ = x[OpENTER-(42)]
	_ = x[OpES-(43)]
	_ = x[OpF2XM1-(44)]
	_ = x[OpFABS-(45)]
	_ = x[OpFADDP-(46)]
	_ = x[OpFCHS-(47)]
	_ = x[OpFCLEX-(48)]
	_ = x[OpFCOM-(49)]
	_ = x[OpFCOMP-(50)]
	_ = x[OpFCOMPP-(51)]
	_ = x[OpFCOS-(52)]
	_ = x[OpFDECSTP-(53)]
	_ = x[OpFDISI-(54)]
	_ = x[OpFDIVP-(55)]
	_ = x[OpFDIVRP-(56)]
	_ = x[OpFENI-(57)]
	_ = x[OpFINCSTP-(58)]
	_ = x[OpFINIT-(59)]
	_ = x[OpFLD1-(60)]
	_ = x[OpFLDL2E-(61)]
	_ = x[OpFLDL2T-(62)]
	_ = x[OpFLDLG2-(63)]
	_ = x[OpFLDLN2-(64)]
	_ = x[OpFLDPI-(65)]
	_ = x[OpFLDZ-(66)]
	_ = x[OpFMULP-(67)]
	_ = x[OpFNCLEX-(68)]
	_ = x[OpFNDISI-(69)]
	_ = x[OpFNENI-(70)]
	_ = x[OpFNINIT-(71)]
	_ = x[OpFNOP-(72)]
	_ = x[OpFNSETPM-(73)]
	_ = x[OpFPATAN-(74)]
	_ = x[OpFPREM-(75)]
	_ = x[OpFPREM1-(76)]
	_ = x[OpFPTAN-(77)]
	_ = x[OpFRNDINT-(78)]
	_ = x[OpFRSTOR-(79)]
	_ = x[OpFS-(80)]
	_ = x[OpFSCALE-(81)]
	_ = x[OpFSETPM-(82)]
	_ = x[OpFSIN-(83)]
	_ = x[OpFSINCOS-(84)]
	_ = x[OpFSQRT-(85)]
	_ = x[OpFSUBP-(86)]
	_ = x[OpFSUBRP-(87)]
	_ = x[OpFTST-(88)]
	_ = x[OpFUCOM-(89)]
	_ = x[OpFUCOMP-(90)]
	_ = x[OpFUCOMPP-(91)]
	_ = x[OpFXAM-(92)]
	_ = x[OpFXCH-(93)]
	_ = x[OpFXRSTOR-(94)]
	_ = x[OpFXTRACT-(95)]
	_ = x[OpFYL2X-(96)]
	_ = x[OpFYL2XP1-(97)]
	_ = x[OpGETSEC-(98)]
	_ = x[OpGS-(99)]
	_ = x[OpHLT-(100)]
	_ = x[OpICEBP-(101)]
	_ = x[OpIDIV-(102)]
	_ = x[OpIMUL-(103)]
	_ = x[OpINT-(104)]
	_ = x[OpCMP-(105)]
	_ = x[OpCALL-(106)]
	_ = x[OpLGDT-(107)]
	_ = x[OpJMP-(108)]
	_ = x[OpJA-(109)]
	_ = x[OpJAE-(110)]
	_ = x[OpJB-(111)]
	_ = x[OpJBE-(112)]
	_ = x[OpJC-(113)]
	_ = x[OpJE-(114)]
	_ = x[OpJG-(115)]
	_ = x[OpJGE-(116)]
	_ = x[OpJL-(117)]
	_ = x[OpJLE-(118)]
	_ = x[OpJNA-(119)]
	_ = x[OpJNAE-(120)]
	_ = x[OpJNB-(121)]
	_ = x[OpJNBE-(122)]
	_ = x[OpJNC-(123)]
	_ = x[OpJNE-(124)]
	_ = x[OpJNG-(125)]
	_ = x[OpJNGE-(126)]
	_ = x[OpJNL-(127)]
	_ = x[OpJNLE-(128)]
	_ = x[OpJNO-(129)]
	_ = x[OpJNP-(130)]
	_ = x[OpJNS-(131)]
	_ = x[OpJNZ-(132)]
	_ = x[OpJO-(133)]
	_ = x[OpJP-(134)]
	_ = x[OpJPE-(135)]
	_ = x[OpJPO-(136)]
	_ = x[OpJS-(137)]
	_ = x[OpJZ-(138)]
	_ = x[OpINTO-(139)]
	_ = x[OpINVD-(140)]
	_ = x[OpIRET-(141)]
	_ = x[OpIRETD-(142)]
	_ = x[OpIRETQ-(143)]
	_ = x[OpJMPE-(144)]
	_ = x[OpLAHF-(145)]
	_ = x[OpLEAVE-(146)]
	_ = x[OpLFENCE-(147)]
	_ = x[OpLOADALL-(148)]
	_ = x[OpLOCK-(149)]
	_ = x[OpMFENCE-(150)]
	_ = x[OpMONITOR-(151)]
	_ = x[OpMUL-(152)]
	_ = x[OpMWAIT-(153)]
	_ = x[OpNOP-(154)]
	_ = x[OpNTAKEN-(155)]
	_ = x[OpPAUSE-(156)]
	_ = x[OpPOPA-(157)]
	_ = x[OpPOPAD-(158)]
	_ = x[OpPOPF-(159)]
	_ = x[OpPOPFD-(160)]
	_ = x[OpPOPFQ-(161)]
	_ = x[OpPUSHA-(162)]
	_ = x[OpPUSHAD-(163)]
	_ = x[OpPUSHF-(164)]
	_ = x[OpPUSHFD-(165)]
	_ = x[OpPUSHFQ-(166)]
	_ = x[OpRDMSR-(167)]
	_ = x[OpRDPMC-(168)]
	_ = x[OpRDTSC-(169)]
	_ = x[OpRDTSCP-(170)]
	_ = x[OpREP-(171)]
	_ = x[OpREPE-(172)]
	_ = x[OpREPNE-(173)]
	_ = x[OpRET-(174)]
	_ = x[OpRETF-(175)]
	_ = x[OpRETN-(176)]
	_ = x[OpRSM-(177)]
	_ = x[OpSAHF-(178)]
	_ = x[OpSETALC-(179)]
	_ = x[OpSFENCE-(180)]
	_ = x[OpSS-(181)]
	_ = x[OpSTC-(182)]
	_ = x[OpSTD-(183)]
	_ = x[OpSTI-(184)]
	_ = x[OpSWAPGS-(185)]
	_ = x[OpSYSCALL-(186)]
	_ = x[OpSYSENTER-(187)]
	_ = x[OpSYSEXIT-(188)]
	_ = x[OpSYSRET-(189)]
	_ = x[OpTAKEN-(190)]
	_ = x[OpUD2-(191)]
	_ = x[OpVMCALL-(192)]
	_ = x[OpVMLAUNCH-(193)]
	_ = x[OpVMRESUME-(194)]
	_ = x[OpVMXOFF-(195)]
	_ = x[OpWAIT-(196)]
	_ = x[OpWBINVD-(197)]
	_ = x[OpWRMSR-(198)]
	_ = x[OpXGETBV-(199)]
	_ = x[OpXRSTOR-(200)]
	_ = x[OpXSETBV-(201)]
}

var _OcodeKindValues = []OcodeKind{OpL, OpOUT, OpDB, OpDW, OpDD, OpRESB, OpMOV, OpADD, OpSUB, OpAND, OpOR, OpXOR, OpNOT, OpSHR, OpSHL, OpSAR, OpIN, OpAAA, OpAAD, OpAAM, OpAAS, OpADX, OpALTER, OpAMX, OpCBW, OpCDQ, OpCDQE, OpCLC, OpCLD, OpCLI, OpCLTS, OpCMC, OpCPUID, OpCQO, OpCS, OpCWD, OpCWDE, OpDAA, OpDAS, OpDIV, OpDS, OpEMMS, OpENTER, OpES, OpF2XM1, OpFABS, OpFADDP, OpFCHS, OpFCLEX, OpFCOM, OpFCOMP, OpFCOMPP, OpFCOS, OpFDECSTP, OpFDISI, OpFDIVP, OpFDIVRP, OpFENI, OpFINCSTP, OpFINIT, OpFLD1, OpFLDL2E, OpFLDL2T, OpFLDLG2, OpFLDLN2, OpFLDPI, OpFLDZ, OpFMULP, OpFNCLEX, OpFNDISI, OpFNENI, OpFNINIT, OpFNOP, OpFNSETPM, OpFPATAN, OpFPREM, OpFPREM1, OpFPTAN, OpFRNDINT, OpFRSTOR, OpFS, OpFSCALE, OpFSETPM, OpFSIN, OpFSINCOS, OpFSQRT, OpFSUBP, OpFSUBRP, OpFTST, OpFUCOM, OpFUCOMP, OpFUCOMPP, OpFXAM, OpFXCH, OpFXRSTOR, OpFXTRACT, OpFYL2X, OpFYL2XP1, OpGETSEC, OpGS, OpHLT, OpICEBP, OpIDIV, OpIMUL, OpINT, OpCMP, OpCALL, OpLGDT, OpJMP, OpJA, OpJAE, OpJB, OpJBE, OpJC, OpJE, OpJG, OpJGE, OpJL, OpJLE, OpJNA, OpJNAE, OpJNB, OpJNBE, OpJNC, OpJNE, OpJNG, OpJNGE, OpJNL, OpJNLE, OpJNO, OpJNP, OpJNS, OpJNZ, OpJO, OpJP, OpJPE, OpJPO, OpJS, OpJZ, OpINTO, OpINVD, OpIRET, OpIRETD, OpIRETQ, OpJMPE, OpLAHF, OpLEAVE, OpLFENCE, OpLOADALL, OpLOCK, OpMFENCE, OpMONITOR, OpMUL, OpMWAIT, OpNOP, OpNTAKEN, OpPAUSE, OpPOPA, OpPOPAD, OpPOPF, OpPOPFD, OpPOPFQ, OpPUSHA, OpPUSHAD, OpPUSHF, OpPUSHFD, OpPUSHFQ, OpRDMSR, OpRDPMC, OpRDTSC, OpRDTSCP, OpREP, OpREPE, OpREPNE, OpRET, OpRETF, OpRETN, OpRSM, OpSAHF, OpSETALC, OpSFENCE, OpSS, OpSTC, OpSTD, OpSTI, OpSWAPGS, OpSYSCALL, OpSYSENTER, OpSYSEXIT, OpSYSRET, OpTAKEN, OpUD2, OpVMCALL, OpVMLAUNCH, OpVMRESUME, OpVMXOFF, OpWAIT, OpWBINVD, OpWRMSR, OpXGETBV, OpXRSTOR, OpXSETBV}

var _OcodeKindNameToValueMap = map[string]OcodeKind{
	_OcodeKindName[0:3]:            OpL,
	_OcodeKindLowerName[0:3]:       OpL,
	_OcodeKindName[3:8]:            OpOUT,
	_OcodeKindLowerName[3:8]:       OpOUT,
	_OcodeKindName[8:12]:           OpDB,
	_OcodeKindLowerName[8:12]:      OpDB,
	_OcodeKindName[12:16]:          OpDW,
	_OcodeKindLowerName[12:16]:     OpDW,
	_OcodeKindName[16:20]:          OpDD,
	_OcodeKindLowerName[16:20]:     OpDD,
	_OcodeKindName[20:26]:          OpRESB,
	_OcodeKindLowerName[20:26]:     OpRESB,
	_OcodeKindName[26:31]:          OpMOV,
	_OcodeKindLowerName[26:31]:     OpMOV,
	_OcodeKindName[31:36]:          OpADD,
	_OcodeKindLowerName[31:36]:     OpADD,
	_OcodeKindName[36:41]:          OpSUB,
	_OcodeKindLowerName[36:41]:     OpSUB,
	_OcodeKindName[41:46]:          OpAND,
	_OcodeKindLowerName[41:46]:     OpAND,
	_OcodeKindName[46:50]:          OpOR,
	_OcodeKindLowerName[46:50]:     OpOR,
	_OcodeKindName[50:55]:          OpXOR,
	_OcodeKindLowerName[50:55]:     OpXOR,
	_OcodeKindName[55:60]:          OpNOT,
	_OcodeKindLowerName[55:60]:     OpNOT,
	_OcodeKindName[60:65]:          OpSHR,
	_OcodeKindLowerName[60:65]:     OpSHR,
	_OcodeKindName[65:70]:          OpSHL,
	_OcodeKindLowerName[65:70]:     OpSHL,
	_OcodeKindName[70:75]:          OpSAR,
	_OcodeKindLowerName[70:75]:     OpSAR,
	_OcodeKindName[75:79]:          OpIN,
	_OcodeKindLowerName[75:79]:     OpIN,
	_OcodeKindName[79:84]:          OpAAA,
	_OcodeKindLowerName[79:84]:     OpAAA,
	_OcodeKindName[84:89]:          OpAAD,
	_OcodeKindLowerName[84:89]:     OpAAD,
	_OcodeKindName[89:94]:          OpAAM,
	_OcodeKindLowerName[89:94]:     OpAAM,
	_OcodeKindName[94:99]:          OpAAS,
	_OcodeKindLowerName[94:99]:     OpAAS,
	_OcodeKindName[99:104]:         OpADX,
	_OcodeKindLowerName[99:104]:    OpADX,
	_OcodeKindName[104:111]:        OpALTER,
	_OcodeKindLowerName[104:111]:   OpALTER,
	_OcodeKindName[111:116]:        OpAMX,
	_OcodeKindLowerName[111:116]:   OpAMX,
	_OcodeKindName[116:121]:        OpCBW,
	_OcodeKindLowerName[116:121]:   OpCBW,
	_OcodeKindName[121:126]:        OpCDQ,
	_OcodeKindLowerName[121:126]:   OpCDQ,
	_OcodeKindName[126:132]:        OpCDQE,
	_OcodeKindLowerName[126:132]:   OpCDQE,
	_OcodeKindName[132:137]:        OpCLC,
	_OcodeKindLowerName[132:137]:   OpCLC,
	_OcodeKindName[137:142]:        OpCLD,
	_OcodeKindLowerName[137:142]:   OpCLD,
	_OcodeKindName[142:147]:        OpCLI,
	_OcodeKindLowerName[142:147]:   OpCLI,
	_OcodeKindName[147:153]:        OpCLTS,
	_OcodeKindLowerName[147:153]:   OpCLTS,
	_OcodeKindName[153:158]:        OpCMC,
	_OcodeKindLowerName[153:158]:   OpCMC,
	_OcodeKindName[158:165]:        OpCPUID,
	_OcodeKindLowerName[158:165]:   OpCPUID,
	_OcodeKindName[165:170]:        OpCQO,
	_OcodeKindLowerName[165:170]:   OpCQO,
	_OcodeKindName[170:174]:        OpCS,
	_OcodeKindLowerName[170:174]:   OpCS,
	_OcodeKindName[174:179]:        OpCWD,
	_OcodeKindLowerName[174:179]:   OpCWD,
	_OcodeKindName[179:185]:        OpCWDE,
	_OcodeKindLowerName[179:185]:   OpCWDE,
	_OcodeKindName[185:190]:        OpDAA,
	_OcodeKindLowerName[185:190]:   OpDAA,
	_OcodeKindName[190:195]:        OpDAS,
	_OcodeKindLowerName[190:195]:   OpDAS,
	_OcodeKindName[195:200]:        OpDIV,
	_OcodeKindLowerName[195:200]:   OpDIV,
	_OcodeKindName[200:204]:        OpDS,
	_OcodeKindLowerName[200:204]:   OpDS,
	_OcodeKindName[204:210]:        OpEMMS,
	_OcodeKindLowerName[204:210]:   OpEMMS,
	_OcodeKindName[210:217]:        OpENTER,
	_OcodeKindLowerName[210:217]:   OpENTER,
	_OcodeKindName[217:221]:        OpES,
	_OcodeKindLowerName[217:221]:   OpES,
	_OcodeKindName[221:228]:        OpF2XM1,
	_OcodeKindLowerName[221:228]:   OpF2XM1,
	_OcodeKindName[228:234]:        OpFABS,
	_OcodeKindLowerName[228:234]:   OpFABS,
	_OcodeKindName[234:241]:        OpFADDP,
	_OcodeKindLowerName[234:241]:   OpFADDP,
	_OcodeKindName[241:247]:        OpFCHS,
	_OcodeKindLowerName[241:247]:   OpFCHS,
	_OcodeKindName[247:254]:        OpFCLEX,
	_OcodeKindLowerName[247:254]:   OpFCLEX,
	_OcodeKindName[254:260]:        OpFCOM,
	_OcodeKindLowerName[254:260]:   OpFCOM,
	_OcodeKindName[260:267]:        OpFCOMP,
	_OcodeKindLowerName[260:267]:   OpFCOMP,
	_OcodeKindName[267:275]:        OpFCOMPP,
	_OcodeKindLowerName[267:275]:   OpFCOMPP,
	_OcodeKindName[275:281]:        OpFCOS,
	_OcodeKindLowerName[275:281]:   OpFCOS,
	_OcodeKindName[281:290]:        OpFDECSTP,
	_OcodeKindLowerName[281:290]:   OpFDECSTP,
	_OcodeKindName[290:297]:        OpFDISI,
	_OcodeKindLowerName[290:297]:   OpFDISI,
	_OcodeKindName[297:304]:        OpFDIVP,
	_OcodeKindLowerName[297:304]:   OpFDIVP,
	_OcodeKindName[304:312]:        OpFDIVRP,
	_OcodeKindLowerName[304:312]:   OpFDIVRP,
	_OcodeKindName[312:318]:        OpFENI,
	_OcodeKindLowerName[312:318]:   OpFENI,
	_OcodeKindName[318:327]:        OpFINCSTP,
	_OcodeKindLowerName[318:327]:   OpFINCSTP,
	_OcodeKindName[327:334]:        OpFINIT,
	_OcodeKindLowerName[327:334]:   OpFINIT,
	_OcodeKindName[334:340]:        OpFLD1,
	_OcodeKindLowerName[334:340]:   OpFLD1,
	_OcodeKindName[340:348]:        OpFLDL2E,
	_OcodeKindLowerName[340:348]:   OpFLDL2E,
	_OcodeKindName[348:356]:        OpFLDL2T,
	_OcodeKindLowerName[348:356]:   OpFLDL2T,
	_OcodeKindName[356:364]:        OpFLDLG2,
	_OcodeKindLowerName[356:364]:   OpFLDLG2,
	_OcodeKindName[364:372]:        OpFLDLN2,
	_OcodeKindLowerName[364:372]:   OpFLDLN2,
	_OcodeKindName[372:379]:        OpFLDPI,
	_OcodeKindLowerName[372:379]:   OpFLDPI,
	_OcodeKindName[379:385]:        OpFLDZ,
	_OcodeKindLowerName[379:385]:   OpFLDZ,
	_OcodeKindName[385:392]:        OpFMULP,
	_OcodeKindLowerName[385:392]:   OpFMULP,
	_OcodeKindName[392:400]:        OpFNCLEX,
	_OcodeKindLowerName[392:400]:   OpFNCLEX,
	_OcodeKindName[400:408]:        OpFNDISI,
	_OcodeKindLowerName[400:408]:   OpFNDISI,
	_OcodeKindName[408:415]:        OpFNENI,
	_OcodeKindLowerName[408:415]:   OpFNENI,
	_OcodeKindName[415:423]:        OpFNINIT,
	_OcodeKindLowerName[415:423]:   OpFNINIT,
	_OcodeKindName[423:429]:        OpFNOP,
	_OcodeKindLowerName[423:429]:   OpFNOP,
	_OcodeKindName[429:438]:        OpFNSETPM,
	_OcodeKindLowerName[429:438]:   OpFNSETPM,
	_OcodeKindName[438:446]:        OpFPATAN,
	_OcodeKindLowerName[438:446]:   OpFPATAN,
	_OcodeKindName[446:453]:        OpFPREM,
	_OcodeKindLowerName[446:453]:   OpFPREM,
	_OcodeKindName[453:461]:        OpFPREM1,
	_OcodeKindLowerName[453:461]:   OpFPREM1,
	_OcodeKindName[461:468]:        OpFPTAN,
	_OcodeKindLowerName[461:468]:   OpFPTAN,
	_OcodeKindName[468:477]:        OpFRNDINT,
	_OcodeKindLowerName[468:477]:   OpFRNDINT,
	_OcodeKindName[477:485]:        OpFRSTOR,
	_OcodeKindLowerName[477:485]:   OpFRSTOR,
	_OcodeKindName[485:489]:        OpFS,
	_OcodeKindLowerName[485:489]:   OpFS,
	_OcodeKindName[489:497]:        OpFSCALE,
	_OcodeKindLowerName[489:497]:   OpFSCALE,
	_OcodeKindName[497:505]:        OpFSETPM,
	_OcodeKindLowerName[497:505]:   OpFSETPM,
	_OcodeKindName[505:511]:        OpFSIN,
	_OcodeKindLowerName[505:511]:   OpFSIN,
	_OcodeKindName[511:520]:        OpFSINCOS,
	_OcodeKindLowerName[511:520]:   OpFSINCOS,
	_OcodeKindName[520:527]:        OpFSQRT,
	_OcodeKindLowerName[520:527]:   OpFSQRT,
	_OcodeKindName[527:534]:        OpFSUBP,
	_OcodeKindLowerName[527:534]:   OpFSUBP,
	_OcodeKindName[534:542]:        OpFSUBRP,
	_OcodeKindLowerName[534:542]:   OpFSUBRP,
	_OcodeKindName[542:548]:        OpFTST,
	_OcodeKindLowerName[542:548]:   OpFTST,
	_OcodeKindName[548:555]:        OpFUCOM,
	_OcodeKindLowerName[548:555]:   OpFUCOM,
	_OcodeKindName[555:563]:        OpFUCOMP,
	_OcodeKindLowerName[555:563]:   OpFUCOMP,
	_OcodeKindName[563:572]:        OpFUCOMPP,
	_OcodeKindLowerName[563:572]:   OpFUCOMPP,
	_OcodeKindName[572:578]:        OpFXAM,
	_OcodeKindLowerName[572:578]:   OpFXAM,
	_OcodeKindName[578:584]:        OpFXCH,
	_OcodeKindLowerName[578:584]:   OpFXCH,
	_OcodeKindName[584:593]:        OpFXRSTOR,
	_OcodeKindLowerName[584:593]:   OpFXRSTOR,
	_OcodeKindName[593:602]:        OpFXTRACT,
	_OcodeKindLowerName[593:602]:   OpFXTRACT,
	_OcodeKindName[602:609]:        OpFYL2X,
	_OcodeKindLowerName[602:609]:   OpFYL2X,
	_OcodeKindName[609:618]:        OpFYL2XP1,
	_OcodeKindLowerName[609:618]:   OpFYL2XP1,
	_OcodeKindName[618:626]:        OpGETSEC,
	_OcodeKindLowerName[618:626]:   OpGETSEC,
	_OcodeKindName[626:630]:        OpGS,
	_OcodeKindLowerName[626:630]:   OpGS,
	_OcodeKindName[630:635]:        OpHLT,
	_OcodeKindLowerName[630:635]:   OpHLT,
	_OcodeKindName[635:642]:        OpICEBP,
	_OcodeKindLowerName[635:642]:   OpICEBP,
	_OcodeKindName[642:648]:        OpIDIV,
	_OcodeKindLowerName[642:648]:   OpIDIV,
	_OcodeKindName[648:654]:        OpIMUL,
	_OcodeKindLowerName[648:654]:   OpIMUL,
	_OcodeKindName[654:659]:        OpINT,
	_OcodeKindLowerName[654:659]:   OpINT,
	_OcodeKindName[659:664]:        OpCMP,
	_OcodeKindLowerName[659:664]:   OpCMP,
	_OcodeKindName[664:670]:        OpCALL,
	_OcodeKindLowerName[664:670]:   OpCALL,
	_OcodeKindName[670:676]:        OpLGDT,
	_OcodeKindLowerName[670:676]:   OpLGDT,
	_OcodeKindName[676:681]:        OpJMP,
	_OcodeKindLowerName[676:681]:   OpJMP,
	_OcodeKindName[681:685]:        OpJA,
	_OcodeKindLowerName[681:685]:   OpJA,
	_OcodeKindName[685:690]:        OpJAE,
	_OcodeKindLowerName[685:690]:   OpJAE,
	_OcodeKindName[690:694]:        OpJB,
	_OcodeKindLowerName[690:694]:   OpJB,
	_OcodeKindName[694:699]:        OpJBE,
	_OcodeKindLowerName[694:699]:   OpJBE,
	_OcodeKindName[699:703]:        OpJC,
	_OcodeKindLowerName[699:703]:   OpJC,
	_OcodeKindName[703:707]:        OpJE,
	_OcodeKindLowerName[703:707]:   OpJE,
	_OcodeKindName[707:711]:        OpJG,
	_OcodeKindLowerName[707:711]:   OpJG,
	_OcodeKindName[711:716]:        OpJGE,
	_OcodeKindLowerName[711:716]:   OpJGE,
	_OcodeKindName[716:720]:        OpJL,
	_OcodeKindLowerName[716:720]:   OpJL,
	_OcodeKindName[720:725]:        OpJLE,
	_OcodeKindLowerName[720:725]:   OpJLE,
	_OcodeKindName[725:730]:        OpJNA,
	_OcodeKindLowerName[725:730]:   OpJNA,
	_OcodeKindName[730:736]:        OpJNAE,
	_OcodeKindLowerName[730:736]:   OpJNAE,
	_OcodeKindName[736:741]:        OpJNB,
	_OcodeKindLowerName[736:741]:   OpJNB,
	_OcodeKindName[741:747]:        OpJNBE,
	_OcodeKindLowerName[741:747]:   OpJNBE,
	_OcodeKindName[747:752]:        OpJNC,
	_OcodeKindLowerName[747:752]:   OpJNC,
	_OcodeKindName[752:757]:        OpJNE,
	_OcodeKindLowerName[752:757]:   OpJNE,
	_OcodeKindName[757:762]:        OpJNG,
	_OcodeKindLowerName[757:762]:   OpJNG,
	_OcodeKindName[762:768]:        OpJNGE,
	_OcodeKindLowerName[762:768]:   OpJNGE,
	_OcodeKindName[768:773]:        OpJNL,
	_OcodeKindLowerName[768:773]:   OpJNL,
	_OcodeKindName[773:779]:        OpJNLE,
	_OcodeKindLowerName[773:779]:   OpJNLE,
	_OcodeKindName[779:784]:        OpJNO,
	_OcodeKindLowerName[779:784]:   OpJNO,
	_OcodeKindName[784:789]:        OpJNP,
	_OcodeKindLowerName[784:789]:   OpJNP,
	_OcodeKindName[789:794]:        OpJNS,
	_OcodeKindLowerName[789:794]:   OpJNS,
	_OcodeKindName[794:799]:        OpJNZ,
	_OcodeKindLowerName[794:799]:   OpJNZ,
	_OcodeKindName[799:803]:        OpJO,
	_OcodeKindLowerName[799:803]:   OpJO,
	_OcodeKindName[803:807]:        OpJP,
	_OcodeKindLowerName[803:807]:   OpJP,
	_OcodeKindName[807:812]:        OpJPE,
	_OcodeKindLowerName[807:812]:   OpJPE,
	_OcodeKindName[812:817]:        OpJPO,
	_OcodeKindLowerName[812:817]:   OpJPO,
	_OcodeKindName[817:821]:        OpJS,
	_OcodeKindLowerName[817:821]:   OpJS,
	_OcodeKindName[821:825]:        OpJZ,
	_OcodeKindLowerName[821:825]:   OpJZ,
	_OcodeKindName[825:831]:        OpINTO,
	_OcodeKindLowerName[825:831]:   OpINTO,
	_OcodeKindName[831:837]:        OpINVD,
	_OcodeKindLowerName[831:837]:   OpINVD,
	_OcodeKindName[837:843]:        OpIRET,
	_OcodeKindLowerName[837:843]:   OpIRET,
	_OcodeKindName[843:850]:        OpIRETD,
	_OcodeKindLowerName[843:850]:   OpIRETD,
	_OcodeKindName[850:857]:        OpIRETQ,
	_OcodeKindLowerName[850:857]:   OpIRETQ,
	_OcodeKindName[857:863]:        OpJMPE,
	_OcodeKindLowerName[857:863]:   OpJMPE,
	_OcodeKindName[863:869]:        OpLAHF,
	_OcodeKindLowerName[863:869]:   OpLAHF,
	_OcodeKindName[869:876]:        OpLEAVE,
	_OcodeKindLowerName[869:876]:   OpLEAVE,
	_OcodeKindName[876:884]:        OpLFENCE,
	_OcodeKindLowerName[876:884]:   OpLFENCE,
	_OcodeKindName[884:893]:        OpLOADALL,
	_OcodeKindLowerName[884:893]:   OpLOADALL,
	_OcodeKindName[893:899]:        OpLOCK,
	_OcodeKindLowerName[893:899]:   OpLOCK,
	_OcodeKindName[899:907]:        OpMFENCE,
	_OcodeKindLowerName[899:907]:   OpMFENCE,
	_OcodeKindName[907:916]:        OpMONITOR,
	_OcodeKindLowerName[907:916]:   OpMONITOR,
	_OcodeKindName[916:921]:        OpMUL,
	_OcodeKindLowerName[916:921]:   OpMUL,
	_OcodeKindName[921:928]:        OpMWAIT,
	_OcodeKindLowerName[921:928]:   OpMWAIT,
	_OcodeKindName[928:933]:        OpNOP,
	_OcodeKindLowerName[928:933]:   OpNOP,
	_OcodeKindName[933:941]:        OpNTAKEN,
	_OcodeKindLowerName[933:941]:   OpNTAKEN,
	_OcodeKindName[941:948]:        OpPAUSE,
	_OcodeKindLowerName[941:948]:   OpPAUSE,
	_OcodeKindName[948:954]:        OpPOPA,
	_OcodeKindLowerName[948:954]:   OpPOPA,
	_OcodeKindName[954:961]:        OpPOPAD,
	_OcodeKindLowerName[954:961]:   OpPOPAD,
	_OcodeKindName[961:967]:        OpPOPF,
	_OcodeKindLowerName[961:967]:   OpPOPF,
	_OcodeKindName[967:974]:        OpPOPFD,
	_OcodeKindLowerName[967:974]:   OpPOPFD,
	_OcodeKindName[974:981]:        OpPOPFQ,
	_OcodeKindLowerName[974:981]:   OpPOPFQ,
	_OcodeKindName[981:988]:        OpPUSHA,
	_OcodeKindLowerName[981:988]:   OpPUSHA,
	_OcodeKindName[988:996]:        OpPUSHAD,
	_OcodeKindLowerName[988:996]:   OpPUSHAD,
	_OcodeKindName[996:1003]:       OpPUSHF,
	_OcodeKindLowerName[996:1003]:  OpPUSHF,
	_OcodeKindName[1003:1011]:      OpPUSHFD,
	_OcodeKindLowerName[1003:1011]: OpPUSHFD,
	_OcodeKindName[1011:1019]:      OpPUSHFQ,
	_OcodeKindLowerName[1011:1019]: OpPUSHFQ,
	_OcodeKindName[1019:1026]:      OpRDMSR,
	_OcodeKindLowerName[1019:1026]: OpRDMSR,
	_OcodeKindName[1026:1033]:      OpRDPMC,
	_OcodeKindLowerName[1026:1033]: OpRDPMC,
	_OcodeKindName[1033:1040]:      OpRDTSC,
	_OcodeKindLowerName[1033:1040]: OpRDTSC,
	_OcodeKindName[1040:1048]:      OpRDTSCP,
	_OcodeKindLowerName[1040:1048]: OpRDTSCP,
	_OcodeKindName[1048:1053]:      OpREP,
	_OcodeKindLowerName[1048:1053]: OpREP,
	_OcodeKindName[1053:1059]:      OpREPE,
	_OcodeKindLowerName[1053:1059]: OpREPE,
	_OcodeKindName[1059:1066]:      OpREPNE,
	_OcodeKindLowerName[1059:1066]: OpREPNE,
	_OcodeKindName[1066:1071]:      OpRET,
	_OcodeKindLowerName[1066:1071]: OpRET,
	_OcodeKindName[1071:1077]:      OpRETF,
	_OcodeKindLowerName[1071:1077]: OpRETF,
	_OcodeKindName[1077:1083]:      OpRETN,
	_OcodeKindLowerName[1077:1083]: OpRETN,
	_OcodeKindName[1083:1088]:      OpRSM,
	_OcodeKindLowerName[1083:1088]: OpRSM,
	_OcodeKindName[1088:1094]:      OpSAHF,
	_OcodeKindLowerName[1088:1094]: OpSAHF,
	_OcodeKindName[1094:1102]:      OpSETALC,
	_OcodeKindLowerName[1094:1102]: OpSETALC,
	_OcodeKindName[1102:1110]:      OpSFENCE,
	_OcodeKindLowerName[1102:1110]: OpSFENCE,
	_OcodeKindName[1110:1114]:      OpSS,
	_OcodeKindLowerName[1110:1114]: OpSS,
	_OcodeKindName[1114:1119]:      OpSTC,
	_OcodeKindLowerName[1114:1119]: OpSTC,
	_OcodeKindName[1119:1124]:      OpSTD,
	_OcodeKindLowerName[1119:1124]: OpSTD,
	_OcodeKindName[1124:1129]:      OpSTI,
	_OcodeKindLowerName[1124:1129]: OpSTI,
	_OcodeKindName[1129:1137]:      OpSWAPGS,
	_OcodeKindLowerName[1129:1137]: OpSWAPGS,
	_OcodeKindName[1137:1146]:      OpSYSCALL,
	_OcodeKindLowerName[1137:1146]: OpSYSCALL,
	_OcodeKindName[1146:1156]:      OpSYSENTER,
	_OcodeKindLowerName[1146:1156]: OpSYSENTER,
	_OcodeKindName[1156:1165]:      OpSYSEXIT,
	_OcodeKindLowerName[1156:1165]: OpSYSEXIT,
	_OcodeKindName[1165:1173]:      OpSYSRET,
	_OcodeKindLowerName[1165:1173]: OpSYSRET,
	_OcodeKindName[1173:1180]:      OpTAKEN,
	_OcodeKindLowerName[1173:1180]: OpTAKEN,
	_OcodeKindName[1180:1185]:      OpUD2,
	_OcodeKindLowerName[1180:1185]: OpUD2,
	_OcodeKindName[1185:1193]:      OpVMCALL,
	_OcodeKindLowerName[1185:1193]: OpVMCALL,
	_OcodeKindName[1193:1203]:      OpVMLAUNCH,
	_OcodeKindLowerName[1193:1203]: OpVMLAUNCH,
	_OcodeKindName[1203:1213]:      OpVMRESUME,
	_OcodeKindLowerName[1203:1213]: OpVMRESUME,
	_OcodeKindName[1213:1221]:      OpVMXOFF,
	_OcodeKindLowerName[1213:1221]: OpVMXOFF,
	_OcodeKindName[1221:1227]:      OpWAIT,
	_OcodeKindLowerName[1221:1227]: OpWAIT,
	_OcodeKindName[1227:1235]:      OpWBINVD,
	_OcodeKindLowerName[1227:1235]: OpWBINVD,
	_OcodeKindName[1235:1242]:      OpWRMSR,
	_OcodeKindLowerName[1235:1242]: OpWRMSR,
	_OcodeKindName[1242:1250]:      OpXGETBV,
	_OcodeKindLowerName[1242:1250]: OpXGETBV,
	_OcodeKindName[1250:1258]:      OpXRSTOR,
	_OcodeKindLowerName[1250:1258]: OpXRSTOR,
	_OcodeKindName[1258:1266]:      OpXSETBV,
	_OcodeKindLowerName[1258:1266]: OpXSETBV,
}

var _OcodeKindNames = []string{
	_OcodeKindName[0:3],
	_OcodeKindName[3:8],
	_OcodeKindName[8:12],
	_OcodeKindName[12:16],
	_OcodeKindName[16:20],
	_OcodeKindName[20:26],
	_OcodeKindName[26:31],
	_OcodeKindName[31:36],
	_OcodeKindName[36:41],
	_OcodeKindName[41:46],
	_OcodeKindName[46:50],
	_OcodeKindName[50:55],
	_OcodeKindName[55:60],
	_OcodeKindName[60:65],
	_OcodeKindName[65:70],
	_OcodeKindName[70:75],
	_OcodeKindName[75:79],
	_OcodeKindName[79:84],
	_OcodeKindName[84:89],
	_OcodeKindName[89:94],
	_OcodeKindName[94:99],
	_OcodeKindName[99:104],
	_OcodeKindName[104:111],
	_OcodeKindName[111:116],
	_OcodeKindName[116:121],
	_OcodeKindName[121:126],
	_OcodeKindName[126:132],
	_OcodeKindName[132:137],
	_OcodeKindName[137:142],
	_OcodeKindName[142:147],
	_OcodeKindName[147:153],
	_OcodeKindName[153:158],
	_OcodeKindName[158:165],
	_OcodeKindName[165:170],
	_OcodeKindName[170:174],
	_OcodeKindName[174:179],
	_OcodeKindName[179:185],
	_OcodeKindName[185:190],
	_OcodeKindName[190:195],
	_OcodeKindName[195:200],
	_OcodeKindName[200:204],
	_OcodeKindName[204:210],
	_OcodeKindName[210:217],
	_OcodeKindName[217:221],
	_OcodeKindName[221:228],
	_OcodeKindName[228:234],
	_OcodeKindName[234:241],
	_OcodeKindName[241:247],
	_OcodeKindName[247:254],
	_OcodeKindName[254:260],
	_OcodeKindName[260:267],
	_OcodeKindName[267:275],
	_OcodeKindName[275:281],
	_OcodeKindName[281:290],
	_OcodeKindName[290:297],
	_OcodeKindName[297:304],
	_OcodeKindName[304:312],
	_OcodeKindName[312:318],
	_OcodeKindName[318:327],
	_OcodeKindName[327:334],
	_OcodeKindName[334:340],
	_OcodeKindName[340:348],
	_OcodeKindName[348:356],
	_OcodeKindName[356:364],
	_OcodeKindName[364:372],
	_OcodeKindName[372:379],
	_OcodeKindName[379:385],
	_OcodeKindName[385:392],
	_OcodeKindName[392:400],
	_OcodeKindName[400:408],
	_OcodeKindName[408:415],
	_OcodeKindName[415:423],
	_OcodeKindName[423:429],
	_OcodeKindName[429:438],
	_OcodeKindName[438:446],
	_OcodeKindName[446:453],
	_OcodeKindName[453:461],
	_OcodeKindName[461:468],
	_OcodeKindName[468:477],
	_OcodeKindName[477:485],
	_OcodeKindName[485:489],
	_OcodeKindName[489:497],
	_OcodeKindName[497:505],
	_OcodeKindName[505:511],
	_OcodeKindName[511:520],
	_OcodeKindName[520:527],
	_OcodeKindName[527:534],
	_OcodeKindName[534:542],
	_OcodeKindName[542:548],
	_OcodeKindName[548:555],
	_OcodeKindName[555:563],
	_OcodeKindName[563:572],
	_OcodeKindName[572:578],
	_OcodeKindName[578:584],
	_OcodeKindName[584:593],
	_OcodeKindName[593:602],
	_OcodeKindName[602:609],
	_OcodeKindName[609:618],
	_OcodeKindName[618:626],
	_OcodeKindName[626:630],
	_OcodeKindName[630:635],
	_OcodeKindName[635:642],
	_OcodeKindName[642:648],
	_OcodeKindName[648:654],
	_OcodeKindName[654:659],
	_OcodeKindName[659:664],
	_OcodeKindName[664:670],
	_OcodeKindName[670:676],
	_OcodeKindName[676:681],
	_OcodeKindName[681:685],
	_OcodeKindName[685:690],
	_OcodeKindName[690:694],
	_OcodeKindName[694:699],
	_OcodeKindName[699:703],
	_OcodeKindName[703:707],
	_OcodeKindName[707:711],
	_OcodeKindName[711:716],
	_OcodeKindName[716:720],
	_OcodeKindName[720:725],
	_OcodeKindName[725:730],
	_OcodeKindName[730:736],
	_OcodeKindName[736:741],
	_OcodeKindName[741:747],
	_OcodeKindName[747:752],
	_OcodeKindName[752:757],
	_OcodeKindName[757:762],
	_OcodeKindName[762:768],
	_OcodeKindName[768:773],
	_OcodeKindName[773:779],
	_OcodeKindName[779:784],
	_OcodeKindName[784:789],
	_OcodeKindName[789:794],
	_OcodeKindName[794:799],
	_OcodeKindName[799:803],
	_OcodeKindName[803:807],
	_OcodeKindName[807:812],
	_OcodeKindName[812:817],
	_OcodeKindName[817:821],
	_OcodeKindName[821:825],
	_OcodeKindName[825:831],
	_OcodeKindName[831:837],
	_OcodeKindName[837:843],
	_OcodeKindName[843:850],
	_OcodeKindName[850:857],
	_OcodeKindName[857:863],
	_OcodeKindName[863:869],
	_OcodeKindName[869:876],
	_OcodeKindName[876:884],
	_OcodeKindName[884:893],
	_OcodeKindName[893:899],
	_OcodeKindName[899:907],
	_OcodeKindName[907:916],
	_OcodeKindName[916:921],
	_OcodeKindName[921:928],
	_OcodeKindName[928:933],
	_OcodeKindName[933:941],
	_OcodeKindName[941:948],
	_OcodeKindName[948:954],
	_OcodeKindName[954:961],
	_OcodeKindName[961:967],
	_OcodeKindName[967:974],
	_OcodeKindName[974:981],
	_OcodeKindName[981:988],
	_OcodeKindName[988:996],
	_OcodeKindName[996:1003],
	_OcodeKindName[1003:1011],
	_OcodeKindName[1011:1019],
	_OcodeKindName[1019:1026],
	_OcodeKindName[1026:1033],
	_OcodeKindName[1033:1040],
	_OcodeKindName[1040:1048],
	_OcodeKindName[1048:1053],
	_OcodeKindName[1053:1059],
	_OcodeKindName[1059:1066],
	_OcodeKindName[1066:1071],
	_OcodeKindName[1071:1077],
	_OcodeKindName[1077:1083],
	_OcodeKindName[1083:1088],
	_OcodeKindName[1088:1094],
	_OcodeKindName[1094:1102],
	_OcodeKindName[1102:1110],
	_OcodeKindName[1110:1114],
	_OcodeKindName[1114:1119],
	_OcodeKindName[1119:1124],
	_OcodeKindName[1124:1129],
	_OcodeKindName[1129:1137],
	_OcodeKindName[1137:1146],
	_OcodeKindName[1146:1156],
	_OcodeKindName[1156:1165],
	_OcodeKindName[1165:1173],
	_OcodeKindName[1173:1180],
	_OcodeKindName[1180:1185],
	_OcodeKindName[1185:1193],
	_OcodeKindName[1193:1203],
	_OcodeKindName[1203:1213],
	_OcodeKindName[1213:1221],
	_OcodeKindName[1221:1227],
	_OcodeKindName[1227:1235],
	_OcodeKindName[1235:1242],
	_OcodeKindName[1242:1250],
	_OcodeKindName[1250:1258],
	_OcodeKindName[1258:1266],
}

// OcodeKindString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OcodeKindString(s string) (OcodeKind, error) {
	if val, ok := _OcodeKindNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _OcodeKindNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to OcodeKind values", s)
}

// OcodeKindValues returns all values of the enum
func OcodeKindValues() []OcodeKind {
	return _OcodeKindValues
}

// OcodeKindStrings returns a slice of all String values of the enum
func OcodeKindStrings() []string {
	strs := make([]string, len(_OcodeKindNames))
	copy(strs, _OcodeKindNames)
	return strs
}

// IsAOcodeKind returns "true" if the value is listed in the enum definition. "false" otherwise
func (i OcodeKind) IsAOcodeKind() bool {
	for _, v := range _OcodeKindValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for OcodeKind
func (i OcodeKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for OcodeKind
func (i *OcodeKind) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("OcodeKind should be a string, got %s", data)
	}

	var err error
	*i, err = OcodeKindString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for OcodeKind
func (i OcodeKind) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for OcodeKind
func (i *OcodeKind) UnmarshalText(text []byte) error {
	var err error
	*i, err = OcodeKindString(string(text))
	return err
}
