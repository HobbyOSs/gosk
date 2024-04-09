#!/bin/bash

yq -r '.[] | select(.op1 == null and .op2 == null) | .mnem' x86reference.yml | sort
