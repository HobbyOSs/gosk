; haribote-os
; TAB=4

		ORG		0xc200			; 縺薙繝励Ο繧ー繝ゥ繝縺後←縺薙↓隱ュ縺ソ霎シ縺セ繧後ｋ縺ョ縺

		MOV		AL,0x13			; VGA繧ー繝ゥ繝輔ぅ繝け繧ケ縲20x200x8bit繧ォ繝ゥ繝シ
		MOV		AH,0x00
		INT		0x10
fin:
		HLT
		JMP		fin
