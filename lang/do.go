package lang

type Do []Node

var _ Node = Do(nil)

func (do Do) Eval(env *Env) (o *Object, err error) {
	o = Nil()
	for _, n := range do {
		o, err = n.Eval(env)
		if err != nil {
			return nil, err
		}
	}
	return o, nil
}

type KeyDoItem int

func (do Do) Walk(v Visitor) error {
	for i, n := range do {
		skip, err := parseVisitorError(v.Push(KeyDoItem(i), n))
		if err != nil {
			return err
		}
		if !skip {
			if err := n.Walk(v); err != nil {
				return err
			}
		}
		if err := v.Pop(KeyDoItem(i), n); err != nil {
			return err
		}
	}
	return nil
}
