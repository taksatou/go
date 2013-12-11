package functional

import (
	//	"fmt"
	"errors"
	"reflect"
)

type F0 func(interface{}) func() interface{}
type F1 func(interface{}) func(interface{}) interface{}
type F2 func(interface{}) func(interface{}, interface{}) interface{}

func applyLeft(x, y interface{}) (interface{}, error) {
	var e error
	xv := reflect.ValueOf(x)

	switch xv.Type().NumIn() {
	case 1:
		var f0 F0
		e = Curry(x, &f0)
		if e != nil {
			return nil, e
		}
		return f0(y)(), nil // reduction!
	case 2:
		var f1 F1
		e = Curry(x, &f1)
		if e != nil {
			return nil, e
		}
		return f1(y), nil
	case 3:
		var f2 F2
		e = Curry(x, &f2)
		if e != nil {
			return nil, e
		}
		return f2(y), nil
	default:
		return nil, errors.New("invalid args")
	}

	return nil, nil
}

func applyLeftOrPanic(x interface{}, y ...interface{}) interface{} {
	var e error
	for _, a := range y {
		x, e = applyLeft(x, a)
		if e != nil {
			panic(e)
		}
	}
	return x
}

func K(x, y interface{}) interface{} {
	return x
}

func S(x, y, z interface{}) interface{} {
	xz := applyLeftOrPanic(x, z)
	yz := applyLeftOrPanic(y, z)
	return applyLeftOrPanic(xz, yz)
}

var (
	_i = func() interface{} {
		return applyLeftOrPanic(S, K, K)
	}()

	_c = func() interface{} {
		kk := applyLeftOrPanic(K, K)
		bbs := applyLeftOrPanic(B, B, S)
		return applyLeftOrPanic(S, bbs, kk)
	}()

	_b = func() interface{} {
		ks := applyLeftOrPanic(K, S)
		return applyLeftOrPanic(S, ks, K)
	}()

	_m = func() interface{} {
		return applyLeftOrPanic(S, I, I)
	}()

	_l = func() interface{} {
		return applyLeftOrPanic(C, B, M)
	}()

	_y = func() interface{} {
		return applyLeftOrPanic(S, L, L)
	}()
)

func I(x interface{}) interface{} {
	return applyLeftOrPanic(_i, x)
}

func C(x, y, z interface{}) interface{} {
	return applyLeftOrPanic(_c, x, y, z)
}

func B(x, y, z interface{}) interface{} {
	return applyLeftOrPanic(_b, x, y, z)
}

func M(x interface{}) interface{} {
	return applyLeftOrPanic(_m, x)
}

func L(x, y interface{}) interface{} {
	return applyLeftOrPanic(_l, x, y)
}

// // Y combinator loops forever...
// func Y(x interface{}) interface{} {
// 	return applyLeftOrPanic(_y, x)
// }

// func y(x interface{}) interface{} {
// 	return reflect.ValueOf(x).Call([]reflect.Value{reflect.ValueOf(y(x))})[0].Interface()
// }
