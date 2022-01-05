package calc

import (
	"strconv"

	"github.com/jncornett/calc/lang"
	"github.com/jncornett/calc/syntax"
)

func Eval(source []byte, env *lang.Env) (interface{}, error) {
	node, err := Parse(source)
	if err != nil {
		return nil, err
	}
	o, err := node.Eval(env)
	if err != nil {
		return nil, err
	}
	return o.GoValue, nil
}

func DefaultEnv() *lang.Env {
	return &lang.Env{SymbolTable: lang.Chained{builtins, lang.Map{}}}
}

func Parse(source []byte) (lang.Node, error) {
	prog, err := ParseSyntax(source)
	if err != nil {
		return nil, err
	}
	return ParseNode(prog), nil
}

func ParseSyntax(source []byte) (*syntax.Program, error) {
	var out syntax.Program
	err := syntax.Parser.ParseBytes("", source, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func ParseNode(p *syntax.Program) lang.Node {
	return parseProgram(p)
}

func parseProgram(p *syntax.Program) lang.Do {
	var out lang.Do
	for _, expr := range p.Do {
		got := parseExpression(&expr)
		out = append(out, got)
	}
	return out
}

func parseAssignment(a *syntax.Assignment) lang.Node {
	return &lang.Assign{
		Name:  lang.Ref(a.Save),
		Value: parseExpression(&a.Value),
	}
}

func parseExpression(e *syntax.Expression) lang.Node {
	lhs := parseTerm(e.LHS)
	for _, opTerm := range e.RHS {
		rhs := parseTerm(opTerm.Term)
		lhs = &lang.Call{
			Func: lang.Ref(opTerm.Op),
			Args: []lang.Node{lhs, rhs},
		}
	}
	return lhs
}

func parseTerm(t *syntax.Term) lang.Node {
	lhs := parseFactor(t.LHS)
	for _, opFactor := range t.RHS {
		rhs := parseFactor(opFactor.Factor)
		lhs = &lang.Call{
			Func: lang.Ref(opFactor.Op),
			Args: []lang.Node{lhs, rhs},
		}
	}
	return lhs
}

func parseFactor(f *syntax.Factor) lang.Node {
	lhs := parseValue(f.LHS)
	if f.RHS != nil {
		rhs := parseValue(f.RHS.Value)
		lhs = &lang.Call{
			Func: lang.Ref(f.RHS.Op),
			Args: []lang.Node{lhs, rhs},
		}
	}
	return lhs
}

func parseValue(v *syntax.Value) lang.Node {
	switch {
	case v.Nested != nil:
		return parseExpression(v.Nested)
	case v.Assignment != nil:
		return parseAssignment(v.Assignment)
	case v.Call != nil:
		return parseCall(v.Call)
	case v.Symbol != "":
		return parseSymbol(v.Symbol)
	case v.Number != "":
		return parseNumber(v.Number)
	case v.Vector != nil:
		return parseVector(v.Vector)
	default:
		panic("invariant")
	}
}

func parseCall(c *syntax.Call) *lang.Call {
	fn := lang.Ref(c.Func)
	var args []lang.Node
	for _, expr := range c.Args {
		args = append(args, parseExpression(&expr))
	}
	return &lang.Call{
		Func: fn,
		Args: args,
	}
}

func parseSymbol(s string) lang.Ref { return lang.Ref(s) }

func parseNumber(n string) *lang.Object {
	f, err := strconv.ParseFloat(n, 64)
	if err != nil {
		panic(err)
	}
	return &lang.Object{GoValue: f}
}

func parseVector(v *syntax.Vector) lang.Vector {
	var out lang.Vector
	for _, expr := range v.Values {
		out = append(out, parseExpression(&expr))
	}
	return out
}
