% CodeOwner: @erlang_owner

% This is a single-line comment

%% This is a conventional single-line comment in Erlang.

-module(example).
-export([greet/1]).

X = 42,  % Inline comment

greet(Name) ->
    %% Another single-line comment
    "Hello, " ++ Name ++ "!".

main() ->
    io:format("~s~n", [greet("world")]).
