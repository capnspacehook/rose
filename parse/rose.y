%{
package parse

import (
    "strconv"

    "github.com/capnspacehook/rose/ast"
    "github.com/capnspacehook/rose/token"
)
%}

%union {
    stmt ast.Statement
    stmtlist []ast.Statement
    expr ast.Expression
    exprlist []ast.Expression
    node ast.Node
    tok token.Token
}

// literals
%token <tok> IDENT INT FLOAT CHAR STRING RAW_STRING
// operators and punctuation
%token <tok> ASSIGN
// keywords
%token <tok> VAR

%type <stmtlist> statements
%type <stmt> statement varDecl
//%type <node> type
%type <expr> expression

%%

main: 
    statements
    {
        yylex.(*Lexer).Statements = $1
    }
;

statements:
    { $$ = nil }
|   statements statement
    {
        $$ = append($1, $2)
    }
;

statement:
    varDecl
;

varDecl:
    VAR IDENT ASSIGN expression
    {
        $$ = &ast.VarDeclStatement{
            Token: $1,
            Name:  &ast.Identifier{Token: $2},
            Value: $4,
        }
    }
;



expression:
    IDENT
    {
        if v, ok := boolConsts[$1.Literal]; ok {
            $$ = &ast.Boolean{
                Token: $1,
                Value: v,
            }
        } else if $1.Literal == "nil" {
            $$ = &ast.Nil{
                Token: $1,
            }
        } else {
            $$ = &ast.Identifier{
                Token: $1,
            }
        }
    }
|   INT
    {
        i, err := strconv.ParseInt($1.Literal, 0, 64)
        if err != nil {
            yylex.(*Lexer).Error(err.Error())
        }

        $$ = &ast.IntegerLiteral{
            Token: $1,
            Value: i,
        }
    }
|   FLOAT
    {
        f, err := strconv.ParseFloat($1.Literal, 64)
        if err != nil {
            yylex.(*Lexer).Error(err.Error())
        }

        $$ = &ast.FloatLiteral{
            Token: $1,
            Value: f,
        }
    }
;
