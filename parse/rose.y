%{
package parse

import (
    "fmt"
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
    typename *ast.TypeName
    tok token.Token
}

// identifiers and literals
%token <tok> IDENT INT FLOAT CHAR STRING RAW_STRING
// operators and punctuation
%token <tok> ADD SUB MUL QUO REM EXP ASSIGN
%token <tok> LPAREN LBRACK LBRACE COMMA PERIOD 
%token <tok> RPAREN RBRACK RBRACE SEMICOLON COLON
// keywords
%token <tok> VAR

%type <stmtlist> statements
%type <stmt> statement assignment varDecl
%type <typename> type
%type <expr> expression

%%

main: 
    statements
    {
        yylex.(*lexer).Statements = $1
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
    assignment
|   varDecl 
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
|   VAR IDENT type ASSIGN expression
    {
        $$ = &ast.VarDeclStatement{
            Token: $1,
            Name:  &ast.Identifier{Token: $2},
            Type:  $3,
            Value: $5,
        }
    }
;

type:
    IDENT
    {
        if _, ok := typeNames[$1.Literal]; ok {
            $$ = &ast.TypeName{
                Token: $1,
            }
        } else {
            yylex.Error(fmt.Sprintf("%q is not a valid type", $1.Literal))
        }
    }
;

assignment:
    IDENT ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
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
            yylex.(*lexer).Error(err.Error())
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
            yylex.(*lexer).Error(err.Error())
        }

        $$ = &ast.FloatLiteral{
            Token: $1,
            Value: f,
        }
    }
;
