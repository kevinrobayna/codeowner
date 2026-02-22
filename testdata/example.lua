-- CodeOwner: @lua_owner

-- This is a single-line comment

--[[
This is a multi-line
block comment in Lua.
--]]

local x = 42  -- Inline comment

function greet(name)
    -- Another single-line comment
    return "Hello, " .. name .. "!"
end

local result = greet("world")
print(result)
