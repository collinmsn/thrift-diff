package main

import (
	"fmt"
	"github.com/samuel/go-thrift/parser"
	"os"
)

func getField(id int, s parser.Struct) *parser.Field {
	for _, f := range s.Fields {
		if f.ID == id {
			return f
		}
	}

	return nil
}

func compareStruct(from, to parser.Struct) error {
	for _, fromField := range from.Fields {
		toField := getField(fromField.ID, to)
		if toField == nil {
			return fmt.Errorf("field '%s' was removed from struct '%s'", to.Name)
		}

		if err := compareField(*fromField, *toField); err != nil {
			return err
		}
	}

	return nil
}

func compareField(from, to parser.Field) error {
	if from.ID != to.ID {
		return fmt.Errorf("field ID was changed. %d => %d", from.ID, to.ID)
	}

	if err := compareType(*from.Type, *to.Type); err != nil {
		return err
	}

	if from.Optional != to.Optional {
		if from.Optional == true && to.Optional == false {
			return fmt.Errorf("field cannot be made optional once required. %t => %t", from.Optional, to.Optional)
		}
	}

	return nil
}

func compareType(from, to parser.Type) error {
	if from.Name != to.Name {
		return fmt.Errorf("type was changed. %s => %s", from.Name, to.Name)
	}

	if from.KeyType != nil {
		if err := compareType(*from.KeyType, *to.KeyType); err != nil {
			return fmt.Errorf("type key was changed. %s => %s", from.Name, to.Name)
		}
	}

	if from.ValueType != nil {
		if err := compareType(*from.ValueType, *to.ValueType); err != nil {
			return fmt.Errorf("type value was changed. %s => %s", from.Name, to.Name)
		}
	}

	return nil
}

func printUsage() {
	fmt.Printf("Usage: %s [FROM_THRIFT_FILE] [TO_THRIFT_FILE]\n", os.Args[0])
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "No Thrift files specified!")
		printUsage()
		os.Exit(1)
	}

	fromFile := os.Args[1]
	toFile := os.Args[2]

	if _, err := os.Stat(fromFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File %s does not exist!", fromFile)
		os.Exit(1)
	}

	if _, err := os.Stat(toFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File %s does not exist!", toFile)
		os.Exit(1)
	}

	p := &parser.Parser{}

	fromThrifts, _, err := p.ParseFile(fromFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	toThrifts, _, err := p.ParseFile(toFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	// Hard coded hack for testing
	var fromThrift *parser.Thrift
	for _, t := range fromThrifts {
		fromThrift = t
		break
	}

	var toThrift *parser.Thrift
	for _, t := range toThrifts {
		toThrift = t
		break
	}

	err = compareStruct(*fromThrift.Structs["User"], *toThrift.Structs["User"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "not backwards compatible: %s\n", err.Error())
	}
}
