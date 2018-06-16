package main

import (
	p "github.com/samuel/go-thrift/parser"
	"testing"
)

var emptyAnno = []*p.Annotation{}
var i64Type = p.Type{"i64", nil, nil, []*p.Annotation{}}
var i32Type = p.Type{"i32", nil, nil, []*p.Annotation{}}
var strType = p.Type{"string", nil, nil, []*p.Annotation{}}

var typeTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  p.Type
	to                    p.Type
}{
	{
		"identical", true,
		p.Type{"i64", nil, nil, emptyAnno},
		p.Type{"i64", nil, nil, emptyAnno},
	},
	{
		"type name changed", false,
		p.Type{"i64", nil, nil, emptyAnno},
		p.Type{"string", nil, nil, emptyAnno}},
	{
		"identical map", true,
		p.Type{"map", &strType, &i64Type, emptyAnno},
		p.Type{"map", &strType, &i64Type, emptyAnno},
	},
	{
		"map key changed", false,
		p.Type{"map", &i32Type, &i64Type, emptyAnno},
		p.Type{"map", &strType, &i64Type, emptyAnno},
	},
	{
		"map value changed", false,
		p.Type{"map", &strType, &i32Type, emptyAnno},
		p.Type{"map", &strType, &i64Type, emptyAnno},
	},
}

func field(id int, name string, type_ *p.Type) p.Field {
	return p.Field{id, name, false, type_, nil, emptyAnno}
}

var fieldTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  p.Field
	to                    p.Field
}{
	{
		"identical", true,
		field(1, "foo", &i64Type),
		field(1, "foo", &i64Type),
	},
	{
		"ID changed", false,
		field(1, "foo", &i64Type),
		field(2, "foo", &i64Type),
	},
	{
		"name changed", true,
		field(1, "foo", &i64Type),
		field(1, "bar", &i64Type),
	},
	{
		"type changed", false,
		field(1, "foo", &i64Type),
		field(1, "foo", &strType),
	},
	{
		"optional to required", false,
		p.Field{1, "foo", true, &i64Type, nil, emptyAnno},
		p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
	},
	{
		"required to optional", true,
		p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
		p.Field{1, "foo", true, &i64Type, nil, emptyAnno},
	},
}

var structTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  p.Struct
	to                    p.Struct
}{
	{
		"identical", true,
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
			&p.Field{2, "bar", false, &i32Type, nil, emptyAnno},
		}, emptyAnno},
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
			&p.Field{2, "bar", false, &i32Type, nil, emptyAnno},
		}, emptyAnno},
	},
	{
		"field added", true,
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
		}, emptyAnno},
		p.Struct{"foo", []*p.Field{
			&p.Field{6, "bar", false, &i32Type, nil, emptyAnno},
			&p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
		}, emptyAnno},
	},
	{
		"field removed", false,
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
			&p.Field{2, "bar", false, &i32Type, nil, emptyAnno},
		}, emptyAnno},
		p.Struct{"foo", []*p.Field{
			&p.Field{2, "bar", false, &i32Type, nil, emptyAnno},
		}, emptyAnno},
	},
	{
		"field name changed", true,
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
		}, emptyAnno},
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "bar", false, &i64Type, nil, emptyAnno},
		}, emptyAnno},
	},
	{
		"field type changed", false,
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "foo", false, &i64Type, nil, emptyAnno},
		}, emptyAnno},
		p.Struct{"foo", []*p.Field{
			&p.Field{1, "foo", false, &i32Type, nil, emptyAnno},
		}, emptyAnno},
	},
	{
		"field name changed", true,
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i64Type, nil, emptyAnno},
			{1, "foo", false, &i64Type, nil, emptyAnno},
		}, emptyAnno},
		p.Struct{"foo", []*p.Field{
			{1, "bar", false, &i64Type, nil, emptyAnno},
		}, emptyAnno},
	},
}

func TestType(t *testing.T) {
	for _, tt := range typeTests {
		err := compareType(tt.from, tt.to)
		isBackwardsCompatible := err == nil
		if isBackwardsCompatible != tt.isBackwardsCompatible {
			t.Errorf("error in type test '%s': %v", tt.name, err)
		}
	}
}

func TestField(t *testing.T) {
	for _, ft := range fieldTests {
		err := compareField(ft.from, ft.to)
		isBackwardsCompatible := err == nil
		if isBackwardsCompatible != ft.isBackwardsCompatible {
			t.Errorf("error in field test '%s': %v", ft.name, err)
		}
	}
}

func TestStruct(t *testing.T) {
	for _, st := range structTests {
		err := compareStruct(st.from, st.to)
		isBackwardsCompatible := err == nil
		if isBackwardsCompatible != st.isBackwardsCompatible {
			t.Errorf("error in struct test '%s': %v", st.name, err)
		}
	}
}
