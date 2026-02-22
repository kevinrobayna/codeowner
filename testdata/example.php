<?php
// CodeOwner: @php_owner

// This is a single-line comment

# This is also a single-line comment (hash style)

/*
 * This is a multi-line
 * block comment in PHP.
 */

/** PHPDoc comment. */
$x = 42;  // Inline comment

function greet(string $name): string {
    // Another single-line comment
    return "Hello, {$name}!";
}

echo greet("world");  # Inline hash comment
echo "\n";  /* Inline block comment */
