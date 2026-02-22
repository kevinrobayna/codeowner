// CodeOwner: @ts_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in TypeScript.
 */

const x: number = 42;  // Inline comment

function greet(name: string): string {
    // Another single-line comment
    return `Hello, ${name}!`;
}

const result: string = greet("world");  /* Inline block comment */
console.log(result);
