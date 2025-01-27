// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package high contains a set of high-level models that represent OpenAPI 2 and 3 documents.
// These high-level models (porcelain) are used by applications directly, rather than the low-level models
// plumbing) that are used to compose high level models.
//
// High level models are simple to navigate, strongly typed, precise representations of the OpenAPI schema
// that are created from an OpenAPI specification.
//
// All high level objects contains a 'GoLow' method. This 'GoLow' method will return the low-level model that
// was used to create it, which provides an engineer as much low level detail about the raw spec used to create
// those models, things like key/value breakdown of each value, lines, column, source comments etc.
package high

import "github.com/pb33f/libopenapi/datamodel/low"

// GoesLow is used to represent any high-level model. All high level models meet this interface and can be used to
// extract low-level models from any high-level model.
type GoesLow[T any] interface {

	// GoLow returns the low-level object that was used to create the high-level object. This allows consumers
	// to dive-down into the plumbing API at any point in the model.
	GoLow() T
}

// ExtractExtensions is a convenience method for converting low-level extension definitions, to a high level map[string]any
// definition that is easier to consume in applications.
func ExtractExtensions(extensions map[low.KeyReference[string]]low.ValueReference[any]) map[string]any {
	extracted := make(map[string]any)
	for k, v := range extensions {
		extracted[k.Value] = v.Value
	}
	return extracted
}

// UnpackExtensions is a convenience function that makes it easy and simple to unpack an objects extensions
// into a complex type, provided as a generic. This function is for high-level models that implement `GoesLow()`
// and for low-level models that support extensions via `HasExtensions`.
//
// This feature will be upgraded at some point to hold a registry of types and extension mappings to make this
// functionality available a little more automatically.
// You can read more about the discussion here: https://github.com/pb33f/libopenapi/issues/8
//
// `T` represents the Type you want to unpack into
// `R` represents the LOW type of the object that contains the extensions (not the high)
// `low` represents the HIGH type of the object that contains the extensions.
//
// to use:
//  schema := schemaProxy.Schema() // any high-level object that has extensions
//  extensions, err := UnpackExtensions[MyComplexType, low.Schema](schema)
func UnpackExtensions[T any, R low.HasExtensions[T]](low GoesLow[R]) (map[string]*T, error) {
	m := make(map[string]*T)
	ext := low.GoLow().GetExtensions()
	for i := range ext {
		key := i.Value
		g := new(T)
		valueNode := ext[i].ValueNode
		err := valueNode.Decode(g)
		if err != nil {
			return nil, err
		}
		m[key] = g
	}
	return m, nil
}