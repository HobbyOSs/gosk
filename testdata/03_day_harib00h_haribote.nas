; haribote-os
; TAB=4

; BOOT_INFO髢「菫
CYLS	EQU		0x0ff0			; 繝悶繝医そ繧ッ繧ソ縺瑚ィュ螳壹☆繧
LEDS	EQU		0x0ff1
VMODE	EQU		0x0ff2			; 濶イ謨ー縺ォ髢「縺吶ｋ諠⒦ア縲ゆス輔ン繝ヨ繧ォ繝ゥ繝シ縺具シ
SCRNX	EQU		0x0ff4			; 隗」蜒丞コヲ縺ョX
SCRNY	EQU		0x0ff6			; 隗」蜒丞コヲ縺ョY
VRAM	EQU		0x0ff8			; 繧ー繝ゥ繝輔ぅ繝け繝舌ャ繝輔ぃ縺ョ髢句ァ狗分蝨ー

		ORG		0xc200			; 縺薙繝励Ο繧ー繝ゥ繝縺後←縺薙↓隱ュ縺ソ霎シ縺セ繧後ｋ縺ョ縺

		MOV		AL,0x13			; VGA繧ー繝ゥ繝輔ぅ繝け繧ケ縲20x200x8bit繧ォ繝ゥ繝シ
		MOV		AH,0x00
		INT		0x10
		MOV		BYTE [VMODE],8	; 逕サ髱「繝「繝シ繝峨ｒ繝。繝「縺吶ｋ
		MOV		WORD [SCRNX],320
		MOV		WORD [SCRNY],200
		MOV		DWORD [VRAM],0x000a0000

; 繧ュ繝シ繝懊繝峨LED迥カ諷九ｒBIOS縺ォ謨吶∴縺ヲ繧ゅｉ縺

		MOV		AH,0x02
		INT		0x16 			; keyboard BIOS
		MOV		[LEDS],AL

fin:
		HLT
		JMP		fin
