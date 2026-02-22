# CodeOwner: @r_owner

# This is a single-line comment

# R does not have multi-line block comments,
# so we use consecutive hash signs.

x <- 42  # Inline comment

greet <- function(name) {
  # Another single-line comment
  paste0("Hello, ", name, "!")
}

result <- greet("world")
print(result)
