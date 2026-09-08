package staticmethod

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func sameMethod(a, b reflect.Method) bool {
	if a.Name != b.Name || a.PkgPath != b.PkgPath || a.Type != b.Type || a.Index != b.Index ||
		a.Func.IsValid() != b.Func.IsValid() {
		return false
	}
	return !a.Func.IsValid() || a.Func.Type() == b.Func.Type() && a.Func.Pointer() == b.Func.Pointer()
}

func TestCoreMethodMetadata(t *testing.T) {
	types := []reflect.Type{
		reflect.TypeOf(time.Time{}), reflect.TypeOf((*time.Time)(nil)),
		reflect.TypeOf(time.Duration(0)), reflect.TypeOf((*time.Duration)(nil)),
		reflect.TypeOf(time.Month(0)), reflect.TypeOf((*time.Month)(nil)),
		reflect.TypeOf(time.Weekday(0)), reflect.TypeOf((*time.Weekday)(nil)),
		reflect.TypeOf((*time.Location)(nil)), reflect.TypeOf((*error)(nil)).Elem(),
	}
	count := 0
	for _, typ := range types {
		for i := 0; i < typ.NumMethod(); i++ {
			want := typ.Method(i)
			got := TypeByIndex(typ, i)
			byName, ok := TypeByName(typ, want.Name)
			if !ok || !sameMethod(want, got) || !sameMethod(want, byName) {
				t.Fatalf("method metadata differs: type=%v method=%s", typ, want.Name)
			}
			if typ.Kind() != reflect.Interface {
				v := reflect.Zero(typ)
				if v.Method(i).Type() != ValueByIndex(v, i).Type() || v.Method(i).Type() != ValueByName(v, want.Name).Type() {
					t.Fatalf("bound signature differs: type=%v method=%s", typ, want.Name)
				}
			}
			count++
		}
		if _, ok := TypeByName(typ, "NotAnExistingMethod"); ok {
			t.Fatalf("unexpected method on %v", typ)
		}
	}
	t.Logf("checked %d method descriptors", count)
}

type futureHostType struct{}

func (futureHostType) FutureCardMethod() string { return "unchanged" }

func TestUnknownHostMethod(t *testing.T) {
	v := reflect.ValueOf(futureHostType{})
	if !Enabled {
		m := ValueByName(v, "FutureCardMethod")
		if !m.IsValid() || m.Call(nil)[0].String() != "unchanged" {
			t.Fatal("default mode changed host method behavior")
		}
		return
	}
	defer func() {
		p := recover()
		if p == nil || !strings.Contains(fmt.Sprint(p), "unvalidated method set") {
			t.Fatalf("static mode silently accepted unvalidated host type: %v", p)
		}
	}()
	ValueByName(v, "FutureCardMethod")
}
