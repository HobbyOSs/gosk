// Code generated by "enumer -type=OcodeKind -json -text"; DO NOT EDIT.

package ocode

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _OcodeKindName = "OpLOpDBOpDWOpDDOpRESBOpMOVOpADDOpAAAOpAADOpAAMOpAASOpADXOpALTEROpAMXOpCBWOpCDQOpCDQEOpCLCOpCLDOpCLIOpCLTSOpCMCOpCPUIDOpCQOOpCSOpCWDOpCWDEOpDAAOpDASOpDIVOpDSOpEMMSOpENTEROpESOpF2XM1OpFABSOpFADDPOpFCHSOpFCLEXOpFCOMOpFCOMPOpFCOMPPOpFCOSOpFDECSTPOpFDISIOpFDIVPOpFDIVRPOpFENIOpFINCSTPOpFINITOpFLD1OpFLDL2EOpFLDL2TOpFLDLG2OpFLDLN2OpFLDPIOpFLDZOpFMULPOpFNCLEXOpFNDISIOpFNENIOpFNINITOpFNOPOpFNSETPMOpFPATANOpFPREMOpFPREM1OpFPTANOpFRNDINTOpFRSTOROpFSOpFSCALEOpFSETPMOpFSINOpFSINCOSOpFSQRTOpFSUBPOpFSUBRPOpFTSTOpFUCOMOpFUCOMPOpFUCOMPPOpFXAMOpFXCHOpFXRSTOROpFXTRACTOpFYL2XOpFYL2XP1OpGETSECOpGSOpHLTOpICEBPOpIDIVOpIMULOpINTOpCMPOpJMPOpINTOOpINVDOpIRETOpIRETDOpIRETQOpJMPEOpLAHFOpLEAVEOpLFENCEOpLOADALLOpLOCKOpMFENCEOpMONITOROpMULOpMWAITOpNOPOpNTAKENOpPAUSEOpPOPAOpPOPADOpPOPFOpPOPFDOpPOPFQOpPUSHAOpPUSHADOpPUSHFOpPUSHFDOpPUSHFQOpRDMSROpRDPMCOpRDTSCOpRDTSCPOpREPOpREPEOpREPNEOpRETFOpRETNOpRSMOpSAHFOpSETALCOpSFENCEOpSSOpSTCOpSTDOpSTIOpSWAPGSOpSYSCALLOpSYSENTEROpSYSEXITOpSYSRETOpTAKENOpUD2OpVMCALLOpVMLAUNCHOpVMRESUMEOpVMXOFFOpWAITOpWBINVDOpWRMSROpXGETBVOpXRSTOROpXSETBV"

var _OcodeKindIndex = [...]uint16{0, 3, 7, 11, 15, 21, 26, 31, 36, 41, 46, 51, 56, 63, 68, 73, 78, 84, 89, 94, 99, 105, 110, 117, 122, 126, 131, 137, 142, 147, 152, 156, 162, 169, 173, 180, 186, 193, 199, 206, 212, 219, 227, 233, 242, 249, 256, 264, 270, 279, 286, 292, 300, 308, 316, 324, 331, 337, 344, 352, 360, 367, 375, 381, 390, 398, 405, 413, 420, 429, 437, 441, 449, 457, 463, 472, 479, 486, 494, 500, 507, 515, 524, 530, 536, 545, 554, 561, 570, 578, 582, 587, 594, 600, 606, 611, 616, 621, 627, 633, 639, 646, 653, 659, 665, 672, 680, 689, 695, 703, 712, 717, 724, 729, 737, 744, 750, 757, 763, 770, 777, 784, 792, 799, 807, 815, 822, 829, 836, 844, 849, 855, 862, 868, 874, 879, 885, 893, 901, 905, 910, 915, 920, 928, 937, 947, 956, 964, 971, 976, 984, 994, 1004, 1012, 1018, 1026, 1033, 1041, 1049, 1057}

const _OcodeKindLowerName = "oplopdbopdwopddopresbopmovopaddopaaaopaadopaamopaasopadxopalteropamxopcbwopcdqopcdqeopclcopcldopcliopcltsopcmcopcpuidopcqoopcsopcwdopcwdeopdaaopdasopdivopdsopemmsopenteropesopf2xm1opfabsopfaddpopfchsopfclexopfcomopfcompopfcomppopfcosopfdecstpopfdisiopfdivpopfdivrpopfeniopfincstpopfinitopfld1opfldl2eopfldl2topfldlg2opfldln2opfldpiopfldzopfmulpopfnclexopfndisiopfneniopfninitopfnopopfnsetpmopfpatanopfpremopfprem1opfptanopfrndintopfrstoropfsopfscaleopfsetpmopfsinopfsincosopfsqrtopfsubpopfsubrpopftstopfucomopfucompopfucomppopfxamopfxchopfxrstoropfxtractopfyl2xopfyl2xp1opgetsecopgsophltopicebpopidivopimulopintopcmpopjmpopintoopinvdopiretopiretdopiretqopjmpeoplahfopleaveoplfenceoploadalloplockopmfenceopmonitoropmulopmwaitopnopopntakenoppauseoppopaoppopadoppopfoppopfdoppopfqoppushaoppushadoppushfoppushfdoppushfqoprdmsroprdpmcoprdtscoprdtscpoprepoprepeoprepneopretfopretnoprsmopsahfopsetalcopsfenceopssopstcopstdopstiopswapgsopsyscallopsysenteropsysexitopsysretoptakenopud2opvmcallopvmlaunchopvmresumeopvmxoffopwaitopwbinvdopwrmsropxgetbvopxrstoropxsetbv"

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
	_ = x[OpCMP-(95)]
	_ = x[OpJMP-(96)]
	_ = x[OpINTO-(97)]
	_ = x[OpINVD-(98)]
	_ = x[OpIRET-(99)]
	_ = x[OpIRETD-(100)]
	_ = x[OpIRETQ-(101)]
	_ = x[OpJMPE-(102)]
	_ = x[OpLAHF-(103)]
	_ = x[OpLEAVE-(104)]
	_ = x[OpLFENCE-(105)]
	_ = x[OpLOADALL-(106)]
	_ = x[OpLOCK-(107)]
	_ = x[OpMFENCE-(108)]
	_ = x[OpMONITOR-(109)]
	_ = x[OpMUL-(110)]
	_ = x[OpMWAIT-(111)]
	_ = x[OpNOP-(112)]
	_ = x[OpNTAKEN-(113)]
	_ = x[OpPAUSE-(114)]
	_ = x[OpPOPA-(115)]
	_ = x[OpPOPAD-(116)]
	_ = x[OpPOPF-(117)]
	_ = x[OpPOPFD-(118)]
	_ = x[OpPOPFQ-(119)]
	_ = x[OpPUSHA-(120)]
	_ = x[OpPUSHAD-(121)]
	_ = x[OpPUSHF-(122)]
	_ = x[OpPUSHFD-(123)]
	_ = x[OpPUSHFQ-(124)]
	_ = x[OpRDMSR-(125)]
	_ = x[OpRDPMC-(126)]
	_ = x[OpRDTSC-(127)]
	_ = x[OpRDTSCP-(128)]
	_ = x[OpREP-(129)]
	_ = x[OpREPE-(130)]
	_ = x[OpREPNE-(131)]
	_ = x[OpRETF-(132)]
	_ = x[OpRETN-(133)]
	_ = x[OpRSM-(134)]
	_ = x[OpSAHF-(135)]
	_ = x[OpSETALC-(136)]
	_ = x[OpSFENCE-(137)]
	_ = x[OpSS-(138)]
	_ = x[OpSTC-(139)]
	_ = x[OpSTD-(140)]
	_ = x[OpSTI-(141)]
	_ = x[OpSWAPGS-(142)]
	_ = x[OpSYSCALL-(143)]
	_ = x[OpSYSENTER-(144)]
	_ = x[OpSYSEXIT-(145)]
	_ = x[OpSYSRET-(146)]
	_ = x[OpTAKEN-(147)]
	_ = x[OpUD2-(148)]
	_ = x[OpVMCALL-(149)]
	_ = x[OpVMLAUNCH-(150)]
	_ = x[OpVMRESUME-(151)]
	_ = x[OpVMXOFF-(152)]
	_ = x[OpWAIT-(153)]
	_ = x[OpWBINVD-(154)]
	_ = x[OpWRMSR-(155)]
	_ = x[OpXGETBV-(156)]
	_ = x[OpXRSTOR-(157)]
	_ = x[OpXSETBV-(158)]
}

var _OcodeKindValues = []OcodeKind{OpL, OpDB, OpDW, OpDD, OpRESB, OpMOV, OpADD, OpAAA, OpAAD, OpAAM, OpAAS, OpADX, OpALTER, OpAMX, OpCBW, OpCDQ, OpCDQE, OpCLC, OpCLD, OpCLI, OpCLTS, OpCMC, OpCPUID, OpCQO, OpCS, OpCWD, OpCWDE, OpDAA, OpDAS, OpDIV, OpDS, OpEMMS, OpENTER, OpES, OpF2XM1, OpFABS, OpFADDP, OpFCHS, OpFCLEX, OpFCOM, OpFCOMP, OpFCOMPP, OpFCOS, OpFDECSTP, OpFDISI, OpFDIVP, OpFDIVRP, OpFENI, OpFINCSTP, OpFINIT, OpFLD1, OpFLDL2E, OpFLDL2T, OpFLDLG2, OpFLDLN2, OpFLDPI, OpFLDZ, OpFMULP, OpFNCLEX, OpFNDISI, OpFNENI, OpFNINIT, OpFNOP, OpFNSETPM, OpFPATAN, OpFPREM, OpFPREM1, OpFPTAN, OpFRNDINT, OpFRSTOR, OpFS, OpFSCALE, OpFSETPM, OpFSIN, OpFSINCOS, OpFSQRT, OpFSUBP, OpFSUBRP, OpFTST, OpFUCOM, OpFUCOMP, OpFUCOMPP, OpFXAM, OpFXCH, OpFXRSTOR, OpFXTRACT, OpFYL2X, OpFYL2XP1, OpGETSEC, OpGS, OpHLT, OpICEBP, OpIDIV, OpIMUL, OpINT, OpCMP, OpJMP, OpINTO, OpINVD, OpIRET, OpIRETD, OpIRETQ, OpJMPE, OpLAHF, OpLEAVE, OpLFENCE, OpLOADALL, OpLOCK, OpMFENCE, OpMONITOR, OpMUL, OpMWAIT, OpNOP, OpNTAKEN, OpPAUSE, OpPOPA, OpPOPAD, OpPOPF, OpPOPFD, OpPOPFQ, OpPUSHA, OpPUSHAD, OpPUSHF, OpPUSHFD, OpPUSHFQ, OpRDMSR, OpRDPMC, OpRDTSC, OpRDTSCP, OpREP, OpREPE, OpREPNE, OpRETF, OpRETN, OpRSM, OpSAHF, OpSETALC, OpSFENCE, OpSS, OpSTC, OpSTD, OpSTI, OpSWAPGS, OpSYSCALL, OpSYSENTER, OpSYSEXIT, OpSYSRET, OpTAKEN, OpUD2, OpVMCALL, OpVMLAUNCH, OpVMRESUME, OpVMXOFF, OpWAIT, OpWBINVD, OpWRMSR, OpXGETBV, OpXRSTOR, OpXSETBV}

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
	_OcodeKindName[611:616]:        OpCMP,
	_OcodeKindLowerName[611:616]:   OpCMP,
	_OcodeKindName[616:621]:        OpJMP,
	_OcodeKindLowerName[616:621]:   OpJMP,
	_OcodeKindName[621:627]:        OpINTO,
	_OcodeKindLowerName[621:627]:   OpINTO,
	_OcodeKindName[627:633]:        OpINVD,
	_OcodeKindLowerName[627:633]:   OpINVD,
	_OcodeKindName[633:639]:        OpIRET,
	_OcodeKindLowerName[633:639]:   OpIRET,
	_OcodeKindName[639:646]:        OpIRETD,
	_OcodeKindLowerName[639:646]:   OpIRETD,
	_OcodeKindName[646:653]:        OpIRETQ,
	_OcodeKindLowerName[646:653]:   OpIRETQ,
	_OcodeKindName[653:659]:        OpJMPE,
	_OcodeKindLowerName[653:659]:   OpJMPE,
	_OcodeKindName[659:665]:        OpLAHF,
	_OcodeKindLowerName[659:665]:   OpLAHF,
	_OcodeKindName[665:672]:        OpLEAVE,
	_OcodeKindLowerName[665:672]:   OpLEAVE,
	_OcodeKindName[672:680]:        OpLFENCE,
	_OcodeKindLowerName[672:680]:   OpLFENCE,
	_OcodeKindName[680:689]:        OpLOADALL,
	_OcodeKindLowerName[680:689]:   OpLOADALL,
	_OcodeKindName[689:695]:        OpLOCK,
	_OcodeKindLowerName[689:695]:   OpLOCK,
	_OcodeKindName[695:703]:        OpMFENCE,
	_OcodeKindLowerName[695:703]:   OpMFENCE,
	_OcodeKindName[703:712]:        OpMONITOR,
	_OcodeKindLowerName[703:712]:   OpMONITOR,
	_OcodeKindName[712:717]:        OpMUL,
	_OcodeKindLowerName[712:717]:   OpMUL,
	_OcodeKindName[717:724]:        OpMWAIT,
	_OcodeKindLowerName[717:724]:   OpMWAIT,
	_OcodeKindName[724:729]:        OpNOP,
	_OcodeKindLowerName[724:729]:   OpNOP,
	_OcodeKindName[729:737]:        OpNTAKEN,
	_OcodeKindLowerName[729:737]:   OpNTAKEN,
	_OcodeKindName[737:744]:        OpPAUSE,
	_OcodeKindLowerName[737:744]:   OpPAUSE,
	_OcodeKindName[744:750]:        OpPOPA,
	_OcodeKindLowerName[744:750]:   OpPOPA,
	_OcodeKindName[750:757]:        OpPOPAD,
	_OcodeKindLowerName[750:757]:   OpPOPAD,
	_OcodeKindName[757:763]:        OpPOPF,
	_OcodeKindLowerName[757:763]:   OpPOPF,
	_OcodeKindName[763:770]:        OpPOPFD,
	_OcodeKindLowerName[763:770]:   OpPOPFD,
	_OcodeKindName[770:777]:        OpPOPFQ,
	_OcodeKindLowerName[770:777]:   OpPOPFQ,
	_OcodeKindName[777:784]:        OpPUSHA,
	_OcodeKindLowerName[777:784]:   OpPUSHA,
	_OcodeKindName[784:792]:        OpPUSHAD,
	_OcodeKindLowerName[784:792]:   OpPUSHAD,
	_OcodeKindName[792:799]:        OpPUSHF,
	_OcodeKindLowerName[792:799]:   OpPUSHF,
	_OcodeKindName[799:807]:        OpPUSHFD,
	_OcodeKindLowerName[799:807]:   OpPUSHFD,
	_OcodeKindName[807:815]:        OpPUSHFQ,
	_OcodeKindLowerName[807:815]:   OpPUSHFQ,
	_OcodeKindName[815:822]:        OpRDMSR,
	_OcodeKindLowerName[815:822]:   OpRDMSR,
	_OcodeKindName[822:829]:        OpRDPMC,
	_OcodeKindLowerName[822:829]:   OpRDPMC,
	_OcodeKindName[829:836]:        OpRDTSC,
	_OcodeKindLowerName[829:836]:   OpRDTSC,
	_OcodeKindName[836:844]:        OpRDTSCP,
	_OcodeKindLowerName[836:844]:   OpRDTSCP,
	_OcodeKindName[844:849]:        OpREP,
	_OcodeKindLowerName[844:849]:   OpREP,
	_OcodeKindName[849:855]:        OpREPE,
	_OcodeKindLowerName[849:855]:   OpREPE,
	_OcodeKindName[855:862]:        OpREPNE,
	_OcodeKindLowerName[855:862]:   OpREPNE,
	_OcodeKindName[862:868]:        OpRETF,
	_OcodeKindLowerName[862:868]:   OpRETF,
	_OcodeKindName[868:874]:        OpRETN,
	_OcodeKindLowerName[868:874]:   OpRETN,
	_OcodeKindName[874:879]:        OpRSM,
	_OcodeKindLowerName[874:879]:   OpRSM,
	_OcodeKindName[879:885]:        OpSAHF,
	_OcodeKindLowerName[879:885]:   OpSAHF,
	_OcodeKindName[885:893]:        OpSETALC,
	_OcodeKindLowerName[885:893]:   OpSETALC,
	_OcodeKindName[893:901]:        OpSFENCE,
	_OcodeKindLowerName[893:901]:   OpSFENCE,
	_OcodeKindName[901:905]:        OpSS,
	_OcodeKindLowerName[901:905]:   OpSS,
	_OcodeKindName[905:910]:        OpSTC,
	_OcodeKindLowerName[905:910]:   OpSTC,
	_OcodeKindName[910:915]:        OpSTD,
	_OcodeKindLowerName[910:915]:   OpSTD,
	_OcodeKindName[915:920]:        OpSTI,
	_OcodeKindLowerName[915:920]:   OpSTI,
	_OcodeKindName[920:928]:        OpSWAPGS,
	_OcodeKindLowerName[920:928]:   OpSWAPGS,
	_OcodeKindName[928:937]:        OpSYSCALL,
	_OcodeKindLowerName[928:937]:   OpSYSCALL,
	_OcodeKindName[937:947]:        OpSYSENTER,
	_OcodeKindLowerName[937:947]:   OpSYSENTER,
	_OcodeKindName[947:956]:        OpSYSEXIT,
	_OcodeKindLowerName[947:956]:   OpSYSEXIT,
	_OcodeKindName[956:964]:        OpSYSRET,
	_OcodeKindLowerName[956:964]:   OpSYSRET,
	_OcodeKindName[964:971]:        OpTAKEN,
	_OcodeKindLowerName[964:971]:   OpTAKEN,
	_OcodeKindName[971:976]:        OpUD2,
	_OcodeKindLowerName[971:976]:   OpUD2,
	_OcodeKindName[976:984]:        OpVMCALL,
	_OcodeKindLowerName[976:984]:   OpVMCALL,
	_OcodeKindName[984:994]:        OpVMLAUNCH,
	_OcodeKindLowerName[984:994]:   OpVMLAUNCH,
	_OcodeKindName[994:1004]:       OpVMRESUME,
	_OcodeKindLowerName[994:1004]:  OpVMRESUME,
	_OcodeKindName[1004:1012]:      OpVMXOFF,
	_OcodeKindLowerName[1004:1012]: OpVMXOFF,
	_OcodeKindName[1012:1018]:      OpWAIT,
	_OcodeKindLowerName[1012:1018]: OpWAIT,
	_OcodeKindName[1018:1026]:      OpWBINVD,
	_OcodeKindLowerName[1018:1026]: OpWBINVD,
	_OcodeKindName[1026:1033]:      OpWRMSR,
	_OcodeKindLowerName[1026:1033]: OpWRMSR,
	_OcodeKindName[1033:1041]:      OpXGETBV,
	_OcodeKindLowerName[1033:1041]: OpXGETBV,
	_OcodeKindName[1041:1049]:      OpXRSTOR,
	_OcodeKindLowerName[1041:1049]: OpXRSTOR,
	_OcodeKindName[1049:1057]:      OpXSETBV,
	_OcodeKindLowerName[1049:1057]: OpXSETBV,
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
	_OcodeKindName[611:616],
	_OcodeKindName[616:621],
	_OcodeKindName[621:627],
	_OcodeKindName[627:633],
	_OcodeKindName[633:639],
	_OcodeKindName[639:646],
	_OcodeKindName[646:653],
	_OcodeKindName[653:659],
	_OcodeKindName[659:665],
	_OcodeKindName[665:672],
	_OcodeKindName[672:680],
	_OcodeKindName[680:689],
	_OcodeKindName[689:695],
	_OcodeKindName[695:703],
	_OcodeKindName[703:712],
	_OcodeKindName[712:717],
	_OcodeKindName[717:724],
	_OcodeKindName[724:729],
	_OcodeKindName[729:737],
	_OcodeKindName[737:744],
	_OcodeKindName[744:750],
	_OcodeKindName[750:757],
	_OcodeKindName[757:763],
	_OcodeKindName[763:770],
	_OcodeKindName[770:777],
	_OcodeKindName[777:784],
	_OcodeKindName[784:792],
	_OcodeKindName[792:799],
	_OcodeKindName[799:807],
	_OcodeKindName[807:815],
	_OcodeKindName[815:822],
	_OcodeKindName[822:829],
	_OcodeKindName[829:836],
	_OcodeKindName[836:844],
	_OcodeKindName[844:849],
	_OcodeKindName[849:855],
	_OcodeKindName[855:862],
	_OcodeKindName[862:868],
	_OcodeKindName[868:874],
	_OcodeKindName[874:879],
	_OcodeKindName[879:885],
	_OcodeKindName[885:893],
	_OcodeKindName[893:901],
	_OcodeKindName[901:905],
	_OcodeKindName[905:910],
	_OcodeKindName[910:915],
	_OcodeKindName[915:920],
	_OcodeKindName[920:928],
	_OcodeKindName[928:937],
	_OcodeKindName[937:947],
	_OcodeKindName[947:956],
	_OcodeKindName[956:964],
	_OcodeKindName[964:971],
	_OcodeKindName[971:976],
	_OcodeKindName[976:984],
	_OcodeKindName[984:994],
	_OcodeKindName[994:1004],
	_OcodeKindName[1004:1012],
	_OcodeKindName[1012:1018],
	_OcodeKindName[1018:1026],
	_OcodeKindName[1026:1033],
	_OcodeKindName[1033:1041],
	_OcodeKindName[1041:1049],
	_OcodeKindName[1049:1057],
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
