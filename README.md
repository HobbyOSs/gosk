# GOSK ![](https://github.com/HobbyOSs/gosk/actions/workflows/go.yml/badge.svg)

GOSK is an assembler project written in Go, designed to parse an assembly language similar to NASK (a language used in the "30 days self-made OS" book) and generate machine code (currently supporting COFF format). It utilizes PEG (Parsing Expression Grammar) for parsing.

## Features

*   Parses NASK-like assembly code.
*   Generates COFF object files.
*   Includes a test suite for verification.

## Build & Run

Requires Go and Make installed.

```bash
# Clone the repository (if you haven't already)
# git clone https://github.com/HobbyOSs/gosk.git
# cd gosk

# Build the executable
make build

# Or install directly (might require setting up GOPATH/GOBIN)
# go install github.com/HobbyOSs/gosk@latest
```

## Usage

```bash
# Display help message
./gosk --help

# Assemble a source file (e.g., input.nas) into an object file (e.g., output.o)
./gosk input.nas output.o

# Assemble and generate a listing file (e.g., output.lst)
./gosk input.nas output.o output.lst
```

The command-line interface is:

```
usage: ./gosk [--help | -v] source [object/binary] [list]
  source: Input assembly file (.nas)
  object/binary: Output object file (e.g., .o, .obj) or raw binary file (Optional, defaults based on source name)
  list:   Output listing file (Optional)
  --help: Show this help message
  -v:     Show version and license information
```

## Contributing

Contributions are welcome! Please refer to the project's issue tracker and coding guidelines. (Details might be found in the `memory-bank` directory).

## License

This project is licensed under the [LICENSE](LICENSE) file.
