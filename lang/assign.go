package lang

type Assign struct {
	Name  Ref
	Value Node
}

const (
	KeyAssignName  = "AssignName"
	KeyAssignValue = "AssignValue"
)

var _ Node = (*Assign)(nil)

func (a *Assign) Eval(env *Env) (*Object, error) {
	v, err := a.Value.Eval(env)
	if err != nil {
		return nil, err
	}
	env.Set(Ref(a.Name), v)
	return v, nil
}

func (a *Assign) Walk(v Visitor) error {
	skip, err := parseVisitorError(v.Push(KeyAssignName, a.Name))
	if err != nil {
		return err
	}
	if !skip {
		if err := a.Name.Walk(v); err != nil {
			return err
		}
	}
	if err := v.Pop(KeyAssignName, a.Name); err != nil {
		return err
	}
	skip, err = parseVisitorError(v.Push(KeyAssignValue, a.Value))
	if err != nil {
		return err
	}
	if !skip {
		if err := a.Value.Walk(v); err != nil {
			return err
		}
	}
	return v.Pop(KeyAssignValue, a.Value)
}
