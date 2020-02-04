# The Rose Programming Language Specification

*Alpha Version*

## Introduction
This is a reference manual for the Rose programming language. Rose is a general-purpose language designed to allow users to quickly produce safe, simple code that isn't hard to understand. It is strongly typed and garbage collected, and has explicit support for concurrent programming. 

## Notation
The syntax is specified using Extended Backus-Naur Form (EBNF).

## Source code representation
Source code is Unicode text encoded in UTF-8. Each code point is distinct; for example upper and lower case letters are different characters.

### Letters and digits
The underscore character `_` (U+005F) is considered a letter.
```
letter        = "A" ... "Z" | "a" ... "z" | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
```

## Lexical elements
### Comments
Comments serve as program documentation. There are two forms:
1.  *Line comments* start with the character sequence `//` and stop at the end of the line.
2.  *General comments* start with the character sequence `/*` and stop with the first subsequent character sequence `*/`.
A comment cannot start inside a character or string literal, or inside a comment. A general comment containing no newlines acts like a space. Any other comment acts like a newline.

### Tokens TODO

### Semicolons
The formal grammar uses semicolons `;` as terminatiors in a number of productions. Rose programs may omit most of these semicolons using the following two rules: 

1. When the input is broken into tokens, a semicolon is automatically inserted into the token stream immediately after a line's final token if that token is
-   an identifier
-   an integer, floating-point, character, or string literal
-   one of the keywords] `break`, `continue`, `fallthrough`, or `return` **TODO: add other keywords**
-   one of the operators and punctuation `++`, `--`, `)`, `]`, or `}` **TODO: add other operators**
2. To allow complex statements to occupy a single line, a semicolon may be omitted before a closing `")"` or `"}"`.

To reflect idiomatic use, code examples in this document elide semicolons using these rules.

### Identifiers
Identifiers name program entities such as variables and types. An identifier is a sequence of one or more letters and digits. The first character in an identifier must be a letter.
```
identifier = letter { letter | digit } .
```
```
a
_x86
HiImAVariable
αβ
```

### Keywords
The following keywords are reserved and may not be used as identifiers.
```
and
or
not
const
let
fn
return
defer
is
as
in
if
for
import
break
continue
switch
select
go
chan
del
guard
```

### Operators and punctuation
