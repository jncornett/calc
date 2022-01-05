package lang

import "errors"

type Node interface {
	Eval(env *Env) (*Object, error)
	Walk(v Visitor) error
}

var ErrSkip = errors.New("skip")

type Visitor interface {
	Push(key interface{}, n Node) error
	Pop(key interface{}, n Node) error
}

type VisitorFuncs struct {
	PushFunc func(key interface{}, n Node) error
	PopFunc  func(key interface{}, n Node) error
}

var _ Visitor = VisitorFuncs{}

func (vf VisitorFuncs) Push(key interface{}, n Node) error {
	if vf.PushFunc == nil {
		return nil
	}
	return vf.PushFunc(key, n)
}

func (vf VisitorFuncs) Pop(key interface{}, n Node) error {
	if vf.PopFunc == nil {
		return nil
	}
	return vf.PopFunc(key, n)
}

type NativeFunc func(args ...*Object) (*Object, error)

func parseVisitorError(e error) (skip bool, err error) {
	if e != nil {
		if errors.Is(err, ErrSkip) {
			return true, nil
		}
		return false, e
	}
	return false, nil
}
