#!/bin/bash

# CodeOwner: @shell_owner

# This is a single-line comment

: '
This is a multi-line
block comment in Bash.
'

X=42  # Inline comment

greet() {
    # Another single-line comment
    echo "Hello, $1!"
}

greet "world"
