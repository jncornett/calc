package lang

var nilValue = new(Object)

type Object struct {
	Type    []string
	GoValue interface{}
}

var _ Node = (*Object)(nil)

func Nil() *Object { return nilValue }

func ObjectOf(v interface{}) *Object { return &Object{GoValue: v} }

func (o *Object) Eval(*Env) (*Object, error) { return o, nil }

func (*Object) Walk(Visitor) error { return nil }
