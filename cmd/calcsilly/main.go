package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"

	"github.com/jncornett/calc"
	"github.com/jncornett/calc/lang"
)

//go:embed silly.calc
var source []byte

func main() {
	env := calc.DefaultEnv()
	env.SetFunc("+", func(args ...*lang.Object) (*lang.Object, error) {
		var params [2]float64
		if err := decode(toIfaces(args...), &params); err != nil {
			return nil, err
		}
		return lang.ObjectOf(params[0] + params[1]), nil
	})
	env.SetFunc("sin", func(args ...*lang.Object) (*lang.Object, error) {
		var params [1]float64
		if err := decode(toIfaces(args...), &params); err != nil {
			return nil, err
		}
		return lang.ObjectOf(math.Sin(params[0])), nil
	})
	env.SetFunc("cos", func(args ...*lang.Object) (*lang.Object, error) {
		var params [1]float64
		if err := decode(toIfaces(args...), &params); err != nil {
			return nil, err
		}
		return lang.ObjectOf(math.Cos(params[0])), nil
	})
	got, err := calc.Eval(source, env)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(got)
}

func match(args []interface{}, alternatives ...interface{}) int {
	for i, alt := range alternatives {
		if err := decode(args, alt); err != nil {
			continue
		}
		return i
	}
	return -1
}

func toIfaces(args ...*lang.Object) []interface{} {
	var out []interface{}
	for _, arg := range args {
		out = append(out, arg.GoValue)
	}
	return out
}

func decode(args []interface{}, out interface{}) error {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Errorf("decode: out must be a pointer"))
	}
	rv = rv.Elem()
	switch rv.Kind() {
	default:
		panic(fmt.Errorf("decode: out must be a pointer to an array, struct, or slice"))
	case reflect.Array:
		return decodeArray(args, rv)
	case reflect.Struct:
		return decodeStruct(args, rv)
	case reflect.Slice:
		return decodeSlice(args, rv)
	}
}

func decodeArray(args []interface{}, ary reflect.Value) error {
	if len(args) != ary.Type().Len() {
		return errors.New("length mismatch")
	}
	elem := ary.Type().Elem()
	for i, arg := range args {
		rva := reflect.ValueOf(arg)
		if !rva.Type().AssignableTo(elem) {
			return fmt.Errorf("at arg %d: type mismatch: want %v, got %v", i, elem, rva.Type())
		}
		ary.Index(i).Set(rva)
	}
	return nil
}

func decodeSlice(args []interface{}, sl reflect.Value) error {
	elem := sl.Type().Elem()
	for _, arg := range args {
		rva := reflect.ValueOf(arg)
		if !rva.Type().AssignableTo(elem) {
			return errors.New("type mismatch")
		}
		sl.Set(reflect.Append(sl, rva))
	}
	return nil
}

func decodeStruct(args []interface{}, st reflect.Value) error {
	fields := reflect.VisibleFields(st.Type())
	var exported []reflect.StructField
	for _, f := range fields {
		if !f.IsExported() {
			continue
		}
		exported = append(exported, f)
	}
	if len(args) != len(exported) {
		return errors.New("length mismatch")
	}
	for i, arg := range args {
		sf := exported[i]
		rva := reflect.ValueOf(arg)
		if !rva.Type().AssignableTo(sf.Type) {
			return errors.New("type mismatch")
		}
		st.FieldByIndex(sf.Index).Set(rva)
	}
	return nil
}
