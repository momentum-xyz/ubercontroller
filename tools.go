//go:build tools
// +build tools

// https://github.com/golang/go/issues/48429
package main

import (
	_ "github.com/swaggo/swag/cmd/swag"
)
