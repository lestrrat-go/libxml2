//go:build static_build
// +build static_build

package clib

// #cgo pkg-config: --static libxml-2.0
// #cgo LDFLAGS: -static
import "C"
