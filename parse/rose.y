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
%token <tok> ADD SUB MUL QUO REM EXP
%token <tok> AND OR XOR SHL SHR AND_NOT
%token <tok> ADD_ASSIGN SUB_ASSIGN MUL_ASSIGN QUO_ASSIGN REM_ASSIGN EXP_ASSIGN
%token <tok> AND_ASSIGN OR_ASSIGN XOR_ASSIGN SHL_ASSIGN SHR_ASSIGN AND_NOT_ASSIGN
%token <tok> LAND LOR ARROW INC DEC
%token <tok> EQL LSS GTR ASSIGN NOT
%token <tok> NEQ LEQ GEQ ELLIPSIS
%token <tok> LPAREN LBRACK LBRACE COMMA PERIOD 
%token <tok> RPAREN RBRACK RBRACE SEMICOLON COLON
%token <tok> QUES EXCLM
// keywords
%token <tok> CONST LET VAR

%type <stmtlist> statements
%type <stmt> statement assignment varDecl constDecl
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
|   statements statement SEMICOLON
    {
        $$ = append($1, $2)
    }
;

statement:
    assignment
|   varDecl
|   constDecl
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

constDecl:
    CONST IDENT ASSIGN expression
    {
        $$ = &ast.ConstDeclStatement{
            Token: $1,
            Name:  &ast.Identifier{Token: $2},
            Value: $4,
        }
    }
|   CONST IDENT type ASSIGN expression
    {
        $$ = &ast.ConstDeclStatement{
            Token: $1,
            Name:  &ast.Identifier{Token: $2},
            Type:  $3,
            Value: $5,
        }
    }
|   LET IDENT ASSIGN expression
    {
        $$ = &ast.ConstDeclStatement{
            Token: $1,
            Name:  &ast.Identifier{Token: $2},
            Value: $4,
        }
    }
|   LET IDENT type ASSIGN expression
    {
        $$ = &ast.ConstDeclStatement{
            Token: $1,
            Name:  &ast.Identifier{Token: $2},
            Type:  $3,
            Value: $5,
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
|   CHAR
    {
        if n := len($1.Literal); n >= 2 {
			c, _, _, err := strconv.UnquoteChar($1.Literal[1:n-1], '\'')
            if err != nil {
				yylex.(*lexer).Error(err.Error())
			} else {
                $$ = &ast.CharLiteral{
                    Token: $1,
                    Value: c,
                }
            }
		}
    }
|   STRING
    {
        $$ = &ast.StringLiteral{
            Token: $1,
        }
    }
|   RAW_STRING
    {
        $$ = &ast.RawStringLiteral{
            Token: $1,
        }
    }
;
