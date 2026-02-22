-- CodeOwner: @haskell_owner

-- This is a single-line comment

{-
This is a multi-line
block comment in Haskell.
-}

x :: Int
x = 42  -- Inline comment

greet :: String -> String
greet name = "Hello, " ++ name ++ "!"  -- Inline comment

main :: IO ()
main = putStrLn (greet "world")
