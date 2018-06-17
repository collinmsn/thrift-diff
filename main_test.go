package main

import (
	p "github.com/samuel/go-thrift/parser"
	"testing"
)

var noAnno = []*p.Annotation{}
var noFields = []*p.Field{}
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
		p.Type{"i64", nil, nil, noAnno},
		p.Type{"i64", nil, nil, noAnno},
	},
	{
		"type name changed", false,
		p.Type{"i64", nil, nil, noAnno},
		p.Type{"string", nil, nil, noAnno}},
	{
		"identical map", true,
		p.Type{"map", &strType, &i64Type, noAnno},
		p.Type{"map", &strType, &i64Type, noAnno},
	},
	{
		"map key changed", false,
		p.Type{"map", &i32Type, &i64Type, noAnno},
		p.Type{"map", &strType, &i64Type, noAnno},
	},
	{
		"map value changed", false,
		p.Type{"map", &strType, &i32Type, noAnno},
		p.Type{"map", &strType, &i64Type, noAnno},
	},
}

func field(id int, name string, type_ *p.Type) *p.Field {
	return &p.Field{id, name, false, type_, nil, noAnno}
}

var fieldTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  p.Field
	to                    p.Field
}{
	{
		"identical", true,
		*field(1, "foo", &i64Type),
		*field(1, "foo", &i64Type),
	},
	{
		"ID changed", false,
		*field(1, "foo", &i64Type),
		*field(2, "foo", &i64Type),
	},
	{
		"name changed", true,
		*field(1, "foo", &i64Type),
		*field(1, "bar", &i64Type),
	},
	{
		"type changed", false,
		*field(1, "foo", &i64Type),
		*field(1, "foo", &strType),
	},
	{
		"optional to required", false,
		p.Field{1, "foo", true, &i64Type, nil, noAnno},
		p.Field{1, "foo", false, &i64Type, nil, noAnno},
	},
	{
		"required to optional", true,
		p.Field{1, "foo", false, &i64Type, nil, noAnno},
		p.Field{1, "foo", true, &i64Type, nil, noAnno},
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
			{1, "foo", false, &i64Type, nil, noAnno},
			{2, "bar", false, &i32Type, nil, noAnno},
		}, noAnno},
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i64Type, nil, noAnno},
			{2, "bar", false, &i32Type, nil, noAnno},
		}, noAnno},
	},
	{
		"field added", true,
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i64Type, nil, noAnno},
		}, noAnno},
		p.Struct{"foo", []*p.Field{
			{6, "bar", false, &i32Type, nil, noAnno},
			{1, "foo", false, &i64Type, nil, noAnno},
		}, noAnno},
	},
	{
		"field removed", false,
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i64Type, nil, noAnno},
			{2, "bar", false, &i32Type, nil, noAnno},
		}, noAnno},
		p.Struct{"foo", []*p.Field{
			{2, "bar", false, &i32Type, nil, noAnno},
		}, noAnno},
	},
	{
		"field name changed", true,
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i64Type, nil, noAnno},
		}, noAnno},
		p.Struct{"foo", []*p.Field{
			{1, "bar", false, &i64Type, nil, noAnno},
		}, noAnno},
	},
	{
		"field type changed", false,
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i64Type, nil, noAnno},
		}, noAnno},
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i32Type, nil, noAnno},
		}, noAnno},
	},
	{
		"field name changed", true,
		p.Struct{"foo", []*p.Field{
			{1, "foo", false, &i64Type, nil, noAnno},
			{1, "foo", false, &i64Type, nil, noAnno},
		}, noAnno},
		p.Struct{"foo", []*p.Field{
			{1, "bar", false, &i64Type, nil, noAnno},
		}, noAnno},
	},
}

var methodTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  p.Method
	to                    p.Method
}{
	{
		"identical", true,
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
	},
	{
		"name changed", false,
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		p.Method{"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
	},
	{
		"return type changed", false,
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		p.Method{"", "foo", false, &i32Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
	},
	{
		"argument added", true,
		p.Method{"", "foo", false, &i64Type, []*p.Field{
			field(1, "foo", &strType),
		}, noFields, noAnno},
		p.Method{"", "foo", false, &i64Type, []*p.Field{
			field(43, "bar", &strType),
			field(1, "baz", &strType),
		}, noFields, noAnno},
	},
	{
		"argument removed", false,
		p.Method{"", "foo", false, &i64Type, []*p.Field{
			field(1, "foo", &strType),
			field(5, "bar", &strType),
		}, noFields, noAnno},
		p.Method{"", "foo", false, &i64Type, []*p.Field{
			field(5, "bar", &strType),
			field(2, "foo", &strType),
		}, noFields, noAnno},
	},
	{
		"argument name changed", true,
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "quuz", &strType)}, noFields, noAnno},
	},
	{
		"argument type changed", false,
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		p.Method{"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &i64Type)}, noFields, noAnno},
	},
}

var serviceTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  *p.Service
	to                    *p.Service
}{
	{
		"identical", true,
		&p.Service{"foo", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
		&p.Service{"foo", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
	},
	{
		"name changed", false,
		&p.Service{"foo", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
		&p.Service{"quuz", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
	},
	{
		"method name changed", false,
		&p.Service{"foo", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
		&p.Service{"foo", "", map[string]*p.Method{
			"quuz": {"", "quuz", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
	},
	{
		"method added", true,
		&p.Service{"foo", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
		&p.Service{"foo", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
			"quuz": {"", "quuz", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
	},
	{
		"method removed", false,
		&p.Service{"foo", "", map[string]*p.Method{
			"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
			"quuz": {"", "quuz", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
		&p.Service{"foo", "", map[string]*p.Method{
			"quuz": {"", "quuz", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
		}, noAnno},
	},
}

var thriftTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  *p.Thrift
	to                    *p.Thrift
}{
	{
		"identical", true,
		&p.Thrift{
			Services: map[string]*p.Service{},
		},
		&p.Thrift{
			Services: map[string]*p.Service{},
		},
	},
	{
		"service removed", false,
		&p.Thrift{
			Services: map[string]*p.Service{
				"MyService": {
					Name: "MyService",
					Methods: map[string]*p.Method{
						"foo": {"", "foo", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
						"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
					},
				},
			},
		},
		&p.Thrift{
			Services: map[string]*p.Service{
				"MyService": {
					Name: "MyService",
					Methods: map[string]*p.Method{
						"bar": {"", "bar", false, &i64Type, []*p.Field{field(1, "baz", &strType)}, noFields, noAnno},
					},
				},
			},
		},
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

func TestMethod(t *testing.T) {
	for _, mt := range methodTests {
		err := compareMethod(mt.from, mt.to)
		isBackwardsCompatible := err == nil
		if isBackwardsCompatible != mt.isBackwardsCompatible {
			t.Errorf("error in method test '%s': %v", mt.name, err)
		}
	}
}

func TestService(t *testing.T) {
	for _, st := range serviceTests {
		err := compareService(st.from, st.to)
		isBackwardsCompatible := err == nil
		if isBackwardsCompatible != st.isBackwardsCompatible {
			t.Errorf("error in service test '%s': %v", st.name, err)
		}
	}
}

func TestThrift(t *testing.T) {
	for _, tt := range thriftTests {
		err := compareThrift(tt.from, tt.to)
		isBackwardsCompatible := err == nil
		if isBackwardsCompatible != tt.isBackwardsCompatible {
			t.Errorf("error in service test '%s': %v", tt.name, err)
		}
	}
}
