# CodeOwner: @elixir_owner

# This is a single-line comment

@moduledoc """
This is a multi-line
doc comment in Elixir.
"""

defmodule Example do
  @doc "Greets the given name."
  def greet(name) do
    # Another single-line comment
    "Hello, #{name}!"  # Inline comment
  end
end

IO.puts(Example.greet("world"))
