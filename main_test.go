package main

import (
	"testing"
	"strings"
)

var rawStructTests = []struct {
	name                  string
	isBackwardsCompatible bool
	from                  string
	to                    string
}{
	{
		"field name change", true,
		"struct Foo { 1: i32 bar }",
		"struct Foo { 1: i32 quuz }",
	},
	{
		"field type change", false,
		"struct Foo { 1: i32 bar }",
		"struct Foo { 1: i64 bar }",
	},
	{
		"field map key type change", false,
		"struct Foo { 1: map<string, string> bar }",
		"struct Foo { 1: map<i64, string> bar }",
	},
	{
		"field map value type change", false,
		"struct Foo { 1: map<string, string> bar }",
		"struct Foo { 1: map<string, i64> bar }",
	},
	{
		"field id change", false,
		"struct Foo { 1: i32 bar }",
		"struct Foo { 2: i32 bar }",
	},
	{
		"struct rename", false,
		"struct Foo { 1: i32 bar }",
		"struct Baz { 1: i32 bar }",
	},
	{
		"struct remove", false,
		"struct Foo { 1: i32 bar }; struct Baz { 1: i32 bar }",
		"struct Baz { 1: i32 bar }",
	},
	{
		"add field", true,
		"struct Foo { 1: i32 bar }",
		"struct Foo { 1: i32 bar, 2: i32 baz }",
	},
	{
		"field in substruct changed", false,
		"struct Foo { 1: i32 bar }; struct Baz { 1: Foo bar }",
		"struct Foo { 1: i64 bar }; struct Baz { 1: Foo bar }",
	},
	{
		"field in substruct changed", false,
		"struct Foo { 1: i32 bar }; struct Baz { 1: Foo bar }",
		"struct Foo { 1: i64 bar }; struct Baz { 1: Foo bar }",
	},
	{
		"field optional -> required", true,
		"struct Foo { 1: optional i32 bar }",
		"struct Foo { 1: required i32 bar }",
	},
	{
		"field required -> optional", false,
		"struct Foo { 1: required i32 bar }",
		"struct Foo { 1: optional i32 bar }",
	},
	// Not handled by current Thrift IDL parser
	//{
	//	"field required -> optional (default)", false,
	//	"struct Foo { 1: required i32 bar }",
	//	"struct Foo { 1: i32 bar }",
	//},
	{
		"method name change", false,
		"service Foo { void ping() }",
		"service Foo { void hello() }",
	},
	{
		"method add", true,
		"service Foo { void hello() }",
		"service Foo { void ping(); void hello() }",
	},
	{
		"method remove", false,
		"service Foo { void ping(); void hello() }",
		"service Foo { void hello() }",
	},
	{
		"method return type change", false,
		"service Foo { void ping() }",
		"service Foo { i32 ping() }",
	},
	{
		"method argument add", true,
		"service Foo { void ping(1: i32 foo) }",
		"service Foo { void ping(1: i32 foo, 2: i64 bar) }",
	},
	{
		"method argument remove", false,
		"service Foo { void ping(1: i32 foo, 2: i64 bar) }",
		"service Foo { void ping(1: i32 foo) }",
	},
	{
		"method argument name change", true,
		"service Foo { void ping(1: i32 foo) }",
		"service Foo { void ping(1: i32 bar) }",
	},
	{
		"method argument type change", false,
		"service Foo { void ping(1: i32 foo) }",
		"service Foo { void ping(1: i64 foo) }",
	},
	{
		"service add", true,
		"service Foo { void ping() }",
		"service Foo { void ping() }; service Bar { void ping() }",
	},
	{
		"service remove", false,
		"service Foo { void ping() }; service Bar { void ping() }",
		"service Foo { void ping() }",
	},
	{
		"service name change", true,
		"service Foo { void ping() }",
		"service Bar { void ping() }",
	},
}

func TestThrift(t *testing.T) {
	for i := 0; i < 100; i++ {
		for _, tt := range rawStructTests {
			fromThrift := parse(tt.from)
			toThrift := parse(tt.to)

			err := compareThrift(fromThrift, toThrift)
			isBackwardsCompatible := err == nil
			if isBackwardsCompatible != tt.isBackwardsCompatible {
				t.Errorf("error in raw Thrift test '%s': %v", tt.name, err)
				divider := strings.Repeat("-", len(tt.from) + 7)
				t.Errorf("\nFrom: %s\n%s\n  To: %s", tt.from, divider, tt.to)
				t.FailNow()
			}
		}
	}
}
