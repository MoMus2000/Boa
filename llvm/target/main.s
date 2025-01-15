	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 14, 0
	.globl	_main                           ; -- Begin function main
	.p2align	2
_main:                                  ; @main
	.cfi_startproc
; %bb.0:                                ; %main_entry
	mov	w0, #420                        ; =0x1a4
	ret
	.cfi_endproc
                                        ; -- End function
.subsections_via_symbols
