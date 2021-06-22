package main

import (
	"fmt"
	"go/types"
	"os"
	"regexp"
	"strings"
	"time"

	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

/*type generatorData struct {
	sourceStruct     interface{}
	targetStruct     interface{}
	sourceStructType string
	targetStructType string
}*/

type info struct {
	typeName string
	typePackage string
	structType *types.Struct
}

func main() {
	//var data generatorData
	//flag.StringVar(&data.sourceStructType, "source", "", "The source struct type to be generated for mapperx")
	//flag.StringVar(&data.targetStructType, "target", "", "The target struct type to be generated for mappex")
	//flag.Parse()

	// TODO we should have source and target type then validate those
	if len(os.Args) != 3 {
		exitWithError(fmt.Errorf("expected exactly two arguments: <source type> <target type>"))
	}
	sourceType := os.Args[1]
	sourceTypePackage, sourceTypeName := splitSourceType(sourceType)

	targetType := os.Args[2]
	targetTypePackage, targetTypeName := splitSourceType(targetType)

	sourceStructType := checkTypeEligibility(sourceTypePackage, sourceTypeName)
	targetStructType := checkTypeEligibility(targetTypePackage, targetTypeName)

	sourceInfo := &info{
		typeName:    sourceTypeName,
		typePackage: sourceTypePackage,
		structType:  sourceStructType,
	}
	targetInfo := &info{
		typeName: targetTypeName,
		typePackage: targetTypePackage,
		structType: targetStructType,
	}

	err := generateMapperx(sourceInfo, targetInfo)
	if err != nil {
		exitWithError(err)
	}
}

// Use a simple regexp pattern to match tag values
var structColPattern = regexp.MustCompile(`mapperx:"([^"]+)"`)

func camelCase(s string) string {
	f := strings.ToLower(string(s[0]))
	return f + s[1:]
}

func generateMapperx(s *info, t *info) error {
	f := NewFile("mapperx")
	comment := fmt.Sprintf("Code generated by mapperx, PLEASE DO NOT EVER EDIT. GENERATATED AT %s", time.Now())
	f.PackageComment(comment)

	// struct-to-struct type check
	functionName := fmt.Sprintf("map%sTo%s", s.typeName, t.typeName)
	if s.typeName == t.typeName || s.structType.String() == t.structType.String() {
		return fmt.Errorf("expected 2 different struct types")
	}

	// form the function signature
	signature := f.Func().Id(functionName).Params(
		Id(camelCase(s.typeName)).Op("*").Qual(s.typePackage, t.typeName),
		Id(camelCase(t.typeName)).Op("*").Qual(t.typePackage, t.typeName),
	)

	var mappedFields []Code

	for i := 0; i < s.structType.NumFields(); i++ {
		sourceField := s.structType.Field(i)
		for j := 0; j < t.structType.NumFields(); j++ {
			targetField := t.structType.Field(j)
			if sourceField.Name() == targetField.Name() {
				fmt.Println("MATCHED")
				fmt.Printf("source: %s , target: %s \n", sourceField.Name(), targetField.Name())

				switch sourceField.Type().(type) {
				case *types.Basic:
					line := Id(camelCase(t.typeName)).Dot(targetField.Name()).
						Op("=").
						Id(camelCase(s.typeName)).Dot(sourceField.Name())
					mappedFields = append(mappedFields, line)
				case *types.Named:
					// TODO what if a struct is named?
					line := Id(camelCase(t.typeName)).Dot(targetField.Name()).
						Op("=").
						Id(camelCase(s.typeName)).Dot(sourceField.Name())
					mappedFields = append(mappedFields, line)
				case *types.Slice:
					// TODO what if slice has structs instead of primitives?
					allocLine := Id(camelCase(t.typeName)).Dot(targetField.Name()).
						Op("=").Make(Id(camelCase(targetField.Type().String())),Len(Id(camelCase(s.typeName)).Dot(sourceField.Name())))
					copyLine := Copy(Id(camelCase(t.typeName)).Dot(targetField.Name()), Id(camelCase(s.typeName)).Dot(sourceField.Name()))
					mappedFields = append(mappedFields, allocLine, copyLine)
				//case *types.Array:
				//case *types.Map:
				//case *types.Struct:
				default:
					fmt.Printf("%v type is not supported, so it is ignored for now \n", sourceField.Type())
				}

				break
			}
		}
	}

	// Iterate over struct fields
	/*for i := 0; i < sourceStructType.NumFields(); i++ {
		field := sourceStructType.Field(i)

		// Generate code for each changeset field
		code := Id(field.Name())
		switch v := field.Type().(type) {

		case *types.Named:
			typeName := v.Obj()
			// Qual automatically imports packages
			code.Op("*").Qual(
				typeName.Pkg().Path(),
				typeName.Name(),
			)
		default:
			return fmt.Errorf("struct field type not handled: %T", v)
		}
		changeSetFields = append(changeSetFields, code)
	}*/

	// form the body
	signature.Block(
		mappedFields...,
	)


	fmt.Printf("%#v", f)
	return nil
	// 1. Collect code in toMap() block
	/*var toMapBlock []Code

	// 2. Build "m := make(map[string]interface{})"
	toMapBlock = append(toMapBlock, Id("m").Op(":=").Make(Map(String()).Interface()))

	for i := 0; i < sourceStructType.NumFields(); i++ {
		field := sourceStructType.Field(i)
		tagValue := sourceStructType.Tag(i)

		matches := structColPattern.FindStringSubmatch(tagValue)
		if matches == nil {
			continue
		}
		col := matches[1]

		// 3. Build "if c.Field != nil { m["col"] = *c.Field }"
		code := If(Id("c").Dot(field.Name()).Op("!=").Nil()).Block(
			Id("m").Index(Lit(col)).Op("=").Op("*").Id("c").Dot(field.Name()),
		)
		toMapBlock = append(toMapBlock, code)
	}

	// 4. Build return statement
	toMapBlock = append(toMapBlock, Return(Id("m")))

	// 5. Build toMap method
	f.Func().Params(
		Id("c").Id(changeSetName),
	).Id("toMap").Params().Map(String()).Interface().Block(
		toMapBlock...,
	)
	*/
	/*
		// Build the target file name
		goFile := os.Getenv("GOFILE")
		ext := filepath.Ext(goFile)
		baseFilename := goFile[0 : len(goFile)-len(ext)]
		targetFilename := baseFilename + "_" + strings.ToLower(sourceTypeName) + "_gen.go"

		// Write generated file
		return f.Save(targetFilename)*/
}

func checkTypeEligibility(typePackage string, typeName string) *types.Struct {
	// Inspect package and use type checker to infer imported types
	pkg := loadPackage(typePackage)

	// Lookup the given source type name in the package declarations
	obj := pkg.Types.Scope().Lookup(typeName)
	if obj == nil {
		exitWithError(fmt.Errorf("%s not found in declared types of %s", typeName, pkg))
	}

	// We check if it is a declared type
	if _, ok := obj.(*types.TypeName); !ok {
		exitWithError(fmt.Errorf("%v is not a named type", obj))
	}
	// We expect the underlying type to be a struct
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		exitWithError(fmt.Errorf("type %v is not a struct", obj))
	}
	return structType
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		exitWithError(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func splitSourceType(sourceType string) (string, string) {
	index := strings.LastIndexByte(sourceType, '.')
	if index == -1 {
		exitWithError(fmt.Errorf(`expected qualified type as "filepath/filename.MyType"`))
	}
	sourceTypePackage := sourceType[0:index]
	sourceTypeName := sourceType[index+1:]
	return sourceTypePackage, sourceTypeName
}

func exitWithError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
