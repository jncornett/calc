package lang

type Call struct {
	Func Node
	Args []Node
}

var _ Node = (*Call)(nil)

func (c *Call) Eval(env *Env) (*Object, error) {
	lhs, err := c.Func.Eval(env)
	if err != nil {
		return nil, err
	}
	fn, ok := lhs.GoValue.(NativeFunc)
	if !ok {
		return nil, NewError(InvalidType, "object of type %T does not support function calls", lhs.GoValue)
	}
	var args []*Object
	for _, n := range c.Args {
		v, err := n.Eval(env)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	got, err := fn(args...)
	if err != nil {
		return nil, err
	}
	return got, nil
}

const KeyCallName = "CallName"

type KeyCallArg int

func (c *Call) Walk(v Visitor) error {
	skip, err := parseVisitorError(v.Push(KeyCallName, c.Func))
	if err != nil {
		return err
	}
	if !skip {
		if err := c.Func.Walk(v); err != nil {
			return err
		}
	}
	if err := v.Pop(KeyCallName, c.Func); err != nil {
		return err
	}
	for i, n := range c.Args {
		skip, err := parseVisitorError(v.Push(KeyCallArg(i), n))
		if err != nil {
			return err
		}
		if !skip {
			if err := n.Walk(v); err != nil {
				return err
			}
		}
		if err := v.Pop(KeyCallArg(i), n); err != nil {
			return err
		}
	}
	return nil
}
