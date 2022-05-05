//go:build tools
// +build tools

// This is used to ensure tools we use for code generation are added to our modules file.
//
// It ensures we can pin the tool version and not worry about dynamic updates to the modules
// we use to build our software.
package main

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/golang/mock/mockgen"
)
