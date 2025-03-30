#!/bin/bash

NASK_PATH="/home4/wine/drive_c/MinGW/msys/1.0/bin/nask.exe"

show_help() {
    cat <<EOF
Usage: echo 'assembly code' | $0
       echo 'assembly code' | $0 -l

Assemble x86 code via nask (NASM-like assembler via Wine) and show hexdump.

Options:
  -l       Also show nask list file (.lst) output
  -h       Show this help

Example:
  echo "mov al, 0x42\nhlt" | $0
  echo "mov al, 0x90\nnop\nhlt" | $0 -l
EOF
}

# Check help flag or no input
if [[ "$1" == "-h" || -z "$1" && -t 0 ]]; then
    show_help
    exit 0
fi

# Option parsing
SHOW_LIST=0
if [[ "$1" == "-l" ]]; then
    SHOW_LIST=1
    shift
fi

# Prepare temp files
tmpdir=$(mktemp -d /tmp/naskXXXX)
input="$tmpdir/input.asm"
output="$tmpdir/output.bin"
list="$tmpdir/output.lst"

# Read stdin to input.asm
cat - > "$input"

# Run nask via wine
wine "$NASK_PATH" "$input" "$output" "$list" > /dev/null 2>&1

# Check output
if [ -f "$output" ]; then
    echo "Hexdump of binary output:"
    hexdump -C "$output"
    if [[ $SHOW_LIST -eq 1 && -f "$list" ]]; then
        echo -e "\nAssembler list output:"
        cat "$list"
    fi
else
    echo "nask failed or did not produce output."
fi

# Clean up
rm -rf "$tmpdir"

