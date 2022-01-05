package lang

type Ref string

var _ Node = Ref("")

func (ref Ref) Eval(env *Env) (*Object, error) {
	o, ok := env.Get(ref)
	if !ok {
		return nil, NewError(NotFound, "name not found: %q", string(ref))
	}
	return o, nil
}

func (ref Ref) Walk(v Visitor) error { return nil }
