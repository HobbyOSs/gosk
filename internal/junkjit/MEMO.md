# Design memo

- Dead copy of asmjit
- https://asmjit.com/doc/classasmjit_1_1x86_1_1Gp.html

```
junkjit/
├── assembler.go      # Assemblerインターフェイスと基本実装
├── codeholder.go     # CodeHolderインターフェイス
├── operand.go        # Operandインターフェイス
└── x86/              # x86アーキテクチャ向けの命令実装
    ├── assembler_impl.go
    ├── pseudo.go
    ├── jcc.go
    ├── no_param.go
    └── add.go
```
