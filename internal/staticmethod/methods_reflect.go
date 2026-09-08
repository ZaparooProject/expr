//go:build !expr_static_methods

// Package staticmethod isolates method lookup so an explicitly validated host
// can opt into a generated constant-name profile with expr_static_methods.
// Without that tag, every lookup delegates to reflect unchanged.
package staticmethod

import "reflect"

const Enabled = false

func TypeByName(t reflect.Type, name string) (reflect.Method, bool) {
	return t.MethodByName(name)
}

func ValueByName(v reflect.Value, name string) reflect.Value {
	return v.MethodByName(name)
}

func TypeByIndex(t reflect.Type, index int) reflect.Method {
	return t.Method(index)
}

func ValueByIndex(v reflect.Value, index int) reflect.Value {
	return v.Method(index)
}
