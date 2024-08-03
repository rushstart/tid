package tid

import (
	"reflect"
	"strings"
)

type ID struct {
	typ reflect.Type
	tag string
	ptr bool
}

func (id ID) String() string {
	s := make([]string, 0, 6)
	if id.ptr {
		s = append(s, "*")
	}
	pkgPath := id.typ.PkgPath()
	if pkgPath != "" {
		s = append(s, pkgPath, "/")
	}
	name := id.typ.Name()
	if name == "" {
		name = id.typ.String()
	}
	s = append(s, name)
	if id.tag != "" {
		s = append(s, "#", id.tag)
	}
	return strings.Join(s, "")
}

func (id ID) Type() reflect.Type {
	return id.typ
}

func (id ID) Tag() string {
	return id.tag
}

func (id ID) IsPtr() bool {
	return id.ptr
}

func (id ID) asPtr() ID {
	id.ptr = true
	return id
}
func (id ID) asNotPtr() ID {
	id.ptr = false
	return id
}

func From[T any](tag ...string) ID {
	var t T
	tTyp := reflect.TypeOf(t)
	if tTyp == nil {
		return FromType(reflect.TypeOf(&t), first(tag)).asNotPtr()
	}

	return FromType(tTyp, first(tag))
}

func FromType(typ reflect.Type, tag ...string) ID {
	if typ.Kind() == reflect.Pointer {
		return FromType(typ.Elem(), first(tag)).asPtr()
	}

	return ID{typ: typ, tag: first(tag)}
}
