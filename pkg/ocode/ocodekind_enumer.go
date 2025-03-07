// Code generated by "enumer -type=OcodeKind -json -text"; DO NOT EDIT.

package ocode

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _OcodeKindName = "OpLOpDBOpDWOpDDOpRESBOpMOVOpADDOpAAAOpAADOpAAMOpAASOpADXOpALTEROpAMXOpCBWOpCDQOpCDQEOpCLCOpCLDOpCLIOpCLTSOpCMCOpCPUIDOpCQOOpCSOpCWDOpCWDEOpDAAOpDASOpDIVOpDSOpEMMSOpENTEROpESOpF2XM1OpFABSOpFADDPOpFCHSOpFCLEXOpFCOMOpFCOMPOpFCOMPPOpFCOSOpFDECSTPOpFDISIOpFDIVPOpFDIVRPOpFENIOpFINCSTPOpFINITOpFLD1OpFLDL2EOpFLDL2TOpFLDLG2OpFLDLN2OpFLDPIOpFLDZOpFMULPOpFNCLEXOpFNDISIOpFNENIOpFNINITOpFNOPOpFNSETPMOpFPATANOpFPREMOpFPREM1OpFPTANOpFRNDINTOpFRSTOROpFSOpFSCALEOpFSETPMOpFSINOpFSINCOSOpFSQRTOpFSUBPOpFSUBRPOpFTSTOpFUCOMOpFUCOMPOpFUCOMPPOpFXAMOpFXCHOpFXRSTOROpFXTRACTOpFYL2XOpFYL2XP1OpGETSECOpGSOpHLTOpICEBPOpIDIVOpIMULOpINTOpINTOOpINVDOpIRETOpIRETDOpIRETQOpJMPEOpLAHFOpLEAVEOpLFENCEOpLOADALLOpLOCKOpMFENCEOpMONITOROpMULOpMWAITOpNOPOpNTAKENOpPAUSEOpPOPAOpPOPADOpPOPFOpPOPFDOpPOPFQOpPUSHAOpPUSHADOpPUSHFOpPUSHFDOpPUSHFQOpRDMSROpRDPMCOpRDTSCOpRDTSCPOpREPOpREPEOpREPNEOpRETFOpRETNOpRSMOpSAHFOpSETALCOpSFENCEOpSSOpSTCOpSTDOpSTIOpSWAPGSOpSYSCALLOpSYSENTEROpSYSEXITOpSYSRETOpTAKENOpUD2OpVMCALLOpVMLAUNCHOpVMRESUMEOpVMXOFFOpWAITOpWBINVDOpWRMSROpXGETBVOpXRSTOROpXSETBV"

var _OcodeKindIndex = [...]uint16{0, 3, 7, 11, 15, 21, 26, 31, 36, 41, 46, 51, 56, 63, 68, 73, 78, 84, 89, 94, 99, 105, 110, 117, 122, 126, 131, 137, 142, 147, 152, 156, 162, 169, 173, 180, 186, 193, 199, 206, 212, 219, 227, 233, 242, 249, 256, 264, 270, 279, 286, 292, 300, 308, 316, 324, 331, 337, 344, 352, 360, 367, 375, 381, 390, 398, 405, 413, 420, 429, 437, 441, 449, 457, 463, 472, 479, 486, 494, 500, 507, 515, 524, 530, 536, 545, 554, 561, 570, 578, 582, 587, 594, 600, 606, 611, 617, 623, 629, 636, 643, 649, 655, 662, 670, 679, 685, 693, 702, 707, 714, 719, 727, 734, 740, 747, 753, 760, 767, 774, 782, 789, 797, 805, 812, 819, 826, 834, 839, 845, 852, 858, 864, 869, 875, 883, 891, 895, 900, 905, 910, 918, 927, 937, 946, 954, 961, 966, 974, 984, 994, 1002, 1008, 1016, 1023, 1031, 1039, 1047}

const _OcodeKindLowerName = "oplopdbopdwopddopresbopmovopaddopaaaopaadopaamopaasopadxopalteropamxopcbwopcdqopcdqeopclcopcldopcliopcltsopcmcopcpuidopcqoopcsopcwdopcwdeopdaaopdasopdivopdsopemmsopenteropesopf2xm1opfabsopfaddpopfchsopfclexopfcomopfcompopfcomppopfcosopfdecstpopfdisiopfdivpopfdivrpopfeniopfincstpopfinitopfld1opfldl2eopfldl2topfldlg2opfldln2opfldpiopfldzopfmulpopfnclexopfndisiopfneniopfninitopfnopopfnsetpmopfpatanopfpremopfprem1opfptanopfrndintopfrstoropfsopfscaleopfsetpmopfsinopfsincosopfsqrtopfsubpopfsubrpopftstopfucomopfucompopfucomppopfxamopfxchopfxrstoropfxtractopfyl2xopfyl2xp1opgetsecopgsophltopicebpopidivopimulopintopintoopinvdopiretopiretdopiretqopjmpeoplahfopleaveoplfenceoploadalloplockopmfenceopmonitoropmulopmwaitopnopopntakenoppauseoppopaoppopadoppopfoppopfdoppopfqoppushaoppushadoppushfoppushfdoppushfqoprdmsroprdpmcoprdtscoprdtscpoprepoprepeoprepneopretfopretnoprsmopsahfopsetalcopsfenceopssopstcopstdopstiopswapgsopsyscallopsysenteropsysexitopsysretoptakenopud2opvmcallopvmlaunchopvmresumeopvmxoffopwaitopwbinvdopwrmsropxgetbvopxrstoropxsetbv"

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
	_ = x[OpDB-(1)]
	_ = x[OpDW-(2)]
	_ = x[OpDD-(3)]
	_ = x[OpRESB-(4)]
	_ = x[OpMOV-(5)]
	_ = x[OpADD-(6)]
	_ = x[OpAAA-(7)]
	_ = x[OpAAD-(8)]
	_ = x[OpAAM-(9)]
	_ = x[OpAAS-(10)]
	_ = x[OpADX-(11)]
	_ = x[OpALTER-(12)]
	_ = x[OpAMX-(13)]
	_ = x[OpCBW-(14)]
	_ = x[OpCDQ-(15)]
	_ = x[OpCDQE-(16)]
	_ = x[OpCLC-(17)]
	_ = x[OpCLD-(18)]
	_ = x[OpCLI-(19)]
	_ = x[OpCLTS-(20)]
	_ = x[OpCMC-(21)]
	_ = x[OpCPUID-(22)]
	_ = x[OpCQO-(23)]
	_ = x[OpCS-(24)]
	_ = x[OpCWD-(25)]
	_ = x[OpCWDE-(26)]
	_ = x[OpDAA-(27)]
	_ = x[OpDAS-(28)]
	_ = x[OpDIV-(29)]
	_ = x[OpDS-(30)]
	_ = x[OpEMMS-(31)]
	_ = x[OpENTER-(32)]
	_ = x[OpES-(33)]
	_ = x[OpF2XM1-(34)]
	_ = x[OpFABS-(35)]
	_ = x[OpFADDP-(36)]
	_ = x[OpFCHS-(37)]
	_ = x[OpFCLEX-(38)]
	_ = x[OpFCOM-(39)]
	_ = x[OpFCOMP-(40)]
	_ = x[OpFCOMPP-(41)]
	_ = x[OpFCOS-(42)]
	_ = x[OpFDECSTP-(43)]
	_ = x[OpFDISI-(44)]
	_ = x[OpFDIVP-(45)]
	_ = x[OpFDIVRP-(46)]
	_ = x[OpFENI-(47)]
	_ = x[OpFINCSTP-(48)]
	_ = x[OpFINIT-(49)]
	_ = x[OpFLD1-(50)]
	_ = x[OpFLDL2E-(51)]
	_ = x[OpFLDL2T-(52)]
	_ = x[OpFLDLG2-(53)]
	_ = x[OpFLDLN2-(54)]
	_ = x[OpFLDPI-(55)]
	_ = x[OpFLDZ-(56)]
	_ = x[OpFMULP-(57)]
	_ = x[OpFNCLEX-(58)]
	_ = x[OpFNDISI-(59)]
	_ = x[OpFNENI-(60)]
	_ = x[OpFNINIT-(61)]
	_ = x[OpFNOP-(62)]
	_ = x[OpFNSETPM-(63)]
	_ = x[OpFPATAN-(64)]
	_ = x[OpFPREM-(65)]
	_ = x[OpFPREM1-(66)]
	_ = x[OpFPTAN-(67)]
	_ = x[OpFRNDINT-(68)]
	_ = x[OpFRSTOR-(69)]
	_ = x[OpFS-(70)]
	_ = x[OpFSCALE-(71)]
	_ = x[OpFSETPM-(72)]
	_ = x[OpFSIN-(73)]
	_ = x[OpFSINCOS-(74)]
	_ = x[OpFSQRT-(75)]
	_ = x[OpFSUBP-(76)]
	_ = x[OpFSUBRP-(77)]
	_ = x[OpFTST-(78)]
	_ = x[OpFUCOM-(79)]
	_ = x[OpFUCOMP-(80)]
	_ = x[OpFUCOMPP-(81)]
	_ = x[OpFXAM-(82)]
	_ = x[OpFXCH-(83)]
	_ = x[OpFXRSTOR-(84)]
	_ = x[OpFXTRACT-(85)]
	_ = x[OpFYL2X-(86)]
	_ = x[OpFYL2XP1-(87)]
	_ = x[OpGETSEC-(88)]
	_ = x[OpGS-(89)]
	_ = x[OpHLT-(90)]
	_ = x[OpICEBP-(91)]
	_ = x[OpIDIV-(92)]
	_ = x[OpIMUL-(93)]
	_ = x[OpINT-(94)]
	_ = x[OpINTO-(95)]
	_ = x[OpINVD-(96)]
	_ = x[OpIRET-(97)]
	_ = x[OpIRETD-(98)]
	_ = x[OpIRETQ-(99)]
	_ = x[OpJMPE-(100)]
	_ = x[OpLAHF-(101)]
	_ = x[OpLEAVE-(102)]
	_ = x[OpLFENCE-(103)]
	_ = x[OpLOADALL-(104)]
	_ = x[OpLOCK-(105)]
	_ = x[OpMFENCE-(106)]
	_ = x[OpMONITOR-(107)]
	_ = x[OpMUL-(108)]
	_ = x[OpMWAIT-(109)]
	_ = x[OpNOP-(110)]
	_ = x[OpNTAKEN-(111)]
	_ = x[OpPAUSE-(112)]
	_ = x[OpPOPA-(113)]
	_ = x[OpPOPAD-(114)]
	_ = x[OpPOPF-(115)]
	_ = x[OpPOPFD-(116)]
	_ = x[OpPOPFQ-(117)]
	_ = x[OpPUSHA-(118)]
	_ = x[OpPUSHAD-(119)]
	_ = x[OpPUSHF-(120)]
	_ = x[OpPUSHFD-(121)]
	_ = x[OpPUSHFQ-(122)]
	_ = x[OpRDMSR-(123)]
	_ = x[OpRDPMC-(124)]
	_ = x[OpRDTSC-(125)]
	_ = x[OpRDTSCP-(126)]
	_ = x[OpREP-(127)]
	_ = x[OpREPE-(128)]
	_ = x[OpREPNE-(129)]
	_ = x[OpRETF-(130)]
	_ = x[OpRETN-(131)]
	_ = x[OpRSM-(132)]
	_ = x[OpSAHF-(133)]
	_ = x[OpSETALC-(134)]
	_ = x[OpSFENCE-(135)]
	_ = x[OpSS-(136)]
	_ = x[OpSTC-(137)]
	_ = x[OpSTD-(138)]
	_ = x[OpSTI-(139)]
	_ = x[OpSWAPGS-(140)]
	_ = x[OpSYSCALL-(141)]
	_ = x[OpSYSENTER-(142)]
	_ = x[OpSYSEXIT-(143)]
	_ = x[OpSYSRET-(144)]
	_ = x[OpTAKEN-(145)]
	_ = x[OpUD2-(146)]
	_ = x[OpVMCALL-(147)]
	_ = x[OpVMLAUNCH-(148)]
	_ = x[OpVMRESUME-(149)]
	_ = x[OpVMXOFF-(150)]
	_ = x[OpWAIT-(151)]
	_ = x[OpWBINVD-(152)]
	_ = x[OpWRMSR-(153)]
	_ = x[OpXGETBV-(154)]
	_ = x[OpXRSTOR-(155)]
	_ = x[OpXSETBV-(156)]
}

var _OcodeKindValues = []OcodeKind{OpL, OpDB, OpDW, OpDD, OpRESB, OpMOV, OpADD, OpAAA, OpAAD, OpAAM, OpAAS, OpADX, OpALTER, OpAMX, OpCBW, OpCDQ, OpCDQE, OpCLC, OpCLD, OpCLI, OpCLTS, OpCMC, OpCPUID, OpCQO, OpCS, OpCWD, OpCWDE, OpDAA, OpDAS, OpDIV, OpDS, OpEMMS, OpENTER, OpES, OpF2XM1, OpFABS, OpFADDP, OpFCHS, OpFCLEX, OpFCOM, OpFCOMP, OpFCOMPP, OpFCOS, OpFDECSTP, OpFDISI, OpFDIVP, OpFDIVRP, OpFENI, OpFINCSTP, OpFINIT, OpFLD1, OpFLDL2E, OpFLDL2T, OpFLDLG2, OpFLDLN2, OpFLDPI, OpFLDZ, OpFMULP, OpFNCLEX, OpFNDISI, OpFNENI, OpFNINIT, OpFNOP, OpFNSETPM, OpFPATAN, OpFPREM, OpFPREM1, OpFPTAN, OpFRNDINT, OpFRSTOR, OpFS, OpFSCALE, OpFSETPM, OpFSIN, OpFSINCOS, OpFSQRT, OpFSUBP, OpFSUBRP, OpFTST, OpFUCOM, OpFUCOMP, OpFUCOMPP, OpFXAM, OpFXCH, OpFXRSTOR, OpFXTRACT, OpFYL2X, OpFYL2XP1, OpGETSEC, OpGS, OpHLT, OpICEBP, OpIDIV, OpIMUL, OpINT, OpINTO, OpINVD, OpIRET, OpIRETD, OpIRETQ, OpJMPE, OpLAHF, OpLEAVE, OpLFENCE, OpLOADALL, OpLOCK, OpMFENCE, OpMONITOR, OpMUL, OpMWAIT, OpNOP, OpNTAKEN, OpPAUSE, OpPOPA, OpPOPAD, OpPOPF, OpPOPFD, OpPOPFQ, OpPUSHA, OpPUSHAD, OpPUSHF, OpPUSHFD, OpPUSHFQ, OpRDMSR, OpRDPMC, OpRDTSC, OpRDTSCP, OpREP, OpREPE, OpREPNE, OpRETF, OpRETN, OpRSM, OpSAHF, OpSETALC, OpSFENCE, OpSS, OpSTC, OpSTD, OpSTI, OpSWAPGS, OpSYSCALL, OpSYSENTER, OpSYSEXIT, OpSYSRET, OpTAKEN, OpUD2, OpVMCALL, OpVMLAUNCH, OpVMRESUME, OpVMXOFF, OpWAIT, OpWBINVD, OpWRMSR, OpXGETBV, OpXRSTOR, OpXSETBV}

var _OcodeKindNameToValueMap = map[string]OcodeKind{
	_OcodeKindName[0:3]:            OpL,
	_OcodeKindLowerName[0:3]:       OpL,
	_OcodeKindName[3:7]:            OpDB,
	_OcodeKindLowerName[3:7]:       OpDB,
	_OcodeKindName[7:11]:           OpDW,
	_OcodeKindLowerName[7:11]:      OpDW,
	_OcodeKindName[11:15]:          OpDD,
	_OcodeKindLowerName[11:15]:     OpDD,
	_OcodeKindName[15:21]:          OpRESB,
	_OcodeKindLowerName[15:21]:     OpRESB,
	_OcodeKindName[21:26]:          OpMOV,
	_OcodeKindLowerName[21:26]:     OpMOV,
	_OcodeKindName[26:31]:          OpADD,
	_OcodeKindLowerName[26:31]:     OpADD,
	_OcodeKindName[31:36]:          OpAAA,
	_OcodeKindLowerName[31:36]:     OpAAA,
	_OcodeKindName[36:41]:          OpAAD,
	_OcodeKindLowerName[36:41]:     OpAAD,
	_OcodeKindName[41:46]:          OpAAM,
	_OcodeKindLowerName[41:46]:     OpAAM,
	_OcodeKindName[46:51]:          OpAAS,
	_OcodeKindLowerName[46:51]:     OpAAS,
	_OcodeKindName[51:56]:          OpADX,
	_OcodeKindLowerName[51:56]:     OpADX,
	_OcodeKindName[56:63]:          OpALTER,
	_OcodeKindLowerName[56:63]:     OpALTER,
	_OcodeKindName[63:68]:          OpAMX,
	_OcodeKindLowerName[63:68]:     OpAMX,
	_OcodeKindName[68:73]:          OpCBW,
	_OcodeKindLowerName[68:73]:     OpCBW,
	_OcodeKindName[73:78]:          OpCDQ,
	_OcodeKindLowerName[73:78]:     OpCDQ,
	_OcodeKindName[78:84]:          OpCDQE,
	_OcodeKindLowerName[78:84]:     OpCDQE,
	_OcodeKindName[84:89]:          OpCLC,
	_OcodeKindLowerName[84:89]:     OpCLC,
	_OcodeKindName[89:94]:          OpCLD,
	_OcodeKindLowerName[89:94]:     OpCLD,
	_OcodeKindName[94:99]:          OpCLI,
	_OcodeKindLowerName[94:99]:     OpCLI,
	_OcodeKindName[99:105]:         OpCLTS,
	_OcodeKindLowerName[99:105]:    OpCLTS,
	_OcodeKindName[105:110]:        OpCMC,
	_OcodeKindLowerName[105:110]:   OpCMC,
	_OcodeKindName[110:117]:        OpCPUID,
	_OcodeKindLowerName[110:117]:   OpCPUID,
	_OcodeKindName[117:122]:        OpCQO,
	_OcodeKindLowerName[117:122]:   OpCQO,
	_OcodeKindName[122:126]:        OpCS,
	_OcodeKindLowerName[122:126]:   OpCS,
	_OcodeKindName[126:131]:        OpCWD,
	_OcodeKindLowerName[126:131]:   OpCWD,
	_OcodeKindName[131:137]:        OpCWDE,
	_OcodeKindLowerName[131:137]:   OpCWDE,
	_OcodeKindName[137:142]:        OpDAA,
	_OcodeKindLowerName[137:142]:   OpDAA,
	_OcodeKindName[142:147]:        OpDAS,
	_OcodeKindLowerName[142:147]:   OpDAS,
	_OcodeKindName[147:152]:        OpDIV,
	_OcodeKindLowerName[147:152]:   OpDIV,
	_OcodeKindName[152:156]:        OpDS,
	_OcodeKindLowerName[152:156]:   OpDS,
	_OcodeKindName[156:162]:        OpEMMS,
	_OcodeKindLowerName[156:162]:   OpEMMS,
	_OcodeKindName[162:169]:        OpENTER,
	_OcodeKindLowerName[162:169]:   OpENTER,
	_OcodeKindName[169:173]:        OpES,
	_OcodeKindLowerName[169:173]:   OpES,
	_OcodeKindName[173:180]:        OpF2XM1,
	_OcodeKindLowerName[173:180]:   OpF2XM1,
	_OcodeKindName[180:186]:        OpFABS,
	_OcodeKindLowerName[180:186]:   OpFABS,
	_OcodeKindName[186:193]:        OpFADDP,
	_OcodeKindLowerName[186:193]:   OpFADDP,
	_OcodeKindName[193:199]:        OpFCHS,
	_OcodeKindLowerName[193:199]:   OpFCHS,
	_OcodeKindName[199:206]:        OpFCLEX,
	_OcodeKindLowerName[199:206]:   OpFCLEX,
	_OcodeKindName[206:212]:        OpFCOM,
	_OcodeKindLowerName[206:212]:   OpFCOM,
	_OcodeKindName[212:219]:        OpFCOMP,
	_OcodeKindLowerName[212:219]:   OpFCOMP,
	_OcodeKindName[219:227]:        OpFCOMPP,
	_OcodeKindLowerName[219:227]:   OpFCOMPP,
	_OcodeKindName[227:233]:        OpFCOS,
	_OcodeKindLowerName[227:233]:   OpFCOS,
	_OcodeKindName[233:242]:        OpFDECSTP,
	_OcodeKindLowerName[233:242]:   OpFDECSTP,
	_OcodeKindName[242:249]:        OpFDISI,
	_OcodeKindLowerName[242:249]:   OpFDISI,
	_OcodeKindName[249:256]:        OpFDIVP,
	_OcodeKindLowerName[249:256]:   OpFDIVP,
	_OcodeKindName[256:264]:        OpFDIVRP,
	_OcodeKindLowerName[256:264]:   OpFDIVRP,
	_OcodeKindName[264:270]:        OpFENI,
	_OcodeKindLowerName[264:270]:   OpFENI,
	_OcodeKindName[270:279]:        OpFINCSTP,
	_OcodeKindLowerName[270:279]:   OpFINCSTP,
	_OcodeKindName[279:286]:        OpFINIT,
	_OcodeKindLowerName[279:286]:   OpFINIT,
	_OcodeKindName[286:292]:        OpFLD1,
	_OcodeKindLowerName[286:292]:   OpFLD1,
	_OcodeKindName[292:300]:        OpFLDL2E,
	_OcodeKindLowerName[292:300]:   OpFLDL2E,
	_OcodeKindName[300:308]:        OpFLDL2T,
	_OcodeKindLowerName[300:308]:   OpFLDL2T,
	_OcodeKindName[308:316]:        OpFLDLG2,
	_OcodeKindLowerName[308:316]:   OpFLDLG2,
	_OcodeKindName[316:324]:        OpFLDLN2,
	_OcodeKindLowerName[316:324]:   OpFLDLN2,
	_OcodeKindName[324:331]:        OpFLDPI,
	_OcodeKindLowerName[324:331]:   OpFLDPI,
	_OcodeKindName[331:337]:        OpFLDZ,
	_OcodeKindLowerName[331:337]:   OpFLDZ,
	_OcodeKindName[337:344]:        OpFMULP,
	_OcodeKindLowerName[337:344]:   OpFMULP,
	_OcodeKindName[344:352]:        OpFNCLEX,
	_OcodeKindLowerName[344:352]:   OpFNCLEX,
	_OcodeKindName[352:360]:        OpFNDISI,
	_OcodeKindLowerName[352:360]:   OpFNDISI,
	_OcodeKindName[360:367]:        OpFNENI,
	_OcodeKindLowerName[360:367]:   OpFNENI,
	_OcodeKindName[367:375]:        OpFNINIT,
	_OcodeKindLowerName[367:375]:   OpFNINIT,
	_OcodeKindName[375:381]:        OpFNOP,
	_OcodeKindLowerName[375:381]:   OpFNOP,
	_OcodeKindName[381:390]:        OpFNSETPM,
	_OcodeKindLowerName[381:390]:   OpFNSETPM,
	_OcodeKindName[390:398]:        OpFPATAN,
	_OcodeKindLowerName[390:398]:   OpFPATAN,
	_OcodeKindName[398:405]:        OpFPREM,
	_OcodeKindLowerName[398:405]:   OpFPREM,
	_OcodeKindName[405:413]:        OpFPREM1,
	_OcodeKindLowerName[405:413]:   OpFPREM1,
	_OcodeKindName[413:420]:        OpFPTAN,
	_OcodeKindLowerName[413:420]:   OpFPTAN,
	_OcodeKindName[420:429]:        OpFRNDINT,
	_OcodeKindLowerName[420:429]:   OpFRNDINT,
	_OcodeKindName[429:437]:        OpFRSTOR,
	_OcodeKindLowerName[429:437]:   OpFRSTOR,
	_OcodeKindName[437:441]:        OpFS,
	_OcodeKindLowerName[437:441]:   OpFS,
	_OcodeKindName[441:449]:        OpFSCALE,
	_OcodeKindLowerName[441:449]:   OpFSCALE,
	_OcodeKindName[449:457]:        OpFSETPM,
	_OcodeKindLowerName[449:457]:   OpFSETPM,
	_OcodeKindName[457:463]:        OpFSIN,
	_OcodeKindLowerName[457:463]:   OpFSIN,
	_OcodeKindName[463:472]:        OpFSINCOS,
	_OcodeKindLowerName[463:472]:   OpFSINCOS,
	_OcodeKindName[472:479]:        OpFSQRT,
	_OcodeKindLowerName[472:479]:   OpFSQRT,
	_OcodeKindName[479:486]:        OpFSUBP,
	_OcodeKindLowerName[479:486]:   OpFSUBP,
	_OcodeKindName[486:494]:        OpFSUBRP,
	_OcodeKindLowerName[486:494]:   OpFSUBRP,
	_OcodeKindName[494:500]:        OpFTST,
	_OcodeKindLowerName[494:500]:   OpFTST,
	_OcodeKindName[500:507]:        OpFUCOM,
	_OcodeKindLowerName[500:507]:   OpFUCOM,
	_OcodeKindName[507:515]:        OpFUCOMP,
	_OcodeKindLowerName[507:515]:   OpFUCOMP,
	_OcodeKindName[515:524]:        OpFUCOMPP,
	_OcodeKindLowerName[515:524]:   OpFUCOMPP,
	_OcodeKindName[524:530]:        OpFXAM,
	_OcodeKindLowerName[524:530]:   OpFXAM,
	_OcodeKindName[530:536]:        OpFXCH,
	_OcodeKindLowerName[530:536]:   OpFXCH,
	_OcodeKindName[536:545]:        OpFXRSTOR,
	_OcodeKindLowerName[536:545]:   OpFXRSTOR,
	_OcodeKindName[545:554]:        OpFXTRACT,
	_OcodeKindLowerName[545:554]:   OpFXTRACT,
	_OcodeKindName[554:561]:        OpFYL2X,
	_OcodeKindLowerName[554:561]:   OpFYL2X,
	_OcodeKindName[561:570]:        OpFYL2XP1,
	_OcodeKindLowerName[561:570]:   OpFYL2XP1,
	_OcodeKindName[570:578]:        OpGETSEC,
	_OcodeKindLowerName[570:578]:   OpGETSEC,
	_OcodeKindName[578:582]:        OpGS,
	_OcodeKindLowerName[578:582]:   OpGS,
	_OcodeKindName[582:587]:        OpHLT,
	_OcodeKindLowerName[582:587]:   OpHLT,
	_OcodeKindName[587:594]:        OpICEBP,
	_OcodeKindLowerName[587:594]:   OpICEBP,
	_OcodeKindName[594:600]:        OpIDIV,
	_OcodeKindLowerName[594:600]:   OpIDIV,
	_OcodeKindName[600:606]:        OpIMUL,
	_OcodeKindLowerName[600:606]:   OpIMUL,
	_OcodeKindName[606:611]:        OpINT,
	_OcodeKindLowerName[606:611]:   OpINT,
	_OcodeKindName[611:617]:        OpINTO,
	_OcodeKindLowerName[611:617]:   OpINTO,
	_OcodeKindName[617:623]:        OpINVD,
	_OcodeKindLowerName[617:623]:   OpINVD,
	_OcodeKindName[623:629]:        OpIRET,
	_OcodeKindLowerName[623:629]:   OpIRET,
	_OcodeKindName[629:636]:        OpIRETD,
	_OcodeKindLowerName[629:636]:   OpIRETD,
	_OcodeKindName[636:643]:        OpIRETQ,
	_OcodeKindLowerName[636:643]:   OpIRETQ,
	_OcodeKindName[643:649]:        OpJMPE,
	_OcodeKindLowerName[643:649]:   OpJMPE,
	_OcodeKindName[649:655]:        OpLAHF,
	_OcodeKindLowerName[649:655]:   OpLAHF,
	_OcodeKindName[655:662]:        OpLEAVE,
	_OcodeKindLowerName[655:662]:   OpLEAVE,
	_OcodeKindName[662:670]:        OpLFENCE,
	_OcodeKindLowerName[662:670]:   OpLFENCE,
	_OcodeKindName[670:679]:        OpLOADALL,
	_OcodeKindLowerName[670:679]:   OpLOADALL,
	_OcodeKindName[679:685]:        OpLOCK,
	_OcodeKindLowerName[679:685]:   OpLOCK,
	_OcodeKindName[685:693]:        OpMFENCE,
	_OcodeKindLowerName[685:693]:   OpMFENCE,
	_OcodeKindName[693:702]:        OpMONITOR,
	_OcodeKindLowerName[693:702]:   OpMONITOR,
	_OcodeKindName[702:707]:        OpMUL,
	_OcodeKindLowerName[702:707]:   OpMUL,
	_OcodeKindName[707:714]:        OpMWAIT,
	_OcodeKindLowerName[707:714]:   OpMWAIT,
	_OcodeKindName[714:719]:        OpNOP,
	_OcodeKindLowerName[714:719]:   OpNOP,
	_OcodeKindName[719:727]:        OpNTAKEN,
	_OcodeKindLowerName[719:727]:   OpNTAKEN,
	_OcodeKindName[727:734]:        OpPAUSE,
	_OcodeKindLowerName[727:734]:   OpPAUSE,
	_OcodeKindName[734:740]:        OpPOPA,
	_OcodeKindLowerName[734:740]:   OpPOPA,
	_OcodeKindName[740:747]:        OpPOPAD,
	_OcodeKindLowerName[740:747]:   OpPOPAD,
	_OcodeKindName[747:753]:        OpPOPF,
	_OcodeKindLowerName[747:753]:   OpPOPF,
	_OcodeKindName[753:760]:        OpPOPFD,
	_OcodeKindLowerName[753:760]:   OpPOPFD,
	_OcodeKindName[760:767]:        OpPOPFQ,
	_OcodeKindLowerName[760:767]:   OpPOPFQ,
	_OcodeKindName[767:774]:        OpPUSHA,
	_OcodeKindLowerName[767:774]:   OpPUSHA,
	_OcodeKindName[774:782]:        OpPUSHAD,
	_OcodeKindLowerName[774:782]:   OpPUSHAD,
	_OcodeKindName[782:789]:        OpPUSHF,
	_OcodeKindLowerName[782:789]:   OpPUSHF,
	_OcodeKindName[789:797]:        OpPUSHFD,
	_OcodeKindLowerName[789:797]:   OpPUSHFD,
	_OcodeKindName[797:805]:        OpPUSHFQ,
	_OcodeKindLowerName[797:805]:   OpPUSHFQ,
	_OcodeKindName[805:812]:        OpRDMSR,
	_OcodeKindLowerName[805:812]:   OpRDMSR,
	_OcodeKindName[812:819]:        OpRDPMC,
	_OcodeKindLowerName[812:819]:   OpRDPMC,
	_OcodeKindName[819:826]:        OpRDTSC,
	_OcodeKindLowerName[819:826]:   OpRDTSC,
	_OcodeKindName[826:834]:        OpRDTSCP,
	_OcodeKindLowerName[826:834]:   OpRDTSCP,
	_OcodeKindName[834:839]:        OpREP,
	_OcodeKindLowerName[834:839]:   OpREP,
	_OcodeKindName[839:845]:        OpREPE,
	_OcodeKindLowerName[839:845]:   OpREPE,
	_OcodeKindName[845:852]:        OpREPNE,
	_OcodeKindLowerName[845:852]:   OpREPNE,
	_OcodeKindName[852:858]:        OpRETF,
	_OcodeKindLowerName[852:858]:   OpRETF,
	_OcodeKindName[858:864]:        OpRETN,
	_OcodeKindLowerName[858:864]:   OpRETN,
	_OcodeKindName[864:869]:        OpRSM,
	_OcodeKindLowerName[864:869]:   OpRSM,
	_OcodeKindName[869:875]:        OpSAHF,
	_OcodeKindLowerName[869:875]:   OpSAHF,
	_OcodeKindName[875:883]:        OpSETALC,
	_OcodeKindLowerName[875:883]:   OpSETALC,
	_OcodeKindName[883:891]:        OpSFENCE,
	_OcodeKindLowerName[883:891]:   OpSFENCE,
	_OcodeKindName[891:895]:        OpSS,
	_OcodeKindLowerName[891:895]:   OpSS,
	_OcodeKindName[895:900]:        OpSTC,
	_OcodeKindLowerName[895:900]:   OpSTC,
	_OcodeKindName[900:905]:        OpSTD,
	_OcodeKindLowerName[900:905]:   OpSTD,
	_OcodeKindName[905:910]:        OpSTI,
	_OcodeKindLowerName[905:910]:   OpSTI,
	_OcodeKindName[910:918]:        OpSWAPGS,
	_OcodeKindLowerName[910:918]:   OpSWAPGS,
	_OcodeKindName[918:927]:        OpSYSCALL,
	_OcodeKindLowerName[918:927]:   OpSYSCALL,
	_OcodeKindName[927:937]:        OpSYSENTER,
	_OcodeKindLowerName[927:937]:   OpSYSENTER,
	_OcodeKindName[937:946]:        OpSYSEXIT,
	_OcodeKindLowerName[937:946]:   OpSYSEXIT,
	_OcodeKindName[946:954]:        OpSYSRET,
	_OcodeKindLowerName[946:954]:   OpSYSRET,
	_OcodeKindName[954:961]:        OpTAKEN,
	_OcodeKindLowerName[954:961]:   OpTAKEN,
	_OcodeKindName[961:966]:        OpUD2,
	_OcodeKindLowerName[961:966]:   OpUD2,
	_OcodeKindName[966:974]:        OpVMCALL,
	_OcodeKindLowerName[966:974]:   OpVMCALL,
	_OcodeKindName[974:984]:        OpVMLAUNCH,
	_OcodeKindLowerName[974:984]:   OpVMLAUNCH,
	_OcodeKindName[984:994]:        OpVMRESUME,
	_OcodeKindLowerName[984:994]:   OpVMRESUME,
	_OcodeKindName[994:1002]:       OpVMXOFF,
	_OcodeKindLowerName[994:1002]:  OpVMXOFF,
	_OcodeKindName[1002:1008]:      OpWAIT,
	_OcodeKindLowerName[1002:1008]: OpWAIT,
	_OcodeKindName[1008:1016]:      OpWBINVD,
	_OcodeKindLowerName[1008:1016]: OpWBINVD,
	_OcodeKindName[1016:1023]:      OpWRMSR,
	_OcodeKindLowerName[1016:1023]: OpWRMSR,
	_OcodeKindName[1023:1031]:      OpXGETBV,
	_OcodeKindLowerName[1023:1031]: OpXGETBV,
	_OcodeKindName[1031:1039]:      OpXRSTOR,
	_OcodeKindLowerName[1031:1039]: OpXRSTOR,
	_OcodeKindName[1039:1047]:      OpXSETBV,
	_OcodeKindLowerName[1039:1047]: OpXSETBV,
}

var _OcodeKindNames = []string{
	_OcodeKindName[0:3],
	_OcodeKindName[3:7],
	_OcodeKindName[7:11],
	_OcodeKindName[11:15],
	_OcodeKindName[15:21],
	_OcodeKindName[21:26],
	_OcodeKindName[26:31],
	_OcodeKindName[31:36],
	_OcodeKindName[36:41],
	_OcodeKindName[41:46],
	_OcodeKindName[46:51],
	_OcodeKindName[51:56],
	_OcodeKindName[56:63],
	_OcodeKindName[63:68],
	_OcodeKindName[68:73],
	_OcodeKindName[73:78],
	_OcodeKindName[78:84],
	_OcodeKindName[84:89],
	_OcodeKindName[89:94],
	_OcodeKindName[94:99],
	_OcodeKindName[99:105],
	_OcodeKindName[105:110],
	_OcodeKindName[110:117],
	_OcodeKindName[117:122],
	_OcodeKindName[122:126],
	_OcodeKindName[126:131],
	_OcodeKindName[131:137],
	_OcodeKindName[137:142],
	_OcodeKindName[142:147],
	_OcodeKindName[147:152],
	_OcodeKindName[152:156],
	_OcodeKindName[156:162],
	_OcodeKindName[162:169],
	_OcodeKindName[169:173],
	_OcodeKindName[173:180],
	_OcodeKindName[180:186],
	_OcodeKindName[186:193],
	_OcodeKindName[193:199],
	_OcodeKindName[199:206],
	_OcodeKindName[206:212],
	_OcodeKindName[212:219],
	_OcodeKindName[219:227],
	_OcodeKindName[227:233],
	_OcodeKindName[233:242],
	_OcodeKindName[242:249],
	_OcodeKindName[249:256],
	_OcodeKindName[256:264],
	_OcodeKindName[264:270],
	_OcodeKindName[270:279],
	_OcodeKindName[279:286],
	_OcodeKindName[286:292],
	_OcodeKindName[292:300],
	_OcodeKindName[300:308],
	_OcodeKindName[308:316],
	_OcodeKindName[316:324],
	_OcodeKindName[324:331],
	_OcodeKindName[331:337],
	_OcodeKindName[337:344],
	_OcodeKindName[344:352],
	_OcodeKindName[352:360],
	_OcodeKindName[360:367],
	_OcodeKindName[367:375],
	_OcodeKindName[375:381],
	_OcodeKindName[381:390],
	_OcodeKindName[390:398],
	_OcodeKindName[398:405],
	_OcodeKindName[405:413],
	_OcodeKindName[413:420],
	_OcodeKindName[420:429],
	_OcodeKindName[429:437],
	_OcodeKindName[437:441],
	_OcodeKindName[441:449],
	_OcodeKindName[449:457],
	_OcodeKindName[457:463],
	_OcodeKindName[463:472],
	_OcodeKindName[472:479],
	_OcodeKindName[479:486],
	_OcodeKindName[486:494],
	_OcodeKindName[494:500],
	_OcodeKindName[500:507],
	_OcodeKindName[507:515],
	_OcodeKindName[515:524],
	_OcodeKindName[524:530],
	_OcodeKindName[530:536],
	_OcodeKindName[536:545],
	_OcodeKindName[545:554],
	_OcodeKindName[554:561],
	_OcodeKindName[561:570],
	_OcodeKindName[570:578],
	_OcodeKindName[578:582],
	_OcodeKindName[582:587],
	_OcodeKindName[587:594],
	_OcodeKindName[594:600],
	_OcodeKindName[600:606],
	_OcodeKindName[606:611],
	_OcodeKindName[611:617],
	_OcodeKindName[617:623],
	_OcodeKindName[623:629],
	_OcodeKindName[629:636],
	_OcodeKindName[636:643],
	_OcodeKindName[643:649],
	_OcodeKindName[649:655],
	_OcodeKindName[655:662],
	_OcodeKindName[662:670],
	_OcodeKindName[670:679],
	_OcodeKindName[679:685],
	_OcodeKindName[685:693],
	_OcodeKindName[693:702],
	_OcodeKindName[702:707],
	_OcodeKindName[707:714],
	_OcodeKindName[714:719],
	_OcodeKindName[719:727],
	_OcodeKindName[727:734],
	_OcodeKindName[734:740],
	_OcodeKindName[740:747],
	_OcodeKindName[747:753],
	_OcodeKindName[753:760],
	_OcodeKindName[760:767],
	_OcodeKindName[767:774],
	_OcodeKindName[774:782],
	_OcodeKindName[782:789],
	_OcodeKindName[789:797],
	_OcodeKindName[797:805],
	_OcodeKindName[805:812],
	_OcodeKindName[812:819],
	_OcodeKindName[819:826],
	_OcodeKindName[826:834],
	_OcodeKindName[834:839],
	_OcodeKindName[839:845],
	_OcodeKindName[845:852],
	_OcodeKindName[852:858],
	_OcodeKindName[858:864],
	_OcodeKindName[864:869],
	_OcodeKindName[869:875],
	_OcodeKindName[875:883],
	_OcodeKindName[883:891],
	_OcodeKindName[891:895],
	_OcodeKindName[895:900],
	_OcodeKindName[900:905],
	_OcodeKindName[905:910],
	_OcodeKindName[910:918],
	_OcodeKindName[918:927],
	_OcodeKindName[927:937],
	_OcodeKindName[937:946],
	_OcodeKindName[946:954],
	_OcodeKindName[954:961],
	_OcodeKindName[961:966],
	_OcodeKindName[966:974],
	_OcodeKindName[974:984],
	_OcodeKindName[984:994],
	_OcodeKindName[994:1002],
	_OcodeKindName[1002:1008],
	_OcodeKindName[1008:1016],
	_OcodeKindName[1016:1023],
	_OcodeKindName[1023:1031],
	_OcodeKindName[1031:1039],
	_OcodeKindName[1039:1047],
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
