// CodeOwner: @rust_owner

// This is a single-line comment

/*
 * This is a multi-line
 * block comment in Rust.
 */

/// Doc comment for the function.
fn greet(name: &str) -> String {
    // Another single-line comment
    format!("Hello, {}!", name)
}

fn main() {
    let x = 42;  // Inline comment
    let result = greet("world");  /* Inline block comment */
    println!("{}", result);
}
