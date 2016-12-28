; haribote-os boot asm
; TAB=4

[INSTRSET "i486p"]

VBEMODE	EQU		0x105			; 1024 x  768 x 8bit繧ォ繝ゥ繝シ
; 育判髱「繝「繝シ繝我クヲァ;	0x100 :  640 x  400 x 8bit繧ォ繝ゥ繝シ
;	0x101 :  640 x  480 x 8bit繧ォ繝ゥ繝シ
;	0x103 :  800 x  600 x 8bit繧ォ繝ゥ繝シ
;	0x105 : 1024 x  768 x 8bit繧ォ繝ゥ繝シ
;	0x107 : 1280 x 1024 x 8bit繧ォ繝ゥ繝シ

BOTPAK	EQU		0x00280000		; bootpack縺ョ繝ュ繝シ繝牙
DSKCAC	EQU		0x00100000		; 繝ぅ繧ケ繧ッ繧ュ繝」繝す繝・縺ョ蝣エ謇DSKCAC0	EQU		0x00008000		; 繝ぅ繧ケ繧ッ繧ュ繝」繝す繝・縺ョ蝣エ謇シ医Μ繧「繝ォ繝「繝シ繝会シ
; BOOT_INFO髢「菫CYLS	EQU		0x0ff0			; 繝悶繝医そ繧ッ繧ソ縺瑚ィュ螳壹☆繧LEDS	EQU		0x0ff1
VMODE	EQU		0x0ff2			; 濶イ謨ー縺ォ髢「縺吶ｋ諠⒦ア縲ゆス輔ン繝ヨ繧ォ繝ゥ繝シ縺具シSCRNX	EQU		0x0ff4			; 隗」蜒丞コヲ縺ョX
SCRNY	EQU		0x0ff6			; 隗」蜒丞コヲ縺ョY
VRAM	EQU		0x0ff8			; 繧ー繝ゥ繝輔ぅ繝け繝舌ャ繝輔ぃ縺ョ髢句ァ狗分蝨ー

		ORG		0xc200			; 縺薙繝励Ο繧ー繝ゥ繝縺後←縺薙↓隱ュ縺ソ霎シ縺セ繧後ｋ縺ョ縺
; VBE蟄伜惠遒コ隱
		MOV		AX,0x9000
		MOV		ES,AX
		MOV		DI,0
		MOV		AX,0x4f00
		INT		0x10
		CMP		AX,0x004f
		JNE		scrn320

; VBE縺ョ繝舌繧ク繝ァ繝ウ繝√ぉ繝け

		MOV		AX,[ES:DI+4]
		CMP		AX,0x0200
		JB		scrn320			; if (AX < 0x0200) goto scrn320

; 逕サ髱「繝「繝シ繝画ュ蝣ア繧貞セ励ｋ

		MOV		CX,VBEMODE
		MOV		AX,0x4f01
		INT		0x10
		CMP		AX,0x004f
		JNE		scrn320

; 逕サ髱「繝「繝シ繝画ュ蝣ア縺ョ遒コ隱
		CMP		BYTE [ES:DI+0x19],8
		JNE		scrn320
		CMP		BYTE [ES:DI+0x1b],4
		JNE		scrn320
		MOV		AX,[ES:DI+0x00]
		AND		AX,0x0080
		JZ		scrn320			; 繝「繝シ繝牙ア樊縺ョbit7縺縺縺」縺溘縺ァ縺ゅ″繧峨ａ繧
; 逕サ髱「繝「繝シ繝峨蛻ｊ譖ソ縺
		MOV		BX,VBEMODE+0x4000
		MOV		AX,0x4f02
		INT		0x10
		MOV		BYTE [VMODE],8	; 逕サ髱「繝「繝シ繝峨ｒ繝。繝「縺吶ｋ險ェ槭′蜿らⅨ縺吶ｋ		MOV		AX,[ES:DI+0x12]
		MOV		[SCRNX],AX
		MOV		AX,[ES:DI+0x14]
		MOV		[SCRNY],AX
		MOV		EAX,[ES:DI+0x28]
		MOV		[VRAM],EAX
		JMP		keystatus

scrn320:
		MOV		AL,0x13			; VGA繧ー繝ゥ繝輔ぅ繝け繧ケ縲20x200x8bit繧ォ繝ゥ繝シ
		MOV		AH,0x00
		INT		0x10
		MOV		BYTE [VMODE],8	; 逕サ髱「繝「繝シ繝峨ｒ繝。繝「縺吶ｋ險ェ槭′蜿らⅨ縺吶ｋ		MOV		WORD [SCRNX],320
		MOV		WORD [SCRNY],200
		MOV		DWORD [VRAM],0x000a0000

; 繧ュ繝シ繝懊繝峨LED迥カ諷九ｒBIOS縺ォ謨吶∴縺ヲ繧ゅｉ縺
keystatus:
		MOV		AH,0x02
		INT		0x16 			; keyboard BIOS
		MOV		[LEDS],AL

; PIC縺御ク縺ョ蜑イ繧願セシ縺ソ繧貞女縺台サ倥¢縺ェ縺ｈ縺↓縺吶ｋ
;	AT莠呈鋤讖溘莉墓ァ倥〒縺ッ縲‾IC縺ョ蛻晄悄蛹悶ｒ縺吶ｋ縺ェ繧峨
;	縺薙＞縺、繧辰LI蜑阪↓繧▲縺ヲ縺翫°縺ェ縺→縲√◆縺セ縺ォ繝上Φ繧ー繧「繝縺吶ｋ
;	PIC縺ョ蛻晄悄蛹悶縺ゅ→縺ァ繧ｋ

		MOV		AL,0xff
		OUT		0x21,AL
		NOP						; OUT蜻ス莉、繧帝邯壹＆縺帙ｋ縺ィ縺∪縺上＞縺九↑縺ゥ溽ィョ縺後≠繧九ｉ縺励＞縺ョ縺ァ
		OUT		0xa1,AL

		CLI						; 縺輔ｉ縺ォCPU繝ャ繝吶Ν縺ァ繧ょ牡繧願セシ縺ソ遖∵ュ「

; CPU縺九ｉ1MB莉・荳翫繝。繝「繝ェ縺ォ繧「繧ッ繧サ繧ケ縺ァ縺阪ｋ繧医≧縺ォ縲、20GATE繧定ィュ螳
		CALL	waitkbdout
		MOV		AL,0xd1
		OUT		0x64,AL
		CALL	waitkbdout
		MOV		AL,0xdf			; enable A20
		OUT		0x60,AL
		CALL	waitkbdout

; 繝励Ο繝け繝医Δ繝シ繝臥ァサ陦
		LGDT	[GDTR0]			; 證ォ螳哦DT繧定ィュ螳		MOV		EAX,CR0
		AND		EAX,0x7fffffff	; bit31繧縺ォ縺吶ｋ医繝シ繧ク繝ウ繧ー遖∵ュ「縺ョ縺溘ａ		OR		EAX,0x00000001	; bit0繧縺ォ縺吶ｋ医繝ュ繝け繝医Δ繝シ繝臥ァサ陦後縺溘ａ		MOV		CR0,EAX
		JMP		pipelineflush
pipelineflush:
		MOV		AX,1*8			;  隱ュ縺ソ譖ク縺榊庄閭ス繧サ繧ー繝。繝ウ繝2bit
		MOV		DS,AX
		MOV		ES,AX
		MOV		FS,AX
		MOV		GS,AX
		MOV		SS,AX

; bootpack縺ョ霆「騾
		MOV		ESI,bootpack	; 霆「騾∝❼
		MOV		EDI,BOTPAK		; 霆「騾∝
		MOV		ECX,512*1024/4
		CALL	memcpy

; 縺、縺〒縺ォ繝ぅ繧ケ繧ッ繝繧ソ繧よ悽譚・縺ョ菴咲スョ縺ク霆「騾
; 縺セ縺壹繝悶繝医そ繧ッ繧ソ縺九ｉ

		MOV		ESI,0x7c00		; 霆「騾∝❼
		MOV		EDI,DSKCAC		; 霆「騾∝
		MOV		ECX,512/4
		CALL	memcpy

; 谿九ｊ蜈ィ驛ィ

		MOV		ESI,DSKCAC0+512	; 霆「騾∝❼
		MOV		EDI,DSKCAC+512	; 霆「騾∝
		MOV		ECX,0
		MOV		CL,BYTE [CYLS]
		IMUL	ECX,512*18*2/4	; 繧キ繝ェ繝ウ繝焚縺九ｉ繝舌う繝域焚/4縺ォ螟画鋤
		SUB		ECX,512/4		; IPL縺ョ蛻□縺大キョ縺怜シ輔¥
		CALL	memcpy

; asmhead縺ァ縺励↑縺代ｌ縺ー縺¢縺ェ縺％縺ィ縺ッ蜈ィ驛ィ縺礼オゅｏ縺」縺溘縺ァ縲;	縺ゅ→縺ッbootpack縺ォ莉サ縺帙ｋ

; bootpack縺ョ襍キ蜍
		MOV		EBX,BOTPAK
		MOV		ECX,[EBX+16]
		ADD		ECX,3			; ECX += 3;
		SHR		ECX,2			; ECX /= 4;
		JZ		skip			; 霆「騾√☆繧九∋縺阪ｂ縺ョ縺後↑縺		MOV		ESI,[EBX+20]	; 霆「騾∝❼
		ADD		ESI,EBX
		MOV		EDI,[EBX+12]	; 霆「騾∝
		CALL	memcpy
skip:
		MOV		ESP,[EBX+12]	; 繧ケ繧ソ繝け蛻晄悄蛟、
		JMP		DWORD 2*8:0x0000001b

waitkbdout:
		IN		 AL,0x64
		AND		 AL,0x02
		JNZ		waitkbdout		; AND縺ョ邨先棡縺縺ァ縺ェ縺代ｌ縺ーwaitkbdout縺ク
		RET

memcpy:
		MOV		EAX,[ESI]
		ADD		ESI,4
		MOV		[EDI],EAX
		ADD		EDI,4
		SUB		ECX,1
		JNZ		memcpy			; 蠑輔″邂励＠縺溽オ先棡縺縺ァ縺ェ縺代ｌ縺ーmemcpy縺ク
		RET
; memcpy縺ッ繧「繝峨Ξ繧ケ繧オ繧、繧コ繝励Μ繝輔ぅ繧ッ繧ケ繧貞Ⅶ繧悟ソ倥ｌ縺ェ縺代ｌ縺ー縲√せ繝医Μ繝ウ繧ー蜻ス莉、縺ァ繧よ嶌縺代ｋ

		ALIGNB	16
GDT0:
		RESB	8				; 繝後Ν繧サ繝ャ繧ッ繧ソ
		DW		0xffff,0x0000,0x9200,0x00cf	; 隱ュ縺ソ譖ク縺榊庄閭ス繧サ繧ー繝。繝ウ繝2bit
		DW		0xffff,0x0000,0x9a28,0x0047	; 螳溯。悟庄閭ス繧サ繧ー繝。繝ウ繝2bitootpack逕ィ
		DW		0
GDTR0:
		DW		8*3-1
		DD		GDT0

		ALIGNB	16
bootpack: