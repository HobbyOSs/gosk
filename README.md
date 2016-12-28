# Gosk [![Build Status](https://travis-ci.org/hangingman/gosk.svg?branch=master)](https://travis-ci.org/hangingman/gosk)

This is a yet another assembly interpreter gosk!

## Build & Run

* You need to install Go and Make

```
$ go get -u github.com/hangingman/gosk
$ cd $GO_HOME/src/github.com/hangingman/gosk
$ make
```

## How to run gosk

* REPL mode
  * You can run REPL mode with no option

```
$ ./gosk
Hello user! This is yet another assembly gosk!
Feel free to type in commands
>> DB 0x00
[  info ] parser_parse_stmt.go:49: { OPCODE:{ DB: DB,0x00 } }
[  info ] eval.go:136: [OPCODE: DB, HEX_LIT: 0x00]
00
>> RESB 10
[  info ] parser_parse_stmt.go:49: { OPCODE:{ RESB: RESB,10 } }
[  info ] eval.go:228: [OPCODE: RESB, INT: 10]
00000000000000000000
>>
```

* Normal assembly mode
  * You can generate an object file from an assembly source (*.nas format)

```
./gosk --help
usage:  [--help | -v] source [object/binary] [list]
  -v	バージョンとライセンス情報を表示する
```
