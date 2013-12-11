package functional

import (
	"errors"
	"fmt"
	"reflect"
)

func assertIterable(a reflect.Value) error {
	if k := a.Kind(); k != reflect.Array && k != reflect.Slice && k != reflect.String && k != reflect.Map { // TODO: support channels
		return errors.New("invalid argument")
	}
	return nil
}

// Curry makes a function that takes first argument of original function
// and returns a function that takes remaining arguments.
func Curry(f interface{}, curried interface{}) error {
	foutv := reflect.ValueOf(curried).Elem()
	foutt := foutv.Type()

	fv := reflect.ValueOf(f)
	if fv.Type().NumIn() <= 0 {
		return errors.New("f must take at least one argument")
	}

	if foutt.NumOut() != 1 || foutt.Out(0).Kind() != reflect.Func {
		return errors.New("curried must return a function")
	}

	foutv.Set(reflect.MakeFunc(foutt, func(args1 []reflect.Value) []reflect.Value {
		f := reflect.MakeFunc(foutt.Out(0), func(args2 []reflect.Value) []reflect.Value {
			return fv.Call(append(args1, args2...))
		})
		return []reflect.Value{f}
	}))
	return nil
}

// Compose makes a function that is a composition of given two functions.
// `f` and `g` can take and return arbitrary type of variables,
// as long as they have the same type respectively.
func Compose(f, g, composed interface{}) error {
	foutv := reflect.ValueOf(composed).Elem()

	fv, gv := reflect.ValueOf(f), reflect.ValueOf(g)
	ft, gt := fv.Type(), gv.Type()

	if ft.NumIn() != gt.NumOut() {
		return errors.New(fmt.Sprintf("f takes %d args, but g returns %d args", ft.NumIn(), gt.NumOut()))
	}
	for i, n := 0, ft.NumIn(); i < n; i++ {
		if ft.In(i) != gt.Out(i) {
			return errors.New(fmt.Sprintf("f takes <%s>, but g returns <%s>", ft.In(i), gt.Out(i)))
		}
	}

	foutv.Set(reflect.MakeFunc(foutv.Type(), func(args []reflect.Value) []reflect.Value {
		return fv.Call(gv.Call(args))
	}))
	return nil
}

// Flip makes a function that takes first two arguments in
// reverse order of original function.
func Flip(f interface{}, flipped interface{}) error {
	foutv := reflect.ValueOf(flipped).Elem()

	fv := reflect.ValueOf(f)
	if fv.Type().NumIn() < 2 {
		return errors.New("f must take at least two arguments")
	}

	foutv.Set(reflect.MakeFunc(foutv.Type(), func(args []reflect.Value) []reflect.Value {
		a := make([]reflect.Value, len(args))
		a[0], a[1] = args[1], args[0]
		copy(a[2:], args[2:])
		return fv.Call(a)
	}))
	return nil
}

// Fold applys f of two arguments cumulatively to the items of lis.
// initial must be the same type as contents of lis or nil.
func Fold(f interface{}, lis interface{}, initial interface{}) (interface{}, error) {
	lisv := reflect.ValueOf(lis)
	e := assertIterable(lisv)
	if e != nil {
		return nil, e
	}
	if lisv.Kind() == reflect.Map {
		return nil, errors.New("map is not supported now")
	}

	fn := reflect.ValueOf(f)
	if fn.Type().NumIn() != 2 {
		return nil, errors.New("f must take exactly two argument")
	}

	if fn.Type().NumOut() != 1 {
		return nil, errors.New("curried must return exactly one value")
	}

	if initial != nil {
		iniv := reflect.ValueOf(initial)
		if lisv.Len() > 0 && lisv.Index(0).Type() != iniv.Type() {
			return nil, errors.New("type mismatch")
		}
		lisv = reflect.Append(lisv, iniv)
	}

	l := lisv.Len()
	if l <= 0 {
		return nil, nil
	}

	v0 := lisv.Index(0)
	if !(fn.Type().NumOut() == 1 && fn.Type().Out(0) == v0.Type()) {
		return nil, errors.New(fmt.Sprintf("function should return single <%s>, but got <%s>", v0.Type(), fn.Type().Out(0)))
	}
	for i := 1; i < l; i++ {
		v0 = fn.Call([]reflect.Value{v0, lisv.Index(i)})[0]
	}
	return v0.Interface(), nil
}

// Map applys f to all items of lis and returns the results
func Map(lis interface{}, f interface{}) (interface{}, error) {
	a := reflect.ValueOf(lis)
	e := assertIterable(a)
	if e != nil {
		return nil, e
	}

	fn := reflect.ValueOf(f)
	if a.Kind() == reflect.Map {
		if fn.Type().NumIn() != 2 {
			return nil, errors.New("f must take exactly one argument")
		}

		if fn.Type().NumOut() != 2 {
			return nil, errors.New("curried must return exactly one value")
		}
	} else {
		if fn.Type().NumIn() != 1 {
			return nil, errors.New("f must take exactly one argument")
		}

		if fn.Type().NumOut() != 1 {
			return nil, errors.New("curried must return exactly one value")
		}
	}
	if a.Kind() == reflect.String {
		buf := make([]rune, a.Len())
		f, ok := f.(func(rune) rune)
		if !ok {
			return nil, errors.New("invalid filter function")
		}
		for i, c := range lis.(string) {
			buf[i] = f(c)
		}
		return string(buf), nil
	} else if a.Kind() == reflect.Map {
		res := reflect.MakeMap(a.Type())
		for _, k := range a.MapKeys() {
			v := a.MapIndex(k)
			out := fn.Call([]reflect.Value{k, v})
			res.SetMapIndex(out[0], out[1])
		}
		return res.Interface(), nil
	} else {
		res := reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), a.Len(), a.Cap())
		for i, l := 0, a.Len(); i < l; i++ {
			v := a.Index(i)
			res.Index(i).Set(fn.Call([]reflect.Value{v})[0])
		}
		return res.Interface(), nil
	}
	return nil, errors.New("never happen")
}

// Map applys f to all items of lis and update them
// func Map2(lis interface{}, f interface{}) error {
// 	a := reflect.ValueOf(lis)
// 	e := assertIterable(a)
// 	if e != nil {
// 		return e
// 	}
// 	fn := reflect.ValueOf(f)
// 	if a.Kind() == reflect.Map {
// 		for _, k := range a.MapKeys() {
// 			v := a.MapIndex(k)
// 			res := fn.Call([]reflect.Value{k, v})
// 			k.Set(res[0])
// 			v.Set(res[1])
// 		}
// 	} else {
// 		for i, l := 0, a.Len(); i < l; i++ {
// 			v := a.Index(i)
// 			r := fn.Call([]reflect.Value{v})
// 			v.Set(r[0])
// 		}
// 	}
// 	return nil
// }

func Each(args interface{}, f interface{}) error {
	a := reflect.ValueOf(args)
	e := assertIterable(a)
	if e != nil {
		return e
	}

	fn := reflect.ValueOf(f)
	if a.Kind() == reflect.Map {
		for _, k := range a.MapKeys() {
			fn.Call([]reflect.Value{k, a.MapIndex(k)})
		}
	} else {
		for i, l := 0, a.Len(); i < l; i++ {
			fn.Call([]reflect.Value{a.Index(i)})
		}
	}
	return nil
}

func Filter(arg interface{}, f interface{}) (interface{}, error) {
	a := reflect.ValueOf(arg)
	e := assertIterable(a)
	if e != nil {
		return nil, e
	}

	fn := reflect.ValueOf(f)
	if a.Kind() == reflect.String {
		buf := []rune{}
		f, ok := f.(func(rune) bool)
		if !ok {
			return nil, errors.New("invalid filter function")
		}
		for _, c := range arg.(string) {
			if f(c) {
				buf = append(buf, c)
			}
		}
		return string(buf), nil
	} else if a.Kind() == reflect.Map {
		res := reflect.MakeMap(a.Type())
		for _, k := range a.MapKeys() {
			v := a.MapIndex(k)
			r := fn.Call([]reflect.Value{k, v})
			if r[0].Bool() {
				res.SetMapIndex(k, v)
			}
		}
		return res.Interface(), nil
	} else {
		res := reflect.New(a.Type())
		e := res.Elem()
		for i, l := 0, a.Len(); i < l; i++ {
			v := a.Index(i)
			r := fn.Call([]reflect.Value{v})
			if r[0].Bool() {
				res = reflect.Append(e, v)
			}
		}
		return res.Interface(), nil
	}
	return nil, errors.New("never happen")
}

// func Filter2(arg interface{}, f interface{}) error {
// 	a := reflect.ValueOf(arg)
// 	if a.Kind() != reflect.Ptr {
// 		return errors.New("arg should be settable")
// 	}
// 	a = a.Elem()
// 	e := assertIterable(a)
// 	if e != nil {
// 		return e
// 	}

// 	fn := reflect.ValueOf(f)
// 	if a.Kind() == reflect.String {
// 		buf := []rune{}
// 		f, ok := f.(func(rune) bool)
// 		if !ok {
// 			return errors.New("invalid filter function")
// 		}
// 		for _, c := range *(arg.(*string)) {
// 			if f(c) {
// 				buf = append(buf, c)
// 			}
// 		}
// 		a.Set(reflect.ValueOf(string(buf)))
// 	} else if a.Kind() == reflect.Map {
// 		res := reflect.MakeMap(a.Type())
// 		for _, k := range a.MapKeys() {
// 			v := a.MapIndex(k)
// 			r := fn.Call([]reflect.Value{k, v})
// 			if r[0].Bool() {
// 				res.SetMapIndex(k, v)
// 			}
// 		}
// 		a.Set(res)
// 	} else {
// 		res := reflect.New(a.Type())
// 		e := res.Elem()
// 		for i, l := 0, a.Len(); i < l; i++ {
// 			v := a.Index(i)
// 			r := fn.Call([]reflect.Value{v})
// 			if r[0].Bool() {
// 				res = reflect.Append(e, v)
// 			}
// 		}
// 		a.Set(res)
// 	}
// 	return nil
// }
