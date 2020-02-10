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
%token <tok> AND OR XOR INVT SHL SHR AND_NOT
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

// precedence
%left MUL QUO REM SHL SHR AND AND_NOT
%right EXP
%left ADD SUB OR XOR
%left EQL NEQ LSS LEQ GTR GEQ
%left LAND
%left LOR 

%type <stmtlist> statements
%type <stmt> statement declaration varDecl constDecl simpleStatement expressionStatement assignment
%type <typename> type
%type <expr> expression unary_expression primary_expression operand conversion basic_lit
%type <tok> binary_op rel_op add_op mul_op unary_op

%%

main: 
    statements
    {
        yylex.(*lexer).Program = &ast.Program{
            Statements: $1,
        }
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
    declaration
|   simpleStatement
;

declaration:
    varDecl
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

simpleStatement:
    expressionStatement
|   assignment
;

expressionStatement:
    expression
    {
        $$ = &ast.ExprStatement{
            Expr: $1,
        }
    }

assignment:
    IDENT ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT ADD_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT SUB_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT MUL_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT QUO_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT REM_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT EXP_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT AND_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT OR_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT XOR_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT SHL_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT SHR_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
|   IDENT AND_NOT_ASSIGN expression
    {
        $$ = &ast.AssignmentStatement{
            Token: $2,
            Name:  &ast.Identifier{Token: $1},
            Value: $3,
        }
    }
;

expression:
    unary_expression
|   expression binary_op expression
    {
        $$ = &ast.BinaryExpression{
            Lhs:   $1,
            Token: $2,
            Rhs:   $3,
        }
    }
;

unary_expression:
    primary_expression
|   unary_op unary_expression
    {
        $$ = &ast.UnaryExpression{
            Token: $1,
            Value: $2,
        }
    }
;

primary_expression:
    operand
|   conversion
;

operand:
    IDENT
    {
        if v, ok := boolConsts[$1.Literal]; ok {
            $$ = &ast.BooleanLiteral{
                Token: $1,
                Value: v,
            }
        } else if $1.Literal == "nil" {
            $$ = &ast.NilLiteral{
                Token: $1,
            }
        } else {
            $$ = &ast.Identifier{
                Token: $1,
            }
        }
    }
|   basic_lit
|   LPAREN expression RPAREN
    { 
        $$ = &ast.ParenExpression{
            Lparen: $1,
            Expr:   $2,
            Rparen: $3,
        }
    }
;

conversion:
    type LPAREN expression RPAREN
    {
        $$ = &ast.Conversion{
            Type:  $1,
            Value: $3,
        }
    }
;

basic_lit:
    INT
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

binary_op:
    LAND
|   LOR
|   rel_op
|   add_op
|   mul_op
;

rel_op:
    EQL
|   NEQ
|   LSS
|   LEQ
|   GTR
|   GEQ
;

add_op:
    ADD
|   SUB
|   OR
|   XOR
;

mul_op:
    MUL
|   QUO
|   REM
|   EXP
|   SHL
|   SHR
|   AND
|   AND_NOT
;

unary_op:
    ADD
|   SUB
|   NOT
|   INVT
|   ARROW
;
