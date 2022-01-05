package lang

type Vector []Node

var _ Node = Vector(nil)

func (vec Vector) Eval(env *Env) (*Object, error) {
	var args []interface{}
	for _, n := range vec {
		v, err := n.Eval(env)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	return &Object{GoValue: args}, nil
}

type KeyVectorItem int

func (vec Vector) Walk(v Visitor) error {
	for i, n := range vec {
		skip, err := parseVisitorError(v.Push(KeyVectorItem(i), n))
		if err != nil {
			return err
		}
		if !skip {
			if err := n.Walk(v); err != nil {
				return err
			}
		}
		if err := v.Pop(KeyVectorItem(i), n); err != nil {
			return err
		}
	}
	return nil
}
