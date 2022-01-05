package lang

type Env struct {
	SymbolTable
}

func (env *Env) SetFunc(ref Ref, fn NativeFunc) {
	env.Set(ref, ObjectOf(fn))
}

type SymbolTable interface {
	Get(ref Ref) (o *Object, ok bool)
	Set(ref Ref, o *Object)
}

type Chained []SymbolTable

var _ SymbolTable = Chained(nil)

func (ch Chained) Get(ref Ref) (o *Object, ok bool) {
	for i := len(ch) - 1; i >= 0; i-- {
		o, ok = ch[i].Get(ref)
		if ok {
			return o, true
		}
	}
	return nil, false
}

func (ch Chained) Set(ref Ref, o *Object) {
	if len(ch) == 0 {
		return
	}
	ch[len(ch)-1].Set(ref, o)
}

type Map map[Ref]*Object

func (m Map) Get(ref Ref) (o *Object, ok bool) {
	o, ok = m[ref]
	return o, ok
}

func (m Map) Set(ref Ref, o *Object) {
	m[ref] = o
}
