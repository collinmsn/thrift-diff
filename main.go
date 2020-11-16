package main

import (
	"fmt"
	"github.com/collinmsn/go-thrift/parser"
	"os"
	"path"
	"path/filepath"
)

func getField(id int, fields []*parser.Field) *parser.Field {
	for _, f := range fields {
		if f.ID == id {
			return f
		}
	}

	return nil
}

// The compare functions return an error when a change is not backwards compatible

func compareFields(from, to []*parser.Field) error {
	for _, fromField := range from {
		toField := getField(fromField.ID, to)
		if toField == nil {
			return fmt.Errorf("field '%s' was removed", fromField.Name)
		}

		if err := compareField(*fromField, *toField); err != nil {
			return err
		}
	}

	return nil
}

func compareStruct(from, to parser.Struct) error {
	return compareFields(from.Fields, to.Fields)
}

func compareMethod(from, to parser.Method) error {
	if from.Name != to.Name {
		return fmt.Errorf("name was changed. %s => %s", from.Name, to.Name)
	}

	if err := compareType(from.ReturnType, to.ReturnType); err != nil {
		return fmt.Errorf("method '%s': %v", from.Name, err)
	}

	if err := compareFields(from.Arguments, to.Arguments); err != nil {
		return fmt.Errorf("method '%s': %v", from.Name, err)
	}

	return nil
}

func compareService(from, to *parser.Service) error {
	if from.Name != to.Name {
		return fmt.Errorf("name was changed. %s => %s", from.Name, to.Name)
	}

	for _, fromMethod := range from.Methods {
		var found = false
		for _, toMethod := range to.Methods {
			if fromMethod.Name == toMethod.Name {
				if err := compareMethod(*fromMethod, *toMethod); err != nil {
					return fmt.Errorf("method '%s' was changed: %v", fromMethod.Name, err)
				}
				found = true
			}
		}

		if !found {
			return fmt.Errorf("method '%s' was removed", fromMethod.Name)
		}
	}

	return nil
}

func compareField(from, to parser.Field) error {
	if from.ID != to.ID {
		return fmt.Errorf("field ID was changed. %d => %d", from.ID, to.ID)
	}

	if err := compareType(from.Type, to.Type); err != nil {
		return err
	}

	if from.Optional == false && to.Optional == true {
		return fmt.Errorf("field cannot be made optional once required. %t => %t", from.Optional, to.Optional)
	}

	return nil
}

func compareType(from, to *parser.Type) error {
	if from == nil && to == nil {
		return nil
	}

	if from == nil || to == nil {
		return fmt.Errorf("type was changed.")
	}

	if from.Name != to.Name {
		return fmt.Errorf("type was changed. %s => %s", from.Name, to.Name)
	}

	if from.KeyType != nil {
		if err := compareType(from.KeyType, to.KeyType); err != nil {
			return fmt.Errorf("type key was changed. %s => %s", from.Name, to.Name)
		}
	}

	if from.ValueType != nil {
		if err := compareType(from.ValueType, to.ValueType); err != nil {
			return fmt.Errorf("type value was changed. %s => %s", from.Name, to.Name)
		}
	}

	return nil
}

func compareThrift(from, to *parser.Thrift) error {
	for _, fromService := range from.Services {
		var found = false
		for _, toService := range to.Services {
			if fromService.Name == toService.Name {
				if err := compareService(fromService, toService); err != nil {
					return fmt.Errorf("service '%s' was changed: %v", fromService.Name, err)
				}
				found = true
			}
		}

		if !found {
			return fmt.Errorf("service '%s' was removed", fromService.Name)
		}
	}

	for _, fromStruct := range from.Structs {
		var found = false
		for _, toStruct := range to.Structs {
			if fromStruct.Name == toStruct.Name {
				if err := compareStruct(*fromStruct, *toStruct); err != nil {
					return fmt.Errorf("struct '%s' was changed: %v", fromStruct.Name, err)
				}
				found = true
			}
		}

		if !found {
			return fmt.Errorf("struct '%s' was removed", fromStruct.Name)
		}
	}

	return nil
}

func compareThrifts(fromIn, toIn map[string]*parser.Thrift) error {
	from, to := make(map[string]*parser.Thrift), make(map[string]*parser.Thrift)
	for k, v := range fromIn {
		from[filepath.Base(k)] = v
	}
	for k, v := range toIn {
		to[filepath.Base(k)] = v
	}

	for k, v := range from {
		if vto, ok := to[k]; !ok {
			err := fmt.Errorf("to missing file: %s\n", k)
			return err
		} else if err := compareThrift(v, vto); err != nil {
			return err
		}
	}

	return nil
}

func printUsage() {
	fmt.Printf("Usage: %s [INCLUDE_ROOT] [FROM_THRIFT_FILE] [TO_THRIFT_FILE]\n", os.Args[0])
	fmt.Printf("Usage: check compatibilites for from -> to and then to -> from\n")
}

func mergeThriftFiles(files map[string]*parser.Thrift) (*parser.Thrift, error) {
	var res = parser.Thrift{
		Typedefs: map[string]*parser.Typedef{},
		Namespaces: map[string]string{},
		Constants: map[string]*parser.Constant{},
		Enums: map[string]*parser.Enum{},
		Structs: map[string]*parser.Struct{},
		Exceptions: map[string]*parser.Struct{},
		Unions: map[string]*parser.Struct{},
		Services: map[string]*parser.Service{},
	}

	for _, t := range files {
		fmt.Printf("file: %#v\n", t)
		for k, v := range t.Typedefs {
			if _, exists := res.Typedefs[k]; exists {
				return nil, fmt.Errorf("key %s already exists", k)
			}
			res.Typedefs[k] = v
		}

		for k, v := range t.Namespaces {
			if _, exists := res.Namespaces[k]; exists {
				return nil, fmt.Errorf("key %s already exists, namespaces: %v", k, t.Namespaces)
			}
			res.Namespaces[k] = v
		}

		for k, v := range t.Constants {
			if _, exists := res.Constants[k]; exists {
				return nil, fmt.Errorf("key %s already exists", k)
			}
			res.Constants[k] = v
		}

		for k, v := range t.Enums {
			if _, exists := res.Enums[k]; exists {
				return nil, fmt.Errorf("key %s already exists", k)
			}
			res.Enums[k] = v
		}

		for k, v := range t.Structs {
			if _, exists := res.Structs[k]; exists {
				return nil, fmt.Errorf("key %s already exists", k)
			}
			res.Structs[k] = v
		}

		for k, v := range t.Exceptions {
			if _, exists := res.Exceptions[k]; exists {
				return nil, fmt.Errorf("key %s already exists", k)
			}
			res.Exceptions[k] = v
		}

		for k, v := range t.Unions {
			if _, exists := res.Unions[k]; exists {
				return nil, fmt.Errorf("key %s already exists", k)
			}
			res.Unions[k] = v
		}

		for k, v := range t.Services {
			if _, exists := res.Services[k]; exists {
				return nil, fmt.Errorf("key %s already exists", k)
			}
			res.Services[k] = v
		}
	}

	return &res, nil
}

func main() {
	var includeRoot string
	var fromFile string
	var toFile string
	if len(os.Args) == 3 {
		fromFile = os.Args[1]
		toFile = os.Args[2]
	} else if len(os.Args) == 4 {
		includeRoot = os.Args[1]
		fromFile = os.Args[2]
		toFile = os.Args[3]
	} else {
		fmt.Fprintf(os.Stderr, "No Thrift files specified!\n")
		printUsage()
		os.Exit(1)
	}
	if includeRoot != "" {
		fromFile = path.Join(includeRoot, fromFile)
		toFile = path.Join(includeRoot, toFile)
	}

	if _, err := os.Stat(fromFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File %s does not exist!", fromFile)
		os.Exit(1)
	}

	if _, err := os.Stat(toFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File %s does not exist!", toFile)
		os.Exit(1)
	}

	p := &parser.Parser{
		IncludeRoot:includeRoot,
	}

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

	//fromThrift, err := mergeThriftFiles(fromThrifts)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "failed merging loaded Thrift files: %v, fromThrifts: %v\n", err, fromThrifts)
	//	os.Exit(1)
	//}
	//
	//toThrift, err := mergeThriftFiles(toThrifts)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "failed merging loaded Thrift files: %v, toThrifts: %v\n", err, toThrifts)
	//	os.Exit(1)
	//}
	//
	//
	//err = compareThrift(fromThrift, toThrift)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "not backwards compatible: %s\n", err.Error())
	//	os.Exit(1)
	//}
	err = compareThrifts(fromThrifts, toThrifts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "not backwards compatible: %s\n", err.Error())
		os.Exit(1)
	}
}
