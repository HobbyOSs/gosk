// Code generated by "enumer -type=OcodeKind -json -text"; DO NOT EDIT.

package ocode

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _OcodeKindName = "OpLOpOUTOpDBOpDWOpDDOpRESBOpMOVOpADDOpSUBOpANDOpOROpXOROpNOTOpSHROpSHLOpSAROpINOpAAAOpAADOpAAMOpAASOpADXOpALTEROpAMXOpCBWOpCDQOpCDQEOpCLCOpCLDOpCLIOpCLTSOpCMCOpCPUIDOpCQOOpCSOpCWDOpCWDEOpDAAOpDASOpDIVOpDSOpEMMSOpENTEROpESOpF2XM1OpFABSOpFADDPOpFCHSOpFCLEXOpFCOMOpFCOMPOpFCOMPPOpFCOSOpFDECSTPOpFDISIOpFDIVPOpFDIVRPOpFENIOpFINCSTPOpFINITOpFLD1OpFLDL2EOpFLDL2TOpFLDLG2OpFLDLN2OpFLDPIOpFLDZOpFMULPOpFNCLEXOpFNDISIOpFNENIOpFNINITOpFNOPOpFNSETPMOpFPATANOpFPREMOpFPREM1OpFPTANOpFRNDINTOpFRSTOROpFSOpFSCALEOpFSETPMOpFSINOpFSINCOSOpFSQRTOpFSUBPOpFSUBRPOpFTSTOpFUCOMOpFUCOMPOpFUCOMPPOpFXAMOpFXCHOpFXRSTOROpFXTRACTOpFYL2XOpFYL2XP1OpGETSECOpGSOpHLTOpICEBPOpIDIVOpIMULOpINTOpCMPOpCALLOpLGDTOpJMPOpJMP_FAROpJAOpJAEOpJBOpJBEOpJCOpJEOpJGOpJGEOpJLOpJLEOpJNAOpJNAEOpJNBOpJNBEOpJNCOpJNEOpJNGOpJNGEOpJNLOpJNLEOpJNOOpJNPOpJNSOpJNZOpJOOpJPOpJPEOpJPOOpJSOpJZOpINTOOpINVDOpIRETOpIRETDOpIRETQOpJMPEOpLAHFOpLEAVEOpLFENCEOpLOADALLOpLOCKOpMFENCEOpMONITOROpMULOpMWAITOpNOPOpNTAKENOpPAUSEOpPOPAOpPOPADOpPOPFOpPOPFDOpPOPFQOpPUSHAOpPUSHADOpPUSHFOpPUSHFDOpPUSHFQOpRDMSROpRDPMCOpRDTSCOpRDTSCPOpREPOpREPEOpREPNEOpRETOpRETFOpRETNOpRSMOpSAHFOpSETALCOpSFENCEOpSSOpSTCOpSTDOpSTIOpSWAPGSOpSYSCALLOpSYSENTEROpSYSEXITOpSYSRETOpTAKENOpUD2OpVMCALLOpVMLAUNCHOpVMRESUMEOpVMXOFFOpWAITOpWBINVDOpWRMSROpXGETBVOpXRSTOROpXSETBV"

var _OcodeKindIndex = [...]uint16{0, 3, 8, 12, 16, 20, 26, 31, 36, 41, 46, 50, 55, 60, 65, 70, 75, 79, 84, 89, 94, 99, 104, 111, 116, 121, 126, 132, 137, 142, 147, 153, 158, 165, 170, 174, 179, 185, 190, 195, 200, 204, 210, 217, 221, 228, 234, 241, 247, 254, 260, 267, 275, 281, 290, 297, 304, 312, 318, 327, 334, 340, 348, 356, 364, 372, 379, 385, 392, 400, 408, 415, 423, 429, 438, 446, 453, 461, 468, 477, 485, 489, 497, 505, 511, 520, 527, 534, 542, 548, 555, 563, 572, 578, 584, 593, 602, 609, 618, 626, 630, 635, 642, 648, 654, 659, 664, 670, 676, 681, 690, 694, 699, 703, 708, 712, 716, 720, 725, 729, 734, 739, 745, 750, 756, 761, 766, 771, 777, 782, 788, 793, 798, 803, 808, 812, 816, 821, 826, 830, 834, 840, 846, 852, 859, 866, 872, 878, 885, 893, 902, 908, 916, 925, 930, 937, 942, 950, 957, 963, 970, 976, 983, 990, 997, 1005, 1012, 1020, 1028, 1035, 1042, 1049, 1057, 1062, 1068, 1075, 1080, 1086, 1092, 1097, 1103, 1111, 1119, 1123, 1128, 1133, 1138, 1146, 1155, 1165, 1174, 1182, 1189, 1194, 1202, 1212, 1222, 1230, 1236, 1244, 1251, 1259, 1267, 1275}

const _OcodeKindLowerName = "oplopoutopdbopdwopddopresbopmovopaddopsubopandoporopxoropnotopshropshlopsaropinopaaaopaadopaamopaasopadxopalteropamxopcbwopcdqopcdqeopclcopcldopcliopcltsopcmcopcpuidopcqoopcsopcwdopcwdeopdaaopdasopdivopdsopemmsopenteropesopf2xm1opfabsopfaddpopfchsopfclexopfcomopfcompopfcomppopfcosopfdecstpopfdisiopfdivpopfdivrpopfeniopfincstpopfinitopfld1opfldl2eopfldl2topfldlg2opfldln2opfldpiopfldzopfmulpopfnclexopfndisiopfneniopfninitopfnopopfnsetpmopfpatanopfpremopfprem1opfptanopfrndintopfrstoropfsopfscaleopfsetpmopfsinopfsincosopfsqrtopfsubpopfsubrpopftstopfucomopfucompopfucomppopfxamopfxchopfxrstoropfxtractopfyl2xopfyl2xp1opgetsecopgsophltopicebpopidivopimulopintopcmpopcalloplgdtopjmpopjmp_faropjaopjaeopjbopjbeopjcopjeopjgopjgeopjlopjleopjnaopjnaeopjnbopjnbeopjncopjneopjngopjngeopjnlopjnleopjnoopjnpopjnsopjnzopjoopjpopjpeopjpoopjsopjzopintoopinvdopiretopiretdopiretqopjmpeoplahfopleaveoplfenceoploadalloplockopmfenceopmonitoropmulopmwaitopnopopntakenoppauseoppopaoppopadoppopfoppopfdoppopfqoppushaoppushadoppushfoppushfdoppushfqoprdmsroprdpmcoprdtscoprdtscpoprepoprepeoprepneopretopretfopretnoprsmopsahfopsetalcopsfenceopssopstcopstdopstiopswapgsopsyscallopsysenteropsysexitopsysretoptakenopud2opvmcallopvmlaunchopvmresumeopvmxoffopwaitopwbinvdopwrmsropxgetbvopxrstoropxsetbv"

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
	_ = x[OpJMP_FAR-(109)]
	_ = x[OpJA-(110)]
	_ = x[OpJAE-(111)]
	_ = x[OpJB-(112)]
	_ = x[OpJBE-(113)]
	_ = x[OpJC-(114)]
	_ = x[OpJE-(115)]
	_ = x[OpJG-(116)]
	_ = x[OpJGE-(117)]
	_ = x[OpJL-(118)]
	_ = x[OpJLE-(119)]
	_ = x[OpJNA-(120)]
	_ = x[OpJNAE-(121)]
	_ = x[OpJNB-(122)]
	_ = x[OpJNBE-(123)]
	_ = x[OpJNC-(124)]
	_ = x[OpJNE-(125)]
	_ = x[OpJNG-(126)]
	_ = x[OpJNGE-(127)]
	_ = x[OpJNL-(128)]
	_ = x[OpJNLE-(129)]
	_ = x[OpJNO-(130)]
	_ = x[OpJNP-(131)]
	_ = x[OpJNS-(132)]
	_ = x[OpJNZ-(133)]
	_ = x[OpJO-(134)]
	_ = x[OpJP-(135)]
	_ = x[OpJPE-(136)]
	_ = x[OpJPO-(137)]
	_ = x[OpJS-(138)]
	_ = x[OpJZ-(139)]
	_ = x[OpINTO-(140)]
	_ = x[OpINVD-(141)]
	_ = x[OpIRET-(142)]
	_ = x[OpIRETD-(143)]
	_ = x[OpIRETQ-(144)]
	_ = x[OpJMPE-(145)]
	_ = x[OpLAHF-(146)]
	_ = x[OpLEAVE-(147)]
	_ = x[OpLFENCE-(148)]
	_ = x[OpLOADALL-(149)]
	_ = x[OpLOCK-(150)]
	_ = x[OpMFENCE-(151)]
	_ = x[OpMONITOR-(152)]
	_ = x[OpMUL-(153)]
	_ = x[OpMWAIT-(154)]
	_ = x[OpNOP-(155)]
	_ = x[OpNTAKEN-(156)]
	_ = x[OpPAUSE-(157)]
	_ = x[OpPOPA-(158)]
	_ = x[OpPOPAD-(159)]
	_ = x[OpPOPF-(160)]
	_ = x[OpPOPFD-(161)]
	_ = x[OpPOPFQ-(162)]
	_ = x[OpPUSHA-(163)]
	_ = x[OpPUSHAD-(164)]
	_ = x[OpPUSHF-(165)]
	_ = x[OpPUSHFD-(166)]
	_ = x[OpPUSHFQ-(167)]
	_ = x[OpRDMSR-(168)]
	_ = x[OpRDPMC-(169)]
	_ = x[OpRDTSC-(170)]
	_ = x[OpRDTSCP-(171)]
	_ = x[OpREP-(172)]
	_ = x[OpREPE-(173)]
	_ = x[OpREPNE-(174)]
	_ = x[OpRET-(175)]
	_ = x[OpRETF-(176)]
	_ = x[OpRETN-(177)]
	_ = x[OpRSM-(178)]
	_ = x[OpSAHF-(179)]
	_ = x[OpSETALC-(180)]
	_ = x[OpSFENCE-(181)]
	_ = x[OpSS-(182)]
	_ = x[OpSTC-(183)]
	_ = x[OpSTD-(184)]
	_ = x[OpSTI-(185)]
	_ = x[OpSWAPGS-(186)]
	_ = x[OpSYSCALL-(187)]
	_ = x[OpSYSENTER-(188)]
	_ = x[OpSYSEXIT-(189)]
	_ = x[OpSYSRET-(190)]
	_ = x[OpTAKEN-(191)]
	_ = x[OpUD2-(192)]
	_ = x[OpVMCALL-(193)]
	_ = x[OpVMLAUNCH-(194)]
	_ = x[OpVMRESUME-(195)]
	_ = x[OpVMXOFF-(196)]
	_ = x[OpWAIT-(197)]
	_ = x[OpWBINVD-(198)]
	_ = x[OpWRMSR-(199)]
	_ = x[OpXGETBV-(200)]
	_ = x[OpXRSTOR-(201)]
	_ = x[OpXSETBV-(202)]
}

var _OcodeKindValues = []OcodeKind{OpL, OpOUT, OpDB, OpDW, OpDD, OpRESB, OpMOV, OpADD, OpSUB, OpAND, OpOR, OpXOR, OpNOT, OpSHR, OpSHL, OpSAR, OpIN, OpAAA, OpAAD, OpAAM, OpAAS, OpADX, OpALTER, OpAMX, OpCBW, OpCDQ, OpCDQE, OpCLC, OpCLD, OpCLI, OpCLTS, OpCMC, OpCPUID, OpCQO, OpCS, OpCWD, OpCWDE, OpDAA, OpDAS, OpDIV, OpDS, OpEMMS, OpENTER, OpES, OpF2XM1, OpFABS, OpFADDP, OpFCHS, OpFCLEX, OpFCOM, OpFCOMP, OpFCOMPP, OpFCOS, OpFDECSTP, OpFDISI, OpFDIVP, OpFDIVRP, OpFENI, OpFINCSTP, OpFINIT, OpFLD1, OpFLDL2E, OpFLDL2T, OpFLDLG2, OpFLDLN2, OpFLDPI, OpFLDZ, OpFMULP, OpFNCLEX, OpFNDISI, OpFNENI, OpFNINIT, OpFNOP, OpFNSETPM, OpFPATAN, OpFPREM, OpFPREM1, OpFPTAN, OpFRNDINT, OpFRSTOR, OpFS, OpFSCALE, OpFSETPM, OpFSIN, OpFSINCOS, OpFSQRT, OpFSUBP, OpFSUBRP, OpFTST, OpFUCOM, OpFUCOMP, OpFUCOMPP, OpFXAM, OpFXCH, OpFXRSTOR, OpFXTRACT, OpFYL2X, OpFYL2XP1, OpGETSEC, OpGS, OpHLT, OpICEBP, OpIDIV, OpIMUL, OpINT, OpCMP, OpCALL, OpLGDT, OpJMP, OpJMP_FAR, OpJA, OpJAE, OpJB, OpJBE, OpJC, OpJE, OpJG, OpJGE, OpJL, OpJLE, OpJNA, OpJNAE, OpJNB, OpJNBE, OpJNC, OpJNE, OpJNG, OpJNGE, OpJNL, OpJNLE, OpJNO, OpJNP, OpJNS, OpJNZ, OpJO, OpJP, OpJPE, OpJPO, OpJS, OpJZ, OpINTO, OpINVD, OpIRET, OpIRETD, OpIRETQ, OpJMPE, OpLAHF, OpLEAVE, OpLFENCE, OpLOADALL, OpLOCK, OpMFENCE, OpMONITOR, OpMUL, OpMWAIT, OpNOP, OpNTAKEN, OpPAUSE, OpPOPA, OpPOPAD, OpPOPF, OpPOPFD, OpPOPFQ, OpPUSHA, OpPUSHAD, OpPUSHF, OpPUSHFD, OpPUSHFQ, OpRDMSR, OpRDPMC, OpRDTSC, OpRDTSCP, OpREP, OpREPE, OpREPNE, OpRET, OpRETF, OpRETN, OpRSM, OpSAHF, OpSETALC, OpSFENCE, OpSS, OpSTC, OpSTD, OpSTI, OpSWAPGS, OpSYSCALL, OpSYSENTER, OpSYSEXIT, OpSYSRET, OpTAKEN, OpUD2, OpVMCALL, OpVMLAUNCH, OpVMRESUME, OpVMXOFF, OpWAIT, OpWBINVD, OpWRMSR, OpXGETBV, OpXRSTOR, OpXSETBV}

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
	_OcodeKindName[681:690]:        OpJMP_FAR,
	_OcodeKindLowerName[681:690]:   OpJMP_FAR,
	_OcodeKindName[690:694]:        OpJA,
	_OcodeKindLowerName[690:694]:   OpJA,
	_OcodeKindName[694:699]:        OpJAE,
	_OcodeKindLowerName[694:699]:   OpJAE,
	_OcodeKindName[699:703]:        OpJB,
	_OcodeKindLowerName[699:703]:   OpJB,
	_OcodeKindName[703:708]:        OpJBE,
	_OcodeKindLowerName[703:708]:   OpJBE,
	_OcodeKindName[708:712]:        OpJC,
	_OcodeKindLowerName[708:712]:   OpJC,
	_OcodeKindName[712:716]:        OpJE,
	_OcodeKindLowerName[712:716]:   OpJE,
	_OcodeKindName[716:720]:        OpJG,
	_OcodeKindLowerName[716:720]:   OpJG,
	_OcodeKindName[720:725]:        OpJGE,
	_OcodeKindLowerName[720:725]:   OpJGE,
	_OcodeKindName[725:729]:        OpJL,
	_OcodeKindLowerName[725:729]:   OpJL,
	_OcodeKindName[729:734]:        OpJLE,
	_OcodeKindLowerName[729:734]:   OpJLE,
	_OcodeKindName[734:739]:        OpJNA,
	_OcodeKindLowerName[734:739]:   OpJNA,
	_OcodeKindName[739:745]:        OpJNAE,
	_OcodeKindLowerName[739:745]:   OpJNAE,
	_OcodeKindName[745:750]:        OpJNB,
	_OcodeKindLowerName[745:750]:   OpJNB,
	_OcodeKindName[750:756]:        OpJNBE,
	_OcodeKindLowerName[750:756]:   OpJNBE,
	_OcodeKindName[756:761]:        OpJNC,
	_OcodeKindLowerName[756:761]:   OpJNC,
	_OcodeKindName[761:766]:        OpJNE,
	_OcodeKindLowerName[761:766]:   OpJNE,
	_OcodeKindName[766:771]:        OpJNG,
	_OcodeKindLowerName[766:771]:   OpJNG,
	_OcodeKindName[771:777]:        OpJNGE,
	_OcodeKindLowerName[771:777]:   OpJNGE,
	_OcodeKindName[777:782]:        OpJNL,
	_OcodeKindLowerName[777:782]:   OpJNL,
	_OcodeKindName[782:788]:        OpJNLE,
	_OcodeKindLowerName[782:788]:   OpJNLE,
	_OcodeKindName[788:793]:        OpJNO,
	_OcodeKindLowerName[788:793]:   OpJNO,
	_OcodeKindName[793:798]:        OpJNP,
	_OcodeKindLowerName[793:798]:   OpJNP,
	_OcodeKindName[798:803]:        OpJNS,
	_OcodeKindLowerName[798:803]:   OpJNS,
	_OcodeKindName[803:808]:        OpJNZ,
	_OcodeKindLowerName[803:808]:   OpJNZ,
	_OcodeKindName[808:812]:        OpJO,
	_OcodeKindLowerName[808:812]:   OpJO,
	_OcodeKindName[812:816]:        OpJP,
	_OcodeKindLowerName[812:816]:   OpJP,
	_OcodeKindName[816:821]:        OpJPE,
	_OcodeKindLowerName[816:821]:   OpJPE,
	_OcodeKindName[821:826]:        OpJPO,
	_OcodeKindLowerName[821:826]:   OpJPO,
	_OcodeKindName[826:830]:        OpJS,
	_OcodeKindLowerName[826:830]:   OpJS,
	_OcodeKindName[830:834]:        OpJZ,
	_OcodeKindLowerName[830:834]:   OpJZ,
	_OcodeKindName[834:840]:        OpINTO,
	_OcodeKindLowerName[834:840]:   OpINTO,
	_OcodeKindName[840:846]:        OpINVD,
	_OcodeKindLowerName[840:846]:   OpINVD,
	_OcodeKindName[846:852]:        OpIRET,
	_OcodeKindLowerName[846:852]:   OpIRET,
	_OcodeKindName[852:859]:        OpIRETD,
	_OcodeKindLowerName[852:859]:   OpIRETD,
	_OcodeKindName[859:866]:        OpIRETQ,
	_OcodeKindLowerName[859:866]:   OpIRETQ,
	_OcodeKindName[866:872]:        OpJMPE,
	_OcodeKindLowerName[866:872]:   OpJMPE,
	_OcodeKindName[872:878]:        OpLAHF,
	_OcodeKindLowerName[872:878]:   OpLAHF,
	_OcodeKindName[878:885]:        OpLEAVE,
	_OcodeKindLowerName[878:885]:   OpLEAVE,
	_OcodeKindName[885:893]:        OpLFENCE,
	_OcodeKindLowerName[885:893]:   OpLFENCE,
	_OcodeKindName[893:902]:        OpLOADALL,
	_OcodeKindLowerName[893:902]:   OpLOADALL,
	_OcodeKindName[902:908]:        OpLOCK,
	_OcodeKindLowerName[902:908]:   OpLOCK,
	_OcodeKindName[908:916]:        OpMFENCE,
	_OcodeKindLowerName[908:916]:   OpMFENCE,
	_OcodeKindName[916:925]:        OpMONITOR,
	_OcodeKindLowerName[916:925]:   OpMONITOR,
	_OcodeKindName[925:930]:        OpMUL,
	_OcodeKindLowerName[925:930]:   OpMUL,
	_OcodeKindName[930:937]:        OpMWAIT,
	_OcodeKindLowerName[930:937]:   OpMWAIT,
	_OcodeKindName[937:942]:        OpNOP,
	_OcodeKindLowerName[937:942]:   OpNOP,
	_OcodeKindName[942:950]:        OpNTAKEN,
	_OcodeKindLowerName[942:950]:   OpNTAKEN,
	_OcodeKindName[950:957]:        OpPAUSE,
	_OcodeKindLowerName[950:957]:   OpPAUSE,
	_OcodeKindName[957:963]:        OpPOPA,
	_OcodeKindLowerName[957:963]:   OpPOPA,
	_OcodeKindName[963:970]:        OpPOPAD,
	_OcodeKindLowerName[963:970]:   OpPOPAD,
	_OcodeKindName[970:976]:        OpPOPF,
	_OcodeKindLowerName[970:976]:   OpPOPF,
	_OcodeKindName[976:983]:        OpPOPFD,
	_OcodeKindLowerName[976:983]:   OpPOPFD,
	_OcodeKindName[983:990]:        OpPOPFQ,
	_OcodeKindLowerName[983:990]:   OpPOPFQ,
	_OcodeKindName[990:997]:        OpPUSHA,
	_OcodeKindLowerName[990:997]:   OpPUSHA,
	_OcodeKindName[997:1005]:       OpPUSHAD,
	_OcodeKindLowerName[997:1005]:  OpPUSHAD,
	_OcodeKindName[1005:1012]:      OpPUSHF,
	_OcodeKindLowerName[1005:1012]: OpPUSHF,
	_OcodeKindName[1012:1020]:      OpPUSHFD,
	_OcodeKindLowerName[1012:1020]: OpPUSHFD,
	_OcodeKindName[1020:1028]:      OpPUSHFQ,
	_OcodeKindLowerName[1020:1028]: OpPUSHFQ,
	_OcodeKindName[1028:1035]:      OpRDMSR,
	_OcodeKindLowerName[1028:1035]: OpRDMSR,
	_OcodeKindName[1035:1042]:      OpRDPMC,
	_OcodeKindLowerName[1035:1042]: OpRDPMC,
	_OcodeKindName[1042:1049]:      OpRDTSC,
	_OcodeKindLowerName[1042:1049]: OpRDTSC,
	_OcodeKindName[1049:1057]:      OpRDTSCP,
	_OcodeKindLowerName[1049:1057]: OpRDTSCP,
	_OcodeKindName[1057:1062]:      OpREP,
	_OcodeKindLowerName[1057:1062]: OpREP,
	_OcodeKindName[1062:1068]:      OpREPE,
	_OcodeKindLowerName[1062:1068]: OpREPE,
	_OcodeKindName[1068:1075]:      OpREPNE,
	_OcodeKindLowerName[1068:1075]: OpREPNE,
	_OcodeKindName[1075:1080]:      OpRET,
	_OcodeKindLowerName[1075:1080]: OpRET,
	_OcodeKindName[1080:1086]:      OpRETF,
	_OcodeKindLowerName[1080:1086]: OpRETF,
	_OcodeKindName[1086:1092]:      OpRETN,
	_OcodeKindLowerName[1086:1092]: OpRETN,
	_OcodeKindName[1092:1097]:      OpRSM,
	_OcodeKindLowerName[1092:1097]: OpRSM,
	_OcodeKindName[1097:1103]:      OpSAHF,
	_OcodeKindLowerName[1097:1103]: OpSAHF,
	_OcodeKindName[1103:1111]:      OpSETALC,
	_OcodeKindLowerName[1103:1111]: OpSETALC,
	_OcodeKindName[1111:1119]:      OpSFENCE,
	_OcodeKindLowerName[1111:1119]: OpSFENCE,
	_OcodeKindName[1119:1123]:      OpSS,
	_OcodeKindLowerName[1119:1123]: OpSS,
	_OcodeKindName[1123:1128]:      OpSTC,
	_OcodeKindLowerName[1123:1128]: OpSTC,
	_OcodeKindName[1128:1133]:      OpSTD,
	_OcodeKindLowerName[1128:1133]: OpSTD,
	_OcodeKindName[1133:1138]:      OpSTI,
	_OcodeKindLowerName[1133:1138]: OpSTI,
	_OcodeKindName[1138:1146]:      OpSWAPGS,
	_OcodeKindLowerName[1138:1146]: OpSWAPGS,
	_OcodeKindName[1146:1155]:      OpSYSCALL,
	_OcodeKindLowerName[1146:1155]: OpSYSCALL,
	_OcodeKindName[1155:1165]:      OpSYSENTER,
	_OcodeKindLowerName[1155:1165]: OpSYSENTER,
	_OcodeKindName[1165:1174]:      OpSYSEXIT,
	_OcodeKindLowerName[1165:1174]: OpSYSEXIT,
	_OcodeKindName[1174:1182]:      OpSYSRET,
	_OcodeKindLowerName[1174:1182]: OpSYSRET,
	_OcodeKindName[1182:1189]:      OpTAKEN,
	_OcodeKindLowerName[1182:1189]: OpTAKEN,
	_OcodeKindName[1189:1194]:      OpUD2,
	_OcodeKindLowerName[1189:1194]: OpUD2,
	_OcodeKindName[1194:1202]:      OpVMCALL,
	_OcodeKindLowerName[1194:1202]: OpVMCALL,
	_OcodeKindName[1202:1212]:      OpVMLAUNCH,
	_OcodeKindLowerName[1202:1212]: OpVMLAUNCH,
	_OcodeKindName[1212:1222]:      OpVMRESUME,
	_OcodeKindLowerName[1212:1222]: OpVMRESUME,
	_OcodeKindName[1222:1230]:      OpVMXOFF,
	_OcodeKindLowerName[1222:1230]: OpVMXOFF,
	_OcodeKindName[1230:1236]:      OpWAIT,
	_OcodeKindLowerName[1230:1236]: OpWAIT,
	_OcodeKindName[1236:1244]:      OpWBINVD,
	_OcodeKindLowerName[1236:1244]: OpWBINVD,
	_OcodeKindName[1244:1251]:      OpWRMSR,
	_OcodeKindLowerName[1244:1251]: OpWRMSR,
	_OcodeKindName[1251:1259]:      OpXGETBV,
	_OcodeKindLowerName[1251:1259]: OpXGETBV,
	_OcodeKindName[1259:1267]:      OpXRSTOR,
	_OcodeKindLowerName[1259:1267]: OpXRSTOR,
	_OcodeKindName[1267:1275]:      OpXSETBV,
	_OcodeKindLowerName[1267:1275]: OpXSETBV,
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
	_OcodeKindName[681:690],
	_OcodeKindName[690:694],
	_OcodeKindName[694:699],
	_OcodeKindName[699:703],
	_OcodeKindName[703:708],
	_OcodeKindName[708:712],
	_OcodeKindName[712:716],
	_OcodeKindName[716:720],
	_OcodeKindName[720:725],
	_OcodeKindName[725:729],
	_OcodeKindName[729:734],
	_OcodeKindName[734:739],
	_OcodeKindName[739:745],
	_OcodeKindName[745:750],
	_OcodeKindName[750:756],
	_OcodeKindName[756:761],
	_OcodeKindName[761:766],
	_OcodeKindName[766:771],
	_OcodeKindName[771:777],
	_OcodeKindName[777:782],
	_OcodeKindName[782:788],
	_OcodeKindName[788:793],
	_OcodeKindName[793:798],
	_OcodeKindName[798:803],
	_OcodeKindName[803:808],
	_OcodeKindName[808:812],
	_OcodeKindName[812:816],
	_OcodeKindName[816:821],
	_OcodeKindName[821:826],
	_OcodeKindName[826:830],
	_OcodeKindName[830:834],
	_OcodeKindName[834:840],
	_OcodeKindName[840:846],
	_OcodeKindName[846:852],
	_OcodeKindName[852:859],
	_OcodeKindName[859:866],
	_OcodeKindName[866:872],
	_OcodeKindName[872:878],
	_OcodeKindName[878:885],
	_OcodeKindName[885:893],
	_OcodeKindName[893:902],
	_OcodeKindName[902:908],
	_OcodeKindName[908:916],
	_OcodeKindName[916:925],
	_OcodeKindName[925:930],
	_OcodeKindName[930:937],
	_OcodeKindName[937:942],
	_OcodeKindName[942:950],
	_OcodeKindName[950:957],
	_OcodeKindName[957:963],
	_OcodeKindName[963:970],
	_OcodeKindName[970:976],
	_OcodeKindName[976:983],
	_OcodeKindName[983:990],
	_OcodeKindName[990:997],
	_OcodeKindName[997:1005],
	_OcodeKindName[1005:1012],
	_OcodeKindName[1012:1020],
	_OcodeKindName[1020:1028],
	_OcodeKindName[1028:1035],
	_OcodeKindName[1035:1042],
	_OcodeKindName[1042:1049],
	_OcodeKindName[1049:1057],
	_OcodeKindName[1057:1062],
	_OcodeKindName[1062:1068],
	_OcodeKindName[1068:1075],
	_OcodeKindName[1075:1080],
	_OcodeKindName[1080:1086],
	_OcodeKindName[1086:1092],
	_OcodeKindName[1092:1097],
	_OcodeKindName[1097:1103],
	_OcodeKindName[1103:1111],
	_OcodeKindName[1111:1119],
	_OcodeKindName[1119:1123],
	_OcodeKindName[1123:1128],
	_OcodeKindName[1128:1133],
	_OcodeKindName[1133:1138],
	_OcodeKindName[1138:1146],
	_OcodeKindName[1146:1155],
	_OcodeKindName[1155:1165],
	_OcodeKindName[1165:1174],
	_OcodeKindName[1174:1182],
	_OcodeKindName[1182:1189],
	_OcodeKindName[1189:1194],
	_OcodeKindName[1194:1202],
	_OcodeKindName[1202:1212],
	_OcodeKindName[1212:1222],
	_OcodeKindName[1222:1230],
	_OcodeKindName[1230:1236],
	_OcodeKindName[1236:1244],
	_OcodeKindName[1244:1251],
	_OcodeKindName[1251:1259],
	_OcodeKindName[1259:1267],
	_OcodeKindName[1267:1275],
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
