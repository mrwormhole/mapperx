package main

import "go/types"

type info struct {
	typeName    string
	typePackage string
	structType  *types.Struct
}

func NewInfo(typeName string, typePackage string, structType *types.Struct) *info {
	return &info{
		typeName:    typeName,
		typePackage: typePackage,
		structType:  structType,
	}
}
