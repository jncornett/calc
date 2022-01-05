package syntax

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type (
	Program struct {
		Do []Expression `@@*`
	}

	Value struct {
		Nested     *Expression `  "(" @@ ")"`
		Assignment *Assignment `| @@`
		Unary      *OpUnary    `| @@`
		Call       *Call       `| @@`
		Symbol     string      `| @Ident`
		Number     string      `| @Number`
		Vector     *Vector     `| @@`
	}

	Assignment struct {
		Save  string     `@Ident ":="`
		Value Expression `@@`
	}

	OpUnary struct {
		Op    string     `@("-")`
		Value Expression `@@`
	}

	Expression struct {
		LHS *Term    `@@`
		RHS []OpTerm `@@*`
	}

	OpTerm struct {
		Op   string `@("+" | "-")`
		Term *Term  `@@`
	}

	Term struct {
		LHS *Factor    `@@`
		RHS []OpFactor `@@*`
	}

	OpFactor struct {
		Op     string  `@("*" | "/")`
		Factor *Factor `@@`
	}

	Factor struct {
		LHS *Value  `@@`
		RHS *OpUtil `@@?`
	}

	OpUtil struct {
		Op    string `@(".")`
		Value *Value `@@`
	}

	Call struct {
		Func string       `@Ident`
		Args []Expression `"(" @@ ("," @@)* ","? ")"`
	}

	Vector struct {
		Values []Expression `"[" @@ ("," @@)* ","? "]"`
	}
)

var Parser = participle.MustBuild(&Program{},
	participle.Lexer(Lexer),
	participle.Elide("Comment", "Whitespace", "EOL"),
	participle.UseLookahead(4))

var Lexer = lexer.MustSimple([]lexer.Rule{
	{Name: "Comment", Pattern: `//[^\n]*|/\*.*?\*/`},
	{Name: "OpenExpr", Pattern: `\(`},
	{Name: "CloseExpr", Pattern: `\)`},
	{Name: "OpenVector", Pattern: `\[`},
	{Name: "CloseVector", Pattern: `\]`},
	{Name: "Assign", Pattern: `:=`},
	{Name: "Comma", Pattern: `,`},
	{Name: "Dot", Pattern: `\.`},
	{Name: "Operator", Pattern: `[-+*/]`},
	{Name: "Number", Pattern: `(0|[1-9][0-9]*)(\.[0-9]+)?([eE][-+]?[0-9]+)?`},
	{Name: "Ident", Pattern: `[\p{L}_][\p{L}\p{N}_]*`},
	{Name: "EOL", Pattern: `[\n\r]+`},
	{Name: "Whitespace", Pattern: `[ \t]+`},
})
