module main

require (
    validators v0.0.0
    table v0.0.0
)

replace (
    validators => ./validators
    table => ./table
)

go 1.18